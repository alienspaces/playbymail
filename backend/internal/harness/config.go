package harness

import "gitlab.com/alienspaces/playbymail/internal/record"

const (
	GameOneRef    = "game-one"
	AccountOneRef = "account-one"
)

// DataConfig -
type DataConfig struct {
	GameConfig    []GameConfig
	AccountConfig []AccountConfig
}

type GameConfig struct {
	Reference string // Reference to the game record
	Record    *record.Game
}

type AccountConfig struct {
	Reference string // Reference to the account record
	Record    *record.Account
}

// DefaultDataConfig -
var DefaultDataConfig = DataConfig{
	GameConfig: []GameConfig{
		{
			Reference: GameOneRef,
			Record:    &record.Game{},
		},
	},
	AccountConfig: []AccountConfig{
		{
			Reference: AccountOneRef,
			Record:    &record.Account{},
		},
	},
}
