-- Remove aggregate weapon/equipment slot counts from mecha_game_chassis.
ALTER TABLE mecha_game_chassis
    DROP CONSTRAINT IF EXISTS mecha_game_chassis_slot_bounds_check;

ALTER TABLE mecha_game_chassis
    DROP COLUMN IF EXISTS small_slots,
    DROP COLUMN IF EXISTS medium_slots,
    DROP COLUMN IF EXISTS large_slots;
