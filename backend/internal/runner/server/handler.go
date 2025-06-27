package runner

import "gitlab.com/alienspaces/playbymail/core/jsonschema"

// Reference schemas commonly used by all param, request and response schemas
var referenceSchemas = []jsonschema.Schema{
	{
		Location: "schema",
		Name:     "query.schema.json",
	},
	{
		Location: "schema",
		Name:     "common.schema.json",
	},
	{
		Location: "schema",
		Name:     "game.schema.json",
	},
}
