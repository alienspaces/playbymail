-- Adventure Game Item Effect Schema Migration

BEGIN;

-- adventure_game_item_effect defines what happens when a player performs an action on an item.
-- Multiple rows for the same (item, action_type, required conditions) are allowed — all fire atomically.
CREATE TABLE public.adventure_game_item_effect (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    adventure_game_item_id UUID NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    required_adventure_game_item_id UUID,
    required_adventure_game_location_id UUID,
    result_description TEXT NOT NULL,
    effect_type VARCHAR(50) NOT NULL,
    result_adventure_game_item_id UUID,
    result_adventure_game_location_link_id UUID,
    result_adventure_game_creature_id UUID,
    result_adventure_game_location_id UUID,
    result_value_min INTEGER,
    result_value_max INTEGER,
    is_repeatable BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT adventure_game_item_effect_action_type_check CHECK (
        action_type IN ('use', 'equip', 'unequip', 'inspect', 'drop')
    ),
    CONSTRAINT adventure_game_item_effect_effect_type_check CHECK (
        effect_type IN (
            'info', 'damage_target', 'damage_wielder', 'heal_target', 'heal_wielder',
            'teleport', 'open_link', 'close_link', 'give_item', 'remove_item',
            'summon_creature', 'nothing'
        )
    ),
    CONSTRAINT adventure_game_item_effect_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT adventure_game_item_effect_item_id_fkey FOREIGN KEY (adventure_game_item_id) REFERENCES public.adventure_game_item(id),
    CONSTRAINT adventure_game_item_effect_required_item_fkey FOREIGN KEY (required_adventure_game_item_id) REFERENCES public.adventure_game_item(id),
    CONSTRAINT adventure_game_item_effect_required_location_fkey FOREIGN KEY (required_adventure_game_location_id) REFERENCES public.adventure_game_location(id),
    CONSTRAINT adventure_game_item_effect_result_item_fkey FOREIGN KEY (result_adventure_game_item_id) REFERENCES public.adventure_game_item(id),
    CONSTRAINT adventure_game_item_effect_result_link_fkey FOREIGN KEY (result_adventure_game_location_link_id) REFERENCES public.adventure_game_location_link(id),
    CONSTRAINT adventure_game_item_effect_result_creature_fkey FOREIGN KEY (result_adventure_game_creature_id) REFERENCES public.adventure_game_creature(id),
    CONSTRAINT adventure_game_item_effect_result_location_fkey FOREIGN KEY (result_adventure_game_location_id) REFERENCES public.adventure_game_location(id)
);
CREATE INDEX idx_adventure_game_item_effect_game_id ON public.adventure_game_item_effect(game_id);
CREATE INDEX idx_adventure_game_item_effect_item_id ON public.adventure_game_item_effect(adventure_game_item_id);

COMMIT;
