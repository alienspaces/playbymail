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

// LocationLinkRequestToRecord maps a LocationLinkRequest to a record.LocationLink
func LocationLinkRequestToRecord(l logger.Logger, r *http.Request, rec *record.LocationLink) (*record.LocationLink, error) {
	l.Debug("mapping location_link request to record")

	var req schema.LocationLinkRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.FromLocationID = req.FromLocationID
		rec.ToLocationID = req.ToLocationID
		rec.Description = req.Description
		rec.Name = req.Name
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.FromLocationID = req.FromLocationID
		rec.ToLocationID = req.ToLocationID
		rec.Description = req.Description
		rec.Name = req.Name
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

// LocationLinkRecordToResponseData maps a record.LocationLink to schema.LocationLinkResponseData
func LocationLinkRecordToResponseData(l logger.Logger, rec *record.LocationLink) (schema.LocationLinkResponseData, error) {
	l.Debug("mapping location_link record to response data")
	data := schema.LocationLinkResponseData{
		ID:             rec.ID,
		FromLocationID: rec.FromLocationID,
		ToLocationID:   rec.ToLocationID,
		Description:    rec.Description,
		Name:           rec.Name,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:      nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

// LocationLinkRecordToResponse maps a record.LocationLink to a schema.LocationLinkResponse
func LocationLinkRecordToResponse(rec *record.LocationLink) *schema.LocationLinkResponse {
	data, _ := LocationLinkRecordToResponseData(nil, rec)
	return &schema.LocationLinkResponse{
		LocationLinkResponseData: &data,
	}
}
