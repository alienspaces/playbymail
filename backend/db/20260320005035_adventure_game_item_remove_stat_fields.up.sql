-- Remove hardcoded stat fields from adventure_game_item.
-- These are replaced by adventure_game_item_effect records with
-- effect_type 'weapon_damage', 'armor_defense', 'heal_wielder', and 'heal_target'.
-- Whether an item is usable is now derived from whether it has any use-action effects.

BEGIN;

ALTER TABLE public.adventure_game_item
    DROP COLUMN IF EXISTS damage,
    DROP COLUMN IF EXISTS defense,
    DROP COLUMN IF EXISTS heal_amount,
    DROP COLUMN IF EXISTS can_be_used;

COMMIT;
