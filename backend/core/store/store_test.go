package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {

	tests := map[string]struct {
		values  map[string]string
		setup   func(*Store)
		wantErr bool
	}{
		"existing pool": {
			values: map[string]string{
				"DATABASE_MAX_OPEN_CONNECTIONS": "1",
			},
			setup: func(s *Store) {
				s.Pool()
			},
			wantErr: false,
		},
		"new pool": {
			values: map[string]string{
				"DATABASE_MAX_OPEN_CONNECTIONS": "1",
			},
			setup: func(s *Store) {
				s.pgxPool = nil
			},
			wantErr: false,
		},
	}

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		t.Run(tcName, func(t *testing.T) {
			s := newTestStore(t)

			tc.setup(s)

			pool, err := s.Pool()
			defer func() {
				// Don't close a pool that failed to create
				if s.pgxPool != nil {
					err := s.ClosePool()
					require.NoError(t, err, "ClosePool returns without error")
				}
			}()

			if tc.wantErr {
				require.Error(t, err, "Pool returns with error")
				return
			}
			require.NoError(t, err, "Pool returns without error")
			require.NotNil(t, pool, "Pool returns a pool")
		})
	}
}

func TestBeginTx(t *testing.T) {
	tests := map[string]struct {
		setup   func(*Store) func()
		wantErr bool
	}{
		"existing db connection": {},
		"no db connection": {
			setup: func(s *Store) func() {
				oldPool := s.pgxPool
				s.pgxPool = nil
				return func() {
					s.pgxPool = oldPool
				}
			},
			wantErr: false,
		},
	}

	s := newTestStore(t)

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		func() {
			if tc.setup != nil {
				teardown := tc.setup(s)
				defer func() {
					teardown()
				}()
			}

			tx, err := s.BeginTx()
			if tc.wantErr {
				require.Error(t, err, "BeginTx returns with error")
				return
			}

			require.NoError(t, err, "BeginTx returns without error")
			require.NotNil(t, tx, "BeginTx returns tx struct")
		}()
	}
}
