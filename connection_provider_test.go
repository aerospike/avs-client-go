package avs

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aerospike/avs-client-go/protos"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestNewConnectionProvider_FailSeedsNil(t *testing.T) {
	seeds := HostPortSlice{}
	listenerName := "listener"
	isLoadBalancer := false
	tlsConfig := &tls.Config{}
	var logger *slog.Logger
	var token tokenManager

	cp, err := newConnectionProvider(context.Background(), seeds, &listenerName, isLoadBalancer, token, tlsConfig, logger)
	defer cp.Close()

	assert.Nil(t, cp)
	assert.Equal(t, err, errors.New("seeds cannot be nil or empty"))
}

func TestNewConnectionProvider_FailNoTLS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	seeds := HostPortSlice{
		&HostPort{
			Host: "host",
			Port: 3000,
		},
	}
	listenerName := "listener"
	isLoadBalancer := false

	var tlsConfig *tls.Config
	var logger *slog.Logger

	token := NewMocktokenManager(ctrl)

	token.
		EXPECT().
		RequireTransportSecurity().
		Return(true)

	cp, err := newConnectionProvider(context.Background(), seeds, &listenerName, isLoadBalancer, token, tlsConfig, logger)
	defer cp.Close()

	assert.Nil(t, cp)
	assert.Equal(t, err, errors.New("tlsConfig is required when username/password authentication"))
}

func TestNewConnectionProvider_FailConnectToSeedConns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	seeds := HostPortSlice{
		&HostPort{
			Host: "host",
			Port: 3000,
		},
	}
	listenerName := "listener"
	isLoadBalancer := false

	var tlsConfig *tls.Config
	var logger *slog.Logger
	var token tokenManager

	cp, err := newConnectionProvider(ctx, seeds, &listenerName, isLoadBalancer, token, tlsConfig, logger)
	defer cp.Close()

	assert.Nil(t, cp)
	assert.Equal(t, "failed to connect to seeds: context deadline exceeded", err.Error())
}

func TestClose_FailsToCloseConns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeedConn := NewMockgrpcClientConn(ctrl)
	mockNodeConn := NewMockgrpcClientConn(ctrl)

	mockSeedConn.
		EXPECT().
		Close().
		Return(fmt.Errorf("foo"))

	mockSeedConn.
		EXPECT().
		Target().
		Return("")

	mockNodeConn.
		EXPECT().
		Close().
		Return(fmt.Errorf("bar"))

	mockNodeConn.
		EXPECT().
		Target().
		Return("")

	cp := &connectionProvider{
		isLoadBalancer: true,
		seedConns: []*connection{
			{
				grpcConn: mockSeedConn,
			},
		},
		nodeConns: map[uint64]*connectionAndEndpoints{
			uint64(1): {
				conn:      &connection{grpcConn: mockNodeConn},
				endpoints: &protos.ServerEndpointList{},
			},
		},
		logger: slog.Default(),
	}

	err := cp.Close()

	assert.Equal(t, fmt.Errorf("foo"), err)
}

func TestGetSeedConn_FailBecauseClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cp := &connectionProvider{
		isLoadBalancer: true,
		closed:         atomic.Bool{},
		logger:         slog.Default(),
	}

	cp.closed.Store(true)

	_, err := cp.GetSeedConn()

	assert.Equal(t, errors.New("connectionProvider is closed, cannot get connection"), err)
}

func TestGetSeedConn_FailSeedConnEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cp := &connectionProvider{
		isLoadBalancer: true,
		closed:         atomic.Bool{},
		logger:         slog.Default(),
	}

	cp.closed.Store(false)

	_, err := cp.GetSeedConn()

	assert.Equal(t, errors.New("no seed connections found"), err)
}
func TestUpdateClusterConns_NoNewClusterID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	cp := &connectionProvider{
		logger:         slog.Default(),
		nodeConns:      make(map[uint64]*connectionAndEndpoints),
		seedConns:      []*connection{},
		tlsConfig:      nil,
		seeds:          HostPortSlice{},
		nodeConnsLock:  &sync.RWMutex{},
		tendInterval:   time.Second * 1,
		clusterID:      123,
		listenerName:   nil,
		isLoadBalancer: false,
		token:          nil,
		stopTendChan:   make(chan struct{}),
		closed:         atomic.Bool{},
	}

	cp.logger = cp.logger.With(slog.String("test", "TestUpdateClusterConns_NoNewClusterID"))

	cp.logger.Debug("Setting up existing node connections")

	grpcConn1 := NewMockgrpcClientConn(ctrl)
	mockClusterInfoClient1 := protos.NewMockClusterInfoServiceClient(ctrl)
	grpcConn2 := NewMockgrpcClientConn(ctrl)
	mockClusterInfoClient2 := protos.NewMockClusterInfoServiceClient(ctrl)

	grpcConn1.
		EXPECT().
		Target().
		Return("")

	mockClusterInfoClient1.
		EXPECT().
		GetClusterId(gomock.Any(), gomock.Any()).
		Return(&protos.ClusterId{
			Id: 123,
		}, nil)

	grpcConn2.
		EXPECT().
		Target().
		Return("")

	mockClusterInfoClient2.
		EXPECT().
		GetClusterId(gomock.Any(), gomock.Any()).
		Return(&protos.ClusterId{
			Id: 123,
		}, nil)

	// Existing node connections
	cp.nodeConns[1] = &connectionAndEndpoints{
		conn: &connection{
			grpcConn:          grpcConn1,
			clusterInfoClient: mockClusterInfoClient1,
		},
		endpoints: &protos.ServerEndpointList{},
	}

	cp.nodeConns[2] = &connectionAndEndpoints{
		conn: &connection{
			grpcConn:          grpcConn2,
			clusterInfoClient: mockClusterInfoClient2,
		},
		endpoints: &protos.ServerEndpointList{},
	}

	cp.logger.Debug("Running updateClusterConns")

	cp.updateClusterConns(ctx)

	assert.Equal(t, uint64(123), cp.clusterID)
	assert.Len(t, cp.nodeConns, 2)
}

// func TestUpdateClusterConns_NewClusterID(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
// 	defer cancel()

// 	cp := &connectionProvider{
// 		logger:         slog.Default(),
// 		nodeConns:      make(map[uint64]*connectionAndEndpoints),
// 		seedConns:      []*connection{},
// 		tlsConfig:      &tls.Config{},
// 		seeds:          HostPortSlice{},
// 		nodeConnsLock:  &sync.RWMutex{},
// 		tendInterval:   time.Second * 1,
// 		clusterID:      123,
// 		listenerName:   nil,
// 		isLoadBalancer: false,
// 		token:          nil,
// 		stopTendChan:   make(chan struct{}),
// 		closed:         atomic.Bool{},
// 	}

// 	cp.logger = cp.logger.With(slog.String("test", "TestUpdateClusterConns_NewClusterID"))

// 	cp.logger.Debug("Setting up existing node connections")

// 	grpcConn1 := NewMockgrpcClientConn(ctrl)
// 	mockClusterInfoClient1 := protos.NewMockClusterInfoServiceClient(ctrl)
// 	grpcConn2 := NewMockgrpcClientConn(ctrl)
// 	mockClusterInfoClient2 := protos.NewMockClusterInfoServiceClient(ctrl)

// 	grpcConn1.
// 		EXPECT().
// 		Target().
// 		Return("")

// 	mockClusterInfoClient1.
// 		EXPECT().
// 		GetClusterId(gomock.Any(), gomock.Any()).
// 		Return(&protos.ClusterId{
// 			Id: 123,
// 		}, nil)

// 	grpcConn1.
// 		EXPECT().
// 		Close().
// 		Return(nil)

// 	grpcConn2.
// 		EXPECT().
// 		Target().
// 		Return("")

// 	mockClusterInfoClient2.
// 		EXPECT().
// 		GetClusterId(gomock.Any(), gomock.Any()).
// 		Return(&protos.ClusterId{
// 			Id: 456,
// 		}, nil)

// 	mockClusterInfoClient2.
// 		EXPECT().
// 		GetClusterEndpoints(gomock.Any(), gomock.Any()).
// 		Return(&protos.ClusterNodeEndpoints{
// 			Endpoints: map[uint64]*protos.ServerEndpointList{
// 				3: {
// 					Endpoints: []*protos.ServerEndpoint{
// 						{
// 							Address: "1.1.1.1",
// 							Port:    3000,
// 						},
// 					},
// 				},
// 				4: {
// 					Endpoints: []*protos.ServerEndpoint{
// 						{
// 							Address: "2.2.2.2",
// 							Port:    3000,
// 						},
// 					},
// 				},
// 			},
// 		}, nil)

// 	grpcConn2.
// 		EXPECT().
// 		Close().
// 		Return(nil)

// 	// Existing node connections
// 	cp.nodeConns[1] = &connectionAndEndpoints{
// 		conn: &connection{
// 			grpcConn:          grpcConn1,
// 			clusterInfoClient: mockClusterInfoClient1,
// 		},
// 		endpoints: &protos.ServerEndpointList{},
// 	}

// 	cp.nodeConns[2] = &connectionAndEndpoints{
// 		conn: &connection{
// 			grpcConn:          grpcConn2,
// 			clusterInfoClient: mockClusterInfoClient2,
// 		},
// 		endpoints: &protos.ServerEndpointList{},
// 	}

// 	cp.logger.Debug("Running updateClusterConns")

// 	// New cluster ID
// 	// newEndpoints := &protos.ServerEndpointList{
// 	// 	Endpoints: []*protos.ServerEndpoint{
// 	// 		{
// 	// 			Address: "host1",
// 	// 			Port:    3000,
// 	// 		},
// 	// 		{
// 	// 			Address: "host2",
// 	// 			Port:    3000,
// 	// 		},
// 	// 	},
// 	// }

// 	cp.updateClusterConns(ctx)

// 	// cp.checkAndSetNodeConns(ctx, map[uint64]*protos.ServerEndpointList{
// 	// 	1: newEndpoints,
// 	// 	2: newEndpoints,
// 	// })

// 	// cp.removeDownNodes(map[uint64]*protos.ServerEndpointList{
// 	// 	1: newEndpoints,
// 	// 	2: newEndpoints,
// 	// })

// 	// cp.updateClusterConns(ctx)

// 	assert.Equal(t, uint64(456), cp.clusterID)
// 	assert.Len(t, cp.nodeConns, 2)
// }
