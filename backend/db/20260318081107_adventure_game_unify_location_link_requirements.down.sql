BEGIN;

DROP INDEX IF EXISTS idx_adventure_game_location_link_requirement_purpose;
DROP INDEX IF EXISTS idx_adventure_game_location_link_requirement_creature_id;

ALTER TABLE public.adventure_game_location_link_requirement
    DROP CONSTRAINT IF EXISTS adventure_game_location_link_requirement_condition_target_check,
    DROP CONSTRAINT IF EXISTS adventure_game_location_link_requirement_condition_check,
    DROP CONSTRAINT IF EXISTS adventure_game_location_link_requirement_purpose_check,
    DROP CONSTRAINT IF EXISTS adventure_game_location_link_requirement_one_target,
    DROP COLUMN IF EXISTS condition,
    DROP COLUMN IF EXISTS purpose,
    DROP COLUMN IF EXISTS adventure_game_creature_id;

ALTER TABLE public.adventure_game_location_link_requirement
    ALTER COLUMN adventure_game_item_id SET NOT NULL;

ALTER TABLE public.adventure_game_location_link
    DROP COLUMN IF EXISTS locked_description;

COMMIT;
