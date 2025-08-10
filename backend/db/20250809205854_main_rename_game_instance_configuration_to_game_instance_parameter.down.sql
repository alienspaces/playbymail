-- Rename game_instance_parameter table back to game_instance_configuration

-- Drop indexes first (they reference the current table name)
DROP INDEX IF EXISTS idx_game_instance_parameter_game_instance_id;
DROP INDEX IF EXISTS idx_game_instance_parameter_config_key;

-- Drop constraints that reference the current table name
ALTER TABLE game_instance_parameter DROP CONSTRAINT IF EXISTS game_instance_parameter_value_type_check;
ALTER TABLE game_instance_parameter DROP CONSTRAINT IF EXISTS game_instance_parameter_unique_key_per_instance;

-- Rename the table back
ALTER TABLE game_instance_parameter RENAME TO game_instance_configuration;

-- Recreate constraints with original names
ALTER TABLE game_instance_configuration ADD CONSTRAINT game_instance_configuration_value_type_check 
    CHECK (value_type IN ('string', 'integer', 'boolean', 'json'));

ALTER TABLE game_instance_configuration ADD CONSTRAINT game_instance_configuration_unique_key_per_instance 
    UNIQUE (game_instance_id, config_key);

-- Recreate indexes with original names
CREATE INDEX idx_game_instance_configuration_game_instance_id ON game_instance_configuration(game_instance_id);
CREATE INDEX idx_game_instance_configuration_config_key ON game_instance_configuration(config_key);

-- Restore original comments
COMMENT ON TABLE game_instance_configuration IS 'Runtime configuration values for specific game instances';
COMMENT ON COLUMN game_instance_configuration.id IS 'Unique identifier for the game instance configuration';
COMMENT ON COLUMN game_instance_configuration.game_instance_id IS 'The game instance this configuration belongs to';
COMMENT ON COLUMN game_instance_configuration.config_key IS 'Configuration key name';
COMMENT ON COLUMN game_instance_configuration.value_type IS 'Data type for this configuration value';
COMMENT ON COLUMN game_instance_configuration.string_value IS 'String value (when value_type is string)';
COMMENT ON COLUMN game_instance_configuration.integer_value IS 'Integer value (when value_type is integer)';
COMMENT ON COLUMN game_instance_configuration.boolean_value IS 'Boolean value (when value_type is boolean)';
COMMENT ON COLUMN game_instance_configuration.json_value IS 'JSON value (when value_type is json)';
