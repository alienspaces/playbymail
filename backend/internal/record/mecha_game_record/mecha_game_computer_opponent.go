package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameComputerOpponent string = "mecha_game_computer_opponent"
)

const (
	FieldMechaGameComputerOpponentID          string = "id"
	FieldMechaGameComputerOpponentGameID      string = "game_id"
	FieldMechaGameComputerOpponentName        string = "name"
	FieldMechaGameComputerOpponentDescription string = "description"
	FieldMechaGameComputerOpponentAggression  string = "aggression"
	FieldMechaGameComputerOpponentIQ          string = "iq"
	FieldMechaGameComputerOpponentCreatedAt   string = "created_at"
	FieldMechaGameComputerOpponentUpdatedAt   string = "updated_at"
	FieldMechaGameComputerOpponentDeletedAt   string = "deleted_at"
)

// MechaGameComputerOpponent is a computer-controlled opposing command in a mecha game.
// It owns one or more squads and holds behaviour configuration used by the
// decision engine during turn processing.
//
// Aggression (1-10): 1 = purely defensive, 10 = all-out assault.
// IQ (1-10): 1 = predictable/random moves, 10 = expert use of terrain and flanking.
type MechaGameComputerOpponent struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Aggression  int    `db:"aggression"`
	IQ          int    `db:"iq"`
}

func (r *MechaGameComputerOpponent) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameComputerOpponentGameID] = r.GameID
	args[FieldMechaGameComputerOpponentName] = r.Name
	args[FieldMechaGameComputerOpponentDescription] = r.Description
	args[FieldMechaGameComputerOpponentAggression] = r.Aggression
	args[FieldMechaGameComputerOpponentIQ] = r.IQ
	return args
}
