package jobworker

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"bytes"
	"html/template"
	"path/filepath"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

type SendAccountVerificationEmailWorkerArgs struct {
	AccountID string
}

func (SendAccountVerificationEmailWorkerArgs) Kind() string { return "send-account-verification-token" }
func (SendAccountVerificationEmailWorkerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: jobqueue.QueueDefault,
	}
}

type SendAccountVerificationEmailWorker struct {
	river.WorkerDefaults[SendAccountVerificationEmailWorkerArgs]
	emailClient emailer.Emailer
	JobWorker
}

func NewSendAccountVerificationEmailWorker(l logger.Logger, cfg config.Config, s storer.Storer, e emailer.Emailer) (*SendAccountVerificationEmailWorker, error) {
	l = l.WithPackageContext("SendAccountVerificationEmailWorker")

	l.Info("Instantiation SendAccountVerificationEmailWorker")

	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	// We need an email client and we need to know where the templates are stored
	// otherwise we cannot do our work.
	if e == nil {
		return nil, fmt.Errorf("email client is nil")
	}

	if cfg.TemplatesPath == "" {
		return nil, fmt.Errorf("templates path is empty")
	}

	// Check the templates path actually exists
	l.Info("templates path >%s<", cfg.TemplatesPath)

	if _, err := os.Stat(cfg.TemplatesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("templates path does not exist >%s<", cfg.TemplatesPath)
	}

	return &SendAccountVerificationEmailWorker{
		JobWorker:   *jw,
		emailClient: e,
	}, nil
}

func (w *SendAccountVerificationEmailWorker) Work(ctx context.Context, j *river.Job[SendAccountVerificationEmailWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("SendAccountVerificationEmailWorker/Work")

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
		l.Error("SendAccountVerificationEmailWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type SendAccountVerificationEmailDoWorkResult struct {
	RecordCount int
}

func (w *SendAccountVerificationEmailWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[SendAccountVerificationEmailWorkerArgs]) (*SendAccountVerificationEmailDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("SendAccountVerificationEmailWorker/DoWork")

	l.Info("SendAccountVerificationEmailWorker report work record ID >%s<", j.Args.AccountID)

	accountRec, err := m.GetAccountRec(j.Args.AccountID, nil)
	if err != nil {
		l.Warn("failed to get account record >%v<", err)
		return nil, err
	}
	token, err := m.GenerateAccountVerificationToken(accountRec)
	if err != nil {
		l.Warn("failed to generate account verification token >%v<", err)
		return nil, err
	}

	// Render the HTML email template
	tmplPath := filepath.Join(w.Config.TemplatesPath, "account_verification.email.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		l.Warn("failed to parse email template >%v<", err)
		return nil, err
	}

	var body bytes.Buffer
	tmplData := struct {
		VerificationCode string
		SupportEmail     string
	}{
		VerificationCode: token,
		SupportEmail:     "support@playbymail.games",
	}

	if err := tmpl.ExecuteTemplate(&body, "account_verification", tmplData); err != nil {
		l.Warn("failed to render email template >%v<", err)
		return nil, err
	}

	emailMsg := &emailer.Message{
		From:    "noreply@playbymail.games",
		To:      []string{accountRec.Email},
		Subject: "Your PlayByMail verification code",
		Body:    body.String(),
	}
	if err := w.emailClient.Send(emailMsg); err != nil {
		l.Warn("failed to send verification email >%v<", err)
		return nil, err
	}

	l.Info("sent verification email to >%s< with code >%s<", accountRec.Email, token)

	return &SendAccountVerificationEmailDoWorkResult{
		RecordCount: 1,
	}, nil
}
