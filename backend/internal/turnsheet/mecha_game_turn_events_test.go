package turnsheet_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

func TestAppendMechaGameTurnEvent(t *testing.T) {
	t.Parallel()

	t.Run("appends to empty squad instance", func(t *testing.T) {
		t.Parallel()
		squad := &mecha_game_record.MechaGameSquadInstance{}

		err := turnsheet.AppendMechaGameTurnEvent(squad, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryMovement,
			Icon:     turnsheet.TurnEventIconMovement,
			Message:  "Hammer moved to Northern Ridge.",
		})
		require.NoError(t, err)

		events, err := turnsheet.ReadMechaGameTurnEvents(squad)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, turnsheet.TurnEventCategoryMovement, events[0].Category)
		require.Equal(t, "Hammer moved to Northern Ridge.", events[0].Message)
	})

	t.Run("appends to squad instance with existing events", func(t *testing.T) {
		t.Parallel()
		squad := &mecha_game_record.MechaGameSquadInstance{}

		require.NoError(t, turnsheet.AppendMechaGameTurnEvent(squad, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryMovement,
			Message:  "First event.",
		}))
		require.NoError(t, turnsheet.AppendMechaGameTurnEvent(squad, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryCombat,
			Message:  "Second event.",
		}))
		require.NoError(t, turnsheet.AppendMechaGameTurnEvent(squad, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategorySystem,
			Message:  "Third event.",
		}))

		events, err := turnsheet.ReadMechaGameTurnEvents(squad)
		require.NoError(t, err)
		require.Len(t, events, 3)
		require.Equal(t, "First event.", events[0].Message)
		require.Equal(t, "Second event.", events[1].Message)
		require.Equal(t, "Third event.", events[2].Message)
	})
}

func TestReadMechaGameTurnEvents(t *testing.T) {
	t.Parallel()

	t.Run("returns empty slice when LastTurnEvents is nil", func(t *testing.T) {
		t.Parallel()
		squad := &mecha_game_record.MechaGameSquadInstance{}

		events, err := turnsheet.ReadMechaGameTurnEvents(squad)
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

		squad := &mecha_game_record.MechaGameSquadInstance{
			LastTurnEvents: json.RawMessage(raw),
		}

		events, err := turnsheet.ReadMechaGameTurnEvents(squad)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, turnsheet.TurnEventCategoryCombat, events[0].Category)
		require.Equal(t, "Hit for 5.", events[0].Message)
	})

	t.Run("does not modify the squad instance", func(t *testing.T) {
		t.Parallel()
		raw, err := json.Marshal([]turnsheet.TurnEvent{
			{Category: turnsheet.TurnEventCategorySystem, Message: "System event."},
		})
		require.NoError(t, err)

		squad := &mecha_game_record.MechaGameSquadInstance{
			LastTurnEvents: json.RawMessage(raw),
		}
		originalRaw := string(squad.LastTurnEvents)

		_, err = turnsheet.ReadMechaGameTurnEvents(squad)
		require.NoError(t, err)
		require.Equal(t, originalRaw, string(squad.LastTurnEvents))
	})
}

func TestReadAndClearMechaGameTurnEvents(t *testing.T) {
	t.Parallel()

	t.Run("returns events and clears the field", func(t *testing.T) {
		t.Parallel()
		squad := &mecha_game_record.MechaGameSquadInstance{}

		require.NoError(t, turnsheet.AppendMechaGameTurnEvent(squad, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryMovement,
			Message:  "Moved to sector A.",
		}))
		require.NoError(t, turnsheet.AppendMechaGameTurnEvent(squad, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryCombat,
			Message:  "Fired at enemy.",
		}))

		events, err := turnsheet.ReadAndClearMechaGameTurnEvents(squad)
		require.NoError(t, err)
		require.Len(t, events, 2)
		require.Equal(t, "Moved to sector A.", events[0].Message)

		// Field should now be empty
		remaining, err := turnsheet.ReadMechaGameTurnEvents(squad)
		require.NoError(t, err)
		require.Len(t, remaining, 0)
	})

	t.Run("returns empty slice and leaves field empty when no events", func(t *testing.T) {
		t.Parallel()
		squad := &mecha_game_record.MechaGameSquadInstance{}

		events, err := turnsheet.ReadAndClearMechaGameTurnEvents(squad)
		require.NoError(t, err)
		require.NotNil(t, events)
		require.Len(t, events, 0)
	})
}
