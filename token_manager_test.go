package avs

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/aerospike/avs-client-go/protos"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestRefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testUser"
	password := "testPassword"
	logger := slog.Default()

	// Create a mock gRPC client connection
	mockConn := NewMockgrpcClientConn(ctrl)
	mockAuthServiceClient := protos.NewMockAuthServiceClient(ctrl)
	mockConnClients := &connection{
		grpcConn:   mockConn,
		authClient: mockAuthServiceClient,
	}

	b64token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJleHAiOjMwMDAwMDAwMDAsImlhdCI6MTcyNzExMDc1NH0.GD01CEWxW6-7lHcyeetM95WKdUlwY85m5lFqzcTCtzs"

	// Set up expectations for AuthServiceClient.Authenticate()
	mockAuthServiceClient.
		EXPECT().
		Authenticate(gomock.Any(), gomock.Any()).
		Return(&protos.AuthResponse{
			Token: b64token,
		}, nil)

	// Create the token manager with the mock connection
	tm := newGrpcJWTToken(username, password, logger)

	// Refresh the token
	err := tm.RefreshToken(context.Background(), mockConnClients)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the token was set correctly
	assert.Equal(t, "Bearer "+b64token, tm.token.Load().(string))
}

func TestRefreshToken_FailedToRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testUser"
	password := "testPassword"
	logger := slog.Default()

	// Create a mock gRPC client connection
	mockConn := NewMockgrpcClientConn(ctrl)
	mockAuthServiceClient := protos.NewMockAuthServiceClient(ctrl)
	mockConnClients := &connection{
		grpcConn:   mockConn,
		authClient: mockAuthServiceClient,
	}

	// Set up expectations for AuthServiceClient.Authenticate()
	mockAuthServiceClient.
		EXPECT().
		Authenticate(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("foo"))

	// Create the token manager with the
	tm := newGrpcJWTToken(username, password, logger)

	// Refresh
	err := tm.RefreshToken(context.Background(), mockConnClients)

	// Assert error occurred
	assert.Equal(t, "failed to authenticate: foo", err.Error())
}

func TestRefreshToken_FailedMissingClaims(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testUser"
	password := "testPassword"
	logger := slog.Default()

	// Create a mock gRPC client connection
	mockConn := NewMockgrpcClientConn(ctrl)
	mockAuthServiceClient := protos.NewMockAuthServiceClient(ctrl)
	mockConnClients := &connection{
		grpcConn:   mockConn,
		authClient: mockAuthServiceClient,
	}

	// Set up expectations for AuthServiceClient.Authenticate()
	mockAuthServiceClient.
		EXPECT().
		Authenticate(gomock.Any(), gomock.Any()).
		Return(&protos.AuthResponse{
			Token: "badToken",
		}, nil)

	// Create the token manager with the
	tm := newGrpcJWTToken(username, password, logger)

	// Refresh
	err := tm.RefreshToken(context.Background(), mockConnClients)

	// Assert error occurred
	assert.Error(t, err)
	assert.Equal(t, "failed to authenticate: missing either header, payload, or signature", err.Error())
}

func TestRefreshToken_FailedToDecodeToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testUser"
	password := "testPassword"
	logger := slog.Default()

	// Create a mock gRPC client connection
	mockConn := NewMockgrpcClientConn(ctrl)
	mockAuthServiceClient := protos.NewMockAuthServiceClient(ctrl)
	mockConnClients := &connection{
		grpcConn:   mockConn,
		authClient: mockAuthServiceClient,
	}

	// Set up expectations for AuthServiceClient.Authenticate()
	mockAuthServiceClient.
		EXPECT().
		Authenticate(gomock.Any(), gomock.Any()).
		Return(&protos.AuthResponse{
			Token: "badToken.blah.foo",
		}, nil)

	// Create the token manager with the
	tm := newGrpcJWTToken(username, password, logger)

	// Refresh
	err := tm.RefreshToken(context.Background(), mockConnClients)

	// Assert error occurred
	assert.Error(t, err)
	assert.Equal(t, "failed to authenticate: invalid character 'V' in literal null (expecting 'u')", err.Error())
}

func TestRefreshToken_FailedInvalidJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testUser"
	password := "testPassword"
	logger := slog.Default()

	// Create a mock gRPC client connection
	mockConn := NewMockgrpcClientConn(ctrl)
	mockAuthServiceClient := protos.NewMockAuthServiceClient(ctrl)
	mockConnClients := &connection{
		grpcConn:   mockConn,
		authClient: mockAuthServiceClient,
	}

	b64token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJleHAiOjMwMDAwMDAwMDAsImlhdCI6MTcyNzExMD.GD01CEWxW6-7lHcyeetM95WKdUlwY85m5lFqzcTCtzs"

	// Set up expectations for AuthServiceClient.Authenticate()
	mockAuthServiceClient.
		EXPECT().
		Authenticate(gomock.Any(), gomock.Any()).
		Return(&protos.AuthResponse{
			Token: b64token,
		}, nil)

	// Create the token manager with the mock connection
	tm := newGrpcJWTToken(username, password, logger)

	// Refresh the token
	err := tm.RefreshToken(context.Background(), mockConnClients)

	// Assert no error occurred
	assert.Error(t, err)
	assert.Equal(t, "failed to authenticate: unexpected end of JSON input", err.Error())
}

func TestRefreshToken_FailedFindExp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testUser"
	password := "testPassword"
	logger := slog.Default()

	// Create a mock gRPC client connection
	mockConn := NewMockgrpcClientConn(ctrl)
	mockAuthServiceClient := protos.NewMockAuthServiceClient(ctrl)
	mockConnClients := &connection{
		grpcConn:   mockConn,
		authClient: mockAuthServiceClient,
	}

	b64token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJpYXQiOjE3MjcxMTA3NTR9.50IZcLoS7mQPzQsKvJZyXNUukvT5FdiqN2tynNIjHuk"

	// Set up expectations for AuthServiceClient.Authenticate()
	mockAuthServiceClient.
		EXPECT().
		Authenticate(gomock.Any(), gomock.Any()).
		Return(&protos.AuthResponse{
			Token: b64token,
		}, nil)

	// Create the token manager with the mock connection
	tm := newGrpcJWTToken(username, password, logger)

	// Refresh the token
	err := tm.RefreshToken(context.Background(), mockConnClients)

	// Assert no error occurred
	assert.Error(t, err)
	assert.Equal(t, "failed to authenticate: unable to find exp in token", err.Error())
}

func TestRefreshToken_FailedFindIat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testUser"
	password := "testPassword"
	logger := slog.Default()

	// Create a mock gRPC client connection
	mockConn := NewMockgrpcClientConn(ctrl)
	mockAuthServiceClient := protos.NewMockAuthServiceClient(ctrl)
	mockConnClients := &connection{
		grpcConn:   mockConn,
		authClient: mockAuthServiceClient,
	}

	b64token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJleHAiOjMwMDAwMDAwMDB9.f5DtjF1sYLH6fz0ThcFKxwngIXkVMLnhJtIrjLi_1p0"

	// Set up expectations for AuthServiceClient.Authenticate()
	mockAuthServiceClient.
		EXPECT().
		Authenticate(gomock.Any(), gomock.Any()).
		Return(&protos.AuthResponse{
			Token: b64token,
		}, nil)

	// Create the token manager with the mock connection
	tm := newGrpcJWTToken(username, password, logger)

	// Refresh the token
	err := tm.RefreshToken(context.Background(), mockConnClients)

	// Assert no error occurred
	assert.Error(t, err)
	assert.Equal(t, "failed to authenticate: unable to find iat in token", err.Error())
}

func TestRefreshToken_FailedTtlLessThan0(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testUser"
	password := "testPassword"
	logger := slog.Default()

	// Create a mock gRPC client connection
	mockConn := NewMockgrpcClientConn(ctrl)
	mockAuthServiceClient := protos.NewMockAuthServiceClient(ctrl)
	mockConnClients := &connection{
		grpcConn:   mockConn,
		authClient: mockAuthServiceClient,
	}

	b64token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJleHAiOjE3MjcxMTA3NTMsImlhdCI6MTcyNzExMDc1NH0.YdH5twU6-LGLtgvD2sktiw1j40MRUe_r4oPN565z4Ok"

	// Set up expectations for AuthServiceClient.Authenticate()
	mockAuthServiceClient.
		EXPECT().
		Authenticate(gomock.Any(), gomock.Any()).
		Return(&protos.AuthResponse{
			Token: b64token,
		}, nil)

	// Create the token manager with the mock connection
	tm := newGrpcJWTToken(username, password, logger)

	// Refresh the token
	err := tm.RefreshToken(context.Background(), mockConnClients)

	// Assert no error occurred
	assert.Error(t, err)
	assert.Equal(t, "failed to authenticate: jwt ttl is less than 0", err.Error())
}
