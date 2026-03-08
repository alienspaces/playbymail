CREATE OR REPLACE VIEW public.account_game_view AS
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
    -- is_designer: account has an active designer game_subscription for this game
    EXISTS (
        SELECT 1 FROM public.game_subscription gs
        WHERE gs.game_id = g.id
          AND gs.account_id = a.id
          AND gs.subscription_type = 'designer'
          AND gs.status = 'active'
          AND gs.deleted_at IS NULL
    ) AS is_designer,
    -- is_manager: account has an active manager game_subscription for this game
    EXISTS (
        SELECT 1 FROM public.game_subscription gs
        WHERE gs.game_id = g.id
          AND gs.account_id = a.id
          AND gs.subscription_type = 'manager'
          AND gs.status = 'active'
          AND gs.deleted_at IS NULL
    ) AS is_manager,
    -- can_manage: game is published, account has a manager account_subscription,
    -- and account is not already managing this game
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
