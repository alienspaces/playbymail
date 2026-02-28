package server

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/runnable"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

const (
	// ConfigKeyValidateSchemaLocation - Directory location of JSON schema's
	ConfigKeyValidateSchemaLocation string = "validateSchemaLocation"
	// ConfigKeyValidateMainSchema - Main schema that can include reference schema's
	ConfigKeyValidateMainSchema string = "validateMainSchema"
	// ConfigKeyValidateReferenceSchemas - Schema referenced from the main schema
	ConfigKeyValidateReferenceSchemas string = "validateReferenceSchemas"
)

// Runner - implements the runnerer interface
type Runner struct {
	// General configuration
	Config config.Config

	// Logger
	Log logger.Logger

	// Store
	Store storer.Storer

	// JobClient is the riverqueue client for registering and processing jobs.
	//
	// It is unsafe to access this job client instance from the handler functions
	// and should only be accessed from the server daemon to start and stop the
	// job worker.
	//
	// A separate job client instance is created in middleware and passed along to
	// handler functions and can be used solely for registering jobs within the
	// handler functions.
	JobClient *river.Client[pgx.Tx]

	// HTTPCORSConfig
	HTTPCORSConfig HTTPCORSConfig

	// Handler and message configuration
	HandlerConfig map[string]HandlerConfig

	// HTTP server
	httpServer *http.Server

	// Assignable functions
	RunHTTPFunc   func(args map[string]any) error
	RunDaemonFunc func(ctx context.Context, args map[string]any) error

	RouterFunc func(router *httprouter.Router) (*httprouter.Router, error)

	// HandlerFunc is the default handler function. It is used for liveness and
	// healthz. Therefore, it should execute quickly.
	HandlerFunc Handle

	// HandlerMiddlewareFuncs returns a list of middleware functions to apply to
	// routes
	HandlerMiddlewareFuncs func() []MiddlewareFunc

	// DomainFunc returns the service specific domainer implementation
	DomainFunc func(l logger.Logger) (domainer.Domainer, error)

	// JobClientFunc returns a new job client instance.
	JobClientFunc func(l logger.Logger, s storer.Storer) (*river.Client[pgx.Tx], error)

	// AuthenticateRequestFunc authenticates a request based on the authentication type
	AuthenticateRequestFunc func(l logger.Logger, m domainer.Domainer, r *http.Request, authType AuthenticationType) (AuthenData, error)

	// RLSFunc is the service specific RLS implementation which is responsible for
	// returning a map of identifiers the current user has access to. This is used
	// to set the RLS constraints on the repositories.
	RLSFunc func(l logger.Logger, m domainer.Domainer, authedReq AuthenData) (RLS, error)
}

type HTTPCORSConfig struct {
	AllowedOrigins   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
}

var _ runnable.Runnable = &Runner{}

// Handle - custom service handle
type Handle func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error

// HandlerConfig - configuration for routes, handlers and middleware
type HandlerConfig struct {
	Name string
	// Method - The HTTP method
	Method string
	// Path - The HTTP request URI including :parameter placeholders
	Path string
	// HandlerFunc - Function to handle requests for this method and path
	HandlerFunc Handle
	// MiddlewareConfig -
	MiddlewareConfig MiddlewareConfig
	// DocumentationConfig -
	DocumentationConfig DocumentationConfig
}

type AuthenData struct {
	Type        AuthenticatedType      `json:"type"`
	RLSType     RLSType                `json:"-"`
	AccountUser AuthenticatedAccountUser `json:"account_user"`
	Permissions []AuthorizedPermission `json:"permissions"`
}

func (a AuthenData) IsAuthenticated() bool {
	return a.Type != ""
}

// HasPermission checks if the AuthenData contains a specific permission
func (a AuthenData) HasPermission(permission AuthorizedPermission) bool {
	for _, p := range a.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

type RLSType string

const (
	RLSTypeOpen       RLSType = "open"
	RLSTypeRestricted RLSType = "restricted"
)

type AuthenticatedType string

const (
	AuthenticatedTypeJWT   AuthenticatedType = "JWT"   // JSON Web Token
	AuthenticatedTypeKey   AuthenticatedType = "Key"   // API Key
	AuthenticatedTypeToken AuthenticatedType = "Token" // Session Token
)

type AuthenticatedAccountUser struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type AuthenticationType string
type AuthorizedPermission string

const (
	AuthenticationTypePublic        AuthenticationType = "Public"
	AuthenticationTypeJWT           AuthenticationType = "JWT"           // JSON Web Token
	AuthenticationTypeKey           AuthenticationType = "Key"           // API Key
	AuthenticationTypeToken         AuthenticationType = "Token"         // Session Token
	AuthenticationTypeOptionalToken AuthenticationType = "OptionalToken" // Optional session token - handler receives auth data if valid token present, nil otherwise
)

// MiddlewareConfig - configuration for global default middleware
type MiddlewareConfig struct {
	AuthenTypes            []AuthenticationType
	AuthzPermissions       []AuthorizedPermission
	ValidateRequestSchema  jsonschema.SchemaWithReferences
	ValidateResponseSchema jsonschema.SchemaWithReferences
	ValidateParamsConfig   *ValidateParamsConfig
}

// ValidateParamsConfig defines how route path parameters should be validated
//
// ExcludePathParamsFromQueryParams - By default path parameters will be added
// as query parameters and validated as part of query parameter validation.
// When disabled, path parameters will need to be validated by the handler.
// Schema - Validate query parameters using this JSON schema set
// QueryParams - Specifies the query parameters expected for the route.
type ValidateParamsConfig struct {
	ExcludePathParamsFromQueryParams bool
	Schema                           jsonschema.SchemaWithReferences
	QueryParams                      QueryParams
	queryParamTypes                  map[string]jsonschema.JSONType
}

// QueryParams is used to ensure that the ValidateParamsConfig.QueryParams is
// actually a QueryParams type, not accidentally some other API type.
type QueryParams interface {
	GetPageNumber() int
	GetPageSize() int
	GetSortColumns() []string
}

// DocumentationConfig - Configuration describing how to document a route
type DocumentationConfig struct {
	Document        bool
	Collection      bool
	Title           string // used for API doc endpoint title
	Description     string // used for API doc endpoint description
	RequestHeaders  []Header
	ResponseHeaders []Header
}

type Header struct {
	Name     string
	Required bool
	Schema   jsonschema.SchemaWithReferences
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

// NewRunnerWithConfig - creates a new runner with provided configuration. This is useful for testing
// with configuration that is not sourced from the environment or defaults.
func NewRunner(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (*Runner, error) {
	l = l.WithPackageContext("runner")

	rnr := &Runner{
		Config:                 cfg,
		Log:                    l,
		Store:                  s,
		JobClient:              j,
		HandlerConfig:          make(map[string]HandlerConfig),
		HandlerMiddlewareFuncs: nil,
	}

	rnr.HandlerFunc = rnr.defaultHandler
	rnr.RouterFunc = rnr.defaultRouter
	rnr.RunHTTPFunc = rnr.runHTTP
	rnr.RunDaemonFunc = rnr.runDaemon

	return rnr, nil
}

// Init -
func (rnr *Runner) Init(s storer.Storer) error {

	rnr.Log.Info("(core) init")

	rnr.Store = s

	// run server
	if rnr.RunHTTPFunc == nil {
		rnr.RunHTTPFunc = rnr.runHTTP
	}

	// run daemon
	if rnr.RunDaemonFunc == nil {
		rnr.RunDaemonFunc = rnr.runDaemon
	}

	// http server - router
	if rnr.RouterFunc == nil {
		rnr.Log.Warn("(core) RouterFunc is nil, using default router function")
		rnr.RouterFunc = rnr.defaultRouter
	}

	// http server - middleware
	if rnr.HandlerMiddlewareFuncs == nil {
		rnr.Log.Warn("(core) HandlerMiddlewareFuncs is nil, using default middleware functions")
		rnr.HandlerMiddlewareFuncs = rnr.defaultMiddlewareFuncs
	}

	// http server - handler
	if rnr.HandlerFunc == nil {
		rnr.Log.Warn("(core) HandlerFunc is nil, using default handler function")
		rnr.HandlerFunc = rnr.defaultHandler
	}

	if err := validateAuthenticationTypes(rnr.HandlerConfig); err != nil {
		rnr.Log.Warn("(core) failed to validate authentication types >%v<", err)
		return err
	}

	return nil
}

// InitJobClient initialises and returns a new job client instance.
func (rnr *Runner) InitJobClient(l logger.Logger) (*river.Client[pgx.Tx], error) {
	l = Logger(l, "InitJobClient")

	l.Info("(core) initialising job client")

	// NOTE: This is called from txmiddleware during request processing.
	// The domain model is created fresh for each HTTP request and passed along
	// the handler chain, ensuring each request has a unique domain instance.
	// This prevents shared state issues when multiple requests are processed
	// concurrently, as each request gets its own isolated domain context.
	if rnr.JobClientFunc == nil {
		return nil, fmt.Errorf("(core) JobClientFunc is nil on runner: %p", rnr)
	}

	return rnr.JobClientFunc(l, rnr.Store)
}

// InitDomain initialises and returns a new domain
func (rnr *Runner) InitDomain(l logger.Logger) (domainer.Domainer, error) {
	l = Logger(l, "InitDomain")

	l.Info("(core) initialising domain")

	// NOTE: This is called from txmiddleware during request processing.
	// The domain model is created fresh for each HTTP request and passed along
	// the handler chain, ensuring each request has a unique domain instance.
	// This prevents shared state issues when multiple requests are processed
	// concurrently, as each request gets its own isolated domain context.
	if rnr.DomainFunc == nil {
		return nil, fmt.Errorf("(core) DomainFunc is nil on runner: %p", rnr)
	}

	l.Info("(core) calling DomainFunc")

	m, err := rnr.DomainFunc(l)
	if err != nil {
		l.Warn("(core) failed DomainFunc >%v<", err)
		return nil, err
	}
	if m == nil {
		err := fmt.Errorf("domainer is nil, cannot initialise domain")
		return nil, err
	}

	l.Info("(core) calling Store.BeginTx")
	tx, err := rnr.Store.BeginTx()
	if err != nil {
		l.Warn("(core) failed store.BeginTx >%v<", err)
		return nil, err
	}

	l.Info("(core) calling m.Init")
	err = m.Init(tx)
	if err != nil {
		l.Warn("(core) failed domain.Init >%v<", err)
		return m, err
	}

	return m, nil
}

// Run starts the HTTP server and daemon processes.
func (rnr *Runner) Run(args map[string]any) error {
	l := Logger(rnr.Log, "Run")

	l.Info("(core) Starting http and daemon processes")

	// signal channel
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var err error

	// Calling on the store pool here ensures there is no concurrent call
	// by the http or daemon processes that might result in multiple pools.
	if rnr.Store != nil {
		_, err := rnr.Store.Pool()
		if err != nil {
			return err
		}
	}

	// run HTTP server
	go func() {
		httperr := rnr.RunHTTPFunc(args)
		if httperr != nil && !errors.Is(httperr, http.ErrServerClosed) {
			err = fmt.Errorf("failed run http: %w", httperr)
			l.Warn(err.Error())
			sigCh <- syscall.SIGTERM
		}
	}()

	// Daemon context and cancellation
	ctx, cancelDaemon := context.WithCancel(context.Background())

	// This is to ensure that the daemon server shuts down before the main goroutine exits.
	// The same is not needed for the HTTP server because of the server shutdown below
	daemonWg := &sync.WaitGroup{}
	daemonWg.Add(1)
	go func() {
		daemonerr := rnr.RunDaemonFunc(ctx, args)
		if daemonerr != nil && !errors.Is(daemonerr, context.Canceled) {
			err = fmt.Errorf("failed run daemon: %w", daemonerr)
			l.Warn(err.Error())
			sigCh <- syscall.SIGTERM
		}
		daemonWg.Done()
	}()

	<-sigCh
	go func() {
		// if SIGTERM is sent on the HTTP server and then the daemon server,
		// we will be stuck waiting for the daemonWg to be done
		<-sigCh
	}()

	rnr.LogMemStats(l)

	l.Info("(core) Shutting down daemon")

	cancelDaemon()

	l.Info("(core) Shutting down http")

	if err := rnr.ShutdownHTTP(); err != nil {
		l.Error("failed shutting down http server >%v<", err)
	}

	l.Info("(core) Waiting for daemon to exit")

	daemonWg.Wait()

	if rnr.Store != nil {
		rnr.Store.ClosePool()
	}

	return err
}

func (rnr *Runner) LogMemStats(l logger.Logger) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	bToKb := func(b uint64) uint64 {
		return b / 1024
	}

	l.Info("(core) Mem alloc >%d Kb< sys >%d kb< numGC >%d<",
		bToKb(m.Alloc),
		bToKb(m.Sys),
		m.NumGC,
	)
}

func (rnr *Runner) resolveHandlerSchemaLocation(hc HandlerConfig) (HandlerConfig, error) {

	schemaPath := rnr.Config.SchemaPath

	// Query and path paramater schemas
	if hc.MiddlewareConfig.ValidateParamsConfig != nil {
		if schemaPath == "" {
			err := fmt.Errorf("missing SchemaPath")
			rnr.Log.Warn(err.Error())
			return hc, err
		}
		schema := hc.MiddlewareConfig.ValidateParamsConfig.Schema
		if len(schema.Main.Name) > 0 {
			hc.MiddlewareConfig.ValidateParamsConfig.Schema = jsonschema.ResolveSchemaLocation(schemaPath, schema)
		}
	}

	// Request schemas
	if len(hc.MiddlewareConfig.ValidateRequestSchema.Main.Name) > 0 {
		if schemaPath == "" {
			err := fmt.Errorf("missing SchemaPath")
			rnr.Log.Warn(err.Error())
			return hc, err
		}
		hc.MiddlewareConfig.ValidateRequestSchema = jsonschema.ResolveSchemaLocation(schemaPath, hc.MiddlewareConfig.ValidateRequestSchema)
	}

	// Response schemas
	if len(hc.MiddlewareConfig.ValidateResponseSchema.Main.Name) > 0 {
		if schemaPath == "" {
			err := fmt.Errorf("missing SchemaPath")
			rnr.Log.Warn(err.Error())
			return hc, err
		}
		hc.MiddlewareConfig.ValidateResponseSchema = jsonschema.ResolveSchemaLocation(schemaPath, hc.MiddlewareConfig.ValidateResponseSchema)
	}

	// Headers for documentation
	if len(hc.DocumentationConfig.RequestHeaders) > 0 {
		if schemaPath == "" {
			err := fmt.Errorf("missing SchemaPath")
			rnr.Log.Warn(err.Error())
			return hc, err
		}
		for i, header := range hc.DocumentationConfig.RequestHeaders {
			hc.DocumentationConfig.RequestHeaders[i].Schema = jsonschema.ResolveSchemaLocation(schemaPath, header.Schema)
		}
	}

	return hc, nil
}

func validateAuthenticationTypes(handlerConfig map[string]HandlerConfig) error {
	for _, cfg := range handlerConfig {
		if len(cfg.MiddlewareConfig.AuthenTypes) == 0 {
			return fmt.Errorf("handler method >%s< path >%s< with undefined authentication type", cfg.Method, cfg.Path)
		}
	}
	return nil
}

// GetHandlerConfig returns the HandlerConfig map
func (r *Runner) GetHandlerConfig() map[string]HandlerConfig {
	return r.HandlerConfig
}

func MergeHandlerConfigs(hc1 map[string]HandlerConfig, hc2 map[string]HandlerConfig) map[string]HandlerConfig {
	if hc1 == nil {
		hc1 = map[string]HandlerConfig{}
	}
	maps.Copy(hc1, hc2)
	return hc1
}
