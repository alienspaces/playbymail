package harness

import "gitlab.com/alienspaces/playbymail/internal/record"

const (
	GameOneRef         = "game-one"
	GameTwoRef         = "game-two"
	AccountOneRef      = "account-one"
	AccountTwoRef      = "account-two"
	GameLocationOneRef = "game-location-one"
	GameLocationTwoRef = "game-location-two"
)

// DataConfig -
type DataConfig struct {
	GameConfigs          []GameConfig
	AccountConfigs       []AccountConfig
	GameCharacterConfigs []GameCharacterConfig // Add this line
}

type GameConfig struct {
	Reference               string // Reference to the game record
	Record                  *record.Game
	GameLocationConfigs     []GameLocationConfig     // Locations associated with this game
	GameLocationLinkConfigs []GameLocationLinkConfig // Links associated with this game
	GameCharacterConfigs    []GameCharacterConfig    // Add this line
}

type GameCharacterConfig struct {
	Reference  string // Reference to the game_character record
	AccountRef string // Reference to the account
	Record     *record.GameCharacter
}

type AccountConfig struct {
	Reference string // Reference to the account record
	Record    *record.Account
}

type GameLocationConfig struct {
	Reference string // Reference to the game_location record
	Record    *record.GameLocation
}

type GameLocationLinkConfig struct {
	Reference       string // Reference to the location link record
	FromLocationRef string // Reference to the from location
	ToLocationRef   string // Reference to the to location
	Record          *record.GameLocationLink
}

// DefaultDataConfig -
func DefaultDataConfig() DataConfig {
	return DataConfig{
		GameConfigs: []GameConfig{
			{
				Reference: GameOneRef,
				Record: &record.Game{
					Name:     "Default Game One",
					GameType: record.GameTypeAdventure,
				},
				GameLocationConfigs: []GameLocationConfig{
					{
						Reference: GameLocationOneRef,
						Record: &record.GameLocation{
							Name:        "Default Location One",
							Description: "Default location one for handler tests",
						},
					},
					{
						Reference: GameLocationTwoRef,
						Record: &record.GameLocation{
							Name:        "Default Location Two",
							Description: "Default location two for handler tests",
						},
					},
				},
				GameLocationLinkConfigs: []GameLocationLinkConfig{
					{
						Reference:       "link-one-two",
						FromLocationRef: GameLocationOneRef,
						ToLocationRef:   GameLocationTwoRef,
						Record: &record.GameLocationLink{
							Description: "Travel by boat to the swamp of the long forgotten Frog God",
							Name:        "The Red Door",
						},
					},
				},
				GameCharacterConfigs: []GameCharacterConfig{
					{
						Reference:  "character-one",
						AccountRef: AccountOneRef,
						Record: &record.GameCharacter{
							Name: "Default Character One",
						},
					},
				},
			},
		},
		AccountConfigs: []AccountConfig{
			{
				Reference: AccountOneRef,
				Record: &record.Account{
					Email: UniqueEmail("default-account-one@example.com"),
				},
			},
			{
				Reference: AccountTwoRef,
				Record: &record.Account{
					Email: UniqueEmail("default-account-two@example.com"),
				},
			},
		},
		GameCharacterConfigs: []GameCharacterConfig{}, // No global characters by default
	}
}
