package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableGameCharacter = "game_character"

const (
	FieldGameCharacterID        = "id"
	FieldGameCharacterGameID    = "game_id"
	FieldGameCharacterAccountID = "account_id"
	FieldGameCharacterName      = "name"
	FieldGameCharacterCreatedAt = "created_at"
	FieldGameCharacterUpdatedAt = "updated_at"
	FieldGameCharacterDeletedAt = "deleted_at"
)

type GameCharacter struct {
	record.Record
	GameID    string `db:"game_id"`
	AccountID string `db:"account_id"`
	Name      string `db:"name"`
}

func (r *GameCharacter) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameCharacterGameID] = r.GameID
	args[FieldGameCharacterAccountID] = r.AccountID
	args[FieldGameCharacterName] = r.Name
	return args
}
