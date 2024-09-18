package avs

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnectionProvider(t *testing.T) {
	seeds := HostPortSlice{}
	listenerName := "listener"
	isLoadBalancer := false
	credentials := &UserPassCredentials{
		username: "admin",
		password: "password",
	}
	tlsConfig := &tls.Config{}
	var logger *slog.Logger

	cp, err := newConnectionProvider(context.Background(), seeds, &listenerName, isLoadBalancer, credentials, tlsConfig, logger)

	assert.Nil(t, cp)
	assert.Error(t, err, errors.New("seeds cannot be nil or empty"))
}
