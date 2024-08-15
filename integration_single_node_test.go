//go:build integration
// +build integration

package avs

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/aerospike/avs-client-go/protos"
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
		// {
		// 	ServerTestBaseSuite: ServerTestBaseSuite{
		// 		ComposeFile: "docker/multi-node/docker-compose.yml", // vanilla
		// 		// SuiteFlags: []string{
		// 		// 	"--log-level debug",
		// 		// 	"--timeout 10s",
		// 		// 	CreateFlagStr(flags.Seeds, avsHostPort.String()),
		// 		// },
		// 		AvsLB:       false,
		// 		AvsHostPort: avsHostPort,
		// 	},
		// },
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
				AvsLB:       true,
			},
		},
		// {
		// 	ServerTestBaseSuite: ServerTestBaseSuite{
		// 		ComposeFile: "docker/auth/docker-compose.yml", // tls + auth (auth requires tls)
		// 		// SuiteFlags: []string{
		// 		// 	"--log-level debug",
		// 		// 	"--timeout 10s",
		// 		// 	CreateFlagStr(flags.Host, avsHostPort.String()),
		// 		// 	CreateFlagStr(flags.TLSCaFile, "docker/auth/config/tls/ca.aerospike.com.crt"),
		// 		// 	CreateFlagStr(flags.AuthUser, "admin"),
		// 		// 	CreateFlagStr(flags.AuthPassword, "admin"),
		// 		// },
		// 		AvsCreds: NewCredentialsFromUserPass("admin", "admin"),
		// 		AvsTLSConfig: &tls.Config{
		// 			Certificates: nil,
		// 			RootCAs:      rootCA,
		// 		},
		// 		AvsHostPort: avsHostPort,
		// 		AvsLB:       true,
		// 	},
		// },
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

func (suite *SingleNodeTestSuite) TestBasicUpsertGetDelete() {
	records := []struct {
		namespace  string
		set        *string
		key        any
		recordData map[string]any
	}{
		{
			"test",
			getUniqueSetName(),
			"key1",
			map[string]any{
				"str":   "str",
				"int":   int64(64),
				"float": 3.14,
				"bool":  false,
				"arr":   []any{int64(0), int64(1), int64(2), int64(3)},
				"map": map[any]any{
					"foo": "bar",
				},
			},
		},
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	ctx := context.Background()

	for _, rec := range records {
		err := suite.AvsClient.Upsert(ctx, rec.namespace, rec.set, rec.key, rec.recordData, false)
		suite.NoError(err)

		if err != nil {
			return
		}

		record, err := suite.AvsClient.Get(ctx, rec.namespace, rec.set, rec.key, nil, nil)
		suite.NoError(err)

		if err != nil {
			return
		}

		suite.Equal(rec.recordData, record.Data)

		err = suite.AvsClient.Delete(ctx, rec.namespace, rec.set, rec.key)
		suite.NoError(err)

		if err != nil {
			return
		}

		_, err = suite.AvsClient.Get(context.TODO(), rec.namespace, rec.set, rec.key, nil, nil)
		suite.Error(err)
	}
}

func getVectorFloat32(length int, last float32) []float32 {
	vector := make([]float32, length)
	for i := 0; i < length-1; i++ {
		vector[i] = 1.0
	}

	vector[length-1] = last

	return vector
}

func getVectorBool(length int, last int) []bool {
	vector := make([]bool, length)
	for i := 0; i < last; i++ {
		vector[i] = true
	}

	return vector
}

func getKey(num int) string {
	return "key" + fmt.Sprintf("%d", num)
}

var indexNameCount = -1

func getUniqueIndexName() string {
	indexNameCount++
	return "index" + fmt.Sprintf("%d", indexNameCount)
}

var setNameCount = -1

func getUniqueSetName() *string {
	setNameCount++
	val := "set" + fmt.Sprintf("%d", setNameCount)
	return &val
}

func createNeighborFloat32(namespace string, set *string, key string, distance float32, vector []float32) *Neighbor {
	return &Neighbor{
		Namespace: namespace,
		Set:       set,
		Key:       key,
		Distance:  distance,
		Record: &Record{
			Data: map[string]any{
				"vector": vector,
			},
		},
	}
}

func createNeighborBool(namespace string, set *string, key string, distance float32, vector []bool) *Neighbor {
	return &Neighbor{
		Namespace: namespace,
		Set:       set,
		Key:       key,
		Distance:  distance,
		Record: &Record{
			Data: map[string]any{
				"vector": vector,
			},
		},
	}
}

func (suite *SingleNodeTestSuite) TestVectorSearchFloat32() {
	testCases := []struct {
		name              string
		namespace         string
		indexType         protos.VectorDistanceMetric
		vectors           [][]float32
		searchLimit       int
		query             []float32
		expectedNeighbors []*Neighbor
	}{
		{
			"",
			"test",
			protos.VectorDistanceMetric_SQUARED_EUCLIDEAN,
			[][]float32{
				getVectorFloat32(10, 0.0),
				getVectorFloat32(10, 1.0),
				getVectorFloat32(10, 2.0),
				getVectorFloat32(10, 3.0),
				getVectorFloat32(10, 4.0),
				getVectorFloat32(10, 5.0),
			},
			3,
			getVectorFloat32(10, 0.0),
			[]*Neighbor{
				createNeighborFloat32("test", nil, "key0", 0.0, getVectorFloat32(10, 0.0)),
				createNeighborFloat32("test", nil, "key1", 1.0, getVectorFloat32(10, 1.0)),
				createNeighborFloat32("test", nil, "key2", 4.0, getVectorFloat32(10, 2.0)),
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			setName := getUniqueSetName()

			for _, neighbor := range tc.expectedNeighbors {
				neighbor.Set = setName
			}

			// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			// defer cancel()
			ctx := context.Background()

			indexName := getUniqueIndexName()

			err := suite.AvsClient.IndexCreate(ctx, tc.namespace, indexName, "vector", 10, tc.indexType, &IndexCreateOpts{Sets: []string{*setName}})
			suite.NoError(err)

			if err != nil {
				return
			}

			for idx, vec := range tc.vectors {
				recordData := map[string]any{
					"vector": vec,
				}
				err := suite.AvsClient.Upsert(ctx, tc.namespace, setName, getKey(idx), recordData, false)
				suite.NoError(err)

				if err != nil {
					return
				}
			}

			err = suite.AvsClient.WaitForIndexCompletion(ctx, tc.namespace, indexName, time.Second*12)
			suite.NoError(err)

			if err != nil {
				return
			}

			time.Sleep(time.Second * 10)

			neighbors, err := suite.AvsClient.VectorSearchFloat32(ctx, tc.namespace, indexName, tc.query, 3, nil, nil, nil)

			suite.Equal(tc.expectedNeighbors, neighbors)

		})
	}
}

func (suite *SingleNodeTestSuite) TestVectorSearchBool() {
	testCases := []struct {
		name              string
		namespace         string
		indexType         protos.VectorDistanceMetric
		vectors           [][]bool
		searchLimit       int
		query             []bool
		expectedNeighbors []*Neighbor
	}{
		{
			"",
			"test",
			protos.VectorDistanceMetric_SQUARED_EUCLIDEAN,
			[][]bool{
				getVectorBool(10, 0),
				getVectorBool(10, 1),
				getVectorBool(10, 2),
				getVectorBool(10, 3),
				getVectorBool(10, 4),
				getVectorBool(10, 5),
			},
			3,
			getVectorBool(10, 0),
			[]*Neighbor{
				createNeighborBool("test", nil, "key0", 0.0, getVectorBool(10, 0)),
				createNeighborBool("test", nil, "key1", 1.0, getVectorBool(10, 1)),
				createNeighborBool("test", nil, "key2", 2.0, getVectorBool(10, 2)),
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			setName := getUniqueSetName()

			for _, neighbor := range tc.expectedNeighbors {
				neighbor.Set = setName
			}

			// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			// defer cancel()
			ctx := context.Background()

			indexName := getUniqueIndexName()

			err := suite.AvsClient.IndexCreate(ctx, tc.namespace, indexName, "vector", 10, tc.indexType, &IndexCreateOpts{Sets: []string{*setName}})
			suite.NoError(err)

			if err != nil {
				return
			}

			for idx, vec := range tc.vectors {
				recordData := map[string]any{
					"vector": vec,
				}
				err := suite.AvsClient.Upsert(ctx, tc.namespace, setName, getKey(idx), recordData, false)
				suite.NoError(err)

				if err != nil {
					return
				}
			}

			err = suite.AvsClient.WaitForIndexCompletion(ctx, tc.namespace, indexName, time.Second*12)
			suite.NoError(err)

			if err != nil {
				return
			}

			time.Sleep(time.Second * 10)

			neighbors, err := suite.AvsClient.VectorSearchBool(ctx, tc.namespace, indexName, tc.query, 3, nil, nil, nil)

			suite.Equal(tc.expectedNeighbors, neighbors)

		})
	}
}
