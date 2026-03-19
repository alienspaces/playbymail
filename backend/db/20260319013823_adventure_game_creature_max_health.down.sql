BEGIN;

ALTER TABLE public.adventure_game_creature
    DROP COLUMN IF EXISTS max_health;

COMMIT;
