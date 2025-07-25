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

// HeaderXTxRollback is used to rollback API transactions during testing.
//
// The only handler tests where the API Server must commit the DB transaction
// are those that contain further DB queries, which must be able to see changes
// made by the API server tx. This is due to the current default transaction
// isolation level (read committed).
const HeaderXTxRollback = "X-Tx-Rollback"

// TxMiddleware -
func (rnr *Runner) TxMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, _ domainer.Domainer, jc *river.Client[pgx.Tx]) error {
		l = Logger(l, "TxMiddleware")

		// NOTE: The domainer is created and initialised with every request instead of
		// creating and assigning to a runner struct "Domain" property at start up.
		// This prevents directly accessing a shared property from with the handler
		// function which is running in a goroutine. Otherwise accessing the "Domain"
		// property would require locking and block simultaneous requests.

		m, err := rnr.InitDomain(l)
		if err != nil {
			l.Warn("(core) failed initialising database transaction >%v<", err)
			return err
		}

		err = h(w, r, pp, qp, l, m, jc)
		if err != nil {
			l.Warn("(core) handler error, rolling back database transaction")
			if err := m.Rollback(); err != nil {
				l.Warn("(core) failed Tx rollback >%v<", err)
				return err
			}
			return err
		}

		l.Debug("(core) committing database transaction")

		if r.Header.Get(HeaderXTxRollback) != "" {
			if err = m.Rollback(); err != nil {
				l.Warn("(core) failed Tx rollback >%v<", err)
				return err
			}
		} else if err = m.Commit(); err != nil {
			l.Warn("(core) failed Tx commit >%v<", err)
			return err
		}

		return nil
	}

	return handle, nil
}
