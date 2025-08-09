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

func AdventureGameItemRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameItem) (*adventure_game_record.AdventureGameItem, error) {
	l.Debug("mapping adventure_game_item request to record")

	var req adventure_game_schema.AdventureGameItemRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Name = req.Name
		rec.Description = req.Description
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.Description = req.Description
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AdventureGameItemRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameItem) (*adventure_game_schema.AdventureGameItemResponseData, error) {
	l.Debug("mapping adventure_game_item record to response data")
	return &adventure_game_schema.AdventureGameItemResponseData{
		ID:          rec.ID,
		GameID:      rec.GameID,
		Name:        rec.Name,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:   nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameItemRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameItem) (*adventure_game_schema.AdventureGameItemResponse, error) {
	l.Debug("mapping adventure_game_item record to response")
	data, err := AdventureGameItemRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameItemResponse{
		Data: data,
	}, nil
}

func AdventureGameItemRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameItem) (adventure_game_schema.AdventureGameItemCollectionResponse, error) {
	l.Debug("mapping adventure_game_item records to collection response")
	data := []*adventure_game_schema.AdventureGameItemResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameItemRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameItemCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameItemCollectionResponse{
		Data: data,
	}, nil
}
