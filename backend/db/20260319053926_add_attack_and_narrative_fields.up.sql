BEGIN;

-- Creature attack flavor and lifecycle fields
ALTER TABLE public.adventure_game_creature
    ADD COLUMN attack_method VARCHAR(50) NOT NULL DEFAULT 'claws',
    ADD COLUMN attack_description VARCHAR(512) NOT NULL DEFAULT '',
    ADD COLUMN body_decay_turns INT NOT NULL DEFAULT 3,
    ADD COLUMN respawn_turns INT NOT NULL DEFAULT 0;

ALTER TABLE public.adventure_game_creature
    ADD CONSTRAINT adventure_game_creature_attack_method_check CHECK (
        attack_method IN ('claws', 'bite', 'sting', 'weapon', 'spell', 'slam', 'touch', 'breath', 'gaze')
    );

-- Track when a creature instance died (NULL = alive)
ALTER TABLE public.adventure_game_creature_instance
    ADD COLUMN died_at_turn INT;

-- Travel narrative on location links
ALTER TABLE public.adventure_game_location_link
    ADD COLUMN traversal_description VARCHAR(1024);

-- Accumulated turn events on character instances, consumed next turn
ALTER TABLE public.adventure_game_character_instance
    ADD COLUMN last_turn_events JSONB NOT NULL DEFAULT '[]';

COMMIT;
