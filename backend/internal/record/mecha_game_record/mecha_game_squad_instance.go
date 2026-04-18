package mecha_game_record

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameSquadInstance string = "mecha_game_squad_instance"
)

const (
	FieldMechaGameSquadInstanceID                         string = "id"
	FieldMechaGameSquadInstanceGameID                     string = "game_id"
	FieldMechaGameSquadInstanceGameInstanceID             string = "game_instance_id"
	FieldMechaGameSquadInstanceMechaGameSquadID               string = "mecha_game_squad_id"
	FieldMechaGameSquadInstanceGameSubscriptionInstanceID string = "game_subscription_instance_id"
	FieldMechaGameSquadInstanceMechaGameComputerOpponentID    string = "mecha_game_computer_opponent_id"
	FieldMechaGameSquadInstanceLastTurnEvents             string = "last_turn_events"
	FieldMechaGameSquadInstanceSupplyPoints               string = "supply_points"
	FieldMechaGameSquadInstanceCreatedAt                  string = "created_at"
	FieldMechaGameSquadInstanceUpdatedAt                  string = "updated_at"
	FieldMechaGameSquadInstanceDeletedAt                  string = "deleted_at"
)

// MechaGameSquadInstance is the runtime squad record for a game instance.
// For player-owned squads, GameSubscriptionInstanceID is set; MechaGameComputerOpponentID is NULL.
// For computer-opponent squads, MechaGameComputerOpponentID is set; GameSubscriptionInstanceID is NULL.
type MechaGameSquadInstance struct {
	record.Record
	GameID                     string          `db:"game_id"`
	GameInstanceID             string          `db:"game_instance_id"`
	MechaGameSquadID               string          `db:"mecha_game_squad_id"`
	GameSubscriptionInstanceID sql.NullString  `db:"game_subscription_instance_id"`
	MechaGameComputerOpponentID    sql.NullString  `db:"mecha_game_computer_opponent_id"`
	LastTurnEvents             json.RawMessage `db:"last_turn_events"`
	SupplyPoints               int             `db:"supply_points"`
}

func (r *MechaGameSquadInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameSquadInstanceGameID] = r.GameID
	args[FieldMechaGameSquadInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechaGameSquadInstanceMechaGameSquadID] = r.MechaGameSquadID
	args[FieldMechaGameSquadInstanceGameSubscriptionInstanceID] = r.GameSubscriptionInstanceID
	args[FieldMechaGameSquadInstanceMechaGameComputerOpponentID] = r.MechaGameComputerOpponentID
	if len(r.LastTurnEvents) == 0 {
		args[FieldMechaGameSquadInstanceLastTurnEvents] = json.RawMessage("[]")
	} else {
		args[FieldMechaGameSquadInstanceLastTurnEvents] = r.LastTurnEvents
	}
	args[FieldMechaGameSquadInstanceSupplyPoints] = r.SupplyPoints
	return args
}
