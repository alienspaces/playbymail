package harness

import "gitlab.com/alienspaces/playbymail/internal/record"

const (
	GameOneRef     = "game-one"
	GameTwoRef     = "game-two"
	AccountOneRef  = "account-one"
	AccountTwoRef  = "account-two"
	LocationOneRef = "location-one"
	LocationTwoRef = "location-two"
)

// DataConfig -
type DataConfig struct {
	GameConfigs          []GameConfig
	AccountConfigs       []AccountConfig
	GameCharacterConfigs []GameCharacterConfig // Add this line
}

type GameConfig struct {
	Reference            string // Reference to the game record
	Record               *record.Game
	LocationConfigs      []LocationConfig      // Locations associated with this game
	LocationLinkConfigs  []LocationLinkConfig  // Links associated with this game
	GameCharacterConfigs []GameCharacterConfig // Add this line
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

type LocationConfig struct {
	Reference string // Reference to the location record
	Record    *record.Location
}

type LocationLinkConfig struct {
	Reference       string // Reference to the location link record
	FromLocationRef string // Reference to the from location
	ToLocationRef   string // Reference to the to location
	Record          *record.LocationLink
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
				LocationConfigs: []LocationConfig{
					{
						Reference: LocationOneRef,
						Record: &record.Location{
							Name:        "Default Location One",
							Description: "Default location one for handler tests",
						},
					},
					{
						Reference: LocationTwoRef,
						Record: &record.Location{
							Name:        "Default Location Two",
							Description: "Default location two for handler tests",
						},
					},
				},
				LocationLinkConfigs: []LocationLinkConfig{
					{
						Reference:       "link-one-two",
						FromLocationRef: LocationOneRef,
						ToLocationRef:   LocationTwoRef,
						Record: &record.LocationLink{
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
