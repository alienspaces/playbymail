package mecha_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaSectorLink string = "mecha_sector_link"
)

const (
	FieldMechaSectorLinkID                    string = "id"
	FieldMechaSectorLinkGameID                string = "game_id"
	FieldMechaSectorLinkFromMechaSectorID string = "from_mecha_sector_id"
	FieldMechaSectorLinkToMechaSectorID   string = "to_mecha_sector_id"
	FieldMechaSectorLinkCoverModifier          string = "cover_modifier"
	FieldMechaSectorLinkCreatedAt              string = "created_at"
	FieldMechaSectorLinkUpdatedAt              string = "updated_at"
	FieldMechaSectorLinkDeletedAt              string = "deleted_at"
)

type MechaSectorLink struct {
	record.Record
	GameID                    string `db:"game_id"`
	FromMechaSectorID   string `db:"from_mecha_sector_id"`
	ToMechaSectorID     string `db:"to_mecha_sector_id"`
	CoverModifier             int    `db:"cover_modifier"`
}

func (r *MechaSectorLink) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaSectorLinkGameID] = r.GameID
	args[FieldMechaSectorLinkFromMechaSectorID] = r.FromMechaSectorID
	args[FieldMechaSectorLinkToMechaSectorID] = r.ToMechaSectorID
	args[FieldMechaSectorLinkCoverModifier] = r.CoverModifier
	return args
}
