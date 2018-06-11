package heimdall

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"bytes"
	"io/ioutil"

	"github.com/gojektech/valkyrie"
	"github.com/pkg/errors"
)

const defaultRetryCount int = 0

type httpClient struct {
	client Doer

	retryCount int
	retrier    Retriable
}

// NewHTTPClient returns a new instance of HTTPClient
func NewHTTPClient(timeout time.Duration) Client {
	return &httpClient{
		client: &http.Client{
			Timeout: timeout,
		},

		retryCount: defaultRetryCount,
		retrier:    NewNoRetrier(),
	}
}

// SetRetryCount sets the retry count for the httpClient
func (c *httpClient) SetRetryCount(count int) {
	c.retryCount = count
}

// SetCustomHTTPClient sets custom HTTP client
func (c *httpClient) SetCustomHTTPClient(customHTTPClient Doer) {
	c.client = customHTTPClient
}

// SetRetrier sets the strategy for retrying
func (c *httpClient) SetRetrier(retrier Retriable) {
	c.retrier = retrier
}

// Get makes a HTTP GET request to provided URL
func (c *httpClient) Get(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (c *httpClient) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response, errors.Wrap(err, "POST - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (c *httpClient) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PUT - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (c *httpClient) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PATCH - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (c *httpClient) Delete(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "DELETE - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Do makes an HTTP request with the native `http.Do` interface
func (c *httpClient) Do(request *http.Request) (*http.Response, error) {
	request.Close = true

	var reqBuffer []byte

	if request != nil && request.Body != nil {
		var err error

		// Storing request buffer to create new reader on each request
		reqBuffer, err = ioutil.ReadAll(request.Body)

		if err != nil {
			return nil, err
		}
	}
	multiErr := &valkyrie.MultiError{}
	var response *http.Response

	for i := 0; i <= c.retryCount; i++ {
		var err error
		request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBuffer))
		response, err = c.client.Do(request)
		if err != nil {
			multiErr.Push(err.Error())

			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		if response.StatusCode >= http.StatusInternalServerError {
			multiErr.Push(fmt.Sprintf("server error: %d", response.StatusCode))

			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		multiErr = &valkyrie.MultiError{} // Clear errors if any iteration succeeds
		break
	}

	return response, multiErr.HasError()
}
