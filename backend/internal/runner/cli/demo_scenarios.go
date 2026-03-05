package runner

import (
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/cli/demo_scenarios"
)

// DemoScenarioEntry holds a demo scenario: a config factory and a short description.
type DemoScenarioEntry struct {
	Config      func() harness.DataConfig
	Description string
}

// DemoScenarios is the registry of scenario name -> config and description.
// Used by db-load-game-data to resolve --scenario and by --list-scenarios.
var DemoScenarios = map[string]DemoScenarioEntry{
	"full": {
		Config:      demo_scenarios.FullAdventureConfig,
		Description: "Full adventure game demo: locations, links, link requirements, items, creatures, characters, instances, accounts. Games loaded as draft. Use --publish to publish.",
	},
}
