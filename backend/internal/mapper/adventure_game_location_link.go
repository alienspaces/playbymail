package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameLocationLinkRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLink) (adventure_game_schema.AdventureGameLocationLinkResponseData, error) {
	data := adventure_game_schema.AdventureGameLocationLinkResponseData{
		ID:                          rec.ID,
		GameID:                      rec.GameID,
		Name:                        rec.Name,
		Description:                 rec.Description,
		FromAdventureGameLocationID: rec.FromAdventureGameLocationID, // Map old field name to new
		ToAdventureGameLocationID:   rec.ToAdventureGameLocationID,   // Map old field name to new
		CreatedAt:                   rec.CreatedAt,
		UpdatedAt:                   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:                   nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameLocationLinkRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationLink) (*adventure_game_schema.AdventureGameLocationLinkCollectionResponse, error) {
	data := []*adventure_game_schema.AdventureGameLocationLinkResponseData{}
	for _, rec := range recs {
		item, err := AdventureGameLocationLinkRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, &item)
	}
	return &adventure_game_schema.AdventureGameLocationLinkCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationLinkRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_schema.AdventureGameLocationLinkResponse, error) {
	data, err := AdventureGameLocationLinkRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameLocationLinkResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationLinkRequestToRecord(l logger.Logger, req *adventure_game_schema.AdventureGameLocationLinkRequest, rec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_record.AdventureGameLocationLink, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocationLink{}
	}
	if req == nil {
		return nil, nil
	}

	rec.Name = req.Name
	rec.Description = req.Description
	rec.FromAdventureGameLocationID = req.FromAdventureGameLocationID
	rec.ToAdventureGameLocationID = req.ToAdventureGameLocationID

	return rec, nil
}
