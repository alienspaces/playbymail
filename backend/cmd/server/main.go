package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/playbymail/core/server"
	runner "gitlab.com/alienspaces/playbymail/internal/runner/server"
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
		os.Exit(0)
	}

	r, err := runner.NewRunner(l, s, j, cfg)
	if err != nil {
		fmt.Printf("(cmd) failed new runner >%v<\n", err)
		os.Exit(0)
	}

	svc, err := server.NewServer(l, s, r)
	if err != nil {
		fmt.Printf("(cmd) failed new server >%v<\n", err)
		os.Exit(0)
	}

	args := make(map[string]any)

	err = svc.Run(args)
	if err != nil {
		fmt.Printf("(cmd) failed server run >%v<\n", err)
		os.Exit(0)
	}

	os.Exit(1)
}
