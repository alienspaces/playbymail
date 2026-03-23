package mecha

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/agent"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
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
	TurnNumber    int
}

type mechState struct {
	Instance *mecha_record.MechaMechInstance
	Sector   *mecha_record.MechaSectorInstance
}

type sectorState struct {
	Instance  *mecha_record.MechaSectorInstance
	Design    *mecha_record.MechaSector
	LinkDestInstanceIDs []string
}

// ComputerOpponentStrategy is the interface both strategies implement.
type ComputerOpponentStrategy interface {
	GenerateOrders(ctx context.Context, l logger.Logger, state *GameStateContext) ([]turnsheet.ScannedMechOrder, error)
}

// ComputerOpponentDecisionEngine selects a strategy and generates orders for
// each computer opponent lance during turn processing.
type ComputerOpponentDecisionEngine struct {
	logger          logger.Logger
	domain          *domain.Domain
	primaryStrategy ComputerOpponentStrategy
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

	return &GameStateContext{
		Opponent:      opponentRec,
		LanceInstance: lanceInstance,
		OwnMechs:      ownMechs,
		EnemyMechs:    enemyMechs,
		Sectors:       sectors,
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

	sb.WriteString(fmt.Sprintf("Turn %d. You command the following mechs:\n", state.TurnNumber))
	for _, m := range state.OwnMechs {
		sectorName := "unknown"
		for _, sec := range state.Sectors {
			if sec.Instance.ID == m.MechaSectorInstanceID {
				sectorName = sec.Design.Name
				break
			}
		}
		sb.WriteString(fmt.Sprintf("  - Mech ID: %s, Callsign: %s, Status: %s, Location: %s (sector_instance_id: %s), Armor: %d, Structure: %d\n",
			m.ID, m.Callsign, m.Status, sectorName, m.MechaSectorInstanceID,
			m.CurrentArmor, m.CurrentStructure))
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
			sb.WriteString(fmt.Sprintf("  - Mech ID: %s, Callsign: %s, Location: %s (sector_instance_id: %s)\n",
				em.Instance.ID, em.Instance.Callsign, sectorName, em.Instance.MechaSectorInstanceID))
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
		sb.WriteString(fmt.Sprintf("  - %s (id: %s, terrain: %s, elevation: %d) → [%s]\n",
			sec.Design.Name, sec.Instance.ID, sec.Design.TerrainType, sec.Design.Elevation,
			strings.Join(linkNames, ", ")))
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
Include one entry per mech. Use exact IDs from above. Only move to directly adjacent sectors.
Attack targets must be in the same sector or one adjacent sector (after movement).
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

		targetSectorInstanceID := s.pickMovementTarget(l, opp, mech, state)

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
func (s *ruleBasedStrategy) pickAttackTarget(opp *mecha_record.MechaComputerOpponent, fromSectorID string, state *GameStateContext) string {
	var candidates []*mechState
	for _, em := range state.EnemyMechs {
		if em.Instance.Status == mecha_record.MechInstanceStatusDestroyed {
			continue
		}
		dist := s.sectorDistance(fromSectorID, em.Instance.MechaSectorInstanceID, state)
		if dist <= 1 {
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
// empty string to stay in place.
func (s *ruleBasedStrategy) pickMovementTarget(_ logger.Logger, opp *mecha_record.MechaComputerOpponent, mech *mecha_record.MechaMechInstance, state *GameStateContext) string {
	adjacentIDs := s.getAdjacentSectorIDs(mech.MechaSectorInstanceID, state)
	if len(adjacentIDs) == 0 {
		return ""
	}

	if opp.Aggression >= 7 && len(state.EnemyMechs) > 0 {
		// Advance toward nearest enemy
		target := s.findNearestEnemy(mech.MechaSectorInstanceID, state)
		if target != "" {
			// Pick the adjacent sector that gets us closer (simple: first adjacent)
			best := s.pickBestAdvanceStep(mech.MechaSectorInstanceID, target, adjacentIDs, opp.IQ, state)
			if best != "" {
				return best
			}
		}
	} else if opp.Aggression <= 3 {
		// Retreat to starting / best-covered sector
		best := s.pickBestDefensiveStep(adjacentIDs, opp.IQ, state)
		if best != "" {
			return best
		}
	}

	// Hold — no movement
	return ""
}

func (s *ruleBasedStrategy) getAdjacentSectorIDs(sectorInstanceID string, state *GameStateContext) []string {
	for _, sec := range state.Sectors {
		if sec.Instance.ID == sectorInstanceID {
			return sec.LinkDestInstanceIDs
		}
	}
	return nil
}

func (s *ruleBasedStrategy) findNearestEnemy(_ string, state *GameStateContext) string {
	for _, em := range state.EnemyMechs {
		if em.Sector != nil {
			return em.Sector.ID
		}
	}
	return ""
}

func (s *ruleBasedStrategy) pickBestAdvanceStep(_, targetID string, adjacentIDs []string, iq int, state *GameStateContext) string {
	// If an adjacent sector IS the target, move there directly.
	for _, adjID := range adjacentIDs {
		if adjID == targetID {
			return adjID
		}
	}

	// Otherwise pick the adjacent sector that is also adjacent to the target.
	for _, adjID := range adjacentIDs {
		adj := s.getAdjacentSectorIDs(adjID, state)
		for _, next := range adj {
			if next == targetID {
				if iq >= 5 {
					// Prefer cover when high IQ
					return s.pickHighCoverSector([]string{adjID}, state)
				}
				return adjID
			}
		}
	}

	// Fall back to any adjacent sector
	if iq >= 5 {
		return s.pickHighCoverSector(adjacentIDs, state)
	}
	return adjacentIDs[0]
}

func (s *ruleBasedStrategy) pickBestDefensiveStep(adjacentIDs []string, iq int, state *GameStateContext) string {
	if iq >= 5 {
		return s.pickHighCoverSector(adjacentIDs, state)
	}
	return ""
}

// pickHighCoverSector returns the sector with the highest elevation (proxy for
// cover) from the given candidates.
func (s *ruleBasedStrategy) pickHighCoverSector(candidates []string, state *GameStateContext) string {
	best := ""
	bestElev := -999
	for _, id := range candidates {
		for _, sec := range state.Sectors {
			if sec.Instance.ID == id && sec.Design.Elevation > bestElev {
				bestElev = sec.Design.Elevation
				best = id
				break
			}
		}
	}
	return best
}
