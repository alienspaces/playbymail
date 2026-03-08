package game_record

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	TableManagerGameInstanceView = "manager_game_instance_view"
)

const (
	FieldManagerGameInstanceVID                 = "id"
	FieldManagerGameInstanceVAccountID          = "account_id"
	FieldManagerGameInstanceVGameID             = "game_id"
	FieldManagerGameInstanceVGameName           = "game_name"
	FieldManagerGameInstanceVGameType           = "game_type"
	FieldManagerGameInstanceVGameDescription    = "game_description"
	FieldManagerGameInstanceVGameSubscriptionID = "game_subscription_id"
	FieldManagerGameInstanceVGameInstanceID     = "game_instance_id"
	FieldManagerGameInstanceVInstanceStatus     = "instance_status"
	FieldManagerGameInstanceVCurrentTurn        = "current_turn"
	FieldManagerGameInstanceVRequiredPlayerCnt  = "required_player_count"
	FieldManagerGameInstanceVDeliveryEmail      = "delivery_email"
	FieldManagerGameInstanceVDeliveryPhysPost   = "delivery_physical_post"
	FieldManagerGameInstanceVDeliveryPhysLocal  = "delivery_physical_local"
	FieldManagerGameInstanceVIsClosedTesting    = "is_closed_testing"
	FieldManagerGameInstanceVStartedAt          = "started_at"
	FieldManagerGameInstanceVNextTurnDueAt      = "next_turn_due_at"
	FieldManagerGameInstanceVInstanceCreatedAt  = "instance_created_at"
	FieldManagerGameInstanceVCreatedAt          = "created_at"
	FieldManagerGameInstanceVUpdatedAt          = "updated_at"
	FieldManagerGameInstanceVDeletedAt          = "deleted_at"
)

type ManagerGameInstanceView struct {
	ID                    string         `db:"id"`
	AccountID             string         `db:"account_id"`
	GameID                string         `db:"game_id"`
	GameName              string         `db:"game_name"`
	GameType              string         `db:"game_type"`
	GameDescription       string         `db:"game_description"`
	GameSubscriptionID    string         `db:"game_subscription_id"`
	GameInstanceID        sql.NullString `db:"game_instance_id"`
	InstanceStatus        sql.NullString `db:"instance_status"`
	CurrentTurn           sql.NullInt32  `db:"current_turn"`
	RequiredPlayerCount   sql.NullInt32  `db:"required_player_count"`
	DeliveryEmail         sql.NullBool   `db:"delivery_email"`
	DeliveryPhysicalPost  sql.NullBool   `db:"delivery_physical_post"`
	DeliveryPhysicalLocal sql.NullBool   `db:"delivery_physical_local"`
	IsClosedTesting       sql.NullBool   `db:"is_closed_testing"`
	StartedAt             sql.NullTime   `db:"started_at"`
	NextTurnDueAt         sql.NullTime   `db:"next_turn_due_at"`
	InstanceCreatedAt     sql.NullTime   `db:"instance_created_at"`
	CreatedAt             time.Time      `db:"created_at"`
	UpdatedAt             sql.NullTime   `db:"updated_at"`
	DeletedAt             sql.NullTime   `db:"deleted_at"`
}

func (r *ManagerGameInstanceView) ToNamedArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		FieldManagerGameInstanceVID:                 r.ID,
		FieldManagerGameInstanceVAccountID:          r.AccountID,
		FieldManagerGameInstanceVGameID:             r.GameID,
		FieldManagerGameInstanceVGameName:           r.GameName,
		FieldManagerGameInstanceVGameType:           r.GameType,
		FieldManagerGameInstanceVGameDescription:    r.GameDescription,
		FieldManagerGameInstanceVGameSubscriptionID: r.GameSubscriptionID,
		FieldManagerGameInstanceVGameInstanceID:     r.GameInstanceID,
		FieldManagerGameInstanceVInstanceStatus:     r.InstanceStatus,
		FieldManagerGameInstanceVCurrentTurn:        r.CurrentTurn,
		FieldManagerGameInstanceVRequiredPlayerCnt:  r.RequiredPlayerCount,
		FieldManagerGameInstanceVDeliveryEmail:      r.DeliveryEmail,
		FieldManagerGameInstanceVDeliveryPhysPost:   r.DeliveryPhysicalPost,
		FieldManagerGameInstanceVDeliveryPhysLocal:  r.DeliveryPhysicalLocal,
		FieldManagerGameInstanceVIsClosedTesting:    r.IsClosedTesting,
		FieldManagerGameInstanceVStartedAt:          r.StartedAt,
		FieldManagerGameInstanceVNextTurnDueAt:      r.NextTurnDueAt,
		FieldManagerGameInstanceVInstanceCreatedAt:  r.InstanceCreatedAt,
		FieldManagerGameInstanceVCreatedAt:          r.CreatedAt,
		FieldManagerGameInstanceVUpdatedAt:          r.UpdatedAt,
		FieldManagerGameInstanceVDeletedAt:          r.DeletedAt,
	}
}
