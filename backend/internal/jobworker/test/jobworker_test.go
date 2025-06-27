package jobworker

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"

	"gitlab.com/alienspaces/playbymail/internal/harness"
)

func newTestHarness(t *testing.T, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], cfg *harness.DataConfig) *harness.Testing {
	var dataConfig harness.DataConfig

	if cfg == nil {
		dataConfig = harness.DataConfig{}
	} else {
		dataConfig = *cfg
	}

	// For all job worker tests the test harness should not commit so that
	// the test harness db transaction can be shared with the worker being
	// tested ensuring test data is isolated to the test harness and worker.
	h, err := harness.NewTesting(l, s, j, dataConfig)
	require.NoError(t, err, "NewTesting returns without error")
	require.NotNil(t, h, "NewTesting returns a new test harness")

	return h
}
