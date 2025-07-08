package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/cli/maintestdata"
)

// loadMainTestData seeds the database with test data
func (rnr *Runner) loadMainTestData(c *cli.Context) error {
	l := loggerWithFunctionContext(rnr.Log, "loadMainTestData")

	l.Info("** Load Main Test Data **")

	// harness
	config := maintestdata.MainTestDataConfig()

	err := rnr.InitDomain()
	if err != nil {
		l.Warn("failed domain init >%v<", err)
		return err
	}

	testHarness, err := harness.NewTesting(rnr.Log, rnr.Store, rnr.JobClient, config)
	if err != nil {
		l.Warn("failed new testing harness >%v<", err)
		return err
	}

	_, err = testHarness.Setup()
	if err != nil {
		l.Warn("failed harness setup >%v<", err)
		return err
	}

	return nil
}
