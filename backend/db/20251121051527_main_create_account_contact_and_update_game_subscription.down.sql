-- Rollback: Restore account name and remove account_contact

-- Add name back to account table
ALTER TABLE public.account
    ADD COLUMN name VARCHAR(255) DEFAULT '';

-- Remove account_contact_id from game_subscription
ALTER TABLE public.game_subscription
    DROP COLUMN IF EXISTS account_contact_id;

-- Drop account_contact table
DROP TABLE IF EXISTS public.account_contact;

