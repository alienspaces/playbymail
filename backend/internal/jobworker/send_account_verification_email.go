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
	"time"

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

	l.Info("instantiating send account verification email worker")

	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	// Job clients and workers will be instantiated under two scenarios.
	// 1. To start a job client server and execute jobs
	// 2. To register jobs with a job client server
	//
	// In the first scenario we need an email client and we need to know where
	// the templates are stored otherwise we cannot do our work.
	//
	// In the second scenario we do not need an email client and we do not need
	// to know where the templates are stored.
	if e == nil {
		l.Warn("email client is nil, assuming instantiation for registration purposes only")
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
	l := w.Log.WithFunctionContext("SendAccountVerificationEmailWorker/Work")

	l.Info("running job ID >%s< args >%#v<", strconv.FormatInt(j.ID, 10), j.Args)

	// We can assume here that we must have an email client as we've been told to perform work.
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
		l.Error("SendAccountVerificationEmailWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type SendAccountVerificationEmailDoWorkResult struct {
	RecordCount int
}

func (w *SendAccountVerificationEmailWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[SendAccountVerificationEmailWorkerArgs]) (*SendAccountVerificationEmailDoWorkResult, error) {
	l := w.Log.WithFunctionContext("SendAccountVerificationEmailWorker/DoWork")

	l.Info("send account verification email worker report work record ID >%s<", j.Args.AccountID)

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
	baseTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "base.email.html")
	specificTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "account_verification.email.html")
	tmpl, err := template.ParseFiles(baseTmplPath, specificTmplPath)
	if err != nil {
		l.Warn("failed to parse email template >%v<", err)
		return nil, err
	}

	var body bytes.Buffer
	tmplData := struct {
		VerificationCode string
		SupportEmail     string
		Year             int
	}{
		VerificationCode: token,
		SupportEmail:     w.Config.SupportEmailAddress,
		Year:             time.Now().Year(),
	}

	if err := tmpl.ExecuteTemplate(&body, "base", tmplData); err != nil {
		l.Warn("failed to render email template >%v<", err)
		return nil, err
	}

	emailMsg := &emailer.Message{
		From:    w.Config.NoReplyEmailAddress,
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
