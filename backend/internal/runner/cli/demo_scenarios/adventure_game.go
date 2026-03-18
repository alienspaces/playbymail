package demo_scenarios

import (
		"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// Harness references (prefixed "demo-" to avoid collisions with test data)
const (
	DemoLocGrandStaircaseRef    = "demo-loc-grand-staircase"
	DemoLocNarrowPassageRef     = "demo-loc-narrow-passage"
	DemoLocWineCellarRef        = "demo-loc-wine-cellar"
	DemoLocCryptRef             = "demo-loc-crypt"
	DemoLocUndergroundChapelRef = "demo-loc-underground-chapel"
	DemoLocFloodedCorridorRef   = "demo-loc-flooded-corridor"
	DemoLocAbbotsStudyRef       = "demo-loc-abbots-study"
	DemoLocWellChamberRef       = "demo-loc-well-chamber"
	DemoLocHerbGardenRef        = "demo-loc-herb-garden"
	DemoLocBellTowerVaultRef    = "demo-loc-bell-tower-vault"

	DemoItemRustyKeyRef      = "demo-item-rusty-key"
	DemoItemTallowCandleRef  = "demo-item-tallow-candle"
	DemoItemAbbotsJournalRef = "demo-item-abbots-journal"
	DemoItemSilverCrossRef   = "demo-item-silver-cross"
	DemoItemCoilOfRopeRef    = "demo-item-coil-of-rope"

	DemoCreatureShadowMonkRef = "demo-creature-shadow-monk"
	DemoCreatureCellarRatRef  = "demo-creature-cellar-rat"

	DemoInstanceOneRef      = "demo-instance-one"
	DemoInstanceParamOneRef = "demo-instance-param-one"

	DemoLocInstanceGrandStaircaseRef = "demo-loc-inst-grand-staircase"
	DemoLocInstanceNarrowPassageRef  = "demo-loc-inst-narrow-passage"
	DemoLocInstanceWineCellarRef     = "demo-loc-inst-wine-cellar"

	DemoItemInstanceRustyKeyRef      = "demo-item-inst-rusty-key"
	DemoCreatureInstanceCellarRatRef = "demo-creature-inst-cellar-rat"

	DemoImageJoinGameRef  = "demo-image-join-game"
	DemoImageInventoryRef = "demo-image-inventory"
)

// Image file names in demo_scenario_images directory
const (
	ImageJoinGame             = "join-game.jpg"
	ImageInventoryManagement  = "inventory-management.jpg"
	ImageLocGrandStaircase    = "location-grand-staircase.jpg"
	ImageLocNarrowPassage     = "location-narrow-passage.jpg"
	ImageLocWineCellar        = "location-wine-cellar.jpg"
	ImageLocCrypt             = "location-crypt.jpg"
	ImageLocUndergroundChapel = "location-underground-chapel.jpg"
	ImageLocFloodedCorridor   = "location-flooded-corridor.jpg"
	ImageLocAbbotsStudy       = "location-abbots-study.jpg"
	ImageLocWellChamber       = "location-well-chamber.jpg"
	ImageLocHerbGarden        = "location-herb-garden.jpg"
	ImageLocBellTowerVault    = "location-bell-tower-vault.jpg"
)

func adventureGameConfigs() []harness.GameConfig {
	return []harness.GameConfig{
		{
			Reference: DemoAdventureGameRef,
			Record: &game_record.Game{
				Name:              DemoAdventureGameName,
				Description:       "An old abbey stands silent on a windswept hill. Beneath its grand staircase, a small wooden door opens onto passages long forgotten. Descend into wine cellars, flooded corridors, and ancient crypts. Uncover the secrets the monks left behind -- if you dare.",
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

			// ── Locations ──────────────────────────────────────────────
			AdventureGameLocationConfigs: []harness.AdventureGameLocationConfig{
				{
					Reference: DemoLocGrandStaircaseRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:               "The Grand Staircase",
						Description:        "A wide stone staircase sweeps upward through the abbey's entrance hall. Dust motes drift in the pale light. Beneath the lowest step, a small wooden door is barely visible, its iron handle dark with age.",
						IsStartingLocation: true,
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocGrandStaircase},
				},
				{
					Reference: DemoLocNarrowPassageRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Narrow Passage",
						Description: "A cramped tunnel of rough-hewn stone stretches ahead, so low you must stoop. A single candle would barely push back the darkness here.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocNarrowPassage},
				},
				{
					Reference: DemoLocWineCellarRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Wine Cellar",
						Description: "Vaulted stone arches line a long cellar. Dusty bottles fill sagging wooden shelves, and cobwebs drape every surface. Something rustles in the shadows between the racks.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocWineCellar},
				},
				{
					Reference: DemoLocCryptRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Crypt",
						Description: "Cold stone sarcophagi rest in alcoves along the walls. Faded inscriptions in an old language cover the floor. A chill hangs in the air that has nothing to do with the depth underground.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocCrypt},
				},
				{
					Reference: DemoLocUndergroundChapelRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Underground Chapel",
						Description: "A forgotten chapel carved from the bedrock. Crumbling wooden pews face a cracked stone altar. Faint light filters down through a narrow shaft in the ceiling far above.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocUndergroundChapel},
				},
				{
					Reference: DemoLocFloodedCorridorRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Flooded Corridor",
						Description: "Dark water covers the floor of this vaulted passage, reflecting distant torchlight. Drips echo from the ceiling. The water is cold and still, hiding whatever lies beneath.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocFloodedCorridor},
				},
				{
					Reference: DemoLocAbbotsStudyRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Abbot's Study",
						Description: "A hidden room behind a false wall. A heavy oak desk is buried under scattered manuscripts, and old leather-bound books line the shelves. A single candle stub sits in a brass holder.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocAbbotsStudy},
				},
				{
					Reference: DemoLocWellChamberRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Well Chamber",
						Description: "A circular stone chamber built around an ancient well. A frayed rope hangs over the lip, and darkness yawns below. The stones are slick with moisture.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocWellChamber},
				},
				{
					Reference: DemoLocHerbGardenRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Herb Garden",
						Description: "An overgrown walled garden behind the abbey. Tangled herbs push through cracked flagstones. Crumbling stone walls frame a grey overcast sky. The air smells of thyme and damp earth.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocHerbGarden},
				},
				{
					Reference: DemoLocBellTowerVaultRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Bell Tower Vault",
						Description: "A small secret vault beneath the bell tower, reached by a hidden stair. Iron-bound chests sit under low stone arches. Dust motes swirl in a thin shaft of light from above.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocBellTowerVault},
				},
			},

			// ── Items ──────────────────────────────────────────────────
			AdventureGameItemConfigs: []harness.AdventureGameItemConfig{
				{
					Reference: DemoItemRustyKeyRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Rusty Key",
						Description: "A heavy iron key, rough with rust. It looks like it might fit an old lock.",
					},
				},
				{
					Reference: DemoItemTallowCandleRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Tallow Candle",
						Description: "A stubby candle of yellowed tallow. Its flame gutters but holds, pushing back the dark.",
					},
				},
				{
					Reference: DemoItemAbbotsJournalRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Abbot's Journal",
						Description: "A cracked leather journal filled with faded ink. The last entries speak of something sealed below the chapel.",
					},
				},
				{
					Reference: DemoItemSilverCrossRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Silver Cross",
						Description: "A small silver cross on a tarnished chain. It feels oddly warm to the touch.",
					},
				},
				{
					Reference: DemoItemCoilOfRopeRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:        "Coil of Rope",
						Description: "A length of sturdy hemp rope, coiled and ready. Useful for climbing or descending into deep places.",
					},
				},
			},

			// ── Creatures ──────────────────────────────────────────────
			AdventureGameCreatureConfigs: []harness.AdventureGameCreatureConfig{
				{
					Reference: DemoCreatureShadowMonkRef,
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Shadow Monk",
						Description: "A spectral figure in a hooded robe drifts silently through the crypt. Its face is hidden, but cold radiates from it like a winter wind.",
					},
				},
				{
					Reference: DemoCreatureCellarRatRef,
					Record: &adventure_game_record.AdventureGameCreature{
						Name:        "Cellar Rat",
						Description: "A large grey rat with bright eyes. It watches from the shadows between the wine racks, unafraid.",
					},
				},
			},

			// ── Location Links ─────────────────────────────────────────
			AdventureGameLocationLinkConfigs: []harness.AdventureGameLocationLinkConfig{
			// Grand Staircase -> Narrow Passage (requires Rusty Key in inventory)
			{
				Reference:       "demo-link-staircase-to-passage",
				FromLocationRef: DemoLocGrandStaircaseRef,
				ToLocationRef:   DemoLocNarrowPassageRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:              "The Small Door",
					Description:       "A low wooden door beneath the staircase. The lock is stiff but yields to the right key.",
					LockedDescription: nullstring.FromString("A low wooden door bound with iron is set beneath the staircase. It is firmly locked."),
				},
				AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
					{
						Reference:   "demo-link-req-staircase-to-passage",
						GameItemRef: DemoItemRustyKeyRef,
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory,
							Quantity:  1,
						},
					},
				},
			},
				// Grand Staircase -> Herb Garden (no requirement)
				{
					Reference:       "demo-link-staircase-to-garden",
					FromLocationRef: DemoLocGrandStaircaseRef,
					ToLocationRef:   DemoLocHerbGardenRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Side Door",
						Description: "A weathered door at the back of the entrance hall opens onto the walled garden.",
					},
				},
				// Narrow Passage -> Wine Cellar (no requirement)
				{
					Reference:       "demo-link-passage-to-cellar",
					FromLocationRef: DemoLocNarrowPassageRef,
					ToLocationRef:   DemoLocWineCellarRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "Cellar Steps",
						Description: "Worn stone steps descend from the passage into the wine cellar below.",
					},
				},
				// Narrow Passage -> Grand Staircase (no requirement, return path)
				{
					Reference:       "demo-link-passage-to-staircase",
					FromLocationRef: DemoLocNarrowPassageRef,
					ToLocationRef:   DemoLocGrandStaircaseRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "Back Through the Door",
						Description: "The small door leads back to the grand staircase.",
					},
				},
			// Wine Cellar -> Crypt (requires Tallow Candle; hidden while Cellar Rat is alive)
			{
				Reference:       "demo-link-cellar-to-crypt",
				FromLocationRef: DemoLocWineCellarRef,
				ToLocationRef:   DemoLocCryptRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:              "The Dark Archway",
					Description:       "A low archway at the far end of the cellar leads into pitch darkness. Only a fool would enter without light.",
					LockedDescription: nullstring.FromString("The archway ahead fades into absolute darkness. Without a light you dare not proceed."),
				},
				AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
					{
						Reference:       "demo-link-req-cellar-rat-dead",
						GameCreatureRef: DemoCreatureCellarRatRef,
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeVisible,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveAtLocation,
							Quantity:  1,
						},
					},
					{
						Reference:   "demo-link-req-cellar-to-crypt",
						GameItemRef: DemoItemTallowCandleRef,
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory,
							Quantity:  1,
						},
					},
				},
			},
				// Wine Cellar -> Narrow Passage (no requirement, return path)
				{
					Reference:       "demo-link-cellar-to-passage",
					FromLocationRef: DemoLocWineCellarRef,
					ToLocationRef:   DemoLocNarrowPassageRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "Stone Steps Up",
						Description: "The steps climb back up to the narrow passage above.",
					},
				},
			// Crypt -> Underground Chapel (requires Silver Cross in inventory)
			{
				Reference:       "demo-link-crypt-to-chapel",
				FromLocationRef: DemoLocCryptRef,
				ToLocationRef:   DemoLocUndergroundChapelRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:              "The Warded Gate",
					Description:       "An iron gate covered in faded holy symbols. The cold intensifies near it. Something silver might grant safe passage.",
					LockedDescription: nullstring.FromString("An iron gate engraved with religious symbols bars the way. Something about it unsettles you."),
				},
				AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
					{
						Reference:   "demo-link-req-crypt-to-chapel",
						GameItemRef: DemoItemSilverCrossRef,
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory,
							Quantity:  1,
						},
					},
				},
			},
				// Crypt -> Wine Cellar (no requirement, return path)
				{
					Reference:       "demo-link-crypt-to-cellar",
					FromLocationRef: DemoLocCryptRef,
					ToLocationRef:   DemoLocWineCellarRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "Back to the Cellar",
						Description: "The archway leads back into the wine cellar's dusty warmth.",
					},
				},
			// Underground Chapel -> Flooded Corridor (requires Tallow Candle in inventory)
			{
				Reference:       "demo-link-chapel-to-flooded",
				FromLocationRef: DemoLocUndergroundChapelRef,
				ToLocationRef:   DemoLocFloodedCorridorRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:              "The Broken Wall",
					Description:       "A gap in the chapel wall reveals a passage sloping downward. Water glints in the darkness ahead.",
					LockedDescription: nullstring.FromString("A gap in the chapel wall leads into darkness below. Without a light source, descending would be madness."),
				},
				AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
					{
						Reference:   "demo-link-req-chapel-to-flooded",
						GameItemRef: DemoItemTallowCandleRef,
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory,
							Quantity:  1,
						},
					},
				},
			},
				// Underground Chapel -> Abbot's Study (no requirement)
				{
					Reference:       "demo-link-chapel-to-study",
					FromLocationRef: DemoLocUndergroundChapelRef,
					ToLocationRef:   DemoLocAbbotsStudyRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Hidden Panel",
						Description: "A stone panel behind the altar slides aside to reveal a small room beyond.",
					},
				},
			// Flooded Corridor -> Well Chamber (requires Coil of Rope in inventory)
			{
				Reference:       "demo-link-flooded-to-well",
				FromLocationRef: DemoLocFloodedCorridorRef,
				ToLocationRef:   DemoLocWellChamberRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:              "The Submerged Ledge",
					Description:       "A narrow ledge runs along the flooded corridor to a shaft leading upward. Rope would make the climb possible.",
					LockedDescription: nullstring.FromString("A submerged stone ledge leads deeper into the water. Without a rope, the descent looks fatal."),
				},
				AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
					{
						Reference:   "demo-link-req-flooded-to-well",
						GameItemRef: DemoItemCoilOfRopeRef,
						Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
							Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
							Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory,
							Quantity:  1,
						},
					},
				},
			},
				// Well Chamber -> Bell Tower Vault (no requirement)
				{
					Reference:       "demo-link-well-to-vault",
					FromLocationRef: DemoLocWellChamberRef,
					ToLocationRef:   DemoLocBellTowerVaultRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Hidden Stair",
						Description: "Behind a loose stone in the well chamber, a narrow staircase spirals upward into the bell tower.",
					},
				},
				// Herb Garden -> Grand Staircase (no requirement, return path)
				{
					Reference:       "demo-link-garden-to-staircase",
					FromLocationRef: DemoLocHerbGardenRef,
					ToLocationRef:   DemoLocGrandStaircaseRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "Back Inside",
						Description: "The side door leads back into the abbey's entrance hall.",
					},
				},
				// Bell Tower Vault -> Well Chamber (no requirement, return path)
				{
					Reference:       "demo-link-vault-to-well",
					FromLocationRef: DemoLocBellTowerVaultRef,
					ToLocationRef:   DemoLocWellChamberRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:        "The Spiral Stair Down",
						Description: "The narrow staircase winds back down to the well chamber.",
					},
				},
			},
		},
	}
}
