-- Reverse: rename squad → lance across the mecha game type.

BEGIN;

-- ============================================================================
-- 1. Revert stored string values
-- ============================================================================

UPDATE public.game_instance_parameter
    SET parameter_key = 'lance_size'
    WHERE parameter_key = 'squad_size';

UPDATE public.game_image
    SET turn_sheet_type = 'mecha_lance_management'
    WHERE turn_sheet_type = 'mecha_squad_management';

UPDATE public.game_turn_sheet
    SET sheet_type = 'mecha_lance_management'
    WHERE sheet_type = 'mecha_squad_management';

-- ============================================================================
-- 2. Revert indexes
-- ============================================================================

ALTER INDEX public.idx_mecha_turn_sheet_squad_instance
    RENAME TO idx_mecha_turn_sheet_lance_instance;

ALTER INDEX public.idx_mecha_mech_instance_squad_instance
    RENAME TO idx_mecha_mech_instance_lance_instance;

ALTER INDEX public.idx_mecha_squad_instance_computer_opponent
    RENAME TO idx_mecha_lance_instance_computer_opponent;
ALTER INDEX public.idx_mecha_squad_instance_squad_id
    RENAME TO idx_mecha_lance_instance_lance_id;
ALTER INDEX public.idx_mecha_squad_instance_game_instance
    RENAME TO idx_mecha_lance_instance_game_instance;

ALTER INDEX public.idx_mecha_squad_mech_chassis_id
    RENAME TO idx_mecha_lance_mech_chassis_id;
ALTER INDEX public.idx_mecha_squad_mech_squad_id
    RENAME TO idx_mecha_lance_mech_lance_id;

ALTER INDEX public.idx_mecha_squad_game_id
    RENAME TO idx_mecha_lance_game_id;

-- ============================================================================
-- 3. Revert constraints on mecha_turn_sheet
-- ============================================================================

ALTER TABLE public.mecha_turn_sheet
    RENAME CONSTRAINT mecha_turn_sheet_squad_instance_id_fkey
                   TO mecha_turn_sheet_lance_instance_id_fkey;

ALTER TABLE public.mecha_turn_sheet
    DROP CONSTRAINT mecha_turn_sheet_unique;
ALTER TABLE public.mecha_turn_sheet
    ADD CONSTRAINT mecha_turn_sheet_unique UNIQUE (mecha_lance_instance_id, game_turn_sheet_id);

-- ============================================================================
-- 4. Revert constraints on mecha_mech_instance
-- ============================================================================

ALTER TABLE public.mecha_mech_instance
    RENAME CONSTRAINT mecha_mech_instance_squad_instance_id_fkey
                   TO mecha_mech_instance_lance_instance_id_fkey;

-- ============================================================================
-- 5. Revert constraints on mecha_squad_instance
-- ============================================================================

ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_computer_opponent_id_fkey
                   TO mecha_lance_instance_computer_opponent_id_fkey;
ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_subscription_instance_id_fkey
                   TO mecha_lance_instance_subscription_instance_id_fkey;
ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_squad_id_fkey
                   TO mecha_lance_instance_lance_id_fkey;
ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_game_instance_id_fkey
                   TO mecha_lance_instance_game_instance_id_fkey;
ALTER TABLE public.mecha_squad_instance
    RENAME CONSTRAINT mecha_squad_instance_game_id_fkey
                   TO mecha_lance_instance_game_id_fkey;

-- ============================================================================
-- 6. Revert constraints on mecha_squad_mech
-- ============================================================================

ALTER TABLE public.mecha_squad_mech
    RENAME CONSTRAINT mecha_squad_mech_chassis_id_fkey TO mecha_lance_mech_chassis_id_fkey;
ALTER TABLE public.mecha_squad_mech
    RENAME CONSTRAINT mecha_squad_mech_squad_id_fkey   TO mecha_lance_mech_lance_id_fkey;
ALTER TABLE public.mecha_squad_mech
    RENAME CONSTRAINT mecha_squad_mech_game_id_fkey    TO mecha_lance_mech_game_id_fkey;
ALTER TABLE public.mecha_squad_mech
    RENAME CONSTRAINT mecha_squad_mech_callsign_unique TO mecha_lance_mech_callsign_unique;

-- ============================================================================
-- 7. Revert constraints on mecha_squad
-- ============================================================================

DROP INDEX IF EXISTS public.idx_mecha_squad_starter_unique;
CREATE UNIQUE INDEX idx_mecha_lance_starter_unique
    ON public.mecha_squad (game_id)
    WHERE lance_type = 'starter' AND deleted_at IS NULL;

ALTER TABLE public.mecha_squad
    RENAME CONSTRAINT mecha_squad_type_check     TO mecha_lance_type_check;
ALTER TABLE public.mecha_squad
    RENAME CONSTRAINT mecha_squad_game_id_fkey   TO mecha_lance_game_id_fkey;
ALTER TABLE public.mecha_squad
    RENAME CONSTRAINT mecha_squad_name_unique    TO mecha_lance_name_unique;

-- ============================================================================
-- 8. Revert column renames
-- ============================================================================

ALTER TABLE public.mecha_turn_sheet     RENAME COLUMN mecha_squad_instance_id TO mecha_lance_instance_id;
ALTER TABLE public.mecha_mech_instance  RENAME COLUMN mecha_squad_instance_id TO mecha_lance_instance_id;
ALTER TABLE public.mecha_squad_instance RENAME COLUMN mecha_squad_id          TO mecha_lance_id;
ALTER TABLE public.mecha_squad_mech     RENAME COLUMN mecha_squad_id          TO mecha_lance_id;
ALTER TABLE public.mecha_squad          RENAME COLUMN squad_type               TO lance_type;

-- ============================================================================
-- 9. Revert table renames
-- ============================================================================

ALTER TABLE public.mecha_squad_instance RENAME TO mecha_lance_instance;
ALTER TABLE public.mecha_squad_mech     RENAME TO mecha_lance_mech;
ALTER TABLE public.mecha_squad          RENAME TO mecha_lance;

COMMIT;
