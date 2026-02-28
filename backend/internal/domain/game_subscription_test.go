package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullint32"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestApproveGameSubscription(t *testing.T) {

	// Create test harness with custom configuration
	dataConfig := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record: &game_record.Game{
					Name:              harness.UniqueName("Test Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
				},
				AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
					{
						Reference: "starting-location",
						Record: &adventure_game_record.AdventureGameLocation{
							Name:               harness.UniqueName("Starting Location"),
							Description:        "A starting location for testing",
							IsStartingLocation: true,
						},
					},
				},
				GameInstanceConfigs: []harness.GameInstanceConfig{
					{
						Reference: "test-instance-1",
					},
					{
						Reference: "test-instance-2",
					},
				},
			},
		},
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "designer-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("designer@example.com"),
					Status: account_record.AccountUserStatusActive,
				},
				GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
					{
						Reference:        "designer-subscription",
						GameRef:          harness.GameOneRef,
						SubscriptionType: game_record.GameSubscriptionTypeDesigner,
						Record: &game_record.GameSubscription{
							Status: game_record.GameSubscriptionStatusActive,
						},
					},
				},
			},
			{
				Reference: "pending-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("pending@example.com"),
					Status: account_record.AccountUserStatusPendingApproval,
				},
				GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
					{
						Reference:        "pending-subscription",
						GameRef:          harness.GameOneRef,
						GameInstanceRefs: []string{"test-instance-1"},
						SubscriptionType: game_record.GameSubscriptionTypePlayer,
						Record: &game_record.GameSubscription{
							Status: game_record.GameSubscriptionStatusPendingApproval,
						},
					},
				},
			},
			{
				Reference: "pending-account-two",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("pending-two@example.com"),
					Status: account_record.AccountUserStatusPendingApproval,
				},
				GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
					{
						Reference:        "active-subscription",
						GameRef:          harness.GameOneRef,
						GameInstanceRefs: []string{"test-instance-2"},
						SubscriptionType: game_record.GameSubscriptionTypePlayer,
						Record: &game_record.GameSubscription{
							Status: game_record.GameSubscriptionStatusActive,
						},
					},
				},
			},
			{
				Reference: "other-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("other@example.com"),
					Status: account_record.AccountUserStatusActive,
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	th, err := harness.NewTesting(cfg, l, s, j, scanner, dataConfig)
	require.NoError(t, err, "NewTesting returns without error")

	// Domain tests use transactions that can be rolled back
	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Teardown returns without error")
	}()

	// Get references from test data
	pendingSubscriptionID, ok := th.Data.Refs.GameSubscriptionRefs["pending-subscription"]
	require.True(t, ok, "pending-subscription reference exists")

	activeSubscriptionID, ok := th.Data.Refs.GameSubscriptionRefs["active-subscription"]
	require.True(t, ok, "active-subscription reference exists")

	pendingAccount, err := th.Data.GetAccountUserRecByRef("pending-account")
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	pendingAccountTwo, err := th.Data.GetAccountUserRecByRef("pending-account-two")
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	testCases := []struct {
		name           string
		subscriptionID string
		email          string
		expectError    bool
		errorCode      coreerror.Code
		expectStatus   string
		validate       func(t *testing.T, rec *game_record.GameSubscription)
	}{
		{
			name:           "successfully approves pending subscription with correct email",
			subscriptionID: pendingSubscriptionID,
			email:          pendingAccount.Email,
			expectError:    false,
			expectStatus:   game_record.GameSubscriptionStatusActive,
			validate: func(t *testing.T, rec *game_record.GameSubscription) {
				require.Equal(t, game_record.GameSubscriptionStatusActive, rec.Status, "Status is active")
				require.Equal(t, pendingAccount.AccountID, rec.AccountID, "Account ID matches")
			},
		},
		{
			name:           "returns error when subscription ID is empty",
			subscriptionID: "",
			email:          pendingAccount.Email,
			expectError:    true,
			errorCode:      coreerror.ErrorCodeInvalidData,
		},
		{
			name:           "returns error when email is empty",
			subscriptionID: pendingSubscriptionID,
			email:          "",
			expectError:    true,
			errorCode:      coreerror.ErrorCodeInvalidData,
		},
		{
			name:           "returns error when subscription does not exist",
			subscriptionID: "00000000-0000-0000-0000-000000000000",
			email:          pendingAccount.Email,
			expectError:    true,
			errorCode:      coreerror.ErrorCodeNotFound,
		},
		{
			name:           "returns error when subscription is not pending approval",
			subscriptionID: activeSubscriptionID,
			email:          pendingAccountTwo.Email,
			expectError:    true,
			errorCode:      coreerror.ErrorCodeInvalidData,
		},
		{
			name:           "returns error when email does not match subscription account",
			subscriptionID: pendingSubscriptionID,
			email:          "wrong-email@example.com",
			expectError:    true,
			errorCode:      coreerror.ErrorCodeInvalidData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := th.Domain.(*domain.Domain)

			rec, err := m.ApproveGameSubscription(tc.subscriptionID, tc.email)

			if tc.expectError {
				require.Error(t, err, "Should return error")
				if tc.errorCode != "" {
					var coreErr coreerror.Error
					require.ErrorAs(t, err, &coreErr, "Error should be a coreerror.Error")
					require.Equal(t, tc.errorCode, coreErr.ErrorCode, "Error code matches expected")
				}
				require.Nil(t, rec, "Record should be nil on error")
			} else {
				require.NoError(t, err, "Should not return error")
				require.NotNil(t, rec, "Record should not be nil")
				require.Equal(t, tc.expectStatus, rec.Status, "Status equals expected")
				if tc.validate != nil {
					tc.validate(t, rec)
				}
			}
		})
	}
}

// Note: Tests for turn sheet token generation/verification have been moved to game_subscription_instance_test.go
// since tokens are now stored on game_subscription_instance rather than game_subscription

func TestGameSubscriptionInstanceLinking(t *testing.T) {
	dataConfig := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record: &game_record.Game{
					Name:              harness.UniqueName("Test Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
				},
				AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
					{
						Reference: "starting-location",
						Record: &adventure_game_record.AdventureGameLocation{
							Name:               harness.UniqueName("Starting Location"),
							Description:        "A starting location for testing",
							IsStartingLocation: true,
						},
					},
				},
				GameInstanceConfigs: []harness.GameInstanceConfig{
					{
						Reference: "game-instance-1",
					},
				},
			},
		},
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "manager-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("manager@example.com"),
					Status: account_record.AccountUserStatusActive,
				},
				GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
					{
						Reference:        "manager-subscription",
						GameRef:          harness.GameOneRef,
						GameInstanceRefs: []string{"game-instance-1"},
						SubscriptionType: game_record.GameSubscriptionTypeManager,
					},
				},
			},
			{
				Reference: "player-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("player@example.com"),
					Status: account_record.AccountUserStatusActive,
				},
				GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
					{
						Reference:        "player-subscription",
						GameRef:          harness.GameOneRef,
						GameInstanceRefs: []string{"game-instance-1"},
						SubscriptionType: game_record.GameSubscriptionTypePlayer,
					},
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	th, err := harness.NewTesting(cfg, l, s, j, scanner, dataConfig)
	require.NoError(t, err, "NewTesting returns without error")

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Teardown returns without error")
	}()

	gameInstanceID, ok := th.Data.Refs.GameInstanceRefs["game-instance-1"]
	require.True(t, ok, "game-instance-1 reference exists")

	managerSubscriptionID, ok := th.Data.Refs.GameSubscriptionRefs["manager-subscription"]
	require.True(t, ok, "manager-subscription reference exists")

	m := th.Domain.(*domain.Domain)

	// Test that subscriptions can be created without instances
	t.Run("create subscription without instance", func(t *testing.T) {
		subscription := &game_record.GameSubscription{
			GameID:           th.Data.GameRecs[0].ID,
			AccountID:        th.Data.AccountRecs[0].ID,
			SubscriptionType: game_record.GameSubscriptionTypeManager,
			Status:           game_record.GameSubscriptionStatusActive,
		}
		_, err := m.CreateGameSubscriptionRec(subscription)
		require.NoError(t, err, "Should create subscription without instance")
	})

	// Test instance linking
	t.Run("link instance to subscription", func(t *testing.T) {
		managerSub, err := m.GetGameSubscriptionRec(managerSubscriptionID, nil)
		require.NoError(t, err)

		// Create subscription-instance link (account_id will be derived from subscription in validation)
		instanceLinkRec := &game_record.GameSubscriptionInstance{
			GameSubscriptionID: managerSub.ID,
			GameInstanceID:     gameInstanceID,
		}
		_, err = m.CreateGameSubscriptionInstanceRec(instanceLinkRec)
		require.NoError(t, err, "Should link instance to subscription")

		// Verify link exists
		instanceLinks, err := m.GetGameSubscriptionInstanceRecsBySubscription(managerSub.ID)
		require.NoError(t, err)
		require.Len(t, instanceLinks, 1, "Should have one instance link")
		require.Equal(t, gameInstanceID, instanceLinks[0].GameInstanceID)
	})

	// Test instance limit validation
	t.Run("instance limit validation", func(t *testing.T) {
		// Create subscription with limit of 1
		limit := int32(1)
		subscription := &game_record.GameSubscription{
			GameID:           th.Data.GameRecs[0].ID,
			AccountID:        th.Data.AccountRecs[0].ID,
			SubscriptionType: game_record.GameSubscriptionTypeManager,
			Status:           game_record.GameSubscriptionStatusActive,
			InstanceLimit:    nullint32.FromInt32(limit),
		}
		subscription, err := m.CreateGameSubscriptionRec(subscription)
		require.NoError(t, err)

		// Link first instance (should succeed)
		instanceLinkRec1 := &game_record.GameSubscriptionInstance{
			GameSubscriptionID: subscription.ID,
			GameInstanceID:     gameInstanceID,
		}
		_, err = m.CreateGameSubscriptionInstanceRec(instanceLinkRec1)
		require.NoError(t, err)

		// Try to link second instance (should fail due to limit)
		// Create another instance first
		gameInstance2 := &game_record.GameInstance{
			GameID: th.Data.GameRecs[0].ID,
		}
		gameInstance2, err = m.CreateGameInstanceRec(gameInstance2)
		require.NoError(t, err)

		instanceLinkRec2 := &game_record.GameSubscriptionInstance{
			GameSubscriptionID: subscription.ID,
			GameInstanceID:     gameInstance2.ID,
		}
		_, err = m.CreateGameSubscriptionInstanceRec(instanceLinkRec2)
		require.Error(t, err, "Should fail when limit is reached")
		require.Contains(t, err.Error(), "instance limit reached")
	})
}
