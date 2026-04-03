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

func MechaSectorLinkRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_record.MechaSectorLink) (*mecha_record.MechaSectorLink, error) {
	l.Debug("mapping mecha_sector_link request to record")

	var req mecha_schema.MechaSectorLinkRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.FromMechaSectorID = req.FromMechaSectorID
		rec.ToMechaSectorID = req.ToMechaSectorID
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaSectorLinkRecordToResponseData(l logger.Logger, rec *mecha_record.MechaSectorLink) (*mecha_schema.MechaSectorLinkResponseData, error) {
	l.Debug("mapping mecha_sector_link record to response data")
	return &mecha_schema.MechaSectorLinkResponseData{
		ID:                rec.ID,
		GameID:            rec.GameID,
		FromMechaSectorID: rec.FromMechaSectorID,
		ToMechaSectorID:   rec.ToMechaSectorID,
		CreatedAt:         rec.CreatedAt,
		UpdatedAt:         nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:         nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaSectorLinkRecordToResponse(l logger.Logger, rec *mecha_record.MechaSectorLink) (*mecha_schema.MechaSectorLinkResponse, error) {
	l.Debug("mapping mecha_sector_link record to response")
	data, err := MechaSectorLinkRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_schema.MechaSectorLinkResponse{
		Data: data,
	}, nil
}

func MechaSectorLinkRecordsToCollectionResponse(l logger.Logger, recs []*mecha_record.MechaSectorLink) (mecha_schema.MechaSectorLinkCollectionResponse, error) {
	l.Debug("mapping mecha_sector_link records to collection response")
	data := []*mecha_schema.MechaSectorLinkResponseData{}
	for _, rec := range recs {
		d, err := MechaSectorLinkRecordToResponseData(l, rec)
		if err != nil {
			return mecha_schema.MechaSectorLinkCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_schema.MechaSectorLinkCollectionResponse{
		Data: data,
	}, nil
}
