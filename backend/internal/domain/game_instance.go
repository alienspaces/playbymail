package domain

import (
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

	// Set initial status and default values if not already set
	if rec.Status == "" {
		rec.Status = game_record.GameInstanceStatusCreated
	}

	if rec.CurrentTurn == 0 {
		rec.CurrentTurn = 0
	}

	if !rec.DeliveryPhysicalPost && !rec.DeliveryPhysicalLocal && !rec.DeliveryEmail {
		rec.DeliveryPhysicalPost = true
	}

	if err := m.validateGameInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_instance record >%v<", err)
		return rec, err
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

	_, err := m.GetGameInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
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

	_, err := m.GetGameInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// Game Runtime Management Functions

// StartGameInstance starts a game instance and sets up the first turn
func (m *Domain) StartGameInstance(instanceID string) (*game_record.GameInstance, error) {
	l := m.Logger("StartGameInstance")

	instance, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status != game_record.GameInstanceStatusCreated {
		return nil, fmt.Errorf("game instance must be in 'created' status to start")
	}

	// Check player count meets required count (only if required_player_count > 0)
	if instance.RequiredPlayerCount > 0 {
		playerCount, err := m.GetPlayerCountForGameInstance(instanceID)
		if err != nil {
			l.Warn("failed to get player count for game instance >%s< >%v<", instanceID, err)
			return nil, err
		}

		if playerCount < instance.RequiredPlayerCount {
			return nil, fmt.Errorf("insufficient players: have %d, need %d", playerCount, instance.RequiredPlayerCount)
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
		return nil, err
	}

	l.Info("started game instance >%s<", instanceID)

	return instance, nil
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

	instance, err := m.GetGameInstanceRec(instanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status != game_record.GameInstanceStatusStarted {
		return nil, fmt.Errorf("game instance must be started to complete turns")
	}

	// Check if we've reached max turns
	// max turns handled via configuration; not tracked directly on instance
	if false {
		instance.Status = game_record.GameInstanceStatusCompleted
		now := time.Now()
		instance.CompletedAt = nulltime.FromTime(now)
		l.Info("game instance >%s< completed", instanceID)
	} else {
		// Advance to next turn
		instance.CurrentTurn++
		// next turn due at computed elsewhere
		instance.NextTurnDueAt = record.NewRecordNullTimestamp()
		l.Info("advanced game instance >%s< to turn >%d<", instanceID, instance.CurrentTurn)
	}

	instance, err = m.UpdateGameInstanceRec(instance)
	if err != nil {
		l.Warn("failed updating game instance after turn completion >%v<", err)
		return nil, err
	}

	return instance, nil
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

// GetPlayerCountForGameInstance counts active Player subscriptions for the game instance's game
func (m *Domain) GetPlayerCountForGameInstance(gameInstanceID string) (int, error) {
	l := m.Logger("GetPlayerCountForGameInstance")

	instance, err := m.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return 0, err
	}

	// Get all active Player subscriptions for this game
	subscriptions, err := m.GetManyGameSubscriptionRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionGameID, Val: instance.GameID},
			{Col: game_record.FieldGameSubscriptionSubscriptionType, Val: game_record.GameSubscriptionTypePlayer},
			{Col: game_record.FieldGameSubscriptionStatus, Val: game_record.GameSubscriptionStatusActive},
		},
	})
	if err != nil {
		l.Warn("failed to get player subscriptions for game ID >%s< >%v<", instance.GameID, err)
		return 0, err
	}

	return len(subscriptions), nil
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

// generateUUID generates a UUID string
func generateUUID() (string, error) {
	uuidVal := uuid.New()
	return uuidVal.String(), nil
}
