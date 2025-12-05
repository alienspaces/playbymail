package jobclient

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobclient "gitlab.com/alienspaces/playbymail/core/jobclient"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// NewJobClient creates a new job client. When no queue names are specified it provides the
// ability to queue jobs only. When one or more queue names are specified it will also process
// jobs for those queues.
func NewJobClient(l logger.Logger, cfg config.Config, s storer.Storer, e emailer.Emailer, queueNames []string) (*river.Client[pgx.Tx], error) {

	var err error

	riverConfig, err := getRiverConfig(l, cfg, s, e, queueNames)
	if err != nil {
		return nil, err
	}

	riverClient, err := corejobclient.NewJobClient(s, riverConfig)
	if err != nil {
		return nil, err
	}

	return riverClient, nil
}

func getRiverConfig(l logger.Logger, cfg config.Config, s storer.Storer, e emailer.Emailer, queueNames []string) (*river.Config, error) {
	l = l.WithFunctionContext("getRiverConfig")

	riverConfig := river.Config{}

	// Add all job workers regardless of queues this client is going to process as river will
	// use the registered job workers to validate registered jobs have an associated worker.
	// This means that every deployed server requires all configuration required for all workers.
	w, err := getWorkers(l, cfg, s, e)
	if err != nil {
		return nil, err
	}
	riverConfig.Workers = w

	if len(queueNames) == 0 {
		return &riverConfig, nil
	}

	// Periodic job configuration functions
	periodicJobFuncs := map[string]func(l logger.Logger, s storer.Storer, p []*river.PeriodicJob) ([]*river.PeriodicJob, error){
		jobqueue.QueueDefault: addDefaultPeriodicJobs,
		jobqueue.QueueGame:    addGamePeriodicJobs,
	}

	// Add queue and periodic job configuration for the queues this client is going to process
	q := make(map[string]river.QueueConfig)
	p := []*river.PeriodicJob{}

	for _, queueName := range queueNames {

		l.Info("adding queue >%s<", queueName)

		// Add queue configuration
		if !jobqueue.Queues.Has(queueName) {
			return nil, fmt.Errorf("failed creating job client configuration, queue name >%s< is not a recognised queue name", queueName)
		}
		q[queueName] = river.QueueConfig{
			MaxWorkers: 1,
		}

		// Add queue periodic jobs
		p, err = periodicJobFuncs[queueName](l, s, p)
		if err != nil {
			return nil, err
		}
	}

	l.Info("added %d periodic jobs", len(p))

	riverConfig.Queues = q
	riverConfig.PeriodicJobs = p

	return &riverConfig, nil
}

func addDefaultPeriodicJobs(l logger.Logger, s storer.Storer, p []*river.PeriodicJob) ([]*river.PeriodicJob, error) {
	l = l.WithFunctionContext("addDefaultPeriodicJobs")

	l.Info("adding default periodic jobs")

	return p, nil
}

func addGamePeriodicJobs(l logger.Logger, s storer.Storer, p []*river.PeriodicJob) ([]*river.PeriodicJob, error) {
	l = l.WithFunctionContext("addGamePeriodicJobs")

	l.Info("adding game periodic jobs")

	// TODO: Add periodic job to check game deadlines every hour
	// Need to check river documentation for correct PeriodicJob structure

	return p, nil
}

func getWorkers(l logger.Logger, cfg config.Config, s storer.Storer, e emailer.Emailer) (*river.Workers, error) {
	w := river.NewWorkers()

	// Add account verification email worker
	// Sends verification emails to users when they sign up or request email verification.
	// Handles email delivery and tracks verification status.
	sendAccountVerificationEmailWorker, err := jobworker.NewSendAccountVerificationEmailWorker(l, cfg, s, e)
	if err != nil {
		return nil, fmt.Errorf("failed NewSendAccountVerificationEmailWorker worker: %w", err)
	}

	if err := river.AddWorkerSafely(w, sendAccountVerificationEmailWorker); err != nil {
		return nil, fmt.Errorf("failed to add NewSendAccountVerificationEmailWorker worker: %w", err)
	}

	// Add game subscription approval email worker
	// Sends confirmation emails when a player submits a join game turn sheet so they can approve their subscription.
	sendGameSubscriptionApprovalEmailWorker, err := jobworker.NewSendGameSubscriptionApprovalEmailWorker(l, cfg, s, e)
	if err != nil {
		return nil, fmt.Errorf("failed NewSendGameSubscriptionApprovalEmailWorker worker: %w", err)
	}

	if err := river.AddWorkerSafely(w, sendGameSubscriptionApprovalEmailWorker); err != nil {
		return nil, fmt.Errorf("failed to add NewSendGameSubscriptionApprovalEmailWorker worker: %w", err)
	}

	// Add tester invitation email worker
	// Sends invitation emails to testers for closed testing game instances with join game links.
	sendTesterInvitationEmailWorker, err := jobworker.NewSendTesterInvitationEmailWorker(l, cfg, s, e)
	if err != nil {
		return nil, fmt.Errorf("failed NewSendTesterInvitationEmailWorker worker: %w", err)
	}

	if err := river.AddWorkerSafely(w, sendTesterInvitationEmailWorker); err != nil {
		return nil, fmt.Errorf("failed to add NewSendTesterInvitationEmailWorker worker: %w", err)
	}

	// Add game turn processing execution worker
	// Executes game turn processing by running game logic, updating game state, and handling player actions.
	// Manages turn progression, game rules enforcement, and state transitions.
	GameTurnProcessingWorker, err := jobworker.NewGameTurnProcessingWorker(l, cfg, s)
	if err != nil {
		return nil, fmt.Errorf("failed NewGameTurnProcessingWorker worker: %w", err)
	}

	if err := river.AddWorkerSafely(w, GameTurnProcessingWorker); err != nil {
		return nil, fmt.Errorf("failed to add NewGameTurnProcessingWorker worker: %w", err)
	}

	// Add game turn processing queue worker
	// Queues turn processing jobs for games that need them based on deadlines and timing.
	// If processing is needed, instantiates a process game turn worker job.
	GameTurnQueueingWorker, err := jobworker.NewGameTurnQueueingWorker(l, cfg, s)
	if err != nil {
		return nil, fmt.Errorf("failed NewGameTurnQueueingWorker worker: %w", err)
	}

	if err := river.AddWorkerSafely(w, GameTurnQueueingWorker); err != nil {
		return nil, fmt.Errorf("failed to add NewGameTurnQueueingWorker worker: %w", err)
	}

	// Add join game turn sheet worker
	// Processes join game turn sheets when a game subscription is approved,
	// creating the necessary game entities (game instance, character, character instance, etc.)
	joinGameTurnSheetWorker, err := jobworker.NewProcessSubscriptionWorker(l, cfg, s)
	if err != nil {
		return nil, fmt.Errorf("failed NewProcessSubscriptionWorker worker: %w", err)
	}

	if err := river.AddWorkerSafely(w, joinGameTurnSheetWorker); err != nil {
		return nil, fmt.Errorf("failed to add NewProcessSubscriptionWorker worker: %w", err)
	}

	return w, nil
}
