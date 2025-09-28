package adventure_game

import (
	"context"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// AdventureGame coordinates turn processing for adventure games
type AdventureGame struct {
	Logger                  logger.Logger
	Domain                  *domain.Domain
	LocationChoiceProcessor *AdventureGameLocationChoiceProcessor
}

// NewAdventureGame creates a new adventure game turn processor
func NewAdventureGame(l logger.Logger, d *domain.Domain) *AdventureGame {
	return &AdventureGame{
		Logger: l,
		Domain: d,
		// Additional processors can be added here
		LocationChoiceProcessor: NewAdventureGameLocationChoiceProcessor(l, d),
	}
}

// getCharacterInstancesForGameInstance retrieves all character instances for a game instance
func (p *AdventureGame) getCharacterInstancesForGameInstance(_ context.Context, gameInstanceRec *game_record.GameInstance) ([]*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := p.Logger.WithFunctionContext("AdventureGame/getCharacterInstancesForGameInstance")

	characterInstanceRecs, err := p.Domain.GetManyAdventureGameCharacterInstanceRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: adventure_game_record.FieldAdventureGameCharacterInstanceGameInstanceID,
					Val: gameInstanceRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Error("failed to get character instances error >%v<", err)
		return nil, err
	}

	return characterInstanceRecs, nil
}
