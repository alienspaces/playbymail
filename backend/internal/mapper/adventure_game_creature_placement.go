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

func AdventureGameCreaturePlacementRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameCreaturePlacement) (*adventure_game_record.AdventureGameCreaturePlacement, error) {
	l.Debug("mapping adventure_game_creature_placement request to record")

	var req adventure_game_schema.AdventureGameCreaturePlacementRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.AdventureGameCreatureID = req.AdventureGameCreatureID
		rec.AdventureGameLocationID = req.AdventureGameLocationID
		rec.InitialCount = req.InitialCount
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.AdventureGameCreatureID = req.AdventureGameCreatureID
		rec.AdventureGameLocationID = req.AdventureGameLocationID
		rec.InitialCount = req.InitialCount
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AdventureGameCreaturePlacementRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCreaturePlacement) (*adventure_game_schema.AdventureGameCreaturePlacementResponseData, error) {
	l.Debug("mapping adventure_game_creature_placement record to response data")
	return &adventure_game_schema.AdventureGameCreaturePlacementResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		AdventureGameCreatureID: rec.AdventureGameCreatureID,
		AdventureGameLocationID: rec.AdventureGameLocationID,
		InitialCount:            rec.InitialCount,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameCreaturePlacementRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCreaturePlacement) (*adventure_game_schema.AdventureGameCreaturePlacementCollectionResponse, error) {
	l.Debug("mapping adventure_game_creature_placement records to collection response")
	data := []*adventure_game_schema.AdventureGameCreaturePlacementResponseData{}
	for _, rec := range recs {
		item, err := AdventureGameCreaturePlacementRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return &adventure_game_schema.AdventureGameCreaturePlacementCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameCreaturePlacementRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCreaturePlacement) (*adventure_game_schema.AdventureGameCreaturePlacementResponse, error) {
	l.Debug("mapping adventure_game_creature_placement record to response")
	data, err := AdventureGameCreaturePlacementRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameCreaturePlacementResponse{
		Data: data,
	}, nil
}
