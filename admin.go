// Package avs provides a client for managing Aerospike Vector Indexes.
package avs

import (
	"context"
	"crypto/tls"
	"log/slog"
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
	username *string,
	password *string,
	tlsConfig *tls.Config,
	logger *slog.Logger,
) (*AdminClient, error) {
	logger = logger.WithGroup("avs.admin")
	logger.Debug("creating new client")

	channelProvider, err := newChannelProvider(
		ctx,
		seeds,
		listenerName,
		isLoadBalancer,
		username,
		password,
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

// IndexCreate creates a new Aerospike Vector Index and blocks until it is created.
// It takes the following parameters:
//   - namespace: The namespace of the index.
//   - sets: The sets to create the index on. Currently, only one set is supported.
//   - name: The name of the index.
//   - vectorField: The field to create the index on.
//   - dimensions: The number of dimensions in the vector.
//   - vectorDistanceMetric: The distance metric to use for the index.
//   - indexParams: Extra options sent to the server to configure behavior of the
//     HNSW algorithm.
//   - indexMetaData: Extra metadata that can be attached to the index.
//   - indexStorage: The storage configuration for the index. This allows you to
//     configure your index and data to be stored in separate namespaces and/or sets.
//
// It returns an error if the index creation fails.
func (c *AdminClient) IndexCreate(
	ctx context.Context,
	namespace string,
	sets []string,
	name string,
	vectorField string,
	dimensions uint32,
	vectorDistanceMetric protos.VectorDistanceMetric,
	indexParams *protos.HnswParams,
	indexMetaData map[string]string,
	indexStorage *protos.IndexStorage,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetConn()
	if err != nil {
		msg := "failed to create index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	var set *string

	if len(sets) > 0 {
		set = &sets[0]

		if len(sets) > 1 {
			logger.Warn(
				"multiple sets not yet supported for index creation, only the first set will be used",
				slog.String("set", *set),
			)
		}
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
		Params:               &protos.IndexDefinition_HnswParams{HnswParams: indexParams},
		Labels:               indexMetaData,
		Storage:              indexStorage,
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

	return c.waitForIndexCreation(ctx, namespace, name, indexWaitDuration)
}

// IndexDrop drops an existing Aerospike Vector Index and blocks until it is.
func (c *AdminClient) IndexDrop(ctx context.Context, namespace, name string) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetConn()
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
	conn, err := c.channelProvider.GetConn()
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

	conn, err := c.channelProvider.GetConn()
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

	conn, err := c.channelProvider.GetConn()
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

// def CreateUser(self, *, username: str, password: str, roles: list[str]) ->
// int:
func (c *AdminClient) CreateUser(ctx context.Context, username, password string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))

	conn, err := c.channelProvider.GetConn()
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

// def UpdateCredentials(self, *, username: str, password: str) -> int:
func (c *AdminClient) UpdateCredentials(ctx context.Context, username, password string) error {
	logger := c.logger.With(slog.String("username", username))

	conn, err := c.channelProvider.GetConn()
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

// def DropUser(self, *, username: str) -> int:
func (c *AdminClient) DropUser(ctx context.Context, username string) error {
	logger := c.logger.With(slog.String("username", username))

	conn, err := c.channelProvider.GetConn()
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

// def GetUser(self, *, username: str) -> int:
func (c *AdminClient) GetUser(ctx context.Context, username string) (*protos.User, error) {
	logger := c.logger.With(slog.String("username", username))

	conn, err := c.channelProvider.GetConn()
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

// def ListUsers(self) -> int:
func (c *AdminClient) ListUsers(ctx context.Context) (*protos.ListUsersResponse, error) {
	conn, err := c.channelProvider.GetConn()
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

// def GrantRoles(self, *, username: str, roles: list[str]) -> int:
func (c *AdminClient) GrantRoles(ctx context.Context, username string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))

	conn, err := c.channelProvider.GetConn()
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

// def RevokeRoles(self, *, username: str, roles: list[str]) -> int:
func (c *AdminClient) RevokeRoles(ctx context.Context, username string, roles []string) error {
	logger := c.logger.With(slog.String("username", username), slog.Any("roles", roles))

	conn, err := c.channelProvider.GetConn()
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

// def ListRoles(self) -> int:
func (c *AdminClient) ListRoles(ctx context.Context) (*protos.ListRolesResponse, error) {
	conn, err := c.channelProvider.GetConn()
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

// waitForIndexCreation waits for an index to be created and blocks until it is.
// The amount of time to wait between each call is defined by waitInterval.
func (c *AdminClient) waitForIndexCreation(ctx context.Context,
	namespace,
	name string,
	waitInterval time.Duration,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetConn()
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

	conn, err := c.channelProvider.GetConn()
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
