-- Revert join_game_key column names
ALTER TABLE game_instance
    RENAME COLUMN closed_testing_join_game_key TO join_game_key;

ALTER TABLE game_instance
    RENAME COLUMN closed_testing_join_game_key_expires_at TO join_game_key_expires_at;

-- Revert index name
DROP INDEX IF EXISTS idx_game_instance_closed_testing_join_game_key;
CREATE INDEX idx_game_instance_join_game_key ON game_instance(join_game_key)
    WHERE join_game_key IS NOT NULL;

