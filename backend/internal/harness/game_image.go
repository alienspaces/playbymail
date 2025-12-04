package harness

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	_ "golang.org/x/image/webp"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// CreateTestImage creates a simple PNG image for testing with visible content
// This is a helper function for creating test image data in harness configs
func CreateTestImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a light gray background so the image is visible
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Light gray color (RGB: 220, 220, 220, fully opaque)
			img.Set(x, y, image.White)
		}
	}

	// Add a simple border pattern to make it obvious it's a test image
	borderWidth := 20
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Draw a dark gray border
			if x < borderWidth || x >= width-borderWidth ||
				y < borderWidth || y >= height-borderWidth {
				img.Set(x, y, image.Black)
			}
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

// LoadTestImageFromPath loads an image from the testdata or seed_images directory
// This is a helper function for loading real test images in harness configs
func LoadTestImageFromPath(imagePath string) ([]byte, string, int, int) {
	// Get current working directory for debugging
	cwd, _ := os.Getwd()

	// Try multiple possible image locations
	// Include paths relative to common project structures
	imagePaths := []string{
		// Seed data images (primary location for seed data)
		"backend/internal/runner/cli/seed_images",
		"internal/runner/cli/seed_images",
		// Test data images (for unit tests)
		"testdata",
		"backend/internal/turn_sheet/testdata",
		"internal/turn_sheet/testdata",
		"turn_sheet/testdata",
	}

	// Also try paths based on TEMPLATES_PATH if set (indicates we're in backend/)
	if templatesPath := os.Getenv("TEMPLATES_PATH"); templatesPath != "" {
		// TEMPLATES_PATH is typically .../backend/templates
		backendDir := filepath.Dir(templatesPath)
		imagePaths = append(imagePaths,
			filepath.Join(backendDir, "internal/runner/cli/seed_images"),
			filepath.Join(backendDir, "internal/turn_sheet/testdata"),
		)
	}

	var fullPath string
	var found bool
	for _, basePath := range imagePaths {
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
		// Log that we couldn't find the image
		fmt.Printf("LoadTestImageFromPath: image not found >%s< cwd >%s< tried paths: %v\n", imagePath, cwd, imagePaths)
		// Return empty data - caller should handle this
		return nil, "", 0, 0
	}

	fmt.Printf("LoadTestImageFromPath: found image at >%s<\n", fullPath)

	// Read image file
	imageData, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, "", 0, 0
	}

	// Determine MIME type from extension
	ext := filepath.Ext(imagePath)
	var mimeType string
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".webp":
		mimeType = "image/webp"
	default:
		mimeType = "image/png"
	}

	// Decode image to get dimensions
	reader := bytes.NewReader(imageData)
	img, _, err := image.Decode(reader)
	if err != nil {
		return imageData, mimeType, 0, 0
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	return imageData, mimeType, width, height
}

// createGameImageRecFromConfig creates a game image record from a GameImageConfig
func (t *Testing) createGameImageRecFromConfig(gameID string, config GameImageConfig) (*game_record.GameImage, error) {
	l := t.Logger("createGameImageRecFromConfig")

	// If ImagePath is specified, load the image from file at runtime
	if config.ImagePath != "" {
		l.Info("loading image from path >%s< for turnSheetType >%s<", config.ImagePath, config.TurnSheetType)
		imageData, mimeType, width, height, err := t.loadImageFromPath(config.ImagePath)
		if err != nil {
			l.Warn("failed loading image from path >%s<: >%v<", config.ImagePath, err)
			return nil, err
		}

		rec := &game_record.GameImage{
			GameID:        gameID,
			Type:          game_record.GameImageTypeTurnSheetBackground,
			TurnSheetType: config.TurnSheetType,
			ImageData:     imageData,
			MimeType:      mimeType,
			FileSize:      len(imageData),
			Width:         width,
			Height:        height,
		}

		l.Info("creating game image record from path: gameID >%s< type >%s< turnSheetType >%s< mimeType >%s< imageDataLen >%d< width >%d< height >%d<",
			gameID, rec.Type, rec.TurnSheetType, rec.MimeType, len(rec.ImageData), rec.Width, rec.Height)

		rec, err = t.Domain.(*domain.Domain).UpsertGameImageRec(rec)
		if err != nil {
			l.Warn("failed creating game image record >%v<", err)
			return nil, err
		}

		l.Info("created game image record ID >%s< with imageDataLen >%d<", rec.ID, len(rec.ImageData))

		// Add to data store
		t.Data.AddGameImageRec(rec)

		// Add to teardown data store
		t.teardownData.AddGameImageRec(rec)

		// Add to references store
		if config.Reference != "" {
			t.Data.Refs.GameImageRefs[config.Reference] = rec.ID
		}

		return rec, nil
	}

	// Otherwise use the Record from config
	var rec *game_record.GameImage
	if config.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *config.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameImage{}
	}

	// Set game ID
	rec.GameID = gameID

	// Set FileSize from ImageData if not already set
	if rec.FileSize == 0 && len(rec.ImageData) > 0 {
		rec.FileSize = len(rec.ImageData)
	}

	// Apply default values if needed
	rec = t.applyGameImageRecDefaultValues(rec)

	// Debug logging for image data
	l.Info("creating game image record: gameID >%s< type >%s< turnSheetType >%s< mimeType >%s< imageDataLen >%d< width >%d< height >%d<",
		gameID, rec.Type, rec.TurnSheetType, rec.MimeType, len(rec.ImageData), rec.Width, rec.Height)

	rec, err := t.Domain.(*domain.Domain).UpsertGameImageRec(rec)
	if err != nil {
		l.Warn("failed creating game image record >%v<", err)
		return nil, err
	}

	l.Info("created game image record ID >%s< with imageDataLen >%d<", rec.ID, len(rec.ImageData))

	// Add to data store
	t.Data.AddGameImageRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameImageRec(rec)

	// Add to references store
	if config.Reference != "" {
		t.Data.Refs.GameImageRefs[config.Reference] = rec.ID
	}

	l.Debug("created game image record ID >%s<", rec.ID)

	return rec, nil
}

func (t *Testing) applyGameImageRecDefaultValues(rec *game_record.GameImage) *game_record.GameImage {
	if rec == nil {
		rec = &game_record.GameImage{}
	}

	// No default values needed - all fields should be set explicitly
	// or loaded from image data if creating from path

	return rec
}

// loadImageFromPath loads an image from a file path and returns the data,
// MIME type, and dimensions
func (t *Testing) loadImageFromPath(imagePath string) ([]byte, string, int, int, error) {
	l := t.Logger("loadImageFromPath")

	// Resolve path relative to known image directories
	// Try multiple possible locations for test/seed images
	imagePaths := []string{
		// Seed data images (primary location for seed data)
		"backend/internal/runner/cli/seed_images",
		"internal/runner/cli/seed_images",
		// Test data images (for unit tests)
		"testdata",
		"backend/internal/turn_sheet/testdata",
		"internal/turn_sheet/testdata",
		"turn_sheet/testdata",
		// Relative paths from harness package to turn_sheet testdata
		"../turn_sheet/testdata",
	}

	var fullPath string
	var found bool
	for _, basePath := range imagePaths {
		candidate := filepath.Join(basePath, imagePath)
		if _, err := os.Stat(candidate); err == nil {
			fullPath = candidate
			found = true
			l.Info("found image at >%s<", fullPath)
			break
		}
	}

	if !found {
		// Try absolute path
		if _, err := os.Stat(imagePath); err == nil {
			fullPath = imagePath
			found = true
			l.Info("found image at absolute path >%s<", fullPath)
		}
	}

	if !found {
		l.Warn("image file not found >%s<", imagePath)
		return nil, "", 0, 0, fmt.Errorf("image file not found: %s", imagePath)
	}

	// Read image file
	imageData, err := os.ReadFile(fullPath)
	if err != nil {
		l.Warn("failed to read image file >%s< >%v<", fullPath, err)
		return nil, "", 0, 0, fmt.Errorf("failed to read image file: %w", err)
	}

	// Detect MIME type
	mimeType := detectMimeType(imageData)
	if mimeType == "" {
		l.Warn("failed to detect MIME type for image >%s<", fullPath)
		return nil, "", 0, 0, fmt.Errorf("failed to detect MIME type for image: %s", fullPath)
	}

	// Decode image to get dimensions
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		l.Warn("failed to decode image >%s< >%v<", fullPath, err)
		return nil, "", 0, 0, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	return imageData, mimeType, width, height, nil
}

// createGameImageRecFromPath creates a game image record from a file path
// DEPRECATED: Use createGameImageRecFromConfig instead
// recordID: if empty, this is a game-level image (defaults to join_game), otherwise it's a location-level image (location_choice)
func (t *Testing) createGameImageRecFromPath(gameID string, recordID string, imagePath string) (*game_record.GameImage, error) {
	l := t.Logger("createGameImageRecFromPath")

	if imagePath == "" {
		return nil, nil
	}

	// Load image data from path
	imageData, mimeType, width, height, err := t.loadImageFromPath(imagePath)
	if err != nil {
		l.Warn("failed to load image from path >%s< >%v<", imagePath, err)
		return nil, err
	}

	// Determine turn_sheet_type based on context
	// If recordID is empty, this is a game-level image (for join_game turn sheets)
	// If recordID is set, this is a location-level image (for location_choice turn sheets)
	turnSheetType := adventure_game_record.AdventureGameTurnSheetTypeJoinGame
	if recordID != "" {
		turnSheetType = adventure_game_record.AdventureGameTurnSheetTypeLocationChoice
	}

	// Use the new config-based approach with Record
	config := GameImageConfig{
		Record: &game_record.GameImage{
			GameID:        gameID,
			RecordID:      nullstring.FromString(recordID),
			Type:          game_record.GameImageTypeTurnSheetBackground,
			TurnSheetType: turnSheetType,
			ImageData:     imageData,
			MimeType:      mimeType,
			FileSize:      len(imageData),
			Width:         width,
			Height:        height,
		},
	}
	return t.createGameImageRecFromConfig(gameID, config)
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
