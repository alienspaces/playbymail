package turn_sheet_processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

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

	// Step 3: Update character's location
	characterInstanceRec.AdventureGameLocationInstanceID = chosenLocationID
	characterInstanceRec, err := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec)
	if err != nil {
		l.Warn("failed to update character location >%v<", err)
		return fmt.Errorf("failed to update character location: %w", err)
	}

	l.Info("successfully updated character >%s< to location >%s< via pathway >%s<", characterInstanceRec.ID, chosenLocationID, chosenLocationOption.LocationLinkName)

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

	// Step 9: Create sheet data with REAL game data
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
		},
		LocationName:        locationRec.Name,
		LocationDescription: locationRec.Description,
		LocationOptions:     locationOptions,
	}

	sheetDataBytes, err := json.Marshal(sheetData)
	if err != nil {
		l.Warn("failed to marshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to marshal sheet data: %w", err)
	}

	// Step 9: Create turn sheet record
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
		return p.characterHasItemInInventory(l, characterInstanceRec.ID, req.AdventureGameItemID.String, req.Quantity)

	case adventure_game_record.AdventureGameLocationLinkRequirementConditionEquipped:
		return p.characterHasItemEquipped(l, characterInstanceRec.ID, req.AdventureGameItemID.String)

	// Creature conditions
	case adventure_game_record.AdventureGameLocationLinkRequirementConditionDeadAtLocation:
		return p.creatureDeadAtLocation(l, gameInstanceRec.ID, fromLocationInstanceRec.ID, req.AdventureGameCreatureID.String, req.Quantity)

	case adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveAtLocation:
		return p.noCreaturesAliveAtLocation(l, gameInstanceRec.ID, fromLocationInstanceRec.ID, req.AdventureGameCreatureID.String)

	case adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveInGame:
		return p.noCreaturesAliveInGame(l, gameInstanceRec.ID, req.AdventureGameCreatureID.String)

	default:
		l.Warn("unknown requirement condition >%s<", req.Condition)
		return false, fmt.Errorf("unknown requirement condition: %s", req.Condition)
	}
}

func (p *AdventureGameLocationChoiceProcessor) characterHasItemInInventory(l logger.Logger, characterInstanceID, itemID string, quantity int) (bool, error) {
	itemInstances, err := p.Domain.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: "adventure_game_character_instance_id", Val: characterInstanceID},
			{Col: "adventure_game_item_id", Val: itemID},
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

func (p *AdventureGameLocationChoiceProcessor) characterHasItemEquipped(l logger.Logger, characterInstanceID, itemID string) (bool, error) {
	itemInstances, err := p.Domain.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: "adventure_game_character_instance_id", Val: characterInstanceID},
			{Col: "adventure_game_item_id", Val: itemID},
			{Col: "is_equipped", Val: true},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query equipped item instances: %w", err)
	}
	return len(itemInstances) > 0, nil
}

func (p *AdventureGameLocationChoiceProcessor) creatureDeadAtLocation(l logger.Logger, gameInstanceID, fromLocationInstanceID, creatureID string, quantity int) (bool, error) {
	allInstances, err := p.Domain.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: "game_instance_id", Val: gameInstanceID},
			{Col: "adventure_game_creature_id", Val: creatureID},
			{Col: "adventure_game_location_instance_id", Val: fromLocationInstanceID},
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

func (p *AdventureGameLocationChoiceProcessor) noCreaturesAliveAtLocation(l logger.Logger, gameInstanceID, fromLocationInstanceID, creatureID string) (bool, error) {
	allInstances, err := p.Domain.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: "game_instance_id", Val: gameInstanceID},
			{Col: "adventure_game_creature_id", Val: creatureID},
			{Col: "adventure_game_location_instance_id", Val: fromLocationInstanceID},
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

func (p *AdventureGameLocationChoiceProcessor) noCreaturesAliveInGame(l logger.Logger, gameInstanceID, creatureID string) (bool, error) {
	allInstances, err := p.Domain.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: "game_instance_id", Val: gameInstanceID},
			{Col: "adventure_game_creature_id", Val: creatureID},
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

// getLocationInstanceForLocation finds the location instance ID for a given game and location
func (p *AdventureGameLocationChoiceProcessor) getLocationInstanceForLocation(gameInstanceID, locationID string) (*adventure_game_record.AdventureGameLocationInstance, error) {
	locationInstances, err := p.Domain.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: "game_instance_id",
				Val: gameInstanceID,
			},
			{
				Col: "adventure_game_location_id",
				Val: locationID,
			},
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
