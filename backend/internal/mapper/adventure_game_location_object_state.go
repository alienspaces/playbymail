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

// AdventureGameLocationObjectStateRequestToRecord maps a request to a record
func AdventureGameLocationObjectStateRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameLocationObjectState) (*adventure_game_record.AdventureGameLocationObjectState, error) {
	l.Debug("mapping adventure_game_location_object_state request to record")

	var req adventure_game_schema.AdventureGameLocationObjectStateRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Name = req.Name
		rec.Description = req.Description
		rec.SortOrder = req.SortOrder
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.Description = req.Description
		rec.SortOrder = req.SortOrder
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AdventureGameLocationObjectStateRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationObjectState) (*adventure_game_schema.AdventureGameLocationObjectStateResponseData, error) {
	l.Debug("mapping adventure_game_location_object_state record to response data")
	return &adventure_game_schema.AdventureGameLocationObjectStateResponseData{
		ID:                            rec.ID,
		GameID:                        rec.GameID,
		AdventureGameLocationObjectID: rec.AdventureGameLocationObjectID,
		Name:                          rec.Name,
		Description:                   rec.Description,
		SortOrder:                     rec.SortOrder,
		CreatedAt:                     rec.CreatedAt,
		UpdatedAt:                     nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:                     nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameLocationObjectStateRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationObjectState) (*adventure_game_schema.AdventureGameLocationObjectStateResponse, error) {
	l.Debug("mapping adventure_game_location_object_state record to response")
	data, err := AdventureGameLocationObjectStateRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameLocationObjectStateResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationObjectStateRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationObjectState) (adventure_game_schema.AdventureGameLocationObjectStateCollectionResponse, error) {
	l.Debug("mapping adventure_game_location_object_state records to collection response")
	data := []*adventure_game_schema.AdventureGameLocationObjectStateResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationObjectStateRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameLocationObjectStateCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameLocationObjectStateCollectionResponse{
		Data: data,
	}, nil
}
