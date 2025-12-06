-- Remove turn sheet key fields from game_subscription table
DROP INDEX IF EXISTS idx_game_subscription_turn_sheet_key;

ALTER TABLE game_subscription
    DROP COLUMN IF EXISTS turn_sheet_key_expires_at,
    DROP COLUMN IF EXISTS turn_sheet_key;

