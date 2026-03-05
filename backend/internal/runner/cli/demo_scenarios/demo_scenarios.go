package demo_scenarios

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// Image file names in demo_scenario_images directory
const (
	ImageJoinGame            = "join-game.png"
	ImageInventoryManagement = "inventory-management.png"
	ImageLocationDarkforest  = "location-darkforest.png"
	ImageLocationDungeon     = "location-dungeon.png"
	ImageLocationCliffpath   = "location-cliffpath.png"
)

// Demo-specific harness references (prefixed "demo-" to avoid collisions with test data)
const (
	DemoGameOneRef   = "demo-game-one"
	DemoGameTwoRef   = "demo-game-two"

	DemoAccountDesignerRef = "demo-account-designer"
	DemoAccountManagerRef  = "demo-account-manager"
	DemoAccountPlayerRef   = "demo-account-player"

	DemoSubscriptionDesignerOneRef = "demo-subscription-designer-one"
	DemoSubscriptionDesignerTwoRef = "demo-subscription-designer-two"
	DemoSubscriptionManagerOneRef  = "demo-subscription-manager-one"
	DemoSubscriptionManagerTwoRef  = "demo-subscription-manager-two"
	DemoSubscriptionPlayerOneRef   = "demo-subscription-player-one"
	DemoSubscriptionPlayerTwoRef   = "demo-subscription-player-two"

	DemoLocationOneRef   = "demo-location-one"
	DemoLocationTwoRef   = "demo-location-two"
	DemoLocationThreeRef = "demo-location-three"
	DemoLocationFourRef  = "demo-location-four"

	DemoLocationLinkOneRef   = "demo-location-link-one"
	DemoLocationLinkTwoRef   = "demo-location-link-two"
	DemoLocationLinkThreeRef = "demo-location-link-three"
	DemoLocationLinkFourRef  = "demo-location-link-four"

	DemoLinkReqOneRef   = "demo-link-req-one"
	DemoLinkReqTwoRef   = "demo-link-req-two"
	DemoLinkReqThreeRef = "demo-link-req-three"
	DemoLinkReqFourRef  = "demo-link-req-four"

	DemoItemOneRef   = "demo-item-one"
	DemoItemTwoRef   = "demo-item-two"
	DemoItemThreeRef = "demo-item-three"
	DemoItemFourRef  = "demo-item-four"

	DemoCreatureOneRef = "demo-creature-one"
	DemoCreatureTwoRef = "demo-creature-two"

	DemoCharacterOneRef   = "demo-character-one"
	DemoCharacterTwoRef   = "demo-character-two"
	DemoCharacterThreeRef = "demo-character-three"

	DemoInstanceOneRef = "demo-instance-one"
	DemoInstanceTwoRef = "demo-instance-two"

	DemoInstanceParamOneRef = "demo-instance-param-one"
	DemoInstanceParamTwoRef = "demo-instance-param-two"

	DemoLocationInstanceOneRef = "demo-location-instance-one"
	DemoLocationInstanceTwoRef = "demo-location-instance-two"

	DemoItemInstanceOneRef = "demo-item-instance-one"

	DemoCreatureInstanceOneRef = "demo-creature-instance-one"

	DemoCharacterInstanceOneRef = "demo-character-instance-one"

	DemoImageJoinGameRef  = "demo-image-join-game"
	DemoImageInventoryRef = "demo-image-inventory"
)

// FullAdventureConfig returns a standalone demo scenario exercising all adventure game features.
// All games are created as draft. Use --publish to publish them after loading.
func FullAdventureConfig() harness.DataConfig {
	return harness.DataConfig{
		AccountConfigs: demoAccountConfigs(),
		GameConfigs:    demoGameConfigs(),
	}
}

func demoAccountConfigs() []harness.AccountConfig {
	return []harness.AccountConfig{
		{
			Reference: DemoAccountDesignerRef,
			Record: &account_record.AccountUser{
				Email: "demo-designer@example.com",
			},
			GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
				{
					Reference:        DemoSubscriptionDesignerOneRef,
					GameRef:          DemoGameOneRef,
					SubscriptionType: game_record.GameSubscriptionTypeDesigner,
					Record:           &game_record.GameSubscription{},
				},
				{
					Reference:        DemoSubscriptionDesignerTwoRef,
					GameRef:          DemoGameTwoRef,
					SubscriptionType: game_record.GameSubscriptionTypeDesigner,
					Record:           &game_record.GameSubscription{},
				},
			},
		},
		{
			Reference: DemoAccountManagerRef,
			Record: &account_record.AccountUser{
				Email: "demo-manager@example.com",
			},
			GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
				{
					Reference:        DemoSubscriptionManagerOneRef,
					GameRef:          DemoGameOneRef,
					SubscriptionType: game_record.GameSubscriptionTypeManager,
					Record:           &game_record.GameSubscription{},
				},
				{
					Reference:        DemoSubscriptionManagerTwoRef,
					GameRef:          DemoGameTwoRef,
					SubscriptionType: game_record.GameSubscriptionTypeManager,
					Record:           &game_record.GameSubscription{},
				},
			},
		},
		{
			Reference: DemoAccountPlayerRef,
			Record: &account_record.AccountUser{
				Email: "demo-player@example.com",
			},
			GameSubscriptionConfigs: []harness.GameSubscriptionConfig{
				{
					Reference:        DemoSubscriptionPlayerOneRef,
					GameRef:          DemoGameOneRef,
					SubscriptionType: game_record.GameSubscriptionTypePlayer,
					Record:           &game_record.GameSubscription{},
				},
				{
					Reference:        DemoSubscriptionPlayerTwoRef,
					GameRef:          DemoGameTwoRef,
					SubscriptionType: game_record.GameSubscriptionTypePlayer,
					Record:           &game_record.GameSubscription{},
				},
			},
		},
	}
}

func demoGameConfigs() []harness.GameConfig {
	return []harness.GameConfig{
		{
			Reference: DemoGameOneRef,
			Record: &game_record.Game{
				Name:              "The Enchanted Forest Adventure",
				Description:       "Step into a world of magic and mystery where ancient forests hold secrets waiting to be discovered. Journey through mystical groves, crystal caverns, and floating islands as you encounter magical creatures, solve puzzles, and uncover the mysteries of this enchanted realm. Your choices matter - every decision shapes your adventure.",
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 168,
				Status:            game_record.GameStatusDraft,
			},
			GameImageConfigs: []harness.GameImageConfig{
				{
					Reference:     DemoImageJoinGameRef,
					ImagePath:     ImageJoinGame,
					TurnSheetType: adventure_game_record.AdventureGameTurnSheetTypeJoinGame,
				},
				{
					Reference:     DemoImageInventoryRef,
					ImagePath:     ImageInventoryManagement,
					TurnSheetType: adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement,
				},
			},
			AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
				{
					Reference: DemoLocationOneRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:               "Mystic Grove",
						Description:        "A peaceful grove filled with ancient trees and magical flowers. The air shimmers with enchantment.",
						IsStartingLocation: true,
					},
					BackgroundImagePath: ImageLocationDarkforest,
				},
				{
					Reference: DemoLocationTwoRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Crystal Caverns",
						Description: "Deep underground caves filled with glowing crystals. Strange sounds echo from the depths.",
					},
					BackgroundImagePath: ImageLocationDungeon,
				},
				{
					Reference: DemoLocationThreeRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Floating Islands",
						Description: "Mysterious islands suspended in the sky by unknown magic. Wind howls between them.",
					},
					BackgroundImagePath: ImageLocationCliffpath,
				},
				{
					Reference: DemoLocationFourRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Shadow Valley",
						Description: "A dark valley shrouded in perpetual shadows. Danger lurks in every corner.",
					},
					BackgroundImagePath: ImageLocationDarkforest,
				},
			},
			AdventureGameItemConfigs: []harness.AdventureGameItemConfig{
				{
					Reference: DemoItemOneRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Crystal Key",
						Description: "A glowing key made of pure crystal. It hums with magical energy.",
					},
				},
				{
					Reference: DemoItemTwoRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Shadow Cloak",
						Description: "A cloak that allows the wearer to blend into shadows and move silently.",
					},
				},
				{
					Reference: DemoItemThreeRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Healing Potion",
						Description: "A bright blue potion that restores health and removes minor ailments.",
					},
				},
				{
					Reference: DemoItemFourRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Wind Charm",
						Description: "A small charm that allows the bearer to control wind currents.",
					},
				},
			},
			AdventureGameCreatureConfigs: []harness.AdventureGameCreatureConfig{
				{
					Reference: DemoCreatureOneRef,
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Forest Guardian",
						Description: "A majestic creature made of living wood and leaves. Protects the grove.",
					},
				},
				{
					Reference: DemoCreatureTwoRef,
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Crystal Spider",
						Description: "A giant spider with a body made of living crystal. Spins webs of light.",
					},
				},
			},
			AdventureGameCharacterConfigs: []harness.AdventureGameCharacterConfig{
				{
					Reference:  DemoCharacterOneRef,
					AccountRef: DemoAccountDesignerRef,
					Record: &adventure_game_record.AdventureGameCharacter{
						Name: "Aria the Mage",
					},
				},
				{
					Reference:  DemoCharacterTwoRef,
					AccountRef: DemoAccountManagerRef,
					Record: &adventure_game_record.AdventureGameCharacter{
						Name: "Thorne the Warrior",
					},
				},
				{
					Reference:  DemoCharacterThreeRef,
					AccountRef: DemoAccountPlayerRef,
					Record: &adventure_game_record.AdventureGameCharacter{
						Name: "Luna the Scout",
					},
				},
			},
			AdventureGameLocationLinkConfigs: []harness.AdventureGameLocationLinkConfig{
				// Grove -> Caverns (requires Crystal Key)
				{
					Reference:       DemoLocationLinkOneRef,
					FromLocationRef: DemoLocationOneRef,
					ToLocationRef:   DemoLocationTwoRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Crystal Path",
						Description: "A winding path that leads from the grove down into the crystal caverns.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   DemoLinkReqOneRef,
							GameItemRef: DemoItemOneRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				// Grove -> Islands (requires Wind Charm)
				{
					Reference:       "demo-link-grove-to-islands",
					FromLocationRef: DemoLocationOneRef,
					ToLocationRef:   DemoLocationThreeRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Sky Vine",
						Description: "A massive enchanted vine that grows from the grove up to the floating islands.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "demo-link-req-grove-to-islands",
							GameItemRef: DemoItemFourRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				// Grove -> Shadow (requires Shadow Cloak)
				{
					Reference:       "demo-link-grove-to-shadow",
					FromLocationRef: DemoLocationOneRef,
					ToLocationRef:   DemoLocationFourRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Dark Tunnel",
						Description: "A hidden tunnel beneath ancient roots that leads directly to the shadow valley.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "demo-link-req-grove-to-shadow",
							GameItemRef: DemoItemTwoRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				// Caverns -> Islands (requires Wind Charm)
				{
					Reference:       DemoLocationLinkTwoRef,
					FromLocationRef: DemoLocationTwoRef,
					ToLocationRef:   DemoLocationThreeRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Wind Lift",
						Description: "A magical elevator that rises from the caverns to the floating islands.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   DemoLinkReqTwoRef,
							GameItemRef: DemoItemFourRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				// Caverns -> Grove (no requirement)
				{
					Reference:       "demo-link-caverns-to-grove",
					FromLocationRef: DemoLocationTwoRef,
					ToLocationRef:   DemoLocationOneRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Glowing Passage",
						Description: "Crystal-lit tunnels that wind back up to the mystic grove.",
					},
				},
				// Caverns -> Shadow (requires Healing Potion)
				{
					Reference:       "demo-link-caverns-to-shadow",
					FromLocationRef: DemoLocationTwoRef,
					ToLocationRef:   DemoLocationFourRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Deep Descent",
						Description: "A treacherous path that plunges from the caverns into the shadow valley.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "demo-link-req-caverns-to-shadow",
							GameItemRef: DemoItemThreeRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				// Islands -> Shadow (requires Shadow Cloak)
				{
					Reference:       DemoLocationLinkThreeRef,
					FromLocationRef: DemoLocationThreeRef,
					ToLocationRef:   DemoLocationFourRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Shadow Bridge",
						Description: "A bridge of pure darkness that connects the floating islands to the shadow valley.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   DemoLinkReqThreeRef,
							GameItemRef: DemoItemTwoRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				// Islands -> Grove (no requirement)
				{
					Reference:       "demo-link-islands-to-grove",
					FromLocationRef: DemoLocationThreeRef,
					ToLocationRef:   DemoLocationOneRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Feather Fall",
						Description: "Enchanted feathers allow a gentle descent from the islands to the grove below.",
					},
				},
				// Islands -> Caverns (no requirement)
				{
					Reference:       "demo-link-islands-to-caverns",
					FromLocationRef: DemoLocationThreeRef,
					ToLocationRef:   DemoLocationTwoRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Crystal Chute",
						Description: "A smooth crystal slide that spirals down into the caverns.",
					},
				},
				// Shadow -> Grove (requires Healing Potion)
				{
					Reference:       DemoLocationLinkFourRef,
					FromLocationRef: DemoLocationFourRef,
					ToLocationRef:   DemoLocationOneRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Return Portal",
						Description: "A magical portal that allows quick return to the mystic grove.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   DemoLinkReqFourRef,
							GameItemRef: DemoItemThreeRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
				// Shadow -> Caverns (no requirement)
				{
					Reference:       "demo-link-shadow-to-caverns",
					FromLocationRef: DemoLocationFourRef,
					ToLocationRef:   DemoLocationTwoRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Echoing Stairs",
						Description: "Ancient stone stairs that climb from the valley up to the crystal caverns.",
					},
				},
				// Shadow -> Islands (requires Shadow Cloak)
				{
					Reference:       "demo-link-shadow-to-islands",
					FromLocationRef: DemoLocationFourRef,
					ToLocationRef:   DemoLocationThreeRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Shadow Ascent",
						Description: "Dark tendrils of shadow that can lift travelers to the floating islands.",
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "demo-link-req-shadow-to-islands",
							GameItemRef: DemoItemTwoRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Quantity: 1,
							},
						},
					},
				},
			},
			GameInstanceConfigs: []harness.GameInstanceConfig{
				{
					Reference: DemoInstanceOneRef,
					Record:    &game_record.GameInstance{},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: DemoInstanceParamOneRef,
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   domain.AdventureGameParameterCharacterLives,
								ParameterValue: nullstring.FromString("5"),
							},
						},
					},
					AdventureGameLocationInstanceConfigs: []harness.AdventureGameLocationInstanceConfig{
						{
							Reference:       DemoLocationInstanceOneRef,
							GameLocationRef: DemoLocationOneRef,
							Record:          &adventure_game_record.AdventureGameLocationInstance{},
						},
						{
							Reference:       DemoLocationInstanceTwoRef,
							GameLocationRef: DemoLocationTwoRef,
							Record:          &adventure_game_record.AdventureGameLocationInstance{},
						},
					},
					AdventureGameItemInstanceConfigs: []harness.AdventureGameItemInstanceConfig{
						{
							Reference:       DemoItemInstanceOneRef,
							GameItemRef:     DemoItemOneRef,
							GameLocationRef: DemoLocationOneRef,
							Record:          &adventure_game_record.AdventureGameItemInstance{},
						},
					},
					AdventureGameCreatureInstanceConfigs: []harness.AdventureGameCreatureInstanceConfig{
						{
							Reference:       DemoCreatureInstanceOneRef,
							GameCreatureRef: DemoCreatureOneRef,
							GameLocationRef: DemoLocationOneRef,
							Record:          &adventure_game_record.AdventureGameCreatureInstance{},
						},
					},
					AdventureGameCharacterInstanceConfigs: []harness.AdventureGameCharacterInstanceConfig{
						{
							Reference:        DemoCharacterInstanceOneRef,
							GameCharacterRef: DemoCharacterOneRef,
							GameLocationRef:  DemoLocationOneRef,
							Record:           &adventure_game_record.AdventureGameCharacterInstance{},
						},
					},
				},
			},
		},
		{
			Reference: DemoGameTwoRef,
			Record: &game_record.Game{
				Name:              "The Desert Kingdom",
				Description:       "Embark on an epic journey across vast sand dunes, ancient ruins, and hidden oases. Navigate treacherous terrain, encounter nomadic tribes, and uncover the lost secrets of a forgotten civilization. In this harsh but beautiful landscape, survival requires wit, courage, and careful planning.",
				GameType:          game_record.GameTypeAdventure,
				TurnDurationHours: 336,
				Status:            game_record.GameStatusDraft,
			},
			AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
				{
					Reference: "demo-location-five",
					Record: &adventure_game_record.AdventureGameLocation{
						Name:               "Oasis Village",
						Description:        "A bustling village built around a life-giving oasis in the desert.",
						IsStartingLocation: true,
					},
				},
				{
					Reference: "demo-location-six",
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "Ancient Ruins",
						Description: "Crumbling ruins of a lost civilization, filled with secrets and danger.",
					},
				},
			},
			AdventureGameItemConfigs: []harness.AdventureGameItemConfig{
				{
					Reference: "demo-item-five",
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Desert Compass",
						Description: "A magical compass that always points to water sources.",
					},
				},
			},
			AdventureGameCreatureConfigs: []harness.AdventureGameCreatureConfig{
				{
					Reference: "demo-creature-three",
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Sand Serpent",
						Description: "A massive serpent that burrows through the desert sands.",
					},
				},
			},
			GameInstanceConfigs: []harness.GameInstanceConfig{
				{
					Reference: DemoInstanceTwoRef,
					Record:    &game_record.GameInstance{},
					GameInstanceParameterConfigs: []harness.GameInstanceParameterConfig{
						{
							Reference: DemoInstanceParamTwoRef,
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
