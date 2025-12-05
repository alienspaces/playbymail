-- Add delivery method flags, required player count, and closed testing fields to game_instance
ALTER TABLE game_instance
    ADD COLUMN delivery_physical_post BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN delivery_physical_local BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN delivery_email BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN required_player_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN is_closed_testing BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN join_game_key VARCHAR(255),
    ADD COLUMN join_game_key_expires_at TIMESTAMPTZ;

-- Add constraints
ALTER TABLE game_instance
    ADD CONSTRAINT game_instance_delivery_methods_check
        CHECK (delivery_physical_post = true OR delivery_physical_local = true OR delivery_email = true),
    ADD CONSTRAINT game_instance_required_player_count_check
        CHECK (required_player_count >= 0);

-- Create index on join_game_key for fast lookups
CREATE INDEX idx_game_instance_join_game_key ON game_instance(join_game_key)
    WHERE join_game_key IS NOT NULL;

-- Add comments
COMMENT ON COLUMN game_instance.delivery_physical_post IS 'Enable physical post delivery (traditional mail-based)';
COMMENT ON COLUMN game_instance.delivery_physical_local IS 'Enable physical local delivery (convention/classroom - game master prints locally, players fill at table, manual scanning/submission)';
COMMENT ON COLUMN game_instance.delivery_email IS 'Enable email delivery (web-based turn sheet viewer via email links)';
COMMENT ON COLUMN game_instance.required_player_count IS 'Minimum number of players required before game can start';
COMMENT ON COLUMN game_instance.is_closed_testing IS 'Whether this game instance is in closed testing mode';
COMMENT ON COLUMN game_instance.join_game_key IS 'Unique secret key for closed testing join game links';
COMMENT ON COLUMN game_instance.join_game_key_expires_at IS 'Optional expiration timestamp for join game key';

