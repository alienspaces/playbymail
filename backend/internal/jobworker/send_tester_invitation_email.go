package jobworker

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// SendTesterInvitationEmailWorkerArgs defines the job payload for sending tester invitation emails
type SendTesterInvitationEmailWorkerArgs struct {
	GameInstanceID string
	Email          string
}

func (SendTesterInvitationEmailWorkerArgs) Kind() string {
	return "send-tester-invitation-email"
}

func (SendTesterInvitationEmailWorkerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: jobqueue.QueueDefault}
}

// SendTesterInvitationEmailWorker sends an email containing a join game link for closed testing
type SendTesterInvitationEmailWorker struct {
	river.WorkerDefaults[SendTesterInvitationEmailWorkerArgs]
	emailClient emailer.Emailer
	JobWorker
}

func NewSendTesterInvitationEmailWorker(l logger.Logger, cfg config.Config, s storer.Storer, e emailer.Emailer) (*SendTesterInvitationEmailWorker, error) {
	l = l.WithPackageContext("SendTesterInvitationEmailWorker")

	l.Info("instantiating SendTesterInvitationEmailWorker")

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

	return &SendTesterInvitationEmailWorker{
		JobWorker:   *jw,
		emailClient: e,
	}, nil
}

func (w *SendTesterInvitationEmailWorker) Work(ctx context.Context, j *river.Job[SendTesterInvitationEmailWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("SendTesterInvitationEmailWorker/Work")

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
		l.Error("SendTesterInvitationEmailWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

// SendTesterInvitationEmailDoWorkResult summarises the work carried out by the worker
type SendTesterInvitationEmailDoWorkResult struct {
	RecordCount int
}

func (w *SendTesterInvitationEmailWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[SendTesterInvitationEmailWorkerArgs]) (*SendTesterInvitationEmailDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("SendTesterInvitationEmailWorker/DoWork")

	l.Info("preparing tester invitation email for game instance ID >%s< email >%s<", j.Args.GameInstanceID, j.Args.Email)

	// Get the game instance
	instanceRec, err := m.GetGameInstanceRec(j.Args.GameInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance record >%v<", err)
		return nil, err
	}

	// Validate it's in closed testing mode
	if !instanceRec.IsClosedTesting {
		return nil, fmt.Errorf("game instance is not in closed testing mode")
	}

	// Get join game key (should already exist, but generate if not)
	if !instanceRec.JoinGameKey.Valid || instanceRec.JoinGameKey.String == "" {
		_, err = m.GenerateJoinGameKey(j.Args.GameInstanceID)
		if err != nil {
			l.Warn("failed to generate join game key >%v<", err)
			return nil, err
		}
		// Re-fetch to get the new key
		instanceRec, err = m.GetGameInstanceRec(j.Args.GameInstanceID, nil)
		if err != nil {
			l.Warn("failed to get game instance after key generation >%v<", err)
			return nil, err
		}
	}

	// Get game record for name
	gameRec, err := m.GetGameRec(instanceRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%v<", err)
		return nil, err
	}

	// Build join game URL
	joinPath := fmt.Sprintf("/player/join-game/%s", instanceRec.JoinGameKey.String)
	joinURL := fmt.Sprintf("%s%s", w.Config.AppHost, joinPath)

	// Render the HTML email template
	baseTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "base.email.html")
	specificTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "tester_invitation.email.html")
	tmpl, err := template.ParseFiles(baseTmplPath, specificTmplPath)
	if err != nil {
		l.Warn("failed to parse email template >%v<", err)
		return nil, err
	}

	var body bytes.Buffer
	tmplData := struct {
		GameName     string
		JoinURL      string
		SupportEmail string
		Year         int
	}{
		GameName:     gameRec.Name,
		JoinURL:      joinURL,
		SupportEmail: "support@playbymail.games",
		Year:         time.Now().Year(),
	}

	if err := tmpl.ExecuteTemplate(&body, "base", tmplData); err != nil {
		l.Warn("failed to render email template >%v<", err)
		return nil, err
	}

	emailMsg := &emailer.Message{
		From:    "noreply@playbymail.games",
		To:      []string{j.Args.Email},
		Subject: fmt.Sprintf("You're invited to test %s", gameRec.Name),
		Body:    body.String(),
	}

	if err := w.emailClient.Send(emailMsg); err != nil {
		l.Warn("failed to send tester invitation email >%v<", err)
		return nil, err
	}

	l.Info("sent tester invitation email to >%s< for game >%s<", j.Args.Email, gameRec.Name)

	return &SendTesterInvitationEmailDoWorkResult{RecordCount: 1}, nil
}

