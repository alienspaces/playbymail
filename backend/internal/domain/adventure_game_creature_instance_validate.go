package domain

import (
	"context"
	"errors"

	"gitlab.com/alienspaces/playbymail/internal/record"
)

func ValidateAdventureGameCreatureInstanceFields(ctx context.Context, rec *record.AdventureGameCreatureInstance) error {
	if rec.ID == "" {
		return errors.New("id is required")
	}
	if rec.GameID == "" {
		return errors.New("game_id is required")
	}
	if rec.AdventureGameCreatureID == "" {
		return errors.New("adventure_game_creature_id is required")
	}
	if rec.AdventureGameInstanceID == "" {
		return errors.New("game_instance_id is required")
	}
	if rec.AdventureGameLocationInstanceID == "" {
		return errors.New("adventure_game_location_instance_id is required")
	}
	return nil
}
