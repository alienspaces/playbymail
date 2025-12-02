-- Remove item_category and equipment_slot from adventure_game_item table

ALTER TABLE adventure_game_item
    DROP COLUMN IF EXISTS item_category,
    DROP COLUMN IF EXISTS equipment_slot;

