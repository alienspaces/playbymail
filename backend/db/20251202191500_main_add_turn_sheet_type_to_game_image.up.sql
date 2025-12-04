-- Add turn_sheet_type column to game_image table
-- When type is 'turn_sheet_background', turn_sheet_type must be provided
-- This allows the same handler to support all turn sheet types

ALTER TABLE public.game_image ADD COLUMN turn_sheet_type VARCHAR(50);

-- Update unique constraint to include turn_sheet_type
-- This allows multiple turn sheet background images per game (one per turn sheet type)
ALTER TABLE public.game_image DROP CONSTRAINT IF EXISTS game_image_game_id_record_id_type_key;

-- Add new unique constraint including turn_sheet_type
-- For turn_sheet_background: unique on (game_id, record_id, type, turn_sheet_type)
-- For other types: unique on (game_id, record_id, type) with turn_sheet_type = NULL
CREATE UNIQUE INDEX game_image_unique_turn_sheet_background 
    ON public.game_image(game_id, record_id, type, turn_sheet_type) 
    WHERE type = 'turn_sheet_background' AND turn_sheet_type IS NOT NULL;

CREATE UNIQUE INDEX game_image_unique_other_types 
    ON public.game_image(game_id, record_id, type) 
    WHERE type != 'turn_sheet_background' OR turn_sheet_type IS NULL;

-- Add check constraint: turn_sheet_type is required when type is 'turn_sheet_background'
ALTER TABLE public.game_image ADD CONSTRAINT game_image_turn_sheet_type_check 
    CHECK (
        (type = 'turn_sheet_background' AND turn_sheet_type IS NOT NULL) OR
        (type != 'turn_sheet_background')
    );

-- Add index for querying by turn_sheet_type
CREATE INDEX idx_game_image_turn_sheet_type ON public.game_image(turn_sheet_type);

-- Add comment
COMMENT ON COLUMN public.game_image.turn_sheet_type IS 'Turn sheet type when type is turn_sheet_background (e.g., adventure_game_join_game, adventure_game_inventory_management). Required when type is turn_sheet_background.';

