-- Remove delivery method flags, required player count, and closed testing fields from game_instance
DROP INDEX IF EXISTS idx_game_instance_closed_testing_join_game_token;

ALTER TABLE game_instance
    DROP CONSTRAINT IF EXISTS game_instance_required_player_count_check,
    DROP CONSTRAINT IF EXISTS game_instance_delivery_methods_check,
    DROP COLUMN IF EXISTS closed_testing_join_game_token_expires_at,
    DROP COLUMN IF EXISTS closed_testing_join_game_token,
    DROP COLUMN IF EXISTS is_closed_testing,
    DROP COLUMN IF EXISTS required_player_count,
    DROP COLUMN IF EXISTS delivery_email,
    DROP COLUMN IF EXISTS delivery_physical_local,
    DROP COLUMN IF EXISTS delivery_physical_post;

