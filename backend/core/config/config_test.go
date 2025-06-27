package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// setupEnv sets environment variables for testing and returns a teardown function to unset them after the test.
func setupEnv(t *testing.T, envs map[string]string) func() {
	t.Helper()
	for k, v := range envs {
		require.NoError(t, os.Setenv(k, v))
	}
	return func() {
		for k := range envs {
			_ = os.Unsetenv(k)
		}
	}
}

type TestConfig struct {
	StrField   string   `env:"TEST_STR" envDefault:"defaultstr"`
	IntField   int      `env:"TEST_INT" envDefault:"42"`
	BoolField  bool     `env:"TEST_BOOL" envDefault:"true"`
	SliceField []string `env:"TEST_SLICE" envDefault:"a,b,c"`
	ReqField   string   `env:"TEST_REQ,required"`
}

func TestGenericParse(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		want    TestConfig
		wantErr bool
	}{
		{
			name: "custom env vars",
			envs: map[string]string{
				"TEST_STR":   "customstr",
				"TEST_INT":   "99",
				"TEST_BOOL":  "false",
				"TEST_SLICE": "x,y,z",
				"TEST_REQ":   "mustset",
			},
			want: TestConfig{
				StrField:   "customstr",
				IntField:   99,
				BoolField:  false,
				SliceField: []string{"x", "y", "z"},
				ReqField:   "mustset",
			},
			wantErr: false,
		},
		{
			name: "defaults",
			envs: map[string]string{"TEST_REQ": "defaultreq"},
			want: TestConfig{
				StrField:   "defaultstr",
				IntField:   42,
				BoolField:  true,
				SliceField: []string{"a", "b", "c"},
				ReqField:   "defaultreq",
			},
			wantErr: false,
		},
		{
			name:    "missing required field",
			envs:    map[string]string{},
			want:    TestConfig{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teardown := setupEnv(t, tc.envs)
			defer teardown()
			var got TestConfig
			err := Parse(&got)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}
