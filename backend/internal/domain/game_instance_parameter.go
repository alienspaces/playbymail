package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) GetGameInstanceParameterRec(recID string, lock *sql.Lock) (*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateGameInstanceParameterRec(rec *game_record.GameInstanceParameter) (*game_record.GameInstanceParameter, error) {
	// Validate before creating
	if err := m.ValidateGameInstanceParameter(rec); err != nil {
		return nil, err
	}

	r := m.GameInstanceParameterRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateGameInstanceParameterRec(next *game_record.GameInstanceParameter) (*game_record.GameInstanceParameter, error) {
	// Validate before updating
	if err := m.ValidateGameInstanceParameter(next); err != nil {
		return nil, err
	}

	r := m.GameInstanceParameterRepository()
	rec, err := r.UpdateOne(next)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) DeleteGameInstanceParameterRec(recID string) error {
	r := m.GameInstanceParameterRepository()
	return r.DeleteOne(recID)
}

func (m *Domain) GetGameInstanceParameterRecs(opts *sql.Options) ([]*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (m *Domain) ValidateGameInstanceParameter(rec *game_record.GameInstanceParameter) error {
	// Validate required fields
	if rec.GameInstanceID == "" {
		return RequiredField("game_instance_id")
	}

	if rec.ParameterKey == "" {
		return RequiredField("parameter_key")
	}

	if rec.ParameterValue.String == "" {
		return RequiredField("parameter_value")
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

func (m *Domain) GetManyGameInstanceParameterRecs(opts *sql.Options) ([]*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

// GetGameInstanceParametersByGameInstanceID gets all parameters for a specific game instance
func (m *Domain) GetGameInstanceParameterRecsByGameInstanceID(gameInstanceID string) ([]*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
	opts := &sql.Options{
		Params: []sql.Param{
			{
				Col: "game_instance_id",
				Val: gameInstanceID,
			},
		},
	}
	return r.GetMany(opts)
}
