ALTER TABLE game_location_instance ADD COLUMN game_id UUID NOT NULL REFERENCES game(id);
