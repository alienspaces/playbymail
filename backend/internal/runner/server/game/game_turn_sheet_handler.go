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
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

const (
	// Turn Sheet Scanning Endpoints
	UploadTurnSheet = "upload-turn-sheet"
	// Turn Sheet Download Endpoints
	DownloadJoinGameTurnSheets = "download-join-game-turn-sheets"
)

// TurnSheetUploadResponse represents the response from uploading and processing a turn sheet
type TurnSheetUploadResponse struct {
	TurnSheetID      string         `json:"turn_sheet_id"`
	TurnSheetCode    string         `json:"turn_sheet_code"`
	SheetType        string         `json:"sheet_type"`
	ScannedData      map[string]any `json:"scanned_data"`
	ProcessingStatus string         `json:"processing_status"`
}

func gameTurnSheetHandlerConfig(l logger.Logger, scanner turnsheet.TurnSheetScanner) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameTurnSheetHandlerConfig")

	l.Debug("Adding game turn sheet handler configuration")

	gameTurnSheetConfig := make(map[string]server.HandlerConfig)

	// Upload turn sheet endpoint - Handles both existing player and join-game submissions
	// Handles both existing player and join-game submissions
	gameTurnSheetConfig[UploadTurnSheet] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/turn-sheets",
		HandlerFunc: uploadTurnSheetHandler(scanner),
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			// Permissions checked in handler (Player or Manager based on turn sheet type)
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Upload and process turn sheet",
			Description: "Upload a scanned turn sheet image and process it in a single pass. " +
				"This extracts the turn sheet code, retrieves the turn sheet record, " +
				"processes the scanned data, and saves the results.",
		},
	}

	// Download join game turn sheets endpoint
	// Supports optional token auth: managers can access without game_subscription_id query param,
	// while unauthenticated requests require game_subscription_id to be provided.
	gameTurnSheetConfig[DownloadJoinGameTurnSheets] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/turn-sheets",
		HandlerFunc: downloadJoinGameTurnSheetsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeOptionalToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameDesign,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Download join game turn sheet",
			Description: "Generate and download a join game turn sheet PDF. " +
				"The same PDF can be printed multiple times for distribution. " +
				"All join sheets for a game use the same turn sheet code.",
		},
	}

	return gameTurnSheetConfig, nil
}

// uploadTurnSheetHandler handles the single-pass turn sheet upload and processing
func uploadTurnSheetHandler(scanner turnsheet.TurnSheetScanner) server.Handle {
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

		// Get the turn sheet code from the uploaded image
		turnSheetCode, err := scanner.GetTurnSheetCodeFromImage(ctx, l, imageData)
		if err != nil {
			l.Warn("failed to extract turn sheet code >%v<", err)
			return coreerror.NewInvalidDataError("failed to extract turn sheet code from image")
		}

		// Get the type of turn sheet that was uploaded
		turnSheetCodeType, err := turnsheetutil.ParseTurnSheetCodeTypeFromCode(turnSheetCode)
		if err != nil {
			l.Warn("failed to parse turn sheet code >%v<", err)
			return coreerror.NewInvalidDataError("invalid turn sheet code format")
		}

		var (
			resp   *TurnSheetUploadResponse
			status int
		)

		switch turnSheetCodeType {
		case turnsheetutil.TurnSheetCodeTypeJoiningGame:
			var joinSheetData *turnsheetutil.JoinGameTurnSheetCodeData
			joinSheetData, err = turnsheetutil.ParseJoinGameTurnSheetCodeData(turnSheetCode)
			if err != nil {
				l.Warn("failed to parse join game turn sheet code >%v<", err)
				return coreerror.NewInvalidDataError("invalid join game turn sheet code format")
			}
			resp, status, err = handleJoinGameTurnSheetUpload(ctx, l, scanner, mm, jc, turnSheetCode, joinSheetData, imageData)
		case turnsheetutil.TurnSheetCodeTypePlayingGame:
			var playSheetData *turnsheetutil.PlayGameTurnSheetCodeData
			playSheetData, err = turnsheetutil.ParsePlayGameTurnSheetCodeData(turnSheetCode)
			if err != nil {
				l.Warn("failed to parse play game turn sheet code >%v<", err)
				return coreerror.NewInvalidDataError("invalid play game turn sheet code format")
			}
			resp, status, err = handlePlayGameTurnSheetUpload(ctx, l, scanner, mm, turnSheetCode, playSheetData, imageData)
		default:
			err = coreerror.NewInvalidDataError("unsupported turn sheet code type: %s", turnSheetCodeType)
		}

		if err != nil {
			l.Warn("failed to process turn sheet upload >%v<", err)
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

// handleJoinGameTurnSheetUpload processes an uploaded join game turn sheet.
func handleJoinGameTurnSheetUpload(ctx context.Context, l logger.Logger, scanner turnsheet.TurnSheetScanner, m *domain.Domain, jc *river.Client[pgx.Tx], turnSheetCode string, turnSheetCodeData *turnsheetutil.JoinGameTurnSheetCodeData, imageData []byte) (*TurnSheetUploadResponse, int, error) {
	l = l.WithFunctionContext("handleJoinGameTurnSheetUpload")

	l.Info("processing join game turn sheet upload to join game subscription ID >%s< turn sheet code >%s<", turnSheetCodeData.GameSubscriptionID, turnSheetCode)

	// Get the game managers game subscription record - bypasses RLS because the turn
	// sheet code itself is the authorization token for this operation.
	gameSubscriptionRec, err := m.GetGameSubscriptionRecByIDForJoinProcess(turnSheetCodeData.GameSubscriptionID)
	if err != nil {
		l.Warn("failed to get game subscription record >%s< for join game turn sheet upload >%v<", turnSheetCodeData.GameSubscriptionID, err)
		return nil, 0, coreerror.NewNotFoundError("game subscription", turnSheetCodeData.GameSubscriptionID)
	}

	l.Info("creating join game turn sheet data for game >%s<", gameSubscriptionRec.GameID)

	gameRec, err := m.GetGameRecByIDForJoinProcess(gameSubscriptionRec.GameID)
	if err != nil {
		l.Warn("failed to get game record >%s< for join game turn sheet upload >%v<", gameSubscriptionRec.GameID, err)
		return nil, 0, coreerror.NewNotFoundError("game", gameSubscriptionRec.GameID)
	}

	// To identify and parse the completed join game turn sheet data from the uploaded image we
	// need to construct the join game data from the game record and turn sheet code
	joinGameData, err := turnsheet.GetTurnSheetJoinGameData(gameRec, turnSheetCode)
	if err != nil {
		l.Warn("failed to create join game data >%v<", err)
		return nil, 0, coreerror.NewInvalidDataError("failed to create join game data: %v", err)
	}

	joinGameDataBytes, err := json.Marshal(joinGameData)
	if err != nil {
		l.Warn("failed to marshal join game sheet data >%v<", err)
		return nil, 0, coreerror.NewInternalError("failed to marshal join game sheet data")
	}

	joinGameScanDataBytes, err := scanner.GetTurnSheetScanData(ctx, l, adventure_game_record.AdventureGameTurnSheetTypeJoinGame, joinGameDataBytes, imageData)
	if err != nil {
		l.Warn("failed to scan join game turn sheet for game >%s< turn sheet code >%s< >%v<", gameRec.ID, turnSheetCode, err)
		return nil, 0, coreerror.NewInvalidDataError("failed to process join game turn sheet: %v", err)
	}

	var joinGameScanData turnsheet.JoinGameScanData
	if err := json.Unmarshal(joinGameScanDataBytes, &joinGameScanData); err != nil {
		l.Warn("failed to unmarshal join game scan data for game >%s< turn sheet code >%s< >%v<", gameRec.ID, turnSheetCode, err)
		return nil, 0, coreerror.NewInvalidDataError("invalid join game turn sheet data")
	}

	accountRec, err := m.GetAccountRecByEmail(joinGameScanData.Email)
	if err != nil {
		l.Warn("failed to get account by email >%s< >%v<", joinGameScanData.Email, err)
		return nil, 0, err
	}

	if accountRec == nil {
		l.Info("creating new pending account for email >%s<", joinGameScanData.Email)
		pendingRec := &account_record.AccountUser{
			Email:  joinGameScanData.Email,
			Status: account_record.AccountUserStatusPendingApproval,
		}
		_, accountRec, _, err = m.CreateAccount(pendingRec)
		if err != nil {
			l.Warn("failed to create pending account >%v<", err)
			return nil, 0, err
		}
	}

	// Create or get account contact
	accountUserContactRec := &account_record.AccountUserContact{
		AccountUserID:      accountRec.ID,
		Name:               joinGameScanData.Name,
		PostalAddressLine1: joinGameScanData.PostalAddressLine1,
		PostalAddressLine2: nullstring.FromString(joinGameScanData.PostalAddressLine2),
		StateProvince:      joinGameScanData.StateProvince,
		Country:            joinGameScanData.Country,
		PostalCode:         joinGameScanData.PostalCode,
	}

	accountUserContactRec, err = m.CreateAccountUserContactRec(accountUserContactRec)
	if err != nil {
		l.Warn("failed to create account contact >%v<", err)
		return nil, 0, err
	}

	subscriptionRec, err := m.UpsertPendingGameSubscriptionForJoinProcess(gameRec.ID, accountRec.AccountID, accountRec.ID, accountUserContactRec.ID, game_record.GameSubscriptionTypePlayer)
	if err != nil {
		l.Warn("failed to upsert game subscription >%v<", err)
		return nil, 0, err
	}

	turnSheetRec := &game_record.GameTurnSheet{
		GameID:           gameRec.ID,
		AccountID:        accountRec.AccountID,
		AccountUserID:    accountRec.ID,
		TurnNumber:       0,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeJoinGame,
		SheetOrder:       1,
		SheetData:        json.RawMessage(joinGameDataBytes),
		ScannedData:      json.RawMessage(joinGameScanDataBytes),
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
	turnSheetRec.ScannedAt = sql.NullTime{Time: time.Now(), Valid: true}

	createdTurnSheetRec, err := m.CreateGameTurnSheetRec(turnSheetRec)
	if err != nil {
		l.Warn("failed to create game turn sheet record >%#v< >%v<", turnSheetRec, err)
		return nil, 0, err
	}

	if _, err := jc.InsertTx(context.Background(), m.Tx, &jobworker.SendGameSubscriptionApprovalEmailWorkerArgs{
		GameSubscriptionID: subscriptionRec.ID,
	}, &river.InsertOpts{Queue: jobqueue.QueueDefault}); err != nil {
		l.Warn("failed to enqueue game subscription approval email job >%v<", err)
	}

	var scannedDataMap map[string]any
	if err := json.Unmarshal(joinGameScanDataBytes, &scannedDataMap); err != nil {
		l.Warn("failed to unmarshal join game scanned data >%v<", err)
		scannedDataMap = make(map[string]any)
	}

	response := &TurnSheetUploadResponse{
		TurnSheetID:      createdTurnSheetRec.ID,
		TurnSheetCode:    turnSheetCode,
		SheetType:        createdTurnSheetRec.SheetType,
		ScannedData:      scannedDataMap,
		ProcessingStatus: createdTurnSheetRec.ProcessingStatus,
	}

	return response, http.StatusAccepted, nil
}

func handlePlayGameTurnSheetUpload(ctx context.Context, l logger.Logger, scanner turnsheet.TurnSheetScanner, m *domain.Domain, turnSheetCode string, turnSheetCodeData *turnsheetutil.PlayGameTurnSheetCodeData, imageData []byte) (*TurnSheetUploadResponse, int, error) {
	l = l.WithFunctionContext("handlePlayGameTurnSheetUpload")

	l.Info("processing turn sheet code >%s< with image data length >%d<", turnSheetCode, len(imageData))

	// Get the turn sheet record
	turnSheetRec, err := m.GetGameTurnSheetRec(turnSheetCodeData.GameTurnSheetID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get turn sheet record >%v<", err)
		return nil, 0, coreerror.NewNotFoundError("turn sheet", turnSheetCodeData.GameTurnSheetID)
	}

	scannedData, err := scanner.GetTurnSheetScanData(ctx, l, turnSheetRec.SheetType, turnSheetRec.SheetData, imageData)
	if err != nil {
		l.Warn("failed to process turn sheet >%v<", err)
		return nil, 0, coreerror.NewInvalidDataError("failed to process turn sheet: %v", err)
	}

	turnSheetRec.ScannedData = json.RawMessage(scannedData)
	turnSheetRec.ScannedAt = sql.NullTime{Time: time.Now(), Valid: true}
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
		TurnSheetID:      turnSheetCodeData.GameTurnSheetID,
		TurnSheetCode:    turnSheetCode,
		SheetType:        turnSheetRec.SheetType,
		ScannedData:      scannedDataMap,
		ProcessingStatus: turnSheetRec.ProcessingStatus,
	}

	return response, http.StatusOK, nil
}

// downloadJoinGameTurnSheetsHandler generates and downloads a join game turn sheet PDF.
//
// When a player joins a game, they are joining a game that has been subscribed to
// by a game manager who will run the game. This handler supports two use cases:
//
//  1. Game manager downloading: If no game_subscription_id query parameter is provided,
//     the handler finds the manager subscription for the authenticated user. This assumes
//     the authenticated user is the game manager downloading the join game turn sheet
//     to print and distribute to players.
//
//  2. Public player browsing: If a game_subscription_id query parameter is provided,
//     the handler uses that manager subscription ID. This allows public players browsing
//     available games to download join game turn sheets for games they want to join.
//
// The join game turn sheet includes a turn sheet code that identifies the game and
// manager subscription, allowing scanned submissions to be properly associated with
// the correct game instance.
func downloadJoinGameTurnSheetsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "downloadJoinGameTurnSheetsHandler")

	l.Info("downloading join game turn sheet with path params >%#v<", pp)

	ctx := r.Context()
	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	// Get game record - bypass RLS since this is a public endpoint
	gameRec, err := mm.GetGameRecByIDForJoinProcess(gameID)
	if err != nil {
		l.Warn("failed to get game record >%s< >%v<", gameID, err)
		return coreerror.NewNotFoundError("game", gameID)
	}

	// Resolve the manager game subscription ID.
	// If provided as a query param, use it directly.
	// Otherwise look up the authenticated user's manager subscription for the game.
	var gameSubscriptionID string
	if gameSubParams, exists := qp.Params["game_subscription_id"]; exists && len(gameSubParams) > 0 {
		gameSubscriptionID, _ = gameSubParams[0].Val.(string)
	}

	if gameSubscriptionID == "" {
		// Try to find manager subscription via authenticated user
		authenData := server.GetRequestAuthenData(l, r)
		if authenData == nil {
			l.Warn("no game_subscription_id provided and no authenticated request data available")
			return coreerror.NewInvalidDataError("game_subscription_id query parameter is required")
		}

		accountID := authenData.AccountUser.AccountID
		subscriptions, err := mm.GetManyGameSubscriptionRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: game_record.FieldGameSubscriptionGameID, Val: gameID},
				{Col: game_record.FieldGameSubscriptionAccountID, Val: accountID},
				{Col: game_record.FieldGameSubscriptionSubscriptionType, Val: game_record.GameSubscriptionTypeManager},
			},
			Limit: 1,
		})
		if err != nil || len(subscriptions) == 0 {
			l.Warn("no manager subscription found for account >%s< and game >%s<", accountID, gameID)
			return coreerror.NewInvalidDataError("game_subscription_id query parameter is required")
		}

		gameSubscriptionID = subscriptions[0].ID
	}

	// Get the manager game subscription record - bypass RLS since this is a public endpoint
	gameSubscriptionRec, err := mm.GetGameSubscriptionRecByIDForJoinProcess(gameSubscriptionID)
	if err != nil {
		l.Warn("failed to get game subscription record >%s< >%v<", gameSubscriptionID, err)
		return coreerror.NewNotFoundError("game subscription", gameSubscriptionID)
	}

	if gameSubscriptionRec.SubscriptionType != game_record.GameSubscriptionTypeManager {
		l.Warn("game subscription >%s< is not a manager subscription", gameSubscriptionID)
		return coreerror.NewInvalidDataError("game subscription must be a manager subscription")
	}

	if gameSubscriptionRec.GameID != gameID {
		l.Warn("game subscription >%s< does not belong to game >%s<", gameSubscriptionID, gameID)
		return coreerror.NewInvalidDataError("game subscription does not belong to this game")
	}

	// Parse config for processor
	cfg, err := config.Parse()
	if err != nil {
		l.Warn("failed to parse config >%v<", err)
		return coreerror.NewInternalError("failed to parse config: %v", err)
	}

	turnSheetCode, err := turnsheetutil.GenerateJoinGameTurnSheetCode(gameSubscriptionID)
	if err != nil {
		l.Warn("failed to generate join turn sheet code >%v<", err)
		return coreerror.NewInternalError("failed to generate turn sheet code: %v", err)
	}

	// Create join game data using mapper function
	joinData, err := turnsheet.GetTurnSheetJoinGameData(gameRec, turnSheetCode)
	if err != nil {
		l.Warn("failed to create join game data >%v<", err)
		return coreerror.NewInvalidDataError("failed to create join game data: %v", err)
	}

	// Get uploaded turn sheet background image and add it to the data
	turnSheetType := adventure_game_record.AdventureGameTurnSheetTypeJoinGame
	backgroundImage, err := mm.GetGameTurnSheetImageDataURL(gameRec.ID, turnSheetType)
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

	// Create join game processor for PDF generation
	processor, err := turnsheet.NewJoinGameProcessor(l, cfg)
	if err != nil {
		l.Warn("failed to create join game processor >%v<", err)
		return coreerror.NewInternalError("failed to create join game processor: %v", err)
	}

	// Generate PDF
	pdfData, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatPDF, sheetDataBytes)
	if err != nil {
		l.Warn("failed to generate join game turn sheet PDF >%v<", err)
		return coreerror.NewInternalError("failed to generate turn sheet PDF: %v", err)
	}

	// Set filename in Content-Disposition header
	filename := fmt.Sprintf("join-game-turn-sheet-%s.pdf", gameRec.Name)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename)) //nolint:gocritic // HTTP Content-Disposition requires specific quote format

	// Return PDF response
	l.Info("responding with join game turn sheet PDF for game >%s< size >%d<", gameRec.ID, len(pdfData))
	if err := server.WritePDFResponse(l, w, http.StatusOK, pdfData); err != nil {
		l.Warn("failed writing PDF response >%v<", err)
		return err
	}

	return nil
}
