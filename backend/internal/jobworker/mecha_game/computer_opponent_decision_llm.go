package mecha_game

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/agent"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

type llmStrategy struct {
	textAgent agent.TextAgent
}

func (s *llmStrategy) GenerateOrders(ctx context.Context, l logger.Logger, state *GameStateContext) ([]turnsheet.ScannedMechOrder, error) {
	l = l.WithFunctionContext("llmStrategy/GenerateOrders")

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
		l.Warn("LLM text generation failed: %v", err)
		return nil, fmt.Errorf("LLM text generation failed: %w", err)
	}

	// Extract JSON from the response (the model may wrap it in markdown)
	jsonStr := extractJSONFromResponse(resp)
	if jsonStr == "" {
		l.Warn("LLM response contained no JSON: %s", resp)
		return nil, fmt.Errorf("LLM response contained no JSON: %s", resp)
	}

	var scanData turnsheet.OrdersScanData
	if err := json.Unmarshal([]byte(jsonStr), &scanData); err != nil {
		l.Warn("failed to parse LLM orders JSON: %v", err)
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
			if sec.Instance.ID == m.MechaGameSectorInstanceID {
				sectorName = sec.Design.Name
				break
			}
		}
		speed := 1
		if state.ChassisCache != nil {
			if cr, ok := state.ChassisCache[m.MechaGameChassisID]; ok {
				speed = cr.Speed
			}
		}
		fmt.Fprintf(&sb, "  - Mech ID: %s, Callsign: %s, Status: %s, Location: %s (sector_instance_id: %s), Armor: %d, Structure: %d, Speed: %d (max sector hops per turn)\n",
			m.ID, m.Callsign, m.Status, sectorName, m.MechaGameSectorInstanceID,
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
				em.Instance.ID, em.Instance.Callsign, sectorName, em.Instance.MechaGameSectorInstanceID)
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
