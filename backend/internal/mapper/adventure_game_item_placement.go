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

func AdventureGameItemPlacementRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameItemPlacement) (*adventure_game_record.AdventureGameItemPlacement, error) {
	l.Debug("mapping adventure_game_item_placement request to record")

	var req adventure_game_schema.AdventureGameItemPlacementRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.AdventureGameItemID = req.AdventureGameItemID
		rec.AdventureGameLocationID = req.AdventureGameLocationID
		rec.InitialCount = req.InitialCount
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.AdventureGameItemID = req.AdventureGameItemID
		rec.AdventureGameLocationID = req.AdventureGameLocationID
		rec.InitialCount = req.InitialCount
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AdventureGameItemPlacementRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameItemPlacement) (*adventure_game_schema.AdventureGameItemPlacementResponseData, error) {
	l.Debug("mapping adventure_game_item_placement record to response data")
	return &adventure_game_schema.AdventureGameItemPlacementResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		AdventureGameItemID:     rec.AdventureGameItemID,
		AdventureGameLocationID: rec.AdventureGameLocationID,
		InitialCount:            rec.InitialCount,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameItemPlacementRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameItemPlacement) (*adventure_game_schema.AdventureGameItemPlacementCollectionResponse, error) {
	l.Debug("mapping adventure_game_item_placement records to collection response")
	data := []*adventure_game_schema.AdventureGameItemPlacementResponseData{}
	for _, rec := range recs {
		item, err := AdventureGameItemPlacementRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return &adventure_game_schema.AdventureGameItemPlacementCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameItemPlacementRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameItemPlacement) (*adventure_game_schema.AdventureGameItemPlacementResponse, error) {
	l.Debug("mapping adventure_game_item_placement record to response")
	data, err := AdventureGameItemPlacementRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameItemPlacementResponse{
		Data: data,
	}, nil
}
