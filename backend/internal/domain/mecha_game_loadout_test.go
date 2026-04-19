package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func TestFitsLoadout_StrictFit(t *testing.T) {
	cap := LoadoutCapacity{Small: 2, Medium: 2, Large: 1}
	items := []Mountable{
		{MountSize: mecha_game_record.WeaponMountSizeSmall, Label: "s1"},
		{MountSize: mecha_game_record.WeaponMountSizeSmall, Label: "s2"},
		{MountSize: mecha_game_record.WeaponMountSizeMedium, Label: "m1"},
		{MountSize: mecha_game_record.WeaponMountSizeMedium, Label: "m2"},
		{MountSize: mecha_game_record.WeaponMountSizeLarge, Label: "l1"},
	}
	require.NoError(t, fitsLoadout(cap, items))
}

func TestFitsLoadout_EmptyItemsAlwaysFits(t *testing.T) {
	require.NoError(t, fitsLoadout(LoadoutCapacity{}, nil))
	require.NoError(t, fitsLoadout(LoadoutCapacity{Small: 1}, nil))
}

func TestFitsLoadout_SmallSpillsUpOneBand(t *testing.T) {
	// No small slots; two small items must spill into the single medium slot
	// and one spills further into the large slot.
	cap := LoadoutCapacity{Small: 0, Medium: 1, Large: 1}
	items := []Mountable{
		{MountSize: mecha_game_record.WeaponMountSizeSmall, Label: "s1"},
		{MountSize: mecha_game_record.WeaponMountSizeSmall, Label: "s2"},
	}
	require.NoError(t, fitsLoadout(cap, items))
}

func TestFitsLoadout_SmallSpillsUpTwoBands(t *testing.T) {
	cap := LoadoutCapacity{Small: 0, Medium: 0, Large: 1}
	items := []Mountable{
		{MountSize: mecha_game_record.WeaponMountSizeSmall, Label: "s1"},
	}
	require.NoError(t, fitsLoadout(cap, items))
}

func TestFitsLoadout_MediumSpillsIntoLarge(t *testing.T) {
	cap := LoadoutCapacity{Small: 0, Medium: 0, Large: 2}
	items := []Mountable{
		{MountSize: mecha_game_record.WeaponMountSizeMedium, Label: "m1"},
		{MountSize: mecha_game_record.WeaponMountSizeMedium, Label: "m2"},
	}
	require.NoError(t, fitsLoadout(cap, items))
}

func TestFitsLoadout_LargeRefusesMediumSlot(t *testing.T) {
	// Medium slot is plentiful but a large item may never use it.
	cap := LoadoutCapacity{Small: 0, Medium: 5, Large: 0}
	items := []Mountable{
		{MountSize: mecha_game_record.WeaponMountSizeLarge, Label: "Big Gun"},
	}
	err := fitsLoadout(cap, items)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Big Gun")
	require.Contains(t, err.Error(), "large")
}

func TestFitsLoadout_MediumRefusesSmallSlot(t *testing.T) {
	cap := LoadoutCapacity{Small: 5, Medium: 0, Large: 0}
	items := []Mountable{
		{MountSize: mecha_game_record.WeaponMountSizeMedium, Label: "Mid Cannon"},
	}
	err := fitsLoadout(cap, items)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Mid Cannon")
}

func TestFitsLoadout_LargePlacedBeforeSmallSpillover(t *testing.T) {
	// Regression guard: with one large and one small item plus a single large
	// slot and no smaller slots, the large must claim the large slot so the
	// small item has somewhere to spill into later. A naive single-pass
	// allocator might hand the large slot to the small item and then fail.
	cap := LoadoutCapacity{Small: 0, Medium: 0, Large: 1}
	items := []Mountable{
		{MountSize: mecha_game_record.WeaponMountSizeSmall, Label: "s1"},
		{MountSize: mecha_game_record.WeaponMountSizeLarge, Label: "l1"},
	}
	err := fitsLoadout(cap, items)
	require.Error(t, err, "small item should not steal the only large slot")
	require.Contains(t, err.Error(), "s1")
}

func TestFitsLoadout_OverflowErrorNamesItem(t *testing.T) {
	cap := LoadoutCapacity{Small: 1}
	items := []Mountable{
		{MountSize: mecha_game_record.WeaponMountSizeSmall, Label: "First"},
		{MountSize: mecha_game_record.WeaponMountSizeSmall, Label: "Second"},
	}
	err := fitsLoadout(cap, items)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Second")
}

func TestFitsLoadout_UnknownMountSizeRejected(t *testing.T) {
	cap := LoadoutCapacity{Small: 5}
	items := []Mountable{{MountSize: "colossal", Label: "Oversized"}}
	err := fitsLoadout(cap, items)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "unknown") || strings.Contains(err.Error(), "colossal"))
}

func TestMountablesFromWeaponConfig_ResolvesNamesAndSizes(t *testing.T) {
	weapons := map[string]*mecha_game_record.MechaGameWeapon{
		"wid-1": {Name: "AC/10", MountSize: mecha_game_record.WeaponMountSizeMedium},
	}
	entries := []mecha_game_record.WeaponConfigEntry{
		{WeaponID: "wid-1", SlotLocation: "right-torso"},
	}
	out := MountablesFromWeaponConfig(entries, weapons)
	require.Len(t, out, 1)
	require.Equal(t, mecha_game_record.WeaponMountSizeMedium, out[0].MountSize)
	require.Equal(t, "AC/10", out[0].Label)
}

func TestMountablesFromWeaponConfig_UnknownWeaponFallsBackToID(t *testing.T) {
	entries := []mecha_game_record.WeaponConfigEntry{
		{WeaponID: "missing-id", SlotLocation: "left-arm"},
	}
	out := MountablesFromWeaponConfig(entries, nil)
	require.Len(t, out, 1)
	require.Equal(t, "missing-id", out[0].Label)
	require.Equal(t, "", out[0].MountSize)
}

func TestValidateWeaponLoadoutFits_NilChassisIsNoop(t *testing.T) {
	err := ValidateWeaponLoadoutFits(nil, []mecha_game_record.WeaponConfigEntry{{WeaponID: "x"}}, nil)
	require.NoError(t, err)
}

func TestValidateWeaponLoadoutFits_ReturnsFitError(t *testing.T) {
	chassis := &mecha_game_record.MechaGameChassis{
		SmallSlots:  0,
		MediumSlots: 0,
		LargeSlots:  0,
	}
	entries := []mecha_game_record.WeaponConfigEntry{{WeaponID: "wid-1"}}
	weapons := map[string]*mecha_game_record.MechaGameWeapon{
		"wid-1": {Name: "Plasma", MountSize: mecha_game_record.WeaponMountSizeLarge},
	}
	err := ValidateWeaponLoadoutFits(chassis, entries, weapons)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Plasma")
}

func TestLoadoutCapacityFromChassis_HandlesNil(t *testing.T) {
	require.Equal(t, LoadoutCapacity{}, LoadoutCapacityFromChassis(nil))
}

func TestMountablesFromEquipmentConfig_ResolvesNamesAndSizes(t *testing.T) {
	equipment := map[string]*mecha_game_record.MechaGameEquipment{
		"eid-1": {Name: "Double Heat Sink", MountSize: mecha_game_record.EquipmentMountSizeSmall},
	}
	entries := []mecha_game_record.EquipmentConfigEntry{
		{EquipmentID: "eid-1", SlotLocation: "center-torso"},
	}
	out := MountablesFromEquipmentConfig(entries, equipment)
	require.Len(t, out, 1)
	require.Equal(t, mecha_game_record.EquipmentMountSizeSmall, out[0].MountSize)
	require.Equal(t, "Double Heat Sink", out[0].Label)
}

func TestMountablesFromEquipmentConfig_UnknownEquipmentFallsBackToID(t *testing.T) {
	entries := []mecha_game_record.EquipmentConfigEntry{
		{EquipmentID: "missing-id", SlotLocation: "left-arm"},
	}
	out := MountablesFromEquipmentConfig(entries, nil)
	require.Len(t, out, 1)
	require.Equal(t, "missing-id", out[0].Label)
}

func TestValidateCombinedLoadoutFits_EmptyInputsAreNoop(t *testing.T) {
	chassis := &mecha_game_record.MechaGameChassis{SmallSlots: 1}
	require.NoError(t, ValidateCombinedLoadoutFits(chassis, nil, nil, nil, nil))
}

func TestValidateCombinedLoadoutFits_NilChassisIsNoop(t *testing.T) {
	weaponEntries := []mecha_game_record.WeaponConfigEntry{{WeaponID: "w1"}}
	require.NoError(t, ValidateCombinedLoadoutFits(nil, weaponEntries, nil, nil, nil))
}

func TestValidateCombinedLoadoutFits_SharesSlotBudget(t *testing.T) {
	chassis := &mecha_game_record.MechaGameChassis{SmallSlots: 1, MediumSlots: 0, LargeSlots: 0}
	weapons := map[string]*mecha_game_record.MechaGameWeapon{
		"w1": {Name: "Small Laser", MountSize: mecha_game_record.WeaponMountSizeSmall},
	}
	equipment := map[string]*mecha_game_record.MechaGameEquipment{
		"e1": {Name: "Heat Sink", MountSize: mecha_game_record.EquipmentMountSizeSmall},
	}
	weaponEntries := []mecha_game_record.WeaponConfigEntry{{WeaponID: "w1"}}
	equipmentEntries := []mecha_game_record.EquipmentConfigEntry{{EquipmentID: "e1"}}
	err := ValidateCombinedLoadoutFits(chassis, weaponEntries, weapons, equipmentEntries, equipment)
	require.Error(t, err, "single small slot must not accommodate both a weapon and equipment")
}

func TestValidateCombinedLoadoutFits_SpilloverWorksAcrossMix(t *testing.T) {
	chassis := &mecha_game_record.MechaGameChassis{SmallSlots: 0, MediumSlots: 0, LargeSlots: 2}
	weapons := map[string]*mecha_game_record.MechaGameWeapon{
		"w1": {Name: "Small Laser", MountSize: mecha_game_record.WeaponMountSizeSmall},
	}
	equipment := map[string]*mecha_game_record.MechaGameEquipment{
		"e1": {Name: "Heat Sink", MountSize: mecha_game_record.EquipmentMountSizeSmall},
	}
	weaponEntries := []mecha_game_record.WeaponConfigEntry{{WeaponID: "w1"}}
	equipmentEntries := []mecha_game_record.EquipmentConfigEntry{{EquipmentID: "e1"}}
	err := ValidateCombinedLoadoutFits(chassis, weaponEntries, weapons, equipmentEntries, equipment)
	require.NoError(t, err, "both small items should spill up into the two large slots")
}
