package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

// validChassisBaseRec returns a record that passes every non-slot field
// validation so individual tests can tweak just the slot values without
// repeating the same boilerplate.
func validChassisBaseRec() *mecha_game_record.MechaGameChassis {
	rec := &mecha_game_record.MechaGameChassis{
		GameID:          "00000000-0000-0000-0000-000000000001",
		Name:            "Test Chassis",
		ChassisClass:    mecha_game_record.ChassisClassMedium,
		ArmorPoints:     100,
		StructurePoints: 50,
		HeatCapacity:    30,
		Speed:           4,
		SmallSlots:      2,
		MediumSlots:     2,
		LargeSlots:      1,
	}
	return rec
}

func runCreateValidator(t *testing.T, rec *mecha_game_record.MechaGameChassis) error {
	t.Helper()
	args := &validateMechaGameChassisArgs{nextRec: rec}
	return validateMechaGameChassisRec(args, false)
}

func TestValidateMechaGameChassis_SlotDefaultsPass(t *testing.T) {
	rec := validChassisBaseRec()
	require.NoError(t, runCreateValidator(t, rec))
}

func TestValidateMechaGameChassis_NegativeSlotsRejected(t *testing.T) {
	rec := validChassisBaseRec()
	rec.SmallSlots = -1
	err := runCreateValidator(t, rec)
	require.Error(t, err)
}

func TestValidateMechaGameChassis_OverMaxSlotsRejected(t *testing.T) {
	for _, tc := range []struct {
		name  string
		tweak func(*mecha_game_record.MechaGameChassis)
	}{
		{"small over max", func(r *mecha_game_record.MechaGameChassis) { r.SmallSlots = 11 }},
		{"medium over max", func(r *mecha_game_record.MechaGameChassis) { r.MediumSlots = 11 }},
		{"large over max", func(r *mecha_game_record.MechaGameChassis) { r.LargeSlots = 11 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := validChassisBaseRec()
			tc.tweak(rec)
			err := runCreateValidator(t, rec)
			require.Error(t, err)
		})
	}
}

func TestValidateMechaGameChassis_AllZeroSlotsRejected(t *testing.T) {
	rec := validChassisBaseRec()
	rec.SmallSlots = 0
	rec.MediumSlots = 0
	rec.LargeSlots = 0
	err := runCreateValidator(t, rec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least one slot")
}

func TestValidateMechaGameChassis_SingleSlotClassIsAllowed(t *testing.T) {
	// A light scout with only a single small hardpoint is legal — the "at
	// least one slot" rule guards against creating a chassis that nothing
	// can be mounted to, not against minimalism.
	rec := validChassisBaseRec()
	rec.ChassisClass = mecha_game_record.ChassisClassLight
	rec.SmallSlots = 1
	rec.MediumSlots = 0
	rec.LargeSlots = 0
	require.NoError(t, runCreateValidator(t, rec))
}

func TestDefaultSlotsForChassisClass_ReturnsExpectedValues(t *testing.T) {
	cases := []struct {
		class                 string
		small, medium, large int
	}{
		{mecha_game_record.ChassisClassLight, 2, 1, 0},
		{mecha_game_record.ChassisClassMedium, 2, 2, 1},
		{mecha_game_record.ChassisClassHeavy, 2, 2, 2},
		{mecha_game_record.ChassisClassAssault, 2, 3, 3},
		{"unknown-class-falls-back-to-medium", 2, 2, 1},
	}
	for _, tc := range cases {
		t.Run(tc.class, func(t *testing.T) {
			s, m, l := mecha_game_record.DefaultSlotsForChassisClass(tc.class)
			require.Equal(t, tc.small, s)
			require.Equal(t, tc.medium, m)
			require.Equal(t, tc.large, l)
		})
	}
}
