BEGIN;

ALTER TABLE public.adventure_game_creature
    ALTER COLUMN max_health DROP DEFAULT;

COMMIT;
