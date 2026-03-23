package mecha

import (
	"testing"

	"github.com/stretchr/testify/require"

	corerecord "gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

// Tests for pure helper functions in combat_resolution.go.
// Integration tests for resolveCombat require a full DB harness and are omitted here.

func TestRangeDistance(t *testing.T) {
	t.Parallel()

	// Build a simple linear map: A <-> B <-> C <-> D
	secA := &sectorState{
		Instance:            &mecha_record.MechaSectorInstance{Record: corerecord.Record{ID: "A"}},
		LinkDestInstanceIDs: []string{"B"},
	}
	secB := &sectorState{
		Instance:            &mecha_record.MechaSectorInstance{Record: corerecord.Record{ID: "B"}},
		LinkDestInstanceIDs: []string{"A", "C"},
	}
	secC := &sectorState{
		Instance:            &mecha_record.MechaSectorInstance{Record: corerecord.Record{ID: "C"}},
		LinkDestInstanceIDs: []string{"B", "D"},
	}
	secD := &sectorState{
		Instance:            &mecha_record.MechaSectorInstance{Record: corerecord.Record{ID: "D"}},
		LinkDestInstanceIDs: []string{"C"},
	}
	sectors := []*sectorState{secA, secB, secC, secD}

	tests := []struct {
		name     string
		fromID   string
		toID     string
		expected int
	}{
		{name: "same sector returns 0", fromID: "A", toID: "A", expected: 0},
		{name: "adjacent returns 1", fromID: "A", toID: "B", expected: 1},
		{name: "two hops returns 2", fromID: "A", toID: "C", expected: 2},
		{name: "three hops returns 3", fromID: "A", toID: "D", expected: 3},
		{name: "four hops (beyond 3) returns 999", fromID: "B", toID: "D", expected: 2},
		{name: "unreachable sector returns 999", fromID: "A", toID: "Z", expected: 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dist := rangeDistance(tt.fromID, tt.toID, sectors)
			require.Equal(t, tt.expected, dist)
		})
	}
}

func TestWeaponCanFire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		rangeBand string
		distance  int
		expected  bool
	}{
		{name: "short range, same sector", rangeBand: mecha_record.WeaponRangeBandShort, distance: 0, expected: true},
		{name: "medium range, same sector", rangeBand: mecha_record.WeaponRangeBandMedium, distance: 0, expected: true},
		{name: "long range, same sector", rangeBand: mecha_record.WeaponRangeBandLong, distance: 0, expected: true},
		{name: "short range, adjacent", rangeBand: mecha_record.WeaponRangeBandShort, distance: 1, expected: false},
		{name: "medium range, adjacent", rangeBand: mecha_record.WeaponRangeBandMedium, distance: 1, expected: true},
		{name: "long range, adjacent", rangeBand: mecha_record.WeaponRangeBandLong, distance: 1, expected: true},
		{name: "short range, 2 sectors", rangeBand: mecha_record.WeaponRangeBandShort, distance: 2, expected: false},
		{name: "medium range, 2 sectors", rangeBand: mecha_record.WeaponRangeBandMedium, distance: 2, expected: false},
		{name: "long range, 2 sectors", rangeBand: mecha_record.WeaponRangeBandLong, distance: 2, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := weaponCanFire(tt.rangeBand, tt.distance)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestHitChance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pilotSkill int
		expected   int
	}{
		{name: "pilot skill 0 gives 50%", pilotSkill: 0, expected: 50},
		{name: "pilot skill 1 gives 55%", pilotSkill: 1, expected: 55},
		{name: "pilot skill 5 gives 75%", pilotSkill: 5, expected: 75},
		{name: "pilot skill 9 gives 95%", pilotSkill: 9, expected: 95},
		{name: "pilot skill 10 caps at 95%", pilotSkill: 10, expected: 95},
		{name: "pilot skill 20 caps at 95%", pilotSkill: 20, expected: 95},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chance := hitChance(tt.pilotSkill)
			require.Equal(t, tt.expected, chance)
		})
	}
}
