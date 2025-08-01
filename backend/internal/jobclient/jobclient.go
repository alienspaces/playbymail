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

		l.Info("Adding queue >%s<", queueName)

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

	l.Info("Added %d periodic jobs", len(p))

	riverConfig.Queues = q
	riverConfig.PeriodicJobs = p

	return &riverConfig, nil
}

func addDefaultPeriodicJobs(l logger.Logger, s storer.Storer, p []*river.PeriodicJob) ([]*river.PeriodicJob, error) {
	l = l.WithFunctionContext("addDefaultPeriodicJobs")

	l.Info("Adding default periodic jobs")

	return p, nil
}

func addGamePeriodicJobs(l logger.Logger, s storer.Storer, p []*river.PeriodicJob) ([]*river.PeriodicJob, error) {
	l = l.WithFunctionContext("addGamePeriodicJobs")

	l.Info("Adding game periodic jobs")

	// TODO: Add periodic job to check game deadlines every hour
	// Need to check river documentation for correct PeriodicJob structure

	return p, nil
}

func getWorkers(l logger.Logger, cfg config.Config, s storer.Storer, e emailer.Emailer) (*river.Workers, error) {
	w := river.NewWorkers()

	SendAccountVerificationEmailWorker, err := jobworker.NewSendAccountVerificationEmailWorker(l, cfg, s, e)
	if err != nil {
		return nil, fmt.Errorf("failed NewSendAccountVerificationEmailWorker worker: %w", err)
	}

	if err := river.AddWorkerSafely(w, SendAccountVerificationEmailWorker); err != nil {
		return nil, fmt.Errorf("failed to add NewSendAccountVerificationEmailWorker worker: %w", err)
	}

	// Add game processing workers
	ProcessGameTurnWorker, err := jobworker.NewProcessGameTurnWorker(l, cfg, s)
	if err != nil {
		return nil, fmt.Errorf("failed NewProcessGameTurnWorker worker: %w", err)
	}

	if err := river.AddWorkerSafely(w, ProcessGameTurnWorker); err != nil {
		return nil, fmt.Errorf("failed to add NewProcessGameTurnWorker worker: %w", err)
	}

	CheckGameDeadlinesWorker, err := jobworker.NewCheckGameDeadlinesWorker(l, cfg, s)
	if err != nil {
		return nil, fmt.Errorf("failed NewCheckGameDeadlinesWorker worker: %w", err)
	}

	if err := river.AddWorkerSafely(w, CheckGameDeadlinesWorker); err != nil {
		return nil, fmt.Errorf("failed to add NewCheckGameDeadlinesWorker worker: %w", err)
	}

	return w, nil
}
