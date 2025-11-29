package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) validateGameInstanceRecForCreate(rec *game_record.GameInstance) error {
	return validateGameInstanceRec(rec, false)
}

func (m *Domain) validateGameInstanceRecForUpdate(rec *game_record.GameInstance) error {
	return validateGameInstanceRec(rec, true)
}

func validateGameInstanceRec(rec *game_record.GameInstance, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(game_record.FieldGameInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(game_record.FieldGameInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(game_record.FieldGameInstanceGameSubscriptionID, rec.GameSubscriptionID); err != nil {
		return err
	}

	if rec.Status == "" {
		rec.Status = game_record.GameInstanceStatusCreated
	}

	if err := validateGameInstanceStatus(rec.Status); err != nil {
		return err
	}

	if rec.CurrentTurn < 0 {
		return InvalidField(
			game_record.FieldGameInstanceCurrentTurn,
			fmt.Sprintf("%d", rec.CurrentTurn),
			"current_turn must be zero or greater",
		)
	}

	return nil
}

func validateGameInstanceStatus(status string) error {
	switch status {
	case game_record.GameInstanceStatusCreated,
		game_record.GameInstanceStatusStarted,
		game_record.GameInstanceStatusPaused,
		game_record.GameInstanceStatusCompleted,
		game_record.GameInstanceStatusCancelled:
		return nil
	default:
		return coreerror.NewInvalidDataError("invalid game instance status >%s<", status)
	}
}
