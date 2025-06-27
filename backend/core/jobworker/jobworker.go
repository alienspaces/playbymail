package jobworker

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

type JobWorker struct {
	Log   logger.Logger
	Store storer.Storer
}

func NewWorker(l logger.Logger, s storer.Storer) (*JobWorker, error) {

	if l == nil {
		return nil, fmt.Errorf("failed new job worker, missing logger")
	}
	if s == nil {
		return nil, fmt.Errorf("failed new job worker, missing storer")
	}

	l, err := l.NewInstance()
	if err != nil {
		return nil, err
	}

	return &JobWorker{
		Log:   l,
		Store: s,
	}, nil
}

// CompleteJob must be called once a job has successfully completed otherwise river will kick the job off again.
func CompleteJob[TArgs river.JobArgs](ctx context.Context, tx pgx.Tx, job *river.Job[TArgs]) error {
	_, err := river.JobCompleteTx[*riverpgxv5.Driver](ctx, tx, job)
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
