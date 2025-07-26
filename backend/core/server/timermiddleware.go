package server

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// TimerMiddleware writes response errors returned from other middle ware or
// handler functions
func (rnr *Runner) TimerMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, _ domainer.Domainer, jc *river.Client[pgx.Tx]) error {

		startTime := time.Now()
		err := h(w, r, pp, qp, l, nil, jc)
		duration := time.Since(startTime)

		l = Logger(l, "TimerMiddleware").WithDurationContext(duration.String())
		l.Info("(core) request method >%s< path >%s<", r.Method, r.RequestURI)

		return err
	}

	return handle, nil
}
