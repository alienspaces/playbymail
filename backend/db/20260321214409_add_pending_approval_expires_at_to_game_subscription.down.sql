ALTER TABLE public.game_subscription
    DROP COLUMN IF EXISTS pending_approval_expires_at;
