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

func MechWargameSectorRequestToRecord(l logger.Logger, r *http.Request, rec *mech_wargame_record.MechWargameSector) (*mech_wargame_record.MechWargameSector, error) {
	l.Debug("mapping mech_wargame_sector request to record")

	var req mech_wargame_schema.MechWargameSectorRequest
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
			rec.TerrainType = mech_wargame_record.SectorTerrainTypeOpen
		}
		rec.Elevation = req.Elevation
		rec.IsStartingSector = req.IsStartingSector
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechWargameSectorRecordToResponseData(l logger.Logger, rec *mech_wargame_record.MechWargameSector) (*mech_wargame_schema.MechWargameSectorResponseData, error) {
	l.Debug("mapping mech_wargame_sector record to response data")
	return &mech_wargame_schema.MechWargameSectorResponseData{
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

func MechWargameSectorRecordToResponse(l logger.Logger, rec *mech_wargame_record.MechWargameSector) (*mech_wargame_schema.MechWargameSectorResponse, error) {
	l.Debug("mapping mech_wargame_sector record to response")
	data, err := MechWargameSectorRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mech_wargame_schema.MechWargameSectorResponse{
		Data: data,
	}, nil
}

func MechWargameSectorRecordsToCollectionResponse(l logger.Logger, recs []*mech_wargame_record.MechWargameSector) (mech_wargame_schema.MechWargameSectorCollectionResponse, error) {
	l.Debug("mapping mech_wargame_sector records to collection response")
	data := []*mech_wargame_schema.MechWargameSectorResponseData{}
	for _, rec := range recs {
		d, err := MechWargameSectorRecordToResponseData(l, rec)
		if err != nil {
			return mech_wargame_schema.MechWargameSectorCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mech_wargame_schema.MechWargameSectorCollectionResponse{
		Data: data,
	}, nil
}
