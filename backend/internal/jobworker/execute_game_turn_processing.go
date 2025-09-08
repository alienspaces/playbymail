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
	"gitlab.com/alienspaces/playbymail/internal/jobworker/adventure"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// ExecuteGameTurnProcessingWorkerArgs defines the arguments for executing game turn processing
type ExecuteGameTurnProcessingWorkerArgs struct {
	GameInstanceID string `json:"game_instance_id"`
	TurnNumber     int    `json:"turn_number"`
}

func (ExecuteGameTurnProcessingWorkerArgs) Kind() string { return "execute_game_turn_processing" }

// ExecuteGameTurnProcessingWorker executes a single game turn processing
type ExecuteGameTurnProcessingWorker struct {
	river.WorkerDefaults[ExecuteGameTurnProcessingWorkerArgs]
	JobWorker
}

func NewExecuteGameTurnProcessingWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*ExecuteGameTurnProcessingWorker, error) {
	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	return &ExecuteGameTurnProcessingWorker{
		JobWorker: *jw,
	}, nil
}

func (w *ExecuteGameTurnProcessingWorker) Work(ctx context.Context, j *river.Job[ExecuteGameTurnProcessingWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("ExecuteGameTurnProcessingWorker/Work")

	l.Info("running job ID >%s< Args >%#v<", strconv.FormatInt(j.ID, 10), j.Args)

	c, m, err := w.beginJob(ctx)
	if err != nil {
		return err
	}
	defer func() {
		m.Tx.Rollback(context.Background())
	}()

	_, err = w.DoWork(ctx, m, c, j)
	if err != nil {
		l.Error("ExecuteGameTurnProcessingWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type ExecuteGameTurnProcessingDoWorkResult struct {
	GameInstanceID string
	TurnNumber     int
	ProcessedAt    time.Time
}

func (w *ExecuteGameTurnProcessingWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[ExecuteGameTurnProcessingWorkerArgs]) (*ExecuteGameTurnProcessingDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("ExecuteGameTurnProcessingWorker/DoWork")

	l.Info("processing game turn for instance >%s< turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber)

	// Get the game instance
	instance, err := m.GetGameInstanceRec(j.Args.GameInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%v<", err)
		return nil, err
	}

	// Verify we're processing the correct turn
	if instance.CurrentTurn != j.Args.TurnNumber {
		l.Warn("turn number mismatch: expected >%d< but instance is at turn >%d<", j.Args.TurnNumber, instance.CurrentTurn)
		return nil, fmt.Errorf("turn number mismatch")
	}

	// Begin turn processing
	instance, err = m.BeginTurnProcessing(j.Args.GameInstanceID)
	if err != nil {
		l.Warn("failed to begin turn processing >%v<", err)
		return nil, err
	}

	// Process player turns based on game type
	l.Info("processing turn logic for game instance >%s< turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber)

	// Get the game to determine the game type
	game, err := m.GetGameRec(instance.GameID, nil)
	if err != nil {
		l.Warn("failed to get game >%v<", err)
		return nil, err
	}

	// Route to appropriate game type processor
	switch game.GameType {
	case game_record.GameTypeAdventure:
		adventureProcessor := adventure.NewAdventureGameTurnProcessor(l, m)
		err = adventureProcessor.ProcessTurn(ctx, j.Args.GameInstanceID, j.Args.TurnNumber)
		if err != nil {
			l.Warn("failed to process adventure game turn >%v<", err)
			return nil, err
		}
	default:
		l.Warn("unsupported game type >%s<", game.GameType)
		return nil, fmt.Errorf("unsupported game type: %s", game.GameType)
	}

	// For now, we'll just complete the turn
	instance, err = m.CompleteTurn(j.Args.GameInstanceID)
	if err != nil {
		l.Warn("failed to complete turn >%v<", err)
		return nil, err
	}

	// TODO: Generate and queue print/mail jobs for turn results
	// This would involve:
	// 1. Creating print jobs for each player's turn results
	// 2. Queueing mail jobs to send the results

	l.Info("completed turn processing for game instance >%s< turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber)

	return &ExecuteGameTurnProcessingDoWorkResult{
		GameInstanceID: j.Args.GameInstanceID,
		TurnNumber:     j.Args.TurnNumber,
		ProcessedAt:    time.Now(),
	}, nil
}
