package turn_sheet_processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

const characterStartingHealth = 50

// AdventureGameCreatureEncounterProcessor processes monster encounter turn sheet business logic.
type AdventureGameCreatureEncounterProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewAdventureGameCreatureEncounterProcessor creates a new monster encounter processor.
func NewAdventureGameCreatureEncounterProcessor(l logger.Logger, d *domain.Domain) (*AdventureGameCreatureEncounterProcessor, error) {
	l = l.WithFunctionContext("NewAdventureGameCreatureEncounterProcessor")
	return &AdventureGameCreatureEncounterProcessor{
		Logger: l,
		Domain: d,
	}, nil
}

// GetSheetType returns the sheet type this processor handles.
func (p *AdventureGameCreatureEncounterProcessor) GetSheetType() string {
	return adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter
}

// ProcessTurnSheetResponse resolves combat for a monster encounter turn sheet.
func (p *AdventureGameCreatureEncounterProcessor) ProcessTurnSheetResponse(
	ctx context.Context,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	turnSheet *game_record.GameTurnSheet,
) error {
	l := p.Logger.WithFunctionContext("AdventureGameCreatureEncounterProcessor/ProcessTurnSheetResponse")
	l.Info("processing monster encounter for turn sheet >%s< character >%s<", turnSheet.ID, characterInstanceRec.ID)

	if turnSheet.SheetType != adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter {
		return fmt.Errorf("invalid sheet type: expected %s, got %s",
			adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter, turnSheet.SheetType)
	}

	// Step 1: Check whether the player took inventory actions this turn.
	// If so, all combat is forfeited for this turn.
	tookInventoryActions, err := p.inventorySheetHadActions(ctx, gameInstanceRec, characterInstanceRec)
	if err != nil {
		l.Warn("failed to check inventory sheet actions >%v< — proceeding with combat", err)
	}
	if tookInventoryActions {
		l.Info("inventory actions taken this turn — forfeiting all combat actions for character >%s<", characterInstanceRec.ID)
		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryCombat,
			Icon:     turnsheet.TurnEventIconInventory,
			Message:  "You managed your inventory — your combat actions were forfeited.",
		})
		_, saveErr := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec)
		if saveErr != nil {
			l.Warn("failed to save character instance events >%v<", saveErr)
		}
		return nil
	}

	// Step 2: Parse scanned actions.
	var scanData turnsheet.MonsterEncounterScanData
	if err := json.Unmarshal(turnSheet.ScannedData, &scanData); err != nil {
		l.Warn("failed to unmarshal scanned data >%v<", err)
		return fmt.Errorf("failed to parse scanned data: %w", err)
	}

	combatActions := scanData.GetActions()
	if len(combatActions) == 0 {
		l.Info("no combat actions in scanned data — nothing to resolve")
		return nil
	}

	// Step 3: Determine equipped weapon damage and armor defense.
	weaponDamage, armorDefense, err := ResolveEquipmentStats(l, p.Domain, characterInstanceRec.ID)
	if err != nil {
		l.Warn("failed to resolve equipment stats >%v< — using unarmed defaults", err)
		weaponDamage = defaultUnarmedAttackDamage
		armorDefense = 0
	}

	// Step 4: Track which non-aggressive creatures have been provoked this encounter.
	provoked := make(map[string]bool)

	// Step 5: Execute each action in order.
	for i, action := range combatActions {
		l.Info("executing action %d: %s target: %s", i+1, action.ActionType, action.TargetCreatureInstanceID)

		switch action.ActionType {
		case turnsheet.CombatActionTypeDoNothing:
			// No effect.

		case turnsheet.CombatActionTypeAttack:
			if action.TargetCreatureInstanceID == "" {
				l.Warn("attack action %d missing target — skipping", i+1)
				continue
			}

			// Get creature instance.
			creatureInstance, err := p.Domain.GetAdventureGameCreatureInstanceRec(action.TargetCreatureInstanceID, nil)
			if err != nil {
				l.Warn("failed to get creature instance >%s< >%v< — skipping", action.TargetCreatureInstanceID, err)
				continue
			}
			if creatureInstance.Health <= 0 {
				l.Info("creature >%s< already dead — skipping action", action.TargetCreatureInstanceID)
				continue
			}

			// Get creature definition for stats and disposition.
			creatureDef, err := p.Domain.GetAdventureGameCreatureRec(creatureInstance.AdventureGameCreatureID, nil)
			if err != nil {
				l.Warn("failed to get creature definition >%v< — skipping", err)
				continue
			}

			// Mark non-aggressive creatures as provoked.
			if creatureDef.Disposition != adventure_game_record.AdventureGameCreatureDispositionAggressive {
				provoked[creatureInstance.ID] = true
			}

			// Calculate player attack damage.
			playerDamage := weaponDamage - creatureDef.Defense
			if playerDamage < 1 {
				playerDamage = 1
			}

			// Apply damage to creature.
			creatureInstance.Health -= playerDamage
			l.Info("player deals %d damage to creature >%s< (health now %d)", playerDamage, creatureDef.Name, creatureInstance.Health)

			_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
				Category: turnsheet.TurnEventCategoryCombat,
				Icon:     turnsheet.TurnEventIconCombat,
				Message:  fmt.Sprintf("You attacked the %s for %d damage.", creatureDef.Name, playerDamage),
			})

			if creatureInstance.Health <= 0 {
				creatureInstance.Health = 0
				creatureInstance.DiedAtTurn = sql.NullInt64{Int64: int64(gameInstanceRec.CurrentTurn), Valid: true}
				_, err = p.Domain.UpdateAdventureGameCreatureInstanceRec(creatureInstance)
				if err != nil {
					l.Warn("failed to update dead creature instance >%v<", err)
				}
				l.Info("creature >%s< has been killed", creatureDef.Name)

				_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
					Category: turnsheet.TurnEventCategoryCombat,
					Icon:     turnsheet.TurnEventIconDeath,
					Message:  fmt.Sprintf("The %s has been slain!", creatureDef.Name),
				})

				// Move creature's item instances to the location, collecting names for narrative.
				droppedItems, err := p.moveCreatureItemsToLocation(l, creatureInstance, characterInstanceRec.AdventureGameLocationInstanceID)
				if err != nil {
					l.Warn("failed to move creature items to location >%v<", err)
				}
				for _, itemName := range droppedItems {
					_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
						Category: turnsheet.TurnEventCategoryCombat,
						Icon:     turnsheet.TurnEventIconInventory,
						Message:  fmt.Sprintf("The %s drops a %s.", creatureDef.Name, itemName),
					})
				}
				// Creature is dead — no retaliation this action.
				continue
			}

			// Creature retaliates.
			// Aggressive creatures always retaliate; inquisitive/indifferent only if provoked.
			willRetaliate := creatureDef.Disposition == adventure_game_record.AdventureGameCreatureDispositionAggressive ||
				provoked[creatureInstance.ID]

			if willRetaliate {
				creatureDamage := creatureDef.AttackDamage - armorDefense
				if creatureDamage < 1 {
					creatureDamage = 1
				}
				characterInstanceRec.Health -= creatureDamage
				l.Info("creature >%s< retaliates for %d damage (character health now %d)",
					creatureDef.Name, creatureDamage, characterInstanceRec.Health)

				attackDesc := creatureDef.AttackDescription
				if attackDesc == "" {
					attackDesc = "attacks you"
				}
				_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
					Category: turnsheet.TurnEventCategoryCombat,
					Icon:     turnsheet.TurnEventIconCombat,
					Message:  fmt.Sprintf("The %s %s for %d damage.", creatureDef.Name, attackDesc, creatureDamage),
				})
			}

			// Save updated creature.
			_, err = p.Domain.UpdateAdventureGameCreatureInstanceRec(creatureInstance)
			if err != nil {
				l.Warn("failed to update creature instance >%v<", err)
			}
		}
	}

	// Step 6: Handle character death.
	if characterInstanceRec.Health <= 0 {
		l.Info("character >%s< has been killed — resetting to starting location", characterInstanceRec.ID)
		characterInstanceRec.Health = 0

		startingLocationName := "your starting location"
		if err := p.resetDeadCharacter(l, gameInstanceRec, characterInstanceRec, &startingLocationName); err != nil {
			l.Warn("failed to reset dead character >%v<", err)
		}

		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryCombat,
			Icon:     turnsheet.TurnEventIconDeath,
			Message:  fmt.Sprintf("You were overwhelmed. You awaken at %s, battered but breathing.", startingLocationName),
		})
	}

	// Step 7: Persist updated character health and events.
	_, err = p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec)
	if err != nil {
		return fmt.Errorf("failed to save character instance: %w", err)
	}

	l.Info("monster encounter processing complete for character >%s<", characterInstanceRec.ID)
	return nil
}

// CreateNextTurnSheet creates a monster encounter turn sheet when creatures are present at a location.
// Returns nil, nil when no alive creatures are at the current location (no sheet needed).
func (p *AdventureGameCreatureEncounterProcessor) CreateNextTurnSheet(
	ctx context.Context,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("AdventureGameCreatureEncounterProcessor/CreateNextTurnSheet")
	l.Info("creating monster encounter turn sheet for character >%s<", characterInstanceRec.ID)

	// Step 1: Find creature instances at current location (both alive and recently dead).
	creatureInstances, err := p.Domain.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceRec.ID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, Val: characterInstanceRec.AdventureGameLocationInstanceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get creature instances: %w", err)
	}

	// Separate alive from recently dead.
	alive := make([]*adventure_game_record.AdventureGameCreatureInstance, 0, len(creatureInstances))
	recentlyDead := make([]*adventure_game_record.AdventureGameCreatureInstance, 0)
	for _, ci := range creatureInstances {
		if ci.Health > 0 {
			alive = append(alive, ci)
			continue
		}
		// Include recently dead within body_decay_turns.
		if ci.DiedAtTurn.Valid {
			creatureDef, err := p.Domain.GetAdventureGameCreatureRec(ci.AdventureGameCreatureID, nil)
			if err == nil {
				turnsDeadFor := int64(gameInstanceRec.CurrentTurn) - ci.DiedAtTurn.Int64
				if turnsDeadFor <= int64(creatureDef.BodyDecayTurns) {
					recentlyDead = append(recentlyDead, ci)
				}
			}
		}
	}

	if len(alive) == 0 && len(recentlyDead) == 0 {
		l.Info("no alive or recently dead creatures at location — skipping monster encounter sheet")
		return nil, nil
	}

	isReadOnly := len(alive) == 0 // No alive creatures means dead-body read-only sheet

	// Step 2: Load character info.
	characterRec, err := p.Domain.GetAdventureGameCharacterRec(characterInstanceRec.AdventureGameCharacterID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	accountUserRec, err := p.Domain.GetAccountUserRec(characterRec.AccountUserID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account user: %w", err)
	}

	gameRec, err := p.Domain.GetGameRec(gameInstanceRec.GameID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Step 3: Resolve equipped weapon and armor.
	_, _, err = ResolveEquipmentStats(l, p.Domain, characterInstanceRec.ID)
	if err != nil {
		l.Warn("failed to resolve equipment stats >%v<", err)
	}

	equippedWeapon, equippedArmor, err := p.resolveEquippedGear(l, characterInstanceRec.ID)
	if err != nil {
		l.Warn("failed to resolve equipped gear details >%v<", err)
	}

	characterAttack := defaultUnarmedAttackDamage
	if equippedWeapon != nil {
		characterAttack = equippedWeapon.Damage
	}
	characterDefense := 0
	if equippedArmor != nil {
		characterDefense = equippedArmor.Defense
	}

	// Step 4: Build encounter creature list (alive first, then recently dead).
	allCreaturesForSheet := make([]*adventure_game_record.AdventureGameCreatureInstance, 0, len(alive)+len(recentlyDead))
	allCreaturesForSheet = append(allCreaturesForSheet, alive...)
	allCreaturesForSheet = append(allCreaturesForSheet, recentlyDead...)

	creatures := make([]turnsheet.EncounterCreature, 0, len(allCreaturesForSheet))
	for _, ci := range allCreaturesForSheet {
		creatureDef, err := p.Domain.GetAdventureGameCreatureRec(ci.AdventureGameCreatureID, nil)
		if err != nil {
			l.Warn("failed to get creature definition >%s< >%v<", ci.AdventureGameCreatureID, err)
			continue
		}

		// Attempt to load creature portrait image.
		var imageDataURL *string
		imgURL, err := p.Domain.GetAdventureGameCreatureImageDataURL(gameRec.ID, ci.AdventureGameCreatureID)
		if err != nil {
			l.Warn("failed to get creature image >%v<", err)
		} else if imgURL != "" {
			imageDataURL = &imgURL
		}

		maxHealth := creatureDef.MaxHealth
		health := ci.Health
		if health > maxHealth {
			health = maxHealth
		}

		creatures = append(creatures, turnsheet.EncounterCreature{
			CreatureInstanceID: ci.ID,
			Name:               creatureDef.Name,
			Description:        creatureDef.Description,
			Health:             health,
			MaxHealth:          maxHealth,
			AttackDamage:       creatureDef.AttackDamage,
			Defense:            creatureDef.Defense,
			Disposition:        creatureDef.Disposition,
			ImageDataURL:       imageDataURL,
			IsDead:             ci.Health <= 0,
		})
	}

	// Step 4b: Multi-player awareness — detect creature state changes by other players.
	// Compare current creature health against what was recorded in the character's previous
	// encounter sheet to identify kills or damage dealt by other adventurers.
	p.detectMultiPlayerCreatureChanges(l, gameInstanceRec, characterInstanceRec, creatures)

	// Step 5: Load background image (same as location choice sheet).
	var backgroundImage *string
	locationInstanceRec, err := p.Domain.GetAdventureGameLocationInstanceRec(characterInstanceRec.AdventureGameLocationInstanceID, nil)
	if err != nil {
		l.Warn("failed to get location instance >%v<", err)
	} else {
		bgURL, err := p.Domain.GetAdventureGameLocationChoiceTurnSheetImageDataURL(gameRec.ID, locationInstanceRec.AdventureGameLocationID)
		if err != nil {
			l.Warn("failed to get background image >%v<", err)
		} else if bgURL != "" {
			backgroundImage = &bgURL
		}
	}

	// Step 6: Generate turn sheet code.
	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	// Step 7: Read combat events for this sheet. Events are cleared after all processors run.
	displayEvents := ReadTurnEventsForCategories(l, p.Domain, characterInstanceRec, turnsheet.TurnEventCategoryCombat)

	// Step 8: Build sheet data.
	maxActions := 3
	sheetTitle := "Creature Encounter"
	if isReadOnly {
		maxActions = 0
		sheetTitle = "Encounter Record"
	}

	sheetData := turnsheet.MonsterEncounterData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:        convert.Ptr(gameRec.Name),
			GameType:        convert.Ptr("adventure"),
			TurnNumber:      convert.Ptr(gameInstanceRec.CurrentTurn),
			AccountName:     convert.Ptr(accountUserRec.Email),
			TurnSheetTitle:  convert.Ptr(sheetTitle),
			TurnSheetCode:   convert.Ptr(turnSheetCode),
			BackgroundImage: backgroundImage,
			TurnEvents:      displayEvents,
		},
		CharacterName:      characterRec.Name,
		CharacterHealth:    characterInstanceRec.Health,
		CharacterMaxHealth: 100,
		CharacterAttack:    characterAttack,
		CharacterDefense:   characterDefense,
		EquippedWeapon:     equippedWeapon,
		EquippedArmor:      equippedArmor,
		Creatures:          creatures,
		MaxActions:         maxActions,
		ReadOnly:           isReadOnly,
	}

	sheetDataBytes, err := json.Marshal(sheetData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sheet data: %w", err)
	}

	// Step 9: Create turn sheet record.
	turnSheetRec := &game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		AccountID:        accountUserRec.AccountID,
		AccountUserID:    characterRec.AccountUserID,
		TurnNumber:       gameInstanceRec.CurrentTurn,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter,
		SheetOrder:       adventure_game_record.AdventureGameSheetOrderForType(adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter),
		SheetData:        json.RawMessage(sheetDataBytes),
		IsCompleted:      false,
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
	turnSheetRec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

	createdTurnSheetRec, err := p.Domain.CreateGameTurnSheetRec(turnSheetRec)
	if err != nil {
		return nil, fmt.Errorf("failed to create turn sheet record: %w", err)
	}

	adventureTurnSheet := &adventure_game_record.AdventureGameTurnSheet{
		GameID:                           gameInstanceRec.GameID,
		AdventureGameCharacterInstanceID: characterInstanceRec.ID,
		GameTurnSheetID:                  createdTurnSheetRec.ID,
	}
	_, err = p.Domain.CreateAdventureGameTurnSheetRec(adventureTurnSheet)
	if err != nil {
		return nil, fmt.Errorf("failed to create adventure game turn sheet record: %w", err)
	}

	l.Info("created monster encounter turn sheet >%s< for character >%s< with %d creature(s)",
		createdTurnSheetRec.ID, characterInstanceRec.ID, len(creatures))

	return createdTurnSheetRec, nil
}

// inventorySheetHadActions checks whether the inventory management sheet for the current turn had any actions.
func (p *AdventureGameCreatureEncounterProcessor) inventorySheetHadActions(
	_ context.Context,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
) (bool, error) {
	// Get the character's adventure turn sheets for this turn
	adventureTurnSheets, err := p.Domain.GetManyAdventureGameTurnSheetRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID, Val: characterInstanceRec.ID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query adventure turn sheets: %w", err)
	}

	for _, ats := range adventureTurnSheets {
		gameTurnSheet, err := p.Domain.GetGameTurnSheetRec(ats.GameTurnSheetID, nil)
		if err != nil {
			continue
		}
		if gameTurnSheet.SheetType != adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement {
			continue
		}
		if gameTurnSheet.TurnNumber != gameInstanceRec.CurrentTurn {
			continue
		}
		if len(gameTurnSheet.ScannedData) == 0 {
			continue
		}

		var scanData turnsheet.InventoryManagementScanData
		if err := json.Unmarshal(gameTurnSheet.ScannedData, &scanData); err != nil {
			continue
		}

		hasActions := len(scanData.PickUp) > 0 ||
			len(scanData.Drop) > 0 ||
			len(scanData.Equip) > 0 ||
			len(scanData.Unequip) > 0

		if hasActions {
			return true, nil
		}
	}

	return false, nil
}

// resolveEquippedGear returns the display structs for the character's equipped weapon and armor.
func (p *AdventureGameCreatureEncounterProcessor) resolveEquippedGear(l logger.Logger, characterInstanceID string) (weapon *turnsheet.EquippedWeapon, armor *turnsheet.EquippedArmor, err error) {
	inventoryItems, err := p.Domain.GetAdventureGameItemInstanceRecsByCharacterInstance(characterInstanceID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	for _, itemInstance := range inventoryItems {
		if !itemInstance.IsEquipped || !itemInstance.EquipmentSlot.Valid {
			continue
		}
		slot := itemInstance.EquipmentSlot.String
		itemDef, err := p.Domain.GetAdventureGameItemRec(itemInstance.AdventureGameItemID, nil)
		if err != nil {
			l.Warn("failed to get item definition >%s< >%v<", itemInstance.AdventureGameItemID, err)
			continue
		}
		if slot == adventure_game_record.AdventureGameItemEquipmentSlotWeapon && weapon == nil {
			weapon = &turnsheet.EquippedWeapon{
				ItemInstanceID: itemInstance.ID,
				Name:           itemDef.Name,
				Damage:         itemDef.Damage,
			}
		} else if slot != adventure_game_record.AdventureGameItemEquipmentSlotWeapon && armor == nil {
			armor = &turnsheet.EquippedArmor{
				ItemInstanceID: itemInstance.ID,
				Name:           itemDef.Name,
				Defense:        itemDef.Defense,
			}
		}
	}

	return weapon, armor, nil
}

// moveCreatureItemsToLocation transfers all item instances owned by a dead creature to the given location.
// Returns the names of items that were dropped.
func (p *AdventureGameCreatureEncounterProcessor) moveCreatureItemsToLocation(l logger.Logger, creatureInstance *adventure_game_record.AdventureGameCreatureInstance, locationInstanceID string) ([]string, error) {
	items, err := p.Domain.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCreatureInstanceID, Val: creatureInstance.ID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get creature item instances: %w", err)
	}

	var droppedNames []string
	for _, item := range items {
		itemDef, err := p.Domain.GetAdventureGameItemRec(item.AdventureGameItemID, nil)
		if err == nil {
			droppedNames = append(droppedNames, itemDef.Name)
		}
		item.AdventureGameCreatureInstanceID = sql.NullString{}
		item.AdventureGameLocationInstanceID = sql.NullString{String: locationInstanceID, Valid: true}
		item.IsEquipped = false
		item.EquipmentSlot = sql.NullString{}
		_, err = p.Domain.UpdateAdventureGameItemInstanceRec(item)
		if err != nil {
			l.Warn("failed to move item >%s< from creature to location >%v<", item.ID, err)
		}
	}

	return droppedNames, nil
}

// detectMultiPlayerCreatureChanges compares the current creature state to the character's previous
// encounter sheet and generates "world" events for kills or damage caused by other players.
func (p *AdventureGameCreatureEncounterProcessor) detectMultiPlayerCreatureChanges(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	currentCreatures []turnsheet.EncounterCreature,
) {
	// Find the character's most recent monster encounter sheet from the previous turn.
	adventureTurnSheets, err := p.Domain.GetManyAdventureGameTurnSheetRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID, Val: characterInstanceRec.ID},
		},
	})
	if err != nil {
		l.Warn("multi-player awareness: failed to get adventure turn sheets >%v<", err)
		return
	}

	// Find the most recent monster encounter sheet.
	var prevSheetData *turnsheet.MonsterEncounterData
	for _, ats := range adventureTurnSheets {
		gameTurnSheet, err := p.Domain.GetGameTurnSheetRec(ats.GameTurnSheetID, nil)
		if err != nil || gameTurnSheet.SheetType != adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter {
			continue
		}
		// We want the sheet from the previous turn (not the current one).
		if gameTurnSheet.TurnNumber >= gameInstanceRec.CurrentTurn {
			continue
		}
		var data turnsheet.MonsterEncounterData
		if err := json.Unmarshal(gameTurnSheet.SheetData, &data); err != nil {
			continue
		}
		prevSheetData = &data
		break
	}

	if prevSheetData == nil {
		// No previous encounter sheet — skip (turn 1 or first encounter).
		return
	}

	// Build a map of previous creature health by instance ID.
	prevHealth := make(map[string]int, len(prevSheetData.Creatures))
	prevAlive := make(map[string]bool, len(prevSheetData.Creatures))
	for _, c := range prevSheetData.Creatures {
		prevHealth[c.CreatureInstanceID] = c.Health
		prevAlive[c.CreatureInstanceID] = !c.IsDead
	}

	// Build a set of creature names the player attacked this turn from
	// combat events already written to LastTurnEvents. This prevents
	// attributing the player's own damage to "another adventurer".
	playerAttackedNames := make(map[string]bool)
	existingEvents, err := turnsheet.ReadTurnEvents(characterInstanceRec)
	if err != nil {
		l.Warn("multi-player awareness: failed to read turn events >%v<", err)
	} else {
		combatEvents := turnsheet.FilterTurnEventsByCategory(
			existingEvents, turnsheet.TurnEventCategoryCombat,
		)
		for _, e := range combatEvents {
			for _, c := range currentCreatures {
				if strings.Contains(e.Message, c.Name) {
					playerAttackedNames[c.Name] = true
				}
			}
		}
	}

	// Compare current state.
	for _, c := range currentCreatures {
		prev, hadBefore := prevHealth[c.CreatureInstanceID]
		wasAlive := prevAlive[c.CreatureInstanceID]
		if !hadBefore || !wasAlive {
			continue
		}

		// Skip "another adventurer" events for creatures the player attacked
		// — the player's own combat events already tell this story.
		if playerAttackedNames[c.Name] {
			continue
		}

		if c.IsDead || c.Health <= 0 {
			// Creature was alive last sheet but is now dead.
			_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
				Category: turnsheet.TurnEventCategoryWorld,
				Icon:     turnsheet.TurnEventIconWorld,
				Message:  fmt.Sprintf("Another adventurer has slain the %s.", c.Name),
			})
		} else if c.Health < prev {
			// Creature took damage not from this player.
			_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
				Category: turnsheet.TurnEventCategoryWorld,
				Icon:     turnsheet.TurnEventIconWorld,
				Message:  fmt.Sprintf("The %s looks wounded — another adventurer has been here.", c.Name),
			})
		}
	}
}

// resetDeadCharacter moves a dead character to the starting location and restores their health.
// If startingLocationName is non-nil, it is set to the name of the starting location for narrative purposes.
func (p *AdventureGameCreatureEncounterProcessor) resetDeadCharacter(l logger.Logger, gameInstanceRec *game_record.GameInstance, characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance, startingLocationName *string) error {
	// Find the starting location instance for this game instance.
	locationInstances, err := p.Domain.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get location instances: %w", err)
	}

	for _, li := range locationInstances {
		locationRec, err := p.Domain.GetAdventureGameLocationRec(li.AdventureGameLocationID, nil)
		if err != nil {
			continue
		}
		if locationRec.IsStartingLocation {
			characterInstanceRec.AdventureGameLocationInstanceID = li.ID
			characterInstanceRec.Health = characterStartingHealth
			if startingLocationName != nil {
				*startingLocationName = locationRec.Name
			}
			l.Info("resetting dead character to starting location >%s< with health >%d<",
				li.ID, characterStartingHealth)
			return nil
		}
	}

	// No starting location found — just restore health without moving.
	characterInstanceRec.Health = characterStartingHealth
	l.Warn("no starting location found for game instance >%s< — restored health only", gameInstanceRec.ID)
	return nil
}
