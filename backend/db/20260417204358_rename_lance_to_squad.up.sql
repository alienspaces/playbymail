-- Rename lance → squad across the mecha game type.
-- Renames three tables, their columns, constraints, indexes, and updates
-- stored string values in game_turn_sheet, game_image, and game_instance_parameter.

BEGIN;

-- ============================================================================
-- 1. Rename tables
-- ============================================================================

ALTER TABLE public.mecha_lance          RENAME TO mecha_squad;
ALTER TABLE public.mecha_lance_mech     RENAME TO mecha_squad_mech;
ALTER TABLE public.mecha_lance_instance RENAME TO mecha_squad_instance;

-- ============================================================================
-- 2. Rename columns
-- ============================================================================

-- mecha_squad: lance_type → squad_type
ALTER TABLE public.mecha_squad RENAME COLUMN lance_type TO squad_type;

-- mecha_squad_mech: mecha_lance_id → mecha_squad_id
ALTER TABLE public.mecha_squad_mech RENAME COLUMN mecha_lance_id TO mecha_squad_id;

-- mecha_squad_instance: mecha_lance_id → mecha_squad_id
ALTER TABLE public.mecha_squad_instance RENAME COLUMN mecha_lance_id TO mecha_squad_id;

-- mecha_mech_instance: mecha_lance_instance_id → mecha_squad_instance_id
ALTER TABLE public.mecha_mech_instance RENAME COLUMN mecha_lance_instance_id TO mecha_squad_instance_id;

-- mecha_turn_sheet: mecha_lance_instance_id → mecha_squad_instance_id
ALTER TABLE public.mecha_turn_sheet RENAME COLUMN mecha_lance_instance_id TO mecha_squad_instance_id;

-- ============================================================================
-- 3. Rename constraints on mecha_squad (formerly mecha_lance)
-- ============================================================================

ALTER TABLE public.mecha_squad
    RENAME CONSTRAINT mecha_lance_name_unique    TO mecha_squad_name_unique;
ALTER TABLE public.mecha_squad
    RENAME CONSTRAINT mecha_lance_game_id_fkey   TO mecha_squad_game_id_fkey;
ALTER TABLE public.mecha_squad
    RENAME CONSTRAINT mecha_lance_type_check     TO mecha_squad_type_check;

-- Partial unique index references the column by name in its WHERE clause text;
-- drop and recreate with new column name so the definition is clean.
DROP INDEX IF EXISTS public.idx_mecha_lance_starter_unique;
CREATE UNIQUE INDEX idx_mecha_squad_starter_unique
    ON public.mecha_squad (game_id)
    WHERE squad_type = 'starter' AND deleted_at IS NULL;

-- ============================================================================
-- 4. Rename constraints on mecha_squad_mech (formerly mecha_lance_mech)
-- ============================================================================

ALTER TABLE public.mecha_squad_mech
    RENAME CONSTRAINT mecha_lance_mech_callsign_unique TO mecha_squad_mech_callsign_unique;
ALTER TABLE public.mecha_squad_mech
    RENAME CONSTRAINT mecha_lance_mech_game_id_fkey    TO mecha_squad_mech_game_id_fkey;
ALTER TABLE public.mecha_squad_mech
    RENAME CONSTRAINT mecha_lance_mech_lance_id_fkey   TO mecha_squad_mech_squad_id_fkey;
ALTER TABLE public.mecha_squad_mech
    RENAME CONSTRAINT mecha_lance_mech_chassis_id_fkey TO mecha_squad_mech_chassis_id_fkey;

-- ============================================================================
-- 5. Rename constraints on mecha_squad_instance (formerly mecha_lance_instance)
-- ============================================================================

ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_lance_instance_game_id_fkey
                   TO mecha_squad_instance_game_id_fkey;
ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_lance_instance_game_instance_id_fkey
                   TO mecha_squad_instance_game_instance_id_fkey;
ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_lance_instance_lance_id_fkey
                   TO mecha_squad_instance_squad_id_fkey;
ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_lance_instance_subscription_instance_id_fkey
                   TO mecha_squad_instance_subscription_instance_id_fkey;
ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_lance_instance_computer_opponent_id_fkey
                   TO mecha_squad_instance_computer_opponent_id_fkey;

-- ============================================================================
-- 6. Rename constraints on mecha_mech_instance
-- ============================================================================

ALTER TABLE public.mecha_mech_instance
    RENAME CONSTRAINT mecha_mech_instance_lance_instance_id_fkey
                   TO mecha_mech_instance_squad_instance_id_fkey;

-- ============================================================================
-- 7. Rename constraints on mecha_turn_sheet
-- ============================================================================

-- Unique constraint references the renamed column; drop and recreate for a clean definition.
ALTER TABLE public.mecha_turn_sheet
    DROP CONSTRAINT mecha_turn_sheet_unique;
ALTER TABLE public.mecha_turn_sheet
    ADD CONSTRAINT mecha_turn_sheet_unique UNIQUE (mecha_squad_instance_id, game_turn_sheet_id);

ALTER TABLE public.mecha_turn_sheet
    RENAME CONSTRAINT mecha_turn_sheet_lance_instance_id_fkey
                   TO mecha_turn_sheet_squad_instance_id_fkey;

-- ============================================================================
-- 8. Rename indexes
-- ============================================================================

ALTER INDEX public.idx_mecha_lance_game_id
    RENAME TO idx_mecha_squad_game_id;

ALTER INDEX public.idx_mecha_lance_mech_lance_id
    RENAME TO idx_mecha_squad_mech_squad_id;
ALTER INDEX public.idx_mecha_lance_mech_chassis_id
    RENAME TO idx_mecha_squad_mech_chassis_id;

ALTER INDEX public.idx_mecha_lance_instance_game_instance
    RENAME TO idx_mecha_squad_instance_game_instance;
ALTER INDEX public.idx_mecha_lance_instance_lance_id
    RENAME TO idx_mecha_squad_instance_squad_id;
ALTER INDEX public.idx_mecha_lance_instance_computer_opponent
    RENAME TO idx_mecha_squad_instance_computer_opponent;

ALTER INDEX public.idx_mecha_mech_instance_lance_instance
    RENAME TO idx_mecha_mech_instance_squad_instance;

ALTER INDEX public.idx_mecha_turn_sheet_lance_instance
    RENAME TO idx_mecha_turn_sheet_squad_instance;

-- ============================================================================
-- 9. Update stored string values
-- ============================================================================

-- Turn sheet type: mecha_lance_management → mecha_squad_management
UPDATE public.game_turn_sheet
    SET sheet_type = 'mecha_squad_management'
    WHERE sheet_type = 'mecha_lance_management';

UPDATE public.game_image
    SET turn_sheet_type = 'mecha_squad_management'
    WHERE turn_sheet_type = 'mecha_lance_management';

-- Game instance parameter key: lance_size → squad_size
UPDATE public.game_instance_parameter
    SET parameter_key = 'squad_size'
    WHERE parameter_key = 'lance_size';

COMMIT;
