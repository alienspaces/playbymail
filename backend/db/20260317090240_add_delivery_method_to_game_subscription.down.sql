BEGIN;

ALTER TABLE public.game_subscription
    DROP CONSTRAINT IF EXISTS game_subscription_player_requires_delivery_method;

ALTER TABLE public.game_subscription
    DROP CONSTRAINT IF EXISTS game_subscription_delivery_method_check;

ALTER TABLE public.game_subscription
    DROP COLUMN IF EXISTS delivery_method;

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
    gs.created_at,
    gs.updated_at,
    gs.deleted_at;

COMMIT;
