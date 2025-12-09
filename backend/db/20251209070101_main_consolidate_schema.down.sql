-- Consolidated schema migration - rollback
-- Drops all tables in reverse dependency order

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
DROP TABLE IF EXISTS public.game_turn_sheet;
DROP TABLE IF EXISTS public.game_instance_parameter;
DROP TABLE IF EXISTS public.game_instance;
DROP TABLE IF EXISTS public.game_subscription;
DROP TABLE IF EXISTS public.game_image;
DROP TABLE IF EXISTS public.game;
DROP TABLE IF EXISTS public.account_contact;
DROP TABLE IF EXISTS public.account;

