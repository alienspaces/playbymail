package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func AdventureGameLocationLinkRequirementRequestToRecord(l logger.Logger, req *api.AdventureGameLocationLinkRequirementRequest, rec *adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocationLinkRequirement{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping adventure_game_location_link_requirement request to record")
	rec.GameID = req.GameID
	rec.AdventureGameLocationLinkID = req.GameLocationLinkID
	rec.AdventureGameItemID = req.GameItemID
	rec.Quantity = req.Quantity
	return rec, nil
}

func AdventureGameLocationLinkRequirementRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLinkRequirement) (api.AdventureGameLocationLinkRequirementResponseData, error) {
	l.Debug("mapping adventure_game_location_link_requirement record to response data")
	data := api.AdventureGameLocationLinkRequirementResponseData{
		ID:                 rec.ID,
		GameID:             rec.GameID,
		GameLocationLinkID: rec.AdventureGameLocationLinkID,
		GameItemID:         rec.AdventureGameItemID,
		Quantity:           rec.Quantity,
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:          nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameLocationLinkRequirementRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLinkRequirement) (api.AdventureGameLocationLinkRequirementResponse, error) {
	data, err := AdventureGameLocationLinkRequirementRecordToResponseData(l, rec)
	if err != nil {
		return api.AdventureGameLocationLinkRequirementResponse{}, err
	}
	return api.AdventureGameLocationLinkRequirementResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationLinkRequirementRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationLinkRequirement) (api.AdventureGameLocationLinkRequirementCollectionResponse, error) {
	data := []*api.AdventureGameLocationLinkRequirementResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationLinkRequirementRecordToResponseData(l, rec)
		if err != nil {
			return api.AdventureGameLocationLinkRequirementCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AdventureGameLocationLinkRequirementCollectionResponse{
		Data: data,
	}, nil
}
