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

func TestCreateGameRec_Validation(t *testing.T) {
	dataConfig := harness.DataConfig{
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "test-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("game-test@example.com"),
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
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	testCases := []struct {
		name        string
		rec         *game_record.Game
		expectError bool
	}{
		{
			name: "succeeds with valid adventure game",
			rec: &game_record.Game{
				Name:              harness.UniqueName("Valid Game"),
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 168,
				Description:       "A valid game description",
			},
			expectError: false,
		},
		{
			name: "succeeds and defaults status to draft",
			rec: &game_record.Game{
				Name:              harness.UniqueName("Draft Game"),
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 24,
				Description:       "A draft game",
			},
			expectError: false,
		},
		{
			name: "fails when name is empty",
			rec: &game_record.Game{
				Name:              "",
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 168,
				Description:       "Missing name",
			},
			expectError: true,
		},
		{
			name: "fails with invalid game type",
			rec: &game_record.Game{
				Name:              harness.UniqueName("Bad Type"),
				GameType:          "invalid_type",
				TurnDurationHours: 168,
				Description:       "Invalid type",
			},
			expectError: true,
		},
		{
			name: "fails when turn duration is zero",
			rec: &game_record.Game{
				Name:              harness.UniqueName("Zero Duration"),
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 0,
				Description:       "No duration",
			},
			expectError: true,
		},
		{
			name: "fails when turn duration is negative",
			rec: &game_record.Game{
				Name:              harness.UniqueName("Negative Duration"),
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: -1,
				Description:       "Negative duration",
			},
			expectError: true,
		},
		{
			name: "fails when description is empty",
			rec: &game_record.Game{
				Name:              harness.UniqueName("No Desc"),
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 168,
				Description:       "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := th.Domain.(*domain.Domain)

			rec, err := m.CreateGameRec(tc.rec)

			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, rec)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rec)
				require.Equal(t, game_record.GameStatusDraft, rec.Status)
			}
		})
	}
}

func TestUpdateGameRec_StatusTransitions(t *testing.T) {
	dataConfig := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record: &game_record.Game{
					Name:              harness.UniqueName("Status Test Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
				},
			},
		},
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "test-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("status-test@example.com"),
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
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	m := th.Domain.(*domain.Domain)
	gameRec := th.Data.GameRecs[0]

	t.Run("allows draft to draft update", func(t *testing.T) {
		updateRec := &game_record.Game{
			Name:              harness.UniqueName("Updated Name"),
			GameType:          game_record.GameTypeAdventure,
			TurnDurationHours: 72,
			Description:       "Updated description",
			Status:            game_record.GameStatusDraft,
		}
		updateRec.ID = gameRec.ID

		rec, err := m.UpdateGameRec(updateRec)
		require.NoError(t, err)
		require.NotNil(t, rec)
		require.Equal(t, game_record.GameStatusDraft, rec.Status)
	})

	t.Run("allows draft to published transition", func(t *testing.T) {
		updateRec := &game_record.Game{
			Name:              gameRec.Name,
			GameType:          game_record.GameTypeAdventure,
			TurnDurationHours: gameRec.TurnDurationHours,
			Description:       gameRec.Description,
			Status:            game_record.GameStatusPublished,
		}
		updateRec.ID = gameRec.ID

		rec, err := m.UpdateGameRec(updateRec)
		require.NoError(t, err)
		require.NotNil(t, rec)
		require.Equal(t, game_record.GameStatusPublished, rec.Status)
	})

	t.Run("prevents modification of published game", func(t *testing.T) {
		updateRec := &game_record.Game{
			Name:              harness.UniqueName("Changed Name"),
			GameType:          game_record.GameTypeAdventure,
			TurnDurationHours: 48,
			Description:       "Changed description",
			Status:            game_record.GameStatusPublished,
		}
		updateRec.ID = gameRec.ID

		_, err := m.UpdateGameRec(updateRec)
		require.Error(t, err)
	})
}

func TestValidateGameReadyForInstance(t *testing.T) {
	dataConfig := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: "game-with-starting-location",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Ready Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
				},
				AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
					{
						Reference: "starting-location",
						Record: &adventure_game_record.AdventureGameLocation{
							Name:               harness.UniqueName("Starting Location"),
							Description:        "A starting location",
							IsStartingLocation: true,
						},
					},
				},
			},
			{
				Reference: "game-without-locations",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Empty Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
				},
			},
			{
				Reference: "game-without-starting-location",
				Record: &game_record.Game{
					Name:              harness.UniqueName("No Start Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
				},
				AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
					{
						Reference: "non-starting-location",
						Record: &adventure_game_record.AdventureGameLocation{
							Name:               harness.UniqueName("Normal Location"),
							Description:        "Not a starting location",
							IsStartingLocation: false,
						},
					},
				},
			},
		},
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "test-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("ready-test@example.com"),
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
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	m := th.Domain.(*domain.Domain)

	t.Run("returns no issues for game with starting location", func(t *testing.T) {
		gameID, ok := th.Data.Refs.GameRefs["game-with-starting-location"]
		require.True(t, ok)

		issues, err := m.ValidateGameReadyForInstance(gameID)
		require.NoError(t, err)
		require.Empty(t, issues)
	})

	t.Run("returns issue when game has no locations", func(t *testing.T) {
		gameID, ok := th.Data.Refs.GameRefs["game-without-locations"]
		require.True(t, ok)

		issues, err := m.ValidateGameReadyForInstance(gameID)
		require.NoError(t, err)
		require.NotEmpty(t, issues)

		foundLocationIssue := false
		for _, issue := range issues {
			if issue.Field == "locations" {
				foundLocationIssue = true
				require.Equal(t, domain.ValidationSeverityError, issue.Severity)
			}
		}
		require.True(t, foundLocationIssue, "should report missing locations issue")
	})

	t.Run("returns issue when game has no starting location", func(t *testing.T) {
		gameID, ok := th.Data.Refs.GameRefs["game-without-starting-location"]
		require.True(t, ok)

		issues, err := m.ValidateGameReadyForInstance(gameID)
		require.NoError(t, err)
		require.NotEmpty(t, issues)

		foundStartingIssue := false
		for _, issue := range issues {
			if issue.Field == "starting_location" {
				foundStartingIssue = true
				require.Equal(t, domain.ValidationSeverityError, issue.Severity)
			}
		}
		require.True(t, foundStartingIssue, "should report missing starting location issue")
	})

	t.Run("returns error for non-existent game", func(t *testing.T) {
		_, err := m.ValidateGameReadyForInstance("00000000-0000-0000-0000-000000000000")
		require.Error(t, err)
	})
}
