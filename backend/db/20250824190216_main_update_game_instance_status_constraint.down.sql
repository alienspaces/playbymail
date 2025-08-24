-- Revert game instance status constraint to original values
-- Drop the new constraint
ALTER TABLE game_instance DROP CONSTRAINT IF EXISTS game_instance_status_check;

-- Add back the original constraint that uses 'starting' and 'running'
ALTER TABLE game_instance ADD CONSTRAINT game_instance_status_check 
    CHECK (status IN ('created', 'starting', 'running', 'paused', 'completed', 'cancelled'));
