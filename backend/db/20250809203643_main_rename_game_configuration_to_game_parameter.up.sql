-- Rename game_configuration table to game_parameter

-- Drop indexes first (they reference the old table name)
DROP INDEX IF EXISTS idx_game_configuration_game_type;
DROP INDEX IF EXISTS idx_game_configuration_config_key;

-- Drop constraints that reference the old table name
ALTER TABLE game_configuration DROP CONSTRAINT IF EXISTS game_configuration_game_type_check;
ALTER TABLE game_configuration DROP CONSTRAINT IF EXISTS game_configuration_value_type_check;
ALTER TABLE game_configuration DROP CONSTRAINT IF EXISTS game_configuration_unique_key_per_type;

-- Rename the table
ALTER TABLE game_configuration RENAME TO game_parameter;

-- Recreate constraints with new names
ALTER TABLE game_parameter ADD CONSTRAINT game_parameter_game_type_check 
    CHECK (game_type IN ('adventure', 'strategy', 'puzzle', 'simulation'));

ALTER TABLE game_parameter ADD CONSTRAINT game_parameter_value_type_check 
    CHECK (value_type IN ('string', 'integer', 'boolean', 'json'));

ALTER TABLE game_parameter ADD CONSTRAINT game_parameter_unique_key_per_type 
    UNIQUE (game_type, config_key);

-- Recreate indexes with new names
CREATE INDEX idx_game_parameter_game_type ON game_parameter(game_type);
CREATE INDEX idx_game_parameter_config_key ON game_parameter(config_key);

-- Remove columns that don't exist in Go record
ALTER TABLE game_parameter DROP COLUMN IF EXISTS ui_hint;
ALTER TABLE game_parameter DROP COLUMN IF EXISTS validation_rules;

-- Add column that exists in Go record but missing from DB
ALTER TABLE game_parameter ADD COLUMN is_global BOOLEAN NOT NULL DEFAULT false;

-- Update comments
COMMENT ON TABLE game_parameter IS 'Defines available parameters for different game types';
COMMENT ON COLUMN game_parameter.id IS 'Unique identifier for the game parameter';
COMMENT ON COLUMN game_parameter.game_type IS 'Type of game this parameter applies to';
COMMENT ON COLUMN game_parameter.config_key IS 'Parameter key name';
COMMENT ON COLUMN game_parameter.value_type IS 'Data type for this parameter value';
COMMENT ON COLUMN game_parameter.default_value IS 'Default value for this parameter';
COMMENT ON COLUMN game_parameter.is_required IS 'Whether this parameter is required';
COMMENT ON COLUMN game_parameter.description IS 'Description of this parameter option';
COMMENT ON COLUMN game_parameter.is_global IS 'Whether this parameter applies globally across all instances';
