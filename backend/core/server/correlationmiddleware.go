package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const HeaderXCorrelationID = "X-Correlation-ID"

// CorrelationMiddleware -
func (rnr *Runner) CorrelationMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, _ domainer.Domainer, jc *river.Client[pgx.Tx]) error {
		l = Logger(l, "CorrelationMiddleware")

		correlationID := r.Header.Get(HeaderXCorrelationID)
		if correlationID == "" {
			correlationID = uuid.NewString()
			l.Debug("(correlationmiddleware) generated correlation ID >%s<", correlationID)
		}

		r, err := SetRequestCorrelationID(l, r, correlationID)
		if err != nil {
			l.Error("(correlationmiddleware) failed to set request correlation ID >%v<", err)
			return err
		}

		// Add correlation ID to response header
		w.Header().Set(HeaderXCorrelationID, correlationID)

		// Add correlation ID to logger context
		l.Context(log.ContextKeyCorrelationID, correlationID)

		return h(w, r, pp, nil, l, nil, jc)
	}

	return handle, nil
}

// SetRequestCorrelationID sets the correlation ID in http request context
func SetRequestCorrelationID(l logger.Logger, r *http.Request, correlationID string) (*http.Request, error) {

	l.Info("(correlationmiddleware) setting request correlation ID >%s<", correlationID)

	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxKeyCorrelationID, correlationID)
	r = r.WithContext(ctx)

	return r, nil
}

// GetRequestCorrelationID returns the correlation ID from http request context
func GetRequestCorrelationID(l logger.Logger, r *http.Request) (string, error) {
	correlationID, ok := (r.Context().Value(ctxKeyCorrelationID)).(string)
	if !ok {
		return "", fmt.Errorf("missing correlation ID")
	}

	l.Info("(correlationmiddleware) returning request correlation ID >%s<", correlationID)

	return correlationID, nil
}
