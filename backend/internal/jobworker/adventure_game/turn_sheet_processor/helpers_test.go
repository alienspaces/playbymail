package turn_sheet_processor_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/adventure_game/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// characterWithEvents returns a minimal AdventureGameCharacterInstance whose
// LastTurnEvents field is populated with the given events.
func characterWithEvents(t *testing.T, events []turnsheet.TurnEvent) *adventure_game_record.AdventureGameCharacterInstance {
	t.Helper()
	data, err := json.Marshal(events)
	require.NoError(t, err)
	return &adventure_game_record.AdventureGameCharacterInstance{
		LastTurnEvents: json.RawMessage(data),
	}
}

func TestReadTurnEventsForCategories(t *testing.T) {
	l := log.NewDefaultLogger()

	combatEvent := turnsheet.TurnEvent{Category: turnsheet.TurnEventCategoryCombat, Icon: turnsheet.TurnEventIconCombat, Message: "You attacked the Goblin for 5 damage."}
	movementEvent := turnsheet.TurnEvent{Category: turnsheet.TurnEventCategoryMovement, Icon: turnsheet.TurnEventIconMovement, Message: "You took the Dark Path to the Forest."}
	inventoryEvent := turnsheet.TurnEvent{Category: turnsheet.TurnEventCategoryInventory, Icon: turnsheet.TurnEventIconInventory, Message: "You picked up an Iron Sword."}
	worldEvent := turnsheet.TurnEvent{Category: turnsheet.TurnEventCategoryWorld, Icon: turnsheet.TurnEventIconWorld, Message: "Another adventurer has slain the Goblin."}
	fleeEvent := turnsheet.TurnEvent{Category: turnsheet.TurnEventCategoryFlee, Icon: turnsheet.TurnEventIconFlee, Message: "As you fled, the Goblin attacks you for 3 damage."}

	allEvents := []turnsheet.TurnEvent{combatEvent, movementEvent, inventoryEvent, worldEvent, fleeEvent}

	tests := []struct {
		name           string
		events         []turnsheet.TurnEvent
		categories     []string
		wantCategories []string
		wantLen        int
	}{
		{
			name:           "given only combat category requested then only combat events are returned",
			events:         allEvents,
			categories:     []string{turnsheet.TurnEventCategoryCombat},
			wantCategories: []string{turnsheet.TurnEventCategoryCombat},
			wantLen:        1,
		},
		{
			name:           "given movement flee world categories requested then only those events are returned",
			events:         allEvents,
			categories:     []string{turnsheet.TurnEventCategoryMovement, turnsheet.TurnEventCategoryFlee, turnsheet.TurnEventCategoryWorld},
			wantCategories: []string{turnsheet.TurnEventCategoryMovement, turnsheet.TurnEventCategoryFlee, turnsheet.TurnEventCategoryWorld},
			wantLen:        3,
		},
		{
			name:           "given inventory category requested then only inventory events are returned",
			events:         allEvents,
			categories:     []string{turnsheet.TurnEventCategoryInventory},
			wantCategories: []string{turnsheet.TurnEventCategoryInventory},
			wantLen:        1,
		},
		{
			name:       "given category not present in events then empty slice is returned",
			events:     []turnsheet.TurnEvent{movementEvent},
			categories: []string{turnsheet.TurnEventCategoryCombat},
			wantLen:    0,
		},
		{
			name:       "given no events then empty slice is returned",
			events:     []turnsheet.TurnEvent{},
			categories: []string{turnsheet.TurnEventCategoryCombat},
			wantLen:    0,
		},
		{
			name:    "given no categories requested then empty slice is returned",
			events:  allEvents,
			wantLen: 0,
		},
		{
			name:           "given combat category requested then movement inventory world and flee events are excluded",
			events:         allEvents,
			categories:     []string{turnsheet.TurnEventCategoryCombat},
			wantCategories: []string{turnsheet.TurnEventCategoryCombat},
			wantLen:        1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charRec := characterWithEvents(t, tt.events)

			got := turn_sheet_processor.ReadTurnEventsForCategories(l, nil, charRec, tt.categories...)

			require.Len(t, got, tt.wantLen)

			// Verify every returned event is in the wanted category set.
			if len(tt.wantCategories) > 0 {
				wantSet := make(map[string]bool, len(tt.wantCategories))
				for _, c := range tt.wantCategories {
					wantSet[c] = true
				}
				for _, e := range got {
					require.True(t, wantSet[e.Category],
						"returned event category %q should be in wanted set %v", e.Category, tt.wantCategories)
				}
			}
		})
	}
}
