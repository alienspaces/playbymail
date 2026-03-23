-- Mech Wargame Schema Rollback

BEGIN;

DROP TABLE IF EXISTS public.mech_wargame_turn_sheet;
DROP TABLE IF EXISTS public.mech_wargame_mech_instance;
DROP TABLE IF EXISTS public.mech_wargame_lance_instance;
DROP TABLE IF EXISTS public.mech_wargame_sector_instance;
DROP TABLE IF EXISTS public.mech_wargame_lance_mech;
DROP TABLE IF EXISTS public.mech_wargame_lance;
DROP TABLE IF EXISTS public.mech_wargame_sector_link;
DROP TABLE IF EXISTS public.mech_wargame_sector;
DROP TABLE IF EXISTS public.mech_wargame_weapon;
DROP TABLE IF EXISTS public.mech_wargame_chassis;

COMMIT;
