-- Revert closed_testing_join_game_key columns back to closed_testing_join_game_token
ALTER TABLE game_instance
    RENAME COLUMN closed_testing_join_game_key TO closed_testing_join_game_token;

ALTER TABLE game_instance
    RENAME COLUMN closed_testing_join_game_key_expires_at TO closed_testing_join_game_token_expires_at;

-- Rename the index back
DROP INDEX IF EXISTS idx_game_instance_closed_testing_join_game_key;
CREATE INDEX idx_game_instance_closed_testing_join_game_token ON game_instance(closed_testing_join_game_token)
    WHERE closed_testing_join_game_token IS NOT NULL;

-- Update comments
COMMENT ON COLUMN game_instance.closed_testing_join_game_token IS 'Unique token for closed testing join game links';
COMMENT ON COLUMN game_instance.closed_testing_join_game_token_expires_at IS 'Optional expiration timestamp for closed testing join game token';


