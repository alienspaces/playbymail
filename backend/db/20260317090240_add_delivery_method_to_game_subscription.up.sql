BEGIN;

ALTER TABLE public.game_subscription
    ADD COLUMN delivery_method VARCHAR(32);

ALTER TABLE public.game_subscription
    ADD CONSTRAINT game_subscription_delivery_method_check
    CHECK (delivery_method IN ('email', 'local', 'post'));

-- Backfill existing player subscriptions to 'email'.
UPDATE public.game_subscription
SET delivery_method = 'email'
WHERE subscription_type = 'player' AND delivery_method IS NULL;

ALTER TABLE public.game_subscription
    ADD CONSTRAINT game_subscription_player_requires_delivery_method
    CHECK (subscription_type != 'player' OR delivery_method IS NOT NULL);

DROP VIEW IF EXISTS public.game_subscription_view;

CREATE VIEW public.game_subscription_view AS
SELECT
    gs.id,
    gs.game_id,
    gs.account_id,
    gs.account_user_id,
    gs.account_contact_id,
    gs.subscription_type,
    gs.status,
    gs.instance_limit,
    gs.delivery_method,
    gs.created_at,
    gs.updated_at,
    gs.deleted_at,
    COALESCE(
        array_agg(gsi.game_instance_id) FILTER (WHERE gsi.game_instance_id IS NOT NULL AND gsi.deleted_at IS NULL),
        ARRAY[]::UUID[]
    ) AS game_instance_ids
FROM public.game_subscription gs
LEFT JOIN public.game_subscription_instance gsi ON gs.id = gsi.game_subscription_id
WHERE gs.deleted_at IS NULL
GROUP BY
    gs.id,
    gs.game_id,
    gs.account_id,
    gs.account_user_id,
    gs.account_contact_id,
    gs.subscription_type,
    gs.status,
    gs.instance_limit,
    gs.delivery_method,
    gs.created_at,
    gs.updated_at,
    gs.deleted_at;

COMMENT ON COLUMN public.game_subscription.delivery_method IS 'Player preferred turn sheet delivery method. Required for player subscriptions.';

COMMIT;
