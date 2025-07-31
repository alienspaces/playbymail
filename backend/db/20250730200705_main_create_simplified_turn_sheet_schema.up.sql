-- Create simplified turn sheet schema
-- This migration creates the final, simplified turn sheet design with:
-- - Single game_turn_sheet table for all game types (includes result fields)
-- - Mapping table for adventure game character instances
-- - No separate result tables or template tables

-- Single turn sheet table for all game types
CREATE TABLE game_turn_sheet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES game(id),
    game_instance_id UUID NOT NULL REFERENCES game_instance(id),
    account_id UUID NOT NULL REFERENCES account(id),
    turn_number INTEGER NOT NULL,
    sheet_type VARCHAR(50) NOT NULL, -- 'location_choice', 'combat', 'inventory', etc.
    sheet_order INTEGER NOT NULL DEFAULT 1,
    sheet_data JSONB NOT NULL, -- Validated against sheet_type schema in code
    is_completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMPTZ,
    
    -- Result fields (NULL until scanned)
    result_data JSONB,
    scanned_at TIMESTAMPTZ,
    scanned_by UUID REFERENCES account(id),
    scan_quality DECIMAL(3,2) CHECK (scan_quality >= 0 AND scan_quality <= 1),
    processing_status VARCHAR(20) DEFAULT 'pending' CHECK (processing_status IN ('pending', 'processed', 'error')),
    error_message TEXT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- Mapping table for adventure game entities
CREATE TABLE adventure_game_turn_sheet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES game(id),
    adventure_game_character_instance_id UUID NOT NULL REFERENCES adventure_game_character_instance(id),
    game_turn_sheet_id UUID NOT NULL REFERENCES game_turn_sheet(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE(adventure_game_character_instance_id, game_turn_sheet_id)
);

-- Indexes for performance
CREATE INDEX idx_game_turn_sheet_game_id ON game_turn_sheet(game_id);
CREATE INDEX idx_game_turn_sheet_game_instance_account ON game_turn_sheet(game_instance_id, account_id);
CREATE INDEX idx_game_turn_sheet_turn_number ON game_turn_sheet(turn_number);
CREATE INDEX idx_game_turn_sheet_sheet_type ON game_turn_sheet(sheet_type);
CREATE INDEX idx_game_turn_sheet_processing_status ON game_turn_sheet(processing_status);
CREATE INDEX idx_adventure_game_turn_sheet_game_id ON adventure_game_turn_sheet(game_id);
CREATE INDEX idx_adventure_game_turn_sheet_character_instance ON adventure_game_turn_sheet(adventure_game_character_instance_id);
CREATE INDEX idx_adventure_game_turn_sheet_turn_sheet ON adventure_game_turn_sheet(game_turn_sheet_id);

-- Comments
COMMENT ON TABLE game_turn_sheet IS 'Single turn sheet table for all game types, includes both sheet data and scanned results';
COMMENT ON TABLE adventure_game_turn_sheet IS 'Mapping table linking adventure game entities (characters, locations, etc.) to turn sheets';
