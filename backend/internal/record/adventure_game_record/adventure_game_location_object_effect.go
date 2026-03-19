package adventure_game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameLocationObjectEffect = "adventure_game_location_object_effect"

const (
	FieldAdventureGameLocationObjectEffectID                                  = "id"
	FieldAdventureGameLocationObjectEffectGameID                              = "game_id"
	FieldAdventureGameLocationObjectEffectAdventureGameLocationObjectID       = "adventure_game_location_object_id"
	FieldAdventureGameLocationObjectEffectActionType                          = "action_type"
	FieldAdventureGameLocationObjectEffectRequiredAdventureGameLocationObjectStateID = "required_adventure_game_location_object_state_id"
	FieldAdventureGameLocationObjectEffectRequiredAdventureGameItemID               = "required_adventure_game_item_id"
	FieldAdventureGameLocationObjectEffectResultDescription                         = "result_description"
	FieldAdventureGameLocationObjectEffectEffectType                                = "effect_type"
	FieldAdventureGameLocationObjectEffectResultAdventureGameLocationObjectStateID  = "result_adventure_game_location_object_state_id"
	FieldAdventureGameLocationObjectEffectResultAdventureGameItemID           = "result_adventure_game_item_id"
	FieldAdventureGameLocationObjectEffectResultAdventureGameLocationLinkID   = "result_adventure_game_location_link_id"
	FieldAdventureGameLocationObjectEffectResultAdventureGameCreatureID       = "result_adventure_game_creature_id"
	FieldAdventureGameLocationObjectEffectResultAdventureGameLocationObjectID = "result_adventure_game_location_object_id"
	FieldAdventureGameLocationObjectEffectResultAdventureGameLocationID       = "result_adventure_game_location_id"
	FieldAdventureGameLocationObjectEffectResultValueMin                      = "result_value_min"
	FieldAdventureGameLocationObjectEffectResultValueMax                      = "result_value_max"
	FieldAdventureGameLocationObjectEffectIsRepeatable                        = "is_repeatable"
)

// Action type constants — values for the action_type CHECK constraint.
const (
	AdventureGameLocationObjectEffectActionTypeInspect = "inspect"
	AdventureGameLocationObjectEffectActionTypeTouch   = "touch"
	AdventureGameLocationObjectEffectActionTypeOpen    = "open"
	AdventureGameLocationObjectEffectActionTypeClose   = "close"
	AdventureGameLocationObjectEffectActionTypeLock    = "lock"
	AdventureGameLocationObjectEffectActionTypeUnlock  = "unlock"
	AdventureGameLocationObjectEffectActionTypeSearch  = "search"
	AdventureGameLocationObjectEffectActionTypeBreak   = "break"
	AdventureGameLocationObjectEffectActionTypePush    = "push"
	AdventureGameLocationObjectEffectActionTypePull    = "pull"
	AdventureGameLocationObjectEffectActionTypeMove    = "move"
	AdventureGameLocationObjectEffectActionTypeBurn    = "burn"
	AdventureGameLocationObjectEffectActionTypeRead    = "read"
	AdventureGameLocationObjectEffectActionTypeTake    = "take"
	AdventureGameLocationObjectEffectActionTypeListen  = "listen"
	AdventureGameLocationObjectEffectActionTypeInsert  = "insert"
	AdventureGameLocationObjectEffectActionTypePour    = "pour"
	AdventureGameLocationObjectEffectActionTypeDisarm  = "disarm"
	AdventureGameLocationObjectEffectActionTypeClimb   = "climb"
	AdventureGameLocationObjectEffectActionTypeUse     = "use"
)

// AdventureGameLocationObjectEffectActionTypes is the set of all valid action_type values.
var AdventureGameLocationObjectEffectActionTypes = set.New(
	AdventureGameLocationObjectEffectActionTypeInspect,
	AdventureGameLocationObjectEffectActionTypeTouch,
	AdventureGameLocationObjectEffectActionTypeOpen,
	AdventureGameLocationObjectEffectActionTypeClose,
	AdventureGameLocationObjectEffectActionTypeLock,
	AdventureGameLocationObjectEffectActionTypeUnlock,
	AdventureGameLocationObjectEffectActionTypeSearch,
	AdventureGameLocationObjectEffectActionTypeBreak,
	AdventureGameLocationObjectEffectActionTypePush,
	AdventureGameLocationObjectEffectActionTypePull,
	AdventureGameLocationObjectEffectActionTypeMove,
	AdventureGameLocationObjectEffectActionTypeBurn,
	AdventureGameLocationObjectEffectActionTypeRead,
	AdventureGameLocationObjectEffectActionTypeTake,
	AdventureGameLocationObjectEffectActionTypeListen,
	AdventureGameLocationObjectEffectActionTypeInsert,
	AdventureGameLocationObjectEffectActionTypePour,
	AdventureGameLocationObjectEffectActionTypeDisarm,
	AdventureGameLocationObjectEffectActionTypeClimb,
	AdventureGameLocationObjectEffectActionTypeUse,
)

// Effect type constants — values for the effect_type CHECK constraint.
const (
	AdventureGameLocationObjectEffectEffectTypeInfo              = "info"
	AdventureGameLocationObjectEffectEffectTypeChangeState       = "change_state"
	AdventureGameLocationObjectEffectEffectTypeChangeObjectState = "change_object_state"
	AdventureGameLocationObjectEffectEffectTypeGiveItem          = "give_item"
	AdventureGameLocationObjectEffectEffectTypeRemoveItem        = "remove_item"
	AdventureGameLocationObjectEffectEffectTypeOpenLink          = "open_link"
	AdventureGameLocationObjectEffectEffectTypeCloseLink         = "close_link"
	AdventureGameLocationObjectEffectEffectTypeRevealObject      = "reveal_object"
	AdventureGameLocationObjectEffectEffectTypeHideObject        = "hide_object"
	AdventureGameLocationObjectEffectEffectTypeDamage            = "damage"
	AdventureGameLocationObjectEffectEffectTypeHeal              = "heal"
	AdventureGameLocationObjectEffectEffectTypeSummonCreature    = "summon_creature"
	AdventureGameLocationObjectEffectEffectTypeTeleport          = "teleport"
	AdventureGameLocationObjectEffectEffectTypeNothing           = "nothing"
	AdventureGameLocationObjectEffectEffectTypeRemoveObject      = "remove_object"
)

// AdventureGameLocationObjectEffectEffectTypes is the set of all valid effect_type values.
var AdventureGameLocationObjectEffectEffectTypes = set.New(
	AdventureGameLocationObjectEffectEffectTypeInfo,
	AdventureGameLocationObjectEffectEffectTypeChangeState,
	AdventureGameLocationObjectEffectEffectTypeChangeObjectState,
	AdventureGameLocationObjectEffectEffectTypeGiveItem,
	AdventureGameLocationObjectEffectEffectTypeRemoveItem,
	AdventureGameLocationObjectEffectEffectTypeOpenLink,
	AdventureGameLocationObjectEffectEffectTypeCloseLink,
	AdventureGameLocationObjectEffectEffectTypeRevealObject,
	AdventureGameLocationObjectEffectEffectTypeHideObject,
	AdventureGameLocationObjectEffectEffectTypeDamage,
	AdventureGameLocationObjectEffectEffectTypeHeal,
	AdventureGameLocationObjectEffectEffectTypeSummonCreature,
	AdventureGameLocationObjectEffectEffectTypeTeleport,
	AdventureGameLocationObjectEffectEffectTypeNothing,
	AdventureGameLocationObjectEffectEffectTypeRemoveObject,
)

// AdventureGameLocationObjectEffect defines what happens when a player performs an action on an object.
// Multiple rows for the same (object, action_type, required_state_id) are allowed and all fire atomically.
type AdventureGameLocationObjectEffect struct {
	record.Record
	GameID                                          string         `db:"game_id"`
	AdventureGameLocationObjectID                   string         `db:"adventure_game_location_object_id"`
	ActionType                                      string         `db:"action_type"`
	RequiredAdventureGameLocationObjectStateID      sql.NullString `db:"required_adventure_game_location_object_state_id"`
	RequiredAdventureGameItemID                     sql.NullString `db:"required_adventure_game_item_id"`
	ResultDescription                               string         `db:"result_description"`
	EffectType                                      string         `db:"effect_type"`
	ResultAdventureGameLocationObjectStateID        sql.NullString `db:"result_adventure_game_location_object_state_id"`
	ResultAdventureGameItemID                       sql.NullString `db:"result_adventure_game_item_id"`
	ResultAdventureGameLocationLinkID               sql.NullString `db:"result_adventure_game_location_link_id"`
	ResultAdventureGameCreatureID                   sql.NullString `db:"result_adventure_game_creature_id"`
	ResultAdventureGameLocationObjectID             sql.NullString `db:"result_adventure_game_location_object_id"`
	ResultAdventureGameLocationID                   sql.NullString `db:"result_adventure_game_location_id"`
	ResultValueMin                                  sql.NullInt32  `db:"result_value_min"`
	ResultValueMax                                  sql.NullInt32  `db:"result_value_max"`
	IsRepeatable                                    bool           `db:"is_repeatable"`
}

func (r *AdventureGameLocationObjectEffect) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationObjectEffectGameID] = r.GameID
	args[FieldAdventureGameLocationObjectEffectAdventureGameLocationObjectID] = r.AdventureGameLocationObjectID
	args[FieldAdventureGameLocationObjectEffectActionType] = r.ActionType
	args[FieldAdventureGameLocationObjectEffectRequiredAdventureGameLocationObjectStateID] = r.RequiredAdventureGameLocationObjectStateID
	args[FieldAdventureGameLocationObjectEffectRequiredAdventureGameItemID] = r.RequiredAdventureGameItemID
	args[FieldAdventureGameLocationObjectEffectResultDescription] = r.ResultDescription
	args[FieldAdventureGameLocationObjectEffectEffectType] = r.EffectType
	args[FieldAdventureGameLocationObjectEffectResultAdventureGameLocationObjectStateID] = r.ResultAdventureGameLocationObjectStateID
	args[FieldAdventureGameLocationObjectEffectResultAdventureGameItemID] = r.ResultAdventureGameItemID
	args[FieldAdventureGameLocationObjectEffectResultAdventureGameLocationLinkID] = r.ResultAdventureGameLocationLinkID
	args[FieldAdventureGameLocationObjectEffectResultAdventureGameCreatureID] = r.ResultAdventureGameCreatureID
	args[FieldAdventureGameLocationObjectEffectResultAdventureGameLocationObjectID] = r.ResultAdventureGameLocationObjectID
	args[FieldAdventureGameLocationObjectEffectResultAdventureGameLocationID] = r.ResultAdventureGameLocationID
	args[FieldAdventureGameLocationObjectEffectResultValueMin] = r.ResultValueMin
	args[FieldAdventureGameLocationObjectEffectResultValueMax] = r.ResultValueMax
	args[FieldAdventureGameLocationObjectEffectIsRepeatable] = r.IsRepeatable
	return args
}
