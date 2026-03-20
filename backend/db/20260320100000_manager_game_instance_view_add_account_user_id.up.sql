-- Add account_user_id to manager_game_instance_view so ownership checks
-- can use account_user_id instead of account_id.
-- DROP + CREATE required because CREATE OR REPLACE cannot add columns.

DROP VIEW IF EXISTS public.manager_game_instance_view;

CREATE VIEW public.manager_game_instance_view AS
SELECT
    gs.id AS id,
    gs.account_id,
    gs.account_user_id,
    g.id AS game_id,
    g.name AS game_name,
    g.game_type,
    g.description AS game_description,
    gs.id AS game_subscription_id,
    gi.id AS game_instance_id,
    gi.status AS instance_status,
    gi.current_turn,
    gi.required_player_count,
    gi.delivery_email,
    gi.delivery_physical_post,
    gi.delivery_physical_local,
    gi.is_closed_testing,
    gi.started_at,
    gi.next_turn_due_at,
    gi.created_at AS instance_created_at,
    g.created_at,
    g.updated_at,
    g.deleted_at
FROM public.game_subscription gs
JOIN public.game g ON g.id = gs.game_id AND g.deleted_at IS NULL
LEFT JOIN public.game_subscription_instance gsi
    ON gsi.game_subscription_id = gs.id AND gsi.deleted_at IS NULL
LEFT JOIN public.game_instance gi
    ON gi.id = gsi.game_instance_id AND gi.deleted_at IS NULL
WHERE gs.subscription_type = 'manager'
  AND gs.status = 'active'
  AND gs.deleted_at IS NULL;
