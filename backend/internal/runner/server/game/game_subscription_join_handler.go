package game

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/record"
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
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
	"gitlab.com/alienspaces/playbymail/schema/api/player_schema"
)

const (
	GetJoinInfo  = "get-join-info"
	GetJoinSheet = "get-join-sheet"
	SubmitJoin   = "submit-join"
)

func gameSubscriptionJoinHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameSubscriptionJoinHandlerConfig")

	l.Debug("Adding game subscription join handler configuration")

	cfg := make(map[string]server.HandlerConfig)

	cfg[GetJoinInfo] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id/join",
		HandlerFunc: getJoinInfoHandler,
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
			Description: "Returns game and subscription information for displaying the join game form.",
		},
	}

	cfg[GetJoinSheet] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id/join/sheet",
		HandlerFunc: getJoinSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeOptionalToken},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Get join game turn sheet",
			Description: "Returns the join game turn sheet as HTML for online completion or printing.",
		},
	}

	cfg[SubmitJoin] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id/join",
		HandlerFunc: submitJoinHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeOptionalToken},
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
			Description: "Submits the join game form. Creates an account if needed, creates a player subscription, and assigns the player to an available game instance.",
		},
	}

	return cfg, nil
}

// resolveManagerSubscription fetches and validates that a subscription is an active manager subscription.
func resolveManagerSubscription(l logger.Logger, pp httprouter.Params, mm *domain.Domain) (*game_record.GameSubscription, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "resolveManagerSubscription")

	l.Info("resolving manager subscription with path params >%#v<", pp)

	gameSubscriptionID := pp.ByName("game_subscription_id")
	if gameSubscriptionID == "" {
		return nil, coreerror.RequiredPathParameter("game_subscription_id")
	}

	gameSubscriptionRec, err := mm.GetGameSubscriptionRec(gameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription >%s< >%v<", gameSubscriptionID, err)
		return nil, err
	}

	if gameSubscriptionRec.SubscriptionType != game_record.GameSubscriptionTypeManager {
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscription, gameSubscriptionID)
	}
	if gameSubscriptionRec.Status != game_record.GameSubscriptionStatusActive {
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscription, gameSubscriptionID)
	}

	return gameSubscriptionRec, nil
}

// aggregateSubscriptionInstances returns available instances and player counts for a subscription.
func aggregateSubscriptionInstances(l logger.Logger, mm *domain.Domain, gameSubscriptionID string) (
	instances []*game_record.GameInstance,
	playerCounts map[string]int,
	err error,
) {
	gameSubscriptionInstanceRecs, err := mm.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameSubscriptionID, Val: gameSubscriptionID},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	instances = make([]*game_record.GameInstance, 0)
	playerCounts = make(map[string]int)

	for _, gameSubscriptionInstanceRec := range gameSubscriptionInstanceRecs {
		gameInstanceRec, err := mm.GetGameInstanceRec(gameSubscriptionInstanceRec.GameInstanceID, nil)
		if err != nil {
			l.Warn("failed getting game instance >%s< >%v<", gameSubscriptionInstanceRec.GameInstanceID, err)
			continue
		}
		if gameInstanceRec.Status != game_record.GameInstanceStatusCreated {
			continue
		}
		hasCapacity, err := mm.GameInstanceHasAvailableCapacity(gameInstanceRec.ID)
		if err != nil {
			l.Warn("failed checking capacity for game instance >%s< >%v<", gameInstanceRec.ID, err)
			continue
		}
		if !hasCapacity {
			continue
		}
		count, err := mm.GetPlayerCountForGameInstance(gameInstanceRec.ID)
		if err != nil {
			l.Warn("failed getting player count for game instance >%s< >%v<", gameInstanceRec.ID, err)
			continue
		}
		instances = append(instances, gameInstanceRec)
		playerCounts[gameInstanceRec.ID] = count
	}

	return instances, playerCounts, nil
}

func getJoinInfoHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getJoinInfoHandler")

	l.Info("getting join info with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	subRec, err := resolveManagerSubscription(l, pp, mm)
	if err != nil {
		return err
	}

	gameRec, err := mm.GetGameRec(subRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game >%s< >%v<", subRec.GameID, err)
		return err
	}

	instances, playerCounts, err := aggregateSubscriptionInstances(l, mm, subRec.ID)
	if err != nil {
		l.Warn("failed aggregating instances for subscription >%s< >%v<", subRec.ID, err)
		return err
	}

	if len(instances) == 0 {
		return coreerror.NewInvalidDataError("this subscription has no game instances accepting new players")
	}

	var totalCapacity, totalPlayers int
	var deliveryPost, deliveryLocal, deliveryEmail bool
	for _, inst := range instances {
		totalCapacity += inst.RequiredPlayerCount
		totalPlayers += playerCounts[inst.ID]
		deliveryPost = deliveryPost || inst.DeliveryPhysicalPost
		deliveryLocal = deliveryLocal || inst.DeliveryPhysicalLocal
		deliveryEmail = deliveryEmail || inst.DeliveryEmail
	}

	res := player_schema.JoinGameInfoResponse{
		Data: &player_schema.JoinGameInfoResponseData{
			GameSubscriptionID:    subRec.ID,
			GameName:              gameRec.Name,
			GameDescription:       gameRec.Description,
			GameType:              gameRec.GameType,
			TurnDurationHours:     gameRec.TurnDurationHours,
			TotalCapacity:         totalCapacity,
			TotalPlayers:          totalPlayers,
			DeliveryPhysicalPost:  deliveryPost,
			DeliveryPhysicalLocal: deliveryLocal,
			DeliveryEmail:         deliveryEmail,
		},
	}

	l.Info("responding with join info for subscription >%s< game >%s< with >%d< available instances", subRec.ID, gameRec.Name, len(instances))

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func getJoinSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getJoinSheetHandler")

	l.Info("getting join sheet with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	authenData := server.GetRequestAuthenData(l, r)
	var accountEmail string
	if authenData != nil && authenData.IsAuthenticated() {
		accountEmail = authenData.AccountUser.Email
	}

	subRec, err := resolveManagerSubscription(l, pp, mm)
	if err != nil {
		return err
	}

	gameRec, err := mm.GetGameRec(subRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game for subscription >%s< >%v<", subRec.ID, err)
		return err
	}

	instances, _, err := aggregateSubscriptionInstances(l, mm, subRec.ID)
	if err != nil {
		return err
	}

	// Union delivery methods across available instances.
	var deliveryPost, deliveryLocal, deliveryEmail bool
	for _, inst := range instances {
		deliveryPost = deliveryPost || inst.DeliveryPhysicalPost
		deliveryLocal = deliveryLocal || inst.DeliveryPhysicalLocal
		deliveryEmail = deliveryEmail || inst.DeliveryEmail
	}

	cfg := mm.Config()

	joinGameSheetType := adventure_game_record.AdventureGameTurnSheetTypeJoinGame

	processor, err := turnsheet.GetDocumentProcessor(l, cfg, joinGameSheetType)
	if err != nil {
		l.Warn("failed to get join game processor >%v<", err)
		return err
	}

	turnSheetCode, err := turnsheetutil.GenerateJoinGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		l.Warn("failed to generate join game turn sheet code >%v<", err)
		return err
	}

	title := "Join Game"
	instructions := turnsheet.DefaultJoinGameInstructions()
	turnNumber := 0

	joinData := &turnsheet.JoinGameData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:              &gameRec.Name,
			GameType:              &gameRec.GameType,
			TurnSheetTitle:        &title,
			TurnSheetDescription:  &gameRec.Description,
			TurnSheetInstructions: &instructions,
			TurnSheetCode:         &turnSheetCode,
			TurnNumber:            &turnNumber,
		},
		GameDescription: gameRec.Description,
		AvailableDeliveryMethods: turnsheet.DeliveryMethods{
			Email:         deliveryEmail,
			PhysicalLocal: deliveryLocal,
			PhysicalPost:  deliveryPost,
		},
		AccountEmail: accountEmail,
	}

	backgroundImage, err := mm.GetGameTurnSheetImageDataURL(gameRec.ID, joinGameSheetType)
	if err != nil {
		l.Warn("failed to get turn sheet background image >%v<", err)
	} else if backgroundImage != "" {
		joinData.BackgroundImage = &backgroundImage
		l.Info("loaded background image for join sheet, length >%d<", len(backgroundImage))
	} else {
		l.Info("no background image found for join sheet")
	}

	sheetData, err := json.Marshal(joinData)
	if err != nil {
		l.Warn("failed to marshal join game sheet data >%v<", err)
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

func submitJoinHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "submitJoinHandler")

	l.Info("submitting join game with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameSubscriptionRec, err := resolveManagerSubscription(l, pp, mm)
	if err != nil {
		return err
	}

	var req player_schema.JoinGameSubmitRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed to read request >%v<", err)
		return err
	}

	if req.Email == "" {
		return coreerror.NewInvalidDataError("email is required")
	}

	if req.DeliveryMethod == game_record.GameSubscriptionDeliveryMethodPost {
		if req.PostalAddressLine1 == "" || req.StateProvince == "" || req.Country == "" || req.PostalCode == "" {
			return coreerror.NewInvalidDataError("postal address is required for post delivery")
		}
	}

	// Source account record IDs from authenticated session if available.
	var accountID, accountUserID, accountUserContactID string
	isAuthenticated := false

	authenData := server.GetRequestAuthenData(l, r)
	if authenData != nil && authenData.IsAuthenticated() {
		isAuthenticated = true
		accountID = authenData.AccountUser.AccountID
		accountUserID = authenData.AccountUser.ID
		accountUserContactID = authenData.AccountUser.AccountUserContactID
		// Override the submitted email with the authenticated account's email so that a
		// logged-in user cannot accidentally create a subscription under a different account.
		if req.Email != authenData.AccountUser.Email {
			l.Warn("submitted email >%s< differs from authenticated account email >%s<, using authenticated email", req.Email, authenData.AccountUser.Email)
		}
		req.Email = authenData.AccountUser.Email
	}

	accountRec := &account_record.Account{
		Record: record.Record{
			ID: accountID,
		},
		Name: req.Name,
	}
	accountUserRec := &account_record.AccountUser{
		Record: record.Record{
			ID: accountUserID,
		},
		AccountID: accountID,
		Email:     req.Email,
	}
	accountUserContactRec := &account_record.AccountUserContact{
		Record: record.Record{
			ID: accountUserContactID,
		},
		AccountUserID:      accountUserID,
		Name:               nullstring.FromString(req.Name),
		PostalAddressLine1: nullstring.FromString(req.PostalAddressLine1),
		StateProvince:      nullstring.FromString(req.StateProvince),
		Country:            nullstring.FromString(req.Country),
		PostalCode:         nullstring.FromString(req.PostalCode),
	}

	// Create or update account, account user, and account user contact.
	accountRec, accountUserRec, accountUserContactRec, _, err = mm.UpsertAccount(accountRec, accountUserRec, accountUserContactRec)
	if err != nil {
		l.Warn("failed to upsert account >%v<", err)
		return err
	}

	// Determine subscription status: authenticated users are active immediately;
	// non-authenticated users start as pending_approval and confirm via email.
	subscriptionStatus := game_record.GameSubscriptionStatusActive
	if !isAuthenticated {
		subscriptionStatus = game_record.GameSubscriptionStatusPendingApproval
	}

	// Reject duplicate joins: check if this account user already has a player subscription
	// for this game. This prevents a unique constraint violation on adventure_game_character
	// and gives a clear error when a player tries to join a second time.
	if accountUserID != "" {
		existingPlayerSubs, err := mm.GetManyGameSubscriptionRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: game_record.FieldGameSubscriptionGameID, Val: gameSubscriptionRec.GameID},
				{Col: game_record.FieldGameSubscriptionAccountUserID, Val: accountUserID},
				{Col: game_record.FieldGameSubscriptionSubscriptionType, Val: game_record.GameSubscriptionTypePlayer},
			},
		})
		if err != nil {
			l.Warn("failed to check for existing player subscription >%v<", err)
			return err
		}
		if len(existingPlayerSubs) > 0 {
			return coreerror.NewInvalidDataError("you have already joined this game")
		}
	}

	// Find an available instance before creating the subscription so we reject
	// early if no capacity is available. Uses the manager subscription's instance
	// links to locate a suitable game instance.
	gameInstanceRec, err := mm.FindAvailableGameInstance(gameSubscriptionRec.ID)
	if err != nil {
		l.Warn("failed to find available instance for subscription >%s< >%v<", gameSubscriptionRec.ID, err)
		return err
	}
	if gameInstanceRec == nil {
		return coreerror.NewInvalidDataError("no game instances are currently accepting new players")
	}

	// Non-authenticated subscriptions expire after 24 hours if not confirmed.
	var pendingApprovalExpiresAt sql.NullTime
	if !isAuthenticated {
		pendingApprovalExpiresAt = nulltime.FromTime(time.Now().Add(24 * time.Hour))
	}

	playerGameSubscriptionRec, err := mm.CreateGameSubscriptionRec(&game_record.GameSubscription{
		GameID:                   gameSubscriptionRec.GameID,
		AccountID:                accountRec.ID,
		AccountUserID:            accountUserRec.ID,
		AccountUserContactID:     nullstring.FromString(accountUserContactRec.ID),
		SubscriptionType:         game_record.GameSubscriptionTypePlayer,
		Status:                   subscriptionStatus,
		DeliveryMethod:           nullstring.FromString(req.DeliveryMethod),
		PendingApprovalExpiresAt: pendingApprovalExpiresAt,
	})
	if err != nil {
		l.Warn("failed to create player game subscription >%v<", err)
		return err
	}

	// For adventure games, create the character definition so it is linked to this player.
	// The character instance is created later by StartGameInstance when the game begins.
	gameRec, err := mm.GetGameRec(gameSubscriptionRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%v<", err)
		return err
	}

	if gameRec.GameType == game_record.GameTypeAdventure {
		characterName := req.CharacterName
		if characterName == "" {
			characterName = req.Name
		}

		characterRec := &adventure_game_record.AdventureGameCharacter{
			GameID:        gameRec.ID,
			AccountID:     accountRec.ID,
			AccountUserID: accountUserRec.ID,
			Name:          characterName,
		}
		characterRec, err = mm.CreateAdventureGameCharacterRec(characterRec)
		if err != nil {
			l.Warn("failed to create adventure game character >%v<", err)
			return err
		}

		l.Info("created adventure game character >%s< for player >%s<", characterRec.ID, accountUserRec.ID)
	}

	// Reserve the slot by linking the player subscription to the game instance
	// for both authenticated and non-authenticated users.
	_, err = mm.AssignPlayerToGameInstance(playerGameSubscriptionRec.ID, gameInstanceRec.ID)
	if err != nil {
		l.Warn("failed to assign player to game instance >%v<", err)
		return err
	}

	l.Info("assigned player subscription >%s< to game instance >%s<", playerGameSubscriptionRec.ID, gameInstanceRec.ID)

	if !isAuthenticated {
		// Non-authenticated: send confirmation email; the subscription stays
		// pending_approval until the player clicks the link.
		if _, err := jc.InsertTx(r.Context(), mm.Tx, &jobworker.SendGameSubscriptionApprovalEmailWorkerArgs{
			GameSubscriptionID: playerGameSubscriptionRec.ID,
		}, &river.InsertOpts{Queue: jobqueue.QueueDefault}); err != nil {
			l.Warn("failed to enqueue approval email job >%v<", err)
			return err
		}

		l.Info("responding with pending subscription >%s< awaiting email confirmation", playerGameSubscriptionRec.ID)

		return server.WriteResponse(l, w, http.StatusCreated, &player_schema.JoinGameSubmitResponse{
			Data: &player_schema.JoinGameSubmitResponseData{
				GameSubscriptionID: playerGameSubscriptionRec.ID,
				GameInstanceID:     gameInstanceRec.ID,
				GameID:             gameSubscriptionRec.GameID,
				Status:             game_record.GameSubscriptionStatusPendingApproval,
			},
		})
	}

	l.Info("responding with created player subscription >%s< assigned to instance >%s<", playerGameSubscriptionRec.ID, gameInstanceRec.ID)

	return server.WriteResponse(l, w, http.StatusCreated, &player_schema.JoinGameSubmitResponse{
		Data: &player_schema.JoinGameSubmitResponseData{
			GameSubscriptionID: playerGameSubscriptionRec.ID,
			GameInstanceID:     gameInstanceRec.ID,
			GameID:             gameSubscriptionRec.GameID,
			Status:             game_record.GameSubscriptionStatusActive,
		},
	})
}
