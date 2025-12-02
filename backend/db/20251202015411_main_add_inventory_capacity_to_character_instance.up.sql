-- Add inventory_capacity to adventure_game_character_instance table
-- Defines how many items a character can carry

ALTER TABLE adventure_game_character_instance
    ADD COLUMN inventory_capacity INTEGER NOT NULL DEFAULT 10 CHECK (inventory_capacity > 0);

-- Add comment
COMMENT ON COLUMN adventure_game_character_instance.inventory_capacity IS 'Maximum number of items the character can carry in inventory. Default is 10.';

