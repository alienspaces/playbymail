package runner

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/store"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func newDefaultDependencies(t *testing.T) (config.Config, *log.Log, *store.Store, *river.Client[pgx.Tx], turn_sheet.TurnSheetScanner) {
	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	// The CLI runner does not need a turn sheet scanner
	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	return cfg, l, s, j, scanner
}

func newTestRunner(t *testing.T, cfg config.Config, l *log.Log, s *store.Store, j *river.Client[pgx.Tx], scanner turn_sheet.TurnSheetScanner) *Runner {

	r, err := NewRunner(cfg, l, j, scanner)
	require.NoError(t, err, "NewRunner returns without error")

	err = r.Init(s)
	require.NoError(t, err, "Init returns without error")

	return r
}

func TestNewRunner(t *testing.T) {
	t.Parallel()

	cfg, l, s, j, scanner := newDefaultDependencies(t)

	r := newTestRunner(t, cfg, l, s, j, scanner)
	require.NotNil(t, r, "newTestRunner returns a new Runner")
}
