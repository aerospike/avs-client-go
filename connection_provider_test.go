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

func TestConnectToSeeds_FailedAlreadyConnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cp := &connectionProvider{
		isLoadBalancer: true,
		closed:         atomic.Bool{},
		logger:         slog.Default(),
	}

	cp.seedConns = []*connection{
		{
			grpcConn: NewMockgrpcClientConn(ctrl),
		},
	}

	err := cp.connectToSeeds(context.Background())

	assert.Equal(t, errors.New("seed connections already exist, close them first"), err)
}

func TestConnectToSeeds_FailedFailedToCreateConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cp := &connectionProvider{
		isLoadBalancer: true,
		closed:         atomic.Bool{},
		logger:         slog.Default(),
		grpcConnFactory: func(hostPort *HostPort) (grpcClientConn, error) {
			return nil, fmt.Errorf("foo")
		},
	}

	cp.seeds = HostPortSlice{
		&HostPort{
			Host: "host",
			Port: 3000,
		},
	}

	err := cp.connectToSeeds(context.Background())

	assert.Equal(t, NewAVSError("failed to connect to seeds", nil), err)
}

func TestConnectToSeeds_FailedToRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockToken := NewMocktokenManager(ctrl)
	mockToken.
		EXPECT().
		RefreshToken(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("foo"))

	cp := &connectionProvider{
		isLoadBalancer: true,
		closed:         atomic.Bool{},
		logger:         slog.Default(),
		grpcConnFactory: func(hostPort *HostPort) (grpcClientConn, error) {
			return nil, nil
		},
		connFactory: func(conn grpcClientConn) *connection {
			return &connection{}
		},
		token: mockToken,
	}

	cp.seeds = HostPortSlice{
		&HostPort{
			Host: "host",
			Port: 3000,
		},
	}

	err := cp.connectToSeeds(context.Background())

	assert.Equal(t, NewAVSError("failed to connect to seeds", fmt.Errorf("foo")), err)
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

func TestUpdateClusterConns_NewClusterIDWithDIFFERENTNodeIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	mockNewGrpcConn1111 := NewMockgrpcClientConn(ctrl)
	mockNewGrpcConn2222 := NewMockgrpcClientConn(ctrl)

	mockClusterInfoClient1111 := protos.NewMockClusterInfoServiceClient(ctrl)
	mockClusterInfoClient2222 := protos.NewMockClusterInfoServiceClient(ctrl)

	mockAboutClient1111 := protos.NewMockAboutServiceClient(ctrl)
	mockAboutClient2222 := protos.NewMockAboutServiceClient(ctrl)

	mockAboutClient1111.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	mockAboutClient2222.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	cp := &connectionProvider{
		logger:         slog.Default(),
		nodeConns:      make(map[uint64]*connectionAndEndpoints),
		seedConns:      []*connection{},
		tlsConfig:      &tls.Config{},
		seeds:          HostPortSlice{},
		nodeConnsLock:  &sync.RWMutex{},
		tendInterval:   time.Second * 1,
		clusterID:      123,
		listenerName:   nil,
		isLoadBalancer: false,
		token:          nil,
		stopTendChan:   make(chan struct{}),
		closed:         atomic.Bool{},
		grpcConnFactory: func(hostPort *HostPort) (grpcClientConn, error) {
			if hostPort.String() == "1.1.1.1:3000" {
				return mockNewGrpcConn1111, nil
			} else if hostPort.String() == "2.2.2.2:3000" {
				return mockNewGrpcConn2222, nil
			}

			return nil, fmt.Errorf("foo")
		},
		connFactory: func(grpcConn grpcClientConn) *connection {
			if grpcConn == mockNewGrpcConn1111 {
				return &connection{
					clusterInfoClient: mockClusterInfoClient1111,
					aboutClient:       mockAboutClient1111,
				}
			} else if grpcConn == mockNewGrpcConn2222 {
				return &connection{
					clusterInfoClient: mockClusterInfoClient2222,
					aboutClient:       mockAboutClient2222,
				}
			}

			return nil
		},
	}

	cp.logger = cp.logger.With(slog.String("test", "TestUpdateClusterConns_NewClusterID"))

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
			Id: 789, // Different cluster id from client 2
		}, nil)

	mockClusterInfoClient1.
		EXPECT().
		GetClusterEndpoints(gomock.Any(), gomock.Any()).
		Return(&protos.ClusterNodeEndpoints{
			Endpoints: map[uint64]*protos.ServerEndpointList{ // Smaller num of endpoints from client 2
				0: {
					Endpoints: []*protos.ServerEndpoint{
						{
							Address: "1.1.1.1",
							Port:    3000,
						},
					},
				},
			},
		}, nil)

	grpcConn1.
		EXPECT().
		Close().
		Return(nil)

	grpcConn2.
		EXPECT().
		Target().
		Return("")

	expectedClusterID := uint64(456)
	// After a new cluster is discovered we expect these to be the new nodeConns
	expectedNewNodeConns := map[uint64]*connectionAndEndpoints{
		3: {
			conn: &connection{
				clusterInfoClient: mockClusterInfoClient1111,
				aboutClient:       mockAboutClient1111,
			},
			endpoints: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "1.1.1.1",
						Port:    3000,
					},
				},
			},
		},
		4: {
			conn: &connection{
				clusterInfoClient: mockClusterInfoClient2222,
				aboutClient:       mockAboutClient2222,
			},
			endpoints: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "2.2.2.2",
						Port:    3000,
					},
				},
			},
		},
	}

	mockClusterInfoClient2.
		EXPECT().
		GetClusterId(gomock.Any(), gomock.Any()).
		Return(&protos.ClusterId{
			Id: expectedClusterID,
		}, nil)

	mockClusterInfoClient2.
		EXPECT().
		GetClusterEndpoints(gomock.Any(), gomock.Any()).
		Return(&protos.ClusterNodeEndpoints{
			Endpoints: map[uint64]*protos.ServerEndpointList{ // larger, so the cluster id 456 will win
				3: {
					Endpoints: []*protos.ServerEndpoint{
						{
							Address: "1.1.1.1",
							Port:    3000,
						},
					},
				},
				4: {
					Endpoints: []*protos.ServerEndpoint{
						{
							Address: "2.2.2.2",
							Port:    3000,
						},
					},
				},
			},
		}, nil)

	grpcConn2.
		EXPECT().
		Close().
		Return(nil)

	// Existing node connections. These will be replaced after a new cluster is found.
	cp.nodeConns = map[uint64]*connectionAndEndpoints{
		1: {
			conn: &connection{
				grpcConn:          grpcConn1,
				clusterInfoClient: mockClusterInfoClient1,
			},
			endpoints: &protos.ServerEndpointList{},
		},
		2: {
			conn: &connection{
				grpcConn:          grpcConn2,
				clusterInfoClient: mockClusterInfoClient2,
			},
			endpoints: &protos.ServerEndpointList{},
		},
	}

	cp.logger.Debug("Running updateClusterConns")

	cp.updateClusterConns(ctx)

	assert.Equal(t, expectedClusterID, cp.clusterID)
	assert.Len(t, cp.nodeConns, 2)

	for k, v := range cp.nodeConns {
		assert.EqualExportedValues(t, expectedNewNodeConns[k].endpoints, v.endpoints)
	}
}

func TestUpdateClusterConns_NewClusterIDWithSAMENodeIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	grpcConn1 := NewMockgrpcClientConn(ctrl)
	mockClusterInfoClient1 := protos.NewMockClusterInfoServiceClient(ctrl)
	grpcConn2 := NewMockgrpcClientConn(ctrl)
	mockClusterInfoClient2 := protos.NewMockClusterInfoServiceClient(ctrl)

	expectedClusterID := uint64(456)

	mockClusterInfoClient2.
		EXPECT().
		GetClusterId(gomock.Any(), gomock.Any()).
		Return(&protos.ClusterId{
			Id: expectedClusterID,
		}, nil)

	mockClusterInfoClient2.
		EXPECT().
		GetClusterEndpoints(gomock.Any(), gomock.Any()).
		Return(&protos.ClusterNodeEndpoints{
			Endpoints: map[uint64]*protos.ServerEndpointList{ // larger, so the cluster id 456 will win
				1: {
					Endpoints: []*protos.ServerEndpoint{
						{
							Address: "1.1.1.1",
							Port:    3000,
						},
					},
				},
				2: {
					Endpoints: []*protos.ServerEndpoint{
						{
							Address: "2.2.2.2",
							Port:    3000,
						},
					},
				},
			},
		}, nil)

	mockNewGrpcConn1111 := NewMockgrpcClientConn(ctrl)
	mockNewGrpcConn2222 := NewMockgrpcClientConn(ctrl)

	mockClusterInfoClient1111 := protos.NewMockClusterInfoServiceClient(ctrl)
	mockClusterInfoClient2222 := protos.NewMockClusterInfoServiceClient(ctrl)

	mockAboutClient1111 := protos.NewMockAboutServiceClient(ctrl)
	mockAboutClient2222 := protos.NewMockAboutServiceClient(ctrl)

	mockAboutClient1111.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	mockAboutClient2222.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	cp := &connectionProvider{
		logger:         slog.Default(),
		seedConns:      []*connection{},
		tlsConfig:      &tls.Config{},
		seeds:          HostPortSlice{},
		nodeConnsLock:  &sync.RWMutex{},
		tendInterval:   time.Second * 1,
		clusterID:      123,
		listenerName:   nil,
		isLoadBalancer: false,
		token:          nil,
		stopTendChan:   make(chan struct{}),
		closed:         atomic.Bool{},
		grpcConnFactory: func(hostPort *HostPort) (grpcClientConn, error) {
			if hostPort.String() == "1.1.1.1:3000" {
				return mockNewGrpcConn1111, nil
			} else if hostPort.String() == "2.2.2.2:3000" {
				return mockNewGrpcConn2222, nil
			}

			return nil, fmt.Errorf("foo")
		},
		connFactory: func(grpcConn grpcClientConn) *connection {
			if grpcConn == mockNewGrpcConn1111 {
				return &connection{
					clusterInfoClient: mockClusterInfoClient1111,
					aboutClient:       mockAboutClient1111,
				}
			} else if grpcConn == mockNewGrpcConn2222 {
				return &connection{
					clusterInfoClient: mockClusterInfoClient2222,
					aboutClient:       mockAboutClient2222,
				}
			}

			return nil
		},
		// Existing node connections. These will be replaced after a new cluster is found.
		nodeConns: map[uint64]*connectionAndEndpoints{
			1: {
				conn: &connection{
					grpcConn:          grpcConn1,
					clusterInfoClient: mockClusterInfoClient1,
				},
				endpoints: &protos.ServerEndpointList{},
			},
			2: {
				conn: &connection{
					grpcConn:          grpcConn2,
					clusterInfoClient: mockClusterInfoClient2,
				},
				endpoints: &protos.ServerEndpointList{},
			},
		},
	}

	cp.logger = cp.logger.With(slog.String("test", "TestUpdateClusterConns_NewClusterID"))

	cp.logger.Debug("Setting up existing node connections")

	grpcConn1.
		EXPECT().
		Target().
		Return("")

	mockClusterInfoClient1.
		EXPECT().
		GetClusterId(gomock.Any(), gomock.Any()).
		Return(&protos.ClusterId{
			Id: 789, // Different cluster id from client 2
		}, nil)

	mockClusterInfoClient1.
		EXPECT().
		GetClusterEndpoints(gomock.Any(), gomock.Any()).
		Return(&protos.ClusterNodeEndpoints{
			Endpoints: map[uint64]*protos.ServerEndpointList{ // Smaller num of endpoints from client 2
				0: {
					Endpoints: []*protos.ServerEndpoint{
						{
							Address: "1.1.1.1",
							Port:    3000,
						},
					},
				},
			},
		}, nil)

	grpcConn1.
		EXPECT().
		Close().
		Return(nil)

	grpcConn2.
		EXPECT().
		Target().
		Return("")

	// After a new cluster is discovered we expect these to be the new nodeConns
	expectedNewNodeConns := map[uint64]*connectionAndEndpoints{
		1: {
			conn: &connection{
				clusterInfoClient: mockClusterInfoClient1111,
				aboutClient:       mockAboutClient1111,
			},
			endpoints: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "1.1.1.1",
						Port:    3000,
					},
				},
			},
		},
		2: {
			conn: &connection{
				clusterInfoClient: mockClusterInfoClient2222,
				aboutClient:       mockAboutClient2222,
			},
			endpoints: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "2.2.2.2",
						Port:    3000,
					},
				},
			},
		},
	}

	grpcConn2.
		EXPECT().
		Close().
		Return(nil)

	cp.logger.Debug("Running updateClusterConns")

	cp.updateClusterConns(ctx)

	assert.Equal(t, expectedClusterID, cp.clusterID)
	assert.Len(t, cp.nodeConns, 2)

	for k, v := range cp.nodeConns {
		assert.EqualExportedValues(t, expectedNewNodeConns[k].endpoints, v.endpoints)
	}
}
