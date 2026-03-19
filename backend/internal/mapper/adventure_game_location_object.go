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

// AdventureGameLocationObjectRequestToRecord maps a request to a record
func AdventureGameLocationObjectRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameLocationObject) (*adventure_game_record.AdventureGameLocationObject, error) {
	l.Debug("mapping adventure_game_location_object request to record")

	var req adventure_game_schema.AdventureGameLocationObjectRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.AdventureGameLocationID = req.AdventureGameLocationID
		rec.Name = req.Name
		rec.Description = req.Description
		if req.InitialState != "" {
			rec.InitialState = req.InitialState
		} else {
			rec.InitialState = "intact"
		}
		rec.IsHidden = req.IsHidden
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.AdventureGameLocationID = req.AdventureGameLocationID
		rec.Name = req.Name
		rec.Description = req.Description
		if req.InitialState != "" {
			rec.InitialState = req.InitialState
		}
		rec.IsHidden = req.IsHidden
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AdventureGameLocationObjectRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationObject) (*adventure_game_schema.AdventureGameLocationObjectResponseData, error) {
	l.Debug("mapping adventure_game_location_object record to response data")
	return &adventure_game_schema.AdventureGameLocationObjectResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		AdventureGameLocationID: rec.AdventureGameLocationID,
		Name:                    rec.Name,
		Description:             rec.Description,
		InitialState:            rec.InitialState,
		IsHidden:                rec.IsHidden,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameLocationObjectRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationObject) (*adventure_game_schema.AdventureGameLocationObjectResponse, error) {
	l.Debug("mapping adventure_game_location_object record to response")
	data, err := AdventureGameLocationObjectRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameLocationObjectResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationObjectRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationObject) (adventure_game_schema.AdventureGameLocationObjectCollectionResponse, error) {
	l.Debug("mapping adventure_game_location_object records to collection response")
	data := []*adventure_game_schema.AdventureGameLocationObjectResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationObjectRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameLocationObjectCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameLocationObjectCollectionResponse{
		Data: data,
	}, nil
}
