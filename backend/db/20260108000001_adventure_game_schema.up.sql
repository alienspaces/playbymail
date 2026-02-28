-- Adventure Game Schema Migration
-- Generated: 2026-01-08

BEGIN;

-- ============================================================================
-- ADVENTURE GAME TABLES
-- ============================================================================

CREATE TABLE public.adventure_game_location (
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
CREATE INDEX idx_adventure_game_location_is_starting_location ON public.adventure_game_location(is_starting_location) WHERE is_starting_location = true;

CREATE TABLE public.adventure_game_location_instance (
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

CREATE TABLE public.adventure_game_location_link (
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

CREATE TABLE public.adventure_game_item (
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

CREATE TABLE public.adventure_game_item_placement (
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

CREATE TABLE public.adventure_game_location_link_requirement (
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

CREATE TABLE public.adventure_game_creature (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_creature_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id)
);

CREATE TABLE public.adventure_game_creature_placement (
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
CREATE UNIQUE INDEX idx_adventure_game_creature_placement_unique ON public.adventure_game_creature_placement(game_id, adventure_game_creature_id, adventure_game_location_id) WHERE deleted_at IS NULL;

CREATE TABLE public.adventure_game_character (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    account_id UUID NOT NULL,
    account_user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_character_name_not_empty CHECK (name IS NOT NULL AND name != ''),
    CONSTRAINT adventure_game_character_unique UNIQUE (game_id, account_id),
    CONSTRAINT adventure_game_character_name_unique UNIQUE (game_id, name),
    CONSTRAINT adventure_game_character_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_character_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id),
    CONSTRAINT adventure_game_character_account_user_id_fkey FOREIGN KEY (account_user_id) REFERENCES public.account_user(id)
);
CREATE INDEX idx_adventure_game_character_game_id ON public.adventure_game_character(game_id);
CREATE INDEX idx_adventure_game_character_account_id ON public.adventure_game_character(account_id);
CREATE INDEX idx_adventure_game_character_account_user_id ON public.adventure_game_character(account_user_id);

CREATE TABLE public.adventure_game_character_instance (
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

CREATE TABLE public.adventure_game_creature_instance (
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

CREATE TABLE public.adventure_game_item_instance (
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

CREATE TABLE public.adventure_game_turn_sheet (
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
CREATE INDEX idx_adventure_game_turn_sheet_game_id ON public.adventure_game_turn_sheet(game_id);
CREATE INDEX idx_adventure_game_turn_sheet_character_instance ON public.adventure_game_turn_sheet(adventure_game_character_instance_id);
CREATE INDEX idx_adventure_game_turn_sheet_turn_sheet ON public.adventure_game_turn_sheet(game_turn_sheet_id);

COMMIT;
