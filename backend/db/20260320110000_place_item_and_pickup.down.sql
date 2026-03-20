BEGIN;

ALTER TABLE public.adventure_game_location_object_effect
    DROP CONSTRAINT adventure_game_location_object_effect_effect_type_check;

ALTER TABLE public.adventure_game_location_object_effect
    ADD CONSTRAINT adventure_game_location_object_effect_effect_type_check CHECK (
        effect_type IN (
            'info', 'change_state', 'change_object_state', 'give_item', 'remove_item',
            'open_link', 'close_link', 'reveal_object', 'hide_object',
            'damage', 'heal', 'summon_creature', 'teleport', 'nothing', 'remove_object'
        )
    );

ALTER TABLE public.adventure_game_item_effect
    DROP CONSTRAINT adventure_game_item_effect_action_type_check;

ALTER TABLE public.adventure_game_item_effect
    ADD CONSTRAINT adventure_game_item_effect_action_type_check CHECK (
        action_type IN ('use', 'equip', 'unequip', 'inspect', 'drop')
    );

COMMIT;
