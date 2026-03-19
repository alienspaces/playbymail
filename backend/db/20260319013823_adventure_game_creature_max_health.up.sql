BEGIN;

ALTER TABLE public.adventure_game_creature
    ADD COLUMN max_health INT NOT NULL DEFAULT 50;

COMMIT;
