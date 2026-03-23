package mech_wargame

import (
	"context"
	"fmt"
	"slices"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

// ProcessTurnSheets processes all turn sheet records for the current turn of a mech wargame instance.
func (p *MechWargame) ProcessTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance) error {
	l := p.Logger.WithFunctionContext("MechWargame/ProcessTurnSheets")

	l.Info("processing mech wargame turn for instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	lanceInstanceRecs, err := p.getLanceInstancesForGameInstance(ctx, gameInstanceRec)
	if err != nil {
		l.Error("failed to get lance instances for game instance >%s< error >%v<", gameInstanceRec.ID, err)
		return err
	}

	l.Info("found >%d< lance instances for game instance >%s<", len(lanceInstanceRecs), gameInstanceRec.ID)

	if len(lanceInstanceRecs) == 0 {
		l.Info("no lance instances found for game instance >%s<", gameInstanceRec.ID)
		return nil
	}

	var errs []error
	for _, lanceInstanceRec := range lanceInstanceRecs {
		if err := p.processLanceTurnSheets(ctx, gameInstanceRec, lanceInstanceRec); err != nil {
			l.Warn("failed to process turn sheets for lance >%s< error >%v<", lanceInstanceRec.ID, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to process turn sheets for some lances: %v", errs)
	}

	return nil
}

// processLanceTurnSheets processes all turn sheets for a specific lance instance.
func (p *MechWargame) processLanceTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mech_wargame_record.MechWargameLanceInstance) error {
	l := p.Logger.WithFunctionContext("MechWargame/processLanceTurnSheets")

	l.Debug("processing turn sheets for lance >%s< turn >%d<", lanceInstance.ID, gameInstanceRec.CurrentTurn)

	turnSheetRecs, err := p.getTurnSheetsForLance(lanceInstance, gameInstanceRec.CurrentTurn)
	if err != nil {
		l.Error("failed to get turn sheets for lance >%s< turn >%d< error >%v<", lanceInstance.ID, gameInstanceRec.CurrentTurn, err)
		return err
	}

	l.Info("found >%d< turn sheets for lance >%s< turn >%d<", len(turnSheetRecs), lanceInstance.ID, gameInstanceRec.CurrentTurn)

	if len(turnSheetRecs) == 0 {
		return nil
	}

	slices.SortFunc(turnSheetRecs, func(a, b *game_record.GameTurnSheet) int {
		return a.SheetOrder - b.SheetOrder
	})

	for _, turnSheet := range turnSheetRecs {
		if err := p.processTurnSheet(ctx, gameInstanceRec, lanceInstance, turnSheet); err != nil {
			l.Warn("failed to process turn sheet >%s< for lance >%s< error >%v<", turnSheet.ID, lanceInstance.ID, err)
			return err
		}
	}

	return nil
}

// processTurnSheet processes a single turn sheet for a lance instance.
func (p *MechWargame) processTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mech_wargame_record.MechWargameLanceInstance, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("MechWargame/processTurnSheet")

	l.Debug("processing turn sheet >%s< type >%s< for lance >%s<", turnSheet.ID, turnSheet.SheetType, lanceInstance.ID)

	if len(turnSheet.ScannedData) == 0 {
		l.Info("skipping turn sheet >%s< — no scanned data (not yet submitted)", turnSheet.ID)
		return nil
	}

	processor, exists := p.Processors[turnSheet.SheetType]
	if !exists {
		l.Warn("unsupported sheet type >%s< for turn sheet >%s<", turnSheet.SheetType, turnSheet.ID)
		return fmt.Errorf("unsupported sheet type: %s", turnSheet.SheetType)
	}

	return processor.ProcessTurnSheetResponse(ctx, gameInstanceRec, lanceInstance, turnSheet)
}

// getTurnSheetsForLance retrieves turn sheets for a specific lance instance and turn.
func (p *MechWargame) getTurnSheetsForLance(lanceInstance *mech_wargame_record.MechWargameLanceInstance, turnNumber int) ([]*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("MechWargame/getTurnSheetsForLance")

	mechWargameTurnSheetRecs, err := p.Domain.GetManyMechWargameTurnSheetRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: mech_wargame_record.FieldMechWargameTurnSheetMechWargameLanceInstanceID,
					Val: lanceInstance.ID,
				},
			},
		},
	)
	if err != nil {
		l.Error("failed to get mech wargame turn sheets for lance >%s< turn >%d< error >%v<", lanceInstance.ID, turnNumber, err)
		return nil, err
	}

	var turnSheetRecs []*game_record.GameTurnSheet
	for _, mwTurnSheet := range mechWargameTurnSheetRecs {
		turnSheetRec, err := p.Domain.GetGameTurnSheetRec(mwTurnSheet.GameTurnSheetID, nil)
		if err != nil {
			l.Error("failed to get game turn sheet >%s< error >%v<", mwTurnSheet.GameTurnSheetID, err)
			return nil, err
		}
		if turnSheetRec.TurnNumber != turnNumber {
			continue
		}
		turnSheetRecs = append(turnSheetRecs, turnSheetRec)
	}

	return turnSheetRecs, nil
}
