package server

import (
	"context"
	"net/http"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

type RLS struct {
	Identifiers map[string][]string
}

func (rnr *Runner) RLSMiddleware(hc HandlerConfig, h Handle) (Handle, error) {
	handlerAuthenTypes := set.New(hc.MiddlewareConfig.AuthenTypes...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
		l = Logger(l, "RLSMiddleware")

		if _, ok := handlerAuthenTypes[AuthenticationTypePublic]; ok {
			l.Info("(rlsmiddleware) handler name >%s< is public, not applying RLS", hc.Name)
			return h(w, r, pp, qp, l, m, jc)
		}

		AuthenData := GetRequestAuthenData(l, r)
		if AuthenData == nil {
			err := coreerror.NewInternalError("failed to read request authen data")
			return err
		}

		if AuthenData.RLSType == RLSTypeOpen {
			l.Info("(rlsmiddleware) handler name >%s< is open, not applying RLS", hc.Name)
			return h(w, r, pp, qp, l, m, jc)
		}

		if AuthenData.RLSType == RLSTypeRestricted {
			if rnr.RLSFunc == nil {
				l.Warn("(rlsmiddleware) failed to set RLS >%v<", "RLSFunc is nil")
				return coreerror.NewInternalError("failed to set RLS >%v<", "RLSFunc is nil")
			}

			rls, err := rnr.RLSFunc(l, m, *AuthenData)
			if err != nil {
				l.Warn("(rlsmiddleware) failed to set RLS >%v<", err)
				return err
			}

			r, err = SetRequestRLSData(l, r, rls)
			if err != nil {
				l.Warn("(rlsmiddleware) failed to set RLS >%v<", err)
				return err
			}

			// Check the path parameters against the RLS identifiers to ensure the user
			// has access to the resource.
			//
			// For example, if the RLS identifiers are `["account_id", "game_id"]` and
			// the path parameters are `["123", "456"]`, then we can check if the user
			// has access to the resource.

			// Validate path parameters against RLS identifiers
			for paramName, allowedValues := range rls.Identifiers {
				paramValue := pp.ByName(paramName)
				if paramValue != "" && !slices.Contains(allowedValues, paramValue) {
					l.Warn("(rlsmiddleware) access denied: path parameter %s=%s not in allowed values %v", paramName, paramValue, allowedValues)
					return coreerror.NewUnauthorizedError()
				}
			}
		}

		return h(w, r, pp, qp, l, m, jc)
	}
	return handle, nil
}

// SetRequestRLSData sets RLS data in http request context
func SetRequestRLSData(l logger.Logger, r *http.Request, rls RLS) (*http.Request, error) {
	ctx := r.Context()

	l.Info("(rlsmiddleware) setting request RLS data >%#v<", rls)

	ctx = context.WithValue(ctx, ctxKeyRLS, rls)
	r = r.WithContext(ctx)

	return r, nil
}

// GetRequestRLSData returns RLS data from http request context
func GetRequestRLSData(l logger.Logger, r *http.Request) *RLS {
	rls, ok := (r.Context().Value(ctxKeyRLS)).(RLS)
	if !ok {
		return nil
	}

	l.Info("(rlsmiddleware) returning request RLS data >%#v<", rls)

	return &rls
}

// GetRequestRLSIdentifierSet returns RLS data as sets from http request context
func GetRequestRLSIdentifierSet(l logger.Logger, r *http.Request) map[string]set.Set[string] {
	rls := GetRequestRLSData(l, r)
	if rls == nil {
		return nil
	}

	identifiers := make(map[string]set.Set[string])
	for k, v := range rls.Identifiers {
		identifiers[k] = set.New(v...)
	}

	return identifiers
}
