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
	"gitlab.com/alienspaces/playbymail/schema"
)

const (
	tagGroupLocation server.TagGroup = "Locations"
	TagLocation      server.Tag      = "Locations"
)

const (
	getManyLocations = "get-locations"
	getOneLocation   = "get-location"
	createLocation   = "create-location"
	updateLocation   = "update-location"
	deleteLocation   = "delete-location"
)

func (rnr *Runner) locationHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "locationHandlerConfig")

	l.Debug("Adding location handler configuration")

	locationConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "location.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "location.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "location.response.schema.json",
		},
		References: referenceSchemas,
	}

	locationConfig[getManyLocations] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/locations",
		HandlerFunc: rnr.getManyLocationsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get location collection",
		},
	}

	locationConfig[getOneLocation] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/locations/:location_id",
		HandlerFunc: rnr.getLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get location",
		},
	}

	locationConfig[createLocation] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/locations",
		HandlerFunc: rnr.createLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create location",
		},
	}

	locationConfig[updateLocation] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/locations/:location_id",
		HandlerFunc: rnr.updateLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update location",
		},
	}

	locationConfig[deleteLocation] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/locations/:location_id",
		HandlerFunc: rnr.deleteLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete location",
		},
	}

	return locationConfig, nil
}

func (rnr *Runner) getManyLocationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyLocationsHandler")

	l.Info("querying many location records with params >%#v<", qp)

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyLocationRecs(opts)
	if err != nil {
		l.Warn("failed getting location records >%v<", err)
		return err
	}

	data, err := server.Paginate(l, recs, mapper.LocationRecordToResponseData, qp.PageSize)
	if err != nil {
		return err
	}

	l.Info("responding with >%d< location records", len(data))

	res := schema.LocationCollectionResponse(data)

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetLocationHandler")

	locationID := pp.ByName("location_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetLocationRec(locationID, nil)
	if err != nil {
		l.Warn("failed getting location record >%v<", err)
		return err
	}

	data, err := mapper.LocationRecordToResponseData(l, rec)
	if err != nil {
		l.Warn("failed mapping location record to response data >%v<", err)
		return err
	}

	l.Info("responding with location record id >%s<", rec.ID)

	res := schema.LocationResponse{
		LocationResponseData: &data,
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) createLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateLocationHandler")

	l.Info("creating location record with path params >%#v<", pp)

	var req schema.LocationRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.LocationRequestToRecord(l, &req, nil)
	if err != nil {
		return err
	}

	mm := m.(*domain.Domain)

	rec, err = mm.CreateLocationRec(rec)
	if err != nil {
		l.Warn("failed creating location record >%v<", err)
		return err
	}

	respData, err := mapper.LocationRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.LocationResponse{
		LocationResponseData: &respData,
	}

	l.Info("responding with created location record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) updateLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateLocationHandler")

	locationID := pp.ByName("location_id")

	l.Info("updating location record with path params >%#v<", pp)

	var req schema.LocationRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetLocationRec(locationID, nil)
	if err != nil {
		return err
	}

	rec, err = mapper.LocationRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateLocationRec(rec)
	if err != nil {
		l.Warn("failed updating location record >%v<", err)
		return err
	}

	data, err := mapper.LocationRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.LocationResponse{
		LocationResponseData: &data,
	}

	l.Info("responding with updated location record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) deleteLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteLocationHandler")

	locationID := pp.ByName("location_id")

	l.Info("deleting location record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	if err := mm.DeleteLocationRec(locationID); err != nil {
		l.Warn("failed deleting location record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
