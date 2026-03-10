-- All game subscriptions must have an account user. Backfill from account (one user per account).
UPDATE public.game_subscription gs
SET account_user_id = (
    SELECT au.id FROM public.account_user au
    WHERE au.account_id = gs.account_id
    LIMIT 1
)
WHERE gs.account_user_id IS NULL AND gs.account_id IS NOT NULL;

-- Remove any rows that could not be backfilled (no account_user for the account)
DELETE FROM public.game_subscription WHERE account_user_id IS NULL;

ALTER TABLE public.game_subscription ALTER COLUMN account_user_id SET NOT NULL;
