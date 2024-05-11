package aerospike_vector_search

type HostPort struct {
	Host  string
	Port  int
	isTLS bool
}

func NewHostPort(host string, port int, isTLS bool) *HostPort {
	return &HostPort{
		Host:  host,
		Port:  port,
		isTLS: isTLS,
	}
}
