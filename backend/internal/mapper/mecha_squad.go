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

func MechaSquadRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_record.MechaSquad) (*mecha_record.MechaSquad, error) {
	l.Debug("mapping mecha_squad request to record")

	var req mecha_schema.MechaSquadRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.SquadType = req.SquadType
		rec.Name = req.Name
		rec.Description = req.Description
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaSquadRecordToResponseData(l logger.Logger, rec *mecha_record.MechaSquad) (*mecha_schema.MechaSquadResponseData, error) {
	l.Debug("mapping mecha_squad record to response data")

	data := &mecha_schema.MechaSquadResponseData{
		ID:          rec.ID,
		GameID:      rec.GameID,
		SquadType:   rec.SquadType,
		Name:        rec.Name,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:   nulltime.ToTimePtr(rec.DeletedAt),
	}

	return data, nil
}

func MechaSquadRecordToResponse(l logger.Logger, rec *mecha_record.MechaSquad) (*mecha_schema.MechaSquadResponse, error) {
	l.Debug("mapping mecha_squad record to response")
	data, err := MechaSquadRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_schema.MechaSquadResponse{
		Data: data,
	}, nil
}

func MechaSquadRecordsToCollectionResponse(l logger.Logger, recs []*mecha_record.MechaSquad) (mecha_schema.MechaSquadCollectionResponse, error) {
	l.Debug("mapping mecha_squad records to collection response")
	data := []*mecha_schema.MechaSquadResponseData{}
	for _, rec := range recs {
		d, err := MechaSquadRecordToResponseData(l, rec)
		if err != nil {
			return mecha_schema.MechaSquadCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_schema.MechaSquadCollectionResponse{
		Data: data,
	}, nil
}
