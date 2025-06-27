package server

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

var unauthErr = coreerror.NewUnauthenticatedError("Authentication failed")

// AuthenMiddleware -
func (rnr *Runner) AuthenMiddleware(hc HandlerConfig, h Handle) (Handle, error) {
	handlerAuthenTypes := set.New(hc.MiddlewareConfig.AuthenTypes...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
		l = Logger(l, "AuthenMiddleware")

		if _, ok := handlerAuthenTypes[AuthenticationTypePublic]; ok {
			l.Debug("(core) handler name >%s< is public, not authenticating", hc.Name)
			return h(w, r, pp, qp, l, m)
		}

		var auth AuthenticatedRequest
		var err error
		hasAuthnType := false

		authnTypes := []AuthenticationType{
			AuthenticationTypeAPIKey,
			AuthenticationTypeJWT,
		}

		// A handler may have multiple authentication types. If supported, every
		// authentication type must be attempted, until one succeeds.
		for _, authnType := range authnTypes {
			if !handlerAuthenTypes.Has(authnType) {
				continue
			}

			hasAuthnType = true

			// AuthenticateRequestFunc expects any returned error to be an
			// UnauthenticatedError or a 500 error.
			auth, err = rnr.AuthenticateRequestFunc(l, m, r, authnType)
			if err != nil && !coreerror.IsUnauthenticatedError(err) {
				l.Error("(core) failed to authenticate request authnType >%s< err >%v<", authnType, err)
				return unauthErr // For security reasons, always respond with ambiguous 401
			}
			if err != nil {
				l.Warn("(core) failed to authenticate request authnType >%s< err >%v<", authnType, err)
			}

			if auth.IsAuthenticated() {
				break
			}
		}

		if !hasAuthnType {
			l.Error("(core) handler name >%s< with non-public authentication type has no registered authentication types", hc.Name)
			return unauthErr
		}

		if !auth.IsAuthenticated() {
			l.Warn("(core) failed to authenticate request >%#v<", err)
			return unauthErr
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyAuth, auth)
		r = r.WithContext(ctx)

		switch id := auth.User.ID.(type) {
		case string:
			l.Context("user-id", id)
		case []byte:
			requesterID := base64.URLEncoding.EncodeToString(id)
			l.Context("user-id", requesterID)
		default:
			err := coreerror.NewInternalError("unknown auth user ID type")
			l.Warn(err.Error())
			return err
		}

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}
