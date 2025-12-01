package harness

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	_ "golang.org/x/image/webp"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// createGameImageRecFromPath creates a game image record from a file path
func (t *Testing) createGameImageRecFromPath(gameID string, recordID string, imagePath string) (*game_record.GameImage, error) {
	l := t.Logger("createGameImageRecFromPath")

	if imagePath == "" {
		return nil, nil
	}

	// TODO: Adopt a standard approach to managing test images
	// The current approach of trying multiple possible testdata locations
	// is fragile and depends on the working directory. A standard approach
	// should be adopted, such as:
	// - Using a well-defined testdata root directory
	// - Using embed.FS to bundle test images
	// - Using a test helper that resolves paths consistently
	// - Defining a standard testdata directory structure
	// Resolve path relative to testdata directory
	// Try multiple possible testdata locations
	testdataPaths := []string{
		"testdata",
		"backend/internal/turn_sheet/testdata",
		"internal/turn_sheet/testdata",
		"turn_sheet/testdata",
	}

	var fullPath string
	var found bool
	for _, basePath := range testdataPaths {
		candidate := filepath.Join(basePath, imagePath)
		if _, err := os.Stat(candidate); err == nil {
			fullPath = candidate
			found = true
			break
		}
	}

	if !found {
		// Try absolute path
		if _, err := os.Stat(imagePath); err == nil {
			fullPath = imagePath
			found = true
		}
	}

	if !found {
		l.Warn("image file not found >%s<", imagePath)
		return nil, fmt.Errorf("image file not found: %s", imagePath)
	}

	// Read image file
	imageData, err := os.ReadFile(fullPath)
	if err != nil {
		l.Warn("failed to read image file >%s< >%v<", fullPath, err)
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	// Detect MIME type
	mimeType := detectMimeType(imageData)
	if mimeType == "" {
		l.Warn("failed to detect MIME type for image >%s<", fullPath)
		return nil, fmt.Errorf("failed to detect MIME type for image: %s", fullPath)
	}

	// Decode image to get dimensions
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		l.Warn("failed to decode image >%s< >%v<", fullPath, err)
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create game image record
	rec := &game_record.GameImage{
		GameID:    gameID,
		RecordID:  nullstring.FromString(recordID),
		Type:      game_record.GameImageTypeTurnSheetBackground,
		ImageData: imageData,
		MimeType:  mimeType,
		FileSize:  len(imageData),
		Width:     width,
		Height:    height,
	}

	l.Debug("creating game image record gameID >%s< recordID >%s< type >%s< size >%d<", gameID, recordID, rec.Type, len(imageData))

	// Upsert game image record
	rec, err = t.Domain.(*domain.Domain).UpsertGameImageRec(rec)
	if err != nil {
		l.Warn("failed creating game image record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameImageRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameImageRec(rec)

	l.Debug("created game image record ID >%s<", rec.ID)

	return rec, nil
}

// detectMimeType detects the MIME type from image data
func detectMimeType(imageData []byte) string {
	if len(imageData) < 12 {
		return ""
	}

	// Check for WebP
	if len(imageData) >= 12 && string(imageData[0:4]) == "RIFF" && string(imageData[8:12]) == "WEBP" {
		return game_record.GameImageMimeTypeWebP
	}

	// Check for PNG
	if len(imageData) >= 8 && string(imageData[0:8]) == "\x89PNG\r\n\x1a\n" {
		return game_record.GameImageMimeTypePNG
	}

	// Check for JPEG
	if len(imageData) >= 3 && string(imageData[0:3]) == "\xff\xd8\xff" {
		return game_record.GameImageMimeTypeJPEG
	}

	return ""
}
