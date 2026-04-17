package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestDomain_StartGameInstance(t *testing.T) {
	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	th, err := harness.NewTesting(cfg, l, s, j, scanner, harness.DefaultDataConfig())
	require.NoError(t, err, "NewTesting returns without error")

	// Keep transaction open so domain can query the data it creates
	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	m := th.Domain.(*domain.Domain)

	// --- Adventure game ---
	// GameInstanceTwoRef has GameSubscriptionPlayerTwoRef linked (1 player, RequiredPlayerCount=1).
	// GameCharacterTwoRef belongs to AccountUserProPlayerRef (same user as PlayerTwoRef).
	// GameOneRef has 3 locations, 1 creature placement, 2 location objects, 1 item placement.
	adventureInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceTwoRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	// GameInstanceCleanRef has no player subscriptions; used to test the insufficient-players error.
	noPlayersInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceCleanRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	// --- Mecha game ---
	// The default harness creates GameMechaRef with 2 sectors and MechaSquadOneRef (1 mech) for
	// AccountUserStandardRef, but no game instances or subscriptions. We create them here.
	mechaGameRec, err := th.Data.GetGameRecByRef(harness.GameMechaRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	accountUserStdRec, err := th.Data.GetAccountUserRecByRef(harness.AccountUserStandardRef)
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	accountUserContactRec, err := th.Data.GetAccountUserContactRecByAccountUserID(accountUserStdRec.ID)
	require.NoError(t, err, "GetAccountUserContactRecByAccountUserID returns without error")

	// A manager subscription is required before a game instance can be created.
	_, err = m.CreateGameSubscriptionRec(&game_record.GameSubscription{
		GameID:           mechaGameRec.ID,
		AccountID:        accountUserStdRec.AccountID,
		AccountUserID:    accountUserStdRec.ID,
		SubscriptionType: game_record.GameSubscriptionTypeManager,
		Status:           game_record.GameSubscriptionStatusActive,
	})
	require.NoError(t, err, "CreateGameSubscriptionRec (mecha manager) returns without error")

	mechaPlayerSubRec, err := m.CreateGameSubscriptionRec(&game_record.GameSubscription{
		GameID:               mechaGameRec.ID,
		AccountID:            accountUserStdRec.AccountID,
		AccountUserID:        accountUserStdRec.ID,
		AccountUserContactID: nullstring.FromString(accountUserContactRec.ID),
		SubscriptionType:     game_record.GameSubscriptionTypePlayer,
		Status:               game_record.GameSubscriptionStatusActive,
		DeliveryMethod:       nullstring.FromString(game_record.GameSubscriptionDeliveryMethodEmail),
	})
	require.NoError(t, err, "CreateGameSubscriptionRec (mecha player) returns without error")

	mechaInstanceRec, err := m.CreateGameInstanceRec(&game_record.GameInstance{
		GameID:              mechaGameRec.ID,
		Status:              game_record.GameInstanceStatusCreated,
		RequiredPlayerCount: 1,
		DeliveryEmail:       true,
	})
	require.NoError(t, err, "CreateGameInstanceRec (mecha) returns without error")

	_, err = m.CreateGameSubscriptionInstanceRec(&game_record.GameSubscriptionInstance{
		AccountID:          accountUserStdRec.AccountID,
		AccountUserID:      accountUserStdRec.ID,
		GameSubscriptionID: mechaPlayerSubRec.ID,
		GameInstanceID:     mechaInstanceRec.ID,
	})
	require.NoError(t, err, "CreateGameSubscriptionInstanceRec (mecha player) returns without error")

	testCases := []struct {
		name        string
		instanceID  string
		expectError bool
		errContains string

		expectStatus string

		// adventure-game assertions (non-zero when expected)
		expectAdventureLocationCount int
		expectAdventureCharCount     int
		expectAdventureCreatureCount int
		expectAdventureItemCount     int
		expectAdventureObjectCount   int

		// mecha-game assertions (non-zero when expected)
		expectMechaSectorCount int
		expectMechaSquadCount  int
		expectMechaMechCount   int
	}{
		{
			name:         "adventure game starts and returns adventure instance data",
			instanceID:   adventureInstanceRec.ID,
			expectStatus: game_record.GameInstanceStatusStarted,
			// GameOneRef: 3 locations, 1 creature placement, 1 item placement, 2 location objects.
			// PlayerTwoRef → CharacterTwoRef → 1 character instance.
			expectAdventureLocationCount: 3,
			expectAdventureCharCount:     1,
			expectAdventureCreatureCount: 1,
			expectAdventureItemCount:     1,
			expectAdventureObjectCount:   2,
		},
		{
			name:         "mecha game starts and returns mecha instance data",
			instanceID:   mechaInstanceRec.ID,
			expectStatus: game_record.GameInstanceStatusStarted,
			// GameMechaRef: 2 sectors; AccountUserStandardRef has MechaSquadOneRef with 1 mech.
			expectMechaSectorCount: 2,
			expectMechaSquadCount:  1,
			expectMechaMechCount:   1,
		},
		{
			// The adventure instance was already started by the first test case.
			name:        "error when instance is not in created status",
			instanceID:  adventureInstanceRec.ID,
			expectError: true,
			errContains: "'created' status",
		},
		{
			// GameInstanceCleanRef has RequiredPlayerCount=1 but 0 players subscribed.
			name:        "error when insufficient players",
			instanceID:  noPlayersInstanceRec.ID,
			expectError: true,
			errContains: "insufficient players",
		},
		{
			name:        "error for invalid instance ID",
			instanceID:  "invalid-uuid",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			instanceRec, instanceData, err := m.StartGameInstance(tc.instanceID)

			if tc.expectError {
				require.Error(t, err, "Should return error")
				if tc.errContains != "" {
					require.Contains(t, err.Error(), tc.errContains, "Error message contains expected text")
				}
				return
			}

			require.NoError(t, err, "Should not return error")
			require.NotNil(t, instanceRec, "Returned instance record should not be nil")
			require.NotNil(t, instanceData, "Returned instance data should not be nil")
			require.Equal(t, tc.expectStatus, instanceRec.Status, "Instance status equals expected")
			require.True(t, instanceRec.StartedAt.Valid, "StartedAt should be set")
			require.Equal(t, 0, instanceRec.CurrentTurn, "CurrentTurn should be 0")

			adventureExpected := tc.expectAdventureLocationCount > 0 ||
				tc.expectAdventureCharCount > 0

			if adventureExpected {
				require.NotNil(t, instanceData.Adventure, "Adventure data should be non-nil for adventure game")
				require.Nil(t, instanceData.Mecha, "Mecha data should be nil for adventure game")
				require.Len(t, instanceData.Adventure.LocationInstances, tc.expectAdventureLocationCount, "Location instance count equals expected")
				require.Len(t, instanceData.Adventure.CharacterInstances, tc.expectAdventureCharCount, "Character instance count equals expected")
				require.Len(t, instanceData.Adventure.CreatureInstances, tc.expectAdventureCreatureCount, "Creature instance count equals expected")
				require.Len(t, instanceData.Adventure.ItemInstances, tc.expectAdventureItemCount, "Item instance count equals expected")
				require.Len(t, instanceData.Adventure.LocationObjectInstances, tc.expectAdventureObjectCount, "Location object instance count equals expected")
			}

			mechaExpected := tc.expectMechaSectorCount > 0 ||
				tc.expectMechaSquadCount > 0

			if mechaExpected {
				require.NotNil(t, instanceData.Mecha, "Mecha data should be non-nil for mecha game")
				require.Nil(t, instanceData.Adventure, "Adventure data should be nil for mecha game")
				require.Len(t, instanceData.Mecha.SectorInstances, tc.expectMechaSectorCount, "Sector instance count equals expected")
				require.Len(t, instanceData.Mecha.SquadInstances, tc.expectMechaSquadCount, "Squad instance count equals expected")
				require.Len(t, instanceData.Mecha.MechInstances, tc.expectMechaMechCount, "Mech instance count equals expected")
			}
		})
	}
}
