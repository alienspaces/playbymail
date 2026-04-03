-- Add cover_modifier to mecha_sector (replaces the field previously on mecha_sector_link).
ALTER TABLE mecha_sector
    ADD COLUMN cover_modifier INTEGER NOT NULL DEFAULT 0;

-- Drop cover_modifier from mecha_sector_link; the field now lives on the sector.
ALTER TABLE mecha_sector_link
    DROP COLUMN IF EXISTS cover_modifier;

-- Add experience_points to mecha_mech_instance for the pilot progression system.
ALTER TABLE mecha_mech_instance
    ADD COLUMN experience_points INTEGER NOT NULL DEFAULT 0;
