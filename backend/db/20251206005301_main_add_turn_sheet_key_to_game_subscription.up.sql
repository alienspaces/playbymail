-- Add turn sheet key fields to game_subscription table
ALTER TABLE game_subscription
    ADD COLUMN turn_sheet_key VARCHAR(255),
    ADD COLUMN turn_sheet_key_expires_at TIMESTAMPTZ;

-- Create index on turn_sheet_key for fast lookups
CREATE INDEX idx_game_subscription_turn_sheet_key ON game_subscription(turn_sheet_key)
    WHERE turn_sheet_key IS NOT NULL;

-- Add comments
COMMENT ON COLUMN game_subscription.turn_sheet_key IS 'Unique secret key for accessing turn sheets via web viewer. Refreshed with every new turn and automatically expired once player submits their latest turn sheets';
COMMENT ON COLUMN game_subscription.turn_sheet_key_expires_at IS 'Expiration timestamp for turn sheet key (3 days from generation). Key is also automatically expired when player submits their latest turn sheets for the related game';

