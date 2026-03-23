package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
	"gitlab.com/alienspaces/playbymail/schema/api/mech_wargame_schema"
)

func MechWargameLanceRequestToRecord(l logger.Logger, r *http.Request, rec *mech_wargame_record.MechWargameLance) (*mech_wargame_record.MechWargameLance, error) {
	l.Debug("mapping mech_wargame_lance request to record")

	var req mech_wargame_schema.MechWargameLanceRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.AccountUserID = req.AccountUserID
		rec.Name = req.Name
		rec.Description = req.Description
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechWargameLanceRecordToResponseData(l logger.Logger, rec *mech_wargame_record.MechWargameLance) (*mech_wargame_schema.MechWargameLanceResponseData, error) {
	l.Debug("mapping mech_wargame_lance record to response data")
	return &mech_wargame_schema.MechWargameLanceResponseData{
		ID:            rec.ID,
		GameID:        rec.GameID,
		AccountID:     rec.AccountID,
		AccountUserID: rec.AccountUserID,
		Name:          rec.Name,
		Description:   rec.Description,
		CreatedAt:     rec.CreatedAt,
		UpdatedAt:     nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:     nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechWargameLanceRecordToResponse(l logger.Logger, rec *mech_wargame_record.MechWargameLance) (*mech_wargame_schema.MechWargameLanceResponse, error) {
	l.Debug("mapping mech_wargame_lance record to response")
	data, err := MechWargameLanceRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mech_wargame_schema.MechWargameLanceResponse{
		Data: data,
	}, nil
}

func MechWargameLanceRecordsToCollectionResponse(l logger.Logger, recs []*mech_wargame_record.MechWargameLance) (mech_wargame_schema.MechWargameLanceCollectionResponse, error) {
	l.Debug("mapping mech_wargame_lance records to collection response")
	data := []*mech_wargame_schema.MechWargameLanceResponseData{}
	for _, rec := range recs {
		d, err := MechWargameLanceRecordToResponseData(l, rec)
		if err != nil {
			return mech_wargame_schema.MechWargameLanceCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mech_wargame_schema.MechWargameLanceCollectionResponse{
		Data: data,
	}, nil
}
