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

func MechWargameSectorLinkRequestToRecord(l logger.Logger, r *http.Request, rec *mech_wargame_record.MechWargameSectorLink) (*mech_wargame_record.MechWargameSectorLink, error) {
	l.Debug("mapping mech_wargame_sector_link request to record")

	var req mech_wargame_schema.MechWargameSectorLinkRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.FromMechWargameSectorID = req.FromMechWargameSectorID
		rec.ToMechWargameSectorID = req.ToMechWargameSectorID
		rec.CoverModifier = req.CoverModifier
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechWargameSectorLinkRecordToResponseData(l logger.Logger, rec *mech_wargame_record.MechWargameSectorLink) (*mech_wargame_schema.MechWargameSectorLinkResponseData, error) {
	l.Debug("mapping mech_wargame_sector_link record to response data")
	return &mech_wargame_schema.MechWargameSectorLinkResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		FromMechWargameSectorID: rec.FromMechWargameSectorID,
		ToMechWargameSectorID:   rec.ToMechWargameSectorID,
		CoverModifier:           rec.CoverModifier,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechWargameSectorLinkRecordToResponse(l logger.Logger, rec *mech_wargame_record.MechWargameSectorLink) (*mech_wargame_schema.MechWargameSectorLinkResponse, error) {
	l.Debug("mapping mech_wargame_sector_link record to response")
	data, err := MechWargameSectorLinkRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mech_wargame_schema.MechWargameSectorLinkResponse{
		Data: data,
	}, nil
}

func MechWargameSectorLinkRecordsToCollectionResponse(l logger.Logger, recs []*mech_wargame_record.MechWargameSectorLink) (mech_wargame_schema.MechWargameSectorLinkCollectionResponse, error) {
	l.Debug("mapping mech_wargame_sector_link records to collection response")
	data := []*mech_wargame_schema.MechWargameSectorLinkResponseData{}
	for _, rec := range recs {
		d, err := MechWargameSectorLinkRecordToResponseData(l, rec)
		if err != nil {
			return mech_wargame_schema.MechWargameSectorLinkCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mech_wargame_schema.MechWargameSectorLinkCollectionResponse{
		Data: data,
	}, nil
}
