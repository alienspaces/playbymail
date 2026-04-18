-- Revert mecha_game_* tables back to mecha_*.

BEGIN;

-- ============================================================================
-- 1. Restore stored string values
-- ============================================================================

UPDATE public.game_turn_sheet
    SET sheet_type = 'mecha_join_game'
    WHERE sheet_type = 'mecha_game_join_game';

UPDATE public.game_turn_sheet
    SET sheet_type = 'mecha_orders'
    WHERE sheet_type = 'mecha_game_orders';

UPDATE public.game_turn_sheet
    SET sheet_type = 'mecha_squad_management'
    WHERE sheet_type = 'mecha_game_squad_management';

UPDATE public.game_image
    SET turn_sheet_type = 'mecha_join_game'
    WHERE turn_sheet_type = 'mecha_game_join_game';

UPDATE public.game_image
    SET turn_sheet_type = 'mecha_orders'
    WHERE turn_sheet_type = 'mecha_game_orders';

UPDATE public.game_image
    SET turn_sheet_type = 'mecha_squad_management'
    WHERE turn_sheet_type = 'mecha_game_squad_management';

-- ============================================================================
-- 2. Restore indexes
-- ============================================================================

ALTER INDEX public.idx_mecha_game_chassis_game_id
    RENAME TO idx_mecha_chassis_game_id;

ALTER INDEX public.idx_mecha_game_weapon_game_id
    RENAME TO idx_mecha_weapon_game_id;

ALTER INDEX public.idx_mecha_game_sector_game_id
    RENAME TO idx_mecha_sector_game_id;
ALTER INDEX public.idx_mecha_game_sector_is_starting
    RENAME TO idx_mecha_sector_is_starting;

ALTER INDEX public.idx_mecha_game_sector_link_from_sector
    RENAME TO idx_mecha_sector_link_from_sector;
ALTER INDEX public.idx_mecha_game_sector_link_to_sector
    RENAME TO idx_mecha_sector_link_to_sector;

DROP INDEX IF EXISTS public.idx_mecha_game_squad_starter_unique;
CREATE UNIQUE INDEX idx_mecha_squad_starter_unique
    ON public.mecha_game_squad (game_id)
    WHERE squad_type = 'starter' AND deleted_at IS NULL;

ALTER INDEX public.idx_mecha_game_squad_game_id
    RENAME TO idx_mecha_squad_game_id;

ALTER INDEX public.idx_mecha_game_squad_mech_squad_id
    RENAME TO idx_mecha_squad_mech_squad_id;
ALTER INDEX public.idx_mecha_game_squad_mech_chassis_id
    RENAME TO idx_mecha_squad_mech_chassis_id;

ALTER INDEX public.idx_mecha_game_squad_instance_game_instance
    RENAME TO idx_mecha_squad_instance_game_instance;
ALTER INDEX public.idx_mecha_game_squad_instance_squad_id
    RENAME TO idx_mecha_squad_instance_squad_id;
ALTER INDEX public.idx_mecha_game_squad_instance_computer_opponent
    RENAME TO idx_mecha_squad_instance_computer_opponent;

ALTER INDEX public.idx_mecha_game_sector_instance_game_instance
    RENAME TO idx_mecha_sector_instance_game_instance;
ALTER INDEX public.idx_mecha_game_sector_instance_sector_id
    RENAME TO idx_mecha_sector_instance_sector_id;

ALTER INDEX public.idx_mecha_game_mech_instance_game_instance
    RENAME TO idx_mecha_mech_instance_game_instance;
ALTER INDEX public.idx_mecha_game_mech_instance_squad_instance
    RENAME TO idx_mecha_mech_instance_squad_instance;
ALTER INDEX public.idx_mecha_game_mech_instance_sector_instance
    RENAME TO idx_mecha_mech_instance_sector_instance;

ALTER INDEX public.idx_mecha_game_turn_sheet_game_id
    RENAME TO idx_mecha_turn_sheet_game_id;
ALTER INDEX public.idx_mecha_game_turn_sheet_squad_instance
    RENAME TO idx_mecha_turn_sheet_squad_instance;
ALTER INDEX public.idx_mecha_game_turn_sheet_game_turn_sheet
    RENAME TO idx_mecha_turn_sheet_game_turn_sheet;

-- ============================================================================
-- 3. Restore constraints on mecha_game_turn_sheet
-- ============================================================================

ALTER TABLE public.mecha_game_turn_sheet
    DROP CONSTRAINT mecha_game_turn_sheet_unique;
ALTER TABLE public.mecha_game_turn_sheet
    ADD CONSTRAINT mecha_turn_sheet_unique
        UNIQUE (mecha_game_squad_instance_id, game_turn_sheet_id);

ALTER TABLE public.mecha_game_turn_sheet
    RENAME CONSTRAINT mecha_game_turn_sheet_game_id_fkey
                   TO mecha_turn_sheet_game_id_fkey;
ALTER TABLE public.mecha_game_turn_sheet
    RENAME CONSTRAINT mecha_game_turn_sheet_squad_instance_id_fkey
                   TO mecha_turn_sheet_squad_instance_id_fkey;
ALTER TABLE public.mecha_game_turn_sheet
    RENAME CONSTRAINT mecha_game_turn_sheet_game_turn_sheet_id_fkey
                   TO mecha_turn_sheet_game_turn_sheet_id_fkey;

-- ============================================================================
-- 4. Restore constraints on mecha_game_mech_instance
-- ============================================================================

ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_game_mech_instance_status_check
                   TO mecha_mech_instance_status_check;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_game_mech_instance_game_id_fkey
                   TO mecha_mech_instance_game_id_fkey;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_game_mech_instance_game_instance_id_fkey
                   TO mecha_mech_instance_game_instance_id_fkey;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_game_mech_instance_squad_instance_id_fkey
                   TO mecha_mech_instance_squad_instance_id_fkey;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_game_mech_instance_sector_instance_id_fkey
                   TO mecha_mech_instance_sector_instance_id_fkey;
ALTER TABLE public.mecha_game_mech_instance
    RENAME CONSTRAINT mecha_game_mech_instance_chassis_id_fkey
                   TO mecha_mech_instance_chassis_id_fkey;

-- ============================================================================
-- 5. Restore constraints on mecha_game_sector_instance
-- ============================================================================

ALTER TABLE public.mecha_game_sector_instance
    RENAME CONSTRAINT mecha_game_sector_instance_game_id_fkey
                   TO mecha_sector_instance_game_id_fkey;
ALTER TABLE public.mecha_game_sector_instance
    RENAME CONSTRAINT mecha_game_sector_instance_game_instance_id_fkey
                   TO mecha_sector_instance_game_instance_id_fkey;
ALTER TABLE public.mecha_game_sector_instance
    RENAME CONSTRAINT mecha_game_sector_instance_sector_id_fkey
                   TO mecha_sector_instance_sector_id_fkey;

-- ============================================================================
-- 6. Restore constraints on mecha_game_squad_instance
-- ============================================================================

ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_game_squad_instance_game_id_fkey
                   TO mecha_squad_instance_game_id_fkey;
ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_game_squad_instance_game_instance_id_fkey
                   TO mecha_squad_instance_game_instance_id_fkey;
ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_game_squad_instance_squad_id_fkey
                   TO mecha_squad_instance_squad_id_fkey;
ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_game_squad_instance_subscription_instance_id_fkey
                   TO mecha_squad_instance_subscription_instance_id_fkey;
ALTER TABLE public.mecha_game_squad_instance
    RENAME CONSTRAINT mecha_game_squad_instance_computer_opponent_id_fkey
                   TO mecha_squad_instance_computer_opponent_id_fkey;

-- ============================================================================
-- 7. Restore constraints on mecha_game_squad_mech
-- ============================================================================

ALTER TABLE public.mecha_game_squad_mech
    RENAME CONSTRAINT mecha_game_squad_mech_callsign_unique TO mecha_squad_mech_callsign_unique;
ALTER TABLE public.mecha_game_squad_mech
    RENAME CONSTRAINT mecha_game_squad_mech_game_id_fkey TO mecha_squad_mech_game_id_fkey;
ALTER TABLE public.mecha_game_squad_mech
    RENAME CONSTRAINT mecha_game_squad_mech_squad_id_fkey TO mecha_squad_mech_squad_id_fkey;
ALTER TABLE public.mecha_game_squad_mech
    RENAME CONSTRAINT mecha_game_squad_mech_chassis_id_fkey TO mecha_squad_mech_chassis_id_fkey;

-- ============================================================================
-- 8. Restore constraints on mecha_game_squad
-- ============================================================================

ALTER TABLE public.mecha_game_squad
    RENAME CONSTRAINT mecha_game_squad_name_unique TO mecha_squad_name_unique;
ALTER TABLE public.mecha_game_squad
    RENAME CONSTRAINT mecha_game_squad_game_id_fkey TO mecha_squad_game_id_fkey;
ALTER TABLE public.mecha_game_squad
    RENAME CONSTRAINT mecha_game_squad_type_check TO mecha_squad_type_check;

-- ============================================================================
-- 9. Restore constraints on mecha_game_computer_opponent
-- ============================================================================

ALTER TABLE public.mecha_game_computer_opponent
    RENAME CONSTRAINT mecha_game_computer_opponent_game_id_fkey TO mecha_computer_opponent_game_id_fkey;

-- ============================================================================
-- 10. Restore constraints on mecha_game_sector_link
-- ============================================================================

ALTER TABLE public.mecha_game_sector_link
    RENAME CONSTRAINT mecha_game_sector_link_unique TO mecha_sector_link_unique;
ALTER TABLE public.mecha_game_sector_link
    RENAME CONSTRAINT mecha_game_sector_link_game_id_fkey TO mecha_sector_link_game_id_fkey;
ALTER TABLE public.mecha_game_sector_link
    RENAME CONSTRAINT mecha_game_sector_link_from_sector_fkey TO mecha_sector_link_from_sector_fkey;
ALTER TABLE public.mecha_game_sector_link
    RENAME CONSTRAINT mecha_game_sector_link_to_sector_fkey TO mecha_sector_link_to_sector_fkey;

-- ============================================================================
-- 11. Restore constraints on mecha_game_sector
-- ============================================================================

ALTER TABLE public.mecha_game_sector
    RENAME CONSTRAINT mecha_game_sector_game_id_fkey TO mecha_sector_game_id_fkey;
ALTER TABLE public.mecha_game_sector
    RENAME CONSTRAINT mecha_game_sector_terrain_type_check TO mecha_sector_terrain_type_check;

-- ============================================================================
-- 12. Restore constraints on mecha_game_weapon
-- ============================================================================

ALTER TABLE public.mecha_game_weapon
    RENAME CONSTRAINT mecha_game_weapon_game_id_fkey TO mecha_weapon_game_id_fkey;
ALTER TABLE public.mecha_game_weapon
    RENAME CONSTRAINT mecha_game_weapon_range_band_check TO mecha_weapon_range_band_check;
ALTER TABLE public.mecha_game_weapon
    RENAME CONSTRAINT mecha_game_weapon_mount_size_check TO mecha_weapon_mount_size_check;

-- ============================================================================
-- 13. Restore constraints on mecha_game_chassis
-- ============================================================================

ALTER TABLE public.mecha_game_chassis
    RENAME CONSTRAINT mecha_game_chassis_game_id_fkey TO mecha_chassis_game_id_fkey;
ALTER TABLE public.mecha_game_chassis
    RENAME CONSTRAINT mecha_game_chassis_class_check TO mecha_chassis_class_check;

-- ============================================================================
-- 14. Restore FK columns
-- ============================================================================

ALTER TABLE public.mecha_game_turn_sheet
    RENAME COLUMN mecha_game_squad_instance_id TO mecha_squad_instance_id;

ALTER TABLE public.mecha_game_mech_instance
    RENAME COLUMN mecha_game_squad_instance_id TO mecha_squad_instance_id;
ALTER TABLE public.mecha_game_mech_instance
    RENAME COLUMN mecha_game_sector_instance_id TO mecha_sector_instance_id;
ALTER TABLE public.mecha_game_mech_instance
    RENAME COLUMN mecha_game_chassis_id TO mecha_chassis_id;

ALTER TABLE public.mecha_game_squad_instance
    RENAME COLUMN mecha_game_squad_id TO mecha_squad_id;
ALTER TABLE public.mecha_game_squad_instance
    RENAME COLUMN mecha_game_computer_opponent_id TO mecha_computer_opponent_id;

ALTER TABLE public.mecha_game_squad_mech
    RENAME COLUMN mecha_game_squad_id TO mecha_squad_id;
ALTER TABLE public.mecha_game_squad_mech
    RENAME COLUMN mecha_game_chassis_id TO mecha_chassis_id;

ALTER TABLE public.mecha_game_sector_instance
    RENAME COLUMN mecha_game_sector_id TO mecha_sector_id;

ALTER TABLE public.mecha_game_sector_link
    RENAME COLUMN from_mecha_game_sector_id TO from_mecha_sector_id;
ALTER TABLE public.mecha_game_sector_link
    RENAME COLUMN to_mecha_game_sector_id TO to_mecha_sector_id;

-- ============================================================================
-- 15. Rename tables back
-- ============================================================================

ALTER TABLE public.mecha_game_sector          RENAME TO mecha_sector;
ALTER TABLE public.mecha_game_weapon          RENAME TO mecha_weapon;
ALTER TABLE public.mecha_game_chassis         RENAME TO mecha_chassis;
ALTER TABLE public.mecha_game_computer_opponent RENAME TO mecha_computer_opponent;
ALTER TABLE public.mecha_game_squad           RENAME TO mecha_squad;
ALTER TABLE public.mecha_game_sector_link     RENAME TO mecha_sector_link;
ALTER TABLE public.mecha_game_sector_instance RENAME TO mecha_sector_instance;
ALTER TABLE public.mecha_game_squad_instance  RENAME TO mecha_squad_instance;
ALTER TABLE public.mecha_game_squad_mech      RENAME TO mecha_squad_mech;
ALTER TABLE public.mecha_game_turn_sheet      RENAME TO mecha_turn_sheet;
ALTER TABLE public.mecha_game_mech_instance   RENAME TO mecha_mech_instance;

COMMIT;
