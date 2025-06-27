package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/runnable"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

type Runner struct {
	Log    logger.Logger
	Store  storer.Storer
	Domain domainer.Domainer

	// The JobClient within CLI context should only be used for registering jobs
	// and should have no registered workers. As such the CLI runnner does not
	// start and stop the the job client workers. Jobs would typically be registered
	// within the model layer however there may be scenarios within CLI context
	// where jobs may need to be created in a non-typical scenario.
	JobClient *river.Client[pgx.Tx]

	// cli configuration - https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	App *cli.App

	// DomainFunc returns the service specific domainer implementation
	DomainFunc func() (domainer.Domainer, error)

	// Initialisation will be deferred, it becomes the responsiblity
	// of the runner implementation to call init
	DeferDomainInitialisation bool

	// General configuration
	config config.Config
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

func NewRunner(l logger.Logger, j *river.Client[pgx.Tx], cfg config.Config) (*Runner, error) {

	r := Runner{
		Log:       l,
		JobClient: j,
		config:    cfg,
	}

	return &r, nil
}

// Init - override to perform custom initialization
func (rnr *Runner) Init(s storer.Storer) error {

	rnr.Log.Debug("init")

	// Storer
	rnr.Store = s

	// Deferring model intialisation provides the ability for the
	// CLI runner to handle connections and model initialisation
	// itself or to forfeit the usage of database storage and a
	// model altogether.
	if !rnr.DeferDomainInitialisation && rnr.Store == nil {
		msg := "store undefined, cannot init runner"
		rnr.Log.Warn(msg)
		return errors.New(msg)
	}

	return nil
}

// Run - Runs the CLI application.
func (rnr *Runner) Run(args map[string]any) (err error) {

	// Deferring model intialisation provides the ability for the
	// CLI runner to handle connections and model initialisation
	// itself or to forfeit the usage of database storage and a
	// model altogether.
	if !rnr.DeferDomainInitialisation {
		err := rnr.InitDomain()
		if err != nil {
			rnr.Log.Warn("failed model init >%v<", err)
			return err
		}
	}

	// Run
	err = rnr.App.Run(os.Args)
	if err != nil {
		rnr.Log.Warn("failed running app >%v<", err)

		// Rollback database transaction on error
		if rnr.Domain != nil {
			rnr.Log.Warn("rolling back database transaction")
			rnr.Domain.Rollback()
		}

		return err
	}

	// Commit database transaction
	if rnr.Domain != nil {
		rnr.Log.Debug("committing database transaction")
		err = rnr.Domain.Commit()
		if err != nil {
			rnr.Log.Warn("failed model commit >%v<", err)
			return err
		}
	}

	return nil
}

// InitDomain iniitialises a database connection, sources a new model and initialises the model
// with a new database transaction.
func (rnr *Runner) InitDomain() error {

	if rnr.DomainFunc == nil {
		rnr.DomainFunc = rnr.defaultDomainFunc
	}

	m, err := rnr.DomainFunc()
	if err != nil {
		return err
	}

	if m == nil {
		err := fmt.Errorf("model is nil, cannot continue")
		return err
	}

	rnr.Domain = m

	err = rnr.InitDomainTx()
	if err != nil {
		rnr.Log.Warn("failed model init tx >%v<", err)
		return err
	}

	return nil
}

// InitDomainTx initialises the model with a new database transaction.
func (rnr *Runner) InitDomainTx() error {

	rnr.Log.Debug("initialising model tx")

	tx, err := rnr.Store.BeginTx()
	if err != nil {
		rnr.Log.Warn("failed getting store transaction >%v<", err)
		return err
	}

	err = rnr.Domain.Init(tx)
	if err != nil {
		rnr.Log.Warn("failed model init >%v<", err)
		return err
	}

	return nil
}

// defaultDomainFunc does not provide a domainer, set the property DomainFunc to
// provide your own custom domainer.
func (rnr *Runner) defaultDomainFunc() (domainer.Domainer, error) {

	rnr.Log.Debug("defaultDomainFunc")

	return nil, nil
}
