-- Consolidated schema migration
-- This migration represents the complete database schema as of 2025-12-09
-- It consolidates all previous migrations into a single file for easier maintenance

-- ============================================================================
-- ACCOUNT TABLES
-- ============================================================================

CREATE TABLE IF NOT EXISTS public.account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    verification_token TEXT,
    verification_token_expires_at TIMESTAMPTZ,
    session_token TEXT,
    session_token_expires_at TIMESTAMPTZ,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT account_status_check CHECK (status IN ('pending_approval', 'active', 'disabled'))
);

COMMENT ON TABLE public.account IS 'User accounts for authentication. Contact information is stored separately in account_contact.';
COMMENT ON COLUMN public.account.status IS 'Current approval status of the account (pending_approval, active, disabled).';

CREATE TABLE IF NOT EXISTS public.account_contact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    postal_address_line1 VARCHAR(255) NOT NULL,
    postal_address_line2 VARCHAR(255),
    state_province VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT account_contact_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.account_contact IS 'Contact information for accounts, including name and postal address. Collected from join game turn sheets.';
COMMENT ON COLUMN public.account_contact.id IS 'Unique identifier for the contact record.';
COMMENT ON COLUMN public.account_contact.account_id IS 'The account this contact belongs to.';
COMMENT ON COLUMN public.account_contact.name IS 'Contact name.';
COMMENT ON COLUMN public.account_contact.postal_address_line1 IS 'Primary postal address line.';
COMMENT ON COLUMN public.account_contact.postal_address_line2 IS 'Secondary postal address line (optional).';
COMMENT ON COLUMN public.account_contact.state_province IS 'State or province.';
COMMENT ON COLUMN public.account_contact.country IS 'Country.';
COMMENT ON COLUMN public.account_contact.postal_code IS 'Postal or ZIP code.';

CREATE INDEX idx_account_contact_account_id ON public.account_contact(account_id);

-- ============================================================================
-- GAME TABLES
-- ============================================================================

CREATE TABLE IF NOT EXISTS public.game (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    game_type VARCHAR(50) NOT NULL,
    description VARCHAR(512) NOT NULL,
    turn_duration_hours INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT game_name_check CHECK (name != ''),
    CONSTRAINT game_type_check CHECK (game_type = 'adventure'),
    CONSTRAINT turn_duration_hours_check CHECK (turn_duration_hours > 0),
    CONSTRAINT game_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id)
);

COMMENT ON COLUMN public.game.description IS 'Game description that appears on the join game turn sheet (required, max 512 characters)';
COMMENT ON COLUMN public.game.account_id IS 'The account that created/owns this game record.';

CREATE INDEX idx_game_account_id ON public.game(account_id);

CREATE TABLE IF NOT EXISTS public.game_image (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    record_id UUID,
    type VARCHAR(50) NOT NULL,
    turn_sheet_type VARCHAR(50),
    image_data BYTEA NOT NULL,
    mime_type VARCHAR(50) NOT NULL,
    file_size INTEGER NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT game_image_file_size_check CHECK (file_size > 0 AND file_size <= 1048576),
    CONSTRAINT game_image_width_check CHECK (width >= 400 AND width <= 4000),
    CONSTRAINT game_image_height_check CHECK (height >= 200 AND height <= 6000),
    CONSTRAINT game_image_mime_type_check CHECK (mime_type IN ('image/webp', 'image/png', 'image/jpeg')),
    CONSTRAINT game_image_type_check CHECK (type IN ('turn_sheet_background', 'asset')),
    CONSTRAINT game_image_turn_sheet_type_check CHECK (
        (type = 'turn_sheet_background' AND turn_sheet_type IS NOT NULL) OR
        (type != 'turn_sheet_background')
    ),
    CONSTRAINT game_image_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id)
);

COMMENT ON TABLE public.game_image IS 'Stores turn sheet artwork images for games. record_id is NULL for game-level turn sheet images, or references a location/asset ID for record-specific images.';
COMMENT ON COLUMN public.game_image.type IS 'Image type: turn_sheet_background for turn sheet background, asset for future use';
COMMENT ON COLUMN public.game_image.file_size IS 'File size in bytes, max 1MB (1048576 bytes)';
COMMENT ON COLUMN public.game_image.width IS 'Image width in pixels, min 400px, max 4000px';
COMMENT ON COLUMN public.game_image.height IS 'Image height in pixels, min 200px, max 6000px';
COMMENT ON COLUMN public.game_image.turn_sheet_type IS 'Turn sheet type when type is turn_sheet_background (e.g., adventure_game_join_game, adventure_game_inventory_management). Required when type is turn_sheet_background.';

CREATE INDEX idx_game_image_game_id ON public.game_image(game_id);
CREATE INDEX idx_game_image_game_record ON public.game_image(game_id, record_id);
CREATE INDEX idx_game_image_type ON public.game_image(type);
CREATE INDEX idx_game_image_turn_sheet_type ON public.game_image(turn_sheet_type);
CREATE UNIQUE INDEX game_image_unique_turn_sheet_background ON public.game_image(game_id, record_id, type, turn_sheet_type)
    WHERE type = 'turn_sheet_background' AND turn_sheet_type IS NOT NULL;
CREATE UNIQUE INDEX game_image_unique_other_types ON public.game_image(game_id, record_id, type)
    WHERE (type != 'turn_sheet_background' OR turn_sheet_type IS NULL);

CREATE TABLE IF NOT EXISTS public.game_subscription (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    account_id UUID NOT NULL,
    account_contact_id UUID,
    subscription_type VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    turn_sheet_token VARCHAR(255),
    turn_sheet_token_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT game_subscription_status_check CHECK (status IN ('pending_approval', 'active', 'revoked')),
    CONSTRAINT game_subscription_subscription_type_check CHECK (subscription_type IN ('Player', 'Manager', 'Designer')),
    CONSTRAINT game_subscription_game_id_account_id_subscription_type_key UNIQUE (game_id, account_id, subscription_type),
    CONSTRAINT game_subscription_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT game_subscription_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id),
    CONSTRAINT game_subscription_account_contact_id_fkey FOREIGN KEY (account_contact_id) REFERENCES public.account_contact(id) ON DELETE SET NULL
);

COMMENT ON TABLE public.game_subscription IS 'Tracks account subscriptions to a game, including Player, Manager, and Collaborator roles.';
COMMENT ON COLUMN public.game_subscription.id IS 'Unique identifier for the subscription.';
COMMENT ON COLUMN public.game_subscription.game_id IS 'The game being subscribed to.';
COMMENT ON COLUMN public.game_subscription.account_id IS 'The subscribing account.';
COMMENT ON COLUMN public.game_subscription.subscription_type IS 'Role: Player, Manager, or Designer.';
COMMENT ON COLUMN public.game_subscription.created_at IS 'When the subscription was created.';
COMMENT ON COLUMN public.game_subscription.updated_at IS 'When the subscription was last updated.';
COMMENT ON COLUMN public.game_subscription.deleted_at IS 'When the subscription was logically deleted.';
COMMENT ON COLUMN public.game_subscription.status IS 'Approval status of the game subscription (pending_approval, active, revoked).';
COMMENT ON COLUMN public.game_subscription.account_contact_id IS 'The contact information used for this subscription. References account_contact collected from join game turn sheets.';
COMMENT ON COLUMN public.game_subscription.turn_sheet_token IS 'Unique token for accessing turn sheets via web viewer. Refreshed with every new turn and automatically expired once player submits their latest turn sheets';
COMMENT ON COLUMN public.game_subscription.turn_sheet_token_expires_at IS 'Expiration timestamp for turn sheet token (3 days from generation). Key is also automatically expired when player submits their latest turn sheets for the related game';

CREATE INDEX idx_game_subscription_turn_sheet_token ON public.game_subscription(turn_sheet_token)
    WHERE turn_sheet_token IS NOT NULL;

CREATE TABLE IF NOT EXISTS public.game_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_subscription_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'created',
    current_turn INTEGER NOT NULL DEFAULT 0,
    required_player_count INTEGER NOT NULL DEFAULT 0,
    is_closed_testing BOOLEAN NOT NULL DEFAULT false,
    closed_testing_join_game_key VARCHAR(255),
    closed_testing_join_game_key_expires_at TIMESTAMPTZ,
    delivery_physical_post BOOLEAN NOT NULL DEFAULT true,
    delivery_physical_local BOOLEAN NOT NULL DEFAULT false,
    delivery_email BOOLEAN NOT NULL DEFAULT false,
    last_turn_processed_at TIMESTAMPTZ,
    next_turn_due_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT game_instance_status_check CHECK (status IN ('created', 'started', 'paused', 'completed', 'cancelled')),
    CONSTRAINT game_instance_turn_check CHECK (current_turn >= 0),
    CONSTRAINT game_instance_required_player_count_check CHECK (required_player_count >= 0),
    CONSTRAINT game_instance_delivery_methods_check CHECK (
        delivery_physical_post = true OR
        delivery_physical_local = true OR
        delivery_email = true
    ),
    CONSTRAINT game_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT game_instance_game_subscription_id_fkey FOREIGN KEY (game_subscription_id) REFERENCES public.game_subscription(id)
);

COMMENT ON TABLE public.game_instance IS 'Tracks a specific play session or instance of a game.';
COMMENT ON COLUMN public.game_instance.id IS 'Unique identifier for the game instance.';
COMMENT ON COLUMN public.game_instance.game_id IS 'The game this instance is based on.';
COMMENT ON COLUMN public.game_instance.status IS 'Current status of the game instance';
COMMENT ON COLUMN public.game_instance.current_turn IS 'Current turn number (0-based)';
COMMENT ON COLUMN public.game_instance.last_turn_processed_at IS 'When the last turn was processed';
COMMENT ON COLUMN public.game_instance.next_turn_due_at IS 'When the next turn is due';
COMMENT ON COLUMN public.game_instance.started_at IS 'When the game instance was started';
COMMENT ON COLUMN public.game_instance.completed_at IS 'When the game instance was completed';
COMMENT ON COLUMN public.game_instance.created_at IS 'When this instance was created.';
COMMENT ON COLUMN public.game_instance.updated_at IS 'When this instance was last updated.';
COMMENT ON COLUMN public.game_instance.deleted_at IS 'When this instance was deleted (soft delete).';
COMMENT ON COLUMN public.game_instance.game_subscription_id IS 'The Manager subscription that created this game instance.';
COMMENT ON COLUMN public.game_instance.delivery_physical_post IS 'Enable physical post delivery (traditional mail-based)';
COMMENT ON COLUMN public.game_instance.delivery_physical_local IS 'Enable physical local delivery (convention/classroom - game master prints locally, players fill at table, manual scanning/submission)';
COMMENT ON COLUMN public.game_instance.delivery_email IS 'Enable email delivery (web-based turn sheet viewer via email links)';
COMMENT ON COLUMN public.game_instance.required_player_count IS 'Minimum number of players required before game can start';
COMMENT ON COLUMN public.game_instance.is_closed_testing IS 'Whether this game instance is in closed testing mode';
COMMENT ON COLUMN public.game_instance.closed_testing_join_game_key IS 'Unique token for closed testing join game links';
COMMENT ON COLUMN public.game_instance.closed_testing_join_game_key_expires_at IS 'Optional expiration timestamp for closed testing join game token';

CREATE INDEX idx_game_instance_game_subscription_id ON public.game_instance(game_subscription_id);
CREATE INDEX idx_game_instance_closed_testing_join_game_key ON public.game_instance(closed_testing_join_game_key)
    WHERE closed_testing_join_game_key IS NOT NULL;

CREATE TABLE IF NOT EXISTS public.game_instance_parameter (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_instance_id UUID NOT NULL,
    parameter_key VARCHAR(100) NOT NULL,
    parameter_value TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT game_instance_parameter_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id)
);

COMMENT ON TABLE public.game_instance_parameter IS 'Runtime parameter values for specific game instances';
COMMENT ON COLUMN public.game_instance_parameter.id IS 'Unique identifier for the game instance parameter';
COMMENT ON COLUMN public.game_instance_parameter.game_instance_id IS 'The game instance this parameter belongs to';
COMMENT ON COLUMN public.game_instance_parameter.parameter_key IS 'Parameter key name';
COMMENT ON COLUMN public.game_instance_parameter.parameter_value IS 'Parameter value';

CREATE INDEX idx_game_instance_parameter_game_instance_id ON public.game_instance_parameter(game_instance_id);
CREATE INDEX idx_game_instance_parameter_parameter_key ON public.game_instance_parameter(parameter_key);
CREATE UNIQUE INDEX idx_game_instance_parameter_unique_key_per_instance ON public.game_instance_parameter(game_instance_id, parameter_key)
    WHERE deleted_at IS NULL;

COMMENT ON INDEX public.idx_game_instance_parameter_unique_key_per_instance IS 'Ensures unique parameter keys per game instance, but only for non-deleted records';

CREATE TABLE IF NOT EXISTS public.game_turn_sheet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_instance_id UUID,
    account_id UUID NOT NULL,
    turn_number INTEGER NOT NULL,
    sheet_type VARCHAR(50) NOT NULL,
    sheet_order INTEGER NOT NULL DEFAULT 1,
    sheet_data JSONB NOT NULL,
    is_completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMPTZ,
    scanned_data JSONB,
    scanned_at TIMESTAMPTZ,
    scanned_by UUID,
    scan_quality NUMERIC(3,2),
    processing_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT game_turn_sheet_processing_status_check CHECK (processing_status IN ('pending', 'processed', 'error')),
    CONSTRAINT game_turn_sheet_scan_quality_check CHECK (scan_quality >= 0 AND scan_quality <= 1),
    CONSTRAINT game_turn_sheet_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT game_turn_sheet_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id),
    CONSTRAINT game_turn_sheet_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id),
    CONSTRAINT game_turn_sheet_scanned_by_fkey FOREIGN KEY (scanned_by) REFERENCES public.account(id)
);

COMMENT ON TABLE public.game_turn_sheet IS 'Single turn sheet table for all game types, includes both sheet data and scanned results';

CREATE INDEX idx_game_turn_sheet_game_id ON public.game_turn_sheet(game_id);
CREATE INDEX idx_game_turn_sheet_game_instance_account ON public.game_turn_sheet(game_instance_id, account_id);
CREATE INDEX idx_game_turn_sheet_turn_number ON public.game_turn_sheet(turn_number);
CREATE INDEX idx_game_turn_sheet_sheet_type ON public.game_turn_sheet(sheet_type);
CREATE INDEX idx_game_turn_sheet_processing_status ON public.game_turn_sheet(processing_status);

-- ============================================================================
-- ADVENTURE GAME TABLES
-- ============================================================================

CREATE TABLE IF NOT EXISTS public.adventure_game_location (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_starting_location BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_location_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id)
);

COMMENT ON COLUMN public.adventure_game_location.is_starting_location IS 'Indicates if this location is a valid starting location for new players joining the game.';

CREATE INDEX idx_adventure_game_location_is_starting_location ON public.adventure_game_location(is_starting_location)
    WHERE is_starting_location = true;

CREATE TABLE IF NOT EXISTS public.adventure_game_location_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_instance_id UUID NOT NULL,
    adventure_game_location_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_location_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_location_instance_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id),
    CONSTRAINT adventure_game_location_instanc_adventure_game_location_id_fkey FOREIGN KEY (adventure_game_location_id) REFERENCES public.adventure_game_location(id)
);

COMMENT ON TABLE public.adventure_game_location_instance IS 'Tracks a specific instance of a location within a game instance.';
COMMENT ON COLUMN public.adventure_game_location_instance.id IS 'Unique identifier for the game location instance.';
COMMENT ON COLUMN public.adventure_game_location_instance.game_id IS 'The game this instance is based on.';
COMMENT ON COLUMN public.adventure_game_location_instance.game_instance_id IS 'The game instance this location instance belongs to.';
COMMENT ON COLUMN public.adventure_game_location_instance.adventure_game_location_id IS 'The base location this instance represents.';
COMMENT ON COLUMN public.adventure_game_location_instance.created_at IS 'When this location instance was created.';
COMMENT ON COLUMN public.adventure_game_location_instance.updated_at IS 'When this location instance was last updated.';
COMMENT ON COLUMN public.adventure_game_location_instance.deleted_at IS 'When this location instance was deleted (soft delete).';

CREATE TABLE IF NOT EXISTS public.adventure_game_location_link (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    from_adventure_game_location_id UUID NOT NULL,
    to_adventure_game_location_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_location_link_unique UNIQUE (from_adventure_game_location_id, to_adventure_game_location_id),
    CONSTRAINT adventure_game_location_link_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_location_link_from_adventure_game_location__fkey FOREIGN KEY (from_adventure_game_location_id) REFERENCES public.adventure_game_location(id),
    CONSTRAINT adventure_game_location_link_to_adventure_game_location_id_fkey FOREIGN KEY (to_adventure_game_location_id) REFERENCES public.adventure_game_location(id)
);

CREATE INDEX idx_adventure_game_location_link_from_adventure_game_location_i ON public.adventure_game_location_link(from_adventure_game_location_id);
CREATE INDEX idx_adventure_game_location_link_to_adventure_game_location_id ON public.adventure_game_location_link(to_adventure_game_location_id);

CREATE TABLE IF NOT EXISTS public.adventure_game_item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    item_category VARCHAR(50),
    equipment_slot VARCHAR(50),
    can_be_equipped BOOLEAN NOT NULL DEFAULT false,
    can_be_used BOOLEAN NOT NULL DEFAULT false,
    is_starting_item BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_item_item_category_check CHECK (
        item_category IN ('weapon', 'armor', 'clothing', 'jewelry', 'consumable', 'misc')
    ),
    CONSTRAINT adventure_game_item_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id)
);

COMMENT ON COLUMN public.adventure_game_item.item_category IS 'Item category: weapon, armor, clothing, jewelry, consumable, or misc';
COMMENT ON COLUMN public.adventure_game_item.equipment_slot IS 'Equipment slot this item occupies when equipped (e.g., weapon, armor_body, jewelry_ring). NULL for non-equippable items.';
COMMENT ON COLUMN public.adventure_game_item.is_starting_item IS 'If true, this item is automatically assigned to characters when they join the game';

CREATE TABLE IF NOT EXISTS public.adventure_game_item_placement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    adventure_game_item_id UUID NOT NULL,
    adventure_game_location_id UUID NOT NULL,
    initial_count INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_item_placement_game_id_adventure_game_item_i_key UNIQUE (game_id, adventure_game_item_id, adventure_game_location_id),
    CONSTRAINT adventure_game_item_placement_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_item_placement_adventure_game_item_id_fkey FOREIGN KEY (adventure_game_item_id) REFERENCES public.adventure_game_item(id),
    CONSTRAINT adventure_game_item_placement_adventure_game_location_id_fkey FOREIGN KEY (adventure_game_location_id) REFERENCES public.adventure_game_location(id)
);

COMMENT ON TABLE public.adventure_game_item_placement IS 'Configuration for placing items in specific locations when game instances are created.';
COMMENT ON COLUMN public.adventure_game_item_placement.id IS 'Unique identifier for the placement configuration.';
COMMENT ON COLUMN public.adventure_game_item_placement.game_id IS 'The game this placement configuration belongs to.';
COMMENT ON COLUMN public.adventure_game_item_placement.adventure_game_item_id IS 'The item type to be placed.';
COMMENT ON COLUMN public.adventure_game_item_placement.adventure_game_location_id IS 'The location where the item should be placed.';
COMMENT ON COLUMN public.adventure_game_item_placement.initial_count IS 'How many of this item should exist at this location when game instance is created.';
COMMENT ON COLUMN public.adventure_game_item_placement.created_at IS 'When this placement configuration was created.';
COMMENT ON COLUMN public.adventure_game_item_placement.updated_at IS 'When this placement configuration was last updated.';
COMMENT ON COLUMN public.adventure_game_item_placement.deleted_at IS 'When this placement configuration was deleted (soft delete).';

CREATE TABLE IF NOT EXISTS public.adventure_game_location_link_requirement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    adventure_game_location_link_id UUID NOT NULL,
    adventure_game_item_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_location_link_requirement_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_location_link__adventure_game_location_link_fkey FOREIGN KEY (adventure_game_location_link_id) REFERENCES public.adventure_game_location_link(id),
    CONSTRAINT adventure_game_location_link_requir_adventure_game_item_id_fkey FOREIGN KEY (adventure_game_item_id) REFERENCES public.adventure_game_item(id)
);

COMMENT ON TABLE public.adventure_game_location_link_requirement IS 'Specifies which items (and how many) are required to traverse a location link.';
COMMENT ON COLUMN public.adventure_game_location_link_requirement.id IS 'Unique identifier for the requirement.';
COMMENT ON COLUMN public.adventure_game_location_link_requirement.game_id IS 'The game this requirement belongs to.';
COMMENT ON COLUMN public.adventure_game_location_link_requirement.adventure_game_location_link_id IS 'The location link this requirement applies to.';
COMMENT ON COLUMN public.adventure_game_location_link_requirement.adventure_game_item_id IS 'The item type required to traverse the link.';
COMMENT ON COLUMN public.adventure_game_location_link_requirement.quantity IS 'How many of this item are required.';
COMMENT ON COLUMN public.adventure_game_location_link_requirement.created_at IS 'When this requirement was created.';
COMMENT ON COLUMN public.adventure_game_location_link_requirement.updated_at IS 'When this requirement was last updated.';
COMMENT ON COLUMN public.adventure_game_location_link_requirement.deleted_at IS 'When this requirement was deleted (soft delete).';

CREATE TABLE IF NOT EXISTS public.adventure_game_creature (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_creature_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id)
);

CREATE TABLE IF NOT EXISTS public.adventure_game_creature_placement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    adventure_game_creature_id UUID NOT NULL,
    adventure_game_location_id UUID NOT NULL,
    initial_count INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_creature_placement_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_creature_placeme_adventure_game_creature_id_fkey FOREIGN KEY (adventure_game_creature_id) REFERENCES public.adventure_game_creature(id),
    CONSTRAINT adventure_game_creature_placeme_adventure_game_location_id_fkey FOREIGN KEY (adventure_game_location_id) REFERENCES public.adventure_game_location(id)
);

CREATE INDEX idx_adventure_game_creature_placement_game_id ON public.adventure_game_creature_placement(game_id);
CREATE INDEX idx_adventure_game_creature_placement_creature_id ON public.adventure_game_creature_placement(adventure_game_creature_id);
CREATE INDEX idx_adventure_game_creature_placement_location_id ON public.adventure_game_creature_placement(adventure_game_location_id);
CREATE INDEX idx_adventure_game_creature_placement_deleted_at ON public.adventure_game_creature_placement(deleted_at);
CREATE UNIQUE INDEX idx_adventure_game_creature_placement_unique ON public.adventure_game_creature_placement(game_id, adventure_game_creature_id, adventure_game_location_id)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS public.adventure_game_character (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    account_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_character_name_not_empty CHECK (name IS NOT NULL AND name != ''),
    CONSTRAINT adventure_game_character_unique UNIQUE (game_id, account_id),
    CONSTRAINT adventure_game_character_name_unique UNIQUE (game_id, name),
    CONSTRAINT adventure_game_character_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_character_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id)
);

CREATE INDEX idx_adventure_game_character_game_id ON public.adventure_game_character(game_id);
CREATE INDEX idx_adventure_game_character_account_id ON public.adventure_game_character(account_id);

CREATE TABLE IF NOT EXISTS public.adventure_game_character_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_instance_id UUID NOT NULL,
    adventure_game_character_id UUID NOT NULL,
    adventure_game_location_instance_id UUID,
    health INTEGER NOT NULL DEFAULT 100,
    inventory_capacity INTEGER NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_character_instance_inventory_capacity_check CHECK (inventory_capacity > 0),
    CONSTRAINT adventure_game_character_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_character_instance_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id),
    CONSTRAINT adventure_game_character_insta_adventure_game_character_id_fkey FOREIGN KEY (adventure_game_character_id) REFERENCES public.adventure_game_character(id),
    CONSTRAINT adventure_game_character_inst_adventure_game_location_inst_fkey FOREIGN KEY (adventure_game_location_instance_id) REFERENCES public.adventure_game_location_instance(id)
);

COMMENT ON TABLE public.adventure_game_character_instance IS 'Tracks a specific instance of a character within a game instance.';
COMMENT ON COLUMN public.adventure_game_character_instance.id IS 'Unique identifier for the character instance.';
COMMENT ON COLUMN public.adventure_game_character_instance.game_id IS 'The game this character instance belongs to.';
COMMENT ON COLUMN public.adventure_game_character_instance.game_instance_id IS 'The game instance this character instance belongs to.';
COMMENT ON COLUMN public.adventure_game_character_instance.adventure_game_character_id IS 'The base character this instance represents.';
COMMENT ON COLUMN public.adventure_game_character_instance.adventure_game_location_instance_id IS 'The location this character instance is at (if applicable).';
COMMENT ON COLUMN public.adventure_game_character_instance.health IS 'The health of the character instance.';
COMMENT ON COLUMN public.adventure_game_character_instance.created_at IS 'When this character instance was created.';
COMMENT ON COLUMN public.adventure_game_character_instance.updated_at IS 'When this character instance was last updated.';
COMMENT ON COLUMN public.adventure_game_character_instance.deleted_at IS 'When this character instance was deleted (soft delete).';
COMMENT ON COLUMN public.adventure_game_character_instance.inventory_capacity IS 'Maximum number of items the character can carry in inventory. Default is 10.';

CREATE TABLE IF NOT EXISTS public.adventure_game_creature_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    adventure_game_creature_id UUID NOT NULL,
    game_instance_id UUID NOT NULL,
    adventure_game_location_instance_id UUID,
    health INTEGER NOT NULL DEFAULT 100,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_creature_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_creature_instance_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id),
    CONSTRAINT adventure_game_creature_instanc_adventure_game_creature_id_fkey FOREIGN KEY (adventure_game_creature_id) REFERENCES public.adventure_game_creature(id),
    CONSTRAINT adventure_game_creature_insta_adventure_game_location_inst_fkey FOREIGN KEY (adventure_game_location_instance_id) REFERENCES public.adventure_game_location_instance(id)
);

CREATE TABLE IF NOT EXISTS public.adventure_game_item_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_instance_id UUID NOT NULL,
    adventure_game_item_id UUID NOT NULL,
    adventure_game_location_instance_id UUID,
    adventure_game_character_instance_id UUID,
    adventure_game_creature_instance_id UUID,
    equipment_slot VARCHAR(50),
    is_equipped BOOLEAN NOT NULL DEFAULT false,
    is_used BOOLEAN NOT NULL DEFAULT false,
    uses_remaining INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT only_one_owner CHECK (
        (adventure_game_location_instance_id IS NOT NULL)::integer +
        (adventure_game_character_instance_id IS NOT NULL)::integer +
        (adventure_game_creature_instance_id IS NOT NULL)::integer = 1
    ),
    CONSTRAINT adventure_game_item_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_item_instance_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id),
    CONSTRAINT adventure_game_item_instance_adventure_game_item_id_fkey FOREIGN KEY (adventure_game_item_id) REFERENCES public.adventure_game_item(id),
    CONSTRAINT adventure_game_item_instance_adventure_game_location_insta_fkey FOREIGN KEY (adventure_game_location_instance_id) REFERENCES public.adventure_game_location_instance(id),
    CONSTRAINT adventure_game_item_instance_adventure_game_character_inst_fkey FOREIGN KEY (adventure_game_character_instance_id) REFERENCES public.adventure_game_character_instance(id),
    CONSTRAINT adventure_game_item_instance_adventure_game_creature_insta_fkey FOREIGN KEY (adventure_game_creature_instance_id) REFERENCES public.adventure_game_creature_instance(id)
);

COMMENT ON TABLE public.adventure_game_item_instance IS 'Tracks a specific instance of a game item, which may be at a location, in a character inventory, or in a creature inventory.';
COMMENT ON COLUMN public.adventure_game_item_instance.id IS 'Unique identifier for the item instance.';
COMMENT ON COLUMN public.adventure_game_item_instance.game_id IS 'The game this item instance belongs to.';
COMMENT ON COLUMN public.adventure_game_item_instance.adventure_game_item_id IS 'The item type this instance represents.';
COMMENT ON COLUMN public.adventure_game_item_instance.game_instance_id IS 'The game instance this item instance belongs to.';
COMMENT ON COLUMN public.adventure_game_item_instance.adventure_game_location_instance_id IS 'If set, the item is at this location.';
COMMENT ON COLUMN public.adventure_game_item_instance.adventure_game_character_instance_id IS 'If set, the item is in this character''s inventory.';
COMMENT ON COLUMN public.adventure_game_item_instance.adventure_game_creature_instance_id IS 'If set, the item is in this creature''s inventory.';
COMMENT ON COLUMN public.adventure_game_item_instance.is_equipped IS 'Whether this item instance is currently equipped.';
COMMENT ON COLUMN public.adventure_game_item_instance.is_used IS 'Whether this item instance has been used (if applicable).';
COMMENT ON COLUMN public.adventure_game_item_instance.uses_remaining IS 'Number of uses left for this instance (if applicable).';
COMMENT ON COLUMN public.adventure_game_item_instance.created_at IS 'When this instance was created.';
COMMENT ON COLUMN public.adventure_game_item_instance.updated_at IS 'When this instance was last updated.';
COMMENT ON COLUMN public.adventure_game_item_instance.deleted_at IS 'When this instance was deleted (soft delete).';
COMMENT ON COLUMN public.adventure_game_item_instance.equipment_slot IS 'Equipment slot this item instance is equipped in (e.g., weapon, armor_body, jewelry_ring). NULL when not equipped or not equippable.';

CREATE TABLE IF NOT EXISTS public.adventure_game_turn_sheet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    adventure_game_character_instance_id UUID NOT NULL,
    game_turn_sheet_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_turn_sheet_adventure_game_character_instance_key UNIQUE (adventure_game_character_instance_id, game_turn_sheet_id),
    CONSTRAINT adventure_game_turn_sheet_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_turn_sheet_adventure_game_character_instanc_fkey FOREIGN KEY (adventure_game_character_instance_id) REFERENCES public.adventure_game_character_instance(id),
    CONSTRAINT adventure_game_turn_sheet_game_turn_sheet_id_fkey FOREIGN KEY (game_turn_sheet_id) REFERENCES public.game_turn_sheet(id)
);

COMMENT ON TABLE public.adventure_game_turn_sheet IS 'Mapping table linking adventure game entities (characters, locations, etc.) to turn sheets';

CREATE INDEX idx_adventure_game_turn_sheet_game_id ON public.adventure_game_turn_sheet(game_id);
CREATE INDEX idx_adventure_game_turn_sheet_character_instance ON public.adventure_game_turn_sheet(adventure_game_character_instance_id);
CREATE INDEX idx_adventure_game_turn_sheet_turn_sheet ON public.adventure_game_turn_sheet(game_turn_sheet_id);

