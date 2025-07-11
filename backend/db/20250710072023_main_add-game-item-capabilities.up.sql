-- Add can_be_equipped and can_be_used columns to game_item
ALTER TABLE game_item
    ADD COLUMN can_be_equipped BOOLEAN NOT NULL DEFAULT FALSE;
COMMENT ON COLUMN game_item.can_be_equipped IS 'Whether this item type can be equipped by a character or creature.';

ALTER TABLE game_item
    ADD COLUMN can_be_used BOOLEAN NOT NULL DEFAULT FALSE;
COMMENT ON COLUMN game_item.can_be_used IS 'Whether this item type can be used (e.g., consumed, activated, etc.).';
