package mecha_game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corerecord "gitlab.com/alienspaces/playbymail/core/record"
	corelog "gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// noopLogger satisfies corelog.Logger for unit tests that exercise logging paths.
type noopLogger struct{}

func (n *noopLogger) NewInstance() (corelog.Logger, error)           { return n, nil }
func (n *noopLogger) Context(_, _ string)                            {}
func (n *noopLogger) WithApplicationContext(_ string) corelog.Logger { return n }
func (n *noopLogger) WithDurationContext(_ string) corelog.Logger    { return n }
func (n *noopLogger) WithPackageContext(_ string) corelog.Logger     { return n }
func (n *noopLogger) WithFunctionContext(_ string) corelog.Logger    { return n }
func (n *noopLogger) Debug(_ string, _ ...any)                       {}
func (n *noopLogger) Info(_ string, _ ...any)                        {}
func (n *noopLogger) Warn(_ string, _ ...any)                        {}
func (n *noopLogger) Error(_ string, _ ...any)                       {}

var testLogger corelog.Logger = &noopLogger{}

// Tests for pure helper functions in combat_resolution.go.
// Integration tests for resolveCombat require a full DB harness and are omitted here.

func TestRangeDistance(t *testing.T) {
	t.Parallel()

	// Build a simple linear map: A <-> B <-> C <-> D
	secA := &sectorState{
		Instance:            &mecha_game_record.MechaGameSectorInstance{Record: corerecord.Record{ID: "A"}},
		LinkDestInstanceIDs: []string{"B"},
	}
	secB := &sectorState{
		Instance:            &mecha_game_record.MechaGameSectorInstance{Record: corerecord.Record{ID: "B"}},
		LinkDestInstanceIDs: []string{"A", "C"},
	}
	secC := &sectorState{
		Instance:            &mecha_game_record.MechaGameSectorInstance{Record: corerecord.Record{ID: "C"}},
		LinkDestInstanceIDs: []string{"B", "D"},
	}
	secD := &sectorState{
		Instance:            &mecha_game_record.MechaGameSectorInstance{Record: corerecord.Record{ID: "D"}},
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
		{name: "two hops B to D returns 2", fromID: "B", toID: "D", expected: 2},
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
		{name: "short range, same sector", rangeBand: mecha_game_record.WeaponRangeBandShort, distance: 0, expected: true},
		{name: "medium range, same sector", rangeBand: mecha_game_record.WeaponRangeBandMedium, distance: 0, expected: true},
		{name: "long range, same sector", rangeBand: mecha_game_record.WeaponRangeBandLong, distance: 0, expected: false},
		{name: "short range, adjacent", rangeBand: mecha_game_record.WeaponRangeBandShort, distance: 1, expected: false},
		{name: "medium range, adjacent", rangeBand: mecha_game_record.WeaponRangeBandMedium, distance: 1, expected: true},
		{name: "long range, adjacent", rangeBand: mecha_game_record.WeaponRangeBandLong, distance: 1, expected: true},
		{name: "short range, 2 sectors", rangeBand: mecha_game_record.WeaponRangeBandShort, distance: 2, expected: false},
		{name: "medium range, 2 sectors", rangeBand: mecha_game_record.WeaponRangeBandMedium, distance: 2, expected: false},
		{name: "long range, 2 sectors", rangeBand: mecha_game_record.WeaponRangeBandLong, distance: 2, expected: true},
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
		name          string
		pilotSkill    int
		coverModifier int
		expected      int
	}{
		{name: "pilot skill 0 gives 50%", pilotSkill: 0, coverModifier: 0, expected: 50},
		{name: "pilot skill 1 gives 55%", pilotSkill: 1, coverModifier: 0, expected: 55},
		{name: "pilot skill 5 gives 75%", pilotSkill: 5, coverModifier: 0, expected: 75},
		{name: "pilot skill 9 gives 95%", pilotSkill: 9, coverModifier: 0, expected: 95},
		{name: "pilot skill 10 caps at 95%", pilotSkill: 10, coverModifier: 0, expected: 95},
		{name: "pilot skill 20 caps at 95%", pilotSkill: 20, coverModifier: 0, expected: 95},
		{name: "cover -10 reduces hit chance", pilotSkill: 5, coverModifier: -10, expected: 65},
		{name: "cover floors at 0%", pilotSkill: 0, coverModifier: -100, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chance := hitChance(tt.pilotSkill, tt.coverModifier)
			require.Equal(t, tt.expected, chance)
		})
	}
}

func TestAccumulateDamage(t *testing.T) {
	t.Parallel()

	t.Run("first attacker creates entry", func(t *testing.T) {
		t.Parallel()
		dm := map[string]*pendingDamage{}
		snaps := map[string]*mechSnapshot{}
		accumulateDamage("mech1", 10, snaps, dm)
		require.NotNil(t, dm["mech1"])
		assert.Equal(t, 10, dm["mech1"].rawTotal)
	})

	t.Run("second attacker accumulates damage on same target", func(t *testing.T) {
		t.Parallel()
		dm := map[string]*pendingDamage{}
		snaps := map[string]*mechSnapshot{}
		accumulateDamage("mech1", 10, snaps, dm)
		accumulateDamage("mech1", 5, snaps, dm)
		assert.Equal(t, 15, dm["mech1"].rawTotal)
	})

	t.Run("independent targets tracked separately", func(t *testing.T) {
		t.Parallel()
		dm := map[string]*pendingDamage{}
		snaps := map[string]*mechSnapshot{}
		accumulateDamage("mech1", 8, snaps, dm)
		accumulateDamage("mech2", 3, snaps, dm)
		assert.Equal(t, 8, dm["mech1"].rawTotal)
		assert.Equal(t, 3, dm["mech2"].rawTotal)
	})
}

func TestApplyPendingDamage(t *testing.T) {
	t.Parallel()

	mechaGame := &MechaGame{}

	makeSnap := func(id, callsign, squadID string, armor, structure int) *mechSnapshot {
		return &mechSnapshot{
			Instance: &mecha_game_record.MechaGameMechInstance{
				Record:           corerecord.Record{ID: id},
				Callsign:         callsign,
				CurrentArmor:     armor,
				CurrentStructure: structure,
				Status:           mecha_game_record.MechInstanceStatusOperational,
			},
			SquadInstanceID: squadID,
		}
	}

	t.Run("damage absorbed fully by armor", func(t *testing.T) {
		t.Parallel()
		inst := &mecha_game_record.MechaGameMechInstance{
			Record:           corerecord.Record{ID: "m1"},
			Callsign:         "Titan",
			CurrentArmor:     20,
			CurrentStructure: 10,
			Status:           mecha_game_record.MechInstanceStatusOperational,
		}
		snap := &mechSnapshot{Instance: inst, SquadInstanceID: "squad1"}
		snapshots := map[string]*mechSnapshot{"m1": snap}
		dm := map[string]*pendingDamage{"m1": {rawTotal: 5}}
		eventsBySquad := map[string][]turnsheet.TurnEvent{}

		mechaGame.applyPendingDamage(testLogger, dm, snapshots, nil, eventsBySquad)

		assert.Equal(t, 15, inst.CurrentArmor)
		assert.Equal(t, 10, inst.CurrentStructure)
		assert.Equal(t, mecha_game_record.MechInstanceStatusOperational, inst.Status)
	})

	t.Run("damage overflows armor into structure causes damaged status", func(t *testing.T) {
		t.Parallel()
		inst := &mecha_game_record.MechaGameMechInstance{
			Record:           corerecord.Record{ID: "m2"},
			Callsign:         "Wraith",
			CurrentArmor:     5,
			CurrentStructure: 10,
			Status:           mecha_game_record.MechInstanceStatusOperational,
		}
		snap := &mechSnapshot{Instance: inst, SquadInstanceID: "squad1"}
		snapshots := map[string]*mechSnapshot{"m2": snap}
		dm := map[string]*pendingDamage{"m2": {rawTotal: 8}}
		eventsBySquad := map[string][]turnsheet.TurnEvent{}

		mechaGame.applyPendingDamage(testLogger, dm, snapshots, nil, eventsBySquad)

		assert.Equal(t, 0, inst.CurrentArmor)
		assert.Equal(t, 7, inst.CurrentStructure)
		assert.Equal(t, mecha_game_record.MechInstanceStatusDamaged, inst.Status)
	})

	t.Run("focus fire from two attackers pools damage before armor split", func(t *testing.T) {
		t.Parallel()
		// Armor = 10; attacker A deals 7, attacker B deals 7; total = 14.
		// With pooled damage: armor absorbs 10, structure takes 4.
		// Without pooling (old bug): each attacker would absorb 7 armor independently → 0 structure damage.
		inst := &mecha_game_record.MechaGameMechInstance{
			Record:           corerecord.Record{ID: "m3"},
			Callsign:         "IronFist",
			CurrentArmor:     10,
			CurrentStructure: 15,
			Status:           mecha_game_record.MechInstanceStatusOperational,
		}
		snap := &mechSnapshot{Instance: inst, SquadInstanceID: "squad1"}
		snapshots := map[string]*mechSnapshot{"m3": snap}
		dm := map[string]*pendingDamage{}
		accumulateDamage("m3", 7, snapshots, dm)
		accumulateDamage("m3", 7, snapshots, dm)
		eventsBySquad := map[string][]turnsheet.TurnEvent{}

		mechaGame.applyPendingDamage(testLogger, dm, snapshots, nil, eventsBySquad)

		assert.Equal(t, 0, inst.CurrentArmor)
		assert.Equal(t, 11, inst.CurrentStructure)
		assert.Equal(t, mecha_game_record.MechInstanceStatusDamaged, inst.Status)
	})

	t.Run("structure reduced to zero causes destroyed status", func(t *testing.T) {
		t.Parallel()
		inst := &mecha_game_record.MechaGameMechInstance{
			Record:           corerecord.Record{ID: "m4"},
			Callsign:         "Rattler",
			CurrentArmor:     2,
			CurrentStructure: 3,
			Status:           mecha_game_record.MechInstanceStatusOperational,
		}
		snap := &mechSnapshot{Instance: inst, SquadInstanceID: "squad1"}
		snapshots := map[string]*mechSnapshot{"m4": snap}
		attacks := []AttackDeclaration{{AttackerMechInstanceID: "atk", TargetMechInstanceID: "m4"}}
		snapshots["atk"] = makeSnap("atk", "Hunter", "squad2", 0, 10)
		dm := map[string]*pendingDamage{"m4": {rawTotal: 20}}
		eventsBySquad := map[string][]turnsheet.TurnEvent{}

		mechaGame.applyPendingDamage(testLogger, dm, snapshots, attacks, eventsBySquad)

		assert.Equal(t, 0, inst.CurrentStructure)
		assert.Equal(t, mecha_game_record.MechInstanceStatusDestroyed, inst.Status)
	})
}

func TestPilotSkillThresholds(t *testing.T) {
	t.Parallel()

	t.Run("thresholds are strictly increasing", func(t *testing.T) {
		t.Parallel()
		for i := 1; i < len(pilotSkillThresholds); i++ {
			assert.Greater(t, pilotSkillThresholds[i], pilotSkillThresholds[i-1],
				"threshold[%d] must be > threshold[%d]", i, i-1)
		}
	})

	t.Run("skill level 0 requires 0 XP", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 0, pilotSkillThresholds[0])
	})

	t.Run("has exactly 10 levels (0-9)", func(t *testing.T) {
		t.Parallel()
		assert.Len(t, pilotSkillThresholds, 10)
	})

	t.Run("pilot advances skill from 0 to 1 at threshold", func(t *testing.T) {
		t.Parallel()
		xp := pilotSkillThresholds[1]
		skill := 0
		for nextSkill := skill + 1; nextSkill < len(pilotSkillThresholds); nextSkill++ {
			if xp >= pilotSkillThresholds[nextSkill] {
				skill = nextSkill
			} else {
				break
			}
		}
		assert.Equal(t, 1, skill)
	})

	t.Run("pilot does not advance skill with XP below threshold", func(t *testing.T) {
		t.Parallel()
		xp := pilotSkillThresholds[1] - 1
		skill := 0
		for nextSkill := skill + 1; nextSkill < len(pilotSkillThresholds); nextSkill++ {
			if xp >= pilotSkillThresholds[nextSkill] {
				skill = nextSkill
			} else {
				break
			}
		}
		assert.Equal(t, 0, skill)
	})
}
