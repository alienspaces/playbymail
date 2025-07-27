-- Add runtime parameters to game instances for state management and turn processing
ALTER TABLE adventure_game_instance ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'created';
ALTER TABLE adventure_game_instance ADD COLUMN current_turn INTEGER NOT NULL DEFAULT 0;
ALTER TABLE adventure_game_instance ADD COLUMN max_turns INTEGER;
ALTER TABLE adventure_game_instance ADD COLUMN turn_deadline_hours INTEGER DEFAULT 168; -- 7 days default
ALTER TABLE adventure_game_instance ADD COLUMN last_turn_processed_at TIMESTAMPTZ;
ALTER TABLE adventure_game_instance ADD COLUMN next_turn_deadline TIMESTAMPTZ;
ALTER TABLE adventure_game_instance ADD COLUMN started_at TIMESTAMPTZ;
ALTER TABLE adventure_game_instance ADD COLUMN completed_at TIMESTAMPTZ;
ALTER TABLE adventure_game_instance ADD COLUMN game_config JSONB; -- Flexible configuration for game-specific parameters

-- Add constraints
ALTER TABLE adventure_game_instance ADD CONSTRAINT adventure_game_instance_status_check 
    CHECK (status IN ('created', 'starting', 'running', 'paused', 'completed', 'cancelled'));

ALTER TABLE adventure_game_instance ADD CONSTRAINT adventure_game_instance_turn_check 
    CHECK (current_turn >= 0);

ALTER TABLE adventure_game_instance ADD CONSTRAINT adventure_game_instance_max_turns_check 
    CHECK (max_turns IS NULL OR max_turns > 0);

ALTER TABLE adventure_game_instance ADD CONSTRAINT adventure_game_instance_turn_deadline_check 
    CHECK (turn_deadline_hours > 0);

-- Add comments
COMMENT ON COLUMN adventure_game_instance.status IS 'Current status of the game instance';
COMMENT ON COLUMN adventure_game_instance.current_turn IS 'Current turn number (0-based)';
COMMENT ON COLUMN adventure_game_instance.max_turns IS 'Maximum number of turns (NULL for unlimited)';
COMMENT ON COLUMN adventure_game_instance.turn_deadline_hours IS 'Hours allowed for each turn';
COMMENT ON COLUMN adventure_game_instance.last_turn_processed_at IS 'When the last turn was processed';
COMMENT ON COLUMN adventure_game_instance.next_turn_deadline IS 'Deadline for the next turn submission';
COMMENT ON COLUMN adventure_game_instance.started_at IS 'When the game instance was started';
COMMENT ON COLUMN adventure_game_instance.completed_at IS 'When the game instance was completed';
COMMENT ON COLUMN adventure_game_instance.game_config IS 'Game-specific configuration parameters';
