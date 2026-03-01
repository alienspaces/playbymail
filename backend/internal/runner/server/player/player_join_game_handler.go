package player

import (
	"database/sql"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/schema/api/player_schema"
)

const (
	GetJoinGameInfo         = "get-join-game-info"
	VerifyJoinGameEmail     = "verify-join-game-email"
	GetJoinGameSheet        = "get-join-game-sheet"
	SaveJoinGameSheet       = "save-join-game-sheet"
	SubmitJoinGame          = "submit-join-game"
)

func playerJoinGameHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "playerJoinGameHandlerConfig")

	l.Debug("Adding player join game handler configuration")

	cfg := make(map[string]server.HandlerConfig)

	cfg[GetJoinGameInfo] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-instances/:game_instance_id/join-game",
		HandlerFunc: getJoinGameInfoHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypePublic},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.join-game-info.response.schema.json",
				},
				References: referenceSchemas,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Get join game info",
			Description: "Returns game and instance information for displaying the join game form.",
		},
	}

	cfg[VerifyJoinGameEmail] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/player/game-instances/:game_instance_id/join-game/verify-email",
		HandlerFunc: verifyJoinGameEmailHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypePublic},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.join-game-verify-email.request.schema.json",
				},
				References: referenceSchemas,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.join-game-verify-email.response.schema.json",
				},
				References: referenceSchemas,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Verify join game email",
			Description: "Checks whether the provided email already has an account on the platform.",
		},
	}

	cfg[GetJoinGameSheet] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/player/game-instances/:game_instance_id/join-game/sheet",
		HandlerFunc: getJoinGameSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypePublic},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Get join game turn sheet",
			Description: "Returns the join game turn sheet as HTML for online completion or printing.",
		},
	}

	cfg[SaveJoinGameSheet] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/player/game-instances/:game_instance_id/join-game/sheet",
		HandlerFunc: saveJoinGameSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypePublic},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Save join game sheet",
			Description: "Saves partial join game sheet data for multi-step form completion.",
		},
	}

	cfg[SubmitJoinGame] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/player/game-instances/:game_instance_id/join-game",
		HandlerFunc: submitJoinGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypePublic},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.join-game-submit.request.schema.json",
				},
				References: referenceSchemas,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/player_schema",
					Name:     "player.join-game-submit.response.schema.json",
				},
				References: referenceSchemas,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Submit join game",
			Description: "Submits the join game form. Creates an account if needed, creates a player subscription, and assigns the player to the game instance.",
		},
	}

	return cfg, nil
}

// getJoinGameInfoHandler returns game and instance information needed to render the join form.
func getJoinGameInfoHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getJoinGameInfoHandler")

	instanceID := pp.ByName("game_instance_id")
	if instanceID == "" {
		return coreerror.RequiredPathParameter("game_instance_id")
	}

	mm := m.(*domain.Domain)

	instanceRec, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%s< >%v<", instanceID, err)
		return err
	}

	if instanceRec.Status != game_record.GameInstanceStatusCreated {
		l.Warn("game instance >%s< is not open for enrollment (status >%s<)", instanceID, instanceRec.Status)
		return coreerror.NewNotFoundError(game_record.TableGameInstance, instanceID)
	}

	hasCapacity, err := mm.HasAvailableCapacity(instanceID)
	if err != nil {
		l.Warn("failed to check capacity for instance >%s< >%v<", instanceID, err)
		return err
	}
	if !hasCapacity {
		l.Warn("game instance >%s< has no available capacity", instanceID)
		return coreerror.NewInvalidDataError("this game instance is no longer accepting new players")
	}

	gameRec, err := mm.GetGameRec(instanceRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game record for instance >%s< >%v<", instanceID, err)
		return err
	}

	playerCount, err := mm.GetPlayerCountForGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to get player count for instance >%s< >%v<", instanceID, err)
		return err
	}

	res := player_schema.JoinGameInfoResponse{
		Data: &player_schema.JoinGameInfoResponseData{
			GameID:            gameRec.ID,
			GameName:          gameRec.Name,
			GameDescription:   gameRec.Description,
			GameType:          gameRec.GameType,
			TurnDurationHours: gameRec.TurnDurationHours,
			Instance: &player_schema.JoinGameInfoInstanceData{
				ID:                    instanceRec.ID,
				RequiredPlayerCount:   instanceRec.RequiredPlayerCount,
				PlayerCount:           playerCount,
				DeliveryPhysicalPost:  instanceRec.DeliveryPhysicalPost,
				DeliveryPhysicalLocal: instanceRec.DeliveryPhysicalLocal,
				DeliveryEmail:         instanceRec.DeliveryEmail,
			},
		},
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

// verifyJoinGameEmailHandler checks whether the provided email already has an account.
func verifyJoinGameEmailHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "verifyJoinGameEmailHandler")

	var req player_schema.JoinGameVerifyEmailRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed to read request >%v<", err)
		return err
	}

	if req.Email == "" {
		return coreerror.NewInvalidDataError("email is required")
	}

	mm := m.(*domain.Domain)

	accountRec, err := mm.GetAccountRecByEmail(req.Email)
	if err != nil {
		l.Warn("failed to check account by email >%v<", err)
		return err
	}

	res := player_schema.JoinGameVerifyEmailResponse{
		Data: &player_schema.JoinGameVerifyEmailResponseData{
			HasAccount: accountRec != nil,
		},
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

// getJoinGameSheetHandler returns the join game turn sheet as HTML.
func getJoinGameSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getJoinGameSheetHandler")

	instanceID := pp.ByName("game_instance_id")
	if instanceID == "" {
		return coreerror.RequiredPathParameter("game_instance_id")
	}

	mm := m.(*domain.Domain)

	instanceRec, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%s< >%v<", instanceID, err)
		return err
	}

	gameRec, err := mm.GetGameRec(instanceRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game for instance >%s< >%v<", instanceID, err)
		return err
	}

	cfg := mm.Config()

	// For now only adventure games are supported; extend when new game types are added.
	joinGameSheetType := adventure_game_record.AdventureGameTurnSheetTypeJoinGame

	processor, err := turnsheet.GetDocumentProcessor(l, cfg, joinGameSheetType)
	if err != nil {
		l.Warn("failed to get join game processor for sheet type >%s< >%v<", joinGameSheetType, err)
		return err
	}

	sheetData, err := processor.GeneratePreviewData(r.Context(), l, gameRec, nil)
	if err != nil {
		l.Warn("failed to generate join game preview data >%v<", err)
		return err
	}

	htmlBytes, err := processor.GenerateTurnSheet(r.Context(), l, turnsheet.DocumentFormatHTML, sheetData)
	if err != nil {
		l.Warn("failed to generate join game turn sheet HTML >%v<", err)
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(htmlBytes)
	return err
}

// saveJoinGameSheetHandler saves partial join game sheet data for multi-step form completion.
func saveJoinGameSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "saveJoinGameSheetHandler")

	// NOTE: Partial save not yet implemented; reserved for multi-step form support.
	l.Info("saveJoinGameSheetHandler called (partial save not yet implemented)")
	return server.WriteResponse(l, w, http.StatusAccepted, nil)
}

// submitJoinGameHandler processes the join game submission: creates an account if needed,
// creates a player subscription, and assigns the player to the game instance.
func submitJoinGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "submitJoinGameHandler")

	instanceID := pp.ByName("game_instance_id")
	if instanceID == "" {
		return coreerror.RequiredPathParameter("game_instance_id")
	}

	var req player_schema.JoinGameSubmitRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed to read request >%v<", err)
		return err
	}

	if req.Email == "" {
		return coreerror.NewInvalidDataError("email is required")
	}

	mm := m.(*domain.Domain)

	// Verify instance exists and is accepting players.
	instanceRec, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%s< >%v<", instanceID, err)
		return err
	}
	if instanceRec.Status != game_record.GameInstanceStatusCreated {
		return coreerror.NewInvalidDataError("this game instance is not open for enrollment")
	}

	hasCapacity, err := mm.HasAvailableCapacity(instanceID)
	if err != nil {
		l.Warn("failed to check capacity for instance >%s< >%v<", instanceID, err)
		return err
	}
	if !hasCapacity {
		return coreerror.NewInvalidDataError("this game instance is no longer accepting new players")
	}

	// Find or create an account for the provided email.
	accountRec, err := mm.GetAccountRecByEmail(req.Email)
	if err != nil {
		l.Warn("failed to find account by email >%v<", err)
		return err
	}

	var accountID string
	var accountUserID string

	if accountRec == nil {
		l.Info("no existing account for email >%s<, creating new account", req.Email)
		_, newUser, _, err := mm.CreateAccount(&account_record.AccountUser{
			Email:  req.Email,
			Status: account_record.AccountUserStatusActive,
		})
		if err != nil {
			l.Warn("failed to create account for email >%s< >%v<", req.Email, err)
			return err
		}
		accountID = newUser.AccountID
		accountUserID = newUser.ID
		l.Info("created new account >%s< for email >%s<", accountID, req.Email)
	} else {
		accountID = accountRec.AccountID
		accountUserID = accountRec.ID
		l.Info("found existing account >%s< for email >%s<", accountID, req.Email)
	}

	// Create a contact record for the player — required for player subscriptions.
	contactRec := &account_record.AccountUserContact{
		AccountUserID:      accountUserID,
		Name:               req.Name,
		PostalAddressLine1: req.PostalAddressLine1,
		StateProvince:      req.StateProvince,
		Country:            req.Country,
		PostalCode:         req.PostalCode,
	}
	if req.PostalAddressLine2 != "" {
		contactRec.PostalAddressLine2 = sql.NullString{String: req.PostalAddressLine2, Valid: true}
	}
	createdContact, err := mm.CreateAccountUserContactRec(contactRec)
	if err != nil {
		l.Warn("failed to create account user contact >%v<", err)
		return err
	}
	l.Info("created account user contact >%s< for account user >%s<", createdContact.ID, accountUserID)

	// Create a player subscription for the game associated with this instance.
	subRec := &game_record.GameSubscription{
		GameID:               instanceRec.GameID,
		AccountID:            accountID,
		AccountUserID:        sql.NullString{String: accountUserID, Valid: true},
		AccountUserContactID: sql.NullString{String: createdContact.ID, Valid: true},
		SubscriptionType:     game_record.GameSubscriptionTypePlayer,
		Status:               game_record.GameSubscriptionStatusActive,
	}

	createdSub, err := mm.CreateGameSubscriptionRec(subRec)
	if err != nil {
		l.Warn("failed to create game subscription >%v<", err)
		return err
	}
	l.Info("created game subscription >%s< for account >%s<", createdSub.ID, accountID)

	// Assign the player to the game instance.
	_, err = mm.AssignPlayerToGameInstance(createdSub.ID, instanceID)
	if err != nil {
		l.Warn("failed to assign player to game instance >%s< >%v<", instanceID, err)
		return err
	}
	l.Info("assigned subscription >%s< to game instance >%s<", createdSub.ID, instanceID)

	res := player_schema.JoinGameSubmitResponse{
		Data: &player_schema.JoinGameSubmitResponseData{
			GameSubscriptionID: createdSub.ID,
			GameInstanceID:     instanceID,
			GameID:             instanceRec.GameID,
		},
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}
