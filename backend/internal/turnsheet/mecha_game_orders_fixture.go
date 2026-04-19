package turnsheet

import (
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// MechaGameOrdersFixture returns the sample rendering fixture for the
// mecha orders turn sheet.
func MechaGameOrdersFixture() DevFixture {
	return DevFixture{
		TemplatePath:   "turnsheet/mecha_game_orders.template",
		OutputBaseName: "mecha_game_orders_turnsheet",
		BackgroundFile: "background-darkforest.png",
		MakeData: func(bg, code string) any {
			deadline := time.Now().Add(7 * 24 * time.Hour)
			return &OrdersData{
				TurnSheetTemplateData: TurnSheetTemplateData{
					GameName:              strPtr("Steel Thunder"),
					GameType:              strPtr("mecha"),
					TurnNumber:            intPtr(3),
					TurnSheetTitle:        strPtr("Mech Orders"),
					TurnSheetInstructions: strPtr(DefaultOrdersInstructions()),
					TurnSheetCode:         strPtr(code),
					TurnSheetDeadline:     &deadline,
					BackgroundImage:       &bg,
					TurnEvents: []TurnEvent{
						{Category: TurnEventCategoryMovement, Icon: TurnEventIconMovement, Message: "Hammer advanced to Central Wastes."},
						{Category: TurnEventCategoryCombat, Icon: TurnEventIconCombat, Message: "Anvil fired Pulse Cannon at Stalker — HIT for 5 damage."},
						{Category: TurnEventCategorySystem, Icon: TurnEventIconSystem, Message: "Hammer field repairs restored 18 armor (72/72)."},
					},
				},
			SquadName: "Alpha Squad",
			SquadMechs: []MechOrderEntry{
				{
					MechInstanceID: "mech-1", MechCallsign: "Hammer", MechStatus: "operational",
					CurrentSectorName: "Central Wastes", ChassisName: "Scout", ChassisClass: "light",
					CurrentArmor: 55, MaxArmor: 72, CurrentStructure: 32, MaxStructure: 32,
					CurrentHeat: 4, HeatCapacity: 18, Speed: 7, EffectiveSpeed: 7, PilotSkill: 4,
					Weapons: []MechWeaponEntry{
						{Name: "Light Pulse Cannon", Damage: 3, HeatCost: 1, RangeBand: "short", SlotLocation: "left-arm"},
						{Name: "Chaingun", Damage: 2, HeatCost: 0, RangeBand: "short", SlotLocation: "right-arm"},
					},
					ReachableSectors: []SectorOption{
						{SectorInstanceID: "sector-1", SectorName: "Northern Ridge"},
						{SectorInstanceID: "sector-2", SectorName: "Southern Flats"},
						{SectorInstanceID: "sector-3", SectorName: "Eastern Pass"},
					},
				},
				{
					MechInstanceID: "mech-2", MechCallsign: "Anvil", MechStatus: "operational",
					CurrentSectorName: "Central Wastes", ChassisName: "Sentinel", ChassisClass: "medium",
					CurrentArmor: 130, MaxArmor: 130, CurrentStructure: 65, MaxStructure: 65,
					CurrentHeat: 6, HeatCapacity: 28, Speed: 4, EffectiveSpeed: 4, PilotSkill: 4,
					Weapons: []MechWeaponEntry{
						{Name: "Pulse Cannon", Damage: 5, HeatCost: 3, RangeBand: "medium", SlotLocation: "left-torso"},
						{Name: "Rocket Pack", Damage: 8, HeatCost: 3, RangeBand: "short", SlotLocation: "right-arm", AmmoCapacity: 2},
					},
					Equipment: []MechEquipmentEntry{
						{Name: "Targeting Computer Mk II", EffectKind: "targeting_computer", Magnitude: 10, HeatCost: 1, MountSize: "medium", SlotLocation: "head"},
						{Name: "Ammo Bin (Standard)", EffectKind: "ammo_bin", Magnitude: 8, HeatCost: 0, MountSize: "small", SlotLocation: "right-torso"},
					},
					AmmoRemaining: 6, AmmoCapacity: 10,
					ReachableSectors: []SectorOption{
						{SectorInstanceID: "sector-1", SectorName: "Northern Ridge"},
						{SectorInstanceID: "sector-2", SectorName: "Southern Flats"},
					},
				},
				{
					MechInstanceID: "mech-3", MechCallsign: "Titan", MechStatus: "damaged",
					CurrentSectorName: "Northern Ridge", ChassisName: "Colossus", ChassisClass: "heavy",
					CurrentArmor: 45, MaxArmor: 180, CurrentStructure: 60, MaxStructure: 80,
					CurrentHeat: 14, HeatCapacity: 35, Speed: 3, EffectiveSpeed: 5, PilotSkill: 5,
					Weapons: []MechWeaponEntry{
						{Name: "Auto-Cannon", Damage: 10, HeatCost: 5, RangeBand: "medium", SlotLocation: "right-arm"},
					},
					Equipment: []MechEquipmentEntry{
						{Name: "Jump Jets", EffectKind: "jump_jets", Magnitude: 2, HeatCost: 2, MountSize: "medium", SlotLocation: "left-leg"},
						{Name: "Heat Sink", EffectKind: "heat_sink", Magnitude: 4, HeatCost: 0, MountSize: "small", SlotLocation: "right-torso"},
					},
					ReachableSectors: []SectorOption{
						{SectorInstanceID: "sector-1", SectorName: "Northern Ridge"},
						{SectorInstanceID: "sector-2", SectorName: "Southern Flats"},
						{SectorInstanceID: "sector-3", SectorName: "Eastern Pass"},
						{SectorInstanceID: "sector-4", SectorName: "Ridge Overlook"},
					},
				},
				{
					MechInstanceID: "mech-4", MechCallsign: "Wrench", MechStatus: "operational",
					IsRefitting: true, ChassisName: "Scout", ChassisClass: "light",
					CurrentArmor: 72, MaxArmor: 72, CurrentStructure: 32, MaxStructure: 32,
					Speed: 7, EffectiveSpeed: 7, PilotSkill: 3,
				},
			},
			AvailableSectors: []SectorOption{
				{SectorInstanceID: "sector-1", SectorName: "Northern Ridge"},
				{SectorInstanceID: "sector-2", SectorName: "Southern Flats"},
				{SectorInstanceID: "sector-3", SectorName: "Eastern Pass"},
				{SectorInstanceID: "sector-4", SectorName: "Ridge Overlook"},
			},
			EnemyMechs: []EnemyMechOption{
				{MechInstanceID: "enemy-1", Callsign: "Stalker", SectorName: "Northern Ridge"},
				{MechInstanceID: "enemy-2", Callsign: "Predator", SectorName: "Southern Flats"},
			},
		}
		},
		NewProcessor: func(l logger.Logger, cfg config.Config) (TurnSheetProcessor, error) {
			return NewMechaGameOrdersProcessor(l, cfg)
		},
	}
}
