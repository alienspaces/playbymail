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
				AccountUserConfigs: []harness.AccountUserConfig{
					{
						Reference: "test-account-user",
						Record: &account_record.AccountUser{
							Email:  harness.UniqueEmail("game-test@example.com"),
							Status: account_record.AccountUserStatusActive,
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
				if tc.rec == nil {
					require.Nil(t, rec)
				}
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
				Reference: "game-draft-1",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Draft Game 1"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
					Status:            game_record.GameStatusDraft,
				},
			},
			{
				Reference: "game-draft-2",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Draft Game 2"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
					Status:            game_record.GameStatusDraft,
				},
			},
			{
				Reference: "game-published",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Published Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
					Status:            game_record.GameStatusPublished,
				},
			},
		},
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "test-account",
				AccountUserConfigs: []harness.AccountUserConfig{
					{
						Reference: "test-account-user",
						Record: &account_record.AccountUser{
							Email:  harness.UniqueEmail("status-test@example.com"),
							Status: account_record.AccountUserStatusActive,
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
	draftGame1, err := th.Data.GetGameRecByRef("game-draft-1")
	require.NoError(t, err)
	draftGame2, err := th.Data.GetGameRecByRef("game-draft-2")
	require.NoError(t, err)
	publishedGame, err := th.Data.GetGameRecByRef("game-published")
	require.NoError(t, err)

	t.Run("allows draft to draft update", func(t *testing.T) {
		updateRec := &game_record.Game{
			Name:              harness.UniqueName("Updated Name"),
			GameType:          game_record.GameTypeAdventure,
			TurnDurationHours: 72,
			Description:       "Updated description",
			Status:            game_record.GameStatusDraft,
		}
		updateRec.ID = draftGame1.ID

		rec, err := m.UpdateGameRec(updateRec)
		require.NoError(t, err)
		require.NotNil(t, rec)
		require.Equal(t, game_record.GameStatusDraft, rec.Status)
	})

	t.Run("allows draft to published transition", func(t *testing.T) {
		updateRec := &game_record.Game{
			Name:              draftGame2.Name,
			GameType:          game_record.GameTypeAdventure,
			TurnDurationHours: draftGame2.TurnDurationHours,
			Description:       draftGame2.Description,
			Status:            game_record.GameStatusPublished,
		}
		updateRec.ID = draftGame2.ID

		rec, err := m.UpdateGameRec(updateRec)
		require.NoError(t, err)
		require.NotNil(t, rec)
		require.Equal(t, game_record.GameStatusPublished, rec.Status)
	})

	t.Run("allows modification of published game", func(t *testing.T) {
		updateRec := &game_record.Game{
			Name:              harness.UniqueName("Changed Name"),
			GameType:          game_record.GameTypeAdventure,
			TurnDurationHours: 48,
			Description:       "Changed description",
			Status:            game_record.GameStatusPublished,
		}
		updateRec.ID = publishedGame.ID

		rec, err := m.UpdateGameRec(updateRec)
		require.NoError(t, err)
		require.NotNil(t, rec)
		require.Equal(t, game_record.GameStatusPublished, rec.Status)
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
				AccountUserConfigs: []harness.AccountUserConfig{
					{
						Reference: "test-account-user",
						Record: &account_record.AccountUser{
							Email:  harness.UniqueEmail("ready-test@example.com"),
							Status: account_record.AccountUserStatusActive,
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

func TestValidateAdventureGameObjectStateGraph(t *testing.T) {
	// These tests exercise the cross-object effect and destructible-object suppression
	// paths in validateAdventureGameObjectStateGraph. The key behaviours under test:
	//
	//  1. A state reachable only via a cross-object change_object_state effect from
	//     another object must NOT be flagged as unreachable.
	//  2. States on an object that is targeted by cross-object effects must NOT be
	//     flagged as dead-ends, because external effects can unconditionally move the
	//     object without requiring a specific state.
	//  3. A terminal state on an object that has a remove_object effect must NOT be
	//     flagged as a dead-end, because the object is intentionally destroyed in that state.

	dataConfig := harness.DataConfig{
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "val-account",
				AccountUserConfigs: []harness.AccountUserConfig{
					{
						Reference: "val-account-user",
						Record: &account_record.AccountUser{
							Email:  harness.UniqueEmail("validate-state-graph@example.com"),
							Status: account_record.AccountUserStatusActive,
						},
					},
				},
			},
		},
		GameConfigs: []harness.GameConfig{
			// ── Game 1: cross-object effects ────────────────────────────────────────
			// Object A (lever) uses change_object_state and reveal_object to control
			// Object B (door). Object B's "revealed" state is only reachable via
			// Object A's cross-object effect; Object B's states are only exited via
			// Object A's cross-object hide_object. Neither should produce warnings.
			{
				Reference: "game-cross-object",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Cross Object Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 24,
				},
				AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
					{
						Reference: "cross-loc-one",
						Record: &adventure_game_record.AdventureGameLocation{
							Name:               harness.UniqueName("Start Room"),
							Description:        "The starting room.",
							IsStartingLocation: true,
						},
					},
				},
				AdventureGameLocationObjectConfigs: []harness.AdventureGameLocationObjectConfig{
					// Object B (door) — created first so Object A can reference its states.
					// Its "revealed" state is produced only by Object A's cross-object
					// change_object_state effect; no own effect produces it.
					{
						Reference:       "cross-obj-door",
						LocationRef:     "cross-loc-one",
						InitialStateRef: "cross-obj-door-state-hidden",
						Record: &adventure_game_record.AdventureGameLocationObject{
							Name:        "Hidden Door",
							Description: "A door hidden in the wall.",
							IsHidden:    true,
						},
						AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
							{
								Reference: "cross-obj-door-state-hidden",
								Record: &adventure_game_record.AdventureGameLocationObjectState{
									Name:      "hidden",
									SortOrder: 0,
								},
							},
							{
								Reference: "cross-obj-door-state-revealed",
								Record: &adventure_game_record.AdventureGameLocationObjectState{
									Name:      "revealed",
									SortOrder: 1,
								},
							},
						},
						// The door has one own effect: push (when revealed) → open.
						// This means requiredStateIDs = {"revealed"}, so "hidden" would
						// normally look like a dead-end. The cross-object suppression must
						// prevent that warning.
						AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
							{
								Reference:        "cross-obj-door-effect-push",
								RequiredStateRef: "cross-obj-door-state-revealed",
								Record: &adventure_game_record.AdventureGameLocationObjectEffect{
									ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
									EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
									ResultDescription: "The door swings open.",
									IsRepeatable:      false,
								},
							},
						},
					},
					// Object A (lever) — owns the cross-object effects that drive Object B.
					{
						Reference:       "cross-obj-lever",
						LocationRef:     "cross-loc-one",
						InitialStateRef: "cross-obj-lever-state-up",
						Record: &adventure_game_record.AdventureGameLocationObject{
							Name:        "Lever",
							Description: "A wall lever.",
						},
						AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
							{
								Reference: "cross-obj-lever-state-up",
								Record: &adventure_game_record.AdventureGameLocationObjectState{
									Name:      "up",
									SortOrder: 0,
								},
							},
							{
								Reference: "cross-obj-lever-state-down",
								Record: &adventure_game_record.AdventureGameLocationObjectState{
									Name:      "down",
									SortOrder: 1,
								},
							},
						},
						AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
							// pull (up) → lever becomes down
							{
								Reference:        "cross-obj-lever-effect-pull-self",
								RequiredStateRef: "cross-obj-lever-state-up",
								ResultStateRef:   "cross-obj-lever-state-down",
								Record: &adventure_game_record.AdventureGameLocationObjectEffect{
									ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
									EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
									ResultDescription: "You pull the lever.",
									IsRepeatable:      false,
								},
							},
							// pull (up) → door becomes revealed (cross-object change_object_state)
							{
								Reference:        "cross-obj-lever-effect-pull-door-state",
								RequiredStateRef: "cross-obj-lever-state-up",
								ResultObjectRef:  "cross-obj-door",
								ResultStateRef:   "cross-obj-door-state-revealed",
								Record: &adventure_game_record.AdventureGameLocationObjectEffect{
									ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
									EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeObjectState,
									ResultDescription: "The wall shifts.",
									IsRepeatable:      false,
								},
							},
							// push (down) → lever becomes up
							{
								Reference:        "cross-obj-lever-effect-push-self",
								RequiredStateRef: "cross-obj-lever-state-down",
								ResultStateRef:   "cross-obj-lever-state-up",
								Record: &adventure_game_record.AdventureGameLocationObjectEffect{
									ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
									EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
									ResultDescription: "You push the lever back.",
									IsRepeatable:      false,
								},
							},
							// push (down) → door hidden (cross-object hide_object)
							{
								Reference:        "cross-obj-lever-effect-push-hide-door",
								RequiredStateRef: "cross-obj-lever-state-down",
								ResultObjectRef:  "cross-obj-door",
								Record: &adventure_game_record.AdventureGameLocationObjectEffect{
									ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
									EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHideObject,
									ResultDescription: "",
									IsRepeatable:      false,
								},
							},
						},
					},
				},
			},
			// ── Game 2: destructible object ──────────────────────────────────────────
			// Object (rack) has a "burned" terminal state followed by remove_object.
			// The "burned" state must NOT be flagged as a dead-end.
			{
				Reference: "game-destructible",
				Record: &game_record.Game{
					Name:              harness.UniqueName("Destructible Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 24,
				},
				AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
					{
						Reference: "dest-loc-one",
						Record: &adventure_game_record.AdventureGameLocation{
							Name:               harness.UniqueName("Infirmary"),
							Description:        "The infirmary.",
							IsStartingLocation: true,
						},
					},
				},
				AdventureGameLocationObjectConfigs: []harness.AdventureGameLocationObjectConfig{
					{
						Reference:       "dest-obj-rack",
						LocationRef:     "dest-loc-one",
						InitialStateRef: "dest-obj-rack-state-intact",
						Record: &adventure_game_record.AdventureGameLocationObject{
							Name:        "Herb Drying Rack",
							Description: "A wooden rack of herbs.",
						},
						AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
							{
								Reference: "dest-obj-rack-state-intact",
								Record: &adventure_game_record.AdventureGameLocationObjectState{
									Name:      "intact",
									SortOrder: 0,
								},
							},
							{
								Reference: "dest-obj-rack-state-burned",
								Record: &adventure_game_record.AdventureGameLocationObjectState{
									Name:      "burned",
									SortOrder: 1,
								},
							},
						},
						AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
							// burn (intact) → change_state to burned
							{
								Reference:        "dest-obj-rack-effect-burn-state",
								RequiredStateRef: "dest-obj-rack-state-intact",
								ResultStateRef:   "dest-obj-rack-state-burned",
								Record: &adventure_game_record.AdventureGameLocationObjectEffect{
									ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBurn,
									EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
									ResultDescription: "The rack burns away.",
									IsRepeatable:      false,
								},
							},
							// burn (intact) → remove_object (rack is consumed; "burned" is the
							// terminal state before removal — must not be flagged as a dead-end)
							{
								Reference:        "dest-obj-rack-effect-burn-remove",
								RequiredStateRef: "dest-obj-rack-state-intact",
								Record: &adventure_game_record.AdventureGameLocationObjectEffect{
									ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBurn,
									EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveObject,
									ResultDescription: "",
									IsRepeatable:      false,
								},
							},
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

	t.Run("no warnings when target object states are driven by cross-object effects", func(t *testing.T) {
		// The Hidden Door's "revealed" state is only produced by the Lever's
		// change_object_state effect, not by any of the door's own effects. Without
		// cross-object awareness the validator would flag "revealed" as unreachable
		// and both "hidden" and "revealed" as dead-ends.
		gameID, ok := th.Data.Refs.GameRefs["game-cross-object"]
		require.True(t, ok)

		issues, err := m.ValidateGameReadyForInstance(gameID)
		require.NoError(t, err)

		for _, issue := range issues {
			require.NotContains(t, issue.Message, "Hidden Door",
				"Hidden Door should not produce any state-graph warnings; got: %s", issue.Message)
		}
	})

	t.Run("no dead-end warning for terminal state on destructible object", func(t *testing.T) {
		// The "burned" state on the Herb Drying Rack has no own effects that require
		// it (the rack is removed when burned). Without destructible-object suppression
		// the validator would flag "burned" as a dead-end.
		gameID, ok := th.Data.Refs.GameRefs["game-destructible"]
		require.True(t, ok)

		issues, err := m.ValidateGameReadyForInstance(gameID)
		require.NoError(t, err)

		for _, issue := range issues {
			require.NotContains(t, issue.Message, "dead-end",
				"Destructible object terminal state should not produce a dead-end warning; got: %s", issue.Message)
		}
	})
}
