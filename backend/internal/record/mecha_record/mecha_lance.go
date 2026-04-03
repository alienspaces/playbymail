package mecha_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaLance string = "mecha_lance"
)

const (
	FieldMechaLanceID          string = "id"
	FieldMechaLanceGameID      string = "game_id"
	FieldMechaLanceLanceType   string = "lance_type"
	FieldMechaLanceName        string = "name"
	FieldMechaLanceDescription string = "description"
	FieldMechaLanceCreatedAt   string = "created_at"
	FieldMechaLanceUpdatedAt   string = "updated_at"
	FieldMechaLanceDeletedAt   string = "deleted_at"
)

const (
	LanceTypeStarter  = "starter"
	LanceTypeOpponent = "opponent"
)

// MechaLance is a design-time lance template. Two types exist:
//   - starter: the loadout cloned for every player when they join a run (at most one per game)
//   - opponent: a template randomly assigned to a computer opponent when a run starts
type MechaLance struct {
	record.Record
	GameID      string `db:"game_id"`
	LanceType   string `db:"lance_type"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *MechaLance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaLanceGameID] = r.GameID
	args[FieldMechaLanceLanceType] = r.LanceType
	args[FieldMechaLanceName] = r.Name
	args[FieldMechaLanceDescription] = r.Description
	return args
}
