package mecha

import (
	"context"
	"fmt"
	"math"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
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
func (p *Mecha) runEndOfTurn(
	_ context.Context,
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	xpMap map[string]int,
) error {
	l = l.WithFunctionContext("Mecha/runEndOfTurn")

	allMechInsts, err := p.Domain.GetManyMechaMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaMechInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to load mech instances: %w", err)
	}

	allSquadInsts, err := p.Domain.GetManyMechaSquadInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSquadInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to load squad instances: %w", err)
	}

	// Events keyed by squad instance ID
	eventsBySquad := make(map[string][]turnsheet.TurnEvent)

	for _, inst := range allMechInsts {
		if inst.Status == mecha_record.MechInstanceStatusDestroyed {
			continue
		}

		chassisRec, err := p.Domain.GetMechaChassisRec(inst.MechaChassisID, nil)
		if err != nil {
			l.Warn("failed to get chassis for mech >%s<: %v", inst.ID, err)
			continue
		}

		// 1. Heat dissipation
		dissipate := chassisRec.HeatCapacity / heatDissipationDenom
		if inst.Status == mecha_record.MechInstanceStatusShutdown {
			// Shutdown resets heat and brings mech back online
			inst.CurrentHeat = 0
			inst.Status = mecha_record.MechInstanceStatusOperational
			appendLifecycleEvent(eventsBySquad, inst.MechaSquadInstanceID,
				fmt.Sprintf("%s emergency shutdown complete — back online.", inst.Callsign))
		} else {
			prev := inst.CurrentHeat
			inst.CurrentHeat -= dissipate
			if inst.CurrentHeat < 0 {
				inst.CurrentHeat = 0
			}
			if prev > 0 && inst.CurrentHeat < prev {
				appendLifecycleEvent(eventsBySquad, inst.MechaSquadInstanceID,
					fmt.Sprintf("%s heat dissipated from %d to %d.", inst.Callsign, prev, inst.CurrentHeat))
			}
		}

		// 2. Field auto-repair (armor only, no structure)
		if !inst.IsRefitting && inst.CurrentArmor < chassisRec.ArmorPoints {
			repairAmt := int(math.Ceil(float64(chassisRec.ArmorPoints) * autoRepairPercent / 100))
			prevArmor := inst.CurrentArmor
			inst.CurrentArmor += repairAmt
			if inst.CurrentArmor > chassisRec.ArmorPoints {
				inst.CurrentArmor = chassisRec.ArmorPoints
			}
			if inst.CurrentArmor > prevArmor {
				if inst.CurrentArmor == chassisRec.ArmorPoints {
					inst.Status = mecha_record.MechInstanceStatusOperational
				}
				appendLifecycleEvent(eventsBySquad, inst.MechaSquadInstanceID,
					fmt.Sprintf("%s field repairs restored %d armor (%d/%d).",
						inst.Callsign, inst.CurrentArmor-prevArmor, inst.CurrentArmor, chassisRec.ArmorPoints))
			}
		}

		// 3. Apply XP and check for pilot skill level-up.
		if xpMap != nil {
			if earned := xpMap[inst.ID]; earned > 0 {
				inst.ExperiencePoints += earned
				for nextSkill := inst.PilotSkill + 1; nextSkill < len(pilotSkillThresholds); nextSkill++ {
					if inst.ExperiencePoints >= pilotSkillThresholds[nextSkill] {
						inst.PilotSkill = nextSkill
						appendLifecycleEvent(eventsBySquad, inst.MechaSquadInstanceID,
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
			appendLifecycleEvent(eventsBySquad, inst.MechaSquadInstanceID,
				fmt.Sprintf("%s refit complete.", inst.Callsign))
		}

		if _, err := p.Domain.UpdateMechaMechInstanceRec(inst); err != nil {
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
			if err := turnsheet.AppendMechaTurnEvent(squadInst, evt); err != nil {
				l.Warn("failed to append end-of-turn event for squad >%s<: %v", squadInst.ID, err)
			}
		}

		if _, err := p.Domain.UpdateMechaSquadInstanceRec(squadInst); err != nil {
			l.Warn("failed to update squad instance >%s< after end-of-turn: %v", squadInst.ID, err)
		}
	}

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
