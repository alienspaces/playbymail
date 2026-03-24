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

// runEndOfTurn runs the end-of-turn lifecycle for all lances in a game instance:
//  1. Heat dissipation per mech
//  2. Auto-repair armor (field repairs)
//  3. Complete refits (apply queued changes, clear is_refitting)
//  4. Supply point accrual for player lances
//  5. Append lifecycle TurnEvents to each lance instance
func (p *Mecha) runEndOfTurn(
	_ context.Context,
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
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

	allLanceInsts, err := p.Domain.GetManyMechaLanceInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaLanceInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to load lance instances: %w", err)
	}

	// Map lance instance ID → record
	lanceMap := make(map[string]*mecha_record.MechaLanceInstance, len(allLanceInsts))
	for _, li := range allLanceInsts {
		lanceMap[li.ID] = li
	}

	// Events keyed by lance instance ID
	eventsByLance := make(map[string][]turnsheet.TurnEvent)

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
			appendLifecycleEvent(eventsByLance, inst.MechaLanceInstanceID,
				fmt.Sprintf("%s emergency shutdown complete — back online.", inst.Callsign))
		} else {
			prev := inst.CurrentHeat
			inst.CurrentHeat -= dissipate
			if inst.CurrentHeat < 0 {
				inst.CurrentHeat = 0
			}
			if prev > 0 && inst.CurrentHeat < prev {
				appendLifecycleEvent(eventsByLance, inst.MechaLanceInstanceID,
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
				appendLifecycleEvent(eventsByLance, inst.MechaLanceInstanceID,
					fmt.Sprintf("%s field repairs restored %d armor (%d/%d).",
						inst.Callsign, inst.CurrentArmor-prevArmor, inst.CurrentArmor, chassisRec.ArmorPoints))
			}
		}

		// 3. Complete refits (clear flag; actual changes are applied by the
		// management processor before this runs)
		if inst.IsRefitting {
			inst.IsRefitting = false
			appendLifecycleEvent(eventsByLance, inst.MechaLanceInstanceID,
				fmt.Sprintf("%s refit complete.", inst.Callsign))
		}

		if _, err := p.Domain.UpdateMechaMechInstanceRec(inst); err != nil {
			l.Warn("failed to update mech instance >%s< after end-of-turn: %v", inst.ID, err)
		}
	}

	// 4. Supply point accrual for player lances + persist events
	for _, lanceInst := range allLanceInsts {
		// Only player-owned lances accrue supply points
		if lanceInst.GameSubscriptionInstanceID.Valid {
			lanceInst.SupplyPoints += supplyPointsPerTurn
			appendLifecycleEvent(eventsByLance, lanceInst.ID,
				fmt.Sprintf("Lance received %d supply points (%d total).",
					supplyPointsPerTurn, lanceInst.SupplyPoints))
		}

		// Append any events collected for this lance
		for _, evt := range eventsByLance[lanceInst.ID] {
			if err := turnsheet.AppendMechaTurnEvent(lanceInst, evt); err != nil {
				l.Warn("failed to append end-of-turn event for lance >%s<: %v", lanceInst.ID, err)
			}
		}

		if _, err := p.Domain.UpdateMechaLanceInstanceRec(lanceInst); err != nil {
			l.Warn("failed to update lance instance >%s< after end-of-turn: %v", lanceInst.ID, err)
		}
	}

	return nil
}

func appendLifecycleEvent(
	eventsByLance map[string][]turnsheet.TurnEvent,
	lanceInstanceID string,
	message string,
) {
	eventsByLance[lanceInstanceID] = append(
		eventsByLance[lanceInstanceID],
		turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategorySystem,
			Icon:     turnsheet.TurnEventIconSystem,
			Message:  message,
		},
	)
}
