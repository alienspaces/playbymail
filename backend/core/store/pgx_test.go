package store

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/config"
)

func newTestStore(t *testing.T) *Store {
	cfg := config.Config{}
	err := config.Parse(&cfg)
	require.NoError(t, err, "config.Parse returns without error")

	s, err := NewStore(cfg)
	require.NoError(t, err, "NewStore returns without error")
	require.NotNil(t, s, "NewStore returns a store")

	return s
}

func TestConnectPgx(t *testing.T) {
	tests := map[string]struct {
		values  map[string]string
		wantErr bool
	}{
		"with all config": {
			values: map[string]string{
				"DATABASE_MAX_OPEN_CONNECTIONS": "1",
			},
			wantErr: false,
		},
	}

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		t.Run(tcName, func(t *testing.T) {
			s := newTestStore(t)

			c, err := connectPgx(s.log, s.config)
			if tc.wantErr == false {
				defer func() {
					err = s.ClosePool()
					require.NoError(t, err, "ClosePool returns without error")
				}()
			}

			if tc.wantErr {
				require.Error(t, err, "connectPgx returns with error")
				return
			}

			require.NoError(t, err, "connectPgx returns without error")
			require.NotNil(t, c, "connectPgx returns a pgx pool")
		})
	}
}
