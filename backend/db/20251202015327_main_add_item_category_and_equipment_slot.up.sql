-- Add item_category and equipment_slot to adventure_game_item table
-- Supports inventory management and equipment system

ALTER TABLE adventure_game_item
    ADD COLUMN item_category VARCHAR(50) CHECK (item_category IN ('weapon', 'armor', 'clothing', 'jewelry', 'consumable', 'misc')),
    ADD COLUMN equipment_slot VARCHAR(50);

-- Add comments
COMMENT ON COLUMN adventure_game_item.item_category IS 'Item category: weapon, armor, clothing, jewelry, consumable, or misc';
COMMENT ON COLUMN adventure_game_item.equipment_slot IS 'Equipment slot this item occupies when equipped (e.g., weapon, armor_body, jewelry_ring). NULL for non-equippable items.';

