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

func AdventureGameInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameInstance) (*adventure_game_schema.AdventureGameInstanceResponseData, error) {
	l.Debug("mapping adventure_game_instance record to response data")
	return &adventure_game_schema.AdventureGameInstanceResponseData{
		ID:        rec.ID,
		GameID:    rec.GameID,
		AccountID: "", // TODO: (agent) Add account_id to game_instance (migration + record + domain), then set AccountID from rec here. If product decision is to not have instance owner, remove AccountID from API schema and this mapping instead.
		Status:    rec.Status,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt: nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameInstance) (*adventure_game_schema.AdventureGameInstanceResponse, error) {
	l.Debug("mapping adventure_game_instance record to response")
	data, err := AdventureGameInstanceRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameInstanceResponse{
		Data: data,
	}, nil
}

func AdventureGameInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameInstance) (adventure_game_schema.AdventureGameInstanceCollectionResponse, error) {
	l.Debug("mapping adventure_game_instance records to collection response")
	data := []*adventure_game_schema.AdventureGameInstanceResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameInstanceRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameInstanceCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameInstanceCollectionResponse{
		Data: data,
	}, nil
}

// AdventureGameInstanceRequestToRecord maps a request to a record for consistency
func AdventureGameInstanceRequestToRecord(l logger.Logger, r *http.Request, rec *adventure_game_record.AdventureGameInstance) (*adventure_game_record.AdventureGameInstance, error) {
	l.Debug("mapping adventure_game_instance request to record")

	var req adventure_game_schema.AdventureGameInstanceRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		// TODO: (agent) When account_id exists on game_instance record, map req or context account to rec.AccountID in request-to-record mapping.
		if req.Status != "" {
			rec.Status = req.Status
		}
	case server.HttpMethodPut, server.HttpMethodPatch:
		if req.Status != "" {
			rec.Status = req.Status
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}
