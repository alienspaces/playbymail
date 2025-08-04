-- Make location description required
ALTER TABLE public.adventure_game_location 
    ALTER COLUMN description SET NOT NULL;

-- Make location link description required
ALTER TABLE public.adventure_game_location_link 
    ALTER COLUMN description SET NOT NULL;
