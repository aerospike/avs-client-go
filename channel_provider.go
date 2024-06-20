package avs

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aerospike/aerospike-proximus-client-go/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChannelAndEndpoints struct {
	Channel   *grpc.ClientConn
	Endpoints *protos.ServerEndpointList
}

func NewChannelAndEndpoints(channel *grpc.ClientConn, endpoints *protos.ServerEndpointList) *ChannelAndEndpoints {
	return &ChannelAndEndpoints{
		Channel:   channel,
		Endpoints: endpoints,
	}
}

//nolint:govet // We will favor readability over field alignment
type ChannelProvider struct {
	logger         *slog.Logger
	nodeConns      map[uint64]*ChannelAndEndpoints
	seedConns      []*grpc.ClientConn
	seeds          HostPortSlice
	nodeConnsLock  *sync.RWMutex
	tendInterval   time.Duration
	clusterID      uint64
	listenerName   *string
	isLoadBalancer bool
	stopTendChan   chan struct{}
	closed         bool
}

func NewChannelProvider(
	ctx context.Context,
	seeds HostPortSlice,
	listenerName *string,
	isLoadBalancer bool,
	logger *slog.Logger,
) (*ChannelProvider, error) {
	if len(seeds) == 0 {
		return nil, fmt.Errorf("seeds cannot be nil or empty")
	}

	logger = logger.WithGroup("cp")

	cp := &ChannelProvider{
		nodeConns:      make(map[uint64]*ChannelAndEndpoints),
		seeds:          seeds,
		listenerName:   listenerName,
		isLoadBalancer: isLoadBalancer,
		tendInterval:   time.Second * 1,
		nodeConnsLock:  &sync.RWMutex{},
		stopTendChan:   make(chan struct{}),
		logger:         logger,
	}

	err := cp.connectToSeeds(ctx)
	if err != nil {
		logger.Error("failed to connect to seeds", slog.Any("error", err))
		return nil, err
	}

	if !isLoadBalancer {
		cp.updateClusterChannels(ctx)    // We want at least one tend to occur before we return
		go cp.tend(context.Background()) // Might add a tend specific timeout in the future?
	} else {
		cp.logger.Debug("load balancer is enabled, not starting tend routine")
	}

	return cp, nil
}

func (cp *ChannelProvider) Close() error {
	if !cp.isLoadBalancer {
		cp.stopTendChan <- struct{}{}
		<-cp.stopTendChan
	}

	var firstErr error

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

func (cp *ChannelProvider) GetConn() (*grpc.ClientConn, error) {
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

	discoverdChannels := make([]*ChannelAndEndpoints, len(cp.nodeConns))

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

func (cp *ChannelProvider) connectToSeeds(ctx context.Context) error {
	if len(cp.seedConns) != 0 {
		msg := "seed channels already exist, close them first"
		cp.logger.Error(msg)
		return errors.New(msg)
	}

	wg := sync.WaitGroup{}
	seedCons := make(chan *grpc.ClientConn)
	cp.seedConns = []*grpc.ClientConn{}

	for _, seed := range cp.seeds {
		wg.Add(1)

		go func(seed *HostPort) {
			defer wg.Done()

			logger := cp.logger.With(slog.String("host", seed.String()))

			logger := cp.logger.With(slog.String("host", seed.String()))

			conn, err := createChannel(ctx, seed)
			if err != nil {
				logger.ErrorContext(ctx, "failed to create channel", slog.Any("error", err))
				return
			}

			client := protos.NewClusterInfoClient(conn)

			_, err = client.GetClusterId(ctx, &emptypb.Empty{})
			if err != nil {
				logger.WarnContext(ctx, "failed to connect to seed", slog.Any("error", err))
				return
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

		if err := ctx.Err(); err != nil {
			msg = fmt.Sprintf("%s: %s", msg, err.Error())
		}

		return NewAVSError(msg)
	}

	return nil
}

func (cp *ChannelProvider) updateNodeConns(
	ctx context.Context,
	node uint64,
	endpoints *protos.ServerEndpointList,
) error {
	newChannel, err := createChannelFromEndpoints(ctx, endpoints)
	if err != nil {
		return err
	}

	cp.nodeConnsLock.Lock()
	cp.nodeConns[node] = NewChannelAndEndpoints(newChannel, endpoints)
	cp.nodeConnsLock.Unlock()

	return nil
}

func (cp *ChannelProvider) checkAndSetClusterID(clusterID uint64) bool {
	if clusterID != cp.clusterID {
		cp.clusterID = clusterID
		return true
	}

	return false
}

func (cp *ChannelProvider) getTendConns() []*grpc.ClientConn {
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

func (cp *ChannelProvider) getUpdatedEndpoints(ctx context.Context) map[uint64]*protos.ServerEndpointList {
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

func (cp *ChannelProvider) checkAndSetNodeConns(
	ctx context.Context,
	newNodeEndpoints map[uint64]*protos.ServerEndpointList,
) {
	wg := sync.WaitGroup{}

	// Find which nodes have a different endpoint list and update their channel
	for node, newEndpoints := range newNodeEndpoints {
		wg.Add(1)

		go func(node uint64, newEndpoints *protos.ServerEndpointList) {
			defer wg.Done()
			logger := cp.logger.With(slog.Uint64("node", node))

			cp.nodeConnsLock.RLock()
			currEndpoints, ok := cp.nodeConns[node]
			cp.nodeConnsLock.Unlock()

			if ok {
				if !endpointListEqual(currEndpoints.Endpoints, newEndpoints) {
					logger.Debug("endpoints for node changed, recreating channel")

					err := currEndpoints.Channel.Close()
					if err != nil {
						logger.Warn("failed to close channel", slog.Any("error", err))
					}

					// Either this is a new node or its endpoints have changed
					err = cp.updateNodeConns(ctx, node, newEndpoints)
					if err != nil {
						logger.Error("failed to create new channel", slog.Any("error", err))
					}
				}
			}
		}(node, newEndpoints)
	}

	wg.Wait()
}

func (cp *ChannelProvider) removeDownNodes(newNodeEndpoints map[uint64]*protos.ServerEndpointList) {
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

func (cp *ChannelProvider) updateClusterChannels(ctx context.Context) {
	updatedEndpoints := cp.getUpdatedEndpoints(ctx)
	if updatedEndpoints == nil {
		cp.logger.Debug("no new cluster ID found, cluster state is unchanged, skipping channel discovery")
		return
	}

	cp.checkAndSetNodeConns(ctx, updatedEndpoints)
	cp.removeDownNodes(updatedEndpoints)
}

func (cp *ChannelProvider) tend(ctx context.Context) {
	timer := time.NewTimer(cp.tendInterval)
	defer timer.Stop()

	for {
		timer.Reset(cp.tendInterval)

		select {
		case <-timer.C:
			ctx, cancel := context.WithTimeout(ctx, cp.tendInterval) // TODO: make configurable?

			cp.updateClusterChannels(ctx)

			if err := ctx.Err(); err != nil {
				cp.logger.Warn("tend context cancelled", slog.Any("error", err))
			}

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
	return NewHostPort(endpoint.Address, int(endpoint.Port), endpoint.IsTls)
}

func createChannelFromEndpoints(
	ctx context.Context,
	endpoints *protos.ServerEndpointList,
) (*grpc.ClientConn, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	for _, endpoint := range endpoints.Endpoints {
		if strings.ContainsRune(endpoint.Address, ':') {
			continue // TODO: Add logging and support for IPv6
		}

		conn, err := grpc.DialContext(
			ctx,
			endpointToHostPort(endpoint).toDialString(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err == nil {
			return conn, nil
		}
	}

	return nil, errors.New("no valid endpoint found")
}

func createChannel(ctx context.Context, hostPort *HostPort) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(
		ctx,
		hostPort.toDialString(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
