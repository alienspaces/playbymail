-- Remove equipment_slot from adventure_game_item_instance table

ALTER TABLE adventure_game_item_instance
    DROP COLUMN IF EXISTS equipment_slot;

