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

// Record represents a record in the database.
type Record struct {
	Data       map[string]any // The data of the record. This includes the vector
	Expiration *time.Time     // The expiration time of the record.
	Generation uint32         // The generation of the record.
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
		if metadata.Expiration != 0 {
			exp := AerospikeEpoch.Add(time.Second * time.Duration(metadata.GetExpiration()))
			record.Expiration = &exp
		}

		record.Generation = metadata.GetGeneration()
	}

	return record
}

// Neighbor represents a record that is a neighbor to a query vector.
type Neighbor struct {
	Record    *Record // The record data.
	Set       *string // The set within the namespace where the record resides.
	Key       any     // The key of the record.
	Namespace string  // The namespace of the record.
	Distance  float32 // The distance from the query vector.
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
