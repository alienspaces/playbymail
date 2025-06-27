package jobworker

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
)

type JobWorkerExampleArgs struct {
	ExampleWorkRecID string
}

func (JobWorkerExampleArgs) Kind() string { return "processingreport" }
func (JobWorkerExampleArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: jobqueue.QueueDefault,
	}
}

type JobWorkerExampleWorker struct {
	river.WorkerDefaults[JobWorkerExampleArgs]
	emailClient emailer.Emailer
	JobWorker
}

func NewJobWorkerExampleWorker(l logger.Logger, s storer.Storer, e emailer.Emailer) (*JobWorkerExampleWorker, error) {
	l = l.WithPackageContext("JobWorkerExampleWorker")

	l.Info("Instantiation JobWorkerExampleWorker")

	jw, err := NewJobWorker(l, s)
	if err != nil {
		return nil, err
	}

	return &JobWorkerExampleWorker{
		JobWorker:   *jw,
		emailClient: e,
	}, nil
}

func (w *JobWorkerExampleWorker) Work(ctx context.Context, j *river.Job[JobWorkerExampleArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("JobWorkerExampleWorker/Work")

	l.Info("Running job ID >%s< Args >%#v<", strconv.FormatInt(j.ID, 10), j.Args)

	c, m, err := w.beginJob(ctx)
	if err != nil {
		return err
	}
	defer func() {
		m.Tx.Rollback(context.Background())
	}()

	_, err = w.DoWork(ctx, m, c, j)
	if err != nil {
		l.Error("JobWorkerExample job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type JobWorkerExampleDoWorkResult struct {
	RecordCount int
}

func (w *JobWorkerExampleWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[JobWorkerExampleArgs]) (*JobWorkerExampleDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("JobWorkerExampleWorker/DoWork")

	l.Info("JobWorkerExample report work record ID >%s<", j.Args.ExampleWorkRecID)

	return &JobWorkerExampleDoWorkResult{
		RecordCount: 1,
	}, nil
}
