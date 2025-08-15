CREATE TABLE IF NOT EXISTS public.game (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    game_type VARCHAR(50) NOT NULL,
    turn_duration_hours INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
	deleted_at TIMESTAMP WITH TIME ZONE,
    -- Constraints
    CONSTRAINT game_name_check CHECK (name != ''),
    CONSTRAINT game_type_check CHECK (game_type = 'adventure'),
    CONSTRAINT turn_duration_hours_check CHECK (turn_duration_hours > 0)
);
