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

func AdventureGameItemEffectRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameItemEffect) (*adventure_game_record.AdventureGameItemEffect, error) {
	l.Debug("mapping adventure_game_item_effect request to record")

	var req adventure_game_schema.AdventureGameItemEffectRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	applyEffectRequest := func() {
		rec.AdventureGameItemID = req.AdventureGameItemID
		rec.ActionType = req.ActionType
		rec.RequiredAdventureGameItemID = nullstring.FromStringPtr(req.RequiredAdventureGameItemID)
		rec.RequiredAdventureGameLocationID = nullstring.FromStringPtr(req.RequiredAdventureGameLocationID)
		rec.ResultDescription = req.ResultDescription
		rec.EffectType = req.EffectType
		rec.ResultAdventureGameItemID = nullstring.FromStringPtr(req.ResultAdventureGameItemID)
		rec.ResultAdventureGameLocationLinkID = nullstring.FromStringPtr(req.ResultAdventureGameLocationLinkID)
		rec.ResultAdventureGameCreatureID = nullstring.FromStringPtr(req.ResultAdventureGameCreatureID)
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

func AdventureGameItemEffectRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameItemEffect) (*adventure_game_schema.AdventureGameItemEffectResponseData, error) {
	l.Debug("mapping adventure_game_item_effect record to response data")

	return &adventure_game_schema.AdventureGameItemEffectResponseData{
		ID:                                rec.ID,
		GameID:                            rec.GameID,
		AdventureGameItemID:               rec.AdventureGameItemID,
		ActionType:                        rec.ActionType,
		RequiredAdventureGameItemID:       nullstring.ToStringPtr(rec.RequiredAdventureGameItemID),
		RequiredAdventureGameLocationID:   nullstring.ToStringPtr(rec.RequiredAdventureGameLocationID),
		ResultDescription:                 rec.ResultDescription,
		EffectType:                        rec.EffectType,
		ResultAdventureGameItemID:         nullstring.ToStringPtr(rec.ResultAdventureGameItemID),
		ResultAdventureGameLocationLinkID: nullstring.ToStringPtr(rec.ResultAdventureGameLocationLinkID),
		ResultAdventureGameCreatureID:     nullstring.ToStringPtr(rec.ResultAdventureGameCreatureID),
		ResultAdventureGameLocationID:     nullstring.ToStringPtr(rec.ResultAdventureGameLocationID),
		ResultValueMin:                    nullint32.ToInt32PtrOrNil(rec.ResultValueMin),
		ResultValueMax:                    nullint32.ToInt32PtrOrNil(rec.ResultValueMax),
		IsRepeatable:                      rec.IsRepeatable,
		CreatedAt:                         rec.CreatedAt,
		UpdatedAt:                         nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:                         nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameItemEffectRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameItemEffect) (*adventure_game_schema.AdventureGameItemEffectResponse, error) {
	l.Debug("mapping adventure_game_item_effect record to response")
	data, err := AdventureGameItemEffectRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameItemEffectResponse{
		Data: data,
	}, nil
}

func AdventureGameItemEffectRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameItemEffect) (adventure_game_schema.AdventureGameItemEffectCollectionResponse, error) {
	l.Debug("mapping adventure_game_item_effect records to collection response")
	data := []*adventure_game_schema.AdventureGameItemEffectResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameItemEffectRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameItemEffectCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameItemEffectCollectionResponse{
		Data: data,
	}, nil
}
