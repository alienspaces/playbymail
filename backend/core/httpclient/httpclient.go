package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const (
	maxRetries int = 5
	// AuthTypeBearer will use the AuthToken as Bearer data
	AuthTypeBearer string = "JWT"
	// AuthTypeBasic will use AuthUser and AuthPass as credentials
	AuthTypeBasic string = "Basic"
)

type AuthType string

// HTTPClient -
type HTTPClient struct {
	Log    logger.Logger
	client *http.Client

	// RequestLogFunc will be called with the request URL, resulting request data and response data
	// to be used by client consumers wanting to store requests and responses for debugging etc
	RequestLogFunc func(url, requestData, responseData string)

	MaxRetries int
	// Path is the base path for all requests
	Path string
	// Host is the host for all requests
	Host string
	// Setting AuthKey will result in the header "Authorisation: [AuthKey]"
	AuthKey string
	// Setting AuthToken will result in the header "Authorisation: Bearer [AuthKey]"
	AuthToken string
	// Setting AuthUser/AuthPass will result in the header "Basic: [AuthKey]"
	AuthUser string
	AuthPass string

	// Verbose will log requests at info log level as opposed to the default debug log level. This
	// helps by limiting noise from other packages when attempting to determine HTTP client interactions
	Verbose bool

	// RequestHeaders provides a way to add additional request headers
	RequestHeaders map[string]string
}

// Error with HTTP response code
type Error struct {
	Body       string
	Err        error
	retryAfter time.Duration
}

func (e Error) Error() string {
	return fmt.Sprintf("Error >%v< Body >%s<", e.Err, e.Body)
}

// NewHTTPClient -
func NewHTTPClient(l logger.Logger) (*HTTPClient, error) {

	cl := HTTPClient{
		Log: l,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	return &cl, nil
}

// RetryRequest always responds with a *http.Response struct if available, even if an error is also returned.
func (c *HTTPClient) RetryRequest(method, path string, params map[string]string, reqData interface{}, respData interface{}, opts ...func(r *http.Request) *http.Request) (*http.Response, error) {
	l := loggerWithFunctionContext(c.Log, "RetryRequest")

	var err error

	// Replace placeholder parameters and add query parameters
	url, err := c.buildURL(path, params)
	if err != nil {
		l.Warn("failed building URL >%v<", err)
		return nil, err
	}

	data, err := c.encodeData(reqData)
	if err != nil {
		l.Warn("failed marshalling request data >%v<", err)
		return nil, err
	}

	if c.MaxRetries == 0 {
		c.MaxRetries = maxRetries
	}

	var resp *http.Response
	var respErr *Error

	l.Info("Retrying request method >%s< host >%s< path >%s< params >%#v<", method, c.Host, path, params)

	// Packages:
	// 1. Currently, https://github.com/hashicorp/go-retryablehttp/ does not implement exponential backoff with jitter,
	//    which is weird, because of thundering herds.
	// 2. backoff does, but does not respect 429s: https://github.com/cenkalti/backoff/issues/134.
	//
	// We can use backoff.RetryNotify to ensure we do not retry until after the Retry-After header, if present.
	// Retries and 429s are likely rare enough that we do not need to worry about it for now.
	err = backoff.RetryNotify(func() error {
		// respErr must be set to nil, so we can distinguish between a successful request with a previous retry,
		// and a failed request.
		respErr = nil
		resp, err = c.Request(method, url, data, opts...)
		if err == nil {
			return nil
		}

		var shouldRetry bool
		respErr, shouldRetry = processError(l, resp, err)
		if !shouldRetry {
			return backoff.Permanent(err)
		}

		return respErr
	}, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
		var respErr *Error
		if errors.As(err, &respErr) && respErr.retryAfter.Seconds() > 0 {
			retryAfterDuration := respErr.retryAfter - duration
			if retryAfterDuration.Seconds() > 0 {
				l.Debug("sleeping for %v", retryAfterDuration)
				time.Sleep(retryAfterDuration)
			}
		}
	})
	if resp != nil {
		defer closeResp(resp)
	}
	if err != nil {
		l.Warn("failed client request err >%s< respErr >%s<", err, respErr)
		if respErr != nil {
			return resp, respErr
		}
		return resp, err
	}

	if resp != nil && respData != nil {
		err = c.decodeData(resp.Body, respData)
		if err != nil {
			respErr = &Error{
				Err: fmt.Errorf("failed decoding response >%v<", err),
			}
			l.Warn(respErr.Error())
			return resp, respErr
		}
	}

	return resp, nil
}

func processError(l logger.Logger, resp *http.Response, err error) (*Error, bool) {
	defer closeResp(resp)

	// With or without a response, assign the response error
	respErr := &Error{
		Err: err,
	}

	const shouldRetry = true

	if resp == nil {
		return respErr, shouldRetry
	}

	if resp.Body != nil {
		buf := new(bytes.Buffer)
		_, rerr := buf.ReadFrom(resp.Body)
		if rerr != nil {
			respErr.Err = fmt.Errorf("response error >%v< with read error >%v<", err, rerr)
			l.Warn(respErr.Error())
			return respErr, false
		}

		respErr.Body = buf.String()
		respErr.Err = fmt.Errorf("response error >%v<", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
		if s, ok := resp.Header["Retry-After"]; ok {
			if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
				respErr.retryAfter = time.Second * time.Duration(sleep)
			}
		}

		return respErr, shouldRetry
	}

	if resp.StatusCode == 0 || resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode != http.StatusNotImplemented {
		return respErr, shouldRetry
	}

	return respErr, false
}

// https://www.reddit.com/r/golang/comments/13fphyz/til_go_response_body_must_be_closed_even_if_you/
func closeResp(r *http.Response) {
	if r != nil && r.Body != nil {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}
}

// Request -
func (c *HTTPClient) Request(method, url string, data []byte, opts ...func(r *http.Request) *http.Request) (*http.Response, error) {
	l := loggerWithFunctionContext(c.Log, "Request")

	var err error

	c.log(l, "HTTPClient request host >%s< method >%s< URL >%s< data length >%d<", c.Host, method, url, len(data))

	var resp *http.Response
	var req *http.Request

	// Request + Response logging
	var requestDump []byte
	var responseDump []byte

	var body io.Reader
	if data != nil {
		body = bytes.NewBuffer(data)
	}

	req, err = http.NewRequest(method, url, body)
	if err != nil {
		l.Warn("failed new client request >%v<", err)
		return nil, err
	}

	for _, opt := range opts {
		if opt != nil {
			opt(req)
		}
	}

	err = c.SetHeaders(req)
	if err != nil {
		l.Warn("failed setting headers >%v<", err)
		return nil, err
	}

	err = c.SetAuthHeaders(req)
	if err != nil {
		l.Warn("failed setting request auth headers >%v<", err)
		return nil, err
	}

	switch method {
	case http.MethodGet:
		if c.RequestLogFunc != nil {
			requestDump, err = httputil.DumpRequest(req, true)
			if err != nil {
				l.Warn("failed request dump >%v<", err)
				return nil, err
			}
		}

		resp, err = c.client.Do(req)
		if err != nil {
			l.Warn("failed client request >%v<", err)
			return resp, err
		}

		if c.RequestLogFunc != nil {
			responseDump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				l.Warn("failed response dump >%v<", err)
				return nil, err
			}
			c.RequestLogFunc(url, string(requestDump), string(responseDump))
		}

	case http.MethodPost, http.MethodPut:
		req.Header.Add("Content-Type", "application/json")

		if c.RequestLogFunc != nil {
			requestDump, err = httputil.DumpRequest(req, true)
			if err != nil {
				l.Warn("failed request dump >%v<", err)
				return nil, err
			}
		}

		resp, err = c.client.Do(req)
		if err != nil {
			l.Warn("failed client request >%#v< >%v<", resp, err)
			return resp, err
		}

		if c.RequestLogFunc != nil {
			responseDump, err = httputil.DumpResponse(resp, true)
			if err != nil {
				l.Warn("failed response dump >%v<", err)
				return nil, err
			}
			c.RequestLogFunc(url, string(requestDump), string(responseDump))
		}

	case http.MethodDelete:
		req.Header.Add("Content-Type", "application/json")

		if c.RequestLogFunc != nil {
			requestDump, err = httputil.DumpRequest(req, true)
			if err != nil {
				l.Warn("failed request dump >%v<", err)
				return nil, err
			}
		}

		resp, err = c.client.Do(req)
		if err != nil {
			l.Warn("failed client request >%#v< >%v<", resp, err)
			return resp, err
		}

		if c.RequestLogFunc != nil {
			responseDump, err = httputil.DumpResponse(resp, true)
			if err != nil {
				l.Warn("failed response dump >%v<", err)
				return nil, err
			}
			c.RequestLogFunc(url, string(requestDump), string(responseDump))
		}
	default:
		// boom
		msg := fmt.Sprintf("method >%s< currently unsupported!", method)
		l.Warn(msg)
		return nil, errors.New(msg)
	}

	if resp == nil {
		return nil, errors.New("empty response")
	}

	c.log(l, "HTTPClient response status >%s<", resp.Status)

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusAccepted &&
		resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("response status >%d<", resp.StatusCode)
	}

	return resp, err
}

// SetHeaders sets various headers as per configuration
func (c *HTTPClient) SetHeaders(req *http.Request) error {
	l := loggerWithFunctionContext(c.Log, "SetHeaders")

	if c.RequestHeaders != nil {
		for key, value := range c.RequestHeaders {
			c.log(l, "Setting header >%s< .%s<", key, value)
			req.Header.Add(key, value)
		}
	}

	return nil
}

// SetAuthHeaders sets authentication headers on an request object based
// on client authentication configuration
func (c *HTTPClient) SetAuthHeaders(req *http.Request) error {

	// Authorization with key
	if c.AuthKey != "" {
		req.Header.Add("Authorization", c.AuthKey)
		return nil
	}

	// Authorization with bearer token
	if c.AuthToken != "" {
		var bearer = "Bearer " + c.AuthToken
		req.Header.Add("Authorization", bearer)
		return nil
	}

	// Authorization with user and pass
	if c.AuthUser != "" && c.AuthPass != "" {
		req.SetBasicAuth(c.AuthUser, c.AuthPass)
		return nil
	}

	return nil
}

// buildURL replaces placeholder parameters and adds query parameters
// The parameter "id" or ":id" has special behaviour. When provided the
// returned URL will have "/:id" appended and replaced with whatever
// the parameter value for "id" or ":id" was.
func (c *HTTPClient) buildURL(requestURL string, params map[string]string) (string, error) {
	l := loggerWithFunctionContext(c.Log, "buildURL")

	// Request URL
	requestURL = c.Host + c.Path + requestURL

	// Replace placeholders and add query parameters
	paramString := ""
	for param, value := range params {

		// do not allow empty param values
		if value == "" {
			return requestURL, fmt.Errorf("param >%s< has empty value", param)
		}

		found := false
		if strings.Contains(requestURL, "/:"+param) {
			requestURL = strings.Replace(requestURL, "/:"+param, "/"+value, 1)
			found = true
		} else if strings.Contains(requestURL, "/"+param) {
			requestURL = strings.Replace(requestURL, "/"+param, "/"+value, 1)
			found = true
		} else if strings.Contains(requestURL, "/"+strings.Replace(param, ":", "", 1)) {
			requestURL = strings.Replace(requestURL, "/"+strings.Replace(param, ":", "", 1), "/"+value, 1)
			found = true
		}
		if !found {
			param = strings.Replace(param, ":", "", 1)
			if paramString != "" {
				paramString += "&"
			}
			paramString = paramString + param + "=" + url.QueryEscape(value)
		}
	}

	if paramString != "" {
		requestURL = requestURL + "?" + paramString
	}

	// do not allow missing parameters
	if strings.Contains(requestURL, "/:") {
		return requestURL, fmt.Errorf("URL >%s< still contains placeholders", requestURL)
	}

	c.log(l, "Request URL >%s<", requestURL)

	return requestURL, nil
}

// RegisterRequestLogFunc -
func (c *HTTPClient) RegisterRequestLogFunc(logFunc func(url, request, response string)) {
	c.RequestLogFunc = logFunc
}

// encodeData is a convenience function that encodes struct data into bytes
func (c *HTTPClient) encodeData(data interface{}) ([]byte, error) {
	l := loggerWithFunctionContext(c.Log, "encodeData")

	dataBytes, err := json.Marshal(data)
	if err != nil {
		l.Warn("failed encoding data >%v<", err)
		return nil, err
	}
	return dataBytes, nil
}

// decodeData is a convenience function that decodes bytes into struct data
func (c *HTTPClient) decodeData(rc io.ReadCloser, data interface{}) error {
	l := loggerWithFunctionContext(c.Log, "decodeData")

	if dataPtr, ok := data.(*[]byte); ok {
		read, err := io.ReadAll(rc)
		if err != nil {
			l.Warn("failed reading data into buffer >%v<", err)
			return err
		}

		*dataPtr = read
	} else {
		err := json.NewDecoder(rc).Decode(&data)
		if err != nil && err.Error() != "EOF" {
			l.Warn("failed decoding data >%v<", err)
			return err
		}
	}

	return nil
}

func (c *HTTPClient) log(l logger.Logger, msg string, args ...interface{}) {
	if c.Verbose {
		l.Info(msg, args...)
	} else {
		l.Debug(msg, args...)
	}
}

// loggerWithFunctionContext - Returns a logger with package context and provided function context
func loggerWithFunctionContext(l logger.Logger, functionName string) logger.Logger {
	return l.WithPackageContext("client").WithFunctionContext(functionName)
}
