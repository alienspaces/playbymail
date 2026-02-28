package game

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// API Resource Search Path
//
// GET (collection) /api/v1/game-instances
//
// # API Resource CRUD Paths
//
//   - GET (collection)  /api/v1/games/{game_id}/instances
//   - GET (document)    /api/v1/games/{game_id}/instances/{instance_id}
//   - POST (document)   /api/v1/games/{game_id}/instances
//   - PUT (document)    /api/v1/games/{game_id}/instances/{instance_id}
//   - DELETE (document) /api/v1/games/{game_id}/instances/{instance_id}
//
// # Runtime Management Endpoints
//
//   - POST (document)   /api/v1/games/{game_id}/instances/{instance_id}/start
//   - POST (document)   /api/v1/games/{game_id}/instances/{instance_id}/pause
//   - POST (document)   /api/v1/games/{game_id}/instances/{instance_id}/resume
//   - POST (document)   /api/v1/games/{game_id}/instances/{instance_id}/cancel
//   - POST (document)   /api/v1/games/{game_id}/instances/{instance_id}/reset
const (
	SearchManyGameInstances = "search-many-game-instances"
	GetManyGameInstances    = "get-many-game-instances"
	GetOneGameInstance      = "get-one-game-instance"
	CreateOneGameInstance   = "create-one-game-instance"
	UpdateOneGameInstance   = "update-one-game-instance"
	DeleteOneGameInstance   = "delete-one-game-instance"
	StartGameInstance       = "start-game-instance"
	PauseGameInstance       = "pause-game-instance"
	ResumeGameInstance      = "resume-game-instance"
	CancelGameInstance      = "cancel-game-instance"
	ResetGameInstance       = "reset-game-instance"
	GetJoinGameLink         = "get-join-game-link"
	InviteTester            = "invite-tester"
)

func gameInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameInstanceHandlerConfig")

	l.Debug("Adding game instance handler configuration")

	gameInstanceConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_instance.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_instance.schema.json",
			},
		}...),
	}

	gameInstanceConfig[SearchManyGameInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-instances",
		HandlerFunc: searchManyGameInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search game instances",
		},
	}

	gameInstanceConfig[GetManyGameInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/instances",
		HandlerFunc: getManyGameInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game instances",
		},
	}

	gameInstanceConfig[GetOneGameInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/instances/:instance_id",
		HandlerFunc: getOneGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game instance",
		},
	}

	gameInstanceConfig[CreateOneGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances",
		HandlerFunc: createOneGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game instance",
		},
	}

	gameInstanceConfig[UpdateOneGameInstance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/games/:game_id/instances/:instance_id",
		HandlerFunc: updateOneGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game instance",
		},
	}

	gameInstanceConfig[DeleteOneGameInstance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/games/:game_id/instances/:instance_id",
		HandlerFunc: deleteOneGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game instance",
		},
	}

	// Runtime management endpoints
	gameInstanceConfig[StartGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/start",
		HandlerFunc: startGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Start game instance",
		},
	}

	gameInstanceConfig[PauseGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/pause",
		HandlerFunc: pauseGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Pause game instance",
		},
	}

	gameInstanceConfig[ResumeGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/resume",
		HandlerFunc: resumeGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Resume game instance",
		},
	}

	gameInstanceConfig[CancelGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/cancel",
		HandlerFunc: cancelGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Cancel game instance",
		},
	}

	gameInstanceConfig[ResetGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/reset",
		HandlerFunc: resetOneGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Reset game instance",
			Description: "Reset a game instance to its initial state. Soft-deletes all instance-level data " +
				"(turn sheets, characters, items, locations, etc.) and resets the instance to status=created " +
				"with turn 0. Subscription links are preserved so players remain joined.",
		},
	}

	joinGameLinkResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "join_game_link.response.schema.json",
		},
		References: referenceSchemas,
	}

	gameInstanceConfig[GetJoinGameLink] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/join-link",
		HandlerFunc: getJoinGameLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateResponseSchema: joinGameLinkResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get join game link",
			Description: "Get the join game link for a closed testing game instance. " +
				"Returns the URL that can be shared with testers to join the game.",
		},
	}

	inviteTesterRequestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "invite_tester.request.schema.json",
		},
		References: referenceSchemas,
	}

	inviteTesterResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "invite_tester.response.schema.json",
		},
		References: referenceSchemas,
	}

	gameInstanceConfig[InviteTester] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/invite-tester",
		HandlerFunc: inviteTesterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateRequestSchema:  inviteTesterRequestSchema,
			ValidateResponseSchema: inviteTesterResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Invite tester",
			Description: "Send an email invitation to a tester for a closed testing game instance. " +
				"The email will include a join game link.",
		},
	}

	return gameInstanceConfig, nil
}

func searchManyGameInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyGameInstancesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GameInstanceRepository().GetMany(opts)
	if err != nil {
		l.Warn("failed getting game instance records >%v<", err)
		return err
	}

	getPlayerCount := func(instanceID string) (int, error) {
		return mm.GetPlayerCountForGameInstance(instanceID)
	}

	res, err := mapper.GameInstanceRecordsToCollectionResponse(l, recs, getPlayerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyGameInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameInstancesHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.RequiredPathParameter("game_id")
	}

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{
		Col: game_record.FieldGameInstanceGameID,
		Val: gameID,
	})

	recs, err := mm.GameInstanceRepository().GetMany(opts)
	if err != nil {
		l.Warn("failed getting game instance records >%v<", err)
		return err
	}

	getPlayerCount := func(instanceID string) (int, error) {
		return mm.GetPlayerCountForGameInstance(instanceID)
	}

	res, err := mapper.GameInstanceRecordsToCollectionResponse(l, recs, getPlayerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneGameInstanceHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.RequiredPathParameter("game_id")
	}
	if instanceID == "" {
		l.Warn("instance id is required")
		return coreerror.RequiredPathParameter("instance_id")
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed getting game instance record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		l.Warn("instance does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	playerCount, err := mm.GetPlayerCountForGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to get player count for game instance >%s< >%v<", instanceID, err)
		// Continue with 0 if we can't get the count
		playerCount = 0
	}

	res, err := mapper.GameInstanceRecordToResponse(l, rec, playerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneGameInstanceHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.RequiredPathParameter("game_id")
	}

	// Get authenticated account ID
	authenData := server.GetRequestAuthenData(l, r)
	if authenData == nil || authenData.AccountUser.ID == "" {
		l.Warn("authenticated account is required to create game instance")
		return coreerror.NewUnauthorizedError()
	}

	mm := m.(*domain.Domain)

	// Verify account has manager subscription for this game (required for creating instances)
	_, err := mm.GetGameSubscriptionRecByAccountAndGame(
		authenData.AccountUser.AccountID,
		gameID,
		game_record.GameSubscriptionTypeManager,
	)
	if err != nil {
		l.Warn("failed to find manager subscription for account >%s< and game >%s<: %v",
			authenData.AccountUser.AccountID, gameID, err)
		return coreerror.NewUnauthorizedError()
	}

	// Create record with GameID
	rec := &game_record.GameInstance{
		GameID: gameID,
	}

	// Map request data to record
	rec, err = mapper.GameInstanceRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping request to record >%v<", err)
		return coreerror.NewInvalidDataError("invalid request data")
	}

	l.Debug("mapped request to record: delivery_physical_post=%v, delivery_physical_local=%v, delivery_email=%v",
		rec.DeliveryPhysicalPost, rec.DeliveryPhysicalLocal, rec.DeliveryEmail)

	rec, err = mm.CreateGameInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game instance record >%v<", err)
		return err
	}

	playerCount, err := mm.GetPlayerCountForGameInstance(rec.ID)
	if err != nil {
		l.Warn("failed to get player count for game instance >%s< >%v<", rec.ID, err)
		playerCount = 0
	}

	res, err := mapper.GameInstanceRecordToResponse(l, rec, playerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneGameInstanceHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		return err
	}

	// Read and parse request body and apply updates using mapper
	updatedRec, err := mapper.GameInstanceRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping request to record >%v<", err)
		return coreerror.NewInvalidDataError("invalid request data")
	}

	// Update the record
	rec, err = mm.UpdateGameInstanceRec(updatedRec)
	if err != nil {
		l.Warn("failed updating game instance record >%v<", err)
		return err
	}

	playerCount, err := mm.GetPlayerCountForGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to get player count for game instance >%s< >%v<", instanceID, err)
		playerCount = 0
	}

	res, err := mapper.GameInstanceRecordToResponse(l, rec, playerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneGameInstanceHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("instance does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	if err := mm.DeleteGameInstanceRec(instanceID); err != nil {
		l.Warn("failed deleting game instance record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// startGameInstanceHandler starts a game instance
func startGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "startGameInstanceHandler")

	l.Info("starting game instance")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Start the game instance
	instance, err := mm.StartGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to start game instance >%v<", err)
		return err
	}

	playerCount, err := mm.GetPlayerCountForGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to get player count for game instance >%s< >%v<", instanceID, err)
		playerCount = 0
	}

	// Convert to response
	res, err := mapper.GameInstanceRecordToResponse(l, instance, playerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// pauseGameInstanceHandler pauses a game instance
func pauseGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "pauseGameInstanceHandler")

	l.Info("pausing game instance")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Pause the game instance
	instance, err := mm.PauseGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to pause game instance >%v<", err)
		return err
	}

	playerCount, err := mm.GetPlayerCountForGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to get player count for game instance >%s< >%v<", instanceID, err)
		playerCount = 0
	}

	// Convert to response
	res, err := mapper.GameInstanceRecordToResponse(l, instance, playerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// resumeGameInstanceHandler resumes a game instance
func resumeGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "resumeGameInstanceHandler")

	l.Info("resuming game instance")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Resume the game instance
	instance, err := mm.ResumeGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to resume game instance >%v<", err)
		return err
	}

	playerCount, err := mm.GetPlayerCountForGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to get player count for game instance >%s< >%v<", instanceID, err)
		playerCount = 0
	}

	// Convert to response
	res, err := mapper.GameInstanceRecordToResponse(l, instance, playerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// cancelGameInstanceHandler cancels a game instance
func cancelGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "cancelGameInstanceHandler")

	l.Info("canceling game instance")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Cancel the game instance
	instance, err := mm.CancelGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to cancel game instance >%v<", err)
		return err
	}

	playerCount, err := mm.GetPlayerCountForGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to get player count for game instance >%s< >%v<", instanceID, err)
		playerCount = 0
	}

	// Convert to response
	res, err := mapper.GameInstanceRecordToResponse(l, instance, playerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// getJoinGameLinkHandler returns the join game link for a closed testing game instance
func getJoinGameLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getJoinGameLinkHandler")

	l.Info("getting join game link")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Get the game instance
	instance, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%v<", err)
		return err
	}

	// Validate it's in closed testing mode
	if !instance.IsClosedTesting {
		return coreerror.NewInvalidDataError("game instance is not in closed testing mode")
	}

	// Generate join game key if it doesn't exist
	if !instance.ClosedTestingJoinGameKey.Valid || instance.ClosedTestingJoinGameKey.String == "" {
		_, err = mm.GenerateClosedTestingJoinGameKey(instanceID)
		if err != nil {
			l.Warn("failed to generate join game key >%v<", err)
			return err
		}
		// Re-fetch to get the new key
		instance, err = mm.GetGameInstanceRec(instanceID, nil)
		if err != nil {
			l.Warn("failed to get game instance after key generation >%v<", err)
			return err
		}
	}

	// Build join game URL
	joinURL := fmt.Sprintf("/player/join-game/%s", instance.ClosedTestingJoinGameKey.String)

	// Map to response
	res, err := mapper.JoinGameLinkToResponse(l, joinURL, instance.ClosedTestingJoinGameKey.String)
	if err != nil {
		l.Warn("failed mapping join game link to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// inviteTesterHandler sends an email invitation to a tester
func inviteTesterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "inviteTesterHandler")

	l.Info("inviting tester")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Get the game instance
	instance, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%v<", err)
		return err
	}

	// Validate it's in closed testing mode
	if !instance.IsClosedTesting {
		return coreerror.NewInvalidDataError("game instance is not in closed testing mode")
	}

	// Read request body for email
	email, err := mapper.InviteTesterRequestToEmail(l, r)
	if err != nil {
		l.Warn("failed to read invite tester request >%v<", err)
		return coreerror.NewInvalidDataError("invalid request: %s", err.Error())
	}

	// Generate join game key if it doesn't exist
	if !instance.ClosedTestingJoinGameKey.Valid || instance.ClosedTestingJoinGameKey.String == "" {
		_, err = mm.GenerateClosedTestingJoinGameKey(instanceID)
		if err != nil {
			l.Warn("failed to generate join game key >%v<", err)
			return err
		}
	}

	// Queue email job
	l.Info("queuing tester invitation email for >%s< game instance >%s<", email, instanceID)

	_, err = jc.InsertTx(context.Background(), mm.Tx, &jobworker.SendTesterInvitationEmailWorkerArgs{
		GameInstanceID: instanceID,
		Email:          email,
	}, &river.InsertOpts{Queue: jobqueue.QueueDefault})
	if err != nil {
		l.Warn("failed to enqueue tester invitation email job >%v<", err)
		return coreerror.NewInternalError("failed to queue tester invitation email: %v", err)
	}

	// Map to response
	res, err := mapper.InviteTesterToResponse(l, email)
	if err != nil {
		l.Warn("failed mapping invite tester to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func resetOneGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "resetOneGameInstanceHandler")

	l.Info("resetting game instance")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	instance, err := mm.ResetGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to reset game instance >%v<", err)
		return err
	}

	playerCount, err := mm.GetPlayerCountForGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to get player count for game instance >%s< >%v<", instanceID, err)
		playerCount = 0
	}

	res, err := mapper.GameInstanceRecordToResponse(l, instance, playerCount)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
