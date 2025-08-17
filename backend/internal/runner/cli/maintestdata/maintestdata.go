package maintestdata

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
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
			Record: &game_record.Game{
				Name:              "Test Game One",
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 168, // 1 week
			},
			GameInstanceConfigs: []harness.GameInstanceConfig{
				{
					Reference: "game-instance-one",
					Record:    &game_record.GameInstance{},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: "game-instance-parameter-one",
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   "character_lives",
								ParameterValue: nullstring.FromString("5"),
							},
						},
					},
				},
			},
		},
		{
			Reference: "game-two",
			Record: &game_record.Game{
				Name:              "Test Game Two",
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 336, // 2 weeks
			},
			GameInstanceConfigs: []harness.GameInstanceConfig{
				{
					Reference: "game-instance-two",
					Record:    &game_record.GameInstance{},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: "game-instance-parameter-two",
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   "character_lives",
								ParameterValue: nullstring.FromString("3"),
							},
						},
					},
				},
			},
		},
	}
}
