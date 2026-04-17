package mecha_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaComputerOpponent string = "mecha_computer_opponent"
)

const (
	FieldMechaComputerOpponentID          string = "id"
	FieldMechaComputerOpponentGameID      string = "game_id"
	FieldMechaComputerOpponentName        string = "name"
	FieldMechaComputerOpponentDescription string = "description"
	FieldMechaComputerOpponentAggression  string = "aggression"
	FieldMechaComputerOpponentIQ          string = "iq"
	FieldMechaComputerOpponentCreatedAt   string = "created_at"
	FieldMechaComputerOpponentUpdatedAt   string = "updated_at"
	FieldMechaComputerOpponentDeletedAt   string = "deleted_at"
)

// MechaComputerOpponent is a computer-controlled opposing command in a mecha game.
// It owns one or more squads and holds behaviour configuration used by the
// decision engine during turn processing.
//
// Aggression (1-10): 1 = purely defensive, 10 = all-out assault.
// IQ (1-10): 1 = predictable/random moves, 10 = expert use of terrain and flanking.
type MechaComputerOpponent struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Aggression  int    `db:"aggression"`
	IQ          int    `db:"iq"`
}

func (r *MechaComputerOpponent) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaComputerOpponentGameID] = r.GameID
	args[FieldMechaComputerOpponentName] = r.Name
	args[FieldMechaComputerOpponentDescription] = r.Description
	args[FieldMechaComputerOpponentAggression] = r.Aggression
	args[FieldMechaComputerOpponentIQ] = r.IQ
	return args
}
