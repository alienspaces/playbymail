package game_record

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	TableCatalogGameInstanceView = "catalog_game_instance_view"
)

const (
	FieldCGIVID                  = "id"
	FieldCGIVGameInstanceID      = "game_instance_id"
	FieldCGIVGameID              = "game_id"
	FieldCGIVGameName            = "game_name"
	FieldCGIVGameType            = "game_type"
	FieldCGIVGameDescription     = "game_description"
	FieldCGIVTurnDurationHours   = "turn_duration_hours"
	FieldCGIVGameSubscriptionID  = "game_subscription_id"
	FieldCGIVAccountName         = "account_name"
	FieldCGIVRequiredPlayerCount = "required_player_count"
	FieldCGIVPlayerCount         = "player_count"
	FieldCGIVRemainingCapacity   = "remaining_capacity"
	FieldCGIVDeliveryEmail       = "delivery_email"
	FieldCGIVDeliveryPhysPost    = "delivery_physical_post"
	FieldCGIVDeliveryPhysLocal   = "delivery_physical_local"
	FieldCGIVIsClosedTesting     = "is_closed_testing"
	FieldCGIVCreatedAt           = "created_at"
	FieldCGIVUpdatedAt           = "updated_at"
	FieldCGIVDeletedAt           = "deleted_at"
)

type CatalogGameInstanceView struct {
	ID                    string       `db:"id"`
	GameInstanceID        string       `db:"game_instance_id"`
	GameID                string       `db:"game_id"`
	GameName              string       `db:"game_name"`
	GameType              string       `db:"game_type"`
	GameDescription       string       `db:"game_description"`
	TurnDurationHours     int          `db:"turn_duration_hours"`
	GameSubscriptionID    string       `db:"game_subscription_id"`
	AccountName           string       `db:"account_name"`
	RequiredPlayerCount   int          `db:"required_player_count"`
	PlayerCount           int          `db:"player_count"`
	RemainingCapacity     int          `db:"remaining_capacity"`
	DeliveryEmail         bool         `db:"delivery_email"`
	DeliveryPhysicalPost  bool         `db:"delivery_physical_post"`
	DeliveryPhysicalLocal bool         `db:"delivery_physical_local"`
	IsClosedTesting       bool         `db:"is_closed_testing"`
	CreatedAt             time.Time    `db:"created_at"`
	UpdatedAt             sql.NullTime `db:"updated_at"`
	DeletedAt             sql.NullTime `db:"deleted_at"`
}

func (r *CatalogGameInstanceView) ToNamedArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		FieldCGIVID:                  r.ID,
		FieldCGIVGameInstanceID:      r.GameInstanceID,
		FieldCGIVGameID:              r.GameID,
		FieldCGIVGameName:            r.GameName,
		FieldCGIVGameType:            r.GameType,
		FieldCGIVGameDescription:     r.GameDescription,
		FieldCGIVTurnDurationHours:   r.TurnDurationHours,
		FieldCGIVGameSubscriptionID:  r.GameSubscriptionID,
		FieldCGIVAccountName:         r.AccountName,
		FieldCGIVRequiredPlayerCount: r.RequiredPlayerCount,
		FieldCGIVPlayerCount:         r.PlayerCount,
		FieldCGIVRemainingCapacity:   r.RemainingCapacity,
		FieldCGIVDeliveryEmail:       r.DeliveryEmail,
		FieldCGIVDeliveryPhysPost:    r.DeliveryPhysicalPost,
		FieldCGIVDeliveryPhysLocal:   r.DeliveryPhysicalLocal,
		FieldCGIVIsClosedTesting:     r.IsClosedTesting,
		FieldCGIVCreatedAt:           r.CreatedAt,
		FieldCGIVUpdatedAt:           r.UpdatedAt,
		FieldCGIVDeletedAt:           r.DeletedAt,
	}
}
