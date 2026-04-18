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

func MechaGameSectorLinkRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_game_record.MechaGameSectorLink) (*mecha_game_record.MechaGameSectorLink, error) {
	l.Debug("mapping mecha_game_sector_link request to record")

	var req mecha_game_schema.MechaGameSectorLinkRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.FromMechaGameSectorID = req.FromMechaGameSectorID
		rec.ToMechaGameSectorID = req.ToMechaGameSectorID
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaGameSectorLinkRecordToResponseData(l logger.Logger, rec *mecha_game_record.MechaGameSectorLink) (*mecha_game_schema.MechaGameSectorLinkResponseData, error) {
	l.Debug("mapping mecha_game_sector_link record to response data")
	return &mecha_game_schema.MechaGameSectorLinkResponseData{
		ID:                rec.ID,
		GameID:            rec.GameID,
		FromMechaGameSectorID: rec.FromMechaGameSectorID,
		ToMechaGameSectorID:   rec.ToMechaGameSectorID,
		CreatedAt:         rec.CreatedAt,
		UpdatedAt:         nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:         nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaGameSectorLinkRecordToResponse(l logger.Logger, rec *mecha_game_record.MechaGameSectorLink) (*mecha_game_schema.MechaGameSectorLinkResponse, error) {
	l.Debug("mapping mecha_game_sector_link record to response")
	data, err := MechaGameSectorLinkRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_game_schema.MechaGameSectorLinkResponse{
		Data: data,
	}, nil
}

func MechaGameSectorLinkRecordsToCollectionResponse(l logger.Logger, recs []*mecha_game_record.MechaGameSectorLink) (mecha_game_schema.MechaGameSectorLinkCollectionResponse, error) {
	l.Debug("mapping mecha_game_sector_link records to collection response")
	data := []*mecha_game_schema.MechaGameSectorLinkResponseData{}
	for _, rec := range recs {
		d, err := MechaGameSectorLinkRecordToResponseData(l, rec)
		if err != nil {
			return mecha_game_schema.MechaGameSectorLinkCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_game_schema.MechaGameSectorLinkCollectionResponse{
		Data: data,
	}, nil
}
