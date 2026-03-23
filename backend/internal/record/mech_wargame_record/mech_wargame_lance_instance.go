package mech_wargame_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameLanceInstance string = "mech_wargame_lance_instance"
)

const (
	FieldMechWargameLanceInstanceID                          string = "id"
	FieldMechWargameLanceInstanceGameID                      string = "game_id"
	FieldMechWargameLanceInstanceGameInstanceID              string = "game_instance_id"
	FieldMechWargameLanceInstanceMechWargameLanceID          string = "mech_wargame_lance_id"
	FieldMechWargameLanceInstanceGameSubscriptionInstanceID  string = "game_subscription_instance_id"
	FieldMechWargameLanceInstanceCreatedAt                   string = "created_at"
	FieldMechWargameLanceInstanceUpdatedAt                   string = "updated_at"
	FieldMechWargameLanceInstanceDeletedAt                   string = "deleted_at"
)

type MechWargameLanceInstance struct {
	record.Record
	GameID                     string `db:"game_id"`
	GameInstanceID             string `db:"game_instance_id"`
	MechWargameLanceID         string `db:"mech_wargame_lance_id"`
	GameSubscriptionInstanceID string `db:"game_subscription_instance_id"`
}

func (r *MechWargameLanceInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameLanceInstanceGameID] = r.GameID
	args[FieldMechWargameLanceInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechWargameLanceInstanceMechWargameLanceID] = r.MechWargameLanceID
	args[FieldMechWargameLanceInstanceGameSubscriptionInstanceID] = r.GameSubscriptionInstanceID
	return args
}
