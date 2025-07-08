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

func GameRequestToRecord(l logger.Logger, r *http.Request, rec *record.Game) (*record.Game, error) {

	l.Debug("mapping game request to record")

	var data schema.GameRequest
	_, err := server.ReadRequest(l, r, &data)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Name = data.Name
		rec.GameType = data.GameType
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = data.Name
		rec.GameType = data.GameType
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameRecordToResponseData(l logger.Logger, rec *record.Game) (schema.GameResponseData, error) {

	l.Debug("mapping game record to response data")

	data := schema.GameResponseData{
		ID:        rec.ID,
		Name:      rec.Name,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
	}

	return data, nil
}
