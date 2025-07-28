-- Drop consolidated instance tables in reverse order

-- Drop location link item requirements
DROP TABLE IF EXISTS adventure_game_location_link_requirement;

-- Drop item instances
DROP TABLE IF EXISTS adventure_game_item_instance;

-- Drop character instances
DROP TABLE IF EXISTS adventure_game_character_instance;

-- Drop creature instances
DROP TABLE IF EXISTS adventure_game_creature_instance;

-- Drop creatures
DROP TABLE IF EXISTS adventure_game_creature;

-- Drop location instances
DROP TABLE IF EXISTS adventure_game_location_instance;

-- Drop game instances
DROP TABLE IF EXISTS game_instance; 