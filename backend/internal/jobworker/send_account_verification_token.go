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

type SendAccountVerificationTokenWorkerArgs struct {
	AccountID string
}

func (SendAccountVerificationTokenWorkerArgs) Kind() string { return "send-account-verification-token" }
func (SendAccountVerificationTokenWorkerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: jobqueue.QueueDefault,
	}
}

type SendAccountVerificationTokenWorker struct {
	river.WorkerDefaults[SendAccountVerificationTokenWorkerArgs]
	emailClient emailer.Emailer
	JobWorker
}

func NewSendAccountVerificationTokenWorker(l logger.Logger, s storer.Storer, e emailer.Emailer) (*SendAccountVerificationTokenWorker, error) {
	l = l.WithPackageContext("SendAccountVerificationTokenWorker")

	l.Info("Instantiation SendAccountVerificationTokenWorker")

	jw, err := NewJobWorker(l, s)
	if err != nil {
		return nil, err
	}

	return &SendAccountVerificationTokenWorker{
		JobWorker:   *jw,
		emailClient: e,
	}, nil
}

func (w *SendAccountVerificationTokenWorker) Work(ctx context.Context, j *river.Job[SendAccountVerificationTokenWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("SendAccountVerificationTokenWorker/Work")

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
		l.Error("SendAccountVerificationTokenWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type SendAccountVerificationTokenDoWorkResult struct {
	RecordCount int
}

func (w *SendAccountVerificationTokenWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[SendAccountVerificationTokenWorkerArgs]) (*SendAccountVerificationTokenDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("SendAccountVerificationTokenWorker/DoWork")

	l.Info("SendAccountVerificationTokenWorker report work record ID >%s<", j.Args.AccountID)

	return &SendAccountVerificationTokenDoWorkResult{
		RecordCount: 1,
	}, nil
}
