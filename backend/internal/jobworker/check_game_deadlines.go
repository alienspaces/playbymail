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

// CheckGameDeadlinesWorkerArgs defines the arguments for checking game deadlines
type CheckGameDeadlinesWorkerArgs struct {
	// No arguments needed - this is a periodic job
}

func (CheckGameDeadlinesWorkerArgs) Kind() string { return "check_game_deadlines" }

// CheckGameDeadlinesWorker checks for games that need turn processing
type CheckGameDeadlinesWorker struct {
	river.WorkerDefaults[CheckGameDeadlinesWorkerArgs]
	JobWorker
}

func NewCheckGameDeadlinesWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*CheckGameDeadlinesWorker, error) {
	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	return &CheckGameDeadlinesWorker{
		JobWorker: *jw,
	}, nil
}

func (w *CheckGameDeadlinesWorker) Work(ctx context.Context, j *river.Job[CheckGameDeadlinesWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("CheckGameDeadlinesWorker/Work")

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
		l.Error("CheckGameDeadlinesWorker job ID >%s< failed >%v<", strconv.FormatInt(j.ID, 10), err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type CheckGameDeadlinesDoWorkResult struct {
	GamesChecked int
	JobsQueued   int
	ProcessedAt  time.Time
}

func (w *CheckGameDeadlinesWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[CheckGameDeadlinesWorkerArgs]) (*CheckGameDeadlinesDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("CheckGameDeadlinesWorker/DoWork")

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

		// Create turn processing job
		turnJob := ProcessGameTurnWorkerArgs{
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

	l.Info("completed deadline check: checked >%d< games, queued >%d< jobs", len(instanceRecs), jobsQueued)

	return &CheckGameDeadlinesDoWorkResult{
		GamesChecked: len(instanceRecs),
		JobsQueued:   jobsQueued,
		ProcessedAt:  time.Now(),
	}, nil
}
