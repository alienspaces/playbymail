package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

type RLS struct {
	Identifiers map[string][]string
}

func (rnr *Runner) RLSMiddleware(hc HandlerConfig, h Handle) (Handle, error) {
	handlerAuthenTypes := set.New(hc.MiddlewareConfig.AuthenTypes...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
		l = Logger(l, "RLSMiddleware")

		if _, ok := handlerAuthenTypes[AuthenticationTypePublic]; ok {
			l.Debug("(core) handler name >%s< is public, not apply RLS", hc.Name)
			return h(w, r, pp, qp, l, m)
		}

		authenticatedRequest := RequestAuthData(l, r)
		if authenticatedRequest == nil {
			err := fmt.Errorf("failed to read auth data")
			return err
		}

		if authenticatedRequest.RLSType == RLSTypeOpen {
			return h(w, r, pp, qp, l, m)
		}

		if authenticatedRequest.RLSType == RLSTypeRestricted {
			if rnr.SetRLSFunc == nil {
				l.Warn("(core) failed to set RLS >%v<", "SetRLSFunc is nil")
				return fmt.Errorf("failed to set RLS >%v<", "SetRLSFunc is nil")
			}

			rls, err := rnr.SetRLSFunc(l, m, *authenticatedRequest)
			if err != nil {
				l.Warn("(core) failed to set RLS >%v<", err)
				return err
			}

			// Set RLS context
			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxKeyRLS, rls)
			r = r.WithContext(ctx)
		}

		return h(w, r, pp, qp, l, m)
	}
	return handle, nil
}
