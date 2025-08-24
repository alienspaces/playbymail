-- Update game instance status constraint to match Go constants
-- First, update any existing 'starting' status values to 'started' to match Go constants
UPDATE game_instance SET status = 'started' WHERE status = 'starting';

-- Drop the old constraint that uses 'starting' and 'running'
ALTER TABLE game_instance DROP CONSTRAINT IF EXISTS game_instance_status_check;

-- Add the new constraint that uses 'started' to match Go constants
ALTER TABLE game_instance ADD CONSTRAINT game_instance_status_check 
    CHECK (status IN ('created', 'started', 'paused', 'completed', 'cancelled'));
