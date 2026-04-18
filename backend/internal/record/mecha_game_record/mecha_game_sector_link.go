package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameSectorLink string = "mecha_game_sector_link"
)

const (
	FieldMechaGameSectorLinkID                string = "id"
	FieldMechaGameSectorLinkGameID            string = "game_id"
	FieldMechaGameSectorLinkFromMechaGameSectorID string = "from_mecha_game_sector_id"
	FieldMechaGameSectorLinkToMechaGameSectorID   string = "to_mecha_game_sector_id"
	FieldMechaGameSectorLinkCreatedAt         string = "created_at"
	FieldMechaGameSectorLinkUpdatedAt         string = "updated_at"
	FieldMechaGameSectorLinkDeletedAt         string = "deleted_at"
)

type MechaGameSectorLink struct {
	record.Record
	GameID            string `db:"game_id"`
	FromMechaGameSectorID string `db:"from_mecha_game_sector_id"`
	ToMechaGameSectorID   string `db:"to_mecha_game_sector_id"`
}

func (r *MechaGameSectorLink) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameSectorLinkGameID] = r.GameID
	args[FieldMechaGameSectorLinkFromMechaGameSectorID] = r.FromMechaGameSectorID
	args[FieldMechaGameSectorLinkToMechaGameSectorID] = r.ToMechaGameSectorID
	return args
}
