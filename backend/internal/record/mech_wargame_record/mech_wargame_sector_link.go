package mech_wargame_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameSectorLink string = "mech_wargame_sector_link"
)

const (
	FieldMechWargameSectorLinkID                    string = "id"
	FieldMechWargameSectorLinkGameID                string = "game_id"
	FieldMechWargameSectorLinkFromMechWargameSectorID string = "from_mech_wargame_sector_id"
	FieldMechWargameSectorLinkToMechWargameSectorID   string = "to_mech_wargame_sector_id"
	FieldMechWargameSectorLinkCoverModifier          string = "cover_modifier"
	FieldMechWargameSectorLinkCreatedAt              string = "created_at"
	FieldMechWargameSectorLinkUpdatedAt              string = "updated_at"
	FieldMechWargameSectorLinkDeletedAt              string = "deleted_at"
)

type MechWargameSectorLink struct {
	record.Record
	GameID                    string `db:"game_id"`
	FromMechWargameSectorID   string `db:"from_mech_wargame_sector_id"`
	ToMechWargameSectorID     string `db:"to_mech_wargame_sector_id"`
	CoverModifier             int    `db:"cover_modifier"`
}

func (r *MechWargameSectorLink) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameSectorLinkGameID] = r.GameID
	args[FieldMechWargameSectorLinkFromMechWargameSectorID] = r.FromMechWargameSectorID
	args[FieldMechWargameSectorLinkToMechWargameSectorID] = r.ToMechWargameSectorID
	args[FieldMechWargameSectorLinkCoverModifier] = r.CoverModifier
	return args
}
