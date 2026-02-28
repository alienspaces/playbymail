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

// SendTurnSheetNotificationEmailWorkerArgs defines the job payload for sending turn sheet notification emails
type SendTurnSheetNotificationEmailWorkerArgs struct {
	GameSubscriptionInstanceID string
	TurnNumber                 int
}

func (SendTurnSheetNotificationEmailWorkerArgs) Kind() string {
	return "send-turn-sheet-notification-email"
}

func (SendTurnSheetNotificationEmailWorkerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: jobqueue.QueueDefault}
}

// SendTurnSheetNotificationEmailWorker sends an email containing a secure link to the turn sheet viewer
type SendTurnSheetNotificationEmailWorker struct {
	river.WorkerDefaults[SendTurnSheetNotificationEmailWorkerArgs]
	emailClient emailer.Emailer
	JobWorker
}

func NewSendTurnSheetNotificationEmailWorker(l logger.Logger, cfg config.Config, s storer.Storer, e emailer.Emailer) (*SendTurnSheetNotificationEmailWorker, error) {
	l = l.WithPackageContext("SendTurnSheetNotificationEmailWorker")

	l.Info("instantiating SendTurnSheetNotificationEmailWorker")

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

	return &SendTurnSheetNotificationEmailWorker{
		JobWorker:   *jw,
		emailClient: e,
	}, nil
}

func (w *SendTurnSheetNotificationEmailWorker) Work(ctx context.Context, j *river.Job[SendTurnSheetNotificationEmailWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("SendTurnSheetNotificationEmailWorker/Work")

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
		l.Error("SendTurnSheetNotificationEmailWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

// SendTurnSheetNotificationEmailDoWorkResult summarises the work carried out by the worker
type SendTurnSheetNotificationEmailDoWorkResult struct {
	RecordCount int
}

func (w *SendTurnSheetNotificationEmailWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[SendTurnSheetNotificationEmailWorkerArgs]) (*SendTurnSheetNotificationEmailDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("SendTurnSheetNotificationEmailWorker/DoWork")

	l.Info("preparing turn sheet notification email for instance ID >%s< turn >%d<", j.Args.GameSubscriptionInstanceID, j.Args.TurnNumber)

	// Get the game subscription instance
	instanceRec, err := m.GetGameSubscriptionInstanceRec(j.Args.GameSubscriptionInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game subscription instance record >%v<", err)
		return nil, err
	}

	// Get the game subscription
	subscriptionRec, err := m.GetGameSubscriptionRec(instanceRec.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription record >%v<", err)
		return nil, err
	}

	// Get the game instance to check delivery method
	gameInstanceRec, err := m.GetGameInstanceRec(instanceRec.GameInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance ID >%s< >%v<", instanceRec.GameInstanceID, err)
		return nil, err
	}

	// Check if email delivery is enabled for this game instance
	if !gameInstanceRec.DeliveryEmail {
		l.Info("email delivery not enabled for game instance >%s<, skipping email notification", gameInstanceRec.ID)
		return &SendTurnSheetNotificationEmailDoWorkResult{RecordCount: 0}, nil
	}

	// Get the account to get the email address
	accountRec, err := m.GetAccountRec(instanceRec.AccountID, nil)
	if err != nil {
		l.Warn("failed to get account record >%v<", err)
		return nil, err
	}

	// Get the game record for the game name
	gameRec, err := m.GetGameRec(subscriptionRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%v<", err)
		return nil, err
	}

	// Generate or get the turn sheet key
	turnSheetToken, err := m.GenerateGameSubscriptionInstanceTurnSheetToken(j.Args.GameSubscriptionInstanceID)
	if err != nil {
		l.Warn("failed to generate game subscription instance turn sheet token >%v<", err)
		return nil, err
	}

	// Get the instance again to get the expiration time
	instanceRec, err = m.GetGameSubscriptionInstanceRec(j.Args.GameSubscriptionInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game subscription instance record after token generation >%v<", err)
		return nil, err
	}

	// Build turn sheet viewer login URL
	turnSheetPath := fmt.Sprintf("/player/game-subscription-instances/%s/login/%s", j.Args.GameSubscriptionInstanceID, turnSheetToken)
	turnSheetURL := fmt.Sprintf("%s%s", w.Config.AppHost, turnSheetPath)

	// Format expiration date/time
	var expirationDate, expirationTime string
	if instanceRec.TurnSheetTokenExpiresAt.Valid {
		expirationTimeVal := instanceRec.TurnSheetTokenExpiresAt.Time
		expirationDate = expirationTimeVal.Format("January 2, 2006")
		expirationTime = expirationTimeVal.Format("3:04 PM MST")
	} else {
		// Fallback: 3 days from now
		expirationTimeVal := time.Now().Add(3 * 24 * time.Hour)
		expirationDate = expirationTimeVal.Format("January 2, 2006")
		expirationTime = expirationTimeVal.Format("3:04 PM MST")
	}

	// Render the HTML email template
	baseTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "base.email.html")
	specificTmplPath := filepath.Join(w.Config.TemplatesPath, "email", "turn_sheet_notification.email.html")
	tmpl, err := template.ParseFiles(baseTmplPath, specificTmplPath)
	if err != nil {
		l.Warn("failed to parse email template >%v<", err)
		return nil, err
	}

	var body bytes.Buffer
	tmplData := struct {
		GameName       string
		TurnNumber     int
		TurnSheetURL   string
		ExpirationDate string
		ExpirationTime string
		SupportEmail   string
		Year           int
	}{
		GameName:       gameRec.Name,
		TurnNumber:     j.Args.TurnNumber,
		TurnSheetURL:   turnSheetURL,
		ExpirationDate: expirationDate,
		ExpirationTime: expirationTime,
		SupportEmail:   w.Config.SupportEmailAddress,
		Year:           time.Now().Year(),
	}

	if err := tmpl.ExecuteTemplate(&body, "base", tmplData); err != nil {
		l.Warn("failed to render email template >%v<", err)
		return nil, err
	}

	emailMsg := &emailer.Message{
		From:    w.Config.NoReplyEmailAddress,
		To:      []string{accountRec.Email},
		Subject: fmt.Sprintf("Turn %d is ready for %s", j.Args.TurnNumber, gameRec.Name),
		Body:    body.String(),
	}

	if err := w.emailClient.Send(emailMsg); err != nil {
		l.Warn("failed to send turn sheet notification email >%v<", err)
		return nil, err
	}

	l.Info("sent turn sheet notification email to >%s< for game >%s< turn >%d<", accountRec.Email, gameRec.Name, j.Args.TurnNumber)

	return &SendTurnSheetNotificationEmailDoWorkResult{RecordCount: 1}, nil
}
