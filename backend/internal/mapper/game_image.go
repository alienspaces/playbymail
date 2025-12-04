package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func GameImageRecordToResponseData(l logger.Logger, rec *game_record.GameImage) (*game_schema.GameImageResponseData, error) {
	l.Debug("mapping game image record to response data")
	return &game_schema.GameImageResponseData{
		ID:            rec.ID,
		GameID:        rec.GameID,
		RecordID:      nullstring.ToString(rec.RecordID),
		Type:          rec.Type,
		TurnSheetType: rec.TurnSheetType,
		MimeType:      rec.MimeType,
		FileSize:      rec.FileSize,
		Width:         rec.Width,
		Height:        rec.Height,
		CreatedAt:     rec.CreatedAt,
		UpdatedAt:     nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func GameImageRecordToResponse(l logger.Logger, rec *game_record.GameImage, warning string) (*game_schema.GameImageResponse, error) {
	l.Debug("mapping game image record to response")
	data, err := GameImageRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	data.Warning = warning
	return &game_schema.GameImageResponse{
		Data: data,
	}, nil
}

func GameImageRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameImage) (game_schema.GameImageCollectionResponse, error) {
	l.Debug("mapping game image records to collection response")
	data := []*game_schema.GameImageResponseData{}
	for _, rec := range recs {
		d, err := GameImageRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameImageCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameImageCollectionResponse{
		Data: data,
	}, nil
}

func GameTurnSheetImageToResponse(l logger.Logger, gameID string, img *game_record.GameImage) (*game_schema.GameTurnSheetImageResponse, error) {
	l.Debug("mapping game turn sheet image to response")

	data := &game_schema.GameTurnSheetImageData{
		GameID: gameID,
	}

	if img != nil {
		responseData, err := GameImageRecordToResponseData(l, img)
		if err != nil {
			return nil, err
		}
		data.Background = responseData
	}

	return &game_schema.GameTurnSheetImageResponse{
		Data: data,
	}, nil
}

func LocationTurnSheetImageToResponse(l logger.Logger, gameID, locationID string, img *game_record.GameImage) (*game_schema.LocationTurnSheetImageResponse, error) {
	l.Debug("mapping location turn sheet image to response")

	data := &game_schema.LocationTurnSheetImageData{
		GameID:     gameID,
		LocationID: locationID,
	}

	if img != nil {
		responseData, err := GameImageRecordToResponseData(l, img)
		if err != nil {
			return nil, err
		}
		data.Background = responseData
	}

	return &game_schema.LocationTurnSheetImageResponse{
		Data: data,
	}, nil
}
