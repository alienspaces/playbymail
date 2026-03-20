package demo_scenarios

import (
	"gitlab.com/alienspaces/playbymail/core/convert"
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
	DemoObjAncientWellRef = "demo-obj-ancient-well"
	DemoObjStoneAltarRef  = "demo-obj-stone-altar"
	DemoObjManuscriptsRef = "demo-obj-manuscripts"
	DemoObjIronChestRef   = "demo-obj-iron-chest"
	DemoObjRustedGateRef  = "demo-obj-rusted-gate"
	DemoObjLeverRef       = "demo-obj-lever"
	DemoObjHiddenDoorRef  = "demo-obj-hidden-door"

	// Location object state references — Hidden Door
	DemoObjHiddenDoorStateHiddenRef   = "demo-obj-hidden-door-state-hidden"
	DemoObjHiddenDoorStateRevealedRef = "demo-obj-hidden-door-state-revealed"
	DemoObjHiddenDoorStateOpenRef     = "demo-obj-hidden-door-state-open"

	// Location object state references — Ancient Well
	DemoObjWellStateIntactRef   = "demo-obj-well-state-intact"
	DemoObjWellStateSearchedRef = "demo-obj-well-state-searched"

	// Location object state references — Stone Altar
	DemoObjAltarStateBareRef    = "demo-obj-altar-state-bare"
	DemoObjAltarStateBlessedRef = "demo-obj-altar-state-blessed"

	// Location object state references — Bundle of Manuscripts
	DemoObjManuscriptsStateUnreadRef = "demo-obj-manuscripts-state-unread"
	DemoObjManuscriptsStateReadRef   = "demo-obj-manuscripts-state-read"

	// Location object state references — Iron-Bound Chest
	DemoObjChestStateLockedRef   = "demo-obj-chest-state-locked"
	DemoObjChestStateUnlockedRef = "demo-obj-chest-state-unlocked"
	DemoObjChestStateOpenRef     = "demo-obj-chest-state-open"

	// Location object state references — Rusted Gate
	DemoObjGateStateClosedRef = "demo-obj-gate-state-closed"
	DemoObjGateStateBrokenRef = "demo-obj-gate-state-broken"

	// Location object state references — Lever on the Wall
	DemoObjLeverStateUpRef   = "demo-obj-lever-state-up"
	DemoObjLeverStateDownRef = "demo-obj-lever-state-down"

	DemoInstanceOneRef      = "demo-instance-one"
	DemoInstanceParamOneRef = "demo-instance-param-one"

	DemoLocInstanceGrandStaircaseRef = "demo-loc-inst-grand-staircase"
	DemoLocInstanceNarrowPassageRef  = "demo-loc-inst-narrow-passage"
	DemoLocInstanceWineCellarRef     = "demo-loc-inst-wine-cellar"

	DemoItemInstanceRustyKeyRef      = "demo-item-inst-rusty-key"
	DemoCreatureInstanceCellarRatRef = "demo-creature-inst-cellar-rat"

	DemoImageJoinGameRef  = "demo-image-join-game"
	DemoImageInventoryRef = "demo-image-inventory"

	// New location refs
	DemoLocSacristyRef  = "demo-loc-sacristy"
	DemoLocOssuaryRef   = "demo-loc-ossuary"
	DemoLocInfirmaryRef = "demo-loc-infirmary"

	// New item refs
	DemoItemRustyDaggerRef     = "demo-item-rusty-dagger"
	DemoItemMonksIronMaceRef   = "demo-item-monks-iron-mace"
	DemoItemLeatherCuirassRef  = "demo-item-leather-cuirass"
	DemoItemHealingDraughtRef  = "demo-item-healing-draught"
	DemoItemBrassThuriblRef    = "demo-item-brass-thurible"
	DemoItemTarnishedLocketRef = "demo-item-tarnished-locket"

	// New creature refs
	DemoCreatureCryptSpiderRef  = "demo-creature-crypt-spider"
	DemoCreatureBoneRevenantRef = "demo-creature-bone-revenant"
	DemoCreatureDrownedMonkRef  = "demo-creature-drowned-monk"

	// New location object refs
	DemoObjHolyWaterFontRef     = "demo-obj-holy-water-font"
	DemoObjVestmentChestRef     = "demo-obj-vestment-chest"
	DemoObjApothecaryRef        = "demo-obj-apothecary-cabinet"
	DemoObjHerbDryingRackRef    = "demo-obj-herb-drying-rack"
	DemoObjTrappedReliquaryRef  = "demo-obj-trapped-reliquary"
	DemoObjRunicCircleRef       = "demo-obj-runic-circle"
	DemoObjCursedSarcophagusRef = "demo-obj-cursed-sarcophagus"
	DemoObjPortcullisWinchRef   = "demo-obj-portcullis-winch"

	// Holy Water Font states
	DemoObjHolyWaterFontStateFlowingRef = "demo-obj-holy-water-font-state-flowing"
	DemoObjHolyWaterFontStateDryRef     = "demo-obj-holy-water-font-state-dry"

	// Vestment Chest states
	DemoObjVestmentChestStateClosedRef = "demo-obj-vestment-chest-state-closed"
	DemoObjVestmentChestStateOpenRef   = "demo-obj-vestment-chest-state-open"

	// Apothecary Cabinet states
	DemoObjApothecaryStateUndisturbedRef = "demo-obj-apothecary-state-undisturbed"
	DemoObjApothecaryStateSearchedRef    = "demo-obj-apothecary-state-searched"

	// Herb Drying Rack states
	DemoObjHerbDryingRackStateIntactRef = "demo-obj-herb-drying-rack-state-intact"
	DemoObjHerbDryingRackStateBurnedRef = "demo-obj-herb-drying-rack-state-burned"

	// Trapped Reliquary states
	DemoObjReliquaryStateTrappedRef  = "demo-obj-reliquary-state-trapped"
	DemoObjReliquaryStateDisarmedRef = "demo-obj-reliquary-state-disarmed"
	DemoObjReliquaryStateOpenRef     = "demo-obj-reliquary-state-open"

	// Runic Circle states
	DemoObjRunicCircleStateDormantRef = "demo-obj-runic-circle-state-dormant"
	DemoObjRunicCircleStateActiveRef  = "demo-obj-runic-circle-state-active"

	// Cursed Sarcophagus states
	DemoObjSarcophagusStateSealedRef = "demo-obj-sarcophagus-state-sealed"
	DemoObjSarcophagusStateOpenedRef = "demo-obj-sarcophagus-state-opened"

	// Portcullis Winch states
	DemoObjPortcullisWinchStateUpRef   = "demo-obj-portcullis-winch-state-up"
	DemoObjPortcullisWinchStateDownRef = "demo-obj-portcullis-winch-state-down"
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

	ImageLocSacristy  = "location-sacristy.jpg"
	ImageLocOssuary   = "location-ossuary.jpg"
	ImageLocInfirmary = "location-infirmary.jpg"

	ImageCreatureCryptSpider  = "creature-crypt-spider.jpg"
	ImageCreatureBoneRevenant = "creature-bone-revenant.jpg"
	ImageCreatureDrownedMonk  = "creature-drowned-monk.jpg"
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
				{
					Reference: DemoLocSacristyRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Sacristy",
						Description: "A narrow preparation room off the underground chapel. Tarnished vestment chests line the walls beside empty candle racks. A shallow font of holy water sits in a stone niche by the door — somehow still brimming.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocSacristy},
				},
				{
					Reference: DemoLocOssuaryRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Ossuary",
						Description: "Walls of neatly stacked bones recede into the gloom. Niches hold skulls arranged in silent rows. A runic circle is carved into the floor, its lines faintly luminous. An iron reliquary rests on a low plinth, its lock set with a mechanism of unusual complexity.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocOssuary},
				},
				{
					Reference: DemoLocInfirmaryRef,
					Record: &adventure_game_record.AdventureGameLocation{
						Name:        "The Infirmary",
						Description: "Cracked plaster walls and bare wooden cots mark the abbey's old healing ward. Shelves of dried herbs and crumbled tincture bottles line one wall. A drying rack of bundled herbs hangs over a stone hearth that has not been lit in decades.",
					},
					BackgroundImage: &harness.GameImageConfig{ImagePath: ImageLocInfirmary},
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
				// ── New items: weapons, armour, consumables, misc ─────────
				{
					Reference: DemoItemRustyDaggerRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:           "Rusty Dagger",
						Description:    "A short-bladed knife, its edge pitted with rust. It will do in a pinch.",
						CanBeEquipped:  true,
						ItemCategory:   convert.PtrStrict("weapon"),
						EquipmentSlot:  convert.PtrStrict(adventure_game_record.AdventureGameItemEquipmentSlotWeapon),
						IsStartingItem: true,
					},
					AdventureGameItemEffectConfigs: []harness.AdventureGameItemEffectConfig{
						{
							Reference: "demo-item-effect-rusty-dagger-weapon-damage",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeEquip,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeWeaponDamage,
								ResultDescription: "",
								ResultValueMin:    nullint32.FromInt32(3),
								ResultValueMax:    nullint32.FromInt32(6),
								IsRepeatable:      true,
							},
						},
						{
							Reference: "demo-item-effect-rusty-dagger-inspect",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeInfo,
								ResultDescription: "A battered iron dagger. Rust pits the blade but the edge still holds.",
								IsRepeatable:      true,
							},
						},
					},
				},
				{
					Reference: DemoItemMonksIronMaceRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:          "Monk's Iron Mace",
						Description:   "A heavy iron mace wrapped in worn leather at the grip. The monks apparently kept one on hand for reasons the manuscripts do not explain.",
						CanBeEquipped: true,
						ItemCategory:  convert.PtrStrict("weapon"),
						EquipmentSlot: convert.PtrStrict(adventure_game_record.AdventureGameItemEquipmentSlotWeapon),
					},
					AdventureGameItemEffectConfigs: []harness.AdventureGameItemEffectConfig{
						{
							Reference: "demo-item-effect-monks-mace-weapon-damage",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeEquip,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeWeaponDamage,
								ResultDescription: "",
								ResultValueMin:    nullint32.FromInt32(8),
								ResultValueMax:    nullint32.FromInt32(14),
								IsRepeatable:      true,
							},
						},
						{
							Reference: "demo-item-effect-monks-mace-inspect",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeInfo,
								ResultDescription: "A solid iron mace, heavy in the hand. Whatever the monks used it for, it would crack bone easily.",
								IsRepeatable:      true,
							},
						},
					},
				},
				{
					Reference: DemoItemLeatherCuirassRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:          "Leather Cuirass",
						Description:   "A stiffened leather breastplate, cracked with age but still serviceable. Stitching marks where it has been repaired more than once.",
						CanBeEquipped: true,
						ItemCategory:  convert.PtrStrict("armor"),
						EquipmentSlot: convert.PtrStrict(adventure_game_record.AdventureGameItemEquipmentSlotArmor),
					},
					AdventureGameItemEffectConfigs: []harness.AdventureGameItemEffectConfig{
						{
							Reference: "demo-item-effect-leather-cuirass-armor-defense",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeEquip,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeArmorDefense,
								ResultDescription: "",
								ResultValueMin:    nullint32.FromInt32(5),
								ResultValueMax:    nullint32.FromInt32(5),
								IsRepeatable:      true,
							},
						},
						{
							Reference: "demo-item-effect-leather-cuirass-inspect",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeInfo,
								ResultDescription: "Old leather armour, stiff but intact. It would blunt a blade or claw.",
								IsRepeatable:      true,
							},
						},
					},
				},
				{
					Reference: DemoItemHealingDraughtRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:         "Healing Draught",
						Description:  "A small clay vial sealed with wax. The liquid inside smells of herbs and honey. One dose remains.",
						ItemCategory: convert.PtrStrict("consumable"),
					},
					AdventureGameItemEffectConfigs: []harness.AdventureGameItemEffectConfig{
						{
							Reference: "demo-item-effect-healing-draught-use",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeHealWielder,
								ResultDescription: "You drink the draught. Warmth spreads through your chest and your wounds close slightly.",
								ResultValueMin:    nullint32.FromInt32(15),
								ResultValueMax:    nullint32.FromInt32(25),
								IsRepeatable:      false,
							},
						},
						{
							Reference: "demo-item-effect-healing-draught-inspect",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeInfo,
								ResultDescription: "A clay vial of healing tincture. A single use remains.",
								IsRepeatable:      true,
							},
						},
					},
				},
				{
					Reference: DemoItemBrassThuriblRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:         "Brass Thurible",
						Description:  "A brass incense burner on a chain, tarnished with age. Unburnt incense still sits in the bowl. Something feels faintly wrong about it.",
						ItemCategory: convert.PtrStrict("misc"),
					},
					AdventureGameItemEffectConfigs: []harness.AdventureGameItemEffectConfig{
						{
							Reference: "demo-item-effect-brass-thurible-use-damage",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeDamageWielder,
								ResultDescription: "The thurible swings erratically and strikes you. A cold pain flares — the incense smoke carries a curse.",
								ResultValueMin:    nullint32.FromInt32(5),
								ResultValueMax:    nullint32.FromInt32(10),
								IsRepeatable:      true,
							},
						},
						{
							Reference:     "demo-item-effect-brass-thurible-use-open-link",
							ResultLinkRef: "demo-link-flooded-to-ossuary",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeOpenLink,
								ResultDescription: "As the cursed smoke clears the portcullis ahead groans and slowly rises.",
								IsRepeatable:      true,
							},
						},
						{
							Reference: "demo-item-effect-brass-thurible-inspect",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeInfo,
								ResultDescription: "A tarnished brass thurible. The unburnt incense smells of something bitter. Using it might have consequences.",
								IsRepeatable:      true,
							},
						},
					},
				},
				{
					Reference: DemoItemTarnishedLocketRef,
					Record: &adventure_game_record.AdventureGameItem{
						Name:          "Tarnished Locket",
						Description:   "A small silver locket on a chain, blackened with age. It opens to reveal a tiny portrait, the face worn to nothing.",
						CanBeEquipped: true,
						ItemCategory:  convert.PtrStrict("jewelry"),
						EquipmentSlot: convert.PtrStrict(adventure_game_record.AdventureGameItemEquipmentSlotJewelry),
					},
					AdventureGameItemEffectConfigs: []harness.AdventureGameItemEffectConfig{
						{
							Reference: "demo-item-effect-tarnished-locket-equip",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeEquip,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeInfo,
								ResultDescription: "As you clasp the locket around your neck, whispers brush the edges of your mind — voices from long ago, fading before you can make out the words.",
								IsRepeatable:      true,
							},
						},
						{
							Reference: "demo-item-effect-tarnished-locket-inspect",
							Record: &adventure_game_record.AdventureGameItemEffect{
								ActionType:        adventure_game_record.AdventureGameItemEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameItemEffectEffectTypeInfo,
								ResultDescription: "A small locket, its surface black with tarnish. The portrait inside is too worn to identify.",
								IsRepeatable:      true,
							},
						},
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
				// ── New creatures ──────────────────────────────────────────
				{
					Reference: DemoCreatureCryptSpiderRef,
					Record: &adventure_game_record.AdventureGameCreature{
						Name:              "Crypt Spider",
						Description:       "A pale, bloated spider the size of a cat. It moves in slow deliberate circles, its eight milky eyes tracking everything that enters its territory.",
						MaxHealth:         30,
						AttackDamage:      8,
						Defense:           2,
						Disposition:       adventure_game_record.AdventureGameCreatureDispositionIndifferent,
						AttackMethod:      adventure_game_record.AdventureGameCreatureAttackMethodSting,
						AttackDescription: "drives a venom-tipped fang into your leg",
						BodyDecayTurns:    2,
						RespawnTurns:      4,
					},
					PortraitImage: &harness.GameImageConfig{ImagePath: ImageCreatureCryptSpider},
				},
				{
					Reference: DemoCreatureBoneRevenantRef,
					Record: &adventure_game_record.AdventureGameCreature{
						Name:              "Bone Revenant",
						Description:       "An animated skeleton clad in the rotted remnants of a monk's habit. Its movements are deliberate and eerily silent. Empty sockets glow with a dull red light.",
						MaxHealth:         60,
						AttackDamage:      12,
						Defense:           8,
						Disposition:       adventure_game_record.AdventureGameCreatureDispositionAggressive,
						AttackMethod:      adventure_game_record.AdventureGameCreatureAttackMethodClaws,
						AttackDescription: "rakes at you with bony fingers",
						BodyDecayTurns:    4,
						RespawnTurns:      0,
					},
					PortraitImage: &harness.GameImageConfig{ImagePath: ImageCreatureBoneRevenant},
				},
				{
					Reference: DemoCreatureDrownedMonkRef,
					Record: &adventure_game_record.AdventureGameCreature{
						Name:              "Drowned Monk",
						Description:       "A waterlogged corpse in a sodden habit, drifting through the flooded passage with slow purposeful steps. Water streams from its sleeves and cowl. Its face is pale and bloated.",
						MaxHealth:         50,
						AttackDamage:      10,
						Defense:           3,
						Disposition:       adventure_game_record.AdventureGameCreatureDispositionAggressive,
						AttackMethod:      adventure_game_record.AdventureGameCreatureAttackMethodSlam,
						AttackDescription: "hurls its sodden bulk into you with tremendous force",
						BodyDecayTurns:    3,
						RespawnTurns:      6,
					},
					PortraitImage: &harness.GameImageConfig{ImagePath: ImageCreatureDrownedMonk},
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
				// ── New links ─────────────────────────────────────────────
				// Underground Chapel -> Sacristy (no requirement)
				{
					Reference:       "demo-link-chapel-to-sacristy",
					FromLocationRef: DemoLocUndergroundChapelRef,
					ToLocationRef:   DemoLocSacristyRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:                 "The Vestry Door",
						Description:          "A low door set into the chapel's side wall opens onto the monks' preparation room.",
						TraversalDescription: nullstring.FromString("You push through the vestry door into the narrow sacristy beyond."),
					},
				},
				// Sacristy -> Underground Chapel (no requirement, return path)
				{
					Reference:       "demo-link-sacristy-to-chapel",
					FromLocationRef: DemoLocSacristyRef,
					ToLocationRef:   DemoLocUndergroundChapelRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:                 "The Vestry Door",
						Description:          "The vestry door leads back into the underground chapel.",
						TraversalDescription: nullstring.FromString("You step back through the vestry door into the flickering quiet of the chapel."),
					},
				},
				// Crypt -> Ossuary (requires Rusty Dagger equipped)
				{
					Reference:       "demo-link-crypt-to-ossuary",
					FromLocationRef: DemoLocCryptRef,
					ToLocationRef:   DemoLocOssuaryRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:              "The Bone Arch",
						Description:       "A low archway framed by fused bones leads deeper into the hill. A palpable unease hangs around it — something stirs on the other side.",
						LockedDescription: nullstring.FromString("The bone archway radiates menace. You feel certain that entering unarmed would be fatal."),
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "demo-link-req-crypt-to-ossuary",
							GameItemRef: DemoItemRustyDaggerRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
								Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionEquipped,
								Quantity:  1,
							},
						},
					},
				},
				// Ossuary -> Crypt (no requirement, return path)
				{
					Reference:       "demo-link-ossuary-to-crypt",
					FromLocationRef: DemoLocOssuaryRef,
					ToLocationRef:   DemoLocCryptRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:                 "The Bone Arch",
						Description:          "The bone archway leads back into the cold silence of the crypt.",
						TraversalDescription: nullstring.FromString("You pass back under the bone arch into the crypt's chill."),
					},
				},
				// Herb Garden -> Infirmary (no requirement)
				{
					Reference:       "demo-link-garden-to-infirmary",
					FromLocationRef: DemoLocHerbGardenRef,
					ToLocationRef:   DemoLocInfirmaryRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:                 "The Infirmary Door",
						Description:          "A heavy oak door in the garden wall leads into the abbey's old infirmary.",
						TraversalDescription: nullstring.FromString("You push through the heavy door into the dusty infirmary beyond."),
					},
				},
				// Infirmary -> Herb Garden (no requirement, return path)
				{
					Reference:       "demo-link-infirmary-to-garden",
					FromLocationRef: DemoLocInfirmaryRef,
					ToLocationRef:   DemoLocHerbGardenRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:                 "The Infirmary Door",
						Description:          "The heavy door leads back out into the walled herb garden.",
						TraversalDescription: nullstring.FromString("You step back out through the heavy door into the herb garden's grey daylight."),
					},
				},
				// Flooded Corridor -> Ossuary (portcullis link — starts locked, requires Brass
				// Thurible in inventory; the Portcullis Winch object can open it permanently)
				{
					Reference:       "demo-link-flooded-to-ossuary",
					FromLocationRef: DemoLocFloodedCorridorRef,
					ToLocationRef:   DemoLocOssuaryRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:              "The Portcullis",
						Description:       "A heavy iron portcullis bars a low archway at the end of the flooded corridor. Chains disappear into the ceiling above it.",
						LockedDescription: nullstring.FromString("The portcullis is sealed by a spiritual ward. You sense the incense from the thurible might break it — or find the winch mechanism."),
					},
					AdventureGameLocationLinkRequirementConfigs: []harness.AdventureGameLocationLinkRequirementConfig{
						{
							Reference:   "demo-link-req-flooded-to-ossuary",
							GameItemRef: DemoItemBrassThuriblRef,
							Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
								Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
								Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory,
								Quantity:  1,
							},
						},
					},
				},
				// Ossuary -> Flooded Corridor (no requirement, return path)
				{
					Reference:       "demo-link-ossuary-to-flooded",
					FromLocationRef: DemoLocOssuaryRef,
					ToLocationRef:   DemoLocFloodedCorridorRef,
					Record: &adventure_game_record.AdventureGameLocationLink{
						Name:                 "The Portcullis",
						Description:          "The iron portcullis leads back into the flooded corridor.",
						TraversalDescription: nullstring.FromString("You pass back through the portcullis into the cold flooded passage."),
					},
				},
			},

			// ── Location Objects ───────────────────────────────────────
			AdventureGameLocationObjectConfigs: []harness.AdventureGameLocationObjectConfig{

				// ── Hidden Door (created first so the Lever can reference its states) ──
				{
					Reference:       DemoObjHiddenDoorRef,
					LocationRef:     DemoLocUndergroundChapelRef,
					InitialStateRef: DemoObjHiddenDoorStateHiddenRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Hidden Door",
						Description: "A section of the chapel wall that seems to blend perfectly with the surrounding stonework. There is something off about the mortar lines.",
						IsHidden:    true,
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjHiddenDoorStateHiddenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "hidden",
								Description: "The door is hidden within the chapel wall.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjHiddenDoorStateRevealedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "revealed",
								Description: "The door's outline is now visible but not yet open.",
								SortOrder:   1,
							},
						},
						{
							Reference: DemoObjHiddenDoorStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "open",
								Description: "The door is open, revealing a dark shaft beyond.",
								SortOrder:   2,
							},
						},
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
							Reference:        "demo-obj-effect-hidden-door-push",
							RequiredStateRef: DemoObjHiddenDoorStateRevealedRef,
							ResultStateRef:   DemoObjHiddenDoorStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "The hidden door swings inward on well-balanced pivots, revealing a dark shaft beyond.",
								IsRepeatable:      false,
							},
						},
					},
				},

				// ── Ancient Well (Well Chamber) ──
				{
					Reference:       DemoObjAncientWellRef,
					LocationRef:     DemoLocWellChamberRef,
					InitialStateRef: DemoObjWellStateIntactRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Ancient Well",
						Description: "A circular stone well of great age. A frayed rope hangs into the darkness below. Water glints faintly far down.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjWellStateIntactRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "intact",
								Description: "The well is undisturbed, its secret hidden beneath the rim.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjWellStateSearchedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "searched",
								Description: "The well has been thoroughly searched. Nothing more remains.",
								SortOrder:   1,
							},
						},
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
							Reference:        "demo-obj-effect-well-search-give",
							RequiredStateRef: DemoObjWellStateIntactRef,
							ResultItemRef:    DemoItemAbbotsJournalRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
								ResultDescription: "Wedged beneath a loose stone on the well's inner rim, you find a cracked leather journal — the Abbot's Journal.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-well-search-state",
							RequiredStateRef: DemoObjWellStateIntactRef,
							ResultStateRef:   DemoObjWellStateSearchedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// search (state=searched) → info
						{
							Reference:        "demo-obj-effect-well-searched",
							RequiredStateRef: DemoObjWellStateSearchedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "You have already searched the well thoroughly. Nothing remains.",
								IsRepeatable:      true,
							},
						},
						// climb → info (repeatable)
						{
							Reference: "demo-obj-effect-well-climb",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeClimb,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "You grip the frayed rope and lower yourself a short distance into the shaft. Cold darkness yawns below. Dripping water echoes from far beneath. Prudence wins out and you climb back up.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Stone Altar (Underground Chapel) ──
				{
					Reference:       DemoObjStoneAltarRef,
					LocationRef:     DemoLocUndergroundChapelRef,
					InitialStateRef: DemoObjAltarStateBareRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Stone Altar",
						Description: "A low cracked altar of dark stone. Wax trails from ancient candles cross its surface. A shallow depression sits at its centre.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjAltarStateBareRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "bare",
								Description: "The altar's depression is empty, waiting for an offering.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjAltarStateBlessedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "blessed",
								Description: "The silver cross rests in the depression, the altar aglow with faint light.",
								SortOrder:   1,
							},
						},
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
							Reference:        "demo-obj-effect-altar-use-state",
							RequiredItemRef:  DemoItemSilverCrossRef,
							RequiredStateRef: DemoObjAltarStateBareRef,
							ResultStateRef:   DemoObjAltarStateBlessedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "You place the silver cross in the depression. The altar shudders. A low grinding resonates through the stone.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-altar-use-reveal",
							RequiredItemRef:  DemoItemSilverCrossRef,
							RequiredStateRef: DemoObjAltarStateBareRef,
							ResultObjectRef:  DemoObjHiddenDoorRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRevealObject,
								ResultDescription: "A section of the chapel wall shifts, revealing what appears to be a hidden door.",
								IsRepeatable:      false,
							},
						},
						// use with silver cross (bare) → remove_item: Silver Cross consumed
						{
							Reference:        "demo-obj-effect-altar-use-remove-cross",
							RequiredItemRef:  DemoItemSilverCrossRef,
							RequiredStateRef: DemoObjAltarStateBareRef,
							ResultItemRef:    DemoItemSilverCrossRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveItem,
								ResultDescription: "The silver cross sinks into the altar's depression, its warmth fading as the stone claims it. The offering is made.",
								IsRepeatable:      false,
							},
						},
						// use (state=blessed) → info
						{
							Reference:        "demo-obj-effect-altar-use-blessed",
							RequiredStateRef: DemoObjAltarStateBlessedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "The silver cross rests in the altar's depression, glowing faintly. Its work is done.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Bundle of Manuscripts (Abbot's Study) ──
				{
					Reference:       DemoObjManuscriptsRef,
					LocationRef:     DemoLocAbbotsStudyRef,
					InitialStateRef: DemoObjManuscriptsStateUnreadRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Bundle of Manuscripts",
						Description: "A stack of aged parchments tied with twine. The writing is cramped and faded, but still legible in places.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjManuscriptsStateUnreadRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "unread",
								Description: "The manuscripts have not yet been read.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjManuscriptsStateReadRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "read",
								Description: "The manuscripts have been read and their secrets known.",
								SortOrder:   1,
							},
						},
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
							Reference:        "demo-obj-effect-manuscripts-read-state",
							RequiredStateRef: DemoObjManuscriptsStateUnreadRef,
							ResultStateRef:   DemoObjManuscriptsStateReadRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeRead,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "You read through the manuscripts carefully. They describe a ritual of sealing, performed beneath the chapel floor. The key, they say, is \"faith made solid\" — placed upon the altar.",
								IsRepeatable:      false,
							},
						},
						// read (state=read) → info (repeatable)
						{
							Reference:        "demo-obj-effect-manuscripts-reread",
							RequiredStateRef: DemoObjManuscriptsStateReadRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeRead,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "You reread the manuscripts. The ritual of sealing requires \"faith made solid\" placed upon the altar. The words do not change.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Iron-Bound Chest (Bell Tower Vault) ──
				{
					Reference:       DemoObjIronChestRef,
					LocationRef:     DemoLocBellTowerVaultRef,
					InitialStateRef: DemoObjChestStateLockedRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Iron-Bound Chest",
						Description: "A heavy oak chest reinforced with iron bands. A large padlock seals it shut.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjChestStateLockedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "locked",
								Description: "The chest is secured by a heavy iron padlock.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjChestStateUnlockedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "unlocked",
								Description: "The padlock has been removed but the chest lid is still closed.",
								SortOrder:   1,
							},
						},
						{
							Reference: DemoObjChestStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "open",
								Description: "The chest lid is raised and its contents are exposed.",
								SortOrder:   2,
							},
						},
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
							Reference:        "demo-obj-effect-chest-unlock",
							RequiredItemRef:  DemoItemRustyKeyRef,
							RequiredStateRef: DemoObjChestStateLockedRef,
							ResultStateRef:   DemoObjChestStateUnlockedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUnlock,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "The rusty key grates in the lock, then catches. With a clunk the padlock falls open.",
								IsRepeatable:      false,
							},
						},
						// open (state=unlocked) → give_item: Monk's Iron Mace, change_state to open
						{
							Reference:        "demo-obj-effect-chest-open-give",
							RequiredStateRef: DemoObjChestStateUnlockedRef,
							ResultItemRef:    DemoItemMonksIronMaceRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
								ResultDescription: "You lift the heavy lid. Inside, wrapped in rotted cloth, is a heavy iron mace — the kind a monastery would keep for emergencies.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-chest-open-state",
							RequiredStateRef: DemoObjChestStateUnlockedRef,
							ResultStateRef:   DemoObjChestStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// inspect (state=open) → info
						{
							Reference:        "demo-obj-effect-chest-inspect-open",
							RequiredStateRef: DemoObjChestStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "The chest stands open. The rotted cloth remains, but the iron mace is gone.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Rusted Gate (Crypt) ──
				{
					Reference:       DemoObjRustedGateRef,
					LocationRef:     DemoLocCryptRef,
					InitialStateRef: DemoObjGateStateClosedRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Rusted Gate",
						Description: "A heavy iron gate bars an alcove at the far end of the crypt. Rust has welded the hinges solid.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjGateStateClosedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "closed",
								Description: "The gate is rusted shut, blocking the alcove.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjGateStateBrokenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "broken",
								Description: "The gate lies twisted on the floor. The alcove beyond is now accessible.",
								SortOrder:   1,
							},
						},
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
							Reference:        "demo-obj-effect-gate-break-state",
							RequiredStateRef: DemoObjGateStateClosedRef,
							ResultStateRef:   DemoObjGateStateBrokenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBreak,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "You hurl yourself against the gate. With a shriek of tortured metal it tears free from the stone, but the jagged edge catches you as it falls.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-gate-break-damage",
							RequiredStateRef: DemoObjGateStateClosedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBreak,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeDamage,
								ResultDescription: "The gate's jagged edge opens a cut across your arm.",
								ResultValueMin:    nullint32.FromInt32(5),
								ResultValueMax:    nullint32.FromInt32(10),
								IsRepeatable:      false,
							},
						},
						// inspect (state=broken) → info
						{
							Reference:        "demo-obj-effect-gate-inspect-broken",
							RequiredStateRef: DemoObjGateStateBrokenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "The rusted gate lies twisted on the crypt floor. The alcove beyond it is now open.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Lever on the Wall (Well Chamber) — demonstrates cross-object reveal ──
				{
					Reference:       DemoObjLeverRef,
					LocationRef:     DemoLocWellChamberRef,
					InitialStateRef: DemoObjLeverStateUpRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Lever on the Wall",
						Description: "An iron lever set into the chamber wall, its purpose obscure. It is in the UP position.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjLeverStateUpRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "up",
								Description: "The lever is in the up position.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjLeverStateDownRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "down",
								Description: "The lever is in the down position.",
								SortOrder:   1,
							},
						},
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
							Reference:        "demo-obj-effect-lever-pull-self",
							RequiredStateRef: DemoObjLeverStateUpRef,
							ResultStateRef:   DemoObjLeverStateDownRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "You pull the lever down. A deep mechanical rumble travels through the stone walls.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-lever-pull-change-door-state",
							RequiredStateRef: DemoObjLeverStateUpRef,
							ResultObjectRef:  DemoObjHiddenDoorRef,
							ResultStateRef:   DemoObjHiddenDoorStateRevealedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeObjectState,
								ResultDescription: "Somewhere above, a section of the chapel wall shifts with a grinding of old stone.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-lever-pull-reveal-door",
							RequiredStateRef: DemoObjLeverStateUpRef,
							ResultObjectRef:  DemoObjHiddenDoorRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRevealObject,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// push (state=down) → change_state to up, change_object_state: Hidden Door → hidden, hide_object: Hidden Door
						{
							Reference:        "demo-obj-effect-lever-push-self",
							RequiredStateRef: DemoObjLeverStateDownRef,
							ResultStateRef:   DemoObjLeverStateUpRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "You push the lever back up. The grinding rumble returns and then silence.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-lever-push-hide-door-state",
							RequiredStateRef: DemoObjLeverStateDownRef,
							ResultObjectRef:  DemoObjHiddenDoorRef,
							ResultStateRef:   DemoObjHiddenDoorStateHiddenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeObjectState,
								ResultDescription: "The section of chapel wall slides back into place.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-lever-push-hide-door",
							RequiredStateRef: DemoObjLeverStateDownRef,
							ResultObjectRef:  DemoObjHiddenDoorRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHideObject,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
					},
				},

				// ── Holy Water Font (Sacristy) — heal, pour ───────────────
				{
					Reference:       DemoObjHolyWaterFontRef,
					LocationRef:     DemoLocSacristyRef,
					InitialStateRef: DemoObjHolyWaterFontStateFlowingRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Holy Water Font",
						Description: "A shallow stone basin set into a niche in the sacristy wall. Clear water fills it almost to the brim, impossibly fresh for a room sealed for decades.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjHolyWaterFontStateFlowingRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "flowing",
								Description: "The font is full of clear, cool water.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjHolyWaterFontStateDryRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "dry",
								Description: "The font is empty, its miracle spent.",
								SortOrder:   1,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
						// inspect → info (repeatable)
						{
							Reference: "demo-obj-effect-font-inspect",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "The water is perfectly clear and cool. No visible source feeds it.",
								IsRepeatable:      true,
							},
						},
						// use (state=flowing) → heal 15-25, change_state to dry
						{
							Reference:        "demo-obj-effect-font-use-heal",
							RequiredStateRef: DemoObjHolyWaterFontStateFlowingRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHeal,
								ResultDescription: "You cup your hands and drink deep. Cool warmth spreads from your throat through your chest. Your wounds knit and your fatigue lifts.",
								ResultValueMin:    nullint32.FromInt32(15),
								ResultValueMax:    nullint32.FromInt32(25),
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-font-use-dry",
							RequiredStateRef: DemoObjHolyWaterFontStateFlowingRef,
							ResultStateRef:   DemoObjHolyWaterFontStateDryRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// pour (flowing) → info (flavour, repeatable)
						{
							Reference:        "demo-obj-effect-font-pour",
							RequiredStateRef: DemoObjHolyWaterFontStateFlowingRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePour,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "You cup a little water and let it run through your fingers. It is startlingly cold and clear.",
								IsRepeatable:      true,
							},
						},
						// use (state=dry) → info
						{
							Reference:        "demo-obj-effect-font-use-dry-state",
							RequiredStateRef: DemoObjHolyWaterFontStateDryRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "The font is empty. The blessing is spent.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Vestment Chest (Sacristy) — give_item: Leather Cuirass ──
				{
					Reference:       DemoObjVestmentChestRef,
					LocationRef:     DemoLocSacristyRef,
					InitialStateRef: DemoObjVestmentChestStateClosedRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Vestment Chest",
						Description: "A long wooden chest with tarnished brass fittings. It holds the sacristy's stored vestments — or held them.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjVestmentChestStateClosedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "closed",
								Description: "The chest is closed.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjVestmentChestStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "open",
								Description: "The chest stands open.",
								SortOrder:   1,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
						// inspect → info
						{
							Reference: "demo-obj-effect-vestment-chest-inspect",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "A long wooden chest, its lid slightly warped. It does not appear to be locked.",
								IsRepeatable:      true,
							},
						},
						// open (closed) → give_item: Leather Cuirass, change_state to open
						{
							Reference:        "demo-obj-effect-vestment-chest-open-give",
							RequiredStateRef: DemoObjVestmentChestStateClosedRef,
							ResultItemRef:    DemoItemLeatherCuirassRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
								ResultDescription: "You lift the lid. Most of the vestments have rotted to dust, but underneath you find a stiffened leather cuirass, still sound.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-vestment-chest-open-state",
							RequiredStateRef: DemoObjVestmentChestStateClosedRef,
							ResultStateRef:   DemoObjVestmentChestStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// inspect (open) → info
						{
							Reference:        "demo-obj-effect-vestment-chest-inspect-open",
							RequiredStateRef: DemoObjVestmentChestStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "The chest stands open, its rotted contents scattered. Nothing remains of use.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Apothecary Cabinet (Infirmary) — give_item: Healing Draught ──
				{
					Reference:       DemoObjApothecaryRef,
					LocationRef:     DemoLocInfirmaryRef,
					InitialStateRef: DemoObjApothecaryStateUndisturbedRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Apothecary Cabinet",
						Description: "A tall wooden cabinet with small labelled drawers. Most are empty or contain crumbled dried matter, but a few clay vials remain sealed.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjApothecaryStateUndisturbedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "undisturbed",
								Description: "The cabinet's contents are undisturbed.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjApothecaryStateSearchedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "searched",
								Description: "The cabinet has been thoroughly searched.",
								SortOrder:   1,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
						// inspect → info
						{
							Reference: "demo-obj-effect-apothecary-inspect",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "Rows of small drawers, most empty. A few sealed clay vials sit on the top shelf.",
								IsRepeatable:      true,
							},
						},
						// search (undisturbed) → give_item: Healing Draught, change_state to searched
						{
							Reference:        "demo-obj-effect-apothecary-search-give",
							RequiredStateRef: DemoObjApothecaryStateUndisturbedRef,
							ResultItemRef:    DemoItemHealingDraughtRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
								ResultDescription: "Among the crumbled contents you find a sealed clay vial — a healing draught, still good.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-apothecary-search-state",
							RequiredStateRef: DemoObjApothecaryStateUndisturbedRef,
							ResultStateRef:   DemoObjApothecaryStateSearchedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// search (searched) → info
						{
							Reference:        "demo-obj-effect-apothecary-searched",
							RequiredStateRef: DemoObjApothecaryStateSearchedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "You have already searched the cabinet. Nothing useful remains.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Herb Drying Rack (Infirmary) — burn, place_item, remove_object ──
				{
					Reference:       DemoObjHerbDryingRackRef,
					LocationRef:     DemoLocInfirmaryRef,
					InitialStateRef: DemoObjHerbDryingRackStateIntactRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Herb Drying Rack",
						Description: "A wooden rack strung with bundles of dried herbs — lavender, rosemary, and others less familiar. It hangs above the cold hearth.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjHerbDryingRackStateIntactRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "intact",
								Description: "The drying rack hangs undisturbed above the hearth.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjHerbDryingRackStateBurnedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "burned",
								Description: "The rack is burned away. Only ash remains on the hearth.",
								SortOrder:   1,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
						// inspect → info
						{
							Reference: "demo-obj-effect-rack-inspect",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "Dry bundles of herbs hang from the rack. Something gleams faintly behind the bundles — an object lodged against the wall.",
								IsRepeatable:      true,
							},
						},
						// search (intact) → info (nothing loose, prompts burning)
						{
							Reference:        "demo-obj-effect-rack-search",
							RequiredStateRef: DemoObjHerbDryingRackStateIntactRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "You feel behind the bundles. Something is definitely lodged back there, but you can't reach it through the rack.",
								IsRepeatable:      true,
							},
						},
						// burn (intact) → change_state to burned
						{
							Reference:        "demo-obj-effect-rack-burn-state",
							RequiredStateRef: DemoObjHerbDryingRackStateIntactRef,
							ResultStateRef:   DemoObjHerbDryingRackStateBurnedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBurn,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "You light the dry herb bundles. They catch instantly, filling the room with bitter smoke. The rack collapses onto the hearth in a shower of ash.",
								IsRepeatable:      false,
							},
						},
						// burn (intact) → place_item: Brass Thurible at Infirmary
						{
							Reference:         "demo-obj-effect-rack-burn-place-item",
							RequiredStateRef:  DemoObjHerbDryingRackStateIntactRef,
							ResultItemRef:     DemoItemBrassThuriblRef,
							ResultLocationRef: DemoLocInfirmaryRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBurn,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypePlaceItem,
								ResultDescription: "As the smoke clears, a brass thurible lies on the hearth, revealed from where it had been hidden behind the rack.",
								IsRepeatable:      false,
							},
						},
						// burn (intact) → remove_object (the rack is consumed)
						{
							Reference:        "demo-obj-effect-rack-burn-remove-object",
							RequiredStateRef: DemoObjHerbDryingRackStateIntactRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBurn,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveObject,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
					},
				},

				// ── Trapped Reliquary (Ossuary) — disarm, damage, give_item ──
				{
					Reference:       DemoObjTrappedReliquaryRef,
					LocationRef:     DemoLocOssuaryRef,
					InitialStateRef: DemoObjReliquaryStateTrappedRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Trapped Reliquary",
						Description: "An iron reliquary on a low plinth, its lid held by a complex locking mechanism set with fine wire. The wire looks deliberately taut.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjReliquaryStateTrappedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "trapped",
								Description: "The reliquary is locked and a hidden trap mechanism is set.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjReliquaryStateDisarmedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "disarmed",
								Description: "The trap mechanism has been carefully neutralised. The reliquary is still locked.",
								SortOrder:   1,
							},
						},
						{
							Reference: DemoObjReliquaryStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "open",
								Description: "The reliquary stands open.",
								SortOrder:   2,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
						// inspect → info
						{
							Reference: "demo-obj-effect-reliquary-inspect",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "A taut wire runs along the lid's edge. Opening it without care would trigger something unpleasant. Disarming it first would be wise.",
								IsRepeatable:      true,
							},
						},
						// disarm (trapped) → change_state to disarmed
						{
							Reference:        "demo-obj-effect-reliquary-disarm",
							RequiredStateRef: DemoObjReliquaryStateTrappedRef,
							ResultStateRef:   DemoObjReliquaryStateDisarmedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeDisarm,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "Working carefully, you trace the wire to its anchor and sever it. The trap mechanism springs loose harmlessly. The reliquary can now be opened safely.",
								IsRepeatable:      false,
							},
						},
						// open (trapped) → damage 15-25, change_state to open
						{
							Reference:        "demo-obj-effect-reliquary-open-trapped-damage",
							RequiredStateRef: DemoObjReliquaryStateTrappedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeDamage,
								ResultDescription: "The wire snaps. A spring mechanism drives a short spike into your forearm.",
								ResultValueMin:    nullint32.FromInt32(15),
								ResultValueMax:    nullint32.FromInt32(25),
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-reliquary-open-trapped-give",
							RequiredStateRef: DemoObjReliquaryStateTrappedRef,
							ResultItemRef:    DemoItemTarnishedLocketRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
								ResultDescription: "Despite the pain, you find a small tarnished locket nestled inside on a bed of rotted silk.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-reliquary-open-trapped-state",
							RequiredStateRef: DemoObjReliquaryStateTrappedRef,
							ResultStateRef:   DemoObjReliquaryStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// open (disarmed) → give_item, change_state to open
						{
							Reference:        "demo-obj-effect-reliquary-open-disarmed-give",
							RequiredStateRef: DemoObjReliquaryStateDisarmedRef,
							ResultItemRef:    DemoItemTarnishedLocketRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
								ResultDescription: "You open the reliquary safely. Inside, on a bed of rotted silk, rests a small tarnished locket.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-reliquary-open-disarmed-state",
							RequiredStateRef: DemoObjReliquaryStateDisarmedRef,
							ResultStateRef:   DemoObjReliquaryStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// inspect (open) → info
						{
							Reference:        "demo-obj-effect-reliquary-inspect-open",
							RequiredStateRef: DemoObjReliquaryStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "The reliquary stands open and empty. The rotted silk lining remains.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Runic Circle (Ossuary) — touch, teleport ──────────────
				{
					Reference:       DemoObjRunicCircleRef,
					LocationRef:     DemoLocOssuaryRef,
					InitialStateRef: DemoObjRunicCircleStateDormantRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Runic Circle",
						Description: "A circle of interlocking runes carved deep into the stone floor. The lines glow with a faint cold light. The symbols are old — older than the abbey.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjRunicCircleStateDormantRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "dormant",
								Description: "The circle's light is dim, its power untapped.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjRunicCircleStateActiveRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "active",
								Description: "The runes blaze with cold white light. The circle hums with energy.",
								SortOrder:   1,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
						// inspect → info (repeatable)
						{
							Reference: "demo-obj-effect-circle-inspect",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "An ancient teleportation circle. Touch it to activate. Use it while active to be transported to the abbey entrance.",
								IsRepeatable:      true,
							},
						},
						// touch (dormant) → change_state to active
						{
							Reference:        "demo-obj-effect-circle-touch-activate",
							RequiredStateRef: DemoObjRunicCircleStateDormantRef,
							ResultStateRef:   DemoObjRunicCircleStateActiveRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeTouch,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "You touch the outermost rune. The circle blazes to life, cold light flooding the chamber.",
								IsRepeatable:      false,
							},
						},
						// use (active) → teleport to Grand Staircase
						{
							Reference:         "demo-obj-effect-circle-use-teleport",
							RequiredStateRef:  DemoObjRunicCircleStateActiveRef,
							ResultLocationRef: DemoLocGrandStaircaseRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeTeleport,
								ResultDescription: "You step into the circle. The light engulfs you. When it fades you are standing in the abbey entrance hall, at the foot of the grand staircase.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Cursed Sarcophagus (Ossuary) — summon_creature ────────
				{
					Reference:       DemoObjCursedSarcophagusRef,
					LocationRef:     DemoLocOssuaryRef,
					InitialStateRef: DemoObjSarcophagusStateSealedRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Cursed Sarcophagus",
						Description: "A stone sarcophagus set apart from the bone walls, sealed with iron clamps and warded with the same runes as the circle. Something shifts inside when you move close.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjSarcophagusStateSealedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "sealed",
								Description: "The sarcophagus is sealed. The clamps hold firm.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjSarcophagusStateOpenedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "opened",
								Description: "The sarcophagus is open. Whatever was inside has stirred.",
								SortOrder:   1,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
						// inspect → info
						{
							Reference: "demo-obj-effect-sarcophagus-inspect",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "The iron clamps are sealed tight. Something scrapes against the stone from within. Opening this would almost certainly release whatever is inside.",
								IsRepeatable:      true,
							},
						},
						// open (sealed) → summon_creature: Bone Revenant
						{
							Reference:         "demo-obj-effect-sarcophagus-open-summon",
							RequiredStateRef:  DemoObjSarcophagusStateSealedRef,
							ResultCreatureRef: DemoCreatureBoneRevenantRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeSummonCreature,
								ResultDescription: "The clamps groan and burst free. The lid crashes to the floor. A skeletal figure in a monk's habit rises from within and turns toward you.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-sarcophagus-open-state",
							RequiredStateRef: DemoObjSarcophagusStateSealedRef,
							ResultStateRef:   DemoObjSarcophagusStateOpenedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// inspect (opened) → info
						{
							Reference:        "demo-obj-effect-sarcophagus-inspect-opened",
							RequiredStateRef: DemoObjSarcophagusStateOpenedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "The empty sarcophagus lies open. Whatever animated its occupant is gone.",
								IsRepeatable:      true,
							},
						},
					},
				},

				// ── Portcullis Winch (Flooded Corridor) — open_link, close_link ──
				{
					Reference:       DemoObjPortcullisWinchRef,
					LocationRef:     DemoLocFloodedCorridorRef,
					InitialStateRef: DemoObjPortcullisWinchStateUpRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        "Portcullis Winch",
						Description: "A heavy iron winch set into the wall of the flooded corridor. A thick chain disappears upward through a hole in the ceiling. The winch handle is in the UP position — portcullis down.",
					},
					AdventureGameLocationObjectStateConfigs: []harness.AdventureGameLocationObjectStateConfig{
						{
							Reference: DemoObjPortcullisWinchStateUpRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "up",
								Description: "The winch is up. The portcullis is lowered and the passage is sealed.",
								SortOrder:   0,
							},
						},
						{
							Reference: DemoObjPortcullisWinchStateDownRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "down",
								Description: "The winch is down. The portcullis is raised and the passage beyond is open.",
								SortOrder:   1,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []harness.AdventureGameLocationObjectEffectConfig{
						// inspect → info
						{
							Reference: "demo-obj-effect-winch-inspect",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								ResultDescription: "A heavy chain-and-ratchet mechanism. Pull the handle down to raise the portcullis. Push it back up to lower it again.",
								IsRepeatable:      true,
							},
						},
						// pull (up) → open_link: Flooded Corridor → Ossuary, change_state to down
						{
							Reference:        "demo-obj-effect-winch-pull-open",
							RequiredStateRef: DemoObjPortcullisWinchStateUpRef,
							ResultLinkRef:    "demo-link-flooded-to-ossuary",
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeOpenLink,
								ResultDescription: "You haul the winch handle down with both hands. The chain rattles and the portcullis grinds upward, opening the passage beyond.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-winch-pull-state",
							RequiredStateRef: DemoObjPortcullisWinchStateUpRef,
							ResultStateRef:   DemoObjPortcullisWinchStateDownRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
						// push (down) → close_link: Flooded Corridor → Ossuary (requires Brass
						// Thurible to re-enter — demonstrates close_link mechanic),
						// change_state to up
						{
							Reference:        "demo-obj-effect-winch-push-close",
							RequiredStateRef: DemoObjPortcullisWinchStateDownRef,
							ResultLinkRef:    "demo-link-flooded-to-ossuary",
							ResultItemRef:    DemoItemBrassThuriblRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeCloseLink,
								ResultDescription: "You push the winch back up. The portcullis crashes down, sealing the passage. The ward snaps back into place — the thurible would break it again.",
								IsRepeatable:      false,
							},
						},
						{
							Reference:        "demo-obj-effect-winch-push-state",
							RequiredStateRef: DemoObjPortcullisWinchStateDownRef,
							ResultStateRef:   DemoObjPortcullisWinchStateUpRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePush,
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								ResultDescription: "",
								IsRepeatable:      false,
							},
						},
					},
				},
			},

			// ── Creature Placements ────────────────────────────────────
			AdventureGameCreaturePlacementConfigs: []harness.AdventureGameCreaturePlacementConfig{
				{
					Reference:       "demo-placement-cellar-rat",
					GameCreatureRef: DemoCreatureCellarRatRef,
					GameLocationRef: DemoLocWineCellarRef,
					InitialCount:    1,
					Record:          &adventure_game_record.AdventureGameCreaturePlacement{},
				},
				{
					Reference:       "demo-placement-shadow-monk",
					GameCreatureRef: DemoCreatureShadowMonkRef,
					GameLocationRef: DemoLocCryptRef,
					InitialCount:    1,
					Record:          &adventure_game_record.AdventureGameCreaturePlacement{},
				},
				{
					Reference:       "demo-placement-drowned-monk",
					GameCreatureRef: DemoCreatureDrownedMonkRef,
					GameLocationRef: DemoLocFloodedCorridorRef,
					InitialCount:    1,
					Record:          &adventure_game_record.AdventureGameCreaturePlacement{},
				},
				{
					Reference:       "demo-placement-crypt-spider",
					GameCreatureRef: DemoCreatureCryptSpiderRef,
					GameLocationRef: DemoLocInfirmaryRef,
					InitialCount:    1,
					Record:          &adventure_game_record.AdventureGameCreaturePlacement{},
				},
			},

			// ── Item Placements ────────────────────────────────────────
			AdventureGameItemPlacementConfigs: []harness.AdventureGameItemPlacementConfig{
				// Rusty Key found in Herb Garden — opens Small Door and Iron-Bound Chest
				{
					Reference:       "demo-placement-rusty-key",
					GameItemRef:     DemoItemRustyKeyRef,
					GameLocationRef: DemoLocHerbGardenRef,
					InitialCount:    1,
					Record:          &adventure_game_record.AdventureGameItemPlacement{},
				},
				// Tallow Candle found in Narrow Passage — needed for Crypt and Flooded Corridor
				{
					Reference:       "demo-placement-tallow-candle",
					GameItemRef:     DemoItemTallowCandleRef,
					GameLocationRef: DemoLocNarrowPassageRef,
					InitialCount:    1,
					Record:          &adventure_game_record.AdventureGameItemPlacement{},
				},
				// Coil of Rope found in Wine Cellar — needed for Flooded Corridor → Well Chamber
				{
					Reference:       "demo-placement-coil-of-rope",
					GameItemRef:     DemoItemCoilOfRopeRef,
					GameLocationRef: DemoLocWineCellarRef,
					InitialCount:    1,
					Record:          &adventure_game_record.AdventureGameItemPlacement{},
				},
				// Silver Cross found in Crypt — needed for Warded Gate and Stone Altar ritual
				{
					Reference:       "demo-placement-silver-cross",
					GameItemRef:     DemoItemSilverCrossRef,
					GameLocationRef: DemoLocCryptRef,
					InitialCount:    1,
					Record:          &adventure_game_record.AdventureGameItemPlacement{},
				},
			},
		},
	}
}
