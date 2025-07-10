package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

const (
	tagGroupGameCharacter server.TagGroup = "GameCharacters"
	TagGameCharacter      server.Tag      = "GameCharacters"
)

const (
	getManyGameCharacters = "get-game-characters"
	getOneGameCharacter   = "get-game-character"
	createGameCharacter   = "create-game-character"
	deleteGameCharacter   = "delete-game-character"
)

func (rnr *Runner) gameCharacterHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameCharacterHandlerConfig")

	l.Debug("Adding game_character handler configuration")

	gameCharacterConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_character.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_character.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_character.response.schema.json",
		},
		References: referenceSchemas,
	}

	gameCharacterConfig[getManyGameCharacters] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-characters",
		HandlerFunc: rnr.getManyGameCharactersHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game character collection",
		},
	}

	gameCharacterConfig[getOneGameCharacter] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-characters/:game_character_id",
		HandlerFunc: rnr.getGameCharacterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game character",
		},
	}

	gameCharacterConfig[createGameCharacter] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/game-characters",
		HandlerFunc: rnr.createGameCharacterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game character",
		},
	}

	gameCharacterConfig[deleteGameCharacter] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/game-characters/:game_character_id",
		HandlerFunc: rnr.deleteGameCharacterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game character",
		},
	}

	return gameCharacterConfig, nil
}

func (rnr *Runner) getManyGameCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyGameCharactersHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyGameCharacterRecs(opts)
	if err != nil {
		l.Warn("failed getting game_character records >%v<", err)
		return err
	}

	res, err := mapper.GameCharacterRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getGameCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetGameCharacterHandler")

	gameCharacterID := pp.ByName("game_character_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetGameCharacterRec(gameCharacterID, nil)
	if err != nil {
		l.Warn("failed getting game_character record >%v<", err)
		return err
	}

	res, err := mapper.GameCharacterRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game_character record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) createGameCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateGameCharacterHandler")

	var req schema.GameCharacterRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.GameCharacterRequestToRecord(l, r, &record.GameCharacter{})
	if err != nil {
		return err
	}

	mm := m.(*domain.Domain)

	rec, err = mm.CreateGameCharacterRec(rec)
	if err != nil {
		l.Warn("failed creating game_character record >%v<", err)
		return err
	}

	res, err := mapper.GameCharacterRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) deleteGameCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteGameCharacterHandler")

	gameCharacterID := pp.ByName("game_character_id")

	mm := m.(*domain.Domain)

	if err := mm.DeleteGameCharacterRec(gameCharacterID); err != nil {
		l.Warn("failed deleting game_character record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
