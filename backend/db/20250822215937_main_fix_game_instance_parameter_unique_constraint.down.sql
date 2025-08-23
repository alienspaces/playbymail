-- Revert game instance parameter unique constraint to original form
-- This restores the constraint that applies to all records (including deleted ones)

-- Drop the partial unique index that only applies to non-deleted records
DROP INDEX IF EXISTS idx_game_instance_parameter_unique_key_per_instance;

-- Restore the original constraint that applies to all records
ALTER TABLE game_instance_parameter ADD CONSTRAINT game_instance_parameter_unique_key_per_instance 
    UNIQUE (game_instance_id, parameter_key);
