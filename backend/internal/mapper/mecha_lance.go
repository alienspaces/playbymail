package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
	"gitlab.com/alienspaces/playbymail/schema/api/mecha_schema"
)

func MechaLanceRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_record.MechaLance) (*mecha_record.MechaLance, error) {
	l.Debug("mapping mecha_lance request to record")

	var req mecha_schema.MechaLanceRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.LanceType = req.LanceType
		rec.Name = req.Name
		rec.Description = req.Description
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaLanceRecordToResponseData(l logger.Logger, rec *mecha_record.MechaLance) (*mecha_schema.MechaLanceResponseData, error) {
	l.Debug("mapping mecha_lance record to response data")

	data := &mecha_schema.MechaLanceResponseData{
		ID:          rec.ID,
		GameID:      rec.GameID,
		LanceType:   rec.LanceType,
		Name:        rec.Name,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:   nulltime.ToTimePtr(rec.DeletedAt),
	}

	return data, nil
}

func MechaLanceRecordToResponse(l logger.Logger, rec *mecha_record.MechaLance) (*mecha_schema.MechaLanceResponse, error) {
	l.Debug("mapping mecha_lance record to response")
	data, err := MechaLanceRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_schema.MechaLanceResponse{
		Data: data,
	}, nil
}

func MechaLanceRecordsToCollectionResponse(l logger.Logger, recs []*mecha_record.MechaLance) (mecha_schema.MechaLanceCollectionResponse, error) {
	l.Debug("mapping mecha_lance records to collection response")
	data := []*mecha_schema.MechaLanceResponseData{}
	for _, rec := range recs {
		d, err := MechaLanceRecordToResponseData(l, rec)
		if err != nil {
			return mecha_schema.MechaLanceCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_schema.MechaLanceCollectionResponse{
		Data: data,
	}, nil
}
