package avs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aerospike/avs-client-go/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type tokenManager struct {
	username         string
	password         string
	token            atomic.Value
	refreshTime      atomic.Value
	lock             sync.Mutex
	logger           *slog.Logger
	stopRefreshChan  chan struct{}
	refreshScheduled bool
}

func newJWTToken(username, password string, logger *slog.Logger) *tokenManager {
	logger.WithGroup("jwt")

	logger.Debug("creating new token manager")

	return &tokenManager{
		username: username,
		password: password,
		lock:     sync.Mutex{},
		logger:   logger,
	}
}

func (tm *tokenManager) Close() {
	if tm.refreshScheduled {
		tm.logger.Debug("stopping scheduled token refresh")
		tm.stopRefreshChan <- struct{}{}
		<-tm.stopRefreshChan
	}

	tm.logger.Debug("closed")
}

func (tm *tokenManager) setRefreshTimeFromTTL(ttl time.Duration) {
	tm.refreshTime.Store(time.Now().Add(ttl))
}

func (tm *tokenManager) expired() bool {
	return time.Now().After(tm.refreshTime.Load().(time.Time))
}

func (tm *tokenManager) RefreshToken(ctx context.Context, conn grpc.ClientConnInterface) (bool, error) {
	// We only want one goroutine to refresh the token at a time
	tm.lock.Lock()
	defer tm.lock.Unlock()

	// If it does not need to be updated return
	if !tm.expired() {
		return false, nil
	}

	client := protos.NewAuthServiceClient(conn)
	resp, err := client.Authenticate(ctx, &protos.AuthRequest{
		Credentials: &protos.Credentials{
			Username: tm.username,
			Credentials: &protos.Credentials_PasswordCredentials{
				PasswordCredentials: &protos.PasswordCredentials{
					Password: tm.password,
				},
			},
		},
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", "failed to authenticate", err)
	}

	claims := strings.Split(resp.GetToken(), ".")
	decClaims, err := base64.RawURLEncoding.DecodeString(claims[1])
	if err != nil {
		return false, fmt.Errorf("%s: %w", "failed to authenticate", err)
	}

	tokenMap := make(map[string]interface{}, 8)
	err = json.Unmarshal(decClaims, &tokenMap)
	if err != nil {
		return false, fmt.Errorf("%s: %w", "failed to authenticate", err)
	}

	expiryToken, ok := tokenMap["exp"].(float64)
	if !ok {
		return false, fmt.Errorf("%s: %w", "failed to authenticate", err)
	}

	iat, ok := tokenMap["iat"].(float64)
	if !ok {
		return false, fmt.Errorf("%s: %w", "failed to authenticate", err)

	}

	ttl := time.Duration(expiryToken-iat) * time.Second
	if ttl <= 0 {
		return false, fmt.Errorf("%s: %w", "failed to authenticate", err)
	}

	// Set expiry based on local clock.
	tm.setRefreshTimeFromTTL(ttl)
	tm.token.Store("Bearer " + resp.GetToken())

	return true, nil
}

func (tm *tokenManager) ScheduleRefresh(getConn func() (*grpc.ClientConn, error)) {
	if tm.refreshScheduled {
		tm.logger.WarnContext(context.Background(), "refresh already scheduled")
	}

	tm.logger.Debug("scheduling token refresh")

	tm.refreshScheduled = true
	timer := time.NewTimer(0)

	defer func() {
		timer.Stop()
		tm.stopRefreshChan <- struct{}{}
		tm.refreshScheduled = false

		tm.logger.Debug("stopped scheduled token refresh")
	}()

	for {
		conn, err := getConn()
		if err != nil {
			tm.logger.WarnContext(context.Background(), "failed to refresh token", slog.Any("error", err))
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

		_, err = tm.RefreshToken(ctx, conn)
		if err != nil {
			tm.logger.WarnContext(context.Background(), "failed to refresh token", slog.Any("error", err))
		}

		cancel()

		waitFor := time.Until(tm.refreshTime.Load().(time.Time)) - time.Second*5

		tm.logger.Debug("waiting to refresh token", slog.Duration("waitTime", waitFor))
		timer.Reset(waitFor)

		select {
		case <-timer.C:
		case <-tm.stopRefreshChan:
			return
		}
	}
}

func (tm *tokenManager) RequireTransportSecurity() bool {
	return true
}

func (tm *tokenManager) UnaryInterceptor() grpc.UnaryClientInterceptor {
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

func (tm *tokenManager) StreamInterceptor() grpc.StreamClientInterceptor {
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

func (jm *tokenManager) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "Authorization", jm.token.Load().(string))
}
