package mech_wargame_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameSectorInstance string = "mech_wargame_sector_instance"
)

const (
	FieldMechWargameSectorInstanceID                  string = "id"
	FieldMechWargameSectorInstanceGameID              string = "game_id"
	FieldMechWargameSectorInstanceGameInstanceID      string = "game_instance_id"
	FieldMechWargameSectorInstanceMechWargameSectorID string = "mech_wargame_sector_id"
	FieldMechWargameSectorInstanceCreatedAt           string = "created_at"
	FieldMechWargameSectorInstanceUpdatedAt           string = "updated_at"
	FieldMechWargameSectorInstanceDeletedAt           string = "deleted_at"
)

type MechWargameSectorInstance struct {
	record.Record
	GameID               string `db:"game_id"`
	GameInstanceID       string `db:"game_instance_id"`
	MechWargameSectorID  string `db:"mech_wargame_sector_id"`
}

func (r *MechWargameSectorInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameSectorInstanceGameID] = r.GameID
	args[FieldMechWargameSectorInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechWargameSectorInstanceMechWargameSectorID] = r.MechWargameSectorID
	return args
}
