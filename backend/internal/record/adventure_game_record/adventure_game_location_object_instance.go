package adventure_game_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameLocationObjectInstance = "adventure_game_location_object_instance"

const (
	FieldAdventureGameLocationObjectInstanceID                              = "id"
	FieldAdventureGameLocationObjectInstanceGameID                          = "game_id"
	FieldAdventureGameLocationObjectInstanceGameInstanceID                  = "game_instance_id"
	FieldAdventureGameLocationObjectInstanceAdventureGameLocationObjectID   = "adventure_game_location_object_id"
	FieldAdventureGameLocationObjectInstanceAdventureGameLocationInstanceID = "adventure_game_location_instance_id"
	FieldAdventureGameLocationObjectInstanceCurrentState                    = "current_state"
	FieldAdventureGameLocationObjectInstanceIsVisible                       = "is_visible"
)

type AdventureGameLocationObjectInstance struct {
	record.Record
	GameID                          string `db:"game_id"`
	GameInstanceID                  string `db:"game_instance_id"`
	AdventureGameLocationObjectID   string `db:"adventure_game_location_object_id"`
	AdventureGameLocationInstanceID string `db:"adventure_game_location_instance_id"`
	CurrentState                    string `db:"current_state"`
	IsVisible                       bool   `db:"is_visible"`
}

func (r *AdventureGameLocationObjectInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationObjectInstanceGameID] = r.GameID
	args[FieldAdventureGameLocationObjectInstanceGameInstanceID] = r.GameInstanceID
	args[FieldAdventureGameLocationObjectInstanceAdventureGameLocationObjectID] = r.AdventureGameLocationObjectID
	args[FieldAdventureGameLocationObjectInstanceAdventureGameLocationInstanceID] = r.AdventureGameLocationInstanceID
	args[FieldAdventureGameLocationObjectInstanceCurrentState] = r.CurrentState
	args[FieldAdventureGameLocationObjectInstanceIsVisible] = r.IsVisible
	return args
}
