package httpclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/valkyrie"
	"github.com/pkg/errors"
)

// Client is the http client implementation
type Client struct {
	client heimdall.Doer

	timeout    time.Duration
	retryCount int
	retrier    heimdall.Retriable
	plugins    []heimdall.Plugin
}

const (
	defaultRetryCount  = 0
	defaultHTTPTimeout = 30 * time.Second
)

var _ heimdall.Client = (*Client)(nil)

// NewClient returns a new instance of http Client
func NewClient(opts ...Option) *Client {
	client := Client{
		timeout:    defaultHTTPTimeout,
		retryCount: defaultRetryCount,
		retrier:    heimdall.NewNoRetrier(),
	}

	for _, opt := range opts {
		opt(&client)
	}

	if client.client == nil {
		client.client = &http.Client{
			Timeout: client.timeout,
		}
	}

	return &client
}

// AddPlugin Adds plugin to client
func (c *Client) AddPlugin(p heimdall.Plugin) {
	c.plugins = append(c.plugins, p)
}

// Get makes a HTTP GET request to provided URL
func (c *Client) Get(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (c *Client) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response, errors.Wrap(err, "POST - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (c *Client) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PUT - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (c *Client) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PATCH - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (c *Client) Delete(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "DELETE - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Do makes an HTTP request with the native `http.Do` interface
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	request.Close = true

	var bodyReader *bytes.Reader

	if request.Body != nil {
		reqData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqData)
		request.Body = ioutil.NopCloser(bodyReader) // prevents closing the body between retries
	}

	multiErr := &valkyrie.MultiError{}
	var response *http.Response

	for i := 0; i <= c.retryCount; i++ {
		if response != nil {
			response.Body.Close()
		}

		c.reportRequestStart(request)
		var err error
		response, err = c.client.Do(request)
		if bodyReader != nil {
			// Reset the body reader after the request since at this point it's already read
			// Note that it's safe to ignore the error here since the 0,0 position is always valid
			_, _ = bodyReader.Seek(0, 0)
		}

		if err != nil {
			multiErr.Push(err.Error())
			c.reportError(request, err)
			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}
		c.reportRequestEnd(request, response)

		if response.StatusCode >= http.StatusInternalServerError {
			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		multiErr = &valkyrie.MultiError{} // Clear errors if any iteration succeeds
		break
	}

	return response, multiErr.HasError()
}

func (c *Client) reportRequestStart(request *http.Request) {
	for _, plugin := range c.plugins {
		plugin.OnRequestStart(request)
	}
}

func (c *Client) reportError(request *http.Request, err error) {
	for _, plugin := range c.plugins {
		plugin.OnError(request, err)
	}
}

func (c *Client) reportRequestEnd(request *http.Request, response *http.Response) {
	for _, plugin := range c.plugins {
		plugin.OnRequestEnd(request, response)
	}
}
