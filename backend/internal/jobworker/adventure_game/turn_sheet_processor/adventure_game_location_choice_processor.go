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
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

// AdventureGameLocationChoiceProcessor implements the TurnSheetProcessor interface
// (defined in the parent adventure_game package)

// AdventureGameLocationChoiceProcessor processes location choice turn sheet business logic for adventure games
type AdventureGameLocationChoiceProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewAdventureGameLocationChoiceProcessor creates a new adventure game location choice processor
func NewAdventureGameLocationChoiceProcessor(l logger.Logger, d *domain.Domain) (*AdventureGameLocationChoiceProcessor, error) {
	l = l.WithFunctionContext("NewAdventureGameLocationChoiceProcessor")

	p := &AdventureGameLocationChoiceProcessor{
		Logger: l,
		Domain: d,
	}
	return p, nil
}

// GetSheetType returns the sheet type this processor handles (implements TurnSheetProcessor interface)
func (p *AdventureGameLocationChoiceProcessor) GetSheetType() string {
	return adventure_game_record.AdventureGameTurnSheetTypeLocationChoice
}

// ProcessTurnSheetResponse processes a single turn sheet response (implements TurnSheetProcessor interface)
func (p *AdventureGameLocationChoiceProcessor) ProcessTurnSheetResponse(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("AdventureGameLocationChoiceProcessor/ProcessTurnSheetResponse")

	l.Info("processing location choice for turn sheet >%s< for character >%s<", turnSheet.ID, characterInstanceRec.ID)

	// Verify this is a location choice sheet
	if turnSheet.SheetType != adventure_game_record.AdventureGameTurnSheetTypeLocationChoice {
		l.Warn("expected location choice sheet type, got >%s<", turnSheet.SheetType)
		return fmt.Errorf("invalid sheet type: expected %s, got %s", adventure_game_record.AdventureGameTurnSheetTypeLocationChoice, turnSheet.SheetType)
	}

	// Step 1: Parse the player's location choice from ScannedData
	var scanData turnsheet.LocationChoiceScanData
	if err := json.Unmarshal(turnSheet.ScannedData, &scanData); err != nil {
		l.Warn("failed to unmarshal scanned data >%v<", err)
		return fmt.Errorf("failed to parse scanned data: %w", err)
	}

	// Handle object interaction if present (mutually exclusive with location choice)
	if scanData.ObjectChoice != "" {
		l.Info("player chose object action >%s<", scanData.ObjectChoice)
		if err := p.processObjectChoice(l, gameInstanceRec, characterInstanceRec, scanData.ObjectChoice); err != nil {
			l.Warn("failed to process object choice >%v<", err)
			return fmt.Errorf("failed to process object choice: %w", err)
		}
		return nil
	}

	choices := scanData.GetChoices()
	if len(choices) == 0 {
		l.Info("no location choice in scanned data — character stays at current location")
		return nil
	}

	chosenLocationID := choices[0]

	l.Info("player chose location >%s<", chosenLocationID)

	// Step 2: Parse SheetData to get the original location options and validate
	var sheetData turnsheet.LocationChoiceData
	if err := json.Unmarshal(turnSheet.SheetData, &sheetData); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return fmt.Errorf("failed to parse sheet data: %w", err)
	}

	// Validate the choice is one of the available non-locked options
	isValidChoice := false
	var chosenLocationOption turnsheet.LocationOption
	for _, option := range sheetData.LocationOptions {
		if option.LocationID == chosenLocationID && !option.IsLocked {
			isValidChoice = true
			chosenLocationOption = option
			break
		}
	}

	if !isValidChoice {
		l.Warn("invalid location choice >%s< not in available options", chosenLocationID)
		return fmt.Errorf("invalid location choice: %s is not an available option", chosenLocationID)
	}

	// Step 3: Apply flee penalty from aggressive creatures at the character's CURRENT location.
	// This happens before moving the character so we can check their current location.
	currentLocationInstanceID := characterInstanceRec.AdventureGameLocationInstanceID
	isMoving := currentLocationInstanceID != chosenLocationID

	// Resolve traversal description for narrative.
	var traversalDescription string
	if isMoving {
		// Resolve the current location instance to its definition ID before querying links,
		// since from_adventure_game_location_id stores definition IDs, not instance IDs.
		currentLocInstanceRec, err := p.Domain.GetAdventureGameLocationInstanceRec(currentLocationInstanceID, nil)
		if err == nil {
			locationLinks, err := p.Domain.GetManyAdventureGameLocationLinkRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: adventure_game_record.FieldAdventureGameLocationLinkFromAdventureGameLocationID, Val: currentLocInstanceRec.AdventureGameLocationID},
				},
			})
			if err == nil {
				for _, lnk := range locationLinks {
					// Find the link whose destination matches.
					destInst, err := p.getLocationInstanceForLocation(gameInstanceRec.ID, lnk.ToAdventureGameLocationID)
					if err == nil && destInst.ID == chosenLocationID {
						if lnk.TraversalDescription.Valid {
							traversalDescription = lnk.TraversalDescription.String
						}
						break
					}
				}
			}
		}

		// Apply flee penalty.
		if err := p.applyFleePenalty(l, gameInstanceRec, characterInstanceRec, currentLocationInstanceID); err != nil {
			l.Warn("failed to apply flee penalty >%v<", err)
			// Non-fatal: continue with movement.
		}
	}

	// Step 4: Update character's location.
	characterInstanceRec.AdventureGameLocationInstanceID = chosenLocationID
	characterInstanceRec, err := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec)
	if err != nil {
		l.Warn("failed to update character location >%v<", err)
		return fmt.Errorf("failed to update character location: %w", err)
	}

	// Step 5: Generate movement narrative event.
	if isMoving {
		// Resolve destination location name.
		destLocInstanceRec, err := p.Domain.GetAdventureGameLocationInstanceRec(chosenLocationID, nil)
		var destLocationName string
		if err == nil {
			destLocRec, err := p.Domain.GetAdventureGameLocationRec(destLocInstanceRec.AdventureGameLocationID, nil)
			if err == nil {
				destLocationName = destLocRec.Name
			}
		}

		linkName := chosenLocationOption.LocationLinkName
		var movementMsg string
		if traversalDescription != "" {
			movementMsg = fmt.Sprintf("You took %s to %s. %s", linkName, destLocationName, traversalDescription)
		} else {
			movementMsg = fmt.Sprintf("You took %s to %s.", linkName, destLocationName)
		}
		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryMovement,
			Icon:     turnsheet.TurnEventIconMovement,
			Message:  movementMsg,
		})

		// Persist events.
		if _, saveErr := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec); saveErr != nil {
			l.Warn("failed to save movement narrative event >%v<", saveErr)
		}
	}

	l.Info("successfully updated character >%s< to location >%s< via pathway >%s<", characterInstanceRec.ID, chosenLocationID, chosenLocationOption.LocationLinkName)

	return nil
}

// applyFleePenalty inflicts free attacks from aggressive creatures on a character who is moving away.
func (p *AdventureGameLocationChoiceProcessor) applyFleePenalty(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	locationInstanceID string,
) error {
	l = l.WithFunctionContext("applyFleePenalty")
	l.Info("applying flee penalty for character >%s< at location >%s<", characterInstanceRec.ID, locationInstanceID)

	creatureInstances, err := p.Domain.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceRec.ID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, Val: locationInstanceID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get creature instances: %w", err)
	}

	// Determine character's armor defense.
	_, armorDefense, err := ResolveEquipmentStats(l, p.Domain, characterInstanceRec.ID)
	if err != nil {
		return fmt.Errorf("failed to resolve equipment stats for flee penalty: %w", err)
	}

	totalFleeDamage := 0
	for _, ci := range creatureInstances {
		if ci.Health <= 0 {
			l.Info("creature >%s< is dead — skipping", ci.AdventureGameCreatureID)
			continue
		}

		creatureDef, err := p.Domain.GetAdventureGameCreatureRec(ci.AdventureGameCreatureID, nil)
		if err != nil {
			return fmt.Errorf("failed to get creature definition >%s< for flee penalty: %w", ci.AdventureGameCreatureID, err)
		}

		// Only aggressive creatures get a free attack on flee.
		if creatureDef.Disposition != adventure_game_record.AdventureGameCreatureDispositionAggressive {
			l.Info("creature >%s< is not aggressive — skipping", ci.AdventureGameCreatureID)
			continue
		}

		damage := creatureDef.AttackDamage - armorDefense
		if damage < 1 {
			l.Info("creature >%s< attack damage < 1 — setting to 1", ci.AdventureGameCreatureID)
			damage = 1
		}
		totalFleeDamage += damage

		// Generate flee narrative event per creature.
		attackDesc := creatureDef.AttackDescription
		if attackDesc == "" {
			attackDesc = "attacks you"
		}
		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryFlee,
			Icon:     turnsheet.TurnEventIconFlee,
			Message:  fmt.Sprintf("As you fled, the %s %s for %d damage.", creatureDef.Name, attackDesc, damage),
		})

		l.Info("aggressive creature >%s< attacks fleeing character for %d damage", creatureDef.Name, damage)
	}

	if totalFleeDamage > 0 {
		characterInstanceRec.Health -= totalFleeDamage
		if characterInstanceRec.Health < 0 {
			characterInstanceRec.Health = 0
		}
		l.Info("flee penalty applied: character >%s< takes %d damage (health now %d)",
			characterInstanceRec.ID, totalFleeDamage, characterInstanceRec.Health)
	}

	return nil
}

// CreateNextTurnSheet creates a new turn sheet for a character (implements TurnSheetProcessor interface)
func (p *AdventureGameLocationChoiceProcessor) CreateNextTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("AdventureGameLocationChoiceProcessor/CreateNextTurnSheet")

	l.Info("creating location choice turn sheet for character >%s<", characterInstanceRec.ID)

	// Step 1: Get character's current location instance
	locationInstanceRec, err := p.Domain.GetAdventureGameLocationInstanceRec(characterInstanceRec.AdventureGameLocationInstanceID, nil)
	if err != nil {
		l.Warn("failed to get character's current location >%v<", err)
		return nil, fmt.Errorf("failed to get character's current location: %w", err)
	}

	// Step 2: Get the location definition with name and description
	locationRec, err := p.Domain.GetAdventureGameLocationRec(locationInstanceRec.AdventureGameLocationID, nil)
	if err != nil {
		l.Warn("failed to get location definition >%v<", err)
		return nil, fmt.Errorf("failed to get location definition: %w", err)
	}

	// Step 3: Get all location links FROM this location (outgoing paths)
	locationLinkRecs, err := p.Domain.GetManyAdventureGameLocationLinkRecs(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: adventure_game_record.FieldAdventureGameLocationLinkFromAdventureGameLocationID,
				Val: locationInstanceRec.AdventureGameLocationID,
			},
		},
	})
	if err != nil {
		l.Warn("failed to get location links >%v<", err)
		return nil, fmt.Errorf("failed to get location links: %w", err)
	}

	// Step 4: Get game for game name
	gameRec, err := p.Domain.GetGameRec(gameInstanceRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game >%v<", err)
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Step 5: Get character for account ID
	characterRec, err := p.Domain.GetAdventureGameCharacterRec(characterInstanceRec.AdventureGameCharacterID, nil)
	if err != nil {
		l.Warn("failed to get character >%v<", err)
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	// Step 6: Get account for name
	accountUserRec, err := p.Domain.GetAccountUserRec(characterRec.AccountUserID, nil)
	if err != nil {
		l.Warn("failed to get account >%v<", err)
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Step 7: Build location options from links, evaluating requirements for each link
	locationOptions := make([]turnsheet.LocationOption, 0, len(locationLinkRecs))
	for _, locationLinkRec := range locationLinkRecs {
		// Get the destination location definition
		toLocationRec, err := p.Domain.GetAdventureGameLocationRec(locationLinkRec.ToAdventureGameLocationID, nil)
		if err != nil {
			l.Warn("failed to get destination location >%s< >%v<", locationLinkRec.ToAdventureGameLocationID, err)
			continue
		}

		// Get the location instance ID for the destination
		toLocationInstanceRec, err := p.getLocationInstanceForLocation(gameInstanceRec.ID, locationLinkRec.ToAdventureGameLocationID)
		if err != nil {
			l.Warn("failed to get location instance for location >%s< >%v<", locationLinkRec.ToAdventureGameLocationID, err)
			continue
		}

		// Evaluate requirements: visibility first, then traverse
		isVisible, isTraversable, err := evaluateLinkRequirements(l, p.Domain, gameInstanceRec, characterInstanceRec, locationInstanceRec, locationLinkRec)
		if err != nil {
			l.Warn("failed to evaluate requirements for link >%s< >%v<", locationLinkRec.ID, err)
			continue
		}

		if !isVisible {
			// Link is entirely hidden — omit from turn sheet
			l.Info("link >%s< to >%s< is hidden — omitting", locationLinkRec.Name, toLocationRec.Name)
			continue
		}

		option := turnsheet.LocationOption{
			LocationID:       toLocationInstanceRec.ID,
			LocationLinkName: locationLinkRec.Name,
		}

		if isTraversable {
			option.LocationLinkDescription = locationLinkRec.Description
			l.Info("added accessible location option: >%s< via >%s<", toLocationRec.Name, locationLinkRec.Name)
		} else {
			option.IsLocked = true
			if locationLinkRec.LockedDescription.Valid {
				option.LockedDescription = locationLinkRec.LockedDescription.String
			} else {
				option.LockedDescription = locationLinkRec.Description
			}
			l.Info("added locked location option: >%s< via >%s<", toLocationRec.Name, locationLinkRec.Name)
		}

		locationOptions = append(locationOptions, option)
	}

	// Step 8: Generate turn sheet code for template rendering
	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		l.Warn("failed to generate turn sheet code >%v<", err)
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	// Step 8a: Load background image for this location (falls back to game-level image)
	var backgroundImage *string
	bgImageURL, err := p.Domain.GetAdventureGameLocationChoiceTurnSheetImageDataURL(gameRec.ID, locationInstanceRec.AdventureGameLocationID)
	if err != nil {
		l.Warn("failed to get turn sheet background image >%v<", err)
	} else if bgImageURL != "" {
		backgroundImage = &bgImageURL
		l.Info("loaded background image for location choice turn sheet, length >%d<", len(bgImageURL))
	} else {
		l.Info("no background image found for location choice turn sheet")
	}

	// Step 8b: Load creatures present at this location
	creatures, err := GetAliveCreaturesAtLocation(l, p.Domain, gameInstanceRec.ID, locationInstanceRec.ID)
	if err != nil {
		l.Warn("failed to get creatures at location >%v<", err)
		creatures = nil
	}

	hasAggressiveCreatures := false
	for _, c := range creatures {
		if c.Disposition == adventure_game_record.AdventureGameCreatureDispositionAggressive {
			hasAggressiveCreatures = true
			break
		}
	}

	// Step 8c: Load visible object instances at this location and build LocationObjects
	locationObjects, err := p.buildLocationObjects(l, gameInstanceRec.ID, characterInstanceRec.ID, locationInstanceRec.ID)
	if err != nil {
		l.Warn("failed to build location objects >%v<", err)
		locationObjects = nil
	}

	// Step 9: Read movement, flee, and world events for this sheet. Events are cleared after all processors run.
	displayEvents, err := ReadTurnEventsForCategories(l, p.Domain, characterInstanceRec,
		turnsheet.TurnEventCategoryMovement,
		turnsheet.TurnEventCategoryFlee,
		turnsheet.TurnEventCategoryWorld,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read location events: %w", err)
	}

	// Step 10: Create sheet data with REAL game data
	sheetData := turnsheet.LocationChoiceData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr("adventure"),
			TurnNumber:            convert.Ptr(gameInstanceRec.CurrentTurn),
			AccountName:           convert.Ptr(accountUserRec.Email),
			TurnSheetTitle:        convert.Ptr(locationRec.Name),
			TurnSheetDescription:  convert.Ptr(locationRec.Description),
			TurnSheetInstructions: convert.Ptr(turnsheet.DefaultLocationChoiceInstructions()),
			TurnSheetCode:         convert.Ptr(turnSheetCode),
			BackgroundImage:       backgroundImage,
			TurnEvents:            displayEvents,
		},
		LocationName:           locationRec.Name,
		LocationDescription:    locationRec.Description,
		Creatures:              creatures,
		HasAggressiveCreatures: hasAggressiveCreatures,
		LocationOptions:        locationOptions,
		LocationObjects:        locationObjects,
	}

	sheetDataBytes, err := json.Marshal(sheetData)
	if err != nil {
		l.Warn("failed to marshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to marshal sheet data: %w", err)
	}

	// Step 11: Create turn sheet record
	turnSheet := &game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		AccountID:        accountUserRec.AccountID,
		AccountUserID:    characterRec.AccountUserID,
		TurnNumber:       gameInstanceRec.CurrentTurn,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeLocationChoice,
		SheetOrder:       adventure_game_record.AdventureGameSheetOrderForType(adventure_game_record.AdventureGameTurnSheetTypeLocationChoice),
		SheetData:        json.RawMessage(sheetDataBytes),
		IsCompleted:      false,
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
	turnSheet.GameInstanceID = nullstring.FromString(gameInstanceRec.ID)

	// Create the turn sheet record
	createdTurnSheetRec, err := p.Domain.CreateGameTurnSheetRec(turnSheet)
	if err != nil {
		l.Warn("failed to create turn sheet record >%v<", err)
		return nil, fmt.Errorf("failed to create turn sheet record: %w", err)
	}

	// Link it to the character via AdventureGameTurnSheet
	adventureTurnSheet := &adventure_game_record.AdventureGameTurnSheet{
		GameID:                           gameInstanceRec.GameID,
		AdventureGameCharacterInstanceID: characterInstanceRec.ID,
		GameTurnSheetID:                  createdTurnSheetRec.ID,
	}

	_, err = p.Domain.CreateAdventureGameTurnSheetRec(adventureTurnSheet)
	if err != nil {
		l.Warn("failed to create adventure game turn sheet record >%v<", err)
		return nil, fmt.Errorf("failed to create adventure game turn sheet record: %w", err)
	}

	l.Info("created turn sheet >%s< for character >%s< at location >%s< with %d pathway options",
		createdTurnSheetRec.ID, characterInstanceRec.ID, locationRec.Name, len(locationOptions))

	return createdTurnSheetRec, nil
}
