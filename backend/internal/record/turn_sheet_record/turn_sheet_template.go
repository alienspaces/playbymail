package turn_sheet_record

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// TurnSheetTemplate
const (
	TableTurnSheetTemplate string = "turn_sheet_template"
)

const (
	FieldTurnSheetTemplateID        string = "id"
	FieldTurnSheetTemplateGameType  string = "game_type"
	FieldTurnSheetTemplateType      string = "template_type"
	FieldTurnSheetTemplateName      string = "template_name"
	FieldTurnSheetTemplateData      string = "template_data"
	FieldTurnSheetTemplateIsActive  string = "is_active"
	FieldTurnSheetTemplateCreatedAt string = "created_at"
	FieldTurnSheetTemplateUpdatedAt string = "updated_at"
	FieldTurnSheetTemplateDeletedAt string = "deleted_at"
)

type TurnSheetTemplate struct {
	record.Record
	GameType     string          `db:"game_type"`
	TemplateType string          `db:"template_type"`
	TemplateName string          `db:"template_name"`
	TemplateData json.RawMessage `db:"template_data"`
	IsActive     bool            `db:"is_active"`
	CreatedAt    sql.NullTime    `db:"created_at"`
	UpdatedAt    sql.NullTime    `db:"updated_at"`
	DeletedAt    sql.NullTime    `db:"deleted_at"`
}

func (r *TurnSheetTemplate) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldTurnSheetTemplateGameType] = r.GameType
	args[FieldTurnSheetTemplateType] = r.TemplateType
	args[FieldTurnSheetTemplateName] = r.TemplateName
	args[FieldTurnSheetTemplateData] = r.TemplateData
	args[FieldTurnSheetTemplateIsActive] = r.IsActive
	args[FieldTurnSheetTemplateCreatedAt] = r.CreatedAt
	args[FieldTurnSheetTemplateUpdatedAt] = r.UpdatedAt
	args[FieldTurnSheetTemplateDeletedAt] = r.DeletedAt
	return args
}
