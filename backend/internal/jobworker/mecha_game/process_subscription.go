package mecha_game

import (
	"context"
	"encoding/json"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// MechaGameJoinGameProcessor handles subscription processing for mecha games.
// Player squads are created at game start by PopulateMechaGameInstanceData;
// this processor exists only to handle the physical-mail scan flow (commander name).
type MechaGameJoinGameProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewMechaGameJoinGameProcessor creates a new mecha join game processor.
func NewMechaGameJoinGameProcessor(l logger.Logger, d *domain.Domain) (*MechaGameJoinGameProcessor, error) {
	return &MechaGameJoinGameProcessor{Logger: l, Domain: d}, nil
}

// ProcessGameSubscriptionProcessing is a no-op for mecha games.
// Squads are created at game-start time from the starter template, not at join time.
// The commander name from a physical join-game scan is currently unused but may be
// stored in a future iteration.
func (p *MechaGameJoinGameProcessor) ProcessGameSubscriptionProcessing(
	ctx context.Context,
	subscriptionRec *game_record.GameSubscription,
	turnSheetRec *game_record.GameTurnSheet,
) error {
	l := p.Logger.WithFunctionContext("MechaGameJoinGameProcessor/ProcessGameSubscriptionProcessing")

	l.Info("processing mecha subscription ID >%s< (no-op: squads created at game start)", subscriptionRec.ID)

	if turnSheetRec != nil && len(turnSheetRec.ScannedData) > 0 {
		var scanData turnsheet.MechaGameJoinGameScanData
		if err := json.Unmarshal(turnSheetRec.ScannedData, &scanData); err == nil {
			if err := scanData.Validate(); err == nil {
				l.Info("physical join scan for subscription >%s<: commander name >%s< (stored for future use)", subscriptionRec.ID, scanData.CommanderName)
			}
		}
	}

	return nil
}
