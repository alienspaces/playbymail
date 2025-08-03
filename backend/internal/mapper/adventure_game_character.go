package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

// AdventureGameCharacterRequestToRecord maps a AdventureGameCharacterRequest to a adventure_game_record.AdventureGameCharacter
func AdventureGameCharacterRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameCharacter) (*adventure_game_record.AdventureGameCharacter, error) {
	l.Debug("mapping adventure_game_character request to record")

	var req adventure_game_schema.AdventureGameCharacterRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, fmt.Errorf("failed reading request: %w", err)
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.AccountID = req.AccountID
		rec.Name = req.Name
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.AccountID = req.AccountID
		rec.Name = req.Name
	}

	return rec, nil
}

// AdventureGameCharacterRecordToResponseData maps a adventure_game_record.AdventureGameCharacter to adventure_game_schema.AdventureGameCharacterResponseData
func AdventureGameCharacterRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCharacter) (adventure_game_schema.AdventureGameCharacterResponseData, error) {
	l.Debug("mapping adventure_game_character record to response data")
	data := adventure_game_schema.AdventureGameCharacterResponseData{
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

// AdventureGameCharacterRecordToResponse wraps response data in adventure_game_schema.AdventureGameCharacterResponse
func AdventureGameCharacterRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCharacter) (adventure_game_schema.AdventureGameCharacterResponse, error) {
	data, err := AdventureGameCharacterRecordToResponseData(l, rec)
	if err != nil {
		return adventure_game_schema.AdventureGameCharacterResponse{}, err
	}
	return adventure_game_schema.AdventureGameCharacterResponse{
		Data: &data,
	}, nil
}

func AdventureGameCharacterRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCharacter) (adventure_game_schema.AdventureGameCharacterCollectionResponse, error) {
	data := []*adventure_game_schema.AdventureGameCharacterResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameCharacterRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameCharacterCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return adventure_game_schema.AdventureGameCharacterCollectionResponse{
		Data: data,
	}, nil
}
