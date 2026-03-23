// Package mech_wargame provides turn sheet processing for the mech_wargame game type.
//
// File layout:
//   - mech_wargame.go           — MechWargame struct, processor registry, lance instance helpers
//   - process_turn_sheets.go    — ProcessTurnSheets entry point (implements GameTurnProcessor)
//   - create_turn_sheets.go     — CreateTurnSheets entry point (implements GameTurnProcessor)
//   - turn_sheet_processor/     — per-sheet-type business logic processors
package mech_wargame

import (
	"context"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/mech_wargame/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

// MechWargame is the turn sheet processor for mech wargame games.
//   - Only existing lance instances are processed during turn processing.
//
// Function/method argument order (enforced throughout this package):
//  1. context.Context  — only on interface method boundaries
//  2. logger.Logger    — all package-level helpers and unexported sub-functions
//  3. *domain.Domain   — all package-level helpers and unexported sub-functions
//  4. remaining domain arguments
type MechWargame struct {
	Logger     logger.Logger
	Domain     *domain.Domain
	Processors map[string]TurnSheetProcessor
}

// TurnSheetProcessor defines the interface for processing turn sheet business logic in mech wargame
type TurnSheetProcessor interface {
	// GetSheetType returns the sheet type this processor handles
	GetSheetType() string

	// ProcessTurnSheetResponse processes a single turn sheet response and updates game state
	ProcessTurnSheetResponse(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mech_wargame_record.MechWargameLanceInstance, turnSheet *game_record.GameTurnSheet) error

	// CreateNextTurnSheet creates a new turn sheet record for the next turn
	CreateNextTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mech_wargame_record.MechWargameLanceInstance) (*game_record.GameTurnSheet, error)
}

// NewMechWargame creates a new mech wargame turn processor.
func NewMechWargame(l logger.Logger, d *domain.Domain) (*MechWargame, error) {
	l = l.WithFunctionContext("NewMechWargame")

	g := &MechWargame{
		Logger: l,
		Domain: d,
	}

	p, err := g.initializeTurnSheetProcessors()
	if err != nil {
		l.Warn("failed to initialize turn sheet processors >%v<", err)
		return nil, err
	}
	g.Processors = p

	return g, nil
}

// initializeTurnSheetProcessors creates and registers all available mech wargame turn sheet business processors.
// To add new turn sheet types: 1) Create processor in turn_sheet_processor/ 2) Register here
func (p *MechWargame) initializeTurnSheetProcessors() (map[string]TurnSheetProcessor, error) {
	l := p.Logger.WithFunctionContext("MechWargame/initializeTurnSheetProcessors")

	processors := make(map[string]TurnSheetProcessor)

	ordersProcessor, err := turn_sheet_processor.NewMechWargameOrdersProcessor(l, p.Domain)
	if err != nil {
		l.Warn("failed to initialize orders processor >%v<", err)
		return nil, err
	}
	processors[mech_wargame_record.MechWargameTurnSheetTypeOrders] = ordersProcessor

	return processors, nil
}

// getLanceInstancesForGameInstance retrieves all lance instances for a game instance.
func (p *MechWargame) getLanceInstancesForGameInstance(_ context.Context, gameInstanceRec *game_record.GameInstance) ([]*mech_wargame_record.MechWargameLanceInstance, error) {
	l := p.Logger.WithFunctionContext("MechWargame/getLanceInstancesForGameInstance")

	lanceInstanceRecs, err := p.Domain.GetManyMechWargameLanceInstanceRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: mech_wargame_record.FieldMechWargameLanceInstanceGameInstanceID,
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
