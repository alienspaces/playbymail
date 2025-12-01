package adventure_game

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"
	_ "golang.org/x/image/webp"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheet"
)

const (
	UploadLocationTurnSheetImage   = "upload-location-turn-sheet-image"
	GetLocationTurnSheetImage      = "get-location-turn-sheet-image"
	DeleteLocationTurnSheetImage   = "delete-location-turn-sheet-image"
	PreviewLocationChoiceTurnSheet = "preview-location-choice-turn-sheet"
)

func adventureGameLocationTurnSheetImageHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameLocationTurnSheetImageHandlerConfig")

	l.Debug("Adding adventure game location turn sheet image handler configuration")

	locationImageConfig := make(map[string]server.HandlerConfig)

	// Upload turn sheet image endpoint for location
	locationImageConfig[UploadLocationTurnSheetImage] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/locations/:location_id/turn-sheet-image",
		HandlerFunc: uploadLocationTurnSheetImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Upload location turn sheet image",
			Description: "Upload a turn sheet background image for an adventure game location. " +
				"Accepts multipart form data with 'image' file. " +
				"Images must be WebP, PNG, or JPEG format, max 1MB. " +
				"Recommended: 2480x3508px (A4 @ 300 DPI) for best print quality. " +
				"Use the preview endpoint to see how images appear in the turn sheet.",
		},
	}

	// Get the turn sheet image for a location
	locationImageConfig[GetLocationTurnSheetImage] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/locations/:location_id/turn-sheet-image",
		HandlerFunc: getLocationTurnSheetImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Get location turn sheet image",
			Description: "Get the turn sheet background image for an adventure game location.",
		},
	}

	// Delete the turn sheet image for a location
	locationImageConfig[DeleteLocationTurnSheetImage] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/locations/:location_id/turn-sheet-image",
		HandlerFunc: deleteLocationTurnSheetImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Delete location turn sheet image",
			Description: "Delete the turn sheet background image for an adventure game location.",
		},
	}

	// Preview location choice turn sheet with uploaded images
	locationImageConfig[PreviewLocationChoiceTurnSheet] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/locations/:location_id/turn-sheets/preview",
		HandlerFunc: previewLocationChoiceTurnSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Preview location choice turn sheet",
			Description: "Generate and preview a location choice turn sheet PDF with uploaded " +
				"background images. Returns PDF with Content-Disposition: inline for browser preview.",
		},
	}

	return locationImageConfig, nil
}

// uploadLocationTurnSheetImageHandler handles uploading a turn sheet background image for a location
func uploadLocationTurnSheetImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "uploadLocationTurnSheetImageHandler")

	l.Info("uploading location turn sheet image with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	locationID := pp.ByName("location_id")
	if locationID == "" {
		l.Warn("location ID is empty")
		return coreerror.RequiredPathParameter("location_id")
	}

	l.Info("uploading turn sheet image for location >%s< in game >%s<", locationID, gameID)

	// Verify game exists and user has access (RLS check)
	_, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	// Verify location exists and belongs to the game
	locationRec, err := mm.GetAdventureGameLocationRec(locationID, nil)
	if err != nil {
		l.Warn("failed to get location record >%s< >%v<", locationID, err)
		return err
	}

	if locationRec.GameID != gameID {
		l.Warn("location does not belong to specified game >%s< != >%s<", locationRec.GameID, gameID)
		return coreerror.NewNotFoundError("location", locationID)
	}

	// Parse multipart form (max 2MB to account for form overhead)
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		l.Warn("failed to parse multipart form >%v<", err)
		return coreerror.NewInvalidDataError("failed to parse multipart form: %v", err)
	}

	l.Info("processing image upload for location >%s<", locationID)

	// Get file from form
	file, header, err := r.FormFile("image")
	if err != nil {
		l.Warn("failed to get image file >%v<", err)
		return coreerror.NewInvalidDataError("image file is required")
	}
	defer file.Close()

	l.Info("received file >%s<", header.Filename)

	// Read file data
	imageData, err := io.ReadAll(file)
	if err != nil {
		l.Warn("failed to read image data >%v<", err)
		return coreerror.NewInvalidDataError("failed to read image data")
	}

	l.Info("read image data size >%d< bytes", len(imageData))

	// Validate file size
	if len(imageData) > game_record.GameImageMaxSize {
		l.Warn("image file too large >%d< bytes, max >%d<", len(imageData), game_record.GameImageMaxSize)
		return coreerror.NewInvalidDataError("image file too large, max 1MB")
	}

	// Detect MIME type
	mimeType := http.DetectContentType(imageData)
	l.Info("detected MIME type >%s< from file >%s<", mimeType, header.Filename)

	// Normalize MIME type (http.DetectContentType may return variations)
	switch mimeType {
	case "image/webp":
		mimeType = game_record.GameImageMimeTypeWebP
	case "image/png":
		mimeType = game_record.GameImageMimeTypePNG
	case "image/jpeg":
		mimeType = game_record.GameImageMimeTypeJPEG
	default:
		l.Warn("invalid MIME type >%s<", mimeType)
		return coreerror.NewInvalidDataError("invalid image format, must be WebP, PNG, or JPEG")
	}

	// Decode image to get dimensions
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		l.Warn("failed to decode image >%v<", err)
		return coreerror.NewInvalidDataError("failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	l.Info("image dimensions: width >%d< height >%d<", width, height)

	// Validate dimensions
	valid, warning := domain.ValidateImageDimensions(width, height)
	if !valid {
		l.Warn("invalid image dimensions: width >%d< height >%d<", width, height)
		return coreerror.NewInvalidDataError(
			"invalid image dimensions: width must be %d-%d pixels, height must be %d-%d pixels",
			game_record.GameImageMinWidth, game_record.GameImageMaxWidth,
			game_record.GameImageMinHeight, game_record.GameImageMaxHeight,
		)
	}

	// Create or update image record with location ID as record_id
	rec := &game_record.GameImage{
		GameID:    gameID,
		RecordID:  nullstring.FromString(locationID),
		Type:      game_record.GameImageTypeTurnSheetBackground,
		ImageData: imageData,
		MimeType:  mimeType,
		FileSize:  len(imageData),
		Width:     width,
		Height:    height,
	}

	l.Info("upserting location image record for location >%s<", locationID)

	rec, err = mm.UpsertGameImageRec(rec)
	if err != nil {
		l.Warn("failed to upsert location image record >%v<", err)
		return err
	}

	l.Info("successfully saved location image record id >%s<", rec.ID)

	res, err := mapper.GameImageRecordToResponse(l, rec, warning)
	if err != nil {
		l.Warn("failed to map location image record to response >%v<", err)
		return err
	}

	l.Info("responding with location image record id >%s< warning >%s<", rec.ID, warning)

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// getLocationTurnSheetImageHandler gets the turn sheet background image for a location
func getLocationTurnSheetImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getLocationTurnSheetImageHandler")

	l.Info("getting location turn sheet image with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	locationID := pp.ByName("location_id")
	if locationID == "" {
		l.Warn("location ID is empty")
		return coreerror.RequiredPathParameter("location_id")
	}

	l.Info("fetching turn sheet image for location >%s< in game >%s<", locationID, gameID)

	// Verify game exists and user has access (RLS check)
	_, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	// Verify location exists and belongs to the game
	locationRec, err := mm.GetAdventureGameLocationRec(locationID, nil)
	if err != nil {
		l.Warn("failed to get location record >%s< >%v<", locationID, err)
		return err
	}

	if locationRec.GameID != gameID {
		l.Warn("location does not belong to specified game >%s< != >%s<", locationRec.GameID, gameID)
		return coreerror.NewNotFoundError("location", locationID)
	}

	// Get image for this specific location
	recordID := nullstring.FromString(locationID)
	img, err := mm.GetGameImageRecByGameAndType(gameID, recordID, game_record.GameImageTypeTurnSheetBackground)
	if err != nil {
		l.Warn("failed to get location turn sheet image >%v<", err)
		return err
	}

	res, err := mapper.LocationTurnSheetImageToResponse(l, gameID, locationID, img)
	if err != nil {
		l.Warn("failed to map location turn sheet image to response >%v<", err)
		return err
	}

	l.Info("responding with location turn sheet image for location >%s<", locationID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// deleteLocationTurnSheetImageHandler deletes the turn sheet background image for a location
func deleteLocationTurnSheetImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteLocationTurnSheetImageHandler")

	l.Info("deleting location turn sheet image with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	locationID := pp.ByName("location_id")
	if locationID == "" {
		l.Warn("location ID is empty")
		return coreerror.RequiredPathParameter("location_id")
	}

	// Verify game exists and user has access (RLS check)
	_, err := mm.GetGameRec(gameID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	// Verify location exists and belongs to the game
	locationRec, err := mm.GetAdventureGameLocationRec(locationID, nil)
	if err != nil {
		l.Warn("failed to get location record >%s< >%v<", locationID, err)
		return err
	}

	if locationRec.GameID != gameID {
		l.Warn("location does not belong to specified game >%s< != >%s<", locationRec.GameID, gameID)
		return coreerror.NewNotFoundError("location", locationID)
	}

	recordID := nullstring.FromString(locationID)
	err = mm.DeleteGameImageByGameAndType(gameID, recordID, game_record.GameImageTypeTurnSheetBackground)
	if err != nil {
		l.Warn("failed to delete location image >%v<", err)
		return err
	}

	l.Info("deleted location turn sheet image for location >%s<", locationID)

	if err = server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// previewLocationChoiceTurnSheetHandler generates and previews a location choice turn sheet PDF
func previewLocationChoiceTurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "previewLocationChoiceTurnSheetHandler")

	l.Info("previewing location choice turn sheet with path params >%#v<", pp)

	ctx := r.Context()
	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	locationID := pp.ByName("location_id")
	if locationID == "" {
		l.Warn("location ID is empty")
		return coreerror.RequiredPathParameter("location_id")
	}

	l.Info("generating turn sheet preview for location >%s< in game >%s<", locationID, gameID)

	// Get game record (RLS check)
	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	l.Info("found game >%s< type >%s<", gameRec.Name, gameRec.GameType)

	// Validate game type supports location choice sheets
	if gameRec.GameType != game_record.GameTypeAdventure {
		l.Warn("location choice turn sheet only supported for adventure games, got >%s<", gameRec.GameType)
		return coreerror.NewInvalidDataError("location choice turn sheet is not supported for game type %s", gameRec.GameType)
	}

	// Get location record
	locationRec, err := mm.GetAdventureGameLocationRec(locationID, nil)
	if err != nil {
		l.Warn("failed to get location record >%s< >%v<", locationID, err)
		return err
	}

	if locationRec.GameID != gameID {
		l.Warn("location does not belong to specified game >%s< != >%s<", locationRec.GameID, gameID)
		return coreerror.NewNotFoundError("location", locationID)
	}

	l.Info("found location >%s<", locationRec.Name)

	// Get location links (available destinations from this location)
	locationOptions, err := getLocationOptionsForPreview(ctx, mm, gameID, locationID)
	if err != nil {
		l.Warn("failed to get location options >%v<", err)
		return err
	}

	// Parse config for processor
	cfg, err := config.Parse()
	if err != nil {
		l.Warn("failed to parse config >%v<", err)
		return coreerror.NewInternalError("failed to parse config: %v", err)
	}

	// Generate a realistic preview turn sheet code
	// Use placeholder UUIDs for game instance, account, and turn sheet since this is a preview
	previewGameInstanceID := "00000000-0000-0000-0000-000000000000"
	previewAccountID := "00000000-0000-0000-0000-000000000000"
	previewTurnSheetID := "00000000-0000-0000-0000-000000000000"
	previewCode, err := turnsheet.GenerateTurnSheetCode(gameID, previewGameInstanceID, previewAccountID, previewTurnSheetID)
	if err != nil {
		l.Warn("failed to generate preview turn sheet code >%v<", err)
		return coreerror.NewInternalError("failed to generate turn sheet code: %v", err)
	}

	// Create location choice data
	gameName := gameRec.Name
	gameType := gameRec.GameType
	instructions := turn_sheet.DefaultLocationChoiceInstructions()
	previewTurnNumber := 1 // Placeholder turn number for preview

	locationChoiceData := &turn_sheet.LocationChoiceData{
		TurnSheetTemplateData: turn_sheet.TurnSheetTemplateData{
			GameName:              &gameName,
			GameType:              &gameType,
			TurnSheetTitle:        &locationRec.Name,
			TurnSheetDescription:  &locationRec.Description,
			TurnSheetInstructions: &instructions,
			TurnNumber:            &previewTurnNumber,
			TurnSheetCode:         &previewCode,
		},
		LocationName:        locationRec.Name,
		LocationDescription: locationRec.Description,
		LocationOptions:     locationOptions,
	}

	// Get uploaded turn sheet background image for this location
	// Falls back to game-level image if location doesn't have one
	backgroundImage, err := mm.GetLocationTurnSheetImageDataURL(gameID, locationID)
	if err != nil {
		l.Warn("failed to get turn sheet background image >%v<", err)
		// Continue without image - not a fatal error
	} else if backgroundImage != "" {
		locationChoiceData.BackgroundImage = &backgroundImage
		l.Info("loaded background image for turn sheet, length >%d<", len(backgroundImage))
	} else {
		l.Info("no background image found for turn sheet")
	}

	// Marshal location choice data to JSON
	sheetDataBytes, err := json.Marshal(locationChoiceData)
	if err != nil {
		l.Warn("failed to marshal location choice sheet data >%v<", err)
		return coreerror.NewInternalError("failed to marshal location choice sheet data")
	}

	l.Info("marshaled location choice data size >%d< bytes", len(sheetDataBytes))

	l.Info("creating PDF processor for turn sheet preview")

	// Create location choice processor for PDF generation
	processor, err := turn_sheet.NewLocationChoiceProcessor(l, cfg)
	if err != nil {
		l.Warn("failed to create location choice processor >%v<", err)
		return coreerror.NewInternalError("failed to create location choice processor: %v", err)
	}

	l.Info("generating PDF for turn sheet preview")

	// Generate PDF
	pdfData, err := processor.GenerateTurnSheet(ctx, l, turn_sheet.DocumentFormatPDF, sheetDataBytes)
	if err != nil {
		l.Warn("failed to generate location choice turn sheet PDF >%v<", err)
		return coreerror.NewInternalError("failed to generate turn sheet PDF: %v", err)
	}

	// Set Content-Disposition to inline for preview (not download)
	filename := fmt.Sprintf("preview-location-choice-%s.pdf", locationRec.Name)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))

	// Return PDF response
	l.Info("responding with location choice turn sheet preview PDF for location >%s< size >%d< bytes", locationID, len(pdfData))
	if err := server.WritePDFResponse(l, w, http.StatusOK, pdfData); err != nil {
		l.Warn("failed writing PDF response >%v<", err)
		return err
	}

	return nil
}

// getLocationOptionsForPreview retrieves location links as location options for the preview
func getLocationOptionsForPreview(ctx context.Context, mm *domain.Domain, gameID, locationID string) ([]turn_sheet.LocationOption, error) {
	// Get location links from this location
	opts := &coresql.Options{
		Params: []coresql.Param{
			{Col: "game_id", Val: gameID},
			{Col: "from_adventure_game_location_id", Val: locationID},
		},
	}

	linkRecs, err := mm.GetManyAdventureGameLocationLinkRecs(opts)
	if err != nil {
		return nil, err
	}

	// Convert links to location options
	var options []turn_sheet.LocationOption
	for _, link := range linkRecs {
		// Get the destination location to get its name
		destLocation, err := mm.GetAdventureGameLocationRec(link.ToAdventureGameLocationID, nil)
		if err != nil {
			continue // Skip if we can't get the destination
		}

		options = append(options, turn_sheet.LocationOption{
			LocationID:              link.ToAdventureGameLocationID,
			LocationLinkName:        destLocation.Name,
			LocationLinkDescription: destLocation.Description,
		})
	}

	return options, nil
}
