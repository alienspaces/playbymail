package mecha_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaSquad string = "mecha_squad"
)

const (
	FieldMechaSquadID          string = "id"
	FieldMechaSquadGameID      string = "game_id"
	FieldMechaSquadSquadType   string = "squad_type"
	FieldMechaSquadName        string = "name"
	FieldMechaSquadDescription string = "description"
	FieldMechaSquadCreatedAt   string = "created_at"
	FieldMechaSquadUpdatedAt   string = "updated_at"
	FieldMechaSquadDeletedAt   string = "deleted_at"
)

const (
	SquadTypeStarter  = "starter"
	SquadTypeOpponent = "opponent"
)

// MechaSquad is a design-time squad template. Two types exist:
//   - starter: the loadout cloned for every player when they join a run (at most one per game)
//   - opponent: a template randomly assigned to a computer opponent when a run starts
type MechaSquad struct {
	record.Record
	GameID      string `db:"game_id"`
	SquadType   string `db:"squad_type"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *MechaSquad) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaSquadGameID] = r.GameID
	args[FieldMechaSquadSquadType] = r.SquadType
	args[FieldMechaSquadName] = r.Name
	args[FieldMechaSquadDescription] = r.Description
	return args
}
