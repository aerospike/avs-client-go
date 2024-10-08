package avs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aerospike/avs-client-go/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// grpcTokenManager is responsible for managing authentication tokens and refreshing
// them when necessary.
//
//nolint:govet // We will favor readability over field alignment
type grpcTokenManager struct {
	username         string
	password         string
	token            atomic.Value
	refreshTime      atomic.Value
	logger           *slog.Logger
	stopRefreshChan  chan struct{}
	refreshScheduled bool
}

// newGrpcJWTToken creates a new tokenManager instance with the provided username, password, and logger.
func newGrpcJWTToken(username, password string, logger *slog.Logger) *grpcTokenManager {
	logger.WithGroup("jwt")

	logger.Debug("creating new token manager")

	return &grpcTokenManager{
		username:        username,
		password:        password,
		logger:          logger,
		stopRefreshChan: make(chan struct{}),
	}
}

// Close stops the scheduled token refresh and closes the token manager.
func (tm *grpcTokenManager) Close() {
	if tm.refreshScheduled {
		tm.logger.Debug("stopping scheduled token refresh")
		tm.stopRefreshChan <- struct{}{}
		<-tm.stopRefreshChan
	}

	tm.logger.Debug("closed")
}

// setRefreshTimeFromTTL sets the refresh time based on the provided time-to-live (TTL) duration.
func (tm *grpcTokenManager) setRefreshTimeFromTTL(ttl time.Duration) {
	tm.refreshTime.Store(time.Now().Add(ttl))
}

// RefreshToken refreshes the authentication token using the provided gRPC client connection.
// It returns a boolean indicating if the token was successfully refreshed and
// an error if any. It is not thread safe.
func (tm *grpcTokenManager) RefreshToken(ctx context.Context, conn *connection) error {
	// We only want one goroutine to refresh the token at a time
	resp, err := conn.authClient.Authenticate(ctx, &protos.AuthRequest{
		Credentials: createUserPassCredential(tm.username, tm.password),
	})

	if err != nil {
		return fmt.Errorf("%s: %w", "failed to authenticate", err)
	}

	claims := strings.Split(resp.GetToken(), ".")

	if len(claims) < 3 {
		return fmt.Errorf("failed to authenticate: missing either header, payload, or signature")
	}

	decClaims, err := base64.RawURLEncoding.DecodeString(claims[1])

	if err != nil {
		return fmt.Errorf("%s: %w", "failed to authenticate", err)
	}

	tokenMap := make(map[string]any, 8)
	err = json.Unmarshal(decClaims, &tokenMap)

	if err != nil {
		return fmt.Errorf("%s: %w", "failed to authenticate", err)
	}

	expiryToken, ok := tokenMap["exp"].(float64)
	if !ok {
		return fmt.Errorf("failed to authenticate: unable to find exp in token")
	}

	iat, ok := tokenMap["iat"].(float64)
	if !ok {
		return fmt.Errorf("failed to authenticate: unable to find iat in token")
	}

	ttl := time.Duration(expiryToken-iat) * time.Second
	if ttl <= 0 {
		return fmt.Errorf("failed to authenticate: jwt ttl is less than 0")
	}

	tm.logger.DebugContext(
		ctx,
		"successfully parsed token",
		slog.Float64("exp", expiryToken),
		slog.Float64("iat", iat),
		slog.Duration("ttl", ttl),
	)

	// Set expiry based on local clock.
	tm.setRefreshTimeFromTTL(ttl)
	tm.token.Store("Bearer " + resp.GetToken())

	return nil
}

// ScheduleRefresh schedules the token refresh using the provided function to
// get the gRPC client connection. This is not threadsafe. It should only be
// called once.
func (tm *grpcTokenManager) ScheduleRefresh(getConn func() (*connection, error)) {
	if tm.refreshScheduled {
		tm.logger.Warn("refresh already scheduled")
	}

	tm.logger.Debug("scheduling token refresh")

	tm.refreshScheduled = true
	timer := time.NewTimer(0)

	go func() {
		for {
			connClients, err := getConn()
			if err != nil {
				tm.logger.Warn("failed to refresh token", slog.Any("error", err))
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

			err = tm.RefreshToken(ctx, connClients)
			if err != nil {
				tm.logger.Warn("failed to refresh token", slog.Any("error", err))
			}

			cancel()

			waitFor := time.Until(tm.refreshTime.Load().(time.Time)) - time.Second*5

			tm.logger.Debug("waiting to refresh token", slog.Duration("waitTime", waitFor))
			timer.Reset(waitFor)

			select {
			case <-timer.C:
			case <-tm.stopRefreshChan:
				tm.refreshScheduled = false

				timer.Stop()
				tm.stopRefreshChan <- struct{}{}
				tm.logger.Debug("stopped scheduled token refresh")

				return
			}
		}
	}()
}

// RequireTransportSecurity returns true to indicate that transport security is required.
func (tm *grpcTokenManager) RequireTransportSecurity() bool {
	return true
}

// UnaryInterceptor returns the grpc unary client interceptor that attaches the token to outgoing requests.
func (tm *grpcTokenManager) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(tm.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

// StreamInterceptor returns the grpc stream client interceptor that attaches the token to outgoing requests.
func (tm *grpcTokenManager) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		return streamer(tm.attachToken(ctx), desc, cc, method, opts...)
	}
}

// attachToken attaches the authentication token to the outgoing context.
func (tm *grpcTokenManager) attachToken(ctx context.Context) context.Context {
	rawToken := tm.token.Load()
	if rawToken == nil {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, "Authorization", tm.token.Load().(string))
}
