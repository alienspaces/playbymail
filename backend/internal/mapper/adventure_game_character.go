package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
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
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.AccountUserID = req.AccountUserID
		rec.Name = req.Name
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.AccountUserID = req.AccountUserID
		rec.Name = req.Name
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

// AdventureGameCharacterRecordToResponseData maps a adventure_game_record.AdventureGameCharacter to adventure_game_schema.AdventureGameCharacterResponseData
func AdventureGameCharacterRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCharacter) (*adventure_game_schema.AdventureGameCharacterResponseData, error) {
	l.Debug("mapping adventure_game_character record to response data")
	return &adventure_game_schema.AdventureGameCharacterResponseData{
		ID:        rec.ID,
		GameID:    rec.GameID,
		AccountID: rec.AccountUserID,
		Name:      rec.Name,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

// AdventureGameCharacterRecordToResponse wraps response data in adventure_game_schema.AdventureGameCharacterResponse
func AdventureGameCharacterRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCharacter) (*adventure_game_schema.AdventureGameCharacterResponse, error) {
	l.Debug("mapping adventure_game_character record to response")
	data, err := AdventureGameCharacterRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameCharacterResponse{
		Data: data,
	}, nil
}

func AdventureGameCharacterRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCharacter) (adventure_game_schema.AdventureGameCharacterCollectionResponse, error) {
	l.Debug("mapping adventure_game_character records to collection response")
	data := []*adventure_game_schema.AdventureGameCharacterResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameCharacterRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameCharacterCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameCharacterCollectionResponse{
		Data: data,
	}, nil
}
