package jobworker

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	nulltime "gitlab.com/alienspaces/playbymail/core/nulltime"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// GameTurnQueueingWorkerArgs defines the arguments for queuing game turn processing
type GameTurnQueueingWorkerArgs struct {
	// No arguments needed - this is a periodic job
}

func (GameTurnQueueingWorkerArgs) Kind() string { return "queue_game_turn_processing" }

// GameTurnQueueingWorker queues turn processing jobs for games that need them
type GameTurnQueueingWorker struct {
	river.WorkerDefaults[GameTurnQueueingWorkerArgs]
	JobWorker
}

func NewGameTurnQueueingWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*GameTurnQueueingWorker, error) {
	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	return &GameTurnQueueingWorker{
		JobWorker: *jw,
	}, nil
}

func (w *GameTurnQueueingWorker) Work(ctx context.Context, j *river.Job[GameTurnQueueingWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("GameTurnQueueingWorker/Work")

	l.Info("running job ID >%s<", strconv.FormatInt(j.ID, 10))

	c, m, err := w.beginJob(ctx)
	if err != nil {
		return err
	}
	defer func() {
		m.Tx.Rollback(context.Background())
	}()

	_, err = w.DoWork(ctx, m, c, j)
	if err != nil {
		l.Error("GameTurnQueueingWorker job ID >%s< failed >%v<", strconv.FormatInt(j.ID, 10), err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type GameTurnQueueingDoWorkResult struct {
	GamesChecked int
	JobsQueued   int
	ProcessedAt  time.Time
}

func (w *GameTurnQueueingWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[GameTurnQueueingWorkerArgs]) (*GameTurnQueueingDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("GameTurnQueueingWorker/DoWork")

	l.Info("checking for games that need turn processing")

	// Get all game instances that need turn processing
	instanceRecs, err := w.getGameInstanceRecsNeedingTurnProcessing(m)
	if err != nil {
		l.Warn("failed to get game instances needing turn processing >%v<", err)
		return nil, err
	}

	l.Info("found >%d< game instances needing turn processing", len(instanceRecs))

	jobsQueued := 0

	// Queue turn processing jobs for each game instance
	for _, instanceRec := range instanceRecs {
		l.Info("queueing turn processing for game instance >%s< turn >%d<", instanceRec.ID, instanceRec.CurrentTurn)

		// Create execute game turn processing job for this game instance
		turnProcessingJob := GameTurnProcessingWorkerArgs{
			GameInstanceID: instanceRec.ID,
			TurnNumber:     instanceRec.CurrentTurn,
		}

		// Queue the job
		_, err := c.Insert(ctx, turnProcessingJob, &river.InsertOpts{
			Queue: jobqueue.QueueGame,
		})
		if err != nil {
			l.Warn("failed to queue turn processing job for instance >%s< >%v<", instanceRec.ID, err)
			continue
		}

		jobsQueued++
		l.Info("queued turn processing job for game instance >%s< turn >%d<", instanceRec.ID, instanceRec.CurrentTurn)
	}

	l.Info("completed turn processing queue check: checked >%d< games, queued >%d< jobs", len(instanceRecs), jobsQueued)

	return &GameTurnQueueingDoWorkResult{
		GamesChecked: len(instanceRecs),
		JobsQueued:   jobsQueued,
		ProcessedAt:  time.Now(),
	}, nil
}

// getGameInstanceRecsNeedingTurnProcessing gets game instances that need turn processing
//
// Returns game instances that meet both criteria:
// 1. Status is "started" (active games only)
// 2. NextTurnDueAt <= current time (deadline has passed)
//
// Used by the queue worker to create individual turn processing jobs.
func (w *GameTurnQueueingWorker) getGameInstanceRecsNeedingTurnProcessing(m *domain.Domain) ([]*game_record.GameInstance, error) {
	l := w.JobWorker.Log.WithFunctionContext("GameTurnQueueingWorker/getGameInstanceRecsNeedingTurnProcessing")

	l.Info("getting game instances needing turn processing")

	r := m.GameInstanceRepository()
	now := time.Now().UTC()

	opts := &coresql.Options{
		Params: []coresql.Param{
			{
				Col: game_record.FieldGameInstanceStatus,
				Val: game_record.GameInstanceStatusStarted,
			},
			{
				Col: game_record.FieldGameInstanceNextTurnDueAt,
				Val: nulltime.FromTime(now),
				Op:  coresql.OpLessThanEqual,
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
