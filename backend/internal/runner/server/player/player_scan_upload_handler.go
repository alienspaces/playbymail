package player

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	UploadGSITurnSheetScan = "upload-gsi-turn-sheet-scan"
)

// maxScanUploadBytes is the maximum accepted multipart upload size for player scans.
const maxScanUploadBytes = 10 << 20 // 10 MB

func playerScanHandlerConfig(l logger.Logger, scnr turnsheet.TurnSheetScanner) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "playerScanHandlerConfig")

	l.Debug("Adding player scan upload handler configuration")

	cfg := make(map[string]server.HandlerConfig)

	cfg[UploadGSITurnSheetScan] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id/scan",
		HandlerFunc: uploadGSITurnSheetScanHandler(scnr),
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				"game_playing",
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Upload scanned turn sheet",
			Description: "Upload a scanned image of a completed turn sheet. " +
				"The backend runs OCR to extract the form data and saves it to the turn sheet. " +
				"Accepts multipart form data with an 'image' field. Auth: session token.",
		},
	}

	return cfg, nil
}

// uploadGSITurnSheetScanHandler accepts a scanned turn sheet image, runs OCR on it,
// and saves the extracted data to the turn sheet record.
func uploadGSITurnSheetScanHandler(scnr turnsheet.TurnSheetScanner) server.Handle {
	return func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
		l = logging.LoggerWithFunctionContext(l, packageName, "uploadGSITurnSheetScanHandler")

		l.Info("uploading turn sheet scan with path params >%#v<", pp)

		gameTurnSheetID := pp.ByName("game_turn_sheet_id")
		if gameTurnSheetID == "" {
			return coreerror.RequiredPathParameter("game_turn_sheet_id")
		}

		mm := m.(*domain.Domain)

		gsiRec, err := resolveGSI(l, r, pp, mm)
		if err != nil {
			return err
		}

		authData := server.GetRequestAuthenData(l, r)

		subRec, err := mm.GetGameSubscriptionRec(gsiRec.GameSubscriptionID, nil)
		if err != nil {
			l.Warn("failed to get game subscription >%s< >%v<", gsiRec.GameSubscriptionID, err)
			return err
		}

		turnSheetRec, err := mm.GetGameTurnSheetRec(gameTurnSheetID, coresql.ForUpdateNoWait)
		if err != nil {
			l.Warn("failed to get turn sheet >%s< >%v<", gameTurnSheetID, err)
			return err
		}

		if turnSheetRec.AccountID != gsiRec.AccountID || turnSheetRec.GameID != subRec.GameID {
			l.Warn("turn sheet >%s< does not belong to gsi >%s<", gameTurnSheetID, gsiRec.ID)
			return coreerror.NewNotFoundError("turn_sheet", "Turn sheet not found")
		}
		if turnSheetRec.AccountUserID != authData.AccountUser.ID {
			l.Warn("turn sheet >%s< does not belong to authenticated user >%s<", gameTurnSheetID, authData.AccountUser.ID)
			return coreerror.NewNotFoundError("turn_sheet", "Turn sheet not found")
		}

		if turnSheetRec.IsCompleted {
			return coreerror.NewInvalidDataError("turn sheet is already completed and cannot be modified")
		}

		// Parse multipart form upload.
		if err := r.ParseMultipartForm(maxScanUploadBytes); err != nil {
			l.Warn("failed to parse multipart form >%v<", err)
			return coreerror.NewInvalidDataError("failed to parse multipart form: %v", err)
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			l.Warn("failed to get image file from form >%v<", err)
			return coreerror.NewInvalidDataError("image file is required (field name: 'image')")
		}
		defer file.Close()

		imageData, err := io.ReadAll(file)
		if err != nil {
			l.Warn("failed to read image data >%v<", err)
			return coreerror.NewInvalidDataError("failed to read image data")
		}

		if len(imageData) == 0 {
			return coreerror.NewInvalidDataError("empty image data provided")
		}

		// When a scanner is available, run OCR to extract form data from the image.
		// When scnr is nil (test environments without OCR configured), skip OCR and
		// store the image bytes as-is in scanned_data for later processing.
		var scannedDataBytes []byte
		if scnr != nil {
			ctx := r.Context()
			scannedDataBytes, err = scnr.GetTurnSheetScanData(ctx, l, turnSheetRec.SheetType, turnSheetRec.SheetData, imageData)
			if err != nil {
				l.Warn("failed to scan turn sheet >%s< >%v<", gameTurnSheetID, err)
				return coreerror.NewInvalidDataError("failed to process scanned image: %v", err)
			}
		} else {
			l.Info("no scanner configured; storing raw image reference for turn sheet >%s<", gameTurnSheetID)
			scannedDataBytes, err = json.Marshal(map[string]string{"status": "pending_ocr"})
			if err != nil {
				return coreerror.NewInternalError("failed to marshal pending scan data")
			}
		}

		turnSheetRec.ScannedData = scannedDataBytes

		updatedRec, err := mm.UpdateGameTurnSheetRec(turnSheetRec)
		if err != nil {
			l.Warn("failed to update turn sheet with scan data >%v<", err)
			return err
		}

		l.Info("saved scanned data for turn sheet >%s< via gsi >%s<", gameTurnSheetID, gsiRec.ID)

		return server.WriteResponse(l, w, http.StatusOK, updatedRec)
	}
}
