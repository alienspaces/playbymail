package server

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// ErrorMiddleware writes response errors returned from other middle ware or
// handler functions
func (rnr *Runner) ErrorMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, _ domainer.Domainer, jc *river.Client[pgx.Tx]) error {
		l = Logger(l, "ErrorMiddleware")

		err := h(w, r, pp, qp, l, nil, jc)
		if err != nil {
			l.Warn("(core) error middleware >%v<", err)
			WriteError(l, w, err)
			return err
		}

		return nil
	}

	return handle, nil
}
