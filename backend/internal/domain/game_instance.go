package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GetManyGameInstanceRecs -
func (m *Domain) GetManyGameInstanceRecs(opts *coresql.Options) ([]*game_record.GameInstance, error) {
	l := m.Logger("GetManyGameInstanceRecs")

	l.Debug("getting many game_instance records opts >%#v<", opts)

	r := m.GameInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameInstanceRec -
func (m *Domain) GetGameInstanceRec(recID string, lock *coresql.Lock) (*game_record.GameInstance, error) {
	l := m.Logger("GetGameInstanceRec")

	l.Debug("getting game_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(game_record.TableGameInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateGameInstanceRec -
func (m *Domain) CreateGameInstanceRec(rec *game_record.GameInstance) (*game_record.GameInstance, error) {
	l := m.Logger("CreateGameInstanceRec")

	l.Debug("creating game_instance record >%#v<", rec)

	// Set initial status and default values when missing (do not overwrite user-provided data)
	if rec.Status == "" {
		rec.Status = game_record.GameInstanceStatusCreated
	}

	if rec.CurrentTurn == 0 {
		rec.CurrentTurn = 0
	}

	if !rec.DeliveryPhysicalPost && !rec.DeliveryPhysicalLocal && !rec.DeliveryEmail {
		rec.DeliveryPhysicalPost = true
	}

	// Default turn_duration_hours and draft-specific settings from the parent game
	if rec.GameID != "" {
		gameRec, getErr := m.GetGameRec(rec.GameID, nil)
		if getErr == nil && gameRec != nil {
			// Inherit the game's turn duration when the caller did not supply one
			if rec.TurnDurationHours == 0 {
				rec.TurnDurationHours = gameRec.TurnDurationHours
			}
			// Draft games default to closed testing with email delivery
			if gameRec.Status == game_record.GameStatusDraft {
				if !rec.IsClosedTesting {
					rec.IsClosedTesting = true
				}
				if !rec.DeliveryEmail {
					rec.DeliveryEmail = true
				}
			}
		}
	}

	// Generate join_game_key if closed testing is enabled
	if rec.IsClosedTesting && !nullstring.IsValid(rec.ClosedTestingJoinGameKey) {
		closedTestingJoinGameKey, err := generateUUID()
		if err != nil {
			l.Warn("failed to generate join game key >%v<", err)
			return rec, err
		}
		rec.ClosedTestingJoinGameKey = nullstring.FromString(closedTestingJoinGameKey)
	}

	if err := m.validateGameInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_instance record >%v<", err)
		return rec, err
	}

	r := m.GameInstanceRepository()

	createdRec, err := r.CreateOne(rec)
	if err != nil {
		l.Warn("failed to create game_instance record >%v<", err)
		return rec, err
	}

	return createdRec, nil
}

// UpdateGameInstanceRec -
func (m *Domain) UpdateGameInstanceRec(rec *game_record.GameInstance) (*game_record.GameInstance, error) {
	l := m.Logger("UpdateGameInstanceRec")

	currRec, err := m.GetGameInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating game_instance record >%#v<", rec)

	if err := m.validateGameInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate game_instance record >%v<", err)
		return rec, err
	}

	r := m.GameInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteGameInstanceRec -
func (m *Domain) DeleteGameInstanceRec(recID string) error {
	l := m.Logger("DeleteGameInstanceRec")

	l.Debug("deleting game_instance record ID >%s<", recID)

	rec, err := m.GetGameInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if rec.Status != game_record.GameInstanceStatusCancelled {
		l.Warn("game instance cannot be deleted in status >%s<, must be cancelled", rec.Status)
		return coreerror.NewInvalidActionError("delete", "game instance can only be deleted when cancelled")
	}

	r := m.GameInstanceRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveGameInstanceRec -
func (m *Domain) RemoveGameInstanceRec(recID string) error {
	l := m.Logger("RemoveGameInstanceRec")

	l.Debug("removing game_instance record ID >%s<", recID)

	r := m.GameInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// Game Runtime Management Functions

// StartGameInstance starts a game instance: populates all world and player data then transitions status to started.
func (m *Domain) StartGameInstance(instanceID string) (*game_record.GameInstance, *AdventureGameInstanceData, error) {
	l := m.Logger("StartGameInstance")

	instance, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, nil, err
	}

	if instance.Status != game_record.GameInstanceStatusCreated {
		return nil, nil, fmt.Errorf("game instance must be in 'created' status to start")
	}

	// Check player count meets required count (only if required_player_count > 0)
	if instance.RequiredPlayerCount > 0 {
		playerCount, err := m.GetPlayerCountForGameInstance(instanceID)
		if err != nil {
			l.Warn("failed to get player count for game instance >%s< >%v<", instanceID, err)
			return nil, nil, err
		}

		if playerCount < instance.RequiredPlayerCount {
			return nil, nil, fmt.Errorf("insufficient players: have %d, need %d", playerCount, instance.RequiredPlayerCount)
		}
	}

	// Populate all world and player instance data before transitioning status
	gameRec, err := m.GetGameRec(instance.GameID, nil)
	if err != nil {
		return nil, nil, err
	}

	var instanceData *AdventureGameInstanceData
	if gameRec.GameType == game_record.GameTypeAdventure {
		instanceData, err = m.PopulateAdventureGameInstanceData(instanceID)
		if err != nil {
			l.Warn("failed to populate adventure game instance data >%v<", err)
			return nil, nil, err
		}
	}

	now := time.Now()
	instance.Status = game_record.GameInstanceStatusStarted
	instance.StartedAt = nulltime.FromTime(now)
	instance.CurrentTurn = 0

	// NOTE: The turn processing job will determine when the next turn is due.

	instance, err = m.UpdateGameInstanceRec(instance)
	if err != nil {
		l.Warn("failed updating game instance to starting status >%v<", err)
		return nil, nil, err
	}

	l.Info("started game instance >%s<", instanceID)

	return instance, instanceData, nil
}

// BeginTurnProcessing starts processing the current turn
func (m *Domain) BeginTurnProcessing(instanceID string) (*game_record.GameInstance, error) {
	l := m.Logger("BeginTurnProcessing")

	instance, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status != game_record.GameInstanceStatusStarted {
		return nil, fmt.Errorf("game instance must be started to process turns")
	}

	instance.Status = game_record.GameInstanceStatusStarted
	now := time.Now()
	instance.LastTurnProcessedAt = nulltime.FromTime(now)

	instance, err = m.UpdateGameInstanceRec(instance)
	if err != nil {
		l.Warn("failed updating game instance for turn processing >%v<", err)
		return nil, err
	}

	l.Info("began turn processing for game instance >%s< turn >%d<", instanceID, instance.CurrentTurn)
	return instance, nil
}

// CompleteTurn advances the game to the next turn
func (m *Domain) CompleteTurn(instanceID string) (*game_record.GameInstance, error) {
	l := m.Logger("CompleteTurn")

	gameInstanceRec, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if gameInstanceRec.Status != game_record.GameInstanceStatusStarted {
		return nil, fmt.Errorf("game instance must be started to complete turns")
	}

	// Check if we've reached max turns
	// max turns handled via configuration; not tracked directly on instance
	if false {
		gameInstanceRec.Status = game_record.GameInstanceStatusCompleted
		now := time.Now()
		gameInstanceRec.CompletedAt = nulltime.FromTime(now)
		l.Info("game instance >%s< completed", instanceID)
	} else {
		// Advance to next turn
		gameInstanceRec.CurrentTurn++

		if gameInstanceRec.ProcessWhenAllSubmitted {
			// Player-driven: leave NextTurnDueAt null so the periodic worker
			// won't queue processing until the player submits all sheets.
			gameInstanceRec.NextTurnDueAt = sql.NullTime{}
			l.Info("advanced game instance >%s< to turn >%d< (process-when-all-submitted; awaiting player)", instanceID, gameInstanceRec.CurrentTurn)
		} else {
			nextTurn := time.Now().UTC().Add(time.Duration(gameInstanceRec.TurnDurationHours) * time.Hour)
			gameInstanceRec.NextTurnDueAt = nulltime.FromTime(nextTurn)
			l.Info("advanced game instance >%s< to turn >%d< (next turn due at >%s<)", instanceID, gameInstanceRec.CurrentTurn, nextTurn.Format(time.RFC3339))
		}
	}

	gameInstanceRec, err = m.UpdateGameInstanceRec(gameInstanceRec)
	if err != nil {
		l.Warn("failed updating game instance after turn completion >%v<", err)
		return nil, err
	}

	return gameInstanceRec, nil
}

// PauseGameInstance pauses a running game instance
func (m *Domain) PauseGameInstance(instanceID string) (*game_record.GameInstance, error) {
	l := m.Logger("PauseGameInstance")

	instanceRec, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instanceRec.Status != game_record.GameInstanceStatusStarted {
		return nil, fmt.Errorf("game instance must be started to pause")
	}

	instanceRec.Status = game_record.GameInstanceStatusPaused

	instanceRec, err = m.UpdateGameInstanceRec(instanceRec)
	if err != nil {
		l.Warn("failed updating game instance to paused status >%v<", err)
		return nil, err
	}

	l.Info("paused game instance >%s<", instanceID)
	return instanceRec, nil
}

// ResumeGameInstance resumes a paused game instance
func (m *Domain) ResumeGameInstance(instanceID string) (*game_record.GameInstance, error) {
	l := m.Logger("ResumeGameInstance")

	instanceRec, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instanceRec.Status != game_record.GameInstanceStatusPaused {
		return nil, fmt.Errorf("game instance must be paused to resume")
	}

	instanceRec.Status = game_record.GameInstanceStatusStarted

	instanceRec, err = m.UpdateGameInstanceRec(instanceRec)
	if err != nil {
		l.Warn("failed updating game instance to started status >%v<", err)
		return nil, err
	}

	l.Info("resumed game instance >%s<", instanceID)
	return instanceRec, nil
}

// CancelGameInstance cancels a game instance
func (m *Domain) CancelGameInstance(instanceID string) (*game_record.GameInstance, error) {
	l := m.Logger("CancelGameInstance")

	instanceRec, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instanceRec.Status == game_record.GameInstanceStatusCompleted || instanceRec.Status == game_record.GameInstanceStatusCancelled {
		return nil, fmt.Errorf("game instance is already completed or cancelled")
	}

	instanceRec.Status = game_record.GameInstanceStatusCancelled
	instanceRec.CompletedAt = record.NewRecordNullTimestamp()

	instanceRec, err = m.UpdateGameInstanceRec(instanceRec)
	if err != nil {
		l.Warn("failed updating game instance to cancelled status >%v<", err)
		return nil, err
	}

	l.Info("cancelled game instance >%s<", instanceID)
	return instanceRec, nil
}

// GetPlayerCountForGameInstance counts player-type subscriptions linked to the game instance.
// Manager and designer subscription instances are excluded from the count.
func (m *Domain) GetPlayerCountForGameInstance(gameInstanceID string) (int, error) {
	l := m.Logger("GetPlayerCountForGameInstance")

	if err := domain.ValidateUUIDField("game_instance_id", gameInstanceID); err != nil {
		return 0, err
	}

	gameSubscriptionInstanceRecs, err := m.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: gameInstanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get subscription instances for game instance ID >%s< >%v<", gameInstanceID, err)
		return 0, err
	}

	count := 0
	for _, gameSubscriptionInstanceRec := range gameSubscriptionInstanceRecs {
		gameSubscriptionRec, err := m.GetGameSubscriptionRec(gameSubscriptionInstanceRec.GameSubscriptionID, nil)
		if err != nil {
			l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionInstanceRec.GameSubscriptionID, err)
			continue
		}
		if gameSubscriptionRec.SubscriptionType == game_record.GameSubscriptionTypePlayer {
			count++
		}
	}

	l.Info("player count for game instance >%s< is >%d<", gameInstanceID, count)

	return count, nil
}

// GameInstanceHasAvailableCapacity checks if a game instance has available player slots
func (m *Domain) GameInstanceHasAvailableCapacity(gameInstanceID string) (bool, error) {
	instance, err := m.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return false, err
	}

	// If required_player_count is 0, there's no capacity limit
	if instance.RequiredPlayerCount == 0 {
		return true, nil
	}

	playerCount, err := m.GetPlayerCountForGameInstance(gameInstanceID)
	if err != nil {
		return false, err
	}

	return playerCount < instance.RequiredPlayerCount, nil
}

// FindAvailableGameInstance finds an available game instance for a game subscription
// Returns the first instance with available capacity, or nil if none available
func (m *Domain) FindAvailableGameInstance(gameSubscriptionID string) (*game_record.GameInstance, error) {
	l := m.Logger("FindAvailableGameInstance")

	// Get the subscription to find the game ID
	subscription, err := m.GetGameSubscriptionRec(gameSubscriptionID, nil)
	if err != nil {
		return nil, err
	}

	// Get all active instances for this game
	instances, err := m.GetManyGameInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameInstanceGameID, Val: subscription.GameID},
			{Col: game_record.FieldGameInstanceStatus, Val: game_record.GameInstanceStatusCreated},
		},
	})
	if err != nil {
		l.Warn("failed to get game instances for game ID >%s< >%v<", subscription.GameID, err)
		return nil, err
	}

	// Find first instance with available capacity
	for _, instance := range instances {
		hasCapacity, err := m.GameInstanceHasAvailableCapacity(instance.ID)
		if err != nil {
			l.Warn("failed to check capacity for instance >%s< >%v<", instance.ID, err)
			continue
		}
		if hasCapacity {
			return instance, nil
		}
	}

	return nil, nil
}

// AssignPlayerToGameInstance assigns a player subscription to a game instance
// Creates a game_subscription_instance record linking them
func (m *Domain) AssignPlayerToGameInstance(gameSubscriptionID, gameInstanceID string) (*game_record.GameSubscriptionInstance, error) {
	l := m.Logger("AssignPlayerToGameInstance")

	// Validate inputs
	if err := domain.ValidateUUIDField("game_subscription_id", gameSubscriptionID); err != nil {
		return nil, err
	}
	if err := domain.ValidateUUIDField("game_instance_id", gameInstanceID); err != nil {
		return nil, err
	}

	// Get subscription to get account_id
	gameSubscriptionRec, err := m.GetGameSubscriptionRec(gameSubscriptionID, nil)
	if err != nil {
		return nil, err
	}

	// Check capacity before assigning
	hasCapacity, err := m.GameInstanceHasAvailableCapacity(gameInstanceID)
	if err != nil {
		return nil, err
	}
	if !hasCapacity {
		return nil, fmt.Errorf("game instance has no available capacity")
	}

	// Create the subscription-instance link
	gameSubscriptionInstanceRec := &game_record.GameSubscriptionInstance{
		AccountID:          gameSubscriptionRec.AccountID,
		AccountUserID:      gameSubscriptionRec.AccountUserID,
		GameSubscriptionID: gameSubscriptionID,
		GameInstanceID:     gameInstanceID,
	}

	createdRec, err := m.CreateGameSubscriptionInstanceRec(gameSubscriptionInstanceRec)
	if err != nil {
		l.Warn("failed to create game subscription instance >%v<", err)
		return nil, err
	}

	return createdRec, nil
}

// GenerateClosedTestingJoinGameKey generates a UUID join game key for closed testing instances.
// Will return the existing key if it exists and has not expired, otherwise it will generate a new one.
func (m *Domain) GenerateClosedTestingJoinGameKey(gameInstanceID string) (string, error) {
	l := m.Logger("GenerateClosedTestingJoinGameKey")

	gameInstanceRec, err := m.GetGameInstanceRec(gameInstanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return "", err
	}

	if !gameInstanceRec.IsClosedTesting {
		return "", fmt.Errorf("game instance is not in closed testing mode")
	}

	// Check if there is already a key and it has not expired
	if nullstring.IsValid(gameInstanceRec.ClosedTestingJoinGameKey) && nulltime.IsValid(gameInstanceRec.ClosedTestingJoinGameKeyExpiresAt) {
		if time.Now().Before(nulltime.ToTime(gameInstanceRec.ClosedTestingJoinGameKeyExpiresAt)) {
			return nullstring.ToString(gameInstanceRec.ClosedTestingJoinGameKey), nil
		}
	}

	// Generate UUID for closed testing join game key
	closedTestingJoinGameKey, err := generateUUID()
	if err != nil {
		l.Warn("failed to generate closed testing join game key >%v<", err)
		return "", err
	}

	// Set closed testing join game key and expiration
	gameInstanceRec.ClosedTestingJoinGameKey = nullstring.FromString(closedTestingJoinGameKey)
	gameInstanceRec.ClosedTestingJoinGameKeyExpiresAt = nulltime.FromTime(time.Now().Add(3 * 24 * time.Hour))

	// Update game instance with closed testing join game key
	_, err = m.UpdateGameInstanceRec(gameInstanceRec)
	if err != nil {
		l.Warn("failed updating game instance with closed testing join game key >%v<", err)
		return "", err
	}

	l.Info("generated closed testing join game key for game instance >%s<", gameInstanceID)

	return closedTestingJoinGameKey, nil
}

// GetGameInstanceByClosedTestingJoinGameKey looks up a game instance by join game key
func (m *Domain) GetGameInstanceByClosedTestingJoinGameKey(closedTestingJoinGameKey string) (*game_record.GameInstance, error) {
	l := m.Logger("GetGameInstanceByClosedTestingJoinGameKey")

	if closedTestingJoinGameKey == "" {
		return nil, coreerror.NewInvalidDataError("join_game_key is required")
	}

	// Get game instance by join_game_key
	instances, err := m.GetManyGameInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameInstanceClosedTestingJoinGameKey, Val: closedTestingJoinGameKey},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get game instance by join game key >%s< >%v<", closedTestingJoinGameKey, err)
		return nil, err
	}

	if len(instances) == 0 {
		return nil, coreerror.NewNotFoundError(game_record.TableGameInstance, closedTestingJoinGameKey)
	}

	instance := instances[0]

	// Validate it's in closed testing mode
	if !instance.IsClosedTesting {
		return nil, fmt.Errorf("game instance is not in closed testing mode")
	}

	// Check expiration if set
	if instance.ClosedTestingJoinGameKeyExpiresAt.Valid {
		if time.Now().After(instance.ClosedTestingJoinGameKeyExpiresAt.Time) {
			return nil, fmt.Errorf("join game key has expired")
		}
	}

	return instance, nil
}

// ResetGameInstance resets a game instance to its initial state. All instance-level
// data (turn sheets, character/creature/item/location instances, parameters) is
// soft-deleted. The game_instance record is reset to status=created with turn 0.
// Subscription links (game_subscription_instance) are preserved so players remain
// joined to the instance.
func (m *Domain) ResetGameInstance(instanceID string) (*game_record.GameInstance, error) {
	l := m.Logger("ResetGameInstance")

	instance, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status == game_record.GameInstanceStatusCompleted {
		return nil, fmt.Errorf("cannot reset a completed game instance")
	}

	// 1. Delete adventure_game_turn_sheet records (linked via character instance)
	charInstances, err := m.GetManyAdventureGameCharacterInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCharacterInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get character instances for reset >%v<", err)
		return nil, databaseError(err)
	}

	for _, charInst := range charInstances {
		turnSheets, err := m.GetManyAdventureGameTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID, Val: charInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get adventure turn sheets for character instance >%s< >%v<", charInst.ID, err)
			return nil, databaseError(err)
		}
		for _, ts := range turnSheets {
			if err := m.AdventureGameTurnSheetRepository().DeleteOne(ts.ID); err != nil {
				l.Warn("failed to delete adventure turn sheet >%s< >%v<", ts.ID, err)
				return nil, databaseError(err)
			}
		}
	}

	// 2. Delete game_turn_sheet records
	gameTurnSheetRepo := m.GameTurnSheetRepository()
	gameTurnSheets, err := gameTurnSheetRepo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameTurnSheetGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get game turn sheets for reset >%v<", err)
		return nil, databaseError(err)
	}
	for _, ts := range gameTurnSheets {
		if err := gameTurnSheetRepo.DeleteOne(ts.ID); err != nil {
			l.Warn("failed to delete game turn sheet >%s< >%v<", ts.ID, err)
			return nil, databaseError(err)
		}
	}

	// 3. Delete adventure_game_item_instance records
	itemInstances, err := m.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get item instances for reset >%v<", err)
		return nil, databaseError(err)
	}
	for _, item := range itemInstances {
		if err := m.AdventureGameItemInstanceRepository().DeleteOne(item.ID); err != nil {
			l.Warn("failed to delete item instance >%s< >%v<", item.ID, err)
			return nil, databaseError(err)
		}
	}

	// 4. Delete adventure_game_creature_instance records
	creatureInstances, err := m.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get creature instances for reset >%v<", err)
		return nil, databaseError(err)
	}
	for _, creature := range creatureInstances {
		if err := m.AdventureGameCreatureInstanceRepository().DeleteOne(creature.ID); err != nil {
			l.Warn("failed to delete creature instance >%s< >%v<", creature.ID, err)
			return nil, databaseError(err)
		}
	}

	// 5. Delete adventure_game_character_instance records
	for _, charInst := range charInstances {
		if err := m.AdventureGameCharacterInstanceRepository().DeleteOne(charInst.ID); err != nil {
			l.Warn("failed to delete character instance >%s< >%v<", charInst.ID, err)
			return nil, databaseError(err)
		}
	}

	// 6. Delete adventure_game_location_instance records
	locationInstances, err := m.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get location instances for reset >%v<", err)
		return nil, databaseError(err)
	}
	for _, loc := range locationInstances {
		if err := m.AdventureGameLocationInstanceRepository().DeleteOne(loc.ID); err != nil {
			l.Warn("failed to delete location instance >%s< >%v<", loc.ID, err)
			return nil, databaseError(err)
		}
	}

	// 7. Delete game_instance_parameter records
	params, err := m.GetManyGameInstanceParameterRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameInstanceParameterGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get instance parameters for reset >%v<", err)
		return nil, databaseError(err)
	}
	for _, p := range params {
		if err := m.GameInstanceParameterRepository().DeleteOne(p.ID); err != nil {
			l.Warn("failed to delete instance parameter >%s< >%v<", p.ID, err)
			return nil, databaseError(err)
		}
	}

	// 8. Reset the game instance record — uses repository directly because
	// the standard update validation prevents current_turn from decreasing.
	instance.Status = game_record.GameInstanceStatusCreated
	instance.CurrentTurn = 0
	instance.StartedAt = nulltime.FromTime(time.Time{})
	instance.CompletedAt = nulltime.FromTime(time.Time{})
	instance.LastTurnProcessedAt = nulltime.FromTime(time.Time{})
	instance.NextTurnDueAt = nulltime.FromTime(time.Time{})

	if instance.IsClosedTesting {
		key, keyErr := generateUUID()
		if keyErr != nil {
			l.Warn("failed to generate new join game key >%v<", keyErr)
			return nil, keyErr
		}
		instance.ClosedTestingJoinGameKey = nullstring.FromString(key)
		instance.ClosedTestingJoinGameKeyExpiresAt = nulltime.FromTime(time.Now().Add(3 * 24 * time.Hour))
	}

	r := m.GameInstanceRepository()
	instance, err = r.UpdateOne(instance)
	if err != nil {
		l.Warn("failed to reset game instance record >%v<", err)
		return nil, databaseError(err)
	}

	l.Info("reset game instance >%s< to initial state", instanceID)

	return instance, nil
}

// AdventureGameInstanceData holds all instance records created when a game instance starts.
type AdventureGameInstanceData struct {
	LocationInstances       []*adventure_game_record.AdventureGameLocationInstance
	CreatureInstances       []*adventure_game_record.AdventureGameCreatureInstance
	LocationObjectInstances []*adventure_game_record.AdventureGameLocationObjectInstance
	ItemInstances           []*adventure_game_record.AdventureGameItemInstance
	CharacterInstances      []*adventure_game_record.AdventureGameCharacterInstance
}

// PopulateAdventureGameInstanceData creates all instance records (locations, creatures, objects,
// items, characters) for an adventure game instance from its definitions and player subscriptions.
func (m *Domain) PopulateAdventureGameInstanceData(instanceID string) (*AdventureGameInstanceData, error) {
	l := m.Logger("PopulateAdventureGameInstanceData")

	l.Info("populating adventure game instance data for instance >%s<", instanceID)

	instanceRec, err := m.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		return nil, err
	}

	gameID := instanceRec.GameID
	out := &AdventureGameInstanceData{}

	// 1. Create location instances and build locationID -> locationInstanceID map
	locationRecs, err := m.GetManyAdventureGameLocationRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID},
		},
	})
	if err != nil {
		l.Warn("failed to get adventure game locations >%v<", err)
		return nil, err
	}

	locationIDToInstanceID := make(map[string]string, len(locationRecs))
	var startingLocationInstanceID string

	for _, loc := range locationRecs {
		locInst, err := m.CreateAdventureGameLocationInstanceRec(&adventure_game_record.AdventureGameLocationInstance{
			GameID:                  gameID,
			GameInstanceID:          instanceID,
			AdventureGameLocationID: loc.ID,
		})
		if err != nil {
			l.Warn("failed to create location instance for location >%s< >%v<", loc.ID, err)
			return nil, err
		}
		locationIDToInstanceID[loc.ID] = locInst.ID
		out.LocationInstances = append(out.LocationInstances, locInst)
		if loc.IsStartingLocation && startingLocationInstanceID == "" {
			startingLocationInstanceID = locInst.ID
		}
	}

	if startingLocationInstanceID == "" && len(locationRecs) > 0 {
		return nil, fmt.Errorf("no starting location found for game >%s<: at least one location must have is_starting_location = true", gameID)
	}

	// 2. Create creature instances from creature placements
	creaturePlacements, err := m.GetManyAdventureGameCreaturePlacementRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreaturePlacementGameID, Val: gameID},
		},
	})
	if err != nil {
		l.Warn("failed to get creature placements >%v<", err)
		return nil, err
	}

	for _, placement := range creaturePlacements {
		locationInstanceID, ok := locationIDToInstanceID[placement.AdventureGameLocationID]
		if !ok {
			l.Warn("no location instance found for creature placement location >%s< - skipping", placement.AdventureGameLocationID)
			continue
		}

		creatureDef, err := m.GetAdventureGameCreatureRec(placement.AdventureGameCreatureID, nil)
		if err != nil {
			l.Warn("failed to get creature definition >%s< >%v<", placement.AdventureGameCreatureID, err)
			return nil, err
		}

		count := placement.InitialCount
		if count <= 0 {
			count = 1
		}
		for i := 0; i < count; i++ {
			creatureInst, err := m.CreateAdventureGameCreatureInstanceRec(&adventure_game_record.AdventureGameCreatureInstance{
				GameID:                          gameID,
				GameInstanceID:                  instanceID,
				AdventureGameCreatureID:         placement.AdventureGameCreatureID,
				AdventureGameLocationInstanceID: locationInstanceID,
				Health:                          creatureDef.MaxHealth,
			})
			if err != nil {
				l.Warn("failed to create creature instance >%v<", err)
				return nil, err
			}
			out.CreatureInstances = append(out.CreatureInstances, creatureInst)
		}
	}

	// 3. Create location object instances
	objectRecs, err := m.GetManyAdventureGameLocationObjectRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectGameID, Val: gameID},
		},
	})
	if err != nil {
		l.Warn("failed to get location objects >%v<", err)
		return nil, err
	}

	for _, obj := range objectRecs {
		locationInstanceID, ok := locationIDToInstanceID[obj.AdventureGameLocationID]
		if !ok {
			l.Warn("no location instance found for object >%s< location >%s< - skipping", obj.ID, obj.AdventureGameLocationID)
			continue
		}
		objInst, err := m.CreateAdventureGameLocationObjectInstanceRec(&adventure_game_record.AdventureGameLocationObjectInstance{
			GameID:                          gameID,
			GameInstanceID:                  instanceID,
			AdventureGameLocationObjectID:   obj.ID,
			AdventureGameLocationInstanceID: locationInstanceID,
			CurrentState:                    obj.InitialState,
			IsVisible:                       !obj.IsHidden,
		})
		if err != nil {
			l.Warn("failed to create location object instance >%v<", err)
			return nil, err
		}
		out.LocationObjectInstances = append(out.LocationObjectInstances, objInst)
	}

	// 4. Create item instances from item placements
	itemPlacements, err := m.GetManyAdventureGameItemPlacementRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemPlacementGameID, Val: gameID},
		},
	})
	if err != nil {
		l.Warn("failed to get item placements >%v<", err)
		return nil, err
	}

	for _, placement := range itemPlacements {
		locationInstanceID, ok := locationIDToInstanceID[placement.AdventureGameLocationID]
		if !ok {
			l.Warn("no location instance found for item placement location >%s< - skipping", placement.AdventureGameLocationID)
			continue
		}

		count := placement.InitialCount
		if count <= 0 {
			count = 1
		}
		for i := 0; i < count; i++ {
			itemInst, err := m.CreateAdventureGameItemInstanceRec(&adventure_game_record.AdventureGameItemInstance{
				GameID:                          gameID,
				GameInstanceID:                  instanceID,
				AdventureGameItemID:             placement.AdventureGameItemID,
				AdventureGameLocationInstanceID: nullstring.FromString(locationInstanceID),
			})
			if err != nil {
				l.Warn("failed to create item instance >%v<", err)
				return nil, err
			}
			out.ItemInstances = append(out.ItemInstances, itemInst)
		}
	}

	// 5. Create character instances for all subscribed players (deduplicated per character)
	subscriptionInstances, err := m.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get subscription instances >%v<", err)
		return nil, err
	}

	createdCharacterIDs := make(map[string]bool)
	for _, subInst := range subscriptionInstances {
		// Get game subscription to find account_user_id
		subRecs, err := m.GetManyGameSubscriptionRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: game_record.FieldGameSubscriptionID, Val: subInst.GameSubscriptionID},
			},
			Limit: 1,
		})
		if err != nil || len(subRecs) == 0 {
			l.Warn("failed to get game subscription >%s< >%v<", subInst.GameSubscriptionID, err)
			continue
		}
		sub := subRecs[0]

		// Skip non-player subscriptions
		if sub.SubscriptionType != game_record.GameSubscriptionTypePlayer {
			continue
		}

		// Find the adventure game character for this player
		charRecs, err := m.GetManyAdventureGameCharacterRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameCharacterGameID, Val: gameID},
				{Col: adventure_game_record.FieldAdventureGameCharacterAccountUserID, Val: sub.AccountUserID},
			},
			Limit: 1,
		})
		if err != nil || len(charRecs) == 0 {
			l.Warn("no adventure game character found for account_user_id >%s< game >%s< - skipping", sub.AccountUserID, gameID)
			continue
		}

		characterID := charRecs[0].ID
		if createdCharacterIDs[characterID] {
			l.Info("character instance already created for character >%s< - skipping duplicate subscription", characterID)
			continue
		}

		if startingLocationInstanceID == "" {
			l.Warn("no starting location instance found for game >%s< - cannot create character instance", gameID)
			return nil, fmt.Errorf("no starting location instance found for game >%s<", gameID)
		}

		charInst, err := m.CreateAdventureGameCharacterInstanceRec(&adventure_game_record.AdventureGameCharacterInstance{
			GameID:                          gameID,
			GameInstanceID:                  instanceID,
			AdventureGameCharacterID:        characterID,
			AdventureGameLocationInstanceID: startingLocationInstanceID,
			Health:                          100,
			InventoryCapacity:               10,
			LastTurnEvents:                  []byte("[]"),
		})
		if err != nil {
			l.Warn("failed to create character instance for character >%s< >%v<", characterID, err)
			return nil, err
		}
		createdCharacterIDs[characterID] = true
		out.CharacterInstances = append(out.CharacterInstances, charInst)

		if err := m.AssignStartingItemsToCharacterInstance(charInst); err != nil {
			l.Warn("failed to assign starting items to character instance >%s< >%v<", charInst.ID, err)
			return nil, err
		}

		l.Info("created character instance >%s< for character >%s< at location >%s<", charInst.ID, characterID, startingLocationInstanceID)
	}

	l.Info("populated adventure game instance data: locations=%d creatures=%d objects=%d items=%d characters=%d",
		len(out.LocationInstances), len(out.CreatureInstances), len(out.LocationObjectInstances),
		len(out.ItemInstances), len(out.CharacterInstances))

	return out, nil
}

// DeleteGameInstance permanently removes a game instance and all associated data.
func (m *Domain) DeleteGameInstance(instanceID string) error {
	l := m.Logger("DeleteGameInstance")

	l.Info("deleting game instance >%s< and all associated data", instanceID)

	// Reuse reset cleanup logic for instance data, then delete the instance record itself.
	// We get the instance to confirm it exists before deletion.
	_, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	// Delete adventure_game_turn_sheet records (linked via character instance)
	charInstances, err := m.GetManyAdventureGameCharacterInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCharacterInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get character instances >%v<", err)
		return databaseError(err)
	}

	for _, charInst := range charInstances {
		turnSheets, err := m.GetManyAdventureGameTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID, Val: charInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for character instance >%s< >%v<", charInst.ID, err)
			return databaseError(err)
		}
		for _, ts := range turnSheets {
			if err := m.AdventureGameTurnSheetRepository().DeleteOne(ts.ID); err != nil {
				l.Warn("failed to delete adventure turn sheet >%s< >%v<", ts.ID, err)
				return databaseError(err)
			}
		}
	}

	// Delete game_turn_sheet records
	gameTurnSheetRepo := m.GameTurnSheetRepository()
	gameTurnSheets, err := gameTurnSheetRepo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameTurnSheetGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get game turn sheets >%v<", err)
		return databaseError(err)
	}
	for _, ts := range gameTurnSheets {
		if err := gameTurnSheetRepo.DeleteOne(ts.ID); err != nil {
			l.Warn("failed to delete game turn sheet >%s< >%v<", ts.ID, err)
			return databaseError(err)
		}
	}

	// Delete item instances
	itemInstances, err := m.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get item instances >%v<", err)
		return databaseError(err)
	}
	for _, item := range itemInstances {
		if err := m.AdventureGameItemInstanceRepository().DeleteOne(item.ID); err != nil {
			l.Warn("failed to delete item instance >%s< >%v<", item.ID, err)
			return databaseError(err)
		}
	}

	// Delete creature instances
	creatureInstances, err := m.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get creature instances >%v<", err)
		return databaseError(err)
	}
	for _, creature := range creatureInstances {
		if err := m.AdventureGameCreatureInstanceRepository().DeleteOne(creature.ID); err != nil {
			l.Warn("failed to delete creature instance >%s< >%v<", creature.ID, err)
			return databaseError(err)
		}
	}

	// Delete character instances
	for _, charInst := range charInstances {
		if err := m.AdventureGameCharacterInstanceRepository().DeleteOne(charInst.ID); err != nil {
			l.Warn("failed to delete character instance >%s< >%v<", charInst.ID, err)
			return databaseError(err)
		}
	}

	// Delete location object instances
	locationObjectInstances, err := m.GetManyAdventureGameLocationObjectInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get location object instances >%v<", err)
		return databaseError(err)
	}
	for _, objInst := range locationObjectInstances {
		if err := m.AdventureGameLocationObjectInstanceRepository().DeleteOne(objInst.ID); err != nil {
			l.Warn("failed to delete location object instance >%s< >%v<", objInst.ID, err)
			return databaseError(err)
		}
	}

	// Delete location instances
	locationInstances, err := m.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get location instances >%v<", err)
		return databaseError(err)
	}
	for _, loc := range locationInstances {
		if err := m.AdventureGameLocationInstanceRepository().DeleteOne(loc.ID); err != nil {
			l.Warn("failed to delete location instance >%s< >%v<", loc.ID, err)
			return databaseError(err)
		}
	}

	// Delete game instance parameters
	params, err := m.GetManyGameInstanceParameterRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameInstanceParameterGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get instance parameters >%v<", err)
		return databaseError(err)
	}
	for _, p := range params {
		if err := m.GameInstanceParameterRepository().DeleteOne(p.ID); err != nil {
			l.Warn("failed to delete instance parameter >%s< >%v<", p.ID, err)
			return databaseError(err)
		}
	}

	// Delete game_subscription_instance links
	subscriptionInstances, err := m.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get subscription instances >%v<", err)
		return databaseError(err)
	}
	for _, subInst := range subscriptionInstances {
		if err := m.GameSubscriptionInstanceRepository().DeleteOne(subInst.ID); err != nil {
			l.Warn("failed to delete subscription instance >%s< >%v<", subInst.ID, err)
			return databaseError(err)
		}
	}

	// Finally delete the game instance record itself
	if err := m.RemoveGameInstanceRec(instanceID); err != nil {
		l.Warn("failed to delete game instance record >%s< >%v<", instanceID, err)
		return err
	}

	l.Info("deleted game instance >%s< and all associated data", instanceID)

	return nil
}

// RemoveGameInstance physically removes (hard delete) a game instance and all its associated data.
// Unlike DeleteGameInstance which sets deleted_at, this permanently removes all rows from the database.
// Used by the test harness for teardown.
func (m *Domain) RemoveGameInstance(instanceID string) error {
	l := m.Logger("RemoveGameInstance")

	l.Info("removing game instance >%s< and all associated data", instanceID)

	_, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	// Remove adventure_game_turn_sheet records (linked via character instance)
	charInstances, err := m.GetManyAdventureGameCharacterInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCharacterInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get character instances >%v<", err)
		return databaseError(err)
	}

	for _, charInst := range charInstances {
		turnSheets, err := m.GetManyAdventureGameTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID, Val: charInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for character instance >%s< >%v<", charInst.ID, err)
			return databaseError(err)
		}
		for _, ts := range turnSheets {
			if err := m.RemoveAdventureGameTurnSheetRec(ts.ID); err != nil {
				l.Warn("failed to remove adventure turn sheet >%s< >%v<", ts.ID, err)
				return err
			}
		}
	}

	// Remove game_turn_sheet records
	gameTurnSheets, err := m.GetManyGameTurnSheetRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameTurnSheetGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get game turn sheets >%v<", err)
		return databaseError(err)
	}
	for _, ts := range gameTurnSheets {
		if err := m.RemoveGameTurnSheetRec(ts.ID); err != nil {
			l.Warn("failed to remove game turn sheet >%s< >%v<", ts.ID, err)
			return err
		}
	}

	// Remove item instances
	itemInstances, err := m.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get item instances >%v<", err)
		return databaseError(err)
	}
	for _, item := range itemInstances {
		if err := m.RemoveAdventureGameItemInstanceRec(item.ID); err != nil {
			l.Warn("failed to remove item instance >%s< >%v<", item.ID, err)
			return err
		}
	}

	// Remove creature instances
	creatureInstances, err := m.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get creature instances >%v<", err)
		return databaseError(err)
	}
	for _, creature := range creatureInstances {
		if err := m.RemoveAdventureGameCreatureInstanceRec(creature.ID); err != nil {
			l.Warn("failed to remove creature instance >%s< >%v<", creature.ID, err)
			return err
		}
	}

	// Remove character instances
	for _, charInst := range charInstances {
		if err := m.RemoveAdventureGameCharacterInstanceRec(charInst.ID); err != nil {
			l.Warn("failed to remove character instance >%s< >%v<", charInst.ID, err)
			return err
		}
	}

	// Remove location object instances
	locationObjectInstances, err := m.GetManyAdventureGameLocationObjectInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get location object instances >%v<", err)
		return databaseError(err)
	}
	for _, objInst := range locationObjectInstances {
		if err := m.RemoveAdventureGameLocationObjectInstanceRec(objInst.ID); err != nil {
			l.Warn("failed to remove location object instance >%s< >%v<", objInst.ID, err)
			return err
		}
	}

	// Remove location instances
	locationInstances, err := m.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get location instances >%v<", err)
		return databaseError(err)
	}
	for _, loc := range locationInstances {
		if err := m.RemoveAdventureGameLocationInstanceRec(loc.ID); err != nil {
			l.Warn("failed to remove location instance >%s< >%v<", loc.ID, err)
			return err
		}
	}

	// Remove game instance parameters
	params, err := m.GetManyGameInstanceParameterRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameInstanceParameterGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get instance parameters >%v<", err)
		return databaseError(err)
	}
	for _, p := range params {
		if err := m.RemoveGameInstanceParameterRec(p.ID); err != nil {
			l.Warn("failed to remove instance parameter >%s< >%v<", p.ID, err)
			return err
		}
	}

	// Remove game_subscription_instance links
	subscriptionInstances, err := m.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get subscription instances >%v<", err)
		return databaseError(err)
	}
	for _, subInst := range subscriptionInstances {
		if err := m.RemoveGameSubscriptionInstanceRec(subInst.ID); err != nil {
			l.Warn("failed to remove subscription instance >%s< >%v<", subInst.ID, err)
			return err
		}
	}

	// Finally remove the game instance record itself
	if err := m.RemoveGameInstanceRec(instanceID); err != nil {
		l.Warn("failed to remove game instance record >%s< >%v<", instanceID, err)
		return err
	}

	l.Info("removed game instance >%s< and all associated data", instanceID)

	return nil
}

// generateUUID generates a UUID string
func generateUUID() (string, error) {
	uuidVal := uuid.New()
	return uuidVal.String(), nil
}
