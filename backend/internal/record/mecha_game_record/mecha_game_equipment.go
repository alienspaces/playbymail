package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameEquipment string = "mecha_game_equipment"
)

const (
	FieldMechaGameEquipmentID          string = "id"
	FieldMechaGameEquipmentGameID      string = "game_id"
	FieldMechaGameEquipmentName        string = "name"
	FieldMechaGameEquipmentDescription string = "description"
	FieldMechaGameEquipmentMountSize   string = "mount_size"
	FieldMechaGameEquipmentEffectKind  string = "effect_kind"
	FieldMechaGameEquipmentMagnitude   string = "magnitude"
	FieldMechaGameEquipmentHeatCost    string = "heat_cost"
	FieldMechaGameEquipmentCreatedAt   string = "created_at"
	FieldMechaGameEquipmentUpdatedAt   string = "updated_at"
	FieldMechaGameEquipmentDeletedAt   string = "deleted_at"
)

// Mount sizes reuse the weapon vocabulary (small / medium / large) so the
// Mountable abstraction can tally weapons and equipment against the same
// chassis slot budget.
const (
	EquipmentMountSizeSmall  string = "small"
	EquipmentMountSizeMedium string = "medium"
	EquipmentMountSizeLarge  string = "large"
)

// Effect kinds are a closed enum. Adding a new kind is a code change because
// each kind has a hardcoded "applied this turn?" predicate and effect path in
// the engine.
const (
	EquipmentEffectKindHeatSink          string = "heat_sink"
	EquipmentEffectKindTargetingComputer string = "targeting_computer"
	EquipmentEffectKindArmorUpgrade      string = "armor_upgrade"
	EquipmentEffectKindJumpJets          string = "jump_jets"
	EquipmentEffectKindECM               string = "ecm"
	EquipmentEffectKindAmmoBin           string = "ammo_bin"
)

// Per-kind sanity bounds for validator and client-side UI. Stacking is
// unbounded at runtime; these caps only apply per equipment row, not across
// a mech's loadout.
const (
	EquipmentMagnitudeMaxHeatSink          int = 20
	EquipmentMagnitudeMaxTargetingComputer int = 30
	EquipmentMagnitudeMaxArmorUpgrade      int = 200
	EquipmentMagnitudeMaxJumpJets          int = 5
	EquipmentMagnitudeMaxECM               int = 50
	EquipmentMagnitudeMaxAmmoBin           int = 200
)

// MagnitudeMaxForEffectKind returns the per-kind cap enforced by the validator.
// Unknown kinds return 0 to force the validator to reject them explicitly.
func MagnitudeMaxForEffectKind(kind string) int {
	switch kind {
	case EquipmentEffectKindHeatSink:
		return EquipmentMagnitudeMaxHeatSink
	case EquipmentEffectKindTargetingComputer:
		return EquipmentMagnitudeMaxTargetingComputer
	case EquipmentEffectKindArmorUpgrade:
		return EquipmentMagnitudeMaxArmorUpgrade
	case EquipmentEffectKindJumpJets:
		return EquipmentMagnitudeMaxJumpJets
	case EquipmentEffectKindECM:
		return EquipmentMagnitudeMaxECM
	case EquipmentEffectKindAmmoBin:
		return EquipmentMagnitudeMaxAmmoBin
	default:
		return 0
	}
}

// ValidEquipmentEffectKind returns true when kind is one of the six supported
// effect kinds.
func ValidEquipmentEffectKind(kind string) bool {
	return MagnitudeMaxForEffectKind(kind) > 0
}

// ValidEquipmentMountSize returns true when size is one of the three slot
// bands.
func ValidEquipmentMountSize(size string) bool {
	switch size {
	case EquipmentMountSizeSmall, EquipmentMountSizeMedium, EquipmentMountSizeLarge:
		return true
	}
	return false
}

type MechaGameEquipment struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	MountSize   string `db:"mount_size"`
	EffectKind  string `db:"effect_kind"`
	Magnitude   int    `db:"magnitude"`
	HeatCost    int    `db:"heat_cost"`
}

func (r *MechaGameEquipment) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameEquipmentGameID] = r.GameID
	args[FieldMechaGameEquipmentName] = r.Name
	args[FieldMechaGameEquipmentDescription] = r.Description
	args[FieldMechaGameEquipmentMountSize] = r.MountSize
	args[FieldMechaGameEquipmentEffectKind] = r.EffectKind
	args[FieldMechaGameEquipmentMagnitude] = r.Magnitude
	args[FieldMechaGameEquipmentHeatCost] = r.HeatCost
	return args
}
