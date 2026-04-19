-- Revert mecha equipment catalog and integration columns.
BEGIN;

ALTER TABLE public.mecha_game_mech_instance
    DROP COLUMN IF EXISTS ammo_remaining,
    DROP COLUMN IF EXISTS equipment_config;

ALTER TABLE public.mecha_game_squad_mech
    DROP COLUMN IF EXISTS equipment_config;

ALTER TABLE public.mecha_game_weapon
    DROP CONSTRAINT IF EXISTS mecha_game_weapon_ammo_capacity_check,
    DROP COLUMN IF EXISTS ammo_capacity;

DROP TABLE IF EXISTS public.mecha_game_equipment;

COMMIT;
