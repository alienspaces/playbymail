-- Drop simplified turn sheet schema
-- This migration removes the simplified turn sheet design

-- Drop mapping table first (due to foreign key constraints)
DROP TABLE IF EXISTS adventure_game_character_instance_turn_sheet;

-- Drop main turn sheet table
DROP TABLE IF EXISTS game_turn_sheet;
