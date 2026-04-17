package mecha_record

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaSquadInstance string = "mecha_squad_instance"
)

const (
	FieldMechaSquadInstanceID                         string = "id"
	FieldMechaSquadInstanceGameID                     string = "game_id"
	FieldMechaSquadInstanceGameInstanceID             string = "game_instance_id"
	FieldMechaSquadInstanceMechaSquadID               string = "mecha_squad_id"
	FieldMechaSquadInstanceGameSubscriptionInstanceID string = "game_subscription_instance_id"
	FieldMechaSquadInstanceMechaComputerOpponentID    string = "mecha_computer_opponent_id"
	FieldMechaSquadInstanceLastTurnEvents             string = "last_turn_events"
	FieldMechaSquadInstanceSupplyPoints               string = "supply_points"
	FieldMechaSquadInstanceCreatedAt                  string = "created_at"
	FieldMechaSquadInstanceUpdatedAt                  string = "updated_at"
	FieldMechaSquadInstanceDeletedAt                  string = "deleted_at"
)

// MechaSquadInstance is the runtime squad record for a game instance.
// For player-owned squads, GameSubscriptionInstanceID is set; MechaComputerOpponentID is NULL.
// For computer-opponent squads, MechaComputerOpponentID is set; GameSubscriptionInstanceID is NULL.
type MechaSquadInstance struct {
	record.Record
	GameID                     string          `db:"game_id"`
	GameInstanceID             string          `db:"game_instance_id"`
	MechaSquadID               string          `db:"mecha_squad_id"`
	GameSubscriptionInstanceID sql.NullString  `db:"game_subscription_instance_id"`
	MechaComputerOpponentID    sql.NullString  `db:"mecha_computer_opponent_id"`
	LastTurnEvents             json.RawMessage `db:"last_turn_events"`
	SupplyPoints               int             `db:"supply_points"`
}

func (r *MechaSquadInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaSquadInstanceGameID] = r.GameID
	args[FieldMechaSquadInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechaSquadInstanceMechaSquadID] = r.MechaSquadID
	args[FieldMechaSquadInstanceGameSubscriptionInstanceID] = r.GameSubscriptionInstanceID
	args[FieldMechaSquadInstanceMechaComputerOpponentID] = r.MechaComputerOpponentID
	if len(r.LastTurnEvents) == 0 {
		args[FieldMechaSquadInstanceLastTurnEvents] = json.RawMessage("[]")
	} else {
		args[FieldMechaSquadInstanceLastTurnEvents] = r.LastTurnEvents
	}
	args[FieldMechaSquadInstanceSupplyPoints] = r.SupplyPoints
	return args
}
