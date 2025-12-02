-- Remove inventory_capacity from adventure_game_character_instance table

ALTER TABLE adventure_game_character_instance
    DROP COLUMN IF EXISTS inventory_capacity;

