package harness

import (
	"context"
	"fmt"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// AdventureGameRecords holds all adventure game records created by processAdventureGameConfig.
type AdventureGameRecords struct {
	Items                    []*adventure_game_record.AdventureGameItem
	Locations                []*adventure_game_record.AdventureGameLocation
	Creatures                []*adventure_game_record.AdventureGameCreature
	LocationLinks            []*adventure_game_record.AdventureGameLocationLink
	LocationLinkRequirements []*adventure_game_record.AdventureGameLocationLinkRequirement
	Characters               []*adventure_game_record.AdventureGameCharacter
}

// processAdventureGameConfig creates adventure game records from config and returns all created records.
func (t *Testing) processAdventureGameConfig(gameConfig GameConfig, gameRec *game_record.Game) (*AdventureGameRecords, error) {
	l := t.Logger("processAdventureGameConfig")

	out := &AdventureGameRecords{}

	for _, itemConfig := range gameConfig.AdventureGameItemConfigs {
		itemRec, err := t.createAdventureGameItemRec(itemConfig, gameRec)
		if err != nil {
			l.Warn("failed creating game_item record >%v<", err)
			return nil, err
		}
		out.Items = append(out.Items, itemRec)
		l.Debug("created game_item record for game >%s<", gameRec.ID)
	}

	for _, locationConfig := range gameConfig.AdventureGameLocationConfigs {
		gameLocationRec, err := t.createAdventureGameLocationRec(locationConfig, gameRec)
		if err != nil {
			l.Warn("failed creating game location record >%v<", err)
			return nil, err
		}
		out.Locations = append(out.Locations, gameLocationRec)
		l.Debug("created game location record ID >%s< Name >%s<", gameLocationRec.ID, gameLocationRec.Name)

		// Create background image for location if specified
		if locationConfig.BackgroundImage != nil {
			cfg := *locationConfig.BackgroundImage
			cfg.RecordID = gameLocationRec.ID
			if cfg.TurnSheetType == "" {
				cfg.TurnSheetType = adventure_game_record.AdventureGameTurnSheetTypeLocationChoice
			}
			_, err = t.processGameImageConfig(cfg, gameRec)
			if err != nil {
				l.Warn("failed creating location background image >%v<", err)
				return nil, err
			}
			l.Debug("created location background image for location >%s<", gameLocationRec.ID)
		}
	}

	for _, creatureConfig := range gameConfig.AdventureGameCreatureConfigs {
		creatureRec, err := t.createAdventureGameCreatureRec(creatureConfig, gameRec)
		if err != nil {
			l.Warn("failed creating game creature record >%v<", err)
			return nil, err
		}
		out.Creatures = append(out.Creatures, creatureRec)
		l.Debug("created game creature record ID >%s< Name >%s<", creatureRec.ID, creatureRec.Name)

		// Create portrait image for creature if specified
		if creatureConfig.PortraitImage != nil {
			cfg := *creatureConfig.PortraitImage
			cfg.RecordID = creatureRec.ID
			_, err = t.processCreaturePortraitImageConfig(cfg, gameRec)
			if err != nil {
				l.Warn("failed creating creature portrait image >%v<", err)
				return nil, err
			}
			l.Debug("created creature portrait image for creature >%s<", creatureRec.ID)
		}
	}

	for _, linkConfig := range gameConfig.AdventureGameLocationLinkConfigs {
		gameLocationLinkRec, err := t.createAdventureGameLocationLinkRec(linkConfig, gameRec)
		if err != nil {
			l.Warn("failed creating location link record >%v<", err)
			return nil, err
		}
		out.LocationLinks = append(out.LocationLinks, gameLocationLinkRec)
		l.Debug("created location link record ID >%s<", gameLocationLinkRec.ID)

		for _, reqConfig := range linkConfig.AdventureGameLocationLinkRequirementConfigs {
			reqRec, err := t.createAdventureGameLocationLinkRequirementRec(reqConfig, gameLocationLinkRec)
			if err != nil {
				l.Warn("failed creating game_location_link_requirement record >%v<", err)
				return nil, err
			}
			out.LocationLinkRequirements = append(out.LocationLinkRequirements, reqRec)
			l.Debug("created game_location_link_requirement record for game >%s<", gameRec.ID)
		}
	}

	for _, charConfig := range gameConfig.AdventureGameCharacterConfigs {
		charRec, err := t.createAdventureGameCharacterRec(charConfig, gameRec)
		if err != nil {
			l.Warn("failed creating game_character record >%v<", err)
			return nil, err
		}
		out.Characters = append(out.Characters, charRec)
		l.Debug("created game_character record for game >%s<", gameRec.ID)
	}

	for _, objectConfig := range gameConfig.AdventureGameLocationObjectConfigs {
		_, err := t.createAdventureGameLocationObjectRec(objectConfig, gameRec)
		if err != nil {
			l.Warn("failed creating adventure_game_location_object record >%v<", err)
			return nil, err
		}
	}

	for _, placementConfig := range gameConfig.AdventureGameCreaturePlacementConfigs {
		_, err := t.createAdventureGameCreaturePlacementRec(placementConfig, gameRec)
		if err != nil {
			l.Warn("failed creating adventure_game_creature_placement record >%v<", err)
			return nil, err
		}
	}

	for _, placementConfig := range gameConfig.AdventureGameItemPlacementConfigs {
		_, err := t.createAdventureGameItemPlacementRec(placementConfig, gameRec)
		if err != nil {
			l.Warn("failed creating adventure_game_item_placement record >%v<", err)
			return nil, err
		}
	}

	return out, nil
}

// createAdventureGameInstanceRecords starts the game instance and creates initial turn sheets when
// ShouldStartGameInstance is true. Instance records (locations, creatures, items, characters) are
// created by StartGameInstance via PopulateAdventureGameInstanceData.
func (t *Testing) createAdventureGameInstanceRecords(gameInstanceConfig GameInstanceConfig, gameInstanceRec *game_record.GameInstance) error {
	l := t.Logger("createAdventureGameInstanceRecords")

	if !gameInstanceConfig.ShouldStartGameInstance {
		return nil
	}

	if gameInstanceConfig.Reference == "" {
		return fmt.Errorf("game_instance config must have a reference when ShouldStartGameInstance is true")
	}

	instanceRec := gameInstanceRec
	if instanceRec.Status != game_record.GameInstanceStatusStarted {
		var instanceData *domain.AdventureGameInstanceData
		var err error
		instanceRec, instanceData, err = t.Domain.(*domain.Domain).StartGameInstance(gameInstanceRec.ID)
		if err != nil {
			l.Warn("failed starting game instance >%v<", err)
			return err
		}

		if instanceData != nil {
			t.addAdventureGameInstanceDataToStores(instanceData, gameInstanceRec.ID)
		}
	}

	ctx := context.Background()

	turnSheets, err := t.generateTurnSheetsForGameInstanceInTx(ctx, instanceRec, gameInstanceConfig.Reference)
	if err != nil {
		l.Warn("failed generating adventure game turn sheets >%v<", err)
		return err
	}

	if len(gameInstanceConfig.TurnSheetRefConfigs) > 0 {
		if err := t.assignTurnSheetRefs(gameInstanceConfig.TurnSheetRefConfigs, turnSheets); err != nil {
			l.Warn("failed assigning turn sheet refs for instance >%s< >%v<", gameInstanceConfig.Reference, err)
			return err
		}
	}

	return nil
}

// addAdventureGameInstanceDataToStores adds instance records returned by StartGameInstance to the
// harness data and teardown stores, and maps the first record of each type to its canonical ref.
func (t *Testing) addAdventureGameInstanceDataToStores(instanceData *domain.AdventureGameInstanceData, gameInstanceID string) {
	l := t.Logger("addAdventureGameInstanceDataToStores")

	for i, rec := range instanceData.LocationInstances {
		t.Data.AddAdventureGameLocationInstanceRec(rec)
		t.teardownData.AddAdventureGameLocationInstanceRec(rec)
		switch i {
		case 0:
			t.Data.Refs.AdventureGameLocationInstanceRefs[GameLocationInstanceOneRef] = rec.ID
		case 1:
			t.Data.Refs.AdventureGameLocationInstanceRefs[GameLocationInstanceTwoRef] = rec.ID
		}
	}

	for i, rec := range instanceData.CreatureInstances {
		t.Data.AddAdventureGameCreatureInstanceRec(rec)
		t.teardownData.AddAdventureGameCreatureInstanceRec(rec)
		if i == 0 {
			t.Data.Refs.AdventureGameCreatureInstanceRefs[GameCreatureInstanceOneRef] = rec.ID
		}
	}

	for _, rec := range instanceData.LocationObjectInstances {
		t.Data.AddAdventureGameLocationObjectInstanceRec(rec)
		t.teardownData.AddAdventureGameLocationObjectInstanceRec(rec)
	}

	for i, rec := range instanceData.ItemInstances {
		t.Data.AddAdventureGameItemInstanceRec(rec)
		t.teardownData.AddAdventureGameItemInstanceRec(rec)
		if i == 0 {
			t.Data.Refs.AdventureGameItemInstanceRefs[GameItemInstanceOneRef] = rec.ID
		}
	}

	for i, rec := range instanceData.CharacterInstances {
		t.Data.AddAdventureGameCharacterInstanceRec(rec)
		t.teardownData.AddAdventureGameCharacterInstanceRec(rec)
		switch i {
		case 0:
			t.Data.Refs.AdventureGameCharacterInstanceRefs[GameCharacterInstanceOneRef] = rec.ID
		case 1:
			t.Data.Refs.AdventureGameCharacterInstanceRefs[GameCharacterInstanceTwoRef] = rec.ID
		}
	}

	l.Debug("added adventure game instance data to stores for game instance >%s<: locations=%d creatures=%d objects=%d items=%d characters=%d",
		gameInstanceID,
		len(instanceData.LocationInstances),
		len(instanceData.CreatureInstances),
		len(instanceData.LocationObjectInstances),
		len(instanceData.ItemInstances),
		len(instanceData.CharacterInstances))
}

// removeAdventureGameRecords removes the adventure game records for a game
func (t *Testing) removeAdventureGameRecords() error {
	l := t.Logger("removeAdventureGameRecords")

	l.Debug("removing adventure game records")

	// Remove location object records before location link requirements and locations
	l.Debug("removing >%d< adventure game location object effect records", len(t.teardownData.AdventureGameLocationObjectEffectRecs))
	for _, effectRec := range t.teardownData.AdventureGameLocationObjectEffectRecs {
		if effectRec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveAdventureGameLocationObjectEffectRec(effectRec.ID); err != nil {
			l.Warn("failed removing adventure game location object effect record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< adventure game location object records", len(t.teardownData.AdventureGameLocationObjectRecs))
	for _, objectRec := range t.teardownData.AdventureGameLocationObjectRecs {
		if objectRec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveAdventureGameLocationObjectRec(objectRec.ID); err != nil {
			l.Warn("failed removing adventure game location object record >%v<", err)
			return err
		}
	}

	// Remove game location link requirements before creatures and items (requirements reference both)
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

	// Remove ALL creature and item placements for each game before removing creatures/items.
	// This includes both harness-created placements (tracked in teardownData) and any additional
	// placements created by API handlers during tests, which would cause FK violations if left behind.
	for _, gameRec := range t.teardownData.GameRecs {
		if gameRec.ID == "" {
			continue
		}
		creaturePlacements, err := t.Domain.(*domain.Domain).GetManyAdventureGameCreaturePlacementRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameCreaturePlacementGameID, Val: gameRec.ID},
			},
		})
		if err != nil {
			l.Warn("failed fetching creature placements for game >%s< >%v<", gameRec.ID, err)
			return err
		}
		l.Debug("removing >%d< adventure game creature placement records for game >%s<", len(creaturePlacements), gameRec.ID)
		for _, placementRec := range creaturePlacements {
			if err := t.Domain.(*domain.Domain).RemoveAdventureGameCreaturePlacementRec(placementRec.ID); err != nil {
				l.Warn("failed removing adventure game creature placement record >%v<", err)
				return err
			}
		}

		itemPlacements, err := t.Domain.(*domain.Domain).GetManyAdventureGameItemPlacementRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameItemPlacementGameID, Val: gameRec.ID},
			},
		})
		if err != nil {
			l.Warn("failed fetching item placements for game >%s< >%v<", gameRec.ID, err)
			return err
		}
		l.Debug("removing >%d< adventure game item placement records for game >%s<", len(itemPlacements), gameRec.ID)
		for _, placementRec := range itemPlacements {
			if err := t.Domain.(*domain.Domain).RemoveAdventureGameItemPlacementRec(placementRec.ID); err != nil {
				l.Warn("failed removing adventure game item placement record >%v<", err)
				return err
			}
		}
	}

	// Remove game creature records (after requirements)
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

