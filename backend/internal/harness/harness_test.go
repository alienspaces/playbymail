package harness_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
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
	// Games: 3 (GameOneRef, GameDraftRef, GameMechaGameRef)
	require.Len(t, h.Data.GameRecs, 3, "Should have exactly 3 game records")

	// Adventure game specific resources

	// Locations: 3 (GameLocationOneRef, GameLocationTwoRef, GameLocationThreeRef)
	require.Len(t, h.Data.AdventureGameLocationRecs, 3, "Should have exactly 3 adventure game location records")
	// Location links: 2 (GameLocationLinkOneRef, GameLocationLinkTwoRef return path)
	require.Len(t, h.Data.AdventureGameLocationLinkRecs, 2, "Should have exactly 2 adventure game location link records")
	// Location link requirements: 2 (GameLocationLinkRequirementOneRef item, GameLocationLinkRequirementTwoRef creature)
	require.Len(t, h.Data.AdventureGameLocationLinkRequirementRecs, 2, "Should have exactly 2 adventure game location link requirement records")
	// Characters: 2 (GameCharacterOneRef, GameCharacterTwoRef)
	require.Len(t, h.Data.AdventureGameCharacterRecs, 2, "Should have exactly 2 adventure game character records")
	// Creatures: 2 (GameCreatureOneRef, GameCreatureTwoRef)
	require.Len(t, h.Data.AdventureGameCreatureRecs, 2, "Should have exactly 2 adventure game creature records")
	// Items: 2 (GameItemOneRef, GameItemTwoRef)
	require.Len(t, h.Data.AdventureGameItemRecs, 2, "Should have exactly 2 adventure game item records")
	// Location instances: 3 (one per location for GameInstanceOneRef; GameInstanceCleanRef stays in created status)
	require.Len(t, h.Data.AdventureGameLocationInstanceRecs, 3, "Should have exactly 3 adventure game location instance records (1 per location for started instance)")
	// Character instances: 1 (GameCharacterInstanceOneRef) - only for GameInstanceOneRef
	require.Len(t, h.Data.AdventureGameCharacterInstanceRecs, 1, "Should have exactly 1 adventure game character instance record")
	// Creature instances: 1 (GameCreatureInstanceOneRef) - only for GameInstanceOneRef
	require.Len(t, h.Data.AdventureGameCreatureInstanceRecs, 1, "Should have exactly 1 adventure game creature instance record")
	// Item instances: 1 (GameItemInstanceOneRef) - only for GameInstanceOneRef
	require.Len(t, h.Data.AdventureGameItemInstanceRecs, 1, "Should have exactly 1 adventure game item instance record")

	// MechaGame game specific resources

	// Chassis: 1 (MechaGameChassisOneRef)
	require.Len(t, h.Data.MechaGameChassisRecs, 1, "Should have exactly 1 mecha chassis record")
	// Weapons: 1 (MechaGameWeaponOneRef)
	require.Len(t, h.Data.MechaGameWeaponRecs, 1, "Should have exactly 1 mecha weapon record")
	// Sectors: 2 (MechaGameSectorOneRef, MechaGameSectorTwoRef)
	require.Len(t, h.Data.MechaGameSectorRecs, 2, "Should have exactly 2 mecha sector records")
	// Sector links: 1 (MechaGameSectorLinkOneRef)
	require.Len(t, h.Data.MechaGameSectorLinkRecs, 1, "Should have exactly 1 mecha sector link record")
	// Squads: 2 (MechaGameSquadStarterRef, MechaGameSquadOneRef)
	require.Len(t, h.Data.MechaGameSquadRecs, 2, "Should have exactly 2 mecha squad records")
	// Squad mechs: 2
	require.Len(t, h.Data.MechaGameSquadMechRecs, 2, "Should have exactly 2 mecha squad mech records")

	// All harness account users should be active by default
	for _, rec := range h.Data.AccountUserRecs {
		require.Equalf(t, account_record.AccountUserStatusActive, rec.Status, "Account user %s should have active status", rec.Email)
	}

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
