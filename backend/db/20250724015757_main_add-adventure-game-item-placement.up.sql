-- Create table for adventure_game_item_placement (configuration for item placement in locations)
CREATE TABLE IF NOT EXISTS public.adventure_game_item_placement (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id                         UUID NOT NULL REFERENCES public.game(id), -- The game this placement configuration belongs to
    adventure_game_item_id          UUID NOT NULL REFERENCES public.adventure_game_item(id), -- The item type to be placed
    adventure_game_location_id      UUID NOT NULL REFERENCES public.adventure_game_location(id), -- The location where the item should be placed
    initial_count                   INTEGER NOT NULL DEFAULT 1, -- How many of this item should exist at this location when game instance is created
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- When this placement configuration was created
    updated_at                      TIMESTAMPTZ, -- When this placement configuration was last updated
    deleted_at                      TIMESTAMPTZ, -- When this placement configuration was deleted (soft delete)
    UNIQUE(game_id, adventure_game_item_id, adventure_game_location_id) -- Ensure unique placement configuration per game/item/location combination
);

COMMENT ON TABLE adventure_game_item_placement IS 'Configuration for placing items in specific locations when game instances are created.';
COMMENT ON COLUMN adventure_game_item_placement.id IS 'Unique identifier for the placement configuration.';
COMMENT ON COLUMN adventure_game_item_placement.game_id IS 'The game this placement configuration belongs to.';
COMMENT ON COLUMN adventure_game_item_placement.adventure_game_item_id IS 'The item type to be placed.';
COMMENT ON COLUMN adventure_game_item_placement.adventure_game_location_id IS 'The location where the item should be placed.';
COMMENT ON COLUMN adventure_game_item_placement.initial_count IS 'How many of this item should exist at this location when game instance is created.';
COMMENT ON COLUMN adventure_game_item_placement.created_at IS 'When this placement configuration was created.';
COMMENT ON COLUMN adventure_game_item_placement.updated_at IS 'When this placement configuration was last updated.';
COMMENT ON COLUMN adventure_game_item_placement.deleted_at IS 'When this placement configuration was deleted (soft delete).';
