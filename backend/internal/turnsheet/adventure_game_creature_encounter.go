package turnsheet

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/scanner"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

// MonsterEncounterData represents the data structure for monster encounter turn sheets.
type MonsterEncounterData struct {
	TurnSheetTemplateData

	// Character combat status
	CharacterName      string          `json:"character_name"`
	CharacterHealth    int             `json:"character_health"`
	CharacterMaxHealth int             `json:"character_max_health"`
	CharacterAttack    int             `json:"character_attack"`
	CharacterDefense   int             `json:"character_defense"`
	EquippedWeapon     *EquippedWeapon `json:"equipped_weapon,omitempty"`
	EquippedArmor      *EquippedArmor  `json:"equipped_armor,omitempty"`

	// Creatures in this encounter
	Creatures []EncounterCreature `json:"creatures"`

	// Maximum combat actions allowed this turn (0 = read-only, no action slots shown)
	MaxActions int `json:"max_actions"`

	// ReadOnly disables action slots (used for dead body and flee result sheets)
	ReadOnly bool `json:"read_only,omitempty"`
}

// EquippedWeapon represents the character's currently equipped weapon (display only).
type EquippedWeapon struct {
	ItemInstanceID string `json:"item_instance_id"`
	Name           string `json:"name"`
	Damage         int    `json:"damage"`
}

// EquippedArmor represents the character's currently equipped armor (display only).
type EquippedArmor struct {
	ItemInstanceID string `json:"item_instance_id"`
	Name           string `json:"name"`
	Defense        int    `json:"defense"`
}

// EncounterCreature represents a single creature in a monster encounter.
type EncounterCreature struct {
	CreatureInstanceID string  `json:"creature_instance_id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Health             int     `json:"health"`
	MaxHealth          int     `json:"max_health"`
	AttackDamage       int     `json:"attack_damage"`
	Defense            int     `json:"defense"`
	Disposition        string  `json:"disposition"` // "aggressive", "inquisitive", "indifferent"
	ImageDataURL       *string `json:"image_data_url,omitempty"`
	IsDead             bool    `json:"is_dead,omitempty"` // true for dead body display
}

// MonsterEncounterScanData represents the scanned data from a monster encounter turn sheet.
// Accepts both the structured format {"actions":[...]} produced by the frontend and the OCR
// scanner, and the flat HTML form format {"action_0":"attack","action_0_target":"<id>",...}.
type MonsterEncounterScanData struct {
	Actions    []CombatAction         `json:"actions"`
	FlatFields map[string]interface{} `json:"-"` // retains raw flat keys for GetActions normalisation
}

// GetActions returns the combat actions, normalising both the structured and flat form formats.
func (d *MonsterEncounterScanData) GetActions() []CombatAction {
	if len(d.Actions) > 0 {
		return d.Actions
	}
	if len(d.FlatFields) == 0 {
		return nil
	}
	// Convert flat action_N / action_N_target fields into CombatAction structs.
	var actions []CombatAction
	for i := 0; ; i++ {
		key := fmt.Sprintf("action_%d", i)
		val, ok := d.FlatFields[key]
		if !ok {
			break
		}
		actionType, _ := val.(string)
		action := CombatAction{ActionType: actionType}
		if actionType == CombatActionTypeAttack {
			targetKey := fmt.Sprintf("action_%d_target", i)
			if t, ok := d.FlatFields[targetKey]; ok {
				action.TargetCreatureInstanceID, _ = t.(string)
			}
		}
		actions = append(actions, action)
	}
	return actions
}

// UnmarshalJSON populates both the structured Actions field and the raw FlatFields map
// so that GetActions() can normalise either format.
func (d *MonsterEncounterScanData) UnmarshalJSON(data []byte) error {
	type alias MonsterEncounterScanData
	var a struct {
		alias
	}
	// Decode known structured fields.
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	d.Actions = a.Actions

	// Decode all fields into a raw map to capture flat action_N keys.
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.FlatFields = raw
	return nil
}

// CombatAction represents a single combat action chosen by the player.
// Combat action type constants.
const (
	CombatActionTypeDoNothing = "do_nothing"
	CombatActionTypeAttack    = "attack"
)

// CombatAction represents a single combat action chosen by the player.
type CombatAction struct {
	ActionType               string `json:"action_type"`                           // CombatActionTypeDoNothing or CombatActionTypeAttack
	TargetCreatureInstanceID string `json:"target_creature_instance_id,omitempty"` // required when ActionType is CombatActionTypeAttack
}

// MonsterEncounterScannedDataSchemaName is the filename of the JSON schema for monster encounter scanned_data.
const MonsterEncounterScannedDataSchemaName = "monster_encounter.schema.json"

const monsterEncounterTemplatePath = "turnsheet/adventure_game_monster_encounter.template"

// MonsterEncounterProcessor implements the DocumentProcessor interface for monster encounter turn sheets.
type MonsterEncounterProcessor struct {
	*BaseProcessor
}

// NewMonsterEncounterProcessor creates a new monster encounter processor.
func NewMonsterEncounterProcessor(l logger.Logger, cfg config.Config) (*MonsterEncounterProcessor, error) {
	baseProcessor, err := NewBaseProcessor(l, cfg)
	if err != nil {
		return nil, err
	}
	return &MonsterEncounterProcessor{
		BaseProcessor: baseProcessor,
	}, nil
}

// GeneratePreviewData generates dummy data for a monster encounter turn sheet preview.
func (p *MonsterEncounterProcessor) GeneratePreviewData(ctx context.Context, l logger.Logger, gameRec *game_record.Game, backgroundImage *string) ([]byte, error) {
	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	weaponDamage := 8
	armorDefense := 3
	data := MonsterEncounterData{
		TurnSheetTemplateData: TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr(gameRec.GameType),
			TurnNumber:            convert.Ptr(5),
			TurnSheetTitle:        convert.Ptr("Creature Encounter"),
			TurnSheetDescription:  convert.Ptr("You are not alone here."),
			TurnSheetInstructions: convert.Ptr("Choose your combat actions below."),
			TurnSheetCode:         convert.Ptr(turnSheetCode),
			BackgroundImage:       backgroundImage,
		},
		CharacterName:      "Aldric",
		CharacterHealth:    65,
		CharacterMaxHealth: 100,
		CharacterAttack:    weaponDamage,
		CharacterDefense:   armorDefense,
		EquippedWeapon: &EquippedWeapon{
			ItemInstanceID: "preview-weapon",
			Name:           "Iron Sword",
			Damage:         weaponDamage,
		},
		EquippedArmor: &EquippedArmor{
			ItemInstanceID: "preview-armor",
			Name:           "Leather Jerkin",
			Defense:        armorDefense,
		},
		Creatures: []EncounterCreature{
			{
				CreatureInstanceID: "preview-creature-1",
				Name:               "Shadow Monk",
				Description:        "A spectral figure in a hooded robe drifts silently. Cold radiates from it like a winter wind.",
				Health:             80,
				MaxHealth:          100,
				AttackDamage:       15,
				Defense:            5,
				Disposition:        "aggressive",
			},
			{
				CreatureInstanceID: "preview-creature-2",
				Name:               "Cellar Rat",
				Description:        "A large grey rat watches from the shadows between the wine racks, unafraid.",
				Health:             20,
				MaxHealth:          50,
				AttackDamage:       5,
				Defense:            0,
				Disposition:        "inquisitive",
			},
		},
		MaxActions: 3,
	}

	return json.Marshal(data)
}

// GenerateTurnSheet generates a monster encounter turn sheet document.
func (p *MonsterEncounterProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, sheetData []byte) ([]byte, error) {
	var data MonsterEncounterData
	if err := json.Unmarshal(sheetData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse sheet data: %w", err)
	}

	if err := p.ValidateBaseTemplateData(&data.TurnSheetTemplateData); err != nil {
		return nil, fmt.Errorf("template data validation failed: %w", err)
	}

	if data.TurnSheetInstructions == nil || strings.TrimSpace(*data.TurnSheetInstructions) == "" {
		instr := buildMonsterEncounterInstructions(data.MaxActions)
		data.TurnSheetInstructions = &instr
	}

	if data.TurnSheetTitle == nil || strings.TrimSpace(*data.TurnSheetTitle) == "" {
		title := "Creature Encounter"
		data.TurnSheetTitle = &title
	}

	if data.MaxActions == 0 {
		data.MaxActions = 3
	}

	if data.CharacterName == "" {
		return nil, fmt.Errorf("character name is required")
	}

	if len(data.Creatures) == 0 && !data.ReadOnly {
		return nil, fmt.Errorf("at least one creature is required")
	}

	return p.GenerateDocument(ctx, format, monsterEncounterTemplatePath, &data)
}

// ScanTurnSheet scans a monster encounter turn sheet and extracts combat actions using OCR.
func (p *MonsterEncounterProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	if len(imageData) == 0 {
		return nil, fmt.Errorf("empty image data provided")
	}

	var data MonsterEncounterData
	if err := json.Unmarshal(sheetData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse sheet data: %w", err)
	}

	templateImage, err := p.renderTemplatePreview(ctx, monsterEncounterTemplatePath, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate template preview: %w", err)
	}

	req := scanner.StructuredScanRequest{
		Instructions:      buildMonsterEncounterScanInstructions(),
		AdditionalContext: buildMonsterEncounterContext(&data),
		TemplateImage:     templateImage,
		TemplateImageMIME: "image/png",
		FilledImage:       imageData,
		ExpectedJSONSchema: map[string]any{
			"actions": []any{},
		},
	}

	raw, err := p.Scanner.ExtractStructuredData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	var scanData MonsterEncounterScanData
	if err := json.Unmarshal(raw, &scanData); err != nil {
		return nil, fmt.Errorf("failed to decode structured response: %w", err)
	}

	if err := validateMonsterEncounterScanData(&data, &scanData); err != nil {
		return nil, err
	}

	return json.Marshal(scanData)
}

func buildMonsterEncounterInstructions(maxActions int) string {
	if maxActions == 0 {
		maxActions = 3
	}
	return fmt.Sprintf(
		"Choose up to %d combat actions below. For each slot select Do Nothing or Attack and write the creature name.\n\nIf you take any actions on your Inventory Management sheet this turn, your combat actions will be forfeited.",
		maxActions,
	)
}

func buildMonsterEncounterScanInstructions() string {
	return `Compare the blank template with the completed turn sheet.
For each action slot, determine whether the player selected "Do Nothing" or "Attack".
When Attack is selected, identify the target creature name and map it to the correct creature_instance_id.
Respond with JSON containing an "actions" array of objects with "action_type" ("do_nothing" or "attack") and optionally "target_creature_instance_id".
Omit trailing do_nothing slots.`
}

func buildMonsterEncounterContext(data *MonsterEncounterData) []string {
	var ctx []string
	if data != nil {
		for _, c := range data.Creatures {
			ctx = append(ctx, fmt.Sprintf("creature_instance_id=%s name=%s", c.CreatureInstanceID, c.Name))
		}
	}
	return ctx
}

func validateMonsterEncounterScanData(sheetData *MonsterEncounterData, scanData *MonsterEncounterScanData) error {
	if scanData == nil {
		return fmt.Errorf("no scan data provided")
	}

	maxActions := sheetData.MaxActions
	if maxActions == 0 {
		maxActions = 3
	}

	if len(scanData.Actions) > maxActions {
		return fmt.Errorf("too many actions: got %d, max %d", len(scanData.Actions), maxActions)
	}

	// Build valid creature instance IDs
	validCreatureIDs := make(map[string]bool)
	for _, c := range sheetData.Creatures {
		validCreatureIDs[c.CreatureInstanceID] = true
	}

	for i, action := range scanData.Actions {
		switch action.ActionType {
		case CombatActionTypeDoNothing:
			// always valid
		case CombatActionTypeAttack:
			if action.TargetCreatureInstanceID == "" {
				return fmt.Errorf("action %d: attack requires target_creature_instance_id", i+1)
			}
			if !validCreatureIDs[action.TargetCreatureInstanceID] {
				return fmt.Errorf("action %d: invalid target creature instance id: %s", i+1, action.TargetCreatureInstanceID)
			}
		default:
			return fmt.Errorf("action %d: unknown action_type: %s", i+1, action.ActionType)
		}
	}

	return nil
}
