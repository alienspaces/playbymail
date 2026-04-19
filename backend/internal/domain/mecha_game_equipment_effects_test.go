package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func eq(kind string, magnitude, heatCost int) *mecha_game_record.MechaGameEquipment {
	return &mecha_game_record.MechaGameEquipment{
		EffectKind: kind,
		Magnitude:  magnitude,
		HeatCost:   heatCost,
	}
}

func TestAggregateMechaGameEquipmentEffects(t *testing.T) {
	// A fully populated catalog covering every effect kind so each case can
	// cherry-pick the entries it needs without redefining equipment records.
	catalog := map[string]*mecha_game_record.MechaGameEquipment{
		"hs1":  eq(mecha_game_record.EquipmentEffectKindHeatSink, 3, 0),
		"hs2":  eq(mecha_game_record.EquipmentEffectKindHeatSink, 2, 0),
		"tc1":  eq(mecha_game_record.EquipmentEffectKindTargetingComputer, 10, 0),
		"au1":  eq(mecha_game_record.EquipmentEffectKindArmorUpgrade, 50, 0),
		"jj1":  eq(mecha_game_record.EquipmentEffectKindJumpJets, 2, 0),
		"ecm1": eq(mecha_game_record.EquipmentEffectKindECM, 15, 0),
		"ab1":  eq(mecha_game_record.EquipmentEffectKindAmmoBin, 100, 0),
	}

	cases := []struct {
		name      string
		entries   []mecha_game_record.EquipmentConfigEntry
		byID      map[string]*mecha_game_record.MechaGameEquipment
		refitting bool
		want      MechaGameEquipmentEffects
	}{
		{
			name: "sums per-kind bonuses across one entry of each kind",
			entries: []mecha_game_record.EquipmentConfigEntry{
				{EquipmentID: "hs1"}, {EquipmentID: "hs2"},
				{EquipmentID: "tc1"},
				{EquipmentID: "au1"},
				{EquipmentID: "jj1"},
				{EquipmentID: "ecm1"},
				{EquipmentID: "ab1"},
			},
			byID: catalog,
			want: MechaGameEquipmentEffects{
				HeatDissipationBonus: 5,
				HitChanceBonus:       10,
				ArmorBonus:           50,
				SpeedBonus:           2,
				CoverBonus:           15,
			},
		},
		{
			name: "refitting mech returns zeroed effects",
			entries: []mecha_game_record.EquipmentConfigEntry{
				{EquipmentID: "hs1"}, {EquipmentID: "au1"},
			},
			byID:      catalog,
			refitting: true,
			want:      MechaGameEquipmentEffects{},
		},
		{
			name:    "unknown equipment ids are skipped silently",
			entries: []mecha_game_record.EquipmentConfigEntry{{EquipmentID: "missing"}},
			byID:    nil,
			want:    MechaGameEquipmentEffects{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := AggregateMechaGameEquipmentEffects(tc.entries, tc.byID, tc.refitting)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestEffectiveMechaGameMaxArmor(t *testing.T) {
	cases := []struct {
		name    string
		chassis *mecha_game_record.MechaGameChassis
		effects MechaGameEquipmentEffects
		want    int
	}{
		{
			name:    "includes armor bonus on top of chassis armor",
			chassis: &mecha_game_record.MechaGameChassis{ArmorPoints: 100},
			effects: MechaGameEquipmentEffects{ArmorBonus: 50},
			want:    150,
		},
		{
			name:    "nil chassis falls back to bonus alone",
			chassis: nil,
			effects: MechaGameEquipmentEffects{ArmorBonus: 10},
			want:    10,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, EffectiveMechaGameMaxArmor(tc.chassis, tc.effects))
		})
	}
}

func TestEffectiveMechaGameSpeed_IncludesJumpJetBonus(t *testing.T) {
	chassis := &mecha_game_record.MechaGameChassis{Speed: 3}
	effects := MechaGameEquipmentEffects{SpeedBonus: 2}
	require.Equal(t, 5, EffectiveMechaGameSpeed(chassis, effects))
}

func TestMaxMechaGameAmmoCapacity_SumsWeaponsAndBins(t *testing.T) {
	weapons := map[string]*mecha_game_record.MechaGameWeapon{
		"w1": {AmmoCapacity: 30},
		"w2": {AmmoCapacity: 20},
		"w3": {AmmoCapacity: 0}, // energy weapon contributes nothing
	}
	equipment := map[string]*mecha_game_record.MechaGameEquipment{
		"ab1": eq(mecha_game_record.EquipmentEffectKindAmmoBin, 100, 0),
		"hs1": eq(mecha_game_record.EquipmentEffectKindHeatSink, 5, 0),
	}
	weaponEntries := []mecha_game_record.WeaponConfigEntry{
		{WeaponID: "w1"}, {WeaponID: "w2"}, {WeaponID: "w3"},
	}
	equipmentEntries := []mecha_game_record.EquipmentConfigEntry{
		{EquipmentID: "ab1"}, {EquipmentID: "hs1"},
	}
	require.Equal(t, 150, MaxMechaGameAmmoCapacity(weaponEntries, weapons, equipmentEntries, equipment))
}

func TestMechaGameEquipmentHeatCostSplit(t *testing.T) {
	// Shared fixture: one entry of every kind with distinct, non-overlapping
	// heat_cost values so each helper can only match on its own kinds.
	entries := []mecha_game_record.EquipmentConfigEntry{
		{EquipmentID: "jj1"}, {EquipmentID: "tc1"}, {EquipmentID: "ab1"},
		{EquipmentID: "hs1"}, {EquipmentID: "au1"}, {EquipmentID: "ecm1"},
	}
	byID := map[string]*mecha_game_record.MechaGameEquipment{
		"jj1":  eq(mecha_game_record.EquipmentEffectKindJumpJets, 2, 3),
		"tc1":  eq(mecha_game_record.EquipmentEffectKindTargetingComputer, 10, 5),
		"ab1":  eq(mecha_game_record.EquipmentEffectKindAmmoBin, 100, 2),
		"hs1":  eq(mecha_game_record.EquipmentEffectKindHeatSink, 3, 4),
		"au1":  eq(mecha_game_record.EquipmentEffectKindArmorUpgrade, 50, 2),
		"ecm1": eq(mecha_game_record.EquipmentEffectKindECM, 15, 3),
	}

	t.Run("jump jet heat cost only counts jump jet entries", func(t *testing.T) {
		cases := []struct {
			name      string
			refitting bool
			want      int
		}{
			{"powered on", false, 3},
			{"refitting pays no equipment heat", true, 0},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.want, MechaGameEquipmentJumpJetHeatCost(entries, byID, tc.refitting))
			})
		}
	})

	t.Run("combat heat cost counts targeting and ammo bin by activity flags", func(t *testing.T) {
		// Flags independently gate targeting-computer and ammo-bin heat. In
		// practice firing a weapon implies attack=true, but the helper takes
		// them separately so tests can pin down each branch.
		cases := []struct {
			name        string
			refitting   bool
			attackFired bool
			ammoFired   bool
			want        int
		}{
			{"no combat activity", false, false, false, 0},
			{"attack only — targeting computer fires", false, true, false, 5},
			{"ammo weapon only — ammo bin fires", false, false, true, 2},
			{"attack and ammo fire — both add", false, true, true, 7},
			{"refitting shuts combat heat off entirely", true, true, true, 0},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.want, MechaGameEquipmentCombatHeatCost(entries, byID, tc.refitting, tc.attackFired, tc.ammoFired))
			})
		}
	})

	t.Run("always-on heat cost counts heat sink, armor upgrade and ECM", func(t *testing.T) {
		cases := []struct {
			name      string
			refitting bool
			want      int
		}{
			{"powered on", false, 9}, // hs(4) + au(2) + ecm(3)
			{"refitting zeros always-on heat", true, 0},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.want, MechaGameEquipmentAlwaysOnHeatCost(entries, byID, tc.refitting))
			})
		}
	})
}
