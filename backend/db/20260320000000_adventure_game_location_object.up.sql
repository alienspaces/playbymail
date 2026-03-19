-- Adventure Game Location Object Schema Migration

BEGIN;

CREATE TABLE public.adventure_game_location_object (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    adventure_game_location_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    initial_state VARCHAR(50) NOT NULL DEFAULT 'intact',
    is_hidden BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_location_object_name_not_empty CHECK (name IS NOT NULL AND name != ''),
    CONSTRAINT adventure_game_location_object_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_location_object_adventure_game_location_id_fkey FOREIGN KEY (adventure_game_location_id) REFERENCES public.adventure_game_location(id)
);
CREATE INDEX idx_adventure_game_location_object_game_id ON public.adventure_game_location_object(game_id);
CREATE INDEX idx_adventure_game_location_object_location_id ON public.adventure_game_location_object(adventure_game_location_id);

-- adventure_game_location_object_effect defines what happens when a player performs an action on an object.
-- Multiple rows for the same (object, action_type, required_state) are allowed — all fire atomically.
CREATE TABLE public.adventure_game_location_object_effect (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    adventure_game_location_object_id UUID NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    required_state VARCHAR(50),
    required_adventure_game_item_id UUID,
    result_description TEXT NOT NULL,
    effect_type VARCHAR(50) NOT NULL,
    result_state VARCHAR(50),
    result_adventure_game_item_id UUID,
    result_adventure_game_location_link_id UUID,
    result_adventure_game_creature_id UUID,
    result_adventure_game_location_object_id UUID,
    result_adventure_game_location_id UUID,
    result_value_min INTEGER,
    result_value_max INTEGER,
    is_repeatable BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_location_object_effect_action_type_check CHECK (
        action_type IN (
            'inspect','touch','open','close','lock','unlock','search',
            'break','push','pull','move','burn','read','take',
            'listen','insert','pour','disarm','climb','use'
        )
    ),
    CONSTRAINT adventure_game_location_object_effect_effect_type_check CHECK (
        effect_type IN (
            'info','change_state','change_object_state','give_item','remove_item',
            'open_link','close_link','reveal_object','hide_object',
            'damage','heal','summon_creature','teleport','nothing','remove_object'
        )
    ),
    CONSTRAINT adventure_game_location_object_effect_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_location_object_effect_object_id_fkey FOREIGN KEY (adventure_game_location_object_id) REFERENCES public.adventure_game_location_object(id),
    CONSTRAINT adventure_game_location_object_effect_item_id_fkey FOREIGN KEY (required_adventure_game_item_id) REFERENCES public.adventure_game_item(id),
    CONSTRAINT adventure_game_location_object_effect_result_item_fkey FOREIGN KEY (result_adventure_game_item_id) REFERENCES public.adventure_game_item(id),
    CONSTRAINT adventure_game_location_object_effect_result_link_fkey FOREIGN KEY (result_adventure_game_location_link_id) REFERENCES public.adventure_game_location_link(id),
    CONSTRAINT adventure_game_location_object_effect_result_creature_fkey FOREIGN KEY (result_adventure_game_creature_id) REFERENCES public.adventure_game_creature(id),
    CONSTRAINT adventure_game_location_object_effect_result_object_fkey FOREIGN KEY (result_adventure_game_location_object_id) REFERENCES public.adventure_game_location_object(id),
    CONSTRAINT adventure_game_location_object_effect_result_location_fkey FOREIGN KEY (result_adventure_game_location_id) REFERENCES public.adventure_game_location(id)
);
CREATE INDEX idx_adventure_game_location_object_effect_game_id ON public.adventure_game_location_object_effect(game_id);
CREATE INDEX idx_adventure_game_location_object_effect_object_id ON public.adventure_game_location_object_effect(adventure_game_location_object_id);

CREATE TABLE public.adventure_game_location_object_instance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    game_instance_id UUID NOT NULL,
    adventure_game_location_object_id UUID NOT NULL,
    adventure_game_location_instance_id UUID NOT NULL,
    current_state VARCHAR(50) NOT NULL,
    is_visible BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_location_object_instance_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_location_object_instance_game_instance_id_fkey FOREIGN KEY (game_instance_id) REFERENCES public.game_instance(id),
    CONSTRAINT adventure_game_location_object_instance_object_id_fkey FOREIGN KEY (adventure_game_location_object_id) REFERENCES public.adventure_game_location_object(id),
    CONSTRAINT adventure_game_location_object_instance_location_inst_fkey FOREIGN KEY (adventure_game_location_instance_id) REFERENCES public.adventure_game_location_instance(id)
);
CREATE INDEX idx_adventure_game_location_object_instance_game_id ON public.adventure_game_location_object_instance(game_id);
CREATE INDEX idx_adventure_game_location_object_instance_game_instance_id ON public.adventure_game_location_object_instance(game_instance_id);
CREATE INDEX idx_adventure_game_location_object_instance_object_id ON public.adventure_game_location_object_instance(adventure_game_location_object_id);
CREATE INDEX idx_adventure_game_location_object_instance_location_inst ON public.adventure_game_location_object_instance(adventure_game_location_instance_id);

COMMIT;
