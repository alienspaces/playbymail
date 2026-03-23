package mecha

import (
	"context"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// MechaJoinGameProcessor handles subscription processing for mecha games.
// For the online join flow the lance is created by submitJoinHandler; this
// processor acts as a safety-net (and handles the physical-mail scan flow).
type MechaJoinGameProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewMechaJoinGameProcessor creates a new mecha join game processor.
func NewMechaJoinGameProcessor(l logger.Logger, d *domain.Domain) (*MechaJoinGameProcessor, error) {
	return &MechaJoinGameProcessor{Logger: l, Domain: d}, nil
}

// ProcessGameSubscriptionProcessing creates a default mecha lance for the player if one
// does not already exist. When a physical join-game turn sheet has been scanned the
// CommanderName from the scan data is used as the lance base name.
func (p *MechaJoinGameProcessor) ProcessGameSubscriptionProcessing(
	ctx context.Context,
	subscriptionRec *game_record.GameSubscription,
	turnSheetRec *game_record.GameTurnSheet,
) error {
	l := p.Logger.WithFunctionContext("MechaJoinGameProcessor/ProcessGameSubscriptionProcessing")

	l.Info("processing mecha subscription ID >%s<", subscriptionRec.ID)

	commanderName := ""
	if turnSheetRec != nil && len(turnSheetRec.ScannedData) > 0 {
		var scanData turnsheet.MechaJoinGameScanData
		if err := scanData.Validate(); err == nil {
			commanderName = scanData.CommanderName
		}
	}

	lanceRec, err := p.Domain.CreateDefaultMechaLanceForPlayer(
		subscriptionRec.GameID,
		subscriptionRec.AccountID,
		subscriptionRec.AccountUserID,
		commanderName,
		"",
	)
	if err != nil {
		l.Warn("failed to create default mecha lance for player >%v<", err)
		return err
	}

	l.Info("ensured mecha lance >%s< for player >%s<", lanceRec.ID, subscriptionRec.AccountUserID)

	return nil
}
