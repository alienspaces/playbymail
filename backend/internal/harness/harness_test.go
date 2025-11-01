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

	l, s, j, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(l, s, j, cfg, dcfg)
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
	require.NotEmpty(t, h.Data.AdventureGameLocationRecs, "Adventure game location records should be created")
	require.NotEmpty(t, h.Data.AdventureGameLocationLinkRecs, "Adventure game location link records should be created")
	require.NotEmpty(t, h.Data.AdventureGameLocationLinkRequirementRecs, "Adventure game location link requirement records should be created")
	require.NotEmpty(t, h.Data.AdventureGameCharacterRecs, "Adventure game character records should be created")
	require.NotEmpty(t, h.Data.AdventureGameCreatureRecs, "Adventure game creature records should be created")
	require.NotEmpty(t, h.Data.AdventureGameItemRecs, "Adventure game item records should be created")
	require.NotEmpty(t, h.Data.AdventureGameLocationLinkRequirementRecs, "Adventure game location link requirement records should be created")
	require.NotEmpty(t, h.Data.AdventureGameLocationInstanceRecs, "Adventure game location instance records should be created")
	require.NotEmpty(t, h.Data.AdventureGameCharacterInstanceRecs, "Adventure game character instance records should be created")
	require.NotEmpty(t, h.Data.AdventureGameCreatureInstanceRecs, "Adventure game creature instance records should be created")
	require.NotEmpty(t, h.Data.AdventureGameItemInstanceRecs, "Adventure game item instance records should be created")

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
	for ref, id := range h.Data.Refs.AdventureGameLocationRefs {
		rec, err := h.Data.GetAdventureGameLocationRecByID(id)
		require.NoErrorf(t, err, "Adventure game location ref %s should resolve to a record", ref)
		require.NotNil(t, rec, "Adventure game location record for ref %s should not be nil", ref)
	}
	for ref, id := range h.Data.Refs.AdventureGameItemInstanceRefs {
		rec, err := h.Data.GetAdventureGameItemInstanceRecByID(id)
		require.NoErrorf(t, err, "Adventure game item instance ref %s should resolve to a record", ref)
		require.NotNil(t, rec, "Adventure game item instance record for ref %s should not be nil", ref)
	}
	for ref, id := range h.Data.Refs.AdventureGameLocationLinkRequirementRefs {
		rec, err := h.Data.GetAdventureGameLocationLinkRequirementRecByID(id)
		require.NoErrorf(t, err, "Adventure game location link requirement ref %s should resolve to a record", ref)
		require.NotNil(t, rec, "Adventure game location link requirement record for ref %s should not be nil", ref)
	}
}
