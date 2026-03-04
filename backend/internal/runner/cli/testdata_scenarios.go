package runner

import (
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/cli/seed_data"
)

// GameDataScenarioEntry holds a game data scenario: a config factory and a short description.
type GameDataScenarioEntry struct {
	Config      func() harness.DataConfig
	Description string
}

// GameDataScenarios is the registry of scenario name -> config and description.
// Used by db-load-game-data to resolve --scenario and by --list-scenarios.
var GameDataScenarios = map[string]GameDataScenarioEntry{
	"full": {
		Config:      seed_data.FullScenarioDataConfig,
		Description: "Full adventure game: locations, links, link requirements, items, creatures, characters, instances, accounts. Games loaded as draft. Use --publish to publish.",
	},
}
