-- Adventure Game Location Object State Migration - Rollback
-- Reverts FK state columns back to freeform text strings.

BEGIN;

-- Restore adventure_game_location_object_instance.current_state
ALTER TABLE public.adventure_game_location_object_instance
    DROP CONSTRAINT IF EXISTS adventure_game_location_object_instance_current_state_id_fkey;
ALTER TABLE public.adventure_game_location_object_instance
    DROP COLUMN IF EXISTS current_adventure_game_location_object_state_id;
ALTER TABLE public.adventure_game_location_object_instance
    ADD COLUMN current_state VARCHAR(50) NOT NULL DEFAULT 'intact';

-- Restore adventure_game_location_object_effect.result_state
ALTER TABLE public.adventure_game_location_object_effect
    DROP CONSTRAINT IF EXISTS adventure_game_location_object_effect_result_state_id_fkey;
ALTER TABLE public.adventure_game_location_object_effect
    DROP COLUMN IF EXISTS result_adventure_game_location_object_state_id;
ALTER TABLE public.adventure_game_location_object_effect
    ADD COLUMN result_state VARCHAR(50);

-- Restore adventure_game_location_object_effect.required_state
ALTER TABLE public.adventure_game_location_object_effect
    DROP CONSTRAINT IF EXISTS adventure_game_location_object_effect_required_state_id_fkey;
ALTER TABLE public.adventure_game_location_object_effect
    DROP COLUMN IF EXISTS required_adventure_game_location_object_state_id;
ALTER TABLE public.adventure_game_location_object_effect
    ADD COLUMN required_state VARCHAR(50);

-- Restore adventure_game_location_object.initial_state
ALTER TABLE public.adventure_game_location_object
    DROP CONSTRAINT IF EXISTS adventure_game_location_object_initial_state_id_fkey;
ALTER TABLE public.adventure_game_location_object
    DROP COLUMN IF EXISTS initial_adventure_game_location_object_state_id;
ALTER TABLE public.adventure_game_location_object
    ADD COLUMN initial_state VARCHAR(50) NOT NULL DEFAULT 'intact';

-- Drop state table (must come after FK columns are removed)
DROP TABLE IF EXISTS public.adventure_game_location_object_state;

COMMIT;
