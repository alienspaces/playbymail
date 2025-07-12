CREATE TABLE public.game_character (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.game(id),
    account_id UUID NOT NULL REFERENCES public.account(id),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_game_character_game_id ON public.game_character(game_id);
CREATE INDEX idx_game_character_account_id ON public.game_character(account_id);

ALTER TABLE public.game_character ADD CONSTRAINT game_character_unique UNIQUE (game_id, account_id);
ALTER TABLE public.game_character ADD CONSTRAINT game_character_name_unique UNIQUE (game_id, name);
ALTER TABLE public.game_character ADD CONSTRAINT game_character_name_not_empty CHECK (name IS NOT NULL AND name != '');