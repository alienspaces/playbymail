-- Backfill account_user_id from game_subscription where missing
UPDATE public.game_subscription_instance gsi
SET account_user_id = gs.account_user_id
FROM public.game_subscription gs
WHERE gsi.game_subscription_id = gs.id
  AND gsi.account_user_id IS NULL
  AND gs.account_user_id IS NOT NULL;

-- Remove any rows that could not be backfilled (subscription had no account_user_id)
DELETE FROM public.game_subscription_instance WHERE account_user_id IS NULL;

ALTER TABLE public.game_subscription_instance ALTER COLUMN account_user_id SET NOT NULL;
