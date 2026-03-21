package jobworker

import (
	"context"
	"strconv"
	"time"

	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

type ExpirePendingSubscriptionsWorkerArgs struct{}

func (ExpirePendingSubscriptionsWorkerArgs) Kind() string {
	return "expire_pending_subscriptions"
}

type ExpirePendingSubscriptionsWorker struct {
	river.WorkerDefaults[ExpirePendingSubscriptionsWorkerArgs]
	JobWorker
}

func NewExpirePendingSubscriptionsWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*ExpirePendingSubscriptionsWorker, error) {
	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	return &ExpirePendingSubscriptionsWorker{
		JobWorker: *jw,
	}, nil
}

func (w *ExpirePendingSubscriptionsWorker) Work(ctx context.Context, j *river.Job[ExpirePendingSubscriptionsWorkerArgs]) error {
	l := w.Log.WithFunctionContext("ExpirePendingSubscriptionsWorker/Work")

	l.Info("running job ID >%s<", strconv.FormatInt(j.ID, 10))

	_, m, err := w.beginJob(ctx)
	if err != nil {
		return err
	}
	defer func() {
		m.Tx.Rollback(context.Background())
	}()

	expired, err := w.expireSubscriptions(m)
	if err != nil {
		l.Error("expire pending subscriptions job ID >%s< failed >%v<", strconv.FormatInt(j.ID, 10), err)
		return err
	}

	if expired > 0 {
		l.Info("expired >%d< pending subscriptions", expired)
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

func (w *ExpirePendingSubscriptionsWorker) expireSubscriptions(m *domain.Domain) (int, error) {
	l := w.Log.WithFunctionContext("ExpirePendingSubscriptionsWorker/expireSubscriptions")

	now := time.Now()

	recs, err := m.GetManyGameSubscriptionRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionStatus, Val: game_record.GameSubscriptionStatusPendingApproval},
			{Col: game_record.FieldGameSubscriptionPendingApprovalExpiresAt, Val: nulltime.FromTime(now), Op: coresql.OpLessThan},
		},
	})
	if err != nil {
		l.Warn("failed to get expired pending subscriptions >%v<", err)
		return 0, err
	}

	if len(recs) == 0 {
		return 0, nil
	}

	l.Info("found >%d< expired pending subscriptions", len(recs))

	expired := 0
	for _, rec := range recs {
		// Soft-delete any subscription instance links
		instanceRecs, err := m.GetGameSubscriptionInstanceRecsBySubscription(rec.ID)
		if err != nil {
			l.Warn("failed to get subscription instances for >%s< >%v<", rec.ID, err)
			continue
		}
		for _, instRec := range instanceRecs {
			if err := m.DeleteGameSubscriptionInstanceRec(instRec.ID); err != nil {
				l.Warn("failed to soft-delete subscription instance >%s< >%v<", instRec.ID, err)
			}
		}

		// Revoke the subscription
		rec.Status = game_record.GameSubscriptionStatusRevoked
		if _, err := m.UpdateGameSubscriptionRec(rec); err != nil {
			l.Warn("failed to revoke subscription >%s< >%v<", rec.ID, err)
			continue
		}

		expired++
		l.Info("expired subscription >%s< and removed >%d< instance links", rec.ID, len(instanceRecs))
	}

	return expired, nil
}
