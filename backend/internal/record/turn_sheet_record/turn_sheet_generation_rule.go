package turn_sheet_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// TurnSheetGenerationRule
const (
	TableTurnSheetGenerationRule string = "turn_sheet_generation_rule"
)

const (
	FieldTurnSheetGenerationRuleID             string = "id"
	FieldTurnSheetGenerationRuleGameInstanceID string = "game_instance_id"
	FieldTurnSheetGenerationRuleRuleName       string = "rule_name"
	FieldTurnSheetGenerationRuleTriggerType    string = "trigger_type"
	FieldTurnSheetGenerationRuleSheetType      string = "sheet_type"
	FieldTurnSheetGenerationRuleSheetOrder     string = "sheet_order"
	FieldTurnSheetGenerationRuleIsActive       string = "is_active"
	FieldTurnSheetGenerationRuleCreatedAt      string = "created_at"
	FieldTurnSheetGenerationRuleUpdatedAt      string = "updated_at"
	FieldTurnSheetGenerationRuleDeletedAt      string = "deleted_at"
)

type TurnSheetGenerationRule struct {
	record.Record
	GameInstanceID string       `db:"game_instance_id"`
	RuleName       string       `db:"rule_name"`
	TriggerType    string       `db:"trigger_type"`
	SheetType      string       `db:"sheet_type"`
	SheetOrder     int          `db:"sheet_order"`
	IsActive       bool         `db:"is_active"`
	CreatedAt      sql.NullTime `db:"created_at"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
	DeletedAt      sql.NullTime `db:"deleted_at"`
}

func (r *TurnSheetGenerationRule) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldTurnSheetGenerationRuleGameInstanceID] = r.GameInstanceID
	args[FieldTurnSheetGenerationRuleRuleName] = r.RuleName
	args[FieldTurnSheetGenerationRuleTriggerType] = r.TriggerType
	args[FieldTurnSheetGenerationRuleSheetType] = r.SheetType
	args[FieldTurnSheetGenerationRuleSheetOrder] = r.SheetOrder
	args[FieldTurnSheetGenerationRuleIsActive] = r.IsActive
	args[FieldTurnSheetGenerationRuleCreatedAt] = r.CreatedAt
	args[FieldTurnSheetGenerationRuleUpdatedAt] = r.UpdatedAt
	args[FieldTurnSheetGenerationRuleDeletedAt] = r.DeletedAt
	return args
}
