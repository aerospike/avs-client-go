package avs

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
var AerospikeEpoch = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

type Record struct {
	Data       map[string]any
	Expiration time.Time
	Generation uint32
}

func newRecordFromProto(rec *protos.Record) *Record {
	metadata := rec.GetAerospikeMetadata()

	recordData, err := protos.ConvertFromFields(rec.GetFields())
	// TODO: Should we ignore unsupported fields by default? Maybe just log
	// instead of return error?
	if err != nil {
		return nil
	}

	record := &Record{
		Data: recordData,
	}

	if metadata != nil {
		record.Expiration = AerospikeEpoch.Add(time.Second * time.Duration(metadata.GetExpiration()))
		record.Generation = metadata.GetGeneration()
	}

	return record
}

type Neighbor struct {
	Namespace string
	Set       *string
	Key       any
	Distance  float32
	Record    *Record
}

func newNeighborFromProto(n *protos.Neighbor) (*Neighbor, error) {
	namespace, set, key, err := protos.ConvertFromKey(n.GetKey())
	if err != nil {
		return nil, fmt.Errorf("error converting neighbor: %w", err)
	}

	return &Neighbor{
		Namespace: namespace,
		Set:       set,
		Key:       key,
		Distance:  n.GetDistance(),
		Record:    newRecordFromProto(n.GetRecord()),
	}, nil
}

type HnswSearchParams struct {
	*protos.HnswSearchParams
}

func newHnswSearchParamsFromProto(params *protos.HnswSearchParams) *HnswSearchParams {
	return &HnswSearchParams{HnswSearchParams: params}
}
