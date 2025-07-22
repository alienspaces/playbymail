package testing

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

type ExpectedErrorResponse struct {
	Err  coreerror.Error
	Errs set.Set[coreerror.Code]
}

type TestCaser interface {
	TestName() string
	TestNewRunner(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (*server.Runner, error)
	TestHandlerConfig(rnr *server.Runner) server.HandlerConfig
	TestRequestHeaders(data harness.Data) map[string]string
	TestRequestPathParams(data harness.Data) map[string]string
	TestRequestQueryParams(data harness.Data) map[string]interface{}
	TestRequestForms(data harness.Data) map[string]interface{}
	TestRequestBody(data harness.Data) interface{}
	TestResponseDecoder(body io.Reader) (interface{}, error)
	TestResponseCode() int
	TestResponseError(data harness.Data) *ExpectedErrorResponse
	TestShouldDecodeResponseCode() int
	TestShouldNotTestResponseBody() bool
	TestShouldSetupTeardown() bool
	TestShouldTxCommit() bool
}

type TestCase struct {
	Skip                      bool
	Name                      string
	NewRunner                 func(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (*server.Runner, error)
	HandlerConfig             func(rnr *server.Runner) server.HandlerConfig
	RequestHeaders            func(d harness.Data) map[string]string
	RequestPathParams         func(d harness.Data) map[string]string
	RequestQueryParams        func(d harness.Data) map[string]interface{}
	RequestForms              func(d harness.Data) map[string]interface{}
	RequestBody               func(d harness.Data) interface{}
	cachedRequestBody         interface{}
	ResponseDecoder           func(body io.Reader) (interface{}, error)
	ResponseCode              int
	ResponseError             func(harness.Data) ExpectedErrorResponse
	ShouldDecodeResponseCode  int
	ShouldNotTestResponseBody bool
	ShouldSetupTeardown       bool
	ShouldTxCommit            bool
}

func (t *TestCase) TestName() string { return t.Name }

func (t *TestCase) TestNewRunner(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (*server.Runner, error) {
	if t.NewRunner != nil {
		return t.NewRunner(l, s, j)
	}
	cfg, err := config.Parse()
	if err != nil {
		return nil, err
	}
	rnr, err := server.NewRunnerWithConfig(l, s, j, cfg.Config)
	if err != nil {
		return nil, err
	}
	return rnr, nil
}

func (t *TestCase) TestHandlerConfig(rnr *server.Runner) server.HandlerConfig {
	return t.HandlerConfig(rnr)
}

func (t *TestCase) TestRequestHeaders(data harness.Data) map[string]string {
	headers := map[string]string{}
	if t.RequestHeaders != nil {
		headers = t.RequestHeaders(data)
	}
	if !t.ShouldTxCommit {
		headers[server.HeaderXTxRollback] = "true"
	}
	return headers
}

func (t *TestCase) TestRequestPathParams(data harness.Data) map[string]string {
	pp := map[string]string{}
	if t.RequestPathParams != nil {
		pp = t.RequestPathParams(data)
	}
	return pp
}

func (t *TestCase) TestRequestQueryParams(data harness.Data) map[string]interface{} {
	qp := map[string]interface{}{}
	if t.RequestQueryParams != nil {
		qp = t.RequestQueryParams(data)
	}
	return qp
}

func (t *TestCase) TestRequestForms(data harness.Data) map[string]interface{} {
	qp := map[string]interface{}{}
	if t.RequestForms != nil {
		qp = t.RequestForms(data)
	}
	return qp
}

func (t *TestCase) TestRequestBody(data harness.Data) interface{} {
	if t.cachedRequestBody != nil {
		return t.cachedRequestBody
	}
	var b interface{}
	if t.RequestBody != nil {
		b = t.RequestBody(data)
	}
	t.cachedRequestBody = b
	return b
}

func (t *TestCase) TestResponseDecoder(body io.Reader) (interface{}, error) {
	var b interface{}
	var err error
	if t.ResponseDecoder != nil {
		b, err = t.ResponseDecoder(body)
	}
	return b, err
}

func (t *TestCase) TestResponseCode() int { return t.ResponseCode }

func (t *TestCase) TestResponseError(data harness.Data) *ExpectedErrorResponse {
	if t.ResponseError != nil {
		respErr := t.ResponseError(data)
		if respErr.Errs == nil && respErr.Err.ErrorCode != "" {
			respErr.Errs = set.New(respErr.Err.ErrorCode)
		}
		return &respErr
	}
	return nil
}

func (t *TestCase) TestShouldDecodeResponseCode() int   { return t.ShouldDecodeResponseCode }
func (t *TestCase) TestShouldNotTestResponseBody() bool { return t.ShouldNotTestResponseBody }
func (t *TestCase) TestShouldSetupTeardown() bool       { return t.ShouldSetupTeardown }
func (t *TestCase) TestShouldTxCommit() bool            { return t.ShouldTxCommit }

func TestCaseResponseDecoderGeneric[T any](body io.Reader) (any, error) {
	var responseBody T
	err := json.NewDecoder(body).Decode(&responseBody)
	return responseBody, err
}

func RunTestCase(t *testing.T, th *harness.Testing, tc TestCaser, tf func(method string, body any)) {
	require.NotNil(t, th, "Test harness is not nil")

	rnr, err := tc.TestNewRunner(th.Log, th.Store, th.JobClient)
	require.NoError(t, err, "TestNewRunner returns without error")
	require.NotNil(t, rnr, "TestNewRunner returns a new Runner")

	if tc.TestShouldSetupTeardown() {
		_, err := th.Setup()
		require.NoError(t, err, "Test data setup returns without error")
		defer func() {
			err = th.Teardown()
			require.NoError(t, err, "Test data teardown returns without error")
		}()
	}

	cfg := tc.TestHandlerConfig(rnr)
	h, err := rnr.ApplyMiddleware(cfg, cfg.HandlerFunc)
	require.NoError(t, err, "ApplyMiddleWare returns without error")

	rtr := httprouter.New()
	switch cfg.Method {
	case http.MethodGet:
		rtr.GET(cfg.Path, h)
	case http.MethodPost:
		rtr.POST(cfg.Path, h)
	case http.MethodPut:
		rtr.PUT(cfg.Path, h)
	case http.MethodDelete:
		rtr.DELETE(cfg.Path, h)
	}

	requestParams := tc.TestRequestPathParams(th.Data)
	requestPath := cfg.Path
	for paramKey, paramValue := range requestParams {
		requestPath = strings.Replace(requestPath, paramKey, paramValue, 1)
	}
	t.Logf("Request path >%s<", requestPath)

	queryParams := tc.TestRequestQueryParams(th.Data)
	if len(queryParams) > 0 {
		requestPath += `?`
		for paramKey, paramValue := range queryParams {
			t.Logf("Adding parameter key >%s< param >%s<", paramKey, paramValue)
			switch v := paramValue.(type) {
			case int:
				requestPath = fmt.Sprintf("%s%s=%d&", requestPath, paramKey, v)
			case string:
				requestPath = fmt.Sprintf("%s%s=%s&", requestPath, paramKey, url.QueryEscape(v))
			case bool:
				requestPath = fmt.Sprintf("%s%s=%t&", requestPath, paramKey, v)
			case []string:
				sb := strings.Builder{}
				for i, s := range v {
					sb.WriteString(url.QueryEscape(s))
					if i+1 != len(v) {
						sb.WriteString(",")
					}
				}
				requestPath = fmt.Sprintf("%s%s[]=%s&", requestPath, paramKey, sb.String())
			default:
				t.Errorf("Unsupported query parameter type for value >%v<", v)
			}
		}
		t.Logf("Request path with query params >%s<", requestPath)
	}
	// ... existing code ...
}

func NewDefaultDependencies(t *testing.T) (logger.Logger, storer.Storer, *river.Client[pgx.Tx], config.Config) {
	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	return l, s, j, cfg
}

func NewTestHarness(t *testing.T) *harness.Testing {
	config := harness.DefaultDataConfig()
	l, s, j, cfg := NewDefaultDependencies(t)
	h, err := harness.NewTesting(l, s, j, cfg, config)
	require.NoError(t, err, "NewTesting returns without error")
	h.ShouldCommitData = true
	err = h.Teardown()
	require.NoError(t, err, "Teardown returns without error")
	return h
}
