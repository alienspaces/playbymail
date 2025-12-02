-- Add equipment_slot to adventure_game_item_instance table
-- Tracks which equipment slot this item instance is equipped in

ALTER TABLE adventure_game_item_instance
    ADD COLUMN equipment_slot VARCHAR(50);

-- Add comment
COMMENT ON COLUMN adventure_game_item_instance.equipment_slot IS 'Equipment slot this item instance is equipped in (e.g., weapon, armor_body, jewelry_ring). NULL when not equipped or not equippable.';

