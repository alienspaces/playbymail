package player

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/schema/api/player_schema"
)

const (
	VerifyGameSubscriptionToken           = "verify-game-subscription-token"
	RequestGameSubscriptionTurnSheetToken = "request-game-subscription-turn-sheet-token"
)

const (
	GetGameSubscriptionTurnSheetList = "get-game-subscription-turn-sheet-list"
	GetGameSubscriptionTurnSheet     = "get-game-subscription-turn-sheet"
	GetGameSubscriptionTurnSheetHTML = "get-game-subscription-turn-sheet-html"
	SaveGameSubscriptionTurnSheet    = "save-game-subscription-turn-sheet"
	SubmitGameSubscriptionTurnSheets = "submit-game-subscription-turn-sheets"
)

// New gsi_id-based endpoint names (additive — existing endpoints stay unchanged).
const (
	GetGSITurnSheetList    = "get-gsi-turn-sheet-list"
	GetGSITurnSheet        = "get-gsi-turn-sheet"
	SaveGSITurnSheet       = "save-gsi-turn-sheet"
	SubmitGSITurnSheets    = "submit-gsi-turn-sheets"
	DownloadGSITurnSheetPDF = "download-gsi-turn-sheet-pdf"
)

func playerTurnSheetHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "playerTurnSheetHandlerConfig")

	l.Debug("Adding player game subscriptionturn sheet handler configuration")

	playerTurnSheetConfig := make(map[string]server.HandlerConfig)

	// POST "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/verify-token" - verify a game subscription instance token and return a session token
	playerTurnSheetConfig[VerifyGameSubscriptionToken] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/verify-token",
		HandlerFunc: verifyGameSubscriptionTokenHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.verify-game-subscription-token.request.schema.json",
				},
				References: referenceSchemas,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.verify-game-subscription-token.response.schema.json",
				},
				References: referenceSchemas,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Verify game subscription token",
			Description: "Verify a game subscription turn sheet token and return a session token. " +
				"Validates the token and email address, then generates a session token for the account.",
		},
	}

	// POST /api/v1/player/game-subscription-instances/:game_subscription_instance_id/request-token - request a new game subscription instance token if expired
	playerTurnSheetConfig[RequestGameSubscriptionTurnSheetToken] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/request-token",
		HandlerFunc: requestGameSubscriptionTokenHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.request-game-subscription-token.request.schema.json",
				},
				References: referenceSchemas,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Request new turn sheet token",
			Description: "Request a new turn sheet token if the current one has expired. " +
				"Validates email matches the account and generates a new token.",
		},
	}

	// GET /api/v1/player/game-subscriptions/:game_subscription_id/turn-sheets - return turn sheet list
	playerTurnSheetConfig[GetGameSubscriptionTurnSheetList] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-subscriptions/:game_subscription_id/game-instances/:game_instance_id/turn-sheets",
		HandlerFunc: getGameSubscriptionTurnSheetListHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGamePlaying,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.get-game-subscription-turn-sheet-list.response.schema.json",
				},
				References: append(referenceSchemas, []jsonschema.Schema{
					{
						Location: "api/player_schema",
						Name:     "game_turn_sheet.schema.json",
					},
				}...),
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game subscription turn sheet list",
			Description: "Returns a list of available turn sheets for the authenticated account's game subscription. " +
				"Requires session token authentication.",
		},
	}

	// GET /api/v1/player/game-subscriptions/:game_subscription_id/turn-sheets/:game_turn_sheet_id - get specific turn sheet data
	playerTurnSheetConfig[GetGameSubscriptionTurnSheet] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-subscriptions/:game_subscription_id/game-instances/:game_instance_id/turn-sheets/:game_turn_sheet_id",
		HandlerFunc: getGameSubscriptionTurnSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGamePlaying,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.get-game-subscription-turn-sheet.response.schema.json",
				},
				References: append(referenceSchemas, []jsonschema.Schema{
					{
						Location: "api/player_schema",
						Name:     "game_turn_sheet.schema.json",
					},
				}...),
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get turn sheet data",
			Description: "Get specific turn sheet data for the authenticated account's game subscription. " +
				"Requires session token authentication.",
		},
	}

	// PUT /api/v1/player/game-subscriptions/:game_subscription_id/turn-sheets/:game_turn_sheet_id - save form data
	playerTurnSheetConfig[SaveGameSubscriptionTurnSheet] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/player/game-subscriptions/:game_subscription_id/game-instances/:game_instance_id/turn-sheets/:game_turn_sheet_id",
		HandlerFunc: updateGameSubscriptionTurnSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGamePlaying,
			},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.update-game-subscription-turn-sheet.request.schema.json",
				},
				References: referenceSchemas,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.update-game-subscription-turn-sheet.response.schema.json",
				},
				References: append(referenceSchemas, []jsonschema.Schema{
					{
						Location: "api/player_schema",
						Name:     "game_turn_sheet.schema.json",
					},
				}...),
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Save turn sheet data",
			Description: "Save form data for a turn sheet. Supports incremental/auto-save. " +
				"Requires session token authentication.",
		},
	}

	// POST /api/v1/player/game-subscriptions/.../turn-sheet-upload - submit all sheets
	playerTurnSheetConfig[SubmitGameSubscriptionTurnSheets] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/player/game-subscriptions/:game_subscription_id/game-instances/:game_instance_id/turn-sheet-upload",
		HandlerFunc: submitGameSubscriptionTurnSheetsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGamePlaying,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.submit-game-subscription-turn-sheets.response.schema.json",
				},
				References: referenceSchemas,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Submit all turn sheets",
			Description: "Submit all turn sheets for the authenticated account's game subscription. " +
				"Locks all sheets and marks them as completed. Requires session token authentication.",
		},
	}

	// --- gsi_id-based endpoints (additive; existing subscription+instance endpoints unchanged) ---

	// GET /api/v1/player/game-subscription-instances/:gsi_id/turn-sheets
	playerTurnSheetConfig[GetGSITurnSheetList] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets",
		HandlerFunc: getGSITurnSheetListHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGamePlaying,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.get-game-subscription-turn-sheet-list.response.schema.json",
				},
				References: append(referenceSchemas, []jsonschema.Schema{
					{Location: "api/player_schema", Name: "game_turn_sheet.schema.json"},
				}...),
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "List turn sheets by game subscription instance",
			Description: "Returns available turn sheets for a game_subscription_instance. Auth: session token.",
		},
	}

	// GET /api/v1/player/game-subscription-instances/:gsi_id/turn-sheets/:game_turn_sheet_id
	playerTurnSheetConfig[GetGSITurnSheet] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id",
		HandlerFunc: getGSITurnSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGamePlaying,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.get-game-subscription-turn-sheet.response.schema.json",
				},
				References: append(referenceSchemas, []jsonschema.Schema{
					{Location: "api/player_schema", Name: "game_turn_sheet.schema.json"},
				}...),
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Get turn sheet by game subscription instance",
			Description: "Get specific turn sheet data via game_subscription_instance_id. Auth: session token.",
		},
	}

	// PUT /api/v1/player/game-subscription-instances/:gsi_id/turn-sheets/:game_turn_sheet_id
	playerTurnSheetConfig[SaveGSITurnSheet] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id",
		HandlerFunc: saveGSITurnSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGamePlaying,
			},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.update-game-subscription-turn-sheet.request.schema.json",
				},
				References: referenceSchemas,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.update-game-subscription-turn-sheet.response.schema.json",
				},
				References: append(referenceSchemas, []jsonschema.Schema{
					{Location: "api/player_schema", Name: "game_turn_sheet.schema.json"},
				}...),
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Save turn sheet by game subscription instance",
			Description: "Save (auto-save) turn sheet data via game_subscription_instance_id. Auth: session token.",
		},
	}

	// POST /api/v1/player/game-subscription-instances/:gsi_id/turn-sheet-upload
	playerTurnSheetConfig[SubmitGSITurnSheets] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheet-upload",
		HandlerFunc: submitGSITurnSheetsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGamePlaying,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.submit-game-subscription-turn-sheets.response.schema.json",
				},
				References: referenceSchemas,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Submit turn sheets by game subscription instance",
			Description: "Submit all turn sheets via game_subscription_instance_id. Auth: session token.",
		},
	}

	// GET /api/v1/player/game-subscription-instances/:gsi_id/turn-sheets/:id/download
	playerTurnSheetConfig[DownloadGSITurnSheetPDF] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id/download",
		HandlerFunc: downloadGSITurnSheetPDFHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGamePlaying,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Download turn sheet PDF",
			Description: "Download a printable PDF of a specific turn sheet via game_subscription_instance_id. Auth: session token.",
		},
	}

	return playerTurnSheetConfig, nil
}

// resolveGSI resolves a game_subscription_instance record and verifies ownership for the
// authenticated account. Returns the gsi record or an error.
func resolveGSI(l logger.Logger, r *http.Request, pp httprouter.Params, mm *domain.Domain) (*game_record.GameSubscriptionInstance, error) {
	gsiID := pp.ByName("game_subscription_instance_id")
	if gsiID == "" {
		return nil, coreerror.RequiredPathParameter("game_subscription_instance_id")
	}

	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		return nil, coreerror.NewUnauthorizedError()
	}

	gsiRec, err := mm.GetGameSubscriptionInstanceRec(gsiID, nil)
	if err != nil {
		l.Warn("failed to get game subscription instance >%s< >%v<", gsiID, err)
		return nil, err
	}

	// Verify ownership — the authenticated account must own this gsi.
	if gsiRec.AccountID != authData.AccountUser.AccountID {
		l.Warn("game subscription instance >%s< does not belong to authenticated account >%s<", gsiID, authData.AccountUser.AccountID)
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscriptionInstance, gsiID)
	}

	return gsiRec, nil
}

// getGSITurnSheetListHandler returns the turn sheet list for a game_subscription_instance.
func getGSITurnSheetListHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGSITurnSheetListHandler")

	mm := m.(*domain.Domain)

	gsiRec, err := resolveGSI(l, r, pp, mm)
	if err != nil {
		return err
	}

	authData := server.GetRequestAuthenData(l, r)

	// Resolve the game subscription to get game_id.
	subRec, err := mm.GetGameSubscriptionRec(gsiRec.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gsiRec.GameSubscriptionID, err)
		return err
	}

	turnSheetRecs, err := mm.GameTurnSheetRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameTurnSheetAccountID, Val: gsiRec.AccountID},
			{Col: game_record.FieldGameTurnSheetAccountUserID, Val: authData.AccountUser.ID},
			{Col: game_record.FieldGameTurnSheetGameID, Val: subRec.GameID},
		},
		OrderBy: []coresql.OrderBy{
			{Col: game_record.FieldGameTurnSheetTurnNumber, Direction: coresql.OrderDirectionASC},
			{Col: game_record.FieldGameTurnSheetSheetOrder, Direction: coresql.OrderDirectionASC},
		},
	})
	if err != nil {
		l.Warn("failed to get turn sheets >%v<", err)
		return err
	}

	l.Info("returning turn sheet list for gsi >%s<", gsiRec.ID)

	return server.WriteResponse(l, w, http.StatusOK, map[string]interface{}{
		"subscription_id": subRec.ID,
		"game_id":         subRec.GameID,
		"account_id":      gsiRec.AccountID,
		"turn_sheets":     turnSheetRecs,
	})
}

// getGSITurnSheetHandler returns a specific turn sheet for a game_subscription_instance.
func getGSITurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGSITurnSheetHandler")

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

	turnSheetRec, err := mm.GetGameTurnSheetRec(gameTurnSheetID, nil)
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

	acceptHeader := r.Header.Get("Accept")
	if strings.Contains(acceptHeader, "text/html") {
		cfg := mm.Config()
		processor, err := turnsheet.GetDocumentProcessor(l, cfg, turnSheetRec.SheetType)
		if err != nil {
			return err
		}
		htmlBytes, err := processor.GenerateTurnSheet(r.Context(), l, turnsheet.DocumentFormatHTML, turnSheetRec.SheetData)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(htmlBytes)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, turnSheetRec)
}

// saveGSITurnSheetHandler auto-saves turn sheet form data for a game_subscription_instance.
func saveGSITurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "saveGSITurnSheetHandler")

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

	var req struct {
		ScannedData map[string]interface{} `json:"scanned_data"`
	}
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	scannedDataBytes, err := json.Marshal(req.ScannedData)
	if err != nil {
		l.Warn("failed to marshal scanned data >%v<", err)
		return coreerror.NewInvalidDataError("invalid scanned_data format")
	}

	turnSheetRec.ScannedData = scannedDataBytes

	updatedRec, err := mm.UpdateGameTurnSheetRec(turnSheetRec)
	if err != nil {
		l.Warn("failed to update turn sheet >%v<", err)
		return err
	}

	l.Info("saved turn sheet >%s< via gsi >%s<", gameTurnSheetID, gsiRec.ID)

	return server.WriteResponse(l, w, http.StatusOK, updatedRec)
}

// submitGSITurnSheetsHandler submits all turn sheets for a game_subscription_instance.
func submitGSITurnSheetsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "submitGSITurnSheetsHandler")

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

	turnSheetRecs, err := mm.GameTurnSheetRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameTurnSheetAccountID, Val: gsiRec.AccountID},
			{Col: game_record.FieldGameTurnSheetAccountUserID, Val: authData.AccountUser.ID},
			{Col: game_record.FieldGameTurnSheetGameID, Val: subRec.GameID},
		},
	})
	if err != nil {
		l.Warn("failed to get game subscription turn sheet records >%v<", err)
		return err
	}

	now := time.Now()
	completedCount := 0
	for _, ts := range turnSheetRecs {
		if ts.IsCompleted {
			continue
		}
		ts.IsCompleted = true
		ts.CompletedAt = sql.NullTime{Time: now, Valid: true}
		if _, err := mm.UpdateGameTurnSheetRec(ts); err != nil {
			l.Warn("failed to update turn sheet >%s< >%v<", ts.ID, err)
			continue
		}
		completedCount++
	}

	l.Info("submitted >%d< turn sheets via gsi >%s<", completedCount, gsiRec.ID)

	return server.WriteResponse(l, w, http.StatusOK, map[string]interface{}{
		"submitted_count": completedCount,
		"total_count":     len(turnSheetRecs),
	})
}

// downloadGSITurnSheetPDFHandler returns a printable PDF for a specific turn sheet so the
// player can fill it in offline and mail it back.
func downloadGSITurnSheetPDFHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "downloadGSITurnSheetPDFHandler")

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

	turnSheetRec, err := mm.GetGameTurnSheetRec(gameTurnSheetID, nil)
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

	cfg := mm.Config()
	processor, err := turnsheet.GetDocumentProcessor(l, cfg, turnSheetRec.SheetType)
	if err != nil {
		l.Warn("failed to get document processor for sheet type >%s< >%v<", turnSheetRec.SheetType, err)
		return err
	}

	pdfBytes, err := processor.GenerateTurnSheet(r.Context(), l, turnsheet.DocumentFormatPDF, turnSheetRec.SheetData)
	if err != nil {
		l.Warn("failed to generate PDF for turn sheet >%s< >%v<", gameTurnSheetID, err)
		return err
	}

	filename := fmt.Sprintf("turn-sheet-%s.pdf", gameTurnSheetID)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(pdfBytes)
	return err
}

// verifyGameSubscriptionTokenHandler verifies a game subscription instance token and returns a session token
func verifyGameSubscriptionTokenHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "verifyGameSubscriptionTokenHandler")

	gameSubscriptionInstanceID := pp.ByName("game_subscription_instance_id")
	if gameSubscriptionInstanceID == "" {
		l.Warn("game subscription instance id is empty")
		return coreerror.RequiredPathParameter("game_subscription_instance_id")
	}

	// Read request body
	var req player_schema.VerifyGameSubscriptionTokenRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	// Verify the turn sheet token
	instanceRec, err := mm.VerifyGameSubscriptionInstanceTurnSheetKey(gameSubscriptionInstanceID, req.TurnSheetToken)
	if err != nil {
		l.Warn("failed to verify turn sheet token >%s< >%v<", req.TurnSheetToken, err)
		return coreerror.NewNotFoundError("turn_sheet_token", "This link is no longer valid")
	}

	// Get the account for the instance
	accountRec, err := mm.GetAccountRec(instanceRec.AccountID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get account >%v<", err)
		return err
	}

	// Verify email matches if provided (optional for direct-link access)
	if req.Email != "" && accountRec.Email != req.Email {
		l.Warn("email >%s< does not match account email >%s<", req.Email, accountRec.Email)
		return coreerror.NewInvalidDataError("email does not match the account for this subscription")
	}

	// Generate session token for the account
	sessionToken, err := mm.GenerateAccountSessionToken(accountRec)
	if err != nil {
		l.Warn("failed to generate session token >%v<", err)
		return err
	}

	l.Info("verified game subscription instance token for account >%s<, session token >%s<", accountRec.Email, sessionToken)

	return server.WriteResponse(l, w, http.StatusOK, mapper.MapVerifyGameSubscriptionTokenResponse(sessionToken))
}

// requestGameSubscriptionTokenHandler requests a new turn sheet token
func requestGameSubscriptionTokenHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "requestGameSubscriptionTokenHandler")

	gameSubscriptionInstanceID := pp.ByName("game_subscription_instance_id")
	if gameSubscriptionInstanceID == "" {
		l.Warn("game subscription instance id is empty")
		return coreerror.RequiredPathParameter("game_subscription_instance_id")
	}

	// Read request body
	var req player_schema.RequestGameSubscriptionTokenRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	// Get the game subscription instance record
	instanceRec, err := mm.GetGameSubscriptionInstanceRec(gameSubscriptionInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game subscription instance >%s< >%v<", gameSubscriptionInstanceID, err)
		return err
	}

	// Get the account record for the instance
	accountRec, err := mm.GetAccountRec(instanceRec.AccountID, nil)
	if err != nil {
		l.Warn("failed to get account >%s< >%v<", instanceRec.AccountID, err)
		return err
	}

	// Verify email matches
	if accountRec.Email != req.Email {
		l.Warn("email >%s< does not match account email >%s<", req.Email, accountRec.Email)
		return coreerror.NewInvalidDataError("email does not match the account for this subscription")
	}

	// Get the game instance to get current turn number
	gameInstanceRec, err := mm.GetGameInstanceRec(instanceRec.GameInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%s< >%v<", instanceRec.GameInstanceID, err)
		return err
	}

	// Register SendTurnSheetNotificationEmailWorkerArgs job to send a notification email to the account
	args := jobworker.SendTurnSheetNotificationEmailWorkerArgs{
		GameSubscriptionInstanceID: gameSubscriptionInstanceID,
		TurnNumber:                 gameInstanceRec.CurrentTurn,
	}

	_, err = jc.Insert(r.Context(), args, nil)
	if err != nil {
		l.Warn("failed to queue turn sheet notification email job for instance >%s< >%v<", gameSubscriptionInstanceID, err)
		return err
	}

	l.Info("queued turn sheet notification email job for instance >%s<", gameSubscriptionInstanceID)

	return server.WriteResponse(l, w, http.StatusOK, nil)
}

// getGameSubscriptionTurnSheetListHandler returns the list of turn sheets for the authenticated account's game subscription
func getGameSubscriptionTurnSheetListHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameSubscriptionTurnSheetListHandler")

	gameSubscriptionID := pp.ByName("game_subscription_id")
	if gameSubscriptionID == "" {
		l.Warn("game subscription id is empty")
		return coreerror.RequiredPathParameter("game_subscription_id")
	}

	// Get the authenticated account from the request context
	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		l.Warn("failed getting authenticated account data")
		return server.WriteResponse(l, w, http.StatusUnauthorized, nil)
	}

	mm := m.(*domain.Domain)

	// Get the game subscription record
	gameSubscriptionRec, err := mm.GetGameSubscriptionRec(gameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionID, err)
		return err
	}

	// Verify the subscription belongs to the authenticated account
	if gameSubscriptionRec.AccountID != authData.AccountUser.AccountID {
		l.Warn("game subscription >%s< does not belong to authenticated account >%s<", gameSubscriptionID, authData.AccountUser.AccountID)
		return coreerror.NewNotFoundError("game_subscription", "Game subscription not found")
	}

	// Get turn sheets for this authenticated user, account, and game
	turnSheetRecs, err := mm.GameTurnSheetRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: game_record.FieldGameTurnSheetAccountID,
				Val: gameSubscriptionRec.AccountID,
			},
			{
				Col: game_record.FieldGameTurnSheetAccountUserID,
				Val: authData.AccountUser.ID,
			},
			{
				Col: game_record.FieldGameTurnSheetGameID,
				Val: gameSubscriptionRec.GameID,
			},
		},
		OrderBy: []coresql.OrderBy{
			{Col: game_record.FieldGameTurnSheetTurnNumber, Direction: coresql.OrderDirectionASC},
			{Col: game_record.FieldGameTurnSheetSheetOrder, Direction: coresql.OrderDirectionASC},
		},
	})
	if err != nil {
		l.Warn("failed to get turn sheets >%v<", err)
		return err
	}

	l.Info("returning turn sheet list for subscription >%s<", gameSubscriptionRec.ID)

	// Build response
	response := map[string]interface{}{
		"subscription_id": gameSubscriptionRec.ID,
		"game_id":         gameSubscriptionRec.GameID,
		"account_id":      gameSubscriptionRec.AccountID,
		"turn_sheets":     turnSheetRecs,
	}

	return server.WriteResponse(l, w, http.StatusOK, response)
}

// getGameSubscriptionTurnSheetHandler gets specific turn sheet data for the authenticated account
func getGameSubscriptionTurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameSubscriptionTurnSheetHandler")

	gameSubscriptionID := pp.ByName("game_subscription_id")
	if gameSubscriptionID == "" {
		l.Warn("game subscription id is empty")
		return coreerror.RequiredPathParameter("game_subscription_id")
	}

	gameTurnSheetID := pp.ByName("game_turn_sheet_id")
	if gameTurnSheetID == "" {
		l.Warn("game turn sheet id is empty")
		return coreerror.RequiredPathParameter("game_turn_sheet_id")
	}

	// Get the authenticated account from the request context
	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		l.Warn("failed getting authenticated account data")
		return server.WriteResponse(l, w, http.StatusUnauthorized, nil)
	}

	mm := m.(*domain.Domain)

	// Get the game subscription record
	gameSubscriptionRec, err := mm.GetGameSubscriptionRec(gameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionID, err)
		return err
	}

	// Verify the subscription belongs to the authenticated account
	if gameSubscriptionRec.AccountID != authData.AccountUser.AccountID {
		l.Warn("game subscription >%s< does not belong to authenticated account >%s<", gameSubscriptionID, authData.AccountUser.AccountID)
		return coreerror.NewNotFoundError("game_subscription", "Game subscription not found")
	}

	// Get the turn sheet
	turnSheetRec, err := mm.GetGameTurnSheetRec(gameTurnSheetID, nil)
	if err != nil {
		l.Warn("failed to get turn sheet >%s< >%v<", gameTurnSheetID, err)
		return err
	}

	// Verify the turn sheet belongs to this subscription and to the authenticated user
	if turnSheetRec.AccountID != gameSubscriptionRec.AccountID || turnSheetRec.GameID != gameSubscriptionRec.GameID {
		l.Warn("turn sheet >%s< does not belong to subscription >%s<", gameTurnSheetID, gameSubscriptionRec.ID)
		return coreerror.NewNotFoundError("turn_sheet", "Turn sheet not found")
	}
	if turnSheetRec.AccountUserID != authData.AccountUser.ID {
		l.Warn("turn sheet >%s< does not belong to authenticated user >%s<", gameTurnSheetID, authData.AccountUser.ID)
		return coreerror.NewNotFoundError("turn_sheet", "Turn sheet not found")
	}

	// Check Accept header to determine response format
	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "text/html" || acceptHeader == "text/html, */*" || strings.Contains(acceptHeader, "text/html") {
		l.Info("returning HTML for turn sheet >%s<", gameTurnSheetID)

		// TODO: This is the only place we are accessing the config from the domain object. Perhaps config should be
		// passed in as a parameter to the handler as an argument consistently before the logger.
		cfg := mm.Config()

		// Return HTML format
		processor, err := turnsheet.GetDocumentProcessor(l, cfg, turnSheetRec.SheetType)
		if err != nil {
			l.Warn("failed to get document processor for sheet type >%s< >%v<", turnSheetRec.SheetType, err)
			return err
		}

		// Generate HTML from the turn sheet data
		ctx := r.Context()
		htmlBytes, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatHTML, turnSheetRec.SheetData)
		if err != nil {
			l.Warn("failed to generate HTML for turn sheet >%s< >%v<", gameTurnSheetID, err)
			return err
		}

		// Set content type and return HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(htmlBytes)
		if err != nil {
			l.Warn("failed to write HTML response >%v<", err)
			return err
		}

		l.Info("returned HTML for turn sheet >%s<", gameTurnSheetID)
		return nil
	}

	// Return JSON format (default)
	l.Info("returning turn sheet >%s< for subscription >%s<", gameTurnSheetID, gameSubscriptionRec.ID)

	return server.WriteResponse(l, w, http.StatusOK, turnSheetRec)
}

// updateGameSubscriptionTurnSheetHandler updates form data for a turn sheet
func updateGameSubscriptionTurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateGameSubscriptionTurnSheetHandler")

	gameSubscriptionID := pp.ByName("game_subscription_id")
	if gameSubscriptionID == "" {
		l.Warn("game subscription id is empty")
		return coreerror.RequiredPathParameter("game_subscription_id")
	}

	gameTurnSheetID := pp.ByName("game_turn_sheet_id")
	if gameTurnSheetID == "" {
		l.Warn("game turn sheet id is empty")
		return coreerror.RequiredPathParameter("game_turn_sheet_id")
	}

	// Get the authenticated account from the request context
	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		l.Warn("failed getting authenticated account data")
		return server.WriteResponse(l, w, http.StatusUnauthorized, nil)
	}

	mm := m.(*domain.Domain)

	// Get the game subscription record
	gameSubscriptionRec, err := mm.GetGameSubscriptionRec(gameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionID, err)
		return err
	}

	// Verify the subscription belongs to the authenticated account
	if gameSubscriptionRec.AccountID != authData.AccountUser.AccountID {
		l.Warn("game subscription >%s< does not belong to authenticated account >%s<", gameSubscriptionID, authData.AccountUser.AccountID)
		return coreerror.NewNotFoundError("game_subscription", "Game subscription not found")
	}

	// Get the turn sheet
	turnSheetRec, err := mm.GetGameTurnSheetRec(gameTurnSheetID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get turn sheet >%s< >%v<", gameTurnSheetID, err)
		return err
	}

	// Verify the turn sheet belongs to this subscription and to the authenticated user
	if turnSheetRec.AccountID != gameSubscriptionRec.AccountID || turnSheetRec.GameID != gameSubscriptionRec.GameID {
		l.Warn("turn sheet >%s< does not belong to subscription >%s<", gameTurnSheetID, gameSubscriptionRec.ID)
		return coreerror.NewNotFoundError("turn_sheet", "Turn sheet not found")
	}
	if turnSheetRec.AccountUserID != authData.AccountUser.ID {
		l.Warn("turn sheet >%s< does not belong to authenticated user >%s<", gameTurnSheetID, authData.AccountUser.ID)
		return coreerror.NewNotFoundError("turn_sheet", "Turn sheet not found")
	}

	// Check if already completed
	if turnSheetRec.IsCompleted {
		return coreerror.NewInvalidDataError("turn sheet is already completed and cannot be modified")
	}

	// Read scanned data from request body
	var req struct {
		ScannedData map[string]interface{} `json:"scanned_data"`
	}
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	// Convert scanned data to JSON
	scannedDataBytes, err := json.Marshal(req.ScannedData)
	if err != nil {
		l.Warn("failed to marshal scanned data >%v<", err)
		return coreerror.NewInvalidDataError("invalid scanned_data format")
	}

	// Update turn sheet with scanned data
	turnSheetRec.ScannedData = scannedDataBytes

	// Update the turn sheet
	updatedRec, err := mm.UpdateGameTurnSheetRec(turnSheetRec)
	if err != nil {
		l.Warn("failed to update turn sheet >%v<", err)
		return err
	}

	l.Info("saved turn sheet >%s<", gameTurnSheetID)

	return server.WriteResponse(l, w, http.StatusOK, updatedRec)
}

// submitGameSubscriptionTurnSheetsHandler submits all turn sheets for the authenticated account's game subscription
func submitGameSubscriptionTurnSheetsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "submitGameSubscriptionTurnSheetsHandler")

	gameSubscriptionID := pp.ByName("game_subscription_id")
	if gameSubscriptionID == "" {
		l.Warn("game subscription id is empty")
		return coreerror.RequiredPathParameter("game_subscription_id")
	}

	// Get the authenticated account from the request context
	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		l.Warn("failed getting authenticated account data")
		return server.WriteResponse(l, w, http.StatusUnauthorized, nil)
	}

	mm := m.(*domain.Domain)

	// Get the game subscription record
	gameSubscriptionRec, err := mm.GetGameSubscriptionRec(gameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionID, err)
		return err
	}

	// Verify the subscription belongs to the authenticated account
	if gameSubscriptionRec.AccountID != authData.AccountUser.AccountID {
		l.Warn("game subscription >%s< does not belong to authenticated account >%s<", gameSubscriptionID, authData.AccountUser.AccountID)
		return coreerror.NewNotFoundError("game_subscription", "Game subscription not found")
	}

	// Get turn sheets for this authenticated user's subscription
	turnSheetRecs, err := mm.GameTurnSheetRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: game_record.FieldGameTurnSheetAccountID,
				Val: gameSubscriptionRec.AccountID,
			},
			{
				Col: game_record.FieldGameTurnSheetAccountUserID,
				Val: authData.AccountUser.ID,
			},
			{
				Col: game_record.FieldGameTurnSheetGameID,
				Val: gameSubscriptionRec.GameID,
			},
		},
	})
	if err != nil {
		l.Warn("failed to get game subscription turn sheet records >%v<", err)
		return err
	}

	// Mark all game subscription turn sheet records as completed
	now := time.Now()
	completedCount := 0
	for _, turnSheetRec := range turnSheetRecs {
		if turnSheetRec.IsCompleted {
			continue // Already completed
		}

		turnSheetRec.IsCompleted = true
		turnSheetRec.CompletedAt = sql.NullTime{Time: now, Valid: true}

		_, err := mm.UpdateGameTurnSheetRec(turnSheetRec)
		if err != nil {
			l.Warn("failed to update turn sheet >%s< >%v<", turnSheetRec.ID, err)
			continue
		}

		completedCount++
	}

	l.Info("submitted >%d< game subscription turn sheet records for subscription >%s<", completedCount, gameSubscriptionRec.ID)

	return server.WriteResponse(l, w, http.StatusOK, map[string]interface{}{
		"submitted_count": completedCount,
		"total_count":     len(turnSheetRecs),
	})
}
