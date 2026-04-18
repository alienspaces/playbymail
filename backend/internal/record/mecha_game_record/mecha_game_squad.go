package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameSquad string = "mecha_game_squad"
)

const (
	FieldMechaGameSquadID          string = "id"
	FieldMechaGameSquadGameID      string = "game_id"
	FieldMechaGameSquadSquadType   string = "squad_type"
	FieldMechaGameSquadName        string = "name"
	FieldMechaGameSquadDescription string = "description"
	FieldMechaGameSquadCreatedAt   string = "created_at"
	FieldMechaGameSquadUpdatedAt   string = "updated_at"
	FieldMechaGameSquadDeletedAt   string = "deleted_at"
)

const (
	SquadTypeStarter  = "starter"
	SquadTypeOpponent = "opponent"
)

// MechaGameSquad is a design-time squad template. Two types exist:
//   - starter: the loadout cloned for every player when they join a run (at most one per game)
//   - opponent: a template randomly assigned to a computer opponent when a run starts
type MechaGameSquad struct {
	record.Record
	GameID      string `db:"game_id"`
	SquadType   string `db:"squad_type"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *MechaGameSquad) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameSquadGameID] = r.GameID
	args[FieldMechaGameSquadSquadType] = r.SquadType
	args[FieldMechaGameSquadName] = r.Name
	args[FieldMechaGameSquadDescription] = r.Description
	return args
}
