BEGIN;

ALTER TABLE public.adventure_game_creature
    ALTER COLUMN max_health SET DEFAULT 50;

COMMIT;
