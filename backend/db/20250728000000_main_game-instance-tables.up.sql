-- Create table for game instances (distinct play sessions of a game)
CREATE TABLE IF NOT EXISTS game_instance (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id    UUID NOT NULL REFERENCES game(id), -- The game this instance is based on
    status     VARCHAR(50) NOT NULL DEFAULT 'created',
    current_turn INTEGER NOT NULL DEFAULT 0,
    max_turns  INTEGER,
    turn_deadline_hours INTEGER DEFAULT 168, -- 7 days default
    last_turn_processed_at TIMESTAMPTZ,
    next_turn_deadline TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    game_config JSONB, -- Flexible configuration for game-specific parameters
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- When this instance was created
    updated_at TIMESTAMPTZ, -- When this instance was last updated
    deleted_at TIMESTAMPTZ -- When this instance was deleted (soft delete)
);

-- Add constraints for game_instance
ALTER TABLE game_instance ADD CONSTRAINT game_instance_status_check 
    CHECK (status IN ('created', 'starting', 'running', 'paused', 'completed', 'cancelled'));

ALTER TABLE game_instance ADD CONSTRAINT game_instance_turn_check 
    CHECK (current_turn >= 0);

ALTER TABLE game_instance ADD CONSTRAINT game_instance_max_turns_check 
    CHECK (max_turns IS NULL OR max_turns > 0);

ALTER TABLE game_instance ADD CONSTRAINT game_instance_turn_deadline_check 
    CHECK (turn_deadline_hours > 0);

COMMENT ON TABLE game_instance IS 'Tracks a specific play session or instance of a game.';
COMMENT ON COLUMN game_instance.id IS 'Unique identifier for the game instance.';
COMMENT ON COLUMN game_instance.game_id IS 'The game this instance is based on.';
COMMENT ON COLUMN game_instance.status IS 'Current status of the game instance';
COMMENT ON COLUMN game_instance.current_turn IS 'Current turn number (0-based)';
COMMENT ON COLUMN game_instance.max_turns IS 'Maximum number of turns (NULL for unlimited)';
COMMENT ON COLUMN game_instance.turn_deadline_hours IS 'Hours allowed for each turn';
COMMENT ON COLUMN game_instance.last_turn_processed_at IS 'When the last turn was processed';
COMMENT ON COLUMN game_instance.next_turn_deadline IS 'Deadline for the next turn submission';
COMMENT ON COLUMN game_instance.started_at IS 'When the game instance was started';
COMMENT ON COLUMN game_instance.completed_at IS 'When the game instance was completed';
COMMENT ON COLUMN game_instance.game_config IS 'Game-specific configuration parameters';
COMMENT ON COLUMN game_instance.created_at IS 'When this instance was created.';
COMMENT ON COLUMN game_instance.updated_at IS 'When this instance was last updated.';
COMMENT ON COLUMN game_instance.deleted_at IS 'When this instance was deleted (soft delete).';

-- Create table for game location instances (placement of a location in a game instance)
CREATE TABLE IF NOT EXISTS adventure_game_location_instance (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id           UUID NOT NULL REFERENCES game(id), -- The game this instance is based on
    game_instance_id  UUID NOT NULL REFERENCES game_instance(id), -- The game instance this location instance belongs to
    adventure_game_location_id  UUID NOT NULL REFERENCES adventure_game_location(id), -- The base location this instance represents
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);

COMMENT ON TABLE adventure_game_location_instance IS 'Tracks a specific instance of a location within a game instance.';
COMMENT ON COLUMN adventure_game_location_instance.id IS 'Unique identifier for the game location instance.';
COMMENT ON COLUMN adventure_game_location_instance.game_id IS 'The game this instance is based on.';
COMMENT ON COLUMN adventure_game_location_instance.game_instance_id IS 'The game instance this location instance belongs to.';
COMMENT ON COLUMN adventure_game_location_instance.adventure_game_location_id IS 'The base location this instance represents.';
COMMENT ON COLUMN adventure_game_location_instance.created_at IS 'When this location instance was created.';
COMMENT ON COLUMN adventure_game_location_instance.updated_at IS 'When this location instance was last updated.';
COMMENT ON COLUMN adventure_game_location_instance.deleted_at IS 'When this location instance was deleted (soft delete).';

-- Add game_creature_instance table
CREATE TABLE IF NOT EXISTS adventure_game_creature_instance (
    id                                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id                             UUID NOT NULL REFERENCES game(id),
    adventure_game_creature_id          UUID NOT NULL REFERENCES adventure_game_creature(id),
    game_instance_id                    UUID NOT NULL REFERENCES game_instance(id),
    adventure_game_location_instance_id UUID REFERENCES adventure_game_location_instance(id),
    health                              INTEGER NOT NULL DEFAULT 100,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create table for adventure_game_character_instance (placement of a character in a game instance)
CREATE TABLE IF NOT EXISTS adventure_game_character_instance (
    id                                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id                             UUID NOT NULL REFERENCES game(id), -- The game this character instance belongs to
    game_instance_id                    UUID NOT NULL REFERENCES game_instance(id), -- The game instance this character instance belongs to
    adventure_game_character_id         UUID NOT NULL REFERENCES adventure_game_character(id), -- The base character this instance represents
    adventure_game_location_instance_id UUID REFERENCES adventure_game_location_instance(id), -- The location this character instance is at (if applicable)
    health                              INTEGER NOT NULL DEFAULT 100, -- The health of the character instance
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);

COMMENT ON TABLE adventure_game_character_instance IS 'Tracks a specific instance of a character within a game instance.';
COMMENT ON COLUMN adventure_game_character_instance.id IS 'Unique identifier for the character instance.';
COMMENT ON COLUMN adventure_game_character_instance.game_id IS 'The game this character instance belongs to.';
COMMENT ON COLUMN adventure_game_character_instance.game_instance_id IS 'The game instance this character instance belongs to.';
COMMENT ON COLUMN adventure_game_character_instance.adventure_game_character_id IS 'The base character this instance represents.';
COMMENT ON COLUMN adventure_game_character_instance.adventure_game_location_instance_id IS 'The location this character instance is at (if applicable).';
COMMENT ON COLUMN adventure_game_character_instance.health IS 'The health of the character instance.';
COMMENT ON COLUMN adventure_game_character_instance.created_at IS 'When this character instance was created.';
COMMENT ON COLUMN adventure_game_character_instance.updated_at IS 'When this character instance was last updated.';
COMMENT ON COLUMN adventure_game_character_instance.deleted_at IS 'When this character instance was deleted (soft delete).';

-- Create table for item instances (can be at a location, in a character, or in a creature)
CREATE TABLE IF NOT EXISTS adventure_game_item_instance (
    id                                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id                              UUID NOT NULL REFERENCES game(id), -- The game this item instance belongs to
    adventure_game_item_id               UUID NOT NULL REFERENCES adventure_game_item(id), -- The item type this instance represents
    game_instance_id                     UUID NOT NULL REFERENCES game_instance(id), -- The game instance this item instance belongs to
    adventure_game_location_instance_id  UUID REFERENCES adventure_game_location_instance(id), -- The location this item instance is at (if applicable)
    adventure_game_character_instance_id UUID REFERENCES adventure_game_character_instance(id), -- The character this item instance is on (if applicable)
    adventure_game_creature_instance_id  UUID REFERENCES adventure_game_creature_instance(id), -- The creature this item instance is on (if applicable)
    is_equipped         BOOLEAN NOT NULL DEFAULT FALSE, -- Whether this item instance is currently equipped
    is_used             BOOLEAN NOT NULL DEFAULT FALSE, -- Whether this item instance has been used (if applicable)
    uses_remaining      INTEGER, -- Number of uses left for this instance (if applicable)
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- When this instance was created
    updated_at          TIMESTAMPTZ, -- When this instance was last updated
    deleted_at          TIMESTAMPTZ, -- When this instance was deleted (soft delete)
    CONSTRAINT only_one_owner CHECK (
        (adventure_game_location_instance_id IS NOT NULL)::int + (adventure_game_character_instance_id IS NOT NULL)::int + (adventure_game_creature_instance_id IS NOT NULL)::int = 1
    ) -- Enforces that the item is only in one place at a time
);

COMMENT ON TABLE adventure_game_item_instance IS 'Tracks a specific instance of a game item, which may be at a location, in a character inventory, or in a creature inventory.';
COMMENT ON COLUMN adventure_game_item_instance.id IS 'Unique identifier for the item instance.';
COMMENT ON COLUMN adventure_game_item_instance.game_id IS 'The game this item instance belongs to.';
COMMENT ON COLUMN adventure_game_item_instance.adventure_game_item_id IS 'The item type this instance represents.';
COMMENT ON COLUMN adventure_game_item_instance.game_instance_id IS 'The game instance this item instance belongs to.';
COMMENT ON COLUMN adventure_game_item_instance.adventure_game_location_instance_id IS 'If set, the item is at this location.';
COMMENT ON COLUMN adventure_game_item_instance.adventure_game_character_instance_id IS 'If set, the item is in this character''s inventory.';
COMMENT ON COLUMN adventure_game_item_instance.adventure_game_creature_instance_id IS 'If set, the item is in this creature''s inventory.';
COMMENT ON COLUMN adventure_game_item_instance.is_equipped IS 'Whether this item instance is currently equipped.';
COMMENT ON COLUMN adventure_game_item_instance.is_used IS 'Whether this item instance has been used (if applicable).';
COMMENT ON COLUMN adventure_game_item_instance.uses_remaining IS 'Number of uses left for this instance (if applicable).';
COMMENT ON COLUMN adventure_game_item_instance.created_at IS 'When this instance was created.';
COMMENT ON COLUMN adventure_game_item_instance.updated_at IS 'When this instance was last updated.';
COMMENT ON COLUMN adventure_game_item_instance.deleted_at IS 'When this instance was deleted (soft delete).';

-- Create table for location link item requirements
CREATE TABLE IF NOT EXISTS adventure_game_location_link_requirement (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id                         UUID NOT NULL REFERENCES game(id), -- The game this requirement belongs to
    adventure_game_location_link_id UUID NOT NULL REFERENCES adventure_game_location_link(id), -- The location link this requirement applies to
    adventure_game_item_id          UUID NOT NULL REFERENCES adventure_game_item(id), -- The item type required to traverse the link
    quantity                INTEGER NOT NULL DEFAULT 1, -- How many of this item are required
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- When this requirement was created
    updated_at              TIMESTAMPTZ, -- When this requirement was last updated
    deleted_at              TIMESTAMPTZ -- When this requirement was deleted (soft delete)
);

COMMENT ON TABLE adventure_game_location_link_requirement IS 'Specifies which items (and how many) are required to traverse a location link.';
COMMENT ON COLUMN adventure_game_location_link_requirement.id IS 'Unique identifier for the requirement.';
COMMENT ON COLUMN adventure_game_location_link_requirement.game_id IS 'The game this requirement belongs to.';
COMMENT ON COLUMN adventure_game_location_link_requirement.adventure_game_location_link_id IS 'The location link this requirement applies to.';
COMMENT ON COLUMN adventure_game_location_link_requirement.adventure_game_item_id IS 'The item type required to traverse the link.';
COMMENT ON COLUMN adventure_game_location_link_requirement.quantity IS 'How many of this item are required.';
COMMENT ON COLUMN adventure_game_location_link_requirement.created_at IS 'When this requirement was created.';
COMMENT ON COLUMN adventure_game_location_link_requirement.updated_at IS 'When this requirement was last updated.';
COMMENT ON COLUMN adventure_game_location_link_requirement.deleted_at IS 'When this requirement was deleted (soft delete).'; 
