package harness

import "gitlab.com/alienspaces/playbymail/internal/record"

const (
	GameOneRef = "game-one"
)

// DataConfig -
type DataConfig struct {
	GameConfig []GameConfig
}

type GameConfig struct {
	Reference string // Reference to the game record
	Record    *record.Game
}

// DefaultDataConfig -
var DefaultDataConfig = DataConfig{
	GameConfig: []GameConfig{
		{
			Reference: GameOneRef,
			Record:    &record.Game{},
		},
	},
}
