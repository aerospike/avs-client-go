//go:build integration
// +build integration

package avs

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
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

var keyStart = 100000

func getUniqueKey() string {
	keyStart++
	return fmt.Sprintf("key%d", keyStart)
}

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
				AvsLB:       true,
				AvsHostPort: avsHostPort,
			},
		},
		{
			ServerTestBaseSuite: ServerTestBaseSuite{
				ComposeFile: "docker/mtls/docker-compose.yml", // mutual tls
				AvsTLSConfig: &tls.Config{
					Certificates: certificates,
					RootCAs:      rootCA,
				},
				AvsHostPort: avsHostPort,
				AvsLB:       true,
			},
		},
		{
			ServerTestBaseSuite: ServerTestBaseSuite{
				ComposeFile: "docker/auth/docker-compose.yml", // tls + auth (auth requires tls)
				AvsCreds:    NewCredentialsFromUserPass("admin", "admin"),
				AvsTLSConfig: &tls.Config{
					Certificates: nil,
					RootCAs:      rootCA,
				},
				AvsHostPort: avsHostPort,
				AvsLB:       true,
			},
		},
	}

	testSuiteEnv := os.Getenv("ASVEC_TEST_SUITES")
	picked_suites := map[int]struct{}{}

	if testSuiteEnv != "" {
		testSuites := strings.Split(testSuiteEnv, ",")

		for _, s := range testSuites {
			i, err := strconv.Atoi(s)
			if err != nil {
				t.Fatalf("unable to convert %s to int: %v", s, err)
			}

			picked_suites[i] = struct{}{}
		}
	}

	logger.Info("Running test suites", slog.Any("suites", picked_suites))

	for i, s := range suites {
		if len(picked_suites) != 0 {
			if _, ok := picked_suites[i]; !ok {
				continue
			}
		}
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
			getUniqueKey(),
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

func (suite *SingleNodeTestSuite) TestBasicUpsertExistsDelete() {
	records := []struct {
		namespace  string
		set        *string
		key        any
		recordData map[string]any
	}{
		{
			"test",
			getUniqueSetName(),
			getUniqueKey(),
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

		exists, err := suite.AvsClient.Exists(ctx, rec.namespace, rec.set, rec.key)
		suite.NoError(err)

		if err != nil {
			return
		}

		suite.True(exists)

		err = suite.AvsClient.Delete(ctx, rec.namespace, rec.set, rec.key)
		suite.NoError(err)

		if err != nil {
			return
		}

		exists, err = suite.AvsClient.Exists(ctx, rec.namespace, rec.set, rec.key)
		suite.NoError(err)

		suite.False(exists)
	}
}

func (suite *SingleNodeTestSuite) TestIndexCreate() {
	testcases := []struct {
		name                 string
		namespace            string
		indexName            string
		vectorField          string
		dimension            uint32
		vectorDistanceMetric protos.VectorDistanceMetric
		opts                 *IndexCreateOpts
		expectedIndex        protos.IndexDefinition
	}{
		{
			"basic",
			"test",
			"index",
			"vector",
			10,
			protos.VectorDistanceMetric_SQUARED_EUCLIDEAN,
			nil,
			protos.IndexDefinition{
				Id: &protos.IndexId{
					Namespace: "test",
					Name:      "index",
				},
				Dimensions:           uint32(10),
				VectorDistanceMetric: protos.VectorDistanceMetric_SQUARED_EUCLIDEAN,
				Type:                 protos.IndexType_HNSW,
				SetFilter:            nil,
				Field:                "vector",
			},
		},
		{
			"with opts",
			"test",
			"index",
			"vector",
			10,
			protos.VectorDistanceMetric_COSINE,
			&IndexCreateOpts{
				Sets: []string{"testset"},
				Storage: &protos.IndexStorage{
					Namespace: GetStrPtr("storage-ns"),
					Set:       GetStrPtr("storage-set"),
				},
				Labels: map[string]string{
					"a": "b",
				},
			},
			protos.IndexDefinition{
				Id: &protos.IndexId{
					Namespace: "test",
					Name:      "index",
				},
				Dimensions:           uint32(10),
				VectorDistanceMetric: protos.VectorDistanceMetric_COSINE,
				Type:                 protos.IndexType_HNSW,
				SetFilter:            GetStrPtr("testset"),
				Field:                "vector",
				Storage: &protos.IndexStorage{
					Namespace: GetStrPtr("storage-ns"),
					Set:       GetStrPtr("storage-set"),
				},
				Labels: map[string]string{
					"a": "b",
				},
				Params: &protos.IndexDefinition_HnswParams{
					HnswParams: &protos.HnswParams{},
				},
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			err := suite.AvsClient.IndexCreate(ctx, tc.namespace, tc.indexName, tc.vectorField, tc.dimension, tc.vectorDistanceMetric, tc.opts)
			suite.NoError(err)

			if err != nil {
				return
			}

			defer suite.AvsClient.IndexDrop(ctx, tc.namespace, tc.indexName)

			index, err := suite.AvsClient.IndexGet(ctx, tc.namespace, tc.indexName, false)
			suite.NoError(err)

			if err != nil {
				return
			}

			suite.EqualExportedValues(tc.expectedIndex, *index)
		})
	}
}

// func (suite *SingleNodeTestSuite) TestIsIndexed() {
// 	records := []struct {
// 		namespace  string
// 		set        *string
// 		key        any
// 		recordData map[string]any
// 	}{
// 		{
// 			namespace: "test",
// 			set:       getUniqueSetName(),
// 			key:       getUniqueKey(),
// 			recordData: map[string]any{
// 				// "str":   "str",
// 				// "int":   int64(64),
// 				// "float": 3.14,
// 				// "bool":  false,
// 				"arr": getVectorFloat32(10, 1.0),
// 				// "map": map[any]any{
// 				// 	"foo": "bar",
// 				// },
// 			},
// 		},
// 	}

// 	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	// defer cancel()
// 	ctx := context.Background()

// 	for _, rec := range records {
// 		err := suite.AvsClient.Upsert(ctx, rec.namespace, rec.set, rec.key, rec.recordData, false)
// 		suite.NoError(err)

// 		for i := 0; i < 10; i++ {
// 			recordData := map[string]any{
// 				// "str":   "str",
// 				// "int":   int64(64),
// 				// "float": 3.14,
// 				// "bool":  false,
// 				"arr": getVectorFloat32(10, float32(i*2)),
// 				// "map": map[any]any{
// 				// 	"foo": "bar",
// 				// },
// 			}

// 			suite.AvsClient.Upsert(ctx, rec.namespace, rec.set, getUniqueKey(), recordData, false)
// 		}

// 		if err != nil {
// 			return
// 		}

// 		isIndexed, err := suite.AvsClient.IsIndexed(ctx, rec.namespace, rec.set, "index", rec.key)
// 		suite.Error(err)

// 		var indexOpts *IndexCreateOpts

// 		if rec.set != nil {
// 			indexOpts = &IndexCreateOpts{
// 				Sets: []string{*rec.set},
// 			}
// 		}

// 		err = suite.AvsClient.IndexCreate(
// 			ctx,
// 			rec.namespace,
// 			"index",
// 			"arr",
// 			10,
// 			protos.VectorDistanceMetric_SQUARED_EUCLIDEAN,
// 			indexOpts,
// 		)
// 		suite.NoError(err)

// 		if err != nil {
// 			return
// 		}

// 		isIndexed, err = suite.AvsClient.IsIndexed(ctx, rec.namespace, rec.set, "index", rec.key)
// 		suite.NoError(err)
// 		suite.False(isIndexed)

// 		defer suite.AvsClient.IndexDrop(ctx, rec.namespace, "index")

// 		suite.AvsClient.WaitForIndexCompletion(ctx, rec.namespace, "index", time.Second*12)

// 		isIndexed, err = suite.AvsClient.IsIndexed(ctx, rec.namespace, rec.set, "index", rec.key)
// 		suite.NoError(err)

// 		// time.Sleep(time.Second * 33330)

// 		suite.True(isIndexed)

// 		err = suite.AvsClient.Delete(ctx, rec.namespace, rec.set, rec.key)
// 		suite.NoError(err)
// 	}
// }

func (suite *SingleNodeTestSuite) TestFailsToInsertAlreadyExists() {
	ctx := context.Background()
	key := getUniqueKey()

	err := suite.AvsClient.Insert(ctx, "test", nil, key, map[string]any{"foo": "bar"}, false)
	suite.NoError(err)

	if err != nil {
		return
	}

	err = suite.AvsClient.Insert(ctx, "test", nil, key, map[string]any{"foo": "bar"}, false)
	suite.Error(err)
}

func (suite *SingleNodeTestSuite) TestFailsToUpdateRecordDNE() {
	ctx := context.Background()
	key := getUniqueKey()

	err := suite.AvsClient.Update(ctx, "test", nil, key, map[string]any{"foo": "bar"}, false)
	suite.Error(err)
}

func (suite *SingleNodeTestSuite) TestCanUpsertRecordTwice() {
	ctx := context.Background()
	key := getUniqueKey()

	err := suite.AvsClient.Upsert(ctx, "test", nil, key, map[string]any{"foo": "bar"}, false)
	suite.NoError(err)

	if err != nil {
		return
	}

	err = suite.AvsClient.Upsert(ctx, "test", nil, key, map[string]any{"foo": "bar"}, false)
	suite.NoError(err)
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

var userNameCount = -1

func getUniqueUserName() string {
	userNameCount++
	return "user" + fmt.Sprintf("%d", userNameCount)
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

			isIndexed, err := suite.AvsClient.IsIndexed(ctx, tc.namespace, setName, indexName, getKey(0))
			suite.NoError(err)
			suite.True(isIndexed)

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

			defer suite.AvsClient.IndexDrop(ctx, tc.namespace, indexName)

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

func (suite *SingleNodeTestSuite) TestConcurrentWrites() {
	wg := sync.WaitGroup{}
	numWrites := 10_000
	keys := []string{}

	for i := 0; i < numWrites; i++ {
		keys = append(keys, fmt.Sprintf("concurrent-key-%d", i))
	}

	for i := 0; i < numWrites; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			suite.AvsClient.Upsert(context.Background(), testNamespace, nil, keys[i], map[string]any{"foo": "bar"}, false)
		}(i)
	}

	wg.Wait()

	wg = sync.WaitGroup{}

	for i := 0; i < numWrites; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := suite.AvsClient.Get(context.Background(), testNamespace, nil, keys[i], nil, nil)
			suite.Assert().NoError(err)
		}(i)
	}

	wg.Wait()
}

func (suite *SingleNodeTestSuite) TestUserCreate() {
	suite.SkipIfUserPassAuthDisabled()

	ctx := context.Background()
	username := getUniqueUserName()

	err := suite.AvsClient.CreateUser(ctx, username, "test-password", []string{"read-write"})
	suite.NoError(err)

	actualUser, err := suite.AvsClient.GetUser(ctx, username)
	suite.NoError(err)

	expectedUser := protos.User{
		Username: username,
		Roles: []string{
			"read-write",
		},
	}

	suite.EqualExportedValues(expectedUser, *actualUser)
	return

}

func (suite *SingleNodeTestSuite) TestUserDelete() {
	suite.SkipIfUserPassAuthDisabled()

	ctx := context.Background()
	username := getUniqueUserName()

	err := suite.AvsClient.CreateUser(ctx, username, "test-password", []string{"read-write"})
	suite.NoError(err)

	err = suite.AvsClient.DropUser(ctx, username)
	suite.NoError(err)

	_, err = suite.AvsClient.GetUser(ctx, username)
	suite.Error(err)
}

func (suite *SingleNodeTestSuite) TestUserGrantRoles() {
	suite.SkipIfUserPassAuthDisabled()

	ctx := context.Background()
	username := getUniqueUserName()

	err := suite.AvsClient.CreateUser(ctx, username, "test-password", []string{"read-write"})
	suite.NoError(err)

	err = suite.AvsClient.GrantRoles(ctx, username, []string{"admin"})
	suite.NoError(err)

	actualUser, err := suite.AvsClient.GetUser(ctx, username)
	suite.NoError(err)

	expectedUser := protos.User{
		Username: username,
		Roles: []string{
			"admin",
			"read-write",
		},
	}

	suite.EqualExportedValues(expectedUser, *actualUser)
}

func (suite *SingleNodeTestSuite) TestUserRevokeRoles() {
	suite.SkipIfUserPassAuthDisabled()

	ctx := context.Background()
	username := getUniqueUserName()

	err := suite.AvsClient.CreateUser(ctx, username, "test-password", []string{"read-write", "admin"})
	suite.NoError(err)

	err = suite.AvsClient.RevokeRoles(ctx, username, []string{"admin"})
	suite.NoError(err)

	actualUser, err := suite.AvsClient.GetUser(ctx, username)
	suite.NoError(err)

	expectedUser := protos.User{
		Username: username,
		Roles: []string{
			"read-write",
		},
	}

	suite.EqualExportedValues(expectedUser, *actualUser)
}

func (suite *SingleNodeTestSuite) TestListUsers() {
	suite.SkipIfUserPassAuthDisabled()

	ctx := context.Background()
	username := getUniqueUserName()

	err := suite.AvsClient.CreateUser(ctx, username, "test-password", []string{"read-write"})
	suite.NoError(err)

	username1 := getUniqueUserName()

	err = suite.AvsClient.CreateUser(ctx, username1, "test-password", []string{"read-write"})
	suite.NoError(err)

	users, err := suite.AvsClient.ListUsers(ctx)
	suite.NoError(err)

	suite.Equal(3, len(users.Users))

	for _, user := range users.Users {
		if user.Username == username {
			suite.Equal([]string{"read-write"}, user.Roles)
		} else if user.Username == username1 {
			suite.Equal([]string{"read-write"}, user.Roles)
		} else {
			suite.Equal("admin", user.Username)
			suite.Equal([]string{"admin", "read-write"}, user.Roles)
		}
	}
}

func (suite *SingleNodeTestSuite) TestListRoles() {
	suite.SkipIfUserPassAuthDisabled()

	ctx := context.Background()
	roles, err := suite.AvsClient.ListRoles(ctx)
	suite.NoError(err)

	suite.Equal(2, len(roles.Roles))
}

func (suite *SingleNodeTestSuite) TestNodeIDs() {
	ctx := context.Background()
	nodeIDs := suite.AvsClient.NodeIDs(ctx)

	if suite.AvsLB {
		suite.Equal(0, len(nodeIDs))
	} else {
		suite.Equal(1, len(nodeIDs))
	}
}

func (suite *SingleNodeTestSuite) TestClusterEndpoints() {
	testCases := []struct {
		name              string
		nodeId            *protos.NodeId
		listenerName      *string
		expectedEndpoints []*protos.ServerEndpoint
		expectedErrMsg    *string
	}{
		{
			name:   "nil-node",
			nodeId: nil,
			expectedEndpoints: []*protos.ServerEndpoint{
				{
					Address: "127.0.0.1",
					Port:    10000,
					IsTls:   suite.AvsTLSConfig != nil,
				},
			},
		},
		{
			name: "node id DNE",
			nodeId: &protos.NodeId{
				Id: 1,
			},
			expectedErrMsg: GetStrPtr("failed to get cluster endpoints"),
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			endpoints, err := suite.AvsClient.ClusterEndpoints(ctx, tc.nodeId, tc.listenerName)

			if tc.expectedErrMsg != nil {
				suite.Error(err)
				suite.Contains(err.Error(), *tc.expectedErrMsg)
				return
			} else {
				suite.NoError(err)
			}

			for id, endpoint := range endpoints.Endpoints {
				// When LB is true we aren't able to validate the node-id since
				// we did not tend the cluster. When LB is false we are able to
				// get the node-id of the single node.
				if !suite.AvsLB {
					nodeId := suite.AvsClient.NodeIDs(ctx)[0]
					suite.Assert().Equal(nodeId.Id, id)
				}

				suite.EqualExportedValues(tc.expectedEndpoints[0], endpoint.Endpoints[0])
			}
		})
	}
}

func (suite *SingleNodeTestSuite) TestAbout() {
	testCases := []struct {
		name            string
		nodeId          *protos.NodeId
		expectedVersion string
		expectedErrMsg  *string
	}{
		{
			name:            "nil-node",
			nodeId:          nil,
			expectedVersion: "0.10.0",
		},
		{
			name: "node id DNE",
			nodeId: &protos.NodeId{
				Id: 1,
			},
			expectedErrMsg: GetStrPtr("failed to make about request"),
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			actualVersion, err := suite.AvsClient.About(ctx, tc.nodeId)

			if tc.expectedErrMsg != nil {
				suite.Error(err)
				suite.Contains(err.Error(), *tc.expectedErrMsg)
				return
			} else {
				suite.NoError(err)
			}

			suite.Equal(actualVersion.GetVersion(), tc.expectedVersion)
		})
	}
}
