package server

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// TimerMiddleware writes response errors returned from other middle ware or
// handler functions
func (rnr *Runner) TimerMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, _ domainer.Domainer) error {

		startTime := time.Now()
		err := h(w, r, pp, qp, l, nil)
		duration := time.Since(startTime)

		l = Logger(l, "TimerMiddleware").WithDurationContext(duration.String())
		l.Info("(core) request method >%s< path >%s<", r.Method, r.RequestURI)

		return err
	}

	return handle, nil
}
