-- Create table for game instances (distinct play sessions of a game)
CREATE TABLE game_instance (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id    UUID NOT NULL REFERENCES game(id), -- The game this instance is based on
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- When this instance was created
    updated_at TIMESTAMPTZ, -- When this instance was last updated
    deleted_at TIMESTAMPTZ -- When this instance was deleted (soft delete)
);

COMMENT ON TABLE game_instance IS 'Tracks a specific play session or instance of a game.';
COMMENT ON COLUMN game_instance.id IS 'Unique identifier for the game instance.';
COMMENT ON COLUMN game_instance.game_id IS 'The game this instance is based on.';
COMMENT ON COLUMN game_instance.created_at IS 'When this instance was created.';
COMMENT ON COLUMN game_instance.updated_at IS 'When this instance was last updated.';
COMMENT ON COLUMN game_instance.deleted_at IS 'When this instance was deleted (soft delete).';

-- Create table for game location instances (placement of a location in a game instance)
CREATE TABLE game_location_instance (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_instance_id  UUID NOT NULL REFERENCES game_instance(id), -- The game instance this location instance belongs to
    game_location_id  UUID NOT NULL REFERENCES game_location(id), -- The base location this instance represents
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);

COMMENT ON TABLE game_location_instance IS 'Tracks a specific instance of a location within a game instance.';
COMMENT ON COLUMN game_location_instance.id IS 'Unique identifier for the game location instance.';
COMMENT ON COLUMN game_location_instance.game_instance_id IS 'The game instance this location instance belongs to.';
COMMENT ON COLUMN game_location_instance.game_location_id IS 'The base location this instance represents.';
COMMENT ON COLUMN game_location_instance.created_at IS 'When this location instance was created.';
COMMENT ON COLUMN game_location_instance.updated_at IS 'When this location instance was last updated.';
COMMENT ON COLUMN game_location_instance.deleted_at IS 'When this location instance was deleted (soft delete).'; 