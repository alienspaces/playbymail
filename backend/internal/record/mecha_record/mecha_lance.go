package mecha_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaLance string = "mecha_lance"
)

const (
	FieldMechaLanceID                      string = "id"
	FieldMechaLanceGameID                  string = "game_id"
	FieldMechaLanceAccountID               string = "account_id"
	FieldMechaLanceAccountUserID           string = "account_user_id"
	FieldMechaLanceMechaComputerOpponentID string = "mecha_computer_opponent_id"
	FieldMechaLanceIsPlayerStarter         string = "is_player_starter"
	FieldMechaLanceName                    string = "name"
	FieldMechaLanceDescription             string = "description"
	FieldMechaLanceCreatedAt               string = "created_at"
	FieldMechaLanceUpdatedAt               string = "updated_at"
	FieldMechaLanceDeletedAt               string = "deleted_at"
)

// MechaLance belongs to one of three ownership modes, enforced by the DB CHECK constraint:
//   - Human player: AccountID + AccountUserID set, IsPlayerStarter false
//   - Computer opponent: MechaComputerOpponentID set, IsPlayerStarter false
//   - Player starter template: IsPlayerStarter true, all owner fields NULL
type MechaLance struct {
	record.Record
	GameID                  string         `db:"game_id"`
	AccountID               sql.NullString `db:"account_id"`
	AccountUserID           sql.NullString `db:"account_user_id"`
	MechaComputerOpponentID sql.NullString `db:"mecha_computer_opponent_id"`
	IsPlayerStarter         bool           `db:"is_player_starter"`
	Name                    string         `db:"name"`
	Description             string         `db:"description"`
}

func (r *MechaLance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaLanceGameID] = r.GameID
	args[FieldMechaLanceAccountID] = r.AccountID
	args[FieldMechaLanceAccountUserID] = r.AccountUserID
	args[FieldMechaLanceMechaComputerOpponentID] = r.MechaComputerOpponentID
	args[FieldMechaLanceIsPlayerStarter] = r.IsPlayerStarter
	args[FieldMechaLanceName] = r.Name
	args[FieldMechaLanceDescription] = r.Description
	return args
}
