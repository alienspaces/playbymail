BEGIN;

ALTER TABLE public.adventure_game_item
    DROP COLUMN IF EXISTS heal_amount,
    DROP COLUMN IF EXISTS defense,
    DROP COLUMN IF EXISTS damage;

ALTER TABLE public.adventure_game_creature
    DROP CONSTRAINT IF EXISTS adventure_game_creature_disposition_check,
    DROP COLUMN IF EXISTS disposition,
    DROP COLUMN IF EXISTS defense,
    DROP COLUMN IF EXISTS attack_damage;

COMMIT;
