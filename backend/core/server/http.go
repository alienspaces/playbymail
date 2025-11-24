package server

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"
	"github.com/rs/cors"

	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

type HttpMethod string

const (
	HttpMethodGet     HttpMethod = http.MethodGet
	HttpMethodHead    HttpMethod = http.MethodHead
	HttpMethodPost    HttpMethod = http.MethodPost
	HttpMethodPut     HttpMethod = http.MethodPut
	HttpMethodPatch   HttpMethod = http.MethodPatch
	HttpMethodDelete  HttpMethod = http.MethodDelete
	HttpMethodConnect HttpMethod = http.MethodConnect
	HttpMethodOptions HttpMethod = http.MethodOptions
	HttpMethodTrace   HttpMethod = http.MethodTrace
)

type WriteResponseOption = func(http.ResponseWriter) error

// runHTTP - Starts the HTTP server process. Override to implement a custom
// HTTP server run function. The server process exposes a REST API and is
// intended for clients to manage resources and perform actions.
func (rnr *Runner) runHTTP(args map[string]any) error {
	l := Logger(rnr.Log, "RunHTTP")

	l.Info("(core) running http")

	r := httprouter.New()

	r, err := rnr.registerRoutes(r)
	if err != nil {
		l.Warn("failed default router >%v<", err)
		return err
	}

	port := rnr.Config.Port
	if port == "" {
		l.Warn("(core) config.Port is empty, using PORT from environment")
		port = os.Getenv("PORT")
	}
	if port == "" {
		l.Warn("(core) missing PORT, cannot start server")
		return fmt.Errorf("missing PORT, cannot start server")
	}

	allowedOrigins := rnr.HTTPCORSConfig.AllowedOrigins
	if len(rnr.Config.CORSAllowedOrigins) > 0 {
		allowedOrigins = append(allowedOrigins, rnr.Config.CORSAllowedOrigins...)
	}
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}

	allowedHeaders := []string{
		"X-ProgramID", "X-ProgramName", "Content-Type",
		"Authorization", "X-Authorization-Token",
		"Origin", "X-Requested-With", "Accept",
		"X-CSRF-Token", HeaderXCorrelationID,
	}

	allowedHeaders = append(allowedHeaders, rnr.HTTPCORSConfig.AllowedHeaders...)
	if len(rnr.Config.CORSAllowedHeaders) > 0 {
		allowedHeaders = append(allowedHeaders, rnr.Config.CORSAllowedHeaders...)
	}

	l.Info("(core) http server allowed headers >%v<", allowedHeaders)
	l.Info("(core) http server allowed origins >%v<", allowedOrigins)

	c := cors.New(cors.Options{
		// Access-Control-Allow-Origin, Access-Control-Allow-Headers and
		// Access-Control-Allow-Methods cannot be wildcard if the CORS request
		// is credentialed.
		Debug:            false,
		AllowedOrigins:   allowedOrigins,
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   append(rnr.HTTPCORSConfig.ExposedHeaders, HeaderXPagination, HeaderXCorrelationID),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
		AllowCredentials: rnr.HTTPCORSConfig.AllowCredentials,
	})
	h := c.Handler(r)

	l.Info("(core) server running at: http://localhost:%s", port)

	srv := &http.Server{
		Handler:      h,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	rnr.httpServer = srv

	return srv.ListenAndServe()
}

func (rnr *Runner) ShutdownHTTP() error {
	l := Logger(rnr.Log, "ShutdownHTTP")

	if rnr.httpServer == nil {
		l.Info("(core) HTTP server not running")
		return nil
	}

	ctx, cancelHTTP := context.WithTimeout(context.Background(), 25*time.Second) // Default k8s termination grace period of 30s

	if err := rnr.httpServer.Shutdown(ctx); err != nil {
		l.Warn("(core) failed shutting down server >%v<", err)
		cancelHTTP()
		return err
	}

	// Since the HTTP server is shutdown in the main goroutine, this is actually
	// a no-op, but the linter complains if we do not cancel the ctx.
	cancelHTTP()

	return nil
}

// registerRoutes - registers routes as implemented by the assigned router
// function
func (rnr *Runner) registerRoutes(r *httprouter.Router) (*httprouter.Router, error) {
	l := Logger(rnr.Log, "registerRoutes")

	l.Info("(core) registering routes")

	return rnr.RouterFunc(r)
}

// defaultHandler is the default HandlerFunc
func (rnr *Runner) defaultHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	WriteResponse(l, w, http.StatusOK, "ok")
	return nil
}

func (rnr *Runner) registerDefaultHealthzRoute(r *httprouter.Router) (*httprouter.Router, error) {
	l := Logger(rnr.Log, "registerDefaultHealthzRoute")

	h, err := rnr.ApplyMiddleware(
		HandlerConfig{
			Path: "/healthz",
			MiddlewareConfig: MiddlewareConfig{
				AuthenTypes: []AuthenticationType{AuthenticationTypePublic},
			},
		},
		rnr.HandlerFunc,
	)
	if err != nil {
		l.Warn("(core) failed default middleware >%v<", err)
		return nil, err
	}
	r.GET("/healthz", h)

	l.Info("(core) registered /healthz")

	return r, nil
}

func (rnr *Runner) registerDefaultLivenessRoute(r *httprouter.Router) (*httprouter.Router, error) {
	l := Logger(rnr.Log, "registerDefaultLivenessRoute")

	// This logger should only be used for the liveness endpoint and is to
	// avoid creating a new logger on every request.
	hl, err := rnr.Log.NewInstance()
	if err != nil {
		l.Warn("(core) failed new log instance >%v<", err)
		return nil, err
	}
	r.GET("/liveness", func(w http.ResponseWriter, r *http.Request, pp httprouter.Params) {
		_ = rnr.HandlerFunc(w, r, pp, nil, hl, nil, nil)
	})

	l.Info("(core) registered /liveness")

	return r, nil
}

// defaultRouter - implements default routes based on runner configuration options
func (rnr *Runner) defaultRouter(r *httprouter.Router) (*httprouter.Router, error) {
	l := Logger(rnr.Log, "defaultRouter")

	var err error
	r, err = rnr.registerDefaultHealthzRoute(r)
	if err != nil {
		return nil, err
	}

	r, err = rnr.registerDefaultLivenessRoute(r)
	if err != nil {
		return nil, err
	}

	r, err = rnr.registerDefaultStaticRoutes(r)
	if err != nil {
		return nil, err
	}

	for _, hc := range rnr.HandlerConfig {

		hc, err := rnr.ResolveHandlerConfig(hc)
		if err != nil {
			l.Warn("(core) failed resolving handler config >%v<", err)
			return nil, err
		}

		h, err := rnr.ApplyMiddleware(hc, hc.HandlerFunc)
		if err != nil {
			l.Warn("(core) failed registering handler >%v<", err)
			return nil, err
		}

		switch hc.Method {
		case http.MethodGet:
			r.GET(hc.Path, h)
		case http.MethodPost:
			r.POST(hc.Path, h)
		case http.MethodPut:
			r.PUT(hc.Path, h)
		case http.MethodPatch:
			r.PATCH(hc.Path, h)
		case http.MethodDelete:
			r.DELETE(hc.Path, h)
		case http.MethodOptions:
			r.OPTIONS(hc.Path, h)
		case http.MethodHead:
			r.HEAD(hc.Path, h)
		default:
			l.Warn("(core) router HTTP method >%s< not supported", hc.Method)
			return nil, fmt.Errorf("router HTTP method >%s< not supported", hc.Method)
		}
	}

	return r, nil
}

func (rnr *Runner) registerDefaultStaticRoutes(r *httprouter.Router) (*httprouter.Router, error) {
	l := Logger(rnr.Log, "registerDefaultStaticRoutes")

	if rnr.Config.AssetsPath != "" {
		l.Info("(core) registering assets file server, serving from >%s<", rnr.Config.AssetsPath)

		// Create assets file server with cache-control headers
		assetsFileServer := http.FileServer(http.Dir(rnr.Config.AssetsPath))
		assetsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Assets can be cached, but not too aggressively in development
			w.Header().Set("Cache-Control", "public, max-age=3600")
			assetsFileServer.ServeHTTP(w, r)
		})

		// Register assets route with httprouter
		r.GET("/assets/*filepath", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Restore the original path for the file server
			r.URL.Path = "/assets" + ps.ByName("filepath")
			assetsHandler.ServeHTTP(w, r)
		})
	}

	if rnr.Config.AppHome != "" {
		l.Info("(core) registering static file server with SPA fallback, serving from >%s<", rnr.Config.AppHome)

		// Create a file server for the static directory
		fileServer := http.FileServer(http.Dir(rnr.Config.AppHome))

		// Wrap file server with cache-control headers
		cachedFileServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			filePath := r.URL.Path

			// Set cache-control headers based on file type
			if strings.HasSuffix(filePath, ".html") || filePath == "/" || filePath == "" {
				// HTML files (especially index.html) should not be cached
				// to ensure SPA always gets the latest version
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
			} else if strings.HasSuffix(filePath, ".js") || strings.HasSuffix(filePath, ".css") {
				// Check if file has Vite content hash (pattern: name-hash.ext)
				// Vite hashes look like: index-abc123.js or main-def456.css
				baseName := strings.TrimSuffix(filePath, ".js")
				baseName = strings.TrimSuffix(baseName, ".css")
				// If filename contains a dash followed by alphanumeric hash, it's hashed
				hasHash := strings.Contains(baseName, "-") && len(baseName) > 8
				if hasHash {
					// Hashed assets can be cached long-term
					w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				} else {
					// Non-hashed assets should not be cached in development
					w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				}
			} else {
				// Other static assets (images, fonts, etc.)
				w.Header().Set("Cache-Control", "public, max-age=86400")
			}

			fileServer.ServeHTTP(w, r)
		})

		// Create a custom NotFound handler that implements SPA fallback
		r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// First, try to serve the requested file
			filePath := r.URL.Path

			// Check if the file exists
			fullPath := rnr.Config.AppHome + filePath
			if _, err := os.Stat(fullPath); err == nil {
				// File exists, serve it with cache headers
				cachedFileServer.ServeHTTP(w, r)
				return
			}

			// File doesn't exist, check if it's a directory
			if _, err := os.Stat(fullPath + "/index.html"); err == nil {
				// Directory with index.html exists, serve it
				r.URL.Path = filePath + "/index.html"
				cachedFileServer.ServeHTTP(w, r)
				return
			}

			// Neither file nor directory with index.html exists
			// Serve index.html for SPA routing (except for API routes)
			if !strings.HasPrefix(filePath, "/api/") && !strings.HasPrefix(filePath, "/v1/") && !strings.HasPrefix(filePath, "/healthz") && !strings.HasPrefix(filePath, "/liveness") {
				l.Debug("(core) serving index.html for SPA route >%s<", filePath)
				r.URL.Path = "/index.html"
				cachedFileServer.ServeHTTP(w, r)
				return
			}

			// For API routes or other non-SPA routes, return 404
			http.NotFound(w, r)
		})
	}

	return r, nil
}

// HttpRouterHandlerWrapper wraps a Handle function in an httprouter.Handle function while also
// providing a new logger for every request. Typically this function should be used to wrap the
// final product of applying all middleware to Handle function.
func (rnr *Runner) HttpRouterHandlerWrapper(h Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, pp httprouter.Params) {
		// Create new logger with its own context (fields map) on every request because each
		// request maintains its own context (fields map). If the same logger is used, when
		// different requests set the logger context, there will be concurrent map read/writes.
		l, err := rnr.Log.NewInstance()
		if err != nil {
			rnr.Log.Warn("(core) failed new log instance >%v<", err)
			return
		}

		_ = h(w, r, pp, nil, l, nil, nil)
	}
}

// ReadRequest -
func ReadRequest[T any](l logger.Logger, r *http.Request, s *T) (*T, error) {

	data, err := GetRequestData(r)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	reader := bytes.NewReader(data)
	err = json.NewDecoder(reader).Decode(s)
	if err != nil {
		return nil, fmt.Errorf("failed decoding request data >%s< >%v<", string(data), err)
	}

	return s, nil
}

// ReadXMLRequest -
func ReadXMLRequest(l logger.Logger, r *http.Request, s interface{}) ([]byte, error) {

	data, err := GetRequestData(r)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	reader := bytes.NewReader(data)
	if err := xml.NewDecoder(reader).Decode(s); err != nil {
		return data, fmt.Errorf("failed decoding XML request data >%s< >%v<", string(data), err)
	}

	return data, nil
}

// WriteResponse -
func WriteResponse(l logger.Logger, w http.ResponseWriter, status int, r interface{}, options ...WriteResponseOption) error {

	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	for _, o := range options {
		if err := o(w); err != nil {
			return err
		}
	}

	w.WriteHeader(status)

	if r != nil {
		return json.NewEncoder(w).Encode(r)
	}

	return nil
}

func WriteCSVResponse(l logger.Logger, w http.ResponseWriter, status int, r string, options ...WriteResponseOption) error {

	w.Header().Set("Content-Type", "text/csv;charset=utf-8")

	for _, o := range options {
		if err := o(w); err != nil {
			return err
		}
	}

	w.WriteHeader(status)

	_, err := w.Write([]byte(r))
	return err
}

func WritePDFResponse(l logger.Logger, w http.ResponseWriter, status int, r []byte, options ...WriteResponseOption) error {

	w.Header().Set("Content-Type", "application/pdf")

	for _, o := range options {
		if err := o(w); err != nil {
			return err
		}
	}

	w.WriteHeader(status)

	if r != nil {
		_, err := w.Write(r)
		return err
	}

	return nil
}

func WritePaginatedResponse[R any, D any](l logger.Logger, w http.ResponseWriter, recs []R, mapper func(R) (D, error), pageSize int) error {
	res := []D{}

	buildPageSize := pageSize
	for _, rec := range recs {
		if buildPageSize == 0 {
			break
		}

		responseData, err := mapper(rec)
		if err != nil {
			WriteSystemError(l, w, err)
			return err
		}
		res = append(res, responseData)

		buildPageSize--
	}

	err := WriteResponse(l, w, http.StatusOK, res, XPaginationHeader(len(recs), pageSize))
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func WriteXMLResponse(l logger.Logger, w http.ResponseWriter, status int, s interface{}) error {
	w.Header().Set("Content-Type", HeaderContentTypeXML+"; charset=utf-8")

	w.WriteHeader(status)

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}

	if s != nil {
		return xml.NewEncoder(w).Encode(s)
	}

	return nil
}
