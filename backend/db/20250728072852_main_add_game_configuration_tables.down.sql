-- Drop game configuration tables

-- Drop indexes
DROP INDEX IF EXISTS idx_game_instance_parameter_parameter_key;
DROP INDEX IF EXISTS idx_game_instance_parameter_game_instance_id;

-- Drop constraints
ALTER TABLE game_instance_parameter DROP CONSTRAINT IF EXISTS game_instance_parameter_unique_key_per_instance;

-- Drop tables
DROP TABLE IF EXISTS game_instance_parameter;
