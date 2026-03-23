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

func MechaSectorRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_record.MechaSector) (*mecha_record.MechaSector, error) {
	l.Debug("mapping mecha_sector request to record")

	var req mecha_schema.MechaSectorRequest
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
			rec.TerrainType = mecha_record.SectorTerrainTypeOpen
		}
		rec.Elevation = req.Elevation
		rec.IsStartingSector = req.IsStartingSector
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaSectorRecordToResponseData(l logger.Logger, rec *mecha_record.MechaSector) (*mecha_schema.MechaSectorResponseData, error) {
	l.Debug("mapping mecha_sector record to response data")
	return &mecha_schema.MechaSectorResponseData{
		ID:               rec.ID,
		GameID:           rec.GameID,
		Name:             rec.Name,
		Description:      rec.Description,
		TerrainType:      rec.TerrainType,
		Elevation:        rec.Elevation,
		IsStartingSector: rec.IsStartingSector,
		CreatedAt:        rec.CreatedAt,
		UpdatedAt:        nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:        nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaSectorRecordToResponse(l logger.Logger, rec *mecha_record.MechaSector) (*mecha_schema.MechaSectorResponse, error) {
	l.Debug("mapping mecha_sector record to response")
	data, err := MechaSectorRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_schema.MechaSectorResponse{
		Data: data,
	}, nil
}

func MechaSectorRecordsToCollectionResponse(l logger.Logger, recs []*mecha_record.MechaSector) (mecha_schema.MechaSectorCollectionResponse, error) {
	l.Debug("mapping mecha_sector records to collection response")
	data := []*mecha_schema.MechaSectorResponseData{}
	for _, rec := range recs {
		d, err := MechaSectorRecordToResponseData(l, rec)
		if err != nil {
			return mecha_schema.MechaSectorCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_schema.MechaSectorCollectionResponse{
		Data: data,
	}, nil
}
