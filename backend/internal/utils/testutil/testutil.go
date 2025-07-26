package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	runner "gitlab.com/alienspaces/playbymail/internal/runner/server"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

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

func NewTestRunner(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (*runner.Runner, error) {

	cfg, err := config.Parse()
	if err != nil {
		return nil, err
	}

	rnr, err := runner.NewRunnerWithConfig(l, s, j, cfg)
	if err != nil {
		return nil, err
	}

	rnr.AuthenticateRequestFunc = func(l logger.Logger, m domainer.Domainer, r *http.Request, authType server.AuthenticationType) (server.AuthenData, error) {
		return server.AuthenData{
			Type: server.AuthenticatedTypeToken,
		}, nil
	}

	rnr.RLSFunc = func(l logger.Logger, m domainer.Domainer, authedReq server.AuthenData) (server.RLS, error) {
		return server.RLS{
			Identifiers: map[string][]string{},
		}, nil
	}

	return rnr, nil
}

type ExpectedErrorResponse struct {
	Err  coreerror.Error
	Errs set.Set[coreerror.Code]
}

// TestRunnerer is a minimal interface for test runners
// It is satisfied by both *server.Runner and custom runners embedding it.
type TestRunnerer interface {
	ApplyMiddleware(hc server.HandlerConfig, h server.Handle) (httprouter.Handle, error)
	GetHandlerConfig() map[string]server.HandlerConfig
}

// TestCaser interface for test cases
// (restored after accidental deletion)
type TestCaser interface {
	TestName() string
	TestNewRunner(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (TestRunnerer, error)
	TestHandlerConfig(rnr TestRunnerer) server.HandlerConfig
	TestRequestHeaders(data harness.Data) map[string]string
	TestRequestPathParams(data harness.Data) map[string]string
	TestRequestQueryParams(data harness.Data) map[string]any
	TestRequestForms(data harness.Data) map[string]any
	TestRequestBody(data harness.Data) any
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
	NewRunner                 func(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (TestRunnerer, error)
	HandlerConfig             func(rnr TestRunnerer) server.HandlerConfig
	RequestHeaders            func(d harness.Data) map[string]string
	RequestPathParams         func(d harness.Data) map[string]string
	RequestQueryParams        func(d harness.Data) map[string]any
	RequestForms              func(d harness.Data) map[string]any
	RequestBody               func(d harness.Data) any
	cachedRequestBody         any
	ResponseDecoder           func(body io.Reader) (interface{}, error)
	ResponseCode              int
	ResponseError             func(harness.Data) ExpectedErrorResponse
	ShouldDecodeResponseCode  int
	ShouldNotTestResponseBody bool
	ShouldSetupTeardown       bool
	ShouldTxCommit            bool
}

// Update TestCase to implement TestCaser with TestRunnerer
func (t *TestCase) TestName() string { return t.Name }

func (t *TestCase) TestNewRunner(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (TestRunnerer, error) {
	if t.NewRunner != nil {
		return t.NewRunner(l, s, j)
	}

	rnr, err := NewTestRunner(l, s, j)
	if err != nil {
		return nil, err
	}
	return rnr, nil

}

func (t *TestCase) TestHandlerConfig(rnr TestRunnerer) server.HandlerConfig {
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

func (t *TestCase) TestRequestQueryParams(data harness.Data) map[string]any {
	qp := map[string]any{}
	if t.RequestQueryParams != nil {
		qp = t.RequestQueryParams(data)
	}
	return qp
}

func (t *TestCase) TestRequestForms(data harness.Data) map[string]any {
	qp := map[string]any{}
	if t.RequestForms != nil {
		qp = t.RequestForms(data)
	}
	return qp
}

func (t *TestCase) TestRequestBody(data harness.Data) any {
	if t.cachedRequestBody != nil {
		return t.cachedRequestBody
	}
	var b any
	if t.RequestBody != nil {
		b = t.RequestBody(data)
	}
	t.cachedRequestBody = b
	return b
}

func (t *TestCase) TestResponseDecoder(body io.Reader) (any, error) {
	var b any
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

	t.Logf("Applying middleware to handler >%s %s<", cfg.Method, cfg.Path)

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
	// request data
	data := tc.TestRequestBody(th.Data)
	multipartForms := tc.TestRequestForms(th.Data)
	var req *http.Request

	if data != nil {

		// Test data can be in the form of []byte or a marshalable struct
		var requestData []byte
		switch testData := data.(type) {
		case []byte:
			requestData = testData
		default:
			requestData, err = json.Marshal(testData)
			require.NoError(t, err, "Marshal returns without error")
		}

		req, err = http.NewRequest(cfg.Method, requestPath, bytes.NewBuffer(requestData))
		require.NoError(t, err, "NewRequest returns without error")

	} else if len(multipartForms) > 0 {
		// Prepare a form that you will submit to that URL.
		var requestBody bytes.Buffer
		w := multipart.NewWriter(&requestBody)
		for key, val := range multipartForms {
			if key == "file" {
				fileContent := val.([]byte)

				// Create a form field in the multipart request with a byte slice
				part, err := w.CreateFormFile(key, "fakefile.csv")
				if err != nil {
					fmt.Println(err)
					return
				}

				_, err = part.Write(fileContent)
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				err = w.WriteField(key, val.(string))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		w.Close()

		req, err = http.NewRequest("POST", requestPath, &requestBody)
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", w.FormDataContentType())
	} else {
		req, err = http.NewRequest(cfg.Method, requestPath, nil)
		require.NoError(t, err, "NewRequest returns without error")
	}

	// request headers
	requestHeaders := tc.TestRequestHeaders(th.Data)
	for headerKey, headerVal := range requestHeaders {
		req.Header.Add(headerKey, headerVal)
	}

	// recorder
	recorder := httptest.NewRecorder()

	// serve
	rtr.ServeHTTP(recorder, req)

	// test status
	if tc.TestResponseCode() != recorder.Code {
		t.Logf("%s", recorder.Body.String())
	}
	require.Equalf(t, tc.TestResponseCode(), recorder.Code, "%s - Response code equals expected", tc.TestName())

	// Test expected error response
	expectedErr := tc.TestResponseError(th.Data)
	if expectedErr != nil {
		var actualErrs []coreerror.Error

		err = json.NewDecoder(recorder.Body).Decode(&actualErrs)
		require.NoError(t, err, "Decode returns without error")

		for _, actual := range actualErrs {
			ok := expectedErr.Errs.Has(actual.ErrorCode)
			require.True(t, ok, "expected >%#v< actual >%v<", expectedErr, actual.ErrorCode)
		}
	}

	var responseBody any

	// Test response body
	if recorder.Code == http.StatusOK || recorder.Code == http.StatusCreated || recorder.Code == tc.TestShouldDecodeResponseCode() {

		responseBody, err = tc.TestResponseDecoder(recorder.Result().Body)
		require.NoError(t, err, fmt.Sprintf("Response body decodes without error >%#v<", err))

		if _, isStr := responseBody.(string); responseBody != nil &&
			!isStr &&
			// When handler configuration does not have a JSON schema defined
			!cfg.MiddlewareConfig.ValidateResponseSchema.IsEmpty() &&
			// When a handler supports JSON and additional content types tests
			// that exercise content types other than JSON are required to set
			// the following to true.
			!tc.TestShouldNotTestResponseBody() {

			jsonData, err := json.Marshal(responseBody)
			require.NoError(t, err, "Marshal returns without error")
			testResponseSchema(t, cfg, jsonData)
		}
	}

	if tf != nil {
		tf(cfg.Method, responseBody)
	}

}

func testResponseSchema(t *testing.T, hc server.HandlerConfig, actualRes any) {

	t.Run("response validates against JSON schema", func(t *testing.T) {
		schema := hc.MiddlewareConfig.ValidateResponseSchema
		schemaMain := schema.Main
		require.NotEmpty(t, schemaMain.Location, "handler >%s %s< ValidateResponseSchema main location path should not be empty", hc.Method, hc.Path)
		require.NotEmpty(t, schemaMain.Name, "handler >%s %s< ValidateResponseSchema main filename should not be empty", hc.Method, hc.Path)

		for _, r := range schema.References {
			require.NotEmpty(t, r.Location, "handler >%s %s< ValidateResponseSchema reference location path should not be empty", hc.Method, hc.Path)
			require.NotEmpty(t, r.Name, "handler >%s %s< ValidateResponseSchema reference filename should not be empty", hc.Method, hc.Path)
		}

		testSchemaHelper(t, schema, actualRes)
	})
}

func testSchemaHelper(t *testing.T, s jsonschema.SchemaWithReferences, actualRes any) {
	result, err := jsonschema.Validate(s, actualRes)
	require.NoError(t, err, "schema validation should not error")
	err = jsonschema.MapError(result)
	if result != nil {
		errs := result.Errors()
		for idx := range errs {
			t.Logf("schema error result >%#v<", errs[idx])
		}
	}
	require.NoError(t, err, "schema validation results should be empty")
}
