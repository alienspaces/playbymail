package turnsheet

import (
	"encoding/json"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

// AppendMechaGameTurnEvent appends an event to the squad instance's
// last_turn_events JSON field.
func AppendMechaGameTurnEvent(
	squadInstance *mecha_game_record.MechaGameSquadInstance,
	event TurnEvent,
) error {
	events, err := readMechaGameTurnEvents(squadInstance)
	if err != nil {
		return err
	}
	events = append(events, event)
	return writeMechaGameTurnEvents(squadInstance, events)
}

// ReadMechaGameTurnEvents returns all turn events without clearing them.
func ReadMechaGameTurnEvents(
	squadInstance *mecha_game_record.MechaGameSquadInstance,
) ([]TurnEvent, error) {
	return readMechaGameTurnEvents(squadInstance)
}

// ReadAndClearMechaGameTurnEvents reads all turn events and clears the field.
func ReadAndClearMechaGameTurnEvents(
	squadInstance *mecha_game_record.MechaGameSquadInstance,
) ([]TurnEvent, error) {
	events, err := readMechaGameTurnEvents(squadInstance)
	if err != nil {
		return nil, err
	}
	if err := writeMechaGameTurnEvents(squadInstance, []TurnEvent{}); err != nil {
		return nil, err
	}
	return events, nil
}

func readMechaGameTurnEvents(
	squadInstance *mecha_game_record.MechaGameSquadInstance,
) ([]TurnEvent, error) {
	if len(squadInstance.LastTurnEvents) == 0 {
		return []TurnEvent{}, nil
	}
	var events []TurnEvent
	if err := json.Unmarshal(squadInstance.LastTurnEvents, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func writeMechaGameTurnEvents(
	squadInstance *mecha_game_record.MechaGameSquadInstance,
	events []TurnEvent,
) error {
	data, err := json.Marshal(events)
	if err != nil {
		return err
	}
	squadInstance.LastTurnEvents = json.RawMessage(data)
	return nil
}
