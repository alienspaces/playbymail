-- Revert location description to nullable
ALTER TABLE public.adventure_game_location 
    ALTER COLUMN description DROP NOT NULL;

-- Revert location link description to nullable
ALTER TABLE public.adventure_game_location_link 
    ALTER COLUMN description DROP NOT NULL;
