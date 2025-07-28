-- Create game configuration tables

-- Create table for game configuration schemas (defines available configuration options for game types)
CREATE TABLE IF NOT EXISTS game_configuration (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_type         VARCHAR(50) NOT NULL,
    config_key        VARCHAR(100) NOT NULL,
    value_type        VARCHAR(20) NOT NULL,
    default_value     TEXT,
    is_required       BOOLEAN NOT NULL DEFAULT false,
    description       TEXT,
    ui_hint           VARCHAR(50),
    validation_rules  TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);

-- Create table for game instance configurations (runtime configuration values)
CREATE TABLE IF NOT EXISTS game_instance_configuration (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_instance_id  UUID NOT NULL REFERENCES game_instance(id),
    config_key        VARCHAR(100) NOT NULL,
    value_type        VARCHAR(20) NOT NULL,
    string_value      TEXT,
    integer_value     INTEGER,
    boolean_value     BOOLEAN,
    json_value        TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);

-- Add constraints for game_configuration
ALTER TABLE game_configuration ADD CONSTRAINT game_configuration_game_type_check 
    CHECK (game_type IN ('adventure', 'strategy', 'puzzle', 'simulation'));

ALTER TABLE game_configuration ADD CONSTRAINT game_configuration_value_type_check 
    CHECK (value_type IN ('string', 'integer', 'boolean', 'json'));

ALTER TABLE game_configuration ADD CONSTRAINT game_configuration_unique_key_per_type 
    UNIQUE (game_type, config_key);

-- Add constraints for game_instance_configuration
ALTER TABLE game_instance_configuration ADD CONSTRAINT game_instance_configuration_value_type_check 
    CHECK (value_type IN ('string', 'integer', 'boolean', 'json'));

ALTER TABLE game_instance_configuration ADD CONSTRAINT game_instance_configuration_unique_key_per_instance 
    UNIQUE (game_instance_id, config_key);

-- Add indexes
CREATE INDEX idx_game_configuration_game_type ON game_configuration(game_type);
CREATE INDEX idx_game_configuration_config_key ON game_configuration(config_key);
CREATE INDEX idx_game_instance_configuration_game_instance_id ON game_instance_configuration(game_instance_id);
CREATE INDEX idx_game_instance_configuration_config_key ON game_instance_configuration(config_key);

-- Add comments
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

COMMENT ON TABLE game_instance_configuration IS 'Runtime configuration values for specific game instances';
COMMENT ON COLUMN game_instance_configuration.id IS 'Unique identifier for the game instance configuration';
COMMENT ON COLUMN game_instance_configuration.game_instance_id IS 'The game instance this configuration belongs to';
COMMENT ON COLUMN game_instance_configuration.config_key IS 'Configuration key name';
COMMENT ON COLUMN game_instance_configuration.value_type IS 'Data type for this configuration value';
COMMENT ON COLUMN game_instance_configuration.string_value IS 'String value (when value_type is string)';
COMMENT ON COLUMN game_instance_configuration.integer_value IS 'Integer value (when value_type is integer)';
COMMENT ON COLUMN game_instance_configuration.boolean_value IS 'Boolean value (when value_type is boolean)';
COMMENT ON COLUMN game_instance_configuration.json_value IS 'JSON value (when value_type is json)';
