package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/cli/seed_data"
)

// loadSeedData seeds the database with seed data
func (rnr *Runner) loadSeedData(c *cli.Context) error {
	l := loggerWithFunctionContext(rnr.Log, "loadSeedData")

	l.Info("** Load Seed Data **")

	// harness
	config := seed_data.SeedDataConfig()

	err := rnr.InitDomain()
	if err != nil {
		l.Warn("failed domain init >%v<", err)
		return err
	}

	testHarness, err := harness.NewTesting(rnr.Config, rnr.Log, rnr.Store, rnr.JobClient, rnr.Scanner, config)
	if err != nil {
		l.Warn("failed new testing harness >%v<", err)
		return err
	}

	// We want to commit data so that it is available for other commands
	// that need to use the data.
	testHarness.ShouldCommitData = true

	_, err = testHarness.Setup()
	if err != nil {
		l.Warn("failed harness setup >%v<", err)
		return err
	}

	return nil
}
