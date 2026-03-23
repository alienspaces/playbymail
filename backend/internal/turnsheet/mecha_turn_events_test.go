package turnsheet_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

func TestAppendMechaTurnEvent(t *testing.T) {
	t.Parallel()

	t.Run("appends to empty lance instance", func(t *testing.T) {
		t.Parallel()
		lance := &mecha_record.MechaLanceInstance{}

		err := turnsheet.AppendMechaTurnEvent(lance, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryMovement,
			Icon:     turnsheet.TurnEventIconMovement,
			Message:  "Hammer moved to Northern Ridge.",
		})
		require.NoError(t, err)

		events, err := turnsheet.ReadMechaTurnEvents(lance)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, turnsheet.TurnEventCategoryMovement, events[0].Category)
		require.Equal(t, "Hammer moved to Northern Ridge.", events[0].Message)
	})

	t.Run("appends to lance instance with existing events", func(t *testing.T) {
		t.Parallel()
		lance := &mecha_record.MechaLanceInstance{}

		require.NoError(t, turnsheet.AppendMechaTurnEvent(lance, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryMovement,
			Message:  "First event.",
		}))
		require.NoError(t, turnsheet.AppendMechaTurnEvent(lance, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryCombat,
			Message:  "Second event.",
		}))
		require.NoError(t, turnsheet.AppendMechaTurnEvent(lance, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategorySystem,
			Message:  "Third event.",
		}))

		events, err := turnsheet.ReadMechaTurnEvents(lance)
		require.NoError(t, err)
		require.Len(t, events, 3)
		require.Equal(t, "First event.", events[0].Message)
		require.Equal(t, "Second event.", events[1].Message)
		require.Equal(t, "Third event.", events[2].Message)
	})
}

func TestReadMechaTurnEvents(t *testing.T) {
	t.Parallel()

	t.Run("returns empty slice when LastTurnEvents is nil", func(t *testing.T) {
		t.Parallel()
		lance := &mecha_record.MechaLanceInstance{}

		events, err := turnsheet.ReadMechaTurnEvents(lance)
		require.NoError(t, err)
		require.NotNil(t, events)
		require.Len(t, events, 0)
	})

	t.Run("returns events from valid JSON", func(t *testing.T) {
		t.Parallel()
		raw, err := json.Marshal([]turnsheet.TurnEvent{
			{Category: turnsheet.TurnEventCategoryCombat, Icon: turnsheet.TurnEventIconCombat, Message: "Hit for 5."},
		})
		require.NoError(t, err)

		lance := &mecha_record.MechaLanceInstance{
			LastTurnEvents: json.RawMessage(raw),
		}

		events, err := turnsheet.ReadMechaTurnEvents(lance)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, turnsheet.TurnEventCategoryCombat, events[0].Category)
		require.Equal(t, "Hit for 5.", events[0].Message)
	})

	t.Run("does not modify the lance instance", func(t *testing.T) {
		t.Parallel()
		raw, err := json.Marshal([]turnsheet.TurnEvent{
			{Category: turnsheet.TurnEventCategorySystem, Message: "System event."},
		})
		require.NoError(t, err)

		lance := &mecha_record.MechaLanceInstance{
			LastTurnEvents: json.RawMessage(raw),
		}
		originalRaw := string(lance.LastTurnEvents)

		_, err = turnsheet.ReadMechaTurnEvents(lance)
		require.NoError(t, err)
		require.Equal(t, originalRaw, string(lance.LastTurnEvents))
	})
}

func TestReadAndClearMechaTurnEvents(t *testing.T) {
	t.Parallel()

	t.Run("returns events and clears the field", func(t *testing.T) {
		t.Parallel()
		lance := &mecha_record.MechaLanceInstance{}

		require.NoError(t, turnsheet.AppendMechaTurnEvent(lance, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryMovement,
			Message:  "Moved to sector A.",
		}))
		require.NoError(t, turnsheet.AppendMechaTurnEvent(lance, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryCombat,
			Message:  "Fired at enemy.",
		}))

		events, err := turnsheet.ReadAndClearMechaTurnEvents(lance)
		require.NoError(t, err)
		require.Len(t, events, 2)
		require.Equal(t, "Moved to sector A.", events[0].Message)

		// Field should now be empty
		remaining, err := turnsheet.ReadMechaTurnEvents(lance)
		require.NoError(t, err)
		require.Len(t, remaining, 0)
	})

	t.Run("returns empty slice and leaves field empty when no events", func(t *testing.T) {
		t.Parallel()
		lance := &mecha_record.MechaLanceInstance{}

		events, err := turnsheet.ReadAndClearMechaTurnEvents(lance)
		require.NoError(t, err)
		require.NotNil(t, events)
		require.Len(t, events, 0)
	})
}
