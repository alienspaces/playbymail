-- Remove ON DELETE CASCADE from foreign key constraints.
-- Deletion is the domain layer's responsibility, not the database's.

-- account_user_contact
ALTER TABLE public.account_user_contact
    DROP CONSTRAINT account_user_contact_account_user_id_fkey,
    ADD CONSTRAINT account_user_contact_account_user_id_fkey
        FOREIGN KEY (account_user_id) REFERENCES public.account_user(id);

-- account_subscription
ALTER TABLE public.account_subscription
    DROP CONSTRAINT account_subscription_account_id_fkey,
    ADD CONSTRAINT account_subscription_account_id_fkey
        FOREIGN KEY (account_id) REFERENCES public.account(id);

ALTER TABLE public.account_subscription
    DROP CONSTRAINT account_subscription_account_user_id_fkey,
    ADD CONSTRAINT account_subscription_account_user_id_fkey
        FOREIGN KEY (account_user_id) REFERENCES public.account_user(id);

-- game_subscription
ALTER TABLE public.game_subscription
    DROP CONSTRAINT game_subscription_account_id_fkey,
    ADD CONSTRAINT game_subscription_account_id_fkey
        FOREIGN KEY (account_id) REFERENCES public.account(id);

ALTER TABLE public.game_subscription
    DROP CONSTRAINT game_subscription_account_user_id_fkey,
    ADD CONSTRAINT game_subscription_account_user_id_fkey
        FOREIGN KEY (account_user_id) REFERENCES public.account_user(id);

-- game_subscription_instance
ALTER TABLE public.game_subscription_instance
    DROP CONSTRAINT game_subscription_instance_account_id_fkey,
    ADD CONSTRAINT game_subscription_instance_account_id_fkey
        FOREIGN KEY (account_id) REFERENCES public.account(id);
