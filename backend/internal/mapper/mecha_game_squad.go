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

func MechaGameSquadRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_game_record.MechaGameSquad) (*mecha_game_record.MechaGameSquad, error) {
	l.Debug("mapping mecha_game_squad request to record")

	var req mecha_game_schema.MechaGameSquadRequest
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

func MechaGameSquadRecordToResponseData(l logger.Logger, rec *mecha_game_record.MechaGameSquad) (*mecha_game_schema.MechaGameSquadResponseData, error) {
	l.Debug("mapping mecha_game_squad record to response data")

	data := &mecha_game_schema.MechaGameSquadResponseData{
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

func MechaGameSquadRecordToResponse(l logger.Logger, rec *mecha_game_record.MechaGameSquad) (*mecha_game_schema.MechaGameSquadResponse, error) {
	l.Debug("mapping mecha_game_squad record to response")
	data, err := MechaGameSquadRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_game_schema.MechaGameSquadResponse{
		Data: data,
	}, nil
}

func MechaGameSquadRecordsToCollectionResponse(l logger.Logger, recs []*mecha_game_record.MechaGameSquad) (mecha_game_schema.MechaGameSquadCollectionResponse, error) {
	l.Debug("mapping mecha_game_squad records to collection response")
	data := []*mecha_game_schema.MechaGameSquadResponseData{}
	for _, rec := range recs {
		d, err := MechaGameSquadRecordToResponseData(l, rec)
		if err != nil {
			return mecha_game_schema.MechaGameSquadCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_game_schema.MechaGameSquadCollectionResponse{
		Data: data,
	}, nil
}
