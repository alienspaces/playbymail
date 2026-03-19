BEGIN;

ALTER TABLE public.adventure_game_character_instance
    DROP COLUMN IF EXISTS last_turn_events;

ALTER TABLE public.adventure_game_location_link
    DROP COLUMN IF EXISTS traversal_description;

ALTER TABLE public.adventure_game_creature_instance
    DROP COLUMN IF EXISTS died_at_turn;

ALTER TABLE public.adventure_game_creature
    DROP CONSTRAINT IF EXISTS adventure_game_creature_attack_method_check,
    DROP COLUMN IF EXISTS respawn_turns,
    DROP COLUMN IF EXISTS body_decay_turns,
    DROP COLUMN IF EXISTS attack_description,
    DROP COLUMN IF EXISTS attack_method;

COMMIT;
