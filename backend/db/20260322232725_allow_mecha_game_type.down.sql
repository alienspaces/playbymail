-- Revert game_type constraint to adventure only

BEGIN;

ALTER TABLE public.game DROP CONSTRAINT IF EXISTS game_type_check;
ALTER TABLE public.game ADD CONSTRAINT game_type_check
    CHECK (game_type = 'adventure');

COMMIT;
