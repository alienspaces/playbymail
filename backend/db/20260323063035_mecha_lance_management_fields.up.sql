-- mecha_lance_instance: turn events and supply point economy
ALTER TABLE mecha_lance_instance
    ADD COLUMN last_turn_events JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN supply_points    INT   NOT NULL DEFAULT 0;

-- mecha_mech_instance: per-instance weapon loadout and refit flag
ALTER TABLE mecha_mech_instance
    ADD COLUMN weapon_config JSONB   NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN is_refitting  BOOLEAN NOT NULL DEFAULT false;
