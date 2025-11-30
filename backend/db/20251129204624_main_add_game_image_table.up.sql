-- Create game_image table for storing turn sheet artwork
-- Supports game-level turn sheet images (record_id = NULL) and
-- location-specific turn sheet images (record_id = location ID)

CREATE TABLE IF NOT EXISTS public.game_image (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES game(id),
    record_id UUID,
    type VARCHAR(50) NOT NULL CHECK (type IN ('turn_sheet_background', 'asset')),
    image_data BYTEA NOT NULL,
    mime_type VARCHAR(50) NOT NULL CHECK (mime_type IN ('image/webp', 'image/png', 'image/jpeg')),
    file_size INTEGER NOT NULL CHECK (file_size > 0 AND file_size <= 1048576),
    width INTEGER NOT NULL CHECK (width >= 400 AND width <= 4000),
    height INTEGER NOT NULL CHECK (height >= 200 AND height <= 6000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE(game_id, record_id, type)
);

-- Add indexes for performance
CREATE INDEX idx_game_image_game_id ON game_image(game_id);
CREATE INDEX idx_game_image_game_record ON game_image(game_id, record_id);
CREATE INDEX idx_game_image_type ON game_image(type);

-- Add comments
COMMENT ON TABLE game_image IS 'Stores turn sheet artwork images for games. record_id is NULL for game-level turn sheet images, or references a location/asset ID for record-specific images.';
COMMENT ON COLUMN game_image.type IS 'Image type: turn_sheet_background for turn sheet background, asset for future use';
COMMENT ON COLUMN game_image.file_size IS 'File size in bytes, max 1MB (1048576 bytes)';
COMMENT ON COLUMN game_image.width IS 'Image width in pixels, min 400px, max 4000px';
COMMENT ON COLUMN game_image.height IS 'Image height in pixels, min 200px, max 6000px';
