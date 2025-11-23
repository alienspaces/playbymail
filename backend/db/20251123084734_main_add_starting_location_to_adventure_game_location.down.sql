-- Remove is_starting_location field from adventure_game_location table
DROP INDEX IF EXISTS public.idx_adventure_game_location_is_starting_location;

ALTER TABLE public.adventure_game_location
    DROP COLUMN IF EXISTS is_starting_location;

