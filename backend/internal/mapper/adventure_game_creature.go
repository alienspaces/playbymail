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

func AdventureGameCreatureRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameCreature) (*adventure_game_record.AdventureGameCreature, error) {
	l.Debug("mapping adventure_game_creature request to record")

	var req adventure_game_schema.AdventureGameCreatureRequest
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

func AdventureGameCreatureRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCreature) (*adventure_game_schema.AdventureGameCreatureResponseData, error) {
	l.Debug("mapping adventure_game_creature record to response data")
	return &adventure_game_schema.AdventureGameCreatureResponseData{
		ID:          rec.ID,
		GameID:      rec.GameID,
		Name:        rec.Name,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:   nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameCreatureRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCreature) (*adventure_game_schema.AdventureGameCreatureResponse, error) {
	l.Debug("mapping adventure_game_creature record to response")
	data, err := AdventureGameCreatureRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameCreatureResponse{
		Data: data,
	}, nil
}

func AdventureGameCreatureRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCreature) (adventure_game_schema.AdventureGameCreatureCollectionResponse, error) {
	l.Debug("mapping adventure_game_creature records to collection response")
	data := []*adventure_game_schema.AdventureGameCreatureResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameCreatureRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameCreatureCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameCreatureCollectionResponse{
		Data: data,
	}, nil
}
