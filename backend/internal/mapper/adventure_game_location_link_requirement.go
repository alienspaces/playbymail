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

// AdventureGameLocationLinkRequirementRequestToRecord maps a request to a record for consistency
func AdventureGameLocationLinkRequirementRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	l.Debug("mapping adventure_game_location_link_requirement request to record")

	var req adventure_game_schema.AdventureGameLocationLinkRequirementRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.AdventureGameLocationLinkID = req.GameLocationLinkID // Map new field name to old
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.AdventureGameLocationLinkID = req.GameLocationLinkID // Map new field name to old
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AdventureGameLocationLinkRequirementRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_schema.AdventureGameLocationLinkRequirementResponseData, error) {
	l.Debug("mapping adventure_game_location_link_requirement record to response data")
	return &adventure_game_schema.AdventureGameLocationLinkRequirementResponseData{
		ID:                 rec.ID,
		GameID:             rec.GameID,
		GameLocationLinkID: rec.AdventureGameLocationLinkID, // Map old field name to new
		RequirementType:    "item",                          // TODO: This field doesn't exist in record, hardcoded for now
		RequirementValue:   rec.AdventureGameItemID,         // Map old field name to new
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:          nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameLocationLinkRequirementRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_schema.AdventureGameLocationLinkRequirementCollectionResponse, error) {
	l.Debug("mapping adventure_game_location_link_requirement records to collection response")
	data := []*adventure_game_schema.AdventureGameLocationLinkRequirementResponseData{}
	for _, rec := range recs {
		item, err := AdventureGameLocationLinkRequirementRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return &adventure_game_schema.AdventureGameLocationLinkRequirementCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationLinkRequirementRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_schema.AdventureGameLocationLinkRequirementResponse, error) {
	l.Debug("mapping adventure_game_location_link_requirement record to response")
	data, err := AdventureGameLocationLinkRequirementRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameLocationLinkRequirementResponse{
		Data: data,
	}, nil
}
