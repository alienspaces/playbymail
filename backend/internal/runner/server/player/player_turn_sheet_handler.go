package player

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
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
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
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

// Deprecated subscription-based endpoint names kept only as documentation reference.
// These were removed in favour of the GSI-based endpoints below.
// const (
// 	GetGameSubscriptionTurnSheetList = "get-game-subscription-turn-sheet-list"
// 	GetGameSubscriptionTurnSheet     = "get-game-subscription-turn-sheet"
// 	SaveGameSubscriptionTurnSheet    = "save-game-subscription-turn-sheet"
// 	SubmitGameSubscriptionTurnSheets = "submit-game-subscription-turn-sheets"
// )

// Game subscription instance–scoped endpoint names.
const (
	GetGameSubscriptionInstanceTurnSheetList     = "get-game-subscription-instance-turn-sheet-list"
	GetGameSubscriptionInstanceTurnSheet         = "get-game-subscription-instance-turn-sheet"
	SaveGameSubscriptionInstanceTurnSheet        = "save-game-subscription-instance-turn-sheet"
	SubmitGameSubscriptionInstanceTurnSheets     = "submit-game-subscription-instance-turn-sheets"
	DownloadGameSubscriptionInstanceTurnSheetPDF = "download-game-subscription-instance-turn-sheet-pdf"
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

	// --- Game subscription instance–scoped endpoints ---

	// GET /api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets
	playerTurnSheetConfig[GetGameSubscriptionInstanceTurnSheetList] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets",
		HandlerFunc: getGameSubscriptionInstanceTurnSheetListHandler,
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

	// GET /api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id
	playerTurnSheetConfig[GetGameSubscriptionInstanceTurnSheet] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id",
		HandlerFunc: getGameSubscriptionInstanceTurnSheetHandler,
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

	// PUT /api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id
	playerTurnSheetConfig[SaveGameSubscriptionInstanceTurnSheet] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id",
		HandlerFunc: saveGameSubscriptionInstanceTurnSheetHandler,
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

	// POST /api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheet-upload
	playerTurnSheetConfig[SubmitGameSubscriptionInstanceTurnSheets] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheet-upload",
		HandlerFunc: submitGameSubscriptionInstanceTurnSheetsHandler,
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

	// GET /api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id/download
	playerTurnSheetConfig[DownloadGameSubscriptionInstanceTurnSheetPDF] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:game_turn_sheet_id/download",
		HandlerFunc: downloadGameSubscriptionInstanceTurnSheetPDFHandler,
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

// resolveGameSubscriptionInstance resolves a game_subscription_instance record and verifies ownership for the
// authenticated account. Returns the record or an error.
func resolveGameSubscriptionInstance(l logger.Logger, r *http.Request, pp httprouter.Params, mm *domain.Domain) (*game_record.GameSubscriptionInstance, error) {
	gameSubscriptionInstanceID := pp.ByName("game_subscription_instance_id")
	if gameSubscriptionInstanceID == "" {
		return nil, coreerror.RequiredPathParameter("game_subscription_instance_id")
	}

	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		return nil, coreerror.NewUnauthorizedError()
	}

	gameSubscriptionInstanceRec, err := mm.GetGameSubscriptionInstanceRec(gameSubscriptionInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game subscription instance >%s< >%v<", gameSubscriptionInstanceID, err)
		return nil, err
	}

	// Verify ownership — the authenticated account must own this game subscription instance.
	if gameSubscriptionInstanceRec.AccountID != authData.AccountUser.AccountID {
		l.Warn("game subscription instance >%s< does not belong to authenticated account >%s<", gameSubscriptionInstanceID, authData.AccountUser.AccountID)
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscriptionInstance, gameSubscriptionInstanceID)
	}

	return gameSubscriptionInstanceRec, nil
}

// getGameSubscriptionInstanceTurnSheetListHandler returns the turn sheet list for a game_subscription_instance.
func getGameSubscriptionInstanceTurnSheetListHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameSubscriptionInstanceTurnSheetListHandler")

	mm := m.(*domain.Domain)

	gameSubscriptionInstanceRec, err := resolveGameSubscriptionInstance(l, r, pp, mm)
	if err != nil {
		return err
	}

	authData := server.GetRequestAuthenData(l, r)

	// Resolve the game subscription to get game_id.
	subRec, err := mm.GetGameSubscriptionRec(gameSubscriptionInstanceRec.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionInstanceRec.GameSubscriptionID, err)
		return err
	}

	turnSheetRecs, err := mm.GameTurnSheetRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameTurnSheetAccountID, Val: gameSubscriptionInstanceRec.AccountID},
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

	// Re-sort within each turn by presentation order, which differs from the
	// processing order stored in sheet_order. Players see location choice first
	// (where am I going?) then inventory (what do I take?).
	slices.SortStableFunc(turnSheetRecs, func(a, b *game_record.GameTurnSheet) int {
		if a.TurnNumber != b.TurnNumber {
			return a.TurnNumber - b.TurnNumber
		}
		return adventure_game_record.AdventureGameSheetPresentationOrderForType(a.SheetType) -
			adventure_game_record.AdventureGameSheetPresentationOrderForType(b.SheetType)
	})

	l.Info("returning turn sheet list for game subscription instance >%s<", gameSubscriptionInstanceRec.ID)

	return server.WriteResponse(l, w, http.StatusOK, map[string]interface{}{
		"subscription_id": subRec.ID,
		"game_id":         subRec.GameID,
		"account_id":      gameSubscriptionInstanceRec.AccountID,
		"turn_sheets":     turnSheetRecs,
	})
}

// getGameSubscriptionInstanceTurnSheetHandler returns a specific turn sheet for a game_subscription_instance.
func getGameSubscriptionInstanceTurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameSubscriptionInstanceTurnSheetHandler")

	l.Info("getting game subscription instance turn sheet with path params >%#v<", pp)

	gameTurnSheetID := pp.ByName("game_turn_sheet_id")
	if gameTurnSheetID == "" {
		return coreerror.RequiredPathParameter("game_turn_sheet_id")
	}

	mm := m.(*domain.Domain)

	gameSubscriptionInstanceRec, err := resolveGameSubscriptionInstance(l, r, pp, mm)
	if err != nil {
		return err
	}

	authData := server.GetRequestAuthenData(l, r)

	subRec, err := mm.GetGameSubscriptionRec(gameSubscriptionInstanceRec.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionInstanceRec.GameSubscriptionID, err)
		return err
	}

	turnSheetRec, err := mm.GetGameTurnSheetRec(gameTurnSheetID, nil)
	if err != nil {
		l.Warn("failed to get turn sheet >%s< >%v<", gameTurnSheetID, err)
		return err
	}

	if turnSheetRec.AccountID != gameSubscriptionInstanceRec.AccountID || turnSheetRec.GameID != subRec.GameID {
		l.Warn("turn sheet >%s< does not belong to game subscription instance >%s<", gameTurnSheetID, gameSubscriptionInstanceRec.ID)
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
		l.Info("responding with HTML turn sheet >%s< size >%d<", gameTurnSheetID, len(htmlBytes))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(htmlBytes)
		return err
	}

	l.Info("responding with turn sheet >%s< for game subscription instance >%s<", gameTurnSheetID, gameSubscriptionInstanceRec.ID)

	return server.WriteResponse(l, w, http.StatusOK, turnSheetRec)
}

// saveGameSubscriptionInstanceTurnSheetHandler auto-saves turn sheet form data for a game_subscription_instance.
func saveGameSubscriptionInstanceTurnSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "saveGameSubscriptionInstanceTurnSheetHandler")

	gameTurnSheetID := pp.ByName("game_turn_sheet_id")
	if gameTurnSheetID == "" {
		return coreerror.RequiredPathParameter("game_turn_sheet_id")
	}

	mm := m.(*domain.Domain)

	gameSubscriptionInstanceRec, err := resolveGameSubscriptionInstance(l, r, pp, mm)
	if err != nil {
		return err
	}

	authData := server.GetRequestAuthenData(l, r)

	subRec, err := mm.GetGameSubscriptionRec(gameSubscriptionInstanceRec.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionInstanceRec.GameSubscriptionID, err)
		return err
	}

	turnSheetRec, err := mm.GetGameTurnSheetRec(gameTurnSheetID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get turn sheet >%s< >%v<", gameTurnSheetID, err)
		return err
	}

	if turnSheetRec.AccountID != gameSubscriptionInstanceRec.AccountID || turnSheetRec.GameID != subRec.GameID {
		l.Warn("turn sheet >%s< does not belong to game subscription instance >%s<", gameTurnSheetID, gameSubscriptionInstanceRec.ID)
		return coreerror.NewNotFoundError("turn_sheet", "Turn sheet not found")
	}
	if turnSheetRec.AccountUserID != authData.AccountUser.ID {
		l.Warn("turn sheet >%s< does not belong to authenticated user >%s<", gameTurnSheetID, authData.AccountUser.ID)
		return coreerror.NewNotFoundError("turn_sheet", "Turn sheet not found")
	}

	if turnSheetRec.IsCompleted {
		return coreerror.NewInvalidDataError("turn sheet is already completed and cannot be modified")
	}

	// Read request body
	var req struct {
		ScannedData map[string]any `json:"scanned_data"`
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

	// Validate scanned_data against per-sheet-type schema when one is defined
	if schemaName := turnsheet.ScannedDataSchemaName(turnSheetRec.SheetType); schemaName != "" {
		schema := jsonschema.SchemaWithReferences{
			Main: jsonschema.Schema{
				Location: turnsheet.ScannedDataSchemaLocationForType(turnSheetRec.SheetType),
				Name:     schemaName,
			},
		}
		schema = jsonschema.ResolveSchemaLocation(mm.Config().SchemaPath, schema)
		if err := jsonschema.ValidateJSON(schema, scannedDataBytes); err != nil {
			l.Warn("scanned_data validation failed for sheet type >%s< >%v<", turnSheetRec.SheetType, err)
			return coreerror.NewInvalidDataError("invalid scanned_data for this turn sheet type: %s", err.Error())
		}
	}

	turnSheetRec.ScannedData = scannedDataBytes

	updatedRec, err := mm.UpdateGameTurnSheetRec(turnSheetRec)
	if err != nil {
		l.Warn("failed to update turn sheet >%v<", err)
		return err
	}

	l.Info("saved turn sheet >%s< via game subscription instance >%s<", gameTurnSheetID, gameSubscriptionInstanceRec.ID)

	return server.WriteResponse(l, w, http.StatusOK, updatedRec)
}

// submitGameSubscriptionInstanceTurnSheetsHandler submits all turn sheets for a game_subscription_instance.
// If the game instance has ProcessWhenAllSubmitted enabled and all players have now
// submitted their turn sheets, a GameTurnProcessingWorkerArgs job is enqueued.
func submitGameSubscriptionInstanceTurnSheetsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "submitGameSubscriptionInstanceTurnSheetsHandler")

	mm := m.(*domain.Domain)

	gameSubscriptionInstanceRec, err := resolveGameSubscriptionInstance(l, r, pp, mm)
	if err != nil {
		return err
	}

	authData := server.GetRequestAuthenData(l, r)

	subRec, err := mm.GetGameSubscriptionRec(gameSubscriptionInstanceRec.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionInstanceRec.GameSubscriptionID, err)
		return err
	}

	// Get the game instance to know the current turn
	gameInstanceRec, err := mm.GetGameInstanceRec(gameSubscriptionInstanceRec.GameInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%s< >%v<", gameSubscriptionInstanceRec.GameInstanceID, err)
		return err
	}

	turnSheetRecs, err := mm.GameTurnSheetRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameTurnSheetAccountID, Val: gameSubscriptionInstanceRec.AccountID},
			{Col: game_record.FieldGameTurnSheetAccountUserID, Val: authData.AccountUser.ID},
			{Col: game_record.FieldGameTurnSheetGameID, Val: subRec.GameID},
			{Col: game_record.FieldGameTurnSheetTurnNumber, Val: gameInstanceRec.CurrentTurn},
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

	l.Info("submitted >%d< turn sheets via game subscription instance >%s<", completedCount, gameSubscriptionInstanceRec.ID)

	// Check if we should trigger early turn processing
	if gameInstanceRec.ProcessWhenAllSubmitted && gameInstanceRec.Status == game_record.GameInstanceStatusStarted {
		allSubmitted, checkErr := checkAllTurnSheetsSubmitted(l, mm, gameInstanceRec)
		if checkErr != nil {
			l.Warn("failed to check all turn sheets submitted >%v<", checkErr)
		} else if allSubmitted {
			l.Info("all players submitted for game instance >%s< turn >%d<, enqueueing early turn processing", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)
			_, insertErr := jc.Insert(r.Context(), jobworker.GameTurnProcessingWorkerArgs{
				GameInstanceID: gameInstanceRec.ID,
				TurnNumber:     gameInstanceRec.CurrentTurn,
			}, nil)
			if insertErr != nil {
				l.Warn("failed to enqueue early turn processing >%v<", insertErr)
			}
		}
	}

	return server.WriteResponse(l, w, http.StatusOK, map[string]interface{}{
		"submitted_count": completedCount,
		"total_count":     len(turnSheetRecs),
	})
}

// checkAllTurnSheetsSubmitted checks whether all active players in a game instance
// have completed all their turn sheets for the current turn.
func checkAllTurnSheetsSubmitted(l logger.Logger, mm *domain.Domain, gi *game_record.GameInstance) (bool, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "checkAllTurnSheetsSubmitted")

	// Get all turn sheets for this game instance and current turn
	allSheets, err := mm.GameTurnSheetRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameTurnSheetGameInstanceID, Val: gi.ID},
			{Col: game_record.FieldGameTurnSheetTurnNumber, Val: gi.CurrentTurn},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to get turn sheets for instance >%s< turn >%d<: %w", gi.ID, gi.CurrentTurn, err)
	}

	if len(allSheets) == 0 {
		return false, nil
	}

	for _, ts := range allSheets {
		if !ts.IsCompleted {
			l.Debug("turn sheet >%s< not yet completed for instance >%s< turn >%d<", ts.ID, gi.ID, gi.CurrentTurn)
			return false, nil
		}
	}

	l.Info("all >%d< turn sheets completed for instance >%s< turn >%d<", len(allSheets), gi.ID, gi.CurrentTurn)
	return true, nil
}

// downloadGameSubscriptionInstanceTurnSheetPDFHandler returns a printable PDF for a specific turn sheet so the
// player can fill it in offline and mail it back.
func downloadGameSubscriptionInstanceTurnSheetPDFHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "downloadGameSubscriptionInstanceTurnSheetPDFHandler")

	l.Info("downloading game subscription instance turn sheet PDF with path params >%#v<", pp)

	gameTurnSheetID := pp.ByName("game_turn_sheet_id")
	if gameTurnSheetID == "" {
		return coreerror.RequiredPathParameter("game_turn_sheet_id")
	}

	mm := m.(*domain.Domain)

	gameSubscriptionInstanceRec, err := resolveGameSubscriptionInstance(l, r, pp, mm)
	if err != nil {
		return err
	}

	authData := server.GetRequestAuthenData(l, r)

	subRec, err := mm.GetGameSubscriptionRec(gameSubscriptionInstanceRec.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionInstanceRec.GameSubscriptionID, err)
		return err
	}

	turnSheetRec, err := mm.GetGameTurnSheetRec(gameTurnSheetID, nil)
	if err != nil {
		l.Warn("failed to get turn sheet >%s< >%v<", gameTurnSheetID, err)
		return err
	}

	if turnSheetRec.AccountID != gameSubscriptionInstanceRec.AccountID || turnSheetRec.GameID != subRec.GameID {
		l.Warn("turn sheet >%s< does not belong to game subscription instance >%s<", gameTurnSheetID, gameSubscriptionInstanceRec.ID)
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

	l.Info("responding with turn sheet PDF >%s< size >%d<", gameTurnSheetID, len(pdfBytes))

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
	gameSubscriptionInstanceRec, err := mm.VerifyGameSubscriptionInstanceTurnSheetKey(gameSubscriptionInstanceID, req.TurnSheetToken)
	if err != nil {
		l.Warn("failed to verify turn sheet token >%s< >%v<", req.TurnSheetToken, err)
		return coreerror.NewNotFoundError("turn_sheet_token", "This link is no longer valid")
	}

	// Get the account user for the instance
	accountUserRec, err := mm.GetAccountUserRec(gameSubscriptionInstanceRec.AccountUserID, nil)
	if err != nil {
		l.Warn("failed to get account user ID >%s< >%v<", gameSubscriptionInstanceRec.AccountUserID, err)
		return err
	}

	// Generate session token for the account
	sessionToken, err := mm.GenerateAccountUserSessionToken(accountUserRec)
	if err != nil {
		l.Warn("failed to generate account user session token >%v<", err)
		return err
	}

	l.Info("verified game subscription instance token for account >%s<, session token >%s<", accountUserRec.Email, sessionToken)

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

	// Get the account user for the instance
	accountUserRec, err := mm.GetAccountUserRec(instanceRec.AccountUserID, nil)
	if err != nil {
		l.Warn("failed to get account user >%s< >%v<", instanceRec.AccountUserID, err)
		return err
	}

	// Verify email matches
	if req.Email == "" || accountUserRec.Email != req.Email {
		l.Warn("email >%s< does not match account email >%s<", req.Email, accountUserRec.Email)
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
