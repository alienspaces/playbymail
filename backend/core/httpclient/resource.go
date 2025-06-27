package client

import (
	"errors"
	"fmt"
	"net/http"
)

// Get is a convenience method wrapping RetryRequest
func (c *HTTPClient) Get(path string, params map[string]string, respData interface{}) error {

	c.Log.Context("function", "Get")
	defer func() {
		c.Log.Context("function", "")
	}()

	c.Log.Debug("Request path >%s< params >%#v< respData >%#v<", path, params, respData)

	_, err := c.RetryRequest(
		http.MethodGet,
		path,
		params,
		nil,
		respData,
	)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("failed request >%v<", err))
		return err
	}

	return nil
}

// Create is a convenience method wrapping RetryRequest
func (c *HTTPClient) Create(path string, params map[string]string, reqData interface{}, respData interface{}) error {

	c.Log.Context("function", "Create")
	defer func() {
		c.Log.Context("function", "")
	}()

	c.Log.Debug("Request path >%s< params >%#v< reqData >%#v< respData >%#v<", path, params, reqData, respData)

	if reqData == nil {
		msg := fmt.Sprintf("request data is nil >%v<, cannot create resource", reqData)
		c.Log.Warn(msg)
		return errors.New(msg)
	}

	_, err := c.RetryRequest(
		http.MethodPost,
		path,
		params,
		reqData,
		respData,
	)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("failed request >%v<", err))
		return err
	}

	return nil
}

// Update is a convenience method wrapping RetryRequest
func (c *HTTPClient) Update(path string, params map[string]string, reqData interface{}, respData interface{}) error {

	c.Log.Context("function", "Update")
	defer func() {
		c.Log.Context("function", "")
	}()

	c.Log.Debug("Request path >%s< params >%#v< reqData >%#v< respData >%#v<", path, params, reqData, respData)

	if reqData == nil {
		msg := fmt.Sprintf("request data is nil >%v<, cannot update resource", reqData)
		c.Log.Warn(msg)
		return errors.New(msg)
	}

	_, err := c.RetryRequest(
		http.MethodPut,
		path,
		params,
		reqData,
		respData,
	)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("failed request >%v<", err))
		return err
	}

	return nil
}

// Delete is a convenience method wrapping RetryRequest
func (c *HTTPClient) Delete(path string, params map[string]string, respData interface{}) error {

	c.Log.Context("function", "Delete")
	defer func() {
		c.Log.Context("function", "")
	}()

	_, err := c.RetryRequest(
		http.MethodDelete,
		path,
		params,
		nil,
		respData,
	)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("failed request >%v<", err))
		return err
	}

	return nil
}
