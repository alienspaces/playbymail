package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func GameLocationLinkRequirementRequestToRecord(l logger.Logger, req *schema.GameLocationLinkRequirementRequest, rec *record.GameLocationLinkRequirement) (*record.GameLocationLinkRequirement, error) {
	if rec == nil {
		rec = &record.GameLocationLinkRequirement{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping game_location_link_requirement request to record")
	rec.GameLocationLinkID = req.GameLocationLinkID
	rec.GameItemID = req.GameItemID
	rec.Quantity = req.Quantity
	return rec, nil
}

func GameLocationLinkRequirementRecordToResponseData(l logger.Logger, rec *record.GameLocationLinkRequirement) (schema.GameLocationLinkRequirementResponseData, error) {
	l.Debug("mapping game_location_link_requirement record to response data")
	data := schema.GameLocationLinkRequirementResponseData{
		ID:                 rec.ID,
		GameLocationLinkID: rec.GameLocationLinkID,
		GameItemID:         rec.GameItemID,
		Quantity:           rec.Quantity,
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:          nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func GameLocationLinkRequirementRecordToResponse(l logger.Logger, rec *record.GameLocationLinkRequirement) (schema.GameLocationLinkRequirementResponse, error) {
	data, err := GameLocationLinkRequirementRecordToResponseData(l, rec)
	if err != nil {
		return schema.GameLocationLinkRequirementResponse{}, err
	}
	return schema.GameLocationLinkRequirementResponse{
		Data: &data,
	}, nil
}

func GameLocationLinkRequirementRecordsToCollectionResponse(l logger.Logger, recs []*record.GameLocationLinkRequirement) (schema.GameLocationLinkRequirementCollectionResponse, error) {
	var data []*schema.GameLocationLinkRequirementResponseData
	for _, rec := range recs {
		d, err := GameLocationLinkRequirementRecordToResponseData(l, rec)
		if err != nil {
			return schema.GameLocationLinkRequirementCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.GameLocationLinkRequirementCollectionResponse{
		Data: data,
	}, nil
}
