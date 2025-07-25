-- Create table for game_subscription (Player, Manager, Collaborator roles)
CREATE TABLE IF NOT EXISTS public.game_subscription (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.game(id),
    account_id UUID NOT NULL REFERENCES public.account(id),
    subscription_type VARCHAR(32) NOT NULL CHECK (subscription_type IN ('Player', 'Manager', 'Collaborator')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE (game_id, account_id, subscription_type)
);

COMMENT ON TABLE public.game_subscription IS 'Tracks account subscriptions to a game, including Player, Manager, and Collaborator roles.';
COMMENT ON COLUMN public.game_subscription.id IS 'Unique identifier for the subscription.';
COMMENT ON COLUMN public.game_subscription.game_id IS 'The game being subscribed to.';
COMMENT ON COLUMN public.game_subscription.account_id IS 'The subscribing account.';
COMMENT ON COLUMN public.game_subscription.subscription_type IS 'Role: Player, Manager, or Collaborator.';
COMMENT ON COLUMN public.game_subscription.created_at IS 'When the subscription was created.';
COMMENT ON COLUMN public.game_subscription.updated_at IS 'When the subscription was last updated.';
COMMENT ON COLUMN public.game_subscription.deleted_at IS 'When the subscription was logically deleted.';

-- Create table for game_administration (delegated admin rights for all instances of a game)
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