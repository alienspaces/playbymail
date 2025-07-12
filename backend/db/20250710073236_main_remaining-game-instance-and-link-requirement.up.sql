-- Add game_creature table
CREATE TABLE IF NOT EXISTS public.game_creature (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.game(id),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add game_creature_instance table

CREATE TABLE IF NOT EXISTS public.game_creature_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.game(id),
    game_creature_id UUID NOT NULL REFERENCES public.game_creature(id),
    game_instance_id UUID NOT NULL REFERENCES public.game_instance(id),
    game_location_instance_id UUID NOT NULL REFERENCES public.game_location_instance(id),
    is_alive BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
); 

-- Create table for game character instances (placement of a character in a game instance)
CREATE TABLE game_character_instance (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_instance_id  UUID NOT NULL REFERENCES game_instance(id), -- The game instance this character instance belongs to
    game_character_id  UUID NOT NULL REFERENCES game_character(id), -- The base character this instance represents
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);

COMMENT ON TABLE game_character_instance IS 'Tracks a specific instance of a character within a game instance.';
COMMENT ON COLUMN game_character_instance.id IS 'Unique identifier for the character instance.';
COMMENT ON COLUMN game_character_instance.game_instance_id IS 'The game instance this character instance belongs to.';
COMMENT ON COLUMN game_character_instance.game_character_id IS 'The base character this instance represents.';
COMMENT ON COLUMN game_character_instance.created_at IS 'When this character instance was created.';
COMMENT ON COLUMN game_character_instance.updated_at IS 'When this character instance was last updated.';
COMMENT ON COLUMN game_character_instance.deleted_at IS 'When this character instance was deleted (soft delete).';


-- Create table for item instances (can be at a location, in a character, or in a creature)
CREATE TABLE game_item_instance (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id             UUID NOT NULL REFERENCES game(id), -- The game this item instance belongs to
    game_item_id        UUID NOT NULL REFERENCES game_item(id), -- The item type this instance represents
    game_instance_id    UUID NOT NULL REFERENCES game_instance(id), -- The game instance this item instance belongs to
    game_location_instance_id UUID REFERENCES game_location_instance(id), -- The location this item instance is at (if applicable)
    game_character_instance_id UUID REFERENCES game_character_instance(id), -- The character this item instance is on (if applicable)
    game_creature_instance_id UUID REFERENCES game_creature_instance(id), -- The creature this item instance is on (if applicable)
    is_equipped         BOOLEAN NOT NULL DEFAULT FALSE, -- Whether this item instance is currently equipped
    is_used             BOOLEAN NOT NULL DEFAULT FALSE, -- Whether this item instance has been used (if applicable)
    uses_remaining      INTEGER, -- Number of uses left for this instance (if applicable)
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- When this instance was created
    updated_at          TIMESTAMPTZ, -- When this instance was last updated
    deleted_at          TIMESTAMPTZ, -- When this instance was deleted (soft delete)
    CONSTRAINT only_one_owner CHECK (
        (game_location_instance_id IS NOT NULL)::int + (game_character_instance_id IS NOT NULL)::int + (game_creature_instance_id IS NOT NULL)::int = 1
    ) -- Enforces that the item is only in one place at a time
);

COMMENT ON TABLE game_item_instance IS 'Tracks a specific instance of a game item, which may be at a location, in a character inventory, or in a creature inventory.';
COMMENT ON COLUMN game_item_instance.id IS 'Unique identifier for the item instance.';
COMMENT ON COLUMN game_item_instance.game_id IS 'The game this item instance belongs to.';
COMMENT ON COLUMN game_item_instance.game_item_id IS 'The item type this instance represents.';
COMMENT ON COLUMN game_item_instance.game_instance_id IS 'The game instance this item instance belongs to.';
COMMENT ON COLUMN game_item_instance.game_location_instance_id IS 'If set, the item is at this location.';
COMMENT ON COLUMN game_item_instance.game_character_instance_id IS 'If set, the item is in this character''s inventory.';
COMMENT ON COLUMN game_item_instance.game_creature_instance_id IS 'If set, the item is in this creature''s inventory.';
COMMENT ON COLUMN game_item_instance.is_equipped IS 'Whether this item instance is currently equipped.';
COMMENT ON COLUMN game_item_instance.is_used IS 'Whether this item instance has been used (if applicable).';
COMMENT ON COLUMN game_item_instance.uses_remaining IS 'Number of uses left for this instance (if applicable).';
COMMENT ON COLUMN game_item_instance.created_at IS 'When this instance was created.';
COMMENT ON COLUMN game_item_instance.updated_at IS 'When this instance was last updated.';
COMMENT ON COLUMN game_item_instance.deleted_at IS 'When this instance was deleted (soft delete).';

-- Create table for location link item requirements
CREATE TABLE game_location_link_requirement (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_location_link_id   UUID NOT NULL REFERENCES game_location_link(id), -- The location link this requirement applies to
    game_item_id            UUID NOT NULL REFERENCES game_item(id), -- The item type required to traverse the link
    quantity                INTEGER NOT NULL DEFAULT 1, -- How many of this item are required
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- When this requirement was created
    updated_at              TIMESTAMPTZ, -- When this requirement was last updated
    deleted_at              TIMESTAMPTZ -- When this requirement was deleted (soft delete)
);

COMMENT ON TABLE game_location_link_requirement IS 'Specifies which items (and how many) are required to traverse a location link.';
COMMENT ON COLUMN game_location_link_requirement.id IS 'Unique identifier for the requirement.';
COMMENT ON COLUMN game_location_link_requirement.game_location_link_id IS 'The location link this requirement applies to.';
COMMENT ON COLUMN game_location_link_requirement.game_item_id IS 'The item type required to traverse the link.';
COMMENT ON COLUMN game_location_link_requirement.quantity IS 'How many of this item are required.';
COMMENT ON COLUMN game_location_link_requirement.created_at IS 'When this requirement was created.';
COMMENT ON COLUMN game_location_link_requirement.updated_at IS 'When this requirement was last updated.';
COMMENT ON COLUMN game_location_link_requirement.deleted_at IS 'When this requirement was deleted (soft delete).'; 