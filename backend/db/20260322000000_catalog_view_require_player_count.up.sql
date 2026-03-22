BEGIN;

DROP VIEW IF EXISTS public.catalog_game_instance_view;
CREATE VIEW public.catalog_game_instance_view AS
SELECT DISTINCT ON (
    gs.account_id,
    g.id,
    gi.turn_duration_hours,
    gi.required_player_count,
    gi.delivery_email,
    gi.delivery_physical_post,
    gi.delivery_physical_local
)
    gi.id AS id,
    gi.id AS game_instance_id,
    g.id AS game_id,
    g.name AS game_name,
    g.game_type,
    g.description AS game_description,
    gi.turn_duration_hours,
    gs.id AS game_subscription_id,
    gs.account_id,
    a.name AS account_name,
    gi.required_player_count,
    COALESCE(pc.player_count, 0) AS player_count,
    gi.required_player_count - COALESCE(pc.player_count, 0) AS remaining_capacity,
    gi.delivery_email,
    gi.delivery_physical_post,
    gi.delivery_physical_local,
    gi.is_closed_testing,
    gi.created_at,
    gi.updated_at,
    gi.deleted_at
FROM public.game_instance gi
JOIN public.game_subscription_instance gsi
    ON gsi.game_instance_id = gi.id AND gsi.deleted_at IS NULL
JOIN public.game_subscription gs
    ON gs.id = gsi.game_subscription_id
    AND gs.subscription_type = 'manager'
    AND gs.status = 'active'
    AND gs.deleted_at IS NULL
JOIN public.account a
    ON a.id = gs.account_id AND a.deleted_at IS NULL
JOIN public.game g
    ON g.id = gi.game_id AND g.deleted_at IS NULL
LEFT JOIN (
    SELECT gsi2.game_instance_id, COUNT(*) AS player_count
    FROM public.game_subscription_instance gsi2
    JOIN public.game_subscription gs2
        ON gs2.id = gsi2.game_subscription_id
        AND gs2.subscription_type = 'player'
        AND gs2.deleted_at IS NULL
        AND (gs2.status = 'active'
             OR (gs2.status = 'pending_approval'
                 AND (gs2.pending_approval_expires_at IS NULL
                      OR gs2.pending_approval_expires_at > now())))
    WHERE gsi2.deleted_at IS NULL
    GROUP BY gsi2.game_instance_id
) pc ON pc.game_instance_id = gi.id
WHERE gi.status = 'created'
  AND gi.deleted_at IS NULL
  AND gi.is_closed_testing = false
  AND gi.required_player_count >= 1
  AND COALESCE(pc.player_count, 0) < gi.required_player_count
ORDER BY
    gs.account_id,
    g.id,
    gi.turn_duration_hours,
    gi.required_player_count,
    gi.delivery_email,
    gi.delivery_physical_post,
    gi.delivery_physical_local,
    gi.created_at ASC;

COMMIT;
