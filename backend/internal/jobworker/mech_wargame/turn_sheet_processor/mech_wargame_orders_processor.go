package turn_sheet_processor

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

// MechWargameOrdersProcessor implements the TurnSheetProcessor interface
// (defined in the parent mech_wargame package)

// MechWargameOrdersProcessor processes orders turn sheet business logic for mech wargame
type MechWargameOrdersProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewMechWargameOrdersProcessor creates a new mech wargame orders processor.
func NewMechWargameOrdersProcessor(l logger.Logger, d *domain.Domain) (*MechWargameOrdersProcessor, error) {
	l = l.WithFunctionContext("NewMechWargameOrdersProcessor")

	p := &MechWargameOrdersProcessor{
		Logger: l,
		Domain: d,
	}
	return p, nil
}

// GetSheetType returns the sheet type this processor handles (implements TurnSheetProcessor interface).
func (p *MechWargameOrdersProcessor) GetSheetType() string {
	return mech_wargame_record.MechWargameTurnSheetTypeOrders
}

// ProcessTurnSheetResponse processes a single orders turn sheet response (implements TurnSheetProcessor interface).
func (p *MechWargameOrdersProcessor) ProcessTurnSheetResponse(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mech_wargame_record.MechWargameLanceInstance, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("MechWargameOrdersProcessor/ProcessTurnSheetResponse")

	l.Info("processing orders for turn sheet >%s< for lance instance >%s<", turnSheet.ID, lanceInstance.ID)

	if turnSheet.SheetType != mech_wargame_record.MechWargameTurnSheetTypeOrders {
		return fmt.Errorf("invalid sheet type: expected %s, got %s", mech_wargame_record.MechWargameTurnSheetTypeOrders, turnSheet.SheetType)
	}

	var scanData turnsheet.OrdersScanData
	if err := json.Unmarshal(turnSheet.ScannedData, &scanData); err != nil {
		l.Warn("failed to unmarshal scanned data >%v<", err)
		return fmt.Errorf("failed to parse scanned data: %w", err)
	}

	if len(scanData.MechOrders) == 0 {
		l.Info("no mech orders in scanned data — lance stays in place")
		return nil
	}

	for _, order := range scanData.MechOrders {
		if err := p.processMechOrder(l, gameInstanceRec, order); err != nil {
			l.Warn("failed to process order for mech >%s< >%v<", order.MechInstanceID, err)
			// Non-fatal: continue processing other mechs
		}
	}

	return nil
}

// processMechOrder applies a single mech's movement order.
func (p *MechWargameOrdersProcessor) processMechOrder(l logger.Logger, gameInstanceRec *game_record.GameInstance, order turnsheet.ScannedMechOrder) error {
	if order.MechInstanceID == "" {
		return nil
	}

	if order.MoveToSectorInstanceID == "" {
		l.Debug("no movement order for mech >%s< — staying in place", order.MechInstanceID)
		return nil
	}

	mechInstanceRec, err := p.Domain.GetMechWargameMechInstanceRec(order.MechInstanceID, nil)
	if err != nil {
		l.Warn("failed to get mech instance >%s< >%v<", order.MechInstanceID, err)
		return fmt.Errorf("failed to get mech instance: %w", err)
	}

	if mechInstanceRec.GameInstanceID != gameInstanceRec.ID {
		l.Warn("mech instance >%s< does not belong to game instance >%s<", order.MechInstanceID, gameInstanceRec.ID)
		return fmt.Errorf("mech instance does not belong to this game instance")
	}

	if mechInstanceRec.Status == mech_wargame_record.MechInstanceStatusDestroyed {
		l.Info("mech >%s< is destroyed — ignoring movement order", order.MechInstanceID)
		return nil
	}

	// Validate the target sector instance belongs to this game instance
	sectorInstanceRec, err := p.Domain.GetMechWargameSectorInstanceRec(order.MoveToSectorInstanceID, nil)
	if err != nil {
		l.Warn("failed to get sector instance >%s< >%v<", order.MoveToSectorInstanceID, err)
		return fmt.Errorf("failed to get sector instance: %w", err)
	}

	if sectorInstanceRec.GameInstanceID != gameInstanceRec.ID {
		l.Warn("sector instance >%s< does not belong to game instance >%s<", order.MoveToSectorInstanceID, gameInstanceRec.ID)
		return fmt.Errorf("sector instance does not belong to this game instance")
	}

	mechInstanceRec.MechWargameSectorInstanceID = order.MoveToSectorInstanceID

	if _, err := p.Domain.UpdateMechWargameMechInstanceRec(mechInstanceRec); err != nil {
		l.Warn("failed to update mech instance >%s< >%v<", order.MechInstanceID, err)
		return fmt.Errorf("failed to update mech instance: %w", err)
	}

	l.Info("moved mech >%s< to sector instance >%s<", order.MechInstanceID, order.MoveToSectorInstanceID)
	return nil
}

// CreateNextTurnSheet creates a new orders turn sheet for a lance instance (implements TurnSheetProcessor interface).
func (p *MechWargameOrdersProcessor) CreateNextTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, lanceInstance *mech_wargame_record.MechWargameLanceInstance) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("MechWargameOrdersProcessor/CreateNextTurnSheet")

	l.Info("creating orders turn sheet for lance instance >%s<", lanceInstance.ID)

	// Step 1: Get the lance design record
	lanceRec, err := p.Domain.GetMechWargameLanceRec(lanceInstance.MechWargameLanceID, nil)
	if err != nil {
		l.Warn("failed to get lance >%v<", err)
		return nil, fmt.Errorf("failed to get lance: %w", err)
	}

	// Step 2: Get the account user for the lance owner
	accountUserRec, err := p.Domain.GetAccountUserRec(lanceRec.AccountUserID, nil)
	if err != nil {
		l.Warn("failed to get account user >%v<", err)
		return nil, fmt.Errorf("failed to get account user: %w", err)
	}

	// Step 3: Get the game record
	gameRec, err := p.Domain.GetGameRec(gameInstanceRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game >%v<", err)
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Step 4: Get all mech instances for this lance instance
	mechInstances, err := p.Domain.GetManyMechWargameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mech_wargame_record.FieldMechWargameMechInstanceMechWargameLanceInstanceID, Val: lanceInstance.ID},
		},
	})
	if err != nil {
		l.Warn("failed to get mech instances >%v<", err)
		return nil, fmt.Errorf("failed to get mech instances: %w", err)
	}

	// Step 5: Build mech order entries
	var lanceMechs []turnsheet.MechOrderEntry
	var sectorInstanceIDs []string
	sectorInstanceIDSet := make(map[string]bool)

	for _, mechInst := range mechInstances {
		sectorInstanceRec, sErr := p.Domain.GetMechWargameSectorInstanceRec(mechInst.MechWargameSectorInstanceID, nil)
		var sectorName string
		if sErr == nil {
			sectorRec, sErr2 := p.Domain.GetMechWargameSectorRec(sectorInstanceRec.MechWargameSectorID, nil)
			if sErr2 == nil {
				sectorName = sectorRec.Name
			}
		}

		lanceMechs = append(lanceMechs, turnsheet.MechOrderEntry{
			MechInstanceID:    mechInst.ID,
			MechCallsign:      mechInst.Callsign,
			MechStatus:        mechInst.Status,
			CurrentSectorName: sectorName,
		})

		if !sectorInstanceIDSet[mechInst.MechWargameSectorInstanceID] {
			sectorInstanceIDSet[mechInst.MechWargameSectorInstanceID] = true
			sectorInstanceIDs = append(sectorInstanceIDs, mechInst.MechWargameSectorInstanceID)
		}
	}

	// Step 6: Get adjacent sector options from the sectors occupied by this lance's mechs
	availableSectors, err := p.getAdjacentSectorOptions(l, gameInstanceRec.ID, sectorInstanceIDs)
	if err != nil {
		l.Warn("failed to get adjacent sectors >%v<", err)
		// Non-fatal: continue with no movement options
	}

	// Step 7: Get enemy mech instances visible to this lance
	enemyMechs, err := p.getEnemyMechOptions(l, gameInstanceRec, lanceInstance)
	if err != nil {
		l.Warn("failed to get enemy mechs >%v<", err)
		// Non-fatal: continue with no attack options
	}

	// Step 8: Generate turn sheet code
	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		l.Warn("failed to generate turn sheet code >%v<", err)
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	turnNumber := gameInstanceRec.CurrentTurn
	title := "Mech Orders"
	instructions := turnsheet.DefaultOrdersInstructions()

	sheetData := turnsheet.OrdersData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr(gameRec.GameType),
			TurnNumber:            &turnNumber,
			AccountName:           convert.Ptr(accountUserRec.Email),
			TurnSheetTitle:        &title,
			TurnSheetDescription:  convert.Ptr(gameRec.Description),
			TurnSheetInstructions: &instructions,
			TurnSheetCode:         convert.Ptr(turnSheetCode),
		},
		LanceName:        lanceRec.Name,
		LanceMechs:       lanceMechs,
		AvailableSectors: availableSectors,
		EnemyMechs:       enemyMechs,
	}

	sheetDataBytes, err := json.Marshal(sheetData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sheet data: %w", err)
	}

	// Step 9: Create the game_turn_sheet record
	gameTurnSheet := &game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		AccountID:        accountUserRec.AccountID,
		AccountUserID:    lanceRec.AccountUserID,
		TurnNumber:       gameInstanceRec.CurrentTurn,
		SheetType:        mech_wargame_record.MechWargameTurnSheetTypeOrders,
		SheetOrder:       mech_wargame_record.MechWargameSheetOrderForType(mech_wargame_record.MechWargameTurnSheetTypeOrders),
		SheetData:        json.RawMessage(sheetDataBytes),
		IsCompleted:      false,
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
	gameTurnSheet.GameInstanceID = nullstring.FromString(gameInstanceRec.ID)

	createdTurnSheetRec, err := p.Domain.CreateGameTurnSheetRec(gameTurnSheet)
	if err != nil {
		return nil, fmt.Errorf("failed to create turn sheet record: %w", err)
	}

	// Step 10: Link the game_turn_sheet to the lance instance via mech_wargame_turn_sheet
	mechWargameTurnSheet := &mech_wargame_record.MechWargameTurnSheet{
		GameID:                     gameInstanceRec.GameID,
		MechWargameLanceInstanceID: lanceInstance.ID,
		GameTurnSheetID:            createdTurnSheetRec.ID,
	}

	_, err = p.Domain.CreateMechWargameTurnSheetRec(mechWargameTurnSheet)
	if err != nil {
		return nil, fmt.Errorf("failed to create mech wargame turn sheet record: %w", err)
	}

	l.Info("created orders turn sheet >%s< for lance instance >%s< turn >%d<", createdTurnSheetRec.ID, lanceInstance.ID, gameInstanceRec.CurrentTurn)
	return createdTurnSheetRec, nil
}

// getAdjacentSectorOptions collects the sectors adjacent to the given sector instance IDs.
func (p *MechWargameOrdersProcessor) getAdjacentSectorOptions(l logger.Logger, gameInstanceID string, sectorInstanceIDs []string) ([]turnsheet.SectorOption, error) {
	var options []turnsheet.SectorOption
	seen := make(map[string]bool)

	for _, sectorInstID := range sectorInstanceIDs {
		sectorInstRec, err := p.Domain.GetMechWargameSectorInstanceRec(sectorInstID, nil)
		if err != nil {
			l.Warn("failed to get sector instance >%s< >%v<", sectorInstID, err)
			continue
		}

		// Get sector links from this sector
		sectorLinks, err := p.Domain.GetManyMechWargameSectorLinkRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mech_wargame_record.FieldMechWargameSectorLinkFromMechWargameSectorID, Val: sectorInstRec.MechWargameSectorID},
			},
		})
		if err != nil {
			l.Warn("failed to get sector links for sector >%s< >%v<", sectorInstRec.MechWargameSectorID, err)
			continue
		}

		for _, link := range sectorLinks {
			// Find the sector instance for the linked sector in this game instance
			linkedInstances, err := p.Domain.GetManyMechWargameSectorInstanceRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: mech_wargame_record.FieldMechWargameSectorInstanceGameInstanceID, Val: gameInstanceID},
					{Col: mech_wargame_record.FieldMechWargameSectorInstanceMechWargameSectorID, Val: link.ToMechWargameSectorID},
				},
				Limit: 1,
			})
			if err != nil || len(linkedInstances) == 0 {
				continue
			}

			linkedInstID := linkedInstances[0].ID
			if seen[linkedInstID] {
				continue
			}
			seen[linkedInstID] = true

			// Get sector name
			sectorRec, err := p.Domain.GetMechWargameSectorRec(link.ToMechWargameSectorID, nil)
			if err != nil {
				l.Warn("failed to get sector >%s< >%v<", link.ToMechWargameSectorID, err)
				continue
			}

			options = append(options, turnsheet.SectorOption{
				SectorInstanceID: linkedInstID,
				SectorName:       sectorRec.Name,
			})
		}
	}

	return options, nil
}

// getEnemyMechOptions collects all enemy mech instances visible to the given lance.
func (p *MechWargameOrdersProcessor) getEnemyMechOptions(l logger.Logger, gameInstanceRec *game_record.GameInstance, lanceInstance *mech_wargame_record.MechWargameLanceInstance) ([]turnsheet.EnemyMechOption, error) {
	// Get all mech instances for this game instance
	allMechInstances, err := p.Domain.GetManyMechWargameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mech_wargame_record.FieldMechWargameMechInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get mech instances: %w", err)
	}

	var enemies []turnsheet.EnemyMechOption
	for _, mechInst := range allMechInstances {
		if mechInst.MechWargameLanceInstanceID == lanceInstance.ID {
			continue
		}
		if mechInst.Status == mech_wargame_record.MechInstanceStatusDestroyed {
			continue
		}

		sectorName := ""
		sectorInstRec, err := p.Domain.GetMechWargameSectorInstanceRec(mechInst.MechWargameSectorInstanceID, nil)
		if err == nil {
			sectorRec, err2 := p.Domain.GetMechWargameSectorRec(sectorInstRec.MechWargameSectorID, nil)
			if err2 == nil {
				sectorName = sectorRec.Name
			}
		}

		enemies = append(enemies, turnsheet.EnemyMechOption{
			MechInstanceID: mechInst.ID,
			Callsign:       mechInst.Callsign,
			SectorName:     sectorName,
		})
	}

	return enemies, nil
}
