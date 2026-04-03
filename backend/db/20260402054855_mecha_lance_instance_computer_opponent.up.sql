-- Add mecha_computer_opponent_id to mecha_lance_instance so the AI turn
-- processor can find each opponent's assigned lance instance without needing
-- to store ownership on the design-time mecha_lance table.

BEGIN;

ALTER TABLE public.mecha_lance_instance
    ADD COLUMN mecha_computer_opponent_id UUID,
    ADD CONSTRAINT mecha_lance_instance_computer_opponent_id_fkey
        FOREIGN KEY (mecha_computer_opponent_id) REFERENCES public.mecha_computer_opponent(id);

CREATE INDEX idx_mecha_lance_instance_computer_opponent
    ON public.mecha_lance_instance(mecha_computer_opponent_id)
    WHERE mecha_computer_opponent_id IS NOT NULL;

COMMIT;
