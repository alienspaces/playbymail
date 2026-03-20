-- Adventure Game Location Object Schema Teardown

BEGIN;

DROP TABLE IF EXISTS public.adventure_game_location_object_instance;
DROP TABLE IF EXISTS public.adventure_game_location_object_effect;
DROP TABLE IF EXISTS public.adventure_game_location_object;

COMMIT;
