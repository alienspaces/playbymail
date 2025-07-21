package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameInstance = "adventure_game_instance"

const (
	FieldAdventureGameInstanceID        = "id"
	FieldAdventureGameInstanceGameID    = "game_id"
	FieldAdventureGameInstanceCreatedAt = "created_at"
	FieldAdventureGameInstanceUpdatedAt = "updated_at"
	FieldAdventureGameInstanceDeletedAt = "deleted_at"
)

type AdventureGameInstance struct {
	record.Record
	GameID string `db:"game_id"`
}

func (r *AdventureGameInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameInstanceGameID] = r.GameID
	return args
}
