-- Remove runtime parameters from game instances
ALTER TABLE adventure_game_instance DROP CONSTRAINT IF EXISTS adventure_game_instance_status_check;
ALTER TABLE adventure_game_instance DROP CONSTRAINT IF EXISTS adventure_game_instance_turn_check;
ALTER TABLE adventure_game_instance DROP CONSTRAINT IF EXISTS adventure_game_instance_max_turns_check;
ALTER TABLE adventure_game_instance DROP CONSTRAINT IF EXISTS adventure_game_instance_turn_deadline_check;

ALTER TABLE adventure_game_instance DROP COLUMN IF EXISTS status;
ALTER TABLE adventure_game_instance DROP COLUMN IF EXISTS current_turn;
ALTER TABLE adventure_game_instance DROP COLUMN IF EXISTS max_turns;
ALTER TABLE adventure_game_instance DROP COLUMN IF EXISTS turn_deadline_hours;
ALTER TABLE adventure_game_instance DROP COLUMN IF EXISTS last_turn_processed_at;
ALTER TABLE adventure_game_instance DROP COLUMN IF EXISTS next_turn_deadline;
ALTER TABLE adventure_game_instance DROP COLUMN IF EXISTS started_at;
ALTER TABLE adventure_game_instance DROP COLUMN IF EXISTS completed_at;
ALTER TABLE adventure_game_instance DROP COLUMN IF EXISTS game_config;
