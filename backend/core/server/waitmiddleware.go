package server

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// HeaderXWaitSeconds causes the request to take at least some amount of time to
// complete. The value cannot exceed the WriteTimeout set on the server. This
// can be used to:
//
// 1) simulate server load; or
//
// 2) test locking behaviour.
const HeaderXWaitSeconds = "X-Wait-Seconds"

// HeaderXTxLockWaitTimeoutSeconds is used to specify the lock wait timeout for
// database locks. The value cannot exceed the WriteTimeout set on the server.
// This can be used to test locking behaviour by first making an API request
// that explicitly or implicitly acquires a lock (i.e., UPDATE, DELETE) with
// HeaderXWaitSeconds with a value greater than HeaderXTxLockWaitTimeoutSeconds,
// and, in that intervening period, make another request that explicitly or
// implicitly attempts to acquire a lock on the same rows.
const HeaderXTxLockWaitTimeoutSeconds = "X-Tx-Lock-Wait-Timeout-Seconds"

func (rnr *Runner) WaitMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
		l = Logger(l, "WaitMiddleware")

		if r.Header.Get(HeaderXTxLockWaitTimeoutSeconds) != "" {
			timeoutSecs, err := strconv.ParseFloat(r.Header.Get(HeaderXTxLockWaitTimeoutSeconds), 64)
			if err != nil {
				l.Warn("(core) header %s is not an int or float >%v<", HeaderXTxLockWaitTimeoutSeconds, err)
				err = coreerror.NewHeaderError("%s header value must be an int or float greater than 0 >%f<", HeaderXTxLockWaitTimeoutSeconds, timeoutSecs)
				return err
			}

			if timeoutSecs <= 0 {
				err = coreerror.NewHeaderError("%s header value must be an int or float greater than 0 >%f<", HeaderXTxLockWaitTimeoutSeconds, timeoutSecs)
				return err
			}

			if err = m.SetTxLockTimeout(timeoutSecs); err != nil {
				err = coreerror.NewInternalError("failed to set transaction lock timeout >%s<", err.Error())
				l.Warn(err.Error())
				return err
			}
		}

		err := h(w, r, pp, qp, l, m, jc)

		if r.Header.Get(HeaderXWaitSeconds) != "" {
			waitSecs, err := strconv.ParseFloat(r.Header.Get(HeaderXWaitSeconds), 32)
			if err != nil {
				l.Warn("(core) header %s is not an int >%v<", HeaderXWaitSeconds, err)
				err = coreerror.NewHeaderError("%s header value must be an int or float greater than 0 >%f<", HeaderXWaitSeconds, waitSecs)
				return err
			}

			if waitSecs <= 0 {
				err = coreerror.NewHeaderError("%s header value must be an int or float greater than 0 >%f<", HeaderXWaitSeconds, waitSecs)
				return err
			}

			l.Debug("(core) waiting >%fs<", waitSecs)

			waitMs := math.Round(waitSecs * 1000)
			time.Sleep(time.Duration(waitMs) * time.Millisecond)
		}

		return err
	}

	return handle, nil
}
