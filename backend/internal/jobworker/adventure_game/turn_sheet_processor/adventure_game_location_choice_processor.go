package turn_sheet_processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
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
	currentLocationInstanceID := characterInstanceRec.AdventureGameLocationInstanceID.String
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
	characterInstanceRec.AdventureGameLocationInstanceID = sql.NullString{String: chosenLocationID, Valid: true}
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
	locationInstanceRec, err := p.Domain.GetAdventureGameLocationInstanceRec(characterInstanceRec.AdventureGameLocationInstanceID.String, nil)
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
		isVisible, isTraversable, err := p.evaluateLinkRequirements(l, gameInstanceRec, characterInstanceRec, locationInstanceRec, locationLinkRec)
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
	displayEvents := ReadTurnEventsForCategories(l, p.Domain, characterInstanceRec,
		turnsheet.TurnEventCategoryMovement,
		turnsheet.TurnEventCategoryFlee,
		turnsheet.TurnEventCategoryWorld,
	)

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
	turnSheet.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

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

// evaluateLinkRequirements returns (isVisible, isTraversable, error) for a location link.
// isVisible=false means the link must not appear on the sheet at all.
// isVisible=true, isTraversable=false means it appears locked.
// isVisible=true, isTraversable=true means it appears with a radio button.
func (p *AdventureGameLocationChoiceProcessor) evaluateLinkRequirements(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	fromLocationInstanceRec *adventure_game_record.AdventureGameLocationInstance,
	linkRec *adventure_game_record.AdventureGameLocationLink,
) (isVisible bool, isTraversable bool, err error) {
	requirements, err := p.Domain.GetManyAdventureGameLocationLinkRequirementRecs(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: adventure_game_record.FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID,
				Val: linkRec.ID,
			},
		},
	})
	if err != nil {
		return false, false, fmt.Errorf("failed to get link requirements: %w", err)
	}

	// No requirements — link is always visible and traversable
	if len(requirements) == 0 {
		return true, true, nil
	}

	// Evaluate all visible requirements first (AND logic — all must pass)
	for _, req := range requirements {
		if req.Purpose != adventure_game_record.AdventureGameLocationLinkRequirementPurposeVisible {
			continue
		}
		met, err := p.evaluateSingleRequirement(l, gameInstanceRec, characterInstanceRec, fromLocationInstanceRec, req)
		if err != nil {
			return false, false, err
		}
		if !met {
			return false, false, nil // hidden
		}
	}

	// Evaluate all traverse requirements (AND logic)
	for _, req := range requirements {
		if req.Purpose != adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse {
			continue
		}
		met, err := p.evaluateSingleRequirement(l, gameInstanceRec, characterInstanceRec, fromLocationInstanceRec, req)
		if err != nil {
			return true, false, err
		}
		if !met {
			return true, false, nil // visible but locked
		}
	}

	return true, true, nil
}

// evaluateSingleRequirement returns whether a single link requirement is satisfied.
func (p *AdventureGameLocationChoiceProcessor) evaluateSingleRequirement(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	fromLocationInstanceRec *adventure_game_record.AdventureGameLocationInstance,
	req *adventure_game_record.AdventureGameLocationLinkRequirement,
) (bool, error) {
	switch req.Condition {

	// Item conditions
	case adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory:
		return p.characterHasItemInInventory(characterInstanceRec.ID, req.AdventureGameItemID.String, req.Quantity)

	case adventure_game_record.AdventureGameLocationLinkRequirementConditionEquipped:
		return p.characterHasItemEquipped(characterInstanceRec.ID, req.AdventureGameItemID.String)

	// Creature conditions
	case adventure_game_record.AdventureGameLocationLinkRequirementConditionDeadAtLocation:
		return p.creatureDeadAtLocation(gameInstanceRec.ID, fromLocationInstanceRec.ID, req.AdventureGameCreatureID.String, req.Quantity)

	case adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveAtLocation:
		return p.noCreaturesAliveAtLocation(gameInstanceRec.ID, fromLocationInstanceRec.ID, req.AdventureGameCreatureID.String)

	case adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveInGame:
		return p.noCreaturesAliveInGame(gameInstanceRec.ID, req.AdventureGameCreatureID.String)

	default:
		l.Warn("unknown requirement condition >%s<", req.Condition)
		return false, fmt.Errorf("unknown requirement condition: %s", req.Condition)
	}
}

func (p *AdventureGameLocationChoiceProcessor) characterHasItemInInventory(characterInstanceID, itemID string, quantity int) (bool, error) {
	itemInstances, err := p.Domain.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceID},
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: itemID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query item instances: %w", err)
	}
	count := 0
	for _, inst := range itemInstances {
		if !inst.IsUsed {
			count++
		}
	}
	return count >= quantity, nil
}

func (p *AdventureGameLocationChoiceProcessor) characterHasItemEquipped(characterInstanceID, itemID string) (bool, error) {
	itemInstances, err := p.Domain.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceID},
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: itemID},
			{Col: adventure_game_record.FieldAdventureGameItemInstanceIsEquipped, Val: true},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query equipped item instances: %w", err)
	}
	return len(itemInstances) > 0, nil
}

func (p *AdventureGameLocationChoiceProcessor) creatureDeadAtLocation(gameInstanceID, fromLocationInstanceID, creatureID string, quantity int) (bool, error) {
	allInstances, err := p.Domain.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameCreatureID, Val: creatureID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, Val: fromLocationInstanceID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query creature instances: %w", err)
	}
	deadCount := 0
	for _, inst := range allInstances {
		if inst.Health <= 0 {
			deadCount++
		}
	}
	return deadCount >= quantity, nil
}

func (p *AdventureGameLocationChoiceProcessor) noCreaturesAliveAtLocation(gameInstanceID, fromLocationInstanceID, creatureID string) (bool, error) {
	allInstances, err := p.Domain.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameCreatureID, Val: creatureID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, Val: fromLocationInstanceID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query creature instances: %w", err)
	}
	for _, inst := range allInstances {
		if inst.Health > 0 {
			return false, nil
		}
	}
	return true, nil
}

func (p *AdventureGameLocationChoiceProcessor) noCreaturesAliveInGame(gameInstanceID, creatureID string) (bool, error) {
	allInstances, err := p.Domain.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameCreatureID, Val: creatureID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query creature instances: %w", err)
	}
	for _, inst := range allInstances {
		if inst.Health > 0 {
			return false, nil
		}
	}
	return true, nil
}

// buildLocationObjects queries visible object instances at a location and builds the LocationObjectOption slice.
func (p *AdventureGameLocationChoiceProcessor) buildLocationObjects(
	l logger.Logger,
	gameInstanceID, characterInstanceID, locationInstanceID string,
) ([]turnsheet.LocationObjectOption, error) {
	l = l.WithFunctionContext("buildLocationObjects")

	objectInstances, err := p.Domain.GetManyAdventureGameLocationObjectInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceAdventureGameLocationInstanceID, Val: locationInstanceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object instances: %w", err)
	}

	var result []turnsheet.LocationObjectOption
	for _, inst := range objectInstances {
		if !inst.IsVisible {
			continue
		}

		objectDef, err := p.Domain.GetAdventureGameLocationObjectRec(inst.AdventureGameLocationObjectID, nil)
		if err != nil {
			l.Warn("failed to get object definition >%s< >%v<", inst.AdventureGameLocationObjectID, err)
			continue
		}

		// Load effects that match current state (required_state = current_state OR required_state IS NULL)
		allEffects, err := p.Domain.GetManyAdventureGameLocationObjectEffectRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameLocationObjectEffectAdventureGameLocationObjectID, Val: objectDef.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get object effects >%v<", err)
			continue
		}

		// Collect unique action types available in the current state
		actionMap := map[string]*turnsheet.LocationObjectActionOption{}
		for _, effect := range allEffects {
			if effect.RequiredState.Valid && effect.RequiredState.String != inst.CurrentState {
				continue
			}
			at := effect.ActionType
			if _, exists := actionMap[at]; exists {
				continue
			}

			action := &turnsheet.LocationObjectActionOption{
				ActionType:      at,
				IsAvailable:     true,
				HasRequiredItem: true,
			}

			if effect.RequiredAdventureGameItemID.Valid && effect.RequiredAdventureGameItemID.String != "" {
				hasItem, itemName, err := p.characterHasItemForEffect(characterInstanceID, effect.RequiredAdventureGameItemID.String)
				if err != nil {
					l.Warn("failed to check required item >%v<", err)
				}
				action.RequiredItemName = itemName
				action.HasRequiredItem = hasItem
				action.IsAvailable = hasItem
			}

			actionMap[at] = action
		}

		// Convert map to ordered slice
		var actions []turnsheet.LocationObjectActionOption
		for _, action := range actionMap {
			actions = append(actions, *action)
		}

		result = append(result, turnsheet.LocationObjectOption{
			ObjectInstanceID: inst.ID,
			Name:             objectDef.Name,
			Description:      objectDef.Description,
			CurrentState:     inst.CurrentState,
			Actions:          actions,
		})
	}

	return result, nil
}

// characterHasItemForEffect returns (hasItem, itemName, error).
func (p *AdventureGameLocationChoiceProcessor) characterHasItemForEffect(characterInstanceID, itemID string) (bool, string, error) {
	itemDef, err := p.Domain.GetAdventureGameItemRec(itemID, nil)
	if err != nil {
		return false, "", fmt.Errorf("failed to get item definition: %w", err)
	}

	itemInstances, err := p.Domain.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceID},
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: itemID},
		},
	})
	if err != nil {
		return false, itemDef.Name, fmt.Errorf("failed to query item instances: %w", err)
	}

	for _, inst := range itemInstances {
		if !inst.IsUsed {
			return true, itemDef.Name, nil
		}
	}
	return false, itemDef.Name, nil
}

// processObjectChoice parses "{instance_id}:{action_type}", loads all matching effects, and applies them atomically.
func (p *AdventureGameLocationChoiceProcessor) processObjectChoice(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	objectChoice string,
) error {
	l = l.WithFunctionContext("processObjectChoice")

	parts := strings.SplitN(objectChoice, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid object_choice format: %q", objectChoice)
	}
	instanceID := parts[0]
	actionType := parts[1]

	// Load object instance
	objectInstance, err := p.Domain.GetAdventureGameLocationObjectInstanceRec(instanceID, nil)
	if err != nil {
		return fmt.Errorf("failed to get object instance >%s<: %w", instanceID, err)
	}

	// Load all effects for this object
	allEffects, err := p.Domain.GetManyAdventureGameLocationObjectEffectRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectEffectAdventureGameLocationObjectID, Val: objectInstance.AdventureGameLocationObjectID},
			{Col: adventure_game_record.FieldAdventureGameLocationObjectEffectActionType, Val: actionType},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get effects for object: %w", err)
	}

	// Filter effects matching current state
	var matchingEffects []*adventure_game_record.AdventureGameLocationObjectEffect
	for _, effect := range allEffects {
		if effect.RequiredState.Valid && effect.RequiredState.String != objectInstance.CurrentState {
			continue
		}
		matchingEffects = append(matchingEffects, effect)
	}

	if len(matchingEffects) == 0 {
		l.Info("no matching effects for object >%s< action >%s< state >%s<", instanceID, actionType, objectInstance.CurrentState)
		return nil
	}

	// Validate required items — any effect's required item blocks all effects
	for _, effect := range matchingEffects {
		if !effect.RequiredAdventureGameItemID.Valid || effect.RequiredAdventureGameItemID.String == "" {
			continue
		}
		hasItem, _, err := p.characterHasItemForEffect(characterInstanceRec.ID, effect.RequiredAdventureGameItemID.String)
		if err != nil {
			return fmt.Errorf("failed to check required item: %w", err)
		}
		if !hasItem {
			l.Info("character does not have required item for object interaction")
			return fmt.Errorf("required item not in inventory")
		}
	}

	// Apply all matching effects atomically
	var resultDescriptions []string
	for _, effect := range matchingEffects {
		desc, err := p.applyObjectEffect(l, gameInstanceRec, characterInstanceRec, objectInstance, effect)
		if err != nil {
			l.Warn("failed to apply effect >%s< >%v<", effect.ID, err)
			return fmt.Errorf("failed to apply object effect: %w", err)
		}
		if desc != "" {
			resultDescriptions = append(resultDescriptions, desc)
		}
	}

	// Append a combined world event with all result descriptions
	if len(resultDescriptions) > 0 {
		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryWorld,
			Icon:     turnsheet.TurnEventIconWorld,
			Message:  strings.Join(resultDescriptions, " "),
		})
		if _, saveErr := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec); saveErr != nil {
			l.Warn("failed to save object interaction events >%v<", saveErr)
		}
	}

	return nil
}

// applyObjectEffect applies a single effect and returns its result_description.
func (p *AdventureGameLocationChoiceProcessor) applyObjectEffect(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	objectInstance *adventure_game_record.AdventureGameLocationObjectInstance,
	effect *adventure_game_record.AdventureGameLocationObjectEffect,
) (string, error) {
	l = l.WithFunctionContext("applyObjectEffect")
	l.Info("applying effect >%s< type >%s<", effect.ID, effect.EffectType)

	switch effect.EffectType {
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeNothing:
		// No state change — just return the description

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState:
		if effect.ResultState.Valid && effect.ResultState.String != "" {
			objectInstance.CurrentState = effect.ResultState.String
			updatedInst, err := p.Domain.UpdateAdventureGameLocationObjectInstanceRec(objectInstance)
			if err != nil {
				return "", fmt.Errorf("failed to update object instance state: %w", err)
			}
			*objectInstance = *updatedInst
			l.Info("object instance >%s< state changed to >%s<", objectInstance.ID, objectInstance.CurrentState)
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeObjectState:
		if effect.ResultAdventureGameLocationObjectID.Valid && effect.ResultState.Valid {
			targetInstances, err := p.Domain.GetManyAdventureGameLocationObjectInstanceRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceGameInstanceID, Val: gameInstanceRec.ID},
					{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceAdventureGameLocationObjectID, Val: effect.ResultAdventureGameLocationObjectID.String},
				},
			})
			if err != nil {
				return "", fmt.Errorf("failed to get target object instances: %w", err)
			}
			for _, targetInst := range targetInstances {
				targetInst.CurrentState = effect.ResultState.String
				if _, err := p.Domain.UpdateAdventureGameLocationObjectInstanceRec(targetInst); err != nil {
					l.Warn("failed to update target object instance state >%v<", err)
				}
				l.Info("target object instance >%s< state changed to >%s<", targetInst.ID, targetInst.CurrentState)
			}
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRevealObject:
		if effect.ResultAdventureGameLocationObjectID.Valid {
			targetInstances, err := p.Domain.GetManyAdventureGameLocationObjectInstanceRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceGameInstanceID, Val: gameInstanceRec.ID},
					{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceAdventureGameLocationObjectID, Val: effect.ResultAdventureGameLocationObjectID.String},
				},
			})
			if err != nil {
				return "", fmt.Errorf("failed to get target object instances: %w", err)
			}
			for _, targetInst := range targetInstances {
				targetInst.IsVisible = true
				if _, err := p.Domain.UpdateAdventureGameLocationObjectInstanceRec(targetInst); err != nil {
					l.Warn("failed to reveal target object instance >%v<", err)
				}
			}
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHideObject:
		if effect.ResultAdventureGameLocationObjectID.Valid {
			targetInstances, err := p.Domain.GetManyAdventureGameLocationObjectInstanceRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceGameInstanceID, Val: gameInstanceRec.ID},
					{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceAdventureGameLocationObjectID, Val: effect.ResultAdventureGameLocationObjectID.String},
				},
			})
			if err != nil {
				return "", fmt.Errorf("failed to get target object instances: %w", err)
			}
			for _, targetInst := range targetInstances {
				targetInst.IsVisible = false
				if _, err := p.Domain.UpdateAdventureGameLocationObjectInstanceRec(targetInst); err != nil {
					l.Warn("failed to hide target object instance >%v<", err)
				}
			}
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem:
		if effect.ResultAdventureGameItemID.Valid && effect.ResultAdventureGameItemID.String != "" {
			itemInstance := &adventure_game_record.AdventureGameItemInstance{
				GameID:                         gameInstanceRec.GameID,
				GameInstanceID:                 gameInstanceRec.ID,
				AdventureGameItemID:            effect.ResultAdventureGameItemID.String,
				AdventureGameCharacterInstanceID: sql.NullString{String: characterInstanceRec.ID, Valid: true},
			}
			if _, err := p.Domain.CreateAdventureGameItemInstanceRec(itemInstance); err != nil {
				return "", fmt.Errorf("failed to give item: %w", err)
			}
			l.Info("gave item >%s< to character >%s<", effect.ResultAdventureGameItemID.String, characterInstanceRec.ID)
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveItem:
		if effect.ResultAdventureGameItemID.Valid && effect.ResultAdventureGameItemID.String != "" {
			itemInstances, err := p.Domain.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceRec.ID},
					{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: effect.ResultAdventureGameItemID.String},
				},
			})
			if err != nil {
				return "", fmt.Errorf("failed to query item instances: %w", err)
			}
			for _, inst := range itemInstances {
				if !inst.IsUsed {
					if err := p.Domain.DeleteAdventureGameItemInstanceRec(inst.ID); err != nil {
						l.Warn("failed to remove item instance >%v<", err)
					}
					break
				}
			}
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeOpenLink:
		if effect.ResultAdventureGameLocationLinkID.Valid && effect.ResultAdventureGameLocationLinkID.String != "" {
			// Remove all requirements for this link (making it traversable)
			requirements, err := p.Domain.GetManyAdventureGameLocationLinkRequirementRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: adventure_game_record.FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID, Val: effect.ResultAdventureGameLocationLinkID.String},
					{Col: adventure_game_record.FieldAdventureGameLocationLinkRequirementPurpose, Val: adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse},
				},
			})
			if err != nil {
				return "", fmt.Errorf("failed to get link requirements: %w", err)
			}
			for _, req := range requirements {
				if err := p.Domain.DeleteAdventureGameLocationLinkRequirementRec(req.ID); err != nil {
					l.Warn("failed to remove link requirement >%v<", err)
				}
			}
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeDamage:
		if effect.ResultValueMin.Valid && effect.ResultValueMax.Valid {
			minDmg := int(effect.ResultValueMin.Int32)
			maxDmg := int(effect.ResultValueMax.Int32)
			damage := minDmg
			if maxDmg > minDmg {
				damage = minDmg + rand.Intn(maxDmg-minDmg+1)
			}
			characterInstanceRec.Health -= damage
			if characterInstanceRec.Health < 0 {
				characterInstanceRec.Health = 0
			}
			l.Info("object dealt %d damage to character >%s< (health now %d)", damage, characterInstanceRec.ID, characterInstanceRec.Health)
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHeal:
		if effect.ResultValueMin.Valid && effect.ResultValueMax.Valid {
			minHeal := int(effect.ResultValueMin.Int32)
			maxHeal := int(effect.ResultValueMax.Int32)
			heal := minHeal
			if maxHeal > minHeal {
				heal = minHeal + rand.Intn(maxHeal-minHeal+1)
			}
			characterInstanceRec.Health += heal
			l.Info("object healed %d for character >%s< (health now %d)", heal, characterInstanceRec.ID, characterInstanceRec.Health)
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeSummonCreature:
		if effect.ResultAdventureGameCreatureID.Valid && effect.ResultAdventureGameCreatureID.String != "" {
			creatureInstance := &adventure_game_record.AdventureGameCreatureInstance{
				GameID:                          gameInstanceRec.GameID,
				GameInstanceID:                  gameInstanceRec.ID,
				AdventureGameCreatureID:         effect.ResultAdventureGameCreatureID.String,
				AdventureGameLocationInstanceID: objectInstance.AdventureGameLocationInstanceID,
				Health:                          100,
			}
			if _, err := p.Domain.CreateAdventureGameCreatureInstanceRec(creatureInstance); err != nil {
				return "", fmt.Errorf("failed to summon creature: %w", err)
			}
			l.Info("summoned creature >%s< at location >%s<", effect.ResultAdventureGameCreatureID.String, objectInstance.AdventureGameLocationInstanceID)
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeTeleport:
		if effect.ResultAdventureGameLocationID.Valid && effect.ResultAdventureGameLocationID.String != "" {
			destLocInst, err := p.getLocationInstanceForLocation(gameInstanceRec.ID, effect.ResultAdventureGameLocationID.String)
			if err != nil {
				return "", fmt.Errorf("failed to get destination location instance: %w", err)
			}
			characterInstanceRec.AdventureGameLocationInstanceID = sql.NullString{String: destLocInst.ID, Valid: true}
			updatedChar, err := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec)
			if err != nil {
				return "", fmt.Errorf("failed to teleport character: %w", err)
			}
			*characterInstanceRec = *updatedChar
			l.Info("teleported character >%s< to location instance >%s<", characterInstanceRec.ID, destLocInst.ID)
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveObject:
		if err := p.Domain.DeleteAdventureGameLocationObjectInstanceRec(objectInstance.ID); err != nil {
			return "", fmt.Errorf("failed to remove object instance: %w", err)
		}
		l.Info("removed object instance >%s<", objectInstance.ID)
	}

	return effect.ResultDescription, nil
}

// getLocationInstanceForLocation finds the location instance ID for a given game and location
func (p *AdventureGameLocationChoiceProcessor) getLocationInstanceForLocation(gameInstanceID, locationID string) (*adventure_game_record.AdventureGameLocationInstance, error) {
	locationInstances, err := p.Domain.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceAdventureGameLocationID, Val: locationID},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(locationInstances) == 0 {
		return nil, fmt.Errorf("no location instance found for location >%s<", locationID)
	}
	return locationInstances[0], nil
}
