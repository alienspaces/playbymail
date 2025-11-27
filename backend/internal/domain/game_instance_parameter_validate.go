package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) validateGameInstanceParameterRecForCreate(rec *game_record.GameInstanceParameter) error {
	return validateGameInstanceParameterRec(rec, false)
}

func (m *Domain) validateGameInstanceParameterRecForUpdate(rec *game_record.GameInstanceParameter) error {
	return validateGameInstanceParameterRec(rec, true)
}

func validateGameInstanceParameterRec(rec *game_record.GameInstanceParameter, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(game_record.FieldGameInstanceParameterID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(game_record.FieldGameInstanceParameterGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(game_record.FieldGameInstanceParameterParameterKey, rec.ParameterKey); err != nil {
		return err
	}

	if err := domain.ValidateNullStringField(game_record.FieldGameInstanceParameterParameterValue, rec.ParameterValue); err != nil {
		return err
	}

	return nil
}

// ValidateGameInstanceParameter is a public method that performs comprehensive validation
// including business logic checks. This is kept for backward compatibility.
func (m *Domain) ValidateGameInstanceParameter(rec *game_record.GameInstanceParameter) error {
	// Basic field validation
	if err := validateGameInstanceParameterRec(rec, false); err != nil {
		return err
	}

	// Validate that the game instance exists
	gameInstance, err := m.GetGameInstanceRec(rec.GameInstanceID, nil)
	if err != nil {
		return NotFound("game instance", rec.GameInstanceID)
	}

	// Get the game to find its type
	game, err := m.GetGameRec(gameInstance.GameID, nil)
	if err != nil {
		return NotFound("game", gameInstance.GameID)
	}

	// Validate that the parameter key is valid for this game type
	gameParameters := GetGameParametersByGameType(game.GameType)
	validParameter := false
	var expectedValueType string

	for _, gp := range gameParameters {
		if gp.ConfigKey == rec.ParameterKey {
			validParameter = true
			expectedValueType = gp.ValueType
			break
		}
	}

	if !validParameter {
		return InvalidField(game_record.FieldGameInstanceParameterParameterKey, rec.ParameterKey, "parameter key is not valid for game type")
	}

	// Validate that the parameter value matches the expected type
	if err := validateParameterValue(rec.ParameterValue.String, expectedValueType); err != nil {
		return InvalidField(game_record.FieldGameInstanceParameterParameterValue, rec.ParameterValue.String, err.Error())
	}

	return nil
}

// validateParameterValue validates that a parameter value matches its expected type
func validateParameterValue(value, valueType string) error {
	switch valueType {
	case "string":
		// String values are always valid
		return nil
	case "integer":
		// Check if the string can be parsed as an integer
		if _, err := parseInt(value); err != nil {
			return fmt.Errorf("value '%s' is not a valid integer", value)
		}
		return nil
	case "boolean":
		// Check if the string represents a valid boolean
		if value != "true" && value != "false" {
			return fmt.Errorf("value '%s' is not a valid boolean (must be 'true' or 'false')", value)
		}
		return nil
	default:
		return fmt.Errorf("unknown value type '%s'", valueType)
	}
}

// parseInt is a helper function to parse integers
func parseInt(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
