package mecha_game

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

type ruleBasedStrategy struct{}

func (s *ruleBasedStrategy) GenerateOrders(ctx context.Context, l logger.Logger, state *GameStateContext) ([]turnsheet.ScannedMechOrder, error) {
	l = l.WithFunctionContext("ruleBasedStrategy/GenerateOrders")

	opp := state.Opponent
	orders := make([]turnsheet.ScannedMechOrder, 0, len(state.OwnMechs))

	for _, mech := range state.OwnMechs {
		if mech.Status == mecha_game_record.MechInstanceStatusDestroyed ||
			mech.Status == mecha_game_record.MechInstanceStatusShutdown {
			orders = append(orders, turnsheet.ScannedMechOrder{
				MechInstanceID: mech.ID,
			})
			continue
		}

		// Look up chassis speed for multi-hop movement. Layer on any
		// jump-jets bonus (zeroed while refitting) so the AI budgets match
		// what the orders processor will actually accept.
		mechSpeed := 1
		if mech.MechaGameChassisID != "" {
			if chassisRec, err := s.lookupChassis(state, mech.MechaGameChassisID); chassisRec != nil && err == nil {
				effects := state.EffectsByMechID[mech.ID]
				mechSpeed = chassisRec.Speed + effects.SpeedBonus
			}
		}

		targetSectorInstanceID := s.pickMovementTarget(l, opp, mech, mechSpeed, state)

		// After movement (use new sector if moving), pick an attack target
		postMoveSectorID := targetSectorInstanceID
		if postMoveSectorID == "" {
			postMoveSectorID = mech.MechaGameSectorInstanceID
		}
		attackTargetID := s.pickAttackTarget(opp, postMoveSectorID, mech, state)
		orders = append(orders, turnsheet.ScannedMechOrder{
			MechInstanceID:             mech.ID,
			MoveToSectorInstanceID:     targetSectorInstanceID,
			AttackTargetMechInstanceID: attackTargetID,
		})
	}

	l.Info("rule-based strategy generated %d mech orders for opponent %s", len(orders), opp.Name)

	return orders, nil
}

// pickAttackTarget selects an enemy mech to attack from the given sector.
// Returns empty string if no valid target exists.
// Uses max range 2 (long-range weapon reach) as the engagement envelope.
// Combat resolution will only fire weapons that can actually reach the target.
//
// The attacker's targeting-computer HitChanceBonus and the defender's ECM
// CoverBonus are used purely as tie-breakers between otherwise equally
// attractive targets (same CurrentStructure bucket). This keeps the AI
// honest about the same modifiers the combat resolver applies, without
// overhauling the aggression-driven primary ranking.
func (s *ruleBasedStrategy) pickAttackTarget(
	opp *mecha_game_record.MechaGameComputerOpponent,
	fromSectorID string,
	attacker *mecha_game_record.MechaGameMechInstance,
	state *GameStateContext,
) string {
	const maxEngagementRange = 2

	var candidates []*mechState
	for _, em := range state.EnemyMechs {
		if em.Instance.Status == mecha_game_record.MechInstanceStatusDestroyed {
			continue
		}
		dist := s.sectorDistance(fromSectorID, em.Instance.MechaGameSectorInstanceID, state)
		if dist <= maxEngagementRange {
			candidates = append(candidates, em)
		}
	}

	if len(candidates) == 0 {
		return ""
	}

	attackerHitBonus := 0
	if attacker != nil {
		attackerHitBonus = state.EffectsByMechID[attacker.ID].HitChanceBonus
	}

	// Primary ranking: aggression-driven structure bucket
	// (high aggression → weakest; low aggression → strongest).
	// Tie-break: prefer targets with the smaller *effective* cover
	// (chassis-level cover is not tracked per mech; we use ECM CoverBonus
	// as the per-mech cover signal, adjusted by the attacker's hit bonus).
	best := candidates[0]
	bestScore := targetScore(opp, best, attackerHitBonus, state)
	for _, em := range candidates[1:] {
		score := targetScore(opp, em, attackerHitBonus, state)
		if score > bestScore {
			best = em
			bestScore = score
		}
	}

	return best.Instance.ID
}

// targetScore produces a single comparable score that orders candidates by
// the primary aggression criterion (structure) and breaks ties using the
// defender's ECM cover bonus net of the attacker's targeting-computer hit
// bonus. Higher score == more attractive target.
func targetScore(
	opp *mecha_game_record.MechaGameComputerOpponent,
	em *mechState,
	attackerHitBonus int,
	state *GameStateContext,
) int {
	defenderCover := state.EffectsByMechID[em.Instance.ID].CoverBonus
	// netCover is how hard the defender is to hit after accounting for
	// the attacker's targeting advantage. Larger positive values mean a
	// harder-to-hit target (less attractive); larger negative values mean
	// an easier target (more attractive). We flip the sign so the score
	// increases as the target becomes easier.
	netCover := attackerHitBonus - defenderCover
	structure := em.Instance.CurrentStructure
	// High aggression: prefer low structure (bigger primary term = more
	// attractive when structure is low). Encode as -structure * 1000 so
	// the primary signal dominates the tie-breaker.
	// Low aggression: prefer high structure (opposite sign).
	if opp.Aggression >= 6 {
		return -structure*1000 + netCover
	}
	return structure*1000 + netCover
}

// sectorDistance returns the BFS distance between two sector instance IDs
// using the sector graph, capped at 3 (out of range).
func (s *ruleBasedStrategy) sectorDistance(fromID, toID string, state *GameStateContext) int {
	if fromID == toID {
		return 0
	}
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
		for _, sec := range state.Sectors {
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

// pickMovementTarget returns the sector instance ID the mech should move to, or
// empty string to stay in place. Considers the mech's speed for multi-hop movement.
func (s *ruleBasedStrategy) pickMovementTarget(_ logger.Logger, opp *mecha_game_record.MechaGameComputerOpponent, mech *mecha_game_record.MechaGameMechInstance, mechSpeed int, state *GameStateContext) string {
	reachableIDs := s.getReachableSectorIDs(mech.MechaGameSectorInstanceID, mechSpeed, state)
	if len(reachableIDs) == 0 {
		return ""
	}

	if opp.Aggression >= 7 && len(state.EnemyMechs) > 0 {
		// Advance toward nearest enemy
		nearestEnemySectorID := s.findNearestEnemy(mech.MechaGameSectorInstanceID, state)
		if nearestEnemySectorID != "" {
			best := s.pickBestAdvanceStep(mech.MechaGameSectorInstanceID, nearestEnemySectorID, reachableIDs, opp.IQ, state)
			if best != "" {
				return best
			}
		}
	} else if opp.Aggression <= 3 {
		// Retreat to starting / best-covered sector
		best := s.pickBestDefensiveStep(reachableIDs, opp.IQ, state)
		if best != "" {
			return best
		}
	}

	// Hold — no movement
	return ""
}

// getReachableSectorIDs returns all sector instance IDs reachable within the given
// number of hops from the given sector, using BFS over the sector graph.
func (s *ruleBasedStrategy) getReachableSectorIDs(fromSectorID string, speed int, state *GameStateContext) []string {
	if speed <= 0 {
		return nil
	}

	type bfsNode struct {
		id    string
		depth int
	}

	seen := map[string]bool{fromSectorID: true}
	queue := []bfsNode{{id: fromSectorID, depth: 0}}
	var reachable []string

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.depth >= speed {
			continue
		}

		for _, sec := range state.Sectors {
			if sec.Instance.ID != cur.id {
				continue
			}
			for _, destID := range sec.LinkDestInstanceIDs {
				if seen[destID] {
					continue
				}
				seen[destID] = true
				reachable = append(reachable, destID)
				queue = append(queue, bfsNode{id: destID, depth: cur.depth + 1})
			}
		}
	}

	return reachable
}

// findNearestEnemy returns the sector instance ID of the nearest enemy mech using
// actual BFS distance through the sector graph. Returns empty string if no enemies.
func (s *ruleBasedStrategy) findNearestEnemy(fromSectorID string, state *GameStateContext) string {
	enemySectorIDs := make(map[string]bool)
	for _, em := range state.EnemyMechs {
		if em.Sector != nil {
			enemySectorIDs[em.Sector.ID] = true
		}
	}
	if len(enemySectorIDs) == 0 {
		return ""
	}

	type bfsNode struct {
		id    string
		depth int
	}

	seen := map[string]bool{fromSectorID: true}
	queue := []bfsNode{{id: fromSectorID, depth: 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if enemySectorIDs[cur.id] {
			return cur.id
		}

		if cur.depth >= 20 {
			continue
		}

		for _, sec := range state.Sectors {
			if sec.Instance.ID != cur.id {
				continue
			}
			for _, destID := range sec.LinkDestInstanceIDs {
				if !seen[destID] {
					seen[destID] = true
					queue = append(queue, bfsNode{id: destID, depth: cur.depth + 1})
				}
			}
		}
	}

	return ""
}

func (s *ruleBasedStrategy) pickBestAdvanceStep(_ string, targetID string, reachableIDs []string, iq int, state *GameStateContext) string {
	// Among reachable sectors, pick the one with the shortest distance to the target.
	bestID := ""
	bestDist := 999

	for _, candID := range reachableIDs {
		dist := s.sectorDistance(candID, targetID, state)
		if dist < bestDist {
			bestDist = dist
			bestID = candID
		}
	}

	if bestID == "" {
		return ""
	}

	// When IQ is high, among the equally best sectors prefer higher cover.
	if iq >= 5 {
		var equally []string
		for _, candID := range reachableIDs {
			if s.sectorDistance(candID, targetID, state) == bestDist {
				equally = append(equally, candID)
			}
		}
		return s.pickHighCoverSector(equally, state)
	}

	return bestID
}

func (s *ruleBasedStrategy) pickBestDefensiveStep(reachableIDs []string, iq int, state *GameStateContext) string {
	if iq >= 5 {
		return s.pickHighCoverSector(reachableIDs, state)
	}
	return ""
}

// lookupChassis retrieves a chassis record from the context's cache, or nil if not found.
func (s *ruleBasedStrategy) lookupChassis(state *GameStateContext, chassisID string) (*mecha_game_record.MechaGameChassis, error) {
	if state.ChassisCache != nil {
		if cr, ok := state.ChassisCache[chassisID]; ok {
			return cr, nil
		}
	}
	return nil, fmt.Errorf("chassis >%s< not found in cache", chassisID)
}

// pickHighCoverSector returns the sector with the highest cover_modifier from
// the given candidates. Elevation is used as a tiebreaker when cover is equal.
func (s *ruleBasedStrategy) pickHighCoverSector(candidates []string, state *GameStateContext) string {
	best := ""
	bestCover := -999
	bestElev := -999
	for _, id := range candidates {
		for _, sec := range state.Sectors {
			if sec.Instance.ID != id {
				continue
			}
			cover := sec.Design.CoverModifier
			elev := sec.Design.Elevation
			if cover > bestCover || (cover == bestCover && elev > bestElev) {
				bestCover = cover
				bestElev = elev
				best = id
			}
			break
		}
	}
	return best
}
