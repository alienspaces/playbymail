ALTER TABLE public.location RENAME TO game_location;

-- Update foreign keys in location_link
ALTER TABLE public.location_link DROP CONSTRAINT IF EXISTS location_link_from_location_id_fkey;
ALTER TABLE public.location_link DROP CONSTRAINT IF EXISTS location_link_to_location_id_fkey;
ALTER TABLE public.location_link
    RENAME COLUMN from_location_id TO from_game_location_id;
ALTER TABLE public.location_link
    RENAME COLUMN to_location_id TO to_game_location_id;
ALTER TABLE public.location_link
    ADD CONSTRAINT location_link_from_game_location_id_fkey FOREIGN KEY (from_game_location_id) REFERENCES public.game_location(id);
ALTER TABLE public.location_link
    ADD CONSTRAINT location_link_to_game_location_id_fkey FOREIGN KEY (to_game_location_id) REFERENCES public.game_location(id);

-- Update indexes
DROP INDEX IF EXISTS idx_location_link_from_location_id;
DROP INDEX IF EXISTS idx_location_link_to_location_id;
CREATE INDEX idx_location_link_from_game_location_id ON public.location_link(from_game_location_id);
CREATE INDEX idx_location_link_to_game_location_id ON public.location_link(to_game_location_id);
