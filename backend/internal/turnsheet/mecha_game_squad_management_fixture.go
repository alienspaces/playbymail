package turnsheet

import (
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// MechaGameSquadManagementFixture returns the sample rendering fixture for the
// mecha squad management turn sheet.
func MechaGameSquadManagementFixture() DevFixture {
	return DevFixture{
		TemplatePath:   "turnsheet/mecha_game_squad_management.template",
		OutputBaseName: "mecha_game_squad_management_turnsheet",
		BackgroundFile: "background-darkforest.png",
		MakeData: func(bg, code string) any {
			deadline := time.Now().Add(7 * 24 * time.Hour)
			return &SquadManagementData{
				TurnSheetTemplateData: TurnSheetTemplateData{
					GameName:              strPtr("Steel Thunder"),
					GameType:              strPtr("mecha"),
					TurnNumber:            intPtr(2),
					TurnSheetTitle:        strPtr("Squad Management"),
					TurnSheetInstructions: strPtr(DefaultManagementInstructions()),
					TurnSheetCode:         strPtr(code),
					TurnSheetDeadline:     &deadline,
					BackgroundImage:       &bg,
					TurnEvents: []TurnEvent{
						{Category: TurnEventCategorySystem, Icon: TurnEventIconSystem, Message: "Hammer field repairs restored 18 armor (72/72)."},
						{Category: TurnEventCategorySystem, Icon: TurnEventIconSystem, Message: "Squad received 2 supply points (6 total)."},
					},
				},
			SquadName:    "Alpha Squad",
			SupplyPoints: 6,
			Mechs: []ManagementMechEntry{
				{
					MechInstanceID: "mech-1", Callsign: "Hammer", ChassisName: "Scout", ChassisClass: "light",
					Status: "operational", IsAtDepot: true,
					CurrentArmor: 72, MaxArmor: 72, CurrentStructure: 28, MaxStructure: 32, StructureDamage: 4,
					CurrentHeat: 2, HeatCapacity: 18,
					Weapons: []MechWeaponSlot{
						{SlotLocation: "left-arm", CurrentWeaponID: "wpn-1", CurrentWeaponName: "Light Pulse Cannon"},
						{SlotLocation: "right-arm", CurrentWeaponID: "wpn-2", CurrentWeaponName: "Chaingun"},
					},
					Equipment: []MechEquipmentEntry{
						{Name: "Heat Sink", EffectKind: "heat_sink", Magnitude: 3, HeatCost: 0, MountSize: "small", SlotLocation: "right-torso"},
					},
				},
				{
					MechInstanceID: "mech-2", Callsign: "Anvil", ChassisName: "Sentinel", ChassisClass: "medium",
					Status: "operational", IsAtDepot: false,
					CurrentArmor: 100, MaxArmor: 130, CurrentStructure: 65, MaxStructure: 65, StructureDamage: 0,
					CurrentHeat: 6, HeatCapacity: 28,
					Weapons: []MechWeaponSlot{
						{SlotLocation: "left-torso", CurrentWeaponID: "wpn-3", CurrentWeaponName: "Pulse Cannon"},
						{SlotLocation: "right-arm", CurrentWeaponID: "wpn-4", CurrentWeaponName: "Rocket Pack"},
					},
					Equipment: []MechEquipmentEntry{
						{Name: "Targeting Computer Mk II", EffectKind: "targeting_computer", Magnitude: 10, HeatCost: 1, MountSize: "medium", SlotLocation: "head"},
						{Name: "Ammo Bin (Standard)", EffectKind: "ammo_bin", Magnitude: 8, HeatCost: 0, MountSize: "small", SlotLocation: "right-torso"},
					},
					AmmoRemaining: 6, AmmoCapacity: 10,
				},
				{
					MechInstanceID: "mech-3", Callsign: "Wrench", ChassisName: "Scout", ChassisClass: "light",
					Status: "operational", IsAtDepot: true, IsRefitting: true,
					CurrentArmor: 72, MaxArmor: 72, CurrentStructure: 32, MaxStructure: 32,
					CurrentHeat: 0, HeatCapacity: 18,
				},
			},
			WeaponCatalog: []CatalogWeapon{
				{WeaponID: "cat-1", Name: "Light Pulse Cannon", Damage: 3, HeatCost: 1, RangeBand: "short"},
				{WeaponID: "cat-2", Name: "Chaingun", Damage: 2, HeatCost: 0, RangeBand: "short"},
				{WeaponID: "cat-3", Name: "Pulse Cannon", Damage: 5, HeatCost: 3, RangeBand: "medium"},
				{WeaponID: "cat-4", Name: "Rocket Pack", Damage: 8, HeatCost: 3, RangeBand: "short", AmmoCapacity: 2},
				{WeaponID: "cat-5", Name: "Auto-Cannon", Damage: 10, HeatCost: 5, RangeBand: "medium", AmmoCapacity: 1},
			},
		}
		},
		NewProcessor: func(l logger.Logger, cfg config.Config) (TurnSheetProcessor, error) {
			return NewMechaGameSquadManagementProcessor(l, cfg)
		},
	}
}
