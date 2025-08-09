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

// AdventureGameLocationLinkRequestToRecord maps a request to a record for consistency
func AdventureGameLocationLinkRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_record.AdventureGameLocationLink, error) {
	l.Debug("mapping adventure_game_location_link request to record")

	var req adventure_game_schema.AdventureGameLocationLinkRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Name = req.Name
		rec.Description = req.Description
		rec.FromAdventureGameLocationID = req.FromAdventureGameLocationID
		rec.ToAdventureGameLocationID = req.ToAdventureGameLocationID
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.Description = req.Description
		rec.FromAdventureGameLocationID = req.FromAdventureGameLocationID
		rec.ToAdventureGameLocationID = req.ToAdventureGameLocationID
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AdventureGameLocationLinkRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_schema.AdventureGameLocationLinkResponseData, error) {
	l.Debug("mapping adventure_game_location_link record to response data")
	return &adventure_game_schema.AdventureGameLocationLinkResponseData{
		ID:                          rec.ID,
		GameID:                      rec.GameID,
		Name:                        rec.Name,
		Description:                 rec.Description,
		FromAdventureGameLocationID: rec.FromAdventureGameLocationID, // Map old field name to new
		ToAdventureGameLocationID:   rec.ToAdventureGameLocationID,   // Map old field name to new
		CreatedAt:                   rec.CreatedAt,
		UpdatedAt:                   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:                   nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameLocationLinkRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationLink) (*adventure_game_schema.AdventureGameLocationLinkCollectionResponse, error) {
	l.Debug("mapping adventure_game_location_link records to collection response")
	data := []*adventure_game_schema.AdventureGameLocationLinkResponseData{}
	for _, rec := range recs {
		item, err := AdventureGameLocationLinkRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return &adventure_game_schema.AdventureGameLocationLinkCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationLinkRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_schema.AdventureGameLocationLinkResponse, error) {
	l.Debug("mapping adventure_game_location_link record to response")
	data, err := AdventureGameLocationLinkRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameLocationLinkResponse{
		Data: data,
	}, nil
}
