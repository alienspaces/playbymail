package domain

import (
	"fmt"
	"time"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) GetGameInstanceRec(recID string, lock *sql.Lock) (*game_record.GameInstance, error) {
	r := m.GameInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateGameInstanceRec(rec *game_record.GameInstance) (*game_record.GameInstance, error) {
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
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateGameInstanceRec(next *game_record.GameInstance) (*game_record.GameInstance, error) {
	r := m.GameInstanceRepository()
	next, err := r.UpdateOne(next)
	if err != nil {
		return next, err
	}
	return next, nil
}

func (m *Domain) DeleteGameInstanceRec(recID string) error {
	r := m.GameInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) ValidateGameInstance(rec *game_record.GameInstance) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) RemoveGameInstanceRec(recID string) error {
	r := m.GameInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return err
	}
	return nil
}

// Game Runtime Management Functions

// StartGameInstance starts a game instance and sets up the first turn
func (m *Domain) StartGameInstance(instanceID string) (*game_record.GameInstance, error) {
	l := m.Logger("StartGameInstance")

	instance, err := m.GetGameInstanceRec(instanceID, sql.ForUpdateNoWait)
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

	instance, err := m.GetGameInstanceRec(instanceID, sql.ForUpdateNoWait)
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

	instance, err := m.GetGameInstanceRec(instanceID, sql.ForUpdateNoWait)
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

	instanceRec, err := m.GetGameInstanceRec(instanceID, sql.ForUpdateNoWait)
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

	instanceRec, err := m.GetGameInstanceRec(instanceID, sql.ForUpdateNoWait)
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

	instanceRec, err := m.GetGameInstanceRec(instanceID, sql.ForUpdateNoWait)
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

// GetGameInstanceRecsByStatus gets all game instances with a specific status
func (m *Domain) GetGameInstanceRecsByStatus(status string) ([]*game_record.GameInstance, error) {
	l := m.Logger("GetGameInstanceRecsByStatus")

	l.Info("getting game instances with status >%s<", status)

	r := m.GameInstanceRepository()
	opts := &sql.Options{
		Params: []sql.Param{
			{
				Col: game_record.FieldGameInstanceStatus,
				Val: status,
			},
		},
	}
	instanceRecs, err := r.GetMany(opts)
	if err != nil {
		l.Warn("failed to get game instances with status >%s< >%v<", status, err)
		return nil, err
	}

	l.Info("returning >%d< game instances with status >%s<", len(instanceRecs), status)

	return instanceRecs, nil
}

// GetGameInstanceRecsNeedingTurnProcessing gets game instances that need turn processing
func (m *Domain) GetGameInstanceRecsNeedingTurnProcessing() ([]*game_record.GameInstance, error) {
	l := m.Logger("GetGameInstanceRecsNeedingTurnProcessing")

	l.Info("getting game instances needing turn processing")

	r := m.GameInstanceRepository()
	now := time.Now().UTC()

	opts := &sql.Options{
		Params: []sql.Param{
			{
				Col: game_record.FieldGameInstanceStatus,
				Val: game_record.GameInstanceStatusStarted,
			},
			{
				Col: game_record.FieldGameInstanceNextTurnDueAt,
				Val: nulltime.FromTime(now),
				Op:  sql.OpLessThanEqual,
			},
		},
	}

	instanceRecs, err := r.GetMany(opts)
	if err != nil {
		l.Warn("failed to get game instances needing turn processing >%v<", err)
		return nil, err
	}

	l.Info("returning >%d< game instances needing turn processing", len(instanceRecs))

	return instanceRecs, nil
}
