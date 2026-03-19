package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nullint32"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
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
		rec.RequiredAdventureGameLocationObjectStateID = nullstring.FromStringPtr(req.RequiredAdventureGameLocationObjectStateID)
		rec.RequiredAdventureGameItemID = nullstring.FromStringPtr(req.RequiredAdventureGameItemID)
		rec.ResultDescription = req.ResultDescription
		rec.EffectType = req.EffectType
		rec.ResultAdventureGameLocationObjectStateID = nullstring.FromStringPtr(req.ResultAdventureGameLocationObjectStateID)
		rec.ResultAdventureGameItemID = nullstring.FromStringPtr(req.ResultAdventureGameItemID)
		rec.ResultAdventureGameLocationLinkID = nullstring.FromStringPtr(req.ResultAdventureGameLocationLinkID)
		rec.ResultAdventureGameCreatureID = nullstring.FromStringPtr(req.ResultAdventureGameCreatureID)
		rec.ResultAdventureGameLocationObjectID = nullstring.FromStringPtr(req.ResultAdventureGameLocationObjectID)
		rec.ResultAdventureGameLocationID = nullstring.FromStringPtr(req.ResultAdventureGameLocationID)
		rec.ResultValueMin = nullint32.FromInt32Ptr(req.ResultValueMin)
		rec.ResultValueMax = nullint32.FromInt32Ptr(req.ResultValueMax)
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

	return &adventure_game_schema.AdventureGameLocationObjectEffectResponseData{
		ID:                                              rec.ID,
		GameID:                                          rec.GameID,
		AdventureGameLocationObjectID:                   rec.AdventureGameLocationObjectID,
		ActionType:                                      rec.ActionType,
		RequiredAdventureGameLocationObjectStateID:      nullstring.ToStringPtr(rec.RequiredAdventureGameLocationObjectStateID),
		RequiredAdventureGameItemID:                     nullstring.ToStringPtr(rec.RequiredAdventureGameItemID),
		ResultDescription:                               rec.ResultDescription,
		EffectType:                                      rec.EffectType,
		ResultAdventureGameLocationObjectStateID:        nullstring.ToStringPtr(rec.ResultAdventureGameLocationObjectStateID),
		ResultAdventureGameItemID:                       nullstring.ToStringPtr(rec.ResultAdventureGameItemID),
		ResultAdventureGameLocationLinkID:               nullstring.ToStringPtr(rec.ResultAdventureGameLocationLinkID),
		ResultAdventureGameCreatureID:                   nullstring.ToStringPtr(rec.ResultAdventureGameCreatureID),
		ResultAdventureGameLocationObjectID:             nullstring.ToStringPtr(rec.ResultAdventureGameLocationObjectID),
		ResultAdventureGameLocationID:                   nullstring.ToStringPtr(rec.ResultAdventureGameLocationID),
		ResultValueMin:                                  nullint32.ToInt32PtrOrNil(rec.ResultValueMin),
		ResultValueMax:                                  nullint32.ToInt32PtrOrNil(rec.ResultValueMax),
		IsRepeatable:                                    rec.IsRepeatable,
		CreatedAt:                                       rec.CreatedAt,
		UpdatedAt:                                       nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:                                       nulltime.ToTimePtr(rec.DeletedAt),
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
