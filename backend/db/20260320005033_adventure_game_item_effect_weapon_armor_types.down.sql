BEGIN;

ALTER TABLE public.adventure_game_item_effect
    DROP CONSTRAINT adventure_game_item_effect_effect_type_check;

ALTER TABLE public.adventure_game_item_effect
    ADD CONSTRAINT adventure_game_item_effect_effect_type_check CHECK (
        effect_type IN (
            'info', 'damage_target', 'damage_wielder', 'heal_target', 'heal_wielder',
            'teleport', 'open_link', 'close_link', 'give_item', 'remove_item',
            'summon_creature', 'nothing'
        )
    );

COMMIT;
