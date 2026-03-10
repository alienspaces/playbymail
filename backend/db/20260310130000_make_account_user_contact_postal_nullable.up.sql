-- Allow contact record to be created with no contact info (e.g. on new login).
-- Only account_user_id is required.
ALTER TABLE public.account_user_contact ALTER COLUMN postal_address_line1 DROP NOT NULL;
ALTER TABLE public.account_user_contact ALTER COLUMN state_province DROP NOT NULL;
ALTER TABLE public.account_user_contact ALTER COLUMN country DROP NOT NULL;
ALTER TABLE public.account_user_contact ALTER COLUMN postal_code DROP NOT NULL;
