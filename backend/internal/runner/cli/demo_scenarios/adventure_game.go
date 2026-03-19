package demo_scenarios

import (
	"gitlab.com/alienspaces/playbymail/core/nullint32"
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

	// Location object references
	DemoObjAncientWellRef     = "demo-obj-ancient-well"
	DemoObjStoneAltarRef      = "demo-obj-stone-altar"
	DemoObjManuscriptsRef     = "demo-obj-manuscripts"
	DemoObjIronChestRef       = "demo-obj-iron-chest"
	DemoObjRustedGateRef      = "demo-obj-rusted-gate"
	DemoObjLeverRef           = "demo-obj-lever"
	DemoObjHiddenDoorRef      = "demo-obj-hidden-door"

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

	ImageCreatureShadowMonk = "creature-shadow-monk.jpg"
	ImageCreatureCellarRat  = "creature-cellar-rat.jpg"
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
					Name:              "Shadow Monk",
					Description:       "A spectral figure in a hooded robe drifts silently through the crypt. Its face is hidden, but cold radiates from it like a winter wind.",
					MaxHealth:         80,
					AttackDamage:      15,
					Defense:           5,
					Disposition:       adventure_game_record.AdventureGameCreatureDispositionAggressive,
					AttackMethod:      adventure_game_record.AdventureGameCreatureAttackMethodTouch,
					AttackDescription: "reaches through you with a spectral hand",
					BodyDecayTurns:    3,
					RespawnTurns:      5,
				},
					PortraitImage: &harness.GameImageConfig{ImagePath: ImageCreatureShadowMonk},
				},
				{
			Reference: DemoCreatureCellarRatRef,
				Record: &adventure_game_record.AdventureGameCreature{
					Name:              "Cellar Rat",
					Description:       "A large grey rat with bright eyes. It watches from the shadows between the wine racks, unafraid.",
					MaxHealth:         20,
					AttackDamage:      5,
					Defense:           0,
					Disposition:       adventure_game_record.AdventureGameCreatureDispositionInquisitive,
					AttackMethod:      adventure_game_record.AdventureGameCreatureAttackMethodBite,
					AttackDescription: "snaps at your ankles with sharp teeth",
					BodyDecayTurns:    2,
					RespawnTurns:      3,
				},
					PortraitImage: &harness.GameImageConfig{ImagePath: ImageCreatureCellarRat},
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
						Name:                 "The Small Door",
						Description:          "A low wooden door beneath the staircase. The lock is stiff but yields to the right key.",
						LockedDescription:    nullstring.FromString("A low wooden door bound with iron is set beneath the staircase. It is firmly locked."),
						TraversalDescription: nullstring.FromString("You push open the small wooden door. Cold air rises from the darkness below as you descend a narrow stone staircase."),
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
						Name:                 "The Side Door",
						Description:          "A weathered door at the back of the entrance hall opens onto the walled garden.",
						TraversalDescription: nullstring.FromString("You push through the weathered door and step into the walled herb garden. The scent of rosemary and thyme fills the air."),
					},
				},
			// Narrow Passage -> Wine Cellar (no requirement)
			{
				Reference:       "demo-link-passage-to-cellar",
				FromLocationRef: DemoLocNarrowPassageRef,
				ToLocationRef:   DemoLocWineCellarRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:                 "The Cellar Steps",
					Description:          "Worn stone steps descend from the passage into the wine cellar below.",
					TraversalDescription: nullstring.FromString("You descend the worn stone steps carefully, one hand trailing along the damp wall. The smell of old wine and earth surrounds you."),
				},
			},
			// Narrow Passage -> Grand Staircase (no requirement, return path)
			{
				Reference:       "demo-link-passage-to-staircase",
				FromLocationRef: DemoLocNarrowPassageRef,
				ToLocationRef:   DemoLocGrandStaircaseRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:                 "The Small Door",
					Description:          "The small door opens back onto the foot of the grand staircase.",
					TraversalDescription: nullstring.FromString("You duck back through the small door and climb the narrow staircase into the entrance hall."),
				},
			},
				// Wine Cellar -> Crypt (requires Tallow Candle; hidden while Cellar Rat is alive)
				{
					Reference:       "demo-link-cellar-to-crypt",
					FromLocationRef: DemoLocWineCellarRef,
					ToLocationRef:   DemoLocCryptRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:                 "The Dark Archway",
						Description:          "A low archway at the far end of the cellar leads into pitch darkness. Only a fool would enter without light.",
						LockedDescription:    nullstring.FromString("The archway ahead fades into absolute darkness. Without a light you dare not proceed."),
						TraversalDescription: nullstring.FromString("Holding your candle before you, you step through the dark archway. The flame gutters as you descend into the crypt's cold silence."),
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
					Name:                 "The Cellar Steps",
					Description:          "The worn stone steps climb back up to the narrow passage above.",
					TraversalDescription: nullstring.FromString("You climb the worn stone steps back up into the narrow passage, leaving the earthy smell of the cellar behind."),
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
			// Underground Chapel -> Crypt (no requirement, return path)
			{
				Reference:       "demo-link-chapel-to-crypt",
				FromLocationRef: DemoLocUndergroundChapelRef,
				ToLocationRef:   DemoLocCryptRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:                 "The Warded Gate",
					Description:          "The iron gate bearing holy symbols leads back into the cold silence of the crypt.",
					TraversalDescription: nullstring.FromString("You pass back through the warded gate into the crypt, its chill settling around you once more."),
				},
			},
			// Crypt -> Wine Cellar (no requirement, return path)
			{
				Reference:       "demo-link-crypt-to-cellar",
				FromLocationRef: DemoLocCryptRef,
				ToLocationRef:   DemoLocWineCellarRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:                 "The Dark Archway",
					Description:          "The dark archway leads back into the dusty warmth of the wine cellar.",
					TraversalDescription: nullstring.FromString("You pass back through the dark archway and ascend into the wine cellar's dusty warmth, grateful to leave the cold silence behind."),
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
			// Flooded Corridor -> Underground Chapel (no requirement, return path)
			{
				Reference:       "demo-link-flooded-to-chapel",
				FromLocationRef: DemoLocFloodedCorridorRef,
				ToLocationRef:   DemoLocUndergroundChapelRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:                 "The Broken Wall",
					Description:          "The gap in the broken wall leads back up to the underground chapel above.",
					TraversalDescription: nullstring.FromString("You clamber back up through the gap in the broken wall and into the chapel's flickering quiet."),
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
			// Abbot's Study -> Underground Chapel (no requirement, return path)
			{
				Reference:       "demo-link-study-to-chapel",
				FromLocationRef: DemoLocAbbotsStudyRef,
				ToLocationRef:   DemoLocUndergroundChapelRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:                 "The Hidden Panel",
					Description:          "The stone panel slides aside to reveal the underground chapel beyond.",
					TraversalDescription: nullstring.FromString("You press the stone panel and step back through into the underground chapel."),
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
			// Well Chamber -> Flooded Corridor (no requirement, return path)
			{
				Reference:       "demo-link-well-to-flooded",
				FromLocationRef: DemoLocWellChamberRef,
				ToLocationRef:   DemoLocFloodedCorridorRef,
				Record: &adventure_game_record.AdventureGameLocationLink{
					Name:                 "The Submerged Ledge",
					Description:          "The submerged stone ledge leads back down into the flooded corridor below.",
					TraversalDescription: nullstring.FromString("You lower yourself back along the submerged ledge into the flooded corridor, cold water rising around you."),
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
					Name:                 "The Side Door",
					Description:          "The weathered side door leads back into the abbey's entrance hall.",
					TraversalDescription: nullstring.FromString("You push back through the weathered door into the cool stone entrance hall."),
				},
			},
		// Bell Tower Vault -> Well Chamber (no requirement, return path)
		{
			Reference:       "demo-link-vault-to-well",
			FromLocationRef: DemoLocBellTowerVaultRef,
			ToLocationRef:   DemoLocWellChamberRef,
			Record: &adventure_game_record.AdventureGameLocationLink{
				Name:                 "The Hidden Stair",
				Description:          "The hidden staircase winds back down through the stone to the well chamber below.",
				TraversalDescription: nullstring.FromString("You descend the narrow spiral stair, the bell tower's cold air following you down into the well chamber."),
			},
		},
		},

		// ── Location Objects ───────────────────────────────────────
		AdventureGameLocationObjectConfigs: []harness.AdventureGameLocationObjectConfig{

			// ── Hidden Door (created first so the Lever can reference it) ──
			{
				Reference:   DemoObjHiddenDoorRef,
				LocationRef: DemoLocUndergroundChapelRef,
				Record: &adventure_game_record.AdventureGameLocationObject{
					Name:         "Hidden Door",
					Description:  "A section of the chapel wall that seems to blend perfectly with the surrounding stonework. There is something off about the mortar lines.",
					InitialState: "hidden",
					IsHidden:     true,
				},
				AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
					// inspect (any state) → info
					{
						Reference: "demo-obj-effect-hidden-door-inspect",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "The door's stone face is cold and utterly silent. Whatever mechanism opens it is not obvious from this side.",
							IsRepeatable:      true,
						},
					},
					// push (state=revealed) → change_state to open
					{
						Reference: "demo-obj-effect-hidden-door-push",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
							RequiredState:     nullstring.FromString("revealed"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							ResultState:       nullstring.FromString("open"),
							ResultDescription: "The hidden door swings inward on well-balanced pivots, revealing a dark shaft beyond.",
							IsRepeatable:      false,
						},
					},
				},
			},

			// ── Ancient Well (Well Chamber) ──
			{
				Reference:   DemoObjAncientWellRef,
				LocationRef: DemoLocWellChamberRef,
				Record: &adventure_game_record.AdventureGameLocationObject{
					Name:         "Ancient Well",
					Description:  "A circular stone well of great age. A frayed rope hangs into the darkness below. Water glints faintly far down.",
					InitialState: "intact",
				},
				AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
					// inspect → info
					{
						Reference: "demo-obj-effect-well-inspect",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "The well plunges into darkness. Cold air rises from below. The rope looks old, but the knots are tight.",
							IsRepeatable:      true,
						},
					},
					// listen → info
					{
						Reference: "demo-obj-effect-well-listen",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeListen,
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "Faint dripping echoes rise from the well's depths, slow and rhythmic, like a heartbeat in stone.",
							IsRepeatable:      true,
						},
					},
					// search (state=intact) → give_item: Abbot's Journal, change_state to searched
					{
						Reference:   "demo-obj-effect-well-search-give",
						ResultItemRef: DemoItemAbbotsJournalRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
							RequiredState:     nullstring.FromString("intact"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
							ResultDescription: "Wedged beneath a loose stone on the well's inner rim, you find a cracked leather journal — the Abbot's Journal.",
							IsRepeatable:      false,
						},
					},
					{
						Reference: "demo-obj-effect-well-search-state",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
							RequiredState:     nullstring.FromString("intact"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							ResultState:       nullstring.FromString("searched"),
							ResultDescription: "",
							IsRepeatable:      false,
						},
					},
					// search (state=searched) → info
					{
						Reference: "demo-obj-effect-well-searched",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
							RequiredState:     nullstring.FromString("searched"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "You have already searched the well thoroughly. Nothing remains.",
							IsRepeatable:      true,
						},
					},
				},
			},

			// ── Stone Altar (Underground Chapel) ──
			{
				Reference:   DemoObjStoneAltarRef,
				LocationRef: DemoLocUndergroundChapelRef,
				Record: &adventure_game_record.AdventureGameLocationObject{
					Name:         "Stone Altar",
					Description:  "A low cracked altar of dark stone. Wax trails from ancient candles cross its surface. A shallow depression sits at its centre.",
					InitialState: "bare",
				},
				AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
					// inspect → info
					{
						Reference: "demo-obj-effect-altar-inspect",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "The altar bears a shallow depression, perfectly shaped to hold a cross or similar offering. Faint words are carved around its rim in a language you do not know.",
							IsRepeatable:      true,
						},
					},
					// use with silver cross (state=bare) → change_state to blessed, reveal_object: Hidden Door
					{
						Reference:       "demo-obj-effect-altar-use-state",
						RequiredItemRef: DemoItemSilverCrossRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
							RequiredState:     nullstring.FromString("bare"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							ResultState:       nullstring.FromString("blessed"),
							ResultDescription: "You place the silver cross in the depression. The altar shudders. A low grinding resonates through the stone.",
							IsRepeatable:      false,
						},
					},
					{
						Reference:       "demo-obj-effect-altar-use-reveal",
						RequiredItemRef: DemoItemSilverCrossRef,
						ResultObjectRef:  DemoObjHiddenDoorRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
							RequiredState:     nullstring.FromString("bare"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRevealObject,
							ResultDescription: "A section of the chapel wall shifts, revealing what appears to be a hidden door.",
							IsRepeatable:      false,
						},
					},
					// use (state=blessed) → info
					{
						Reference: "demo-obj-effect-altar-use-blessed",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
							RequiredState:     nullstring.FromString("blessed"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "The silver cross rests in the altar's depression, glowing faintly. Its work is done.",
							IsRepeatable:      true,
						},
					},
				},
			},

			// ── Bundle of Manuscripts (Abbot's Study) ──
			{
				Reference:   DemoObjManuscriptsRef,
				LocationRef: DemoLocAbbotsStudyRef,
				Record: &adventure_game_record.AdventureGameLocationObject{
					Name:         "Bundle of Manuscripts",
					Description:  "A stack of aged parchments tied with twine. The writing is cramped and faded, but still legible in places.",
					InitialState: "unread",
				},
				AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
					// inspect → info (repeatable)
					{
						Reference: "demo-obj-effect-manuscripts-inspect",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "Tightly wound parchment scrolls, covered in cramped script. The ink has faded in places but much remains legible.",
							IsRepeatable:      true,
						},
					},
					// read (state=unread) → change_state to read
					{
						Reference: "demo-obj-effect-manuscripts-read-state",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeRead,
							RequiredState:     nullstring.FromString("unread"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							ResultState:       nullstring.FromString("read"),
							ResultDescription: "You read through the manuscripts carefully. They describe a ritual of sealing, performed beneath the chapel floor. The key, they say, is \"faith made solid\" — placed upon the altar.",
							IsRepeatable:      false,
						},
					},
					// read (state=read) → info (repeatable)
					{
						Reference: "demo-obj-effect-manuscripts-reread",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeRead,
							RequiredState:     nullstring.FromString("read"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "You reread the manuscripts. The ritual of sealing requires \"faith made solid\" placed upon the altar. The words do not change.",
							IsRepeatable:      true,
						},
					},
				},
			},

			// ── Iron-Bound Chest (Bell Tower Vault) ──
			{
				Reference:   DemoObjIronChestRef,
				LocationRef: DemoLocBellTowerVaultRef,
				Record: &adventure_game_record.AdventureGameLocationObject{
					Name:         "Iron-Bound Chest",
					Description:  "A heavy oak chest reinforced with iron bands. A large padlock seals it shut.",
					InitialState: "locked",
				},
				AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
					// inspect → info
					{
						Reference: "demo-obj-effect-chest-inspect",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "A heavy iron padlock secures the chest's hasp. The lock looks old but solid. You might need a key.",
							IsRepeatable:      true,
						},
					},
					// unlock with rusty key (state=locked) → change_state to unlocked
					{
						Reference:       "demo-obj-effect-chest-unlock",
						RequiredItemRef: DemoItemRustyKeyRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUnlock,
							RequiredState:     nullstring.FromString("locked"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							ResultState:       nullstring.FromString("unlocked"),
							ResultDescription: "The rusty key grates in the lock, then catches. With a clunk the padlock falls open.",
							IsRepeatable:      false,
						},
					},
					// open (state=unlocked) → give_item: Silver Cross, change_state to open
					{
						Reference:   "demo-obj-effect-chest-open-give",
						ResultItemRef: DemoItemSilverCrossRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
							RequiredState:     nullstring.FromString("unlocked"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
							ResultDescription: "You lift the heavy lid. Inside, resting on rotted velvet, is a silver cross on a tarnished chain.",
							IsRepeatable:      false,
						},
					},
					{
						Reference: "demo-obj-effect-chest-open-state",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
							RequiredState:     nullstring.FromString("unlocked"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							ResultState:       nullstring.FromString("open"),
							ResultDescription: "",
							IsRepeatable:      false,
						},
					},
					// inspect (state=open) → info
					{
						Reference: "demo-obj-effect-chest-inspect-open",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							RequiredState:     nullstring.FromString("open"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "The chest stands open and empty, its velvet lining rotted away. Nothing more remains inside.",
							IsRepeatable:      true,
						},
					},
				},
			},

			// ── Rusted Gate (Crypt) ──
			{
				Reference:   DemoObjRustedGateRef,
				LocationRef: DemoLocCryptRef,
				Record: &adventure_game_record.AdventureGameLocationObject{
					Name:         "Rusted Gate",
					Description:  "A heavy iron gate bars an alcove at the far end of the crypt. Rust has welded the hinges solid.",
					InitialState: "closed",
				},
				AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
					// inspect → info
					{
						Reference: "demo-obj-effect-gate-inspect",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "The iron gate is solid and immovable. Rust has fused the hinges to the stone frame. Brute force might break it.",
							IsRepeatable:      true,
						},
					},
					// break (state=closed) → change_state to broken, damage character
					{
						Reference: "demo-obj-effect-gate-break-state",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBreak,
							RequiredState:     nullstring.FromString("closed"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							ResultState:       nullstring.FromString("broken"),
							ResultDescription: "You hurl yourself against the gate. With a shriek of tortured metal it tears free from the stone, but the jagged edge catches you as it falls.",
							IsRepeatable:      false,
						},
					},
					{
						Reference: "demo-obj-effect-gate-break-damage",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBreak,
							RequiredState:     nullstring.FromString("closed"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeDamage,
							ResultDescription: "The gate's jagged edge opens a cut across your arm.",
							ResultValueMin:    nullint32.FromInt32(5),
							ResultValueMax:    nullint32.FromInt32(10),
							IsRepeatable:      false,
						},
					},
					// inspect (state=broken) → info
					{
						Reference: "demo-obj-effect-gate-inspect-broken",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							RequiredState:     nullstring.FromString("broken"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "The rusted gate lies twisted on the crypt floor. The alcove beyond it is now open.",
							IsRepeatable:      true,
						},
					},
				},
			},

			// ── Lever on the Wall (Well Chamber) — demonstrates cross-object reveal ──
			{
				Reference:   DemoObjLeverRef,
				LocationRef: DemoLocWellChamberRef,
				Record: &adventure_game_record.AdventureGameLocationObject{
					Name:         "Lever on the Wall",
					Description:  "An iron lever set into the chamber wall, its purpose obscure. It is in the UP position.",
					InitialState: "up",
				},
				AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
					// inspect → info
					{
						Reference: "demo-obj-effect-lever-inspect",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
							ResultDescription: "A heavy iron lever, cold to the touch. It is currently UP. Pulling it down might do something.",
							IsRepeatable:      true,
						},
					},
					// pull (state=up) → change_state to down, change_object_state: Hidden Door → revealed, reveal_object: Hidden Door
					{
						Reference: "demo-obj-effect-lever-pull-self",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
							RequiredState:     nullstring.FromString("up"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							ResultState:       nullstring.FromString("down"),
							ResultDescription: "You pull the lever down. A deep mechanical rumble travels through the stone walls.",
							IsRepeatable:      false,
						},
					},
					{
						Reference:       "demo-obj-effect-lever-pull-change-door-state",
						ResultObjectRef: DemoObjHiddenDoorRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
							RequiredState:     nullstring.FromString("up"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeObjectState,
							ResultState:       nullstring.FromString("revealed"),
							ResultDescription: "Somewhere above, a section of the chapel wall shifts with a grinding of old stone.",
							IsRepeatable:      false,
						},
					},
					{
						Reference:       "demo-obj-effect-lever-pull-reveal-door",
						ResultObjectRef: DemoObjHiddenDoorRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
							RequiredState:     nullstring.FromString("up"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRevealObject,
							ResultDescription: "",
							IsRepeatable:      false,
						},
					},
					// push (state=down) → change_state to up, change_object_state: Hidden Door → hidden, hide_object: Hidden Door
					{
						Reference: "demo-obj-effect-lever-push-self",
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
							RequiredState:     nullstring.FromString("down"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
							ResultState:       nullstring.FromString("up"),
							ResultDescription: "You push the lever back up. The grinding rumble returns and then silence.",
							IsRepeatable:      false,
						},
					},
					{
						Reference:       "demo-obj-effect-lever-push-hide-door-state",
						ResultObjectRef: DemoObjHiddenDoorRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
							RequiredState:     nullstring.FromString("down"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeObjectState,
							ResultState:       nullstring.FromString("hidden"),
							ResultDescription: "The section of chapel wall slides back into place.",
							IsRepeatable:      false,
						},
					},
					{
						Reference:       "demo-obj-effect-lever-push-hide-door",
						ResultObjectRef: DemoObjHiddenDoorRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
							RequiredState:     nullstring.FromString("down"),
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHideObject,
							ResultDescription: "",
							IsRepeatable:      false,
						},
					},
				},
			},
		},
	},
}
}
