package seed_data

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// Image file names in seed_images directory
const (
	ImageJoinGame            = "join-game.png"
	ImageInventoryManagement = "inventory-management.png"
	ImageLocationDarkforest  = "location-darkforest.png"
	ImageLocationDungeon     = "location-dungeon.png"
	ImageLocationCliffpath   = "location-cliffpath.png"
)

// SeedDataConfig returns the seed data configuration for
// test data that can be used for setting up automated tests in
// the public space.
func SeedDataConfig() harness.DataConfig {
	return harness.DataConfig{
		GameConfigs:    GameConfig(),
		AccountConfigs: AccountConfig(),
	}
}

// AccountConfig returns test account configurations
func AccountConfig() []harness.AccountConfig {
	return []harness.AccountConfig{
		{
			Reference: harness.AccountOneRef,
			Record: &account_record.Account{
				Email: "test-account-one@example.com",
			},
		},
		{
			Reference: harness.AccountTwoRef,
			Record: &account_record.Account{
				Email: "test-account-two@example.com",
			},
		},
		{
			Reference: harness.AccountThreeRef,
			Record: &account_record.Account{
				Email: "test-account-three@example.com",
			},
		},
	}
}

// GameConfig returns the seed data configuration for games
func GameConfig() []harness.GameConfig {
	return []harness.GameConfig{
		{
			Reference: harness.GameOneRef,
			Record: &game_record.Game{
				Name:              "The Enchanted Forest Adventure",
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 168, // 1 week
			},
			// Game images for turn sheet backgrounds (loaded from seed_images/)
			GameImageConfigs: []harness.GameImageConfig{
				{
					Reference:     harness.GameImageJoinGameRef,
					ImagePath:     ImageJoinGame,
					TurnSheetType: adventure_game_record.AdventureGameTurnSheetTypeJoinGame,
				},
				{
					Reference:     harness.GameImageInventoryRef,
					ImagePath:     ImageInventoryManagement,
					TurnSheetType: adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement,
				},
			},
			// Rich world with multiple interconnected locations
			AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
				{
					Reference: harness.GameLocationOneRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:               "Mystic Grove",
						Description:        "A peaceful grove filled with ancient trees and magical flowers. The air shimmers with enchantment.",
						IsStartingLocation: true,
					},
					BackgroundImagePath: ImageLocationDarkforest,
				},
				{
					Reference: harness.GameLocationTwoRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Crystal Caverns",
						Description: "Deep underground caves filled with glowing crystals. Strange sounds echo from the depths.",
					},
					BackgroundImagePath: ImageLocationDungeon,
				},
				{
					Reference: harness.GameLocationThreeRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Floating Islands",
						Description: "Mysterious islands suspended in the sky by unknown magic. Wind howls between them.",
					},
					BackgroundImagePath: ImageLocationCliffpath,
				},
				{
					Reference: harness.GameLocationFourRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Shadow Valley",
						Description: "A dark valley shrouded in perpetual shadows. Danger lurks in every corner.",
					},
					BackgroundImagePath: ImageLocationDarkforest,
				},
			},
			// Items that can be found and used
			AdventureGameItemConfigs: []harness.AdventureGameItemConfig{
				{
					Reference: harness.GameItemOneRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Crystal Key",
						Description: "A glowing key made of pure crystal. It hums with magical energy.",
					},
				},
				{
					Reference: harness.GameItemTwoRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Shadow Cloak",
						Description: "A cloak that allows the wearer to blend into shadows and move silently.",
					},
				},
				{
					Reference: harness.GameItemThreeRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Healing Potion",
						Description: "A bright blue potion that restores health and removes minor ailments.",
					},
				},
				{
					Reference: harness.GameItemFourRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Wind Charm",
						Description: "A small charm that allows the bearer to control wind currents.",
					},
				},
			},
			// Creatures that inhabit the world
			AdventureGameCreatureConfigs: []harness.AdventureGameCreatureConfig{
				{
					Reference: harness.GameCreatureOneRef,
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Forest Guardian",
						Description: "A majestic creature made of living wood and leaves. Protects the grove.",
					},
				},
				{
					Reference: harness.GameCreatureTwoRef,
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Crystal Spider",
						Description: "A giant spider with a body made of living crystal. Spins webs of light.",
					},
				},
			},
			// Characters that players can control
			AdventureGameCharacterConfigs: []harness.AdventureGameCharacterConfig{
				{
					Reference:  harness.GameCharacterOneRef,
					AccountRef: harness.AccountOneRef,
					Record: &adventure_game_record.AdventureGameCharacter{
						Name: "Aria the Mage",
					},
				},
				{
					Reference:  harness.GameCharacterTwoRef,
					AccountRef: harness.AccountTwoRef,
					Record: &adventure_game_record.AdventureGameCharacter{
						Name: "Thorne the Warrior",
					},
				},
				{
					Reference:  harness.GameCharacterThreeRef,
					AccountRef: harness.AccountThreeRef,
					Record: &adventure_game_record.AdventureGameCharacter{
						Name: "Luna the Scout",
					},
				},
			},
			// Location links that connect the world
			// Each location has 2-3 travel options for a richer interconnected world
			AdventureGameLocationLinkConfigs: []harness.AdventureGameLocationLinkConfig{
				// From Mystic Grove (starting location)
				{
					Reference:       harness.GameLocationLinkOneRef,
					FromLocationRef: harness.GameLocationOneRef,
					ToLocationRef:   harness.GameLocationTwoRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Crystal Path",
						Description: "A winding path that leads from the grove down into the crystal caverns.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   harness.GameLocationLinkRequirementOneRef,
							GameItemRef: harness.GameItemOneRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				{
					Reference:       "location-link-grove-to-islands",
					FromLocationRef: harness.GameLocationOneRef,
					ToLocationRef:   harness.GameLocationThreeRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Sky Vine",
						Description: "A massive enchanted vine that grows from the grove up to the floating islands.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "link-req-grove-to-islands",
							GameItemRef: harness.GameItemFourRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				{
					Reference:       "location-link-grove-to-shadow",
					FromLocationRef: harness.GameLocationOneRef,
					ToLocationRef:   harness.GameLocationFourRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Dark Tunnel",
						Description: "A hidden tunnel beneath ancient roots that leads directly to the shadow valley.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "link-req-grove-to-shadow",
							GameItemRef: harness.GameItemTwoRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				// From Crystal Caverns
				{
					Reference:       harness.GameLocationLinkTwoRef,
					FromLocationRef: harness.GameLocationTwoRef,
					ToLocationRef:   harness.GameLocationThreeRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Wind Lift",
						Description: "A magical elevator that rises from the caverns to the floating islands.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   harness.GameLocationLinkRequirementTwoRef,
							GameItemRef: harness.GameItemFourRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				{
					Reference:       "location-link-caverns-to-grove",
					FromLocationRef: harness.GameLocationTwoRef,
					ToLocationRef:   harness.GameLocationOneRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Glowing Passage",
						Description: "Crystal-lit tunnels that wind back up to the mystic grove.",
					},
				},
				{
					Reference:       "location-link-caverns-to-shadow",
					FromLocationRef: harness.GameLocationTwoRef,
					ToLocationRef:   harness.GameLocationFourRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Deep Descent",
						Description: "A treacherous path that plunges from the caverns into the shadow valley.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "link-req-caverns-to-shadow",
							GameItemRef: harness.GameItemThreeRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				// From Floating Islands
				{
					Reference:       harness.GameLocationLinkThreeRef,
					FromLocationRef: harness.GameLocationThreeRef,
					ToLocationRef:   harness.GameLocationFourRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Shadow Bridge",
						Description: "A bridge of pure darkness that connects the floating islands to the shadow valley.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   harness.GameLocationLinkRequirementThreeRef,
							GameItemRef: harness.GameItemTwoRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				{
					Reference:       "location-link-islands-to-grove",
					FromLocationRef: harness.GameLocationThreeRef,
					ToLocationRef:   harness.GameLocationOneRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Feather Fall",
						Description: "Enchanted feathers allow a gentle descent from the islands to the grove below.",
					},
				},
				{
					Reference:       "location-link-islands-to-caverns",
					FromLocationRef: harness.GameLocationThreeRef,
					ToLocationRef:   harness.GameLocationTwoRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Crystal Chute",
						Description: "A smooth crystal slide that spirals down into the caverns.",
					},
				},
				// From Shadow Valley
				{
					Reference:       harness.GameLocationLinkFourRef,
					FromLocationRef: harness.GameLocationFourRef,
					ToLocationRef:   harness.GameLocationOneRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Return Portal",
						Description: "A magical portal that allows quick return to the mystic grove.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   harness.GameLocationLinkRequirementFourRef,
							GameItemRef: harness.GameItemThreeRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				{
					Reference:       "location-link-shadow-to-caverns",
					FromLocationRef: harness.GameLocationFourRef,
					ToLocationRef:   harness.GameLocationTwoRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Echoing Stairs",
						Description: "Ancient stone stairs that climb from the valley up to the crystal caverns.",
					},
				},
				{
					Reference:       "location-link-shadow-to-islands",
					FromLocationRef: harness.GameLocationFourRef,
					ToLocationRef:   harness.GameLocationThreeRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Shadow Ascent",
						Description: "Dark tendrils of shadow that can lift travelers to the floating islands.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "link-req-shadow-to-islands",
							GameItemRef: harness.GameItemTwoRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
			},
			// Game subscriptions - Designer and Manager subscriptions for game ownership
			GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
				{
					Reference:        harness.GameSubscriptionOneRef,
					AccountRef:       harness.AccountOneRef,
					SubscriptionType: game_record.GameSubscriptionTypeManager,
					Record:           &game_record.GameSubscription{},
				},
				{
					Reference:        "game-subscription-one-designer",
					AccountRef:       harness.AccountOneRef,
					SubscriptionType: game_record.GameSubscriptionTypeDesigner,
					Record:           &game_record.GameSubscription{},
				},
			},
			// Game instances with all the resources
			GameInstanceConfigs: []harness.GameInstanceConfig{
				{
					Reference: harness.GameInstanceOneRef,
					Record:    &game_record.GameInstance{},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: "game-instance-parameter-one",
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   domain.AdventureGameParameterCharacterLives,
								ParameterValue: nullstring.FromString("5"),
							},
						},
					},
					// Location instances for this game instance
					AdventureGameLocationInstanceConfigs: []harness.AdventureGameLocationInstanceConfig{
						{
							Reference:       harness.GameLocationInstanceOneRef,
							GameLocationRef: harness.GameLocationOneRef,
							Record:          &adventure_game_record.AdventureGameLocationInstance{},
						},
						{
							Reference:       harness.GameLocationInstanceTwoRef,
							GameLocationRef: harness.GameLocationTwoRef,
							Record:          &adventure_game_record.AdventureGameLocationInstance{},
						},
					},
					// Item instances placed in the world
					AdventureGameItemInstanceConfigs: []harness.AdventureGameItemInstanceConfig{
						{
							Reference:       harness.GameItemInstanceOneRef,
							GameItemRef:     harness.GameItemOneRef,
							GameLocationRef: harness.GameLocationOneRef,
							Record:          &adventure_game_record.AdventureGameItemInstance{},
						},
					},
					// Creature instances in the world
					AdventureGameCreatureInstanceConfigs: []harness.AdventureGameCreatureInstanceConfig{
						{
							Reference:       harness.GameCreatureInstanceOneRef,
							GameCreatureRef: harness.GameCreatureOneRef,
							GameLocationRef: harness.GameLocationOneRef,
							Record:          &adventure_game_record.AdventureGameCreatureInstance{},
						},
					},
					// Character instances
					AdventureGameCharacterInstanceConfigs: []harness.AdventureGameCharacterInstanceConfig{
						{
							Reference:        harness.GameCharacterInstanceOneRef,
							GameCharacterRef: harness.GameCharacterOneRef,
							GameLocationRef:  harness.GameLocationOneRef,
							Record:           &adventure_game_record.AdventureGameCharacterInstance{},
						},
					},
				},
			},
		},
		{
			Reference: harness.GameTwoRef,
			Record: &game_record.Game{
				Name:              "The Desert Kingdom",
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 336, // 2 weeks
			},
			// Designer and Manager subscriptions for game ownership
			GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
				{
					Reference:        "game-subscription-two",
					AccountRef:       harness.AccountTwoRef,
					SubscriptionType: game_record.GameSubscriptionTypeManager,
					Record:           &game_record.GameSubscription{},
				},
				{
					Reference:        "game-subscription-two-designer",
					AccountRef:       harness.AccountTwoRef,
					SubscriptionType: game_record.GameSubscriptionTypeDesigner,
					Record:           &game_record.GameSubscription{},
				},
			},
			// Simpler world for the second game
			AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
				{
					Reference: "game-location-five",
					Record: &adventure_game_record.AdventureGameLocation{
						Name:               "Oasis Village",
						Description:        "A bustling village built around a life-giving oasis in the desert.",
						IsStartingLocation: true,
					},
				},
				{
					Reference: "game-location-six",
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Ancient Ruins",
						Description: "Crumbling ruins of a lost civilization, filled with secrets and danger.",
					},
				},
			},
			AdventureGameItemConfigs: []harness.AdventureGameItemConfig{
				{
					Reference: "game-item-five",
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Desert Compass",
						Description: "A magical compass that always points to water sources.",
					},
				},
			},
			AdventureGameCreatureConfigs: []harness.AdventureGameCreatureConfig{
				{
					Reference: "game-creature-three",
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Sand Serpent",
						Description: "A massive serpent that burrows through the desert sands.",
					},
				},
			},
			GameInstanceConfigs: []harness.GameInstanceConfig{
				{
					Reference: harness.GameInstanceTwoRef,
					Record:    &game_record.GameInstance{},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: "game-instance-parameter-three",
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   domain.AdventureGameParameterCharacterLives,
								ParameterValue: nullstring.FromString("3"),
							},
						},
					},
				},
			},
		},
	}
}
