package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
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

	if err := m.validateGameInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_instance record >%v<", err)
		return rec, err
	}

	// Validate adventure game has starting location
	if err := m.validateAdventureGameInstanceCreation(rec.GameID); err != nil {
		l.Warn("failed adventure game instance validation >%v<", err)
		return rec, err
	}

	r := m.GameInstanceRepository()

	// Set initial status and default values if not already set
	if rec.Status == "" {
		rec.Status = game_record.GameInstanceStatusCreated
	}
	if rec.CurrentTurn == 0 {
		rec.CurrentTurn = 0
	}
	// no per-instance turn deadline config here; configs are handled separately

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// validateAdventureGameInstanceCreation ensures adventure games have at least one starting location
func (m *Domain) validateAdventureGameInstanceCreation(gameID string) error {
	l := m.Logger("validateAdventureGameInstanceCreation")

	// Get the game to check if it's an adventure game
	gameRec, err := m.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed to get game record for game ID >%s< >%v<", gameID, err)
		return err
	}

	// Only validate for adventure games
	if gameRec.GameType != game_record.GameTypeAdventure {
		l.Info("game ID >%s< is not an adventure game, skipping starting location validation", gameID)
		return nil
	}

	// Check for starting locations
	startingLocationRecs, err := m.GetManyAdventureGameLocationRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID},
			{Col: adventure_game_record.FieldAdventureGameLocationIsStartingLocation, Val: true},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get starting locations for game ID >%s< >%v<", gameID, err)
		return err
	}

	if len(startingLocationRecs) == 0 {
		l.Warn("no starting locations found for game ID >%s<", gameID)
		return InvalidField(adventure_game_record.FieldAdventureGameLocationGameID, gameID, "adventure game must have at least one starting location before creating an instance")
	}

	return nil
}

// UpdateGameInstanceRec -
func (m *Domain) UpdateGameInstanceRec(rec *game_record.GameInstance) (*game_record.GameInstance, error) {
	l := m.Logger("UpdateGameInstanceRec")

	_, err := m.GetGameInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating game_instance record >%#v<", rec)

	if err := m.validateGameInstanceRecForUpdate(rec); err != nil {
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
