-- Adventure Game Location Object State Migration
-- Promotes object states from freeform strings to first-class database records.

BEGIN;

-- New table for object state definitions
CREATE TABLE public.adventure_game_location_object_state (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    adventure_game_location_object_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_location_object_state_name_not_empty CHECK (name IS NOT NULL AND name != ''),
    CONSTRAINT adventure_game_location_object_state_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_location_object_state_object_id_fkey FOREIGN KEY (adventure_game_location_object_id) REFERENCES public.adventure_game_location_object(id),
    CONSTRAINT adventure_game_location_object_state_unique_name UNIQUE (adventure_game_location_object_id, name, deleted_at)
);
CREATE INDEX idx_adventure_game_location_object_state_game_id ON public.adventure_game_location_object_state(game_id);
CREATE INDEX idx_adventure_game_location_object_state_object_id ON public.adventure_game_location_object_state(adventure_game_location_object_id);
COMMENT ON TABLE public.adventure_game_location_object_state IS 'Defines the discrete named states an object can occupy, scoped per object definition';

-- Alter adventure_game_location_object: replace initial_state text with FK
ALTER TABLE public.adventure_game_location_object
    DROP COLUMN initial_state;
ALTER TABLE public.adventure_game_location_object
    ADD COLUMN initial_adventure_game_location_object_state_id UUID;
ALTER TABLE public.adventure_game_location_object
    ADD CONSTRAINT adventure_game_location_object_initial_state_id_fkey
    FOREIGN KEY (initial_adventure_game_location_object_state_id)
    REFERENCES public.adventure_game_location_object_state(id);

-- Alter adventure_game_location_object_effect: replace required_state and result_state with FKs
ALTER TABLE public.adventure_game_location_object_effect
    DROP COLUMN required_state;
ALTER TABLE public.adventure_game_location_object_effect
    ADD COLUMN required_adventure_game_location_object_state_id UUID;
ALTER TABLE public.adventure_game_location_object_effect
    ADD CONSTRAINT adventure_game_location_object_effect_required_state_id_fkey
    FOREIGN KEY (required_adventure_game_location_object_state_id)
    REFERENCES public.adventure_game_location_object_state(id);

ALTER TABLE public.adventure_game_location_object_effect
    DROP COLUMN result_state;
ALTER TABLE public.adventure_game_location_object_effect
    ADD COLUMN result_adventure_game_location_object_state_id UUID;
ALTER TABLE public.adventure_game_location_object_effect
    ADD CONSTRAINT adventure_game_location_object_effect_result_state_id_fkey
    FOREIGN KEY (result_adventure_game_location_object_state_id)
    REFERENCES public.adventure_game_location_object_state(id);

-- Alter adventure_game_location_object_instance: replace current_state text with FK
ALTER TABLE public.adventure_game_location_object_instance
    DROP COLUMN current_state;
ALTER TABLE public.adventure_game_location_object_instance
    ADD COLUMN current_adventure_game_location_object_state_id UUID;
ALTER TABLE public.adventure_game_location_object_instance
    ADD CONSTRAINT adventure_game_location_object_instance_current_state_id_fkey
    FOREIGN KEY (current_adventure_game_location_object_state_id)
    REFERENCES public.adventure_game_location_object_state(id);

COMMIT;
