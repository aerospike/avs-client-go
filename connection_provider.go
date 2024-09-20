// Package avs provides a connection provider for connecting to Aerospike servers.
package avs

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aerospike/avs-client-go/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	grpcCreds "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var errConnectionProviderClosed = errors.New("connectionProvider is closed, cannot get connection")

type grpcClientConn interface {
	grpc.ClientConnInterface
	Target() string
	Close() error
}

type tokenManager interface {
	RequireTransportSecurity() bool
	ScheduleRefresh(func() (*connection, error))
	RefreshToken(context.Context, grpcClientConn) error
	UnaryInterceptor() grpc.UnaryClientInterceptor
	StreamInterceptor() grpc.StreamClientInterceptor
	Close()
}

// connection represents a gRPC client connection and all the clients (stubs)
// for the various AVS services. It's main purpose to remove the need to create
// multiple clients for the same connection. This follows the documented grpc
// best practice of reusing connections.
type connection struct {
	grpcConn          grpcClientConn
	transactClient    protos.TransactServiceClient
	authClient        protos.AuthServiceClient
	userAdminClient   protos.UserAdminServiceClient
	indexClient       protos.IndexServiceClient
	aboutClient       protos.AboutServiceClient
	clusterInfoClient protos.ClusterInfoServiceClient
}

// newConnection creates a new connection instance.
func newConnection(conn grpcClientConn) *connection {
	return &connection{
		grpcConn:          conn,
		transactClient:    protos.NewTransactServiceClient(conn),
		authClient:        protos.NewAuthServiceClient(conn),
		userAdminClient:   protos.NewUserAdminServiceClient(conn),
		indexClient:       protos.NewIndexServiceClient(conn),
		aboutClient:       protos.NewAboutServiceClient(conn),
		clusterInfoClient: protos.NewClusterInfoServiceClient(conn),
	}
}

// connectionAndEndpoints represents a combination of a gRPC client connection and server endpoints.
type connectionAndEndpoints struct {
	conn      *connection
	endpoints *protos.ServerEndpointList
}

// newConnAndEndpoints creates a new connectionAndEndpoints instance.
func newConnAndEndpoints(conn *connection, endpoints *protos.ServerEndpointList) *connectionAndEndpoints {
	return &connectionAndEndpoints{
		conn:      conn,
		endpoints: endpoints,
	}
}

// connectionProvider is responsible for managing gRPC client connections to
// Aerospike servers.
//
//nolint:govet // We will favor readability over field alignment
type connectionProvider struct {
	logger         *slog.Logger
	nodeConns      map[uint64]*connectionAndEndpoints
	seedConns      []*connection
	tlsConfig      *tls.Config
	seeds          HostPortSlice
	nodeConnsLock  *sync.RWMutex
	tendInterval   time.Duration
	clusterID      uint64
	listenerName   *string
	isLoadBalancer bool
	token          tokenManager
	stopTendChan   chan struct{}
	closed         atomic.Bool
}

// newConnectionProvider creates a new connectionProvider instance.
func newConnectionProvider(
	ctx context.Context,
	seeds HostPortSlice,
	listenerName *string,
	isLoadBalancer bool,
	token tokenManager,
	tlsConfig *tls.Config,
	logger *slog.Logger,
) (*connectionProvider, error) {
	// Initialize the logger.
	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.WithGroup("cp")

	// Validate the seeds.
	if len(seeds) == 0 {
		msg := "seeds cannot be nil or empty"
		logger.Error(msg)

		return nil, errors.New(msg)
	}

	if token != nil {
		if token.RequireTransportSecurity() && tlsConfig == nil {
			msg := "tlsConfig is required when username/password authentication"
			logger.Error(msg)

			return nil, errors.New(msg)
		}
	}

	// Create the connectionProvider instance.
	cp := &connectionProvider{
		nodeConns:      make(map[uint64]*connectionAndEndpoints),
		seeds:          seeds,
		listenerName:   listenerName,
		isLoadBalancer: isLoadBalancer,
		token:          token,
		tlsConfig:      tlsConfig,
		tendInterval:   time.Second * 1,
		nodeConnsLock:  &sync.RWMutex{},
		stopTendChan:   make(chan struct{}),
		logger:         logger,
		closed:         atomic.Bool{},
	}

	// Connect to the seed nodes.
	err := cp.connectToSeeds(ctx)
	if err != nil {
		cp.Close()
		logger.Error("failed to connect to seeds", slog.Any("error", err))
		return nil, err
	}

	// Schedule token refresh if token manager is present.
	if token != nil {
		cp.token.ScheduleRefresh(cp.GetRandomConn)
	}

	// Start the tend routine if load balancing is disabled.
	if !isLoadBalancer {
		cp.updateClusterConns(ctx) // We want at least one tend to occur before we return

		cp.logger.Debug("starting tend routine")
		go cp.tend(context.Background()) // Might add a tend specific timeout in the future?
	} else {
		cp.logger.Debug("load balancer is enabled, not starting tend routine")
	}

	return cp, nil
}

// Close closes the connectionProvider and releases all resources.
func (cp *connectionProvider) Close() error {
	if cp == nil {
		return nil
	}

	if !cp.isLoadBalancer {
		cp.stopTendChan <- struct{}{}
		<-cp.stopTendChan
	}

	var firstErr error

	for _, conn := range cp.seedConns {
		err := conn.grpcConn.Close()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}

			cp.logger.Error("failed to close seed connection",
				slog.Any("error", err),
				slog.String("seed", conn.grpcConn.Target()),
			)
		}
	}

	for _, conn := range cp.nodeConns {
		err := conn.conn.grpcConn.Close()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}

			cp.logger.Error("failed to close node connection",
				slog.Any("error", err),
				slog.String("node", conn.conn.grpcConn.Target()),
			)
		}
	}

	cp.logger.Debug("closed")
	cp.closed.Store(true)

	return firstErr
}

// GetSeedConn returns a gRPC client connection to a seed node.
func (cp *connectionProvider) GetSeedConn() (*connection, error) {
	if cp.closed.Load() {
		cp.logger.Warn("ConnectionProvider is closed, cannot get connection")
		return nil, errConnectionProviderClosed
	}

	if len(cp.seedConns) == 0 {
		msg := "no seed connections found"
		cp.logger.Warn(msg)

		return nil, errors.New(msg)
	}

	idx := rand.Intn(len(cp.seedConns)) //nolint:gosec // Security is not an issue here

	return cp.seedConns[idx], nil
}

// GetRandomConn returns a gRPC client connection to an Aerospike server. If
// isLoadBalancer is enabled, it will return the seed connection.
func (cp *connectionProvider) GetRandomConn() (*connection, error) {
	if cp.closed.Load() {
		cp.logger.Warn("ConnectionProvider is closed, cannot get connection")
		return nil, errors.New("ConnectionProvider is closed")
	}

	if cp.isLoadBalancer {
		cp.logger.Debug("load balancer is enabled, using seed connection")
		return cp.GetSeedConn()
	}

	cp.nodeConnsLock.RLock()
	defer cp.nodeConnsLock.RUnlock()

	discoverdConns := make([]*connectionAndEndpoints, len(cp.nodeConns))

	i := 0

	for _, conn := range cp.nodeConns {
		discoverdConns[i] = conn
		i++
	}

	if len(discoverdConns) == 0 {
		cp.logger.Warn("no node connections found, using seed connection")
		return cp.GetSeedConn()
	}

	idx := rand.Intn(len(discoverdConns)) //nolint:gosec // Security is not an issue here

	return discoverdConns[idx].conn, nil
}

// GetNodeConn returns a gRPC client connection to a specific node. If the node
// ID cannot be found an error is returned.
func (cp *connectionProvider) GetNodeConn(nodeID uint64) (*connection, error) {
	if cp.closed.Load() {
		cp.logger.Warn("ConnectionProvider is closed, cannot get connection")
		return nil, errors.New("ConnectionProvider is closed")
	}

	if cp.isLoadBalancer {
		cp.logger.Error("load balancer is enabled, using seed connection")
		return nil, errors.New("load balancer is enabled, cannot get specific node connection")
	}

	cp.nodeConnsLock.RLock()
	defer cp.nodeConnsLock.RUnlock()

	conn, ok := cp.nodeConns[nodeID]
	if !ok {
		msg := "connection not found for specified node id"
		cp.logger.Error(msg, slog.Uint64("node", nodeID))

		return nil, errors.New(msg)
	}

	return conn.conn, nil
}

// GetNodeIDs returns the node IDs of all nodes discovered during cluster
// tending. If tending is disabled (LB true) then no node IDs are returned.
func (cp *connectionProvider) GetNodeIDs() []uint64 {
	cp.nodeConnsLock.RLock()
	defer cp.nodeConnsLock.RUnlock()

	nodeIDs := make([]uint64, 0, len(cp.nodeConns))

	for node := range cp.nodeConns {
		nodeIDs = append(nodeIDs, node)
	}

	return nodeIDs
}

// connectToSeeds connects to the seed nodes and creates gRPC client connections.
func (cp *connectionProvider) connectToSeeds(ctx context.Context) error {
	if len(cp.seedConns) != 0 {
		msg := "seed connections already exist, close them first"
		cp.logger.Error(msg)

		return errors.New(msg)
	}

	var authErr error

	wg := sync.WaitGroup{}
	seedGrpcConns := make(chan grpcClientConn)
	cp.seedConns = []*connection{}
	tokenLock := sync.Mutex{} // Ensures only one thread attempts to update token at a time
	tokenUpdated := false     // Ensures token update only occurs once

	for _, seed := range cp.seeds {
		wg.Add(1)

		go func(seed *HostPort) {
			defer wg.Done()

			logger := cp.logger.With(slog.String("host", seed.String()))

			grpcConn, err := cp.createGrcpConn(seed)
			if err != nil {
				logger.ErrorContext(ctx, "failed to create connection", slog.Any("error", err))
				grpcConn.Close()
				return
			}

			extraCheck := true

			if cp.token != nil {
				// Only one thread needs to refresh the token. Only first will
				// succeed others will block
				tokenLock.Lock()
				if !tokenUpdated {
					err := cp.token.RefreshToken(ctx, grpcConn)
					if err != nil {
						logger.WarnContext(ctx, "failed to refresh token", slog.Any("error", err))
						authErr = err
						tokenLock.Unlock()
						grpcConn.Close()
						return
					}

					// No need to check this conn again for successful connectivity
					extraCheck = false
					tokenUpdated = true
				}
				tokenLock.Unlock()
			}

			if extraCheck {
				client := protos.NewAboutServiceClient(grpcConn)

				about, err := client.Get(ctx, &protos.AboutRequest{})
				if err != nil {
					logger.WarnContext(ctx, "failed to connect to seed", slog.Any("error", err))
					grpcConn.Close()
					return
				}

				if newVersion(about.Version).lt(minimumFullySupportedAVSVersion) {
					logger.WarnContext(ctx, "incompatible server version", slog.String("version", about.Version))
				}
			}

			seedGrpcConns <- grpcConn
		}(seed)
	}

	go func() {
		wg.Wait()
		close(seedGrpcConns)
	}()

	for conn := range seedGrpcConns {
		cp.seedConns = append(cp.seedConns, newConnection(conn))
	}

	if len(cp.seedConns) == 0 {
		msg := "failed to connect to seeds"

		if authErr != nil {
			return NewAVSErrorFromGrpc(msg, authErr)
		}

		if err := ctx.Err(); err != nil {
			msg = fmt.Sprintf("%s: %s", msg, err)
		}

		return NewAVSError(msg, nil)
	}

	return nil
}

// updateNodeConns updates the gRPC client connection for a specific node.
func (cp *connectionProvider) updateNodeConns(
	ctx context.Context,
	node uint64,
	endpoints *protos.ServerEndpointList,
) error {
	newConn, err := cp.createConnFromEndpoints(endpoints)
	if err != nil {
		return err
	}

	_, err = newConn.aboutClient.Get(ctx, &protos.AboutRequest{})
	if err != nil {
		return err
	}

	cp.nodeConnsLock.Lock()
	cp.nodeConns[node] = newConnAndEndpoints(newConn, endpoints)
	cp.nodeConnsLock.Unlock()

	return nil
}

// checkAndSetClusterID checks if the cluster ID has changed and updates it if necessary.
func (cp *connectionProvider) checkAndSetClusterID(clusterID uint64) bool {
	if clusterID != cp.clusterID {
		cp.clusterID = clusterID
		return true
	}

	return false
}

// getTendConns returns all the gRPC client connections for tend operations.
func (cp *connectionProvider) getTendConns() []grpcClientConn {
	cp.nodeConnsLock.RLock()
	defer cp.nodeConnsLock.RUnlock()

	conns := make([]grpcClientConn, len(cp.seedConns)+len(cp.nodeConns))
	i := 0

	for _, conn := range cp.seedConns {
		conns[i] = conn.grpcConn
		i++
	}

	for _, conn := range cp.nodeConns {
		conns[i] = conn.conn.grpcConn
		i++
	}

	return conns
}

// getUpdatedEndpoints retrieves the updated server endpoints from the Aerospike cluster.
func (cp *connectionProvider) getUpdatedEndpoints(ctx context.Context) map[uint64]*protos.ServerEndpointList {
	conns := cp.getTendConns()
	endpointsChan := make(chan map[uint64]*protos.ServerEndpointList)
	endpointsReq := &protos.ClusterNodeEndpointsRequest{ListenerName: cp.listenerName}
	wg := sync.WaitGroup{}

	for _, conn := range conns {
		wg.Add(1)

		go func(conn grpcClientConn) {
			defer wg.Done()

			logger := cp.logger.With(slog.String("host", conn.Target()))
			client := protos.NewClusterInfoServiceClient(conn)

			clusterID, err := client.GetClusterId(ctx, &emptypb.Empty{})
			if err != nil {
				logger.WarnContext(ctx, "failed to get cluster ID", slog.Any("error", err))
			}

			if !cp.checkAndSetClusterID(clusterID.GetId()) {
				logger.DebugContext(
					ctx,
					"old cluster ID found, skipping connection discovery",
					slog.Uint64("clusterID", clusterID.GetId()),
				)

				return
			}

			logger.DebugContext(ctx, "new cluster ID found", slog.Uint64("clusterID", clusterID.GetId()))

			endpointsResp, err := client.GetClusterEndpoints(ctx, endpointsReq)
			if err != nil {
				logger.ErrorContext(ctx, "failed to get cluster endpoints", slog.Any("error", err))
				return
			}

			endpointsChan <- endpointsResp.Endpoints
		}(conn)
	}

	go func() {
		wg.Wait()
		close(endpointsChan)
	}()

	// Stores the endpoints from the node with the largest view of the cluster
	var maxTempEndpoints map[uint64]*protos.ServerEndpointList
	for endpoints := range endpointsChan {
		if maxTempEndpoints == nil || len(endpoints) > len(maxTempEndpoints) {
			maxTempEndpoints = endpoints
			cp.logger.DebugContext(ctx, "found new cluster ID", slog.Any("endpoints", maxTempEndpoints))
		}
	}

	return maxTempEndpoints
}

// Checks if the node connections need to be updated and updates them if necessary.
func (cp *connectionProvider) checkAndSetNodeConns(
	ctx context.Context,
	newNodeEndpoints map[uint64]*protos.ServerEndpointList,
) {
	wg := sync.WaitGroup{}
	// Find which nodes have a different endpoint list and update their connection
	for node, newEndpoints := range newNodeEndpoints {
		wg.Add(1)

		go func(node uint64, newEndpoints *protos.ServerEndpointList) {
			defer wg.Done()

			logger := cp.logger.With(slog.Uint64("node", node))

			cp.nodeConnsLock.RLock()
			currEndpoints, ok := cp.nodeConns[node]
			cp.nodeConnsLock.RUnlock()

			if ok {
				if !endpointListEqual(currEndpoints.endpoints, newEndpoints) {
					logger.Debug("endpoints for node changed, recreating connection")

					err := currEndpoints.conn.grpcConn.Close()
					if err != nil {
						logger.Warn("failed to close connection", slog.Any("error", err))
					}

					// Either this is a new node or its endpoints have changed
					err = cp.updateNodeConns(ctx, node, newEndpoints)
					if err != nil {
						logger.Error("failed to create new connection", slog.Any("error", err))
					}
				} else {
					cp.logger.Debug("endpoints for node unchanged")
				}
			} else {
				logger.Debug("new node found, creating new connection")

				err := cp.updateNodeConns(ctx, node, newEndpoints)
				if err != nil {
					logger.Error("failed to create new connection", slog.Any("error", err))
				}
			}
		}(node, newEndpoints)
	}

	wg.Wait()
}

// removeDownNodes removes the gRPC client connections for nodes in nodeConns
// that aren't apart of newNodeEndpoints
func (cp *connectionProvider) removeDownNodes(newNodeEndpoints map[uint64]*protos.ServerEndpointList) {
	cp.nodeConnsLock.Lock()
	defer cp.nodeConnsLock.Unlock()

	// The cluster state changed. Remove old connections.
	for node, connEndpoints := range cp.nodeConns {
		if _, ok := newNodeEndpoints[node]; !ok {
			err := connEndpoints.conn.grpcConn.Close()
			if err != nil {
				cp.logger.Warn("failed to close connection", slog.Uint64("node", node), slog.Any("error", err))
			}

			delete(cp.nodeConns, node)
		}
	}
}

// updateClusterConns updates the gRPC client connections for the Aerospike
// cluster if the cluster state has changed.
func (cp *connectionProvider) updateClusterConns(ctx context.Context) {
	updatedEndpoints := cp.getUpdatedEndpoints(ctx)
	if updatedEndpoints == nil {
		cp.logger.Debug("no new cluster ID found, cluster state is unchanged, skipping connection discovery")
		return
	}

	cp.logger.Debug("new cluster id found, updating connections")

	cp.checkAndSetNodeConns(ctx, updatedEndpoints)
	cp.removeDownNodes(updatedEndpoints)
}

// tend starts a thread to periodically update the cluster connections.
func (cp *connectionProvider) tend(ctx context.Context) {
	timer := time.NewTimer(cp.tendInterval)
	defer timer.Stop()

	for {
		timer.Reset(cp.tendInterval)

		select {
		case <-timer.C:
			cp.logger.Debug("tending . . .")

			ctx, cancel := context.WithTimeout(ctx, cp.tendInterval) // TODO: make configurable?

			cp.updateClusterConns(ctx)

			if err := ctx.Err(); err != nil {
				cp.logger.Warn("tend context cancelled", slog.Any("error", err))
			}

			cp.logger.Debug("finished tend")

			cancel()
		case <-cp.stopTendChan:
			cp.stopTendChan <- struct{}{}
			return
		}
	}
}

func endpointEqual(a, b *protos.ServerEndpoint) bool {
	return a.Address == b.Address && a.Port == b.Port && a.IsTls == b.IsTls
}

func endpointListEqual(a, b *protos.ServerEndpointList) bool {
	if len(a.Endpoints) != len(b.Endpoints) {
		return false
	}

	aEndpoints := make([]*protos.ServerEndpoint, len(a.Endpoints))
	copy(aEndpoints, a.Endpoints)

	bEndpoints := make([]*protos.ServerEndpoint, len(b.Endpoints))
	copy(bEndpoints, b.Endpoints)

	sortFunc := func(endpoints []*protos.ServerEndpoint) func(int, int) bool {
		return func(i, j int) bool {
			if endpoints[i].Address < endpoints[j].Address {
				return true
			} else if endpoints[i].Address > endpoints[j].Address {
				return false
			}

			return endpoints[i].Port < endpoints[j].Port
		}
	}

	sort.Slice(aEndpoints, sortFunc(aEndpoints))
	sort.Slice(bEndpoints, sortFunc(bEndpoints))

	for i, endpoint := range aEndpoints {
		if !endpointEqual(endpoint, bEndpoints[i]) {
			return false
		}
	}

	return true
}

func endpointToHostPort(endpoint *protos.ServerEndpoint) *HostPort {
	return NewHostPort(endpoint.Address, int(endpoint.Port))
}

// createGrpcConnFromEndpoints creates a gRPC client connection from the first
// successful endpoint in endpoints.
func (cp *connectionProvider) createGrpcConnFromEndpoints(
	endpoints *protos.ServerEndpointList,
) (grpcClientConn, error) {
	for _, endpoint := range endpoints.Endpoints {
		if strings.ContainsRune(endpoint.Address, ':') {
			continue // TODO: Add logging and support for IPv6
		}

		conn, err := cp.createGrcpConn(endpointToHostPort(endpoint))

		if err == nil {
			return conn, nil
		}
	}

	return nil, errors.New("no valid endpoint found")
}

// createGrcpConn creates a gRPC client connection to a host. This handles adding
// credential and configuring tls.
func (cp *connectionProvider) createGrcpConn(hostPort *HostPort) (grpcClientConn, error) {
	opts := []grpc.DialOption{}

	if cp.tlsConfig == nil {
		cp.logger.Info("using insecure connection to host", slog.String("host", hostPort.String()))

		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		cp.logger.Info("using secure tls connection to host", slog.String("host", hostPort.String()))

		opts = append(opts, grpc.WithTransportCredentials(grpcCreds.NewTLS(cp.tlsConfig)))
	}

	if cp.token != nil {
		opts = append(opts,
			grpc.WithUnaryInterceptor(cp.token.UnaryInterceptor()),
			grpc.WithStreamInterceptor(cp.token.StreamInterceptor()),
		)
	}

	conn, err := grpc.NewClient(
		hostPort.toDialString(),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (cp *connectionProvider) createConnFromEndpoints(endpoints *protos.ServerEndpointList) (*connection, error) {
	conn, err := cp.createGrpcConnFromEndpoints(endpoints)
	if err != nil {
		return nil, err
	}

	return newConnection(conn), nil
}
