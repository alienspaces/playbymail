package game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// Table and field constants
const (
	TableGameAdministration = "game_administration"
)

const (
	FieldGameAdministrationID                 = "id"
	FieldGameAdministrationGameID             = "game_id"
	FieldGameAdministrationAccountID          = "account_id"
	FieldGameAdministrationGrantedByAccountID = "granted_by_account_id"
	FieldGameAdministrationCreatedAt          = "created_at"
)

// GameAdministration represents admin rights for all instances of a game
type GameAdministration struct {
	record.Record
	GameID             string `db:"game_id"`
	AccountID          string `db:"account_id"`
	GrantedByAccountID string `db:"granted_by_account_id"`
}

func (r *GameAdministration) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameAdministrationGameID] = r.GameID
	args[FieldGameAdministrationAccountID] = r.AccountID
	args[FieldGameAdministrationGrantedByAccountID] = r.GrantedByAccountID
	return args
}
