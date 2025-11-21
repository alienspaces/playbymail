-- Remove name column from account table
-- Name is now part of account_contact, not account
ALTER TABLE public.account
    DROP COLUMN IF EXISTS name;

