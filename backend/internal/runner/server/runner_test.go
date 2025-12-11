package runner_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestNewRunner(t *testing.T) {

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	r, err := testutil.NewTestRunner(cfg, l, s, j, scanner)
	require.NoError(t, err, "NewTestRunner returns without error")

	err = r.Init(s)
	require.NoError(t, err, "Init returns without error")
}
