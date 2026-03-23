package mech_wargame_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameLance string = "mech_wargame_lance"
)

const (
	FieldMechWargameLanceID            string = "id"
	FieldMechWargameLanceGameID        string = "game_id"
	FieldMechWargameLanceAccountID     string = "account_id"
	FieldMechWargameLanceAccountUserID string = "account_user_id"
	FieldMechWargameLanceName          string = "name"
	FieldMechWargameLanceDescription   string = "description"
	FieldMechWargameLanceCreatedAt     string = "created_at"
	FieldMechWargameLanceUpdatedAt     string = "updated_at"
	FieldMechWargameLanceDeletedAt     string = "deleted_at"
)

type MechWargameLance struct {
	record.Record
	GameID        string `db:"game_id"`
	AccountID     string `db:"account_id"`
	AccountUserID string `db:"account_user_id"`
	Name          string `db:"name"`
	Description   string `db:"description"`
}

func (r *MechWargameLance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameLanceGameID] = r.GameID
	args[FieldMechWargameLanceAccountID] = r.AccountID
	args[FieldMechWargameLanceAccountUserID] = r.AccountUserID
	args[FieldMechWargameLanceName] = r.Name
	args[FieldMechWargameLanceDescription] = r.Description
	return args
}
