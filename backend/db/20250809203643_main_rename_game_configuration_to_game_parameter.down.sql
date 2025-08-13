-- Rename game_parameter table back to game_configuration

-- Drop indexes first (they reference the current table name)
DROP INDEX IF EXISTS idx_game_parameter_game_type;
DROP INDEX IF EXISTS idx_game_parameter_config_key;

-- Drop constraints that reference the current table name
ALTER TABLE game_parameter DROP CONSTRAINT IF EXISTS game_parameter_game_type_check;
ALTER TABLE game_parameter DROP CONSTRAINT IF EXISTS game_parameter_value_type_check;
ALTER TABLE game_parameter DROP CONSTRAINT IF EXISTS game_parameter_unique_key_per_type;

-- Rename the table back
ALTER TABLE game_parameter RENAME TO game_configuration;

-- Recreate constraints with original names
ALTER TABLE game_configuration ADD CONSTRAINT game_configuration_game_type_check 
    CHECK (game_type IN ('adventure', 'strategy', 'puzzle', 'simulation'));

ALTER TABLE game_configuration ADD CONSTRAINT game_configuration_value_type_check 
    CHECK (value_type IN ('string', 'integer', 'boolean', 'json'));

ALTER TABLE game_configuration ADD CONSTRAINT game_configuration_unique_key_per_type 
    UNIQUE (game_type, config_key);

-- Remove the column that was added in the up migration
ALTER TABLE game_configuration DROP COLUMN IF EXISTS is_global;

-- Add back the columns that were removed in the up migration
ALTER TABLE game_configuration ADD COLUMN ui_hint VARCHAR(50);
ALTER TABLE game_configuration ADD COLUMN validation_rules TEXT;

-- Recreate indexes with original names
CREATE INDEX idx_game_configuration_game_type ON game_configuration(game_type);
CREATE INDEX idx_game_configuration_config_key ON game_configuration(config_key);

-- Restore original comments
COMMENT ON TABLE game_configuration IS 'Defines available configuration options for different game types';
COMMENT ON COLUMN game_configuration.id IS 'Unique identifier for the game configuration';
COMMENT ON COLUMN game_configuration.game_type IS 'Type of game this configuration applies to';
COMMENT ON COLUMN game_configuration.config_key IS 'Configuration key name';
COMMENT ON COLUMN game_configuration.value_type IS 'Data type for this configuration value';
COMMENT ON COLUMN game_configuration.default_value IS 'Default value for this configuration';
COMMENT ON COLUMN game_configuration.is_required IS 'Whether this configuration is required';
COMMENT ON COLUMN game_configuration.description IS 'Description of this configuration option';
COMMENT ON COLUMN game_configuration.ui_hint IS 'UI hint for form generation';
COMMENT ON COLUMN game_configuration.validation_rules IS 'JSON validation rules for this configuration';
