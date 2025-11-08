-- Add status column to account for approval workflow
ALTER TABLE public.account
    ADD COLUMN status VARCHAR(32) NOT NULL DEFAULT 'active'
        CHECK (status IN ('pending_approval', 'active', 'disabled'));

COMMENT ON COLUMN public.account.status IS 'Current approval status of the account (pending_approval, active, disabled).';

-- Backfill existing account rows to active status
UPDATE public.account
SET status = 'active'
WHERE status IS NULL;

-- Add status column to game_subscription for approval workflow
ALTER TABLE public.game_subscription
    ADD COLUMN status VARCHAR(32) NOT NULL DEFAULT 'active'
        CHECK (status IN ('pending_approval', 'active', 'revoked'));

COMMENT ON COLUMN public.game_subscription.status IS 'Approval status of the game subscription (pending_approval, active, revoked).';

-- Backfill existing subscription rows to active status
UPDATE public.game_subscription
SET status = 'active'
WHERE status IS NULL;

-- Allow game turn sheets to be created before a game instance is assigned
ALTER TABLE public.game_turn_sheet
    ALTER COLUMN game_instance_id DROP NOT NULL;
