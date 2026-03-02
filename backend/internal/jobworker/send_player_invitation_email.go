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

// SendPlayerInvitationEmailWorkerArgs defines the job payload for sending player invitation emails.
// GameSubscriptionID is the manager's subscription ID for the game; the join link uses an available
// instance from that manager's subscription. Instance assignment happens when the player accepts.
type SendPlayerInvitationEmailWorkerArgs struct {
	GameSubscriptionID string
	Email              string
}

func (SendPlayerInvitationEmailWorkerArgs) Kind() string {
	return "send-player-invitation-email"
}

func (SendPlayerInvitationEmailWorkerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: jobqueue.QueueDefault}
}

// SendPlayerInvitationEmailWorker sends an invitation email to a prospective player.
// The join link uses any available (non-closed-testing) game instance as the destination.
type SendPlayerInvitationEmailWorker struct {
	river.WorkerDefaults[SendPlayerInvitationEmailWorkerArgs]
	emailClient emailer.Emailer
	JobWorker
}

func NewSendPlayerInvitationEmailWorker(l logger.Logger, cfg config.Config, s storer.Storer, e emailer.Emailer) (*SendPlayerInvitationEmailWorker, error) {
	l = l.WithPackageContext("SendPlayerInvitationEmailWorker")

	l.Info("instantiating SendPlayerInvitationEmailWorker")

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

	if _, err := os.Stat(cfg.TemplatesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("templates path does not exist >%s<", cfg.TemplatesPath)
	}

	return &SendPlayerInvitationEmailWorker{
		JobWorker:   *jw,
		emailClient: e,
	}, nil
}

func (w *SendPlayerInvitationEmailWorker) Work(ctx context.Context, j *river.Job[SendPlayerInvitationEmailWorkerArgs]) error {
	l := w.Log.WithFunctionContext("SendPlayerInvitationEmailWorker/Work")

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
		l.Error("SendPlayerInvitationEmailWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

// SendPlayerInvitationEmailDoWorkResult summarises the work carried out by the worker.
type SendPlayerInvitationEmailDoWorkResult struct {
	RecordCount int
}

func (w *SendPlayerInvitationEmailWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[SendPlayerInvitationEmailWorkerArgs]) (*SendPlayerInvitationEmailDoWorkResult, error) {
	l := w.Log.WithFunctionContext("SendPlayerInvitationEmailWorker/DoWork")

	l.Info("preparing player invitation email for game subscription >%s< email >%s<", j.Args.GameSubscriptionID, j.Args.Email)

	// Validate the subscription exists and has capacity before sending.
	subscriptionRec, err := m.GetGameSubscriptionRec(j.Args.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription record >%v<", err)
		return nil, err
	}

	instanceRec, err := m.FindAvailableGameInstance(j.Args.GameSubscriptionID)
	if err != nil {
		l.Warn("failed to find available game instance for subscription >%s< >%v<", j.Args.GameSubscriptionID, err)
		return nil, err
	}
	if instanceRec == nil {
		return nil, fmt.Errorf("no available game instance found for subscription >%s<", j.Args.GameSubscriptionID)
	}

	gameRec, err := m.GetGameRec(subscriptionRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%v<", err)
		return nil, err
	}

	// Build the join URL using the subscription ID (instance is auto-assigned at join time).
	joinURL := fmt.Sprintf("%s/player/join-game/%s", w.Config.AppHost, j.Args.GameSubscriptionID)

	// Render the HTML email template.
	baseTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "base.email.html")
	specificTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "player_invitation.email.html")
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
		SupportEmail: w.Config.SupportEmailAddress,
		Year:         time.Now().Year(),
	}

	if err := tmpl.ExecuteTemplate(&body, "base", tmplData); err != nil {
		l.Warn("failed to render email template >%v<", err)
		return nil, err
	}

	emailMsg := &emailer.Message{
		From:    w.Config.NoReplyEmailAddress,
		To:      []string{j.Args.Email},
		Subject: fmt.Sprintf("You're invited to play %s", gameRec.Name),
		Body:    body.String(),
	}

	if err := w.emailClient.Send(emailMsg); err != nil {
		l.Warn("failed to send player invitation email >%v<", err)
		return nil, err
	}

	l.Info("sent player invitation email to >%s< for game >%s<", j.Args.Email, gameRec.Name)

	return &SendPlayerInvitationEmailDoWorkResult{RecordCount: 1}, nil
}
