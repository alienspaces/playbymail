package adventure_game

import (
	"context"
	"encoding/json"
	"fmt"

	corerecord "gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
)

// AdventureGameProcessSubscriptionProcessor processes process subscription for adventure games
type AdventureGameProcessSubscriptionProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewAdventureGameJoinGameProcessor creates a new adventure game join game processor
func NewAdventureGameProcessSubscriptionProcessor(l logger.Logger, d *domain.Domain) (*AdventureGameProcessSubscriptionProcessor, error) {
	l = l.WithFunctionContext("NewAdventureGameProcessSubscriptionProcessor")

	return &AdventureGameProcessSubscriptionProcessor{
		Logger: l,
		Domain: d,
	}, nil
}

// ProcessProcessSubscription processes a join game turn sheet and creates the necessary
// game entities (game instance, character, character instance, etc.)
func (p *AdventureGameProcessSubscriptionProcessor) ProcessProcessSubscription(ctx context.Context, subscriptionRec *game_record.GameSubscription, turnSheetRec *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("AdventureGameProcessSubscriptionProcessor/ProcessProcessSubscription")

	l.Info("processing join game turn sheet for subscription ID >%s< turn sheet ID >%s<", subscriptionRec.ID, turnSheetRec.ID)

	// Parse the scanned data to get character name
	var scanData turn_sheet.AdventureGameJoinGameScanData
	if err := json.Unmarshal(turnSheetRec.ScannedData, &scanData); err != nil {
		l.Warn("failed to unmarshal join game scan data >%v<", err)
		return fmt.Errorf("failed to parse join game scan data: %w", err)
	}

	if scanData.CharacterName == "" {
		l.Warn("character name is empty in scan data")
		return fmt.Errorf("character name is required")
	}

	// Get or create game instance for this game
	gameInstanceRec, err := p.getOrCreateGameInstance(subscriptionRec.GameID)
	if err != nil {
		l.Warn("failed to get or create game instance >%v<", err)
		return fmt.Errorf("failed to get or create game instance: %w", err)
	}

	// Create or get adventure game character
	characterRec, err := p.getOrCreateAdventureGameCharacter(subscriptionRec.GameID, subscriptionRec.AccountID, scanData.CharacterName)
	if err != nil {
		l.Warn("failed to get or create adventure game character >%v<", err)
		return fmt.Errorf("failed to get or create adventure game character: %w", err)
	}

	// Check if character instance already exists
	existingCharacterInstanceRecs, err := p.Domain.GetManyAdventureGameCharacterInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCharacterInstanceGameInstanceID, Val: gameInstanceRec.ID},
			{Col: adventure_game_record.FieldAdventureGameCharacterInstanceAdventureGameCharacterID, Val: characterRec.ID},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to check for existing character instance >%v<", err)
		return fmt.Errorf("failed to check for existing character instance: %w", err)
	}

	if len(existingCharacterInstanceRecs) > 0 {
		l.Info("character instance already exists for character ID >%s< game instance ID >%s<", characterRec.ID, gameInstanceRec.ID)
		return nil
	}

	// Get starting location for the character
	startingLocationInstanceID, err := p.getStartingLocationInstance(subscriptionRec.GameID, gameInstanceRec.ID)
	if err != nil {
		l.Warn("failed to get starting location instance >%v<", err)
		return fmt.Errorf("failed to get starting location instance: %w", err)
	}

	// Create character instance
	characterInstanceRec := &adventure_game_record.AdventureGameCharacterInstance{
		GameID:                          subscriptionRec.GameID,
		GameInstanceID:                  gameInstanceRec.ID,
		AdventureGameCharacterID:        characterRec.ID,
		AdventureGameLocationInstanceID: startingLocationInstanceID,
		Health:                          100,
	}

	characterInstanceRec, err = p.Domain.CreateAdventureGameCharacterInstanceRec(characterInstanceRec)
	if err != nil {
		l.Warn("failed to create character instance >%v<", err)
		return fmt.Errorf("failed to create character instance: %w", err)
	}

	l.Info("created character instance ID >%s< for character ID >%s< game instance ID >%s<", characterInstanceRec.ID, characterRec.ID, gameInstanceRec.ID)

	// Assign starting items to the character instance
	err = p.Domain.AssignStartingItemsToCharacterInstance(characterInstanceRec)
	if err != nil {
		l.Warn("failed to assign starting items to character instance >%s< >%v<", characterInstanceRec.ID, err)
		// Don't fail the entire process if starting items assignment fails
		// Log the error but continue
		l.Warn("continuing despite starting items assignment failure")
	} else {
		l.Info("assigned starting items to character instance >%s<", characterInstanceRec.ID)
	}

	return nil
}

// getOrCreateGameInstance gets the first active game instance for a game, or creates a new one
func (p *AdventureGameProcessSubscriptionProcessor) getOrCreateGameInstance(gameID string) (*game_record.GameInstance, error) {
	l := p.Logger.WithFunctionContext("getOrCreateGameInstance")

	// Try to get an existing active game instance that has started
	gameInstanceRecs, err := p.Domain.GetManyGameInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameInstanceGameID, Val: gameID},
			{Col: game_record.FieldGameInstanceStatus, Val: game_record.GameInstanceStatusStarted},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get game instances by game ID >%s< status >%s< >%v<", gameID, game_record.GameInstanceStatusStarted, err)
		return nil, err
	}

	if len(gameInstanceRecs) > 0 {
		l.Info("found existing game instance ID >%s< for game ID >%s<", gameInstanceRecs[0].ID, gameID)
		return gameInstanceRecs[0], nil
	}

	// Also check for created status instances that have not started
	createdInstanceRecs, err := p.Domain.GetManyGameInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameInstanceGameID, Val: gameID},
			{Col: game_record.FieldGameInstanceStatus, Val: game_record.GameInstanceStatusCreated},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get game instances by game ID >%s< status >%s< >%v<", gameID, game_record.GameInstanceStatusCreated, err)
		return nil, err
	}

	if len(createdInstanceRecs) > 0 {
		l.Info("found existing created game instance ID >%s< for game ID >%s<", createdInstanceRecs[0].ID, gameID)
		return createdInstanceRecs[0], nil
	}

	// No existing instance found, create a new one
	l.Info("creating new game instance for game ID >%s<", gameID)
	gameInstanceRec := &game_record.GameInstance{
		GameID: gameID,
		Status: game_record.GameInstanceStatusCreated,
	}

	gameInstanceRec, err = p.Domain.CreateGameInstanceRec(gameInstanceRec)
	if err != nil {
		l.Warn("failed to create game instance >%v<", err)
		return nil, err
	}

	l.Info("created new game instance ID >%s< for game ID >%s<", gameInstanceRec.ID, gameID)

	return gameInstanceRec, nil
}

// getOrCreateAdventureGameCharacter gets or creates an adventure game character
func (p *AdventureGameProcessSubscriptionProcessor) getOrCreateAdventureGameCharacter(gameID, accountID, characterName string) (*adventure_game_record.AdventureGameCharacter, error) {
	l := p.Logger.WithFunctionContext("getOrCreateAdventureGameCharacter")

	// Check if character already exists
	characterRecs, err := p.Domain.GetManyAdventureGameCharacterRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCharacterGameID, Val: gameID},
			{Col: adventure_game_record.FieldAdventureGameCharacterAccountID, Val: accountID},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(characterRecs) > 0 {
		l.Info("found existing character ID >%s< for game ID >%s< account ID >%s<", characterRecs[0].ID, gameID, accountID)
		return characterRecs[0], nil
	}

	// Create new character
	l.Info("creating new character for game ID >%s< account ID >%s< name >%s<", gameID, accountID, characterName)
	characterRec := &adventure_game_record.AdventureGameCharacter{
		GameID:    gameID,
		AccountID: accountID,
		Name:      characterName,
	}

	characterRec, err = p.Domain.CreateAdventureGameCharacterRec(characterRec)
	if err != nil {
		return nil, err
	}

	l.Info("created new character ID >%s< for game ID >%s< account ID >%s<", characterRec.ID, gameID, accountID)

	return characterRec, nil
}

// getStartingLocationInstance gets the starting location instance for a game instance
// It finds starting locations for the game and then finds the corresponding location instance
func (p *AdventureGameProcessSubscriptionProcessor) getStartingLocationInstance(gameID, gameInstanceID string) (string, error) {
	l := p.Logger.WithFunctionContext("getStartingLocationInstance")

	// Get starting locations for this game
	startingLocationRecs, err := p.Domain.GetManyAdventureGameLocationRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID},
			{Col: adventure_game_record.FieldAdventureGameLocationIsStartingLocation, Val: true},
		},
		Limit: 1,
		OrderBy: []coresql.OrderBy{
			{Col: corerecord.FieldCreatedAt, Direction: coresql.OrderDirectionASC},
		},
	})
	if err != nil {
		l.Warn("failed to get starting locations for game ID >%s< >%v<", gameID, err)
		return "", err
	}

	if len(startingLocationRecs) == 0 {
		l.Warn("no starting locations found for game ID >%s<", gameID)
		// Return empty string - location is nullable
		return "", nil
	}

	startingLocationID := startingLocationRecs[0].ID
	l.Info("found starting location ID >%s< for game ID >%s<", startingLocationID, gameID)

	// Find the location instance for this starting location in the game instance
	locationInstanceRecs, err := p.Domain.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceAdventureGameLocationID, Val: startingLocationID},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get location instance for starting location ID >%s< game instance ID >%s< >%v<", startingLocationID, gameInstanceID, err)
		return "", err
	}

	if len(locationInstanceRecs) == 0 {
		l.Warn("no location instance found for starting location ID >%s< in game instance ID >%s<", startingLocationID, gameInstanceID)
		// Return empty string - location is nullable
		return "", nil
	}

	l.Info("found starting location instance ID >%s< for game instance ID >%s<", locationInstanceRecs[0].ID, gameInstanceID)

	return locationInstanceRecs[0].ID, nil
}
