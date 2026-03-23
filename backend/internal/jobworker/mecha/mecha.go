// Package mecha provides turn sheet processing for the mecha game type.
//
// File layout:
//   - mecha.go           — Mecha struct, processor registry, lance instance helpers
//   - process_turn_sheets.go    — ProcessTurnSheets entry point (implements GameTurnProcessor)
//   - create_turn_sheets.go     — CreateTurnSheets entry point (implements GameTurnProcessor)
//   - turn_sheet_processor/     — per-sheet-type business logic processors
package mecha

import (
	"context"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/mecha/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// Mecha is the turn sheet processor for mecha games.
//   - Only existing lance instances are processed during turn processing.
//
// Function/method argument order (enforced throughout this package):
//  1. context.Context  — only on interface method boundaries
//  2. logger.Logger    — all package-level helpers and unexported sub-functions
//  3. *domain.Domain   — all package-level helpers and unexported sub-functions
//  4. remaining domain arguments
type Mecha struct {
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
	ProcessTurnSheetResponse(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mecha_record.MechaLanceInstance, turnSheet *game_record.GameTurnSheet) error

	// CreateNextTurnSheet creates a new turn sheet record for the next turn
	CreateNextTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mecha_record.MechaLanceInstance) (*game_record.GameTurnSheet, error)
}

// NewMecha creates a new mecha turn processor.
func NewMecha(l logger.Logger, d *domain.Domain, cfg config.Config) (*Mecha, error) {
	l = l.WithFunctionContext("NewMecha")

	g := &Mecha{
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
func (p *Mecha) initializeTurnSheetProcessors(cfg config.Config) (map[string]TurnSheetProcessor, error) {
	l := p.Logger.WithFunctionContext("Mecha/initializeTurnSheetProcessors")

	processors := make(map[string]TurnSheetProcessor)

	ordersProcessor, err := turn_sheet_processor.NewMechaOrdersProcessor(l, p.Domain)
	if err != nil {
		l.Warn("failed to initialize orders processor >%v<", err)
		return nil, err
	}
	processors[mecha_record.MechaTurnSheetTypeOrders] = ordersProcessor

	managementProcessor := turn_sheet_processor.NewMechaLanceManagementProcessor(l, p.Domain, cfg)
	processors[mecha_record.MechaTurnSheetTypeLanceManagement] = managementProcessor

	return processors, nil
}

// getLanceInstancesForGameInstance retrieves all lance instances for a game instance.
func (p *Mecha) getLanceInstancesForGameInstance(_ context.Context, gameInstanceRec *game_record.GameInstance) ([]*mecha_record.MechaLanceInstance, error) {
	l := p.Logger.WithFunctionContext("Mecha/getLanceInstancesForGameInstance")

	lanceInstanceRecs, err := p.Domain.GetManyMechaLanceInstanceRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: mecha_record.FieldMechaLanceInstanceGameInstanceID,
					Val: gameInstanceRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Error("failed to get lance instances error >%v<", err)
		return nil, err
	}

	return lanceInstanceRecs, nil
}
