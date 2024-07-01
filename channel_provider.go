// Package avs provides a channel provider for connecting to Aerospike servers.
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
	"time"

	"github.com/aerospike/avs-client-go/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// channelAndEndpoints represents a combination of a gRPC client connection and server endpoints.
type channelAndEndpoints struct {
	Channel   *grpc.ClientConn
	Endpoints *protos.ServerEndpointList
}

// newChannelAndEndpoints creates a new channelAndEndpoints instance.
func newChannelAndEndpoints(channel *grpc.ClientConn, endpoints *protos.ServerEndpointList) *channelAndEndpoints {
	return &channelAndEndpoints{
		Channel:   channel,
		Endpoints: endpoints,
	}
}

// channelProvider is responsible for managing gRPC client connections to Aerospike servers.
type channelProvider struct {
	logger         *slog.Logger
	nodeConns      map[uint64]*channelAndEndpoints
	seedConns      []*grpc.ClientConn
	tlsConfig      *tls.Config
	seeds          HostPortSlice
	nodeConnsLock  *sync.RWMutex
	tendInterval   time.Duration
	clusterID      uint64
	listenerName   *string
	isLoadBalancer bool
	token          *tokenManager
	stopTendChan   chan struct{}
	closed         bool
}

// newChannelProvider creates a new channelProvider instance.
func newChannelProvider(
	ctx context.Context,
	seeds HostPortSlice,
	listenerName *string,
	isLoadBalancer bool,
	username *string,
	password *string,
	tlsConfig *tls.Config,
	logger *slog.Logger,
) (*channelProvider, error) {
	// Initialize the logger.
	logger = logger.WithGroup("cp")

	// Validate the seeds.
	if len(seeds) == 0 {
		msg := "seeds cannot be nil or empty"
		logger.Error(msg)
		return nil, errors.New(msg)
	}

	// Create a token manager if username and password are provided.
	var token *tokenManager
	if username != nil || password != nil {
		if username == nil || password == nil {
			// Either both are set or neither are set
			msg := "username and password must both be set"
			logger.Error(msg)
			return nil, errors.New(msg)
		}

		token = newJWTToken(*username, *password, logger)

		if token.RequireTransportSecurity() && tlsConfig == nil {
			msg := "tlsConfig is required when username/password authentication"
			logger.Error(msg)
			return nil, errors.New(msg)
		}
	}

	// Create the channelProvider instance.
	cp := &channelProvider{
		nodeConns:      make(map[uint64]*channelAndEndpoints),
		seeds:          seeds,
		listenerName:   listenerName,
		isLoadBalancer: isLoadBalancer,
		token:          token,
		tlsConfig:      tlsConfig,
		tendInterval:   time.Second * 1,
		nodeConnsLock:  &sync.RWMutex{},
		stopTendChan:   make(chan struct{}),
		logger:         logger,
	}

	// Connect to the seed nodes.
	err := cp.connectToSeeds(ctx)
	if err != nil {
		logger.Error("failed to connect to seeds", slog.Any("error", err))
		return nil, err
	}

	// Schedule token refresh if token manager is present.
	if token != nil {
		cp.token.ScheduleRefresh(cp.GetConn)
	}

	// Start the tend routine if load balancing is disabled.
	if !isLoadBalancer {
		cp.logger.Debug("starting tend routine")
		cp.updateClusterChannels(ctx)    // We want at least one tend to occur before we return
		go cp.tend(context.Background()) // Might add a tend specific timeout in the future?
	} else {
		cp.logger.Debug("load balancer is enabled, not starting tend routine")
	}

	return cp, nil
}

// Close closes the channelProvider and releases all resources.
func (cp *channelProvider) Close() error {
	if !cp.isLoadBalancer {
		cp.stopTendChan <- struct{}{}
		<-cp.stopTendChan
	}

	var firstErr error

	if cp.token != nil {
		cp.token.Close()
	}

	for _, channel := range cp.seedConns {
		err := channel.Close()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}

			cp.logger.Error("failed to close seed channel",
				slog.Any("error", err),
				slog.String("seed", channel.Target()),
			)
		}
	}

	for _, channel := range cp.nodeConns {
		err := channel.Channel.Close()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}

			cp.logger.Error("failed to close node channel",
				slog.Any("error", err),
				slog.String("node", channel.Channel.Target()),
			)
		}
	}

	cp.logger.Debug("closed")
	cp.closed = true

	return firstErr
}

// GetConn returns a gRPC client connection to an Aerospike server.
func (cp *channelProvider) GetConn() (*grpc.ClientConn, error) {
	if cp.closed {
		cp.logger.Warn("ChannelProvider is closed, cannot get channel")
		return nil, errors.New("ChannelProvider is closed")
	}

	if cp.isLoadBalancer {
		cp.logger.Debug("load balancer is enabled, using seed channel")
		return cp.seedConns[0], nil
	}

	cp.nodeConnsLock.RLock()
	defer cp.nodeConnsLock.RUnlock()

	discoverdChannels := make([]*channelAndEndpoints, len(cp.nodeConns))

	for i, channel := range cp.nodeConns {
		discoverdChannels[i] = channel
	}

	if len(discoverdChannels) == 0 {
		cp.logger.Warn("no node channels found, using seed channel")
		return cp.seedConns[0], nil
	}

	idx := rand.Intn(len(discoverdChannels)) //nolint:gosec // Security is not an issue here

	return discoverdChannels[idx].Channel, nil
}

// connectToSeeds connects to the seed nodes and creates gRPC client connections.
func (cp *channelProvider) connectToSeeds(ctx context.Context) error {
	if len(cp.seedConns) != 0 {
		msg := "seed channels already exist, close them first"
		cp.logger.Error(msg)
		return errors.New(msg)
	}

	var authErr error

	wg := sync.WaitGroup{}
	seedCons := make(chan *grpc.ClientConn)
	cp.seedConns = []*grpc.ClientConn{}
	tokenLock := sync.Mutex{} // Ensures only one thread attempts to update token at a time
	tokenUpdated := false     // Ensures token update only occurs once

	for _, seed := range cp.seeds {
		wg.Add(1)

		go func(seed *HostPort) {
			defer wg.Done()
			logger := cp.logger.With(slog.String("host", seed.String()))

			conn, err := createChannel(ctx, seed)
			if err != nil {
				logger.ErrorContext(ctx, "failed to create channel", slog.Any("error", err))
				return
			}

			extraCheck := true

			if cp.token != nil {
				// Only one thread needs to refresh the token. Only first will
				// succeed others will block
				tokenLock.Lock()
				if !tokenUpdated {
					err := cp.token.RefreshToken(ctx, conn)
					if err != nil {
						logger.WarnContext(ctx, "failed to refresh token", slog.Any("error", err))
						authErr = err
						return
					}

					// No need to check this conn again for successful connectivity
					extraCheck = false
					tokenUpdated = true

				}
				tokenLock.Unlock()
			}

			// TODO: Check compatible client/server version here
			if extraCheck {
				client := protos.NewClusterInfoClient(conn)

				_, err = client.GetClusterId(ctx, &emptypb.Empty{})
				if err != nil {
					logger.WarnContext(ctx, "failed to connect to seed", slog.Any("error", err))
					return
				}
			}

			seedCons <- conn
		}(seed)
	}

	go func() {
		wg.Wait()
		close(seedCons)
	}()

	for conn := range seedCons {
		cp.seedConns = append(cp.seedConns, conn)
	}

	if len(cp.seedConns) == 0 {
		msg := "failed to connect to seeds"

		if authErr != nil {
			return NewAVSErrorFromGrpc(msg, authErr)
		}

		if err := ctx.Err(); err != nil {
			msg = fmt.Sprintf("%s: %s", msg, err)
		}

		return NewAVSError(msg)
	}

	return nil
}

// updateNodeConns updates the gRPC client connection for a specific node.
func (cp *channelProvider) updateNodeConns(
	node uint64,
	endpoints *protos.ServerEndpointList,
) error {
	newChannel, err := cp.createChannelFromEndpoints(endpoints)
	if err != nil {
		return err
	}

	cp.nodeConnsLock.Lock()
	cp.nodeConns[node] = newChannelAndEndpoints(newChannel, endpoints)
	cp.nodeConnsLock.Unlock()

	return nil
}

// checkAndSetClusterID checks if the cluster ID has changed and updates it if necessary.
func (cp *channelProvider) checkAndSetClusterID(clusterID uint64) bool {
	if clusterID != cp.clusterID {
		cp.clusterID = clusterID
		return true
	}

	return false
}

// getTendConns returns all the gRPC client connections for tend operations.
func (cp *channelProvider) getTendConns() []*grpc.ClientConn {
	cp.nodeConnsLock.RLock()
	defer cp.nodeConnsLock.RUnlock()

	channels := make([]*grpc.ClientConn, len(cp.seedConns)+len(cp.nodeConns))
	i := 0

	for _, channel := range cp.seedConns {
		channels[i] = channel
		i++
	}

	for _, channel := range cp.nodeConns {
		channels[i] = channel.Channel
		i++
	}

	return channels
}

// getUpdatedEndpoints retrieves the updated server endpoints from the Aerospike cluster.
func (cp *channelProvider) getUpdatedEndpoints(ctx context.Context) map[uint64]*protos.ServerEndpointList {
	conns := cp.getTendConns()
	endpointsChan := make(chan map[uint64]*protos.ServerEndpointList)
	endpointsReq := &protos.ClusterNodeEndpointsRequest{ListenerName: cp.listenerName}
	wg := sync.WaitGroup{}

	for _, conn := range conns {
		wg.Add(1)

		go func(conn *grpc.ClientConn) {
			defer wg.Done()

			logger := cp.logger.With(slog.String("host", conn.Target()))
			client := protos.NewClusterInfoClient(conn)

			clusterID, err := client.GetClusterId(ctx, &emptypb.Empty{})
			if err != nil {
				logger.Warn("failed to get cluster ID", slog.Any("error", err))
			}

			if !cp.checkAndSetClusterID(clusterID.GetId()) {
				logger.Debug("old cluster ID found, skipping channel discovery")
				return
			}

			endpointsResp, err := client.GetClusterEndpoints(ctx, endpointsReq)
			if err != nil {
				logger.Error("failed to get cluster endpoints", slog.Any("error", err))
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
		}
	}

	return maxTempEndpoints
}

// checkAndSetNodeConns checks if the node connections need to be updated and updates them if necessary.
func (cp *channelProvider) checkAndSetNodeConns(newNodeEndpoints map[uint64]*protos.ServerEndpointList) {
	wg := sync.WaitGroup{}

	// Find which nodes have a different endpoint list and update their channel
	for node, newEndpoints := range newNodeEndpoints {
		wg.Add(1)

		go func(node uint64, newEndpoints *protos.ServerEndpointList) {
			defer wg.Done()

			logger := cp.logger.With(slog.Uint64("node", node))

			cp.nodeConnsLock.RLock()
			currEndpoints, ok := cp.nodeConns[node]
			cp.nodeConnsLock.RUnlock()

			if ok {
				if !endpointListEqual(currEndpoints.Endpoints, newEndpoints) {
					logger.Debug("endpoints for node changed, recreating channel")

					err := currEndpoints.Channel.Close()
					if err != nil {
						logger.Warn("failed to close channel", slog.Any("error", err))
					}

					// Either this is a new node or its endpoints have changed
					err = cp.updateNodeConns(node, newEndpoints)
					if err != nil {
						logger.Error("failed to create new channel", slog.Any("error", err))
					}
				}
			} else {
				cp.logger.Debug("endpoints for node unchanged", slog.Uint64("node", node))
			}
		}(node, newEndpoints)
	}

	wg.Wait()
}

// removeDownNodes removes the gRPC client connections for nodes in nodeConns
// that aren't apart of newNodeEndpoints
func (cp *channelProvider) removeDownNodes(newNodeEndpoints map[uint64]*protos.ServerEndpointList) {
	cp.nodeConnsLock.Lock()
	defer cp.nodeConnsLock.Unlock()

	// The cluster state changed. Remove old channels.
	for node, channelEndpoints := range cp.nodeConns {
		if _, ok := newNodeEndpoints[node]; !ok {
			err := channelEndpoints.Channel.Close()
			if err != nil {
				cp.logger.Warn("failed to close channel", slog.Uint64("node", node), slog.Any("error", err))
			}

			delete(cp.nodeConns, node)
		}
	}
}

// updateClusterChannels updates the gRPC client connections for the Aerospike
// cluster if the cluster state has changed.
func (cp *channelProvider) updateClusterChannels(ctx context.Context) {
	updatedEndpoints := cp.getUpdatedEndpoints(ctx)
	if updatedEndpoints == nil {
		cp.logger.Debug("no new cluster ID found, cluster state is unchanged, skipping channel discovery")
		return
	}

	cp.logger.Debug("new endpoints found, updating channels", slog.Any("endpoints", updatedEndpoints))

	cp.checkAndSetNodeConns(updatedEndpoints)
	cp.removeDownNodes(updatedEndpoints)
}

// tend starts a thread to periodically update the cluster channels.
func (cp *channelProvider) tend(ctx context.Context) {
	timer := time.NewTimer(cp.tendInterval)
	defer timer.Stop()

	for {
		timer.Reset(cp.tendInterval)

		select {
		case <-timer.C:
			cp.logger.Debug("tending . . .")

			ctx, cancel := context.WithTimeout(ctx, cp.tendInterval) // TODO: make configurable?

			cp.updateClusterChannels(ctx)

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

// createChannelFromEndpoints creates a gRPC client connection from the first
// successful endpoint in endpoints.
func (cp *channelProvider) createChannelFromEndpoints(
	endpoints *protos.ServerEndpointList,
) (*grpc.ClientConn, error) {
	for _, endpoint := range endpoints.Endpoints {
		if strings.ContainsRune(endpoint.Address, ':') {
			continue // TODO: Add logging and support for IPv6
		}

		conn, err := cp.createChannel(endpointToHostPort(endpoint))

		if err == nil {
			return conn, nil
		}
	}

	return nil, errors.New("no valid endpoint found")
}

// createChannel creates a gRPC client connection to a host. This handles adding
// credential and configuring tls.
func (cp *channelProvider) createChannel(hostPort *HostPort) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}

	if cp.tlsConfig == nil {
		cp.logger.Debug("using insecure connection to host", slog.String("host", hostPort.String()))

		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		cp.logger.Debug("using secure tls connection to host", slog.String("host", hostPort.String()))

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(cp.tlsConfig)))
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
