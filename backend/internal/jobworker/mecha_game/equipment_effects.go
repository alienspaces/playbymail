package mecha_game

import (
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

// Effects is a package-local alias for the shared
// domain.MechaGameEquipmentEffects type, kept so the jobworker call sites
// can use a short name. All semantics are identical to the domain type.
type Effects = domain.MechaGameEquipmentEffects

// AggregateEffects is a thin wrapper over
// domain.AggregateMechaGameEquipmentEffects. See that function for full
// semantics (including the "refitting returns zero" rule).
func AggregateEffects(
	entries []mecha_game_record.EquipmentConfigEntry,
	byID map[string]*mecha_game_record.MechaGameEquipment,
	refitting bool,
) Effects {
	return domain.AggregateMechaGameEquipmentEffects(entries, byID, refitting)
}

// EffectiveMaxArmor wraps domain.EffectiveMechaGameMaxArmor.
func EffectiveMaxArmor(chassis *mecha_game_record.MechaGameChassis, effects Effects) int {
	return domain.EffectiveMechaGameMaxArmor(chassis, effects)
}

// EffectiveSpeed wraps domain.EffectiveMechaGameSpeed.
func EffectiveSpeed(chassis *mecha_game_record.MechaGameChassis, effects Effects) int {
	return domain.EffectiveMechaGameSpeed(chassis, effects)
}

// MaxAmmoCapacity wraps domain.MaxMechaGameAmmoCapacity.
func MaxAmmoCapacity(
	weaponEntries []mecha_game_record.WeaponConfigEntry,
	weaponByID map[string]*mecha_game_record.MechaGameWeapon,
	equipmentEntries []mecha_game_record.EquipmentConfigEntry,
	equipmentByID map[string]*mecha_game_record.MechaGameEquipment,
) int {
	return domain.MaxMechaGameAmmoCapacity(weaponEntries, weaponByID, equipmentEntries, equipmentByID)
}

// CombatHeatCost wraps domain.MechaGameEquipmentCombatHeatCost.
func CombatHeatCost(
	entries []mecha_game_record.EquipmentConfigEntry,
	byID map[string]*mecha_game_record.MechaGameEquipment,
	refitting bool,
	didAttack, didFireAmmoWeapon bool,
) int {
	return domain.MechaGameEquipmentCombatHeatCost(entries, byID, refitting, didAttack, didFireAmmoWeapon)
}

// AlwaysOnHeatCost wraps domain.MechaGameEquipmentAlwaysOnHeatCost.
func AlwaysOnHeatCost(
	entries []mecha_game_record.EquipmentConfigEntry,
	byID map[string]*mecha_game_record.MechaGameEquipment,
	refitting bool,
) int {
	return domain.MechaGameEquipmentAlwaysOnHeatCost(entries, byID, refitting)
}
