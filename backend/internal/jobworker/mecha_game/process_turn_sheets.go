package mecha_game

import (
	"context"
	"fmt"
	"slices"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/mecha_game/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// ProcessTurnSheets processes all turn sheet records for the current turn of a mecha instance.
func (p *MechaGame) ProcessTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance) error {
	l := p.Logger.WithFunctionContext("MechaGame/ProcessTurnSheets")

	l.Info("processing mecha turn for instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	// Reset accumulated attacks for this processing run
	p.pendingAttacks = nil

	squadInstanceRecs, err := p.getSquadInstancesForGameInstance(ctx, gameInstanceRec)
	if err != nil {
		l.Error("failed to get squad instances for game instance >%s< error >%v<", gameInstanceRec.ID, err)
		return err
	}

	l.Info("found >%d< squad instances for game instance >%s<", len(squadInstanceRecs), gameInstanceRec.ID)

	if len(squadInstanceRecs) == 0 {
		l.Info("no squad instances found for game instance >%s<", gameInstanceRec.ID)
		return nil
	}

	var errs []error
	for _, squadInstanceRec := range squadInstanceRecs {
		if err := p.processSquadTurnSheets(ctx, gameInstanceRec, squadInstanceRec); err != nil {
			l.Warn("failed to process turn sheets for squad >%s< error >%v<", squadInstanceRec.ID, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to process turn sheets for some squads: %v", errs)
	}

	if err := p.processComputerOpponentOrders(ctx, gameInstanceRec); err != nil {
		l.Warn("failed to process computer opponent orders >%v< — continuing (non-fatal)", err)
	}

	// Resolve combat from all collected attack declarations
	xpMap, err := p.resolveCombat(ctx, l, gameInstanceRec, p.pendingAttacks)
	if err != nil {
		l.Warn("failed to resolve combat >%v< — continuing (non-fatal)", err)
		xpMap = nil
	}
	p.pendingAttacks = nil

	// Run end-of-turn lifecycle (heat dissipation, auto-repair, XP/level-up, supply accrual)
	if err := p.runEndOfTurn(ctx, l, gameInstanceRec, xpMap); err != nil {
		l.Warn("failed to run end-of-turn lifecycle >%v< — continuing (non-fatal)", err)
	}

	return nil
}

// processComputerOpponentOrders generates and applies orders for all computer
// opponent squads for the current turn. Errors are non-fatal to avoid blocking
// the human player's turn.
func (p *MechaGame) processComputerOpponentOrders(ctx context.Context, gameInstanceRec *game_record.GameInstance) error {
	l := p.Logger.WithFunctionContext("MechaGame/processComputerOpponentOrders")

	l.Info("generating computer opponent orders for game instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	// Get all computer opponents for this game.
	opponentRecs, err := p.Domain.GetManyMechaGameComputerOpponentRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameComputerOpponentGameID, Val: gameInstanceRec.GameID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get computer opponents: %w", err)
	}

	if len(opponentRecs) == 0 {
		l.Debug("no computer opponents for game >%s<", gameInstanceRec.GameID)
		return nil
	}

	// Get all squad instances for this game instance.
	allSquadInstances, err := p.getSquadInstancesForGameInstance(ctx, gameInstanceRec)
	if err != nil {
		return fmt.Errorf("failed to get squad instances: %w", err)
	}

	for _, opponentRec := range opponentRecs {
		// Find squad instances assigned to this opponent.
		var opponentSquadInstances []*mecha_game_record.MechaGameSquadInstance
		for _, si := range allSquadInstances {
			if si.MechaGameComputerOpponentID.Valid && si.MechaGameComputerOpponentID.String == opponentRec.ID {
				opponentSquadInstances = append(opponentSquadInstances, si)
			}
		}

		if len(opponentSquadInstances) == 0 {
			l.Warn("no squad instances found for computer opponent >%s< — skipping", opponentRec.ID)
			continue
		}

		for _, squadInstance := range opponentSquadInstances {
			orders, err := p.DecisionEngine.GenerateOrdersForSquad(ctx, gameInstanceRec.ID, squadInstance, opponentRec, gameInstanceRec.CurrentTurn)
			if err != nil {
				l.Warn("decision engine failed for opponent >%s< squad instance >%s<: %v", opponentRec.Name, squadInstance.ID, err)
				continue
			}

			for _, order := range orders {
				if err := p.applyComputerOpponentOrder(gameInstanceRec, order); err != nil {
					l.Warn("failed to apply order for mech >%s< opponent >%s<: %v", order.MechInstanceID, opponentRec.Name, err)
					continue
				}
				// Collect attack declarations from AI orders
				if order.MechInstanceID != "" && order.AttackTargetMechInstanceID != "" {
					p.pendingAttacks = append(p.pendingAttacks, turn_sheet_processor.AttackDeclaration{
						AttackerMechInstanceID: order.MechInstanceID,
						TargetMechInstanceID:   order.AttackTargetMechInstanceID,
					})
				}
			}

			l.Info("applied %d orders for computer opponent >%s< squad instance >%s<", len(orders), opponentRec.Name, squadInstance.ID)
		}
	}

	return nil
}

// applyComputerOpponentOrder applies a single mech movement order generated by
// the decision engine, enforcing the same rules as the human orders processor:
// destroyed/shutdown/refitting mechs cannot move, and the destination must be
// within the mech's speed budget.
func (p *MechaGame) applyComputerOpponentOrder(gameInstanceRec *game_record.GameInstance, order turnsheet.ScannedMechOrder) error {
	l := p.Logger.WithFunctionContext("MechaGame/applyComputerOpponentOrder")

	if order.MechInstanceID == "" || order.MoveToSectorInstanceID == "" {
		return nil
	}

	mechInstanceRec, err := p.Domain.GetMechaGameMechInstanceRec(order.MechInstanceID, nil)
	if err != nil {
		return fmt.Errorf("failed to get mech instance >%s<: %w", order.MechInstanceID, err)
	}

	if mechInstanceRec.GameInstanceID != gameInstanceRec.ID {
		return fmt.Errorf("mech instance >%s< does not belong to game instance >%s<", order.MechInstanceID, gameInstanceRec.ID)
	}

	if mechInstanceRec.Status == mecha_game_record.MechInstanceStatusDestroyed ||
		mechInstanceRec.Status == mecha_game_record.MechInstanceStatusShutdown {
		l.Info("mech >%s< is %s — ignoring movement order", order.MechInstanceID, mechInstanceRec.Status)
		return nil
	}

	if mechInstanceRec.IsRefitting {
		l.Info("mech >%s< is refitting — ignoring movement order", order.MechInstanceID)
		return nil
	}

	sectorInstanceRec, err := p.Domain.GetMechaGameSectorInstanceRec(order.MoveToSectorInstanceID, nil)
	if err != nil {
		return fmt.Errorf("failed to get sector instance >%s<: %w", order.MoveToSectorInstanceID, err)
	}

	if sectorInstanceRec.GameInstanceID != gameInstanceRec.ID {
		return fmt.Errorf("sector instance >%s< does not belong to game instance >%s<", order.MoveToSectorInstanceID, gameInstanceRec.ID)
	}

	// Validate destination is within the mech's speed budget.
	chassisRec, err := p.Domain.GetMechaGameChassisRec(mechInstanceRec.MechaGameChassisID, nil)
	if err != nil {
		return fmt.Errorf("failed to get chassis >%s< for movement validation: %w", mechInstanceRec.MechaGameChassisID, err)
	}

	ordersProc, ok := p.Processors[mecha_game_record.MechaGameTurnSheetTypeOrders].(*turn_sheet_processor.MechaGameOrdersProcessor)
	if !ok {
		return fmt.Errorf("orders processor unavailable — cannot validate movement for mech >%s<", order.MechInstanceID)
	}

	_, reachable := ordersProc.IsSectorReachableWithinSpeed(l, gameInstanceRec.ID, mechInstanceRec.MechaGameSectorInstanceID, order.MoveToSectorInstanceID, chassisRec.Speed)
	if !reachable {
		l.Warn("mech >%s< cannot reach sector >%s< within speed budget %d", order.MechInstanceID, order.MoveToSectorInstanceID, chassisRec.Speed)
		return nil
	}

	mechInstanceRec.MechaGameSectorInstanceID = order.MoveToSectorInstanceID
	if _, err := p.Domain.UpdateMechaGameMechInstanceRec(mechInstanceRec); err != nil {
		return fmt.Errorf("failed to update mech instance >%s<: %w", order.MechInstanceID, err)
	}

	return nil
}

// processSquadTurnSheets processes all turn sheets for a specific squad instance.
func (p *MechaGame) processSquadTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance, squadInstance *mecha_game_record.MechaGameSquadInstance) error {
	l := p.Logger.WithFunctionContext("MechaGame/processSquadTurnSheets")

	l.Debug("processing turn sheets for squad >%s< turn >%d<", squadInstance.ID, gameInstanceRec.CurrentTurn)

	turnSheetRecs, err := p.getTurnSheetsForSquad(squadInstance, gameInstanceRec.CurrentTurn)
	if err != nil {
		l.Error("failed to get turn sheets for squad >%s< turn >%d< error >%v<", squadInstance.ID, gameInstanceRec.CurrentTurn, err)
		return err
	}

	l.Info("found >%d< turn sheets for squad >%s< turn >%d<", len(turnSheetRecs), squadInstance.ID, gameInstanceRec.CurrentTurn)

	if len(turnSheetRecs) == 0 {
		return nil
	}

	slices.SortFunc(turnSheetRecs, func(a, b *game_record.GameTurnSheet) int {
		return a.SheetOrder - b.SheetOrder
	})

	for _, turnSheet := range turnSheetRecs {
		if err := p.processTurnSheet(ctx, gameInstanceRec, squadInstance, turnSheet); err != nil {
			l.Warn("failed to process turn sheet >%s< for squad >%s< error >%v<", turnSheet.ID, squadInstance.ID, err)
			return err
		}
	}

	return nil
}

// collectAttacksFromOrdersSheet extracts attack declarations from an orders
// turn sheet and appends them to p.pendingAttacks.
func (p *MechaGame) collectAttacksFromOrdersSheet(l logger.Logger, turnSheet *game_record.GameTurnSheet) {
	if turnSheet.SheetType != mecha_game_record.MechaGameTurnSheetTypeOrders {
		return
	}
	ordersProc, ok := p.Processors[mecha_game_record.MechaGameTurnSheetTypeOrders]
	if !ok {
		return
	}
	op, ok := ordersProc.(*turn_sheet_processor.MechaGameOrdersProcessor)
	if !ok {
		return
	}
	attacks, err := op.ExtractAttackDeclarations(turnSheet)
	if err != nil {
		l.Warn("failed to extract attack declarations: %v", err)
		return
	}
	p.pendingAttacks = append(p.pendingAttacks, attacks...)
}

// processTurnSheet processes a single turn sheet for a squad instance.
func (p *MechaGame) processTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, squadInstance *mecha_game_record.MechaGameSquadInstance, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("MechaGame/processTurnSheet")

	l.Debug("processing turn sheet >%s< type >%s< for squad >%s<", turnSheet.ID, turnSheet.SheetType, squadInstance.ID)

	if len(turnSheet.ScannedData) == 0 {
		l.Info("skipping turn sheet >%s< — no scanned data (not yet submitted)", turnSheet.ID)
		return nil
	}

	processor, exists := p.Processors[turnSheet.SheetType]
	if !exists {
		l.Warn("unsupported sheet type >%s< for turn sheet >%s<", turnSheet.SheetType, turnSheet.ID)
		return fmt.Errorf("unsupported sheet type: %s", turnSheet.SheetType)
	}

	if err := processor.ProcessTurnSheetResponse(ctx, gameInstanceRec, squadInstance, turnSheet); err != nil {
		return err
	}

	// Collect attack declarations from orders sheets for combat resolution
	p.collectAttacksFromOrdersSheet(l, turnSheet)

	return nil
}

// getTurnSheetsForSquad retrieves turn sheets for a specific squad instance and turn.
func (p *MechaGame) getTurnSheetsForSquad(squadInstance *mecha_game_record.MechaGameSquadInstance, turnNumber int) ([]*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("MechaGame/getTurnSheetsForSquad")

	mechaGameTurnSheetRecs, err := p.Domain.GetManyMechaGameTurnSheetRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: mecha_game_record.FieldMechaGameTurnSheetMechaGameSquadInstanceID,
					Val: squadInstance.ID,
				},
			},
		},
	)
	if err != nil {
		l.Error("failed to get mecha turn sheets for squad >%s< turn >%d< error >%v<", squadInstance.ID, turnNumber, err)
		return nil, err
	}

	var turnSheetRecs []*game_record.GameTurnSheet
	for _, mwTurnSheet := range mechaGameTurnSheetRecs {
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
