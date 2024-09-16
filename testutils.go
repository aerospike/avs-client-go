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

type ServerTestBaseSuite struct {
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

func (suite *ServerTestBaseSuite) SetupSuite() {
	suite.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if suite.Name != "" {
		suite.Logger = suite.Logger.With("suite", suite.Name)
	}

	suite.CoverFile = path.Join(wd, "../coverage/client-coverage.cov")

	err := DockerComposeUp(suite.ComposeFile)

	time.Sleep(time.Second * 10)

	if err != nil {
		suite.FailNowf("unable to start docker compose up", "%v", err)
	}

	suite.Assert().NoError(err)

	suite.AvsClient, err = GetClient(suite.AvsHostPort, suite.AvsLB, suite.AvsCreds, suite.AvsTLSConfig, suite.Logger)
	if err != nil {
		suite.FailNowf("unable to create admin client", "%v", err)
	}
}

func (suite *ServerTestBaseSuite) TearDownSuite() {
	err := suite.AvsClient.Close()
	suite.Assert().NoError(err)

	err = DockerComposeDown(suite.ComposeFile)
	if err != nil {
		fmt.Println("unable to stop docker compose down")
	}

	goleak.VerifyNone(suite.T())
}

func GetStrPtr(str string) *string {
	ptr := str
	return &ptr
}

func GetUint32Ptr(i int) *uint32 {
	ptr := uint32(i)
	return &ptr
}

func GetUint64Ptr(i int) *uint64 {
	ptr := uint64(i)
	return &ptr
}

func GetFloat32Ptr(i float32) *float32 {
	ptr := float32(i)
	return &ptr
}

func GetBoolPtr(b bool) *bool {
	ptr := b
	return &ptr
}

func CreateFlagStr(name, value string) string {
	return fmt.Sprintf("--%s %s", name, value)
}

type IndexDefinitionBuilder struct {
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
	hnswCacheExpiry                *uint64
	hnswCacheMaxEntries            *uint64
	hnswHealerMaxScanPageSize      *uint32
	hnswHealerMaxScanRatePerSecond *uint32
	hnswHealerParallelism          *uint32
	HnswHealerReindexPercent       *float32
	HnswHealerSchedule             *string
	hnswMergeIndexParallelism      *uint32
	hnswMergeReIndexParallelism    *uint32
}

func NewIndexDefinitionBuilder(
	indexName,
	namespace string,
	dimension int,
	distanceMetric protos.VectorDistanceMetric,
	vectorField string,
) *IndexDefinitionBuilder {
	return &IndexDefinitionBuilder{
		indexName:            indexName,
		namespace:            namespace,
		dimension:            dimension,
		vectorDistanceMetric: distanceMetric,
		vectorField:          vectorField,
	}
}

func (idb *IndexDefinitionBuilder) WithSet(set string) *IndexDefinitionBuilder {
	idb.set = &set
	return idb
}

func (idb *IndexDefinitionBuilder) WithLabels(labels map[string]string) *IndexDefinitionBuilder {
	idb.labels = labels
	return idb
}

func (idb *IndexDefinitionBuilder) WithStorageNamespace(storageNamespace string) *IndexDefinitionBuilder {
	idb.storageNamespace = &storageNamespace
	return idb
}

func (idb *IndexDefinitionBuilder) WithStorageSet(storageSet string) *IndexDefinitionBuilder {
	idb.storageSet = &storageSet
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswM(m uint32) *IndexDefinitionBuilder {
	idb.hnsfM = &m
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswEf(ef uint32) *IndexDefinitionBuilder {
	idb.hnsfEf = &ef
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswEfConstruction(efConstruction uint32) *IndexDefinitionBuilder {
	idb.hnsfEfC = &efConstruction
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswMaxMemQueueSize(maxMemQueueSize uint32) *IndexDefinitionBuilder {
	idb.hnswMemQueueSize = &maxMemQueueSize
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswBatchingMaxRecord(maxRecord uint32) *IndexDefinitionBuilder {
	idb.hnsfBatchingMaxRecord = &maxRecord
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswBatchingInterval(interval uint32) *IndexDefinitionBuilder {
	idb.hnsfBatchingInterval = &interval
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswCacheExpiry(expiry uint64) *IndexDefinitionBuilder {
	idb.hnswCacheExpiry = &expiry
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswCacheMaxEntries(maxEntries uint64) *IndexDefinitionBuilder {
	idb.hnswCacheMaxEntries = &maxEntries
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswHealerMaxScanPageSize(maxScanPageSize uint32) *IndexDefinitionBuilder {
	idb.hnswHealerMaxScanPageSize = &maxScanPageSize
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswHealerMaxScanRatePerNode(maxScanRatePerSecond uint32) *IndexDefinitionBuilder {
	idb.hnswHealerMaxScanRatePerSecond = &maxScanRatePerSecond
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswHealerParallelism(parallelism uint32) *IndexDefinitionBuilder {
	idb.hnswHealerParallelism = &parallelism
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswHealerReindexPercent(reindexPercent float32) *IndexDefinitionBuilder {
	idb.HnswHealerReindexPercent = &reindexPercent
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswHealerScheduleDelay(schedule string) *IndexDefinitionBuilder {
	idb.HnswHealerSchedule = &schedule
	return idb
}

func (idb *IndexDefinitionBuilder) WithHnswMergeParallelism(mergeParallelism uint32) *IndexDefinitionBuilder {
	idb.hnswMergeIndexParallelism = &mergeParallelism
	return idb
}

func (idb *IndexDefinitionBuilder) Build() *protos.IndexDefinition {
	indexDef := &protos.IndexDefinition{
		Id: &protos.IndexId{
			Name:      idb.indexName,
			Namespace: idb.namespace,
		},
		Dimensions:           uint32(idb.dimension),
		VectorDistanceMetric: idb.vectorDistanceMetric,
		Field:                idb.vectorField,
		Type:                 protos.IndexType_HNSW,
		Storage: &protos.IndexStorage{
			Namespace: &idb.namespace,
			Set:       &idb.indexName,
		},
		Params: &protos.IndexDefinition_HnswParams{
			HnswParams: &protos.HnswParams{
				M:              GetUint32Ptr(16),
				EfConstruction: GetUint32Ptr(100),
				Ef:             GetUint32Ptr(100),
				BatchingParams: &protos.HnswBatchingParams{
					MaxRecords: GetUint32Ptr(100000),
					Interval:   GetUint32Ptr(30000),
				},
				CachingParams: &protos.HnswCachingParams{},
				HealerParams:  &protos.HnswHealerParams{},
				MergeParams:   &protos.HnswIndexMergeParams{},
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
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.BatchingParams.MaxRecords = idb.hnsfBatchingMaxRecord
	}

	if idb.hnsfBatchingInterval != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.BatchingParams.Interval = idb.hnsfBatchingInterval
	}

	if idb.hnswMemQueueSize != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.MaxMemQueueSize = idb.hnswMemQueueSize
	}

	if idb.hnswCacheExpiry != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.CachingParams.Expiry = idb.hnswCacheExpiry
	}

	if idb.hnswCacheMaxEntries != nil {
		indexDef.Params.(*protos.IndexDefinition_HnswParams).HnswParams.CachingParams.MaxEntries = idb.hnswCacheMaxEntries
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

func DockerComposeUp(composeFile string) error {
	fmt.Println("Starting docker containers " + composeFile)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "-lDEBUG", "compose", fmt.Sprintf("-f%s", composeFile), "up", "-d")
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

func DockerComposeDown(composeFile string) error {
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

func GetClient(
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

func GetCACert(cert string) (*x509.CertPool, error) {
	// read in file
	certBytes, err := os.ReadFile(cert)
	if err != nil {
		log.Fatalf("unable to read cert file %v", err)
		return nil, err
	}

	return client.LoadCACerts([][]byte{certBytes}), nil
}

func GetCertificates(certFile string, keyFile string) ([]tls.Certificate, error) {
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
