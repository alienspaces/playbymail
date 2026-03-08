package game

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
	"gitlab.com/alienspaces/playbymail/schema/api/player_schema"
)

const (
	GetJoinInfo         = "get-join-info"
	VerifyJoinEmail     = "verify-join-email"
	GetJoinSheet        = "get-join-sheet"
	SaveJoinSheet       = "save-join-sheet"
	SubmitJoin          = "submit-join"
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

	cfg[VerifyJoinEmail] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id/join/verify-email",
		HandlerFunc: verifyJoinEmailHandler,
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

	cfg[SaveJoinSheet] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id/join/sheet",
		HandlerFunc: saveJoinSheetHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypePublic},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Save join game sheet",
			Description: "Saves partial join game sheet data for multi-step form completion.",
		},
	}

	cfg[SubmitJoin] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id/join",
		HandlerFunc: submitJoinHandler,
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
			Description: "Submits the join game form. Creates an account if needed, creates a player subscription, and assigns the player to an available game instance.",
		},
	}

	return cfg, nil
}

// resolveManagerSubscription fetches and validates that a subscription is an active manager subscription.
func resolveManagerSubscription(l logger.Logger, pp httprouter.Params, mm *domain.Domain) (*game_record.GameSubscription, error) {
	gameSubscriptionID := pp.ByName("game_subscription_id")
	if gameSubscriptionID == "" {
		return nil, coreerror.RequiredPathParameter("game_subscription_id")
	}

	subRec, err := mm.GetGameSubscriptionRec(gameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get subscription >%s< >%v<", gameSubscriptionID, err)
		return nil, err
	}

	if subRec.SubscriptionType != game_record.GameSubscriptionTypeManager {
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscription, gameSubscriptionID)
	}
	if subRec.Status != game_record.GameSubscriptionStatusActive {
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscription, gameSubscriptionID)
	}

	return subRec, nil
}

// aggregateSubscriptionInstances returns available instances and player counts for a subscription.
func aggregateSubscriptionInstances(l logger.Logger, mm *domain.Domain, subID string) (
	instances []*game_record.GameInstance,
	playerCounts map[string]int,
	err error,
) {
	gsiRecs, err := mm.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameSubscriptionID, Val: subID},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	instances = make([]*game_record.GameInstance, 0)
	playerCounts = make(map[string]int)

	for _, gsi := range gsiRecs {
		instRec, err := mm.GetGameInstanceRec(gsi.GameInstanceID, nil)
		if err != nil {
			l.Warn("failed getting instance >%s< >%v<", gsi.GameInstanceID, err)
			continue
		}
		if instRec.Status != game_record.GameInstanceStatusCreated {
			continue
		}
		hasCapacity, err := mm.HasAvailableCapacity(instRec.ID)
		if err != nil {
			l.Warn("failed checking capacity for instance >%s< >%v<", instRec.ID, err)
			continue
		}
		if !hasCapacity {
			continue
		}
		count, err := mm.GetPlayerCountForGameInstance(instRec.ID)
		if err != nil {
			l.Warn("failed getting player count for instance >%s< >%v<", instRec.ID, err)
			continue
		}
		instances = append(instances, instRec)
		playerCounts[instRec.ID] = count
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

func verifyJoinEmailHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "verifyJoinEmailHandler")

	l.Info("verifying join email with path params >%#v<", pp)

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

	l.Info("responding with join email verification, has_account >%v<", accountRec != nil)

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

func saveJoinSheetHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "saveJoinSheetHandler")

	l.Info("saveJoinSheetHandler called (partial save not yet implemented)")
	return server.WriteResponse(l, w, http.StatusAccepted, nil)
}

func submitJoinHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "submitJoinHandler")

	l.Info("submitting join game with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	subRec, err := resolveManagerSubscription(l, pp, mm)
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

	// Resolve the account: use session if authenticated, otherwise find-or-create by email
	var accountID, accountUserID string

	authenData := server.GetRequestAuthenData(l, r)
	if authenData != nil && authenData.IsAuthenticated() {
		accountID = authenData.AccountUser.AccountID
		accountUserID = authenData.AccountUser.ID
		l.Info("authenticated join game submit for account >%s< user >%s<", accountID, accountUserID)
	} else {
		accountUserRec, err := mm.GetAccountRecByEmail(req.Email)
		if err != nil {
			l.Warn("failed to look up account by email >%v<", err)
			return err
		}
		if accountUserRec != nil {
			accountID = accountUserRec.AccountID
			accountUserID = accountUserRec.ID
			l.Info("found existing account for email >%s< account >%s< user >%s<", req.Email, accountID, accountUserID)
		} else {
			newAccountUserRec := &account_record.AccountUser{
				Email: req.Email,
			}
			accountRec, createdUserRec, _, err := mm.CreateAccount(newAccountUserRec)
			if err != nil {
				l.Warn("failed to create account for email >%s< >%v<", req.Email, err)
				return err
			}
			accountID = accountRec.ID
			accountUserID = createdUserRec.ID
			l.Info("created new account for email >%s< account >%s< user >%s<", req.Email, accountID, accountUserID)
		}
	}

	// Auto-assign: find an available instance under this subscription.
	instanceRec, err := mm.FindAvailableGameInstance(subRec.ID)
	if err != nil {
		l.Warn("failed to find available instance for subscription >%s< >%v<", subRec.ID, err)
		return err
	}
	if instanceRec == nil {
		return coreerror.NewInvalidDataError("no game instances are currently accepting new players")
	}

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

	playerSubRec := &game_record.GameSubscription{
		GameID:               subRec.GameID,
		AccountID:            accountID,
		AccountUserID:        sql.NullString{String: accountUserID, Valid: true},
		AccountUserContactID: sql.NullString{String: createdContact.ID, Valid: true},
		SubscriptionType:     game_record.GameSubscriptionTypePlayer,
		Status:               game_record.GameSubscriptionStatusActive,
	}

	createdSub, err := mm.CreateGameSubscriptionRec(playerSubRec)
	if err != nil {
		l.Warn("failed to create game subscription >%v<", err)
		return err
	}

	_, err = mm.AssignPlayerToGameInstance(createdSub.ID, instanceRec.ID)
	if err != nil {
		l.Warn("failed to assign player to game instance >%s< >%v<", instanceRec.ID, err)
		return err
	}

	res := player_schema.JoinGameSubmitResponse{
		Data: &player_schema.JoinGameSubmitResponseData{
			GameSubscriptionID: createdSub.ID,
			GameInstanceID:     instanceRec.ID,
			GameID:             subRec.GameID,
		},
	}

	l.Info("responding with created player subscription >%s< assigned to instance >%s<", createdSub.ID, instanceRec.ID)

	return server.WriteResponse(l, w, http.StatusCreated, res)
}
