package maintestdata

import (
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GameConfig returns the main test data configuration for games
func GameConfig() []harness.GameConfig {
	return []harness.GameConfig{
		{
			Reference: "game-one",
			Record: &record.Game{
				Name: "Test Game One",
			},
		},
		{
			Reference: "game-two",
			Record: &record.Game{
				Name: "Test Game Two",
			},
		},
	}
}
