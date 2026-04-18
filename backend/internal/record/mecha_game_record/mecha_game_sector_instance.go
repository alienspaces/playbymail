package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameSectorInstance string = "mecha_game_sector_instance"
)

const (
	FieldMechaGameSectorInstanceID                  string = "id"
	FieldMechaGameSectorInstanceGameID              string = "game_id"
	FieldMechaGameSectorInstanceGameInstanceID      string = "game_instance_id"
	FieldMechaGameSectorInstanceMechaGameSectorID string = "mecha_game_sector_id"
	FieldMechaGameSectorInstanceCreatedAt           string = "created_at"
	FieldMechaGameSectorInstanceUpdatedAt           string = "updated_at"
	FieldMechaGameSectorInstanceDeletedAt           string = "deleted_at"
)

type MechaGameSectorInstance struct {
	record.Record
	GameID               string `db:"game_id"`
	GameInstanceID       string `db:"game_instance_id"`
	MechaGameSectorID  string `db:"mecha_game_sector_id"`
}

func (r *MechaGameSectorInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameSectorInstanceGameID] = r.GameID
	args[FieldMechaGameSectorInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechaGameSectorInstanceMechaGameSectorID] = r.MechaGameSectorID
	return args
}
