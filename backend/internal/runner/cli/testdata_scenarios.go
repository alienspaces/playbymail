package runner

import (
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/cli/seed_data"
)

// ScenarioEntry holds a test data scenario: a config factory and a short description.
type ScenarioEntry struct {
	Config      func() harness.DataConfig
	Description string
}

// TestDataScenarios is the registry of scenario name -> config and description.
// Used by db-load-test-data to resolve --scenario and by --list-scenarios.
var TestDataScenarios = map[string]ScenarioEntry{
	"seed": {
		Config:      seed_data.SeedDataConfig,
		Description: "Richer seed data: two games, multiple accounts, locations, links. Good for QA and local manual testing.",
	},
	"default": {
		Config:      harness.DefaultDataConfig,
		Description: "Minimal adventure game data from harness default config. Good for unit-test style setups.",
	},
}
