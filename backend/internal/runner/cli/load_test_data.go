package runner

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/playbymail/internal/harness"
)

// loadTestData loads test game data for the selected scenario into the target database.
func (rnr *Runner) loadTestData(c *cli.Context) error {
	l := loggerWithFunctionContext(rnr.Log, "loadTestData")

	if c.Bool("list-scenarios") {
		return rnr.listTestDataScenarios(c)
	}

	scenario := c.String("scenario")
	entry, ok := TestDataScenarios[scenario]
	if !ok {
		l.Warn("unknown scenario >%s<", scenario)
		return fmt.Errorf("unknown scenario %q (use --list-scenarios to see available scenarios)", scenario)
	}

	l.Info("** Load Test Data (scenario: %s) **", scenario)
	config := entry.Config()

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

	testHarness.ShouldCommitData = true

	_, err = testHarness.Setup()
	if err != nil {
		l.Warn("failed harness setup >%v<", err)
		return err
	}

	l.Info("test data loaded successfully")
	return nil
}

// listTestDataScenarios prints registered scenario names and descriptions to stdout.
func (rnr *Runner) listTestDataScenarios(_ *cli.Context) error {
	names := make([]string, 0, len(TestDataScenarios))
	for name := range TestDataScenarios {
		names = append(names, name)
	}
	sort.Strings(names)
	fmt.Println("Available test data scenarios:")
	fmt.Println()
	for _, name := range names {
		entry := TestDataScenarios[name]
		fmt.Printf("  %s\n    %s\n\n", name, entry.Description)
	}
	return nil
}
