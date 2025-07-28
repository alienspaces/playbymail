package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

// AdventureGameCharacterRequestToRecord maps a AdventureGameCharacterRequest to a adventure_game_record.AdventureGameCharacter
func AdventureGameCharacterRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameCharacter) (*adventure_game_record.AdventureGameCharacter, error) {
	l.Debug("mapping adventure_game_character request to record")

	var req api.AdventureGameCharacterRequest
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

// AdventureGameCharacterRecordToResponseData maps a adventure_game_record.AdventureGameCharacter to api.AdventureGameCharacterResponseData
func AdventureGameCharacterRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCharacter) (api.AdventureGameCharacterResponseData, error) {
	l.Debug("mapping adventure_game_character record to response data")
	data := api.AdventureGameCharacterResponseData{
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

// AdventureGameCharacterRecordToResponse wraps response data in api.AdventureGameCharacterResponse
func AdventureGameCharacterRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCharacter) (api.AdventureGameCharacterResponse, error) {
	data, err := AdventureGameCharacterRecordToResponseData(l, rec)
	if err != nil {
		return api.AdventureGameCharacterResponse{}, err
	}
	return api.AdventureGameCharacterResponse{
		Data: &data,
	}, nil
}

func AdventureGameCharacterRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCharacter) (api.AdventureGameCharacterCollectionResponse, error) {
	data := []*api.AdventureGameCharacterResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameCharacterRecordToResponseData(l, rec)
		if err != nil {
			return api.AdventureGameCharacterCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AdventureGameCharacterCollectionResponse{
		Data: data,
	}, nil
}
