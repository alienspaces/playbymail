package mecha

import (
	"context"
	"fmt"
	"slices"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/mecha/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// ProcessTurnSheets processes all turn sheet records for the current turn of a mecha instance.
func (p *Mecha) ProcessTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance) error {
	l := p.Logger.WithFunctionContext("Mecha/ProcessTurnSheets")

	l.Info("processing mecha turn for instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	// Reset accumulated attacks for this processing run
	p.pendingAttacks = nil

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
// opponent lances for the current turn. Errors are non-fatal to avoid blocking
// the human player's turn.
func (p *Mecha) processComputerOpponentOrders(ctx context.Context, gameInstanceRec *game_record.GameInstance) error {
	l := p.Logger.WithFunctionContext("Mecha/processComputerOpponentOrders")

	l.Info("generating computer opponent orders for game instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	// Get all computer opponents for this game.
	opponentRecs, err := p.Domain.GetManyMechaComputerOpponentRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaComputerOpponentGameID, Val: gameInstanceRec.GameID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get computer opponents: %w", err)
	}

	if len(opponentRecs) == 0 {
		l.Debug("no computer opponents for game >%s<", gameInstanceRec.GameID)
		return nil
	}

	// Get all lance instances for this game instance.
	allLanceInstances, err := p.getLanceInstancesForGameInstance(ctx, gameInstanceRec)
	if err != nil {
		return fmt.Errorf("failed to get lance instances: %w", err)
	}

	for _, opponentRec := range opponentRecs {
		// Find lance instances assigned to this opponent.
		var opponentLanceInstances []*mecha_record.MechaLanceInstance
		for _, li := range allLanceInstances {
			if li.MechaComputerOpponentID.Valid && li.MechaComputerOpponentID.String == opponentRec.ID {
				opponentLanceInstances = append(opponentLanceInstances, li)
			}
		}

		if len(opponentLanceInstances) == 0 {
			l.Warn("no lance instances found for computer opponent >%s< — skipping", opponentRec.ID)
			continue
		}

		for _, lanceInstance := range opponentLanceInstances {
			orders, err := p.DecisionEngine.GenerateOrdersForLance(ctx, gameInstanceRec.ID, lanceInstance, opponentRec, gameInstanceRec.CurrentTurn)
			if err != nil {
				l.Warn("decision engine failed for opponent >%s< lance instance >%s<: %v", opponentRec.Name, lanceInstance.ID, err)
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

			l.Info("applied %d orders for computer opponent >%s< lance instance >%s<", len(orders), opponentRec.Name, lanceInstance.ID)
		}
	}

	return nil
}

// applyComputerOpponentOrder applies a single mech movement order generated by
// the decision engine, enforcing the same rules as the human orders processor:
// destroyed/shutdown/refitting mechs cannot move, and the destination must be
// within the mech's speed budget.
func (p *Mecha) applyComputerOpponentOrder(gameInstanceRec *game_record.GameInstance, order turnsheet.ScannedMechOrder) error {
	l := p.Logger.WithFunctionContext("Mecha/applyComputerOpponentOrder")

	if order.MechInstanceID == "" || order.MoveToSectorInstanceID == "" {
		return nil
	}

	mechInstanceRec, err := p.Domain.GetMechaMechInstanceRec(order.MechInstanceID, nil)
	if err != nil {
		return fmt.Errorf("failed to get mech instance >%s<: %w", order.MechInstanceID, err)
	}

	if mechInstanceRec.GameInstanceID != gameInstanceRec.ID {
		return fmt.Errorf("mech instance >%s< does not belong to game instance >%s<", order.MechInstanceID, gameInstanceRec.ID)
	}

	if mechInstanceRec.Status == mecha_record.MechInstanceStatusDestroyed ||
		mechInstanceRec.Status == mecha_record.MechInstanceStatusShutdown {
		l.Info("mech >%s< is %s — ignoring movement order", order.MechInstanceID, mechInstanceRec.Status)
		return nil
	}

	if mechInstanceRec.IsRefitting {
		l.Info("mech >%s< is refitting — ignoring movement order", order.MechInstanceID)
		return nil
	}

	sectorInstanceRec, err := p.Domain.GetMechaSectorInstanceRec(order.MoveToSectorInstanceID, nil)
	if err != nil {
		return fmt.Errorf("failed to get sector instance >%s<: %w", order.MoveToSectorInstanceID, err)
	}

	if sectorInstanceRec.GameInstanceID != gameInstanceRec.ID {
		return fmt.Errorf("sector instance >%s< does not belong to game instance >%s<", order.MoveToSectorInstanceID, gameInstanceRec.ID)
	}

	// Validate destination is within the mech's speed budget.
	chassisRec, err := p.Domain.GetMechaChassisRec(mechInstanceRec.MechaChassisID, nil)
	if err != nil {
		return fmt.Errorf("failed to get chassis >%s< for movement validation: %w", mechInstanceRec.MechaChassisID, err)
	}

	ordersProc, ok := p.Processors[mecha_record.MechaTurnSheetTypeOrders].(*turn_sheet_processor.MechaOrdersProcessor)
	if !ok {
		return fmt.Errorf("orders processor unavailable — cannot validate movement for mech >%s<", order.MechInstanceID)
	}

	_, reachable := ordersProc.IsSectorReachableWithinSpeed(l, gameInstanceRec.ID, mechInstanceRec.MechaSectorInstanceID, order.MoveToSectorInstanceID, chassisRec.Speed)
	if !reachable {
		l.Warn("mech >%s< cannot reach sector >%s< within speed budget %d", order.MechInstanceID, order.MoveToSectorInstanceID, chassisRec.Speed)
		return nil
	}

	mechInstanceRec.MechaSectorInstanceID = order.MoveToSectorInstanceID
	if _, err := p.Domain.UpdateMechaMechInstanceRec(mechInstanceRec); err != nil {
		return fmt.Errorf("failed to update mech instance >%s<: %w", order.MechInstanceID, err)
	}

	return nil
}

// processLanceTurnSheets processes all turn sheets for a specific lance instance.
func (p *Mecha) processLanceTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mecha_record.MechaLanceInstance) error {
	l := p.Logger.WithFunctionContext("Mecha/processLanceTurnSheets")

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

// collectAttacksFromOrdersSheet extracts attack declarations from an orders
// turn sheet and appends them to p.pendingAttacks.
func (p *Mecha) collectAttacksFromOrdersSheet(l logger.Logger, turnSheet *game_record.GameTurnSheet) {
	if turnSheet.SheetType != mecha_record.MechaTurnSheetTypeOrders {
		return
	}
	ordersProc, ok := p.Processors[mecha_record.MechaTurnSheetTypeOrders]
	if !ok {
		return
	}
	op, ok := ordersProc.(*turn_sheet_processor.MechaOrdersProcessor)
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

// processTurnSheet processes a single turn sheet for a lance instance.
func (p *Mecha) processTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mecha_record.MechaLanceInstance, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("Mecha/processTurnSheet")

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

	if err := processor.ProcessTurnSheetResponse(ctx, gameInstanceRec, lanceInstance, turnSheet); err != nil {
		return err
	}

	// Collect attack declarations from orders sheets for combat resolution
	p.collectAttacksFromOrdersSheet(l, turnSheet)

	return nil
}

// getTurnSheetsForLance retrieves turn sheets for a specific lance instance and turn.
func (p *Mecha) getTurnSheetsForLance(lanceInstance *mecha_record.MechaLanceInstance, turnNumber int) ([]*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("Mecha/getTurnSheetsForLance")

	mechaTurnSheetRecs, err := p.Domain.GetManyMechaTurnSheetRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: mecha_record.FieldMechaTurnSheetMechaLanceInstanceID,
					Val: lanceInstance.ID,
				},
			},
		},
	)
	if err != nil {
		l.Error("failed to get mecha turn sheets for lance >%s< turn >%d< error >%v<", lanceInstance.ID, turnNumber, err)
		return nil, err
	}

	var turnSheetRecs []*game_record.GameTurnSheet
	for _, mwTurnSheet := range mechaTurnSheetRecs {
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
