-- Add turn sheet token fields to game_subscription table
ALTER TABLE game_subscription
    ADD COLUMN turn_sheet_token VARCHAR(255),
    ADD COLUMN turn_sheet_token_expires_at TIMESTAMPTZ;

-- Create index on turn_sheet_token for fast lookups
CREATE INDEX idx_game_subscription_turn_sheet_token ON game_subscription(turn_sheet_token)
    WHERE turn_sheet_token IS NOT NULL;

-- Add comments
COMMENT ON COLUMN game_subscription.turn_sheet_token IS 'Unique token for accessing turn sheets via web viewer. Refreshed with every new turn and automatically expired once player submits their latest turn sheets';
COMMENT ON COLUMN game_subscription.turn_sheet_token_expires_at IS 'Expiration timestamp for turn sheet token (3 days from generation). Key is also automatically expired when player submits their latest turn sheets for the related game';

