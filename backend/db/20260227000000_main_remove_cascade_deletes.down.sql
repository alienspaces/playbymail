-- Restore ON DELETE CASCADE foreign key constraints.

-- account_user_contact
ALTER TABLE public.account_user_contact
    DROP CONSTRAINT account_user_contact_account_user_id_fkey,
    ADD CONSTRAINT account_user_contact_account_user_id_fkey
        FOREIGN KEY (account_user_id) REFERENCES public.account_user(id) ON DELETE CASCADE;

-- account_subscription
ALTER TABLE public.account_subscription
    DROP CONSTRAINT account_subscription_account_id_fkey,
    ADD CONSTRAINT account_subscription_account_id_fkey
        FOREIGN KEY (account_id) REFERENCES public.account(id) ON DELETE CASCADE;

ALTER TABLE public.account_subscription
    DROP CONSTRAINT account_subscription_account_user_id_fkey,
    ADD CONSTRAINT account_subscription_account_user_id_fkey
        FOREIGN KEY (account_user_id) REFERENCES public.account_user(id) ON DELETE CASCADE;

-- game_subscription
ALTER TABLE public.game_subscription
    DROP CONSTRAINT game_subscription_account_id_fkey,
    ADD CONSTRAINT game_subscription_account_id_fkey
        FOREIGN KEY (account_id) REFERENCES public.account(id) ON DELETE CASCADE;

ALTER TABLE public.game_subscription
    DROP CONSTRAINT game_subscription_account_user_id_fkey,
    ADD CONSTRAINT game_subscription_account_user_id_fkey
        FOREIGN KEY (account_user_id) REFERENCES public.account_user(id) ON DELETE CASCADE;

-- game_subscription_instance
ALTER TABLE public.game_subscription_instance
    DROP CONSTRAINT game_subscription_instance_account_id_fkey,
    ADD CONSTRAINT game_subscription_instance_account_id_fkey
        FOREIGN KEY (account_id) REFERENCES public.account(id) ON DELETE CASCADE;
