package harness_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestHarnessSetupTeardown_DefaultDataConfig(t *testing.T) {
	// Use the default data config from the harness package
	dcfg := harness.DefaultDataConfig()

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(l, s, j, dcfg)
	require.NoError(t, err, "NewTesting returns without error")

	_, err = h.Setup()
	require.NoError(t, err, "Setup returns without error")
	defer func() {
		err = h.Teardown()
		require.NoError(t, err, "Teardown returns without error")
	}()

	// Check that the default data config created the expected records
	require.NotEmpty(t, h.Data.AccountRecs, "Account records should be created")
	require.NotEmpty(t, h.Data.GameRecs, "Game records should be created")
	require.NotEmpty(t, h.Data.GameLocationRecs, "Game location records should be created")
	require.NotEmpty(t, h.Data.LocationLinkRecs, "Location link records should be created")

	// Check that references are set
	for ref, id := range h.Data.Refs.AccountRefs {
		rec, err := h.Data.GetAccountRecByID(id)
		require.NoErrorf(t, err, "Account ref %s should resolve to a record", ref)
		require.NotNil(t, rec, "Account record for ref %s should not be nil", ref)
	}
	for ref, id := range h.Data.Refs.GameRefs {
		rec, err := h.Data.GetGameRecByID(id)
		require.NoErrorf(t, err, "Game ref %s should resolve to a record", ref)
		require.NotNil(t, rec, "Game record for ref %s should not be nil", ref)
	}
	for ref, id := range h.Data.Refs.GameLocationRefs {
		rec, err := h.Data.GetGameLocationRecByID(id)
		require.NoErrorf(t, err, "Game location ref %s should resolve to a record", ref)
		require.NotNil(t, rec, "Game location record for ref %s should not be nil", ref)
	}
}
