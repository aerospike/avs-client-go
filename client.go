package avs

import (
	"context"
	"log/slog"

	"github.com/aerospike/aerospike-proximus-client-go/protos"
)

//nolint:govet // We will favor readability over field alignment
type Client struct {
	logger         *slog.Logger
	seeds          []HostPort
	listenerName   string
	isLoadBalancer bool
}

//nolint:revive // TODO
func NewClient(
	ctx context.Context,
	seeds []HostPort,
	listenerName string,
	isLoadBalancer bool,
	logger *slog.Logger,
) *Client {
	logger = logger.WithGroup("aerospike_vector_search")

	return &Client{
		seeds:          seeds,
		listenerName:   listenerName,
		isLoadBalancer: isLoadBalancer,
		logger:         logger,
	}
}

//nolint:revive // TODO
func (c *Client) Get(ctx context.Context,
	namespace,
	setName string,
	key interface{},
	binNames []string,
) (*protos.Record, error) {
	return nil, ErrNotImplemented
}

//nolint:revive // TODO
func (c *Client) Delete(ctx context.Context, namespace, setName string, key interface{}) (*protos.Record, error) {
	return nil, ErrNotImplemented
}

//nolint:revive // TODO
func (c *Client) Exists(
	ctx context.Context,
	namespace,
	setName string,
	key interface{},
) (bool, error) {
	return false, ErrNotImplemented
}

//nolint:revive // TODO
func (c *Client) IsIndexed(ctx context.Context, namespace, setName, indexName string, key interface{}) (bool, error) {
	return false, ErrNotImplemented
}

//nolint:revive // TODO
func (c *Client) VectorSearch(ctx context.Context,
	namespace,
	indexName string,
	query []float32,
	limit int,
	searchParams *protos.HnswSearchParams,
	binNames []string,
) ([]*protos.Neighbor, error) {
	return nil, ErrNotImplemented
}

//nolint:revive // TODO
func (c *Client) WaitForIndexCompletion(ctx context.Context, namespace, indexName string, timeout int) error {
	return ErrNotImplemented
}

func (c *Client) Close() {
}
