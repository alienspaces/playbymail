package cli

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/runnable"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

// CLI -
type CLI struct {
	Log    logger.Logger
	Store  storer.Storer
	Runner runnable.Runnable
}

// NewCLI -
func NewCLI(l logger.Logger, s storer.Storer, r runnable.Runnable) (*CLI, error) {

	cli := CLI{
		Log:    l,
		Store:  s,
		Runner: r,
	}

	err := cli.Init()
	if err != nil {
		return nil, err
	}

	return &cli, nil
}

// Init -
func (cli *CLI) Init() error {

	// TODO: alerting, retries
	return cli.Runner.Init(cli.Store)
}

// Run -
func (cli *CLI) Run(args map[string]interface{}) error {

	// TODO:
	// - alerting on errors
	// - retries on start up
	// - reload  on config changes
	return cli.Runner.Run(args)
}
