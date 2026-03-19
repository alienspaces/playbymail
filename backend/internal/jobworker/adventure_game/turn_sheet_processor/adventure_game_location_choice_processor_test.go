package turn_sheet_processor_test

import (
	"context"
	"database/sql"
	"encoding/json"
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

// objectHarness returns a default test harness for location choice / object interaction tests.
func objectHarness(t *testing.T) *harness.Testing {
	t.Helper()

	cfg, err := config.Parse()
	require.NoError(t, err)

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err)

	th, err := harness.NewTesting(cfg, l, s, j, scanner, harness.DefaultDataConfig())
	require.NoError(t, err)

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")

	t.Cleanup(func() {
		require.NoError(t, th.Teardown(), "Test data teardown returns without error")
	})

	return th
}

// buildLocationChoiceTurnSheet constructs a minimal GameTurnSheet with the given scanned data.
func buildLocationChoiceTurnSheet(t *testing.T, gameInstanceRec *game_record.GameInstance, scanData turnsheet.LocationChoiceScanData) *game_record.GameTurnSheet {
	t.Helper()

	scannedBytes, err := json.Marshal(scanData)
	require.NoError(t, err)

	return &game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		GameInstanceID:   sql.NullString{String: gameInstanceRec.ID, Valid: true},
		TurnNumber:       gameInstanceRec.CurrentTurn,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeLocationChoice,
		SheetData:        json.RawMessage(`{}`),
		ScannedData:      scannedBytes,
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
}

// TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_Info verifies that an
// info effect returns without modifying any persistent state.
func TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_Info(t *testing.T) {
	th := objectHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameLocationChoiceProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	// Ancient Shrine with action=inspect → effect_type=info
	objectInstance, err := th.Data.GetAdventureGameLocationObjectInstanceByObjectRef(harness.GameLocationObjectOneRef)
	require.NoError(t, err)

	objectInstance.CurrentState = "intact"
	_, err = m.UpdateAdventureGameLocationObjectInstanceRec(objectInstance)
	require.NoError(t, err)

	scanData := turnsheet.LocationChoiceScanData{
		ObjectChoice: objectInstance.ID + ":" + adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
	}
	turnSheet := buildLocationChoiceTurnSheet(t, gameInstanceRec, scanData)

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	require.NoError(t, err)

	// Object state should be unchanged (info effects don't change state)
	refreshed, err := m.GetAdventureGameLocationObjectInstanceRec(objectInstance.ID, nil)
	require.NoError(t, err)
	require.Equal(t, "intact", refreshed.CurrentState, "info effect should not change object state")
}

// TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_ChangeState verifies that a
// change_state effect transitions the object to the result state.
func TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_ChangeState(t *testing.T) {
	th := objectHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameLocationChoiceProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	// Ancient Shrine with action=touch → effect_type=change_state (intact→activated)
	objectInstance, err := th.Data.GetAdventureGameLocationObjectInstanceByObjectRef(harness.GameLocationObjectOneRef)
	require.NoError(t, err)

	objectInstance.CurrentState = "intact"
	_, err = m.UpdateAdventureGameLocationObjectInstanceRec(objectInstance)
	require.NoError(t, err)

	scanData := turnsheet.LocationChoiceScanData{
		ObjectChoice: objectInstance.ID + ":" + adventure_game_record.AdventureGameLocationObjectEffectActionTypeTouch,
	}
	turnSheet := buildLocationChoiceTurnSheet(t, gameInstanceRec, scanData)

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	require.NoError(t, err)

	refreshed, err := m.GetAdventureGameLocationObjectInstanceRec(objectInstance.ID, nil)
	require.NoError(t, err)
	require.Equal(t, "activated", refreshed.CurrentState, "change_state effect should transition object to activated")
}

// TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_GiveItem verifies that a
// give_item effect creates a new item instance in the character's inventory.
func TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_GiveItem(t *testing.T) {
	th := objectHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameLocationChoiceProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	// The give_item effect on the shrine requires state=activated and action=search.
	objectInstance, err := th.Data.GetAdventureGameLocationObjectInstanceByObjectRef(harness.GameLocationObjectOneRef)
	require.NoError(t, err)

	objectInstance.CurrentState = "activated"
	_, err = m.UpdateAdventureGameLocationObjectInstanceRec(objectInstance)
	require.NoError(t, err)

	// Count items before
	itemRec, err := th.Data.GetAdventureGameItemRecByRef(harness.GameItemTwoRef)
	require.NoError(t, err)

	beforeInstances, err := m.GetManyAdventureGameItemInstanceRecs(nil)
	require.NoError(t, err)
	beforeCount := 0
	for _, inst := range beforeInstances {
		if inst.AdventureGameItemID == itemRec.ID && inst.AdventureGameCharacterInstanceID.String == charInstanceRec.ID {
			beforeCount++
		}
	}

	scanData := turnsheet.LocationChoiceScanData{
		ObjectChoice: objectInstance.ID + ":" + adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
	}
	turnSheet := buildLocationChoiceTurnSheet(t, gameInstanceRec, scanData)

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	require.NoError(t, err)

	afterInstances, err := m.GetManyAdventureGameItemInstanceRecs(nil)
	require.NoError(t, err)
	afterCount := 0
	for _, inst := range afterInstances {
		if inst.AdventureGameItemID == itemRec.ID && inst.AdventureGameCharacterInstanceID.String == charInstanceRec.ID {
			afterCount++
		}
	}

	require.Equal(t, beforeCount+1, afterCount, "give_item effect should add one item instance to character inventory")
}

// TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_Damage verifies that a
// damage effect reduces the character's health.
func TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_Damage(t *testing.T) {
	th := objectHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameLocationChoiceProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	// Ancient Shrine with action=break → effect_type=damage (min=5 max=5 so always 5)
	objectInstance, err := th.Data.GetAdventureGameLocationObjectInstanceByObjectRef(harness.GameLocationObjectOneRef)
	require.NoError(t, err)

	objectInstance.CurrentState = "intact"
	_, err = m.UpdateAdventureGameLocationObjectInstanceRec(objectInstance)
	require.NoError(t, err)

	const startingHealth = 100
	charInstanceRec.Health = startingHealth
	charInstanceRec.LastTurnEvents = []byte("[]")
	_, err = m.UpdateAdventureGameCharacterInstanceRec(charInstanceRec)
	require.NoError(t, err)

	scanData := turnsheet.LocationChoiceScanData{
		ObjectChoice: objectInstance.ID + ":" + adventure_game_record.AdventureGameLocationObjectEffectActionTypeBreak,
	}
	turnSheet := buildLocationChoiceTurnSheet(t, gameInstanceRec, scanData)

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	require.NoError(t, err)

	// The damage effect always deals exactly 5 (min=max=5)
	require.Equal(t, startingHealth-5, charInstanceRec.Health, "damage effect should reduce character health by 5")
}

// TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_Heal verifies that a
// heal effect increases the character's health.
func TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_Heal(t *testing.T) {
	th := objectHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameLocationChoiceProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	// Ancient Shrine with action=touch + required_state=activated → effect_type=heal (min=max=10)
	objectInstance, err := th.Data.GetAdventureGameLocationObjectInstanceByObjectRef(harness.GameLocationObjectOneRef)
	require.NoError(t, err)

	objectInstance.CurrentState = "activated"
	_, err = m.UpdateAdventureGameLocationObjectInstanceRec(objectInstance)
	require.NoError(t, err)

	const startingHealth = 50
	charInstanceRec.Health = startingHealth
	charInstanceRec.LastTurnEvents = []byte("[]")
	_, err = m.UpdateAdventureGameCharacterInstanceRec(charInstanceRec)
	require.NoError(t, err)

	scanData := turnsheet.LocationChoiceScanData{
		ObjectChoice: objectInstance.ID + ":" + adventure_game_record.AdventureGameLocationObjectEffectActionTypeTouch,
	}
	turnSheet := buildLocationChoiceTurnSheet(t, gameInstanceRec, scanData)

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	require.NoError(t, err)

	// The heal effect always heals exactly 10 (min=max=10) when state is activated
	require.Equal(t, startingHealth+10, charInstanceRec.Health, "heal effect should increase character health by 10")
}

// TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_RevealObject verifies that a
// reveal_object effect makes the target object instance visible.
func TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_RevealObject(t *testing.T) {
	th := objectHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameLocationChoiceProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	// Ancient Shrine with action=use → effect_type=reveal_object targets Hidden Passage
	shrineInstance, err := th.Data.GetAdventureGameLocationObjectInstanceByObjectRef(harness.GameLocationObjectOneRef)
	require.NoError(t, err)

	shrineInstance.CurrentState = "intact"
	_, err = m.UpdateAdventureGameLocationObjectInstanceRec(shrineInstance)
	require.NoError(t, err)

	// Ensure the Hidden Passage starts hidden
	passageInstance, err := th.Data.GetAdventureGameLocationObjectInstanceByObjectRef(harness.GameLocationObjectTwoRef)
	require.NoError(t, err)

	passageInstance.IsVisible = false
	_, err = m.UpdateAdventureGameLocationObjectInstanceRec(passageInstance)
	require.NoError(t, err)

	scanData := turnsheet.LocationChoiceScanData{
		ObjectChoice: shrineInstance.ID + ":" + adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
	}
	turnSheet := buildLocationChoiceTurnSheet(t, gameInstanceRec, scanData)

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	require.NoError(t, err)

	refreshedPassage, err := m.GetAdventureGameLocationObjectInstanceRec(passageInstance.ID, nil)
	require.NoError(t, err)
	require.True(t, refreshedPassage.IsVisible, "reveal_object effect should make the Hidden Passage visible")
}

// TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_WrongState verifies that an
// effect whose required_state does not match the object's current state has no effect.
func TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_WrongState(t *testing.T) {
	th := objectHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameLocationChoiceProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	// The give_item effect on the shrine requires state=activated; set state to intact so it won't apply.
	objectInstance, err := th.Data.GetAdventureGameLocationObjectInstanceByObjectRef(harness.GameLocationObjectOneRef)
	require.NoError(t, err)

	objectInstance.CurrentState = "intact"
	_, err = m.UpdateAdventureGameLocationObjectInstanceRec(objectInstance)
	require.NoError(t, err)

	// Reset character events
	charInstanceRec.LastTurnEvents = []byte("[]")
	_, err = m.UpdateAdventureGameCharacterInstanceRec(charInstanceRec)
	require.NoError(t, err)

	scanData := turnsheet.LocationChoiceScanData{
		// search is only effective when state=activated; current state=intact → no matching effects
		ObjectChoice: objectInstance.ID + ":" + adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
	}
	turnSheet := buildLocationChoiceTurnSheet(t, gameInstanceRec, scanData)

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	// Should succeed without error — no matching effects is treated as a no-op
	require.NoError(t, err)

	// Object state should remain intact
	refreshed, err := m.GetAdventureGameLocationObjectInstanceRec(objectInstance.ID, nil)
	require.NoError(t, err)
	require.Equal(t, "intact", refreshed.CurrentState, "state should remain intact when no effects match")
}

// TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_NoObjectChoiceIsLocationChoice
// verifies that a turn sheet with a location choice (not an object choice) still applies
// the standard location movement logic.
func TestAdventureGameLocationChoiceProcessor_ProcessObjectChoice_NoObjectChoiceIsLocationChoice(t *testing.T) {
	th := objectHarness(t)
	m := th.Domain.(*domain.Domain)

	proc, err := turn_sheet_processor.NewAdventureGameLocationChoiceProcessor(m.Log, m)
	require.NoError(t, err)

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err)

	charInstanceRec, err := th.Data.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
	require.NoError(t, err)

	locationInstanceRec, err := th.Data.GetAdventureGameLocationInstanceRecByRef(harness.GameLocationInstanceTwoRef)
	require.NoError(t, err)

	// Place character at location one
	locationOneInstanceRec, err := th.Data.GetAdventureGameLocationInstanceRecByRef(harness.GameLocationInstanceOneRef)
	require.NoError(t, err)
	charInstanceRec.AdventureGameLocationInstanceID = locationOneInstanceRec.ID
	charInstanceRec.LastTurnEvents = []byte("[]")
	_, err = m.UpdateAdventureGameCharacterInstanceRec(charInstanceRec)
	require.NoError(t, err)

	// Build a location choice (not object choice)
	scanData := turnsheet.LocationChoiceScanData{
		Choices: []string{locationInstanceRec.ID},
	}
	turnSheet := buildLocationChoiceTurnSheet(t, gameInstanceRec, scanData)
	// SheetData must include the location option for validation
	sheetData := turnsheet.LocationChoiceData{
		LocationOptions: []turnsheet.LocationOption{
			{LocationID: locationInstanceRec.ID, LocationLinkName: "Forest Path", IsLocked: false},
		},
	}
	sheetDataBytes, err := json.Marshal(sheetData)
	require.NoError(t, err)
	turnSheet.SheetData = sheetDataBytes

	err = proc.ProcessTurnSheetResponse(context.Background(), gameInstanceRec, charInstanceRec, turnSheet)
	require.NoError(t, err)

	require.Equal(t, locationInstanceRec.ID, charInstanceRec.AdventureGameLocationInstanceID,
		"character should have moved to the chosen location")
}
