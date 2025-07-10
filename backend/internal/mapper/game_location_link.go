package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

// GameLocationLinkRequestToRecord maps a GameLocationLinkRequest to a record.GameLocationLink
func GameLocationLinkRequestToRecord(l logger.Logger, r *http.Request, rec *record.GameLocationLink) (*record.GameLocationLink, error) {
	l.Debug("mapping location_link request to record")

	var req schema.GameLocationLinkRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.FromGameLocationID = req.FromGameLocationID
		rec.ToGameLocationID = req.ToGameLocationID
		rec.Description = req.Description
		rec.Name = req.Name
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

// GameLocationLinkRecordToResponseData maps a record.GameLocationLink to schema.GameLocationLinkResponseData
func GameLocationLinkRecordToResponseData(l logger.Logger, rec *record.GameLocationLink) (schema.GameLocationLinkResponseData, error) {
	l.Debug("mapping location_link record to response data")
	data := schema.GameLocationLinkResponseData{
		ID:                 rec.ID,
		FromGameLocationID: rec.FromGameLocationID,
		ToGameLocationID:   rec.ToGameLocationID,
		Description:        rec.Description,
		Name:               rec.Name,
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:          nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

// GameLocationLinkRecordToResponse maps a record.GameLocationLink to a schema.GameLocationLinkResponse
func GameLocationLinkRecordToResponse(l logger.Logger, rec *record.GameLocationLink) (schema.GameLocationLinkResponse, error) {
	data, err := GameLocationLinkRecordToResponseData(l, rec)
	if err != nil {
		return schema.GameLocationLinkResponse{}, err
	}
	return schema.GameLocationLinkResponse{
		Data: &data,
	}, nil
}

func GameLocationLinkRecordsToCollectionResponse(l logger.Logger, recs []*record.GameLocationLink) (schema.GameLocationLinkCollectionResponse, error) {
	var data []*schema.GameLocationLinkResponseData
	for _, rec := range recs {
		d, err := GameLocationLinkRecordToResponseData(l, rec)
		if err != nil {
			return schema.GameLocationLinkCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.GameLocationLinkCollectionResponse{
		Data: data,
	}, nil
}
