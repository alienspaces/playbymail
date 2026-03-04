package runner

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// loadGameData loads game data for the selected scenario into the target database.
// Scenario is required (use --list-scenarios to see options). Loaded games are draft unless --publish is set.
func (rnr *Runner) loadGameData(c *cli.Context) error {
	l := loggerWithFunctionContext(rnr.Log, "loadGameData")

	if c.Bool("list-scenarios") {
		return rnr.listGameDataScenarios(c)
	}

	scenario := c.String("scenario")
	if scenario == "" {
		return fmt.Errorf("--scenario is required (use --list-scenarios to see available scenarios)")
	}

	entry, ok := GameDataScenarios[scenario]
	if !ok {
		l.Warn("unknown scenario >%s<", scenario)
		return fmt.Errorf("unknown scenario %q (use --list-scenarios to see available scenarios)", scenario)
	}

	l.Info("** Load Game Data (scenario: %s) **", scenario)
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

	if c.Bool("publish") {
		mm, ok := rnr.Domain.(*domain.Domain)
		if !ok {
			l.Warn("domain is not *domain.Domain, cannot publish games")
			return fmt.Errorf("cannot publish: domain type assertion failed")
		}
		for _, rec := range testHarness.Data.GameRecs {
			rec.Status = game_record.GameStatusPublished
			_, err = mm.UpdateGameRec(rec)
			if err != nil {
				l.Warn("failed publishing game %s >%v<", rec.ID, err)
				return err
			}
			l.Info("published game %s", rec.ID)
		}
	}

	l.Info("game data loaded successfully")
	return nil
}

// listGameDataScenarios prints registered scenario names and descriptions to stdout.
func (rnr *Runner) listGameDataScenarios(_ *cli.Context) error {
	names := make([]string, 0, len(GameDataScenarios))
	for name := range GameDataScenarios {
		names = append(names, name)
	}
	sort.Strings(names)
	fmt.Println("Available game data scenarios:")
	fmt.Println()
	for _, name := range names {
		entry := GameDataScenarios[name]
		fmt.Printf("  %s\n    %s\n\n", name, entry.Description)
	}
	return nil
}
