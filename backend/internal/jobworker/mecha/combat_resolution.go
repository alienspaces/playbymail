package mecha

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/mecha/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// AttackDeclaration is an alias for the type defined in turn_sheet_processor.
type AttackDeclaration = turn_sheet_processor.AttackDeclaration

// mechSnapshot captures a mech's state before combat begins, so all attacks
// resolve against the same starting state (simultaneous resolution).
type mechSnapshot struct {
	Instance         *mecha_record.MechaMechInstance
	LanceInstanceID  string
	SectorInstanceID string
	Weapons          []mecha_record.WeaponConfigEntry
}

// rangeDistance returns the number of sector hops between two sector instances
// using the cached sector graph. Returns 0 for same sector, 1 for adjacent,
// and 999 if no path is found within 3 hops (effectively out of range).
func rangeDistance(fromID, toID string, sectors []*sectorState) int {
	if fromID == toID {
		return 0
	}
	// BFS up to depth 3
	type node struct {
		id    string
		depth int
	}
	visited := map[string]bool{fromID: true}
	queue := []node{{id: fromID, depth: 0}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.depth >= 3 {
			continue
		}
		for _, sec := range sectors {
			if sec.Instance.ID != cur.id {
				continue
			}
			for _, dest := range sec.LinkDestInstanceIDs {
				if dest == toID {
					return cur.depth + 1
				}
				if !visited[dest] {
					visited[dest] = true
					queue = append(queue, node{id: dest, depth: cur.depth + 1})
				}
			}
		}
	}
	return 999
}

// weaponCanFire returns whether a weapon's range band is valid for the given
// sector distance.
//   - same sector (0): all range bands fire
//   - adjacent (1): medium and long range bands only
//   - 2+ sectors: nothing fires
func weaponCanFire(rangeBand string, distance int) bool {
	switch distance {
	case 0:
		return true
	case 1:
		return rangeBand == mecha_record.WeaponRangeBandMedium ||
			rangeBand == mecha_record.WeaponRangeBandLong
	default:
		return false
	}
}

// hitChance returns the probability of a single weapon hit (0–100).
// Base 50% + 5% per pilot skill point.
func hitChance(pilotSkill int) int {
	chance := 50 + pilotSkill*5
	if chance > 95 {
		chance = 95
	}
	return chance
}

// pendingDamage accumulates damage from all attacks before applying (simultaneous resolution).
type pendingDamage struct {
	armorDmg     int
	structureDmg int
}

// resolveCombat runs combat resolution for a game instance. It must be called
// after all movement (player and AI) has been applied. The sector graph from
// the decision engine context is rebuilt here for range calculations.
func (p *Mecha) resolveCombat(
	_ context.Context,
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	attacks []AttackDeclaration,
) error {
	if len(attacks) == 0 {
		l.Info("no attack declarations for game instance >%s< — skipping combat", gameInstanceRec.ID)
		return nil
	}

	allMechInsts, err := p.Domain.GetManyMechaMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaMechInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		l.Warn("failed to load mech instances: %v", err)
		return err
	}

	snapshots := buildMechSnapshots(l, allMechInsts)

	sectors, err := p.buildSectorGraph(gameInstanceRec.ID)
	if err != nil {
		l.Warn("failed to build sector graph: %v", err)
		return err
	}

	seed := int64(gameInstanceRec.CurrentTurn)
	for _, b := range []byte(gameInstanceRec.ID) {
		seed += int64(b)
	}
	rng := rand.New(rand.NewSource(seed)) //nolint:gosec

	damageMap := make(map[string]*pendingDamage, len(allMechInsts))
	heatMap := make(map[string]int, len(allMechInsts))
	eventsByLance := make(map[string][]turnsheet.TurnEvent)

	p.resolveAttacks(l, attacks, snapshots, sectors, rng, damageMap, heatMap, eventsByLance)
	p.applyPendingDamage(l, damageMap, snapshots, attacks, eventsByLance)
	p.applyPendingHeat(l, heatMap, snapshots, eventsByLance)
	p.persistMechChanges(l, snapshots, damageMap, heatMap)

	if err := p.appendCombatEventsToLances(l, eventsByLance); err != nil {
		l.Warn("failed to persist combat events: %v", err)
		return err
	}

	return nil
}

func buildMechSnapshots(l logger.Logger, insts []*mecha_record.MechaMechInstance) map[string]*mechSnapshot {
	snapshots := make(map[string]*mechSnapshot, len(insts))
	for _, inst := range insts {
		var weapons []mecha_record.WeaponConfigEntry
		if len(inst.WeaponConfigJSON) > 0 {
			if err := json.Unmarshal(inst.WeaponConfigJSON, &weapons); err != nil {
				l.Warn("failed to unmarshal weapon config for mech >%s<: %v", inst.ID, err)
			}
		}
		snapshots[inst.ID] = &mechSnapshot{
			Instance:         inst,
			LanceInstanceID:  inst.MechaLanceInstanceID,
			SectorInstanceID: inst.MechaSectorInstanceID,
			Weapons:          weapons,
		}
	}
	return snapshots
}

func (p *Mecha) buildSectorGraph(gameInstanceID string) ([]*sectorState, error) {
	sectorInsts, err := p.Domain.GetManyMechaSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSectorInstanceGameInstanceID, Val: gameInstanceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load sector instances: %w", err)
	}
	var sectors []*sectorState
	for _, si := range sectorInsts {
		links, _ := p.Domain.GetManyMechaSectorLinkRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaSectorLinkFromMechaSectorID, Val: si.MechaSectorID},
			},
		})
		var linkDestIDs []string
		for _, lnk := range links {
			for _, other := range sectorInsts {
				if other.MechaSectorID == lnk.ToMechaSectorID {
					linkDestIDs = append(linkDestIDs, other.ID)
					break
				}
			}
		}
		sectors = append(sectors, &sectorState{
			Instance:            si,
			LinkDestInstanceIDs: linkDestIDs,
		})
	}
	return sectors, nil
}

func (p *Mecha) resolveAttacks(
	l logger.Logger,
	attacks []AttackDeclaration,
	snapshots map[string]*mechSnapshot,
	sectors []*sectorState,
	rng *rand.Rand,
	damageMap map[string]*pendingDamage,
	heatMap map[string]int,
	eventsByLance map[string][]turnsheet.TurnEvent,
) {
	for _, atk := range attacks {
		attacker, ok := snapshots[atk.AttackerMechInstanceID]
		if !ok {
			l.Warn("attacker mech >%s< not found in snapshot — skipping", atk.AttackerMechInstanceID)
			continue
		}
		target, ok := snapshots[atk.TargetMechInstanceID]
		if !ok {
			l.Warn("target mech >%s< not found in snapshot — skipping", atk.TargetMechInstanceID)
			continue
		}
		if attacker.Instance.Status == mecha_record.MechInstanceStatusDestroyed ||
			attacker.Instance.Status == mecha_record.MechInstanceStatusShutdown {
			l.Debug("attacker >%s< is %s — skipping attack", attacker.Instance.Callsign, attacker.Instance.Status)
			continue
		}
		if target.Instance.Status == mecha_record.MechInstanceStatusDestroyed {
			l.Debug("target >%s< is already destroyed — skipping attack", target.Instance.Callsign)
			continue
		}
		dist := rangeDistance(attacker.SectorInstanceID, target.SectorInstanceID, sectors)
		if dist >= 2 {
			l.Info("%s fired at %s but target is out of range (distance %d)",
				attacker.Instance.Callsign, target.Instance.Callsign, dist)
			appendCombatEvent(eventsByLance, attacker.LanceInstanceID,
				fmt.Sprintf("%s fired at %s — out of range.",
					attacker.Instance.Callsign, target.Instance.Callsign))
			continue
		}
		if len(attacker.Weapons) == 0 {
			l.Debug("attacker >%s< has no weapons — skipping attack", attacker.Instance.Callsign)
			continue
		}
		totalDmg := p.fireWeapons(l, atk, attacker, target, dist, rng, heatMap, eventsByLance)
		if totalDmg > 0 {
			accumulateDamage(atk.TargetMechInstanceID, totalDmg, snapshots, damageMap)
		}
	}
}

func (p *Mecha) fireWeapons(
	l logger.Logger,
	atk AttackDeclaration,
	attacker, target *mechSnapshot,
	dist int,
	rng *rand.Rand,
	heatMap map[string]int,
	eventsByLance map[string][]turnsheet.TurnEvent,
) int {
	chance := hitChance(attacker.Instance.PilotSkill)
	totalDmg := 0
	for _, slot := range attacker.Weapons {
		if slot.WeaponID == "" {
			continue
		}
		weaponRec, err := p.Domain.GetMechaWeaponRec(slot.WeaponID, nil)
		if err != nil {
			l.Warn("failed to load weapon >%s<: %v", slot.WeaponID, err)
			continue
		}
		if !weaponCanFire(weaponRec.RangeBand, dist) {
			l.Debug("%s weapon %s cannot reach target at distance %d",
				attacker.Instance.Callsign, weaponRec.Name, dist)
			continue
		}
		heatMap[atk.AttackerMechInstanceID] += weaponRec.HeatCost
		roll := rng.Intn(100)
		if roll < chance {
			totalDmg += weaponRec.Damage
			l.Info("%s: %s hit %s with %s for %d damage (roll %d < %d%%)",
				attacker.Instance.Callsign, weaponRec.Name,
				target.Instance.Callsign, weaponRec.Name,
				weaponRec.Damage, roll, chance)
			appendCombatEvent(eventsByLance, attacker.LanceInstanceID,
				fmt.Sprintf("%s fired %s at %s — HIT for %d damage.",
					attacker.Instance.Callsign, weaponRec.Name,
					target.Instance.Callsign, weaponRec.Damage))
			appendCombatEvent(eventsByLance, target.LanceInstanceID,
				fmt.Sprintf("%s hit by %s from %s — %d damage.",
					target.Instance.Callsign, weaponRec.Name,
					attacker.Instance.Callsign, weaponRec.Damage))
		} else {
			l.Info("%s: %s missed %s (roll %d >= %d%%)",
				attacker.Instance.Callsign, weaponRec.Name,
				target.Instance.Callsign, roll, chance)
			appendCombatEvent(eventsByLance, attacker.LanceInstanceID,
				fmt.Sprintf("%s fired %s at %s — missed.",
					attacker.Instance.Callsign, weaponRec.Name,
					target.Instance.Callsign))
		}
	}
	return totalDmg
}

func accumulateDamage(
	targetID string,
	totalDmg int,
	snapshots map[string]*mechSnapshot,
	damageMap map[string]*pendingDamage,
) {
	dm := damageMap[targetID]
	if dm == nil {
		dm = &pendingDamage{}
		damageMap[targetID] = dm
	}
	currentArmor := snapshots[targetID].Instance.CurrentArmor
	if totalDmg <= currentArmor {
		dm.armorDmg += totalDmg
	} else {
		dm.armorDmg += currentArmor
		dm.structureDmg += totalDmg - currentArmor
	}
}

func (p *Mecha) applyPendingDamage(
	l logger.Logger,
	damageMap map[string]*pendingDamage,
	snapshots map[string]*mechSnapshot,
	attacks []AttackDeclaration,
	eventsByLance map[string][]turnsheet.TurnEvent,
) {
	for mechID, dm := range damageMap {
		snap, ok := snapshots[mechID]
		if !ok {
			continue
		}
		inst := snap.Instance
		inst.CurrentArmor -= dm.armorDmg
		if inst.CurrentArmor < 0 {
			inst.CurrentArmor = 0
		}
		inst.CurrentStructure -= dm.structureDmg
		if inst.CurrentStructure < 0 {
			inst.CurrentStructure = 0
		}
		if inst.CurrentStructure <= 0 {
			inst.Status = mecha_record.MechInstanceStatusDestroyed
			l.Info("mech >%s< destroyed", inst.Callsign)
			appendCombatEvent(eventsByLance, snap.LanceInstanceID,
				fmt.Sprintf("%s has been DESTROYED!", inst.Callsign))
			for _, atk := range attacks {
				if atk.TargetMechInstanceID == mechID {
					if attSnap, ok := snapshots[atk.AttackerMechInstanceID]; ok {
						appendCombatEvent(eventsByLance, attSnap.LanceInstanceID,
							fmt.Sprintf("%s has been DESTROYED by your fire!", inst.Callsign))
					}
					break
				}
			}
		} else if inst.CurrentStructure < snap.Instance.CurrentStructure {
			inst.Status = mecha_record.MechInstanceStatusDamaged
		}
	}
}

func (p *Mecha) applyPendingHeat(
	l logger.Logger,
	heatMap map[string]int,
	snapshots map[string]*mechSnapshot,
	eventsByLance map[string][]turnsheet.TurnEvent,
) {
	for mechID, heat := range heatMap {
		snap, ok := snapshots[mechID]
		if !ok {
			continue
		}
		inst := snap.Instance
		inst.CurrentHeat += heat
		chassisRec, err := p.Domain.GetMechaChassisRec(inst.MechaChassisID, nil)
		if err != nil {
			l.Warn("failed to get chassis for heat check >%s<: %v", mechID, err)
			continue
		}
		if inst.CurrentHeat > chassisRec.HeatCapacity {
			if inst.Status != mecha_record.MechInstanceStatusDestroyed {
				inst.Status = mecha_record.MechInstanceStatusShutdown
				inst.CurrentHeat = chassisRec.HeatCapacity
				l.Info("mech >%s< overheated and shut down", inst.Callsign)
				appendCombatEvent(eventsByLance, snap.LanceInstanceID,
					fmt.Sprintf("%s has SHUT DOWN from overheating!", inst.Callsign))
			}
		}
	}
}

func (p *Mecha) persistMechChanges(
	l logger.Logger,
	snapshots map[string]*mechSnapshot,
	damageMap map[string]*pendingDamage,
	heatMap map[string]int,
) {
	for _, snap := range snapshots {
		if _, changed := damageMap[snap.Instance.ID]; changed {
			if _, err := p.Domain.UpdateMechaMechInstanceRec(snap.Instance); err != nil {
				l.Warn("failed to update mech instance >%s< after combat: %v", snap.Instance.ID, err)
			}
		} else if _, changed := heatMap[snap.Instance.ID]; changed {
			if _, err := p.Domain.UpdateMechaMechInstanceRec(snap.Instance); err != nil {
				l.Warn("failed to update mech instance >%s< heat after combat: %v", snap.Instance.ID, err)
			}
		}
	}
}

func appendCombatEvent(
	eventsByLance map[string][]turnsheet.TurnEvent,
	lanceInstanceID string,
	message string,
) {
	eventsByLance[lanceInstanceID] = append(
		eventsByLance[lanceInstanceID],
		turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryCombat,
			Icon:     turnsheet.TurnEventIconCombat,
			Message:  message,
		},
	)
}

func (p *Mecha) appendCombatEventsToLances(
	l logger.Logger,
	eventsByLance map[string][]turnsheet.TurnEvent,
) error {
	for lanceID, events := range eventsByLance {
		lanceInst, err := p.Domain.GetMechaLanceInstanceRec(lanceID, nil)
		if err != nil {
			l.Warn("failed to get lance instance >%s< for events: %v", lanceID, err)
			continue
		}
		for _, evt := range events {
			if err := turnsheet.AppendMechaTurnEvent(lanceInst, evt); err != nil {
				l.Warn("failed to append turn event for lance >%s<: %v", lanceID, err)
				continue
			}
		}
		if _, err := p.Domain.UpdateMechaLanceInstanceRec(lanceInst); err != nil {
			l.Warn("failed to persist turn events for lance >%s<: %v", lanceID, err)
		}
	}
	return nil
}
