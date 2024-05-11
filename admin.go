package aerospike_vector_search

import "github.com/aerospike/aerospike-proximus-client-go/protos"

type AdminClient struct {
	seeds          []HostPort
	listenerName   string
	isLoadBalancer bool
}

func NewAdminClient(seeds []HostPort, listenerName string, isLoadBalancer bool) *Client {
	return &Client{
		seeds:          seeds,
		listenerName:   listenerName,
		isLoadBalancer: isLoadBalancer,
	}
}

func (c *AdminClient) IndexCreate(
	namespace string, 
	sets []string, 
	name string, 
	vectorField string, 
	dimensions int, 
	vectorDistanceMetric protos.VectorDistanceMetric, 
	indexParams protos.HnswParams, 
	indexMetaData map[string]string
) error {
	return errors.UnsupportedError("IndexCreate() is not implemented")
}




