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
	tagGroupLocationLink server.TagGroup = "LocationLinks"
	TagLocationLink      server.Tag      = "LocationLinks"
)

const (
	getManyLocationLinks = "get-location-links"
	getOneLocationLink   = "get-location-link"
	createLocationLink   = "create-location-link"
	deleteLocationLink   = "delete-location-link"
)

func (rnr *Runner) locationLinkHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "locationLinkHandlerConfig")

	l.Debug("Adding location_link handler configuration")

	locationLinkConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "location_link.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "location_link.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "location_link.response.schema.json",
		},
		References: referenceSchemas,
	}

	locationLinkConfig[getManyLocationLinks] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/location-links",
		HandlerFunc: rnr.getManyLocationLinksHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get location link collection",
		},
	}

	locationLinkConfig[getOneLocationLink] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/location-links/:location_link_id",
		HandlerFunc: rnr.getLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get location link",
		},
	}

	locationLinkConfig[createLocationLink] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/location-links",
		HandlerFunc: rnr.createLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create location link",
		},
	}

	locationLinkConfig[deleteLocationLink] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/location-links/:location_link_id",
		HandlerFunc: rnr.deleteLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete location link",
		},
	}

	return locationLinkConfig, nil
}

func (rnr *Runner) getManyLocationLinksHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyLocationLinksHandler")

	l.Info("querying many location_link records with params >%#v<", qp)

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyLocationLinkRecs(opts)
	if err != nil {
		l.Warn("failed getting location_link records >%v<", err)
		return err
	}

	data := make(schema.LocationLinkCollectionResponse, 0, len(recs))
	for _, rec := range recs {
		respData, err := mapper.LocationLinkRecordToResponseData(l, rec)
		if err != nil {
			return err
		}
		data = append(data, &respData)
	}

	l.Info("responding with >%d< location_link records", len(data))

	if err = server.WriteResponse(l, w, http.StatusOK, data, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetLocationLinkHandler")

	locationLinkID := pp.ByName("location_link_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetLocationLinkRec(locationLinkID, nil)
	if err != nil {
		l.Warn("failed getting location_link record >%v<", err)
		return err
	}

	data, err := mapper.LocationLinkRecordToResponseData(l, rec)
	if err != nil {
		l.Warn("failed mapping location_link record to response data >%v<", err)
		return err
	}

	l.Info("responding with location_link record id >%s<", rec.ID)

	res := schema.LocationLinkResponse{
		LocationLinkResponseData: &data,
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) createLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateLocationLinkHandler")

	l.Info("creating location_link record with path params >%#v<", pp)

	rec, err := mapper.LocationLinkRequestToRecord(l, r, &record.LocationLink{})
	if err != nil {
		return err
	}

	mm := m.(*domain.Domain)

	rec, err = mm.CreateLocationLinkRec(rec)
	if err != nil {
		l.Warn("failed creating location_link record >%v<", err)
		return err
	}

	respData, err := mapper.LocationLinkRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.LocationLinkResponse{
		LocationLinkResponseData: &respData,
	}

	l.Info("responding with created location_link record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) deleteLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteLocationLinkHandler")

	locationLinkID := pp.ByName("location_link_id")

	l.Info("deleting location_link record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	if err := mm.DeleteLocationLinkRec(locationLinkID); err != nil {
		l.Warn("failed deleting location_link record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
