package avs

import (
	"testing"

	"github.com/aerospike/avs-client-go/protos"
)

func TestVersionLTGT(t *testing.T) {
	testCases := []struct {
		name   string
		v1     version
		v2     version
		wantLT bool
		wantGT bool
	}{
		{
			name:   "v1 is less than v2",
			v1:     newVersion("1.0.0"),
			v2:     newVersion("2.0.0"),
			wantLT: true,
			wantGT: false,
		},
		{
			name:   "v1 is greater than v2",
			v1:     newVersion("2.0.0"),
			v2:     newVersion("1.0.0"),
			wantLT: false,
			wantGT: true,
		},
		{
			name:   "v1 is less than v2",
			v1:     newVersion("1.0.0"),
			v2:     newVersion("1.0.1"),
			wantLT: true,
			wantGT: false,
		},
		{
			name:   "v1 is less than v2",
			v1:     newVersion("1.0.0"),
			v2:     newVersion("1.0.0.dev0"),
			wantLT: true,
			wantGT: false,
		},
		{
			name:   "v1 is greater than v2",
			v1:     newVersion("1.0.0.dev0"),
			v2:     newVersion("1.0.0"),
			wantLT: false,
			wantGT: true,
		},
		{
			name:   "v1 is less than v2",
			v1:     newVersion("1.0.0.dev0"),
			v2:     newVersion("1.0.0.dev1"),
			wantLT: true,
			wantGT: false,
		},
		{
			name:   "v1 is equal to v2",
			v1:     newVersion("1.0.0"),
			v2:     newVersion("1.0.0"),
			wantLT: false,
			wantGT: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.v1.lt(tc.v2)
			if got != tc.wantLT {
				t.Errorf("expected %v, got %v", tc.wantLT, got)
			}

			got = tc.v1.gt(tc.v2)
			if got != tc.wantGT {
				t.Errorf("expected %v, got %v", !tc.wantLT, got)
			}
		})
	}
}

func TestEndpointListEqual(t *testing.T) {
	testCases := []struct {
		name       string
		endpoints1 *protos.ServerEndpointList
		endpoints2 *protos.ServerEndpointList
		want       bool
	}{
		{
			name: "equal endpoints",
			endpoints1: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "localhost",
						Port:    8080,
						IsTls:   false,
					},
					{
						Address: "127.0.0.1",
						Port:    9090,
						IsTls:   true,
					},
				},
			},
			endpoints2: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "localhost",
						Port:    8080,
						IsTls:   false,
					},
					{
						Address: "127.0.0.1",
						Port:    9090,
						IsTls:   true,
					},
				},
			},
			want: true,
		},
		{
			name: "different endpoints",
			endpoints1: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "localhost",
						Port:    8080,
						IsTls:   false,
					},
					{
						Address: "127.0.0.1",
						Port:    9090,
						IsTls:   true,
					},
				},
			},
			endpoints2: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "localhost",
						Port:    8080,
						IsTls:   false,
					},
					{
						Address: "127.0.0.1",
						Port:    9091,
						IsTls:   true,
					},
				},
			},
			want: false,
		},
		{
			name: "different number of endpoints",
			endpoints1: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "localhost",
						Port:    8080,
						IsTls:   false,
					},
					{
						Address: "127.0.0.1",
						Port:    9090,
						IsTls:   true,
					},
				},
			},
			endpoints2: &protos.ServerEndpointList{
				Endpoints: []*protos.ServerEndpoint{
					{
						Address: "localhost",
						Port:    8080,
						IsTls:   false,
					},
				},
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := endpointListEqual(tc.endpoints1, tc.endpoints2)
			if got != tc.want {
				t.Errorf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
