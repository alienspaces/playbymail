package mecha_game_record

import "encoding/json"

// marshalConfigForWrite resolves the JSON bytes to persist for a loadout
// column (weapon_config or equipment_config). The decoded struct slice is
// preferred when it carries entries — that represents fresh data from an
// explicit mutation. When the struct is empty, the raw JSON bytes loaded
// from the DB are used as a fallback so that read-modify-write cycles (e.g.
// updating a mech's position or heat) don't wipe the persisted loadout.
//
// The two loadout fields on MechaGameSquadMech and MechaGameMechInstance are
// intentionally dual-surfaced: the `db:"weapon_config"` []byte field is
// populated by pgx scans, while the `db:"-"` struct slice is populated only
// when calling code decodes it explicitly. Any write path that only mutated
// non-loadout fields would previously marshal nil → "null" into the jsonb
// column, clearing the stored loadout.
func marshalConfigForWrite[T any](decoded []T, raw []byte) string {
	if len(decoded) > 0 {
		b, _ := json.Marshal(decoded)
		return string(b)
	}
	if len(raw) > 0 {
		return string(raw)
	}
	return "[]"
}
