package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func AdventureGameLocationLinkRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_record.AdventureGameLocationLink, error) {
	l.Debug("mapping adventure_game_location_link request to record")

	var req api.AdventureGameLocationLinkRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.FromAdventureGameLocationID = req.FromGameLocationID
		rec.ToAdventureGameLocationID = req.ToGameLocationID
		rec.Description = req.Description
		rec.Name = req.Name
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AdventureGameLocationLinkRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLink) (api.AdventureGameLocationLinkResponseData, error) {
	l.Debug("mapping adventure_game_location_link record to response data")
	data := api.AdventureGameLocationLinkResponseData{
		ID:                 rec.ID,
		GameID:             rec.GameID,
		FromGameLocationID: rec.FromAdventureGameLocationID,
		ToGameLocationID:   rec.ToAdventureGameLocationID,
		Description:        rec.Description,
		Name:               rec.Name,
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:          nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameLocationLinkRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLink) (api.AdventureGameLocationLinkResponse, error) {
	data, err := AdventureGameLocationLinkRecordToResponseData(l, rec)
	if err != nil {
		return api.AdventureGameLocationLinkResponse{}, err
	}
	return api.AdventureGameLocationLinkResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationLinkRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationLink) (api.AdventureGameLocationLinkCollectionResponse, error) {
	data := []*api.AdventureGameLocationLinkResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationLinkRecordToResponseData(l, rec)
		if err != nil {
			return api.AdventureGameLocationLinkCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AdventureGameLocationLinkCollectionResponse{
		Data: data,
	}, nil
}
