-- Drop game configuration tables

-- Drop indexes
DROP INDEX IF EXISTS idx_game_instance_configuration_config_key;
DROP INDEX IF EXISTS idx_game_instance_configuration_game_instance_id;
DROP INDEX IF EXISTS idx_game_configuration_config_key;
DROP INDEX IF EXISTS idx_game_configuration_game_type;

-- Drop constraints
ALTER TABLE game_instance_configuration DROP CONSTRAINT IF EXISTS game_instance_configuration_unique_key_per_instance;
ALTER TABLE game_instance_configuration DROP CONSTRAINT IF EXISTS game_instance_configuration_value_type_check;
ALTER TABLE game_configuration DROP CONSTRAINT IF EXISTS game_configuration_unique_key_per_type;
ALTER TABLE game_configuration DROP CONSTRAINT IF EXISTS game_configuration_value_type_check;
ALTER TABLE game_configuration DROP CONSTRAINT IF EXISTS game_configuration_game_type_check;

-- Drop tables
DROP TABLE IF EXISTS game_instance_configuration;
DROP TABLE IF EXISTS game_configuration;
