package domain

import (
	"context"
	"errors"

	"gitlab.com/alienspaces/playbymail/internal/record"
)

func ValidateGameCreatureInstanceFields(ctx context.Context, rec *record.GameCreatureInstance) error {
	if rec.ID == "" {
		return errors.New("id is required")
	}
	if rec.GameID == "" {
		return errors.New("game_id is required")
	}
	if rec.GameCreatureID == "" {
		return errors.New("game_creature_id is required")
	}
	if rec.GameInstanceID == "" {
		return errors.New("game_instance_id is required")
	}
	if rec.GameLocationInstanceID == "" {
		return errors.New("game_location_instance_id is required")
	}
	return nil
}
