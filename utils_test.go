package avs

import (
	reflect "reflect"
	"testing"

	"github.com/aerospike/avs-client-go/protos"
)

func TestCreateProjectionSpec(t *testing.T) {
	testCases := []struct {
		name                   string
		includeFields          []string
		excludeFields          []string
		expectedProjectionSpec *protos.ProjectionSpec
	}{
		{
			name:          "include fields",
			includeFields: []string{"field1", "field2"},
			excludeFields: nil,
			expectedProjectionSpec: &protos.ProjectionSpec{
				Include: &protos.ProjectionFilter{
					Type:   protos.ProjectionType_SPECIFIED,
					Fields: []string{"field1", "field2"},
				},
				Exclude: &protos.ProjectionFilter{
					Type: protos.ProjectionType_NONE,
				},
			},
		},
		{
			name:          "exclude fields",
			includeFields: nil,
			excludeFields: []string{"field3", "field4"},
			expectedProjectionSpec: &protos.ProjectionSpec{
				Include: &protos.ProjectionFilter{
					Type: protos.ProjectionType_ALL,
				},
				Exclude: &protos.ProjectionFilter{
					Type:   protos.ProjectionType_SPECIFIED,
					Fields: []string{"field3", "field4"},
				},
			},
		},
		{
			name:          "default fields",
			includeFields: nil,
			excludeFields: nil,
			expectedProjectionSpec: &protos.ProjectionSpec{
				Include: &protos.ProjectionFilter{
					Type: protos.ProjectionType_ALL,
				},
				Exclude: &protos.ProjectionFilter{
					Type: protos.ProjectionType_NONE,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spec := createProjectionSpec(tc.includeFields, tc.excludeFields)

			if !reflect.DeepEqual(spec, tc.expectedProjectionSpec) {
				t.Errorf("expected projection spec %v, got %v", tc.expectedProjectionSpec, spec)
			}
		})
	}
}
func TestVersionString(t *testing.T) {
	testCases := []struct {
		name     string
		v        version
		expected string
	}{
		{
			name:     "valid version",
			v:        newVersion("1.2.3"),
			expected: "1.2.3",
		},
		{
			name:     "valid version with suffix",
			v:        newVersion("1.2.3-dev"),
			expected: "1.2.3-dev",
		},
		{
			name:     "valid version with multiple suffixes",
			v:        newVersion("1.2.3-dev.1"),
			expected: "1.2.3-dev.1",
		},
		{
			name:     "valid version with pre-release",
			v:        newVersion("1.2.3-alpha"),
			expected: "1.2.3-alpha",
		},
		{
			name:     "valid version with pre-release and build metadata",
			v:        newVersion("1.2.3-alpha+build123"),
			expected: "1.2.3-alpha+build123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.v.String()
			if got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

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
						Address: "127.0.0.1",
						Port:    9090,
						IsTls:   true,
					},
					{
						Address: "localhost",
						Port:    8080,
						IsTls:   false,
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
