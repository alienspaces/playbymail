package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MiddlewareFunc func(hc HandlerConfig, h Handle) (Handle, error)

// defaultMiddleware provides a list of default middleware
func (rnr *Runner) defaultMiddlewareFuncs() []MiddlewareFunc {
	return []MiddlewareFunc{
		rnr.WaitMiddleware,
		rnr.ParamMiddleware,
		rnr.DataMiddleware,
		rnr.RLSMiddleware,
		rnr.JobClientMiddleware,
		rnr.AuthzMiddleware,
		rnr.AuthenMiddleware,
		rnr.TxMiddleware,
		rnr.TimerMiddleware,
		rnr.CorrelationMiddleware,
		rnr.ErrorMiddleware,
	}
}

// ResolveHandlerConfig validates and resolves handler configuration.
func (rnr *Runner) ResolveHandlerConfig(hc HandlerConfig) (HandlerConfig, error) {
	l := Logger(rnr.Log, "ResolveHandlerConfig")

	l.Info("(core) resolving handler config method >%s< path >%s<", hc.Method, hc.Path)

	hc, err := rnr.resolveHandlerSchemaLocation(hc)
	if err != nil {
		return hc, err
	}

	hc, err = rnr.resolveHandlerQueryParamsConfig(hc)
	if err != nil {
		return hc, err
	}

	return hc, nil
}

// ApplyMiddleware applies middleware by the assigned middleware function
func (rnr *Runner) ApplyMiddleware(hc HandlerConfig, h Handle) (httprouter.Handle, error) {
	l := Logger(rnr.Log, "ApplyMiddleware")

	l.Info("(core) applying middleware to handler method >%s< path >%s<", hc.Method, hc.Path)

	if rnr.HandlerMiddlewareFuncs == nil {
		l.Warn("(core) HandlerMiddlewareFuncs is nil in ApplyMiddleware, using defaultMiddlewareFuncs")
		rnr.HandlerMiddlewareFuncs = rnr.defaultMiddlewareFuncs
	}
	middlewareFuncs := rnr.HandlerMiddlewareFuncs()

	var err error
	for idx := range middlewareFuncs {
		h, err = middlewareFuncs[idx](hc, h)
		if err != nil {
			l.Warn("(core) failed adding middleware >%v<", err)
			return nil, err
		}
	}

	return rnr.httpRouterHandlerWrapper(h), nil
}

// httpRouterHandlerWrapper wraps a Handle function in an httprouter.Handle
// function while also providing a new logger for every request. Typically
// this function should be used to wrap the final product of applying all
// middleware to Handle function.
func (rnr *Runner) httpRouterHandlerWrapper(h Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, pp httprouter.Params) {
		// Create new logger with its own context (fields map) on every request
		// because each request maintains its own context (fields map). If the
		// same logger is used, when different requests set the logger context,
		// there will be concurrent map read/writes.
		l, err := rnr.Log.NewInstance()
		if err != nil {
			return
		}
		l = l.WithPackageContext("core/runner")

		// delegate
		_ = h(w, r, pp, nil, l, nil, nil)
	}
}
