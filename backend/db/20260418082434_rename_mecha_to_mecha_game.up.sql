-- Rename mecha_* tables to mecha_game_* to align with the adventure_game_* naming convention.
-- Also renames FK columns, constraints, indexes, and stored string values.

BEGIN;

-- ============================================================================
-- 1. Rename tables (order matters: leaf tables first to avoid FK conflicts)
-- ============================================================================

ALTER TABLE public.mecha_mech_instance     RENAME TO mecha_game_mech_instance;
ALTER TABLE public.mecha_turn_sheet        RENAME TO mecha_game_turn_sheet;
ALTER TABLE public.mecha_squad_mech        RENAME TO mecha_game_squad_mech;
ALTER TABLE public.mecha_squad_instance    RENAME TO mecha_game_squad_instance;
ALTER TABLE public.mecha_sector_instance   RENAME TO mecha_game_sector_instance;
ALTER TABLE public.mecha_sector_link       RENAME TO mecha_game_sector_link;
ALTER TABLE public.mecha_squad             RENAME TO mecha_game_squad;
ALTER TABLE public.mecha_computer_opponent RENAME TO mecha_game_computer_opponent;
ALTER TABLE public.mecha_chassis           RENAME TO mecha_game_chassis;
ALTER TABLE public.mecha_weapon            RENAME TO mecha_game_weapon;
ALTER TABLE public.mecha_sector            RENAME TO mecha_game_sector;

-- ============================================================================
-- 2. Rename FK columns
-- ============================================================================

-- mecha_game_sector_link
ALTER TABLE public.mecha_game_sector_link
    RENAME COLUMN from_mecha_sector_id TO from_mecha_game_sector_id;
ALTER TABLE public.mecha_game_sector_link
    RENAME COLUMN to_mecha_sector_id TO to_mecha_game_sector_id;

-- mecha_game_sector_instance
ALTER TABLE public.mecha_game_sector_instance
    RENAME COLUMN mecha_sector_id TO mecha_game_sector_id;

-- mecha_game_squad_mech
ALTER TABLE public.mecha_game_squad_mech
    RENAME COLUMN mecha_squad_id TO mecha_game_squad_id;
ALTER TABLE public.mecha_game_squad_mech
    RENAME COLUMN mecha_chassis_id TO mecha_game_chassis_id;

-- mecha_game_squad_instance
ALTER TABLE public.mecha_game_squad_instance
    RENAME COLUMN mecha_squad_id TO mecha_game_squad_id;
ALTER TABLE public.mecha_game_squad_instance
    RENAME COLUMN mecha_computer_opponent_id TO mecha_game_computer_opponent_id;

-- mecha_game_mech_instance
ALTER TABLE public.mecha_game_mech_instance
    RENAME COLUMN mecha_squad_instance_id TO mecha_game_squad_instance_id;
ALTER TABLE public.mecha_game_mech_instance
    RENAME COLUMN mecha_sector_instance_id TO mecha_game_sector_instance_id;
ALTER TABLE public.mecha_game_mech_instance
    RENAME COLUMN mecha_chassis_id TO mecha_game_chassis_id;

-- mecha_game_turn_sheet
ALTER TABLE public.mecha_game_turn_sheet
    RENAME COLUMN mecha_squad_instance_id TO mecha_game_squad_instance_id;

-- ============================================================================
-- 3. Rename constraints on mecha_game_chassis (formerly mecha_chassis)
-- ============================================================================

ALTER TABLE public.mecha_game_chassis
    RENAME CONSTRAINT mecha_chassis_game_id_fkey TO mecha_game_chassis_game_id_fkey;
ALTER TABLE public.mecha_game_chassis
    RENAME CONSTRAINT mecha_chassis_class_check TO mecha_game_chassis_class_check;

-- ============================================================================
-- 4. Rename constraints on mecha_game_weapon (formerly mecha_weapon)
-- ============================================================================

ALTER TABLE public.mecha_game_weapon
    RENAME CONSTRAINT mecha_weapon_game_id_fkey TO mecha_game_weapon_game_id_fkey;
ALTER TABLE public.mecha_game_weapon
    RENAME CONSTRAINT mecha_weapon_range_band_check TO mecha_game_weapon_range_band_check;
ALTER TABLE public.mecha_game_weapon
    RENAME CONSTRAINT mecha_weapon_mount_size_check TO mecha_game_weapon_mount_size_check;

-- ============================================================================
-- 5. Rename constraints on mecha_game_sector (formerly mecha_sector)
-- ============================================================================

ALTER TABLE public.mecha_game_sector
    RENAME CONSTRAINT mecha_sector_game_id_fkey TO mecha_game_sector_game_id_fkey;
ALTER TABLE public.mecha_game_sector
    RENAME CONSTRAINT mecha_sector_terrain_type_check TO mecha_game_sector_terrain_type_check;

-- ============================================================================
-- 6. Rename constraints on mecha_game_sector_link (formerly mecha_sector_link)
-- ============================================================================

ALTER TABLE public.mecha_game_sector_link
    RENAME CONSTRAINT mecha_sector_link_unique TO mecha_game_sector_link_unique;
ALTER TABLE public.mecha_game_sector_link
    RENAME CONSTRAINT mecha_sector_link_game_id_fkey TO mecha_game_sector_link_game_id_fkey;
ALTER TABLE public.mecha_game_sector_link
    RENAME CONSTRAINT mecha_sector_link_from_sector_fkey TO mecha_game_sector_link_from_sector_fkey;
ALTER TABLE public.mecha_game_sector_link
    RENAME CONSTRAINT mecha_sector_link_to_sector_fkey TO mecha_game_sector_link_to_sector_fkey;

-- ============================================================================
-- 7. Rename constraints on mecha_game_computer_opponent
-- ============================================================================

ALTER TABLE public.mecha_game_computer_opponent
    RENAME CONSTRAINT mecha_computer_opponent_game_id_fkey TO mecha_game_computer_opponent_game_id_fkey;

-- ============================================================================
-- 8. Rename constraints on mecha_game_squad (formerly mecha_squad)
-- ============================================================================

ALTER TABLE public.mecha_game_squad
    RENAME CONSTRAINT mecha_squad_name_unique TO mecha_game_squad_name_unique;
ALTER TABLE public.mecha_game_squad
    RENAME CONSTRAINT mecha_squad_game_id_fkey TO mecha_game_squad_game_id_fkey;
ALTER TABLE public.mecha_game_squad
    RENAME CONSTRAINT mecha_squad_type_check TO mecha_game_squad_type_check;

-- ============================================================================
-- 9. Rename constraints on mecha_game_squad_mech (formerly mecha_squad_mech)
-- ============================================================================

ALTER TABLE public.mecha_game_squad_mech
    RENAME CONSTRAINT mecha_squad_mech_callsign_unique TO mecha_game_squad_mech_callsign_unique;
ALTER TABLE public.mecha_game_squad_mech
    RENAME CONSTRAINT mecha_squad_mech_game_id_fkey TO mecha_game_squad_mech_game_id_fkey;
ALTER TABLE public.mecha_game_squad_mech
    RENAME CONSTRAINT mecha_squad_mech_squad_id_fkey TO mecha_game_squad_mech_squad_id_fkey;
ALTER TABLE public.mecha_game_squad_mech
    RENAME CONSTRAINT mecha_squad_mech_chassis_id_fkey TO mecha_game_squad_mech_chassis_id_fkey;

-- ============================================================================
-- 10. Rename constraints on mecha_game_squad_instance (formerly mecha_squad_instance)
-- ============================================================================

ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_game_id_fkey
                   TO mecha_game_squad_instance_game_id_fkey;
ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_game_instance_id_fkey
                   TO mecha_game_squad_instance_game_instance_id_fkey;
ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_squad_id_fkey
                   TO mecha_game_squad_instance_squad_id_fkey;
ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_subscription_instance_id_fkey
                   TO mecha_game_squad_instance_subscription_instance_id_fkey;
ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_computer_opponent_id_fkey
                   TO mecha_game_squad_instance_computer_opponent_id_fkey;

-- ============================================================================
-- 11. Rename constraints on mecha_game_sector_instance
-- ============================================================================

ALTER TABLE public.mecha_game_sector_instance
    RENAME CONSTRAINT mecha_sector_instance_game_id_fkey
                   TO mecha_game_sector_instance_game_id_fkey;
ALTER TABLE public.mecha_game_sector_instance
    RENAME CONSTRAINT mecha_sector_instance_game_instance_id_fkey
                   TO mecha_game_sector_instance_game_instance_id_fkey;
ALTER TABLE public.mecha_game_sector_instance
    RENAME CONSTRAINT mecha_sector_instance_sector_id_fkey
                   TO mecha_game_sector_instance_sector_id_fkey;

-- ============================================================================
-- 12. Rename constraints on mecha_game_mech_instance
-- ============================================================================

ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_mech_instance_status_check
                   TO mecha_game_mech_instance_status_check;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_mech_instance_game_id_fkey
                   TO mecha_game_mech_instance_game_id_fkey;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_mech_instance_game_instance_id_fkey
                   TO mecha_game_mech_instance_game_instance_id_fkey;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_mech_instance_squad_instance_id_fkey
                   TO mecha_game_mech_instance_squad_instance_id_fkey;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_mech_instance_sector_instance_id_fkey
                   TO mecha_game_mech_instance_sector_instance_id_fkey;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_mech_instance_chassis_id_fkey
                   TO mecha_game_mech_instance_chassis_id_fkey;

-- ============================================================================
-- 13. Rename constraints on mecha_game_turn_sheet
-- ============================================================================

-- Unique constraint references renamed column; drop and recreate for a clean definition.
ALTER TABLE public.mecha_game_turn_sheet
    DROP CONSTRAINT mecha_turn_sheet_unique;
ALTER TABLE public.mecha_game_turn_sheet
    ADD CONSTRAINT mecha_game_turn_sheet_unique
        UNIQUE (mecha_game_squad_instance_id, game_turn_sheet_id);

ALTER TABLE public.mecha_game_turn_sheet
    RENAME CONSTRAINT mecha_turn_sheet_game_id_fkey
                   TO mecha_game_turn_sheet_game_id_fkey;
ALTER TABLE public.mecha_game_turn_sheet
    RENAME CONSTRAINT mecha_turn_sheet_squad_instance_id_fkey
                   TO mecha_game_turn_sheet_squad_instance_id_fkey;
ALTER TABLE public.mecha_game_turn_sheet
    RENAME CONSTRAINT mecha_turn_sheet_game_turn_sheet_id_fkey
                   TO mecha_game_turn_sheet_game_turn_sheet_id_fkey;

-- ============================================================================
-- 14. Rename indexes
-- ============================================================================

ALTER INDEX public.idx_mecha_chassis_game_id
    RENAME TO idx_mecha_game_chassis_game_id;

ALTER INDEX public.idx_mecha_weapon_game_id
    RENAME TO idx_mecha_game_weapon_game_id;

ALTER INDEX public.idx_mecha_sector_game_id
    RENAME TO idx_mecha_game_sector_game_id;
ALTER INDEX public.idx_mecha_sector_is_starting
    RENAME TO idx_mecha_game_sector_is_starting;

ALTER INDEX public.idx_mecha_sector_link_from_sector
    RENAME TO idx_mecha_game_sector_link_from_sector;
ALTER INDEX public.idx_mecha_sector_link_to_sector
    RENAME TO idx_mecha_game_sector_link_to_sector;

-- Partial unique index on mecha_game_squad references column; drop and recreate.
DROP INDEX IF EXISTS public.idx_mecha_squad_starter_unique;
CREATE UNIQUE INDEX idx_mecha_game_squad_starter_unique
    ON public.mecha_game_squad (game_id)
    WHERE squad_type = 'starter' AND deleted_at IS NULL;

ALTER INDEX public.idx_mecha_squad_game_id
    RENAME TO idx_mecha_game_squad_game_id;

ALTER INDEX public.idx_mecha_squad_mech_squad_id
    RENAME TO idx_mecha_game_squad_mech_squad_id;
ALTER INDEX public.idx_mecha_squad_mech_chassis_id
    RENAME TO idx_mecha_game_squad_mech_chassis_id;

ALTER INDEX public.idx_mecha_squad_instance_game_instance
    RENAME TO idx_mecha_game_squad_instance_game_instance;
ALTER INDEX public.idx_mecha_squad_instance_squad_id
    RENAME TO idx_mecha_game_squad_instance_squad_id;
ALTER INDEX public.idx_mecha_squad_instance_computer_opponent
    RENAME TO idx_mecha_game_squad_instance_computer_opponent;

ALTER INDEX public.idx_mecha_sector_instance_game_instance
    RENAME TO idx_mecha_game_sector_instance_game_instance;
ALTER INDEX public.idx_mecha_sector_instance_sector_id
    RENAME TO idx_mecha_game_sector_instance_sector_id;

ALTER INDEX public.idx_mecha_mech_instance_game_instance
    RENAME TO idx_mecha_game_mech_instance_game_instance;
ALTER INDEX public.idx_mecha_mech_instance_squad_instance
    RENAME TO idx_mecha_game_mech_instance_squad_instance;
ALTER INDEX public.idx_mecha_mech_instance_sector_instance
    RENAME TO idx_mecha_game_mech_instance_sector_instance;

ALTER INDEX public.idx_mecha_turn_sheet_game_id
    RENAME TO idx_mecha_game_turn_sheet_game_id;
ALTER INDEX public.idx_mecha_turn_sheet_squad_instance
    RENAME TO idx_mecha_game_turn_sheet_squad_instance;
ALTER INDEX public.idx_mecha_turn_sheet_game_turn_sheet
    RENAME TO idx_mecha_game_turn_sheet_game_turn_sheet;

-- ============================================================================
-- 15. Update stored string values (turn sheet types)
-- ============================================================================

UPDATE public.game_turn_sheet
    SET sheet_type = 'mecha_game_join_game'
    WHERE sheet_type = 'mecha_join_game';

UPDATE public.game_turn_sheet
    SET sheet_type = 'mecha_game_orders'
    WHERE sheet_type = 'mecha_orders';

UPDATE public.game_turn_sheet
    SET sheet_type = 'mecha_game_squad_management'
    WHERE sheet_type = 'mecha_squad_management';

UPDATE public.game_image
    SET turn_sheet_type = 'mecha_game_join_game'
    WHERE turn_sheet_type = 'mecha_join_game';

UPDATE public.game_image
    SET turn_sheet_type = 'mecha_game_orders'
    WHERE turn_sheet_type = 'mecha_orders';

UPDATE public.game_image
    SET turn_sheet_type = 'mecha_game_squad_management'
    WHERE turn_sheet_type = 'mecha_squad_management';

COMMIT;
