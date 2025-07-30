-- Create turn sheet tables

-- Turn sheet templates (reusable templates for different game types)
CREATE TABLE IF NOT EXISTS turn_sheet_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_type VARCHAR(50) NOT NULL, -- 'adventure', 'strategy', etc.
    template_type VARCHAR(50) NOT NULL, -- 'multiple_choice', 'inventory', 'combat', 'conversation'
    template_name VARCHAR(100) NOT NULL,
    template_data JSONB NOT NULL, -- Form structure, fields, validation rules
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- Player turn sheets (generated for specific players/turns)
CREATE TABLE IF NOT EXISTS player_turn_sheet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_instance_id UUID NOT NULL REFERENCES game_instance(id),
    player_id UUID NOT NULL REFERENCES account(id),
    turn_number INTEGER NOT NULL,
    sheet_type VARCHAR(50) NOT NULL, -- 'multiple_choice', 'inventory', 'combat', etc.
    sheet_order INTEGER NOT NULL, -- Order within the turn (1, 2, 3...)
    sheet_data JSONB NOT NULL, -- Personalized sheet content
    is_required BOOLEAN DEFAULT true,
    is_completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- Player turn sheet results (scanned from physical turn sheets)
CREATE TABLE IF NOT EXISTS player_turn_sheet_result (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    turn_sheet_id UUID NOT NULL REFERENCES player_turn_sheet(id),
    result_data JSONB NOT NULL, -- Scanned and processed player responses
    scanned_at TIMESTAMPTZ DEFAULT NOW(),
    scanned_by UUID REFERENCES account(id), -- Who scanned the sheet
    scan_quality DECIMAL(3,2), -- OCR confidence score (0.00-1.00)
    processing_status VARCHAR(20) DEFAULT 'scanned', -- 'scanned', 'processed', 'error'
    error_message TEXT,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- Add constraints for turn_sheet_template
ALTER TABLE turn_sheet_template ADD CONSTRAINT turn_sheet_template_game_type_check 
    CHECK (game_type IN ('adventure', 'strategy', 'puzzle', 'simulation'));

ALTER TABLE turn_sheet_template ADD CONSTRAINT turn_sheet_template_template_type_check 
    CHECK (template_type IN ('multiple_choice', 'combat_selection', 'inventory_management', 'conversation', 'text_input', 'number_input'));

ALTER TABLE turn_sheet_template ADD CONSTRAINT turn_sheet_template_unique_name_per_type 
    UNIQUE (game_type, template_name);

-- Add constraints for player_turn_sheet
ALTER TABLE player_turn_sheet ADD CONSTRAINT player_turn_sheet_sheet_type_check 
    CHECK (sheet_type IN ('multiple_choice', 'combat_selection', 'inventory_management', 'conversation', 'text_input', 'number_input'));

ALTER TABLE player_turn_sheet ADD CONSTRAINT player_turn_sheet_sheet_order_check 
    CHECK (sheet_order > 0);

ALTER TABLE player_turn_sheet ADD CONSTRAINT player_turn_sheet_turn_number_check 
    CHECK (turn_number >= 0);

-- Add constraints for player_turn_sheet_result
ALTER TABLE player_turn_sheet_result ADD CONSTRAINT player_turn_sheet_result_processing_status_check 
    CHECK (processing_status IN ('scanned', 'processed', 'error'));

ALTER TABLE player_turn_sheet_result ADD CONSTRAINT player_turn_sheet_result_scan_quality_check 
    CHECK (scan_quality >= 0.00 AND scan_quality <= 1.00);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_turn_sheet_template_game_type ON turn_sheet_template(game_type);
CREATE INDEX IF NOT EXISTS idx_turn_sheet_template_template_type ON turn_sheet_template(template_type);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_game_instance ON player_turn_sheet(game_instance_id);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_player ON player_turn_sheet(player_id);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_turn ON player_turn_sheet(turn_number);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_completed ON player_turn_sheet(is_completed);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_result_sheet ON player_turn_sheet_result(turn_sheet_id);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_result_status ON player_turn_sheet_result(processing_status);

-- Add comments
COMMENT ON TABLE turn_sheet_template IS 'Reusable turn sheet templates for different game types';
COMMENT ON TABLE player_turn_sheet IS 'Generated turn sheets for specific players and turns';
COMMENT ON TABLE player_turn_sheet_result IS 'Scanned and processed player responses from physical turn sheets';
