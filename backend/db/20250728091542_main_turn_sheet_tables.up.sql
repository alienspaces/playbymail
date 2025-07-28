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

-- Player turn sheet responses (submitted answers)
CREATE TABLE IF NOT EXISTS player_turn_sheet_response (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    turn_sheet_id UUID NOT NULL REFERENCES player_turn_sheet(id),
    response_data JSONB NOT NULL, -- Player's answers
    submitted_at TIMESTAMPTZ DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    error_message TEXT
);

-- Turn sheet generation rules (when to generate which sheets)
CREATE TABLE IF NOT EXISTS turn_sheet_generation_rule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_instance_id UUID NOT NULL REFERENCES game_instance(id),
    rule_name VARCHAR(100) NOT NULL,
    trigger_condition JSONB NOT NULL, -- When to apply this rule
    sheet_type VARCHAR(50) NOT NULL,
    sheet_order INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
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

-- Add constraints for turn_sheet_generation_rule
ALTER TABLE turn_sheet_generation_rule ADD CONSTRAINT turn_sheet_generation_rule_sheet_type_check 
    CHECK (sheet_type IN ('multiple_choice', 'combat_selection', 'inventory_management', 'conversation', 'text_input', 'number_input'));

ALTER TABLE turn_sheet_generation_rule ADD CONSTRAINT turn_sheet_generation_rule_sheet_order_check 
    CHECK (sheet_order > 0);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_turn_sheet_template_game_type ON turn_sheet_template(game_type);
CREATE INDEX IF NOT EXISTS idx_turn_sheet_template_template_type ON turn_sheet_template(template_type);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_game_instance ON player_turn_sheet(game_instance_id);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_player ON player_turn_sheet(player_id);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_turn ON player_turn_sheet(turn_number);
CREATE INDEX IF NOT EXISTS idx_player_turn_sheet_completed ON player_turn_sheet(is_completed);
CREATE INDEX IF NOT EXISTS idx_turn_sheet_response_sheet ON player_turn_sheet_response(turn_sheet_id);
CREATE INDEX IF NOT EXISTS idx_turn_sheet_generation_rule_game_instance ON turn_sheet_generation_rule(game_instance_id);

-- Add comments
COMMENT ON TABLE turn_sheet_template IS 'Reusable turn sheet templates for different game types';
COMMENT ON TABLE player_turn_sheet IS 'Generated turn sheets for specific players and turns';
COMMENT ON TABLE player_turn_sheet_response IS 'Player responses to turn sheets';
COMMENT ON TABLE turn_sheet_generation_rule IS 'Rules for when to generate which turn sheets';
