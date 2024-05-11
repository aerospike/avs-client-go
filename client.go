package aerospike_vector_search

import (
	"github.com/aerospike/aerospike-proximus-client-go/protos"
	"golang.org/x/crypto/openpgp/errors"
)

type Client struct {
	seeds          []HostPort
	listenerName   string
	isLoadBalancer bool
}

func NewClient(seeds []HostPort, listenerName string, isLoadBalancer bool) *Client {
	return &Client{
		seeds:          seeds,
		listenerName:   listenerName,
		isLoadBalancer: isLoadBalancer,
	}
}

func (c *Client) Get(namespace string, setName string, key interface{}, binNames []string) (*protos.Record, error) {
	return nil, errors.UnsupportedError("get() is not implemented")
}

func (c *Client) Delete(namespace string, setName string, key interface{}) (*protos.Record, error) {
	return nil, errors.UnsupportedError("delete() is not implemented")
}

func (c *Client) Exists(namespace string, setName string, key interface{}) (bool, error) {
	return false, errors.UnsupportedError("exists() is not implemented")
}

func (c *Client) IsIndexed(namespace string, setName string, indexName string, key interface{}) (bool, error) {
	return false, errors.UnsupportedError("IsIndexed() is not implemented")
}

func (c *Client) VectorSearch(namespace string, indexName string, query []float32, limit int, searchParams *protos.HnswSearchParams, binNames []string) ([]*protos.Neighbor, error) {
	return nil, errors.UnsupportedError("vectorSearch() is not implemented")
}

func (c *Client) WaitForIndexCompletion(namespace string, indexName string, timeout int) error {
	return errors.UnsupportedError("WaitForIndexCompletion() is not implemented")
}

func (c *Client) Close() {
}
