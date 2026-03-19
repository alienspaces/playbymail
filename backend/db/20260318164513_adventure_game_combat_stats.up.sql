BEGIN;

-- Add combat stats to adventure_game_creature
ALTER TABLE public.adventure_game_creature
    ADD COLUMN attack_damage INT NOT NULL DEFAULT 10,
    ADD COLUMN defense INT NOT NULL DEFAULT 0,
    ADD COLUMN disposition VARCHAR(20) NOT NULL DEFAULT 'aggressive';

ALTER TABLE public.adventure_game_creature
    ADD CONSTRAINT adventure_game_creature_disposition_check CHECK (
        disposition IN ('aggressive', 'inquisitive', 'indifferent')
    );

-- Add combat/usage fields to adventure_game_item
ALTER TABLE public.adventure_game_item
    ADD COLUMN damage INT NOT NULL DEFAULT 0,
    ADD COLUMN defense INT NOT NULL DEFAULT 0,
    ADD COLUMN heal_amount INT NOT NULL DEFAULT 0;

COMMIT;
