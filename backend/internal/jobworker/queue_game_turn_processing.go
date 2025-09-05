package jobworker

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// QueueGameTurnProcessingWorkerArgs defines the arguments for queuing game turn processing
type QueueGameTurnProcessingWorkerArgs struct {
	// No arguments needed - this is a periodic job
}

func (QueueGameTurnProcessingWorkerArgs) Kind() string { return "queue_game_turn_processing" }

// QueueGameTurnProcessingWorker queues turn processing jobs for games that need them
type QueueGameTurnProcessingWorker struct {
	river.WorkerDefaults[QueueGameTurnProcessingWorkerArgs]
	JobWorker
}

func NewQueueGameTurnProcessingWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*QueueGameTurnProcessingWorker, error) {
	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	return &QueueGameTurnProcessingWorker{
		JobWorker: *jw,
	}, nil
}

func (w *QueueGameTurnProcessingWorker) Work(ctx context.Context, j *river.Job[QueueGameTurnProcessingWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("QueueGameTurnProcessingWorker/Work")

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
		l.Error("QueueGameTurnProcessingWorker job ID >%s< failed >%v<", strconv.FormatInt(j.ID, 10), err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type QueueGameTurnProcessingDoWorkResult struct {
	GamesChecked int
	JobsQueued   int
	ProcessedAt  time.Time
}

func (w *QueueGameTurnProcessingWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[QueueGameTurnProcessingWorkerArgs]) (*QueueGameTurnProcessingDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("QueueGameTurnProcessingWorker/DoWork")

	l.Info("checking for games that need turn processing")

	// Get all game instances that need turn processing
	instanceRecs, err := m.GetGameInstanceRecsNeedingTurnProcessing()
	if err != nil {
		l.Warn("failed to get game instances needing turn processing >%v<", err)
		return nil, err
	}

	l.Info("found >%d< game instances needing turn processing", len(instanceRecs))

	jobsQueued := 0

	// Queue turn processing jobs for each instance
	for _, instanceRec := range instanceRecs {
		l.Info("queueing turn processing for game instance >%s< turn >%d<", instanceRec.ID, instanceRec.CurrentTurn)

		// Create execute game turn processing job for this instance
		turnJob := ExecuteGameTurnProcessingWorkerArgs{
			GameInstanceID: instanceRec.ID,
			TurnNumber:     instanceRec.CurrentTurn,
		}

		// Queue the job
		_, err := c.Insert(ctx, turnJob, &river.InsertOpts{
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

	return &QueueGameTurnProcessingDoWorkResult{
		GamesChecked: len(instanceRecs),
		JobsQueued:   jobsQueued,
		ProcessedAt:  time.Now(),
	}, nil
}
