package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

// GameCharacterRequestToRecord maps a GameCharacterRequest to a record.GameCharacter
func GameCharacterRequestToRecord(l logger.Logger, r *http.Request, rec *record.GameCharacter) (*record.GameCharacter, error) {
	l.Debug("mapping game_character request to record")

	var req schema.GameCharacterRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, fmt.Errorf("failed reading request: %w", err)
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameID = req.GameID
		rec.AccountID = req.AccountID
		rec.Name = req.Name
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.GameID = req.GameID
		rec.AccountID = req.AccountID
		rec.Name = req.Name
	}

	return rec, nil
}

// GameCharacterRecordToResponseData maps a record.GameCharacter to schema.GameCharacterResponseData
func GameCharacterRecordToResponseData(l logger.Logger, rec *record.GameCharacter) (schema.GameCharacterResponseData, error) {
	l.Debug("mapping game_character record to response data")
	data := schema.GameCharacterResponseData{
		ID:        rec.ID,
		GameID:    rec.GameID,
		AccountID: rec.AccountID,
		Name:      rec.Name,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nil,
		DeletedAt: nil,
	}
	if rec.UpdatedAt.Valid {
		data.UpdatedAt = &rec.UpdatedAt.Time
	}
	if rec.DeletedAt.Valid {
		data.DeletedAt = &rec.DeletedAt.Time
	}
	return data, nil
}

// GameCharacterRecordToResponse wraps response data in schema.GameCharacterResponse
func GameCharacterRecordToResponse(l logger.Logger, rec *record.GameCharacter) (schema.GameCharacterResponse, error) {
	data, err := GameCharacterRecordToResponseData(l, rec)
	if err != nil {
		return schema.GameCharacterResponse{}, err
	}
	return schema.GameCharacterResponse{
		Response:                  schema.Response{},
		GameCharacterResponseData: &data,
	}, nil
}
