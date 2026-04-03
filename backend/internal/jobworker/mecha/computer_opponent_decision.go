package mecha

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/agent"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// GameStateContext bundles all information both decision strategies need for a
// single computer opponent's turn.
type GameStateContext struct {
	Opponent      *mecha_record.MechaComputerOpponent
	LanceInstance *mecha_record.MechaLanceInstance
	OwnMechs      []*mecha_record.MechaMechInstance
	EnemyMechs    []*mechState
	Sectors       []*sectorState
	// ChassisCache maps chassis ID to chassis design record for speed lookups.
	ChassisCache map[string]*mecha_record.MechaChassis
	TurnNumber   int
}

type mechState struct {
	Instance *mecha_record.MechaMechInstance
	Sector   *mecha_record.MechaSectorInstance
}

type sectorState struct {
	Instance            *mecha_record.MechaSectorInstance
	Design              *mecha_record.MechaSector
	LinkDestInstanceIDs []string
}

// ComputerOpponentStrategy is the interface both strategies implement.
type ComputerOpponentStrategy interface {
	GenerateOrders(ctx context.Context, l logger.Logger, state *GameStateContext) ([]turnsheet.ScannedMechOrder, error)
}

// ComputerOpponentDecisionEngine selects a strategy and generates orders for
// each computer opponent lance during turn processing.
type ComputerOpponentDecisionEngine struct {
	logger           logger.Logger
	domain           *domain.Domain
	primaryStrategy  ComputerOpponentStrategy
	fallbackStrategy ComputerOpponentStrategy
}

// NewComputerOpponentDecisionEngine creates the decision engine. If an OpenAI API
// key is configured the engine uses LLM-based orders as the primary strategy with
// a rule-based fallback; otherwise only the rule-based strategy is used.
func NewComputerOpponentDecisionEngine(l logger.Logger, d *domain.Domain, cfg config.Config) *ComputerOpponentDecisionEngine {
	l = l.WithFunctionContext("NewComputerOpponentDecisionEngine")

	rbStrategy := &ruleBasedStrategy{}

	var primary ComputerOpponentStrategy = rbStrategy
	var fallback ComputerOpponentStrategy = rbStrategy

	if cfg.OpenAIAPIKey != "" {
		l.Info("OpenAI API key configured — using LLM strategy with rule-based fallback")
		textAgent := agent.NewOpenAITextAgent(l, cfg)
		primary = &llmStrategy{textAgent: textAgent}
		fallback = rbStrategy
	} else {
		l.Info("no OpenAI API key — using rule-based strategy only")
	}

	return &ComputerOpponentDecisionEngine{
		logger:           l,
		domain:           d,
		primaryStrategy:  primary,
		fallbackStrategy: fallback,
	}
}

// GenerateOrdersForLance builds the GameStateContext for the given lance instance
// and generates orders using the configured strategy.
func (e *ComputerOpponentDecisionEngine) GenerateOrdersForLance(
	ctx context.Context,
	gameInstanceID string,
	lanceInstance *mecha_record.MechaLanceInstance,
	opponentRec *mecha_record.MechaComputerOpponent,
	turnNumber int,
) ([]turnsheet.ScannedMechOrder, error) {
	l := e.logger.WithFunctionContext("ComputerOpponentDecisionEngine/GenerateOrdersForLance")

	state, err := e.buildGameStateContext(ctx, gameInstanceID, lanceInstance, opponentRec, turnNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to build game state context: %w", err)
	}

	orders, err := e.primaryStrategy.GenerateOrders(ctx, l, state)
	if err != nil {
		l.Warn("primary strategy failed, falling back to rule-based: %v", err)
		orders, err = e.fallbackStrategy.GenerateOrders(ctx, l, state)
		if err != nil {
			return nil, fmt.Errorf("fallback strategy also failed: %w", err)
		}
	}

	return orders, nil
}

// buildGameStateContext queries the domain for all data needed by both strategies.
func (e *ComputerOpponentDecisionEngine) buildGameStateContext(
	_ context.Context,
	gameInstanceID string,
	lanceInstance *mecha_record.MechaLanceInstance,
	opponentRec *mecha_record.MechaComputerOpponent,
	turnNumber int,
) (*GameStateContext, error) {
	// Own mechs
	ownMechs, err := e.domain.GetManyMechaMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaMechInstanceMechaLanceInstanceID, Val: lanceInstance.ID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get own mech instances: %w", err)
	}

	// All mech instances in this game instance (for enemy detection)
	allMechs, err := e.domain.GetManyMechaMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaMechInstanceGameInstanceID, Val: gameInstanceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get all mech instances: %w", err)
	}

	ownMechIDs := make(map[string]bool, len(ownMechs))
	for _, m := range ownMechs {
		ownMechIDs[m.ID] = true
	}

	var enemyMechs []*mechState
	for _, m := range allMechs {
		if ownMechIDs[m.ID] {
			continue
		}
		if m.Status == mecha_record.MechInstanceStatusDestroyed {
			continue
		}
		sectorInst, err := e.domain.GetMechaSectorInstanceRec(m.MechaSectorInstanceID, nil)
		if err != nil {
			continue
		}
		enemyMechs = append(enemyMechs, &mechState{Instance: m, Sector: sectorInst})
	}

	// Sector graph
	sectorInstances, err := e.domain.GetManyMechaSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSectorInstanceGameInstanceID, Val: gameInstanceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get sector instances: %w", err)
	}

	var sectors []*sectorState
	for _, si := range sectorInstances {
		design, err := e.domain.GetMechaSectorRec(si.MechaSectorID, nil)
		if err != nil {
			continue
		}

		links, err := e.domain.GetManyMechaSectorLinkRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaSectorLinkFromMechaSectorID, Val: si.MechaSectorID},
			},
		})
		if err != nil {
			links = nil
		}

		var linkDestIDs []string
		for _, lnk := range links {
			for _, destSI := range sectorInstances {
				if destSI.MechaSectorID == lnk.ToMechaSectorID {
					linkDestIDs = append(linkDestIDs, destSI.ID)
					break
				}
			}
		}

		sectors = append(sectors, &sectorState{
			Instance:            si,
			Design:              design,
			LinkDestInstanceIDs: linkDestIDs,
		})
	}

	// Build chassis cache for speed lookups.
	chassisCache := make(map[string]*mecha_record.MechaChassis)
	allMechsForChassis := append(ownMechs, func() []*mecha_record.MechaMechInstance {
		var ms []*mecha_record.MechaMechInstance
		for _, em := range enemyMechs {
			ms = append(ms, em.Instance)
		}
		return ms
	}()...)
	for _, m := range allMechsForChassis {
		if m.MechaChassisID == "" || chassisCache[m.MechaChassisID] != nil {
			continue
		}
		if cr, err := e.domain.GetMechaChassisRec(m.MechaChassisID, nil); err == nil {
			chassisCache[m.MechaChassisID] = cr
		}
	}

	return &GameStateContext{
		Opponent:      opponentRec,
		LanceInstance: lanceInstance,
		OwnMechs:      ownMechs,
		EnemyMechs:    enemyMechs,
		Sectors:       sectors,
		ChassisCache:  chassisCache,
		TurnNumber:    turnNumber,
	}, nil
}

// ─── LLM Strategy ────────────────────────────────────────────────────────────

type llmStrategy struct {
	textAgent agent.TextAgent
}

func (s *llmStrategy) GenerateOrders(ctx context.Context, l logger.Logger, state *GameStateContext) ([]turnsheet.ScannedMechOrder, error) {
	opp := state.Opponent

	// Temperature: low IQ → higher temperature (more random), high IQ → lower (more precise).
	temperature := 1.0 - float64(opp.IQ-1)/9.0*0.8 // Range: 0.2 (IQ=10) to 1.0 (IQ=1)

	systemPrompt := fmt.Sprintf(
		"You are a mecha combat commander. You control the %s force. "+
			"Your aggression level is %d/10 (1=purely defensive, 10=all-out attack). "+
			"Your tactical IQ is %d/10 (1=predictable random moves, 10=expert flanking and terrain use). "+
			"Generate movement orders for your mechs based on current positions and enemy positions.",
		opp.Name, opp.Aggression, opp.IQ,
	)

	userPrompt := s.buildGameStatePrompt(state)

	resp, err := s.textAgent.GenerateContent(ctx, agent.ContentGenerationRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  temperature,
		MaxTokens:    512,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM text generation failed: %w", err)
	}

	// Extract JSON from the response (the model may wrap it in markdown)
	jsonStr := extractJSONFromResponse(resp)
	if jsonStr == "" {
		return nil, fmt.Errorf("LLM response contained no JSON: %s", resp)
	}

	var scanData turnsheet.OrdersScanData
	if err := json.Unmarshal([]byte(jsonStr), &scanData); err != nil {
		return nil, fmt.Errorf("failed to parse LLM orders JSON: %w", err)
	}

	l.Info("LLM strategy generated %d mech orders for opponent %s", len(scanData.MechOrders), opp.Name)
	return scanData.MechOrders, nil
}

func (s *llmStrategy) buildGameStatePrompt(state *GameStateContext) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Turn %d. You command the following mechs:\n", state.TurnNumber)
	for _, m := range state.OwnMechs {
		sectorName := "unknown"
		for _, sec := range state.Sectors {
			if sec.Instance.ID == m.MechaSectorInstanceID {
				sectorName = sec.Design.Name
				break
			}
		}
		speed := 1
		if state.ChassisCache != nil {
			if cr, ok := state.ChassisCache[m.MechaChassisID]; ok {
				speed = cr.Speed
			}
		}
		fmt.Fprintf(&sb, "  - Mech ID: %s, Callsign: %s, Status: %s, Location: %s (sector_instance_id: %s), Armor: %d, Structure: %d, Speed: %d (max sector hops per turn)\n",
			m.ID, m.Callsign, m.Status, sectorName, m.MechaSectorInstanceID,
			m.CurrentArmor, m.CurrentStructure, speed)
	}

	if len(state.EnemyMechs) > 0 {
		sb.WriteString("\nEnemy mechs (targets):\n")
		for _, em := range state.EnemyMechs {
			sectorName := "unknown"
			if em.Sector != nil {
				for _, sec := range state.Sectors {
					if sec.Instance.ID == em.Sector.ID {
						sectorName = sec.Design.Name
						break
					}
				}
			}
			fmt.Fprintf(&sb, "  - Mech ID: %s, Callsign: %s, Location: %s (sector_instance_id: %s)\n",
				em.Instance.ID, em.Instance.Callsign, sectorName, em.Instance.MechaSectorInstanceID)
		}
	}

	sb.WriteString("\nAvailable sectors (movement graph):\n")
	for _, sec := range state.Sectors {
		var linkNames []string
		for _, destID := range sec.LinkDestInstanceIDs {
			for _, other := range state.Sectors {
				if other.Instance.ID == destID {
					linkNames = append(linkNames, fmt.Sprintf("%s (id: %s)", other.Design.Name, other.Instance.ID))
					break
				}
			}
		}
		fmt.Fprintf(&sb, "  - %s (id: %s, terrain: %s, elevation: %d) → [%s]\n",
			sec.Design.Name, sec.Instance.ID, sec.Design.TerrainType, sec.Design.Elevation,
			strings.Join(linkNames, ", "))
	}

	sb.WriteString(`
Return your orders as JSON with this exact structure:
{
  "mech_orders": [
    {
      "mech_instance_id": "<exact mech ID from above>",
      "move_to_sector_instance_id": "<exact sector_instance_id to move to, or empty string to stay>",
      "attack_target_mech_instance_id": "<exact enemy mech ID to attack, or empty string for no attack>"
    }
  ]
}
Include one entry per mech. Use exact IDs from above.
Each mech can move up to its Speed number of connected sector hops per turn (follow the movement graph).
Attack targets must be within weapon range after movement (short-range: same sector, medium-range: same or adjacent, long-range: adjacent or 2 sectors away).
Respond with JSON only — no markdown, no commentary.`)

	return sb.String()
}

// extractJSONFromResponse strips markdown code fences if present.
func extractJSONFromResponse(resp string) string {
	resp = strings.TrimSpace(resp)

	// Strip markdown code fences
	if idx := strings.Index(resp, "```json"); idx != -1 {
		resp = resp[idx+7:]
		if end := strings.Index(resp, "```"); end != -1 {
			resp = resp[:end]
		}
	} else if idx := strings.Index(resp, "```"); idx != -1 {
		resp = resp[idx+3:]
		if end := strings.Index(resp, "```"); end != -1 {
			resp = resp[:end]
		}
	}

	resp = strings.TrimSpace(resp)

	// Find the JSON object bounds
	start := strings.Index(resp, "{")
	end := strings.LastIndex(resp, "}")
	if start == -1 || end == -1 || end <= start {
		return ""
	}

	return resp[start : end+1]
}

// ─── Rule-Based Strategy ─────────────────────────────────────────────────────

type ruleBasedStrategy struct{}

func (s *ruleBasedStrategy) GenerateOrders(ctx context.Context, l logger.Logger, state *GameStateContext) ([]turnsheet.ScannedMechOrder, error) {
	opp := state.Opponent
	orders := make([]turnsheet.ScannedMechOrder, 0, len(state.OwnMechs))

	for _, mech := range state.OwnMechs {
		if mech.Status == mecha_record.MechInstanceStatusDestroyed ||
			mech.Status == mecha_record.MechInstanceStatusShutdown {
			orders = append(orders, turnsheet.ScannedMechOrder{
				MechInstanceID: mech.ID,
			})
			continue
		}

		// Look up chassis speed for multi-hop movement.
		mechSpeed := 1
		if mech.MechaChassisID != "" {
			if chassisRec, err := s.lookupChassis(state, mech.MechaChassisID); chassisRec != nil && err == nil {
				mechSpeed = chassisRec.Speed
			}
		}

		targetSectorInstanceID := s.pickMovementTarget(l, opp, mech, mechSpeed, state)

		// After movement (use new sector if moving), pick an attack target
		postMoveSectorID := targetSectorInstanceID
		if postMoveSectorID == "" {
			postMoveSectorID = mech.MechaSectorInstanceID
		}
		attackTargetID := s.pickAttackTarget(opp, postMoveSectorID, state)

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
func (s *ruleBasedStrategy) pickAttackTarget(opp *mecha_record.MechaComputerOpponent, fromSectorID string, state *GameStateContext) string {
	const maxEngagementRange = 2

	var candidates []*mechState
	for _, em := range state.EnemyMechs {
		if em.Instance.Status == mecha_record.MechInstanceStatusDestroyed {
			continue
		}
		dist := s.sectorDistance(fromSectorID, em.Instance.MechaSectorInstanceID, state)
		if dist <= maxEngagementRange {
			candidates = append(candidates, em)
		}
	}

	if len(candidates) == 0 {
		return ""
	}

	// High aggression: target weakest (lowest structure) to maximise kills
	// Low aggression: target strongest (highest structure) as deterrent
	best := candidates[0]
	for _, em := range candidates[1:] {
		if opp.Aggression >= 6 {
			if em.Instance.CurrentStructure < best.Instance.CurrentStructure {
				best = em
			}
		} else {
			if em.Instance.CurrentStructure > best.Instance.CurrentStructure {
				best = em
			}
		}
	}

	return best.Instance.ID
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
func (s *ruleBasedStrategy) pickMovementTarget(_ logger.Logger, opp *mecha_record.MechaComputerOpponent, mech *mecha_record.MechaMechInstance, mechSpeed int, state *GameStateContext) string {
	reachableIDs := s.getReachableSectorIDs(mech.MechaSectorInstanceID, mechSpeed, state)
	if len(reachableIDs) == 0 {
		return ""
	}

	if opp.Aggression >= 7 && len(state.EnemyMechs) > 0 {
		// Advance toward nearest enemy
		nearestEnemySectorID := s.findNearestEnemy(mech.MechaSectorInstanceID, state)
		if nearestEnemySectorID != "" {
			best := s.pickBestAdvanceStep(mech.MechaSectorInstanceID, nearestEnemySectorID, reachableIDs, opp.IQ, state)
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
func (s *ruleBasedStrategy) lookupChassis(state *GameStateContext, chassisID string) (*mecha_record.MechaChassis, error) {
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
