package avs

import "testing"

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
