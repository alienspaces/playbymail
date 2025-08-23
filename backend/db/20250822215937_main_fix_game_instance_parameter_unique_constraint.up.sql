-- Fix game instance parameter unique constraint to only apply to non-deleted records
-- This prevents the constraint violation when trying to recreate a parameter after deletion

-- Drop the old constraint that doesn't consider deleted_at
ALTER TABLE game_instance_parameter DROP CONSTRAINT IF EXISTS game_instance_parameter_unique_key_per_instance;

-- Create a partial unique index that only applies to non-deleted records
-- This effectively creates a unique constraint that ignores deleted records
CREATE UNIQUE INDEX idx_game_instance_parameter_unique_key_per_instance 
    ON game_instance_parameter (game_instance_id, parameter_key) 
    WHERE deleted_at IS NULL;

-- Add a comment explaining the index
COMMENT ON INDEX idx_game_instance_parameter_unique_key_per_instance 
    IS 'Ensures unique parameter keys per game instance, but only for non-deleted records';
