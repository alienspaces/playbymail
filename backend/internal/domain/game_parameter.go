package domain

import (
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GameParameter - Different types of games may require different parameters.
// Rather than creating a new table for each game type, we can manage the
// available parameters per game type in code.
//
// Available parameters do not set the default values, whether the parameter is
// required or whether the paramter is global and cannot be overriden at the
// instance level. Those values are set when the parameters are set for the
// specific game and then potentially overriden when an instance is created.
//
// The majority of available game parmeters should have default values applied
// by the specific game type engine when not otherwise provided as best practice.
const (
	GameParameterValueTypeString  = "string"
	GameParameterValueTypeInteger = "integer"
	GameParameterValueTypeBoolean = "boolean"
)

const (
	AdventureGameParameterCharacterLives = "character_lives"
)

var gameParameters = []game_record.GameParameter{
	// Adventure game parameters
	{
		GameType:     game_record.GameTypeAdventure,
		ConfigKey:    AdventureGameParameterCharacterLives,
		Description:  "The number of lives a character has.",
		ValueType:    GameParameterValueTypeInteger,
		DefaultValue: "3",
	},
}

// GetGameParameters returns all game parameters
func GetGameParameters() []*game_record.GameParameter {
	configs := make([]*game_record.GameParameter, len(gameParameters))
	for i, config := range gameParameters {
		// Create a copy to avoid modifying the original
		configCopy := config
		configs[i] = &configCopy
	}
	return configs
}

// GetGameParametersByGameType returns parameters filtered by game type
func GetGameParametersByGameType(gameType string) []*game_record.GameParameter {
	var filtered []*game_record.GameParameter
	for _, config := range gameParameters {
		if config.GameType == gameType {
			// Create a copy to avoid modifying the original
			configCopy := config
			filtered = append(filtered, &configCopy)
		}
	}
	return filtered
}
