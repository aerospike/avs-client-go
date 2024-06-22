package avs

import (
	"errors"
	"fmt"
	"strings"
)

type HostPort struct {
	Host string
	Port int
}

func NewHostPort(host string, port int) *HostPort {
	return &HostPort{
		Host: host,
		Port: port,
	}
}

func (hp *HostPort) String() string {
	return hp.toDialString()
}

func (hp *HostPort) toDialString() string {
	return fmt.Sprintf("%s:%d", hp.Host, hp.Port)
}

type HostPortSlice []*HostPort

func (hps HostPortSlice) String() string {
	s := make([]string, len(hps))

	for i, hp := range hps {
		s[i] = hp.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(s, ", "))
}

type HostPortSlice []*HostPort

func (hps HostPortSlice) String() string {
	s := make([]string, len(hps))

	for i, hp := range hps {
		s[i] = hp.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(s, ", "))
}

var ErrNotImplemented = errors.New("not implemented")
