-- Adventure Game Schema Teardown
-- Generated: 2026-01-08

BEGIN;

DROP TABLE IF EXISTS public.adventure_game_turn_sheet;
DROP TABLE IF EXISTS public.adventure_game_item_instance;
DROP TABLE IF EXISTS public.adventure_game_creature_instance;
DROP TABLE IF EXISTS public.adventure_game_character_instance;
DROP TABLE IF EXISTS public.adventure_game_character;
DROP TABLE IF EXISTS public.adventure_game_creature_placement;
DROP TABLE IF EXISTS public.adventure_game_creature;
DROP TABLE IF EXISTS public.adventure_game_location_link_requirement;
DROP TABLE IF EXISTS public.adventure_game_item_placement;
DROP TABLE IF EXISTS public.adventure_game_item;
DROP TABLE IF EXISTS public.adventure_game_location_link;
DROP TABLE IF EXISTS public.adventure_game_location_instance;
DROP TABLE IF EXISTS public.adventure_game_location;

COMMIT;
