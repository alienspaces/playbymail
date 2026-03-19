package test_data

import (
	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
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
			{
				Reference:       "desert-item-placement-scarab",
				GameItemRef:     "desert-item-scarab-key",
				GameLocationRef: "desert-location-canyon",
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
						Name:        "The Trade Road",
						Description: "The main road west winds back to the oasis village.",
					},
				},
				// Oasis Village <-> Sandstone Canyon (free)
				{
					Reference:       "desert-link-oasis-to-canyon",
					FromLocationRef: "desert-location-oasis",
					ToLocationRef:   "desert-location-canyon",
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Canyon Descent",
						Description: "A narrow switchback path descends south into the sandstone canyon.",
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
					Name:        "The Crumbling Steps",
					Description: "Worn stone steps lead back out to the ruins.",
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
					Name:              "The Shadow Path",
					Description:       "A hidden passage through the canyon wall, visible only to those cloaked in sand magic.",
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
		},
	}
}
