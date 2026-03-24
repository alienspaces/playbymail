-- Mecha Schema Migration

BEGIN;

-- ============================================================================
-- MECHA DESIGN-TIME TABLES
-- ============================================================================

CREATE TABLE public.mecha_chassis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    chassis_class VARCHAR(20) NOT NULL DEFAULT 'medium',
    armor_points INTEGER NOT NULL DEFAULT 100,
    structure_points INTEGER NOT NULL DEFAULT 50,
    heat_capacity INTEGER NOT NULL DEFAULT 30,
    speed INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_chassis_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_chassis_class_check CHECK (chassis_class IN ('light', 'medium', 'heavy', 'assault'))
);
CREATE INDEX idx_mecha_chassis_game_id ON public.mecha_chassis(game_id);
COMMENT ON TABLE public.mecha_chassis IS 'Mech type definitions including combat stats.';

CREATE TABLE public.mecha_weapon (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    damage INTEGER NOT NULL DEFAULT 5,
    heat_cost INTEGER NOT NULL DEFAULT 3,
    range_band VARCHAR(20) NOT NULL DEFAULT 'medium',
    mount_size VARCHAR(20) NOT NULL DEFAULT 'medium',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_weapon_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_weapon_range_band_check CHECK (range_band IN ('short', 'medium', 'long')),
    CONSTRAINT mecha_weapon_mount_size_check CHECK (mount_size IN ('small', 'medium', 'large'))
);
CREATE INDEX idx_mecha_weapon_game_id ON public.mecha_weapon(game_id);
COMMENT ON TABLE public.mecha_weapon IS 'Weapon definitions including damage and heat cost.';

CREATE TABLE public.mecha_sector (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    terrain_type VARCHAR(50) NOT NULL DEFAULT 'open',
    elevation INTEGER NOT NULL DEFAULT 0,
    is_starting_sector BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_sector_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_sector_terrain_type_check CHECK (terrain_type IN ('open', 'urban', 'forest', 'rough', 'water'))
);
CREATE INDEX idx_mecha_sector_game_id ON public.mecha_sector(game_id);
CREATE INDEX idx_mecha_sector_is_starting ON public.mecha_sector(is_starting_sector) WHERE is_starting_sector = true;
COMMENT ON TABLE public.mecha_sector IS 'Named battlefield zones (area-based map model).';

CREATE TABLE public.mecha_sector_link (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    from_mecha_sector_id UUID NOT NULL,
    to_mecha_sector_id UUID NOT NULL,
    cover_modifier INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_sector_link_unique UNIQUE (from_mecha_sector_id, to_mecha_sector_id),
    CONSTRAINT mecha_sector_link_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_sector_link_from_sector_fkey FOREIGN KEY (from_mecha_sector_id) REFERENCES public.mecha_sector(id),
    CONSTRAINT mecha_sector_link_to_sector_fkey FOREIGN KEY (to_mecha_sector_id) REFERENCES public.mecha_sector(id)
);
CREATE INDEX idx_mecha_sector_link_from_sector ON public.mecha_sector_link(from_mecha_sector_id);
CREATE INDEX idx_mecha_sector_link_to_sector ON public.mecha_sector_link(to_mecha_sector_id);
COMMENT ON TABLE public.mecha_sector_link IS 'Directed connections between battlefield sectors.';

CREATE TABLE public.mecha_lance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    account_id UUID NOT NULL,
    account_user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_lance_unique UNIQUE (game_id, account_id),
    CONSTRAINT mecha_lance_name_unique UNIQUE (game_id, name),
    CONSTRAINT mecha_lance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_lance_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id),
    CONSTRAINT mecha_lance_account_user_id_fkey FOREIGN KEY (account_user_id) REFERENCES public.account_user(id)
);
CREATE INDEX idx_mecha_lance_game_id ON public.mecha_lance(game_id);
CREATE INDEX idx_mecha_lance_account_id ON public.mecha_lance(account_id);
COMMENT ON TABLE public.mecha_lance IS 'Player lance slot: one lance per player in the game design.';

CREATE TABLE public.mecha_lance_mech (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    mecha_lance_id UUID NOT NULL,
    mecha_chassis_id UUID NOT NULL,
    callsign VARCHAR(50) NOT NULL,
    weapon_config JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_lance_mech_callsign_unique UNIQUE (mecha_lance_id, callsign),
    CONSTRAINT mecha_lance_mech_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_lance_mech_lance_id_fkey FOREIGN KEY (mecha_lance_id) REFERENCES public.mecha_lance(id),
    CONSTRAINT mecha_lance_mech_chassis_id_fkey FOREIGN KEY (mecha_chassis_id) REFERENCES public.mecha_chassis(id)
);
CREATE INDEX idx_mecha_lance_mech_lance_id ON public.mecha_lance_mech(mecha_lance_id);
CREATE INDEX idx_mecha_lance_mech_chassis_id ON public.mecha_lance_mech(mecha_chassis_id);
COMMENT ON TABLE public.mecha_lance_mech IS 'Mech assignments within a lance. weapon_config is a JSONB array of weapon_id + slot_location pairs.';

-- ============================================================================
-- MECHA RUNTIME TABLES
-- ============================================================================

CREATE TABLE public.mecha_sector_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_instance_id UUID NOT NULL,
    mecha_sector_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_sector_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_sector_instance_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id),
    CONSTRAINT mecha_sector_instance_sector_id_fkey FOREIGN KEY (mecha_sector_id) REFERENCES public.mecha_sector(id)
);
CREATE INDEX idx_mecha_sector_instance_game_instance ON public.mecha_sector_instance(game_instance_id);
CREATE INDEX idx_mecha_sector_instance_sector_id ON public.mecha_sector_instance(mecha_sector_id);
COMMENT ON TABLE public.mecha_sector_instance IS 'Runtime sector state per game instance.';

CREATE TABLE public.mecha_lance_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_instance_id UUID NOT NULL,
    mecha_lance_id UUID NOT NULL,
    game_subscription_instance_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_lance_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_lance_instance_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id),
    CONSTRAINT mecha_lance_instance_lance_id_fkey FOREIGN KEY (mecha_lance_id) REFERENCES public.mecha_lance(id),
    CONSTRAINT mecha_lance_instance_subscription_instance_id_fkey FOREIGN KEY (game_subscription_instance_id) REFERENCES public.game_subscription_instance(id)
);
CREATE INDEX idx_mecha_lance_instance_game_instance ON public.mecha_lance_instance(game_instance_id);
CREATE INDEX idx_mecha_lance_instance_lance_id ON public.mecha_lance_instance(mecha_lance_id);
COMMENT ON TABLE public.mecha_lance_instance IS 'Runtime lance state linking a player subscription to their lance.';

CREATE TABLE public.mecha_mech_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_instance_id UUID NOT NULL,
    mecha_lance_instance_id UUID NOT NULL,
    mecha_sector_instance_id UUID NOT NULL,
    mecha_chassis_id UUID NOT NULL,
    callsign VARCHAR(50) NOT NULL,
    current_armor INTEGER NOT NULL DEFAULT 100,
    current_structure INTEGER NOT NULL DEFAULT 50,
    current_heat INTEGER NOT NULL DEFAULT 0,
    pilot_skill INTEGER NOT NULL DEFAULT 4,
    status VARCHAR(20) NOT NULL DEFAULT 'operational',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_mech_instance_status_check CHECK (status IN ('operational', 'damaged', 'destroyed', 'shutdown')),
    CONSTRAINT mecha_mech_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_mech_instance_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id),
    CONSTRAINT mecha_mech_instance_lance_instance_id_fkey FOREIGN KEY (mecha_lance_instance_id) REFERENCES public.mecha_lance_instance(id),
    CONSTRAINT mecha_mech_instance_sector_instance_id_fkey FOREIGN KEY (mecha_sector_instance_id) REFERENCES public.mecha_sector_instance(id),
    CONSTRAINT mecha_mech_instance_chassis_id_fkey FOREIGN KEY (mecha_chassis_id) REFERENCES public.mecha_chassis(id)
);
CREATE INDEX idx_mecha_mech_instance_game_instance ON public.mecha_mech_instance(game_instance_id);
CREATE INDEX idx_mecha_mech_instance_lance_instance ON public.mecha_mech_instance(mecha_lance_instance_id);
CREATE INDEX idx_mecha_mech_instance_sector_instance ON public.mecha_mech_instance(mecha_sector_instance_id);
COMMENT ON TABLE public.mecha_mech_instance IS 'Runtime mech state including current combat stats and position.';

CREATE TABLE public.mecha_turn_sheet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    mecha_lance_instance_id UUID NOT NULL,
    game_turn_sheet_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_turn_sheet_unique UNIQUE (mecha_lance_instance_id, game_turn_sheet_id),
    CONSTRAINT mecha_turn_sheet_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_turn_sheet_lance_instance_id_fkey FOREIGN KEY (mecha_lance_instance_id) REFERENCES public.mecha_lance_instance(id),
    CONSTRAINT mecha_turn_sheet_game_turn_sheet_id_fkey FOREIGN KEY (game_turn_sheet_id) REFERENCES public.game_turn_sheet(id)
);
CREATE INDEX idx_mecha_turn_sheet_game_id ON public.mecha_turn_sheet(game_id);
CREATE INDEX idx_mecha_turn_sheet_lance_instance ON public.mecha_turn_sheet(mecha_lance_instance_id);
CREATE INDEX idx_mecha_turn_sheet_game_turn_sheet ON public.mecha_turn_sheet(game_turn_sheet_id);
COMMENT ON TABLE public.mecha_turn_sheet IS 'Bridge between mecha lance instances and core game turn sheets.';

COMMIT;
