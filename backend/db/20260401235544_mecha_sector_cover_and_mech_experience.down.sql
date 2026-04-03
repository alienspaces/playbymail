-- Restore cover_modifier on mecha_sector_link.
ALTER TABLE mecha_sector_link
    ADD COLUMN IF NOT EXISTS cover_modifier INTEGER NOT NULL DEFAULT 0;

-- Remove cover_modifier from mecha_sector.
ALTER TABLE mecha_sector
    DROP COLUMN IF EXISTS cover_modifier;

-- Remove experience_points from mecha_mech_instance.
ALTER TABLE mecha_mech_instance
    DROP COLUMN IF EXISTS experience_points;
