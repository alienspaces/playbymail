package jobworker

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// JobWorker provides the base functionality for all job workers in the system
type JobWorker struct {
	corejobworker.JobWorker
	Config config.Config
}

func NewJobWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*JobWorker, error) {

	jw, err := corejobworker.NewWorker(l, s)
	if err != nil {
		return nil, err
	}

	return &JobWorker{
		JobWorker: *jw,
		Config:    cfg,
	}, nil
}

// beginJob returns the river client from context and a fresh domain model
// initialised with a new database transaction.
func (w *JobWorker) beginJob(ctx context.Context) (*river.Client[pgx.Tx], *domain.Domain, error) {
	l := w.Log.WithFunctionContext("beginJob")

	c, err := river.ClientFromContextSafely[pgx.Tx](ctx)
	if err != nil {
		l.Warn("failed getting client from context safely >%v<", err)
		return nil, nil, err
	}

	tx, err := w.Store.BeginTx()
	if err != nil {
		return nil, nil, err
	}

	m, err := domain.NewDomain(l, w.Config)
	if err != nil {
		return nil, nil, err
	}

	err = m.Init(tx)
	if err != nil {
		return nil, nil, err
	}

	return c, m, nil
}
