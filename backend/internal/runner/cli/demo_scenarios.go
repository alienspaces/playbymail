package runner

import (
	"sort"
	"strings"

	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/cli/demo_scenarios"
)

// DemoGameEntry holds a demo game: a config factory, game type, and description.
type DemoGameEntry struct {
	Config      func() harness.DataConfig
	GameType    string
	Description string
}

// DemoGames is the registry of game name -> entry.
// Used by db-load-demo-game to resolve --game and by --list.
var DemoGames = map[string]DemoGameEntry{
	demo_scenarios.DemoAdventureGameName: {
		Config:      demo_scenarios.AdventureGameConfig,
		GameType:    game_record.GameTypeAdventure,
		Description: "Locations, links, link requirements, items, creatures, instances, accounts. Loaded as draft. Use --publish to publish.",
	},
}

// DemoGameSummary is a read-only view of a registered demo game.
type DemoGameSummary struct {
	Name        string
	GameType    string
	Description string
}

// ListDemoGames returns a sorted slice of summaries for all registered demo games.
func ListDemoGames() []DemoGameSummary {
	summaries := make([]DemoGameSummary, 0, len(DemoGames))
	for name, entry := range DemoGames {
		summaries = append(summaries, DemoGameSummary{
			Name:        name,
			GameType:    entry.GameType,
			Description: entry.Description,
		})
	}
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Name < summaries[j].Name
	})
	return summaries
}

// LookupDemoGame finds a demo game entry by name (case-insensitive).
func LookupDemoGame(name string) (DemoGameEntry, bool) {
	if entry, ok := DemoGames[name]; ok {
		return entry, true
	}
	lower := strings.ToLower(name)
	for key, entry := range DemoGames {
		if strings.ToLower(key) == lower {
			return entry, true
		}
	}
	return DemoGameEntry{}, false
}
