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
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

// mechOrderEntryFromInstance builds a MechOrderEntry from a mech instance,
// resolving chassis stats, sector name, and weapon loadout. When the mech
// has equipment installed and is not refitting, Speed and MaxArmor use the
// effective (bonus-applied) values so the turn sheet shows the actual
// reachable distance and armor ceiling the engine will enforce.
func mechOrderEntryFromInstance(
	l logger.Logger,
	d *domain.Domain,
	mechInst *mecha_game_record.MechaGameMechInstance,
) turnsheet.MechOrderEntry {
	entry := turnsheet.MechOrderEntry{
		MechInstanceID:   mechInst.ID,
		MechCallsign:     mechInst.Callsign,
		MechStatus:       mechInst.Status,
		CurrentArmor:     mechInst.CurrentArmor,
		CurrentStructure: mechInst.CurrentStructure,
		CurrentHeat:      mechInst.CurrentHeat,
		PilotSkill:       mechInst.PilotSkill,
		IsRefitting:      mechInst.IsRefitting,
	}

	// Resolve equipment effects so Speed / MaxArmor below reflect the
	// chassis *plus* applicable bonuses (zeroed out for refitting mechs).
	var equipmentEntries []mecha_game_record.EquipmentConfigEntry
	if len(mechInst.EquipmentConfigJSON) > 0 {
		if err := json.Unmarshal(mechInst.EquipmentConfigJSON, &equipmentEntries); err != nil {
			l.Warn("failed to unmarshal equipment config for mech >%s< >%v<", mechInst.ID, err)
		}
	}
	equipmentByID, err := d.LoadMechaGameEquipmentByID(equipmentEntries)
	if err != nil {
		l.Warn("failed to load equipment for mech >%s< >%v<", mechInst.ID, err)
	}
	effects := domain.AggregateMechaGameEquipmentEffects(equipmentEntries, equipmentByID, mechInst.IsRefitting)

	// Resolve chassis stats
	chassisRec, err := d.GetMechaGameChassisRec(mechInst.MechaGameChassisID, nil)
	if err != nil {
		l.Warn("failed to get chassis >%s< for mech >%s< >%v<",
			mechInst.MechaGameChassisID, mechInst.ID, err)
	} else {
		entry.ChassisName = chassisRec.Name
		entry.ChassisClass = chassisRec.ChassisClass
		entry.MaxArmor = domain.EffectiveMechaGameMaxArmor(chassisRec, effects)
		entry.MaxStructure = chassisRec.StructurePoints
		entry.HeatCapacity = chassisRec.HeatCapacity
		entry.Speed = domain.EffectiveMechaGameSpeed(chassisRec, effects)
	}

	// Resolve sector name
	sectorInstRec, err := d.GetMechaGameSectorInstanceRec(mechInst.MechaGameSectorInstanceID, nil)
	if err != nil {
		l.Warn("failed to get sector instance >%s< >%v<",
			mechInst.MechaGameSectorInstanceID, err)
	} else {
		sectorRec, err := d.GetMechaGameSectorRec(sectorInstRec.MechaGameSectorID, nil)
		if err != nil {
			l.Warn("failed to get sector >%s< >%v<", sectorInstRec.MechaGameSectorID, err)
		} else {
			entry.CurrentSectorName = sectorRec.Name
		}
	}

	// Resolve weapon entries from instance weapon config. We also build a
	// lookup map so ammo capacity can be summed alongside ammo_bin
	// magnitudes further down.
	var weaponConfig []mecha_game_record.WeaponConfigEntry
	if len(mechInst.WeaponConfigJSON) > 0 {
		if err := json.Unmarshal(mechInst.WeaponConfigJSON, &weaponConfig); err != nil {
			l.Warn("failed to unmarshal weapon config for mech >%s< >%v<", mechInst.ID, err)
		}
	}
	weaponByID := make(map[string]*mecha_game_record.MechaGameWeapon, len(weaponConfig))
	hasAmmoWeapon := false
	for _, slot := range weaponConfig {
		if slot.WeaponID == "" {
			continue
		}
		weaponRec, err := d.GetMechaGameWeaponRec(slot.WeaponID, nil)
		if err != nil {
			l.Warn("failed to get weapon >%s< >%v<", slot.WeaponID, err)
			continue
		}
		weaponByID[weaponRec.ID] = weaponRec
		if weaponRec.AmmoCapacity > 0 {
			hasAmmoWeapon = true
		}
		entry.Weapons = append(entry.Weapons, turnsheet.MechWeaponEntry{
			WeaponID:     weaponRec.ID,
			Name:         weaponRec.Name,
			Damage:       weaponRec.Damage,
			HeatCost:     weaponRec.HeatCost,
			RangeBand:    weaponRec.RangeBand,
			SlotLocation: slot.SlotLocation,
			AmmoCapacity: weaponRec.AmmoCapacity,
		})
	}

	// Equipment entries — mirror the weapon rendering so players can see
	// both mounts on the same chassis diagram. We already resolved
	// equipmentEntries / equipmentByID above when computing effects.
	for _, entryCfg := range equipmentEntries {
		if entryCfg.EquipmentID == "" {
			continue
		}
		eq := equipmentByID[entryCfg.EquipmentID]
		if eq == nil {
			continue
		}
		entry.Equipment = append(entry.Equipment, turnsheet.MechEquipmentEntry{
			EquipmentID:  eq.ID,
			Name:         eq.Name,
			EffectKind:   eq.EffectKind,
			Magnitude:    eq.Magnitude,
			HeatCost:     eq.HeatCost,
			MountSize:    eq.MountSize,
			SlotLocation: entryCfg.SlotLocation,
		})
	}

	// Ammo summary: populated only for mechs carrying at least one
	// ammo-consuming weapon, otherwise the "x/y" display just clutters
	// the sheet.
	if hasAmmoWeapon {
		entry.AmmoRemaining = mechInst.AmmoRemaining
		entry.AmmoCapacity = domain.MaxMechaGameAmmoCapacity(
			weaponConfig, weaponByID, equipmentEntries, equipmentByID,
		)
	}

	return entry
}

// MechaGameOrdersProcessor implements the TurnSheetProcessor interface
// (defined in the parent mecha package)

// MechaGameOrdersProcessor processes orders turn sheet business logic for mecha
type MechaGameOrdersProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewMechaGameOrdersProcessor creates a new mecha orders processor.
func NewMechaGameOrdersProcessor(l logger.Logger, d *domain.Domain) (*MechaGameOrdersProcessor, error) {
	l = l.WithFunctionContext("NewMechaGameOrdersProcessor")

	p := &MechaGameOrdersProcessor{
		Logger: l,
		Domain: d,
	}
	return p, nil
}

// GetSheetType returns the sheet type this processor handles (implements TurnSheetProcessor interface).
func (p *MechaGameOrdersProcessor) GetSheetType() string {
	return mecha_game_record.MechaGameTurnSheetTypeOrders
}

// ProcessTurnSheetResponse processes a single orders turn sheet response (implements TurnSheetProcessor interface).
func (p *MechaGameOrdersProcessor) ProcessTurnSheetResponse(ctx context.Context, gameInstanceRec *game_record.GameInstance, squadInstance *mecha_game_record.MechaGameSquadInstance, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("MechaGameOrdersProcessor/ProcessTurnSheetResponse")

	l.Info("processing orders for turn sheet >%s< for squad instance >%s<", turnSheet.ID, squadInstance.ID)

	if turnSheet.SheetType != mecha_game_record.MechaGameTurnSheetTypeOrders {
		return fmt.Errorf("invalid sheet type: expected %s, got %s", mecha_game_record.MechaGameTurnSheetTypeOrders, turnSheet.SheetType)
	}

	var scanData turnsheet.OrdersScanData
	if err := json.Unmarshal(turnSheet.ScannedData, &scanData); err != nil {
		l.Warn("failed to unmarshal scanned data >%v<", err)
		return fmt.Errorf("failed to parse scanned data: %w", err)
	}

	if len(scanData.MechOrders) == 0 {
		l.Info("no mech orders in scanned data — squad stays in place")
		return nil
	}

	// Reload squadInstance to get latest state for event appending
	freshSquad, err := p.Domain.GetMechaGameSquadInstanceRec(squadInstance.ID, nil)
	if err != nil {
		l.Warn("failed to reload squad instance for events >%v<", err)
		freshSquad = squadInstance
	}

	for _, order := range scanData.MechOrders {
		if moveEvent := p.processMechOrderWithEvent(l, gameInstanceRec, order); moveEvent != nil {
			if err := turnsheet.AppendMechaGameTurnEvent(freshSquad, *moveEvent); err != nil {
				l.Warn("failed to append movement event for mech >%s<: %v", order.MechInstanceID, err)
			}
		}
	}

	if _, err := p.Domain.UpdateMechaGameSquadInstanceRec(freshSquad); err != nil {
		l.Warn("failed to persist movement events for squad >%s<: %v", squadInstance.ID, err)
	}

	return nil
}

// processMechOrderWithEvent applies movement and returns a TurnEvent if movement
// occurred, or nil if no event should be generated.
func (p *MechaGameOrdersProcessor) processMechOrderWithEvent(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	order turnsheet.ScannedMechOrder,
) *turnsheet.TurnEvent {
	if order.MechInstanceID == "" {
		return nil
	}

	if order.MoveToSectorInstanceID == "" {
		l.Debug("no movement order for mech >%s< — staying in place", order.MechInstanceID)
		return nil
	}

	mechInstanceRec, err := p.Domain.GetMechaGameMechInstanceRec(order.MechInstanceID, nil)
	if err != nil {
		l.Warn("failed to get mech instance >%s< >%v<", order.MechInstanceID, err)
		return nil
	}

	if mechInstanceRec.GameInstanceID != gameInstanceRec.ID {
		l.Warn("mech instance >%s< does not belong to game instance >%s<", order.MechInstanceID, gameInstanceRec.ID)
		return nil
	}

	if mechInstanceRec.Status == mecha_game_record.MechInstanceStatusDestroyed {
		l.Info("mech >%s< is destroyed — ignoring movement order", order.MechInstanceID)
		return nil
	}

	if mechInstanceRec.IsRefitting {
		l.Info("mech >%s< is refitting — ignoring movement order", order.MechInstanceID)
		return nil
	}

	sectorInstanceRec, err := p.Domain.GetMechaGameSectorInstanceRec(order.MoveToSectorInstanceID, nil)
	if err != nil {
		l.Warn("failed to get sector instance >%s< >%v<", order.MoveToSectorInstanceID, err)
		return nil
	}

	if sectorInstanceRec.GameInstanceID != gameInstanceRec.ID {
		l.Warn("sector instance >%s< does not belong to game instance >%s<", order.MoveToSectorInstanceID, gameInstanceRec.ID)
		return nil
	}

	// Validate that the destination is within the mech's *effective* speed
	// budget (chassis.Speed + jump-jets SpeedBonus, zeroed while refitting).
	chassisRec, err := p.Domain.GetMechaGameChassisRec(mechInstanceRec.MechaGameChassisID, nil)
	if err != nil {
		l.Warn("failed to get chassis >%s< for movement validation: %v", mechInstanceRec.MechaGameChassisID, err)
		return nil
	}

	var equipmentEntries []mecha_game_record.EquipmentConfigEntry
	if len(mechInstanceRec.EquipmentConfigJSON) > 0 {
		if err := json.Unmarshal(mechInstanceRec.EquipmentConfigJSON, &equipmentEntries); err != nil {
			l.Warn("failed to unmarshal equipment config for mech >%s< >%v<", mechInstanceRec.ID, err)
		}
	}
	equipmentByID, err := p.Domain.LoadMechaGameEquipmentByID(equipmentEntries)
	if err != nil {
		l.Warn("failed to load equipment for mech >%s< >%v<", mechInstanceRec.ID, err)
	}
	effects := domain.AggregateMechaGameEquipmentEffects(equipmentEntries, equipmentByID, mechInstanceRec.IsRefitting)
	effectiveSpeed := domain.EffectiveMechaGameSpeed(chassisRec, effects)

	hops, reachable := p.IsSectorReachableWithinSpeed(l, gameInstanceRec.ID, mechInstanceRec.MechaGameSectorInstanceID, order.MoveToSectorInstanceID, effectiveSpeed)
	if !reachable {
		l.Warn("mech >%s< cannot reach sector >%s< within speed budget %d (distance > %d hops)",
			order.MechInstanceID, order.MoveToSectorInstanceID, effectiveSpeed, effectiveSpeed)
		return nil
	}

	sectorRec, err := p.Domain.GetMechaGameSectorRec(sectorInstanceRec.MechaGameSectorID, nil)
	if err != nil {
		l.Warn("failed to get sector design >%s<: %v", sectorInstanceRec.MechaGameSectorID, err)
	}

	// Jump-jets heat predicate: if the mech traveled more hops than its
	// base chassis speed, any jump_jets equipment with heat_cost > 0 fires
	// this turn. Apply the heat immediately to the mech's running heat so
	// combat resolution and end-of-turn both see the updated value.
	if hops > chassisRec.Speed {
		jumpHeat := domain.MechaGameEquipmentJumpJetHeatCost(
			equipmentEntries, equipmentByID, mechInstanceRec.IsRefitting,
		)
		if jumpHeat > 0 {
			mechInstanceRec.CurrentHeat += jumpHeat
		}
	}

	mechInstanceRec.MechaGameSectorInstanceID = order.MoveToSectorInstanceID
	if _, err := p.Domain.UpdateMechaGameMechInstanceRec(mechInstanceRec); err != nil {
		l.Warn("failed to update mech instance >%s< >%v<", order.MechInstanceID, err)
		return nil
	}

	l.Info("moved mech >%s< to sector instance >%s<", order.MechInstanceID, order.MoveToSectorInstanceID)

	sectorName := order.MoveToSectorInstanceID
	if sectorRec != nil {
		sectorName = sectorRec.Name
	}
	evt := turnsheet.TurnEvent{
		Category: turnsheet.TurnEventCategoryMovement,
		Icon:     turnsheet.TurnEventIconMovement,
		Message:  fmt.Sprintf("%s moved to %s.", mechInstanceRec.Callsign, sectorName),
	}
	return &evt
}

// AttackDeclaration represents a declared attack from one mech to another.
type AttackDeclaration struct {
	AttackerMechInstanceID string
	TargetMechInstanceID   string
}

// ExtractAttackDeclarations reads scanned data from an orders turn sheet and
// returns all declared attack orders. Called after ProcessTurnSheetResponse
// to collect attacks for combat resolution.
func (p *MechaGameOrdersProcessor) ExtractAttackDeclarations(
	turnSheet *game_record.GameTurnSheet,
) ([]AttackDeclaration, error) {
	if len(turnSheet.ScannedData) == 0 {
		return nil, nil
	}

	var scanData turnsheet.OrdersScanData
	if err := json.Unmarshal(turnSheet.ScannedData, &scanData); err != nil {
		return nil, fmt.Errorf("failed to parse scanned data: %w", err)
	}

	var attacks []AttackDeclaration
	for _, order := range scanData.MechOrders {
		if order.MechInstanceID != "" && order.AttackTargetMechInstanceID != "" {
			attacks = append(attacks, AttackDeclaration{
				AttackerMechInstanceID: order.MechInstanceID,
				TargetMechInstanceID:   order.AttackTargetMechInstanceID,
			})
		}
	}

	return attacks, nil
}

// CreateNextTurnSheet creates a new orders turn sheet for a squad instance (implements TurnSheetProcessor interface).
func (p *MechaGameOrdersProcessor) CreateNextTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, squadInstance *mecha_game_record.MechaGameSquadInstance) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("MechaGameOrdersProcessor/CreateNextTurnSheet")

	l.Info("creating orders turn sheet for squad instance >%s<", squadInstance.ID)

	// Step 1: Get the squad design record
	squadRec, err := p.Domain.GetMechaGameSquadRec(squadInstance.MechaGameSquadID, nil)
	if err != nil {
		l.Warn("failed to get squad >%v<", err)
		return nil, fmt.Errorf("failed to get squad: %w", err)
	}

	// Step 2: Get the account user for the squad owner via subscription chain
	accountUserRec, err := squadInstanceAccountUser(p.Domain, squadInstance)
	if err != nil {
		l.Warn("failed to get account user for squad instance >%s< >%v<", squadInstance.ID, err)
		return nil, fmt.Errorf("failed to get account user: %w", err)
	}

	// Step 3: Get the game record
	gameRec, err := p.Domain.GetGameRec(gameInstanceRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game >%v<", err)
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Step 4: Load background image
	var backgroundImage *string
	bgImageURL, err := p.Domain.GetGameTurnSheetImageDataURL(gameRec.ID, mecha_game_record.MechaGameTurnSheetTypeOrders)
	if err != nil {
		l.Warn("failed to get turn sheet background image >%v<", err)
	} else if bgImageURL != "" {
		backgroundImage = &bgImageURL
		l.Info("loaded background image for mecha orders turn sheet, length >%d<", len(bgImageURL))
	}

	// Step 5: Get all mech instances for this squad instance
	mechInstances, err := p.Domain.GetManyMechaGameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameMechInstanceMechaGameSquadInstanceID, Val: squadInstance.ID},
		},
	})
	if err != nil {
		l.Warn("failed to get mech instances >%v<", err)
		return nil, fmt.Errorf("failed to get mech instances: %w", err)
	}

	// Step 6: Build mech order entries with full stats
	var squadMechs []turnsheet.MechOrderEntry
	var sectorInstanceIDs []string
	sectorInstanceIDSet := make(map[string]bool)

	for _, mechInst := range mechInstances {
		entry := mechOrderEntryFromInstance(l, p.Domain, mechInst)
		squadMechs = append(squadMechs, entry)

		if !sectorInstanceIDSet[mechInst.MechaGameSectorInstanceID] {
			sectorInstanceIDSet[mechInst.MechaGameSectorInstanceID] = true
			sectorInstanceIDs = append(sectorInstanceIDs, mechInst.MechaGameSectorInstanceID)
		}
	}

	// Step 7: Compute per-mech reachable sectors within each mech's speed budget.
	// We also build a union for the legacy AvailableSectors field.
	availableSectorsSeen := make(map[string]bool)
	var availableSectors []turnsheet.SectorOption
	for i := range squadMechs {
		mech := &squadMechs[i]
		// Find the matching mech instance to get the current sector and chassis speed.
		var mechInst *mecha_game_record.MechaGameMechInstance
		for _, mi := range mechInstances {
			if mi.ID == mech.MechInstanceID {
				mechInst = mi
				break
			}
		}
		if mechInst == nil {
			continue
		}
		reachable, err := p.getReachableSectorOptions(l, gameInstanceRec.ID, mechInst.MechaGameSectorInstanceID, mech.Speed)
		if err != nil {
			l.Warn("failed to get reachable sectors for mech >%s< >%v<", mech.MechInstanceID, err)
		}
		mech.ReachableSectors = reachable
		for _, opt := range reachable {
			if !availableSectorsSeen[opt.SectorInstanceID] {
				availableSectorsSeen[opt.SectorInstanceID] = true
				availableSectors = append(availableSectors, opt)
			}
		}
	}
	_ = sectorInstanceIDs

	// Step 8: Get enemy mech instances visible to this squad
	enemyMechs, err := p.getEnemyMechOptions(l, gameInstanceRec, squadInstance)
	if err != nil {
		l.Warn("failed to get enemy mechs >%v<", err)
		// Non-fatal: continue with no attack options
	}

	// Step 9: Generate turn sheet code
	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		l.Warn("failed to generate turn sheet code >%v<", err)
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	turnNumber := gameInstanceRec.CurrentTurn
	title := "Mech Orders"
	instructions := turnsheet.DefaultOrdersInstructions()

	// Read and clear turn events accumulated during previous turn processing
	turnEvents, err := turnsheet.ReadAndClearMechaGameTurnEvents(squadInstance)
	if err != nil {
		l.Warn("failed to read turn events for squad >%s< >%v<", squadInstance.ID, err)
		turnEvents = nil
	} else if len(turnEvents) > 0 {
		if _, err := p.Domain.UpdateMechaGameSquadInstanceRec(squadInstance); err != nil {
			l.Warn("failed to persist cleared turn events for squad >%s< >%v<", squadInstance.ID, err)
		}
	}

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
			BackgroundImage:       backgroundImage,
			TurnEvents:            turnEvents,
		},
		SquadName:        squadRec.Name,
		SquadMechs:       squadMechs,
		AvailableSectors: availableSectors,
		EnemyMechs:       enemyMechs,
	}

	sheetDataBytes, err := json.Marshal(sheetData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sheet data: %w", err)
	}

	// Step 10: Create the game_turn_sheet record
	gameTurnSheet := &game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		AccountID:        accountUserRec.AccountID,
		AccountUserID:    accountUserRec.ID,
		TurnNumber:       gameInstanceRec.CurrentTurn,
		SheetType:        mecha_game_record.MechaGameTurnSheetTypeOrders,
		SheetOrder:       mecha_game_record.MechaGameSheetOrderForType(mecha_game_record.MechaGameTurnSheetTypeOrders),
		SheetData:        json.RawMessage(sheetDataBytes),
		IsCompleted:      false,
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
	gameTurnSheet.GameInstanceID = nullstring.FromString(gameInstanceRec.ID)

	createdTurnSheetRec, err := p.Domain.CreateGameTurnSheetRec(gameTurnSheet)
	if err != nil {
		return nil, fmt.Errorf("failed to create turn sheet record: %w", err)
	}

	// Step 11: Link the game_turn_sheet to the squad instance via mecha_game_turn_sheet
	mechaGameTurnSheet := &mecha_game_record.MechaGameTurnSheet{
		GameID:               gameInstanceRec.GameID,
		MechaGameSquadInstanceID: squadInstance.ID,
		GameTurnSheetID:      createdTurnSheetRec.ID,
	}

	_, err = p.Domain.CreateMechaGameTurnSheetRec(mechaGameTurnSheet)
	if err != nil {
		return nil, fmt.Errorf("failed to create mecha turn sheet record: %w", err)
	}

	l.Info("created orders turn sheet >%s< for squad instance >%s< turn >%d<", createdTurnSheetRec.ID, squadInstance.ID, gameInstanceRec.CurrentTurn)
	return createdTurnSheetRec, nil
}

// getReachableSectorOptions returns all sector instances reachable from the given
// starting sector within the given speed (number of hops), using BFS over sector links.
func (p *MechaGameOrdersProcessor) getReachableSectorOptions(l logger.Logger, gameInstanceID string, startSectorInstanceID string, speed int) ([]turnsheet.SectorOption, error) {
	if speed <= 0 {
		return nil, nil
	}

	type bfsNode struct {
		sectorInstanceID string
		depth            int
	}

	seen := map[string]bool{startSectorInstanceID: true}
	queue := []bfsNode{{sectorInstanceID: startSectorInstanceID, depth: 0}}
	var options []turnsheet.SectorOption

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.depth >= speed {
			continue
		}

		sectorInstRec, err := p.Domain.GetMechaGameSectorInstanceRec(cur.sectorInstanceID, nil)
		if err != nil {
			l.Warn("failed to get sector instance >%s< >%v<", cur.sectorInstanceID, err)
			continue
		}

		links, err := p.Domain.GetManyMechaGameSectorLinkRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_game_record.FieldMechaGameSectorLinkFromMechaGameSectorID, Val: sectorInstRec.MechaGameSectorID},
			},
		})
		if err != nil {
			l.Warn("failed to get sector links for sector >%s< >%v<", sectorInstRec.MechaGameSectorID, err)
			continue
		}

		for _, link := range links {
			linkedInstances, err := p.Domain.GetManyMechaGameSectorInstanceRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: mecha_game_record.FieldMechaGameSectorInstanceGameInstanceID, Val: gameInstanceID},
					{Col: mecha_game_record.FieldMechaGameSectorInstanceMechaGameSectorID, Val: link.ToMechaGameSectorID},
				},
				Limit: 1,
			})
			if err != nil || len(linkedInstances) == 0 {
				continue
			}

			destID := linkedInstances[0].ID
			if seen[destID] {
				continue
			}
			seen[destID] = true

			sectorRec, err := p.Domain.GetMechaGameSectorRec(link.ToMechaGameSectorID, nil)
			if err != nil {
				l.Warn("failed to get sector design >%s< >%v<", link.ToMechaGameSectorID, err)
				continue
			}

			options = append(options, turnsheet.SectorOption{
				SectorInstanceID: destID,
				SectorName:       sectorRec.Name,
			})
			queue = append(queue, bfsNode{sectorInstanceID: destID, depth: cur.depth + 1})
		}
	}

	return options, nil
}

// IsSectorReachableWithinSpeed returns (hops, true) if destID is reachable from fromID
// within the given number of hops using BFS, or (0, false) if not reachable.
func (p *MechaGameOrdersProcessor) IsSectorReachableWithinSpeed(l logger.Logger, gameInstanceID, fromSectorInstanceID, destSectorInstanceID string, speed int) (int, bool) {
	if fromSectorInstanceID == destSectorInstanceID {
		return 0, true
	}

	type bfsNode struct {
		sectorInstanceID string
		depth            int
	}

	seen := map[string]bool{fromSectorInstanceID: true}
	queue := []bfsNode{{sectorInstanceID: fromSectorInstanceID, depth: 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.depth >= speed {
			continue
		}

		sectorInstRec, err := p.Domain.GetMechaGameSectorInstanceRec(cur.sectorInstanceID, nil)
		if err != nil {
			l.Warn("failed to get sector instance >%s< for speed check: %v", cur.sectorInstanceID, err)
			continue
		}

		links, err := p.Domain.GetManyMechaGameSectorLinkRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_game_record.FieldMechaGameSectorLinkFromMechaGameSectorID, Val: sectorInstRec.MechaGameSectorID},
			},
		})
		if err != nil {
			continue
		}

		for _, link := range links {
			linkedInstances, err := p.Domain.GetManyMechaGameSectorInstanceRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: mecha_game_record.FieldMechaGameSectorInstanceGameInstanceID, Val: gameInstanceID},
					{Col: mecha_game_record.FieldMechaGameSectorInstanceMechaGameSectorID, Val: link.ToMechaGameSectorID},
				},
				Limit: 1,
			})
			if err != nil || len(linkedInstances) == 0 {
				continue
			}

			destID := linkedInstances[0].ID
			if destID == destSectorInstanceID {
				return cur.depth + 1, true
			}
			if !seen[destID] {
				seen[destID] = true
				queue = append(queue, bfsNode{sectorInstanceID: destID, depth: cur.depth + 1})
			}
		}
	}

	return 0, false
}

// getEnemyMechOptions collects all enemy mech instances visible to the given squad.
func (p *MechaGameOrdersProcessor) getEnemyMechOptions(_ logger.Logger, gameInstanceRec *game_record.GameInstance, squadInstance *mecha_game_record.MechaGameSquadInstance) ([]turnsheet.EnemyMechOption, error) {
	// Get all mech instances for this game instance
	allMechInstances, err := p.Domain.GetManyMechaGameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameMechInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get mech instances: %w", err)
	}

	var enemies []turnsheet.EnemyMechOption
	for _, mechInst := range allMechInstances {
		if mechInst.MechaGameSquadInstanceID == squadInstance.ID {
			continue
		}
		if mechInst.Status == mecha_game_record.MechInstanceStatusDestroyed {
			continue
		}

		sectorName := ""
		sectorInstRec, err := p.Domain.GetMechaGameSectorInstanceRec(mechInst.MechaGameSectorInstanceID, nil)
		if err == nil {
			sectorRec, err2 := p.Domain.GetMechaGameSectorRec(sectorInstRec.MechaGameSectorID, nil)
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
