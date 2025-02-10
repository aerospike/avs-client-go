// Package avs provides a client for managing Aerospike Vector Indexes.
package avs

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/aerospike/avs-client-go/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	indexTimeoutDuration = time.Second * 100
	indexWaitDuration    = time.Millisecond * 100
)
const (
	failedToInsertRecord           = "failed to insert record"
	failedToGetRecord              = "failed to get record"
	failedToDeleteRecord           = "failed to delete record"
	failedToCheckRecordExists      = "failed to check if record exists"
	failedToCheckIsIndexed         = "failed to check if record is indexed"
	failedToWaitForIndexCompletion = "failed to wait for index completion"
)

type connProvider interface {
	GetNodeIDs() []uint64
	GetRandomConn() (*connection, error)
	GetSeedConn() (*connection, error)
	GetNodeConn(id uint64) (*connection, error)
	Close() error
}

// Client is a client for managing Aerospike Vector Indexes.
type Client struct {
	logger             *slog.Logger
	connectionProvider connProvider
	token              tokenManager
}

// NewClient creates a new Client instance.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	seeds (HostPortSlice): A list of seed hosts to connect to.
//	listenerName (*string): The name of the listener to connect to as configured on the server.
//	isLoadBalancer (bool): Whether the client should consider the seed a load balancer.
//	Only the first seed is considered. Subsequent seeds are ignored.
//	credentials (*UserPassCredentials): The credentials to authenticate with.
//	tlsConfig (*tls.Config): The TLS configuration to use for the connection.
//	logger (*slog.Logger): The logger to use for logging.
//
// Returns:
//
//	*Client: A new Client instance.
//	error: An error if the client creation fails, otherwise nil.
func NewClient(
	ctx context.Context,
	seeds HostPortSlice,
	listenerName *string,
	isLoadBalancer bool,
	credentials *UserPassCredentials,
	tlsConfig *tls.Config,
	logger *slog.Logger,
) (*Client, error) {
	logger = logger.WithGroup("avs")
	logger.Info("creating new client")

	var grpcToken tokenManager

	if credentials != nil {
		grpcToken = newGrpcJWTToken(credentials.username, credentials.password, logger)
	}

	connectionProvider, err := newConnectionProvider(
		ctx,
		seeds,
		listenerName,
		isLoadBalancer,
		grpcToken,
		tlsConfig,
		logger,
	)
	if err != nil {
		if grpcToken != nil {
			grpcToken.Close()
		}

		logger.Error("failed to create connection provider", slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc("failed to connect to server", err)
	}

	return newClient(connectionProvider, grpcToken, logger)
}

func newClient(
	connectionProvider connProvider,
	token tokenManager,
	logger *slog.Logger,
) (*Client, error) {
	return &Client{
		logger:             logger,
		token:              token,
		connectionProvider: connectionProvider,
	}, nil
}

// Close closes the Client and releases any resources associated with it.
//
// Returns:
//
//	error: An error if the closure fails, otherwise nil.
func (c *Client) Close() error {
	c.logger.Info("Closing client")

	if c.token != nil {
		c.token.Close()
	}

	return c.connectionProvider.Close()
}

func (c *Client) put(
	ctx context.Context,
	writeType protos.WriteType,
	namespace string,
	set *string,
	key any,
	recordData map[string]any,
	ignoreMemQueueFull bool,
) error {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.Any("key", key),
	)

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		logger.Error(failedToInsertRecord, slog.Any("error", err))
		return NewAVSError(failedToInsertRecord, err)
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		logger.Error(failedToInsertRecord, slog.Any("error", err))
		return NewAVSError(failedToInsertRecord, err)
	}

	fields, err := protos.ConvertToFields(recordData)
	if err != nil {
		logger.Error(failedToInsertRecord, slog.Any("error", err))
		return NewAVSError(failedToInsertRecord, err)
	}

	putReq := &protos.PutRequest{
		Key:                protoKey,
		WriteType:          &writeType,
		Fields:             fields,
		IgnoreMemQueueFull: ignoreMemQueueFull,
	}

	_, err = conn.transactClient.Put(ctx, putReq)
	if err != nil {
		logger.ErrorContext(ctx, failedToInsertRecord, slog.Any("error", err))
		return NewAVSErrorFromGrpc(failedToInsertRecord, err)
	}

	return err
}

// Insert inserts a new record into the specified namespace and set. If the
// record already exists, it will fail.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the record to insert.
//	set (*string): The set within the namespace where the record resides.
//	key (any): The key of the record to insert.
//	recordData (map[string]any): The data of the record to insert.
//	ignoreMemQueueFull (bool): Whether to ignore the in-memory queue full error.
//
// Returns:
//
//	error: An error if the insertion fails, otherwise nil.
func (c *Client) Insert(
	ctx context.Context,
	namespace string,
	set *string,
	key any,
	recordData map[string]any,
	ignoreMemQueueFull bool,
) error {
	c.logger.DebugContext(ctx, "inserting record", slog.String("namespace", namespace), slog.Any("key", key))
	return c.put(ctx, protos.WriteType_INSERT_ONLY, namespace, set, key, recordData, ignoreMemQueueFull)
}

// Update updates a record in the specified namespace and set. If the record
// does not already exist, it will fail.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the record to update.
//	set (*string): The set within the namespace where the record resides.
//	key (any): The key of the record to update.
//	recordData (map[string]any): The data of the record to update.
//	ignoreMemQueueFull (bool): Whether to ignore the in-memory queue full error.
//
// Returns:
//
//	error: An error if the update fails, otherwise nil.
func (c *Client) Update(
	ctx context.Context,
	namespace string,
	set *string,
	key any,
	recordData map[string]any,
	ignoreMemQueueFull bool,
) error {
	c.logger.DebugContext(ctx, "updating record", slog.String("namespace", namespace), slog.Any("key", key))
	return c.put(ctx, protos.WriteType_UPDATE_ONLY, namespace, set, key, recordData, ignoreMemQueueFull)
}

// Upsert inserts or updates a record into the specified namespace and set,
// regardless of whether the record exists.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the record to upsert.
//	set (*string): The set within the namespace where the record resides.
//	key (any): The key of the record to upsert.
//	recordData (map[string]any): The data of the record to upsert.
//	ignoreMemQueueFull (bool): Whether to ignore the in-memory queue full error.
//
// Returns:
//
//	error: An error if the upsert fails, otherwise nil.
func (c *Client) Upsert(
	ctx context.Context,
	namespace string,
	set *string,
	key any,
	recordData map[string]any,
	ignoreMemQueueFull bool,
) error {
	c.logger.DebugContext(
		ctx,
		"upserting record",
		slog.String("namespace", namespace),
		slog.Any("set", set),
		slog.Any("key", key),
	)

	return c.put(ctx, protos.WriteType_UPSERT, namespace, set, key, recordData, ignoreMemQueueFull)
}

// Get retrieves a record from the specified namespace and set. If the record
// does not exist, an error is returned.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the record to retrieve.
//	set (*string): The set within the namespace where the record resides.
//	key (any): The key of the record to retrieve.
//	includeFields ([]string): Fields to include in the response. Default is all.
//	excludeFields ([]string): Fields to exclude from the response. Default is none.
//
// Returns:
//
//	*Record: The retrieved record.
//	error: An error if the retrieval fails, otherwise nil.
func (c *Client) Get(ctx context.Context,
	namespace string,
	set *string,
	key any,
	includeFields []string,
	excludeFields []string,
) (*Record, error) {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.Any("key", key),
	)
	logger.DebugContext(ctx, "getting record")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		logger.Error(failedToGetRecord, slog.Any("error", err))
		return nil, NewAVSError(failedToGetRecord, err)
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		logger.Error(failedToGetRecord, slog.Any("error", err))
		return nil, NewAVSError(failedToGetRecord, err)
	}

	getReq := &protos.GetRequest{
		Key:        protoKey,
		Projection: createProjectionSpec(includeFields, excludeFields),
	}

	record, err := conn.transactClient.Get(ctx, getReq)
	if err != nil {
		logger.ErrorContext(ctx, failedToGetRecord, slog.Any("error", err))
		return nil, NewAVSErrorFromGrpc(failedToGetRecord, err)
	}

	logger.Debug("received record", slog.Any("record", record.Fields))

	return newRecordFromProto(record), nil
}

// Delete removes a record from the database.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the record to delete.
//	set (*string): The set within the namespace where the record resides.
//	key (any): The key of the record to delete.
//
// Returns:
//
//	error: An error if the deletion fails, otherwise nil.
func (c *Client) Delete(ctx context.Context, namespace string, set *string, key any) error {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.Any("key", key),
	)
	logger.DebugContext(ctx, "deleting record")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		logger.Error(failedToDeleteRecord, slog.Any("error", err))
		return NewAVSError(failedToDeleteRecord, err)
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		logger.Error(failedToDeleteRecord, slog.Any("error", err))
		return NewAVSError(failedToDeleteRecord, err)
	}

	getReq := &protos.DeleteRequest{
		Key: protoKey,
	}

	_, err = conn.transactClient.Delete(ctx, getReq)
	if err != nil {
		logger.Error(failedToDeleteRecord, slog.Any("error", err))
		return NewAVSErrorFromGrpc(failedToDeleteRecord, err)
	}

	return nil
}

// Exists checks if a record exists in the specified namespace and set.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the record to check.
//	set (*string): The set within the namespace where the record resides.
//	key (any): The key of the record to check.
//
// Returns:
//
//	bool: True if the record exists, false otherwise.
//	error: An error if the existence check fails, otherwise nil.
func (c *Client) Exists(
	ctx context.Context,
	namespace string,
	set *string,
	key any,
) (bool, error) {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.Any("key", key),
	)
	logger.DebugContext(ctx, "checking if record exists")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		logger.Error(failedToCheckRecordExists, slog.Any("error", err))
		return false, NewAVSError(failedToCheckRecordExists, err)
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		logger.Error(failedToCheckRecordExists, slog.Any("error", err))
		return false, NewAVSError(failedToCheckRecordExists, err)
	}

	existsReq := &protos.ExistsRequest{
		Key: protoKey,
	}

	boolean, err := conn.transactClient.Exists(ctx, existsReq)
	if err != nil {
		logger.Error(failedToCheckRecordExists, slog.Any("error", err))
		return false, NewAVSErrorFromGrpc(failedToCheckRecordExists, err)
	}

	return boolean.GetValue(), nil
}

// IsIndexed checks if a record is indexed in the specified namespace and index.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the record to check.
//	set (*string): The set within the namespace where the record resides.
//	indexName (string): The name of the index.
//	key (any): The key of the record to check.
//
// Returns:
//
//	bool: True if the record is indexed, false otherwise.
//	error: An error if the index check fails, otherwise nil.
func (c *Client) IsIndexed(
	ctx context.Context,
	namespace string,
	set *string,
	indexName string,
	key any,
) (bool, error) {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.String("indexName", indexName),
		slog.Any("key", key),
	)
	logger.DebugContext(ctx, "checking if record is indexed")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		logger.Error(failedToCheckIsIndexed, slog.Any("error", err))
		return false, NewAVSError(failedToCheckIsIndexed, err)
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		logger.Error(failedToCheckIsIndexed, slog.Any("error", err))
		return false, NewAVSError(failedToCheckIsIndexed, err)
	}

	isIndexedReq := &protos.IsIndexedRequest{
		IndexId: &protos.IndexId{
			Namespace: namespace,
			Name:      indexName,
		},
		Key: protoKey,
	}

	boolean, err := conn.transactClient.IsIndexed(ctx, isIndexedReq)
	if err != nil {
		logger.Error(failedToCheckIsIndexed, slog.Any("error", err))
		return false, NewAVSErrorFromGrpc(failedToCheckIsIndexed, err)
	}

	return boolean.GetValue(), nil
}

// vectorSearch searches for the nearest neighbors of the query vector in the
// specified index. It only deals with the protos.Vector type and is not
// intended to be exported.
func (c *Client) vectorSearch(ctx context.Context,
	namespace,
	indexName string,
	vector *protos.Vector,
	limit uint32,
	searchParams *protos.HnswSearchParams,
	includeFields []string,
	excludeFields []string,
) ([]*Neighbor, error) {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))
	logger.DebugContext(ctx, "searching for vector")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to search for vector"
		logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSError(msg, err)
	}

	vectorSearchReq := createVectorSearchRequest(
		namespace,
		indexName,
		vector,
		limit,
		searchParams,
		createProjectionSpec(includeFields, excludeFields),
	)

	resp, err := conn.transactClient.VectorSearch(ctx, vectorSearchReq)
	if err != nil {
		msg := "failed to search for vector"
		logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSError(msg, err)
	}

	neighbors := make([]*Neighbor, 0, limit)

	for {
		protoNeigh, err := resp.Recv()

		if err != nil {
			if !errors.Is(err, io.EOF) {
				msg := "failed to receive all neighbors"
				logger.Error(msg, slog.Any("error", err))

				return neighbors, NewAVSError(msg, err)
			}

			logger.DebugContext(ctx, "received all neighbors", slog.Any("neighbors", neighbors))

			return neighbors, nil
		}

		n, err := newNeighborFromProto(protoNeigh)
		if err != nil {
			msg := "failed to convert neighbor"
			logger.Error(msg, slog.Any("error", err))

			return neighbors, NewAVSError(msg, err)
		}

		neighbors = append(neighbors, n)
	}
}

// VectorSearchFloat32 searches for the nearest neighbors of the query vector in
// the specified index.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the index.
//	indexName (string): The name of the index.
//	query ([]float32): The query vector.
//	limit (uint32): The maximum number of neighbors to return.
//	searchParams (*protos.HnswSearchParams): Extra options to configure the behavior of the HNSW algorithm.
//	includeFields ([]string): Fields to include in the response. Default is all.
//	excludeFields ([]string): Fields to exclude from the response. Default is none.
//
// Returns:
//
//	[]*Neighbor: A list of the nearest neighbors.
//	error: An error if the search fails, otherwise nil.
func (c *Client) VectorSearchFloat32(
	ctx context.Context,
	namespace,
	indexName string,
	query []float32,
	limit uint32,
	searchParams *protos.HnswSearchParams,
	includeFields,
	excludeFields []string,
) ([]*Neighbor, error) {
	c.logger.DebugContext(ctx, "searching for float vector")

	vector := protos.CreateFloat32Vector(query)

	return c.vectorSearch(
		ctx,
		namespace,
		indexName,
		vector,
		limit,
		searchParams,
		includeFields,
		excludeFields,
	)
}

// VectorSearchBool searches for the nearest neighbors of the query vector in
// the specified index.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the index.
//	indexName (string): The name of the index.
//	query ([]bool): The query vector.
//	limit (uint32): The maximum number of neighbors to return.
//	searchParams (*protos.HnswSearchParams): Extra options to configure the behavior of the HNSW algorithm.
//	includeFields ([]string): Fields to include in the response. Default is all.
//	excludeFields ([]string): Fields to exclude from the response. Default is none.
//
// Returns:
//
//	[]*Neighbor: A list of the nearest neighbors.
//	error: An error if the search fails, otherwise nil.
func (c *Client) VectorSearchBool(ctx context.Context,
	namespace,
	indexName string,
	query []bool,
	limit uint32,
	searchParams *protos.HnswSearchParams,
	includeFields,
	excludeFields []string,
) ([]*Neighbor, error) {
	c.logger.DebugContext(ctx, "searching for bool vector")

	vector := protos.CreateBoolVector(query)

	return c.vectorSearch(
		ctx,
		namespace,
		indexName,
		vector,
		limit,
		searchParams,
		includeFields,
		excludeFields,
	)
}

// WaitForIndexCompletion waits for an index to be fully created. USE WITH
// CAUTION. This API will likely break in the future as a mechanism for
// correctly determining index completion is figured out. The wait interval
// determines the number of times the unmerged record queue must be of size 0.
// If the unmerged queue size is not zero, it will immediately return the next time a
// the unmerged queue size is zero. We have found that a good waitInterval is
// roughly 12 seconds. If you are expected to have many concurrent writers or a
// consistent stream of writes this will cause your application to hang.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the index.
//	indexName (string): The name of the index.
//	waitInterval (time.Duration): The interval to wait between checks.
//
// Returns:
//
//	error: An error if waiting for the index to complete fails, otherwise nil.
func (c *Client) WaitForIndexCompletion(
	ctx context.Context,
	namespace,
	indexName string,
	waitInterval time.Duration,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))
	logger.DebugContext(ctx, "waiting for index completion")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := failedToWaitForIndexCompletion
		logger.Error(msg, slog.Any("error", err))

		return NewAVSError(msg, err)
	}

	indexStatusReq := createIndexStatusRequest(namespace, indexName)
	timer := time.NewTimer(waitInterval)
	startTime := time.Now()
	unmergedZeroCount := 0
	unmergedNotZeroCount := 0

	defer timer.Stop()

	for {
		indexStatus, err := conn.indexClient.GetStatus(ctx, indexStatusReq)
		if err != nil {
			msg := failedToWaitForIndexCompletion
			logger.ErrorContext(ctx, msg, slog.Any("error", err))

			return NewAVSError(msg, err)
		}

		// We consider the index completed when unmerged record count == 0 for
		// 2 cycles OR unmerged record count != 0 and now it == 0.
		unmerged := indexStatus.GetUnmergedRecordCount()
		if unmerged == 0 {
			if unmergedZeroCount >= 2 || unmergedNotZeroCount >= 1 {
				logger.DebugContext(ctx, "index completed", slog.Duration("duration", time.Since(startTime)))
				return nil
			}

			unmergedZeroCount++

			logger.DebugContext(ctx, "index not yet completed", slog.Int64("unmerged", unmerged))
		} else {
			logger.DebugContext(ctx, "index not yet completed", slog.Int64("unmerged", unmerged))

			unmergedNotZeroCount++
		}

		timer.Reset(waitInterval)

		select {
		case <-timer.C:
		case <-ctx.Done():
			msg := "waiting for index completion canceled"

			logger.ErrorContext(ctx, msg)

			return NewAVSError(msg, ctx.Err())
		}
	}
}

// IndexCreateOpts are optional fields to further configure the behavior of your index.
//   - Sets: The sets to create the index on. Currently, only one set is supported.
//   - Metadata: Extra metadata that can be attached to the index.
//   - Storage: The storage configuration for the index. This allows you to
//     configure your index and data to be stored in separate namespaces and/or sets.
//   - HNSWParams: Extra options sent to the server to configure behavior of the
//     HNSW algorithm.
type IndexCreateOpts struct {
	Storage    *protos.IndexStorage
	HnswParams *protos.HnswParams
	Mode       *protos.IndexMode
	Labels     map[string]string
	Sets       []string
}

// IndexCreate creates a new Aerospike Vector Index.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the index.
//	indexName (string): The name of the index.
//	vectorField (string): The field to create the index on.
//	dimensions (uint32): The number of dimensions in the vector.
//	vectorDistanceMetric (protos.VectorDistanceMetric): The distance metric to use for the index.
//	opts (*IndexCreateOpts): Optional fields to configure the index.
//
// Returns:
//
//	error: An error if the index creation fails, otherwise nil.
func (c *Client) IndexCreate(
	ctx context.Context,
	namespace string,
	indexName string,
	vectorField string,
	dimensions uint32,
	vectorDistanceMetric protos.VectorDistanceMetric,
	opts *IndexCreateOpts,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))
	logger.DebugContext(ctx, "creating index")

	var (
		set     *string
		params  *protos.IndexDefinition_HnswParams
		labels  map[string]string
		storage *protos.IndexStorage
		mode    *protos.IndexMode
	)

	if opts != nil {
		if len(opts.Sets) > 0 {
			set = &opts.Sets[0]

			if len(opts.Sets) > 1 {
				logger.Warn(
					"multiple sets not yet supported for index creation, only the first set will be used",
					slog.String("set", *set),
				)
			}
		}

		params = &protos.IndexDefinition_HnswParams{HnswParams: opts.HnswParams}
		labels = opts.Labels
		storage = opts.Storage
		mode = opts.Mode
	}

	indexDef := &protos.IndexDefinition{
		Id: &protos.IndexId{
			Namespace: namespace,
			Name:      indexName,
		},
		Dimensions:           dimensions,
		VectorDistanceMetric: &vectorDistanceMetric,
		Field:                vectorField,
		SetFilter:            set,
		Params:               params,
		Labels:               labels,
		Storage:              storage,
		Type:                 nil, // defaults to protos.IndexType_HNSW
		Mode:                 mode,
	}

	return c.IndexCreateFromIndexDef(ctx, indexDef)
}

// IndexCreateFromIndexDef creates an HNSW index in AVS from an index definition.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	indexDef (*protos.IndexDefinition): The index definition to create the index from.
//
// Returns:
//
//	error: An error if the index creation fails, otherwise nil.
func (c *Client) IndexCreateFromIndexDef(
	ctx context.Context,
	indexDef *protos.IndexDefinition,
) error {
	logger := c.logger.With(slog.Any("definition", indexDef))
	logger.DebugContext(ctx, "creating index from definition")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to create index from definition"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSError(msg, err)
	}

	indexCreateReq := &protos.IndexCreateRequest{
		Definition: indexDef,
	}

	_, err = conn.indexClient.Create(ctx, indexCreateReq)
	if err != nil {
		msg := "failed to create index from definition"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	ctx, cancel := context.WithTimeout(ctx, indexTimeoutDuration)
	defer cancel()

	return c.waitForIndexCreation(ctx, indexDef.Id.Namespace, indexDef.Id.Name, indexWaitDuration)
}

// IndexUpdate updates an existing Aerospike Vector Index's dynamic configuration parameters.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the index.
//	name (string): The name of the index.
//	labels (map[string]string): Labels to update on the index.
//	hnswParams (*protos.HnswIndexUpdate): The HNSW parameters to update.
//	mode (*protos.IndexMode): The mode to change the index to.
//
// Returns:
//
//	error: An error if the update fails, otherwise nil.
func (c *Client) IndexUpdate(
	ctx context.Context,
	namespace string,
	indexName string,
	labels map[string]string,
	hnswParams *protos.HnswIndexUpdate,
	mode *protos.IndexMode,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))

	logger.DebugContext(ctx, "updating index")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to update index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	indexUpdate := &protos.IndexUpdateRequest{
		IndexId: &protos.IndexId{
			Namespace: namespace,
			Name:      indexName,
		},
		Labels: labels,
		Update: &protos.IndexUpdateRequest_HnswIndexUpdate{
			HnswIndexUpdate: hnswParams,
		},
		Mode: mode,
	}

	_, err = conn.indexClient.Update(ctx, indexUpdate)
	if err != nil {
		msg := "failed to update index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// IndexDrop drops an existing Aerospike Vector Index and blocks until it is fully removed.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the index to drop.
//	indexName (string): The name of the index to drop.
//
// Returns:
//
//	error: An error if the index drop fails, otherwise nil.
func (c *Client) IndexDrop(ctx context.Context, namespace, indexName string) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))
	logger.DebugContext(ctx, "dropping index")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to drop index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	indexDropReq := &protos.IndexDropRequest{
		IndexId: &protos.IndexId{
			Namespace: namespace,
			Name:      indexName,
		},
	}

	_, err = conn.indexClient.Drop(ctx, indexDropReq)
	if err != nil {
		msg := "failed to drop index"

		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	ctx, cancel := context.WithTimeout(ctx, indexTimeoutDuration)
	defer cancel()

	return c.waitForIndexDrop(ctx, namespace, indexName, indexWaitDuration)
}

// IndexList returns a list of all Aerospike Vector Indexes.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	applyDefaults (bool): Whether to apply server default values to the index definitions.
//
// Returns:
//
//	*protos.IndexDefinitionList: A list of index definitions.
//	error: An error if the list retrieval fails, otherwise nil.
func (c *Client) IndexList(ctx context.Context, applyDefaults bool) (*protos.IndexDefinitionList, error) {
	c.logger.DebugContext(ctx, "listing indexes")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get indexes"

		c.logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	indexListReq := &protos.IndexListRequest{
		ApplyDefaults: &applyDefaults,
	}

	indexList, err := conn.indexClient.List(ctx, indexListReq)
	if err != nil {
		msg := "failed to get indexes"

		c.logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexList, nil
}

// IndexGet returns the definition of an Aerospike Vector Index.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the index.
//	indexName (string): The name of the index.
//	applyDefaults (bool): Whether to apply server default values to the index definition.
//
// Returns:
//
//	*protos.IndexDefinition: The index definition.
//	error: An error if the retrieval fails, otherwise nil.
func (c *Client) IndexGet(
	ctx context.Context,
	namespace,
	indexName string,
	applyDefaults bool,
) (*protos.IndexDefinition, error) {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))
	logger.DebugContext(ctx, "getting index")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get index"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	indexGetReq := &protos.IndexGetRequest{
		IndexId: &protos.IndexId{
			Namespace: namespace,
			Name:      indexName,
		},
		ApplyDefaults: &applyDefaults,
	}

	indexDef, err := conn.indexClient.Get(ctx, indexGetReq)
	if err != nil {
		msg := "failed to get index"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexDef, nil
}

// IndexGetStatus returns the status of an Aerospike Vector Index.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the index.
//	indexName (string): The name of the index.
//
// Returns:
//
//	*protos.IndexStatusResponse: The index status.
//	error: An error if the retrieval fails, otherwise nil.
func (c *Client) IndexGetStatus(ctx context.Context, namespace, indexName string) (*protos.IndexStatusResponse, error) {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))
	logger.DebugContext(ctx, "getting index status")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get index status"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	indexStatusReq := createIndexStatusRequest(namespace, indexName)

	indexStatus, err := conn.indexClient.GetStatus(ctx, indexStatusReq)
	if err != nil {
		msg := "failed to get index status"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexStatus, nil
}

// GcInvalidVertices garbage collects invalid vertices in an Aerospike Vector Index.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	namespace (string): The namespace of the index.
//	indexName (string): The name of the index.
//	cutoffTime (time.Time): The cutoff time for the garbage collection.
//
// Returns:
//
//	error: An error if the garbage collection fails, otherwise nil.
func (c *Client) GcInvalidVertices(ctx context.Context, namespace, indexName string, cutoffTime time.Time) error {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.String("indexName", indexName),
		slog.Any("cutoffTime", cutoffTime),
	)

	logger.DebugContext(ctx, "garbage collection invalid vertices")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to garbage collect invalid vertices"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	gcRequest := &protos.GcInvalidVerticesRequest{
		IndexId: &protos.IndexId{
			Namespace: namespace,
			Name:      indexName,
		},
		CutoffTimestamp: cutoffTime.Unix(),
	}

	_, err = conn.indexClient.GcInvalidVertices(ctx, gcRequest)
	if err != nil {
		msg := "failed to garbage collect invalid vertices"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// CreateUser creates a new user with the provided username, password, and roles.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	username (string): The username of the new user.
//	password (string): The password for the new user.
//	roles ([]string): The roles to assign to the new user.
//
// Returns:
//
//	error: An error if the user creation fails, otherwise nil.
func (c *Client) CreateUser(ctx context.Context, username, password string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))
	logger.DebugContext(ctx, "creating user")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to create user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	addUserRequest := &protos.AddUserRequest{
		Credentials: createUserPassCredential(username, password),
		Roles:       roles,
	}

	_, err = conn.userAdminClient.AddUser(ctx, addUserRequest)
	if err != nil {
		msg := "failed to create user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// UpdateCredentials updates the password for the provided username.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	username (string): The username of the user.
//	password (string): The new password for the user.
//
// Returns:
//
//	error: An error if the credential update fails, otherwise nil.
func (c *Client) UpdateCredentials(ctx context.Context, username, password string) error {
	logger := c.logger.With(slog.String("username", username))
	logger.DebugContext(ctx, "updating user credentials")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to update user credentials"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	updatedCredRequest := &protos.UpdateCredentialsRequest{
		Credentials: createUserPassCredential(username, password),
	}

	_, err = conn.userAdminClient.UpdateCredentials(ctx, updatedCredRequest)
	if err != nil {
		msg := "failed to update user credentials"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// DropUser deletes the user with the provided username.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	username (string): The username of the user to delete.
//
// Returns:
//
//	error: An error if the user deletion fails, otherwise nil.
func (c *Client) DropUser(ctx context.Context, username string) error {
	logger := c.logger.With(slog.String("username", username))
	logger.DebugContext(ctx, "dropping user")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to drop user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	dropUserRequest := &protos.DropUserRequest{
		Username: username,
	}

	_, err = conn.userAdminClient.DropUser(ctx, dropUserRequest)
	if err != nil {
		msg := "failed to drop user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// GetUser returns the user with the provided username.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	username (string): The username of the user to retrieve.
//
// Returns:
//
//	*protos.User: The retrieved user.
//	error: An error if the user retrieval fails, otherwise nil.
func (c *Client) GetUser(ctx context.Context, username string) (*protos.User, error) {
	logger := c.logger.With(slog.String("username", username))
	logger.DebugContext(ctx, "getting user")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	getUserRequest := &protos.GetUserRequest{
		Username: username,
	}

	userResp, err := conn.userAdminClient.GetUser(ctx, getUserRequest)
	if err != nil {
		msg := "failed to get user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return userResp, nil
}

// ListUsers returns a list of all users.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//
// Returns:
//
//	*protos.ListUsersResponse: A list of all users.
//	error: An error if the list retrieval fails, otherwise nil.
func (c *Client) ListUsers(ctx context.Context) (*protos.ListUsersResponse, error) {
	c.logger.DebugContext(ctx, "listing users")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to list users"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	usersResp, err := conn.userAdminClient.ListUsers(ctx, &emptypb.Empty{})
	if err != nil {
		msg := "failed to list users"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return usersResp, nil
}

// GrantRoles grants the provided roles to the user with the provided username.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	username (string): The username of the user.
//	roles ([]string): The roles to grant to the user.
//
// Returns:
//
//	error: An error if the role grant fails, otherwise nil.
func (c *Client) GrantRoles(ctx context.Context, username string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))
	logger.DebugContext(ctx, "granting user roles")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to grant user roles"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	grantRolesRequest := &protos.GrantRolesRequest{
		Username: username,
		Roles:    roles,
	}

	_, err = conn.userAdminClient.GrantRoles(ctx, grantRolesRequest)
	if err != nil {
		msg := "failed to grant user roles"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// RevokeRoles revokes the provided roles from the user with the provided username.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	username (string): The username of the user.
//	roles ([]string): The roles to revoke from the user.
//
// Returns:
//
//	error: An error if the role revocation fails, otherwise nil.
func (c *Client) RevokeRoles(ctx context.Context, username string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))
	logger.DebugContext(ctx, "revoking user roles")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to revoke user roles"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	revokeRolesReq := &protos.RevokeRolesRequest{
		Username: username,
		Roles:    roles,
	}

	_, err = conn.userAdminClient.RevokeRoles(ctx, revokeRolesReq)
	if err != nil {
		msg := "failed to revoke user roles"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// ListRoles returns a list of all roles.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//
// Returns:
//
//	*protos.ListRolesResponse: A list of all roles.
//	error: An error if the list retrieval fails, otherwise nil.
func (c *Client) ListRoles(ctx context.Context) (*protos.ListRolesResponse, error) {
	c.logger.DebugContext(ctx, "listing roles")

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to list roles"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	rolesResp, err := conn.userAdminClient.ListRoles(ctx, &emptypb.Empty{})
	if err != nil {
		msg := "failed to list roles"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return rolesResp, nil
}

// NodeIDs returns a list of all the node IDs that the client is connected to.
// If load-balancer is set true no NodeIDs will be returned.
// If a node is accessible but not a part of the cluster it will not be
// returned.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//
// Returns:
//
//	[]*protos.NodeId: A list of node IDs.
func (c *Client) NodeIDs(ctx context.Context) []*protos.NodeId {
	c.logger.DebugContext(ctx, "getting cluster info")

	ids := c.connectionProvider.GetNodeIDs()
	nodeIDs := make([]*protos.NodeId, len(ids))

	for i, id := range ids {
		nodeIDs[i] = &protos.NodeId{Id: id}
	}

	c.logger.Debug("got node ids", slog.Any("nodeIDs", nodeIDs))

	return nodeIDs
}

// ConnectedNodeEndpoint returns the endpoint used to connect to a node.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	nodeID (*protos.NodeId): The ID of the node. If nil, the seed (or LB) node is used.
//
// Returns:
//
//	*protos.ServerEndpoint: The server endpoint.
//	error: An error if the retrieval fails, otherwise nil.
func (c *Client) ConnectedNodeEndpoint(
	ctx context.Context,
	nodeID *protos.NodeId,
) (*protos.ServerEndpoint, error) {
	c.logger.DebugContext(ctx, "getting connected endpoint for node", slog.Any("nodeID", nodeID))

	conn, err := c.getConnection(nodeID)
	if err != nil {
		msg := "failed to get connected endpoint"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSError(msg, err)
	}

	splitEndpoint := strings.Split(conn.grpcConn.Target(), ":")

	resp := protos.ServerEndpoint{
		Address: splitEndpoint[0],
	}

	if len(splitEndpoint) > 1 {
		port, err := strconv.ParseUint(splitEndpoint[1], 10, 32)
		if err != nil {
			msg := "failed to parse port"
			c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

			return nil, NewAVSError(msg, err)
		}

		resp.Port = uint32(port)
	}

	return &resp, nil
}

// ClusteringState returns the state of the cluster according to the given node.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	nodeID (*protos.NodeId): The ID of the node. If nil, the seed node is used.
//
// Returns:
//
//	*protos.ClusteringState: The clustering state.
//	error: An error if the retrieval fails, otherwise nil.
func (c *Client) ClusteringState(ctx context.Context, nodeID *protos.NodeId) (*protos.ClusteringState, error) {
	c.logger.DebugContext(ctx, "getting clustering state for node", slog.Any("nodeID", nodeID))

	conn, err := c.getConnection(nodeID)
	if err != nil {
		msg := "failed to get clustering state"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSError(msg, err)
	}

	state, err := conn.clusterInfoClient.GetClusteringState(ctx, &emptypb.Empty{})
	if err != nil {
		msg := "failed to get clustering state"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return state, nil
}

// ClusterEndpoints returns the endpoints of all the nodes in the cluster according to the specified node.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	nodeID (*protos.NodeId): The ID of the node. If nil, the seed node is used.
//	listenerName (*string): The name of the listener. If nil, the default listener is used.
//
// Returns:
//
//	*protos.ClusterNodeEndpoints: The cluster node endpoints.
//	error: An error if the retrieval fails, otherwise nil.
func (c *Client) ClusterEndpoints(
	ctx context.Context,
	nodeID *protos.NodeId,
	listenerName *string,
) (*protos.ClusterNodeEndpoints, error) {
	c.logger.DebugContext(ctx, "getting cluster endpoints for node", slog.Any("nodeID", nodeID))

	conn, err := c.getConnection(nodeID)
	if err != nil {
		msg := "failed to get cluster endpoints"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSError(msg, err)
	}

	endpoints, err := conn.clusterInfoClient.GetClusterEndpoints(ctx,
		&protos.ClusterNodeEndpointsRequest{
			ListenerName: listenerName,
		},
	)
	if err != nil {
		msg := "failed to get cluster endpoints"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return endpoints, nil
}

// About returns information about the provided node.
//
// Args:
//
//	ctx (context.Context): The context for the operation.
//	nodeID (*protos.NodeId): The ID of the node. If nil, the seed node is used.
//
// Returns:
//
//	*protos.AboutResponse: Information about the node.
//	error: An error if the retrieval fails, otherwise nil.
func (c *Client) About(ctx context.Context, nodeID *protos.NodeId) (*protos.AboutResponse, error) {
	c.logger.DebugContext(ctx, "getting \"about\" info from nodes")

	conn, err := c.getConnection(nodeID)
	if err != nil {
		msg := "failed to make about request"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSError(msg, err)
	}

	resp, err := conn.aboutClient.Get(ctx, &protos.AboutRequest{})
	if err != nil {
		msg := "failed to make about request"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return resp, nil
}

func (c *Client) getConnection(nodeID *protos.NodeId) (*connection, error) {
	if nodeID == nil {
		return c.connectionProvider.GetSeedConn()
	}

	return c.connectionProvider.GetNodeConn(nodeID.GetId())
}

// waitForIndexCreation waits for an index to be created and blocks until it is.
// The amount of time to wait between each call is defined by waitInterval.
func (c *Client) waitForIndexCreation(ctx context.Context,
	namespace,
	indexName string,
	waitInterval time.Duration,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to wait for index creation"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	indexStatusReq := createIndexStatusRequest(namespace, indexName)
	timer := time.NewTimer(waitInterval)

	defer timer.Stop()

	for {
		_, err := conn.indexClient.GetStatus(ctx, indexStatusReq)
		if err != nil {
			code := status.Code(err)
			if code == codes.Unavailable || code == codes.NotFound {
				logger.Debug("index does not exist, waiting...")

				timer.Reset(waitInterval)

				select {
				case <-timer.C:
				case <-ctx.Done():
					msg := "waiting for index creation canceled"
					logger.ErrorContext(ctx, msg)

					return NewAVSError(msg, ctx.Err())
				}
			} else {
				msg := "unable to wait for index creation, an unexpected error occurred"

				logger.Error(msg, slog.Any("error", err))

				return NewAVSErrorFromGrpc(msg, err)
			}
		} else {
			logger.Info("index has been created")
			break
		}
	}

	return nil
}

// waitForIndexDrop waits for an index to be dropped and blocks until it is. The
// amount of time to wait between each call is defined by waitInterval.
func (c *Client) waitForIndexDrop(ctx context.Context, namespace, indexName string, waitInterval time.Duration) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))

	conn, err := c.connectionProvider.GetRandomConn()
	if err != nil {
		msg := "failed to wait for index deletion"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	indexStatusReq := createIndexStatusRequest(namespace, indexName)
	timer := time.NewTimer(waitInterval)

	defer timer.Stop()

	for {
		_, err := conn.indexClient.GetStatus(ctx, indexStatusReq)
		if err != nil {
			code := status.Code(err)
			if code == codes.Unavailable || code == codes.NotFound {
				logger.Info("index is deleted")
				return nil
			}

			msg := "unable to wait for index deletion, an unexpected error occurred"
			logger.Error(msg, slog.Any("error", err))

			return NewAVSErrorFromGrpc(msg, err)
		}

		c.logger.Debug("index still exists, waiting...")
		timer.Reset(waitInterval)

		select {
		case <-timer.C:
		case <-ctx.Done():
			logger.ErrorContext(ctx, "waiting for index deletion canceled")
			return ctx.Err()
		}
	}
}
