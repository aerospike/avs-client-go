package avs

import (
	"errors"
	"fmt"
)

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

func (hp *HostPort) String() string {
	return hp.Host + ":" + fmt.Sprintf("%d", hp.Port)
}

var ErrNotImplemented = errors.New("not implemented")
