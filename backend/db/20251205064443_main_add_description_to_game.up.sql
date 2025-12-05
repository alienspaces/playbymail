-- Add description column to game table
-- This field stores the game description that appears on the join game turn sheet
-- First add as nullable, then set default for existing records, then make NOT NULL

ALTER TABLE game
ADD COLUMN description TEXT;

-- Set default description for existing games
UPDATE game
SET description = 'Welcome to ' || name || '! Welcome to the PlayByMail Adventure!'
WHERE description IS NULL;

-- Now make it NOT NULL
ALTER TABLE game
ALTER COLUMN description SET NOT NULL;

COMMENT ON COLUMN game.description IS 'Game description that appears on the join game turn sheet (required)';

