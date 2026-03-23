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

func MechaComputerOpponentRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_record.MechaComputerOpponent) (*mecha_record.MechaComputerOpponent, error) {
	l.Debug("mapping mecha_computer_opponent request to record")

	var req mecha_schema.MechaComputerOpponentRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.Description = req.Description
		rec.Aggression = req.Aggression
		rec.IQ = req.IQ
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaComputerOpponentRecordToResponseData(l logger.Logger, rec *mecha_record.MechaComputerOpponent) (*mecha_schema.MechaComputerOpponentResponseData, error) {
	l.Debug("mapping mecha_computer_opponent record to response data")
	return &mecha_schema.MechaComputerOpponentResponseData{
		ID:          rec.ID,
		GameID:      rec.GameID,
		Name:        rec.Name,
		Description: rec.Description,
		Aggression:  rec.Aggression,
		IQ:          rec.IQ,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:   nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaComputerOpponentRecordToResponse(l logger.Logger, rec *mecha_record.MechaComputerOpponent) (*mecha_schema.MechaComputerOpponentResponse, error) {
	l.Debug("mapping mecha_computer_opponent record to response")
	data, err := MechaComputerOpponentRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_schema.MechaComputerOpponentResponse{
		Data: data,
	}, nil
}

func MechaComputerOpponentRecordsToCollectionResponse(l logger.Logger, recs []*mecha_record.MechaComputerOpponent) (mecha_schema.MechaComputerOpponentCollectionResponse, error) {
	l.Debug("mapping mecha_computer_opponent records to collection response")
	data := []*mecha_schema.MechaComputerOpponentResponseData{}
	for _, rec := range recs {
		d, err := MechaComputerOpponentRecordToResponseData(l, rec)
		if err != nil {
			return mecha_schema.MechaComputerOpponentCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_schema.MechaComputerOpponentCollectionResponse{
		Data: data,
	}, nil
}
