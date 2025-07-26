package runner

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// API Resource Paths
//
// All API resources will have a top level path for collection GET requests. These
// paths must generally acept all properties of the resources as query parameters
// and act as "search" API endpoints.
//
// For basic CRUD operations, the resource will be located beneath a hierarchical
// path based on its relationships to parent resources.
//
// For example, a resource called "game" would have the following paths:
//
// GET (collection) /api/v1/games
//
// GET (document) /api/v1/games/{game_id}
// POST (document) /api/v1/games
// PUT (document) /api/v1/games/{game_id}
// DELETE (document) /api/v1/games/{game_id}
//
// And then, a resource "character" that belongs to a "game" would have the following
// path:
//
// GET (collection) /api/v1/game-characters
//
// With its related CRUD operations nested beneath the parent "game" resource:
//
// GET (collection) /api/v1/games/{game_id}/characters
// GET (document) /api/v1/games/{game_id}/characters/{character_id}
// POST (document) /api/v1/games/{game_id}/characters
// PUT (document) /api/v1/games/{game_id}/characters/{character_id}
// DELETE (document) /api/v1/games/{game_id}/characters/{character_id}

// Common reference schemas used by all param, request and response schemas. These
// schemas must be loaded for all schemas to be validated. Schemas are loaded from
// the `./schema` directory.
var referenceSchemas = []jsonschema.Schema{
	{
		Location: "schema",
		Name:     "query.schema.json",
	},
	{
		Location: "schema",
		Name:     "common.schema.json",
	},
}

// handlerFunc - default handler
func (rnr *Runner) handlerFunc(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = loggerWithFunctionContext(l, "Handler")

	l.Info("(playbymail) using playbymail handler")

	fmt.Fprint(w, "Hello from playbymail!\n")

	return nil
}
