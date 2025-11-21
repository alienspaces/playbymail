-- Restore name column to account table
ALTER TABLE public.account
    ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '';

