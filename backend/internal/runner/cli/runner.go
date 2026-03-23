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
				Usage:   "Load or delete a demo game by type (requires --type)",
				Description: `
Loads a demo game into the target database and publishes it.
Use --list to print available game types and names.
Use --replace to remove an existing demo game before loading.
Use --delete to remove an existing demo game without loading a new one.`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "type",
						Aliases: []string{"t"},
						Usage:   "Game type (required; use --list to see options)",
					},
					&cli.BoolFlag{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "Print available demo game types and names, then exit",
					},
					&cli.BoolFlag{
						Name:  "replace",
						Usage: "Remove existing demo game before loading",
					},
					&cli.BoolFlag{
						Name:    "delete",
						Aliases: []string{"d"},
						Usage:   "Delete the demo game for this type (removes all dependents)",
					},
				},
				Action: r.loadGameData,
			},
		// Account / user inspection
		{
			Name:    "db-list-users",
			Aliases: []string{"lu"},
			Usage:   "List all accounts and users with status and session info",
			Action:  r.listUsers,
		},
		// Game instance management
		{
			Name:    "list-game-instances",
			Aliases: []string{"lgi"},
			Usage:   "List all game instances with status, current turn, and player count",
			Action:  r.listGameInstances,
		},
		{
			Name:    "resend-turn-sheet-email",
			Aliases: []string{"rtse"},
			Usage:   "Resend turn sheet notification emails for all players in a game instance",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "game-instance-id",
					Aliases:  []string{"i"},
					Usage:    "ID of the game instance to resend emails for (required)",
					Required: true,
				},
			},
			Action: r.resendTurnSheetEmail,
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
