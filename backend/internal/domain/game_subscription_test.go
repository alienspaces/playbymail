package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
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

func TestGenerateTurnSheetKey(t *testing.T) {
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
						Reference:        "active-subscription",
						AccountRef:       "test-account",
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
				Reference: "test-account",
				Record: &account_record.Account{
					Email:  harness.UniqueEmail("test@example.com"),
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

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Teardown returns without error")
	}()

	subscriptionID, ok := th.Data.Refs.GameSubscriptionRefs["active-subscription"]
	require.True(t, ok, "active-subscription reference exists")

	testCases := []struct {
		name           string
		subscriptionID string
		expectError    bool
		errorCode      coreerror.Code
		validate       func(t *testing.T, key string, rec *game_record.GameSubscription)
	}{
		{
			name:           "successfully generates turn sheet key",
			subscriptionID: subscriptionID,
			expectError:    false,
			validate: func(t *testing.T, key string, rec *game_record.GameSubscription) {
				require.NotEmpty(t, key, "Key should not be empty")
				require.True(t, nulltime.IsValid(rec.TurnSheetKeyExpiresAt), "Expiration should be set")
				expirationTime := nulltime.ToTime(rec.TurnSheetKeyExpiresAt)
				expectedExpiration := time.Now().Add(3 * 24 * time.Hour)
				// Allow 1 minute tolerance for test execution time
				require.WithinDuration(t, expectedExpiration, expirationTime, 1*time.Minute, "Expiration should be approximately 3 days from now")
			},
		},
		{
			name:           "returns error when subscription does not exist",
			subscriptionID: "00000000-0000-0000-0000-000000000000",
			expectError:    true,
			errorCode:      coreerror.ErrorCodeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := th.Domain.(*domain.Domain)

			key, err := m.GenerateTurnSheetKey(tc.subscriptionID)

			if tc.expectError {
				require.Error(t, err, "Should return error")
				if tc.errorCode != "" {
					var coreErr coreerror.Error
					require.ErrorAs(t, err, &coreErr, "Error should be a coreerror.Error")
					require.Equal(t, tc.errorCode, coreErr.ErrorCode, "Error code matches expected")
				}
				require.Empty(t, key, "Key should be empty on error")
			} else {
				require.NoError(t, err, "Should not return error")
				require.NotEmpty(t, key, "Key should not be empty")

				// Verify the key was saved
				rec, err := m.GetGameSubscriptionRec(tc.subscriptionID, nil)
				require.NoError(t, err, "Should be able to retrieve subscription")
				if tc.validate != nil {
					tc.validate(t, key, rec)
				}
			}
		})
	}
}

func TestVerifyTurnSheetKey(t *testing.T) {
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
						Reference:        "active-subscription",
						AccountRef:       "test-account",
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
				Reference: "test-account",
				Record: &account_record.Account{
					Email:  harness.UniqueEmail("test@example.com"),
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

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Teardown returns without error")
	}()

	subscriptionID, ok := th.Data.Refs.GameSubscriptionRefs["active-subscription"]
	require.True(t, ok, "active-subscription reference exists")

	m := th.Domain.(*domain.Domain)

	// Generate a valid key
	validKey, err := m.GenerateTurnSheetKey(subscriptionID)
	require.NoError(t, err, "GenerateTurnSheetKey returns without error")

	testCases := []struct {
		name         string
		turnSheetKey string
		expectError  bool
		validate     func(t *testing.T, rec *game_record.GameSubscription)
	}{
		{
			name:         "successfully validates valid turn sheet key",
			turnSheetKey: validKey,
			expectError:  false,
			validate: func(t *testing.T, rec *game_record.GameSubscription) {
				require.Equal(t, subscriptionID, rec.ID, "Subscription ID matches")
			},
		},
		{
			name:         "returns error when turn sheet key is empty",
			turnSheetKey: "",
			expectError:  true,
		},
		{
			name:         "returns error when turn sheet key does not exist",
			turnSheetKey: "00000000-0000-0000-0000-000000000000",
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, err := m.VerifyTurnSheetKey(tc.turnSheetKey)

			if tc.expectError {
				require.Error(t, err, "Should return error")
				require.Nil(t, rec, "Record should be nil on error")
			} else {
				require.NoError(t, err, "Should not return error")
				require.NotNil(t, rec, "Record should not be nil")
				if tc.validate != nil {
					tc.validate(t, rec)
				}
			}
		})
	}
}

func TestGenerateTurnSheetKeyInvalidatesOldKey(t *testing.T) {
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
						Reference:        "active-subscription",
						AccountRef:       "test-account",
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
				Reference: "test-account",
				Record: &account_record.Account{
					Email:  harness.UniqueEmail("test@example.com"),
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

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Teardown returns without error")
	}()

	subscriptionID, ok := th.Data.Refs.GameSubscriptionRefs["active-subscription"]
	require.True(t, ok, "active-subscription reference exists")

	m := th.Domain.(*domain.Domain)

	// Generate initial key
	initialKey, err := m.GenerateTurnSheetKey(subscriptionID)
	require.NoError(t, err, "GenerateTurnSheetKey returns without error")
	require.NotEmpty(t, initialKey, "Initial key should not be empty")

	// Wait a moment to ensure different timestamp
	time.Sleep(100 * time.Millisecond)

	// Generate a new key (this invalidates the old one)
	newKey, err := m.GenerateTurnSheetKey(subscriptionID)
	require.NoError(t, err, "GenerateTurnSheetKey returns without error")
	require.NotEmpty(t, newKey, "New key should not be empty")
	require.NotEqual(t, initialKey, newKey, "New key should be different from initial key")

	// Verify old key is invalidated (should not validate)
	_, err = m.VerifyTurnSheetKey(initialKey)
	require.Error(t, err, "Old key should not validate")

	// Verify new key works
	rec, err := m.VerifyTurnSheetKey(newKey)
	require.NoError(t, err, "New key should validate")
	require.Equal(t, subscriptionID, rec.ID, "Subscription ID matches")
}
