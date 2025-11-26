// Game turn processing job workers handle the processing of game turns for all game types.
//
// IMPORTANT: This system assumes all players have already joined the game and have game assets created.
// New player onboarding (join game) should be handled separately through API endpoints, not job workers.
//
// To add new game types:
//  1. Create processor in internal/jobworker/[game_type]/
//  2. Register in initializeProcessors() function below
//
// To add new turn sheet types (per game type):
//  1. Create processor in internal/jobworker/[game_type]/turn_sheet_processor/
//  2. Register in that game type's initializeTurnSheetProcessors() function
package jobworker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// GameTurnProcessingWorker processes a game instance turn
//
// To add new game types:
//   - Create processor in internal/jobworker/[game_type]/
//   - Register in initializeProcessors() function below

// GameTurnProcessingWorkerArgs defines the arguments for processing a game instance turn
type GameTurnProcessingWorkerArgs struct {
	GameInstanceID string `json:"game_instance_id"`
	TurnNumber     int    `json:"turn_number"`
}

func (GameTurnProcessingWorkerArgs) Kind() string { return "game_turn_processing" }

// GameTurnProcessor defines the interface for processing and generating turn sheets for different game types
type GameTurnProcessor interface {
	// ProcessTurnSheets processes all turn sheets for a game instance
	ProcessTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance) error

	// CreateTurnSheets generates all turn sheets for a game instance
	CreateTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance) ([]*game_record.GameTurnSheet, error)
}

// GameTurnProcessingWorker processes a game instance turn
type GameTurnProcessingWorker struct {
	river.WorkerDefaults[GameTurnProcessingWorkerArgs]
	JobWorker
}

func NewGameTurnProcessingWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*GameTurnProcessingWorker, error) {
	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	return &GameTurnProcessingWorker{
		JobWorker: *jw,
	}, nil
}

func (w *GameTurnProcessingWorker) Work(ctx context.Context, j *river.Job[GameTurnProcessingWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("GameTurnProcessingWorker/Work")

	l.Info("running job ID >%s< Args >%#v<", strconv.FormatInt(j.ID, 10), j.Args)

	c, d, err := w.beginJob(ctx)
	if err != nil {
		return err
	}
	defer func() {
		d.Tx.Rollback(context.Background())
	}()

	_, err = w.DoWork(ctx, d, c, j)
	if err != nil {
		l.Error("GameTurnProcessingWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, d.Tx, j)
}

type GameTurnProcessingDoWorkResult struct {
	GameInstanceID string
	TurnNumber     int
	ProcessedAt    time.Time
}

func (w *GameTurnProcessingWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[GameTurnProcessingWorkerArgs]) (*GameTurnProcessingDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("GameTurnProcessingWorker/DoWork")

	l.Info("processing game turn for instance >%s< turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber)

	// Initialize all game type processors
	processors, err := w.initializeProcessors(l, m)
	if err != nil {
		l.Warn("failed to initialize processors >%v<", err)
		return nil, err
	}

	// Get the game instance
	gameInstanceRec, err := m.GetGameInstanceRec(j.Args.GameInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance ID >%s<; cannot process game turn >%v<", j.Args.GameInstanceID, err)
		return nil, err
	}

	// Verify we're processing the correct turn
	if gameInstanceRec.CurrentTurn != j.Args.TurnNumber {
		l.Warn("turn number mismatch for game instance ID >%s<: expected >%d< but instance is at turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber, gameInstanceRec.CurrentTurn)
		return nil, fmt.Errorf("turn number mismatch for game instance ID >%s<", j.Args.GameInstanceID)
	}

	// Begin turn processing
	gameInstanceRec, err = m.BeginTurnProcessing(j.Args.GameInstanceID)
	if err != nil {
		l.Warn("failed to begin turn processing for game instance ID >%s<; cannot process game turn >%v<", j.Args.GameInstanceID, err)
		return nil, err
	}

	// Process player turns based on game type
	l.Info("processing turn logic for game instance ID >%s< turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber)

	// Get the game to determine the game type
	gameRec, err := m.GetGameRec(gameInstanceRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game ID >%s< for game instance ID >%s<; cannot process game turn >%v<", gameInstanceRec.GameID, j.Args.GameInstanceID, err)
		return nil, err
	}

	// Get the appropriate processor for this game type
	processor, exists := processors[gameRec.GameType]
	if !exists {
		l.Warn("unsupported game type >%s< for game instance ID >%s<; cannot process game turn >%v<", gameRec.GameType, j.Args.GameInstanceID, err)
		return nil, fmt.Errorf("unsupported game type: %s for game instance ID >%s<", gameRec.GameType, j.Args.GameInstanceID)
	}

	// Process turn sheets using the game-specific processor
	err = processor.ProcessTurnSheets(ctx, gameInstanceRec)
	if err != nil {
		l.Warn("failed to process game turn for game instance ID >%s< turn >%d<; cannot process game turn >%v<", j.Args.GameInstanceID, j.Args.TurnNumber, err)
		return nil, err
	}

	// TODO: Generate post run processing jobs such as notifications, etc.

	// Complete the turn
	gameInstanceRec, err = m.CompleteTurn(j.Args.GameInstanceID)
	if err != nil {
		l.Warn("failed to complete turn for game instance ID >%s< turn >%d<; cannot process game turn >%v<", j.Args.GameInstanceID, j.Args.TurnNumber, err)
		return nil, err
	}

	l.Info("completed turn processing for game instance >%s< turn >%d<", gameInstanceRec.ID, j.Args.TurnNumber)

	// Generate new turn sheets for the next turn
	l.Info("generating new turn sheets for game instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	// Generate turn sheets using the same processor
	_, err = processor.CreateTurnSheets(ctx, gameInstanceRec)
	if err != nil {
		l.Warn("failed to generate new turn sheets for game instance ID >%s< turn >%d<; cannot process game turn >%v<", j.Args.GameInstanceID, j.Args.TurnNumber, err)
		return nil, err
	}

	return &GameTurnProcessingDoWorkResult{
		GameInstanceID: j.Args.GameInstanceID,
		TurnNumber:     j.Args.TurnNumber,
		ProcessedAt:    time.Now(),
	}, nil
}

// initializeProcessors creates and registers all available game type processors
func (w *GameTurnProcessingWorker) initializeProcessors(l logger.Logger, d *domain.Domain) (map[string]GameTurnProcessor, error) {
	processors := make(map[string]GameTurnProcessor)

	// Register adventure game processor
	adventureProcessor, err := adventure_game.NewAdventureGame(l, d)
	if err != nil {
		return nil, err
	}
	processors[game_record.GameTypeAdventure] = adventureProcessor

	// TODO: Add new game type processors here
	// Example: processors[game_record.GameTypeStrategy] = strategyProcessor
	// Example: processors[game_record.GameTypePuzzle] = puzzleProcessor

	return processors, nil
}
