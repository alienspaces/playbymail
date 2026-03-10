-- Restore NOT NULL; existing NULLs would need to be updated before rolling back
ALTER TABLE public.account_user_contact ALTER COLUMN postal_address_line1 SET NOT NULL;
ALTER TABLE public.account_user_contact ALTER COLUMN state_province SET NOT NULL;
ALTER TABLE public.account_user_contact ALTER COLUMN country SET NOT NULL;
ALTER TABLE public.account_user_contact ALTER COLUMN postal_code SET NOT NULL;
