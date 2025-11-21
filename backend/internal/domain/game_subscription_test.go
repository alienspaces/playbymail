package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
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
				GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
					{
						Reference:        "pending-subscription",
						AccountRef:       "pending-account",
						SubscriptionType: game_record.GameSubscriptionTypePlayer,
						Record: &game_record.GameSubscription{
							Status: game_record.GameSubscriptionStatusPendingApproval,
						},
					},
					{
						Reference:        "active-subscription",
						AccountRef:       "pending-account-two",
						SubscriptionType: game_record.GameSubscriptionTypePlayer,
						Record: &game_record.GameSubscription{
							Status: game_record.GameSubscriptionStatusActive,
						},
					},
				},
			},
		},
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "pending-account",
				Record: &account_record.Account{
					Email:  harness.UniqueEmail("pending@example.com"),
					Status: account_record.AccountStatusPendingApproval,
				},
			},
			{
				Reference: "pending-account-two",
				Record: &account_record.Account{
					Email:  harness.UniqueEmail("pending-two@example.com"),
					Status: account_record.AccountStatusPendingApproval,
				},
			},
			{
				Reference: "other-account",
				Record: &account_record.Account{
					Email:  harness.UniqueEmail("other@example.com"),
					Status: account_record.AccountStatusActive,
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	th, err := harness.NewTesting(l, s, j, cfg, dataConfig)
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

	pendingAccount, err := th.Data.GetAccountRecByRef("pending-account")
	require.NoError(t, err, "GetAccountRecByRef returns without error")

	pendingAccountTwo, err := th.Data.GetAccountRecByRef("pending-account-two")
	require.NoError(t, err, "GetAccountRecByRef returns without error")

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
				require.Equal(t, pendingAccount.ID, rec.AccountID, "Account ID matches")
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
