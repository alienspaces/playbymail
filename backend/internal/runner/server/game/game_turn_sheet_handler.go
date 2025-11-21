package game

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/convert"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheet"
)

const (
	// Turn Sheet Scanning Endpoints
	UploadTurnSheet = "upload-turn-sheet"
)

// TurnSheetUploadResponse represents the response from uploading and processing a turn sheet
type TurnSheetUploadResponse struct {
	TurnSheetID      string         `json:"turn_sheet_id"`
	TurnSheetCodeCI  string         `json:"turn_sheet_code"`
	SheetType        string         `json:"sheet_type"`
	ScannedData      map[string]any `json:"scanned_data"`
	ProcessingStatus string         `json:"processing_status"`
	ScanQuality      float64        `json:"scan_quality"`
	Message          string         `json:"message"`
}

func gameTurnSheetHandlerConfig(l logger.Logger, scanner TurnSheetScanner) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameTurnSheetHandlerConfig")

	l.Debug("Adding game turn sheet handler configuration")

	gameTurnSheetConfig := make(map[string]server.HandlerConfig)

	// Upload turn sheet endpoint - Handles both existing player and join-game submissions
	gameTurnSheetConfig[UploadTurnSheet] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/turn-sheets",
		HandlerFunc: uploadTurnSheetHandler(scanner),
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
func uploadTurnSheetHandler(scanner TurnSheetScanner) server.Handle {
	return func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
		l = logging.LoggerWithFunctionContext(l, packageName, "uploadTurnSheetHandler")

		l.Info("uploading and processing turn sheet with path params >%#v<", pp)

		ctx := r.Context()
		mm := m.(*domain.Domain)

		imageData, err := io.ReadAll(r.Body)
		if err != nil {
			l.Warn("failed to read image data >%v<", err)
			return coreerror.NewInvalidDataError("failed to read image data")
		}

		if len(imageData) == 0 {
			l.Warn("empty image data provided")
			return coreerror.NewInvalidDataError("empty image data provided")
		}

		turnSheetCode, err := scanner.GetTurnSheetCodeFromImage(ctx, l, imageData)
		if err != nil {
			l.Warn("failed to extract turn sheet code >%v<", err)
			return coreerror.NewInvalidDataError("failed to extract turn sheet code from image")
		}

		identifier, err := turnsheet.ParseTurnSheetCode(turnSheetCode)
		if err != nil {
			l.Warn("failed to parse turn sheet code >%v<", err)
			return coreerror.NewInvalidDataError("invalid turn sheet code format")
		}

		var (
			resp   *TurnSheetUploadResponse
			status int
		)

		switch identifier.CodeType {
		case turnsheet.TurnSheetCodeTypeJoiningGame:
			resp, status, err = handleJoinTurnSheetUpload(ctx, l, scanner, mm, jc, turnSheetCode, identifier, imageData)
		case turnsheet.TurnSheetCodeTypePlayingGame, "":
			resp, status, err = handleStandardTurnSheetUpload(ctx, l, scanner, mm, turnSheetCode, identifier, imageData)
		default:
			err = coreerror.NewInvalidDataError("unsupported turn sheet code type: %s", identifier.CodeType)
		}
		if err != nil {
			return err
		}

		l.Info("responding with turn sheet upload result >%+v<", resp)

		if err := server.WriteResponse(l, w, status, resp); err != nil {
			l.Warn("failed writing response >%v<", err)
			return err
		}

		return nil
	}
}

// handleJoinTurnSheetUpload processes an uploaded join game turn sheet.
//
// It fetches the game record to validate the game type, constructs default turn sheet metadata,
// and invokes the sheet scanner to extract player registration data from the submission image.
// Using the extracted registration information, it attempts to associate the joining player with
// an account, creates a new pending account if necessary, and upserts a game subscription record.
// Finally, it creates a new turn sheet record for the join game turn sheet and returns the upload
// status and processed data for further handling.
func handleJoinTurnSheetUpload(ctx context.Context, l logger.Logger, scanner TurnSheetScanner, m *domain.Domain, jc *river.Client[pgx.Tx], turnSheetCode string, identifier *turnsheet.TurnSheetIdentifier, imageData []byte) (*TurnSheetUploadResponse, int, error) {
	l = l.WithFunctionContext("handleJoinTurnSheetUpload")

	gameRec, err := m.GetGameRec(identifier.GameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", identifier.GameID, err)
		return nil, 0, coreerror.NewNotFoundError("game", identifier.GameID)
	}

	if gameRec.GameType != game_record.GameTypeAdventure {
		l.Warn("join turn sheet only supported for adventure games, got >%s<", gameRec.GameType)
		return nil, 0, coreerror.NewInvalidDataError("join turn sheet is not supported for game type %s", gameRec.GameType)
	}

	joinData := turn_sheet.JoinGameData{
		TurnSheetTemplateData: turn_sheet.TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr(gameRec.GameType),
			TurnNumber:            convert.Ptr(0),
			AccountName:           nil,
			TurnSheetTitle:        convert.Ptr("Join Game"),
			TurnSheetDescription:  convert.Ptr(fmt.Sprintf("Welcome to %s! Welcome to the PlayByMail Adventure!", gameRec.Name)),
			TurnSheetInstructions: convert.Ptr(turn_sheet.DefaultJoinGameInstructions()),
			TurnSheetDeadline:     nil,
			TurnSheetCode:         convert.Ptr(turnSheetCode),
		},
		GameDescription: fmt.Sprintf("Welcome to %s! Welcome to the PlayByMail Adventure!", gameRec.Name),
	}

	sheetDataBytes, err := json.Marshal(joinData)
	if err != nil {
		l.Warn("failed to marshal join game sheet data >%v<", err)
		return nil, 0, coreerror.NewInternalError("failed to marshal join game sheet data")
	}

	scannedData, err := scanner.GetTurnSheetScanData(ctx, l, adventure_game_record.AdventureSheetTypeJoinGame, sheetDataBytes, imageData)
	if err != nil {
		l.Warn("failed to scan join game turn sheet >%v<", err)
		return nil, 0, coreerror.NewInvalidDataError("failed to process join game turn sheet: %v", err)
	}

	var scanData turn_sheet.JoinGameScanData
	if err := json.Unmarshal(scannedData, &scanData); err != nil {
		l.Warn("failed to unmarshal join game scan data >%v<", err)
		return nil, 0, coreerror.NewInvalidDataError("invalid join game turn sheet data")
	}

	accountRec, err := m.GetAccountRecByEmail(scanData.Email)
	if err != nil {
		l.Warn("failed to get account by email >%s< >%v<", scanData.Email, err)
		return nil, 0, err
	}

	if accountRec == nil {
		l.Info("creating new pending account for email >%s<", scanData.Email)
		accountRec = &account_record.Account{
			Email:  scanData.Email,
			Status: account_record.AccountStatusPendingApproval,
		}

		accountRec, err = m.CreateAccountRec(accountRec)
		if err != nil {
			l.Warn("failed to create account >%v<", err)
			return nil, 0, err
		}
	}

	// Create or get account contact
	accountContactRec := &account_record.AccountContact{
		AccountID:          accountRec.ID,
		Name:               scanData.Name,
		PostalAddressLine1: scanData.PostalAddressLine1,
		PostalAddressLine2: nullstring.FromString(scanData.PostalAddressLine2),
		StateProvince:      scanData.StateProvince,
		Country:            scanData.Country,
		PostalCode:         scanData.PostalCode,
	}

	accountContactRec, err = m.CreateAccountContactRec(accountContactRec)
	if err != nil {
		l.Warn("failed to create account contact >%v<", err)
		return nil, 0, err
	}

	subscriptionRec, err := m.UpsertPendingGameSubscription(gameRec.ID, accountRec.ID, accountContactRec.ID, game_record.GameSubscriptionTypePlayer)
	if err != nil {
		l.Warn("failed to upsert game subscription >%v<", err)
		return nil, 0, err
	}

	turnSheetRec := &game_record.GameTurnSheet{
		GameID:           gameRec.ID,
		AccountID:        accountRec.ID,
		TurnNumber:       0,
		SheetType:        adventure_game_record.AdventureSheetTypeJoinGame,
		SheetOrder:       1,
		SheetData:        json.RawMessage(sheetDataBytes),
		ScannedData:      json.RawMessage(scannedData),
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
	turnSheetRec.ScannedAt = sql.NullTime{Time: time.Now(), Valid: true}

	turnSheetRec, err = m.CreateGameTurnSheetRec(turnSheetRec)
	if err != nil {
		l.Warn("failed to create game turn sheet record >%v<", err)
		return nil, 0, err
	}

	if _, err := jc.InsertTx(context.Background(), m.Tx, &jobworker.SendGameSubscriptionApprovalEmailWorkerArgs{
		GameSubscriptionID: subscriptionRec.ID,
	}, &river.InsertOpts{Queue: jobqueue.QueueDefault}); err != nil {
		l.Warn("failed to enqueue game subscription approval email job >%v<", err)
	}

	var scannedDataMap map[string]any
	if err := json.Unmarshal(scannedData, &scannedDataMap); err != nil {
		l.Warn("failed to unmarshal join game scanned data >%v<", err)
		scannedDataMap = make(map[string]any)
	}

	response := &TurnSheetUploadResponse{
		TurnSheetID:      turnSheetRec.ID,
		TurnSheetCodeCI:  turnSheetCode,
		SheetType:        turnSheetRec.SheetType,
		ScannedData:      scannedDataMap,
		ProcessingStatus: turnSheetRec.ProcessingStatus,
		ScanQuality:      0,
		Message:          "Join game turn sheet received. Please confirm the subscription email to continue.",
	}

	return response, http.StatusAccepted, nil
}

func handleStandardTurnSheetUpload(ctx context.Context, l logger.Logger, scanner TurnSheetScanner, m *domain.Domain, turnSheetCode string, identifier *turnsheet.TurnSheetIdentifier, imageData []byte) (*TurnSheetUploadResponse, int, error) {
	l = l.WithFunctionContext("handleStandardTurnSheetUpload")

	turnSheetRec, err := m.GetGameTurnSheetRec(identifier.GameTurnSheetID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get turn sheet record >%v<", err)
		return nil, 0, coreerror.NewNotFoundError("turn sheet", identifier.GameTurnSheetID)
	}

	if !nullstring.IsValid(turnSheetRec.GameInstanceID) || identifier.GameInstanceID == "" ||
		nullstring.ToString(turnSheetRec.GameInstanceID) != identifier.GameInstanceID {
		l.Warn("turn sheet does not belong to expected game instance >%s<", identifier.GameInstanceID)
		return nil, 0, coreerror.NewInvalidDataError("turn sheet does not belong to expected game instance")
	}

	scannedData, err := scanner.GetTurnSheetScanData(ctx, l, turnSheetRec.SheetType, turnSheetRec.SheetData, imageData)
	if err != nil {
		l.Warn("failed to process turn sheet >%v<", err)
		return nil, 0, coreerror.NewInvalidDataError("failed to process turn sheet: %v", err)
	}

	turnSheetRec.ScannedData = json.RawMessage(scannedData)
	turnSheetRec.ScannedAt = sql.NullTime{Time: time.Now(), Valid: true}
	turnSheetRec.ScanQuality = sql.NullFloat64{Float64: 0.95, Valid: true} // TODO: Calculate actual scan quality
	turnSheetRec.ProcessingStatus = game_record.TurnSheetProcessingStatusProcessed

	if _, err := m.UpdateGameTurnSheetRec(turnSheetRec); err != nil {
		l.Warn("failed to update turn sheet record >%v<", err)
		return nil, 0, coreerror.NewInternalError("failed to update turn sheet record >%v<", err)
	}

	var scannedDataMap map[string]any
	if err := json.Unmarshal(scannedData, &scannedDataMap); err != nil {
		l.Warn("failed to unmarshal scanned data for response >%v<", err)
		return nil, 0, coreerror.NewInternalError("failed to unmarshal scanned data for response >%v<", err)
	}

	response := &TurnSheetUploadResponse{
		TurnSheetID:      identifier.GameTurnSheetID,
		TurnSheetCodeCI:  turnSheetCode,
		SheetType:        turnSheetRec.SheetType,
		ScannedData:      scannedDataMap,
		ProcessingStatus: turnSheetRec.ProcessingStatus,
		ScanQuality:      0.95,
		Message:          "Turn sheet uploaded and processed successfully",
	}

	return response, http.StatusOK, nil
}
