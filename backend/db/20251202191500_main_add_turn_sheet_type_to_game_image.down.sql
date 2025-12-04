-- Remove turn_sheet_type column and related constraints/indexes
DROP INDEX IF EXISTS public.game_image_unique_turn_sheet_background;
DROP INDEX IF EXISTS public.game_image_unique_other_types;
DROP INDEX IF EXISTS public.idx_game_image_turn_sheet_type;
ALTER TABLE public.game_image DROP CONSTRAINT IF EXISTS game_image_turn_sheet_type_check;
ALTER TABLE public.game_image DROP COLUMN IF EXISTS turn_sheet_type;

-- Restore original unique constraint
ALTER TABLE public.game_image ADD CONSTRAINT game_image_game_id_record_id_type_key UNIQUE(game_id, record_id, type);

