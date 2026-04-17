package test_data

import (
	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/nullint32"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

// Image file names in test_data_images directory
const (
	ImageJoinGame            = "join-game.png"
	ImageInventoryManagement = "inventory-management.png"
	ImageLocationDarkforest  = "location-darkforest.png"
	ImageLocationDungeon     = "location-dungeon.png"
	ImageLocationCliffpath   = "location-cliffpath.png"

	// Desert Kingdom images
	ImageDesertJoinGame  = "desert-join-game.jpg"
	ImageDesertInventory = "desert-inventory.jpg"

	// Desert Kingdom location choice turn sheet background images
	ImageDesertOasis  = "desert-oasis.jpg"
	ImageDesertRuins  = "desert-ruins.jpg"
	ImageDesertCanyon = "desert-canyon.jpg"
	ImageDesertTemple = "desert-temple.jpg"

	// Steel Thunder (mecha) images
	ImageSteelJoinGame = "steel-join-game.jpg"
	ImageSteelOrders   = "steel-orders.jpg"

	// Iron Vanguard (mecha) images
	ImageIronJoinGame = "iron-join-game.jpg"
	ImageIronOrders   = "iron-orders.jpg"
)

// TestDataConfig returns the test data configuration for E2E and Playwright tests in the public space.
func TestDataConfig() harness.DataConfig {
	return harness.DataConfig{
		AccountConfigs:                     AccountConfig(),
		GameConfigs:                        GameConfig(),
		AccountUserGameSubscriptionConfigs: AccountUserGameSubscriptionConfig(),
	}
}

// AccountConfig returns test account configurations that will later be subscribed to games
func AccountConfig() []harness.AccountConfig {
	return []harness.AccountConfig{
		{
			Reference: harness.AccountStandardRef,
			AccountUserConfigs: []harness.AccountUserConfig{
				{
					Reference: harness.AccountUserStandardRef,
				Record: &account_record.AccountUser{
					Email: "test-player@example.com",
				},
				},
			},
		},
		{
			Reference: harness.AccountProPlayerRef,
			AccountUserConfigs: []harness.AccountUserConfig{
				{
					Reference: harness.AccountUserProPlayerRef,
				Record: &account_record.AccountUser{
					Email: "test-pro-player@example.com",
				},
				},
			},
		},
		{
			Reference: harness.AccountProDesignerRef,
			AccountUserConfigs: []harness.AccountUserConfig{
				{
					Reference: harness.AccountUserProDesignerRef,
				Record: &account_record.AccountUser{
					Email: "test-pro-designer@example.com",
				},
				},
			},
		},
		{
			Reference: harness.AccountProManagerRef,
			AccountUserConfigs: []harness.AccountUserConfig{
				{
					Reference: harness.AccountUserProManagerRef,
				Record: &account_record.AccountUser{
					Email: "test-pro-manager@example.com",
				},
				},
			},
		},
	}
}

// Once accounts have been created and games have been created the following account
// user game subscription configurations can be created.
func AccountUserGameSubscriptionConfig() []harness.AccountUserGameSubscriptionConfig {

	return []harness.AccountUserGameSubscriptionConfig{
		// Game One
		{
			Reference:                          harness.GameSubscriptionPlayerOneRef,
			GameRef:                            harness.GameOneRef,
			AccountUserRef:                     harness.AccountUserStandardRef,
			SubscriptionType:                   game_record.GameSubscriptionTypePlayer,
			AccountUserManagerGameSubscriptionRef: harness.GameSubscriptionManagerOneRef,
			Record:                             &game_record.GameSubscription{},
		},
		{
			Reference:                          harness.GameSubscriptionPlayerTwoRef,
			GameRef:                            harness.GameOneRef,
			AccountUserRef:                     harness.AccountUserProPlayerRef,
			SubscriptionType:                   game_record.GameSubscriptionTypePlayer,
			AccountUserManagerGameSubscriptionRef: harness.GameSubscriptionManagerOneRef,
			Record:                             &game_record.GameSubscription{},
		},
		{
			Reference:        harness.GameSubscriptionDesignerOneRef,
			GameRef:          harness.GameOneRef,
			AccountUserRef:   harness.AccountUserProDesignerRef,
			SubscriptionType: game_record.GameSubscriptionTypeDesigner,
			Record:           &game_record.GameSubscription{},
		},
		{
			Reference:        harness.GameSubscriptionManagerOneRef,
			GameRef:          harness.GameOneRef,
			AccountUserRef:   harness.AccountUserProManagerRef,
			SubscriptionType: game_record.GameSubscriptionTypeManager,
			Record:           &game_record.GameSubscription{},
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
				},
			},
		},
		// Game Two
		{
			Reference:                          harness.GameSubscriptionPlayerThreeRef,
			GameRef:                            harness.GameTwoRef,
			AccountUserRef:                     harness.AccountUserProPlayerRef,
			SubscriptionType:                   game_record.GameSubscriptionTypePlayer,
			AccountUserManagerGameSubscriptionRef: harness.GameSubscriptionManagerTwoRef,
			Record:                             &game_record.GameSubscription{},
		},
		{
			Reference:        harness.GameSubscriptionManagerTwoRef,
			GameRef:          harness.GameTwoRef,
			AccountUserRef:   harness.AccountUserProManagerRef,
			SubscriptionType: game_record.GameSubscriptionTypeManager,
			Record:           &game_record.GameSubscription{},
			GameInstanceConfigs: []harness.GameInstanceConfig{
				{
					Reference: harness.GameInstanceTwoRef,
					Record: &game_record.GameInstance{
						DeliveryEmail:           true,
						DeliveryPhysicalPost:    true,
						RequiredPlayerCount:     1,
						ProcessWhenAllSubmitted: true,
					},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: "desert-instance-param-lives",
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   domain.AdventureGameParameterCharacterLives,
								ParameterValue: nullstring.FromString("5"),
							},
						},
					},
				},
				// Second instance so parallel E2E tests each get their own isolated game world.
				{
					Reference: "game-instance-desert-b",
					Record: &game_record.GameInstance{
						DeliveryEmail:           true,
						DeliveryPhysicalPost:    true,
						RequiredPlayerCount:     1,
						ProcessWhenAllSubmitted: true,
					},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: "desert-b-instance-param-lives",
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   domain.AdventureGameParameterCharacterLives,
								ParameterValue: nullstring.FromString("5"),
							},
						},
					},
				},
			},
		},
		{
			Reference:        harness.GameSubscriptionDesignerTwoRef,
			GameRef:          harness.GameTwoRef,
			AccountUserRef:   harness.AccountUserProDesignerRef,
			SubscriptionType: game_record.GameSubscriptionTypeDesigner,
			Record:           &game_record.GameSubscription{},
		},
		// Steel Thunder (Mecha)
		{
			Reference:        harness.GameSubscriptionDesignerThreeRef,
			GameRef:          harness.GameThreeRef,
			AccountUserRef:   harness.AccountUserProDesignerRef,
			SubscriptionType: game_record.GameSubscriptionTypeDesigner,
			Record:           &game_record.GameSubscription{},
		},
		{
			Reference:                             harness.GameSubscriptionPlayerFourRef,
			GameRef:                               harness.GameThreeRef,
			AccountUserRef:                        harness.AccountUserStandardRef,
			SubscriptionType:                      game_record.GameSubscriptionTypePlayer,
			AccountUserManagerGameSubscriptionRef: harness.GameSubscriptionManagerThreeRef,
			Record:                                &game_record.GameSubscription{},
		},
		{
			Reference:                             harness.GameSubscriptionPlayerFiveRef,
			GameRef:                               harness.GameThreeRef,
			AccountUserRef:                        harness.AccountUserProPlayerRef,
			SubscriptionType:                      game_record.GameSubscriptionTypePlayer,
			AccountUserManagerGameSubscriptionRef: harness.GameSubscriptionManagerThreeRef,
			Record:                                &game_record.GameSubscription{},
		},
		{
			Reference:        harness.GameSubscriptionManagerThreeRef,
			GameRef:          harness.GameThreeRef,
			AccountUserRef:   harness.AccountUserProManagerRef,
			SubscriptionType: game_record.GameSubscriptionTypeManager,
			Record:           &game_record.GameSubscription{},
			GameInstanceConfigs: []harness.GameInstanceConfig{
				{
					Reference: harness.GameInstanceThreeRef,
					Record: &game_record.GameInstance{
						RequiredPlayerCount:     2,
						ProcessWhenAllSubmitted: true,
					},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: "mech-instance-param-squad-size",
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   domain.MechaParameterSquadSize,
								ParameterValue: nullstring.FromString("2"),
							},
						},
					},
				},
			},
		},
		// Iron Vanguard (Mecha, email only) — single-player scenario with explicit email delivery,
		// mirroring the adventure email-only pattern alongside Steel Thunder's two-player default.
		{
			Reference:        harness.GameSubscriptionDesignerFourRef,
			GameRef:          harness.GameFourRef,
			AccountUserRef:   harness.AccountUserProDesignerRef,
			SubscriptionType: game_record.GameSubscriptionTypeDesigner,
			Record:           &game_record.GameSubscription{},
		},
		{
			Reference:                             harness.GameSubscriptionPlayerSixRef,
			GameRef:                               harness.GameFourRef,
			AccountUserRef:                        harness.AccountUserStandardRef,
			SubscriptionType:                      game_record.GameSubscriptionTypePlayer,
			AccountUserManagerGameSubscriptionRef: harness.GameSubscriptionManagerFourRef,
			Record:                                &game_record.GameSubscription{},
		},
		{
			Reference:        harness.GameSubscriptionManagerFourRef,
			GameRef:          harness.GameFourRef,
			AccountUserRef:   harness.AccountUserProManagerRef,
			SubscriptionType: game_record.GameSubscriptionTypeManager,
			Record:           &game_record.GameSubscription{},
			GameInstanceConfigs: []harness.GameInstanceConfig{
				{
					Reference: harness.GameInstanceFourRef,
					Record: &game_record.GameInstance{
						DeliveryEmail:           true,
						RequiredPlayerCount:     1,
						ProcessWhenAllSubmitted: true,
					},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: "iron-instance-param-squad-size",
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   domain.MechaParameterSquadSize,
								ParameterValue: nullstring.FromString("2"),
							},
						},
					},
				},
			},
		},
	}
}

// GameConfig returns the test data configuration for games
func GameConfig() []harness.GameConfig {
	return []harness.GameConfig{
		{
			Reference: harness.GameOneRef,
			Record: &game_record.Game{
				Name:              "The Enchanted Forest Adventure",
				Description:       "Welcome to The Enchanted Forest Adventure! Step into a world of magic and mystery where ancient forests hold secrets waiting to be discovered. Journey through mystical groves, crystal caverns, and floating islands as you encounter magical creatures, solve puzzles, and uncover the mysteries of this enchanted realm. Your choices matter - every decision shapes your adventure. Whether you're a brave warrior, a cunning mage, or a stealthy scout, this world offers endless possibilities. Join us!",
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 168, // 1 week
			},
			// Game images for turn sheet backgrounds (loaded from test_data_images/)
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
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocationDarkforest},
				},
				{
					Reference: harness.GameLocationTwoRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Crystal Caverns",
						Description: "Deep underground caves filled with glowing crystals. Strange sounds echo from the depths.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocationDungeon},
				},
				{
					Reference: harness.GameLocationThreeRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Floating Islands",
						Description: "Mysterious islands suspended in the sky by unknown magic. Wind howls between them.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocationCliffpath},
				},
				{
					Reference: harness.GameLocationFourRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Shadow Valley",
						Description: "A dark valley shrouded in perpetual shadows. Danger lurks in every corner.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocationDarkforest},
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
					Name:        "The Crystal Path",
					Description: "The crystal-lined path winds back up from the caverns to the mystic grove.",
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
					Name:        "The Sky Vine",
					Description: "The massive enchanted vine provides a path down from the islands to the grove below.",
				},
			},
			{
				Reference:       "location-link-islands-to-caverns",
				FromLocationRef: harness.GameLocationThreeRef,
				ToLocationRef:   harness.GameLocationTwoRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:        "The Wind Lift",
					Description: "The magical wind lift descends from the islands back down into the crystal caverns.",
				},
			},
				// From Shadow Valley
			{
				Reference:       harness.GameLocationLinkFourRef,
				FromLocationRef: harness.GameLocationFourRef,
				ToLocationRef:   harness.GameLocationOneRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:        "The Dark Tunnel",
					Description: "The hidden tunnel beneath ancient roots leads back through the darkness to the mystic grove.",
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
					Name:        "The Deep Descent",
					Description: "The steep path climbs from the shadow valley back up to the crystal caverns.",
				},
			},
			{
				Reference:       "location-link-shadow-to-islands",
				FromLocationRef: harness.GameLocationFourRef,
				ToLocationRef:   harness.GameLocationThreeRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:        "The Shadow Bridge",
					Description: "The bridge of pure darkness stretches from the valley up to the floating islands.",
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
		// Character configuration will be used during the account user game subscription creation
		AdventureGameCharacterConfigs: []harness.AdventureGameCharacterConfig{
			{
				Reference:  harness.GameCharacterOneRef,
				AccountRef: harness.AccountUserStandardRef,
				Record: &adventure_game_record.AdventureGameCharacter{
					Name: "Aria the Mage",
				},
			},
			{
				Reference:  harness.GameCharacterTwoRef,
				AccountRef: harness.AccountUserProPlayerRef,
				Record: &adventure_game_record.AdventureGameCharacter{
					Name: "Thorne the Warrior",
				},
			},
			{
				Reference:  harness.GameCharacterThreeRef,
				AccountRef: harness.AccountUserProDesignerRef,
				Record: &adventure_game_record.AdventureGameCharacter{
					Name: "Luna the Scout",
				},
			},
			{
				Reference:  harness.GameCharacterFourRef,
				AccountRef: harness.AccountUserProManagerRef,
				Record: &adventure_game_record.AdventureGameCharacter{
					Name: "Max the Manager",
				},
			},
		},
		AdventureGameCreaturePlacementConfigs: []harness.AdventureGameCreaturePlacementConfig{
			{
				Reference:       harness.GameCreaturePlacementOneRef,
				GameCreatureRef: harness.GameCreatureOneRef,
				GameLocationRef: harness.GameLocationOneRef,
				InitialCount:    1,
				Record:          &adventure_game_record.AdventureGameCreaturePlacement{},
			},
		},
		AdventureGameItemPlacementConfigs: []harness.AdventureGameItemPlacementConfig{
			{
				Reference:       harness.GameItemPlacementOneRef,
				GameItemRef:     harness.GameItemOneRef,
				GameLocationRef: harness.GameLocationOneRef,
				InitialCount:    1,
				Record:          &adventure_game_record.AdventureGameItemPlacement{},
			},
		},
	},
	{
		Reference: harness.GameTwoRef,
			Record: &game_record.Game{
				Name:              "The Desert Kingdom",
				Description:       "Welcome to The Desert Kingdom! Embark on a solo quest across scorching sand dunes, crumbling ancient ruins, and hidden oases in this single-player email adventure. Survive the elements, gather legendary artefacts, and unlock the sealed passages of a forgotten civilisation. Every turn challenges you to choose your path carefully — equip the right gear, explore boldly, and uncover the secrets buried beneath the shifting sands. Your destiny awaits!",
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 1,
			},
			GameImageConfigs: []harness.GameImageConfig{
				{
					Reference:     "desert-image-join-game",
					ImagePath:     ImageDesertJoinGame,
					TurnSheetType: adventure_game_record.AdventureGameTurnSheetTypeJoinGame,
				},
				{
					Reference:     "desert-image-inventory",
					ImagePath:     ImageDesertInventory,
					TurnSheetType: adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement,
				},
			},
			// 4 locations forming a rich, interconnected desert world
			AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
				{
					Reference: "desert-location-oasis",
					Record: &adventure_game_record.AdventureGameLocation{
						Name:               "Oasis Village",
						Description:        "A bustling village built around a life-giving oasis. Palm trees sway over turquoise pools while merchants hawk exotic wares under colourful canopies. Travellers rest here before venturing into the unforgiving desert beyond.",
						IsStartingLocation: true,
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageDesertOasis},
				},
				{
					Reference: "desert-location-ruins",
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Ancient Ruins",
						Description: "Crumbling columns and weathered arches rise from the sand, remnants of a once-great civilisation. Hieroglyphs cover every surface, and the air is thick with the dust of ages. Something valuable — and dangerous — surely lies deeper within.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageDesertRuins},
				},
				{
					Reference: "desert-location-canyon",
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Sandstone Canyon",
						Description: "Towering walls of red and gold sandstone hem in a narrow path that winds ever deeper. Shafts of sunlight pierce the gloom, and the wind sings eerie melodies through natural arches overhead. Creatures lurk in the shadowed crevices.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageDesertCanyon},
				},
				{
					Reference: "desert-location-temple",
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Hidden Temple",
						Description: "Half-swallowed by the dunes, an ancient temple glows with an otherworldly light. Serpent carvings frame the entrance, and golden runes pulse with forgotten magic. Only those who carry the right tokens may pass the sealed doors within.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageDesertTemple},
				},
			},
			// 4 items exercising all inventory actions (pick up, drop, equip, unequip)
			AdventureGameItemConfigs: []harness.AdventureGameItemConfig{
				{
					Reference: "desert-item-compass",
					Record: &adventure_game_record.AdventureGameItem{
						Name:          "Desert Compass",
						Description:   "A gleaming brass compass inlaid with sapphire. Its needle always points toward the nearest water source — an invaluable tool in the wastes.",
						CanBeEquipped: true,
						EquipmentSlot: convert.PtrStrict("jewelry"),
					},
				},
			{
				Reference: "desert-item-flask",
				Record: &adventure_game_record.AdventureGameItem{
					Name:        "Water Flask",
					Description: "A sturdy leather flask filled with cool, clear water from the oasis. Essential for desert survival, but heavy to carry.",
				},
				AdventureGameItemEffectConfigs: []harness.AdventureGameItemEffectConfig{
					{
						Reference: "desert-item-flask-effect-use-heal",
						Record: &adventure_game_record.AdventureGameItemEffect{
							ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeUse,
							ResultDescription: "You drink deeply from the cool water, feeling refreshed and restored.",
							EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeHealWielder,
							ResultValueMin:    nullint32.FromInt32(5),
							ResultValueMax:    nullint32.FromInt32(10),
							IsRepeatable:      true,
						},
					},
				},
			},
				{
					Reference: "desert-item-cloak",
					Record: &adventure_game_record.AdventureGameItem{
						Name:          "Sand Cloak",
						Description:   "A shimmering cloak woven from enchanted desert silk. It bends light around the wearer, granting near-invisibility in sandy terrain.",
						CanBeEquipped: true,
						EquipmentSlot: convert.PtrStrict("clothing"),
					},
				},
				{
					Reference: "desert-item-scarab-key",
					Record: &adventure_game_record.AdventureGameItem{
						Name:          "Ancient Scarab Key",
						Description:   "A golden scarab amulet etched with temple hieroglyphs. It hums with power when brought near the sealed passages of the Hidden Temple.",
						CanBeEquipped: true,
						EquipmentSlot: convert.PtrStrict("jewelry"),
					},
				},
			},
			// 2 creatures placed in the world
			AdventureGameCreatureConfigs: []harness.AdventureGameCreatureConfig{
				{
					Reference: "desert-creature-serpent",
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Sand Serpent",
						Description: "A massive serpent that burrows through the desert sands, ambushing unwary travellers with lightning speed.",
					},
				},
				{
					Reference: "desert-creature-guardian",
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Temple Guardian",
						Description: "An ancient stone golem animated by temple magic. It guards the inner sanctum and will challenge any who enter.",
					},
				},
			},
		// Placement configs define where creatures/items start in each game instance
		AdventureGameCreaturePlacementConfigs: []harness.AdventureGameCreaturePlacementConfig{
			{
				Reference:       "desert-creature-placement-serpent",
				GameCreatureRef: "desert-creature-serpent",
				GameLocationRef: "desert-location-canyon",
				InitialCount:    1,
				Record:          &adventure_game_record.AdventureGameCreaturePlacement{},
			},
			{
				Reference:       "desert-creature-placement-guardian",
				GameCreatureRef: "desert-creature-guardian",
				GameLocationRef: "desert-location-temple",
				InitialCount:    1,
				Record:          &adventure_game_record.AdventureGameCreaturePlacement{},
			},
		},
		AdventureGameItemPlacementConfigs: []harness.AdventureGameItemPlacementConfig{
			{
				Reference:       "desert-item-placement-compass",
				GameItemRef:     "desert-item-compass",
				GameLocationRef: "desert-location-oasis",
				InitialCount:    1,
				Record:          &adventure_game_record.AdventureGameItemPlacement{},
			},
		{
			Reference:       "desert-item-placement-flask",
			GameItemRef:     "desert-item-flask",
			GameLocationRef: "desert-location-oasis",
			InitialCount:    1,
			Record:          &adventure_game_record.AdventureGameItemPlacement{},
		},
		{
			Reference:       "desert-item-placement-cloak",
			GameItemRef:     "desert-item-cloak",
			GameLocationRef: "desert-location-ruins",
			InitialCount:    1,
			Record:          &adventure_game_record.AdventureGameItemPlacement{},
		},
	},
		// Location links — bidirectional, some with item requirements
		AdventureGameLocationLinkConfigs: []harness.AdventureGameLocationLinkConfig{
			// Oasis Village <-> Ancient Ruins (free)
			{
				Reference:       "desert-link-oasis-to-ruins",
					FromLocationRef: "desert-location-oasis",
					ToLocationRef:   "desert-location-ruins",
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Dusty Trail",
						Description: "A well-worn trail leads east through the dunes toward the ancient ruins.",
					},
				},
			{
				Reference:       "desert-link-ruins-to-oasis",
				FromLocationRef: "desert-location-ruins",
				ToLocationRef:   "desert-location-oasis",
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:        "The Dusty Trail",
					Description: "A well-worn trail leads west through the dunes back to the oasis village.",
				},
			},
				// Oasis Village <-> Sandstone Canyon (free)
			{
				Reference:       "desert-link-oasis-to-canyon",
				FromLocationRef: "desert-location-oasis",
				ToLocationRef:   "desert-location-canyon",
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:        "The Rope Climb",
					Description: "Knotted ropes and iron spikes mark the way down into the sandstone canyon.",
				},
			},
				{
					Reference:       "desert-link-canyon-to-oasis",
					FromLocationRef: "desert-location-canyon",
					ToLocationRef:   "desert-location-oasis",
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Rope Climb",
						Description: "Knotted ropes and iron spikes mark the climb back up to the oasis.",
					},
				},
			// Ancient Ruins -> Hidden Temple
			// traverse: Scarab Key must be equipped (item equipped condition)
			// traverse: Temple Guardian must not be alive at ruins (creature none_alive_at_location condition) -- combined AND
			{
				Reference:       "desert-link-ruins-to-temple",
				FromLocationRef: "desert-location-ruins",
				ToLocationRef:   "desert-location-temple",
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:              "The Sealed Passage",
					Description:       "A massive stone door blocks the way. Scarab-shaped indentations line its surface.",
					LockedDescription: nullstring.FromString("The stone door is sealed fast. Scarab-shaped indentations line its surface, but without the right artefact it will not yield."),
				},
				AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
					{
						Reference:   "desert-link-req-scarab-equipped",
						GameItemRef: "desert-item-scarab-key",
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionEquipped,
							Quantity:  1,
						},
					},
					{
						Reference:       "desert-link-req-guardian-dead",
						GameCreatureRef: "desert-creature-guardian",
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveAtLocation,
							Quantity:  1,
						},
					},
				},
			},
		// Hidden Temple -> Ancient Ruins (free)
		{
			Reference:       "desert-link-temple-to-ruins",
			FromLocationRef: "desert-location-temple",
			ToLocationRef:   "desert-location-ruins",
			Record: &adventure_game_record.AdventureGameLocationLink{
				Name:        "The Sealed Passage",
				Description: "The great stone door opens from this side, revealing a passage back to the ruins.",
			},
		},
			// Sandstone Canyon -> Hidden Temple
			// visible: Sand Serpent must be dead at the canyon (creature dead_at_location visibility condition)
			// traverse: Sand Cloak must be in inventory (item in_inventory traverse condition)
			// -- combined: link hidden while serpent lives; locked (without cloak) once serpent is dead
			{
			Reference:       "desert-link-canyon-to-temple",
			FromLocationRef: "desert-location-canyon",
			ToLocationRef:   "desert-location-temple",
			Record: &adventure_game_record.AdventureGameLocationLink{
				Name:              "The Wind Tunnel",
				Description:       "A narrow tunnel carved by wind connects the canyon to the hidden temple. Those cloaked in sand magic can slip through.",
				LockedDescription: nullstring.FromString("A faint outline of a passage shimmers in the canyon wall. Only one wrapped in desert silk could slip through."),
			},
				AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
					{
						Reference:       "desert-link-req-serpent-dead",
						GameCreatureRef: "desert-creature-serpent",
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeVisible,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionDeadAtLocation,
							Quantity:  1,
						},
					},
					{
						Reference:   "desert-link-req-cloak-inventory",
						GameItemRef: "desert-item-cloak",
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory,
							Quantity:  1,
						},
					},
				},
			},
		// Hidden Temple -> Sandstone Canyon (free)
		{
			Reference:       "desert-link-temple-to-canyon",
			FromLocationRef: "desert-location-temple",
			ToLocationRef:   "desert-location-canyon",
			Record: &adventure_game_record.AdventureGameLocationLink{
				Name:        "The Wind Tunnel",
				Description: "A blast of dry wind funnels through a narrow tunnel back to the canyon.",
			},
		},
		},
		AdventureGameLocationObjectConfigs: []harness.AdventureGameLocationObjectConfig{
			{
				Reference:       "desert-obj-sarcophagus",
				LocationRef:     "desert-location-ruins",
				InitialStateRef: "desert-obj-sarcophagus-state-sealed",
				Record: &adventure_game_record.AdventureGameLocationObject{
					Name:        "Sealed Sarcophagus",
					Description: "A massive stone sarcophagus covered in hieroglyphs. The lid appears to be movable.",
					IsHidden:    false,
				},
				AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
					{
						Reference: "desert-obj-sarcophagus-state-sealed",
						Record: &adventure_game_record.AdventureGameLocationObjectState{
							Name:        "sealed",
							Description: "The sarcophagus is sealed shut.",
							SortOrder:   0,
						},
					},
					{
						Reference: "desert-obj-sarcophagus-state-open",
						Record: &adventure_game_record.AdventureGameLocationObjectState{
							Name:        "open",
							Description: "The sarcophagus lid has been removed.",
							SortOrder:   1,
						},
					},
				},
				AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
					{
						Reference: "desert-obj-effect-sarcophagus-inspect",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							ResultDescription: "A stone sarcophagus covered in hieroglyphs. It looks like it can be opened.",
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							IsRepeatable:      true,
						},
					},
				{
					Reference:         "desert-obj-effect-sarcophagus-open",
					RequiredStateRef:  "desert-obj-sarcophagus-state-sealed",
					ResultItemRef:     "desert-item-scarab-key",
					ResultLocationRef: "desert-location-ruins",
					Record: &adventure_game_record.AdventureGameLocationObjectEffect{
						ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
						ResultDescription: "The lid grinds open. A golden scarab key rests on dusty linen inside.",
						EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypePlaceItem,
						IsRepeatable:      false,
					},
				},
					{
						Reference:        "desert-obj-effect-sarcophagus-open-state",
						RequiredStateRef: "desert-obj-sarcophagus-state-sealed",
						ResultStateRef:   "desert-obj-sarcophagus-state-open",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
							ResultDescription: "",
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							IsRepeatable:      false,
						},
					},
					{
						Reference:        "desert-obj-effect-sarcophagus-search",
						RequiredStateRef: "desert-obj-sarcophagus-state-open",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
							ResultDescription: "Nothing remains inside but dust and scraps of ancient linen.",
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							IsRepeatable:      true,
						},
					},
				},
			},
			{
				Reference:       "desert-obj-obelisk",
				LocationRef:     "desert-location-oasis",
				InitialStateRef: "desert-obj-obelisk-state-intact",
				Record: &adventure_game_record.AdventureGameLocationObject{
					Name:        "Sand-Carved Obelisk",
					Description: "A towering obelisk carved from desert sandstone, covered in ancient symbols.",
					IsHidden:    false,
				},
				AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
					{
						Reference: "desert-obj-obelisk-state-intact",
						Record: &adventure_game_record.AdventureGameLocationObjectState{
							Name:        "intact",
							Description: "The obelisk stands as it has for millennia.",
							SortOrder:   0,
						},
					},
					{
						Reference: "desert-obj-obelisk-state-glowing",
						Record: &adventure_game_record.AdventureGameLocationObjectState{
							Name:        "glowing",
							Description: "Ancient symbols glow with warm golden light.",
							SortOrder:   1,
						},
					},
				},
				AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
					{
						Reference: "desert-obj-effect-obelisk-inspect",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							ResultDescription: "Ancient symbols cover the obelisk. They seem to describe a path through the canyon.",
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							IsRepeatable:      true,
						},
					},
					{
						Reference:        "desert-obj-effect-obelisk-touch",
						RequiredStateRef: "desert-obj-obelisk-state-intact",
						ResultStateRef:   "desert-obj-obelisk-state-glowing",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeTouch,
							ResultDescription: "The symbols begin to glow with a warm golden light.",
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							IsRepeatable:      false,
						},
					},
			{
				Reference:        "desert-obj-effect-obelisk-read",
				RequiredStateRef: "desert-obj-obelisk-state-glowing",
				Record: &adventure_game_record.AdventureGameLocationObjectEffect{
					ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeRead,
					ResultDescription: "As you trace the glowing symbols a surge of ancient energy courses through you, burning your hand and leaving a searing mark.",
					EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeDamage,
					ResultValueMin:    nullint32.FromInt32(10),
					ResultValueMax:    nullint32.FromInt32(20),
					IsRepeatable:      false,
				},
			},
		},
			},
		},
	},
	{
		Reference: harness.GameThreeRef,
		Record: &game_record.Game{
			Name:              "Steel Thunder",
			Description:       "Welcome to Steel Thunder! Command a squad of powerful war mechs across contested industrial sectors in this hard-hitting tactical wargame. Scout enemy positions, manage heat buildup, and coordinate fire to destroy opposing squads before they reach your stronghold. Whether you pilot a nimble light mech or an armoured assault chassis, every battle decision matters. Gear up and take the field!",
			GameType:          game_record.GameTypeMecha,
			TurnDurationHours: 168, // 1 week
		},
		GameImageConfigs: []harness.GameImageConfig{
			{
				Reference:     "steel-image-join-game",
				ImagePath:     ImageSteelJoinGame,
				TurnSheetType: mecha_record.MechaTurnSheetTypeJoinGame,
			},
			{
				Reference:     "steel-image-management",
				ImagePath:     ImageSteelOrders,
				TurnSheetType: mecha_record.MechaTurnSheetTypeSquadManagement,
			},
			{
				Reference:     "steel-image-orders",
				ImagePath:     ImageSteelOrders,
				TurnSheetType: mecha_record.MechaTurnSheetTypeOrders,
			},
		},
		MechaChassisConfigs: []harness.MechaChassisConfig{
			{
				Reference: "steel-chassis-viper",
				Record: &mecha_record.MechaChassis{
					Name:            "Viper",
					Description:     "A fast light recon mech favoured for scouting and harassment.",
					ChassisClass:    mecha_record.ChassisClassLight,
					ArmorPoints:     80,
					StructurePoints: 40,
					HeatCapacity:    20,
					Speed:           6,
				},
			},
			{
				Reference: "steel-chassis-ranger",
				Record: &mecha_record.MechaChassis{
					Name:            "Ranger",
					Description:     "A versatile medium mech with solid armour and good mobility.",
					ChassisClass:    mecha_record.ChassisClassMedium,
					ArmorPoints:     140,
					StructurePoints: 70,
					HeatCapacity:    30,
					Speed:           4,
				},
			},
			{
				Reference: "steel-chassis-crusher",
				Record: &mecha_record.MechaChassis{
					Name:            "Crusher",
					Description:     "A feared heavy mech with devastating long-range firepower.",
					ChassisClass:    mecha_record.ChassisClassHeavy,
					ArmorPoints:     200,
					StructurePoints: 100,
					HeatCapacity:    40,
					Speed:           3,
				},
			},
		},
		MechaWeaponConfigs: []harness.MechaWeaponConfig{
			{
				Reference: "steel-weapon-pulse-cannon",
				Record: &mecha_record.MechaWeapon{
					Name:        "Pulse Cannon",
					Description: "A reliable direct-fire energy weapon with moderate range.",
					Damage:      5,
					HeatCost:    3,
					RangeBand:   mecha_record.WeaponRangeBandMedium,
					MountSize:   mecha_record.WeaponMountSizeMedium,
				},
			},
			{
				Reference: "steel-weapon-heavy-pulse-cannon",
				Record: &mecha_record.MechaWeapon{
					Name:        "Heavy Pulse Cannon",
					Description: "A powerful energy weapon that generates significant heat.",
					Damage:      8,
					HeatCost:    8,
					RangeBand:   mecha_record.WeaponRangeBandLong,
					MountSize:   mecha_record.WeaponMountSizeLarge,
				},
			},
			{
				Reference: "steel-weapon-rocket-pack",
				Record: &mecha_record.MechaWeapon{
					Name:        "Rocket Pack",
					Description: "A short-range missile rack ideal for close-in brawling.",
					Damage:      8,
					HeatCost:    3,
					RangeBand:   mecha_record.WeaponRangeBandShort,
					MountSize:   mecha_record.WeaponMountSizeMedium,
				},
			},
			{
				Reference: "steel-weapon-light-pulse-cannon",
				Record: &mecha_record.MechaWeapon{
					Name:        "Light Pulse Cannon",
					Description: "A light back-up weapon for point-blank defence.",
					Damage:      3,
					HeatCost:    1,
					RangeBand:   mecha_record.WeaponRangeBandShort,
					MountSize:   mecha_record.WeaponMountSizeSmall,
				},
			},
		},
		MechaSectorConfigs: []harness.MechaSectorConfig{
			{
				Reference: "steel-sector-deployment-zone",
				Record: &mecha_record.MechaSector{
					Name:             "Deployment Zone",
					Description:      "An open staging area where squads muster before the battle.",
					TerrainType:      mecha_record.SectorTerrainTypeOpen,
					Elevation:        0,
					IsStartingSector: true,
				},
			},
			{
				Reference: "steel-sector-refinery",
				Record: &mecha_record.MechaSector{
					Name:        "Refinery Complex",
					Description: "A dense industrial site offering plentiful cover but restricted movement.",
					TerrainType: mecha_record.SectorTerrainTypeUrban,
					Elevation:   0,
				},
			},
			{
				Reference: "steel-sector-ridge",
				Record: &mecha_record.MechaSector{
					Name:        "Crater Ridge",
					Description: "A rocky elevated ridge commanding views of the surrounding lowlands.",
					TerrainType: mecha_record.SectorTerrainTypeRough,
					Elevation:   3,
				},
			},
			{
				Reference: "steel-sector-river-crossing",
				Record: &mecha_record.MechaSector{
					Name:        "River Crossing",
					Description: "A shallow ford crossing that slows heavy mechs significantly.",
					TerrainType: mecha_record.SectorTerrainTypeWater,
					Elevation:   -1,
				},
			},
		},
		MechaSectorLinkConfigs: []harness.MechaSectorLinkConfig{
			{
				Reference:     "steel-link-deploy-to-refinery",
				FromSectorRef: "steel-sector-deployment-zone",
				ToSectorRef:   "steel-sector-refinery",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "steel-link-refinery-to-deploy",
				FromSectorRef: "steel-sector-refinery",
				ToSectorRef:   "steel-sector-deployment-zone",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "steel-link-deploy-to-ridge",
				FromSectorRef: "steel-sector-deployment-zone",
				ToSectorRef:   "steel-sector-ridge",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "steel-link-ridge-to-deploy",
				FromSectorRef: "steel-sector-ridge",
				ToSectorRef:   "steel-sector-deployment-zone",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "steel-link-refinery-to-river",
				FromSectorRef: "steel-sector-refinery",
				ToSectorRef:   "steel-sector-river-crossing",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "steel-link-river-to-refinery",
				FromSectorRef: "steel-sector-river-crossing",
				ToSectorRef:   "steel-sector-refinery",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "steel-link-ridge-to-river",
				FromSectorRef: "steel-sector-ridge",
				ToSectorRef:   "steel-sector-river-crossing",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "steel-link-river-to-ridge",
				FromSectorRef: "steel-sector-river-crossing",
				ToSectorRef:   "steel-sector-ridge",
				Record:        &mecha_record.MechaSectorLink{},
			},
		},
		MechaSquadConfigs: []harness.MechaSquadConfig{
			{
				Reference: "steel-squad-starter",
				SquadType: mecha_record.SquadTypeStarter,
				Record: &mecha_record.MechaSquad{
					Name:        "Thunder Squad",
					Description: "Standard issue starter squad for incoming Steel Thunder commanders.",
				},
				SquadMechConfigs: []harness.MechaSquadMechConfig{
					{
						Reference:  "steel-mech-starter-1",
						ChassisRef: "steel-chassis-viper",
						Record:     &mecha_record.MechaSquadMech{Callsign: "Thunder-1"},
						WeaponConfigRefs: []harness.MechaSquadMechWeaponRef{
							{WeaponRef: "steel-weapon-light-pulse-cannon", SlotLocation: "right-arm"},
						},
					},
					{
						Reference:  "steel-mech-starter-2",
						ChassisRef: "steel-chassis-ranger",
						Record:     &mecha_record.MechaSquadMech{Callsign: "Thunder-2"},
						WeaponConfigRefs: []harness.MechaSquadMechWeaponRef{
							{WeaponRef: "steel-weapon-pulse-cannon", SlotLocation: "right-torso"},
							{WeaponRef: "steel-weapon-rocket-pack", SlotLocation: "left-torso"},
						},
					},
				},
			},
			{
				Reference: "steel-squad-opponent",
				SquadType: mecha_record.SquadTypeOpponent,
				Record: &mecha_record.MechaSquad{
					Name:        "Alpha Squad",
					Description: "An opponent squad template for Steel Thunder.",
				},
				SquadMechConfigs: []harness.MechaSquadMechConfig{
					{
						Reference:  "steel-mech-opponent-1",
						ChassisRef: "steel-chassis-ranger",
						Record:     &mecha_record.MechaSquadMech{Callsign: "Opponent-1"},
						WeaponConfigRefs: []harness.MechaSquadMechWeaponRef{
							{WeaponRef: "steel-weapon-pulse-cannon", SlotLocation: "right-torso"},
							{WeaponRef: "steel-weapon-rocket-pack", SlotLocation: "left-torso"},
						},
					},
					{
						Reference:  "steel-mech-opponent-2",
						ChassisRef: "steel-chassis-viper",
						Record:     &mecha_record.MechaSquadMech{Callsign: "Opponent-2"},
						WeaponConfigRefs: []harness.MechaSquadMechWeaponRef{
							{WeaponRef: "steel-weapon-light-pulse-cannon", SlotLocation: "right-arm"},
						},
					},
				},
			},
		},
	},
	// Mecha email-only scenario — mirrors the adventure email-only pattern with a single
	// player and explicit DeliveryEmail delivery. Distinct from Steel Thunder (GameThree)
	// which uses two players and no explicit delivery setting.
	{
		Reference: harness.GameFourRef,
		Record: &game_record.Game{
			Name:              "Iron Vanguard",
			Description:       "Welcome to Iron Vanguard! Pilot a squad of war mechs through three contested sectors in this focused single-player tactical wargame. Manage heat, exploit cover, and push through enemy lines to seize the forward command post. Fast turns, decisive choices — engage!",
			GameType:          game_record.GameTypeMecha,
			TurnDurationHours: 168, // 1 week
		},
		GameImageConfigs: []harness.GameImageConfig{
			{
				Reference:     "iron-image-join-game",
				ImagePath:     ImageIronJoinGame,
				TurnSheetType: mecha_record.MechaTurnSheetTypeJoinGame,
			},
			{
				Reference:     "iron-image-management",
				ImagePath:     ImageIronOrders,
				TurnSheetType: mecha_record.MechaTurnSheetTypeSquadManagement,
			},
			{
				Reference:     "iron-image-orders",
				ImagePath:     ImageIronOrders,
				TurnSheetType: mecha_record.MechaTurnSheetTypeOrders,
			},
		},
		MechaChassisConfigs: []harness.MechaChassisConfig{
			{
				Reference: "iron-chassis-scout",
				Record: &mecha_record.MechaChassis{
					Name:            "Scout",
					Description:     "A nimble light mech optimised for rapid advances and flanking manoeuvres.",
					ChassisClass:    mecha_record.ChassisClassLight,
					ArmorPoints:     72,
					StructurePoints: 32,
					HeatCapacity:    18,
					Speed:           7,
				},
			},
			{
				Reference: "iron-chassis-sentinel",
				Record: &mecha_record.MechaChassis{
					Name:            "Sentinel",
					Description:     "A dependable medium mech that holds ground while laying down sustained fire.",
					ChassisClass:    mecha_record.ChassisClassMedium,
					ArmorPoints:     130,
					StructurePoints: 65,
					HeatCapacity:    28,
					Speed:           4,
				},
			},
		},
		MechaWeaponConfigs: []harness.MechaWeaponConfig{
			{
				Reference: "iron-weapon-pulse-cannon",
				Record: &mecha_record.MechaWeapon{
					Name:        "Pulse Cannon",
					Description: "A reliable medium-range direct-fire energy weapon.",
					Damage:      5,
					HeatCost:    3,
					RangeBand:   mecha_record.WeaponRangeBandMedium,
					MountSize:   mecha_record.WeaponMountSizeMedium,
				},
			},
			{
				Reference: "iron-weapon-chaingun",
				Record: &mecha_record.MechaWeapon{
					Name:        "Chaingun",
					Description: "A rapid-fire ballistic weapon effective against light armour at close range.",
					Damage:      2,
					HeatCost:    0,
					RangeBand:   mecha_record.WeaponRangeBandShort,
					MountSize:   mecha_record.WeaponMountSizeSmall,
				},
			},
		},
		MechaSectorConfigs: []harness.MechaSectorConfig{
			{
				Reference: "iron-sector-staging-area",
				Record: &mecha_record.MechaSector{
					Name:             "Staging Area",
					Description:      "A flat open zone where the squad assembles before pushing into contested ground.",
					TerrainType:      mecha_record.SectorTerrainTypeOpen,
					Elevation:        0,
					IsStartingSector: true,
				},
			},
			{
				Reference: "iron-sector-industrial-district",
				Record: &mecha_record.MechaSector{
					Name:        "Industrial District",
					Description: "A sprawling complex of warehouses and storage tanks offering excellent cover.",
					TerrainType: mecha_record.SectorTerrainTypeUrban,
					Elevation:   0,
				},
			},
			{
				Reference: "iron-sector-command-post",
				Record: &mecha_record.MechaSector{
					Name:        "Forward Command Post",
					Description: "A fortified command bunker on elevated ground — the primary objective.",
					TerrainType: mecha_record.SectorTerrainTypeUrban,
					Elevation:   2,
				},
			},
		},
		MechaSectorLinkConfigs: []harness.MechaSectorLinkConfig{
			{
				Reference:     "iron-link-staging-to-industrial",
				FromSectorRef: "iron-sector-staging-area",
				ToSectorRef:   "iron-sector-industrial-district",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "iron-link-industrial-to-staging",
				FromSectorRef: "iron-sector-industrial-district",
				ToSectorRef:   "iron-sector-staging-area",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "iron-link-industrial-to-command",
				FromSectorRef: "iron-sector-industrial-district",
				ToSectorRef:   "iron-sector-command-post",
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     "iron-link-command-to-industrial",
				FromSectorRef: "iron-sector-command-post",
				ToSectorRef:   "iron-sector-industrial-district",
				Record:        &mecha_record.MechaSectorLink{},
			},
		},
		MechaComputerOpponentConfigs: []harness.MechaComputerOpponentConfig{
			{
				Reference: "iron-computer-opponent",
				Record: &mecha_record.MechaComputerOpponent{
					Name:        "Iron Vanguard Defenders",
					Description: "The entrenched garrison holding the Forward Command Post. Aggressive defenders with solid tactical IQ.",
					Aggression:  7,
					IQ:          5,
				},
			},
		},
		MechaSquadConfigs: []harness.MechaSquadConfig{
			{
				Reference: "iron-squad-starter",
				SquadType: mecha_record.SquadTypeStarter,
				Record: &mecha_record.MechaSquad{
					Name:        "Vanguard Squad",
					Description: "Standard issue starter squad for incoming Iron Vanguard commanders.",
				},
				SquadMechConfigs: []harness.MechaSquadMechConfig{
					{
						Reference:  "iron-mech-starter-1",
						ChassisRef: "iron-chassis-scout",
						Record:     &mecha_record.MechaSquadMech{Callsign: "Vanguard-1"},
						WeaponConfigRefs: []harness.MechaSquadMechWeaponRef{
							{WeaponRef: "iron-weapon-chaingun", SlotLocation: "right-arm"},
						},
					},
					{
						Reference:  "iron-mech-starter-2",
						ChassisRef: "iron-chassis-sentinel",
						Record:     &mecha_record.MechaSquadMech{Callsign: "Vanguard-2"},
						WeaponConfigRefs: []harness.MechaSquadMechWeaponRef{
							{WeaponRef: "iron-weapon-pulse-cannon", SlotLocation: "right-torso"},
							{WeaponRef: "iron-weapon-chaingun", SlotLocation: "left-arm"},
						},
					},
				},
			},
			{
				Reference: "iron-computer-opponent-squad",
				SquadType: mecha_record.SquadTypeOpponent,
				Record: &mecha_record.MechaSquad{
					Name:        "Defender Squad",
					Description: "The garrison's reaction force, defending the Command Post.",
				},
				SquadMechConfigs: []harness.MechaSquadMechConfig{
					{
						Reference:  "iron-defender-1",
						ChassisRef: "iron-chassis-sentinel",
						Record:     &mecha_record.MechaSquadMech{Callsign: "Defender-1"},
						WeaponConfigRefs: []harness.MechaSquadMechWeaponRef{
							{WeaponRef: "iron-weapon-pulse-cannon", SlotLocation: "right-torso"},
							{WeaponRef: "iron-weapon-chaingun", SlotLocation: "left-arm"},
						},
					},
					{
						Reference:  "iron-defender-2",
						ChassisRef: "iron-chassis-scout",
						Record:     &mecha_record.MechaSquadMech{Callsign: "Defender-2"},
						WeaponConfigRefs: []harness.MechaSquadMechWeaponRef{
							{WeaponRef: "iron-weapon-chaingun", SlotLocation: "right-arm"},
						},
					},
				},
			},
		},
	},
}
}
