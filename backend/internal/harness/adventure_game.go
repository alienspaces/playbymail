package harness

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// createAdventureGameRecords creates the adventure game records for a game
func (t *Testing) createAdventureGameRecords(gameConfig GameConfig, gameRec *game_record.Game) error {
	l := t.Logger("createAdventureGameRecords")

	for _, itemConfig := range gameConfig.AdventureGameItemConfigs {
		_, err := t.createAdventureGameItemRec(itemConfig, gameRec)
		if err != nil {
			l.Warn("failed creating game_item record >%v<", err)
			return err
		}
		l.Debug("created game_item record for game >%s<", gameRec.ID)
	}

	for _, locationConfig := range gameConfig.AdventureGameLocationConfigs {
		gameLocationRec, err := t.createAdventureGameLocationRec(locationConfig, gameRec)
		if err != nil {
			l.Warn("failed creating game location record >%v<", err)
			return err
		}
		l.Debug("created game location record ID >%s< Name >%s<", gameLocationRec.ID, gameLocationRec.Name)

		// Create background image for location if specified
		if locationConfig.BackgroundImagePath != "" {
			_, err = t.createGameImageRecFromPath(gameRec.ID, gameLocationRec.ID, locationConfig.BackgroundImagePath)
			if err != nil {
				l.Warn("failed creating location background image >%v<", err)
				return err
			}
			l.Debug("created location background image for location >%s<", gameLocationRec.ID)
		}
	}

	for _, creatureConfig := range gameConfig.AdventureGameCreatureConfigs {
		creatureRec, err := t.createAdventureGameCreatureRec(creatureConfig, gameRec)
		if err != nil {
			l.Warn("failed creating game creature record >%v<", err)
			return err
		}
		l.Debug("created game creature record ID >%s< Name >%s<", creatureRec.ID, creatureRec.Name)
	}

	for _, linkConfig := range gameConfig.AdventureGameLocationLinkConfigs {
		gameLocationLinkRec, err := t.createAdventureGameLocationLinkRec(linkConfig, gameRec)
		if err != nil {
			l.Warn("failed creating location link record >%v<", err)
			return err
		}
		l.Debug("created location link record ID >%s<", gameLocationLinkRec.ID)

		for _, reqConfig := range linkConfig.AdventureGameLocationLinkRequirementConfigs {
			_, err = t.createAdventureGameLocationLinkRequirementRec(reqConfig, gameLocationLinkRec)
			if err != nil {
				l.Warn("failed creating game_location_link_requirement record >%v<", err)
				return err
			}
			l.Debug("created game_location_link_requirement record for game >%s<", gameRec.ID)
		}
	}

	for _, charConfig := range gameConfig.AdventureGameCharacterConfigs {
		_, err := t.createAdventureGameCharacterRec(charConfig, gameRec)
		if err != nil {
			l.Warn("failed creating game_character record >%v<", err)
			return err
		}
		l.Debug("created game_character record for game >%s<", gameRec.ID)
	}

	return nil
}

// createAdventureGameInstanceRecords creates the adventure game instance records for a game instance
func (t *Testing) createAdventureGameInstanceRecords(gameInstanceConfig GameInstanceConfig, gameInstanceRec *game_record.GameInstance) error {
	l := t.Logger("createAdventureGameInstanceRecords")

	// Create game location instance records for this game instance
	for _, locationInstanceConfig := range gameInstanceConfig.AdventureGameLocationInstanceConfigs {
		locationInstanceRec, err := t.createAdventureGameLocationInstanceRec(locationInstanceConfig, gameInstanceRec)
		if err != nil {
			l.Warn("failed creating adventure game location instance record >%v<", err)
			return err
		}
		l.Debug("created adventure game location instance record ID >%s<", locationInstanceRec.ID)
	}

	// Create game creature instance records for this game instance
	for _, creatureInstanceConfig := range gameInstanceConfig.AdventureGameCreatureInstanceConfigs {
		creatureInstanceRec, err := t.createAdventureGameCreatureInstanceRec(creatureInstanceConfig, gameInstanceRec)
		if err != nil {
			l.Warn("failed creating adventure game creature instance record >%v<", err)
			return err
		}
		l.Debug("created adventure game creature instance record ID >%s<", creatureInstanceRec.ID)
	}

	// Create game character instance records for this game instance
	for _, characterInstanceConfig := range gameInstanceConfig.AdventureGameCharacterInstanceConfigs {
		characterInstanceRec, err := t.createAdventureGameCharacterInstanceRec(characterInstanceConfig, gameInstanceRec)
		if err != nil {
			l.Warn("failed creating adventure game character instance record >%v<", err)
			return err
		}
		l.Debug("created adventure game character instance record ID >%s<", characterInstanceRec.ID)
	}

	// Create game item instance records for this game instance
	for _, itemInstanceConfig := range gameInstanceConfig.AdventureGameItemInstanceConfigs {
		itemInstanceRec, err := t.createAdventureGameItemInstanceRec(itemInstanceConfig, gameInstanceRec)
		if err != nil {
			l.Warn("failed creating adventure game item instance record >%v<", err)
			return err
		}
		l.Debug("created adventure game item instance record ID >%s<", itemInstanceRec.ID)
	}

	if len(gameInstanceConfig.GameTurnConfigs) == 0 {
		return nil
	}

	if gameInstanceConfig.Reference == "" {
		return fmt.Errorf("game_instance config must have a reference when configuring turns")
	}

	instanceRec := gameInstanceRec
	if instanceRec.Status != game_record.GameInstanceStatusStarted {
		var err error
		instanceRec, err = t.Domain.(*domain.Domain).StartGameInstance(gameInstanceRec.ID)
		if err != nil {
			l.Warn("failed starting game instance >%v<", err)
			return err
		}
	}

	turnConfigs := normalizeTurnConfigs(gameInstanceConfig.GameTurnConfigs)

	ctx := context.Background()

	turnSheets, err := t.generateTurnSheetsForGameInstanceInTx(ctx, instanceRec, gameInstanceConfig.Reference)
	if err != nil {
		l.Warn("failed generating adventure game turn sheets >%v<", err)
		return err
	}

	turnSheetsCache := map[int][]*game_record.GameTurnSheet{}
	if len(turnSheets) > 0 {
		turnSheetsCache[turnSheets[0].TurnNumber] = turnSheets
	}

	for idx, turnCfg := range turnConfigs {
		turnSheetsForConfig, ok := turnSheetsCache[turnCfg.TurnNumber]
		if !ok {
			turnSheetsForConfig, err = t.getTurnSheetsForTurn(instanceRec.ID, turnCfg.TurnNumber)
			if err != nil {
				l.Warn("failed fetching turn sheets for turn >%d< >%v<", turnCfg.TurnNumber, err)
				return err
			}
			turnSheetsCache[turnCfg.TurnNumber] = turnSheetsForConfig
		}

		if err := t.assignTurnSheetReferencesForTurn(turnCfg, turnSheetsForConfig); err != nil {
			l.Warn("failed assigning turn sheet refs for turn >%d< >%v<", turnCfg.TurnNumber, err)
			return err
		}

		readyForProcessing, err := t.applyConfiguredScanData(ctx, turnCfg)
		if err != nil {
			l.Warn("failed applying scan data for turn >%d< >%v<", turnCfg.TurnNumber, err)
			return err
		}

		isLastTurn := idx == len(turnConfigs)-1
		if !readyForProcessing {
			if !isLastTurn {
				return fmt.Errorf("turn %d is missing scan data but subsequent turns are configured", turnCfg.TurnNumber)
			}
			continue
		}

		if err := t.processGameTurnForInstanceInTx(ctx, instanceRec.ID); err != nil {
			l.Warn("failed processing turn >%d< for instance >%s< >%v<", turnCfg.TurnNumber, instanceRec.ID, err)
			return err
		}

		nextTurn := turnCfg.TurnNumber + 1
		nextSheets, err := t.getTurnSheetsForTurn(instanceRec.ID, nextTurn)
		if err != nil {
			l.Warn("failed fetching next turn sheets >%v<", err)
			return err
		}
		turnSheetsCache[nextTurn] = nextSheets
	}

	return nil
}

// removeAdventureGameRecords removes the adventure game records for a game
func (t *Testing) removeAdventureGameRecords() error {
	l := t.Logger("removeAdventureGameRecords")

	l.Debug("removing adventure game records")

	// Remove game creature records
	l.Debug("removing >%d< game creature records", len(t.teardownData.AdventureGameCreatureRecs))
	for _, creatureRec := range t.teardownData.AdventureGameCreatureRecs {
		l.Debug("[teardown] game creature ID: >%s<", creatureRec.ID)
		if creatureRec.ID == "" {
			l.Warn("[teardown] skipping game creature with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameCreatureRec(creatureRec.ID)
		if err != nil {
			l.Warn("failed removing game creature record >%v<", err)
			return err
		}
	}

	// Remove game location link requirements
	l.Debug("removing >%d< game location link requirement records", len(t.teardownData.AdventureGameLocationLinkRequirementRecs))
	for _, reqRec := range t.teardownData.AdventureGameLocationLinkRequirementRecs {
		l.Debug("[teardown] game location link requirement ID: >%s<", reqRec.ID)
		if reqRec.ID == "" {
			l.Warn("[teardown] skipping game location link requirement with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameLocationLinkRequirementRec(reqRec.ID)
		if err != nil {
			l.Warn("failed removing game location link requirement record >%v<", err)
			return err
		}
	}

	// Remove game location links
	l.Debug("removing >%d< game location link records", len(t.teardownData.AdventureGameLocationLinkRecs))
	for _, linkRec := range t.teardownData.AdventureGameLocationLinkRecs {
		l.Debug("[teardown] game location link ID: >%s<", linkRec.ID)
		if linkRec.ID == "" {
			l.Warn("[teardown] skipping game location link with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameLocationLinkRec(linkRec.ID)
		if err != nil {
			l.Warn("failed removing game location link record >%v<", err)
			return err
		}
	}

	// Remove game location records before games to avoid FK errors
	l.Debug("removing >%d< game location records", len(t.teardownData.AdventureGameLocationRecs))
	for _, gameLocationRec := range t.teardownData.AdventureGameLocationRecs {
		l.Debug("[teardown] game location ID: >%s<", gameLocationRec.ID)
		if gameLocationRec.ID == "" {
			l.Warn("[teardown] skipping game location with empty ID")
			continue
		}
		l.Debug("removing game location record ID >%s<", gameLocationRec.ID)
		err := t.Domain.(*domain.Domain).RemoveAdventureGameLocationRec(gameLocationRec.ID)
		if err != nil {
			l.Warn("failed removing game location record >%v<", err)
			return err
		}
	}

	// Remove game item records before games to avoid FK errors
	l.Debug("removing >%d< game item records", len(t.teardownData.AdventureGameItemRecs))
	for _, itemRec := range t.teardownData.AdventureGameItemRecs {
		l.Debug("[teardown] game item ID: >%s<", itemRec.ID)
		if itemRec.ID == "" {
			l.Warn("[teardown] skipping game item with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameItemRec(itemRec.ID)
		if err != nil {
			l.Warn("failed removing game item record >%v<", err)
			return err
		}
	}

	// Remove game character records
	l.Debug("removing >%d< game character records", len(t.teardownData.AdventureGameCharacterRecs))
	for _, charRec := range t.teardownData.AdventureGameCharacterRecs {
		l.Debug("[teardown] game character ID: >%s<", charRec.ID)
		if charRec.ID == "" {
			l.Warn("[teardown] skipping game character with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameCharacterRec(charRec.ID)
		if err != nil {
			l.Warn("failed removing game character record >%v<", err)
			return err
		}
	}

	return nil
}

// removeAdventureGameInstanceRecords removes the adventure game instance records for a game instance
func (t *Testing) removeAdventureGameInstanceRecords() error {
	l := t.Logger("removeAdventureGameInstanceRecords")

	l.Debug("removing adventure game instance records")

	// Remove adventure game turn sheet records first (they reference character instances)
	l.Debug("removing >%d< adventure game turn sheet records", len(t.teardownData.AdventureGameTurnSheetRecs))
	for _, turnSheetRec := range t.teardownData.AdventureGameTurnSheetRecs {
		l.Debug("[teardown] adventure game turn sheet ID: >%s<", turnSheetRec.ID)
		if turnSheetRec.ID == "" {
			l.Warn("[teardown] skipping adventure game turn sheet with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameTurnSheetRec(turnSheetRec.ID)
		if err != nil {
			l.Warn("failed removing adventure game turn sheet record >%v<", err)
			return err
		}
	}

	// Remove adventure game creature instances
	l.Debug("removing >%d< adventure game creature instance records", len(t.teardownData.AdventureGameCreatureInstanceRecs))
	for _, creatureInstanceRec := range t.teardownData.AdventureGameCreatureInstanceRecs {
		l.Debug("[teardown] adventure game creature instance ID: >%s<", creatureInstanceRec.ID)
		if creatureInstanceRec.ID == "" {
			l.Warn("[teardown] skipping adventure game creature instance with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameCreatureInstanceRec(creatureInstanceRec.ID)
		if err != nil {
			l.Warn("failed removing game creature instance record >%v<", err)
			return err
		}
	}

	// Remove adventure game character instances
	l.Debug("removing >%d< adventure game character instance records", len(t.teardownData.AdventureGameCharacterInstanceRecs))
	for _, characterInstanceRec := range t.teardownData.AdventureGameCharacterInstanceRecs {
		l.Debug("[teardown] adventure game character instance ID: >%s<", characterInstanceRec.ID)
		if characterInstanceRec.ID == "" {
			l.Warn("[teardown] skipping adventure game character instance with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameCharacterInstanceRec(characterInstanceRec.ID)
		if err != nil {
			l.Warn("failed removing game character instance record >%v<", err)
			return err
		}
	}

	// Remove adventure game item instances
	l.Debug("removing >%d< adventure game item instance records", len(t.teardownData.AdventureGameItemInstanceRecs))
	for _, itemInstanceRec := range t.teardownData.AdventureGameItemInstanceRecs {
		l.Debug("[teardown] adventure game item instance ID: >%s<", itemInstanceRec.ID)
		if itemInstanceRec.ID == "" {
			l.Warn("[teardown] skipping adventure game item instance with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameItemInstanceRec(itemInstanceRec.ID)
		if err != nil {
			l.Warn("failed removing adventure game item instance record >%v<", err)
			return err
		}
	}

	// Remove adventure game location instances
	l.Debug("removing >%d< adventure game location instance records", len(t.teardownData.AdventureGameLocationInstanceRecs))
	for _, locationInstanceRec := range t.teardownData.AdventureGameLocationInstanceRecs {
		l.Debug("[teardown] adventure game location instance ID: >%s<", locationInstanceRec.ID)
		if locationInstanceRec.ID == "" {
			l.Warn("[teardown] skipping adventure game location instance with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAdventureGameLocationInstanceRec(locationInstanceRec.ID)
		if err != nil {
			l.Warn("failed removing adventure game location instance record >%v<", err)
			return err
		}
	}

	return nil
}
