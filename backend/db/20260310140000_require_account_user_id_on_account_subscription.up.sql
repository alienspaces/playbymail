-- All account subscriptions must have an account user. Backfill from account (one user per account).
UPDATE public.account_subscription acs
SET account_user_id = (
    SELECT au.id FROM public.account_user au
    WHERE au.account_id = acs.account_id
    LIMIT 1
)
WHERE acs.account_user_id IS NULL AND acs.account_id IS NOT NULL;

-- Remove any rows that could not be backfilled (no account_user for the account)
DELETE FROM public.account_subscription WHERE account_user_id IS NULL;

ALTER TABLE public.account_subscription ALTER COLUMN account_user_id SET NOT NULL;
