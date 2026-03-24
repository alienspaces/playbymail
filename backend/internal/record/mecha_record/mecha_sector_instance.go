package mecha_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaSectorInstance string = "mecha_sector_instance"
)

const (
	FieldMechaSectorInstanceID                  string = "id"
	FieldMechaSectorInstanceGameID              string = "game_id"
	FieldMechaSectorInstanceGameInstanceID      string = "game_instance_id"
	FieldMechaSectorInstanceMechaSectorID string = "mecha_sector_id"
	FieldMechaSectorInstanceCreatedAt           string = "created_at"
	FieldMechaSectorInstanceUpdatedAt           string = "updated_at"
	FieldMechaSectorInstanceDeletedAt           string = "deleted_at"
)

type MechaSectorInstance struct {
	record.Record
	GameID               string `db:"game_id"`
	GameInstanceID       string `db:"game_instance_id"`
	MechaSectorID  string `db:"mecha_sector_id"`
}

func (r *MechaSectorInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaSectorInstanceGameID] = r.GameID
	args[FieldMechaSectorInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechaSectorInstanceMechaSectorID] = r.MechaSectorID
	return args
}
