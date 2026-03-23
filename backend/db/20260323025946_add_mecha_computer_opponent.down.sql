-- Reverse: add_mecha_computer_opponent

BEGIN;

ALTER TABLE public.mecha_lance_instance
    ALTER COLUMN game_subscription_instance_id SET NOT NULL;

DROP INDEX IF EXISTS idx_mecha_lance_game_account_unique;
ALTER TABLE public.mecha_lance ADD CONSTRAINT mecha_lance_unique UNIQUE (game_id, account_id);

ALTER TABLE public.mecha_lance
    DROP CONSTRAINT IF EXISTS mecha_lance_owner_check,
    DROP CONSTRAINT IF EXISTS mecha_lance_computer_opponent_id_fkey,
    DROP COLUMN IF EXISTS mecha_computer_opponent_id,
    ALTER COLUMN account_id SET NOT NULL,
    ALTER COLUMN account_user_id SET NOT NULL;

DROP INDEX IF EXISTS idx_mecha_computer_opponent_game_id;
DROP TABLE IF EXISTS public.mecha_computer_opponent;

COMMIT;
