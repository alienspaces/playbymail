-- Add mecha_computer_opponent table and wire it to mecha_lance

BEGIN;

CREATE TABLE public.mecha_computer_opponent (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    aggression INTEGER NOT NULL DEFAULT 5,
    iq INTEGER NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT mecha_computer_opponent_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id),
    CONSTRAINT mecha_computer_opponent_name_unique UNIQUE (game_id, name),
    CONSTRAINT mecha_computer_opponent_aggression_check CHECK (aggression BETWEEN 1 AND 10),
    CONSTRAINT mecha_computer_opponent_iq_check CHECK (iq BETWEEN 1 AND 10)
);
CREATE INDEX idx_mecha_computer_opponent_game_id ON public.mecha_computer_opponent(game_id);
COMMENT ON TABLE public.mecha_computer_opponent IS 'Computer-controlled opponent command owning one or more lances. Holds behaviour config (aggression, iq) used by the decision engine during turn processing.';

ALTER TABLE public.mecha_lance
    ADD COLUMN mecha_computer_opponent_id UUID,
    ALTER COLUMN account_id DROP NOT NULL,
    ALTER COLUMN account_user_id DROP NOT NULL,
    ADD CONSTRAINT mecha_lance_computer_opponent_id_fkey
        FOREIGN KEY (mecha_computer_opponent_id) REFERENCES public.mecha_computer_opponent(id),
    ADD CONSTRAINT mecha_lance_owner_check CHECK (
        (mecha_computer_opponent_id IS NOT NULL AND account_id IS NULL AND account_user_id IS NULL)
        OR
        (mecha_computer_opponent_id IS NULL AND account_id IS NOT NULL AND account_user_id IS NOT NULL)
    );
CREATE INDEX idx_mecha_lance_computer_opponent_id ON public.mecha_lance(mecha_computer_opponent_id);

-- Replace (game_id, account_id) unique constraint with a partial index so NULLs are handled cleanly
ALTER TABLE public.mecha_lance DROP CONSTRAINT mecha_lance_unique;
CREATE UNIQUE INDEX idx_mecha_lance_game_account_unique
    ON public.mecha_lance(game_id, account_id) WHERE account_id IS NOT NULL;

ALTER TABLE public.mecha_lance_instance
    ALTER COLUMN game_subscription_instance_id DROP NOT NULL;

COMMIT;
