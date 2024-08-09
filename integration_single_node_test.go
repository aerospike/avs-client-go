//go:build integration || integration_large

package avs

import (
	"context"
	"crypto/tls"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

var (
	testNamespace = "test"
	testSet       = "testset"
	barNamespace  = "bar"
)

type SingleNodeTestSuite struct {
	ServerTestBaseSuite
}

func TestSingleNodeSuite(t *testing.T) {
	logger = logger.With(slog.Bool("test-logger", true)) // makes it easy to see which logger is which
	rootCA, err := GetCACert("docker/tls/config/tls/ca.aerospike.com.crt")
	if err != nil {
		t.Fatalf("unable to read root ca %v", err)
		t.FailNow()
		logger.Error("Failed to read cert")
	}

	certificates, err := GetCertificates("docker/mtls/config/tls/localhost.crt", "docker/mtls/config/tls/localhost.key")
	if err != nil {
		t.Fatalf("unable to read certificates %v", err)
		t.FailNow()
		logger.Error("Failed to read cert")
	}

	avsSeed := "localhost"
	avsPort := 10000
	avsHostPort := NewHostPort(avsSeed, avsPort)

	logger.Info("%v", slog.Any("cert", rootCA))
	suites := []*SingleNodeTestSuite{
		{
			ServerTestBaseSuite: ServerTestBaseSuite{
				ComposeFile: "docker/vanilla/docker-compose.yml", // vanilla
				// SuiteFlags: []string{
				// 	"--log-level debug",
				// 	"--timeout 10s",
				// 	CreateFlagStr(flags.Seeds, avsHostPort.String()),
				// },
				AvsHostPort: avsHostPort,
			},
		},
		{
			ServerTestBaseSuite: ServerTestBaseSuite{
				ComposeFile: "docker/tls/docker-compose.yml", // tls
				// SuiteFlags: []string{
				// 	"--log-level debug",
				// 	"--timeout 10s",
				// 	CreateFlagStr(flags.Seeds, avsHostPort.String()),
				// 	CreateFlagStr(flags.TLSCaFile, "docker/tls/config/tls/ca.aerospike.com.crt"),
				// },
				AvsTLSConfig: &tls.Config{
					Certificates: nil,
					RootCAs:      rootCA,
				},
				AvsHostPort: avsHostPort,
			},
		},
		{
			ServerTestBaseSuite: ServerTestBaseSuite{
				ComposeFile: "docker/mtls/docker-compose.yml", // mutual tls
				// SuiteFlags: []string{
				// 	"--log-level debug",
				// 	"--timeout 10s",
				// 	CreateFlagStr(flags.Host, avsHostPort.String()),
				// 	CreateFlagStr(flags.TLSCaFile, "docker/mtls/config/tls/ca.aerospike.com.crt"),
				// 	CreateFlagStr(flags.TLSCertFile, "docker/mtls/config/tls/localhost.crt"),
				// 	CreateFlagStr(flags.TLSKeyFile, "docker/mtls/config/tls/localhost.key"),
				// },
				AvsTLSConfig: &tls.Config{
					Certificates: certificates,
					RootCAs:      rootCA,
				},
				AvsHostPort: avsHostPort,
			},
		},
		{
			ServerTestBaseSuite: ServerTestBaseSuite{
				ComposeFile: "docker/auth/docker-compose.yml", // tls + auth (auth requires tls)
				// SuiteFlags: []string{
				// 	"--log-level debug",
				// 	"--timeout 10s",
				// 	CreateFlagStr(flags.Host, avsHostPort.String()),
				// 	CreateFlagStr(flags.TLSCaFile, "docker/auth/config/tls/ca.aerospike.com.crt"),
				// 	CreateFlagStr(flags.AuthUser, "admin"),
				// 	CreateFlagStr(flags.AuthPassword, "admin"),
				// },
				AvsCreds: NewCredentialsFromUserPass("admin", "admin"),
				AvsTLSConfig: &tls.Config{
					Certificates: nil,
					RootCAs:      rootCA,
				},
				AvsHostPort: avsHostPort,
			},
		},
	}

	for _, s := range suites {
		suite.Run(t, s)
	}
}

func (suite *SingleNodeTestSuite) SkipIfUserPassAuthDisabled() {
	if suite.AvsCreds == nil {
		suite.T().Skip("authentication is disabled. skipping test")
	}
}

func (suite *SingleNodeTestSuite) TestBasicUpsertGet() {
	records := []struct {
		namespace  string
		set        *string
		key        any
		recordData map[string]any
	}{}

	for _, rec := range records {
		err := suite.AvsClient.Upsert(context.TODO(), rec.namespace, rec.set, rec.key, rec.recordData, false)
		suite.NoError(err)
	}

	for _, rec := range records {
		record, err := suite.AvsClient.Get(context.TODO(), rec.namespace, rec.set, rec.key, nil, nil)

	}
}
