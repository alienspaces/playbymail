package turn_sheet_processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/convert"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
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
	return adventure_game_record.AdventureSheetTypeLocationChoice
}

// ProcessTurnSheetResponse processes a single turn sheet response (implements TurnSheetProcessor interface)
func (p *AdventureGameLocationChoiceProcessor) ProcessTurnSheetResponse(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("AdventureGameLocationChoiceProcessor/ProcessTurnSheetResponse")

	l.Info("processing location choice for turn sheet >%s< for character >%s<", turnSheet.ID, characterInstanceRec.ID)

	// Verify this is a location choice sheet
	if turnSheet.SheetType != adventure_game_record.AdventureSheetTypeLocationChoice {
		l.Warn("expected location choice sheet type, got >%s<", turnSheet.SheetType)
		return fmt.Errorf("invalid sheet type: expected %s, got %s", adventure_game_record.AdventureSheetTypeLocationChoice, turnSheet.SheetType)
	}

	// Step 1: Parse the player's location choice from ScannedData
	var scanData turn_sheet.LocationChoiceScanData
	if err := json.Unmarshal(turnSheet.ScannedData, &scanData); err != nil {
		l.Warn("failed to unmarshal scanned data >%v<", err)
		return fmt.Errorf("failed to parse scanned data: %w", err)
	}

	if len(scanData.Choices) == 0 {
		l.Warn("no location choice found in scanned data")
		return fmt.Errorf("no location choice found")
	}

	// For now, we'll use the first choice
	// In the future, we might support multiple simultaneous choices
	chosenLocationID := scanData.Choices[0]
	l.Info("player chose location >%s<", chosenLocationID)

	// Step 2: Parse SheetData to get the original location options and validate
	var sheetData turn_sheet.LocationChoiceData
	if err := json.Unmarshal(turnSheet.SheetData, &sheetData); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return fmt.Errorf("failed to parse sheet data: %w", err)
	}

	// Validate the choice is one of the available options
	isValidChoice := false
	var chosenLocationOption turn_sheet.LocationOption
	for _, option := range sheetData.LocationOptions {
		if option.LocationID == chosenLocationID {
			isValidChoice = true
			chosenLocationOption = option
			break
		}
	}

	if !isValidChoice {
		l.Warn("invalid location choice >%s< not in available options", chosenLocationID)
		return fmt.Errorf("invalid location choice: %s is not an available option", chosenLocationID)
	}

	// Step 3: Update character's location in the game state
	characterInstanceRec.AdventureGameLocationInstanceID = chosenLocationID
	updatedCharacter, err := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec)
	if err != nil {
		l.Warn("failed to update character location >%v<", err)
		return fmt.Errorf("failed to update character location: %w", err)
	}

	l.Info("successfully updated character >%s< to location >%s< via pathway >%s<",
		updatedCharacter.ID, chosenLocationID, chosenLocationOption.LocationLinkName)

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
	accountRec, err := p.Domain.GetAccountRec(characterRec.AccountID, nil)
	if err != nil {
		l.Warn("failed to get account >%v<", err)
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Step 7: Build location options from links (these are the pathways!)
	locationOptions := make([]turn_sheet.LocationOption, 0, len(locationLinkRecs))
	for _, locationLinkRec := range locationLinkRecs {
		// Get the destination location's name
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

		// The link.Name is the pathway name (e.g., "The River of Blood")
		// The link.Description is the pathway description (e.g., "A river of blood flows to the next location")
		locationOptions = append(locationOptions, turn_sheet.LocationOption{
			LocationID:              toLocationInstanceRec.ID,
			LocationLinkName:        locationLinkRec.Name,
			LocationLinkDescription: locationLinkRec.Description,
		})

		l.Info("added location option: >%s< via >%s< (%s)", toLocationRec.Name, locationLinkRec.Name, locationLinkRec.Description)
	}

	// Step 8: Create sheet data with REAL game data
	sheetData := turn_sheet.LocationChoiceData{
		TurnSheetTemplateData: turn_sheet.TurnSheetTemplateData{
			GameName:    convert.Ptr(gameRec.Name),
			GameType:    convert.Ptr("adventure"),
			TurnNumber:  convert.Ptr(gameInstanceRec.CurrentTurn + 1),
			AccountName: convert.Ptr(accountRec.Email),
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
		AccountID:        characterRec.AccountID,
		TurnNumber:       gameInstanceRec.CurrentTurn + 1,
		SheetType:        adventure_game_record.AdventureSheetTypeLocationChoice,
		SheetOrder:       1,
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
