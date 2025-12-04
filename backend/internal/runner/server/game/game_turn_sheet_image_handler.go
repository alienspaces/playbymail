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
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheet"
)

const (
	UploadGameTurnSheetImage    = "upload-game-turn-sheet-image"
	GetManyGameTurnSheetImages  = "get-many-game-turn-sheet-images"
	GetOneGameTurnSheetImage    = "get-one-game-turn-sheet-image"
	DeleteOneGameTurnSheetImage = "delete-one-game-turn-sheet-image"
	PreviewGameTurnSheet        = "preview-game-turn-sheet"
)

func gameImageHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameImageHandlerConfig")

	l.Debug("Adding game image handler configuration")

	gameImageConfig := make(map[string]server.HandlerConfig)

	// Upload turn sheet image endpoint
	gameImageConfig[UploadGameTurnSheetImage] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/turn-sheet-images",
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
				"Query parameter 'turn_sheet_type' specifies the turn sheet type (e.g., 'adventure_game_join_game', 'adventure_game_inventory_management'). " +
				"Defaults to 'adventure_game_join_game' if not provided for backward compatibility. " +
				"Images must be WebP, PNG, or JPEG format, max 1MB. " +
				"Recommended: 2480x3508px (A4 @ 300 DPI) for best print quality. " +
				"Use the preview endpoint to see how images appear in the turn sheet.",
		},
	}

	// Get many turn sheet images for a game
	gameImageConfig[GetManyGameTurnSheetImages] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/turn-sheet-images",
		HandlerFunc: getManyGameTurnSheetImagesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get many game turn sheet images",
			Description: "Get all turn sheet background images for a game. " +
				"Optional query parameter 'turn_sheet_type' filters by turn sheet type (e.g., 'adventure_game_join_game', 'adventure_game_inventory_management'). " +
				"If not provided, returns all turn sheet images for the game.",
		},
	}

	// Get one turn sheet image by ID
	gameImageConfig[GetOneGameTurnSheetImage] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/turn-sheet-images/:game_image_id",
		HandlerFunc: getOneGameTurnSheetImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Get one game turn sheet image",
			Description: "Get a specific turn sheet background image by ID for a game.",
		},
	}

	// Delete one turn sheet image by ID
	gameImageConfig[DeleteOneGameTurnSheetImage] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/games/:game_id/turn-sheet-images/:game_image_id",
		HandlerFunc: deleteOneGameTurnSheetImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Delete one game turn sheet image",
			Description: "Delete a specific turn sheet background image by ID for a game.",
		},
	}

	// Preview game turn sheet with uploaded images
	gameImageConfig[PreviewGameTurnSheet] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/turn-sheets/preview",
		HandlerFunc: previewGameTurnSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Preview game turn sheet",
			Description: "Generate and preview a game turn sheet PDF with uploaded background images. " +
				"Query parameter 'turn_sheet_type' is required and specifies the turn sheet type (e.g., 'adventure_game_join_game', 'adventure_game_inventory_management'). " +
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

	// Get turn_sheet_type from query parameter (default to join_game for backward compatibility)
	turnSheetType := r.URL.Query().Get("turn_sheet_type")
	if turnSheetType == "" {
		turnSheetType = adventure_game_record.AdventureGameTurnSheetTypeJoinGame
		l.Info("no turn_sheet_type provided, defaulting to >%s<", turnSheetType)
	}

	l.Info("using turn_sheet_type >%s<", turnSheetType)

	// Create or update image record
	rec := &game_record.GameImage{
		GameID:        gameID,
		RecordID:      sql.NullString{Valid: false}, // NULL for game-level images
		Type:          game_record.GameImageTypeTurnSheetBackground,
		TurnSheetType: turnSheetType,
		ImageData:     imageData,
		MimeType:      mimeType,
		FileSize:      len(imageData),
		Width:         width,
		Height:        height,
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

// getManyGameTurnSheetImagesHandler gets all turn sheet background images for a game
func getManyGameTurnSheetImagesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameTurnSheetImagesHandler")

	l.Info("getting many game turn sheet images with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	l.Info("fetching turn sheet images for game >%s<", gameID)

	// Verify game exists and user has access (RLS check)
	_, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	// Build query options
	opts := &coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameImageGameID, Val: gameID},
			{Col: game_record.FieldGameImageType, Val: game_record.GameImageTypeTurnSheetBackground},
		},
	}

	// Optional filter by turn_sheet_type
	turnSheetType := r.URL.Query().Get("turn_sheet_type")
	if turnSheetType != "" {
		l.Info("filtering by turn_sheet_type >%s<", turnSheetType)
		opts.Params = append(opts.Params, coresql.Param{
			Col: game_record.FieldGameImageTurnSheetType,
			Val: turnSheetType,
		})
	}

	images, err := mm.GetManyGameImageRecs(opts)
	if err != nil {
		l.Warn("failed to get game turn sheet images >%v<", err)
		return err
	}

	res, err := mapper.GameImageRecordsToCollectionResponse(l, images)
	if err != nil {
		l.Warn("failed to map game turn sheet images to response >%v<", err)
		return err
	}

	l.Info("responding with >%d< game turn sheet images for game >%s<", len(images), gameID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// getOneGameTurnSheetImageHandler gets a specific turn sheet background image by ID
func getOneGameTurnSheetImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneGameTurnSheetImageHandler")

	l.Info("getting one game turn sheet image with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	gameImageID := pp.ByName("game_image_id")
	if gameImageID == "" {
		l.Warn("game image ID is empty")
		return coreerror.RequiredPathParameter("game_image_id")
	}

	l.Info("fetching turn sheet image >%s< for game >%s<", gameImageID, gameID)

	// Verify game exists and user has access (RLS check)
	_, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	// Get the image record
	img, err := mm.GetGameImageRec(gameImageID, nil)
	if err != nil {
		l.Warn("failed to get game image record >%s< >%v<", gameImageID, err)
		return err
	}

	// Verify the image belongs to the game
	if img.GameID != gameID {
		l.Warn("game image >%s< does not belong to game >%s<", gameImageID, gameID)
		return coreerror.NewNotFoundError("game_image", gameImageID)
	}

	// Verify it's a turn sheet background image
	if img.Type != game_record.GameImageTypeTurnSheetBackground {
		l.Warn("game image >%s< is not a turn sheet background image", gameImageID)
		return coreerror.NewNotFoundError("game_image", gameImageID)
	}

	res, err := mapper.GameImageRecordToResponse(l, img, "")
	if err != nil {
		l.Warn("failed to map game turn sheet image to response >%v<", err)
		return err
	}

	l.Info("responding with game turn sheet image >%s< for game >%s<", gameImageID, gameID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// deleteOneGameTurnSheetImageHandler deletes a specific turn sheet background image by ID
func deleteOneGameTurnSheetImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneGameTurnSheetImageHandler")

	l.Info("deleting one game turn sheet image with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	gameImageID := pp.ByName("game_image_id")
	if gameImageID == "" {
		l.Warn("game image ID is empty")
		return coreerror.RequiredPathParameter("game_image_id")
	}

	l.Info("deleting turn sheet image >%s< for game >%s<", gameImageID, gameID)

	// Verify game exists and user has access (RLS check)
	_, err := mm.GetGameRec(gameID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return err
	}

	// Get the image record to verify it belongs to the game
	img, err := mm.GetGameImageRec(gameImageID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get game image record >%s< >%v<", gameImageID, err)
		return err
	}

	// Verify the image belongs to the game
	if img.GameID != gameID {
		l.Warn("game image >%s< does not belong to game >%s<", gameImageID, gameID)
		return coreerror.NewNotFoundError("game_image", gameImageID)
	}

	// Verify it's a turn sheet background image
	if img.Type != game_record.GameImageTypeTurnSheetBackground {
		l.Warn("game image >%s< is not a turn sheet background image", gameImageID)
		return coreerror.NewNotFoundError("game_image", gameImageID)
	}

	// Delete the image
	err = mm.DeleteGameImageRec(gameImageID)
	if err != nil {
		l.Warn("failed to delete game image >%v<", err)
		return err
	}

	l.Info("successfully deleted game turn sheet image >%s< for game >%s<", gameImageID, gameID)

	if err = server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// previewGameTurnSheetHandler generates and previews a game turn sheet PDF
// Supports all game-level turn sheet types (e.g., join_game, inventory_management)
func previewGameTurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "previewGameTurnSheetHandler")

	l.Info("previewing game turn sheet with path params >%#v<", pp)

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

	// Get turn_sheet_type from query parameter (required)
	turnSheetType := r.URL.Query().Get("turn_sheet_type")
	if turnSheetType == "" {
		l.Warn("turn_sheet_type query parameter is required for turn sheet preview")
		return coreerror.RequiredQueryParameter("turn_sheet_type")
	}
	l.Info("using turn_sheet_type >%s< for preview", turnSheetType)

	// Parse config for processor
	cfg, err := config.Parse()
	if err != nil {
		l.Warn("failed to parse config >%v<", err)
		return coreerror.NewInternalError("failed to parse config: %v", err)
	}

	// Get the appropriate processor for the turn sheet type
	processor, err := turn_sheet.GetDocumentProcessor(l, cfg, turnSheetType)
	if err != nil {
		l.Warn("failed to get document processor for turn sheet type >%s< >%v<", turnSheetType, err)
		return coreerror.NewInvalidDataError("turn sheet type %s is not supported for preview", turnSheetType)
	}

	// Generate turn sheet data based on type
	var sheetDataBytes []byte
	switch turnSheetType {
	case adventure_game_record.AdventureGameTurnSheetTypeJoinGame:
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
		backgroundImage, err := mm.GetGameTurnSheetImageDataURL(gameID, turnSheetType)
		if err != nil {
			l.Warn("failed to get turn sheet background image >%v<", err)
			// Continue without image - not a fatal error
		} else if backgroundImage != "" {
			joinData.BackgroundImage = &backgroundImage
			prefixLen := 50
			if len(backgroundImage) < prefixLen {
				prefixLen = len(backgroundImage)
			}
			l.Info("loaded background image for turn sheet, length >%d< prefix >%s<", len(backgroundImage), backgroundImage[:prefixLen])
		} else {
			l.Info("no background image found for turn sheet")
		}

		// Debug: Verify BackgroundImage is set before marshaling
		if joinData.BackgroundImage != nil {
			l.Info("background image set in joinData before marshaling, length >%d<", len(*joinData.BackgroundImage))
		} else {
			l.Info("background image NOT set in joinData before marshaling")
		}

		// Marshal join data to JSON
		sheetDataBytes, err = json.Marshal(joinData)
		if err != nil {
			l.Warn("failed to marshal join game sheet data >%v<", err)
			return coreerror.NewInternalError("failed to marshal join game sheet data")
		}

	case adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement:
		// Generate a preview turn sheet code
		previewTurnSheetCode, err := turnsheet.GeneratePreviewTurnSheetCode(gameID)
		if err != nil {
			l.Warn("failed to generate preview turn sheet code >%v<", err)
			return coreerror.NewInternalError("failed to generate preview turn sheet code: %v", err)
		}

		// Create sample inventory management data for preview
		inventoryData := turn_sheet.InventoryManagementData{
			TurnSheetTemplateData: turn_sheet.TurnSheetTemplateData{
				GameName: &gameRec.Name,
				TurnNumber: func() *int {
					n := 1
					return &n
				}(),
				TurnSheetCode: &previewTurnSheetCode,
			},
			CharacterName:       "Preview Character",
			CurrentLocationName: "Preview Location",
			InventoryCapacity:   10,
			InventoryCount:      3,
			CurrentInventory: []turn_sheet.InventoryItem{
				{
					ItemInstanceID:  "preview-item-1",
					ItemName:        "Sample Sword",
					ItemDescription: "A sturdy iron sword",
					IsEquipped:      true,
					EquipmentSlot:   "weapon",
					CanEquip:        true,
				},
				{
					ItemInstanceID:  "preview-item-2",
					ItemName:        "Leather Armor",
					ItemDescription: "Basic leather protection",
					IsEquipped:      false,
					CanEquip:        true,
				},
				{
					ItemInstanceID:  "preview-item-3",
					ItemName:        "Health Potion",
					ItemDescription: "Restores health when consumed",
					IsEquipped:      false,
					CanEquip:        false,
				},
			},
			EquipmentSlots: turn_sheet.EquipmentSlots{
				Weapon: &turn_sheet.EquippedItem{
					ItemInstanceID: "preview-item-1",
					ItemName:       "Sample Sword",
					SlotName:       "Weapon",
				},
			},
			LocationItems: []turn_sheet.LocationItem{
				{
					ItemInstanceID:  "location-item-1",
					ItemName:        "Gold Coin",
					ItemDescription: "A shiny gold coin",
				},
			},
		}

		// Get uploaded turn sheet background image and add it to the data
		backgroundImage, err := mm.GetGameTurnSheetImageDataURL(gameID, turnSheetType)
		if err != nil {
			l.Warn("failed to get turn sheet background image >%v<", err)
			// Continue without image - not a fatal error
		} else if backgroundImage != "" {
			inventoryData.BackgroundImage = &backgroundImage
			prefixLen := 50
			if len(backgroundImage) < prefixLen {
				prefixLen = len(backgroundImage)
			}
			l.Info("loaded background image for inventory turn sheet, length >%d< prefix >%s<", len(backgroundImage), backgroundImage[:prefixLen])
		} else {
			l.Info("no background image found for inventory turn sheet")
		}

		// Marshal inventory data to JSON
		sheetDataBytes, err = json.Marshal(inventoryData)
		if err != nil {
			l.Warn("failed to marshal inventory management sheet data >%v<", err)
			return coreerror.NewInternalError("failed to marshal inventory management sheet data")
		}

	default:
		l.Warn("turn sheet type >%s< is not supported for preview (requires character instance data)", turnSheetType)
		return coreerror.NewInvalidDataError("turn sheet type %s is not supported for preview", turnSheetType)
	}

	l.Info("marshaled turn sheet data size >%d< bytes", len(sheetDataBytes))

	l.Info("generating PDF for turn sheet preview")

	// Generate PDF using the processor
	pdfData, err := processor.GenerateTurnSheet(ctx, l, turn_sheet.DocumentFormatPDF, sheetDataBytes)
	if err != nil {
		l.Warn("failed to generate turn sheet PDF >%v<", err)
		return coreerror.NewInternalError("failed to generate turn sheet PDF: %v", err)
	}

	// Set Content-Disposition to inline for preview (not download)
	filename := fmt.Sprintf("preview-turn-sheet-%s-%s.pdf", turnSheetType, gameRec.Name)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))

	// Return PDF response
	l.Info("responding with turn sheet preview PDF for game >%s< type >%s< size >%d< bytes", gameID, turnSheetType, len(pdfData))
	if err := server.WritePDFResponse(l, w, http.StatusOK, pdfData); err != nil {
		l.Warn("failed writing PDF response >%v<", err)
		return err
	}

	return nil
}
