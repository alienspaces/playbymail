
-- Create table for game instance parameters (runtime parameter values)
CREATE TABLE IF NOT EXISTS game_instance_parameter (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_instance_id  UUID NOT NULL REFERENCES game_instance(id),
    parameter_key     VARCHAR(100) NOT NULL,
    parameter_value   TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);

-- Add constraints for game_instance_parameter
ALTER TABLE game_instance_parameter ADD CONSTRAINT game_instance_parameter_unique_key_per_instance 
    UNIQUE (game_instance_id, parameter_key);

-- Add indexes
CREATE INDEX idx_game_instance_parameter_game_instance_id ON game_instance_parameter(game_instance_id);
CREATE INDEX idx_game_instance_parameter_parameter_key ON game_instance_parameter(parameter_key);

-- Add comments
COMMENT ON TABLE game_instance_parameter IS 'Runtime parameter values for specific game instances';
COMMENT ON COLUMN game_instance_parameter.id IS 'Unique identifier for the game instance parameter';
COMMENT ON COLUMN game_instance_parameter.game_instance_id IS 'The game instance this parameter belongs to';
COMMENT ON COLUMN game_instance_parameter.parameter_key IS 'Parameter key name';
COMMENT ON COLUMN game_instance_parameter.parameter_value IS 'Parameter value';
