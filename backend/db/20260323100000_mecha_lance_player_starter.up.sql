-- Add is_player_starter to mecha_lance for designer-configurable player starter templates

BEGIN;

ALTER TABLE public.mecha_lance
    ADD COLUMN is_player_starter BOOLEAN NOT NULL DEFAULT FALSE;

-- Update ownership constraint to allow a third mode: starter template (no owner)
ALTER TABLE public.mecha_lance
    DROP CONSTRAINT mecha_lance_owner_check,
    ADD CONSTRAINT mecha_lance_owner_check CHECK (
        (is_player_starter = true  AND mecha_computer_opponent_id IS NULL AND account_id IS NULL AND account_user_id IS NULL)
        OR
        (is_player_starter = false AND mecha_computer_opponent_id IS NOT NULL AND account_id IS NULL AND account_user_id IS NULL)
        OR
        (is_player_starter = false AND mecha_computer_opponent_id IS NULL AND account_id IS NOT NULL AND account_user_id IS NOT NULL)
    );

-- At most one starter lance per game
CREATE UNIQUE INDEX idx_mecha_lance_player_starter_unique
    ON public.mecha_lance (game_id)
    WHERE is_player_starter = true AND deleted_at IS NULL;

COMMIT;
