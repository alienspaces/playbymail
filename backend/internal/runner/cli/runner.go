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
			// Test data operations (E2E / Playwright)
			{
				Name:    "db-load-test-data",
				Aliases: []string{"ltd"},
				Usage:   "Load E2E test data (accounts + games for Playwright)",
				Description: `
Loads E2E test data: accounts and games used by Playwright tests.
Typically used when setting up QA or local environments for E2E.`,
				Action: r.loadTestData,
			},
			{
				Name:    "db-load-test-reference-data",
				Aliases: []string{"ltrd"},
				Usage:   "Load test reference data",
				Description: `
Loads static reference data expected to exist on any test environment.`,
				Action: r.loadTestReferenceData,
			},
			// Demo scenario operations (game data for players to try)
			{
				Name:    "db-load-demo-game",
				Aliases: []string{"lddg"},
				Usage:   "Load a demo game by name (requires --game)",
				Description: `
Loads a demo game into the target database.
Use --list to print available demo games. Games are loaded as draft unless --publish is set.`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "game",
						Aliases: []string{"g"},
						Usage:   "Game name (required; use --list to see options)",
					},
					&cli.BoolFlag{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "Print available demo games, then exit",
					},
					&cli.BoolFlag{
						Name:  "publish",
						Usage: "Publish loaded games (default: games remain draft)",
					},
					&cli.BoolFlag{
						Name:  "replace",
						Usage: "Remove existing game with same name before loading",
					},
				},
				Action: r.loadGameData,
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
