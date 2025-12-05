-- Revert game description column back to TEXT
-- This removes the 512 character limit

ALTER TABLE game
ALTER COLUMN description TYPE TEXT;

COMMENT ON COLUMN game.description IS 'Game description that appears on the join game turn sheet (required)';

