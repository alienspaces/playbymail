-- Refactor mecha_lance to a pure template with a lance_type enum.
-- Player-owned lances (account_user_id IS NOT NULL) are runtime artifacts
-- created at game start; they do not belong in the design table.

BEGIN;

-- 1. Add lance_type column defaulting to 'opponent' so existing computer-opponent
--    rows satisfy the NOT NULL constraint before we enforce it.
ALTER TABLE public.mecha_lance
    ADD COLUMN lance_type VARCHAR(20) NOT NULL DEFAULT 'opponent';

-- 2. Stamp starter templates
UPDATE public.mecha_lance
    SET lance_type = 'starter'
    WHERE is_player_starter = true;

-- 3. Delete player-owned lances (they are runtime artifacts)
DELETE FROM public.mecha_lance
    WHERE account_user_id IS NOT NULL;

-- 4. Drop the default now that the column is populated
ALTER TABLE public.mecha_lance
    ALTER COLUMN lance_type DROP DEFAULT;

-- 5. Enforce the enum
ALTER TABLE public.mecha_lance
    ADD CONSTRAINT mecha_lance_type_check CHECK (lance_type IN ('starter', 'opponent'));

-- 6. Drop old ownership constraints before dropping columns
ALTER TABLE public.mecha_lance
    DROP CONSTRAINT IF EXISTS mecha_lance_owner_check,
    DROP CONSTRAINT IF EXISTS mecha_lance_account_id_fkey,
    DROP CONSTRAINT IF EXISTS mecha_lance_account_user_id_fkey,
    DROP CONSTRAINT IF EXISTS mecha_lance_computer_opponent_id_fkey;

-- 7. Drop old ownership indexes
DROP INDEX IF EXISTS public.idx_mecha_lance_account_id;
DROP INDEX IF EXISTS public.idx_mecha_lance_computer_opponent_id;
DROP INDEX IF EXISTS public.idx_mecha_lance_game_account_unique;
DROP INDEX IF EXISTS public.idx_mecha_lance_player_starter_unique;

-- 8. Drop old ownership columns
ALTER TABLE public.mecha_lance
    DROP COLUMN IF EXISTS is_player_starter,
    DROP COLUMN IF EXISTS account_id,
    DROP COLUMN IF EXISTS account_user_id,
    DROP COLUMN IF EXISTS mecha_computer_opponent_id;

-- 9. Add partial unique index: at most one starter lance per game
CREATE UNIQUE INDEX idx_mecha_lance_starter_unique
    ON public.mecha_lance (game_id)
    WHERE lance_type = 'starter' AND deleted_at IS NULL;

COMMIT;
