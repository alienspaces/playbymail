package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameRec(gameConfig GameConfig) (*game_record.Game, error) {
	l := t.Logger("createGameRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *game_record.Game
	if gameConfig.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *gameConfig.Record
		rec = &recCopy
	} else {
		rec = &game_record.Game{}
	}

	rec = t.applyGameRecDefaultValues(rec)

	// Create record
	l.Debug("creating game record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameRec(rec)
	if err != nil {
		l.Warn("failed creating game record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameRec(rec)

	// Add to references store
	if gameConfig.Reference != "" {
		t.Data.Refs.GameRefs[gameConfig.Reference] = rec.ID
	}

	// Create background image for game if specified (legacy support)
	if gameConfig.BackgroundImagePath != "" {
		_, err = t.createGameImageRecFromPath(rec.ID, "", gameConfig.BackgroundImagePath)
		if err != nil {
			l.Warn("failed creating game background image >%v<", err)
			return nil, err
		}
		l.Debug("created game background image for game >%s<", rec.ID)
	}

	// Create game images from config
	for _, imageConfig := range gameConfig.GameImageConfigs {
		_, err = t.createGameImageRecFromConfig(rec.ID, imageConfig)
		if err != nil {
			l.Warn("failed creating game image from config >%v<", err)
			return nil, err
		}
		l.Debug("created game image from config for game >%s<", rec.ID)
	}

	return rec, nil
}

func (t *Testing) applyGameRecDefaultValues(rec *game_record.Game) *game_record.Game {
	if rec == nil {
		rec = &game_record.Game{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}

	if rec.GameType == "" {
		rec.GameType = game_record.GameTypeAdventure
	}

	if rec.TurnDurationHours == 0 {
		rec.TurnDurationHours = 168 // Default to 1 week
	}

	if rec.Description == "" {
		rec.Description = fmt.Sprintf("Welcome to %s! Welcome to the PlayByMail Adventure!", rec.Name)
	}

	return rec
}
