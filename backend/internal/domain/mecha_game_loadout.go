package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

// Loadout capacity and fit-checking.
//
// A chassis exposes three aggregate slot bands: small, medium, and large. Any
// item that needs a slot (today: weapons) declares a mount size in the same
// vocabulary. fitsLoadout greedily places items into the chassis's slot budget
// with "upward spillover": smaller items can occupy larger slots when their
// native band is full, but never the reverse. Concretely:
//
//   - large items can only use large slots
//   - medium items prefer medium slots, spill into large
//   - small items prefer small slots, spill into medium, then large
//
// The helper is deliberately decoupled from any specific record type so future
// equipment types (ammunition, heat sinks, etc.) can feed into the same budget
// without touching this file — they just need a mount size and a label.

// LoadoutCapacity is the aggregate slot budget exposed by a chassis.
type LoadoutCapacity struct {
	Small  int
	Medium int
	Large  int
}

// LoadoutCapacityFromChassis extracts the slot budget from a chassis record.
func LoadoutCapacityFromChassis(chassis *mecha_game_record.MechaGameChassis) LoadoutCapacity {
	if chassis == nil {
		return LoadoutCapacity{}
	}
	return LoadoutCapacity{
		Small:  chassis.SmallSlots,
		Medium: chassis.MediumSlots,
		Large:  chassis.LargeSlots,
	}
}

// Mountable represents a single item that must occupy one slot. Label is used
// solely for error messages so players/designers can identify which item did
// not fit.
type Mountable struct {
	MountSize string
	Label     string
}

// MountablesFromWeaponConfig converts a persisted weapon loadout into the
// mountable shape used by fitsLoadout. Entries whose weapon_id is missing from
// weaponsByID are treated as unknown-size small items so they are not silently
// dropped; this should only happen in practice if a weapon has been hard
// deleted out from under a loadout, which the upstream domain guards against.
func MountablesFromWeaponConfig(entries []mecha_game_record.WeaponConfigEntry, weaponsByID map[string]*mecha_game_record.MechaGameWeapon) []Mountable {
	if len(entries) == 0 {
		return nil
	}
	out := make([]Mountable, 0, len(entries))
	for _, e := range entries {
		w := weaponsByID[e.WeaponID]
		size := ""
		name := e.WeaponID
		if w != nil {
			size = w.MountSize
			if w.Name != "" {
				name = w.Name
			}
		}
		out = append(out, Mountable{MountSize: size, Label: name})
	}
	return out
}

// ValidateWeaponLoadoutFits is the public entry point used outside this
// package (notably from the turn-sheet processors) to check a proposed weapon
// loadout against a chassis before persisting it. Returns nil when the
// loadout fits, or a descriptive error when an item cannot be placed.
func ValidateWeaponLoadoutFits(chassis *mecha_game_record.MechaGameChassis, entries []mecha_game_record.WeaponConfigEntry, weaponsByID map[string]*mecha_game_record.MechaGameWeapon) error {
	if chassis == nil || len(entries) == 0 {
		return nil
	}
	items := MountablesFromWeaponConfig(entries, weaponsByID)
	return fitsLoadout(LoadoutCapacityFromChassis(chassis), items)
}

// fitsLoadout reports whether every item fits the chassis capacity under the
// upward-spillover rule. Items are placed in order large → medium → small so a
// small item never steals a large slot that a subsequent large item would have
// needed. Returns a descriptive error naming the first item that cannot be
// placed and the slot size it needed.
func fitsLoadout(cap LoadoutCapacity, items []Mountable) error {
	remaining := cap

	// Pass 1: large items must take large slots.
	for _, it := range items {
		if it.MountSize != mecha_game_record.WeaponMountSizeLarge {
			continue
		}
		if remaining.Large <= 0 {
			return fmt.Errorf("%s does not fit: chassis has no available large slots", describeMountable(it))
		}
		remaining.Large--
	}

	// Pass 2: medium items prefer medium, spill into large.
	for _, it := range items {
		if it.MountSize != mecha_game_record.WeaponMountSizeMedium {
			continue
		}
		switch {
		case remaining.Medium > 0:
			remaining.Medium--
		case remaining.Large > 0:
			remaining.Large--
		default:
			return fmt.Errorf("%s does not fit: chassis has no available medium or large slots", describeMountable(it))
		}
	}

	// Pass 3: small items prefer small, then medium, then large.
	for _, it := range items {
		if it.MountSize != mecha_game_record.WeaponMountSizeSmall {
			continue
		}
		switch {
		case remaining.Small > 0:
			remaining.Small--
		case remaining.Medium > 0:
			remaining.Medium--
		case remaining.Large > 0:
			remaining.Large--
		default:
			return fmt.Errorf("%s does not fit: chassis has no available slots", describeMountable(it))
		}
	}

	// Pass 4: anything with an unrecognised mount size is refused so bad data
	// surfaces early rather than silently consuming a slot.
	for _, it := range items {
		switch it.MountSize {
		case mecha_game_record.WeaponMountSizeSmall,
			mecha_game_record.WeaponMountSizeMedium,
			mecha_game_record.WeaponMountSizeLarge:
			continue
		default:
			return fmt.Errorf("%s has unknown mount size >%s<", describeMountable(it), it.MountSize)
		}
	}

	return nil
}

func describeMountable(m Mountable) string {
	if m.Label != "" {
		return fmt.Sprintf("item %q", m.Label)
	}
	return "item"
}
