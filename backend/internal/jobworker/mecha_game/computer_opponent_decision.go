package mecha_game

import (
	"context"
	"fmt"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/agent"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// GameStateContext bundles all information both decision strategies need for a
// single computer opponent's turn.
type GameStateContext struct {
	Opponent      *mecha_game_record.MechaGameComputerOpponent
	SquadInstance *mecha_game_record.MechaGameSquadInstance
	OwnMechs      []*mecha_game_record.MechaGameMechInstance
	EnemyMechs    []*mechState
	Sectors       []*sectorState
	// ChassisCache maps chassis ID to chassis design record for speed lookups.
	ChassisCache map[string]*mecha_game_record.MechaGameChassis
	TurnNumber   int
}

type mechState struct {
	Instance *mecha_game_record.MechaGameMechInstance
	Sector   *mecha_game_record.MechaGameSectorInstance
}

type sectorState struct {
	Instance            *mecha_game_record.MechaGameSectorInstance
	Design              *mecha_game_record.MechaGameSector
	LinkDestInstanceIDs []string
}

// ComputerOpponentStrategy is the interface both strategies implement.
type ComputerOpponentStrategy interface {
	GenerateOrders(ctx context.Context, l logger.Logger, state *GameStateContext) ([]turnsheet.ScannedMechOrder, error)
}

// ComputerOpponentDecisionEngine selects a strategy and generates orders for
// each computer opponent squad during turn processing.
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

// GenerateOrdersForSquad builds the GameStateContext for the given squad instance
// and generates orders using the configured strategy.
func (e *ComputerOpponentDecisionEngine) GenerateOrdersForSquad(
	ctx context.Context,
	gameInstanceID string,
	squadInstance *mecha_game_record.MechaGameSquadInstance,
	opponentRec *mecha_game_record.MechaGameComputerOpponent,
	turnNumber int,
) ([]turnsheet.ScannedMechOrder, error) {
	l := e.logger.WithFunctionContext("ComputerOpponentDecisionEngine/GenerateOrdersForSquad")

	state, err := e.buildGameStateContext(ctx, gameInstanceID, squadInstance, opponentRec, turnNumber)
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
	squadInstance *mecha_game_record.MechaGameSquadInstance,
	opponentRec *mecha_game_record.MechaGameComputerOpponent,
	turnNumber int,
) (*GameStateContext, error) {
	l := e.logger.WithFunctionContext("buildGameStateContext")

	// Get all mech instances for the squad.
	ownMechs, err := e.domain.GetManyMechaGameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameMechInstanceMechaGameSquadInstanceID, Val: squadInstance.ID},
		},
	})
	if err != nil {
		l.Warn("failed to get own mech instances: %v", err)
		return nil, fmt.Errorf("failed to get own mech instances: %w", err)
	}

	// Get all mech instances in this game instance (for enemy detection).
	allMechs, err := e.domain.GetManyMechaGameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameMechInstanceGameInstanceID, Val: gameInstanceID},
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
		if m.Status == mecha_game_record.MechInstanceStatusDestroyed {
			continue
		}
		sectorInst, err := e.domain.GetMechaGameSectorInstanceRec(m.MechaGameSectorInstanceID, nil)
		if err != nil {
			continue
		}
		enemyMechs = append(enemyMechs, &mechState{Instance: m, Sector: sectorInst})
	}

	// Sector graph
	sectorInstances, err := e.domain.GetManyMechaGameSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSectorInstanceGameInstanceID, Val: gameInstanceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get sector instances: %w", err)
	}

	var sectors []*sectorState
	for _, si := range sectorInstances {
		design, err := e.domain.GetMechaGameSectorRec(si.MechaGameSectorID, nil)
		if err != nil {
			continue
		}

		links, err := e.domain.GetManyMechaGameSectorLinkRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_game_record.FieldMechaGameSectorLinkFromMechaGameSectorID, Val: si.MechaGameSectorID},
			},
		})
		if err != nil {
			links = nil
		}

		var linkDestIDs []string
		for _, lnk := range links {
			for _, destSI := range sectorInstances {
				if destSI.MechaGameSectorID == lnk.ToMechaGameSectorID {
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
	chassisCache := make(map[string]*mecha_game_record.MechaGameChassis)
	allMechsForChassis := make([]*mecha_game_record.MechaGameMechInstance, 0, len(ownMechs)+len(enemyMechs))
	allMechsForChassis = append(allMechsForChassis, ownMechs...)
	for _, em := range enemyMechs {
		allMechsForChassis = append(allMechsForChassis, em.Instance)
	}

	for _, m := range allMechsForChassis {
		if m.MechaGameChassisID == "" || chassisCache[m.MechaGameChassisID] != nil {
			continue
		}
		if cr, err := e.domain.GetMechaGameChassisRec(m.MechaGameChassisID, nil); err == nil {
			chassisCache[m.MechaGameChassisID] = cr
		}
	}

	return &GameStateContext{
		Opponent:      opponentRec,
		SquadInstance: squadInstance,
		OwnMechs:      ownMechs,
		EnemyMechs:    enemyMechs,
		Sectors:       sectors,
		ChassisCache:  chassisCache,
		TurnNumber:    turnNumber,
	}, nil
}
