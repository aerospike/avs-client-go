// Package avs provides a client for managing Aerospike Vector Indexes.
package avs

import (
	"context"
	"crypto/tls"
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

// Client is a client for managing Aerospike Vector Indexes.
type Client struct {
	logger          *slog.Logger
	channelProvider *channelProvider
}

// NewAdminClient creates a new AdminClient instance.
//   - seeds: A list of seed hosts to connect to.
//   - listenerName: The name of the listener to connect to as configured on the
//     server.
//   - isLoadBalancer: Whether the client should consider the seed a load balancer.
//     Only the first seed is considered. Subsequent seeds are ignored.
//   - username: The username to authenticate with.
//   - password: The password to authenticate with.
//   - tlsConfig: The TLS configuration to use for the connection.
//   - logger: The logger to use for logging.
func NewClient(
	ctx context.Context,
	seeds HostPortSlice,
	listenerName *string,
	isLoadBalancer bool,
	credentials *UserPassCredentials,
	tlsConfig *tls.Config,
	logger *slog.Logger,
) (*Client, error) {
	logger = logger.WithGroup("avs.admin")
	logger.Info("creating new client")

	channelProvider, err := newChannelProvider(
		ctx,
		seeds,
		listenerName,
		isLoadBalancer,
		credentials,
		tlsConfig,
		logger,
	)
	if err != nil {
		logger.Error("failed to create channel provider", slog.Any("error", err))
		return nil, NewAVSErrorFromGrpc("failed to connect to server", err)
	}

	return &Client{
		logger:          logger,
		channelProvider: channelProvider,
	}, nil
}

// Close closes the AdminClient and releases any resources associated with it.
func (c *Client) Close() error {
	c.logger.Info("Closing client")
	return c.channelProvider.Close()
}

func (c *Client) put(ctx context.Context, writeType protos.WriteType, namespace string, set *string, key any, recordData map[string]any, ignoreMemQueueFull bool) error {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.Any("key", key),
	)

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to insert record"
		logger.Error(msg, slog.Any("error", err))
		return NewAVSError(msg, err)
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		msg := "failed to insert record"
		logger.Error(msg, slog.Any("error", err))
		return NewAVSError(msg, err)
	}

	fields, err := protos.ConvertToFields(recordData)
	if err != nil {
		msg := "failed to insert record"
		logger.Error(msg, slog.Any("error", err))
		return NewAVSError(msg, err)
	}

	putReq := &protos.PutRequest{
		Key:                protoKey,
		WriteType:          writeType,
		Fields:             fields,
		IgnoreMemQueueFull: ignoreMemQueueFull,
	}

	_, err = conn.transactClient.Put(ctx, putReq)
	if err != nil {
		msg := "failed to insert record"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))
		return NewAVSErrorFromGrpc(msg, err)
	}

	return err
}

func (c *Client) Insert(ctx context.Context, namespace string, set *string, key any, recordData map[string]any, ignoreMemQueueFull bool) error {
	c.logger.InfoContext(ctx, "inserting record", slog.String("namespace", namespace), slog.Any("key", key))
	return c.put(ctx, protos.WriteType_INSERT_ONLY, namespace, set, key, recordData, ignoreMemQueueFull)
}

func (c *Client) Update(ctx context.Context, namespace string, set *string, key any, recordData map[string]any, ignoreMemQueueFull bool) error {
	c.logger.InfoContext(ctx, "updating record", slog.String("namespace", namespace), slog.Any("key", key))
	return c.put(ctx, protos.WriteType_UPDATE_ONLY, namespace, set, key, recordData, ignoreMemQueueFull)
}

func (c *Client) Upsert(ctx context.Context, namespace string, set *string, key any, recordData map[string]any, ignoreMemQueueFull bool) error {
	c.logger.InfoContext(ctx, "upserting record", slog.String("namespace", namespace), slog.Any("set", set), slog.Any("key", key))
	return c.put(ctx, protos.WriteType_UPSERT, namespace, set, key, recordData, ignoreMemQueueFull)
}

//nolint:revive // TODO
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
	logger.InfoContext(ctx, "getting record")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get record"
		logger.Error(msg, slog.Any("error", err))
		return nil, NewAVSError(msg, err)
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		msg := "failed to insert record"
		logger.Error(msg, slog.Any("error", err))
		return nil, NewAVSError(msg, err)
	}

	getReq := &protos.GetRequest{
		Key:        protoKey,
		Projection: createProjectionSpec(includeFields, excludeFields),
	}

	record, err := conn.transactClient.Get(ctx, getReq)
	if err != nil {
		msg := "failed to get record"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))
		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	logger.Debug("received record", slog.Any("record", record.Fields))

	return newRecordFromProto(record), nil
}

//nolint:revive // TODO
func (c *Client) Delete(ctx context.Context, namespace string, set *string, key any) error {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.Any("key", key),
	)
	logger.InfoContext(ctx, "deleting record")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		logger.Error("failed to delete record", slog.Any("error", err))
		return err
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		msg := "failed to insert record"
		logger.Error(msg, slog.Any("error", err))
		return NewAVSError(msg, err)
	}

	getReq := &protos.DeleteRequest{
		Key: protoKey,
	}

	_, err = conn.transactClient.Delete(ctx, getReq)

	return err
}

//nolint:revive // TODO
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
	logger.InfoContext(ctx, "checking if record exists")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		logger.Error("failed to check if record exists", slog.Any("error", err))
		return false, err
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		msg := "failed to insert record"
		logger.Error(msg, slog.Any("error", err))
		return false, NewAVSError(msg, err)
	}

	existsReq := &protos.ExistsRequest{
		Key: protoKey,
	}

	boolean, err := conn.transactClient.Exists(ctx, existsReq)

	return boolean.GetValue(), err
}

//nolint:revive // TODO
func (c *Client) IsIndexed(ctx context.Context, namespace string, set *string, indexName string, key any) (bool, error) {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.String("indexName", indexName),
		slog.Any("key", key),
	)
	logger.InfoContext(ctx, "checking if record is indexed")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		logger.Error("failed to check if record is indexed", slog.Any("error", err))
		return false, err
	}

	protoKey, err := protos.ConvertToKey(namespace, set, key)
	if err != nil {
		msg := "failed to insert record"
		logger.Error(msg, slog.Any("error", err))
		return false, NewAVSError(msg, err)
	}

	isIndexedReq := &protos.IsIndexedRequest{
		Key: protoKey,
	}

	boolean, err := conn.transactClient.IsIndexed(ctx, isIndexedReq)

	return boolean.GetValue(), err
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
	logger.InfoContext(ctx, "searching for vector")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		logger.Error("failed to search for vector", slog.Any("error", err))
		return nil, err
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
			if err != io.EOF {
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

// VectorSearchFloat32 searches for the nearest neighbors of the query vector in the
// specified index. It takes the following parameters:
//   - namespace: The namespace of the index.
//   - indexName: The name of the index.
//   - query: The query vector.
//   - limit: The maximum number of neighbors to return.
//   - searchParams: Extra options sent to the server to configure behavior of the
//     HNSW algorithm.
//   - projections: Extra options to configure the behavior of the which keys
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
	c.logger.InfoContext(ctx, "searching for float vector")

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

// VectorSearchBool searches for the nearest neighbors of the query vector in the
// specified index. It takes the following parameters:
//   - namespace: The namespace of the index.
//   - indexName: The name of the index.
//   - query: The query vector.
//   - limit: The maximum number of neighbors to return.
//   - searchParams: Extra options sent to the server to configure behavior of the
//     HNSW algorithm.
//   - projections: Extra options to configure the behavior of the which keys
func (c *Client) VectorSearchBool(ctx context.Context,
	namespace,
	indexName string,
	query []bool,
	limit uint32,
	searchParams *protos.HnswSearchParams,
	includeFields,
	excludeFields []string,
) ([]*Neighbor, error) {
	c.logger.InfoContext(ctx, "searching for bool vector")

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

//nolint:revive // TODO
func (c *Client) WaitForIndexCompletion(ctx context.Context, namespace, indexName string, waitInterval time.Duration) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("indexName", indexName))
	logger.InfoContext(ctx, "waiting for index completion")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		logger.Error("failed to wait for index completion", slog.Any("error", err))
		return err
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      indexName,
	}

	timer := time.NewTimer(waitInterval)
	startTime := time.Now()
	unmergedZeroCount := 0
	unmergedNotZeroCount := 0

	defer timer.Stop()

	for {
		status, err := conn.indexClient.GetStatus(ctx, indexID)
		if err != nil {
			logger.ErrorContext(ctx, "failed to wait for index completion", slog.Any("error", err))
			return err
		}

		// We consider the index completed when unmerged record count == 0 for
		// 2 cycles OR unmerged record count != 0 and now it == 0.
		unmerged := status.GetUnmergedRecordCount()
		if unmerged == 0 {
			if unmergedZeroCount >= 2 || unmergedNotZeroCount >= 1 {
				logger.InfoContext(ctx, "index completed", slog.Duration("duration", time.Since(startTime)))
				return nil
			}

			unmergedZeroCount += 1
			logger.DebugContext(ctx, "index not yet completed", slog.Int64("unmerged", unmerged))
		} else {
			logger.DebugContext(ctx, "index not yet completed", slog.Int64("unmerged", unmerged))
			unmergedNotZeroCount += 1
		}

		timer.Reset(waitInterval)

		select {
		case <-timer.C:
		case <-ctx.Done():
			logger.ErrorContext(ctx, "waiting for index completion canceled")
			return ctx.Err()
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
	Labels     map[string]string
	Sets       []string
}

// IndexCreate creates a new Aerospike Vector Index and blocks until it is created.
// It takes the following parameters:
//   - namespace: The namespace of the index.
//   - name: The name of the index.
//   - vectorField: The field to create the index on.
//   - dimensions: The number of dimensions in the vector.
//   - vectorDistanceMetric: The distance metric to use for the index.
//   - opts: Optional fields to configure the index
//
// It returns an error if the index creation fails.
func (c *Client) IndexCreate(
	ctx context.Context,
	namespace string,
	name string,
	vectorField string,
	dimensions uint32,
	vectorDistanceMetric protos.VectorDistanceMetric,
	opts *IndexCreateOpts,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))
	logger.InfoContext(ctx, "creating index")

	var (
		set     *string
		params  *protos.IndexDefinition_HnswParams
		labels  map[string]string
		storage *protos.IndexStorage
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
	}

	indexDef := &protos.IndexDefinition{
		Id: &protos.IndexId{
			Namespace: namespace,
			Name:      name,
		},
		Dimensions:           dimensions,
		VectorDistanceMetric: vectorDistanceMetric,
		Field:                vectorField,
		SetFilter:            set,
		Params:               params,
		Labels:               labels,
		Storage:              storage,
	}

	return c.IndexCreateFromIndexDef(ctx, indexDef)
}

// IndexCreateFromIndexDef creates a new Aerospike Vector Index from a provided
// IndexDefinition and blocks until it is created. It can be easily used in
// conjunction with IndexList and IndexGet to create a new index using the
// returned IndexDefinitions.
func (c *Client) IndexCreateFromIndexDef(
	ctx context.Context,
	indexDef *protos.IndexDefinition,
) error {
	logger := c.logger.With(slog.Any("definition", indexDef))
	logger.InfoContext(ctx, "creating index from definition")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to create index from definition"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSError(msg, err)
	}

	_, err = conn.indexClient.Create(ctx, indexDef)
	if err != nil {
		msg := "failed to create index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	ctx, cancel := context.WithTimeout(ctx, indexTimeoutDuration)
	defer cancel()

	return c.waitForIndexCreation(ctx, indexDef.Id.Namespace, indexDef.Id.Name, indexWaitDuration)
}

// IndexUpdate updates an existing Aerospike Vector Index's dynamic
// configuration params
func (c *Client) IndexUpdate(
	ctx context.Context,
	namespace string,
	name string,
	metadata map[string]string,
	hnswParams *protos.HnswIndexUpdate,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	logger.InfoContext(ctx, "updating index")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to update index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	indexUpdate := &protos.IndexUpdateRequest{
		IndexId: &protos.IndexId{
			Namespace: namespace,
			Name:      name,
		},
		Labels: metadata,
		Update: &protos.IndexUpdateRequest_HnswIndexUpdate{
			HnswIndexUpdate: hnswParams,
		},
	}

	_, err = conn.indexClient.Update(ctx, indexUpdate)
	if err != nil {
		msg := "failed to update index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// IndexDrop drops an existing Aerospike Vector Index and blocks until it is.
func (c *Client) IndexDrop(ctx context.Context, namespace, name string) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))
	logger.InfoContext(ctx, "dropping index")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to drop index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}

	_, err = conn.indexClient.Drop(ctx, indexID)
	if err != nil {
		msg := "failed to drop index"

		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	ctx, cancel := context.WithTimeout(ctx, indexTimeoutDuration)
	defer cancel()

	return c.waitForIndexDrop(ctx, namespace, name, indexWaitDuration)
}

// IndexList returns a list of all Aerospike Vector Indexes. To get a single
// index use IndexGet.
func (c *Client) IndexList(ctx context.Context) (*protos.IndexDefinitionList, error) {
	c.logger.InfoContext(ctx, "listing indexes")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get indexes"

		c.logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	indexList, err := conn.indexClient.List(ctx, nil)
	if err != nil {
		msg := "failed to get indexes"

		c.logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexList, nil
}

// IndexGet returns the definition of an Aerospike Vector Index. To get all
// indexes use IndexList.
func (c *Client) IndexGet(ctx context.Context, namespace, name string) (*protos.IndexDefinition, error) {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))
	logger.InfoContext(ctx, "getting index")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get index"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}

	indexDef, err := conn.indexClient.Get(ctx, indexID)
	if err != nil {
		msg := "failed to get index"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexDef, nil
}

// IndexGetStatus returns the status of an Aerospike Vector Index.
func (c *Client) IndexGetStatus(ctx context.Context, namespace, name string) (*protos.IndexStatusResponse, error) {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))
	logger.InfoContext(ctx, "getting index status")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get index status"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}

	indexStatus, err := conn.indexClient.GetStatus(ctx, indexID)
	if err != nil {
		msg := "failed to get index status"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexStatus, nil
}

// GcInvalidVertices garbage collects invalid vertices in an Aerospike Vector Index.
func (c *Client) GcInvalidVertices(ctx context.Context, namespace, name string, cutoffTime time.Time) error {
	logger := c.logger.With(
		slog.String("namespace", namespace),
		slog.String("name", name),
		slog.Any("cutoffTime", cutoffTime),
	)

	logger.InfoContext(ctx, "garbage collection invalid vertices")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to garbage collect invalid vertices"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	gcRequest := &protos.GcInvalidVerticesRequest{
		IndexId: &protos.IndexId{
			Namespace: namespace,
			Name:      name,
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
func (c *Client) CreateUser(ctx context.Context, username, password string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))
	logger.InfoContext(ctx, "creating user")

	conn, err := c.channelProvider.GetRandomConn()
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
func (c *Client) UpdateCredentials(ctx context.Context, username, password string) error {
	logger := c.logger.With(slog.String("username", username))
	logger.InfoContext(ctx, "updating user credentials")

	conn, err := c.channelProvider.GetRandomConn()
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
func (c *Client) DropUser(ctx context.Context, username string) error {
	logger := c.logger.With(slog.String("username", username))
	logger.InfoContext(ctx, "dropping user")

	conn, err := c.channelProvider.GetRandomConn()
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
func (c *Client) GetUser(ctx context.Context, username string) (*protos.User, error) {
	logger := c.logger.With(slog.String("username", username))
	logger.InfoContext(ctx, "getting user")

	conn, err := c.channelProvider.GetRandomConn()
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
func (c *Client) ListUsers(ctx context.Context) (*protos.ListUsersResponse, error) {
	c.logger.InfoContext(ctx, "listing users")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to list users"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	usersResp, err := conn.userAdminClient.ListUsers(ctx, &emptypb.Empty{})
	if err != nil {
		msg := "failed to lists users"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return usersResp, nil
}

// GrantRoles grants the provided roles to the user with the provided username.
func (c *Client) GrantRoles(ctx context.Context, username string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))
	logger.InfoContext(ctx, "granting user roles")

	conn, err := c.channelProvider.GetRandomConn()
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
func (c *Client) RevokeRoles(ctx context.Context, username string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))
	logger.InfoContext(ctx, "revoking user roles")

	conn, err := c.channelProvider.GetRandomConn()
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
func (c *Client) ListRoles(ctx context.Context) (*protos.ListRolesResponse, error) {
	c.logger.InfoContext(ctx, "listing roles")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to list roles"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	rolesResp, err := conn.userAdminClient.ListRoles(ctx, &emptypb.Empty{})
	if err != nil {
		msg := "failed to lists roles"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return rolesResp, nil
}

// NodeIds returns a list of all the node ids that the client is connected to.
// If a node is accessible but not a part of the cluster it will not be returned.
func (c *Client) NodeIDs(ctx context.Context) []*protos.NodeId {
	c.logger.InfoContext(ctx, "getting cluster info")

	ids := c.channelProvider.GetNodeIDs()
	nodeIDs := make([]*protos.NodeId, len(ids))

	for i, id := range ids {
		nodeIDs[i] = &protos.NodeId{Id: id}
	}

	c.logger.Debug("got node ids", slog.Any("nodeIDs", nodeIDs))

	return nodeIDs
}

// ConnectedNodeEndpoint returns the endpoint used to connect to a node. If
// nodeID is nil then an endpoint used to connect to your seed (or
// load-balancer) is used.
func (c *Client) ConnectedNodeEndpoint(
	ctx context.Context,
	nodeID *protos.NodeId,
) (*protos.ServerEndpoint, error) {
	c.logger.InfoContext(ctx, "getting connected endpoint for node", slog.Any("nodeID", nodeID))

	var (
		conn *connection
		err  error
	)

	if nodeID == nil {
		conn, err = c.channelProvider.GetSeedConn()
	} else {
		conn, err = c.channelProvider.GetNodeConn(nodeID.Id)
	}

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

// ClusteringState returns the state of the cluster according the
// given node.  If nodeID is nil then the seed node is used.
func (c *Client) ClusteringState(ctx context.Context, nodeID *protos.NodeId) (*protos.ClusteringState, error) {
	c.logger.InfoContext(ctx, "getting clustering state for node", slog.Any("nodeID", nodeID))

	var (
		conn *connection
		err  error
	)

	if nodeID == nil {
		conn, err = c.channelProvider.GetSeedConn()
	} else {
		conn, err = c.channelProvider.GetNodeConn(nodeID.GetId())
	}

	if err != nil {
		msg := "failed to list roles"
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

// ClusterEndpoints returns the endpoints of all the nodes in the cluster
// according to the specified node. If nodeID is nil then the seed node is used.
// If listenerName is nil then the default listener name is used.
func (c *Client) ClusterEndpoints(
	ctx context.Context,
	nodeID *protos.NodeId,
	listenerName *string,
) (*protos.ClusterNodeEndpoints, error) {
	c.logger.InfoContext(ctx, "getting cluster endpoints for node", slog.Any("nodeID", nodeID))

	var (
		conn *connection
		err  error
	)

	if nodeID == nil {
		conn, err = c.channelProvider.GetSeedConn()
	} else {
		conn, err = c.channelProvider.GetNodeConn(nodeID.GetId())
	}

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

// About returns information about the provided node. If nodeID is nil
// then the seed node is used.
func (c *Client) About(ctx context.Context, nodeID *protos.NodeId) (*protos.AboutResponse, error) {
	c.logger.InfoContext(ctx, "getting \"about\" info from nodes")

	var (
		conn *connection
		err  error
	)

	if nodeID == nil {
		conn, err = c.channelProvider.GetSeedConn()
	} else {
		conn, err = c.channelProvider.GetNodeConn(nodeID.GetId())
	}

	if err != nil {
		msg := "failed to make about request"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	resp, err := conn.aboutClient.Get(ctx, &protos.AboutRequest{})
	if err != nil {
		msg := "failed to make about request"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return resp, nil
}

// waitForIndexCreation waits for an index to be created and blocks until it is.
// The amount of time to wait between each call is defined by waitInterval.
func (c *Client) waitForIndexCreation(ctx context.Context,
	namespace,
	name string,
	waitInterval time.Duration,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to wait for index creation"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}

	timer := time.NewTimer(waitInterval)

	defer timer.Stop()

	for {
		_, err := conn.indexClient.GetStatus(ctx, indexID)
		if err != nil {
			code := status.Code(err)
			if code == codes.Unavailable || code == codes.NotFound {
				logger.Debug("index does not exist, waiting...")

				timer.Reset(waitInterval)

				select {
				case <-timer.C:
				case <-ctx.Done():
					logger.ErrorContext(ctx, "waiting for index creation canceled")
					return ctx.Err()
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
func (c *Client) waitForIndexDrop(ctx context.Context, namespace, name string, waitInterval time.Duration) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to wait for index deletion"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}

	timer := time.NewTimer(waitInterval)

	defer timer.Stop()

	for {
		_, err := conn.indexClient.GetStatus(ctx, indexID)
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
