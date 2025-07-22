package adventure_game

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/schema"
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-characters

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/characters
// GET (document)    /api/v1/adventure-games/{game_id}/characters/{character_id}
// POST (document)   /api/v1/adventure-games/{game_id}/characters
// PUT (document)    /api/v1/adventure-games/{game_id}/characters/{character_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/characters/{character_id}

const (
	// API Resource Search Path
	searchManyAdventureGameCharacters = "search-many-adventure-game-characters"

	// API Resource CRUD Paths
	getManyAdventureGameCharacters  = "get-many-adventure-game-characters"
	getOneAdventureGameCharacter    = "get-one-adventure-game-character"
	createOneAdventureGameCharacter = "create-one-adventure-game-character"
	updateOneAdventureGameCharacter = "update-one-adventure-game-character"
	deleteOneAdventureGameCharacter = "delete-one-adventure-game-character"
)

func adventureGameCharacterHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameCharacterHandlerConfig")

	l.Debug("Adding adventure_game_character handler configuration")

	gameCharacterConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_character.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_character.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_character.response.schema.json",
		},
		References: referenceSchemas,
	}

	// New Adventure Game Character API paths
	gameCharacterConfig[searchManyAdventureGameCharacters] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-characters",
		HandlerFunc: searchManyAdventureGameCharactersHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game characters",
		},
	}

	gameCharacterConfig[getManyAdventureGameCharacters] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/characters",
		HandlerFunc: getManyAdventureGameCharactersHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game characters",
		},
	}

	gameCharacterConfig[getOneAdventureGameCharacter] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/characters/:character_id",
		HandlerFunc: getOneAdventureGameCharacterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game character",
		},
	}

	gameCharacterConfig[createOneAdventureGameCharacter] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/characters",
		HandlerFunc: createOneAdventureGameCharacterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game character",
		},
	}

	gameCharacterConfig[updateOneAdventureGameCharacter] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/characters/:character_id",
		HandlerFunc: updateOneAdventureGameCharacterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game character",
		},
	}

	gameCharacterConfig[deleteOneAdventureGameCharacter] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/characters/:character_id",
		HandlerFunc: deleteOneAdventureGameCharacterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game character",
		},
	}

	return gameCharacterConfig, nil
}

func searchManyAdventureGameCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameCharactersHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameCharacterRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game character records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCharacterRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameCharactersHandler")

	// Create SQL options from query parameters
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts.Params = append(opts.Params, sql.Param{
		Col: "game_id",
		Val: gameID,
	})

	mm := m.(*domain.Domain)

	recs, err := mm.GetManyAdventureGameCharacterRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game character records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCharacterRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameCharacterHandler")

	gameID := pp.ByName("game_id")
	characterID := pp.ByName("character_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameCharacterRec(characterID, nil)
	if err != nil {
		l.Warn("failed getting adventure game character record >%v<", err)
		return err
	}

	// Verify the character belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("character does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("character", characterID)
	}

	res, err := mapper.AdventureGameCharacterRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game character record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameCharacterHandler")

	gameID := pp.ByName("game_id")

	var req schema.AdventureGameCharacterRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.AdventureGameCharacterRequestToRecord(l, r, &record.AdventureGameCharacter{})
	if err != nil {
		return err
	}

	// Set the game ID from the path parameter
	rec.GameID = gameID

	mm := m.(*domain.Domain)

	rec, err = mm.CreateAdventureGameCharacterRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game character record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCharacterRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameCharacterHandler")

	gameID := pp.ByName("game_id")
	characterID := pp.ByName("character_id")

	l.Info("updating adventure game character record with path params >%#v<", pp)

	var req schema.AdventureGameCharacterRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameCharacterRec(characterID, nil)
	if err != nil {
		return err
	}

	// Verify the character belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("character does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("character", characterID)
	}

	rec, err = mapper.AdventureGameCharacterRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameCharacterRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game character record >%v<", err)
		return err
	}

	data, err := mapper.AdventureGameCharacterRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.AdventureGameCharacterResponse{
		Data: &data,
	}

	l.Info("responding with updated adventure game character record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameCharacterHandler")

	gameID := pp.ByName("game_id")
	characterID := pp.ByName("character_id")

	l.Info("deleting adventure game character record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	// First get the record to verify it belongs to the specified game
	rec, err := mm.GetAdventureGameCharacterRec(characterID, nil)
	if err != nil {
		return err
	}

	// Verify the character belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("character does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("character", characterID)
	}

	if err := mm.DeleteAdventureGameCharacterRec(characterID); err != nil {
		l.Warn("failed deleting adventure game character record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
