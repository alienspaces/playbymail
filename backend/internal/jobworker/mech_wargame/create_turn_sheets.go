package mech_wargame

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

// CreateTurnSheets creates all turn sheet records for the current turn of a mech wargame instance.
func (p *MechWargame) CreateTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance) ([]*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("MechWargame/CreateTurnSheets")

	l.Info("creating mech wargame turn sheets for instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	lanceInstanceRecs, err := p.getLanceInstancesForGameInstance(ctx, gameInstanceRec)
	if err != nil {
		l.Error("failed to get lance instances for game instance >%s< error >%v<", gameInstanceRec.ID, err)
		return nil, err
	}

	l.Info("found >%d< lance instances for game instance >%s<", len(lanceInstanceRecs), gameInstanceRec.ID)

	if len(lanceInstanceRecs) == 0 {
		l.Info("no lance instances found for game instance >%s<", gameInstanceRec.ID)
		return nil, nil
	}

	var errs []error
	var createdTurnSheets []*game_record.GameTurnSheet
	for _, lanceInstanceRec := range lanceInstanceRecs {
		lanceTurnSheets, err := p.createLanceTurnSheets(ctx, gameInstanceRec, lanceInstanceRec)
		if err != nil {
			l.Warn("failed to create turn sheets for lance >%s< error >%v<", lanceInstanceRec.ID, err)
			errs = append(errs, err)
			continue
		}
		createdTurnSheets = append(createdTurnSheets, lanceTurnSheets...)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to create turn sheets for some lances: %v", errs)
	}

	return createdTurnSheets, nil
}

// createLanceTurnSheets creates all of the current game turn's turn sheets for a lance.
func (p *MechWargame) createLanceTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mech_wargame_record.MechWargameLanceInstance) ([]*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("MechWargame/createLanceTurnSheets")

	l.Debug("creating turn sheets for lance instance >%s< turn number >%d<", lanceInstance.ID, gameInstanceRec.CurrentTurn)

	var createdTurnSheets []*game_record.GameTurnSheet

	for turnSheetType := range p.Processors {
		// Join game turn sheets are handled through the subscription workflow, not turn processing.
		if turnSheetType == mech_wargame_record.MechWargameTurnSheetTypeJoinGame {
			continue
		}

		turnSheetRec, err := p.createTurnSheet(ctx, gameInstanceRec, lanceInstance, turnSheetType)
		if err != nil {
			l.Warn("failed to create turn sheet >%s< for lance instance >%s< error >%v<", turnSheetType, lanceInstance.ID, err)
			return nil, err
		}

		if turnSheetRec == nil {
			l.Debug("no turn sheet generated for type >%s< lance instance ID >%s< — skipping", turnSheetType, lanceInstance.ID)
			continue
		}

		l.Debug("created turn sheet >%s< for lance instance ID >%s< turn sheet type >%s< turn number >%d<", turnSheetRec.ID, lanceInstance.ID, turnSheetType, gameInstanceRec.CurrentTurn)
		createdTurnSheets = append(createdTurnSheets, turnSheetRec)
	}

	return createdTurnSheets, nil
}

// createTurnSheet creates a single turn sheet for a lance instance.
func (p *MechWargame) createTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mech_wargame_record.MechWargameLanceInstance, turnSheetType string) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("MechWargame/createTurnSheet")

	l.Debug("creating turn sheet type >%s< for game instance ID >%s< lance instance ID >%s<", turnSheetType, gameInstanceRec.ID, lanceInstance.ID)

	processor, exists := p.Processors[turnSheetType]
	if !exists {
		l.Warn("unsupported sheet type >%s< for lance instance >%s<", turnSheetType, lanceInstance.ID)
		return nil, fmt.Errorf("unsupported sheet type: %s", turnSheetType)
	}

	return processor.CreateNextTurnSheet(ctx, gameInstanceRec, lanceInstance)
}
