package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func GameAdministrationRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameAdministration) (*game_record.GameAdministration, error) {
	l.Debug("mapping game_administration request to record")

	var req api.GameAdministrationRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameID = req.GameID
		rec.AccountID = req.AccountID
		rec.GrantedByAccountID = req.GrantedByAccountID
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.GameID = req.GameID
		rec.AccountID = req.AccountID
		rec.GrantedByAccountID = req.GrantedByAccountID
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameAdministrationRecordToResponseData(l logger.Logger, rec *game_record.GameAdministration) (api.GameAdministrationResponseData, error) {
	l.Debug("mapping game_administration record to response data")
	data := api.GameAdministrationResponseData{
		ID:                 rec.ID,
		GameID:             rec.GameID,
		AccountID:          rec.AccountID,
		GrantedByAccountID: rec.GrantedByAccountID,
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:          nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func GameAdministrationRecordToResponse(l logger.Logger, rec *game_record.GameAdministration) (api.GameAdministrationResponse, error) {
	data, err := GameAdministrationRecordToResponseData(l, rec)
	if err != nil {
		return api.GameAdministrationResponse{}, err
	}
	return api.GameAdministrationResponse{
		Data: &data,
	}, nil
}

func GameAdministrationRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameAdministration) (api.GameAdministrationCollectionResponse, error) {
	data := []*api.GameAdministrationResponseData{}
	for _, rec := range recs {
		d, err := GameAdministrationRecordToResponseData(l, rec)
		if err != nil {
			return api.GameAdministrationCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.GameAdministrationCollectionResponse{
		Data: data,
	}, nil
}
