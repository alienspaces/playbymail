package domain

import (
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

// LoadMechaGameEquipmentByID resolves every distinct EquipmentID referenced
// by the supplied equipment config entries to its MechaGameEquipment record,
// returning a lookup map suitable for passing to
// AggregateMechaGameEquipmentEffects / MechaGameEquipmentHeatCostsThisTurn.
//
// Entries whose equipment lookup fails are skipped with the error captured
// in the returned error so callers can decide whether to treat the record
// as missing (skip) or abort. Missing records are not fatal for the map;
// callers that already tolerate a nil entry in byID will still see that
// equipment ignored by the aggregate helpers.
func (m *Domain) LoadMechaGameEquipmentByID(
	entries []mecha_game_record.EquipmentConfigEntry,
) (map[string]*mecha_game_record.MechaGameEquipment, error) {
	byID := make(map[string]*mecha_game_record.MechaGameEquipment, len(entries))
	for _, entry := range entries {
		if entry.EquipmentID == "" {
			continue
		}
		if _, ok := byID[entry.EquipmentID]; ok {
			continue
		}
		eq, err := m.GetMechaGameEquipmentRec(entry.EquipmentID, nil)
		if err != nil {
			return byID, err
		}
		byID[entry.EquipmentID] = eq
	}
	return byID, nil
}

// MechaGameEquipmentEffects is the aggregated per-mech equipment effect
// state. All fields are strictly additive — a zero value has no influence on
// any engine calculation, which is what AggregateMechaGameEquipmentEffects
// returns when the mech is refitting (powered down).
type MechaGameEquipmentEffects struct {
	// HeatDissipationBonus is added to the chassis's baseline heat
	// dissipation at end-of-turn.
	HeatDissipationBonus int
	// HitChanceBonus (targeting computer) is added to attacker hit chance
	// after the chassis/pilot baseline. The combat resolver still clamps the
	// final value to the [0, 95] playability band.
	HitChanceBonus int
	// ArmorBonus (armor upgrade) raises effective max armor used for
	// initialization, auto-repair ceiling, and the 25% auto-repair base.
	ArmorBonus int
	// SpeedBonus (jump jets) adds hops to the chassis's base speed used by
	// the orders processor and AI movement budget.
	SpeedBonus int
	// CoverBonus (ECM) is added to the cover-modifier term applied to attacks
	// *against* this mech. Combined with sector cover at the defender.
	CoverBonus int
}

// AggregateMechaGameEquipmentEffects sums the per-kind bonuses from the
// mech's equipment config. If the mech is refitting it is considered powered
// down and returns a zero value — all equipment enhancements go offline.
// Only ammo refill at a depot (handled at end-of-turn) still runs for a
// refitting mech, treated as a crew action rather than an equipment effect.
func AggregateMechaGameEquipmentEffects(
	entries []mecha_game_record.EquipmentConfigEntry,
	byID map[string]*mecha_game_record.MechaGameEquipment,
	refitting bool,
) MechaGameEquipmentEffects {
	if refitting {
		return MechaGameEquipmentEffects{}
	}

	var e MechaGameEquipmentEffects
	for _, entry := range entries {
		eq := byID[entry.EquipmentID]
		if eq == nil {
			continue
		}
		switch eq.EffectKind {
		case mecha_game_record.EquipmentEffectKindHeatSink:
			e.HeatDissipationBonus += eq.Magnitude
		case mecha_game_record.EquipmentEffectKindTargetingComputer:
			e.HitChanceBonus += eq.Magnitude
		case mecha_game_record.EquipmentEffectKindArmorUpgrade:
			e.ArmorBonus += eq.Magnitude
		case mecha_game_record.EquipmentEffectKindJumpJets:
			e.SpeedBonus += eq.Magnitude
		case mecha_game_record.EquipmentEffectKindECM:
			e.CoverBonus += eq.Magnitude
		case mecha_game_record.EquipmentEffectKindAmmoBin:
			// ammo_bin contributes to AmmoRemaining at initial clone and on
			// depot refill; it is not an aggregate effects bonus.
		}
	}
	return e
}

// EffectiveMechaGameMaxArmor returns the chassis's base armor plus the
// aggregated armor-upgrade bonus. Used for initial clone, auto-repair
// ceiling, and the 25%-of-max repair base so all three stay in sync with the
// armor upgrade.
func EffectiveMechaGameMaxArmor(chassis *mecha_game_record.MechaGameChassis, effects MechaGameEquipmentEffects) int {
	if chassis == nil {
		return effects.ArmorBonus
	}
	return chassis.ArmorPoints + effects.ArmorBonus
}

// EffectiveMechaGameSpeed returns the chassis's base speed plus the
// aggregated jump jets bonus. Used by the orders processor for entry
// display, reachable-sector BFS, and move validation, and by the AI
// movement planner.
func EffectiveMechaGameSpeed(chassis *mecha_game_record.MechaGameChassis, effects MechaGameEquipmentEffects) int {
	if chassis == nil {
		return effects.SpeedBonus
	}
	return chassis.Speed + effects.SpeedBonus
}

// MaxMechaGameAmmoCapacity returns the mech's full ammo pool capacity: the
// sum of ammo_capacity from all equipped weapons plus the magnitude of every
// ammo_bin equipment entry. This is the value used to seed
// MechaGameMechInstance.AmmoRemaining at game start and to refill it when
// the mech is at a depot at end-of-turn.
//
// Refit state does not affect capacity — refills at a depot still work
// because refilling is a crew action, not an equipment effect.
func MaxMechaGameAmmoCapacity(
	weaponEntries []mecha_game_record.WeaponConfigEntry,
	weaponByID map[string]*mecha_game_record.MechaGameWeapon,
	equipmentEntries []mecha_game_record.EquipmentConfigEntry,
	equipmentByID map[string]*mecha_game_record.MechaGameEquipment,
) int {
	capacity := 0
	for _, entry := range weaponEntries {
		w := weaponByID[entry.WeaponID]
		if w == nil {
			continue
		}
		capacity += w.AmmoCapacity
	}
	for _, entry := range equipmentEntries {
		eq := equipmentByID[entry.EquipmentID]
		if eq == nil {
			continue
		}
		if eq.EffectKind == mecha_game_record.EquipmentEffectKindAmmoBin {
			capacity += eq.Magnitude
		}
	}
	return capacity
}

// MechaGameEquipmentJumpJetHeatCost returns the heat_cost contribution of
// any jump_jets equipment that fires when the mech moves more hops than
// the chassis base speed this turn. Returns zero when the mech is
// refitting or carries no jump_jets equipment with heat_cost > 0.
//
// This is split out from MechaGameEquipmentHeatCostsThisTurn so the orders
// processor (which is the only place movement happens) can apply jump-jets
// heat independently of the always-on / combat predicates, which are
// accumulated separately in combat resolution.
func MechaGameEquipmentJumpJetHeatCost(
	entries []mecha_game_record.EquipmentConfigEntry,
	byID map[string]*mecha_game_record.MechaGameEquipment,
	refitting bool,
) int {
	if refitting {
		return 0
	}
	total := 0
	for _, entry := range entries {
		eq := byID[entry.EquipmentID]
		if eq == nil || eq.HeatCost == 0 {
			continue
		}
		if eq.EffectKind == mecha_game_record.EquipmentEffectKindJumpJets {
			total += eq.HeatCost
		}
	}
	return total
}

// MechaGameEquipmentCombatHeatCost returns the sum of heat_cost from
// equipment whose "fires this turn" predicate is tied to a combat action:
//
//   - targeting_computer  : mech declared an attack this turn
//   - ammo_bin            : mech fired a weapon with ammo_capacity > 0 this
//     turn
//
// Returns zero when the mech is refitting. This helper is intended for
// combat_resolution which already knows whether a mech attacked and whether
// it fired an ammo-consuming weapon.
func MechaGameEquipmentCombatHeatCost(
	entries []mecha_game_record.EquipmentConfigEntry,
	byID map[string]*mecha_game_record.MechaGameEquipment,
	refitting bool,
	didAttack, didFireAmmoWeapon bool,
) int {
	if refitting {
		return 0
	}
	total := 0
	for _, entry := range entries {
		eq := byID[entry.EquipmentID]
		if eq == nil || eq.HeatCost == 0 {
			continue
		}
		switch eq.EffectKind {
		case mecha_game_record.EquipmentEffectKindTargetingComputer:
			if didAttack {
				total += eq.HeatCost
			}
		case mecha_game_record.EquipmentEffectKindAmmoBin:
			if didFireAmmoWeapon {
				total += eq.HeatCost
			}
		}
	}
	return total
}

// MechaGameEquipmentAlwaysOnHeatCost returns the sum of heat_cost from
// equipment whose "fires this turn" predicate is always-on (while not
// refitting):
//
//   - heat_sink     : always
//   - armor_upgrade : always
//   - ecm           : always
//
// Returns zero when the mech is refitting. This helper is intended for
// end_of_turn, which applies always-on heat once per turn to every mech
// with equipment, regardless of whether that mech entered combat or moved.
func MechaGameEquipmentAlwaysOnHeatCost(
	entries []mecha_game_record.EquipmentConfigEntry,
	byID map[string]*mecha_game_record.MechaGameEquipment,
	refitting bool,
) int {
	if refitting {
		return 0
	}
	total := 0
	for _, entry := range entries {
		eq := byID[entry.EquipmentID]
		if eq == nil || eq.HeatCost == 0 {
			continue
		}
		switch eq.EffectKind {
		case mecha_game_record.EquipmentEffectKindHeatSink,
			mecha_game_record.EquipmentEffectKindArmorUpgrade,
			mecha_game_record.EquipmentEffectKindECM:
			total += eq.HeatCost
		}
	}
	return total
}
