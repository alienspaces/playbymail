-- Remove is_starting_item column from adventure_game_item table

ALTER TABLE adventure_game_item
DROP COLUMN IF EXISTS is_starting_item;

