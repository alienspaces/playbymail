package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestDomain_GetPlayerCountForGameInstance(t *testing.T) {

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

	instanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	testCases := []struct {
		name          string
		instanceID    string
		expectedCount int
		expectError   bool
	}{
		{
			// Default harness links 3 subscriptions to GameInstanceOneRef:
			// GameSubscriptionPlayerThreeRef (StandardAccount),
			// GameSubscriptionPlayerOneRef (ProPlayerAccount),
			// GameSubscriptionManagerOneRef (ProManagerAccount)
			name:          "returns player count for valid instance",
			instanceID:    instanceRec.ID,
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:        "returns error for invalid UUID",
			instanceID:  "invalid-uuid",
			expectError: true,
		},
		{
			// A valid UUID with no linked subscriptions returns 0, not an error
			name:          "returns zero for non-existent instance",
			instanceID:    "00000000-0000-0000-0000-000000000000",
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := th.Domain.(*domain.Domain)

			count, err := m.GetPlayerCountForGameInstance(tc.instanceID)

			if tc.expectError {
				require.Error(t, err, "Should return error")
			} else {
				require.NoError(t, err, "Should not return error")
				require.Equal(t, tc.expectedCount, count, "Player count equals expected")
			}
		})
	}
}

func TestDomain_HasAvailableCapacity(t *testing.T) {

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	th, err := harness.NewTesting(cfg, l, s, j, scanner, harness.DefaultDataConfig())
	require.NoError(t, err, "NewTesting returns without error")

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	m := th.Domain.(*domain.Domain)

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	// Create instance with capacity for 2 players (no players assigned yet)
	instanceWithCapacity, err := m.CreateGameInstanceRec(&game_record.GameInstance{
		GameID:              gameRec.ID,
		Status:              game_record.GameInstanceStatusCreated,
		RequiredPlayerCount: 2,
		DeliveryEmail:       true,
	})
	require.NoError(t, err, "CreateGameInstanceRec returns without error")

	// Create instance with RequiredPlayerCount = 0 (unlimited capacity)
	instanceUnlimited, err := m.CreateGameInstanceRec(&game_record.GameInstance{
		GameID:              gameRec.ID,
		Status:              game_record.GameInstanceStatusCreated,
		RequiredPlayerCount: 0,
		DeliveryEmail:       true,
	})
	require.NoError(t, err, "CreateGameInstanceRec returns without error")

	testCases := []struct {
		name        string
		instanceID  string
		expectValue bool
		expectError bool
	}{
		{
			name:        "returns true when instance has available capacity",
			instanceID:  instanceWithCapacity.ID,
			expectValue: true,
			expectError: false,
		},
		{
			name:        "returns true when instance has no capacity limit",
			instanceID:  instanceUnlimited.ID,
			expectValue: true,
			expectError: false,
		},
		{
			name:        "returns error for invalid UUID",
			instanceID:  "invalid-uuid",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hasCapacity, err := m.HasAvailableCapacity(tc.instanceID)

			if tc.expectError {
				require.Error(t, err, "Should return error")
			} else {
				require.NoError(t, err, "Should not return error")
				require.Equal(t, tc.expectValue, hasCapacity, "HasAvailableCapacity equals expected")
			}
		})
	}
}

func TestDomain_FindAvailableGameInstance(t *testing.T) {

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	th, err := harness.NewTesting(cfg, l, s, j, scanner, harness.DefaultDataConfig())
	require.NoError(t, err, "NewTesting returns without error")

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	m := th.Domain.(*domain.Domain)

	managerSubscriptionRec, err := th.Data.GetGameSubscriptionRecByRef(harness.GameSubscriptionManagerOneRef)
	require.NoError(t, err, "GetGameSubscriptionRecByRef returns without error")

	testCases := []struct {
		name           string
		subscriptionID string
		expectInstance bool
		expectError    bool
	}{
		{
			// The default harness data includes game instances with status "created" for GameOneRef
			name:           "returns available instance for valid subscription",
			subscriptionID: managerSubscriptionRec.ID,
			expectInstance: true,
			expectError:    false,
		},
		{
			name:           "returns error for invalid subscription UUID",
			subscriptionID: "invalid-uuid",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			instance, err := m.FindAvailableGameInstance(tc.subscriptionID)

			if tc.expectError {
				require.Error(t, err, "Should return error")
			} else {
				require.NoError(t, err, "Should not return error")
				if tc.expectInstance {
					require.NotNil(t, instance, "Instance should not be nil")
					require.NotEmpty(t, instance.ID, "Instance should have an ID")
				} else {
					require.Nil(t, instance, "Instance should be nil")
				}
			}
		})
	}
}

func TestDomain_AssignPlayerToGameInstance(t *testing.T) {

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	th, err := harness.NewTesting(cfg, l, s, j, scanner, harness.DefaultDataConfig())
	require.NoError(t, err, "NewTesting returns without error")

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	m := th.Domain.(*domain.Domain)

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	playerSubscriptionRec, err := th.Data.GetGameSubscriptionRecByRef(harness.GameSubscriptionPlayerOneRef)
	require.NoError(t, err, "GetGameSubscriptionRecByRef returns without error")

	// Create a fresh instance with RequiredPlayerCount = 2 (one slot used by filler below)
	instanceWithCapacity, err := m.CreateGameInstanceRec(&game_record.GameInstance{
		GameID:              gameRec.ID,
		Status:              game_record.GameInstanceStatusCreated,
		RequiredPlayerCount: 2,
		DeliveryEmail:       true,
	})
	require.NoError(t, err, "CreateGameInstanceRec returns without error")

	// Create a fresh instance at capacity = 1 and immediately fill it
	instanceFull, err := m.CreateGameInstanceRec(&game_record.GameInstance{
		GameID:              gameRec.ID,
		Status:              game_record.GameInstanceStatusCreated,
		RequiredPlayerCount: 1,
		DeliveryEmail:       true,
	})
	require.NoError(t, err, "CreateGameInstanceRec returns without error")

	_, err = m.CreateGameSubscriptionInstanceRec(&game_record.GameSubscriptionInstance{
		AccountID:          playerSubscriptionRec.AccountID,
		GameSubscriptionID: playerSubscriptionRec.ID,
		GameInstanceID:     instanceFull.ID,
	})
	require.NoError(t, err, "CreateGameSubscriptionInstanceRec returns without error")

	testCases := []struct {
		name           string
		subscriptionID string
		instanceID     string
		expectError    bool
		errContains    string
	}{
		{
			name:           "successfully assigns player to instance with capacity",
			subscriptionID: playerSubscriptionRec.ID,
			instanceID:     instanceWithCapacity.ID,
			expectError:    false,
		},
		{
			name:           "returns error when instance is at capacity",
			subscriptionID: playerSubscriptionRec.ID,
			instanceID:     instanceFull.ID,
			expectError:    true,
			errContains:    "no available capacity",
		},
		{
			name:           "returns error for invalid subscription UUID",
			subscriptionID: "invalid-uuid",
			instanceID:     instanceWithCapacity.ID,
			expectError:    true,
		},
		{
			name:           "returns error for invalid instance UUID",
			subscriptionID: playerSubscriptionRec.ID,
			instanceID:     "invalid-uuid",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gsi, err := m.AssignPlayerToGameInstance(tc.subscriptionID, tc.instanceID)

			if tc.expectError {
				require.Error(t, err, "Should return error")
				if tc.errContains != "" {
					require.Contains(t, err.Error(), tc.errContains, "Error message contains expected text")
				}
			} else {
				require.NoError(t, err, "Should not return error")
				require.NotNil(t, gsi, "Returned record should not be nil")
				require.Equal(t, tc.subscriptionID, gsi.GameSubscriptionID, "GameSubscriptionID equals expected")
				require.Equal(t, tc.instanceID, gsi.GameInstanceID, "GameInstanceID equals expected")
			}
		})
	}
}
