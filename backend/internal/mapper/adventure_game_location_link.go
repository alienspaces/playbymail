package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AdventureGameLocationLinkRequestToRecord(l logger.Logger, r *http.Request, rec *record.AdventureGameLocationLink) (*record.AdventureGameLocationLink, error) {
	l.Debug("mapping adventure_game_location_link request to record")

	var req schema.AdventureGameLocationLinkRequest
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

func AdventureGameLocationLinkRecordToResponseData(l logger.Logger, rec *record.AdventureGameLocationLink) (schema.AdventureGameLocationLinkResponseData, error) {
	l.Debug("mapping adventure_game_location_link record to response data")
	data := schema.AdventureGameLocationLinkResponseData{
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

func AdventureGameLocationLinkRecordToResponse(l logger.Logger, rec *record.AdventureGameLocationLink) (schema.AdventureGameLocationLinkResponse, error) {
	data, err := AdventureGameLocationLinkRecordToResponseData(l, rec)
	if err != nil {
		return schema.AdventureGameLocationLinkResponse{}, err
	}
	return schema.AdventureGameLocationLinkResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationLinkRecordsToCollectionResponse(l logger.Logger, recs []*record.AdventureGameLocationLink) (schema.AdventureGameLocationLinkCollectionResponse, error) {
	var data []*schema.AdventureGameLocationLinkResponseData
	for _, rec := range recs {
		d, err := AdventureGameLocationLinkRecordToResponseData(l, rec)
		if err != nil {
			return schema.AdventureGameLocationLinkCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.AdventureGameLocationLinkCollectionResponse{
		Data: data,
	}, nil
}
