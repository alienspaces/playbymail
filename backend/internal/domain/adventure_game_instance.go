package domain

import (
	"fmt"
	"time"

	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (m *Domain) GetAdventureGameInstanceRec(recID string, lock *sql.Lock) (*record.AdventureGameInstance, error) {
	r := m.AdventureGameInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateAdventureGameInstanceRec(rec *record.AdventureGameInstance) (*record.AdventureGameInstance, error) {
	r := m.AdventureGameInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateAdventureGameInstanceRec(next *record.AdventureGameInstance) (*record.AdventureGameInstance, error) {
	r := m.AdventureGameInstanceRepository()
	next, err := r.UpdateOne(next)
	if err != nil {
		return next, err
	}
	return next, nil
}

func (m *Domain) DeleteAdventureGameInstanceRec(recID string) error {
	r := m.AdventureGameInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) ValidateAdventureGameInstance(rec *record.AdventureGameInstance) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) RemoveAdventureGameInstanceRec(recID string) error {
	r := m.AdventureGameInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return err
	}
	return nil
}

// Game Runtime Management Functions

// StartGameInstance starts a game instance and sets up the first turn
func (m *Domain) StartGameInstance(instanceID string) (*record.AdventureGameInstance, error) {
	l := m.Logger("StartGameInstance")

	instance, err := m.GetAdventureGameInstanceRec(instanceID, sql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status != record.GameInstanceStatusCreated {
		return nil, fmt.Errorf("game instance must be in 'created' status to start")
	}

	now := time.Now()
	instance.Status = record.GameInstanceStatusStarting
	instance.StartedAt = &now
	instance.CurrentTurn = 0

	// Set the first turn deadline
	nextDeadline := now.Add(time.Duration(instance.TurnDeadlineHours) * time.Hour)
	instance.NextTurnDeadline = &nextDeadline

	instance, err = m.UpdateAdventureGameInstanceRec(instance)
	if err != nil {
		l.Warn("failed updating game instance to starting status >%v<", err)
		return nil, err
	}

	l.Info("started game instance >%s< with first turn deadline >%v<", instanceID, nextDeadline)
	return instance, nil
}

// BeginTurnProcessing starts processing the current turn
func (m *Domain) BeginTurnProcessing(instanceID string) (*record.AdventureGameInstance, error) {
	l := m.Logger("BeginTurnProcessing")

	instance, err := m.GetAdventureGameInstanceRec(instanceID, sql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status != record.GameInstanceStatusRunning && instance.Status != record.GameInstanceStatusStarting {
		return nil, fmt.Errorf("game instance must be running or starting to process turns")
	}

	instance.Status = record.GameInstanceStatusRunning
	now := time.Now()
	instance.LastTurnProcessedAt = &now

	instance, err = m.UpdateAdventureGameInstanceRec(instance)
	if err != nil {
		l.Warn("failed updating game instance for turn processing >%v<", err)
		return nil, err
	}

	l.Info("began turn processing for game instance >%s< turn >%d<", instanceID, instance.CurrentTurn)
	return instance, nil
}

// CompleteTurn advances the game to the next turn
func (m *Domain) CompleteTurn(instanceID string) (*record.AdventureGameInstance, error) {
	l := m.Logger("CompleteTurn")

	instance, err := m.GetAdventureGameInstanceRec(instanceID, sql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status != record.GameInstanceStatusRunning {
		return nil, fmt.Errorf("game instance must be running to complete turns")
	}

	// Check if we've reached max turns
	if instance.MaxTurns != nil && instance.CurrentTurn >= *instance.MaxTurns {
		instance.Status = record.GameInstanceStatusCompleted
		now := time.Now()
		instance.CompletedAt = &now
		l.Info("game instance >%s< completed after reaching max turns >%d<", instanceID, *instance.MaxTurns)
	} else {
		// Advance to next turn
		instance.CurrentTurn++
		now := time.Now()
		nextDeadline := now.Add(time.Duration(instance.TurnDeadlineHours) * time.Hour)
		instance.NextTurnDeadline = &nextDeadline
		l.Info("advanced game instance >%s< to turn >%d< with deadline >%v<", instanceID, instance.CurrentTurn, nextDeadline)
	}

	instance, err = m.UpdateAdventureGameInstanceRec(instance)
	if err != nil {
		l.Warn("failed updating game instance after turn completion >%v<", err)
		return nil, err
	}

	return instance, nil
}

// PauseGameInstance pauses a running game instance
func (m *Domain) PauseGameInstance(instanceID string) (*record.AdventureGameInstance, error) {
	l := m.Logger("PauseGameInstance")

	instance, err := m.GetAdventureGameInstanceRec(instanceID, sql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status != record.GameInstanceStatusRunning {
		return nil, fmt.Errorf("game instance must be running to pause")
	}

	instance.Status = record.GameInstanceStatusPaused

	instance, err = m.UpdateAdventureGameInstanceRec(instance)
	if err != nil {
		l.Warn("failed updating game instance to paused status >%v<", err)
		return nil, err
	}

	l.Info("paused game instance >%s<", instanceID)
	return instance, nil
}

// ResumeGameInstance resumes a paused game instance
func (m *Domain) ResumeGameInstance(instanceID string) (*record.AdventureGameInstance, error) {
	l := m.Logger("ResumeGameInstance")

	instance, err := m.GetAdventureGameInstanceRec(instanceID, sql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status != record.GameInstanceStatusPaused {
		return nil, fmt.Errorf("game instance must be paused to resume")
	}

	instance.Status = record.GameInstanceStatusRunning

	instance, err = m.UpdateAdventureGameInstanceRec(instance)
	if err != nil {
		l.Warn("failed updating game instance to running status >%v<", err)
		return nil, err
	}

	l.Info("resumed game instance >%s<", instanceID)
	return instance, nil
}

// CancelGameInstance cancels a game instance
func (m *Domain) CancelGameInstance(instanceID string) (*record.AdventureGameInstance, error) {
	l := m.Logger("CancelGameInstance")

	instance, err := m.GetAdventureGameInstanceRec(instanceID, sql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if instance.Status == record.GameInstanceStatusCompleted || instance.Status == record.GameInstanceStatusCancelled {
		return nil, fmt.Errorf("game instance is already completed or cancelled")
	}

	instance.Status = record.GameInstanceStatusCancelled
	now := time.Now()
	instance.CompletedAt = &now

	instance, err = m.UpdateAdventureGameInstanceRec(instance)
	if err != nil {
		l.Warn("failed updating game instance to cancelled status >%v<", err)
		return nil, err
	}

	l.Info("cancelled game instance >%s<", instanceID)
	return instance, nil
}

// GetGameInstancesByStatus gets all game instances with a specific status
func (m *Domain) GetGameInstancesByStatus(status string) ([]*record.AdventureGameInstance, error) {
	r := m.AdventureGameInstanceRepository()
	opts := &sql.Options{
		Params: []sql.Param{
			{
				Col: "status",
				Val: status,
			},
		},
	}
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

// GetGameInstancesNeedingTurnProcessing gets game instances that need turn processing
func (m *Domain) GetGameInstancesNeedingTurnProcessing() ([]*record.AdventureGameInstance, error) {
	r := m.AdventureGameInstanceRepository()
	now := time.Now()

	opts := &sql.Options{
		Params: []sql.Param{
			{
				Col: "status",
				Val: record.GameInstanceStatusRunning,
			},
			{
				Col: "next_turn_deadline",
				Val: now,
				Op:  sql.OpLessThanEqual,
			},
		},
	}
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}
