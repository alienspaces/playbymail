package jobworker

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// SendGameSubscriptionApprovalEmailWorkerArgs defines the job payload for sending approval emails
// when a player joins a game via a turn sheet upload.
type SendGameSubscriptionApprovalEmailWorkerArgs struct {
	GameSubscriptionID string
}

func (SendGameSubscriptionApprovalEmailWorkerArgs) Kind() string {
	return "send-game-subscription-approval-email"
}

func (SendGameSubscriptionApprovalEmailWorkerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: jobqueue.QueueDefault}
}

// SendGameSubscriptionApprovalEmailWorker sends an email containing the approval link for a
// pending game subscription created from a join game turn sheet upload.
type SendGameSubscriptionApprovalEmailWorker struct {
	river.WorkerDefaults[SendGameSubscriptionApprovalEmailWorkerArgs]
	emailClient emailer.Emailer
	JobWorker
}

func NewSendGameSubscriptionApprovalEmailWorker(l logger.Logger, cfg config.Config, s storer.Storer, e emailer.Emailer) (*SendGameSubscriptionApprovalEmailWorker, error) {
	l = l.WithPackageContext("SendGameSubscriptionApprovalEmailWorker")

	l.Info("instantiating SendGameSubscriptionApprovalEmailWorker")

	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	if e == nil {
		l.Warn("email client is nil, assuming registration-only instantiation")
	}

	if cfg.TemplatesPath == "" {
		return nil, fmt.Errorf("templates path is empty")
	}

	// Check the templates path actually exists
	l.Info("templates path >%s<", cfg.TemplatesPath)

	if _, err := os.Stat(cfg.TemplatesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("templates path does not exist >%s<", cfg.TemplatesPath)
	}

	return &SendGameSubscriptionApprovalEmailWorker{
		JobWorker:   *jw,
		emailClient: e,
	}, nil
}

func (w *SendGameSubscriptionApprovalEmailWorker) Work(ctx context.Context, j *river.Job[SendGameSubscriptionApprovalEmailWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("SendGameSubscriptionApprovalEmailWorker/Work")

	l.Info("running job ID >%s< Args >%#v<", strconv.FormatInt(j.ID, 10), j.Args)

	if w.emailClient == nil {
		return fmt.Errorf("email client is nil")
	}

	c, m, err := w.beginJob(ctx)
	if err != nil {
		return err
	}
	defer func() {
		m.Tx.Rollback(context.Background())
	}()

	_, err = w.DoWork(ctx, m, c, j)
	if err != nil {
		l.Error("SendGameSubscriptionApprovalEmailWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

// SendGameSubscriptionApprovalEmailDoWorkResult summarises the work carried out by the worker.
type SendGameSubscriptionApprovalEmailDoWorkResult struct {
	RecordCount int
}

func (w *SendGameSubscriptionApprovalEmailWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[SendGameSubscriptionApprovalEmailWorkerArgs]) (*SendGameSubscriptionApprovalEmailDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("SendGameSubscriptionApprovalEmailWorker/DoWork")

	l.Info("preparing approval email for game subscription ID >%s<", j.Args.GameSubscriptionID)

	subscriptionRec, err := m.GetGameSubscriptionRec(j.Args.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription record >%v<", err)
		return nil, err
	}

	accountRec, err := m.GetAccountRec(subscriptionRec.AccountID, nil)
	if err != nil {
		l.Warn("failed to get account record >%v<", err)
		return nil, err
	}

	gameRec, err := m.GetGameRec(subscriptionRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%v<", err)
		return nil, err
	}

	approvalPath := fmt.Sprintf("/api/v1/game-subscriptions/%s/approve?email=%s", subscriptionRec.ID, url.QueryEscape(accountRec.Email))

	// Construct full URL using configured app host
	approvalURL := fmt.Sprintf("%s%s", w.Config.AppHost, approvalPath)

	// Get account contact name if available
	accountName := ""
	contactRecs, err := m.GetManyAccountContactRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountContactAccountID, Val: accountRec.ID},
		},
		Limit: 1,
		OrderBy: []coresql.OrderBy{
			{Col: account_record.FieldAccountContactCreatedAt, Direction: coresql.OrderDirectionASC},
		},
	})
	if err == nil && len(contactRecs) > 0 {
		accountName = contactRecs[0].Name
	}

	// Render the HTML email template
	baseTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "base.email.html")
	specificTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "game_subscription_approval.email.html")
	tmpl, err := template.ParseFiles(baseTmplPath, specificTmplPath)
	if err != nil {
		l.Warn("failed to parse email template >%v<", err)
		return nil, err
	}

	var body bytes.Buffer
	tmplData := struct {
		AccountName  string
		GameName     string
		ApprovalURL  string
		SupportEmail string
		Year         int
	}{
		AccountName:  accountName,
		GameName:     gameRec.Name,
		ApprovalURL:  approvalURL,
		SupportEmail: "support@playbymail.games",
		Year:         time.Now().Year(),
	}

	if err := tmpl.ExecuteTemplate(&body, "base", tmplData); err != nil {
		l.Warn("failed to render email template >%v<", err)
		return nil, err
	}

	emailMsg := &emailer.Message{
		From:    "noreply@playbymail.games",
		To:      []string{accountRec.Email},
		Subject: fmt.Sprintf("Confirm your subscription to %s", gameRec.Name),
		Body:    body.String(),
	}

	if err := w.emailClient.Send(emailMsg); err != nil {
		l.Warn("failed to send subscription approval email >%v<", err)
		return nil, err
	}

	l.Info("sent subscription approval email to >%s< for game >%s<", accountRec.Email, gameRec.Name)

	return &SendGameSubscriptionApprovalEmailDoWorkResult{RecordCount: 1}, nil
}
