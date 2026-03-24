package turnsheet

import (
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// AdventureGameMonsterEncounterFixture returns the sample rendering fixture for the
// adventure game monster encounter turn sheet.
func AdventureGameMonsterEncounterFixture() DevFixture {
	return DevFixture{
		TemplatePath:   "turnsheet/adventure_game_monster_encounter.template",
		OutputBaseName: "adventure_game_monster_encounter_turnsheet",
		BackgroundFile: "background-darkforest.png",
		MakeData: func(bg, code string) any {
			deadline := time.Now().Add(7 * 24 * time.Hour)
			return &MonsterEncounterData{
				TurnSheetTemplateData: TurnSheetTemplateData{
					GameName:              strPtr("The Door Beneath the Staircase"),
					GameType:              strPtr("adventure"),
					TurnNumber:            intPtr(3),
					AccountName:           strPtr("Test Player"),
					TurnSheetTitle:        strPtr("Creature Encounter"),
					TurnSheetInstructions: strPtr(DefaultMonsterEncounterInstructions(3)),
					TurnSheetCode:         strPtr(code),
					TurnSheetDeadline:     &deadline,
					BackgroundImage:       &bg,
					TurnEvents: []TurnEvent{
						{Category: TurnEventCategoryMovement, Icon: TurnEventIconMovement, Message: "Aldric entered the shadowed corridor."},
						{Category: TurnEventCategoryCombat, Icon: TurnEventIconCombat, Message: "Aldric attacked the Goblin Scout for 8 damage."},
						{Category: TurnEventCategoryCombat, Icon: TurnEventIconCombat, Message: "Goblin Scout retaliated for 4 damage."},
						{Category: TurnEventCategorySystem, Icon: TurnEventIconSystem, Message: "Aldric's health dropped to 65."},
					},
				},
				CharacterName:      "Aldric",
				CharacterHealth:    65,
				CharacterMaxHealth: 100,
				CharacterAttack:    8,
				CharacterDefense:   3,
				EquippedWeapon: &EquippedWeapon{
					ItemInstanceID: "weapon-1",
					Name:           "Iron Sword",
					Damage:         8,
				},
				EquippedArmor: &EquippedArmor{
					ItemInstanceID: "armor-1",
					Name:           "Leather Jerkin",
					Defense:        3,
				},
				Creatures: []EncounterCreature{
					{
						CreatureInstanceID: "creature-1",
						Name:               "Sand Serpent",
						Description:        "A massive serpent that lurks beneath the desert sands, striking with terrifying speed.",
						Health:             80,
						MaxHealth:          100,
						AttackDamage:       12,
						Defense:            2,
						Disposition:        "aggressive",
					},
					{
						CreatureInstanceID: "creature-2",
						Name:               "Desert Scorpion",
						Description:        "A venomous scorpion the size of a dog. Its tail curls menacingly.",
						Health:             35,
						MaxHealth:          50,
						AttackDamage:       8,
						Defense:            4,
						Disposition:        "aggressive",
					},
				},
				MaxActions: 3,
			}
		},
		NewProcessor: func(l logger.Logger, cfg config.Config) (TurnSheetProcessor, error) {
			return NewMonsterEncounterProcessor(l, cfg)
		},
	}
}
