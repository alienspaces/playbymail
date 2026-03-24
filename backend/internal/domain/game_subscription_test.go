package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/nullint32"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestApproveGameSubscription(t *testing.T) {
	dataConfig := harness.DataConfig{
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "test-account",
				AccountUserConfigs: []harness.AccountUserConfig{
					{
						Reference: "test-account-user",
						Record: &account_record.AccountUser{
							Email:  harness.UniqueEmail("approve-sub@example.com"),
							Status: account_record.AccountUserStatusActive,
						},
					},
				},
			},
		},
		GameConfigs: []harness.GameConfig{
			{
				Reference: "test-game",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Approve Subscription Test Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
					// Status defaults to published in the harness
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err)

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err)

	th, err := harness.NewTesting(cfg, l, s, j, scanner, dataConfig)
	require.NoError(t, err)

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	m := th.Domain.(*domain.Domain)

	accountUserRec, err := th.Data.GetAccountUserRecByRef("test-account-user")
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	accountUserContactRec, err := th.Data.GetAccountUserContactRecByAccountUserID(accountUserRec.ID)
	require.NoError(t, err, "GetAccountUserContactRecByAccountUserID returns without error")

	gameRec, err := th.Data.GetGameRecByRef("test-game")
	require.NoError(t, err, "GetGameRecByRef returns without error")

	// Pending subscription used for the wrong-email and success test cases (in that order).
	pendingSubRec, err := m.CreateGameSubscriptionRec(&game_record.GameSubscription{
		GameID:               gameRec.ID,
		AccountID:            accountUserRec.AccountID,
		AccountUserID:        accountUserRec.ID,
		AccountUserContactID: nullstring.FromString(accountUserContactRec.ID),
		SubscriptionType:     game_record.GameSubscriptionTypePlayer,
		Status:               game_record.GameSubscriptionStatusPendingApproval,
		DeliveryMethod:       nullstring.FromString(game_record.GameSubscriptionDeliveryMethodEmail),
	})
	require.NoError(t, err, "creating pending player subscription returns without error")

	// Active subscription used for the "not pending" test case.
	activeSubRec, err := m.CreateGameSubscriptionRec(&game_record.GameSubscription{
		GameID:               gameRec.ID,
		AccountID:            accountUserRec.AccountID,
		AccountUserID:        accountUserRec.ID,
		AccountUserContactID: nullstring.FromString(accountUserContactRec.ID),
		SubscriptionType:     game_record.GameSubscriptionTypePlayer,
		Status:               game_record.GameSubscriptionStatusActive,
		DeliveryMethod:       nullstring.FromString(game_record.GameSubscriptionDeliveryMethodEmail),
	})
	require.NoError(t, err, "creating active player subscription returns without error")

	testCases := []struct {
		name           string
		subscriptionID string
		email          string
		expectError    bool
		expectStatus   string
	}{
		{
			// Email check fires before any DB access, so subscription state is irrelevant.
			name:           "error when subscription ID is empty",
			subscriptionID: "",
			email:          accountUserRec.Email,
			expectError:    true,
		},
		{
			// Email check fires before subscription is fetched; any valid ID works here.
			name:           "error when email is empty",
			subscriptionID: pendingSubRec.ID,
			email:          "",
			expectError:    true,
		},
		{
			name:           "error when subscription does not exist",
			subscriptionID: "00000000-0000-0000-0000-000000000000",
			email:          accountUserRec.Email,
			expectError:    true,
		},
		{
			name:           "error when subscription is not pending approval",
			subscriptionID: activeSubRec.ID,
			email:          accountUserRec.Email,
			expectError:    true,
		},
		{
			// Must run before the success case so pendingSubRec is still pending.
			name:           "error when email does not match subscription account",
			subscriptionID: pendingSubRec.ID,
			email:          "wrong-email@example.com",
			expectError:    true,
		},
		{
			// Runs last: transitions pendingSubRec from pending_approval to active.
			name:           "succeeds with correct email",
			subscriptionID: pendingSubRec.ID,
			email:          accountUserRec.Email,
			expectError:    false,
			expectStatus:   game_record.GameSubscriptionStatusActive,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, err := m.ApproveGameSubscription(tc.subscriptionID, tc.email)

			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, rec)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, rec)
			require.Equal(t, tc.expectStatus, rec.Status)
		})
	}
}

func TestCreateGameSubscriptionRec_Validation(t *testing.T) {
	dataConfig := harness.DataConfig{
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "test-account",
				AccountUserConfigs: []harness.AccountUserConfig{
					{
						Reference: "test-account-user",
						Record: &account_record.AccountUser{
							Email:  harness.UniqueEmail("sub-validate@example.com"),
							Status: account_record.AccountUserStatusActive,
						},
					},
				},
			},
		},
		GameConfigs: []harness.GameConfig{
			{
				Reference: "published-game",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Published Validation Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
					// Status defaults to published in the harness
				},
			},
			{
				Reference: "draft-game",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Draft Validation Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
					Status:            game_record.GameStatusDraft,
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err)

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err)

	th, err := harness.NewTesting(cfg, l, s, j, scanner, dataConfig)
	require.NoError(t, err)

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	m := th.Domain.(*domain.Domain)

	accountUserRec, err := th.Data.GetAccountUserRecByRef("test-account-user")
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	accountUserContactRec, err := th.Data.GetAccountUserContactRecByAccountUserID(accountUserRec.ID)
	require.NoError(t, err, "GetAccountUserContactRecByAccountUserID returns without error")

	publishedGameRec, err := th.Data.GetGameRecByRef("published-game")
	require.NoError(t, err, "GetGameRecByRef returns without error")

	draftGameRec, err := th.Data.GetGameRecByRef("draft-game")
	require.NoError(t, err, "GetGameRecByRef returns without error")

	validPlayerRec := func() *game_record.GameSubscription {
		return &game_record.GameSubscription{
			GameID:               publishedGameRec.ID,
			AccountID:            accountUserRec.AccountID,
			AccountUserID:        accountUserRec.ID,
			AccountUserContactID: nullstring.FromString(accountUserContactRec.ID),
			SubscriptionType:     game_record.GameSubscriptionTypePlayer,
			Status:               game_record.GameSubscriptionStatusActive,
			DeliveryMethod:       nullstring.FromString(game_record.GameSubscriptionDeliveryMethodEmail),
		}
	}

	validManagerRec := func() *game_record.GameSubscription {
		return &game_record.GameSubscription{
			GameID:           publishedGameRec.ID,
			AccountID:        accountUserRec.AccountID,
			AccountUserID:    accountUserRec.ID,
			SubscriptionType: game_record.GameSubscriptionTypeManager,
			Status:           game_record.GameSubscriptionStatusActive,
		}
	}

	testCases := []struct {
		name        string
		rec         *game_record.GameSubscription
		expectError bool
	}{
		{
			name:        "succeeds creating manager subscription for published game",
			rec:         validManagerRec(),
			expectError: false,
		},
		{
			name: "succeeds creating designer subscription for published game",
			rec: &game_record.GameSubscription{
				GameID:           publishedGameRec.ID,
				AccountID:        accountUserRec.AccountID,
				AccountUserID:    accountUserRec.ID,
				SubscriptionType: game_record.GameSubscriptionTypeDesigner,
				Status:           game_record.GameSubscriptionStatusActive,
			},
			expectError: false,
		},
		{
			name: "succeeds creating manager subscription for draft game",
			rec: &game_record.GameSubscription{
				GameID:           draftGameRec.ID,
				AccountID:        accountUserRec.AccountID,
				AccountUserID:    accountUserRec.ID,
				SubscriptionType: game_record.GameSubscriptionTypeManager,
				Status:           game_record.GameSubscriptionStatusActive,
			},
			expectError: false,
		},
		{
			name:        "succeeds creating player subscription with all required fields",
			rec:         validPlayerRec(),
			expectError: false,
		},
		{
			name: "fails when game_id is empty",
			rec: func() *game_record.GameSubscription {
				r := validManagerRec()
				r.GameID = ""
				return r
			}(),
			expectError: true,
		},
		{
			name: "fails when account_id is empty",
			rec: func() *game_record.GameSubscription {
				r := validManagerRec()
				r.AccountID = ""
				return r
			}(),
			expectError: true,
		},
		{
			name: "fails when account_user_id is empty",
			rec: func() *game_record.GameSubscription {
				r := validManagerRec()
				r.AccountUserID = ""
				return r
			}(),
			expectError: true,
		},
		{
			name: "fails when subscription_type is empty",
			rec: func() *game_record.GameSubscription {
				r := validManagerRec()
				r.SubscriptionType = ""
				return r
			}(),
			expectError: true,
		},
		{
			name: "fails when status is empty",
			rec: func() *game_record.GameSubscription {
				r := validManagerRec()
				r.Status = ""
				return r
			}(),
			expectError: true,
		},
		{
			name: "fails when player subscription is missing delivery method",
			rec: &game_record.GameSubscription{
				GameID:               publishedGameRec.ID,
				AccountID:            accountUserRec.AccountID,
				AccountUserID:        accountUserRec.ID,
				AccountUserContactID: nullstring.FromString(accountUserContactRec.ID),
				SubscriptionType:     game_record.GameSubscriptionTypePlayer,
				Status:               game_record.GameSubscriptionStatusActive,
				// DeliveryMethod intentionally absent (zero sql.NullString)
			},
			expectError: true,
		},
		{
			name: "fails when player subscription is missing account_user_contact_id",
			rec: &game_record.GameSubscription{
				GameID:           publishedGameRec.ID,
				AccountID:        accountUserRec.AccountID,
				AccountUserID:    accountUserRec.ID,
				SubscriptionType: game_record.GameSubscriptionTypePlayer,
				Status:           game_record.GameSubscriptionStatusActive,
				DeliveryMethod:   nullstring.FromString(game_record.GameSubscriptionDeliveryMethodEmail),
				// AccountUserContactID intentionally absent (zero sql.NullString)
			},
			expectError: true,
		},
		{
			name: "fails when player subscription is for a draft game",
			rec: func() *game_record.GameSubscription {
				r := validPlayerRec()
				r.GameID = draftGameRec.ID
				return r
			}(),
			expectError: true,
		},
		{
			name: "fails when instance_limit is zero",
			rec: func() *game_record.GameSubscription {
				r := validManagerRec()
				r.InstanceLimit = nullint32.FromInt32(0)
				return r
			}(),
			expectError: true,
		},
		{
			name: "fails when instance_limit is negative",
			rec: func() *game_record.GameSubscription {
				r := validManagerRec()
				r.InstanceLimit = nullint32.FromInt32(-1)
				return r
			}(),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, err := m.CreateGameSubscriptionRec(tc.rec)

			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, rec)
			require.NotEmpty(t, rec.ID)
		})
	}
}

func TestValidateInstanceLimit(t *testing.T) {
	dataConfig := harness.DataConfig{
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "test-account",
				AccountUserConfigs: []harness.AccountUserConfig{
					{
						Reference: "test-account-user",
						Record: &account_record.AccountUser{
							Email:  harness.UniqueEmail("instance-limit@example.com"),
							Status: account_record.AccountUserStatusActive,
						},
					},
				},
			},
		},
		GameConfigs: []harness.GameConfig{
			{
				Reference: "test-game",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Instance Limit Test Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
				},
				AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
					{
						Reference: "starting-location",
						Record: &adventure_game_record.AdventureGameLocation{
							Name:               harness.UniqueName("Starting Location"),
							Description:        "Starting location for instance limit tests",
							IsStartingLocation: true,
						},
					},
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err)

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err)

	th, err := harness.NewTesting(cfg, l, s, j, scanner, dataConfig)
	require.NoError(t, err)

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	m := th.Domain.(*domain.Domain)

	accountUserRec, err := th.Data.GetAccountUserRecByRef("test-account-user")
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	gameRec, err := th.Data.GetGameRecByRef("test-game")
	require.NoError(t, err, "GetGameRecByRef returns without error")

	newManagerSub := func(limit *int32) *game_record.GameSubscription {
		rec := &game_record.GameSubscription{
			GameID:           gameRec.ID,
			AccountID:        accountUserRec.AccountID,
			AccountUserID:    accountUserRec.ID,
			SubscriptionType: game_record.GameSubscriptionTypeManager,
			Status:           game_record.GameSubscriptionStatusActive,
		}
		if limit != nil {
			rec.InstanceLimit = nullint32.FromInt32(*limit)
		}
		return rec
	}

	// Subscription with no instance limit.
	unlimitedSub, err := m.CreateGameSubscriptionRec(newManagerSub(nil))
	require.NoError(t, err, "creating unlimited manager subscription returns without error")

	limitOne := int32(1)

	// Subscription with limit=1, 0 instances linked.
	belowLimitSub, err := m.CreateGameSubscriptionRec(newManagerSub(&limitOne))
	require.NoError(t, err, "creating limit-1 manager subscription returns without error")

	// Subscription with limit=1, 1 instance linked (at the limit).
	atLimitSub, err := m.CreateGameSubscriptionRec(newManagerSub(&limitOne))
	require.NoError(t, err, "creating limit-1 manager subscription (for at-limit test) returns without error")

	// Create a game instance and link it to atLimitSub.
	gameInst, err := m.CreateGameInstanceRec(&game_record.GameInstance{
		GameID:              gameRec.ID,
		Status:              game_record.GameInstanceStatusCreated,
		RequiredPlayerCount: 1,
		DeliveryEmail:       true,
	})
	require.NoError(t, err, "creating game instance returns without error")

	_, err = m.CreateGameSubscriptionInstanceRec(&game_record.GameSubscriptionInstance{
		AccountID:          atLimitSub.AccountID,
		AccountUserID:      atLimitSub.AccountUserID,
		GameSubscriptionID: atLimitSub.ID,
		GameInstanceID:     gameInst.ID,
	})
	require.NoError(t, err, "linking game instance to subscription returns without error")

	testCases := []struct {
		name           string
		subscriptionID string
		expectError    bool
	}{
		{
			name:           "error for invalid subscription UUID",
			subscriptionID: "not-a-uuid",
			expectError:    true,
		},
		{
			name:           "error when subscription does not exist",
			subscriptionID: "00000000-0000-0000-0000-000000000000",
			expectError:    true,
		},
		{
			name:           "no error for subscription with no limit",
			subscriptionID: unlimitedSub.ID,
			expectError:    false,
		},
		{
			name:           "no error when instance count is below limit",
			subscriptionID: belowLimitSub.ID,
			expectError:    false,
		},
		{
			name:           "error when instance count reaches limit",
			subscriptionID: atLimitSub.ID,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := m.ValidateInstanceLimit(tc.subscriptionID)

			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
