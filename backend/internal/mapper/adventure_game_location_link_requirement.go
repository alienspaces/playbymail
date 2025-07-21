package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AdventureGameLocationLinkRequirementRequestToRecord(l logger.Logger, req *schema.AdventureGameLocationLinkRequirementRequest, rec *record.AdventureGameLocationLinkRequirement) (*record.AdventureGameLocationLinkRequirement, error) {
	if rec == nil {
		rec = &record.AdventureGameLocationLinkRequirement{}
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

func AdventureGameLocationLinkRequirementRecordToResponseData(l logger.Logger, rec *record.AdventureGameLocationLinkRequirement) (schema.AdventureGameLocationLinkRequirementResponseData, error) {
	l.Debug("mapping adventure_game_location_link_requirement record to response data")
	data := schema.AdventureGameLocationLinkRequirementResponseData{
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

func AdventureGameLocationLinkRequirementRecordToResponse(l logger.Logger, rec *record.AdventureGameLocationLinkRequirement) (schema.AdventureGameLocationLinkRequirementResponse, error) {
	data, err := AdventureGameLocationLinkRequirementRecordToResponseData(l, rec)
	if err != nil {
		return schema.AdventureGameLocationLinkRequirementResponse{}, err
	}
	return schema.AdventureGameLocationLinkRequirementResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationLinkRequirementRecordsToCollectionResponse(l logger.Logger, recs []*record.AdventureGameLocationLinkRequirement) (schema.AdventureGameLocationLinkRequirementCollectionResponse, error) {
	var data []*schema.AdventureGameLocationLinkRequirementResponseData
	for _, rec := range recs {
		d, err := AdventureGameLocationLinkRequirementRecordToResponseData(l, rec)
		if err != nil {
			return schema.AdventureGameLocationLinkRequirementCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.AdventureGameLocationLinkRequirementCollectionResponse{
		Data: data,
	}, nil
}
