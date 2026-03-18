BEGIN;

-- Add locked_description to location links (shown when player lacks traverse requirements)
ALTER TABLE public.adventure_game_location_link
    ADD COLUMN locked_description TEXT;

-- Make adventure_game_item_id nullable (was NOT NULL) to allow creature-based requirements
ALTER TABLE public.adventure_game_location_link_requirement
    ALTER COLUMN adventure_game_item_id DROP NOT NULL;

-- Add creature reference
ALTER TABLE public.adventure_game_location_link_requirement
    ADD COLUMN adventure_game_creature_id UUID
        REFERENCES public.adventure_game_creature(id);

-- Add purpose: 'traverse' = must have to walk through, 'visible' = must have to see the link at all
ALTER TABLE public.adventure_game_location_link_requirement
    ADD COLUMN purpose VARCHAR(20) NOT NULL DEFAULT 'traverse';

-- Add condition: how the requirement is evaluated
-- Item conditions:    in_inventory, equipped
-- Creature conditions: dead_at_location, none_alive_at_location, none_alive_in_game
ALTER TABLE public.adventure_game_location_link_requirement
    ADD COLUMN condition VARCHAR(50) NOT NULL DEFAULT 'in_inventory';

-- Backfill existing rows with explicit values matching the previous implicit behaviour
UPDATE public.adventure_game_location_link_requirement
    SET purpose = 'traverse', condition = 'in_inventory'
    WHERE deleted_at IS NULL;

-- Exactly one of item or creature must be set
ALTER TABLE public.adventure_game_location_link_requirement
    ADD CONSTRAINT adventure_game_location_link_requirement_one_target CHECK (
        (adventure_game_item_id IS NOT NULL)::integer +
        (adventure_game_creature_id IS NOT NULL)::integer = 1
    );

-- purpose must be a known value
ALTER TABLE public.adventure_game_location_link_requirement
    ADD CONSTRAINT adventure_game_location_link_requirement_purpose_check CHECK (
        purpose IN ('traverse', 'visible')
    );

-- condition must be a known value
ALTER TABLE public.adventure_game_location_link_requirement
    ADD CONSTRAINT adventure_game_location_link_requirement_condition_check CHECK (
        condition IN ('in_inventory', 'equipped', 'dead_at_location', 'none_alive_at_location', 'none_alive_in_game')
    );

-- Item conditions only allowed when item_id is set; creature conditions only when creature_id is set
ALTER TABLE public.adventure_game_location_link_requirement
    ADD CONSTRAINT adventure_game_location_link_requirement_condition_target_check CHECK (
        (adventure_game_item_id IS NOT NULL AND condition IN ('in_inventory', 'equipped'))
        OR
        (adventure_game_creature_id IS NOT NULL AND condition IN ('dead_at_location', 'none_alive_at_location', 'none_alive_in_game'))
    );

CREATE INDEX idx_adventure_game_location_link_requirement_creature_id
    ON public.adventure_game_location_link_requirement(adventure_game_creature_id)
    WHERE adventure_game_creature_id IS NOT NULL;

CREATE INDEX idx_adventure_game_location_link_requirement_purpose
    ON public.adventure_game_location_link_requirement(purpose);

COMMIT;
