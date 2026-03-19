package turn_sheet_processor_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/adventure_game/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

// combatHarness returns a harness whose game creature has explicit attack/defense stats
// and whose character starts at the same location, making it straightforward to assert
// exact health values after ProcessTurnSheetResponse runs.
func combatHarness(t *testing.T) *harness.Testing {
	t.Helper()

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	dc := harness.DefaultDataConfig()

	// Override the default creature stats for predictable combat arithmetic:
	//   AttackDamage=10, Defense=2  →  player (unarmed, damage=5) deals max(5-2,1)=3 per hit
	//                               →  creature deals max(10-0,1)=10 per retaliation
	for i := range dc.GameConfigs {
		if dc.GameConfigs[i].Reference != harness.GameOneRef {
			continue
		}
		for j := range dc.GameConfigs[i].AdventureGameCreatureConfigs {
			if dc.GameConfigs[i].AdventureGameCreatureConfigs[j].Reference == harness.GameCreatureOneRef {
				dc.GameConfigs[i].AdventureGameCreatureConfigs[j].Record = &adventure_game_record.AdventureGameCreature{
					Name:           harness.UniqueName("Combat Test Creature"),
					AttackDamage:   10,
					Defense:        2,
					Disposition:    adventure_game_record.AdventureGameCreatureDispositionAggressive,
					AttackMethod:   adventure_game_record.AdventureGameCreatureAttackMethodClaws,
					BodyDecayTurns: 3,
				}
			}
		}
	}

	th, err := harness.NewTesting(cfg, l, s, j, scanner, dc)
	require.NoError(t, err, "NewTesting returns without error")

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")

	t.Cleanup(func() {
		require.NoError(t, th.Teardown(), "Test data teardown returns without error")
	})

	return th
}

// buildMonsterTurnSheet constructs a minimal GameTurnSheet with the given scanned combat actions.
func buildMonsterTurnSheet(t *testing.T, gameInstanceRec *game_record.GameInstance, actions []turnsheet.CombatAction) *game_record.GameTurnSheet {
	t.Helper()

	scanData := turnsheet.MonsterEncounterScanData{Actions: actions}
	scannedBytes, err := json.Marshal(scanData)
	require.NoError(t, err)

	return &game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		GameInstanceID:   sql.NullString{String: gameInstanceRec.ID, Valid: true},
		TurnNumber:       gameInstanceRec.CurrentTurn,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter,
		SheetData:        json.RawMessage(`{}`),
		ScannedData:      scannedBytes,
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
}

func TestAdventureGameCreatureEncounterProcessor_ProcessTurnSheetResponse(t *testing.T) {
	th := combatHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameCreatureEncounterProcessor(m.Log, m)
	require.NoError(t, err)

	// Resolve harness records once — sub-tests share these pointers.
	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	creatureInstanceRec, err := th.Data.GetAdventureGameCreatureInstanceRecByRef(harness.GameCreatureInstanceOneRef)
	require.NoError(t, err)

	// unarmedDamage=5, creatureDefense=2  →  playerDealt = max(5-2,1) = 3
	// creatureAttack=10, armorDefense=0   →  creatureDealt = max(10-0,1) = 10
	const (
		unarmedDamage       = 5
		creatureDefense     = 2
		creatureAttack      = 10
		armorDefense        = 0
		playerDealtPerHit   = unarmedDamage - creatureDefense // = 3
		creatureDealtPerHit = creatureAttack - armorDefense   // = 10
		startingHealth      = 100
	)

	tests := []struct {
		name string
		// actions to include in the turn sheet
		actions []turnsheet.CombatAction
		// expected creature health after processing
		wantCreatureHealth int
		// expected character health after processing
		wantCharacterHealth int
		// whether the creature should be marked dead
		wantCreatureDead bool
	}{
		{
			name:                "given no actions when processing then no health changes",
			actions:             []turnsheet.CombatAction{},
			wantCreatureHealth:  startingHealth,
			wantCharacterHealth: startingHealth,
			wantCreatureDead:    false,
		},
		{
			name: "given one attack when processing then creature and character health reduced",
			actions: []turnsheet.CombatAction{
				{ActionType: turnsheet.CombatActionTypeAttack, TargetCreatureInstanceID: creatureInstanceRec.ID},
			},
			// creature loses playerDealtPerHit=3; still alive so retaliates for creatureDealtPerHit=10
			wantCreatureHealth:  startingHealth - playerDealtPerHit,
			wantCharacterHealth: startingHealth - creatureDealtPerHit,
			wantCreatureDead:    false,
		},
		{
			name: "given three attacks when processing then both health values reduced for each round",
			actions: []turnsheet.CombatAction{
				{ActionType: turnsheet.CombatActionTypeAttack, TargetCreatureInstanceID: creatureInstanceRec.ID},
				{ActionType: turnsheet.CombatActionTypeAttack, TargetCreatureInstanceID: creatureInstanceRec.ID},
				{ActionType: turnsheet.CombatActionTypeAttack, TargetCreatureInstanceID: creatureInstanceRec.ID},
			},
			wantCreatureHealth:  startingHealth - 3*playerDealtPerHit,
			wantCharacterHealth: startingHealth - 3*creatureDealtPerHit,
			wantCreatureDead:    false,
		},
		{
			name: "given do_nothing action when processing then no health changes",
			actions: []turnsheet.CombatAction{
				{ActionType: turnsheet.CombatActionTypeDoNothing},
			},
			wantCreatureHealth:  startingHealth,
			wantCharacterHealth: startingHealth,
			wantCreatureDead:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset both health values before each sub-test so tests are independent.
			charInstanceRec.Health = startingHealth
			charInstanceRec.LastTurnEvents = []byte("[]")
			_, err := m.UpdateAdventureGameCharacterInstanceRec(charInstanceRec)
			require.NoError(t, err)

			creatureInstanceRec.Health = startingHealth
			creatureInstanceRec.DiedAtTurn = sql.NullInt64{}
			_, err = m.UpdateAdventureGameCreatureInstanceRec(creatureInstanceRec)
			require.NoError(t, err)

			turnSheet := buildMonsterTurnSheet(t, gameInstanceRec, tt.actions)

			err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
			require.NoError(t, err)

			// Re-fetch creature from DB to verify persisted health.
			updatedCreature, err := m.GetAdventureGameCreatureInstanceRec(creatureInstanceRec.ID, nil)
			require.NoError(t, err)
			require.Equal(t, tt.wantCreatureHealth, updatedCreature.Health,
				"creature health after processing")
			require.Equal(t, tt.wantCreatureDead, updatedCreature.DiedAtTurn.Valid,
				"creature dead-at-turn set correctly")

			// Character health is updated on the in-memory record and persisted inside the processor.
			require.Equal(t, tt.wantCharacterHealth, charInstanceRec.Health,
				"character health after processing")
		})
	}
}

func TestAdventureGameCreatureEncounterProcessor_ProcessTurnSheetResponse_CreatureDeath(t *testing.T) {
	th := combatHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameCreatureEncounterProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	creatureInstanceRec, err := th.Data.GetAdventureGameCreatureInstanceRecByRef(harness.GameCreatureInstanceOneRef)
	require.NoError(t, err)

	// Set creature health low so it dies on the first hit (playerDealtPerHit=3, health=2).
	creatureInstanceRec.Health = 2
	_, err = m.UpdateAdventureGameCreatureInstanceRec(creatureInstanceRec)
	require.NoError(t, err)

	charInstanceRec.Health = 100
	charInstanceRec.LastTurnEvents = []byte("[]")
	_, err = m.UpdateAdventureGameCharacterInstanceRec(charInstanceRec)
	require.NoError(t, err)

	actions := []turnsheet.CombatAction{
		{ActionType: turnsheet.CombatActionTypeAttack, TargetCreatureInstanceID: creatureInstanceRec.ID},
	}
	turnSheet := buildMonsterTurnSheet(t, gameInstanceRec, actions)

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	require.NoError(t, err)

	// Creature should be dead.
	updatedCreature, err := m.GetAdventureGameCreatureInstanceRec(creatureInstanceRec.ID, nil)
	require.NoError(t, err)
	require.Equal(t, 0, updatedCreature.Health, "dead creature health is zero")
	require.True(t, updatedCreature.DiedAtTurn.Valid, "DiedAtTurn is set")

	// Character should NOT have taken retaliation damage (creature died before it could retaliate).
	require.Equal(t, 100, charInstanceRec.Health, "character takes no retaliation damage when creature dies on first hit")
}

// TestAdventureGameCreatureEncounterProcessor_CreateNextTurnSheet_NoAnotherAdventurerWhenPlayerAttacked
// verifies that the "another adventurer" world event is NOT emitted for a creature
// that the current player damaged on the previous turn.
func TestAdventureGameCreatureEncounterProcessor_CreateNextTurnSheet_NoAnotherAdventurerWhenPlayerAttacked(t *testing.T) {
	th := combatHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameCreatureEncounterProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	creatureInstanceRec, err := th.Data.GetAdventureGameCreatureInstanceRecByRef(harness.GameCreatureInstanceOneRef)
	require.NoError(t, err)

	// Get account user record to supply AccountID/AccountUserID for the turn sheet.
	accountUserRec, err := th.Data.GetAccountUserRecByRef(harness.AccountUserStandardRef)
	require.NoError(t, err)

	// Resolve the actual creature definition name (includes harness UUID prefix).
	creatureDef, err := m.GetAdventureGameCreatureRec(creatureInstanceRec.AdventureGameCreatureID, nil)
	require.NoError(t, err)
	creatureName := creatureDef.Name

	// Advance the game to turn 2 so that a sheet from turn 1 is considered "previous".
	gameInstanceRec.CurrentTurn = 2
	gameInstanceRec, err = m.UpdateGameInstanceRec(gameInstanceRec)
	require.NoError(t, err)

	// Build the previous encounter sheet data: creature at full health (100).
	prevCreatureHealth := 100
	prevSheetData := turnsheet.MonsterEncounterData{
		Creatures: []turnsheet.EncounterCreature{
			{
				CreatureInstanceID: creatureInstanceRec.ID,
				Name:               creatureName,
				Health:             prevCreatureHealth,
				MaxHealth:          creatureDef.MaxHealth,
				IsDead:             false,
			},
		},
	}
	prevSheetBytes, err := json.Marshal(prevSheetData)
	require.NoError(t, err)

	// Create a game_turn_sheet for turn 1 (the "previous" turn).
	prevGameTurnSheet, err := m.CreateGameTurnSheetRec(&game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		GameInstanceID:   sql.NullString{String: gameInstanceRec.ID, Valid: true},
		AccountID:        accountUserRec.AccountID,
		AccountUserID:    accountUserRec.ID,
		TurnNumber:       1,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter,
		SheetData:        json.RawMessage(prevSheetBytes),
		ProcessingStatus: game_record.TurnSheetProcessingStatusProcessed,
	})
	require.NoError(t, err)
	th.Data.AddGameTurnSheetRec(prevGameTurnSheet)
	t.Cleanup(func() {
		_ = m.RemoveGameTurnSheetRec(prevGameTurnSheet.ID)
	})

	// Create the adventure_game_turn_sheet link so the processor can find the previous sheet.
	prevLinkRec, err := m.CreateAdventureGameTurnSheetRec(&adventure_game_record.AdventureGameTurnSheet{
		GameID:                           gameInstanceRec.GameID,
		AdventureGameCharacterInstanceID: charInstanceRec.ID,
		GameTurnSheetID:                  prevGameTurnSheet.ID,
	})
	require.NoError(t, err)
	th.Data.AddAdventureGameTurnSheetRec(prevLinkRec)
	t.Cleanup(func() {
		_ = m.RemoveAdventureGameTurnSheetRec(prevLinkRec.ID)
	})

	// Simulate the player attacking the creature this turn: health drops from 100 → 90.
	creatureInstanceRec.Health = 90
	_, err = m.UpdateAdventureGameCreatureInstanceRec(creatureInstanceRec)
	require.NoError(t, err)

	// Write a combat event using the actual creature name (matches what ProcessTurnSheetResponse
	// writes), so detectMultiPlayerCreatureChanges can correlate it correctly.
	charInstanceRec.LastTurnEvents = []byte("[]")
	_ = turnsheet.AppendTurnEvent(charInstanceRec, turnsheet.TurnEvent{
		Category: turnsheet.TurnEventCategoryCombat,
		Icon:     turnsheet.TurnEventIconCombat,
		Message:  fmt.Sprintf("You attacked the %s for 10 damage.", creatureName),
	})
	_, err = m.UpdateAdventureGameCharacterInstanceRec(charInstanceRec)
	require.NoError(t, err)

	// CreateNextTurnSheet should NOT add an "another adventurer" event.
	sheet, err := proc.CreateNextTurnSheet(context.Background(), gameInstanceRec, charInstanceRec)
	require.NoError(t, err)
	require.NotNil(t, sheet, "sheet should be created — creature is alive at location")

	// Parse the sheet data to retrieve turn events.
	var sheetData turnsheet.MonsterEncounterData
	require.NoError(t, json.Unmarshal(sheet.SheetData, &sheetData))

	for _, e := range sheetData.TurnEvents {
		if e.Category == turnsheet.TurnEventCategoryWorld {
			require.NotContains(t, e.Message, "another adventurer",
				"should not emit 'another adventurer' event when the player dealt the damage")
		}
	}
}

// TestAdventureGameCreatureEncounterProcessor_CreateNextTurnSheet_EventCategoryIsolation
// verifies that CreateNextTurnSheet includes only combat events in the encounter sheet,
// even when movement, inventory, and world events are also present in LastTurnEvents.
func TestAdventureGameCreatureEncounterProcessor_CreateNextTurnSheet_EventCategoryIsolation(t *testing.T) {
	th := combatHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameCreatureEncounterProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	// Seed one event of every non-flee_context category into LastTurnEvents.
	charInstanceRec.LastTurnEvents = []byte("[]")
	for _, ev := range []turnsheet.TurnEvent{
		{Category: turnsheet.TurnEventCategoryCombat, Icon: turnsheet.TurnEventIconCombat, Message: "You attacked the Goblin for 5 damage."},
		{Category: turnsheet.TurnEventCategoryMovement, Icon: turnsheet.TurnEventIconMovement, Message: "You took the Dark Path to the Forest."},
		{Category: turnsheet.TurnEventCategoryInventory, Icon: turnsheet.TurnEventIconInventory, Message: "You picked up an Iron Sword."},
		{Category: turnsheet.TurnEventCategoryWorld, Icon: turnsheet.TurnEventIconWorld, Message: "Another adventurer has slain the Goblin."},
		{Category: turnsheet.TurnEventCategoryFlee, Icon: turnsheet.TurnEventIconFlee, Message: "As you fled, the creature attacked you."},
	} {
		require.NoError(t, turnsheet.AppendTurnEvent(charInstanceRec, ev))
	}
	_, err = m.UpdateAdventureGameCharacterInstanceRec(charInstanceRec)
	require.NoError(t, err)

	sheet, err := proc.CreateNextTurnSheet(context.Background(), gameInstanceRec, charInstanceRec)
	require.NoError(t, err)
	require.NotNil(t, sheet)

	var sheetData turnsheet.MonsterEncounterData
	require.NoError(t, json.Unmarshal(sheet.SheetData, &sheetData))

	events := sheetData.TurnEvents

	// Exactly one combat event should be present.
	var combatCount int
	for _, e := range events {
		if e.Category == turnsheet.TurnEventCategoryCombat {
			combatCount++
		}
	}
	require.Equal(t, 1, combatCount, "encounter sheet should contain exactly the one combat event")

	// No non-combat events should appear on the encounter sheet.
	for _, e := range events {
		require.Equal(t, turnsheet.TurnEventCategoryCombat, e.Category,
			"encounter sheet should only contain combat events, got category %q: %q", e.Category, e.Message)
	}
}

func TestAdventureGameCreatureEncounterProcessor_ProcessTurnSheetResponse_FlatScanData(t *testing.T) {
	// Regression test: verifies the flat form format (action_0, action_0_target etc.)
	// that was stored in the DB before the frontend fix is correctly normalised by GetActions().
	th := combatHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameCreatureEncounterProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)
	charInstanceRec.Health = 100
	charInstanceRec.LastTurnEvents = []byte("[]")
	_, err = m.UpdateAdventureGameCharacterInstanceRec(charInstanceRec)
	require.NoError(t, err)

	creatureInstanceRec, err := th.Data.GetAdventureGameCreatureInstanceRecByRef(harness.GameCreatureInstanceOneRef)
	require.NoError(t, err)
	creatureInstanceRec.Health = 100
	creatureInstanceRec.DiedAtTurn = sql.NullInt64{}
	_, err = m.UpdateAdventureGameCreatureInstanceRec(creatureInstanceRec)
	require.NoError(t, err)

	// Build a turn sheet using the flat key format (the pre-fix frontend format).
	flatScanData := map[string]string{
		"action_0":        "attack",
		"action_0_target": creatureInstanceRec.ID,
		"action_1":        "attack",
		"action_1_target": creatureInstanceRec.ID,
		"action_2":        "attack",
		"action_2_target": creatureInstanceRec.ID,
	}
	scannedBytes, err := json.Marshal(flatScanData)
	require.NoError(t, err)

	turnSheet := &game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		GameInstanceID:   sql.NullString{String: gameInstanceRec.ID, Valid: true},
		TurnNumber:       gameInstanceRec.CurrentTurn,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter,
		SheetData:        json.RawMessage(`{}`),
		ScannedData:      scannedBytes,
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	require.NoError(t, err)

	// 3 attacks × playerDealtPerHit(3) = 9 damage to creature.
	updatedCreature, err := m.GetAdventureGameCreatureInstanceRec(creatureInstanceRec.ID, nil)
	require.NoError(t, err)
	require.Equal(t, 100-9, updatedCreature.Health, "flat format: creature health reduced by 3 attacks")

	// 3 retaliations × creatureDealtPerHit(10) = 30 damage to character.
	require.Equal(t, 100-30, charInstanceRec.Health, "flat format: character health reduced by 3 retaliations")
}
