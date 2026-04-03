package mecha

import (
	"context"
	"encoding/json"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// MechaJoinGameProcessor handles subscription processing for mecha games.
// Player lances are created at game start by PopulateMechaGameInstanceData;
// this processor exists only to handle the physical-mail scan flow (commander name).
type MechaJoinGameProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewMechaJoinGameProcessor creates a new mecha join game processor.
func NewMechaJoinGameProcessor(l logger.Logger, d *domain.Domain) (*MechaJoinGameProcessor, error) {
	return &MechaJoinGameProcessor{Logger: l, Domain: d}, nil
}

// ProcessGameSubscriptionProcessing is a no-op for mecha games.
// Lances are created at game-start time from the starter template, not at join time.
// The commander name from a physical join-game scan is currently unused but may be
// stored in a future iteration.
func (p *MechaJoinGameProcessor) ProcessGameSubscriptionProcessing(
	ctx context.Context,
	subscriptionRec *game_record.GameSubscription,
	turnSheetRec *game_record.GameTurnSheet,
) error {
	l := p.Logger.WithFunctionContext("MechaJoinGameProcessor/ProcessGameSubscriptionProcessing")

	l.Info("processing mecha subscription ID >%s< (no-op: lances created at game start)", subscriptionRec.ID)

	if turnSheetRec != nil && len(turnSheetRec.ScannedData) > 0 {
		var scanData turnsheet.MechaJoinGameScanData
		if err := json.Unmarshal(turnSheetRec.ScannedData, &scanData); err == nil {
			if err := scanData.Validate(); err == nil {
				l.Info("physical join scan for subscription >%s<: commander name >%s< (stored for future use)", subscriptionRec.ID, scanData.CommanderName)
			}
		}
	}

	return nil
}
