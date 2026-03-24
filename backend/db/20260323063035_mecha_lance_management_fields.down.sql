ALTER TABLE mecha_lance_instance
    DROP COLUMN IF EXISTS last_turn_events,
    DROP COLUMN IF EXISTS supply_points;

ALTER TABLE mecha_mech_instance
    DROP COLUMN IF EXISTS weapon_config,
    DROP COLUMN IF EXISTS is_refitting;
