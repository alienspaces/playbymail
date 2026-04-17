package turnsheet

import (
	"encoding/json"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

// AppendMechaTurnEvent appends an event to the squad instance's
// last_turn_events JSON field.
func AppendMechaTurnEvent(
	squadInstance *mecha_record.MechaSquadInstance,
	event TurnEvent,
) error {
	events, err := readMechaTurnEvents(squadInstance)
	if err != nil {
		return err
	}
	events = append(events, event)
	return writeMechaTurnEvents(squadInstance, events)
}

// ReadMechaTurnEvents returns all turn events without clearing them.
func ReadMechaTurnEvents(
	squadInstance *mecha_record.MechaSquadInstance,
) ([]TurnEvent, error) {
	return readMechaTurnEvents(squadInstance)
}

// ReadAndClearMechaTurnEvents reads all turn events and clears the field.
func ReadAndClearMechaTurnEvents(
	squadInstance *mecha_record.MechaSquadInstance,
) ([]TurnEvent, error) {
	events, err := readMechaTurnEvents(squadInstance)
	if err != nil {
		return nil, err
	}
	if err := writeMechaTurnEvents(squadInstance, []TurnEvent{}); err != nil {
		return nil, err
	}
	return events, nil
}

func readMechaTurnEvents(
	squadInstance *mecha_record.MechaSquadInstance,
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

func writeMechaTurnEvents(
	squadInstance *mecha_record.MechaSquadInstance,
	events []TurnEvent,
) error {
	data, err := json.Marshal(events)
	if err != nil {
		return err
	}
	squadInstance.LastTurnEvents = json.RawMessage(data)
	return nil
}
