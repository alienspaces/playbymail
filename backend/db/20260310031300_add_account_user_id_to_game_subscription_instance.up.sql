ALTER TABLE public.game_subscription_instance ADD COLUMN account_user_id UUID REFERENCES public.account_user(id);
CREATE INDEX idx_game_subscription_instance_account_user_id ON public.game_subscription_instance(account_user_id) WHERE deleted_at IS NULL;
