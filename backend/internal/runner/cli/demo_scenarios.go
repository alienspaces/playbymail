package runner

import (
	"sort"
	"strings"

	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/cli/demo_scenarios"
)

// DemoGameEntry holds a demo game: a config factory, game name, and description.
type DemoGameEntry struct {
	Config      func() harness.DataConfig
	Name        string
	Description string
}

// DemoGames is the registry of game type -> entry.
// Used by db-load-demo-game to resolve --type and by --list.
var DemoGames = map[string]DemoGameEntry{
	game_record.GameTypeAdventure: {
		Config:      demo_scenarios.AdventureGameConfig,
		Name:        demo_scenarios.DemoAdventureGameName,
		Description: "A solo text adventure set in a mysterious house.",
	},
}

// DemoGameSummary is a read-only view of a registered demo game.
type DemoGameSummary struct {
	Type        string
	Name        string
	Description string
}

// ListDemoGames returns a sorted slice of summaries for all registered demo games.
func ListDemoGames() []DemoGameSummary {
	summaries := make([]DemoGameSummary, 0, len(DemoGames))
	for gameType, entry := range DemoGames {
		summaries = append(summaries, DemoGameSummary{
			Type:        gameType,
			Name:        entry.Name,
			Description: entry.Description,
		})
	}
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Type < summaries[j].Type
	})
	return summaries
}

// LookupDemoGameByType finds a demo game entry by game type (case-insensitive).
func LookupDemoGameByType(gameType string) (DemoGameEntry, bool) {
	if entry, ok := DemoGames[gameType]; ok {
		return entry, true
	}
	for key, entry := range DemoGames {
		if strings.EqualFold(key, gameType) {
			return entry, true
		}
	}
	return DemoGameEntry{}, false
}
