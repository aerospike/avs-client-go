package avs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aerospike/avs-client-go/protos"
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

type UserPassCredentials struct {
	username string
	password string
}

func NewCredentialsFromUserPass(username, password string) *UserPassCredentials {
	return &UserPassCredentials{
		username: username,
		password: password,
	}
}

var ErrNotImplemented = errors.New("not implemented")

type Record struct {
	Data       map[string]any
	Expiration uint32
	Generation uint32
}

func newRecordFromProto(rec *protos.Record) *Record {
	metadata := rec.GetAerospikeMetadata()
	return &Record{
		Data:       protos.ConvertFromFields(rec.Fields),
		Expiration: metadata.Expiration,
		Generation: metadata.Generation,
	}
}
