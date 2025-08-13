package domain

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GameParameter - Different types of games may require different parameters.
// Rather than creating a new table for each game type, we can manage the
// available parameters per game type in code.

const (
	GameParameterValueTypeString  = "string"
	GameParameterValueTypeInteger = "integer"
	GameParameterValueTypeBoolean = "boolean"
	GameParameterValueTypeJSON    = "json"
)

const (
	AdventureGameParameterCharacterLives = "character_lives"
)

var gameParameterConfigurations = []game_record.GameParameter{
	{
		GameType:     game_record.GameTypeAdventure,
		ConfigKey:    AdventureGameParameterCharacterLives,
		Description:  nullstring.FromString("The number of lives a character has."),
		ValueType:    GameParameterValueTypeInteger,
		DefaultValue: nullstring.FromString("3"),
		IsRequired:   true,
		IsGlobal:     false,
	},
}

// GetGameParameterConfigurations returns all game parameter configurations
func GetGameParameterConfigurations() []*game_record.GameParameter {
	configs := make([]*game_record.GameParameter, len(gameParameterConfigurations))
	for i, config := range gameParameterConfigurations {
		// Create a copy to avoid modifying the original
		configCopy := config
		configs[i] = &configCopy
	}
	return configs
}

// GetGameParameterConfigurationsByGameType returns configurations filtered by game type
func GetGameParameterConfigurationsByGameType(gameType string) []*game_record.GameParameter {
	var filtered []*game_record.GameParameter
	for _, config := range gameParameterConfigurations {
		if config.GameType == gameType {
			// Create a copy to avoid modifying the original
			configCopy := config
			filtered = append(filtered, &configCopy)
		}
	}
	return filtered
}
