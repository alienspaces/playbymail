package domain

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GetManyGameImageRecs -
func (m *Domain) GetManyGameImageRecs(opts *coresql.Options) ([]*game_record.GameImage, error) {
	l := m.Logger("GetManyGameImageRecs")

	l.Debug("getting many game image records opts >%#v<", opts)

	r := m.GameImageRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameImageRec -
func (m *Domain) GetGameImageRec(recID string, lock *coresql.Lock) (*game_record.GameImage, error) {
	l := m.Logger("GetGameImageRec")

	l.Debug("getting game image record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameImageRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(game_record.TableGameImage, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// GetGameImageRecByGameAndType retrieves a game image by game ID, record ID (optional), type, and turn_sheet_type (optional)
func (m *Domain) GetGameImageRecByGameAndType(gameID string, recordID sql.NullString, imageType string, turnSheetType string) (*game_record.GameImage, error) {
	l := m.Logger("GetGameImageRecByGameAndType")

	l.Debug("getting game image record gameID >%s< recordID >%v< type >%s< turnSheetType >%s<", gameID, recordID, imageType, turnSheetType)

	if err := domain.ValidateUUIDField("game_id", gameID); err != nil {
		return nil, err
	}

	opts := &coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameImageGameID, Val: gameID},
			{Col: game_record.FieldGameImageType, Val: imageType},
		},
	}

	// Handle nullable record_id
	if recordID.Valid {
		opts.Params = append(opts.Params, coresql.Param{Col: game_record.FieldGameImageRecordID, Val: recordID.String})
	} else {
		opts.Params = append(opts.Params, coresql.Param{Col: game_record.FieldGameImageRecordID, Val: nil, Op: coresql.OpIsNull})
	}

	// Handle turn_sheet_type (required for turn_sheet_background type)
	if imageType == game_record.GameImageTypeTurnSheetBackground {
		if turnSheetType == "" {
			return nil, InvalidField(game_record.FieldGameImageTurnSheetType, "", "turn_sheet_type is required when type is turn_sheet_background")
		}
		opts.Params = append(opts.Params, coresql.Param{Col: game_record.FieldGameImageTurnSheetType, Val: turnSheetType})
	} else if turnSheetType != "" {
		return nil, InvalidField(game_record.FieldGameImageTurnSheetType, turnSheetType, "turn_sheet_type must be empty when type is not turn_sheet_background")
	}

	r := m.GameImageRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	if len(recs) == 0 {
		return nil, nil
	}

	return recs[0], nil
}

// CreateGameImageRec -
func (m *Domain) CreateGameImageRec(rec *game_record.GameImage) (*game_record.GameImage, error) {
	l := m.Logger("CreateGameImageRec")

	l.Debug("creating game image record >%#v<", rec)

	r := m.GameImageRepository()

	if err := m.validateGameImageRecForCreate(rec); err != nil {
		l.Warn("failed to validate game image record >%v<", err)
		return rec, err
	}

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateGameImageRec -
func (m *Domain) UpdateGameImageRec(rec *game_record.GameImage) (*game_record.GameImage, error) {
	l := m.Logger("UpdateGameImageRec")

	curr, err := m.GetGameImageRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating game image record >%#v<", rec)

	if err := m.validateGameImageRecForUpdate(rec, curr); err != nil {
		l.Warn("failed to validate game image record >%v<", err)
		return rec, err
	}

	r := m.GameImageRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// UpsertGameImageRec creates or updates a game image record based on game_id, record_id, type, and turn_sheet_type
func (m *Domain) UpsertGameImageRec(rec *game_record.GameImage) (*game_record.GameImage, error) {
	l := m.Logger("UpsertGameImageRec")

	l.Debug("upserting game image record gameID >%s< recordID >%v< type >%s< turnSheetType >%s<", rec.GameID, rec.RecordID, rec.Type, rec.TurnSheetType)

	// Check if an existing record exists
	existing, err := m.GetGameImageRecByGameAndType(rec.GameID, rec.RecordID, rec.Type, rec.TurnSheetType)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Update existing record
		rec.ID = existing.ID
		rec.CreatedAt = existing.CreatedAt
		return m.UpdateGameImageRec(rec)
	}

	// Create new record
	return m.CreateGameImageRec(rec)
}

// DeleteGameImageRec -
func (m *Domain) DeleteGameImageRec(recID string) error {
	l := m.Logger("DeleteGameImageRec")

	l.Debug("deleting game image record ID >%s<", recID)

	rec, err := m.GetGameImageRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameImageRepository()

	if err := m.validateGameImageRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// DeleteGameImageByGameAndType deletes a game image by game ID, record ID (optional), type, and turn_sheet_type (optional)
func (m *Domain) DeleteGameImageByGameAndType(gameID string, recordID sql.NullString, imageType string, turnSheetType string) error {
	l := m.Logger("DeleteGameImageByGameAndType")

	l.Debug("deleting game image record gameID >%s< recordID >%v< type >%s< turnSheetType >%s<", gameID, recordID, imageType, turnSheetType)

	rec, err := m.GetGameImageRecByGameAndType(gameID, recordID, imageType, turnSheetType)
	if err != nil {
		return err
	}

	if rec == nil {
		return nil // Nothing to delete
	}

	return m.DeleteGameImageRec(rec.ID)
}

// RemoveGameImageRec -
func (m *Domain) RemoveGameImageRec(recID string) error {
	l := m.Logger("RemoveGameImageRec")

	l.Debug("removing game image record ID >%s<", recID)

	rec, err := m.GetGameImageRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameImageRepository()

	if err := m.validateGameImageRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// GetGameImageAsBase64DataURL returns the image as a base64-encoded data URL for embedding in HTML
func (m *Domain) GetGameImageAsBase64DataURL(gameID string, recordID sql.NullString, imageType string, turnSheetType string) (string, error) {
	l := m.Logger("GetGameImageAsBase64DataURL")

	l.Debug("getting game image as base64 data URL gameID >%s< recordID >%v< type >%s< turnSheetType >%s<", gameID, recordID, imageType, turnSheetType)

	rec, err := m.GetGameImageRecByGameAndType(gameID, recordID, imageType, turnSheetType)
	if err != nil {
		return "", err
	}

	if rec == nil {
		return "", nil
	}

	return imageToBase64DataURL(rec.ImageData, rec.MimeType), nil
}

// GetGameTurnSheetBackgroundImage retrieves the turn sheet background image for a game (game-level, record_id is NULL)
// turnSheetType is required when type is turn_sheet_background
func (m *Domain) GetGameTurnSheetBackgroundImage(gameID string, turnSheetType string) (*game_record.GameImage, error) {
	l := m.Logger("GetGameTurnSheetBackgroundImage")

	l.Debug("getting turn sheet background image for game >%s< turnSheetType >%s<", gameID, turnSheetType)

	if err := domain.ValidateUUIDField("game_id", gameID); err != nil {
		return nil, err
	}

	return m.GetGameImageRecByGameAndType(gameID, sql.NullString{}, game_record.GameImageTypeTurnSheetBackground, turnSheetType)
}

// GetGameTurnSheetImageDataURL retrieves the turn sheet background image for a game as base64 data URL
// turnSheetType is required when type is turn_sheet_background
func (m *Domain) GetGameTurnSheetImageDataURL(gameID string, turnSheetType string) (string, error) {
	l := m.Logger("GetGameTurnSheetImageDataURL")

	l.Info("getting turn sheet background image data URL for game >%s< turnSheetType >%s<", gameID, turnSheetType)

	img, err := m.GetGameTurnSheetBackgroundImage(gameID, turnSheetType)
	if err != nil {
		l.Warn("failed to get turn sheet background image >%v<", err)
		return "", err
	}

	if img == nil {
		l.Info("no background image found for game >%s< turnSheetType >%s<", gameID, turnSheetType)
		return "", nil
	}

	l.Info("processing image type >%s< size >%d< bytes mime >%s<", img.Type, len(img.ImageData), img.MimeType)
	dataURL := imageToBase64DataURL(img.ImageData, img.MimeType)
	prefixLen := 100
	if len(dataURL) < prefixLen {
		prefixLen = len(dataURL)
	}
	l.Info("generated data URL length >%d< prefix >%s<", len(dataURL), dataURL[:prefixLen])

	return dataURL, nil
}

// GetAdventureGameLocationChoiceTurnSheetImageDataURL retrieves the turn sheet background image for a location as base64 data URL
func (m *Domain) GetAdventureGameLocationChoiceTurnSheetImageDataURL(gameID, locationID string) (string, error) {
	l := m.Logger("GetAdventureGameLocationChoiceTurnSheetImageDataURL")

	l.Debug("getting turn sheet background image data URL for location >%s< in game >%s<", locationID, gameID)

	recordID := nullstring.FromString(locationID)
	turnSheetType := adventure_game_record.AdventureGameTurnSheetTypeLocationChoice

	img, err := m.GetGameImageRecByGameAndType(gameID, recordID, game_record.GameImageTypeTurnSheetBackground, turnSheetType)
	if err != nil {
		return "", err
	}

	// If location has a background image, use it
	if img != nil {
		dataURL := imageToBase64DataURL(img.ImageData, img.MimeType)
		return dataURL, nil
	}

	// Fall back to game-level background image
	return m.GetGameTurnSheetImageDataURL(gameID, turnSheetType)
}

// GetAdventureGameInventoryTurnSheetImageDataURL retrieves the turn sheet background image for inventory as base64 data URL
func (m *Domain) GetAdventureGameInventoryTurnSheetImageDataURL(gameID string) (string, error) {
	l := m.Logger("GetAdventureGameInventoryTurnSheetImageDataURL")

	l.Debug("getting turn sheet background image data URL for inventory in game >%s<", gameID)

	turnSheetType := adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement

	img, err := m.GetGameImageRecByGameAndType(gameID, sql.NullString{}, game_record.GameImageTypeTurnSheetBackground, turnSheetType)
	if err != nil {
		return "", err
	}

	// If inventory has a background image, use it
	if img != nil {
		dataURL := imageToBase64DataURL(img.ImageData, img.MimeType)
		return dataURL, nil
	}

	// Fall back to game-level background image
	return m.GetGameTurnSheetImageDataURL(gameID, turnSheetType)
}
