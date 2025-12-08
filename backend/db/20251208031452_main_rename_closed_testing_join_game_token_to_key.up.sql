-- Rename closed_testing_join_game_token columns to closed_testing_join_game_key
-- to match the code expectations
ALTER TABLE game_instance
    RENAME COLUMN closed_testing_join_game_token TO closed_testing_join_game_key;

ALTER TABLE game_instance
    RENAME COLUMN closed_testing_join_game_token_expires_at TO closed_testing_join_game_key_expires_at;

-- Rename the index to match the new column name
DROP INDEX IF EXISTS idx_game_instance_closed_testing_join_game_token;
CREATE INDEX idx_game_instance_closed_testing_join_game_key ON game_instance(closed_testing_join_game_key)
    WHERE closed_testing_join_game_key IS NOT NULL;

-- Update comments
COMMENT ON COLUMN game_instance.closed_testing_join_game_key IS 'Unique token for closed testing join game links';
COMMENT ON COLUMN game_instance.closed_testing_join_game_key_expires_at IS 'Optional expiration timestamp for closed testing join game token';



