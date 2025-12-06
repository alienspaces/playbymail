-- Rename join_game_key columns to closed_testing_join_game_key for clarity
ALTER TABLE game_instance
    RENAME COLUMN join_game_key TO closed_testing_join_game_key;

ALTER TABLE game_instance
    RENAME COLUMN join_game_key_expires_at TO closed_testing_join_game_key_expires_at;

-- Rename index to match new column name
DROP INDEX IF EXISTS idx_game_instance_join_game_key;
CREATE INDEX idx_game_instance_closed_testing_join_game_key ON game_instance(closed_testing_join_game_key)
    WHERE closed_testing_join_game_key IS NOT NULL;

-- Update comments
COMMENT ON COLUMN game_instance.closed_testing_join_game_key IS 'Unique secret key for closed testing join game links';
COMMENT ON COLUMN game_instance.closed_testing_join_game_key_expires_at IS 'Optional expiration timestamp for closed testing join game key';

