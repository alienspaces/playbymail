-- Mecha Schema Rollback

BEGIN;

DROP TABLE IF EXISTS public.mecha_turn_sheet;
DROP TABLE IF EXISTS public.mecha_mech_instance;
DROP TABLE IF EXISTS public.mecha_lance_instance;
DROP TABLE IF EXISTS public.mecha_sector_instance;
DROP TABLE IF EXISTS public.mecha_lance_mech;
DROP TABLE IF EXISTS public.mecha_lance;
DROP TABLE IF EXISTS public.mecha_sector_link;
DROP TABLE IF EXISTS public.mecha_sector;
DROP TABLE IF EXISTS public.mecha_weapon;
DROP TABLE IF EXISTS public.mecha_chassis;

COMMIT;
