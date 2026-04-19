// Package mecha provides turn sheet processing for the mecha game type.
//
// File layout:
//   - mecha_game.go           — MechaGame struct, processor registry, squad instance helpers
//   - process_turn_sheets.go    — ProcessTurnSheets entry point (implements GameTurnProcessor)
//   - create_turn_sheets.go     — CreateTurnSheets entry point (implements GameTurnProcessor)
//   - turn_sheet_processor/     — per-sheet-type business logic processors
package mecha_game

import (
	"context"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/mecha_game/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// MechaGame is the turn sheet processor for mecha games.
//   - Only existing squad instances are processed during turn processing.
//
// Function/method argument order (enforced throughout this package):
//  1. context.Context  — only on interface method boundaries
//  2. logger.Logger    — all package-level helpers and unexported sub-functions
//  3. *domain.Domain   — all package-level helpers and unexported sub-functions
//  4. remaining domain arguments
type MechaGame struct {
	Logger         logger.Logger
	Domain         *domain.Domain
	Processors     map[string]TurnSheetProcessor
	DecisionEngine *ComputerOpponentDecisionEngine
	// pendingAttacks accumulates attack declarations from all order processing
	// within a single ProcessTurnSheets call and is reset at the start of each
	// ProcessTurnSheets call.
	pendingAttacks []turn_sheet_processor.AttackDeclaration
}

// TurnSheetProcessor defines the interface for processing turn sheet business logic in mecha
type TurnSheetProcessor interface {
	// GetSheetType returns the sheet type this processor handles
	GetSheetType() string

	// ProcessTurnSheetResponse processes a single turn sheet response and updates game state
	ProcessTurnSheetResponse(ctx context.Context, gameInstanceRec *game_record.GameInstance, squadInstance *mecha_game_record.MechaGameSquadInstance, turnSheet *game_record.GameTurnSheet) error

	// CreateNextTurnSheet creates a new turn sheet record for the next turn
	CreateNextTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, squadInstance *mecha_game_record.MechaGameSquadInstance) (*game_record.GameTurnSheet, error)
}

// NewMechaGame creates a new mecha turn processor.
func NewMechaGame(l logger.Logger, d *domain.Domain, cfg config.Config) (*MechaGame, error) {
	l = l.WithFunctionContext("NewMechaGame")

	g := &MechaGame{
		Logger:         l,
		Domain:         d,
		DecisionEngine: NewComputerOpponentDecisionEngine(l, d, cfg),
	}

	processors, err := g.initializeTurnSheetProcessors(cfg)
	if err != nil {
		l.Warn("failed to initialize turn sheet processors >%v<", err)
		return nil, err
	}
	g.Processors = processors

	return g, nil
}

// initializeTurnSheetProcessors creates and registers all available mecha turn sheet business processors.
// To add new turn sheet types: 1) Create processor in turn_sheet_processor/ 2) Register here
func (p *MechaGame) initializeTurnSheetProcessors(cfg config.Config) (map[string]TurnSheetProcessor, error) {
	l := p.Logger.WithFunctionContext("MechaGame/initializeTurnSheetProcessors")

	processors := make(map[string]TurnSheetProcessor)

	ordersProcessor, err := turn_sheet_processor.NewMechaGameOrdersProcessor(l, p.Domain)
	if err != nil {
		l.Warn("failed to initialize orders processor >%v<", err)
		return nil, err
	}
	processors[mecha_game_record.MechaGameTurnSheetTypeOrders] = ordersProcessor

	managementProcessor := turn_sheet_processor.NewMechaGameSquadManagementProcessor(l, p.Domain, cfg)
	processors[mecha_game_record.MechaGameTurnSheetTypeSquadManagement] = managementProcessor

	return processors, nil
}

// getSquadInstancesForGameInstance retrieves all squad instances for a game instance.
func (p *MechaGame) getSquadInstancesForGameInstance(_ context.Context, gameInstanceRec *game_record.GameInstance) ([]*mecha_game_record.MechaGameSquadInstance, error) {
	l := p.Logger.WithFunctionContext("MechaGame/getSquadInstancesForGameInstance")

	squadInstanceRecs, err := p.Domain.GetManyMechaGameSquadInstanceRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: mecha_game_record.FieldMechaGameSquadInstanceGameInstanceID,
					Val: gameInstanceRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed to get squad instances error >%v<", err)
		return nil, err
	}

	return squadInstanceRecs, nil
}
