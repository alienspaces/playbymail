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
	GameConfigs    []GameConfig
	AccountConfigs []AccountConfig
}

type GameConfig struct {
	Reference       string // Reference to the game record
	Record          *record.Game
	LocationConfigs []LocationConfig // Locations associated with this game
}

type AccountConfig struct {
	Reference string // Reference to the account record
	Record    *record.Account
}

type LocationConfig struct {
	Reference string // Reference to the location record
	Record    *record.Location
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
							Name:        UniqueName("Default Location One"),
							Description: "Default location one for handler tests",
						},
					},
					{
						Reference: LocationTwoRef,
						Record: &record.Location{
							Name:        UniqueName("Default Location Two"),
							Description: "Default location two for handler tests",
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
	}
}
