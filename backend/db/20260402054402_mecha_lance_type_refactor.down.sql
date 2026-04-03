-- Reverse mecha_lance_type_refactor: restore original ownership columns.

BEGIN;

DROP INDEX IF EXISTS public.idx_mecha_lance_starter_unique;

ALTER TABLE public.mecha_lance
    DROP CONSTRAINT IF EXISTS mecha_lance_type_check,
    DROP COLUMN IF EXISTS lance_type,
    ADD COLUMN is_player_starter BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN account_id UUID,
    ADD COLUMN account_user_id UUID,
    ADD COLUMN mecha_computer_opponent_id UUID;

ALTER TABLE public.mecha_lance
    ADD CONSTRAINT mecha_lance_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id),
    ADD CONSTRAINT mecha_lance_account_user_id_fkey FOREIGN KEY (account_user_id) REFERENCES public.account_user(id),
    ADD CONSTRAINT mecha_lance_computer_opponent_id_fkey FOREIGN KEY (mecha_computer_opponent_id) REFERENCES public.mecha_computer_opponent(id),
    ADD CONSTRAINT mecha_lance_owner_check CHECK (
        (is_player_starter = true  AND mecha_computer_opponent_id IS NULL AND account_id IS NULL AND account_user_id IS NULL)
        OR
        (is_player_starter = false AND mecha_computer_opponent_id IS NOT NULL AND account_id IS NULL AND account_user_id IS NULL)
        OR
        (is_player_starter = false AND mecha_computer_opponent_id IS NULL AND account_id IS NOT NULL AND account_user_id IS NOT NULL)
    );

CREATE INDEX idx_mecha_lance_account_id ON public.mecha_lance(account_id);
CREATE INDEX idx_mecha_lance_computer_opponent_id ON public.mecha_lance(mecha_computer_opponent_id);
CREATE UNIQUE INDEX idx_mecha_lance_game_account_unique ON public.mecha_lance(game_id, account_id) WHERE account_id IS NOT NULL;
CREATE UNIQUE INDEX idx_mecha_lance_player_starter_unique ON public.mecha_lance(game_id) WHERE is_player_starter = true AND deleted_at IS NULL;

COMMIT;
