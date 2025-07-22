package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/playbymail/core/cli"
	runner "gitlab.com/alienspaces/playbymail/internal/runner/cli"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func main() {

	// runner config
	cfg, err := config.Parse()
	if err != nil {
		fmt.Printf("(cmd) failed parse config >%v<\n", err)
		os.Exit(1)
	}

	l, s, j, err := deps.NewDefaultDependencies(cfg)
	if err != nil {
		fmt.Printf("(cmd) failed default dependencies >%v<\n", err)
		os.Exit(1)
	}

	r, err := runner.NewRunner(l, j, cfg)
	if err != nil {
		fmt.Printf("(cmd) failed new runner >%v<\n", err)
		os.Exit(1)
	}

	app, err := cli.NewCLI(l, s, r)
	if err != nil {
		fmt.Printf("(cmd) failed new CLI >%v<\n", err)
		os.Exit(1)
	}

	args := make(map[string]any)

	err = app.Run(args)
	if err != nil {
		fmt.Printf("(cmd) failed CLI run >%v<\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
