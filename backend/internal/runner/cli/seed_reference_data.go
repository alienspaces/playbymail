package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/playbymail/internal/harness"
)

// loadSeedReferenceData loads the supported set of seed reference data for
// CI and QA test environments
func (rnr *Runner) loadSeedReferenceData(c *cli.Context) error {
	l := loggerWithFunctionContext(rnr.Log, "loadSeedReferenceData")

	l.Info("** Load Seed Reference Data **")

	// harness
	config := rnr.SeedReferenceDataConfig()

	err := rnr.InitDomain()
	if err != nil {
		l.Warn("failed domain init >%v<", err)
		return err
	}

	testHarness, err := harness.NewTesting(rnr.Log, rnr.Store, rnr.JobClient, rnr.Config, config)
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

func (rnr *Runner) SeedReferenceDataConfig() harness.DataConfig {
	return harness.DataConfig{}
}
