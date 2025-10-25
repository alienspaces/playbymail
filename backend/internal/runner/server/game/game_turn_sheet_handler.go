package game

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet/base"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheet"
)

const (
	// Turn Sheet Scanning Endpoints
	UploadTurnSheet = "upload-turn-sheet"
)

// TurnSheetUploadResponse represents the response from uploading and processing a turn sheet
type TurnSheetUploadResponse struct {
	TurnSheetID      string                 `json:"turn_sheet_id"`
	TurnSheetCodeCI  string                 `json:"turn_sheet_code"`
	SheetType        string                 `json:"sheet_type"`
	ScannedData      map[string]interface{} `json:"scanned_data"`
	ProcessingStatus string                 `json:"processing_status"`
	ScanQuality      float64                `json:"scan_quality"`
	Message          string                 `json:"message"`
}

func gameTurnSheetHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameTurnSheetHandlerConfig")

	l.Debug("Adding game turn sheet handler configuration")

	gameTurnSheetConfig := make(map[string]server.HandlerConfig)

	// Upload turn sheet endpoint - Single pass: Upload, scan, extract code, process, save
	gameTurnSheetConfig[UploadTurnSheet] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/turn-sheets/upload",
		HandlerFunc: uploadTurnSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Upload and process turn sheet",
			Description: "Upload a scanned turn sheet image and process it in a single pass. " +
				"This extracts the turn sheet code, retrieves the turn sheet record, " +
				"processes the scanned data, and saves the results.",
		},
	}

	return gameTurnSheetConfig, nil
}

// uploadTurnSheetHandler handles the single-pass turn sheet upload and processing
func uploadTurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "uploadTurnSheetHandler")

	l.Info("uploading and processing turn sheet with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	if instanceID == "" {
		l.Warn("instance ID is empty")
		return coreerror.RequiredPathParameter("instance_id")
	}

	// Read image data from request body
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		l.Warn("failed to read image data >%v<", err)
		return coreerror.NewInvalidDataError("failed to read image data")
	}

	if len(imageData) == 0 {
		l.Warn("empty image data provided")
		return coreerror.NewInvalidDataError("empty image data provided")
	}

	baseProcessor := base.NewBaseProcessor(l)

	// Step 1: Extract turn sheet code from image
	turnSheetCode, err := baseProcessor.ParseTurnSheetCodeFromImage(r.Context(), imageData)
	if err != nil {
		l.Warn("failed to extract turn sheet code >%v<", err)
		return coreerror.NewInvalidDataError("failed to extract turn sheet code from image")
	}

	// Step 2: Parse the turn sheet code to get identifiers
	identifier, err := turnsheet.ParseTurnSheetCode(turnSheetCode)
	if err != nil {
		l.Warn("failed to parse turn sheet code >%v<", err)
		return coreerror.NewInvalidDataError("invalid turn sheet code format")
	}

	// Step 3: Verify the turn sheet belongs to the specified game instance
	if identifier.GameInstanceID != instanceID {
		l.Warn("turn sheet does not belong to specified game instance")
		return coreerror.NewInvalidDataError("turn sheet does not belong to specified game instance")
	}

	// Step 4: Get the turn sheet record
	turnSheetRec, err := mm.GetGameTurnSheetRec(identifier.GameTurnSheetID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get turn sheet record >%v<", err)
		return coreerror.NewNotFoundError("turn sheet", identifier.GameTurnSheetID)
	}

	// Step 5: Get the appropriate document processor for this turn sheet type
	processor, err := turn_sheet.GetDocumentProcessor(l, turnSheetRec.SheetType)
	if err != nil {
		l.Warn("failed to get processor for turn sheet type >%s< >%v<", turnSheetRec.SheetType, err)
		return coreerror.NewInvalidDataError("unsupported turn sheet type: %s", turnSheetRec.SheetType)
	}

	// Step 6: Parse the sheet data to get the actual data for processing
	var sheetData interface{}
	if err := json.Unmarshal(turnSheetRec.SheetData, &sheetData); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return coreerror.NewInvalidDataError("invalid sheet data format")
	}

	// Step 7: Process the turn sheet using the appropriate processor
	scannedData, err := processor.ScanTurnSheet(r.Context(), l, imageData, sheetData)
	if err != nil {
		l.Warn("failed to process turn sheet >%v<", err)
		return coreerror.NewInvalidDataError("failed to process turn sheet: %v", err)
	}

	// Step 8: Update the turn sheet record with scanned data
	scannedDataBytes, err := json.Marshal(scannedData)
	if err != nil {
		l.Warn("failed to marshal scanned data >%v<", err)
		return coreerror.NewInternalError("failed to marshal scanned data")
	}

	turnSheetRec.ScannedData = json.RawMessage(scannedDataBytes)
	turnSheetRec.ScannedAt = sql.NullTime{Time: time.Now(), Valid: true}
	turnSheetRec.ScanQuality = sql.NullFloat64{Float64: 0.95, Valid: true} // TODO: Calculate actual scan quality
	turnSheetRec.ProcessingStatus = "processed"

	// Step 9: Update the record in the database
	_, err = mm.UpdateGameTurnSheetRec(turnSheetRec)
	if err != nil {
		l.Warn("failed to update turn sheet record >%v<", err)
		return coreerror.NewInternalError("failed to update turn sheet record")
	}

	// Step 10: Create response
	response := TurnSheetUploadResponse{
		TurnSheetID:      identifier.GameTurnSheetID,
		TurnSheetCodeCI:  turnSheetCode,
		SheetType:        turnSheetRec.SheetType,
		ScannedData:      scannedData.(map[string]interface{}),
		ProcessingStatus: "processed",
		ScanQuality:      0.95,
		Message:          "Turn sheet uploaded and processed successfully",
	}

	l.Info("responding with turn sheet upload result >%+v<", response)

	err = server.WriteResponse(l, w, http.StatusOK, response)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
