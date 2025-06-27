package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

type requestData struct {
	Name string
	Age  int
}

type responseData struct{}

type clientTestCase struct {
	name        string
	method      string
	path        string
	params      map[string]string
	requestData *requestData
	serverFunc  func(rw http.ResponseWriter, req *http.Request)
	expectErr   bool
}

// NewDefaultDependencies -
func NewDefaultDependencies() (logger.Logger, error) {
	cfg := config.Config{}
	err := config.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	l, err := log.NewLogger(cfg)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func TestRetryRequest(t *testing.T) {

	l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tests := []clientTestCase{
		{
			name:   "Get resource OK",
			method: http.MethodGet,
			path:   "/api/collections/:collection_id/members",
			params: map[string]string{
				"id":            "52fdfc07-2182-454f-963f-5f0f9a621d72",
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				respData, err := json.Marshal(&responseData{})
				if err != nil {
					l.Warn("Failed encoding data >%v<", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusOK)
				rw.Write(respData)
			},
			expectErr: false,
		},
		{
			name:   "Get resource BadRequest",
			method: http.MethodGet,
			path:   "/api/collections/:collection_id",
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		{
			name:   "Post resource OK",
			method: http.MethodPost,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "John",
				Age:  10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				requestData := requestData{}
				err := json.NewDecoder(req.Body).Decode(&requestData)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				if requestData.Name != "John" {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				if requestData.Age != 10 {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				respData, err := json.Marshal(&responseData{})
				if err != nil {
					l.Warn("Failed encoding data >%v<", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusOK)
				rw.Write(respData)
			},
			expectErr: false,
		},
		{
			name:   "Post resource BadRequest",
			method: http.MethodPost,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "Mary",
				Age:  20,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		{
			name:   "Post resource error - missing request data",
			method: http.MethodPost,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: nil,
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		{
			name:   "Put resource OK",
			method: http.MethodPut,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "John",
				Age:  10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				requestData := requestData{}
				err := json.NewDecoder(req.Body).Decode(&requestData)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				id := strings.TrimPrefix(req.URL.Path, "/api/collections/")
				if id != "52fdfc07-2182-454f-963f-5f0f9a621d72" {
					rw.WriteHeader(http.StatusNotFound)
					return
				}

				if requestData.Name != "John" {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				if requestData.Age != 10 {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				respData, err := json.Marshal(&responseData{})
				if err != nil {
					l.Warn("Failed encoding data >%v<", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusOK)
				rw.Write(respData)
			},
			expectErr: false,
		},
		{
			name:   "Put resource NotFound",
			method: http.MethodPut,
			path:   "/api/collections",
			params: map[string]string{
				"id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "John",
				Age:  10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
			},
			expectErr: true,
		},
		{
			name:   "Put resource error - missing request data",
			method: http.MethodPut,
			path:   "/api/collections",
			params: map[string]string{
				"id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		// TODO delete method not yet supported in core/client/client.go
		{
			name:   "Delete resource OK",
			method: http.MethodDelete,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "John",
				Age:  10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				id := strings.TrimPrefix(req.URL.Path, "/api/collections/")
				if id != "52fdfc07-2182-454f-963f-5f0f9a621d72" {
					rw.WriteHeader(http.StatusNotFound)
					return
				}
				rw.WriteHeader(http.StatusOK)
			},
			expectErr: false,
		},
		{
			name:   "Delete resource NotFound",
			method: http.MethodDelete,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-abcd-abcd-abcde-5f0f9a621d72",
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
			},
			expectErr: true,
		},
	}

	testRequest := func(t *testing.T, tc clientTestCase, methodName string, resp *responseData, err error) {
		if tc.expectErr == true {
			require.Errorf(t, err, "%s returns with error", methodName)
			return
		}
		require.NoErrorf(t, err, "%s returns without error", methodName)
		require.NotNilf(t, resp, "%s returns a response", methodName)
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {
			// Test HTTP server
			server := httptest.NewServer(http.HandlerFunc(tc.serverFunc))
			defer server.Close()

			// HTTPClient
			cl, err := NewHTTPClient(l)
			require.NoError(t, err, "NewHTTPClient returns without error")
			require.NotNil(t, cl, "NewHTTPClient returns a client")

			// Host
			cl.Host = server.URL

			// set max retries to speed tests up
			cl.MaxRetries = 2

			resp := &responseData{}
			_, err = cl.RetryRequest(tc.method, tc.path, tc.params, tc.requestData, resp)
			testRequest(t, tc, "RetryRequest", resp, err)

			switch tc.method {
			case http.MethodGet:
				err = cl.Get(tc.path, tc.params, resp)
				testRequest(t, tc, "Get", resp, err)
			case http.MethodPost:
				if tc.requestData != nil {
					err = cl.Create(tc.path, tc.params, tc.requestData, resp)
				} else {
					err = cl.Create(tc.path, tc.params, nil, resp)
				}
				testRequest(t, tc, "Post", resp, err)
			case http.MethodPut:
				if tc.requestData != nil {
					_ = cl.Update(tc.path, tc.params, tc.requestData, resp)
				} else {
					_ = cl.Update(tc.path, tc.params, nil, resp)
				}
				err = cl.Update(tc.path, tc.params, tc.requestData, resp)
				testRequest(t, tc, "Update", resp, err)
			case http.MethodDelete:
				err = cl.Delete(tc.path, tc.params, resp)
				testRequest(t, tc, "Delete", resp, err)
			}
		})
	}
}

func Test_buildURL(t *testing.T) {

	l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	// Client
	cl, err := NewHTTPClient(l)

	// Host
	cl.Host = "http://example.com"

	require.NoError(t, err, "NewHTTPClient returns without error")
	require.NotNil(t, cl, "NewHTTPClient returns a client")

	// Set base path
	cl.Path = "/api"

	tests := []struct {
		name      string
		path      string
		params    map[string]string
		expectErr bool
		expectURL string
	}{
		{
			name: "Build URL with colon prefixed path parameters and without colon prefixed parameters",
			path: "/collections/:collection_id/members/:member_id",
			params: map[string]string{
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				"member_id":     "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name: "Build URL without colon prefixed path parameters and without colon prefixed parameters",
			path: "/collections/collection_id/members/member_id",
			params: map[string]string{
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				"member_id":     "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name: "Build URL with colon prefixed path parameters and colon prefixed parameters",
			path: "/collections/:collection_id/members/:member_id",
			params: map[string]string{
				":collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				":member_id":     "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name: "Build URL without colon prefixed path params and colon prefixed parameters",
			path: "/collections/collection_id/members/member_id",
			params: map[string]string{
				":collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				":member_id":     "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name: "Build URL with colon prefixed path parameters and without colon prefixed parameters that include additional query params",
			path: "/collections/:collection_id/members",
			params: map[string]string{
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				"qp1":           "1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members?qp1=1",
		},
		{
			name: "Build URL with colon prefixed path parameters and with colon prefixed parameters that include additional query params",
			path: "/collections/:collection_id/members",
			params: map[string]string{
				":collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				":qp1":           "1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members?qp1=1",
		},
		{
			name: "Build URL with colon prefixed path parameters and without colon prefixed parameters that have empty values",
			path: "/collections/:collection_id/members",
			params: map[string]string{
				"collection_id": "",
			},
			expectErr: true,
		},
		{
			name: "Build URL with colon prefixed path parameters and with colon prefixed parameters that have empty values",
			path: "/collections/:collection_id/members",
			params: map[string]string{
				":collection_id": "",
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			url, err := cl.buildURL(tc.path, tc.params)
			if tc.expectErr == true {
				require.Error(t, err, "buildURL returns with error")
				return
			}
			t.Logf("Resulting URL >%s<", url)
			require.NoError(t, err, "buildURL returns without error")
			require.Equal(t, tc.expectURL, url, "buildURL returns expected URL")
		})
	}
}
