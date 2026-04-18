package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/mecha_game_schema"
)

func MechaGameComputerOpponentRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_game_record.MechaGameComputerOpponent) (*mecha_game_record.MechaGameComputerOpponent, error) {
	l.Debug("mapping mecha_game_computer_opponent request to record")

	var req mecha_game_schema.MechaGameComputerOpponentRequest
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

func MechaGameComputerOpponentRecordToResponseData(l logger.Logger, rec *mecha_game_record.MechaGameComputerOpponent) (*mecha_game_schema.MechaGameComputerOpponentResponseData, error) {
	l.Debug("mapping mecha_game_computer_opponent record to response data")
	return &mecha_game_schema.MechaGameComputerOpponentResponseData{
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

func MechaGameComputerOpponentRecordToResponse(l logger.Logger, rec *mecha_game_record.MechaGameComputerOpponent) (*mecha_game_schema.MechaGameComputerOpponentResponse, error) {
	l.Debug("mapping mecha_game_computer_opponent record to response")
	data, err := MechaGameComputerOpponentRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_game_schema.MechaGameComputerOpponentResponse{
		Data: data,
	}, nil
}

func MechaGameComputerOpponentRecordsToCollectionResponse(l logger.Logger, recs []*mecha_game_record.MechaGameComputerOpponent) (mecha_game_schema.MechaGameComputerOpponentCollectionResponse, error) {
	l.Debug("mapping mecha_game_computer_opponent records to collection response")
	data := []*mecha_game_schema.MechaGameComputerOpponentResponseData{}
	for _, rec := range recs {
		d, err := MechaGameComputerOpponentRecordToResponseData(l, rec)
		if err != nil {
			return mecha_game_schema.MechaGameComputerOpponentCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_game_schema.MechaGameComputerOpponentCollectionResponse{
		Data: data,
	}, nil
}
