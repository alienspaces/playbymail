package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const jsonContentType = "application/json"

// DataMiddleware -
func (rnr *Runner) DataMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
		l = Logger(l, "DataMiddleware")

		if r.Method == http.MethodGet {
			return h(w, r, pp, qp, l, m)
		}

		data, err := GetRequestData(r)
		if err != nil {
			l.Warn("(core) failed reading request data >%v<", err)
			return err
		}

		// The default content type header is application/json. When processing
		// requests with a content type header other than application/json we do
		// not validate the content.
		contentType := jsonContentType
		headerContentTypes := r.Header[http.CanonicalHeaderKey("Content-Type")]
		if len(headerContentTypes) > 0 {
			l.Info("(core) assigning content type from header >%s<", headerContentTypes[0])
			contentType = headerContentTypes[0]
		}

		if contentType != jsonContentType {
			l.Debug("(core) skipping validation of URI >%s< Content-Type >%s<", r.RequestURI, contentType)
			return h(w, r, pp, qp, l, m)
		}

		requestSchema := hc.MiddlewareConfig.ValidateRequestSchema
		schemaMain := requestSchema.Main
		if schemaMain.Name == "" || schemaMain.Location == "" {
			l.Warn("(core) missing schemas, not validating data for URI >%s< method >%s<", r.RequestURI, r.Method)
			return h(w, r, pp, qp, l, m)
		}

		l.Debug("(core) schemas >%#v<", requestSchema)

		result, err := jsonschema.Validate(requestSchema, data)
		if err != nil {
			l.Warn("(core) failed validate >%v<", err)
			var jsonSyntaxError *json.SyntaxError
			if errors.As(err, &jsonSyntaxError) || errors.Is(err, io.ErrUnexpectedEOF) {
				err = coreerror.NewInvalidDataError("")
			} else if errors.Is(err, io.EOF) {
				err = coreerror.NewInvalidDataError("Request body is empty.")
			}
			return err
		}

		if !result.Valid() {
			err := coreerror.NewSchemaValidationError(result.Errors())
			l.Warn("(core) failed validating request data >%#v<", err)
			return err
		}

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}

// GetRequestData returns the request data from the http request, allowing multiple reads from the request body
func GetRequestData(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(d))

	if len(d) == 0 {
		return nil, nil
	}

	return d, nil
}
