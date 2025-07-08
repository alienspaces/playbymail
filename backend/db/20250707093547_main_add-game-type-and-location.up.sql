ALTER TABLE public.game ADD COLUMN game_type VARCHAR(50) NOT NULL DEFAULT 'adventure';
ALTER TABLE public.game ADD CONSTRAINT game_type_check CHECK (game_type = 'adventure');

CREATE TABLE IF NOT EXISTS public.location (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES public.game(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
