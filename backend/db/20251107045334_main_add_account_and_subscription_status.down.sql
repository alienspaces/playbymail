-- Remove status column from game_subscription
ALTER TABLE public.game_subscription
    DROP COLUMN IF EXISTS status;

-- Remove status column from account
ALTER TABLE public.account
    DROP COLUMN IF EXISTS status;

-- Reinstate NOT NULL constraint on game turn sheet game_instance_id
ALTER TABLE public.game_turn_sheet
    ALTER COLUMN game_instance_id SET NOT NULL;
