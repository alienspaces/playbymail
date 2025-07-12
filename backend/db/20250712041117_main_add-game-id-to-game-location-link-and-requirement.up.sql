ALTER TABLE game_location_link ADD COLUMN game_id UUID NOT NULL REFERENCES game(id);
ALTER TABLE game_location_link_requirement ADD COLUMN game_id UUID NOT NULL REFERENCES game(id);
