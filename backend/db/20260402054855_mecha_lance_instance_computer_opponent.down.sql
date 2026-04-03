BEGIN;

DROP INDEX IF EXISTS public.idx_mecha_lance_instance_computer_opponent;

ALTER TABLE public.mecha_lance_instance
    DROP CONSTRAINT IF EXISTS mecha_lance_instance_computer_opponent_id_fkey,
    DROP COLUMN IF EXISTS mecha_computer_opponent_id;

COMMIT;
