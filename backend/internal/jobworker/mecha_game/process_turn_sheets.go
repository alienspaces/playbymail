package mecha_game

import (
	"context"
	"fmt"
	"slices"
	"strings"

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
		l.Warn("failed to get squad instances for game instance >%s< error >%v<", gameInstanceRec.ID, err)
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
		l.Info("no computer opponents for game >%s<", gameInstanceRec.GameID)
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

			// LLM-driven strategies occasionally return human-readable
			// sector / mech *names* instead of UUIDs despite explicit
			// prompt instructions. Defensively map them back to IDs here
			// so a single hallucinated field doesn't discard the turn's
			// movement or attack orders.
			orders = p.resolveAIOrderIDs(l, gameInstanceRec.ID, opponentRec.Name, orders)

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
// within the mech's speed budget. When the move succeeds, a sighting-report
// movement event is appended to every player squad in the game so the turn
// sheet's "What Happened" panel surfaces enemy positioning changes.
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

	fromSectorInstanceID := mechInstanceRec.MechaGameSectorInstanceID
	mechInstanceRec.MechaGameSectorInstanceID = order.MoveToSectorInstanceID
	if _, err := p.Domain.UpdateMechaGameMechInstanceRec(mechInstanceRec); err != nil {
		return fmt.Errorf("failed to update mech instance >%s<: %w", order.MechInstanceID, err)
	}

	l.Info("opponent mech >%s< moved from sector >%s< to sector >%s<",
		mechInstanceRec.Callsign, fromSectorInstanceID, order.MoveToSectorInstanceID)

	p.recordOpponentMovement(l, gameInstanceRec, mechInstanceRec, fromSectorInstanceID, order.MoveToSectorInstanceID)

	return nil
}

// recordOpponentMovement appends a movement event to every player-owned squad
// instance in the game so the player's "What Happened" panel surfaces enemy
// positioning changes. Failure to write events is logged but not fatal — the
// movement itself has already been persisted.
func (p *MechaGame) recordOpponentMovement(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	mechInstanceRec *mecha_game_record.MechaGameMechInstance,
	fromSectorInstanceID, toSectorInstanceID string,
) {
	fromName := sectorDisplayName(p, fromSectorInstanceID)
	toName := sectorDisplayName(p, toSectorInstanceID)

	var message string
	switch {
	case fromName != "" && toName != "" && fromName != toName:
		message = fmt.Sprintf("Enemy mech %s moved from %s to %s.", mechInstanceRec.Callsign, fromName, toName)
	case toName != "":
		message = fmt.Sprintf("Enemy mech %s moved to %s.", mechInstanceRec.Callsign, toName)
	default:
		message = fmt.Sprintf("Enemy mech %s changed position.", mechInstanceRec.Callsign)
	}

	squads, err := p.Domain.GetManyMechaGameSquadInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSquadInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		l.Warn("failed to load squad instances for opponent movement event: %v", err)
		return
	}

	evt := turnsheet.TurnEvent{
		Category: turnsheet.TurnEventCategoryMovement,
		Icon:     turnsheet.TurnEventIconMovement,
		Message:  message,
	}

	for _, squad := range squads {
		// Skip AI-owned squads — events only surface on player turn sheets.
		if squad.MechaGameComputerOpponentID.Valid {
			continue
		}
		if err := turnsheet.AppendMechaGameTurnEvent(squad, evt); err != nil {
			l.Warn("failed to append opponent movement event to squad >%s<: %v", squad.ID, err)
			continue
		}
		if _, err := p.Domain.UpdateMechaGameSquadInstanceRec(squad); err != nil {
			l.Warn("failed to persist opponent movement event on squad >%s<: %v", squad.ID, err)
		}
	}
}

// sectorDisplayName resolves a sector instance ID to its design name for
// player-facing events. Returns an empty string on any lookup failure so
// the caller can fall back to a generic phrasing instead of crashing.
func sectorDisplayName(p *MechaGame, sectorInstanceID string) string {
	if sectorInstanceID == "" {
		return ""
	}
	inst, err := p.Domain.GetMechaGameSectorInstanceRec(sectorInstanceID, nil)
	if err != nil || inst == nil {
		return ""
	}
	design, err := p.Domain.GetMechaGameSectorRec(inst.MechaGameSectorID, nil)
	if err != nil || design == nil {
		return ""
	}
	return design.Name
}

// resolveAIOrderIDs inspects each AI-generated order and, when a sector or
// mech reference looks like a display name rather than a UUID, resolves it
// back to the correct instance ID. LLM-driven strategies occasionally emit
// sector names ("Drop Zone Alpha") or mech callsigns in fields the engine
// expects to contain UUIDs; without this rescue the order is dropped and the
// opponent effectively forfeits the turn.
func (p *MechaGame) resolveAIOrderIDs(l logger.Logger, gameInstanceID, opponentName string, orders []turnsheet.ScannedMechOrder) []turnsheet.ScannedMechOrder {
	if len(orders) == 0 {
		return orders
	}

	// Build lookup tables lazily — skip the queries entirely when every
	// field already looks like a UUID.
	needsResolve := false
	for _, o := range orders {
		if !looksLikeUUID(o.MechInstanceID) ||
			(o.MoveToSectorInstanceID != "" && !looksLikeUUID(o.MoveToSectorInstanceID)) ||
			(o.AttackTargetMechInstanceID != "" && !looksLikeUUID(o.AttackTargetMechInstanceID)) {
			needsResolve = true
			break
		}
	}
	if !needsResolve {
		return orders
	}

	sectorByName, err := p.buildSectorInstanceNameIndex(gameInstanceID)
	if err != nil {
		l.Warn("failed to build sector name index for opponent >%s<: %v", opponentName, err)
	}
	mechByCallsign, err := p.buildMechInstanceCallsignIndex(gameInstanceID)
	if err != nil {
		l.Warn("failed to build mech callsign index for opponent >%s<: %v", opponentName, err)
	}

	for i := range orders {
		o := &orders[i]
		if o.MechInstanceID != "" && !looksLikeUUID(o.MechInstanceID) {
			if id, ok := mechByCallsign[strings.ToLower(strings.TrimSpace(o.MechInstanceID))]; ok {
				l.Info("resolved opponent >%s< mech reference >%s< to >%s<", opponentName, o.MechInstanceID, id)
				o.MechInstanceID = id
			} else {
				l.Warn("opponent >%s< returned unresolvable mech reference >%s<", opponentName, o.MechInstanceID)
			}
		}
		if o.MoveToSectorInstanceID != "" && !looksLikeUUID(o.MoveToSectorInstanceID) {
			if id, ok := sectorByName[strings.ToLower(strings.TrimSpace(o.MoveToSectorInstanceID))]; ok {
				l.Info("resolved opponent >%s< sector reference >%s< to >%s<", opponentName, o.MoveToSectorInstanceID, id)
				o.MoveToSectorInstanceID = id
			} else {
				l.Warn("opponent >%s< returned unresolvable sector reference >%s< — dropping movement", opponentName, o.MoveToSectorInstanceID)
				o.MoveToSectorInstanceID = ""
			}
		}
		if o.AttackTargetMechInstanceID != "" && !looksLikeUUID(o.AttackTargetMechInstanceID) {
			if id, ok := mechByCallsign[strings.ToLower(strings.TrimSpace(o.AttackTargetMechInstanceID))]; ok {
				l.Info("resolved opponent >%s< target reference >%s< to >%s<", opponentName, o.AttackTargetMechInstanceID, id)
				o.AttackTargetMechInstanceID = id
			} else {
				l.Warn("opponent >%s< returned unresolvable target reference >%s< — dropping attack", opponentName, o.AttackTargetMechInstanceID)
				o.AttackTargetMechInstanceID = ""
			}
		}
	}
	return orders
}

func (p *MechaGame) buildSectorInstanceNameIndex(gameInstanceID string) (map[string]string, error) {
	sectorInsts, err := p.Domain.GetManyMechaGameSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSectorInstanceGameInstanceID, Val: gameInstanceID},
		},
	})
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(sectorInsts))
	for _, si := range sectorInsts {
		sectorRec, err := p.Domain.GetMechaGameSectorRec(si.MechaGameSectorID, nil)
		if err != nil || sectorRec == nil {
			continue
		}
		out[strings.ToLower(strings.TrimSpace(sectorRec.Name))] = si.ID
	}
	return out, nil
}

func (p *MechaGame) buildMechInstanceCallsignIndex(gameInstanceID string) (map[string]string, error) {
	mechInsts, err := p.Domain.GetManyMechaGameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameMechInstanceGameInstanceID, Val: gameInstanceID},
		},
	})
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(mechInsts))
	for _, mi := range mechInsts {
		out[strings.ToLower(strings.TrimSpace(mi.Callsign))] = mi.ID
	}
	return out, nil
}

// looksLikeUUID is a cheap heuristic: length 36 with dashes in the
// canonical positions. We deliberately avoid a strict regex/parse so
// mixed-case IDs and any valid UUID string passes through untouched.
func looksLikeUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	return s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-'
}

// processSquadTurnSheets processes all turn sheets for a specific squad instance.
func (p *MechaGame) processSquadTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance, squadInstance *mecha_game_record.MechaGameSquadInstance) error {
	l := p.Logger.WithFunctionContext("MechaGame/processSquadTurnSheets")

	l.Info("processing turn sheets for squad >%s< turn >%d<", squadInstance.ID, gameInstanceRec.CurrentTurn)

	turnSheetRecs, err := p.getTurnSheetsForSquad(squadInstance, gameInstanceRec.CurrentTurn)
	if err != nil {
		l.Warn("failed to get turn sheets for squad >%s< turn >%d< error >%v<", squadInstance.ID, gameInstanceRec.CurrentTurn, err)
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

	l.Info("processing turn sheet >%s< type >%s< for squad >%s<", turnSheet.ID, turnSheet.SheetType, squadInstance.ID)

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
		l.Warn("failed to get mecha turn sheets for squad >%s< turn >%d< error >%v<", squadInstance.ID, turnNumber, err)
		return nil, err
	}

	var turnSheetRecs []*game_record.GameTurnSheet
	for _, mwTurnSheet := range mechaGameTurnSheetRecs {
		turnSheetRec, err := p.Domain.GetGameTurnSheetRec(mwTurnSheet.GameTurnSheetID, nil)
		if err != nil {
			l.Warn("failed to get game turn sheet >%s< error >%v<", mwTurnSheet.GameTurnSheetID, err)
			return nil, err
		}
		if turnSheetRec.TurnNumber != turnNumber {
			continue
		}
		turnSheetRecs = append(turnSheetRecs, turnSheetRec)
	}

	return turnSheetRecs, nil
}
