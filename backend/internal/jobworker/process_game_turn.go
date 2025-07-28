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
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// ProcessGameTurnWorkerArgs defines the arguments for processing a game turn
type ProcessGameTurnWorkerArgs struct {
	GameInstanceID string `json:"game_instance_id"`
	TurnNumber     int    `json:"turn_number"`
}

func (ProcessGameTurnWorkerArgs) Kind() string { return "process_game_turn" }

// ProcessGameTurnWorker processes a single game turn
type ProcessGameTurnWorker struct {
	river.WorkerDefaults[ProcessGameTurnWorkerArgs]
	JobWorker
}

func NewProcessGameTurnWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*ProcessGameTurnWorker, error) {
	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	return &ProcessGameTurnWorker{
		JobWorker: *jw,
	}, nil
}

func (w *ProcessGameTurnWorker) Work(ctx context.Context, j *river.Job[ProcessGameTurnWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("ProcessGameTurnWorker/Work")

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
		l.Error("ProcessGameTurnWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type ProcessGameTurnDoWorkResult struct {
	GameInstanceID string
	TurnNumber     int
	ProcessedAt    time.Time
}

func (w *ProcessGameTurnWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[ProcessGameTurnWorkerArgs]) (*ProcessGameTurnDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("ProcessGameTurnWorker/DoWork")

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

	// TODO: Process player turns here
	// This would involve:
	// 1. Getting all player submissions for this turn
	// 2. Processing the game logic based on the game type
	// 3. Updating game state (character positions, items, etc.)
	// 4. Generating turn results for each player

	l.Info("processing turn logic for game instance >%s< turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber)

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

	return &ProcessGameTurnDoWorkResult{
		GameInstanceID: j.Args.GameInstanceID,
		TurnNumber:     j.Args.TurnNumber,
		ProcessedAt:    time.Now(),
	}, nil
}
