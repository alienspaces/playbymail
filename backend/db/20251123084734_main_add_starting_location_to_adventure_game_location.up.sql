-- Add is_starting_location field to adventure_game_location table
ALTER TABLE public.adventure_game_location
    ADD COLUMN is_starting_location BOOLEAN NOT NULL DEFAULT false;

-- Add index for querying starting locations
CREATE INDEX idx_adventure_game_location_is_starting_location
    ON public.adventure_game_location(is_starting_location)
    WHERE is_starting_location = true;

COMMENT ON COLUMN public.adventure_game_location.is_starting_location IS 'Indicates if this location is a valid starting location for new players joining the game.';

