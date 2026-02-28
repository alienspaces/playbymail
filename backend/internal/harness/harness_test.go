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

	// The harness does not need a turn sheet scanner
	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(cfg, l, s, j, scanner, dcfg)
	require.NoError(t, err, "NewTesting returns without error")

	_, err = h.Setup()
	require.NoError(t, err, "Setup returns without error")
	defer func() {
		err = h.Teardown()
		require.NoError(t, err, "Teardown returns without error")
	}()

	// Check that the default data config created the expected records with exact counts
	// Accounts: 4 (one per account user, 1:1 mapping)
	require.Len(t, h.Data.AccountRecs, 4, "Should have exactly 4 account records")
	// AccountUsers: 4 (StandardAccountRef, ProPlayerAccountRef, ProDesignerAccountRef, ProManagerAccountRef)
	require.Len(t, h.Data.AccountUserRecs, 4, "Should have exactly 4 account user records")
	// AccountUserContacts: 4 (one per account user)
	require.Len(t, h.Data.AccountUserContactRecs, 4, "Should have exactly 4 account user contact records")
	// Games: 2 (GameOneRef, GameDraftRef)
	require.Len(t, h.Data.GameRecs, 2, "Should have exactly 2 game records")
	// Locations: 2 (GameLocationOneRef, GameLocationTwoRef)
	require.Len(t, h.Data.AdventureGameLocationRecs, 2, "Should have exactly 2 adventure game location records")
	// Location links: 1 (GameLocationLinkOneRef)
	require.Len(t, h.Data.AdventureGameLocationLinkRecs, 1, "Should have exactly 1 adventure game location link record")
	// Location link requirements: 1 (GameLocationLinkRequirementOneRef)
	require.Len(t, h.Data.AdventureGameLocationLinkRequirementRecs, 1, "Should have exactly 1 adventure game location link requirement record")
	// Characters: 2 (GameCharacterOneRef, GameCharacterTwoRef)
	require.Len(t, h.Data.AdventureGameCharacterRecs, 2, "Should have exactly 2 adventure game character records")
	// Creatures: 2 (GameCreatureOneRef, GameCreatureTwoRef)
	require.Len(t, h.Data.AdventureGameCreatureRecs, 2, "Should have exactly 2 adventure game creature records")
	// Items: 2 (GameItemOneRef, GameItemTwoRef)
	require.Len(t, h.Data.AdventureGameItemRecs, 2, "Should have exactly 2 adventure game item records")
	// Location instances: 2 locations * 2 game instances = 4 (auto-generated for each game instance)
	require.Len(t, h.Data.AdventureGameLocationInstanceRecs, 4, "Should have exactly 4 adventure game location instance records (2 locations * 2 game instances)")
	// Character instances: 1 (GameCharacterInstanceOneRef) - only for GameInstanceOneRef
	require.Len(t, h.Data.AdventureGameCharacterInstanceRecs, 1, "Should have exactly 1 adventure game character instance record")
	// Creature instances: 1 (GameCreatureInstanceOneRef) - only for GameInstanceOneRef
	require.Len(t, h.Data.AdventureGameCreatureInstanceRecs, 1, "Should have exactly 1 adventure game creature instance record")
	// Item instances: 1 (GameItemInstanceOneRef) - only for GameInstanceOneRef
	require.Len(t, h.Data.AdventureGameItemInstanceRecs, 1, "Should have exactly 1 adventure game item instance record")

	// Check that references are set
	for ref, id := range h.Data.Refs.AccountUserRefs {
		rec, err := h.Data.GetAccountUserRecByID(id)
		require.NoErrorf(t, err, "Account user ref %s should resolve to a record", ref)
		require.NotNil(t, rec, "Account user record for ref %s should not be nil", ref)
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
