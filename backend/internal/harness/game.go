package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) processGameConfig(gameConfig GameConfig) (*game_record.Game, []*game_record.GameImage, error) {
	l := t.Logger("processGameConfig")

	var allGameImageRecs []*game_record.GameImage

	// Create a new record instance to avoid reusing the same record across tests
	var rec *game_record.Game
	if gameConfig.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *gameConfig.Record
		rec = &recCopy
	}

	rec = t.applyGameRecDefaultValues(rec)

	// Create record
	l.Debug("creating game record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameRec(rec)
	if err != nil {
		l.Warn("failed creating game record >%v<", err)
		return nil, nil, err
	}

	// Add to data store
	t.Data.AddGameRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameRec(rec)

	// Add to references store
	if gameConfig.Reference != "" {
		t.Data.Refs.GameRefs[gameConfig.Reference] = rec.ID
	}

	// Create game images from config
	for _, imageConfig := range gameConfig.GameImageConfigs {
		gameImageRec, err := t.processGameImageConfig(imageConfig, rec)
		if err != nil {
			l.Warn("failed processing game image config >%v<", err)
			return nil, nil, err
		}
		l.Debug("created game image from config for game >%s<", rec.ID)
		allGameImageRecs = append(allGameImageRecs, gameImageRec)
	}

	return rec, allGameImageRecs, nil
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
		rec.TurnDurationHours = 48 // Default to 2 days
	}

	if rec.Description == "" {
		rec.Description = fmt.Sprintf("Welcome to %s! Welcome to the PlayByMail Adventure!", rec.Name)
	}

	if rec.Status == "" {
		rec.Status = game_record.GameStatusPublished
	}

	return rec
}
