-- Remove turn sheet token fields from game_subscription table
DROP INDEX IF EXISTS idx_game_subscription_turn_sheet_token;

ALTER TABLE game_subscription
    DROP COLUMN IF EXISTS turn_sheet_token_expires_at,
    DROP COLUMN IF EXISTS turn_sheet_token;

