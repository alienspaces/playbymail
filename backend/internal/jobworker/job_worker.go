package jobworker

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
)

type JobWorker struct {
	corejobworker.JobWorker
}

func NewJobWorker(l logger.Logger, s storer.Storer) (*JobWorker, error) {

	jw, err := corejobworker.NewWorker(l, s)
	if err != nil {
		return nil, err
	}

	return &JobWorker{
		JobWorker: *jw,
	}, nil
}

// beginJob returns the river client from context and a fresh domain initialised with a new database transaction.
func (w *JobWorker) beginJob(ctx context.Context) (*river.Client[pgx.Tx], *domain.Domain, error) {
	l := w.JobWorker.Log.WithFunctionContext("beginJob")

	c, err := river.ClientFromContextSafely[pgx.Tx](ctx)
	if err != nil {
		l.Warn("failed getting client from context safely >%v<", err)
		return nil, nil, err
	}

	tx, err := w.JobWorker.Store.BeginTx()
	if err != nil {
		return nil, nil, err
	}

	m, err := domain.NewDomain(l, c)
	if err != nil {
		return nil, nil, err
	}

	err = m.Init(tx)
	if err != nil {
		return nil, nil, err
	}

	return c, m, nil
}
