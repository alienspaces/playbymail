package demo_scenarios

import (
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

const (
	DemoMechaImageJoinGameRef       = "demo-mecha-image-join-game"
	DemoMechaImageOrdersRef         = "demo-mecha-image-orders"
	DemoMechaImageManagementRef     = "demo-mecha-image-management"

	ImageMechaJoinGame = "mecha-join-game.jpg"
	ImageMechaOrders   = "mecha-orders.jpg"
)

const (
	DemoMechaGameName = "Operation Scorched Ridge"
	DemoMechaGameRef  = "demo-mecha-game"

	DemoMechaInstanceRef = "demo-mecha-instance-one"

	// Chassis refs
	DemoMechChassisViperRef    = "demo-mech-chassis-viper"
	DemoMechChassisHornetRef   = "demo-mech-chassis-hornet"
	DemoMechChassisRangerRef   = "demo-mech-chassis-ranger"
	DemoMechChassisWardenRef   = "demo-mech-chassis-warden"
	DemoMechChassisCrusherRef  = "demo-mech-chassis-crusher"
	DemoMechChassisTitanRef    = "demo-mech-chassis-titan"

	// Weapon refs
	DemoMechWeaponLightPulseRef   = "demo-mech-weapon-light-pulse-cannon"
	DemoMechWeaponPulseRef        = "demo-mech-weapon-pulse-cannon"
	DemoMechWeaponHeavyPulseRef   = "demo-mech-weapon-heavy-pulse-cannon"
	DemoMechWeaponPlasmaRef       = "demo-mech-weapon-plasma-accelerator"
	DemoMechWeaponRocketPackRef   = "demo-mech-weapon-rocket-pack"
	DemoMechWeaponMissileRef      = "demo-mech-weapon-missile-battery"
	DemoMechWeaponRotaryCannonRef = "demo-mech-weapon-rotary-cannon"
	DemoMechWeaponChaingunRef     = "demo-mech-weapon-chaingun"

	// Sector refs
	DemoMechSectorDropzoneRef    = "demo-mech-sector-dropzone"
	DemoMechSectorRidgeNorthRef  = "demo-mech-sector-ridge-north"
	DemoMechSectorRidgeSouthRef  = "demo-mech-sector-ridge-south"
	DemoMechSectorValleyRef      = "demo-mech-sector-valley"
	DemoMechSectorRefineryRef    = "demo-mech-sector-refinery"
	DemoMechSectorCrossroadsRef  = "demo-mech-sector-crossroads"
	DemoMechSectorForestRef      = "demo-mech-sector-forest"
	DemoMechSectorCitadelRef     = "demo-mech-sector-citadel"

	// Computer opponent refs
	DemoMechaComputerOpponentRef      = "demo-mecha-computer-opponent-garrison"
	DemoMechaComputerOpponentLanceRef = "demo-mecha-computer-opponent-lance"
	DemoMechaComputerOpponentMech1Ref = "demo-mecha-computer-opponent-mech-1"
	DemoMechaComputerOpponentMech2Ref = "demo-mecha-computer-opponent-mech-2"

	// Player starter lance refs
	DemoMechaPlayerStarterLanceRef  = "demo-mecha-player-starter-lance"
	DemoMechaPlayerStarterMech1Ref  = "demo-mecha-player-starter-mech-1"
	DemoMechaPlayerStarterMech2Ref  = "demo-mecha-player-starter-mech-2"

	// Sector link refs
	DemoMechLinkDropzoneToRidgeNorthRef  = "demo-mech-link-dropzone-to-ridge-north"
	DemoMechLinkRidgeNorthToDropzoneRef  = "demo-mech-link-ridge-north-to-dropzone"
	DemoMechLinkDropzoneToValleyRef      = "demo-mech-link-dropzone-to-valley"
	DemoMechLinkValleyToDropzoneRef      = "demo-mech-link-valley-to-dropzone"
	DemoMechLinkRidgeNorthToRidgeSouthRef = "demo-mech-link-ridge-north-to-ridge-south"
	DemoMechLinkRidgeSouthToRidgeNorthRef = "demo-mech-link-ridge-south-to-ridge-north"
	DemoMechLinkRidgeSouthToRefineryRef  = "demo-mech-link-ridge-south-to-refinery"
	DemoMechLinkRefineryToRidgeSouthRef  = "demo-mech-link-refinery-to-ridge-south"
	DemoMechLinkValleyToCrossroadsRef    = "demo-mech-link-valley-to-crossroads"
	DemoMechLinkCrossroadsToValleyRef    = "demo-mech-link-crossroads-to-valley"
	DemoMechLinkCrossroadsToForestRef    = "demo-mech-link-crossroads-to-forest"
	DemoMechLinkForestToCrossroadsRef    = "demo-mech-link-forest-to-crossroads"
	DemoMechLinkRefinerytoCitadelRef     = "demo-mech-link-refinery-to-citadel"
	DemoMechLinkCitadelToRefineryRef     = "demo-mech-link-citadel-to-refinery"
	DemoMechLinkForestToCitadelRef       = "demo-mech-link-forest-to-citadel"
	DemoMechLinkCitadelToForestRef       = "demo-mech-link-citadel-to-forest"

)

// MechaGameConfig returns a standalone demo scenario for the mecha type,
// showcasing all designer-facing content: chassis, weapons, sectors, and sector links.
// Lances are player-owned and would be created when players subscribe.
// Accounts are managed by the CLI: subscription[0] uses demoRecs.AccountUsers[0] (designer),
// subscription[1] uses [1] (manager).
func MechaGameConfig() harness.DataConfig {
	return harness.DataConfig{
		GameConfigs: mechaGameConfigs(),
		AccountUserGameSubscriptionConfigs: []harness.AccountUserGameSubscriptionConfig{
			{
				Reference:        DemoSubscriptionDesignerTwoRef,
				GameRef:          DemoMechaGameRef,
				SubscriptionType: game_record.GameSubscriptionTypeDesigner,
				Record:           &game_record.GameSubscription{},
			},
			{
				Reference:        DemoSubscriptionManagerTwoRef,
				GameRef:          DemoMechaGameRef,
				SubscriptionType: game_record.GameSubscriptionTypeManager,
				Record:           &game_record.GameSubscription{},
				GameInstanceConfigs: []harness.GameInstanceConfig{
					{
						Reference: DemoMechaInstanceRef,
					Record: &game_record.GameInstance{
						DeliveryEmail:           true,
						TurnDurationHours:       168,
						RequiredPlayerCount:     1,
						ProcessWhenAllSubmitted: true,
					},
					},
				},
			},
		},
	}
}

func mechaGameConfigs() []harness.GameConfig {
	return []harness.GameConfig{
		{
			Reference: DemoMechaGameRef,
			Record: &game_record.Game{
				Name:              DemoMechaGameName,
				Description:       "Two rival commands clash over the strategic Scorched Ridge industrial complex. Command lances of war mechs across eight contested sectors — from the rugged northern ridges to the fortified citadel at the heart of the complex. Manage heat, exploit terrain cover, and outmanoeuvre your opponent to seize and hold the industrial prize. Every decision shapes the battle. Engage!",
			GameType:          game_record.GameTypeMecha,
			TurnDurationHours: 168,
			Status:            game_record.GameStatusDraft,
		},
		GameImageConfigs: []harness.GameImageConfig{
			{
				Reference:     DemoMechaImageJoinGameRef,
				ImagePath:     ImageMechaJoinGame,
				TurnSheetType: mecha_record.MechaTurnSheetTypeJoinGame,
			},
			{
				Reference:     DemoMechaImageOrdersRef,
				ImagePath:     ImageMechaOrders,
				TurnSheetType: mecha_record.MechaTurnSheetTypeOrders,
			},
			{
				Reference:     DemoMechaImageManagementRef,
				ImagePath:     ImageMechaOrders,
				TurnSheetType: mecha_record.MechaTurnSheetTypeLanceManagement,
			},
		},
		MechaChassisConfigs: []harness.MechaChassisConfig{
				{
					Reference: DemoMechChassisViperRef,
					Record: &mecha_record.MechaChassis{
						Name:            "Viper",
						Description:     "An ultra-light recon mech built for speed. Its thin armour means it cannot stand and trade blows, but nothing in its weight class can catch it.",
						ChassisClass:    mecha_record.ChassisClassLight,
						ArmorPoints:     56,
						StructurePoints: 24,
						HeatCapacity:    16,
						Speed:           8,
					},
				},
				{
					Reference: DemoMechChassisHornetRef,
					Record: &mecha_record.MechaChassis{
						Name:            "Hornet",
						Description:     "A fast light mech equipped for both scouting and raiding. Its agility lets it traverse difficult terrain with ease.",
						ChassisClass:    mecha_record.ChassisClassLight,
						ArmorPoints:     72,
						StructurePoints: 32,
						HeatCapacity:    20,
						Speed:           7,
					},
				},
				{
					Reference: DemoMechChassisRangerRef,
					Record: &mecha_record.MechaChassis{
						Name:            "Ranger",
						Description:     "A balanced medium mech. Its missile battery and plasma accelerator give it strong fire support capability while it still moves faster than most heavies.",
						ChassisClass:    mecha_record.ChassisClassMedium,
						ArmorPoints:     120,
						StructurePoints: 60,
						HeatCapacity:    28,
						Speed:           5,
					},
				},
				{
					Reference: DemoMechChassisWardenRef,
					Record: &mecha_record.MechaChassis{
						Name:            "Warden",
						Description:     "A hard-hitting medium mech that can mix it up in any range bracket. Its rotary cannon, rocket pack, and missile battery cover every engagement distance.",
						ChassisClass:    mecha_record.ChassisClassMedium,
						ArmorPoints:     136,
						StructurePoints: 68,
						HeatCapacity:    30,
						Speed:           5,
					},
				},
				{
					Reference: DemoMechChassisCrusherRef,
					Record: &mecha_record.MechaChassis{
						Name:            "Crusher",
						Description:     "A feared heavy mech whose twin plasma accelerators and rotary cannon give it devastating long-range punch. Opponents ignore the Crusher at their peril.",
						ChassisClass:    mecha_record.ChassisClassHeavy,
						ArmorPoints:     200,
						StructurePoints: 100,
						HeatCapacity:    38,
						Speed:           4,
					},
				},
				{
					Reference: DemoMechChassisTitanRef,
					Record: &mecha_record.MechaChassis{
						Name:            "Titan",
						Description:     "The most fearsome assault mech ever deployed. Its staggering armour and weaponry make it a walking fortress. Slow but effectively unstoppable.",
						ChassisClass:    mecha_record.ChassisClassAssault,
						ArmorPoints:     304,
						StructurePoints: 152,
						HeatCapacity:    42,
						Speed:           3,
					},
				},
			},
			MechaWeaponConfigs: []harness.MechaWeaponConfig{
				{
					Reference: DemoMechWeaponLightPulseRef,
					Record: &mecha_record.MechaWeapon{
						Name:        "Light Pulse Cannon",
						Description: "A compact close-range energy weapon used as a back-up or on light mechs with limited capacity.",
						Damage:      3,
						HeatCost:    1,
						RangeBand:   mecha_record.WeaponRangeBandShort,
						MountSize:   mecha_record.WeaponMountSizeSmall,
					},
				},
				{
					Reference: DemoMechWeaponPulseRef,
					Record: &mecha_record.MechaWeapon{
						Name:        "Pulse Cannon",
						Description: "The workhorse direct-fire energy weapon. Reliable, accurate, and found on mechs of every weight class.",
						Damage:      5,
						HeatCost:    3,
						RangeBand:   mecha_record.WeaponRangeBandMedium,
						MountSize:   mecha_record.WeaponMountSizeMedium,
					},
				},
				{
					Reference: DemoMechWeaponHeavyPulseRef,
					Record: &mecha_record.MechaWeapon{
						Name:        "Heavy Pulse Cannon",
						Description: "A powerful long-range energy weapon that deals serious damage but generates substantial heat.",
						Damage:      8,
						HeatCost:    8,
						RangeBand:   mecha_record.WeaponRangeBandLong,
						MountSize:   mecha_record.WeaponMountSizeLarge,
					},
				},
				{
					Reference: DemoMechWeaponPlasmaRef,
					Record: &mecha_record.MechaWeapon{
						Name:        "Plasma Accelerator",
						Description: "A heavy energy weapon firing superheated plasma bolts. Its combination of high damage and long range makes it a favourite of heavy and assault commanders.",
						Damage:      10,
						HeatCost:    10,
						RangeBand:   mecha_record.WeaponRangeBandLong,
						MountSize:   mecha_record.WeaponMountSizeLarge,
					},
				},
				{
					Reference: DemoMechWeaponRocketPackRef,
					Record: &mecha_record.MechaWeapon{
						Name:        "Rocket Pack",
						Description: "A short-range unguided rocket launcher ideal for close-in brawling. Each salvo can cripple a light mech outright.",
						Damage:      8,
						HeatCost:    3,
						RangeBand:   mecha_record.WeaponRangeBandShort,
						MountSize:   mecha_record.WeaponMountSizeMedium,
					},
				},
				{
					Reference: DemoMechWeaponMissileRef,
					Record: &mecha_record.MechaWeapon{
						Name:        "Missile Battery",
						Description: "A long-range guided missile launcher suited to indirect fire support. Can engage targets beyond visual range.",
						Damage:      10,
						HeatCost:    4,
						RangeBand:   mecha_record.WeaponRangeBandLong,
						MountSize:   mecha_record.WeaponMountSizeLarge,
					},
				},
				{
					Reference: DemoMechWeaponRotaryCannonRef,
					Record: &mecha_record.MechaWeapon{
						Name:        "Rotary Cannon",
						Description: "A medium rotary cannon firing bursts of armour-piercing rounds. Effective at medium range with manageable heat generation.",
						Damage:      5,
						HeatCost:    1,
						RangeBand:   mecha_record.WeaponRangeBandMedium,
						MountSize:   mecha_record.WeaponMountSizeLarge,
					},
				},
				{
					Reference: DemoMechWeaponChaingunRef,
					Record: &mecha_record.MechaWeapon{
						Name:        "Chaingun",
						Description: "A rapid-fire ballistic weapon devastating against light armour but of limited effect against heavy combat-mech plating.",
						Damage:      2,
						HeatCost:    0,
						RangeBand:   mecha_record.WeaponRangeBandShort,
						MountSize:   mecha_record.WeaponMountSizeSmall,
					},
				},
			},
			MechaSectorConfigs: []harness.MechaSectorConfig{
				{
					Reference: DemoMechSectorDropzoneRef,
					Record: &mecha_record.MechaSector{
						Name:             "Drop Zone Alpha",
						Description:      "The forward staging area where both commands begin the operation. Open ground, no cover — get moving fast.",
						TerrainType:      mecha_record.SectorTerrainTypeOpen,
						Elevation:        0,
						IsStartingSector: true,
					},
				},
				{
					Reference: DemoMechSectorRidgeNorthRef,
					Record: &mecha_record.MechaSector{
						Name:        "North Ridge",
						Description: "The northern flank of Scorched Ridge, a series of steep rocky outcrops that offer superb fields of fire into the valley below.",
						TerrainType: mecha_record.SectorTerrainTypeRough,
						Elevation:   4,
					},
				},
				{
					Reference: DemoMechSectorRidgeSouthRef,
					Record: &mecha_record.MechaSector{
						Name:        "South Ridge",
						Description: "The southern spur of the ridge, slightly lower than the north but still dominating the refinery approaches. Key high ground.",
						TerrainType: mecha_record.SectorTerrainTypeRough,
						Elevation:   3,
					},
				},
				{
					Reference: DemoMechSectorValleyRef,
					Record: &mecha_record.MechaSector{
						Name:        "Ash Valley",
						Description: "A broad shallow valley blanketed in volcanic ash. Visibility is poor and movement is slowed but the ground is flat.",
						TerrainType: mecha_record.SectorTerrainTypeOpen,
						Elevation:   0,
					},
				},
				{
					Reference: DemoMechSectorRefineryRef,
					Record: &mecha_record.MechaSector{
						Name:        "Refinery Complex",
						Description: "The industrial heart of the operation. Massive processing towers and pipework create a maze of cover — and fire hazards.",
						TerrainType: mecha_record.SectorTerrainTypeUrban,
						Elevation:   0,
					},
				},
				{
					Reference: DemoMechSectorCrossroadsRef,
					Record: &mecha_record.MechaSector{
						Name:        "Crossroads Junction",
						Description: "A vital road junction controlling access to every sector of the battlefield. Whoever holds Crossroads controls the tempo of the battle.",
						TerrainType: mecha_record.SectorTerrainTypeOpen,
						Elevation:   1,
					},
				},
				{
					Reference: DemoMechSectorForestRef,
					Record: &mecha_record.MechaSector{
						Name:        "Scorched Forest",
						Description: "Once dense woodland, now a skeletal tangle of burnt trees. Still provides moderate cover for approaching mechs.",
						TerrainType: mecha_record.SectorTerrainTypeForest,
						Elevation:   1,
					},
				},
				{
					Reference: DemoMechSectorCitadelRef,
					Record: &mecha_record.MechaSector{
						Name:        "Citadel Garrison",
						Description: "A heavily reinforced command bunker at the centre of the complex. The ultimate objective — taking or holding the Citadel wins the campaign.",
						TerrainType: mecha_record.SectorTerrainTypeUrban,
						Elevation:   2,
					},
				},
			},
			MechaSectorLinkConfigs: []harness.MechaSectorLinkConfig{
				// Drop Zone <-> North Ridge
				{
					Reference:     DemoMechLinkDropzoneToRidgeNorthRef,
					FromSectorRef: DemoMechSectorDropzoneRef,
					ToSectorRef:   DemoMechSectorRidgeNorthRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				{
					Reference:     DemoMechLinkRidgeNorthToDropzoneRef,
					FromSectorRef: DemoMechSectorRidgeNorthRef,
					ToSectorRef:   DemoMechSectorDropzoneRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				// Drop Zone <-> Ash Valley
				{
					Reference:     DemoMechLinkDropzoneToValleyRef,
					FromSectorRef: DemoMechSectorDropzoneRef,
					ToSectorRef:   DemoMechSectorValleyRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				{
					Reference:     DemoMechLinkValleyToDropzoneRef,
					FromSectorRef: DemoMechSectorValleyRef,
					ToSectorRef:   DemoMechSectorDropzoneRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				// North Ridge <-> South Ridge
				{
					Reference:     DemoMechLinkRidgeNorthToRidgeSouthRef,
					FromSectorRef: DemoMechSectorRidgeNorthRef,
					ToSectorRef:   DemoMechSectorRidgeSouthRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				{
					Reference:     DemoMechLinkRidgeSouthToRidgeNorthRef,
					FromSectorRef: DemoMechSectorRidgeSouthRef,
					ToSectorRef:   DemoMechSectorRidgeNorthRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				// South Ridge <-> Refinery
				{
					Reference:     DemoMechLinkRidgeSouthToRefineryRef,
					FromSectorRef: DemoMechSectorRidgeSouthRef,
					ToSectorRef:   DemoMechSectorRefineryRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				{
					Reference:     DemoMechLinkRefineryToRidgeSouthRef,
					FromSectorRef: DemoMechSectorRefineryRef,
					ToSectorRef:   DemoMechSectorRidgeSouthRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				// Ash Valley <-> Crossroads
				{
					Reference:     DemoMechLinkValleyToCrossroadsRef,
					FromSectorRef: DemoMechSectorValleyRef,
					ToSectorRef:   DemoMechSectorCrossroadsRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				{
					Reference:     DemoMechLinkCrossroadsToValleyRef,
					FromSectorRef: DemoMechSectorCrossroadsRef,
					ToSectorRef:   DemoMechSectorValleyRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				// Crossroads <-> Scorched Forest
				{
					Reference:     DemoMechLinkCrossroadsToForestRef,
					FromSectorRef: DemoMechSectorCrossroadsRef,
					ToSectorRef:   DemoMechSectorForestRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				{
					Reference:     DemoMechLinkForestToCrossroadsRef,
					FromSectorRef: DemoMechSectorForestRef,
					ToSectorRef:   DemoMechSectorCrossroadsRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				// Refinery <-> Citadel
				{
					Reference:     DemoMechLinkRefinerytoCitadelRef,
					FromSectorRef: DemoMechSectorRefineryRef,
					ToSectorRef:   DemoMechSectorCitadelRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
				{
					Reference:     DemoMechLinkCitadelToRefineryRef,
					FromSectorRef: DemoMechSectorCitadelRef,
					ToSectorRef:   DemoMechSectorRefineryRef,
					Record:        &mecha_record.MechaSectorLink{},
				},
			// Scorched Forest <-> Citadel
			{
				Reference:     DemoMechLinkForestToCitadelRef,
				FromSectorRef: DemoMechSectorForestRef,
				ToSectorRef:   DemoMechSectorCitadelRef,
				Record:        &mecha_record.MechaSectorLink{},
			},
			{
				Reference:     DemoMechLinkCitadelToForestRef,
				FromSectorRef: DemoMechSectorCitadelRef,
				ToSectorRef:   DemoMechSectorForestRef,
				Record:        &mecha_record.MechaSectorLink{},
			},
		},
		MechaComputerOpponentConfigs: []harness.MechaComputerOpponentConfig{
			{
				Reference: DemoMechaComputerOpponentRef,
				Record: &mecha_record.MechaComputerOpponent{
					Name:        "Scorched Ridge Garrison",
					Description: "The defending garrison force that has held the Citadel for years. Aggressive, experienced, and fighting on home ground.",
					Aggression:  6,
					IQ:          5,
				},
			},
		},
		MechaLanceConfigs: []harness.MechaLanceConfig{
			{
				Reference: DemoMechaPlayerStarterLanceRef,
				LanceType: mecha_record.LanceTypeStarter,
				Record: &mecha_record.MechaLance{
					Name:        "Strike Lance",
					Description: "Standard assault lance issued to incoming commanders. One light recon mech and one medium fire-support mech.",
				},
				LanceMechConfigs: []harness.MechaLanceMechConfig{
					{
						Reference:  DemoMechaPlayerStarterMech1Ref,
						ChassisRef: DemoMechChassisViperRef,
						WeaponConfigRefs: []harness.MechaLanceMechWeaponRef{
							{WeaponRef: DemoMechWeaponLightPulseRef, SlotLocation: "left-arm"},
							{WeaponRef: DemoMechWeaponChaingunRef, SlotLocation: "right-arm"},
						},
						Record: &mecha_record.MechaLanceMech{
							Callsign: "Strike-1",
						},
					},
					{
						Reference:  DemoMechaPlayerStarterMech2Ref,
						ChassisRef: DemoMechChassisRangerRef,
						WeaponConfigRefs: []harness.MechaLanceMechWeaponRef{
							{WeaponRef: DemoMechWeaponPulseRef, SlotLocation: "left-torso"},
							{WeaponRef: DemoMechWeaponRocketPackRef, SlotLocation: "right-arm"},
						},
						Record: &mecha_record.MechaLanceMech{
							Callsign: "Strike-2",
						},
					},
				},
			},
			{
				Reference: DemoMechaComputerOpponentLanceRef,
				LanceType: mecha_record.LanceTypeOpponent,
				Record: &mecha_record.MechaLance{
					Name:        "Garrison Heavy Lance",
					Description: "The Citadel's primary defensive lance. Well-armoured heavies backed by lighter recon elements.",
				},
				LanceMechConfigs: []harness.MechaLanceMechConfig{
					{
						Reference:  DemoMechaComputerOpponentMech1Ref,
						ChassisRef: DemoMechChassisCrusherRef,
						WeaponConfigRefs: []harness.MechaLanceMechWeaponRef{
							{WeaponRef: DemoMechWeaponPlasmaRef, SlotLocation: "left-torso"},
							{WeaponRef: DemoMechWeaponRotaryCannonRef, SlotLocation: "right-torso"},
						},
						Record: &mecha_record.MechaLanceMech{
							Callsign: "Garrison-1",
						},
					},
					{
						Reference:  DemoMechaComputerOpponentMech2Ref,
						ChassisRef: DemoMechChassisRangerRef,
						WeaponConfigRefs: []harness.MechaLanceMechWeaponRef{
							{WeaponRef: DemoMechWeaponPulseRef, SlotLocation: "left-torso"},
							{WeaponRef: DemoMechWeaponRocketPackRef, SlotLocation: "right-arm"},
						},
						Record: &mecha_record.MechaLanceMech{
							Callsign: "Garrison-2",
						},
					},
				},
			},
		},
	},
}
}
