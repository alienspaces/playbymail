-- Add is_starting_item column to adventure_game_item table
-- This flag marks items that should be automatically assigned to characters when they join a game

ALTER TABLE adventure_game_item
ADD COLUMN is_starting_item BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN adventure_game_item.is_starting_item IS 'If true, this item is automatically assigned to characters when they join the game';

