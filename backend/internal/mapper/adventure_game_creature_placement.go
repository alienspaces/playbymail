package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameCreaturePlacementRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCreaturePlacement) (adventure_game_schema.AdventureGameCreaturePlacementResponseData, error) {
	data := adventure_game_schema.AdventureGameCreaturePlacementResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		AdventureGameCreatureID: rec.AdventureGameCreatureID,
		AdventureGameLocationID: rec.AdventureGameLocationID,
		InitialCount:            rec.InitialCount,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameCreaturePlacementRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCreaturePlacement) (*adventure_game_schema.AdventureGameCreaturePlacementCollectionResponse, error) {
	data := []*adventure_game_schema.AdventureGameCreaturePlacementResponseData{}
	for _, rec := range recs {
		item, err := AdventureGameCreaturePlacementRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, &item)
	}
	return &adventure_game_schema.AdventureGameCreaturePlacementCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameCreaturePlacementRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCreaturePlacement) (*adventure_game_schema.AdventureGameCreaturePlacementResponse, error) {
	data, err := AdventureGameCreaturePlacementRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameCreaturePlacementResponse{
		Data: &data,
	}, nil
}

func AdventureGameCreaturePlacementRequestToRecord(l logger.Logger, req *adventure_game_schema.AdventureGameCreaturePlacementRequest, rec *adventure_game_record.AdventureGameCreaturePlacement) (*adventure_game_record.AdventureGameCreaturePlacement, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameCreaturePlacement{}
	}
	if req == nil {
		return nil, nil
	}
	rec.AdventureGameCreatureID = req.AdventureGameCreatureID
	rec.AdventureGameLocationID = req.AdventureGameLocationID
	rec.InitialCount = req.InitialCount

	return rec, nil
}
