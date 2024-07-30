// Package avs provides a client for managing Aerospike Vector Indexes.
package avs

import (
	"context"
	"crypto/tls"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/aerospike/avs-client-go/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	indexTimeoutDuration = time.Second * 100
	indexWaitDuration    = time.Millisecond * 100
)

// AdminClient is a client for managing Aerospike Vector Indexes.
type AdminClient struct {
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
func NewAdminClient(
	ctx context.Context,
	seeds HostPortSlice,
	listenerName *string,
	isLoadBalancer bool,
	credentials *UserPassCredentials,
	tlsConfig *tls.Config,
	logger *slog.Logger,
) (*AdminClient, error) {
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

	return &AdminClient{
		logger:          logger,
		channelProvider: channelProvider,
	}, nil
}

// Close closes the AdminClient and releases any resources associated with it.
func (c *AdminClient) Close() {
	c.logger.Info("Closing client")
	c.channelProvider.Close()
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
func (c *AdminClient) IndexCreate(
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
func (c *AdminClient) IndexCreateFromIndexDef(
	ctx context.Context,
	indexDef *protos.IndexDefinition,
) error {
	logger := c.logger.With(slog.Any("definition", indexDef))
	logger.InfoContext(ctx, "creating index from definition")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to create index from definition"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewIndexServiceClient(conn)

	_, err = client.Create(ctx, indexDef)
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
func (c *AdminClient) IndexUpdate(
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

	client := protos.NewIndexServiceClient(conn)

	_, err = client.Update(ctx, indexUpdate)
	if err != nil {
		msg := "failed to update index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// IndexDrop drops an existing Aerospike Vector Index and blocks until it is.
func (c *AdminClient) IndexDrop(ctx context.Context, namespace, name string) error {
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

	client := protos.NewIndexServiceClient(conn)

	_, err = client.Drop(ctx, indexID)
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
func (c *AdminClient) IndexList(ctx context.Context) (*protos.IndexDefinitionList, error) {
	c.logger.InfoContext(ctx, "listing indexes")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get indexes"

		c.logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewIndexServiceClient(conn)

	indexList, err := client.List(ctx, nil)
	if err != nil {
		msg := "failed to get indexes"

		c.logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexList, nil
}

// IndexGet returns the definition of an Aerospike Vector Index. To get all
// indexes use IndexList.
func (c *AdminClient) IndexGet(ctx context.Context, namespace, name string) (*protos.IndexDefinition, error) {
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
	client := protos.NewIndexServiceClient(conn)

	indexDef, err := client.Get(ctx, indexID)
	if err != nil {
		msg := "failed to get index"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexDef, nil
}

// IndexGetStatus returns the status of an Aerospike Vector Index.
func (c *AdminClient) IndexGetStatus(ctx context.Context, namespace, name string) (*protos.IndexStatusResponse, error) {
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
	client := protos.NewIndexServiceClient(conn)

	indexStatus, err := client.GetStatus(ctx, indexID)
	if err != nil {
		msg := "failed to get index status"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexStatus, nil
}

// GcInvalidVertices garbage collects invalid vertices in an Aerospike Vector Index.
func (c *AdminClient) GcInvalidVertices(ctx context.Context, namespace, name string, cutoffTime time.Time) error {
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
	client := protos.NewIndexServiceClient(conn)

	_, err = client.GcInvalidVertices(ctx, gcRequest)
	if err != nil {
		msg := "failed to garbage collect invalid vertices"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// CreateUser creates a new user with the provided username, password, and roles.
func (c *AdminClient) CreateUser(ctx context.Context, username, password string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))
	logger.InfoContext(ctx, "creating user")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to create user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewUserAdminServiceClient(conn)

	addUserRequest := &protos.AddUserRequest{
		Credentials: createUserPassCredential(username, password),
		Roles:       roles,
	}

	_, err = client.AddUser(ctx, addUserRequest)
	if err != nil {
		msg := "failed to create user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// UpdateCredentials updates the password for the provided username.
func (c *AdminClient) UpdateCredentials(ctx context.Context, username, password string) error {
	logger := c.logger.With(slog.String("username", username))
	logger.InfoContext(ctx, "updating user credentials")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to update user credentials"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewUserAdminServiceClient(conn)

	updatedCredRequest := &protos.UpdateCredentialsRequest{
		Credentials: createUserPassCredential(username, password),
	}

	_, err = client.UpdateCredentials(ctx, updatedCredRequest)
	if err != nil {
		msg := "failed to update user credentials"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// DropUser deletes the user with the provided username.
func (c *AdminClient) DropUser(ctx context.Context, username string) error {
	logger := c.logger.With(slog.String("username", username))
	logger.InfoContext(ctx, "dropping user")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to drop user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewUserAdminServiceClient(conn)

	dropUserRequest := &protos.DropUserRequest{
		Username: username,
	}

	_, err = client.DropUser(ctx, dropUserRequest)
	if err != nil {
		msg := "failed to drop user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// GetUser returns the user with the provided username.
func (c *AdminClient) GetUser(ctx context.Context, username string) (*protos.User, error) {
	logger := c.logger.With(slog.String("username", username))
	logger.InfoContext(ctx, "getting user")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to get user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewUserAdminServiceClient(conn)

	getUserRequest := &protos.GetUserRequest{
		Username: username,
	}

	userResp, err := client.GetUser(ctx, getUserRequest)
	if err != nil {
		msg := "failed to get user"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return userResp, nil
}

// ListUsers returns a list of all users.
func (c *AdminClient) ListUsers(ctx context.Context) (*protos.ListUsersResponse, error) {
	c.logger.InfoContext(ctx, "listing users")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to list users"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewUserAdminServiceClient(conn)

	usersResp, err := client.ListUsers(ctx, &emptypb.Empty{})
	if err != nil {
		msg := "failed to lists users"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return usersResp, nil
}

// GrantRoles grants the provided roles to the user with the provided username.
func (c *AdminClient) GrantRoles(ctx context.Context, username string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))
	logger.InfoContext(ctx, "granting user roles")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to grant user roles"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewUserAdminServiceClient(conn)

	grantRolesRequest := &protos.GrantRolesRequest{
		Username: username,
		Roles:    roles,
	}

	_, err = client.GrantRoles(ctx, grantRolesRequest)
	if err != nil {
		msg := "failed to grant user roles"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// RevokeRoles revokes the provided roles from the user with the provided username.
func (c *AdminClient) RevokeRoles(ctx context.Context, username string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))
	logger.InfoContext(ctx, "revoking user roles")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to revoke user roles"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewUserAdminServiceClient(conn)

	revokeRolesReq := &protos.RevokeRolesRequest{
		Username: username,
		Roles:    roles,
	}

	_, err = client.RevokeRoles(ctx, revokeRolesReq)
	if err != nil {
		msg := "failed to revoke user roles"
		logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	return nil
}

// ListRoles returns a list of all roles.
func (c *AdminClient) ListRoles(ctx context.Context) (*protos.ListRolesResponse, error) {
	c.logger.InfoContext(ctx, "listing roles")

	conn, err := c.channelProvider.GetRandomConn()
	if err != nil {
		msg := "failed to list roles"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewUserAdminServiceClient(conn)

	rolesResp, err := client.ListRoles(ctx, &emptypb.Empty{})
	if err != nil {
		msg := "failed to lists roles"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return rolesResp, nil
}

// NodeIds returns a list of all the node ids that the client is connected to.
// If a node is accessible but not apart of the cluster it will not be returned.
func (c *AdminClient) NodeIds(ctx context.Context) []*protos.NodeId {
	c.logger.InfoContext(ctx, "getting cluster info")

	ids := c.channelProvider.GetNodeIds()
	nodeIds := make([]*protos.NodeId, len(ids))

	for i, id := range ids {
		nodeIds[i] = &protos.NodeId{Id: id}
	}

	c.logger.Debug("got node ids", slog.Any("nodeIds", nodeIds))

	return nodeIds
}

// ConnectedNodeEndpoint returns the endpoint used to connect to a node. If
// nodeId is nil then an endpoint used to connect to your seed (or
// load-balancer) is used.
func (c *AdminClient) ConnectedNodeEndpoint(ctx context.Context, nodeId *protos.NodeId) (*protos.ServerEndpoint, error) {
	c.logger.InfoContext(ctx, "getting connected endpoint for node", slog.Any("nodeId", nodeId))

	var (
		conn *grpc.ClientConn
		err  error
	)

	if nodeId == nil {
		conn, err = c.channelProvider.GetSeedConn()
	} else {
		conn, err = c.channelProvider.GetNodeConn(nodeId.Id)
	}

	if err != nil {
		msg := "failed to get connected endpoint"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSError(msg)
	}

	splitEndpoint := strings.Split(conn.Target(), ":")

	resp := protos.ServerEndpoint{
		Address: splitEndpoint[0],
	}

	if len(splitEndpoint) > 1 {
		port, err := strconv.ParseUint(splitEndpoint[1], 10, 32)
		if err != nil {
			msg := "failed to parse port"
			c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

			return nil, NewAVSErrorFromGrpc(msg, err)
		}

		resp.Port = uint32(port)
	}

	return &resp, nil
}

// ClusteringState returns the state of the cluster according the
// given node.  If nodeId is nil then the seed node is used.
func (c *AdminClient) ClusteringState(ctx context.Context, nodeId *protos.NodeId) (*protos.ClusteringState, error) {
	c.logger.InfoContext(ctx, "getting clustering state for node", slog.Any("nodeId", nodeId))

	var (
		conn *grpc.ClientConn
		err  error
	)

	if nodeId == nil {
		conn, err = c.channelProvider.GetSeedConn()
	} else {
		conn, err = c.channelProvider.GetNodeConn(nodeId.GetId())
	}

	if err != nil {
		msg := "failed to list roles"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewClusterInfoServiceClient(conn)

	state, err := client.GetClusteringState(ctx, &emptypb.Empty{})
	if err != nil {
		msg := "failed to get clustering state"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return state, nil
}

// ClusterEndpoints returns the endpoints of all the nodes in the cluster
// according to the specified node. If nodeId is nil then the seed node is used.
// If listenerName is nil then the default listener name is used.
func (c *AdminClient) ClusterEndpoints(ctx context.Context, nodeId *protos.NodeId, listenerName *string) (*protos.ClusterNodeEndpoints, error) {
	c.logger.InfoContext(ctx, "getting cluster endpoints for node", slog.Any("nodeId", nodeId))

	var (
		conn *grpc.ClientConn
		err  error
	)

	if nodeId == nil {
		conn, err = c.channelProvider.GetSeedConn()
	} else {
		conn, err = c.channelProvider.GetNodeConn(nodeId.GetId())
	}

	if err != nil {
		msg := "failed to get cluster endpoints"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewClusterInfoServiceClient(conn)

	endpoints, err := client.GetClusterEndpoints(ctx,
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

// About returns information about the provided node. If nodeId is nil
// then the seed node is used.
func (c *AdminClient) About(ctx context.Context, nodeId *protos.NodeId) (*protos.AboutResponse, error) {
	c.logger.InfoContext(ctx, "getting \"about\" info from nodes")

	var (
		conn *grpc.ClientConn
		err  error
	)

	if nodeId == nil {
		conn, err = c.channelProvider.GetSeedConn()
	} else {
		conn, err = c.channelProvider.GetNodeConn(nodeId.GetId())
	}

	if err != nil {
		msg := "failed to make about request"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewAboutServiceClient(conn)

	resp, err := client.Get(ctx, &protos.AboutRequest{})
	if err != nil {
		msg := "failed to make about request"
		c.logger.ErrorContext(ctx, msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return resp, nil
}

// waitForIndexCreation waits for an index to be created and blocks until it is.
// The amount of time to wait between each call is defined by waitInterval.
func (c *AdminClient) waitForIndexCreation(ctx context.Context,
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

	client := protos.NewIndexServiceClient(conn)
	timer := time.NewTimer(waitInterval)

	defer timer.Stop()

	defer timer.Stop()

	for {
		_, err := client.GetStatus(ctx, indexID)
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
func (c *AdminClient) waitForIndexDrop(ctx context.Context, namespace, name string, waitInterval time.Duration) error {
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

	client := protos.NewIndexServiceClient(conn)
	timer := time.NewTimer(waitInterval)

	defer timer.Stop()

	defer timer.Stop()

	for {
		_, err := client.GetStatus(ctx, indexID)
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
