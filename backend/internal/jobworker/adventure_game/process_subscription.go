package adventure_game

import (
	"context"
	"encoding/json"
	"fmt"

	corerecord "gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
)

// AdventureGameSubscriptionProcessingProcessor processes game subscription processing for adventure games
type AdventureGameSubscriptionProcessingProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewAdventureGameSubscriptionProcessingProcessor creates a new adventure game subscription processing processor
func NewAdventureGameSubscriptionProcessingProcessor(l logger.Logger, d *domain.Domain) (*AdventureGameSubscriptionProcessingProcessor, error) {
	l = l.WithFunctionContext("NewAdventureGameSubscriptionProcessingProcessor")

	return &AdventureGameSubscriptionProcessingProcessor{
		Logger: l,
		Domain: d,
	}, nil
}

// ProcessGameSubscriptionProcessing processes a join game turn sheet and creates the necessary
// game entities (game instance, character, character instance, etc.)
func (p *AdventureGameSubscriptionProcessingProcessor) ProcessGameSubscriptionProcessing(ctx context.Context, subscriptionRec *game_record.GameSubscription, turnSheetRec *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("AdventureGameSubscriptionProcessingProcessor/ProcessGameSubscriptionProcessing")

	l.Info("processing join game turn sheet for subscription ID >%s< turn sheet ID >%s<", subscriptionRec.ID, turnSheetRec.ID)

	// Parse the scanned data to get character name and manager subscription ID
	var scanData turn_sheet.AdventureGameJoinGameScanData
	if err := json.Unmarshal(turnSheetRec.ScannedData, &scanData); err != nil {
		l.Warn("failed to unmarshal join game scan data >%v<", err)
		return fmt.Errorf("failed to parse join game scan data: %w", err)
	}

	if scanData.CharacterName == "" {
		l.Warn("character name is empty in scan data")
		return fmt.Errorf("character name is required")
	}

	// Validate scan data for processing (includes game_subscription_id requirement)
	if err := scanData.ValidateForProcessing(); err != nil {
		l.Warn("failed to validate scan data for processing >%v<", err)
		return fmt.Errorf("failed to validate scan data: %w", err)
	}

	// Validate that the manager subscription exists and is a Manager type
	managerSubscriptionRec, err := p.Domain.GetGameSubscriptionRec(scanData.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get manager subscription >%s< >%v<", scanData.GameSubscriptionID, err)
		return fmt.Errorf("failed to get manager subscription: %w", err)
	}

	if managerSubscriptionRec.SubscriptionType != game_record.GameSubscriptionTypeManager {
		l.Warn("subscription >%s< is not a Manager subscription, type is >%s<", scanData.GameSubscriptionID, managerSubscriptionRec.SubscriptionType)
		return fmt.Errorf("game_subscription_id must reference a Manager subscription")
	}

	if managerSubscriptionRec.GameID != subscriptionRec.GameID {
		l.Warn("manager subscription >%s< is for game >%s< but player subscription is for game >%s<", scanData.GameSubscriptionID, managerSubscriptionRec.GameID, subscriptionRec.GameID)
		return fmt.Errorf("manager subscription must be for the same game")
	}

	// Get or create game instance for this manager subscription
	gameInstanceRec, err := p.getOrCreateGameInstance(subscriptionRec.GameID, scanData.GameSubscriptionID)
	if err != nil {
		l.Warn("failed to get or create game instance >%v<", err)
		return fmt.Errorf("failed to get or create game instance: %w", err)
	}

	// Update player subscription to set game_instance_id
	subscriptionRec.GameInstanceID = nullstring.FromString(gameInstanceRec.ID)
	subscriptionRec, err = p.Domain.UpdateGameSubscriptionRec(subscriptionRec)
	if err != nil {
		l.Warn("failed to update player subscription with game_instance_id >%v<", err)
		return fmt.Errorf("failed to update player subscription: %w", err)
	}

	l.Info("updated player subscription >%s< with game_instance_id >%s<", subscriptionRec.ID, gameInstanceRec.ID)

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

// getOrCreateGameInstance gets an existing game instance for a manager subscription with capacity,
// or creates a new one if none exist or all are full
func (p *AdventureGameSubscriptionProcessingProcessor) getOrCreateGameInstance(gameID, managerSubscriptionID string) (*game_record.GameInstance, error) {
	l := p.Logger.WithFunctionContext("getOrCreateGameInstance")

	// Get all game instances for this manager subscription
	gameInstanceRecs, err := p.Domain.GetManyGameInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameInstanceGameID, Val: gameID},
			{Col: game_record.FieldGameInstanceGameSubscriptionID, Val: managerSubscriptionID},
		},
	})
	if err != nil {
		l.Warn("failed to get game instances by game ID >%s< manager subscription ID >%s< >%v<", gameID, managerSubscriptionID, err)
		return nil, err
	}

	// Filter to only active (non-deleted) instances
	var activeInstances []*game_record.GameInstance
	for _, instance := range gameInstanceRecs {
		if !instance.DeletedAt.Valid {
			activeInstances = append(activeInstances, instance)
		}
	}

	// Try to find an instance with capacity (status created or started)
	// For now, we'll use instances that are created or started and haven't reached required_player_count
	// TODO: Add logic to check actual player count vs capacity if needed
	for _, instance := range activeInstances {
		if instance.Status == game_record.GameInstanceStatusCreated || instance.Status == game_record.GameInstanceStatusStarted {
			// Check if instance has capacity by counting player subscriptions
			playerSubscriptions, err := p.Domain.GetManyGameSubscriptionRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: game_record.FieldGameSubscriptionGameID, Val: gameID},
					{Col: game_record.FieldGameSubscriptionGameInstanceID, Val: instance.ID},
					{Col: game_record.FieldGameSubscriptionSubscriptionType, Val: game_record.GameSubscriptionTypePlayer},
				},
			})
			if err != nil {
				l.Warn("failed to get player subscriptions for game instance >%s< >%v<", instance.ID, err)
				continue
			}

			// Count active player subscriptions
			activePlayerCount := 0
			for _, sub := range playerSubscriptions {
				if !sub.DeletedAt.Valid && sub.Status == game_record.GameSubscriptionStatusActive {
					activePlayerCount++
				}
			}

			// If instance has no required_player_count (0) or hasn't reached it, use this instance
			if instance.RequiredPlayerCount == 0 || activePlayerCount < instance.RequiredPlayerCount {
				l.Info("found game instance ID >%s< with capacity (players: %d/%d) for game ID >%s< manager subscription ID >%s<", instance.ID, activePlayerCount, instance.RequiredPlayerCount, gameID, managerSubscriptionID)
				return instance, nil
			}
		}
	}

	// No instance with capacity found, create a new one
	l.Info("creating new game instance for game ID >%s< manager subscription ID >%s<", gameID, managerSubscriptionID)
	gameInstanceRec := &game_record.GameInstance{
		GameID:            gameID,
		GameSubscriptionID: managerSubscriptionID,
		Status:            game_record.GameInstanceStatusCreated,
	}

	gameInstanceRec, err = p.Domain.CreateGameInstanceRec(gameInstanceRec)
	if err != nil {
		l.Warn("failed to create game instance >%v<", err)
		return nil, err
	}

	l.Info("created new game instance ID >%s< for game ID >%s< manager subscription ID >%s<", gameInstanceRec.ID, gameID, managerSubscriptionID)

	return gameInstanceRec, nil
}

// getOrCreateAdventureGameCharacter gets or creates an adventure game character
func (p *AdventureGameSubscriptionProcessingProcessor) getOrCreateAdventureGameCharacter(gameID, accountID, characterName string) (*adventure_game_record.AdventureGameCharacter, error) {
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
func (p *AdventureGameSubscriptionProcessingProcessor) getStartingLocationInstance(gameID, gameInstanceID string) (string, error) {
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
