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

// JobClientMiddleware creates a new job client instance for each request and passes it to the handler
func (rnr *Runner) JobClientMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, _ *river.Client[pgx.Tx]) error {
		l = Logger(l, "JobClientMiddleware")

		// Create a new job client instance for this request
		jc, err := rnr.InitJobClient(l)
		if err != nil {
			l.Warn("(core) failed initialising job client >%v<", err)
			// If job client initialization fails, pass nil instead of failing the request
			// This allows tests to run without a proper job client setup
			jc = nil
		}

		// Pass the job client to the handler
		return h(w, r, pp, qp, l, m, jc)
	}

	return handle, nil
}
