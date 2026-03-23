package turnsheet

import (
	"encoding/json"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

// AppendMechaTurnEvent appends an event to the lance instance's
// last_turn_events JSON field.
func AppendMechaTurnEvent(
	lanceInstance *mecha_record.MechaLanceInstance,
	event TurnEvent,
) error {
	events, err := readMechaTurnEvents(lanceInstance)
	if err != nil {
		return err
	}
	events = append(events, event)
	return writeMechaTurnEvents(lanceInstance, events)
}

// ReadMechaTurnEvents returns all turn events without clearing them.
func ReadMechaTurnEvents(
	lanceInstance *mecha_record.MechaLanceInstance,
) ([]TurnEvent, error) {
	return readMechaTurnEvents(lanceInstance)
}

// ReadAndClearMechaTurnEvents reads all turn events and clears the field.
func ReadAndClearMechaTurnEvents(
	lanceInstance *mecha_record.MechaLanceInstance,
) ([]TurnEvent, error) {
	events, err := readMechaTurnEvents(lanceInstance)
	if err != nil {
		return nil, err
	}
	if err := writeMechaTurnEvents(lanceInstance, []TurnEvent{}); err != nil {
		return nil, err
	}
	return events, nil
}

func readMechaTurnEvents(
	lanceInstance *mecha_record.MechaLanceInstance,
) ([]TurnEvent, error) {
	if len(lanceInstance.LastTurnEvents) == 0 {
		return []TurnEvent{}, nil
	}
	var events []TurnEvent
	if err := json.Unmarshal(lanceInstance.LastTurnEvents, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func writeMechaTurnEvents(
	lanceInstance *mecha_record.MechaLanceInstance,
	events []TurnEvent,
) error {
	data, err := json.Marshal(events)
	if err != nil {
		return err
	}
	lanceInstance.LastTurnEvents = json.RawMessage(data)
	return nil
}
