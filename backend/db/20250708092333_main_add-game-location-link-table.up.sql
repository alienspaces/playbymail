CREATE TABLE public.game_location_link (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_game_location_id UUID NOT NULL REFERENCES public.game_location(id),
    to_game_location_id UUID NOT NULL REFERENCES public.game_location(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_game_location_link_from_game_location_id ON public.game_location_link(from_game_location_id);
CREATE INDEX idx_game_location_link_to_game_location_id ON public.game_location_link(to_game_location_id);

-- Prevent duplicate links
ALTER TABLE public.game_location_link ADD CONSTRAINT location_link_unique UNIQUE (from_game_location_id, to_game_location_id);
