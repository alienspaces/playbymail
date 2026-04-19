package mecha_game

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

const (
	supplyPointsPerTurn  = 2
	heatDissipationDenom = 3 // dissipate HeatCapacity / 3 per turn
	autoRepairPercent    = 25
)

// pilotSkillThresholds maps pilot skill level (index) to the total XP required.
// pilotSkillThresholds[skill] = minimum total XP to reach that skill level.
var pilotSkillThresholds = []int{0, 3, 8, 15, 24, 35, 48, 63, 80, 99}

// runEndOfTurn runs the end-of-turn lifecycle for all squads in a game instance:
//  1. Heat dissipation per mech
//  2. Auto-repair armor (field repairs)
//  3. XP application and pilot skill level-up
//  4. Complete refits (apply queued changes, clear is_refitting)
//  5. Supply point accrual for player squads
//  6. Append lifecycle TurnEvents to each squad instance
//
// xpMap is the XP earned by each mech this turn (mech instance ID → XP). May be nil.
func (p *MechaGame) runEndOfTurn(
	_ context.Context,
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	xpMap map[string]int,
) error {
	l = l.WithFunctionContext("MechaGame/runEndOfTurn")

	allMechInsts, err := p.Domain.GetManyMechaGameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameMechInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to load mech instances: %w", err)
	}

	allSquadInsts, err := p.Domain.GetManyMechaGameSquadInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSquadInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to load squad instances: %w", err)
	}

	// Events keyed by squad instance ID
	eventsBySquad := make(map[string][]turnsheet.TurnEvent)

	for _, inst := range allMechInsts {
		if inst.Status == mecha_game_record.MechInstanceStatusDestroyed {
			continue
		}

		chassisRec, err := p.Domain.GetMechaGameChassisRec(inst.MechaGameChassisID, nil)
		if err != nil {
			l.Warn("failed to get chassis for mech >%s<: %v", inst.ID, err)
			continue
		}

		// Load equipment for this mech so heat dissipation, armor ceilings,
		// and depot ammo refill can all use the effective values.
		var equipmentEntries []mecha_game_record.EquipmentConfigEntry
		if len(inst.EquipmentConfigJSON) > 0 {
			if err := json.Unmarshal(inst.EquipmentConfigJSON, &equipmentEntries); err != nil {
				l.Warn("failed to unmarshal equipment config for mech >%s<: %v", inst.ID, err)
				equipmentEntries = nil
			}
		}
		equipmentByID := make(map[string]*mecha_game_record.MechaGameEquipment, len(equipmentEntries))
		for _, entry := range equipmentEntries {
			if _, ok := equipmentByID[entry.EquipmentID]; ok {
				continue
			}
			eq, err := p.Domain.GetMechaGameEquipmentRec(entry.EquipmentID, nil)
			if err != nil {
				l.Warn("failed to load equipment >%s< for mech >%s<: %v", entry.EquipmentID, inst.ID, err)
				continue
			}
			equipmentByID[entry.EquipmentID] = eq
		}

		effects := AggregateEffects(equipmentEntries, equipmentByID, inst.IsRefitting)
		effectiveMaxArmor := EffectiveMaxArmor(chassisRec, effects)

		// 1a. Always-on equipment heat (heat_sink / armor_upgrade / ecm).
		// Applied every turn regardless of whether the mech entered combat
		// or moved. Refitting mechs accrue zero (AggregateEffects route).
		alwaysOnHeat := AlwaysOnHeatCost(equipmentEntries, equipmentByID, inst.IsRefitting)
		if alwaysOnHeat > 0 {
			inst.CurrentHeat += alwaysOnHeat
		}

		// 1. Heat dissipation — chassis baseline plus heat-sink bonus
		// (zeroed for refitting mechs by AggregateEffects). Applied after
		// the always-on heat accumulation above so the same turn's upkeep
		// and dissipation net out as a single movement.
		dissipate := chassisRec.HeatCapacity/heatDissipationDenom + effects.HeatDissipationBonus
		if inst.Status == mecha_game_record.MechInstanceStatusShutdown {
			// Shutdown resets heat and brings mech back online
			inst.CurrentHeat = 0
			inst.Status = mecha_game_record.MechInstanceStatusOperational
			appendLifecycleEvent(eventsBySquad, inst.MechaGameSquadInstanceID,
				fmt.Sprintf("%s emergency shutdown complete — back online.", inst.Callsign))
		} else {
			prev := inst.CurrentHeat
			inst.CurrentHeat -= dissipate
			if inst.CurrentHeat < 0 {
				inst.CurrentHeat = 0
			}
			if prev > 0 && inst.CurrentHeat < prev {
				appendLifecycleEvent(eventsBySquad, inst.MechaGameSquadInstanceID,
					fmt.Sprintf("%s heat dissipated from %d to %d.", inst.Callsign, prev, inst.CurrentHeat))
			}
		}

		// 2. Field auto-repair (armor only, no structure). Ceiling and the
		// 25% repair base both use effectiveMaxArmor so armor-upgrade
		// magnitude is honored consistently.
		if !inst.IsRefitting && inst.CurrentArmor < effectiveMaxArmor {
			repairAmt := int(math.Ceil(float64(effectiveMaxArmor) * autoRepairPercent / 100))
			prevArmor := inst.CurrentArmor
			inst.CurrentArmor += repairAmt
			if inst.CurrentArmor > effectiveMaxArmor {
				inst.CurrentArmor = effectiveMaxArmor
			}
			if inst.CurrentArmor > prevArmor {
				if inst.CurrentArmor == effectiveMaxArmor {
					inst.Status = mecha_game_record.MechInstanceStatusOperational
				}
				appendLifecycleEvent(eventsBySquad, inst.MechaGameSquadInstanceID,
					fmt.Sprintf("%s field repairs restored %d armor (%d/%d).",
						inst.Callsign, inst.CurrentArmor-prevArmor, inst.CurrentArmor, effectiveMaxArmor))
			}
		}

		// 2b. Depot ammo refill. Mechs parked on a depot (starting sector)
		// top up their shared ammo pool to MaxAmmoCapacity. This is treated
		// as a crew action, so it runs even for refitting mechs.
		if err := p.refillAmmoAtDepot(l, inst, equipmentEntries, equipmentByID, eventsBySquad); err != nil {
			l.Warn("failed to refill ammo for mech >%s<: %v", inst.ID, err)
		}

		// 3. Apply XP and check for pilot skill level-up.
		if xpMap != nil {
			if earned := xpMap[inst.ID]; earned > 0 {
				inst.ExperiencePoints += earned
				for nextSkill := inst.PilotSkill + 1; nextSkill < len(pilotSkillThresholds); nextSkill++ {
					if inst.ExperiencePoints >= pilotSkillThresholds[nextSkill] {
						inst.PilotSkill = nextSkill
						appendLifecycleEvent(eventsBySquad, inst.MechaGameSquadInstanceID,
							fmt.Sprintf("%s pilot skill increased to %d!", inst.Callsign, inst.PilotSkill))
					} else {
						break
					}
				}
			}
		}

		// 4. Complete refits (clear flag; actual changes are applied by the
		// management processor before this runs)
		if inst.IsRefitting {
			inst.IsRefitting = false
			appendLifecycleEvent(eventsBySquad, inst.MechaGameSquadInstanceID,
				fmt.Sprintf("%s refit complete.", inst.Callsign))
		}

		if _, err := p.Domain.UpdateMechaGameMechInstanceRec(inst); err != nil {
			l.Warn("failed to update mech instance >%s< after end-of-turn: %v", inst.ID, err)
		}
	}

	// 5. Supply point accrual for player squads + persist events
	for _, squadInst := range allSquadInsts {
		// Only player-owned squads accrue supply points
		if squadInst.GameSubscriptionInstanceID.Valid {
			squadInst.SupplyPoints += supplyPointsPerTurn
			appendLifecycleEvent(eventsBySquad, squadInst.ID,
				fmt.Sprintf("Squad received %d supply points (%d total).",
					supplyPointsPerTurn, squadInst.SupplyPoints))
		}

		for _, evt := range eventsBySquad[squadInst.ID] {
			if err := turnsheet.AppendMechaGameTurnEvent(squadInst, evt); err != nil {
				l.Warn("failed to append end-of-turn event for squad >%s<: %v", squadInst.ID, err)
			}
		}

		if _, err := p.Domain.UpdateMechaGameSquadInstanceRec(squadInst); err != nil {
			l.Warn("failed to update squad instance >%s< after end-of-turn: %v", squadInst.ID, err)
		}
	}

	return nil
}

// refillAmmoAtDepot tops a mech's AmmoRemaining back up to its configured
// MaxAmmoCapacity when the mech is parked on a starting (depot) sector.
// The refill is a crew action, so it applies even when the mech is
// refitting. Mechs with no ammo-consuming loadout are no-ops.
func (p *MechaGame) refillAmmoAtDepot(
	l logger.Logger,
	inst *mecha_game_record.MechaGameMechInstance,
	equipmentEntries []mecha_game_record.EquipmentConfigEntry,
	equipmentByID map[string]*mecha_game_record.MechaGameEquipment,
	eventsBySquad map[string][]turnsheet.TurnEvent,
) error {
	sectorInst, err := p.Domain.GetMechaGameSectorInstanceRec(inst.MechaGameSectorInstanceID, nil)
	if err != nil {
		return fmt.Errorf("failed to get sector instance: %w", err)
	}
	sectorDesign, err := p.Domain.GetMechaGameSectorRec(sectorInst.MechaGameSectorID, nil)
	if err != nil {
		return fmt.Errorf("failed to get sector design: %w", err)
	}
	if !sectorDesign.IsStartingSector {
		return nil
	}

	var weaponEntries []mecha_game_record.WeaponConfigEntry
	if len(inst.WeaponConfigJSON) > 0 {
		if err := json.Unmarshal(inst.WeaponConfigJSON, &weaponEntries); err != nil {
			return fmt.Errorf("failed to unmarshal weapon config: %w", err)
		}
	}

	weaponByID := make(map[string]*mecha_game_record.MechaGameWeapon, len(weaponEntries))
	for _, entry := range weaponEntries {
		if entry.WeaponID == "" {
			continue
		}
		if _, ok := weaponByID[entry.WeaponID]; ok {
			continue
		}
		w, err := p.Domain.GetMechaGameWeaponRec(entry.WeaponID, nil)
		if err != nil {
			l.Warn("failed to load weapon >%s< for ammo refill: %v", entry.WeaponID, err)
			continue
		}
		weaponByID[entry.WeaponID] = w
	}

	maxAmmo := MaxAmmoCapacity(weaponEntries, weaponByID, equipmentEntries, equipmentByID)
	if maxAmmo <= 0 || inst.AmmoRemaining >= maxAmmo {
		return nil
	}

	delta := maxAmmo - inst.AmmoRemaining
	inst.AmmoRemaining = maxAmmo
	appendLifecycleEvent(eventsBySquad, inst.MechaGameSquadInstanceID,
		fmt.Sprintf("%s rearmed at depot (+%d ammo, %d total).",
			inst.Callsign, delta, inst.AmmoRemaining))
	return nil
}

func appendLifecycleEvent(
	eventsBySquad map[string][]turnsheet.TurnEvent,
	squadInstanceID string,
	message string,
) {
	eventsBySquad[squadInstanceID] = append(
		eventsBySquad[squadInstanceID],
		turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategorySystem,
			Icon:     turnsheet.TurnEventIconSystem,
			Message:  message,
		},
	)
}
