-- Allow mecha as a valid game type

BEGIN;

ALTER TABLE public.game DROP CONSTRAINT IF EXISTS game_type_check;
ALTER TABLE public.game ADD CONSTRAINT game_type_check
    CHECK (game_type IN ('adventure', 'mecha'));

COMMIT;
