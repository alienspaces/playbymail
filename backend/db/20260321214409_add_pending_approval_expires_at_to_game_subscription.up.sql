ALTER TABLE public.game_subscription
    ADD COLUMN pending_approval_expires_at TIMESTAMPTZ;
