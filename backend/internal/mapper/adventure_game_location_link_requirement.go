package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameLocationLinkRequirementRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLinkRequirement) (adventure_game_schema.AdventureGameLocationLinkRequirementResponseData, error) {
	data := adventure_game_schema.AdventureGameLocationLinkRequirementResponseData{
		ID:                 rec.ID,
		GameID:             rec.GameID,
		GameLocationLinkID: rec.AdventureGameLocationLinkID, // Map old field name to new
		RequirementType:    "item",                          // TODO: This field doesn't exist in record, hardcoded for now
		RequirementValue:   rec.AdventureGameItemID,         // Map old field name to new
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:          nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameLocationLinkRequirementRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_schema.AdventureGameLocationLinkRequirementCollectionResponse, error) {
	data := []*adventure_game_schema.AdventureGameLocationLinkRequirementResponseData{}
	for _, rec := range recs {
		item, err := AdventureGameLocationLinkRequirementRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, &item)
	}
	return &adventure_game_schema.AdventureGameLocationLinkRequirementCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationLinkRequirementRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_schema.AdventureGameLocationLinkRequirementResponse, error) {
	data, err := AdventureGameLocationLinkRequirementRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameLocationLinkRequirementResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationLinkRequirementRequestToRecord(l logger.Logger, req *adventure_game_schema.AdventureGameLocationLinkRequirementRequest, rec *adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocationLinkRequirement{}
	}
	if req == nil {
		return nil, nil
	}
	rec.GameID = ""                                          // TODO: This field doesn't exist in request
	rec.AdventureGameLocationLinkID = req.GameLocationLinkID // Map new field name to old
	// TODO: These fields don't exist in the record but exist in the new schema
	// rec.RequirementType = req.RequirementType
	// rec.RequirementValue = req.RequirementValue
	// TODO: These fields don't exist in the new schema but exist in the record
	// rec.AdventureGameItemID = req.GameItemID
	// rec.Quantity = req.Quantity
	return rec, nil
}
