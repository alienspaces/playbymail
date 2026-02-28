-- Consolidated Schema Migration
-- Generated: 2026-01-08

BEGIN;

-- ============================================================================
-- ACCOUNT TABLES
-- ============================================================================

-- 1. Account
CREATE TABLE public.account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT account_status_check CHECK (status IN ('active', 'disabled'))
);

COMMENT ON TABLE public.account IS 'Account that owns resources.';
COMMENT ON COLUMN public.account.name IS 'Display name for the account.';

-- 2. Account User
CREATE TABLE public.account_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    verification_token TEXT,
    verification_token_expires_at TIMESTAMPTZ,
    session_token TEXT,
    session_token_expires_at TIMESTAMPTZ,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT account_status_check CHECK (status IN ('pending_approval', 'active', 'disabled')),
    CONSTRAINT account_user_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id)
);

COMMENT ON TABLE public.account_user IS 'Individual account users authenticated in the system.';
COMMENT ON COLUMN public.account_user.account_id IS 'The Account this user belongs to.';
COMMENT ON COLUMN public.account_user.status IS 'Current approval status of the account user (pending_approval, active, disabled).';

-- 3. Account User Contact
CREATE TABLE public.account_user_contact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    postal_address_line1 VARCHAR(255) NOT NULL,
    postal_address_line2 VARCHAR(255),
    state_province VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT account_user_contact_account_user_id_fkey FOREIGN KEY (account_user_id) REFERENCES public.account_user(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.account_user_contact IS 'Contact information for accounts, including name and postal address. Collected from join game turn sheets.';
COMMENT ON COLUMN public.account_user_contact.id IS 'Unique identifier for the contact record.';
COMMENT ON COLUMN public.account_user_contact.account_user_id IS 'The account user this contact belongs to.';

CREATE INDEX idx_account_user_contact_account_user_id ON public.account_user_contact(account_user_id);

-- 4. Account Subscription
CREATE TABLE public.account_subscription (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID,
    account_user_id UUID,
    subscription_type VARCHAR(50) NOT NULL,
    subscription_period VARCHAR(32) NOT NULL DEFAULT 'eternal',
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    auto_renew BOOLEAN NOT NULL DEFAULT true,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT account_subscription_subscription_type_check CHECK (subscription_type IN ('basic_game_designer', 'professional_game_designer', 'basic_manager', 'professional_manager', 'basic_player', 'professional_player')),
    CONSTRAINT account_subscription_subscription_period_check CHECK (subscription_period IN ('month', 'year', 'eternal')),
    CONSTRAINT account_subscription_status_check CHECK (status IN ('active', 'expired')),
    CONSTRAINT account_subscription_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id) ON DELETE CASCADE,
    CONSTRAINT account_subscription_account_user_id_fkey FOREIGN KEY (account_user_id) REFERENCES public.account_user(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.account_subscription IS 'Account level subscriptions for game design and game management.';
CREATE INDEX idx_account_subscription_account_id ON public.account_subscription(account_id);
CREATE INDEX idx_account_subscription_status ON public.account_subscription(status) WHERE status = 'active';

-- ============================================================================
-- GAME TABLES
-- ============================================================================

CREATE TABLE public.game (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    game_type VARCHAR(50) NOT NULL,
    description VARCHAR(512) NOT NULL,
    turn_duration_hours INTEGER NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT game_name_check CHECK (name != ''),
    CONSTRAINT game_type_check CHECK (game_type = 'adventure'),
    CONSTRAINT turn_duration_hours_check CHECK (turn_duration_hours > 0),
    CONSTRAINT game_status_check CHECK (status IN ('draft', 'published'))
);

COMMENT ON COLUMN public.game.description IS 'Game description that appears on the join game turn sheet.';
COMMENT ON COLUMN public.game.status IS 'Game status (draft, published).';
CREATE INDEX idx_game_status ON public.game(status);

CREATE TABLE public.game_image (
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
CREATE INDEX idx_game_image_game_record ON public.game_image(game_id, record_id);
CREATE INDEX idx_game_image_type ON public.game_image(type);
CREATE INDEX idx_game_image_turn_sheet_type ON public.game_image(turn_sheet_type);
CREATE UNIQUE INDEX game_image_unique_turn_sheet_background ON public.game_image(game_id, record_id, type, turn_sheet_type)
    WHERE type = 'turn_sheet_background' AND turn_sheet_type IS NOT NULL;
CREATE UNIQUE INDEX game_image_unique_other_types ON public.game_image(game_id, record_id, type)
    WHERE (type != 'turn_sheet_background' OR turn_sheet_type IS NULL);

CREATE TABLE public.game_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
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
    CONSTRAINT game_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id)
);

COMMENT ON TABLE public.game_instance IS 'Tracks a specific play session or instance of a game.';
CREATE INDEX idx_game_instance_closed_testing_join_game_key ON public.game_instance(closed_testing_join_game_key)
    WHERE closed_testing_join_game_key IS NOT NULL;

CREATE TABLE public.game_subscription (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    account_id UUID NOT NULL,
    account_user_id UUID,
    account_contact_id UUID,
    subscription_type VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    instance_limit INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT game_subscription_status_check CHECK (status IN ('pending_approval', 'active', 'revoked')),
    CONSTRAINT game_subscription_subscription_type_check CHECK (subscription_type IN ('player', 'manager', 'designer')),
    CONSTRAINT game_subscription_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT game_subscription_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id) ON DELETE CASCADE,
    CONSTRAINT game_subscription_account_user_id_fkey FOREIGN KEY (account_user_id) REFERENCES public.account_user(id) ON DELETE CASCADE,
    CONSTRAINT game_subscription_account_contact_id_fkey FOREIGN KEY (account_contact_id) REFERENCES public.account_user_contact(id) ON DELETE SET NULL
);

COMMENT ON TABLE public.game_subscription IS 'Tracks account subscriptions to a game, including Player, Manager, and Collaborator roles.';
COMMENT ON COLUMN public.game_subscription.instance_limit IS 'Maximum number of game instances this subscription can manage/play. NULL means unlimited.';

CREATE TABLE public.game_subscription_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    game_subscription_id UUID NOT NULL,
    game_instance_id UUID NOT NULL,
    turn_sheet_token VARCHAR(255),
    turn_sheet_token_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT game_subscription_instance_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id) ON DELETE CASCADE,
    CONSTRAINT game_subscription_instance_game_subscription_id_fkey FOREIGN KEY (game_subscription_id) REFERENCES public.game_subscription(id),
    CONSTRAINT game_subscription_instance_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id)
);

COMMENT ON TABLE public.game_subscription_instance IS 'Links game subscriptions to game instances in a many-to-many relationship.';
CREATE UNIQUE INDEX idx_game_subscription_instance_unique ON public.game_subscription_instance(game_subscription_id, game_instance_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_game_subscription_instance_account_id ON public.game_subscription_instance(account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_game_subscription_instance_game_subscription_id ON public.game_subscription_instance(game_subscription_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_game_subscription_instance_game_instance_id ON public.game_subscription_instance(game_instance_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_game_subscription_instance_turn_sheet_token ON public.game_subscription_instance(turn_sheet_token) WHERE deleted_at IS NULL AND turn_sheet_token IS NOT NULL;

CREATE TABLE public.game_instance_parameter (
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
CREATE INDEX idx_game_instance_parameter_game_instance_id ON public.game_instance_parameter(game_instance_id);
CREATE INDEX idx_game_instance_parameter_parameter_key ON public.game_instance_parameter(parameter_key);
CREATE UNIQUE INDEX idx_game_instance_parameter_unique_key_per_instance ON public.game_instance_parameter(game_instance_id, parameter_key) WHERE deleted_at IS NULL;

CREATE TABLE public.game_turn_sheet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_instance_id UUID,
    account_id UUID NOT NULL, -- Owner
    account_user_id UUID NOT NULL, -- Player
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
    CONSTRAINT game_turn_sheet_account_user_id_fkey FOREIGN KEY (account_user_id) REFERENCES public.account_user(id),
    CONSTRAINT game_turn_sheet_scanned_by_fkey FOREIGN KEY (scanned_by) REFERENCES public.account_user(id)
);

COMMENT ON TABLE public.game_turn_sheet IS 'Single turn sheet table for all game types, includes both sheet data and scanned results';
CREATE INDEX idx_game_turn_sheet_game_id ON public.game_turn_sheet(game_id);
CREATE INDEX idx_game_turn_sheet_game_instance_account ON public.game_turn_sheet(game_instance_id, account_id);
CREATE INDEX idx_game_turn_sheet_account ON public.game_turn_sheet(account_id);
CREATE INDEX idx_game_turn_sheet_account_user ON public.game_turn_sheet(account_user_id);
CREATE INDEX idx_game_turn_sheet_turn_number ON public.game_turn_sheet(turn_number);
CREATE INDEX idx_game_turn_sheet_sheet_type ON public.game_turn_sheet(sheet_type);
CREATE INDEX idx_game_turn_sheet_processing_status ON public.game_turn_sheet(processing_status);



-- ============================================================================
-- VIEWS
-- ============================================================================

CREATE OR REPLACE VIEW public.game_subscription_view AS
SELECT 
    gs.id,
    gs.game_id,
    gs.account_id,
    gs.account_user_id,
    gs.account_contact_id,
    gs.subscription_type,
    gs.status,
    gs.instance_limit,
    gs.created_at,
    gs.updated_at,
    gs.deleted_at,
    COALESCE(
        array_agg(gsi.game_instance_id) FILTER (WHERE gsi.game_instance_id IS NOT NULL AND gsi.deleted_at IS NULL),
        ARRAY[]::UUID[]
    ) AS game_instance_ids
FROM public.game_subscription gs
LEFT JOIN public.game_subscription_instance gsi ON gs.id = gsi.game_subscription_id
WHERE gs.deleted_at IS NULL
GROUP BY 
    gs.id,
    gs.game_id,
    gs.account_id,
    gs.account_user_id,
    gs.account_contact_id,
    gs.subscription_type,
    gs.status,
    gs.instance_limit,
    gs.created_at,
    gs.updated_at,
    gs.deleted_at;

COMMENT ON VIEW public.game_subscription_view IS 'View of game_subscription records with aggregated game_instance_ids array. Read-only view for API queries.';

COMMIT;
