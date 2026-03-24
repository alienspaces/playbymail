package turnsheet

import (
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// AdventureGameLocationChoiceFixture returns the sample rendering fixture for the
// adventure game location choice turn sheet.
func AdventureGameLocationChoiceFixture() DevFixture {
	return DevFixture{
		TemplatePath:   "turnsheet/adventure_game_location_choice.template",
		OutputBaseName: "adventure_game_location_choice_turnsheet",
		BackgroundFile: "background-cliffpath.png",
		MakeData: func(bg, code string) any {
			deadline := time.Now().Add(7 * 24 * time.Hour)
			return &LocationChoiceData{
				TurnSheetTemplateData: TurnSheetTemplateData{
					GameName:              strPtr("The Enchanted Forest Adventure"),
					GameType:              strPtr("adventure"),
					TurnNumber:            intPtr(1),
					AccountName:           strPtr("Test Player"),
					TurnSheetTitle:        strPtr("Location Choice"),
					TurnSheetInstructions: strPtr(DefaultLocationChoiceInstructions()),
					TurnSheetCode:         strPtr(code),
					TurnSheetDeadline:     &deadline,
					BackgroundImage:       &bg,
					TurnEvents: []TurnEvent{
						{Category: TurnEventCategoryMovement, Icon: TurnEventIconMovement, Message: "You arrived at Mystic Grove after a long journey through the forest."},
						{Category: TurnEventCategorySystem, Icon: TurnEventIconSystem, Message: "The ancient trees seem to watch your every move."},
						{Category: TurnEventCategorySystem, Icon: TurnEventIconSystem, Message: "You discovered a hidden path leading north."},
					},
				},
				LocationName:        "Mystic Grove",
				LocationDescription: "You stand at the edge of an ancient forest. The trees whisper secrets of old magic.",
				LocationOptions: []LocationOption{
					{LocationID: "crystal_caverns", LocationLinkName: "Crystal Caverns", LocationLinkDescription: "Enter the glowing caverns where crystals hum with power"},
					{LocationID: "dark_tower", LocationLinkName: "Dark Tower", LocationLinkDescription: "Climb the mysterious tower that pierces the sky"},
					{LocationID: "sunset_plains", LocationLinkName: "Sunset Plains", LocationLinkDescription: "Venture into the vast plains where the sun sets eternally"},
					{LocationID: "mermaid_lagoon", LocationLinkName: "Mermaid Lagoon", LocationLinkDescription: "Dive into the hidden lagoon where mermaids sing"},
				},
			}
		},
		NewProcessor: func(l logger.Logger, cfg config.Config) (TurnSheetProcessor, error) {
			return NewLocationChoiceProcessor(l, cfg)
		},
	}
}
