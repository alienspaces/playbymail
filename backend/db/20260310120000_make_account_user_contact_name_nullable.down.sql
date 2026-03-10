-- Restore NOT NULL; existing NULLs would need to be updated before rolling back
ALTER TABLE public.account_user_contact ALTER COLUMN name SET NOT NULL;
