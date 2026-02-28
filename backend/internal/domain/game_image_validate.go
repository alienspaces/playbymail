package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type validateGameImageArgs struct {
	nextRec *game_record.GameImage
	currRec *game_record.GameImage
}

func (m *Domain) populateGameImageValidateArgs(currRec, nextRec *game_record.GameImage) (*validateGameImageArgs, error) {
	args := &validateGameImageArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateGameImageRecForCreate(rec *game_record.GameImage) error {
	args, err := m.populateGameImageValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateGameImageRecForCreate(args)
}

func (m *Domain) validateGameImageRecForUpdate(currRec, nextRec *game_record.GameImage) error {
	args, err := m.populateGameImageValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateGameImageRecForUpdate(args)
}

func (m *Domain) validateGameImageRecForDelete(rec *game_record.GameImage) error {
	args, err := m.populateGameImageValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateGameImageRecForDelete(args)
}

func validateGameImageRecForCreate(args *validateGameImageArgs) error {
	return validateGameImageRec(args)
}

func validateGameImageRecForUpdate(args *validateGameImageArgs) error {
	return validateGameImageRec(args)
}

func validateGameImageRec(args *validateGameImageArgs) error {
	rec := args.nextRec

	// Validate game_id
	if err := domain.ValidateUUIDField(game_record.FieldGameImageGameID, rec.GameID); err != nil {
		return err
	}

	// Validate record_id if present
	if rec.RecordID.Valid {
		if err := domain.ValidateUUIDField(game_record.FieldGameImageRecordID, rec.RecordID.String); err != nil {
			return err
		}
	}

	// Validate type
	if err := domain.ValidateEnumField(game_record.FieldGameImageType, rec.Type, game_record.GameImageTypes); err != nil {
		return err
	}

	// Validate turn_sheet_type: required when type is turn_sheet_background
	if rec.Type == game_record.GameImageTypeTurnSheetBackground {
		if rec.TurnSheetType == "" {
			return InvalidField(game_record.FieldGameImageTurnSheetType, "", "turn_sheet_type is required when type is turn_sheet_background")
		}
	} else {
		// turn_sheet_type must be empty for non-turn-sheet-background types
		if rec.TurnSheetType != "" {
			return InvalidField(game_record.FieldGameImageTurnSheetType, rec.TurnSheetType, "turn_sheet_type must be empty when type is not turn_sheet_background")
		}
	}

	// Validate mime_type
	if err := domain.ValidateEnumField(game_record.FieldGameImageMimeType, rec.MimeType, game_record.GameImageMimeTypes); err != nil {
		return err
	}

	// Validate image_data is present
	if len(rec.ImageData) == 0 {
		return InvalidField(game_record.FieldGameImageImageData, "", "image data is required")
	}

	// Validate file_size
	if rec.FileSize <= 0 {
		return InvalidField(game_record.FieldGameImageFileSize, fmt.Sprintf("%d", rec.FileSize), "file size must be greater than 0")
	}
	if rec.FileSize > game_record.GameImageMaxSize {
		return InvalidField(game_record.FieldGameImageFileSize, fmt.Sprintf("%d", rec.FileSize), fmt.Sprintf("file size exceeds maximum of %d bytes (1MB)", game_record.GameImageMaxSize))
	}

	// Validate width
	if rec.Width < game_record.GameImageMinWidth {
		return InvalidField(game_record.FieldGameImageWidth, fmt.Sprintf("%d", rec.Width), fmt.Sprintf("width must be at least %d pixels", game_record.GameImageMinWidth))
	}
	if rec.Width > game_record.GameImageMaxWidth {
		return InvalidField(game_record.FieldGameImageWidth, fmt.Sprintf("%d", rec.Width), fmt.Sprintf("width must not exceed %d pixels", game_record.GameImageMaxWidth))
	}

	// Validate height
	if rec.Height < game_record.GameImageMinHeight {
		return InvalidField(game_record.FieldGameImageHeight, fmt.Sprintf("%d", rec.Height), fmt.Sprintf("height must be at least %d pixels", game_record.GameImageMinHeight))
	}
	if rec.Height > game_record.GameImageMaxHeight {
		return InvalidField(game_record.FieldGameImageHeight, fmt.Sprintf("%d", rec.Height), fmt.Sprintf("height must not exceed %d pixels", game_record.GameImageMaxHeight))
	}

	return nil
}

func validateGameImageRecForDelete(args *validateGameImageArgs) error {
	rec := args.nextRec

	if err := domain.ValidateUUIDField(game_record.FieldGameImageID, rec.ID); err != nil {
		return err
	}

	return nil
}

// ValidateImageDimensions checks if image dimensions are within acceptable bounds
// Returns a warning message if aspect ratio differs significantly from A4
func ValidateImageDimensions(width, height int) (valid bool, warning string) {
	if width < game_record.GameImageMinWidth || width > game_record.GameImageMaxWidth {
		return false, ""
	}
	if height < game_record.GameImageMinHeight || height > game_record.GameImageMaxHeight {
		return false, ""
	}

	// Check aspect ratio (A4 is ~1:1.414)
	expectedRatio := 1.414
	actualRatio := float64(height) / float64(width)
	tolerance := 0.15 // 15% tolerance

	if actualRatio < expectedRatio*(1-tolerance) || actualRatio > expectedRatio*(1+tolerance) {
		warning = fmt.Sprintf("Image aspect ratio (%.2f) differs from A4 ratio (1.41). "+
			"Recommended: %dx%d pixels (A4 @ 300 DPI) for best print quality.",
			actualRatio, game_record.GameImageRecommendedWidth, game_record.GameImageRecommendedHeight)
	}

	return true, warning
}
