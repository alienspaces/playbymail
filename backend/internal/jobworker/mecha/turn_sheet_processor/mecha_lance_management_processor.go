package turn_sheet_processor

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

// MechaLanceManagementProcessor handles processing and creation of lance
// management turn sheets.
type MechaLanceManagementProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
	Config config.Config
}

// NewMechaLanceManagementProcessor creates a new management processor.
func NewMechaLanceManagementProcessor(l logger.Logger, d *domain.Domain, cfg config.Config) *MechaLanceManagementProcessor {
	return &MechaLanceManagementProcessor{Logger: l, Domain: d, Config: cfg}
}

// GetSheetType returns the sheet type this processor handles.
func (p *MechaLanceManagementProcessor) GetSheetType() string {
	return mecha_record.MechaTurnSheetTypeLanceManagement
}

// ProcessTurnSheetResponse processes submitted management orders from a player.
// It validates permissions, deducts supply points, flags mechs as refitting,
// and generates TurnEvents.
func (p *MechaLanceManagementProcessor) ProcessTurnSheetResponse(
	ctx context.Context,
	gameInstanceRec *game_record.GameInstance,
	lanceInstance *mecha_record.MechaLanceInstance,
	turnSheet *game_record.GameTurnSheet,
) error {
	l := p.Logger.WithFunctionContext("MechaLanceManagementProcessor/ProcessTurnSheetResponse")

	l.Info("processing management sheet >%s< for lance instance >%s<", turnSheet.ID, lanceInstance.ID)

	if turnSheet.SheetType != mecha_record.MechaTurnSheetTypeLanceManagement {
		return fmt.Errorf("invalid sheet type: expected %s, got %s",
			mecha_record.MechaTurnSheetTypeLanceManagement, turnSheet.SheetType)
	}

	var scanData turnsheet.LanceManagementScanData
	if err := json.Unmarshal(turnSheet.ScannedData, &scanData); err != nil {
		return fmt.Errorf("failed to parse scanned data: %w", err)
	}

	if len(scanData.MechManagementOrders) == 0 {
		l.Info("no management orders — nothing to process")
		return nil
	}

	// Reload lance instance for supply point mutation
	freshLance, err := p.Domain.GetMechaLanceInstanceRec(lanceInstance.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to reload lance instance: %w", err)
	}

	var spCost int
	var events []turnsheet.TurnEvent

	for _, order := range scanData.MechManagementOrders {
		if order.MechInstanceID == "" {
			continue
		}

		mechInst, err := p.Domain.GetMechaMechInstanceRec(order.MechInstanceID, nil)
		if err != nil {
			l.Warn("mech instance >%s< not found: %v", order.MechInstanceID, err)
			continue
		}

		if mechInst.GameInstanceID != gameInstanceRec.ID {
			l.Warn("mech instance >%s< does not belong to game instance >%s<", order.MechInstanceID, gameInstanceRec.ID)
			continue
		}

		if mechInst.IsRefitting {
			l.Info("mech >%s< already refitting — skipping management order", mechInst.Callsign)
			continue
		}

		// Verify mech is at a depot (starting sector)
		isAtDepot, err := p.mechIsAtDepot(mechInst)
		if err != nil {
			l.Warn("failed to check depot for mech >%s<: %v", mechInst.Callsign, err)
			continue
		}
		if !isAtDepot {
			l.Info("mech >%s< is not at a depot sector — ignoring management order", mechInst.Callsign)
			continue
		}

		mechChanged := false

		// Structure repair
		if order.RepairStructure {
			chassisRec, err := p.Domain.GetMechaChassisRec(mechInst.MechaChassisID, nil)
			if err == nil && mechInst.CurrentStructure < chassisRec.StructurePoints {
				damage := chassisRec.StructurePoints - mechInst.CurrentStructure
				// Cost: 1 SP per 25% max structure block to repair
				cost := max((damage*4+chassisRec.StructurePoints-1)/chassisRec.StructurePoints, 1)
				spCost += cost
				// Full structure restored at end of turn (after IsRefitting cleared)
				mechInst.CurrentStructure = chassisRec.StructurePoints
				mechInst.IsRefitting = true
				mechChanged = true
				events = append(events, turnsheet.TurnEvent{
					Category: turnsheet.TurnEventCategorySystem,
					Icon:     turnsheet.TurnEventIconSystem,
					Message: fmt.Sprintf("%s entering depot for structure repairs (%d SP).",
						mechInst.Callsign, cost),
				})
			}
		}

		// Weapon swaps
		if len(order.WeaponSwaps) > 0 {
			var weaponConfig []mecha_record.WeaponConfigEntry
			if len(mechInst.WeaponConfigJSON) > 0 {
				if err := json.Unmarshal(mechInst.WeaponConfigJSON, &weaponConfig); err != nil {
					l.Warn("failed to unmarshal weapon config for mech >%s<: %v", mechInst.Callsign, err)
				}
			}

			for _, swap := range order.WeaponSwaps {
				if swap.SlotLocation == "" || swap.NewWeaponID == "" {
					continue
				}
				// Validate new weapon exists
				newWeapon, err := p.Domain.GetMechaWeaponRec(swap.NewWeaponID, nil)
				if err != nil {
					l.Warn("weapon >%s< not found for swap on mech >%s<: %v",
						swap.NewWeaponID, mechInst.Callsign, err)
					continue
				}

				// Apply the swap (add or replace the slot entry)
				slotFound := false
				for i, entry := range weaponConfig {
					if entry.SlotLocation == swap.SlotLocation {
						weaponConfig[i].WeaponID = swap.NewWeaponID
						slotFound = true
						break
					}
				}
				if !slotFound {
					weaponConfig = append(weaponConfig, mecha_record.WeaponConfigEntry{
						WeaponID:     swap.NewWeaponID,
						SlotLocation: swap.SlotLocation,
					})
				}

				spCost++
				mechInst.IsRefitting = true
				mechChanged = true
				events = append(events, turnsheet.TurnEvent{
					Category: turnsheet.TurnEventCategorySystem,
					Icon:     turnsheet.TurnEventIconSystem,
					Message: fmt.Sprintf("%s installing %s in %s slot (1 SP).",
						mechInst.Callsign, newWeapon.Name, swap.SlotLocation),
				})
			}

			if mechChanged {
				mechInst.WeaponConfig = weaponConfig
				if jsonBytes, err := json.Marshal(weaponConfig); err == nil {
					mechInst.WeaponConfigJSON = jsonBytes
				}
			}
		}

		if mechChanged {
			if _, err := p.Domain.UpdateMechaMechInstanceRec(mechInst); err != nil {
				l.Warn("failed to update mech instance >%s< after management: %v", mechInst.Callsign, err)
			}
		}
	}

	// Deduct supply points
	if spCost > 0 {
		freshLance.SupplyPoints -= spCost
		if freshLance.SupplyPoints < 0 {
			freshLance.SupplyPoints = 0
		}
		events = append(events, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategorySystem,
			Icon:     turnsheet.TurnEventIconSystem,
			Message:  fmt.Sprintf("Spent %d supply points on management orders.", spCost),
		})
	}

	// Persist events
	for _, evt := range events {
		if err := turnsheet.AppendMechaTurnEvent(freshLance, evt); err != nil {
			l.Warn("failed to append management event for lance >%s<: %v", lanceInstance.ID, err)
		}
	}

	if _, err := p.Domain.UpdateMechaLanceInstanceRec(freshLance); err != nil {
		l.Warn("failed to update lance instance after management: %v", err)
	}

	return nil
}

// mechIsAtDepot returns true if the mech's current sector instance is a starting
// (depot-capable) sector.
func (p *MechaLanceManagementProcessor) mechIsAtDepot(
	mechInst *mecha_record.MechaMechInstance,
) (bool, error) {
	sectorInst, err := p.Domain.GetMechaSectorInstanceRec(mechInst.MechaSectorInstanceID, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get sector instance: %w", err)
	}
	sectorDesign, err := p.Domain.GetMechaSectorRec(sectorInst.MechaSectorID, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get sector design: %w", err)
	}
	return sectorDesign.IsStartingSector, nil
}

// CreateNextTurnSheet creates a management sheet for the next turn if any mech
// in this lance is at a depot sector. Returns nil (no sheet) otherwise.
// AI-controlled lances (no AccountUserID) are skipped.
func (p *MechaLanceManagementProcessor) CreateNextTurnSheet(
	ctx context.Context,
	gameInstanceRec *game_record.GameInstance,
	lanceInstance *mecha_record.MechaLanceInstance,
) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("MechaLanceManagementProcessor/CreateNextTurnSheet")

	lanceRec, err := p.Domain.GetMechaLanceRec(lanceInstance.MechaLanceID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get lance: %w", err)
	}

	// Only create management sheets for player-owned lances
	if !lanceRec.AccountUserID.Valid || lanceRec.AccountUserID.String == "" {
		l.Info("lance >%s< has no account user — skipping management sheet (AI lance)", lanceInstance.ID)
		return nil, nil
	}

	mechInstances, err := p.Domain.GetManyMechaMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaMechInstanceMechaLanceInstanceID, Val: lanceInstance.ID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get mech instances: %w", err)
	}

	// Check whether any mech is at a depot sector
	anyAtDepot := false
	for _, mech := range mechInstances {
		atDepot, err := p.mechIsAtDepot(mech)
		if err != nil {
			l.Warn("failed to check depot for mech >%s<: %v", mech.ID, err)
			continue
		}
		if atDepot {
			anyAtDepot = true
			break
		}
	}

	if !anyAtDepot {
		l.Info("no mechs at depot for lance >%s< — skipping management sheet", lanceInstance.ID)
		return nil, nil
	}

	return p.buildManagementSheet(ctx, l, gameInstanceRec, lanceRec, lanceInstance, mechInstances)
}

func (p *MechaLanceManagementProcessor) buildManagementSheet(
	_ context.Context,
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	lanceRec *mecha_record.MechaLance,
	lanceInstance *mecha_record.MechaLanceInstance,
	mechInstances []*mecha_record.MechaMechInstance,
) (*game_record.GameTurnSheet, error) {
	gameRec, err := p.Domain.GetGameRec(gameInstanceRec.GameID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	accountUserRec, err := p.Domain.GetAccountUserRec(lanceRec.AccountUserID.String, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account user: %w", err)
	}

	// Load weapon catalog
	allWeapons, err := p.Domain.GetManyMechaWeaponRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaWeaponGameID, Val: gameRec.ID},
		},
	})
	if err != nil {
		l.Warn("failed to load weapon catalog: %v", err)
	}

	var catalog []turnsheet.CatalogWeapon
	for _, w := range allWeapons {
		catalog = append(catalog, turnsheet.CatalogWeapon{
			WeaponID:  w.ID,
			Name:      w.Name,
			Damage:    w.Damage,
			HeatCost:  w.HeatCost,
			RangeBand: w.RangeBand,
			MountSize: w.MountSize,
		})
	}

	// Build per-mech management entries
	var mechEntries []turnsheet.ManagementMechEntry
	for _, mech := range mechInstances {
		entry := turnsheet.ManagementMechEntry{
			MechInstanceID:   mech.ID,
			Callsign:         mech.Callsign,
			Status:           mech.Status,
			IsRefitting:      mech.IsRefitting,
			CurrentArmor:     mech.CurrentArmor,
			CurrentStructure: mech.CurrentStructure,
		}

		chassisRec, err := p.Domain.GetMechaChassisRec(mech.MechaChassisID, nil)
		if err == nil {
			entry.ChassisName = chassisRec.Name
			entry.ChassisClass = chassisRec.ChassisClass
			entry.MaxArmor = chassisRec.ArmorPoints
			entry.MaxStructure = chassisRec.StructurePoints
			entry.StructureDamage = chassisRec.StructurePoints - mech.CurrentStructure
		}

		isAtDepot, _ := p.mechIsAtDepot(mech)
		entry.IsAtDepot = isAtDepot

		// Build weapon slots from weapon config
		var weaponConfig []mecha_record.WeaponConfigEntry
		if len(mech.WeaponConfigJSON) > 0 {
			_ = json.Unmarshal(mech.WeaponConfigJSON, &weaponConfig)
		}
		// Build a weapon-ID-to-name lookup from catalog
		weaponNames := make(map[string]string, len(allWeapons))
		for _, w := range allWeapons {
			weaponNames[w.ID] = w.Name
		}
		for _, slot := range weaponConfig {
			entry.Weapons = append(entry.Weapons, turnsheet.MechWeaponSlot{
				SlotLocation:      slot.SlotLocation,
				CurrentWeaponID:   slot.WeaponID,
				CurrentWeaponName: weaponNames[slot.WeaponID],
			})
		}

		mechEntries = append(mechEntries, entry)
	}

	turnNumber := gameInstanceRec.CurrentTurn
	title := "Lance Management"
	instructions := turnsheet.DefaultManagementInstructions()

	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	var backgroundImage *string
	bgImageURL, err := p.Domain.GetGameTurnSheetImageDataURL(gameRec.ID, mecha_record.MechaTurnSheetTypeLanceManagement)
	if err == nil && bgImageURL != "" {
		backgroundImage = &bgImageURL
	}

	sheetData := turnsheet.LanceManagementData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:              &gameRec.Name,
			GameType:              &gameRec.GameType,
			TurnNumber:            &turnNumber,
			AccountName:           &accountUserRec.Email,
			TurnSheetTitle:        &title,
			TurnSheetDescription:  &gameRec.Description,
			TurnSheetInstructions: &instructions,
			TurnSheetCode:         &turnSheetCode,
			BackgroundImage:       backgroundImage,
		},
		LanceName:     lanceRec.Name,
		SupplyPoints:  lanceInstance.SupplyPoints,
		Mechs:         mechEntries,
		WeaponCatalog: catalog,
	}

	sheetJSON, err := json.Marshal(sheetData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal management sheet data: %w", err)
	}

	sheetOrder := mecha_record.MechaSheetPresentationOrderForType(mecha_record.MechaTurnSheetTypeLanceManagement)

	gameTurnSheet := &game_record.GameTurnSheet{
		GameID:           gameRec.ID,
		AccountID:        accountUserRec.AccountID,
		AccountUserID:    lanceRec.AccountUserID.String,
		TurnNumber:       turnNumber,
		SheetType:        mecha_record.MechaTurnSheetTypeLanceManagement,
		SheetOrder:       sheetOrder,
		SheetData:        json.RawMessage(sheetJSON),
		IsCompleted:      false,
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
	gameTurnSheet.GameInstanceID = nullstring.FromString(gameInstanceRec.ID)

	sheetRec, err := p.Domain.CreateGameTurnSheetRec(gameTurnSheet)
	if err != nil {
		return nil, fmt.Errorf("failed to create management game turn sheet: %w", err)
	}

	// Link the game_turn_sheet to the lance instance via mecha_turn_sheet
	if _, err := p.Domain.CreateMechaTurnSheetRec(&mecha_record.MechaTurnSheet{
		GameID:               gameRec.ID,
		MechaLanceInstanceID: lanceInstance.ID,
		GameTurnSheetID:      sheetRec.ID,
	}); err != nil {
		return nil, fmt.Errorf("failed to create mecha turn sheet record: %w", err)
	}

	l.Info("created management turn sheet >%s< for lance >%s< turn >%d<",
		sheetRec.ID, lanceInstance.ID, turnNumber)
	return sheetRec, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
