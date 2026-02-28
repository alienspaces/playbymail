package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/schema/api/account_schema"
)

func AccountUserContactRequestToRecord(l logger.Logger, r *http.Request, rec *account_record.AccountUserContact) (*account_record.AccountUserContact, error) {
	l.Debug("mapping account contact request to record")

	var req account_schema.AccountContactRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Name = req.Name
		rec.PostalAddressLine1 = req.PostalAddressLine1
		rec.PostalAddressLine2 = nullstring.FromString(req.PostalAddressLine2)
		rec.StateProvince = req.StateProvince
		rec.Country = req.Country
		rec.PostalCode = req.PostalCode
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.PostalAddressLine1 = req.PostalAddressLine1
		rec.PostalAddressLine2 = nullstring.FromString(req.PostalAddressLine2)
		rec.StateProvince = req.StateProvince
		rec.Country = req.Country
		rec.PostalCode = req.PostalCode
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AccountUserContactRecordToResponseData(l logger.Logger, rec *account_record.AccountUserContact) (*account_schema.AccountContactResponseData, error) {
	l.Debug("mapping account contact record to response data")
	return &account_schema.AccountContactResponseData{
		ID:                 rec.ID,
		AccountUserID:      rec.AccountUserID,
		Name:               rec.Name,
		PostalAddressLine1: rec.PostalAddressLine1,
		PostalAddressLine2: nullstring.ToString(rec.PostalAddressLine2),
		StateProvince:      rec.StateProvince,
		Country:            rec.Country,
		PostalCode:         rec.PostalCode,
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func AccountUserContactRecordToResponse(l logger.Logger, rec *account_record.AccountUserContact) (*account_schema.AccountContactResponse, error) {
	l.Debug("mapping account contact record to response")
	data, err := AccountUserContactRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &account_schema.AccountContactResponse{
		Data: data,
	}, nil
}

func AccountUserContactRecordsToCollectionResponse(l logger.Logger, recs []*account_record.AccountUserContact) (account_schema.AccountContactCollectionResponse, error) {
	l.Debug("mapping account contact records to collection response")
	data := []*account_schema.AccountContactResponseData{}
	for _, rec := range recs {
		d, err := AccountUserContactRecordToResponseData(l, rec)
		if err != nil {
			return account_schema.AccountContactCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return account_schema.AccountContactCollectionResponse{
		Data: data,
	}, nil
}
