package adventure_game_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameCharacter = "adventure_game_character"

const (
	FieldAdventureGameCharacterID            = "id"
	FieldAdventureGameCharacterGameID        = "game_id"
	FieldAdventureGameCharacterAccountID     = "account_id"
	FieldAdventureGameCharacterAccountUserID = "account_user_id"
	FieldAdventureGameCharacterName          = "name"
	FieldAdventureGameCharacterCreatedAt     = "created_at"
	FieldAdventureGameCharacterUpdatedAt     = "updated_at"
	FieldAdventureGameCharacterDeletedAt     = "deleted_at"
)

type AdventureGameCharacter struct {
	record.Record
	GameID        string `db:"game_id"`
	AccountID     string `db:"account_id"`
	AccountUserID string `db:"account_user_id"`
	Name          string `db:"name"`
}

func (r *AdventureGameCharacter) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameCharacterGameID] = r.GameID
	args[FieldAdventureGameCharacterAccountID] = r.AccountID
	args[FieldAdventureGameCharacterAccountUserID] = r.AccountUserID
	args[FieldAdventureGameCharacterName] = r.Name
	return args
}
