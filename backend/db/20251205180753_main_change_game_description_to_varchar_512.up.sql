-- Change game description column from TEXT to VARCHAR(512)
-- This limits the description length to 512 characters

ALTER TABLE game
ALTER COLUMN description TYPE VARCHAR(512);

COMMENT ON COLUMN game.description IS 'Game description that appears on the join game turn sheet (required, max 512 characters)';

