package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func GameParameterRecordToResponseData(l logger.Logger, rec *game_record.GameParameter) (*game_schema.GameParameter, error) {
	l.Debug("mapping game parameter record to response data")
	return &game_schema.GameParameter{
		GameType:     rec.GameType,
		ConfigKey:    rec.ConfigKey,
		Description:  &rec.Description,
		ValueType:    rec.ValueType,
		DefaultValue: &rec.DefaultValue,
	}, nil
}

func GameParameterRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameParameter) (game_schema.GameParameterCollectionResponse, error) {
	l.Debug("mapping game parameter records to collection response")
	data := []*game_schema.GameParameter{}
	for _, rec := range recs {
		d, err := GameParameterRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameParameterCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameParameterCollectionResponse{
		Data: data,
	}, nil
}
