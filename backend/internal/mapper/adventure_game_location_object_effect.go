package mapper

import (
	"database/sql"
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameLocationObjectEffectRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameLocationObjectEffect) (*adventure_game_record.AdventureGameLocationObjectEffect, error) {
	l.Debug("mapping adventure_game_location_object_effect request to record")

	var req adventure_game_schema.AdventureGameLocationObjectEffectRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	applyEffectRequest := func() {
		rec.AdventureGameLocationObjectID = req.AdventureGameLocationObjectID
		rec.ActionType = req.ActionType
		if req.RequiredState != nil {
			rec.RequiredState = sql.NullString{String: *req.RequiredState, Valid: true}
		}
		if req.RequiredAdventureGameItemID != nil {
			rec.RequiredAdventureGameItemID = sql.NullString{String: *req.RequiredAdventureGameItemID, Valid: true}
		}
		rec.ResultDescription = req.ResultDescription
		rec.EffectType = req.EffectType
		if req.ResultState != nil {
			rec.ResultState = sql.NullString{String: *req.ResultState, Valid: true}
		}
		if req.ResultAdventureGameItemID != nil {
			rec.ResultAdventureGameItemID = sql.NullString{String: *req.ResultAdventureGameItemID, Valid: true}
		}
		if req.ResultAdventureGameLocationLinkID != nil {
			rec.ResultAdventureGameLocationLinkID = sql.NullString{String: *req.ResultAdventureGameLocationLinkID, Valid: true}
		}
		if req.ResultAdventureGameCreatureID != nil {
			rec.ResultAdventureGameCreatureID = sql.NullString{String: *req.ResultAdventureGameCreatureID, Valid: true}
		}
		if req.ResultAdventureGameLocationObjectID != nil {
			rec.ResultAdventureGameLocationObjectID = sql.NullString{String: *req.ResultAdventureGameLocationObjectID, Valid: true}
		}
		if req.ResultAdventureGameLocationID != nil {
			rec.ResultAdventureGameLocationID = sql.NullString{String: *req.ResultAdventureGameLocationID, Valid: true}
		}
		if req.ResultValueMin != nil {
			rec.ResultValueMin = sql.NullInt32{Int32: *req.ResultValueMin, Valid: true}
		}
		if req.ResultValueMax != nil {
			rec.ResultValueMax = sql.NullInt32{Int32: *req.ResultValueMax, Valid: true}
		}
		rec.IsRepeatable = req.IsRepeatable
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		applyEffectRequest()
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AdventureGameLocationObjectEffectRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationObjectEffect) (*adventure_game_schema.AdventureGameLocationObjectEffectResponseData, error) {
	l.Debug("mapping adventure_game_location_object_effect record to response data")

	nullStringPtr := func(ns sql.NullString) *string {
		if ns.Valid {
			return &ns.String
		}
		return nil
	}
	nullInt32Ptr := func(ni sql.NullInt32) *int32 {
		if ni.Valid {
			return &ni.Int32
		}
		return nil
	}

	return &adventure_game_schema.AdventureGameLocationObjectEffectResponseData{
		ID:                                  rec.ID,
		GameID:                              rec.GameID,
		AdventureGameLocationObjectID:       rec.AdventureGameLocationObjectID,
		ActionType:                          rec.ActionType,
		RequiredState:                       nullStringPtr(rec.RequiredState),
		RequiredAdventureGameItemID:         nullStringPtr(rec.RequiredAdventureGameItemID),
		ResultDescription:                   rec.ResultDescription,
		EffectType:                          rec.EffectType,
		ResultState:                         nullStringPtr(rec.ResultState),
		ResultAdventureGameItemID:           nullStringPtr(rec.ResultAdventureGameItemID),
		ResultAdventureGameLocationLinkID:   nullStringPtr(rec.ResultAdventureGameLocationLinkID),
		ResultAdventureGameCreatureID:       nullStringPtr(rec.ResultAdventureGameCreatureID),
		ResultAdventureGameLocationObjectID: nullStringPtr(rec.ResultAdventureGameLocationObjectID),
		ResultAdventureGameLocationID:       nullStringPtr(rec.ResultAdventureGameLocationID),
		ResultValueMin:                      nullInt32Ptr(rec.ResultValueMin),
		ResultValueMax:                      nullInt32Ptr(rec.ResultValueMax),
		IsRepeatable:                        rec.IsRepeatable,
		CreatedAt:                           rec.CreatedAt,
		UpdatedAt:                           nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:                           nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameLocationObjectEffectRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationObjectEffect) (*adventure_game_schema.AdventureGameLocationObjectEffectResponse, error) {
	l.Debug("mapping adventure_game_location_object_effect record to response")
	data, err := AdventureGameLocationObjectEffectRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameLocationObjectEffectResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationObjectEffectRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationObjectEffect) (adventure_game_schema.AdventureGameLocationObjectEffectCollectionResponse, error) {
	l.Debug("mapping adventure_game_location_object_effect records to collection response")
	data := []*adventure_game_schema.AdventureGameLocationObjectEffectResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationObjectEffectRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameLocationObjectEffectCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameLocationObjectEffectCollectionResponse{
		Data: data,
	}, nil
}
