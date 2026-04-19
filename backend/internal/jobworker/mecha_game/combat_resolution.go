package mecha_game

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/mecha_game/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// AttackDeclaration is an alias for the type defined in turn_sheet_processor.
type AttackDeclaration = turn_sheet_processor.AttackDeclaration

// mechSnapshot captures a mech's state before combat begins, so all attacks
// resolve against the same starting state (simultaneous resolution).
type mechSnapshot struct {
	Instance         *mecha_game_record.MechaGameMechInstance
	SquadInstanceID  string
	SectorInstanceID string
	Weapons          []mecha_game_record.WeaponConfigEntry
	Equipment        []mecha_game_record.EquipmentConfigEntry
	EquipmentByID    map[string]*mecha_game_record.MechaGameEquipment
	Effects          Effects
	// DidFireAmmoWeapon records whether this mech actually fired at least one
	// weapon with ammo_capacity > 0 this turn, so the equipment heat predicate
	// for ammo_bin can fire correctly at EoT accounting.
	DidFireAmmoWeapon bool
	// DidAttack records whether this mech actually attempted any attack this
	// turn, so the targeting-computer heat predicate can fire correctly.
	DidAttack bool
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
// sector distance. Range band rules:
//
//   - Short  (brawl):    distance 0 only
//   - Medium (versatile): distance 0–1
//   - Long   (standoff):  distance 1–2 only (cannot fire at same sector)
func weaponCanFire(rangeBand string, distance int) bool {
	switch rangeBand {
	case mecha_game_record.WeaponRangeBandShort:
		return distance == 0
	case mecha_game_record.WeaponRangeBandMedium:
		return distance == 0 || distance == 1
	case mecha_game_record.WeaponRangeBandLong:
		return distance == 1 || distance == 2
	default:
		return false
	}
}

// hitChance returns the probability of a single weapon hit (0–100).
// Formula: 50 + pilotSkill*5 + attackerHitBonus + coverModifier, capped to [0, 95].
//
// attackerHitBonus comes from the attacker's targeting-computer equipment
// (zero when absent or the attacker is refitting). coverModifier is the sum
// of the target sector's cover modifier and the target's ECM cover bonus
// (zero when absent or the target is refitting). Negative final values are
// clamped at 0; extreme positive values cap at 95 so there is always some
// miss chance.
func hitChance(pilotSkill int, attackerHitBonus int, coverModifier int) int {
	chance := 50 + pilotSkill*5 + attackerHitBonus + coverModifier
	if chance > 95 {
		chance = 95
	}
	if chance < 0 {
		chance = 0
	}
	return chance
}

// pendingDamage accumulates raw total damage from all attacks before applying (simultaneous
// resolution). Armor/structure split is performed once when all damage is applied so that
// focus-fire from multiple attackers cannot armor-absorb more total damage than the target has.
type pendingDamage struct {
	rawTotal int
}

// resolveCombat runs combat resolution for a game instance. It must be called
// after all movement (player and AI) has been applied. The sector graph from
// the decision engine context is rebuilt here for range calculations.
// Returns a map of mech instance ID → XP earned this turn (for pilot progression).
func (p *MechaGame) resolveCombat(
	_ context.Context,
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	attacks []AttackDeclaration,
) (map[string]int, error) {
	if len(attacks) == 0 {
		l.Info("no attack declarations for game instance >%s< — skipping combat", gameInstanceRec.ID)
		return nil, nil
	}

	allMechInsts, err := p.Domain.GetManyMechaGameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameMechInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		l.Warn("failed to load mech instances: %v", err)
		return nil, err
	}

	snapshots, err := p.buildMechSnapshots(l, allMechInsts)
	if err != nil {
		l.Warn("failed to build mech snapshots: %v", err)
		return nil, err
	}

	sectors, err := p.buildSectorGraph(l, gameInstanceRec.ID)
	if err != nil {
		l.Warn("failed to build sector graph: %v", err)
		return nil, err
	}

	seed := int64(gameInstanceRec.CurrentTurn)
	for _, b := range []byte(gameInstanceRec.ID) {
		seed += int64(b)
	}
	rng := rand.New(rand.NewSource(seed)) //nolint:gosec

	damageMap := make(map[string]*pendingDamage, len(allMechInsts))
	heatMap := make(map[string]int, len(allMechInsts))
	xpMap := make(map[string]int, len(allMechInsts))
	eventsBySquad := make(map[string][]turnsheet.TurnEvent)

	p.resolveAttacks(l, attacks, snapshots, sectors, rng, damageMap, heatMap, xpMap, eventsBySquad)
	p.applyPendingDamage(l, damageMap, snapshots, attacks, eventsBySquad)
	p.applyPendingHeat(l, heatMap, snapshots, eventsBySquad)
	p.persistMechChanges(l, snapshots, damageMap, heatMap)

	if err := p.appendCombatEventsToSquads(eventsBySquad); err != nil {
		l.Warn("failed to persist combat events: %v", err)
		return nil, err
	}

	return xpMap, nil
}

func (p *MechaGame) buildMechSnapshots(l logger.Logger, insts []*mecha_game_record.MechaGameMechInstance) (map[string]*mechSnapshot, error) {

	snapshots := make(map[string]*mechSnapshot, len(insts))

	for _, inst := range insts {
		var weapons []mecha_game_record.WeaponConfigEntry
		if len(inst.WeaponConfigJSON) > 0 {
			if err := json.Unmarshal(inst.WeaponConfigJSON, &weapons); err != nil {
				l.Warn("failed to unmarshal weapon config for mech >%s<: %v", inst.ID, err)
				return nil, err
			}
		}

		var equipment []mecha_game_record.EquipmentConfigEntry
		if len(inst.EquipmentConfigJSON) > 0 {
			if err := json.Unmarshal(inst.EquipmentConfigJSON, &equipment); err != nil {
				l.Warn("failed to unmarshal equipment config for mech >%s<: %v", inst.ID, err)
				return nil, err
			}
		}

		equipmentByID := make(map[string]*mecha_game_record.MechaGameEquipment, len(equipment))
		for _, entry := range equipment {
			if _, ok := equipmentByID[entry.EquipmentID]; ok {
				continue
			}
			eq, err := p.Domain.GetMechaGameEquipmentRec(entry.EquipmentID, nil)
			if err != nil {
				l.Warn("failed to load equipment >%s< for mech >%s<: %v", entry.EquipmentID, inst.ID, err)
				return nil, err
			}
			equipmentByID[entry.EquipmentID] = eq
		}

		snapshots[inst.ID] = &mechSnapshot{
			Instance:         inst,
			SquadInstanceID:  inst.MechaGameSquadInstanceID,
			SectorInstanceID: inst.MechaGameSectorInstanceID,
			Weapons:          weapons,
			Equipment:        equipment,
			EquipmentByID:    equipmentByID,
			Effects:          AggregateEffects(equipment, equipmentByID, inst.IsRefitting),
		}
	}

	return snapshots, nil
}

func (p *MechaGame) buildSectorGraph(l logger.Logger, gameInstanceID string) ([]*sectorState, error) {

	sectorInsts, err := p.Domain.GetManyMechaGameSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSectorInstanceGameInstanceID, Val: gameInstanceID},
		},
	})
	if err != nil {
		l.Warn("failed to load sector instances: %v", err)
		return nil, err
	}

	var sectors []*sectorState

	for _, si := range sectorInsts {

		links, err := p.Domain.GetManyMechaGameSectorLinkRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_game_record.FieldMechaGameSectorLinkFromMechaGameSectorID, Val: si.MechaGameSectorID},
			},
		})
		if err != nil {
			l.Warn("failed to load sector links: %v", err)
			return nil, err
		}

		var linkDestIDs []string
		for _, lnk := range links {
			for _, other := range sectorInsts {
				if other.MechaGameSectorID == lnk.ToMechaGameSectorID {
					linkDestIDs = append(linkDestIDs, other.ID)
					break
				}
			}
		}
			designRec, err := p.Domain.GetMechaGameSectorRec(si.MechaGameSectorID, nil)
		if err != nil {
			l.Warn("failed to load sector design >%s<: %v", si.MechaGameSectorID, err)
			return nil, err
		}

		sectors = append(sectors, &sectorState{
			Instance:            si,
			Design:              designRec,
			LinkDestInstanceIDs: linkDestIDs,
		})
	}

	return sectors, nil
}

const (
	xpPerCombatParticipation = 1
	xpPerKill                = 2
)

func (p *MechaGame) resolveAttacks(
	l logger.Logger,
	attacks []AttackDeclaration,
	snapshots map[string]*mechSnapshot,
	sectors []*sectorState,
	rng *rand.Rand,
	damageMap map[string]*pendingDamage,
	heatMap map[string]int,
	xpMap map[string]int,
	eventsBySquad map[string][]turnsheet.TurnEvent,
) {
	l = l.WithFunctionContext("MechaGame/resolveAttacks")

	participationAwarded := map[string]bool{}

	for _, atk := range attacks {

		l.Info("resolving attack by attacker mech >%s< on target mech >%s<", atk.AttackerMechInstanceID, atk.TargetMechInstanceID)

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

		if attacker.Instance.Status == mecha_game_record.MechInstanceStatusDestroyed ||
			attacker.Instance.Status == mecha_game_record.MechInstanceStatusShutdown {
			l.Info("attacker >%s< is %s — skipping attack", attacker.Instance.Callsign, attacker.Instance.Status)
			continue
		}

		if target.Instance.Status == mecha_game_record.MechInstanceStatusDestroyed {
			l.Info("target >%s< is already destroyed — skipping attack", target.Instance.Callsign)
			continue
		}

		dist := rangeDistance(attacker.SectorInstanceID, target.SectorInstanceID, sectors)
		// Distance 3+ is beyond any weapon's reach (long-range max = 2 hops).
		if dist > 2 {
			l.Info("%s fired at %s but target is out of range (distance %d)",
				attacker.Instance.Callsign, target.Instance.Callsign, dist)
			appendCombatEvent(eventsBySquad, attacker.SquadInstanceID,
				fmt.Sprintf("%s fired at %s — out of range.",
					attacker.Instance.Callsign, target.Instance.Callsign))
			continue
		}

		if len(attacker.Weapons) == 0 {
			l.Info("attacker >%s< has no weapons — skipping attack", attacker.Instance.Callsign)
			continue
		}

		// Award XP for combat participation once per attacker (regardless of hit/miss).
		if !participationAwarded[atk.AttackerMechInstanceID] {
			xpMap[atk.AttackerMechInstanceID] += xpPerCombatParticipation
			participationAwarded[atk.AttackerMechInstanceID] = true
		}

		attacker.DidAttack = true

		totalDmg := p.fireWeapons(l, atk, attacker, target, dist, sectors, rng, heatMap, eventsBySquad)
		if totalDmg > 0 {
			l.Info("total damage: %d", totalDmg)
			accumulateDamage(atk.TargetMechInstanceID, totalDmg, snapshots, damageMap)
		}
	}

	// Combat-predicated equipment heat cost accounting. Targeting-computer
	// heat fires when the attacker declared any attack; ammo-bin heat fires
	// when they actually pulled the trigger on a weapon with
	// ammo_capacity > 0. Always-on equipment heat (heat sink, armor
	// upgrade, ECM) is applied in end_of_turn so it runs even on turns
	// with no combat. Jump-jets heat is applied in the orders processor
	// when movement exceeds chassis base speed.
	for mechID, snap := range snapshots {
		inst := snap.Instance
		if inst.Status == mecha_game_record.MechInstanceStatusDestroyed ||
			inst.Status == mecha_game_record.MechInstanceStatusShutdown {
			continue
		}
		if len(snap.Equipment) == 0 {
			continue
		}
		heat := CombatHeatCost(
			snap.Equipment, snap.EquipmentByID, inst.IsRefitting,
			snap.DidAttack, snap.DidFireAmmoWeapon,
		)
		if heat > 0 {
			heatMap[mechID] += heat
		}
	}

	// Award XP for kills: first project which targets were destroyed, then
	// award kill XP once per unique attacker-target pair to prevent double-counting
	// if duplicate attack rows exist.
	killedTargets := map[string]bool{}
	for targetID, dm := range damageMap {
		snap := snapshots[targetID]
		if snap == nil {
			continue
		}
		projected := snap.Instance.CurrentStructure
		absorbed := dm.rawTotal
		if absorbed > snap.Instance.CurrentArmor {
			absorbed = snap.Instance.CurrentArmor
			projected -= (dm.rawTotal - absorbed)
		}
		if projected <= 0 {
			killedTargets[targetID] = true
		}
	}
	killAwarded := map[string]map[string]bool{}
	for _, atk := range attacks {
		if !killedTargets[atk.TargetMechInstanceID] {
			continue
		}
		if killAwarded[atk.AttackerMechInstanceID] == nil {
			killAwarded[atk.AttackerMechInstanceID] = map[string]bool{}
		}
		if !killAwarded[atk.AttackerMechInstanceID][atk.TargetMechInstanceID] {
			xpMap[atk.AttackerMechInstanceID] += xpPerKill
			killAwarded[atk.AttackerMechInstanceID][atk.TargetMechInstanceID] = true
		}
	}
}

func (p *MechaGame) fireWeapons(
	l logger.Logger,
	atk AttackDeclaration,
	attacker, target *mechSnapshot,
	dist int,
	sectors []*sectorState,
	rng *rand.Rand,
	heatMap map[string]int,
	eventsBySquad map[string][]turnsheet.TurnEvent,
) int {
	l = l.WithFunctionContext("MechaGame/fireWeapons")

	l.Info("firing weapons by attacker mech >%s< on target mech >%s<", attacker.Instance.Callsign, target.Instance.Callsign)

	// Determine the target sector's cover modifier.
	coverModifier := 0
	for _, sec := range sectors {
		if sec.Instance.ID == target.SectorInstanceID && sec.Design != nil {
			coverModifier = sec.Design.CoverModifier
			break
		}
	}

	// Layer on attacker targeting-computer bonus and defender ECM cover
	// bonus. Effects are already zeroed out for refitting mechs via
	// AggregateEffects, so no extra branching needed here.
	attackerHitBonus := attacker.Effects.HitChanceBonus
	effectiveCover := coverModifier + target.Effects.CoverBonus

	chance := hitChance(attacker.Instance.PilotSkill, attackerHitBonus, effectiveCover)

	totalDmg := 0
	for _, slot := range attacker.Weapons {
		if slot.WeaponID == "" {
			continue
		}

		weaponRec, err := p.Domain.GetMechaGameWeaponRec(slot.WeaponID, nil)
		if err != nil {
			l.Warn("failed to load weapon >%s<: %v", slot.WeaponID, err)
			continue
		}

		if !weaponCanFire(weaponRec.RangeBand, dist) {
			l.Debug("%s weapon %s cannot reach target at distance %d",
				attacker.Instance.Callsign, weaponRec.Name, dist)
			continue
		}

		// Ammo gate: weapons with ammo_capacity > 0 draw from the mech's
		// shared ammo pool. When the pool is empty, skip the weapon and emit
		// a visible event so the player knows why nothing happened. No heat,
		// no trigger pull.
		usesAmmo := weaponRec.AmmoCapacity > 0
		if usesAmmo && attacker.Instance.AmmoRemaining <= 0 {
			appendCombatEvent(eventsBySquad, attacker.SquadInstanceID,
				fmt.Sprintf("%s tried to fire %s at %s — OUT OF AMMO.",
					attacker.Instance.Callsign, weaponRec.Name, target.Instance.Callsign))
			continue
		}

		// Decrement ammo pool on trigger-pull regardless of hit/miss. Done
		// before the roll so a miss still consumes one round.
		if usesAmmo {
			attacker.Instance.AmmoRemaining--
			attacker.DidFireAmmoWeapon = true
		}

		heatMap[atk.AttackerMechInstanceID] += weaponRec.HeatCost
		roll := rng.Intn(100)
		if roll < chance {
			totalDmg += weaponRec.Damage
			l.Info("%s: %s hit %s with %s for %d damage (roll %d < %d%%)",
				attacker.Instance.Callsign, weaponRec.Name,
				target.Instance.Callsign, weaponRec.Name,
				weaponRec.Damage, roll, chance)
			appendCombatEvent(eventsBySquad, attacker.SquadInstanceID,
				fmt.Sprintf("%s fired %s at %s — HIT for %d damage.",
					attacker.Instance.Callsign, weaponRec.Name,
					target.Instance.Callsign, weaponRec.Damage))
			appendCombatEvent(eventsBySquad, target.SquadInstanceID,
				fmt.Sprintf("%s hit by %s from %s — %d damage.",
					target.Instance.Callsign, weaponRec.Name,
					attacker.Instance.Callsign, weaponRec.Damage))
		} else {
			l.Info("%s: %s missed %s (roll %d >= %d%%)",
				attacker.Instance.Callsign, weaponRec.Name,
				target.Instance.Callsign, roll, chance)
			appendCombatEvent(eventsBySquad, attacker.SquadInstanceID,
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
	_ map[string]*mechSnapshot,
	damageMap map[string]*pendingDamage,
) {
	dm := damageMap[targetID]
	if dm == nil {
		dm = &pendingDamage{}
		damageMap[targetID] = dm
	}
	dm.rawTotal += totalDmg
}

func (p *MechaGame) applyPendingDamage(
	l logger.Logger,
	damageMap map[string]*pendingDamage,
	snapshots map[string]*mechSnapshot,
	attacks []AttackDeclaration,
	eventsBySquad map[string][]turnsheet.TurnEvent,
) {
	for mechID, dm := range damageMap {
		snap, ok := snapshots[mechID]
		if !ok {
			continue
		}
		inst := snap.Instance

		// Capture pre-damage structure to determine damaged status afterwards.
		origStructure := inst.CurrentStructure

		// Split raw total damage into armor and structure components once (simultaneous
		// resolution: all attacks pool their damage before armor absorbs any of it).
		armorAbsorbed := dm.rawTotal
		structureDmg := 0
		if armorAbsorbed > inst.CurrentArmor {
			armorAbsorbed = inst.CurrentArmor
			structureDmg = dm.rawTotal - armorAbsorbed
		}

		inst.CurrentArmor -= armorAbsorbed
		if inst.CurrentArmor < 0 {
			inst.CurrentArmor = 0
		}
		inst.CurrentStructure -= structureDmg
		if inst.CurrentStructure < 0 {
			inst.CurrentStructure = 0
		}

		if inst.CurrentStructure <= 0 {
			inst.Status = mecha_game_record.MechInstanceStatusDestroyed
			l.Info("mech >%s< destroyed", inst.Callsign)
			appendCombatEvent(eventsBySquad, snap.SquadInstanceID,
				fmt.Sprintf("%s has been DESTROYED!", inst.Callsign))
			for _, atk := range attacks {
				if atk.TargetMechInstanceID == mechID {
					if attSnap, ok := snapshots[atk.AttackerMechInstanceID]; ok {
						appendCombatEvent(eventsBySquad, attSnap.SquadInstanceID,
							fmt.Sprintf("%s has been DESTROYED by your fire!", inst.Callsign))
					}
					break
				}
			}
		} else if inst.CurrentStructure < origStructure {
			inst.Status = mecha_game_record.MechInstanceStatusDamaged
		}
	}
}

func (p *MechaGame) applyPendingHeat(
	l logger.Logger,
	heatMap map[string]int,
	snapshots map[string]*mechSnapshot,
	eventsBySquad map[string][]turnsheet.TurnEvent,
) {
	for mechID, heat := range heatMap {
		snap, ok := snapshots[mechID]
		if !ok {
			continue
		}
		inst := snap.Instance
		inst.CurrentHeat += heat
		chassisRec, err := p.Domain.GetMechaGameChassisRec(inst.MechaGameChassisID, nil)
		if err != nil {
			l.Warn("failed to get chassis for heat check >%s<: %v", mechID, err)
			continue
		}
		if inst.CurrentHeat > chassisRec.HeatCapacity {
			if inst.Status != mecha_game_record.MechInstanceStatusDestroyed {
				inst.Status = mecha_game_record.MechInstanceStatusShutdown
				inst.CurrentHeat = chassisRec.HeatCapacity
				l.Info("mech >%s< overheated and shut down", inst.Callsign)
				appendCombatEvent(eventsBySquad, snap.SquadInstanceID,
					fmt.Sprintf("%s has SHUT DOWN from overheating!", inst.Callsign))
			}
		}
	}
}

func (p *MechaGame) persistMechChanges(
	l logger.Logger,
	snapshots map[string]*mechSnapshot,
	damageMap map[string]*pendingDamage,
	heatMap map[string]int,
) {
	for _, snap := range snapshots {
		// Ammo consumption is also a reason to persist, independent of damage
		// and heat, so a mech that fired but dealt no damage still has its
		// pool decrement saved.
		ammoChanged := snap.DidFireAmmoWeapon
		if _, changed := damageMap[snap.Instance.ID]; changed {
			if _, err := p.Domain.UpdateMechaGameMechInstanceRec(snap.Instance); err != nil {
				l.Warn("failed to update mech instance >%s< after combat: %v", snap.Instance.ID, err)
			}
		} else if _, changed := heatMap[snap.Instance.ID]; changed {
			if _, err := p.Domain.UpdateMechaGameMechInstanceRec(snap.Instance); err != nil {
				l.Warn("failed to update mech instance >%s< heat after combat: %v", snap.Instance.ID, err)
			}
		} else if ammoChanged {
			if _, err := p.Domain.UpdateMechaGameMechInstanceRec(snap.Instance); err != nil {
				l.Warn("failed to update mech instance >%s< ammo after combat: %v", snap.Instance.ID, err)
			}
		}
	}
}

func appendCombatEvent(
	eventsBySquad map[string][]turnsheet.TurnEvent,
	squadInstanceID string,
	message string,
) {
	eventsBySquad[squadInstanceID] = append(
		eventsBySquad[squadInstanceID],
		turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryCombat,
			Icon:     turnsheet.TurnEventIconCombat,
			Message:  message,
		},
	)
}

func (p *MechaGame) appendCombatEventsToSquads(
	eventsBySquad map[string][]turnsheet.TurnEvent,
) error {
	for squadID, events := range eventsBySquad {
		squadInst, err := p.Domain.GetMechaGameSquadInstanceRec(squadID, nil)
		if err != nil {
			return fmt.Errorf("failed to get squad instance >%s< for events: %w", squadID, err)
		}
		for _, evt := range events {
			if err := turnsheet.AppendMechaGameTurnEvent(squadInst, evt); err != nil {
				return fmt.Errorf("failed to append turn event for squad >%s<: %w", squadID, err)
			}
		}
		if _, err := p.Domain.UpdateMechaGameSquadInstanceRec(squadInst); err != nil {
			return fmt.Errorf("failed to persist turn events for squad >%s<: %w", squadID, err)
		}
	}
	return nil
}
