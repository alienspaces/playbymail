package game

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	ValidateGame = "validate-game"
)

type gameValidationResponseData struct {
	Valid  bool                          `json:"valid"`
	Issues []domain.GameValidationIssue  `json:"issues"`
}

type gameValidationResponse struct {
	Data gameValidationResponseData `json:"data"`
}

func gameValidateHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameValidateHandlerConfig")

	l.Debug("Adding game validate handler configuration")

	config := make(map[string]server.HandlerConfig)

	config[ValidateGame] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/validate",
		HandlerFunc: validateGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameDesign,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Validate game",
			Description: "Validate whether a game is ready to create instances. Returns all issues found.",
		},
	}

	return config, nil
}

func validateGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "validateGameHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	mm := m.(*domain.Domain)

	issues, err := mm.ValidateGameReadyForInstance(gameID)
	if err != nil {
		l.Warn("failed validating game >%s< >%v<", gameID, err)
		return err
	}

	if issues == nil {
		issues = []domain.GameValidationIssue{}
	}

	valid := true
	for _, issue := range issues {
		if issue.Severity == domain.ValidationSeverityError {
			valid = false
			break
		}
	}

	res := gameValidationResponse{
		Data: gameValidationResponseData{
			Valid:  valid,
			Issues: issues,
		},
	}

	if err := server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
