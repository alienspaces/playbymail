package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// AuthzMiddleware -
func (rnr *Runner) AuthzMiddleware(hc HandlerConfig, h Handle) (Handle, error) {
	authenTypes := set.New(hc.MiddlewareConfig.AuthenTypes...)
	authzPermissions := set.New(hc.MiddlewareConfig.AuthzPermissions...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
		l = Logger(l, "AuthzMiddleware")

		if _, ok := authenTypes[AuthenticationTypePublic]; ok {
			l.Debug("(authzmiddleware) handler name >%s< is public, not checking permissions", hc.Name)
			return h(w, r, pp, qp, l, m)
		}

		AuthenData := GetRequestAuthenData(l, r)
		if AuthenData == nil {
			err := coreerror.NewInternalError("failed to read request authen data")
			l.Warn(err.Error())
			return err
		}

		// If no permissions are required, allow the request
		if len(authzPermissions) == 0 {
			l.Debug("(authzmiddleware) handler name >%s< requires no permissions, allowing request", hc.Name)
			return h(w, r, pp, qp, l, m)
		}

		for _, permission := range AuthenData.Permissions {
			if _, ok := authzPermissions[permission]; ok {
				return h(w, r, pp, qp, l, m)
			}
		}

		l.Warn("(authzmiddleware) authenticated request >%#v< does not contain any required permissions >%#v<", AuthenData.Permissions, hc.MiddlewareConfig.AuthzPermissions)

		return coreerror.NewUnauthorizedError()
	}

	return handle, nil
}
