package adventure_game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameItemEffect = "adventure_game_item_effect"

const (
	FieldAdventureGameItemEffectID                                = "id"
	FieldAdventureGameItemEffectGameID                            = "game_id"
	FieldAdventureGameItemEffectAdventureGameItemID               = "adventure_game_item_id"
	FieldAdventureGameItemEffectActionType                        = "action_type"
	FieldAdventureGameItemEffectRequiredAdventureGameItemID       = "required_adventure_game_item_id"
	FieldAdventureGameItemEffectRequiredAdventureGameLocationID   = "required_adventure_game_location_id"
	FieldAdventureGameItemEffectResultDescription                 = "result_description"
	FieldAdventureGameItemEffectEffectType                        = "effect_type"
	FieldAdventureGameItemEffectResultAdventureGameItemID         = "result_adventure_game_item_id"
	FieldAdventureGameItemEffectResultAdventureGameLocationLinkID = "result_adventure_game_location_link_id"
	FieldAdventureGameItemEffectResultAdventureGameCreatureID     = "result_adventure_game_creature_id"
	FieldAdventureGameItemEffectResultAdventureGameLocationID     = "result_adventure_game_location_id"
	FieldAdventureGameItemEffectResultValueMin                    = "result_value_min"
	FieldAdventureGameItemEffectResultValueMax                    = "result_value_max"
	FieldAdventureGameItemEffectIsRepeatable                      = "is_repeatable"
)

// Action type constants — values for the action_type CHECK constraint.
const (
	AdventureGameItemEffectActionTypeUse     = "use"
	AdventureGameItemEffectActionTypeEquip   = "equip"
	AdventureGameItemEffectActionTypeUnequip = "unequip"
	AdventureGameItemEffectActionTypeInspect = "inspect"
	AdventureGameItemEffectActionTypeDrop    = "drop"
	AdventureGameItemEffectActionTypePickup  = "pickup"
)

// AdventureGameItemEffectActionTypes is the set of all valid action_type values.
var AdventureGameItemEffectActionTypes = set.New(
	AdventureGameItemEffectActionTypeUse,
	AdventureGameItemEffectActionTypeEquip,
	AdventureGameItemEffectActionTypeUnequip,
	AdventureGameItemEffectActionTypeInspect,
	AdventureGameItemEffectActionTypeDrop,
	AdventureGameItemEffectActionTypePickup,
)

// Effect type constants — values for the effect_type CHECK constraint.
const (
	AdventureGameItemEffectEffectTypeInfo           = "info"
	AdventureGameItemEffectEffectTypeDamageTarget   = "damage_target"
	AdventureGameItemEffectEffectTypeDamageWielder  = "damage_wielder"
	AdventureGameItemEffectEffectTypeHealTarget     = "heal_target"
	AdventureGameItemEffectEffectTypeHealWielder    = "heal_wielder"
	AdventureGameItemEffectEffectTypeTeleport       = "teleport"
	AdventureGameItemEffectEffectTypeOpenLink       = "open_link"
	AdventureGameItemEffectEffectTypeCloseLink      = "close_link"
	AdventureGameItemEffectEffectTypeGiveItem       = "give_item"
	AdventureGameItemEffectEffectTypeRemoveItem     = "remove_item"
	AdventureGameItemEffectEffectTypeSummonCreature = "summon_creature"
	AdventureGameItemEffectEffectTypeNothing        = "nothing"
	// Passive stat effects — applied while an item is equipped.
	AdventureGameItemEffectEffectTypeWeaponDamage = "weapon_damage"
	AdventureGameItemEffectEffectTypeArmorDefense = "armor_defense"
)

// AdventureGameItemEffectEffectTypes is the set of all valid effect_type values.
var AdventureGameItemEffectEffectTypes = set.New(
	AdventureGameItemEffectEffectTypeInfo,
	AdventureGameItemEffectEffectTypeDamageTarget,
	AdventureGameItemEffectEffectTypeDamageWielder,
	AdventureGameItemEffectEffectTypeHealTarget,
	AdventureGameItemEffectEffectTypeHealWielder,
	AdventureGameItemEffectEffectTypeTeleport,
	AdventureGameItemEffectEffectTypeOpenLink,
	AdventureGameItemEffectEffectTypeCloseLink,
	AdventureGameItemEffectEffectTypeGiveItem,
	AdventureGameItemEffectEffectTypeRemoveItem,
	AdventureGameItemEffectEffectTypeSummonCreature,
	AdventureGameItemEffectEffectTypeNothing,
	AdventureGameItemEffectEffectTypeWeaponDamage,
	AdventureGameItemEffectEffectTypeArmorDefense,
)

// AdventureGameItemEffect defines what happens when a player performs an action on an item.
// Multiple rows for the same (item, action_type, required conditions) are allowed and all fire atomically.
type AdventureGameItemEffect struct {
	record.Record
	GameID                            string         `db:"game_id"`
	AdventureGameItemID               string         `db:"adventure_game_item_id"`
	ActionType                        string         `db:"action_type"`
	RequiredAdventureGameItemID       sql.NullString `db:"required_adventure_game_item_id"`
	RequiredAdventureGameLocationID   sql.NullString `db:"required_adventure_game_location_id"`
	ResultDescription                 string         `db:"result_description"`
	EffectType                        string         `db:"effect_type"`
	ResultAdventureGameItemID         sql.NullString `db:"result_adventure_game_item_id"`
	ResultAdventureGameLocationLinkID sql.NullString `db:"result_adventure_game_location_link_id"`
	ResultAdventureGameCreatureID     sql.NullString `db:"result_adventure_game_creature_id"`
	ResultAdventureGameLocationID     sql.NullString `db:"result_adventure_game_location_id"`
	ResultValueMin                    sql.NullInt32  `db:"result_value_min"`
	ResultValueMax                    sql.NullInt32  `db:"result_value_max"`
	IsRepeatable                      bool           `db:"is_repeatable"`
}

func (r *AdventureGameItemEffect) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameItemEffectGameID] = r.GameID
	args[FieldAdventureGameItemEffectAdventureGameItemID] = r.AdventureGameItemID
	args[FieldAdventureGameItemEffectActionType] = r.ActionType
	args[FieldAdventureGameItemEffectRequiredAdventureGameItemID] = r.RequiredAdventureGameItemID
	args[FieldAdventureGameItemEffectRequiredAdventureGameLocationID] = r.RequiredAdventureGameLocationID
	args[FieldAdventureGameItemEffectResultDescription] = r.ResultDescription
	args[FieldAdventureGameItemEffectEffectType] = r.EffectType
	args[FieldAdventureGameItemEffectResultAdventureGameItemID] = r.ResultAdventureGameItemID
	args[FieldAdventureGameItemEffectResultAdventureGameLocationLinkID] = r.ResultAdventureGameLocationLinkID
	args[FieldAdventureGameItemEffectResultAdventureGameCreatureID] = r.ResultAdventureGameCreatureID
	args[FieldAdventureGameItemEffectResultAdventureGameLocationID] = r.ResultAdventureGameLocationID
	args[FieldAdventureGameItemEffectResultValueMin] = r.ResultValueMin
	args[FieldAdventureGameItemEffectResultValueMax] = r.ResultValueMax
	args[FieldAdventureGameItemEffectIsRepeatable] = r.IsRepeatable
	return args
}
