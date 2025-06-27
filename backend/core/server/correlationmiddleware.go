package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const HeaderXCorrelationID = "X-Correlation-ID"

// CorrelationMiddleware -
func (rnr *Runner) CorrelationMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, _ domainer.Domainer) error {
		l = Logger(l, "CorrelationMiddleware")

		correlationID := r.Header.Get(HeaderXCorrelationID)
		if correlationID == "" {
			correlationID = uuid.NewString()
			l.Debug("(core) generated correlation ID >%s<", correlationID)
		}
		w.Header().Set(HeaderXCorrelationID, correlationID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyCorrelationID, correlationID)
		r = r.WithContext(ctx)

		l.Context(log.ContextKeyCorrelationID, correlationID)

		return h(w, r, pp, nil, l, nil)
	}

	return handle, nil
}
