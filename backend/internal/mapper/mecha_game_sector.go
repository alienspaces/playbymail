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

func MechaGameSectorRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_game_record.MechaGameSector) (*mecha_game_record.MechaGameSector, error) {
	l.Debug("mapping mecha_game_sector request to record")

	var req mecha_game_schema.MechaGameSectorRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.Description = req.Description
		rec.TerrainType = req.TerrainType
		if rec.TerrainType == "" {
			rec.TerrainType = mecha_game_record.SectorTerrainTypeOpen
		}
		rec.Elevation = req.Elevation
		rec.CoverModifier = req.CoverModifier
		rec.IsStartingSector = req.IsStartingSector
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaGameSectorRecordToResponseData(l logger.Logger, rec *mecha_game_record.MechaGameSector) (*mecha_game_schema.MechaGameSectorResponseData, error) {
	l.Debug("mapping mecha_game_sector record to response data")
	return &mecha_game_schema.MechaGameSectorResponseData{
		ID:               rec.ID,
		GameID:           rec.GameID,
		Name:             rec.Name,
		Description:      rec.Description,
		TerrainType:      rec.TerrainType,
		Elevation:        rec.Elevation,
		CoverModifier:    rec.CoverModifier,
		IsStartingSector: rec.IsStartingSector,
		CreatedAt:        rec.CreatedAt,
		UpdatedAt:        nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:        nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaGameSectorRecordToResponse(l logger.Logger, rec *mecha_game_record.MechaGameSector) (*mecha_game_schema.MechaGameSectorResponse, error) {
	l.Debug("mapping mecha_game_sector record to response")
	data, err := MechaGameSectorRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_game_schema.MechaGameSectorResponse{
		Data: data,
	}, nil
}

func MechaGameSectorRecordsToCollectionResponse(l logger.Logger, recs []*mecha_game_record.MechaGameSector) (mecha_game_schema.MechaGameSectorCollectionResponse, error) {
	l.Debug("mapping mecha_game_sector records to collection response")
	data := []*mecha_game_schema.MechaGameSectorResponseData{}
	for _, rec := range recs {
		d, err := MechaGameSectorRecordToResponseData(l, rec)
		if err != nil {
			return mecha_game_schema.MechaGameSectorCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_game_schema.MechaGameSectorCollectionResponse{
		Data: data,
	}, nil
}
