package game

import (
	"bytes"
	"database/sql"
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
	UploadGameTurnSheetImage = "upload-game-turn-sheet-image"
	GetGameTurnSheetImage    = "get-game-turn-sheet-image"
	DeleteGameTurnSheetImage = "delete-game-turn-sheet-image"
	PreviewGameJoinTurnSheet = "preview-game-join-turn-sheet"
)

func gameImageHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameImageHandlerConfig")

	l.Debug("Adding game image handler configuration")

	gameImageConfig := make(map[string]server.HandlerConfig)

	// Upload turn sheet image endpoint
	gameImageConfig[UploadGameTurnSheetImage] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/turn-sheet-image",
		HandlerFunc: uploadGameTurnSheetImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Upload game turn sheet image",
			Description: "Upload a turn sheet background image for a game. " +
				"Accepts multipart form data with 'image' file. " +
				"Images must be WebP, PNG, or JPEG format, max 1MB. " +
				"Recommended: 2480x3508px (A4 @ 300 DPI) for best print quality. " +
				"Use the preview endpoint to see how images appear in the turn sheet.",
		},
	}

	// Get the turn sheet image for a game
	gameImageConfig[GetGameTurnSheetImage] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/turn-sheet-image",
		HandlerFunc: getGameTurnSheetImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Get game turn sheet image",
			Description: "Get the turn sheet background image for a game.",
		},
	}

	// Delete the turn sheet image
	gameImageConfig[DeleteGameTurnSheetImage] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/games/:game_id/turn-sheet-image",
		HandlerFunc: deleteGameTurnSheetImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Delete game turn sheet image",
			Description: "Delete the turn sheet background image for a game.",
		},
	}

	// Preview join game turn sheet with uploaded images
	gameImageConfig[PreviewGameJoinTurnSheet] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/turn-sheets/preview",
		HandlerFunc: previewGameJoinTurnSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Preview join game turn sheet",
			Description: "Generate and preview a join game turn sheet PDF with uploaded background images. " +
				"Returns PDF with Content-Disposition: inline for browser preview.",
		},
	}

	return gameImageConfig, nil
}

// uploadGameTurnSheetImageHandler handles uploading a turn sheet background image
func uploadGameTurnSheetImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "uploadGameTurnSheetImageHandler")

	l.Info("uploading game turn sheet image with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	l.Info("uploading turn sheet image for game >%s<", gameID)

	// Verify game exists and user has access (RLS check)
	_, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	// Parse multipart form (max 2MB to account for form overhead)
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		l.Warn("failed to parse multipart form >%v<", err)
		return coreerror.NewInvalidDataError("failed to parse multipart form: %v", err)
	}

	l.Info("processing image upload for game >%s<", gameID)

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

	// Create or update image record
	rec := &game_record.GameImage{
		GameID:    gameID,
		RecordID:  sql.NullString{Valid: false}, // NULL for game-level images
		Type:      game_record.GameImageTypeTurnSheetBackground,
		ImageData: imageData,
		MimeType:  mimeType,
		FileSize:  len(imageData),
		Width:     width,
		Height:    height,
	}

	l.Info("upserting game image record for game >%s<", gameID)

	rec, err = mm.UpsertGameImageRec(rec)
	if err != nil {
		l.Warn("failed to upsert game image record >%v<", err)
		return err
	}

	l.Info("successfully saved game image record id >%s<", rec.ID)

	res, err := mapper.GameImageRecordToResponse(l, rec, warning)
	if err != nil {
		l.Warn("failed to map game image record to response >%v<", err)
		return err
	}

	l.Info("responding with game image record id >%s< warning >%s<", rec.ID, warning)

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// getGameTurnSheetImageHandler gets the turn sheet background image for a game
func getGameTurnSheetImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameTurnSheetImageHandler")

	l.Info("getting game turn sheet image with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	l.Info("fetching turn sheet image for game >%s<", gameID)

	// Verify game exists and user has access (RLS check)
	_, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	img, err := mm.GetGameTurnSheetBackgroundImage(gameID)
	if err != nil {
		l.Warn("failed to get game turn sheet image >%v<", err)
		return err
	}

	res, err := mapper.GameTurnSheetImageToResponse(l, gameID, img)
	if err != nil {
		l.Warn("failed to map game turn sheet image to response >%v<", err)
		return err
	}

	l.Info("responding with game turn sheet image for game >%s<", gameID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// deleteGameTurnSheetImageHandler deletes the turn sheet background image for a game
func deleteGameTurnSheetImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteGameTurnSheetImageHandler")

	l.Info("deleting game turn sheet image with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	// Verify game exists and user has access (RLS check)
	_, err := mm.GetGameRec(gameID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	err = mm.DeleteGameImageByGameAndType(gameID, sql.NullString{Valid: false}, game_record.GameImageTypeTurnSheetBackground)
	if err != nil {
		l.Warn("failed to delete game image >%v<", err)
		return err
	}

	l.Info("deleted game turn sheet image for game >%s<", gameID)

	if err = server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// previewGameJoinTurnSheetHandler generates and previews a join game turn sheet PDF
func previewGameJoinTurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "previewGameJoinTurnSheetHandler")

	l.Info("previewing game join turn sheet with path params >%#v<", pp)

	ctx := r.Context()
	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	l.Info("generating turn sheet preview for game >%s<", gameID)

	// Get game record (RLS check)
	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	l.Info("found game >%s< type >%s<", gameRec.Name, gameRec.GameType)

	// Validate game type supports join sheets
	if gameRec.GameType != game_record.GameTypeAdventure {
		l.Warn("join turn sheet only supported for adventure games, got >%s< for game >%s<", gameRec.GameType, gameID)
		return coreerror.NewInvalidDataError("join turn sheet is not supported for game type %s", gameRec.GameType)
	}

	// Parse config for processor
	cfg, err := config.Parse()
	if err != nil {
		l.Warn("failed to parse config >%v<", err)
		return coreerror.NewInternalError("failed to parse config: %v", err)
	}

	// Generate join turn sheet code for this game
	turnSheetCode, err := turnsheet.GenerateJoinTurnSheetCode(gameID)
	if err != nil {
		l.Warn("failed to generate join turn sheet code >%v<", err)
		return coreerror.NewInternalError("failed to generate turn sheet code: %v", err)
	}

	l.Info("generated turn sheet code for game >%s<", gameID)

	// Create join game data using mapper function
	joinData, err := turn_sheet.GetTurnSheetJoinGameData(gameRec, turnSheetCode)
	if err != nil {
		l.Warn("failed to create join game data >%v<", err)
		return coreerror.NewInvalidDataError("failed to create join game data: %v", err)
	}

	// Get uploaded turn sheet background image and add it to the data
	backgroundImage, err := mm.GetGameTurnSheetImageDataURL(gameID)
	if err != nil {
		l.Warn("failed to get turn sheet background image >%v<", err)
		// Continue without image - not a fatal error
	} else if backgroundImage != "" {
		joinData.BackgroundImage = &backgroundImage
		l.Info("loaded background image for turn sheet, length >%d<", len(backgroundImage))
	} else {
		l.Info("no background image found for turn sheet")
	}

	// Marshal join data to JSON
	sheetDataBytes, err := json.Marshal(joinData)
	if err != nil {
		l.Warn("failed to marshal join game sheet data >%v<", err)
		return coreerror.NewInternalError("failed to marshal join game sheet data")
	}

	l.Info("marshaled join data size >%d< bytes", len(sheetDataBytes))

	l.Info("creating PDF processor for turn sheet preview")

	// Create join game processor for PDF generation
	processor, err := turn_sheet.NewJoinGameProcessor(l, cfg)
	if err != nil {
		l.Warn("failed to create join game processor >%v<", err)
		return coreerror.NewInternalError("failed to create join game processor: %v", err)
	}

	l.Info("generating PDF for turn sheet preview")

	// Generate PDF
	pdfData, err := processor.GenerateTurnSheet(ctx, l, turn_sheet.DocumentFormatPDF, sheetDataBytes)
	if err != nil {
		l.Warn("failed to generate join game turn sheet PDF >%v<", err)
		return coreerror.NewInternalError("failed to generate turn sheet PDF: %v", err)
	}

	// Set Content-Disposition to inline for preview (not download)
	filename := fmt.Sprintf("preview-join-game-turn-sheet-%s.pdf", gameRec.Name)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))

	// Return PDF response
	l.Info("responding with join game turn sheet preview PDF for game >%s< size >%d< bytes", gameID, len(pdfData))
	if err := server.WritePDFResponse(l, w, http.StatusOK, pdfData); err != nil {
		l.Warn("failed writing PDF response >%v<", err)
		return err
	}

	return nil
}
