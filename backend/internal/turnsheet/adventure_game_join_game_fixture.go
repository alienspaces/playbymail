package turnsheet

import (
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// AdventureGameJoinGameFixture returns the sample rendering fixture for the
// adventure game join game turn sheet.
func AdventureGameJoinGameFixture() DevFixture {
	return DevFixture{
		TemplatePath:   "turnsheet/adventure_game_join_game.template",
		OutputBaseName: "adventure_game_join_game_turnsheet",
		BackgroundFile: "background-darkforest.png",
		IsJoinSheet:    true,
		MakeData: func(bg, code string) any {
			deadline := time.Now().Add(7 * 24 * time.Hour)
			return &JoinGameData{
				TurnSheetTemplateData: TurnSheetTemplateData{
					GameName:              strPtr("The Enchanted Forest Adventure"),
					GameType:              strPtr("adventure"),
					TurnNumber:            intPtr(0),
					TurnSheetTitle:        strPtr("Join Game"),
					TurnSheetInstructions: strPtr(DefaultJoinGameInstructions()),
					TurnSheetCode:         strPtr(code),
					TurnSheetDeadline:     &deadline,
					BackgroundImage:       &bg,
					HideNarrative:         true,
				},
				GameDescription: "Welcome to the PlayByMail Adventure!",
			}
		},
		NewProcessor: func(l logger.Logger, cfg config.Config) (TurnSheetProcessor, error) {
			return NewJoinGameProcessor(l, cfg)
		},
	}
}
