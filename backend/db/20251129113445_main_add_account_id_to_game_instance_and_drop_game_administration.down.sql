-- Recreate game_administration table
CREATE TABLE IF NOT EXISTS public.game_administration (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.game(id),
    account_id UUID NOT NULL REFERENCES public.account(id),
    granted_by_account_id UUID NOT NULL REFERENCES public.account(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE (game_id, account_id)
);

COMMENT ON TABLE public.game_administration IS 'Tracks which accounts have been granted admin rights for all instances of a specific game.';
COMMENT ON COLUMN public.game_administration.id IS 'Unique identifier for the admin grant.';
COMMENT ON COLUMN public.game_administration.game_id IS 'The game for which admin rights are granted.';
COMMENT ON COLUMN public.game_administration.account_id IS 'The account being granted admin rights.';
COMMENT ON COLUMN public.game_administration.granted_by_account_id IS 'The account who granted the admin rights.';
COMMENT ON COLUMN public.game_administration.created_at IS 'When the admin rights were granted.';
COMMENT ON COLUMN public.game_administration.updated_at IS 'When the admin rights were last updated.';
COMMENT ON COLUMN public.game_administration.deleted_at IS 'When the admin rights were logically deleted.';

-- Remove game_subscription_id from game_instance
DROP INDEX IF EXISTS idx_game_instance_game_subscription_id;
ALTER TABLE public.game_instance DROP COLUMN IF EXISTS game_subscription_id;

