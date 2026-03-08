-- Revert views to previous definitions (before standard column fix)

DROP VIEW IF EXISTS public.account_game_view;
CREATE VIEW public.account_game_view AS
SELECT
    a.id AS account_id,
    a.name AS account_name,
    g.id AS game_id,
    g.name AS game_name,
    g.game_type,
    g.description,
    g.turn_duration_hours,
    g.status AS game_status,
    g.created_at,
    g.updated_at,
    g.deleted_at,
    EXISTS (
        SELECT 1 FROM public.game_subscription gs
        WHERE gs.game_id = g.id
          AND gs.account_id = a.id
          AND gs.subscription_type = 'designer'
          AND gs.status = 'active'
          AND gs.deleted_at IS NULL
    ) AS is_designer,
    EXISTS (
        SELECT 1 FROM public.game_subscription gs
        WHERE gs.game_id = g.id
          AND gs.account_id = a.id
          AND gs.subscription_type = 'manager'
          AND gs.status = 'active'
          AND gs.deleted_at IS NULL
    ) AS is_manager,
    (
        g.status = 'published'
        AND EXISTS (
            SELECT 1 FROM public.account_subscription acs
            WHERE acs.account_id = a.id
              AND acs.subscription_type IN ('basic_manager', 'professional_manager')
              AND acs.status = 'active'
              AND acs.deleted_at IS NULL
        )
        AND NOT EXISTS (
            SELECT 1 FROM public.game_subscription gs
            WHERE gs.game_id = g.id
              AND gs.account_id = a.id
              AND gs.subscription_type = 'manager'
              AND gs.status = 'active'
              AND gs.deleted_at IS NULL
        )
    ) AS can_manage
FROM public.account a
CROSS JOIN public.game g
WHERE a.deleted_at IS NULL
  AND g.deleted_at IS NULL
  AND (
      g.status = 'published'
      OR EXISTS (
          SELECT 1 FROM public.game_subscription gs
          WHERE gs.game_id = g.id
            AND gs.account_id = a.id
            AND gs.subscription_type = 'designer'
            AND gs.status = 'active'
            AND gs.deleted_at IS NULL
      )
  );

DROP VIEW IF EXISTS public.manager_game_instance_view;
CREATE VIEW public.manager_game_instance_view AS
SELECT
    gs.account_id,
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

DROP VIEW IF EXISTS public.catalog_game_instance_view;
CREATE VIEW public.catalog_game_instance_view AS
SELECT
    gi.id AS game_instance_id,
    g.id AS game_id,
    g.name AS game_name,
    g.game_type,
    g.description AS game_description,
    g.turn_duration_hours,
    gs.id AS game_subscription_id,
    gi.required_player_count,
    gi.delivery_email,
    gi.delivery_physical_post,
    gi.delivery_physical_local,
    gi.is_closed_testing,
    gi.created_at
FROM public.game_instance gi
JOIN public.game_subscription_instance gsi
    ON gsi.game_instance_id = gi.id AND gsi.deleted_at IS NULL
JOIN public.game_subscription gs
    ON gs.id = gsi.game_subscription_id
    AND gs.subscription_type = 'manager'
    AND gs.status = 'active'
    AND gs.deleted_at IS NULL
JOIN public.game g
    ON g.id = gi.game_id AND g.deleted_at IS NULL
WHERE gi.status = 'created'
  AND gi.deleted_at IS NULL
  AND gi.is_closed_testing = false;
