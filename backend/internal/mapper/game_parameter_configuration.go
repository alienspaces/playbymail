package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func GameParameterConfigurationRecordToResponseData(l logger.Logger, rec *game_record.GameParameter) (*game_schema.GameParameterConfiguration, error) {
	l.Debug("mapping game parameter configuration record to response data")
	return &game_schema.GameParameterConfiguration{
		GameType:     rec.GameType,
		ConfigKey:    rec.ConfigKey,
		Description:  nullstring.ToStringPtr(rec.Description),
		ValueType:    rec.ValueType,
		DefaultValue: nullstring.ToStringPtr(rec.DefaultValue),
		IsRequired:   rec.IsRequired,
		IsGlobal:     rec.IsGlobal,
	}, nil
}

func GameParameterConfigurationRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameParameter) (game_schema.GameParameterConfigurationCollectionResponse, error) {
	data := []*game_schema.GameParameterConfiguration{}
	for _, rec := range recs {
		d, err := GameParameterConfigurationRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameParameterConfigurationCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameParameterConfigurationCollectionResponse{
		Data: data,
	}, nil
}
