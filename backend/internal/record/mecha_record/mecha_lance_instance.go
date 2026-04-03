package mecha_record

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaLanceInstance string = "mecha_lance_instance"
)

const (
	FieldMechaLanceInstanceID                         string = "id"
	FieldMechaLanceInstanceGameID                     string = "game_id"
	FieldMechaLanceInstanceGameInstanceID             string = "game_instance_id"
	FieldMechaLanceInstanceMechaLanceID               string = "mecha_lance_id"
	FieldMechaLanceInstanceGameSubscriptionInstanceID string = "game_subscription_instance_id"
	FieldMechaLanceInstanceMechaComputerOpponentID    string = "mecha_computer_opponent_id"
	FieldMechaLanceInstanceLastTurnEvents             string = "last_turn_events"
	FieldMechaLanceInstanceSupplyPoints               string = "supply_points"
	FieldMechaLanceInstanceCreatedAt                  string = "created_at"
	FieldMechaLanceInstanceUpdatedAt                  string = "updated_at"
	FieldMechaLanceInstanceDeletedAt                  string = "deleted_at"
)

// MechaLanceInstance is the runtime lance record for a game instance.
// For player-owned lances, GameSubscriptionInstanceID is set; MechaComputerOpponentID is NULL.
// For computer-opponent lances, MechaComputerOpponentID is set; GameSubscriptionInstanceID is NULL.
type MechaLanceInstance struct {
	record.Record
	GameID                     string          `db:"game_id"`
	GameInstanceID             string          `db:"game_instance_id"`
	MechaLanceID               string          `db:"mecha_lance_id"`
	GameSubscriptionInstanceID sql.NullString  `db:"game_subscription_instance_id"`
	MechaComputerOpponentID    sql.NullString  `db:"mecha_computer_opponent_id"`
	LastTurnEvents             json.RawMessage `db:"last_turn_events"`
	SupplyPoints               int             `db:"supply_points"`
}

func (r *MechaLanceInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaLanceInstanceGameID] = r.GameID
	args[FieldMechaLanceInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechaLanceInstanceMechaLanceID] = r.MechaLanceID
	args[FieldMechaLanceInstanceGameSubscriptionInstanceID] = r.GameSubscriptionInstanceID
	args[FieldMechaLanceInstanceMechaComputerOpponentID] = r.MechaComputerOpponentID
	if len(r.LastTurnEvents) == 0 {
		args[FieldMechaLanceInstanceLastTurnEvents] = json.RawMessage("[]")
	} else {
		args[FieldMechaLanceInstanceLastTurnEvents] = r.LastTurnEvents
	}
	args[FieldMechaLanceInstanceSupplyPoints] = r.SupplyPoints
	return args
}
