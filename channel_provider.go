package avs

import (
	"errors"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aerospike/aerospike-proximus-client-go/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	logger          *slog.Logger
	nodeConns       map[uint64]*ChannelAndEndpoints
	seedConns       []*grpc.ClientConn
	seeds           []*HostPort
	tendLock        *sync.RWMutex
	tendInterval    int // in seconds
	clusterID       uint64
	listenerName    string
	isLoadBalancer  bool
	stopTendChan    chan struct{}
	tendStoppedChan chan struct{}
	closed          bool
}

func NewChannelProvider(
	ctx context.Context,
	seeds []*HostPort,
	listenerName string,
	isLoadBalancer bool,
	logger *slog.Logger,
) (*ChannelProvider, error) {
	if len(seeds) == 0 {
		panic("seeds cannot be nil or empty")
	}

	logger = logger.WithGroup("channel_provider")

	cp := &ChannelProvider{
		nodeConns:       make(map[uint64]*ChannelAndEndpoints),
		seedConns:       nil,
		closed:          false,
		clusterID:       0,
		seeds:           seeds,
		listenerName:    listenerName,
		isLoadBalancer:  isLoadBalancer,
		tendInterval:    1,
		tendLock:        &sync.RWMutex{},
		stopTendChan:    make(chan struct{}),
		tendStoppedChan: make(chan struct{}),
		logger:          logger,
	}

	err := cp.connectToSeeds(ctx)
	if err != nil {
		logger.Error("failed to connect to seeds", slog.Any("error", err))
		return nil, err
	}

	if !isLoadBalancer {
		go cp.tend()
	} else {
		cp.logger.Debug("load balancer is enabled, not starting tend routine")
	}

	return cp, nil
}

func (cp *ChannelProvider) Close() {
	if cp.isLoadBalancer {
		cp.stopTendChan <- struct{}{}
		<-cp.tendStoppedChan
	}

	for _, channel := range cp.seedConns {
		channel.Close()
	}

	for _, channel := range cp.nodeConns {
		channel.Channel.Close()
	}

	cp.logger.Debug("closed")
	cp.closed = true
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

	cp.tendLock.RLock()
	defer cp.tendLock.RUnlock()

	discoverdChannels := make([]*ChannelAndEndpoints, len(cp.nodeConns))
	i := 0

	for _, channel := range cp.nodeConns {
		discoverdChannels[i] = channel
		i++
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
		cp.logger.Error("seed channels already exist, close them first")
		return errors.New("seed channels already exist, close them first")
	}

	success := false
	wg := sync.WaitGroup{}
	seedCons := make(chan *grpc.ClientConn)
	cp.seedConns = []*grpc.ClientConn{}

	defer close(seedCons)

	for _, seed := range cp.seeds {
		wg.Add(1)

		go func(seed *HostPort) {
			defer wg.Done()

			conn, err := createChannel(ctx, seed)
			if err != nil {
				cp.logger.Warn("failed to connect to seed", slog.String("address", seed.String()), slog.Any("error", err))
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
		success = true

		cp.seedConns = append(cp.seedConns, conn)
	}

	if !success {
		cp.logger.Error("failed to connect to any seed")
		return errors.New("failed to connect to any seed")
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

	cp.tendLock.Lock()
	cp.nodeConns[node] = NewChannelAndEndpoints(newChannel, endpoints)
	cp.tendLock.Unlock()

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

func (cp *ChannelProvider) getUpdatedEndpoints() map[uint64]*protos.ServerEndpointList {
	conns := cp.getTendConns()
	endpointsChan := make(chan map[uint64]*protos.ServerEndpointList)
	wg := sync.WaitGroup{}

	for _, conn := range conns {
		wg.Add(1)

		go func(conn *grpc.ClientConn) {
			defer wg.Done()

			client := protos.NewClusterInfoClient(conn)

			clusterID, err := client.GetClusterId(context.Background(), nil)
			if err != nil {
				cp.logger.Warn("failed to get cluster ID", slog.Any("error", err))
			}

			if !cp.checkAndSetClusterID(clusterID.GetId()) {
				cp.logger.Debug("old cluster ID found, skipping channel discovery")
				return
			}

			endpointsReq := &protos.ClusterNodeEndpointsRequest{ListenerName: &cp.listenerName}

			endpointsResp, err := client.GetClusterEndpoints(context.Background(), endpointsReq)
			if err != nil {
				cp.logger.Error("failed to get cluster endpoints", slog.Any("error", err))
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

func (cp *ChannelProvider) checkAndSetNodeConns(newNodeEndpoints map[uint64]*protos.ServerEndpointList) {
	wg := sync.WaitGroup{}

	// Find which nodes have a different endpoint list and update their channel
	for node, newEndpoints := range newNodeEndpoints {
		wg.Add(1)

		go func(node uint64, newEndpoints *protos.ServerEndpointList) {
			defer wg.Done()

			addNewConn := true

			if currEndpoints, ok := cp.nodeConns[node]; ok {
				if endpointListEqual(currEndpoints.Endpoints, newEndpoints) {
					addNewConn = false
				} else {
					cp.logger.Debug("endpoints for node changed, recreating channel", slog.Uint64("node", node))

					addNewConn = true

					currEndpoints.Channel.Close()
				}
			}

			if addNewConn {
				// Either this is a new node or its endpoints have changed
				err := cp.updateNodeConns(context.TODO(), node, newEndpoints)
				if err != nil {
					cp.logger.Error("failed to create new channel", slog.Uint64("node", node), slog.Any("error", err))
				}
			}
		}(node, newEndpoints)
	}

	wg.Wait()
}

func (cp *ChannelProvider) removeDownNodes(newNodeEndpoints map[uint64]*protos.ServerEndpointList) {
	cp.tendLock.Lock()

	// The cluster state changed. Remove old channels.
	for node, channelEndpoints := range cp.nodeConns {
		if _, ok := newNodeEndpoints[node]; !ok {
			channelEndpoints.Channel.Close()
			delete(cp.nodeConns, node)
		}
	}

	cp.tendLock.Unlock()
}

func (cp *ChannelProvider) updateClusterChannels() {
	updatedEndpoints := cp.getUpdatedEndpoints()
	if updatedEndpoints == nil {
		cp.logger.Debug("no new cluster ID found, cluster state is unchanged, skipping channel discovery")
		return
	}

	cp.checkAndSetNodeConns(updatedEndpoints)
	cp.removeDownNodes(updatedEndpoints)
}

func (cp *ChannelProvider) tend() {
	for {
		select {
		case <-time.After(time.Duration(cp.tendInterval) * time.Second):
			cp.updateClusterChannels()
		case <-cp.stopTendChan:
			cp.tendStoppedChan <- struct{}{}
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

	for i, endpoint := range a.Endpoints {
		if !endpointEqual(endpoint, b.Endpoints[i]) {
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
			hostPortToDialString(endpointToHostPort(endpoint)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err == nil {
			return conn, nil
		}
	}

	return nil, errors.New("no valid endpoint found")
}

func createChannel(ctx context.Context, hostPort *HostPort) (*grpc.ClientConn, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	conn, err := grpc.DialContext(
		ctx,
		hostPortToDialString(hostPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func hostPortToDialString(hostPort *HostPort) string {
	return hostPort.Host + ":" + strconv.Itoa(hostPort.Port)
}
