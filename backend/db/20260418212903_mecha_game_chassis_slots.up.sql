-- Add aggregate weapon/equipment slot counts to mecha_game_chassis.
-- Slots are budgeted as three bands (small/medium/large); mounted items use
-- their native slot first, and smaller items can spill up into larger slots
-- when their native band is full. Large items never fit a medium/small slot.
ALTER TABLE mecha_game_chassis
    ADD COLUMN small_slots  INTEGER NOT NULL DEFAULT 2,
    ADD COLUMN medium_slots INTEGER NOT NULL DEFAULT 2,
    ADD COLUMN large_slots  INTEGER NOT NULL DEFAULT 1,
    ADD CONSTRAINT mecha_game_chassis_slot_bounds_check
        CHECK (small_slots  BETWEEN 0 AND 10
           AND medium_slots BETWEEN 0 AND 10
           AND large_slots  BETWEEN 0 AND 10);

-- Backfill per-class defaults. These keep every existing demo loadout within
-- capacity (verified: light mechs carry 1 small weapon, medium mechs carry up
-- to 2 medium weapons). Heavy/assault defaults leave headroom for the kinds of
-- loadouts those classes are expected to field.
UPDATE mecha_game_chassis SET small_slots = 2, medium_slots = 1, large_slots = 0 WHERE chassis_class = 'light';
UPDATE mecha_game_chassis SET small_slots = 2, medium_slots = 2, large_slots = 1 WHERE chassis_class = 'medium';
UPDATE mecha_game_chassis SET small_slots = 2, medium_slots = 2, large_slots = 2 WHERE chassis_class = 'heavy';
UPDATE mecha_game_chassis SET small_slots = 2, medium_slots = 3, large_slots = 3 WHERE chassis_class = 'assault';
