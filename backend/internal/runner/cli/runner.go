package runner

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/urfave/cli/v2"

	corecli "gitlab.com/alienspaces/playbymail/core/cli"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// Runner -
type Runner struct {
	corecli.Runner
	Config  config.Config
	Scanner turnsheet.TurnSheetScanner
}

const (
	applicationName = "cli"
)

func NewRunner(cfg config.Config, l logger.Logger, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner) (*Runner, error) {
	l = l.WithApplicationContext(applicationName)

	cr, err := corecli.NewRunnerWithConfig(l, j, cfg.Config)
	if err != nil {
		err := fmt.Errorf("failed core runner >%v<", err)
		l.Warn(err.Error())
		return nil, err
	}

	r := Runner{
		Runner:  *cr,
		Config:  cfg,
		Scanner: scanner,
	}

	r.DeferDomainInitialisation = true
	r.DomainFunc = r.domainFunc

	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	r.App = &cli.App{
		Commands: []*cli.Command{
			// Seed Data Operations
			{
				Name:    "db-load-seed-data",
				Aliases: []string{"lsd"},
				Usage:   "Load seed data",
				Description: `
Loads seed data. Typically used when deploying QA environments or for local manual testing.`,
				Action: r.loadSeedData,
			},
			{
				Name:    "db-load-seed-reference-data",
				Aliases: []string{"lsrd"},
				Usage:   "Load seed reference data",
				Description: `
Loads static reference data that is expected to exist on any environment.`,
				Action: r.loadSeedReferenceData,
			},
		},
	}

	return &r, nil
}

func (rnr *Runner) domainFunc() (domainer.Domainer, error) {
	m, err := domain.NewDomain(rnr.Log, rnr.Config)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// loggerWithFunctionContext - Returns a logger with package context and provided function context
func loggerWithFunctionContext(l logger.Logger, functionName string) logger.Logger {
	return logging.LoggerWithFunctionContext(l, "runner", functionName)
}
