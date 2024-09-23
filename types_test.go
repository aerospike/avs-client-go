package avs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostPort_String(t *testing.T) {
	host := "localhost"
	port := 8080
	hp := NewHostPort(host, port)

	expected := "localhost:8080"
	result := hp.String()

	assert.Equal(t, expected, result)
}

func TestHostPortSlice_String(t *testing.T) {
	hps := HostPortSlice{
		NewHostPort("localhost", 8080),
		NewHostPort("example.com", 1234),
	}

	expected := "[localhost:8080, example.com:1234]"
	result := hps.String()

	assert.Equal(t, expected, result)
}
