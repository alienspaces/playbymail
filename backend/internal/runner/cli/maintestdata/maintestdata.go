package maintestdata

import (
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// MainTestDataConfig returns the main test data configuration for
// test data that can be used for setting up automated tests in
// the public space.
func MainTestDataConfig() harness.DataConfig {
	return harness.DataConfig{
		GameConfigs: GameConfig(),
	}
}

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
