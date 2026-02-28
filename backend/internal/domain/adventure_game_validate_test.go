package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func setupAdventureGameHarness(t *testing.T) (*harness.Testing, *domain.Domain) {
	t.Helper()

	dataConfig := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record: &game_record.Game{
					Name:              harness.UniqueName("Adventure Test Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
				},
				AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
					{
						Reference: "test-location",
						Record: &adventure_game_record.AdventureGameLocation{
							Name:               harness.UniqueName("Test Location"),
							Description:        "A test location",
							IsStartingLocation: true,
						},
					},
				},
				AdventureGameItemConfigs: []harness.AdventureGameItemConfig{
					{
						Reference: "test-item",
						Record: &adventure_game_record.AdventureGameItem{
							Name:        harness.UniqueName("Test Item"),
							Description: "A test item",
						},
					},
				},
				AdventureGameCreatureConfigs: []harness.AdventureGameCreatureConfig{
					{
						Reference: "test-creature",
						Record: &adventure_game_record.AdventureGameCreature{
							Name:        harness.UniqueName("Test Creature"),
							Description: "A test creature",
						},
					},
				},
				AdventureGameCharacterConfigs: []harness.AdventureGameCharacterConfig{
					{
						Reference:  "test-character",
						AccountRef: "test-account",
						Record: &adventure_game_record.AdventureGameCharacter{
							Name: harness.UniqueName("Test Character"),
						},
					},
				},
			},
		},
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "test-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("adventure-test@example.com"),
					Status: account_record.AccountUserStatusActive,
				},
			},
			{
				Reference: "test-account-2",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("adventure-test-2@example.com"),
					Status: account_record.AccountUserStatusActive,
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

	m := th.Domain.(*domain.Domain)
	return th, m
}

func TestCreateAdventureGameLocationRec_Validation(t *testing.T) {
	th, m := setupAdventureGameHarness(t)
	defer func() {
		err := th.Teardown()
		require.NoError(t, err)
	}()

	gameID := th.Data.GameRecs[0].ID

	testCases := []struct {
		name        string
		rec         *adventure_game_record.AdventureGameLocation
		expectError bool
	}{
		{
			name: "succeeds with valid location",
			rec: &adventure_game_record.AdventureGameLocation{
				GameID:             gameID,
				Name:               harness.UniqueName("New Location"),
				Description:        "A new location",
				IsStartingLocation: false,
			},
			expectError: false,
		},
		{
			name:        "fails when record is nil",
			rec:         nil,
			expectError: true,
		},
		{
			name: "fails when game_id is empty",
			rec: &adventure_game_record.AdventureGameLocation{
				GameID:      "",
				Name:        harness.UniqueName("No Game ID"),
				Description: "Missing game ID",
			},
			expectError: true,
		},
		{
			name: "fails when name is empty",
			rec: &adventure_game_record.AdventureGameLocation{
				GameID:      gameID,
				Name:        "",
				Description: "Missing name",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, err := m.CreateAdventureGameLocationRec(tc.rec)
			if tc.expectError {
				require.Error(t, err)
				if tc.rec == nil {
					require.Nil(t, rec)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, rec)
				require.NotEmpty(t, rec.ID)
			}
		})
	}
}

func TestCreateAdventureGameItemRec_Validation(t *testing.T) {
	th, m := setupAdventureGameHarness(t)
	defer func() {
		err := th.Teardown()
		require.NoError(t, err)
	}()

	gameID := th.Data.GameRecs[0].ID

	testCases := []struct {
		name        string
		rec         *adventure_game_record.AdventureGameItem
		expectError bool
	}{
		{
			name: "succeeds with valid item",
			rec: &adventure_game_record.AdventureGameItem{
				GameID:      gameID,
				Name:        harness.UniqueName("New Item"),
				Description: "A new item",
			},
			expectError: false,
		},
		{
			name:        "fails when record is nil",
			rec:         nil,
			expectError: true,
		},
		{
			name: "fails when game_id is empty",
			rec: &adventure_game_record.AdventureGameItem{
				GameID:      "",
				Name:        harness.UniqueName("No Game"),
				Description: "Missing game ID",
			},
			expectError: true,
		},
		{
			name: "fails when name is empty",
			rec: &adventure_game_record.AdventureGameItem{
				GameID:      gameID,
				Name:        "",
				Description: "Missing name",
			},
			expectError: true,
		},
		{
			name: "fails when description is empty",
			rec: &adventure_game_record.AdventureGameItem{
				GameID:      gameID,
				Name:        harness.UniqueName("No Desc Item"),
				Description: "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, err := m.CreateAdventureGameItemRec(tc.rec)
			if tc.expectError {
				require.Error(t, err)
				if tc.rec == nil {
					require.Nil(t, rec)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, rec)
				require.NotEmpty(t, rec.ID)
			}
		})
	}
}

func TestCreateAdventureGameCreatureRec_Validation(t *testing.T) {
	th, m := setupAdventureGameHarness(t)
	defer func() {
		err := th.Teardown()
		require.NoError(t, err)
	}()

	gameID := th.Data.GameRecs[0].ID

	testCases := []struct {
		name        string
		rec         *adventure_game_record.AdventureGameCreature
		expectError bool
	}{
		{
			name: "succeeds with valid creature",
			rec: &adventure_game_record.AdventureGameCreature{
				GameID:      gameID,
				Name:        harness.UniqueName("New Creature"),
				Description: "A new creature",
			},
			expectError: false,
		},
		{
			name:        "fails when record is nil",
			rec:         nil,
			expectError: true,
		},
		{
			name: "fails when game_id is empty",
			rec: &adventure_game_record.AdventureGameCreature{
				GameID:      "",
				Name:        harness.UniqueName("No Game"),
				Description: "Missing game ID",
			},
			expectError: true,
		},
		{
			name: "fails when name is empty",
			rec: &adventure_game_record.AdventureGameCreature{
				GameID:      gameID,
				Name:        "",
				Description: "Missing name",
			},
			expectError: true,
		},
		{
			name: "fails when description is empty",
			rec: &adventure_game_record.AdventureGameCreature{
				GameID:      gameID,
				Name:        harness.UniqueName("No Desc"),
				Description: "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, err := m.CreateAdventureGameCreatureRec(tc.rec)
			if tc.expectError {
				require.Error(t, err)
				if tc.rec == nil {
					require.Nil(t, rec)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, rec)
				require.NotEmpty(t, rec.ID)
			}
		})
	}
}

func TestCreateAdventureGameCharacterRec_Validation(t *testing.T) {
	th, m := setupAdventureGameHarness(t)
	defer func() {
		err := th.Teardown()
		require.NoError(t, err)
	}()

	gameID := th.Data.GameRecs[0].ID
	// Use test-account-2 for new character: test-account already has a character (test-character) for this game; unique is (game_id, account_id).
	accountRec, err := th.Data.GetAccountUserRecByRef("test-account-2")
	require.NoError(t, err)
	accountRecForFailCase, err := th.Data.GetAccountUserRecByRef("test-account")
	require.NoError(t, err)

	testCases := []struct {
		name        string
		rec         *adventure_game_record.AdventureGameCharacter
		expectError bool
	}{
		{
			name: "succeeds with valid character",
			rec: &adventure_game_record.AdventureGameCharacter{
				GameID:        gameID,
				AccountID:     accountRec.AccountID,
				AccountUserID: accountRec.ID,
				Name:          harness.UniqueName("New Character"),
			},
			expectError: false,
		},
		{
			name: "fails when name is empty",
			rec: &adventure_game_record.AdventureGameCharacter{
				GameID:        gameID,
				AccountID:     accountRecForFailCase.AccountID,
				AccountUserID: accountRecForFailCase.ID,
				Name:          "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, err := m.CreateAdventureGameCharacterRec(tc.rec)
			if tc.expectError {
				require.Error(t, err)
				if tc.rec == nil {
					require.Nil(t, rec)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, rec)
				require.NotEmpty(t, rec.ID)
			}
		})
	}
}
