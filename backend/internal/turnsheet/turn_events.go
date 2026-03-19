package turnsheet

import (
	"encoding/json"

	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// AppendTurnEvent appends a TurnEvent to the character instance's last_turn_events JSON field.
func AppendTurnEvent(characterInstance *adventure_game_record.AdventureGameCharacterInstance, event TurnEvent) error {
	events, err := readTurnEvents(characterInstance)
	if err != nil {
		return err
	}
	events = append(events, event)
	return writeTurnEvents(characterInstance, events)
}

// ReadAndClearTurnEvents reads all turn events from the character instance and clears the field.
func ReadAndClearTurnEvents(characterInstance *adventure_game_record.AdventureGameCharacterInstance) ([]TurnEvent, error) {
	events, err := readTurnEvents(characterInstance)
	if err != nil {
		return nil, err
	}
	if err := writeTurnEvents(characterInstance, []TurnEvent{}); err != nil {
		return nil, err
	}
	return events, nil
}

// ReadTurnEvents reads all turn events without clearing them.
func ReadTurnEvents(characterInstance *adventure_game_record.AdventureGameCharacterInstance) ([]TurnEvent, error) {
	return readTurnEvents(characterInstance)
}

func readTurnEvents(characterInstance *adventure_game_record.AdventureGameCharacterInstance) ([]TurnEvent, error) {
	if len(characterInstance.LastTurnEvents) == 0 {
		return []TurnEvent{}, nil
	}
	var events []TurnEvent
	if err := json.Unmarshal(characterInstance.LastTurnEvents, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func writeTurnEvents(characterInstance *adventure_game_record.AdventureGameCharacterInstance, events []TurnEvent) error {
	data, err := json.Marshal(events)
	if err != nil {
		return err
	}
	characterInstance.LastTurnEvents = json.RawMessage(data)
	return nil
}

// FilterTurnEventsByCategory returns only events matching the given categories.
func FilterTurnEventsByCategory(events []TurnEvent, categories ...string) []TurnEvent {
	catSet := make(map[string]bool, len(categories))
	for _, c := range categories {
		catSet[c] = true
	}
	var result []TurnEvent
	for _, e := range events {
		if catSet[e.Category] {
			result = append(result, e)
		}
	}
	return result
}

// ExcludeTurnEventsByCategory returns events NOT matching the given categories.
func ExcludeTurnEventsByCategory(events []TurnEvent, categories ...string) []TurnEvent {
	catSet := make(map[string]bool, len(categories))
	for _, c := range categories {
		catSet[c] = true
	}
	var result []TurnEvent
	for _, e := range events {
		if !catSet[e.Category] {
			result = append(result, e)
		}
	}
	return result
}
