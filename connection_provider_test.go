package avs

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
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
