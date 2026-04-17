package turnsheet

import (
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// MechaJoinGameFixture returns the sample rendering fixture for the
// mecha join game turn sheet.
func MechaJoinGameFixture() DevFixture {
	return DevFixture{
		TemplatePath:   "turnsheet/mecha_join_game.template",
		OutputBaseName: "mecha_join_game_turnsheet",
		BackgroundFile: "background-darkforest.png",
		IsJoinSheet:    true,
		MakeData: func(bg, code string) any {
			deadline := time.Now().Add(7 * 24 * time.Hour)
			return &JoinGameData{
				TurnSheetTemplateData: TurnSheetTemplateData{
					GameName:              strPtr("Steel Thunder"),
					GameType:              strPtr("mecha"),
					TurnNumber:            intPtr(0),
					TurnSheetTitle:        strPtr("Join Game"),
					TurnSheetInstructions: strPtr(DefaultMechaJoinGameInstructions()),
					TurnSheetCode:         strPtr(code),
					TurnSheetDeadline:     &deadline,
					BackgroundImage:       &bg,
					HideNarrative:         true,
				},
				GameDescription:          "Command a squad of powerful war mechs!",
				AvailableDeliveryMethods: DeliveryMethods{Email: true},
			}
		},
		NewProcessor: func(l logger.Logger, cfg config.Config) (TurnSheetProcessor, error) {
			return NewMechaJoinGameProcessor(l, cfg)
		},
	}
}
