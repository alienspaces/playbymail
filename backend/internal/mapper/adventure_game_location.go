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

// AdventureGameLocationRequestToRecord maps a request to a record for consistency
func AdventureGameLocationRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameLocation) (*adventure_game_record.AdventureGameLocation, error) {
	l.Debug("mapping adventure_game_location request to record")

	var req adventure_game_schema.AdventureGameLocationRequest
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

func AdventureGameLocationRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocation) (*adventure_game_schema.AdventureGameLocationResponseData, error) {
	l.Debug("mapping adventure_game_location record to response data")
	return &adventure_game_schema.AdventureGameLocationResponseData{
		ID:          rec.ID,
		GameID:      rec.GameID,
		Name:        rec.Name,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:   nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameLocationRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocation) (*adventure_game_schema.AdventureGameLocationResponse, error) {
	l.Debug("mapping adventure_game_location record to response")
	data, err := AdventureGameLocationRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameLocationResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocation) (adventure_game_schema.AdventureGameLocationCollectionResponse, error) {
	l.Debug("mapping adventure_game_location records to collection response")
	data := []*adventure_game_schema.AdventureGameLocationResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameLocationCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameLocationCollectionResponse{
		Data: data,
	}, nil
}
