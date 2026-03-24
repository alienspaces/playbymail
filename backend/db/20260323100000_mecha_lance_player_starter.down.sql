-- Reverse: mecha_lance_player_starter

BEGIN;

DROP INDEX IF EXISTS idx_mecha_lance_player_starter_unique;

ALTER TABLE public.mecha_lance
    DROP CONSTRAINT IF EXISTS mecha_lance_owner_check,
    ADD CONSTRAINT mecha_lance_owner_check CHECK (
        (mecha_computer_opponent_id IS NOT NULL AND account_id IS NULL AND account_user_id IS NULL)
        OR
        (mecha_computer_opponent_id IS NULL AND account_id IS NOT NULL AND account_user_id IS NOT NULL)
    );

ALTER TABLE public.mecha_lance
    DROP COLUMN IF EXISTS is_player_starter;

COMMIT;
