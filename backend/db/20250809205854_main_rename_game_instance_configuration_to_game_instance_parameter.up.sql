-- Rename game_instance_configuration table to game_instance_parameter

-- Drop indexes first (they reference the old table name)
DROP INDEX IF EXISTS idx_game_instance_configuration_game_instance_id;
DROP INDEX IF EXISTS idx_game_instance_configuration_config_key;

-- Drop constraints that reference the old table name
ALTER TABLE game_instance_configuration DROP CONSTRAINT IF EXISTS game_instance_configuration_value_type_check;
ALTER TABLE game_instance_configuration DROP CONSTRAINT IF EXISTS game_instance_configuration_unique_key_per_instance;

-- Rename the table
ALTER TABLE game_instance_configuration RENAME TO game_instance_parameter;

-- Recreate constraints with new names
ALTER TABLE game_instance_parameter ADD CONSTRAINT game_instance_parameter_value_type_check 
    CHECK (value_type IN ('string', 'integer', 'boolean', 'json'));

ALTER TABLE game_instance_parameter ADD CONSTRAINT game_instance_parameter_unique_key_per_instance 
    UNIQUE (game_instance_id, config_key);

-- Recreate indexes with new names
CREATE INDEX idx_game_instance_parameter_game_instance_id ON game_instance_parameter(game_instance_id);
CREATE INDEX idx_game_instance_parameter_config_key ON game_instance_parameter(config_key);

-- Update comments
COMMENT ON TABLE game_instance_parameter IS 'Runtime parameter values for specific game instances';
COMMENT ON COLUMN game_instance_parameter.id IS 'Unique identifier for the game instance parameter';
COMMENT ON COLUMN game_instance_parameter.game_instance_id IS 'The game instance this parameter belongs to';
COMMENT ON COLUMN game_instance_parameter.config_key IS 'Parameter key name';
COMMENT ON COLUMN game_instance_parameter.value_type IS 'Data type for this parameter value';
COMMENT ON COLUMN game_instance_parameter.string_value IS 'String value (when value_type is string)';
COMMENT ON COLUMN game_instance_parameter.integer_value IS 'Integer value (when value_type is integer)';
COMMENT ON COLUMN game_instance_parameter.boolean_value IS 'Boolean value (when value_type is boolean)';
COMMENT ON COLUMN game_instance_parameter.json_value IS 'JSON value (when value_type is json)';
