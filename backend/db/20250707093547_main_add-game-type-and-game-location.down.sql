DROP TABLE IF EXISTS public.adventure_game_location;

ALTER TABLE public.game DROP CONSTRAINT IF EXISTS game_type_check;
ALTER TABLE public.game DROP COLUMN IF EXISTS game_type;
