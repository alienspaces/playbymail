-- Remove game_instance_id from game_subscription table

DROP INDEX IF EXISTS public.idx_game_subscription_game_instance_id;

ALTER TABLE public.game_subscription
DROP CONSTRAINT IF EXISTS game_subscription_game_instance_id_fkey;

ALTER TABLE public.game_subscription
DROP COLUMN IF EXISTS game_instance_id;

