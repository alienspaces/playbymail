CREATE TABLE public.adventure_game_location_link (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.game(id), -- The game this link belongs to
    from_adventure_game_location_id UUID NOT NULL REFERENCES public.adventure_game_location(id),
    to_adventure_game_location_id UUID NOT NULL REFERENCES public.adventure_game_location(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_adventure_game_location_link_from_adventure_game_location_id 
    ON public.adventure_game_location_link(from_adventure_game_location_id);
CREATE INDEX idx_adventure_game_location_link_to_adventure_game_location_id 
    ON public.adventure_game_location_link(to_adventure_game_location_id);

-- Prevent duplicate links
ALTER TABLE public.adventure_game_location_link 
    ADD CONSTRAINT adventure_game_location_link_unique 
    UNIQUE (from_adventure_game_location_id, to_adventure_game_location_id);
