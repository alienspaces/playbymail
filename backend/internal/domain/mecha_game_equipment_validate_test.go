package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func newValidEquipment(kind string, magnitude int) *mecha_game_record.MechaGameEquipment {
	return &mecha_game_record.MechaGameEquipment{
		Record:     record.Record{ID: uuid.NewString()},
		GameID:     uuid.NewString(),
		Name:       "Test Equipment",
		MountSize:  mecha_game_record.EquipmentMountSizeSmall,
		EffectKind: kind,
		Magnitude:  magnitude,
		HeatCost:   0,
	}
}

func TestValidateEquipment_AcceptsAllEffectKinds(t *testing.T) {
	cases := []struct {
		kind      string
		magnitude int
	}{
		{mecha_game_record.EquipmentEffectKindHeatSink, 1},
		{mecha_game_record.EquipmentEffectKindTargetingComputer, 10},
		{mecha_game_record.EquipmentEffectKindArmorUpgrade, 50},
		{mecha_game_record.EquipmentEffectKindJumpJets, 3},
		{mecha_game_record.EquipmentEffectKindECM, 25},
		{mecha_game_record.EquipmentEffectKindAmmoBin, 100},
	}
	for _, c := range cases {
		t.Run(c.kind, func(t *testing.T) {
			err := validateMechaGameEquipmentRec(&validateMechaGameEquipmentArgs{nextRec: newValidEquipment(c.kind, c.magnitude)}, false)
			require.NoError(t, err)
		})
	}
}

func TestValidateEquipment_RejectsUnknownEffectKind(t *testing.T) {
	rec := newValidEquipment("bogus_kind", 1)
	err := validateMechaGameEquipmentRec(&validateMechaGameEquipmentArgs{nextRec: rec}, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "effect_kind")
}

func TestValidateEquipment_EnforcesPerKindMagnitudeCap(t *testing.T) {
	// Jump jets cap at 5 — 6 should fail.
	rec := newValidEquipment(mecha_game_record.EquipmentEffectKindJumpJets, 6)
	err := validateMechaGameEquipmentRec(&validateMechaGameEquipmentArgs{nextRec: rec}, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "magnitude")
}

func TestValidateEquipment_RejectsZeroMagnitude(t *testing.T) {
	rec := newValidEquipment(mecha_game_record.EquipmentEffectKindHeatSink, 0)
	err := validateMechaGameEquipmentRec(&validateMechaGameEquipmentArgs{nextRec: rec}, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "magnitude")
}

func TestValidateEquipment_RejectsHeatCostOutOfRange(t *testing.T) {
	rec := newValidEquipment(mecha_game_record.EquipmentEffectKindHeatSink, 1)
	rec.HeatCost = 21
	err := validateMechaGameEquipmentRec(&validateMechaGameEquipmentArgs{nextRec: rec}, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "heat_cost")
}

func TestValidateEquipment_AcceptsHeatCostAtBoundary(t *testing.T) {
	rec := newValidEquipment(mecha_game_record.EquipmentEffectKindHeatSink, 1)
	rec.HeatCost = 20
	err := validateMechaGameEquipmentRec(&validateMechaGameEquipmentArgs{nextRec: rec}, false)
	require.NoError(t, err)
}

func TestValidateEquipment_RejectsNegativeHeatCost(t *testing.T) {
	rec := newValidEquipment(mecha_game_record.EquipmentEffectKindHeatSink, 1)
	rec.HeatCost = -1
	err := validateMechaGameEquipmentRec(&validateMechaGameEquipmentArgs{nextRec: rec}, false)
	require.Error(t, err)
}

func TestValidateEquipment_RejectsInvalidMountSize(t *testing.T) {
	rec := newValidEquipment(mecha_game_record.EquipmentEffectKindHeatSink, 1)
	rec.MountSize = "colossal"
	err := validateMechaGameEquipmentRec(&validateMechaGameEquipmentArgs{nextRec: rec}, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "mount_size")
}

func TestValidateEquipment_DefaultsEmptyMountSizeToMedium(t *testing.T) {
	rec := newValidEquipment(mecha_game_record.EquipmentEffectKindHeatSink, 1)
	rec.MountSize = ""
	err := validateMechaGameEquipmentRec(&validateMechaGameEquipmentArgs{nextRec: rec}, false)
	require.NoError(t, err)
	require.Equal(t, mecha_game_record.EquipmentMountSizeMedium, rec.MountSize)
}
