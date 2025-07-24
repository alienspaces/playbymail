-- Create table for adventure_game_creature_placement (configuration for creature placement in locations)
CREATE TABLE IF NOT EXISTS public.adventure_game_creature_placement (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id                         UUID NOT NULL REFERENCES public.game(id), -- The game this placement configuration belongs to
    adventure_game_creature_id      UUID NOT NULL REFERENCES public.adventure_game_creature(id), -- The creature type to be placed
    adventure_game_location_id      UUID NOT NULL REFERENCES public.adventure_game_location(id), -- The location where the creature should be placed
    initial_count                   INTEGER NOT NULL DEFAULT 1, -- How many of this creature should be placed initially
    created_at                      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMP WITH TIME ZONE,
    deleted_at                      TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_adventure_game_creature_placement_game_id ON public.adventure_game_creature_placement(game_id);
CREATE INDEX IF NOT EXISTS idx_adventure_game_creature_placement_creature_id ON public.adventure_game_creature_placement(adventure_game_creature_id);
CREATE INDEX IF NOT EXISTS idx_adventure_game_creature_placement_location_id ON public.adventure_game_creature_placement(adventure_game_location_id);
CREATE INDEX IF NOT EXISTS idx_adventure_game_creature_placement_deleted_at ON public.adventure_game_creature_placement(deleted_at);

-- Add unique constraint to prevent duplicate placements for the same creature in the same location for the same game
CREATE UNIQUE INDEX IF NOT EXISTS idx_adventure_game_creature_placement_unique 
ON public.adventure_game_creature_placement(game_id, adventure_game_creature_id, adventure_game_location_id) 
WHERE deleted_at IS NULL;
