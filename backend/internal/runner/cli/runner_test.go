package runner

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/store"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func newDefaultDependencies(t *testing.T) (*log.Log, *store.Store, *river.Client[pgx.Tx]) {
	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	return l, s, j
}

func newTestRunner(t *testing.T, l *log.Log, s *store.Store, j *river.Client[pgx.Tx]) *Runner {

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	r, err := NewRunner(l, j, cfg)
	require.NoError(t, err, "NewRunner returns without error")

	err = r.Init(s)
	require.NoError(t, err, "Init returns without error")

	return r
}

func TestNewRunner(t *testing.T) {
	t.Parallel()

	l, s, j := newDefaultDependencies(t)

	r := newTestRunner(t, l, s, j)
	require.NotNil(t, r, "newTestRunner returns a new Runner")
}
