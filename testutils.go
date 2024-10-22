//go:build unit || integration || integration_large

package avs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/aerospike/avs-client-go/protos"
	"github.com/aerospike/tools-common-go/client"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"golang.org/x/net/context"
)

type serverTestBaseSuite struct {
	suite.Suite
	Name         string
	ComposeFile  string
	CoverFile    string
	AvsHostPort  *HostPort
	AvsLB        bool
	AvsTLSConfig *tls.Config
	AvsCreds     *UserPassCredentials
	AvsClient    *Client
	Logger       *slog.Logger
}

var wd, _ = os.Getwd()

func (suite *serverTestBaseSuite) SetupSuite() {
	suite.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if suite.Name != "" {
		suite.Logger = suite.Logger.With("suite", suite.Name)
	}

	suite.CoverFile = path.Join(wd, "../coverage/client-coverage.cov")

	err := dockerComposeUp(suite.ComposeFile)

	time.Sleep(time.Second * 10)

	if err != nil {
		suite.FailNowf("unable to start docker compose up", "%v", err)
	}

	suite.Assert().NoError(err)

	suite.AvsClient, err = getClient(suite.AvsHostPort, suite.AvsLB, suite.AvsCreds, suite.AvsTLSConfig, suite.Logger)
	if err != nil {
		suite.FailNowf("unable to create admin client", "%v", err)
	}
}

func (suite *serverTestBaseSuite) TearDownSuite() {
	err := suite.AvsClient.Close()
	suite.Assert().NoError(err)

	err = dockerComposeDown(suite.ComposeFile)
	if err != nil {
		fmt.Println("unable to stop docker compose down")
	}

	goleak.VerifyNone(suite.T())
}

func createFlagStr(name, value string) string {
	return fmt.Sprintf("--%s %s", name, value)
}

type indexDefinitionBuilder struct {
	indexName                      string
	namespace                      string
	set                            *string
	dimension                      int
	vectorDistanceMetric           protos.VectorDistanceMetric
	vectorField                    string
	storageNamespace               *string
	storageSet                     *string
	labels                         map[string]string
	hnsfM                          *uint32
	hnsfEfC                        *uint32
	hnsfEf                         *uint32
	hnswMemQueueSize               *uint32
	hnsfBatchingMaxRecord          *uint32
	hnsfBatchingInterval           *uint32
	hnswCacheExpiry                *int64
	hnswCacheMaxEntries            *uint64
	hnswHealerMaxScanPageSize      *uint32
	hnswHealerMaxScanRatePerSecond *uint32
	hnswHealerParallelism          *uint32
	HnswHealerReindexPercent       *float32
	HnswHealerSchedule             *string
	hnswMergeIndexParallelism      *uint32
	hnswMergeReIndexParallelism    *uint32
}

func newIndexDefinitionBuilder(
	indexName,
	namespace string,
	dimension int,
	distanceMetric protos.VectorDistanceMetric,
	vectorField string,
) *indexDefinitionBuilder {
	return &indexDefinitionBuilder{
		indexName:            indexName,
		namespace:            namespace,
		dimension:            dimension,
		vectorDistanceMetric: distanceMetric,
		vectorField:          vectorField,
	}
}

func (idb *indexDefinitionBuilder) WithSet(set string) *indexDefinitionBuilder {
	idb.set = &set
	return idb
}

func (idb *indexDefinitionBuilder) WithLabels(labels map[string]string) *indexDefinitionBuilder {
	idb.labels = labels
	return idb
}

func (idb *indexDefinitionBuilder) WithStorageNamespace(storageNamespace string) *indexDefinitionBuilder {
	idb.storageNamespace = &storageNamespace
	return idb
}

func (idb *indexDefinitionBuilder) WithStorageSet(storageSet string) *indexDefinitionBuilder {
	idb.storageSet = &storageSet
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswM(m uint32) *indexDefinitionBuilder {
	idb.hnsfM = &m
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswEf(ef uint32) *indexDefinitionBuilder {
	idb.hnsfEf = &ef
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswEfConstruction(efConstruction uint32) *indexDefinitionBuilder {
	idb.hnsfEfC = &efConstruction
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswMaxMemQueueSize(maxMemQueueSize uint32) *indexDefinitionBuilder {
	idb.hnswMemQueueSize = &maxMemQueueSize
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswBatchingMaxRecord(maxRecord uint32) *indexDefinitionBuilder {
	idb.hnsfBatchingMaxRecord = &maxRecord
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswBatchingInterval(interval uint32) *indexDefinitionBuilder {
	idb.hnsfBatchingInterval = &interval
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswCacheExpiry(expiry int64) *indexDefinitionBuilder {
	idb.hnswCacheExpiry = &expiry
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswCacheMaxEntries(maxEntries uint64) *indexDefinitionBuilder {
	idb.hnswCacheMaxEntries = &maxEntries
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswHealerMaxScanPageSize(maxScanPageSize uint32) *indexDefinitionBuilder {
	idb.hnswHealerMaxScanPageSize = &maxScanPageSize
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswHealerMaxScanRatePerNode(maxScanRatePerSecond uint32) *indexDefinitionBuilder {
	idb.hnswHealerMaxScanRatePerSecond = &maxScanRatePerSecond
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswHealerParallelism(parallelism uint32) *indexDefinitionBuilder {
	idb.hnswHealerParallelism = &parallelism
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswHealerReindexPercent(reindexPercent float32) *indexDefinitionBuilder {
	idb.HnswHealerReindexPercent = &reindexPercent
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswHealerScheduleDelay(schedule string) *indexDefinitionBuilder {
	idb.HnswHealerSchedule = &schedule
	return idb
}

func (idb *indexDefinitionBuilder) WithHnswMergeParallelism(mergeParallelism uint32) *indexDefinitionBuilder {
	idb.hnswMergeIndexParallelism = &mergeParallelism
	return idb
}

func (idb *indexDefinitionBuilder) Build() *protos.IndexDefinition {
	indexDef := &protos.IndexDefinition{
		Id: &protos.IndexId{
			Name:      idb.indexName,
			Namespace: idb.namespace,
		},
		Dimensions:           uint32(idb.dimension),
		VectorDistanceMetric: &idb.vectorDistanceMetric,
		Field:                idb.vectorField,
		Type:                 nil, // defaults to protos.IndexType_HNSW
		Storage: &protos.IndexStorage{
			Namespace: &idb.namespace,
			Set:       &idb.indexName,
		},
		Params: &protos.IndexDefinition_HnswParams{
			HnswParams: &protos.HnswParams{
				M:              ptr(uint32(16)),
				EfConstruction: ptr(uint32(100)),
				Ef:             ptr(uint32(100)),
				BatchingParams: &protos.HnswBatchingParams{
					MaxIndexRecords: ptr(uint32(100000)),
					IndexInterval:   ptr(uint32(30000)),
				},
				IndexCachingParams: &protos.HnswCachingParams{},
				HealerParams:       &protos.HnswHealerParams{},
				MergeParams:        &protos.HnswIndexMergeParams{},
			},
		},
	}

	indexDef.SetFilter = idb.set

	if idb.labels != nil {
		indexDef.Labels = idb.labels
	}

	if idb.storageNamespace != nil {
		indexDef.Storage.Namespace = idb.storageNamespace
	}

	if idb.storageSet != nil {
		indexDef.Storage.Set = idb.storageSet
	}

	if idb.hnsfM != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.M = idb.hnsfM
	}

	if idb.hnsfEf != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.Ef = idb.hnsfEf
	}

	if idb.hnsfEfC != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.EfConstruction = idb.hnsfEfC
	}

	if idb.hnsfBatchingMaxRecord != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.BatchingParams.MaxIndexRecords = idb.hnsfBatchingMaxRecord
	}

	if idb.hnsfBatchingInterval != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.BatchingParams.IndexInterval = idb.hnsfBatchingInterval
	}

	if idb.hnswMemQueueSize != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.MaxMemQueueSize = idb.hnswMemQueueSize
	}

	if idb.hnswCacheExpiry != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.IndexCachingParams.Expiry = idb.hnswCacheExpiry
	}

	if idb.hnswCacheMaxEntries != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.IndexCachingParams.MaxEntries = idb.hnswCacheMaxEntries
	}

	if idb.hnswHealerMaxScanPageSize != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.HealerParams.MaxScanPageSize = idb.hnswHealerMaxScanPageSize
	}

	if idb.hnswHealerMaxScanRatePerSecond != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.HealerParams.MaxScanRatePerNode = idb.hnswHealerMaxScanRatePerSecond
	}

	if idb.hnswHealerParallelism != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.HealerParams.Parallelism = idb.hnswHealerParallelism
	}

	if idb.HnswHealerReindexPercent != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.HealerParams.ReindexPercent = idb.HnswHealerReindexPercent
	}

	if idb.HnswHealerSchedule != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.HealerParams.Schedule = idb.HnswHealerSchedule
	}

	if idb.hnswMergeIndexParallelism != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.MergeParams.IndexParallelism = idb.hnswMergeIndexParallelism
	}

	if idb.hnswMergeReIndexParallelism != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.MergeParams.ReIndexParallelism = idb.hnswMergeReIndexParallelism
	}

	return indexDef
}

func dockerComposeUp(composeFile string) error {
	fmt.Println("Starting docker containers " + composeFile)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "-lDEBUG", "compose", fmt.Sprintf("-f%s", composeFile), "--env-file", "docker/.env", "up", "-d")
	err := cmd.Run()
	cmd.Wait()

	// fmt.Printf("docker compose up output: %s\n", string())

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return err
		}
		return err
	}

	return nil
}

func dockerComposeDown(composeFile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "compose", fmt.Sprintf("-f%s", composeFile), "down")
	_, err := cmd.Output()

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return err
		}
		return err
	}

	return nil
}

func getClient(
	avsHostPort *HostPort,
	avsLB bool,
	avsCreds *UserPassCredentials,
	avsTLSConfig *tls.Config,
	logger *slog.Logger,
) (*Client, error) {
	// Connect avs client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	var (
		avsClient *Client
		err       error
	)

	for {
		avsClient, err = NewClient(
			ctx,
			HostPortSlice{avsHostPort},
			nil,
			avsLB,
			avsCreds,
			avsTLSConfig,
			logger,
		)

		if err == nil {
			break
		}

		fmt.Printf("unable to create avs client %v", err)

		if ctx.Err() != nil {
			return nil, err
		}

		time.Sleep(time.Second)
	}

	// Wait for cluster to be ready
	for {
		_, err := avsClient.IndexList(ctx, false)
		if err == nil {
			break
		}

		fmt.Printf("waiting for the cluster to be ready %v", err)

		if ctx.Err() != nil {
			return nil, err
		}

		time.Sleep(time.Second)
	}

	return avsClient, nil
}

func getCACert(cert string) (*x509.CertPool, error) {
	// read in file
	certBytes, err := os.ReadFile(cert)
	if err != nil {
		log.Fatalf("unable to read cert file %v", err)
		return nil, err
	}

	return client.LoadCACerts([][]byte{certBytes}), nil
}

func getCertificates(certFile string, keyFile string) ([]tls.Certificate, error) {
	cert, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatalf("unable to read cert file %v", err)
		return nil, err
	}

	key, err := os.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("unable to read key file %v", err)
		return nil, err
	}

	return client.LoadServerCertAndKey([]byte(cert), []byte(key), nil)
}
