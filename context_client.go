package heimdall

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gojektech/valkyrie"
	"github.com/pkg/errors"
)

type httpClientWithContext struct {
	client *http.Client

	retryCount int
	retrier    Retriable
}

// NewHTTPClientWithContext returns a new instance of httpClientWithContext
func NewHTTPClientWithContext(timeout time.Duration) ClientWithContext {
	return &httpClientWithContext{
		client: &http.Client{
			Timeout: timeout,
		},

		retryCount: defaultRetryCount,
		retrier:    NewNoRetrier(),
	}
}

// SetRetryCount sets the retry count for the httpClient
func (c *httpClientWithContext) SetRetryCount(count int) {
	c.retryCount = count
}

// SetRetryCount sets the retry count for the httpClient
func (c *httpClientWithContext) SetRetrier(retrier Retriable) {
	c.retrier = retrier
}

// Get makes a HTTP GET request to provided URL with context passed in
func (c *httpClientWithContext) Get(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}

	request.Header = headers

	return c.Do(ctx, request)
}

// Post makes a HTTP POST request to provided URL with context passed in
func (c *httpClientWithContext) Post(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response, errors.Wrap(err, "POST - request creation failed")
	}

	request.Header = headers

	return c.Do(ctx, request)
}

// Put makes a HTTP PUT request to provided URL with context passed in
func (c *httpClientWithContext) Put(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PUT - request creation failed")
	}

	request.Header = headers

	return c.Do(ctx, request)
}

// Patch makes a HTTP PATCH request to provided URL with context passed in
func (c *httpClientWithContext) Patch(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PATCH - request creation failed")
	}

	request.Header = headers

	return c.Do(ctx, request)
}

// Delete makes a HTTP DELETE request to provided URL with context passed in
func (c *httpClientWithContext) Delete(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "DELETE - request creation failed")
	}

	request.Header = headers

	return c.Do(ctx, request)
}

// Do makes an HTTP request with the native `http.Do` interface and context passed in
func (c *httpClientWithContext) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var response *http.Response

	multiErr := &valkyrie.MultiError{}

	for i := 0; i <= c.retryCount; i++ {
		contextCancelled := false
		var err error
		response, err = c.client.Do(req.WithContext(ctx))
		if err != nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				contextCancelled = true
			}

			multiErr.Push(err.Error())
			if contextCancelled {
				break
			}
			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		if response.StatusCode >= http.StatusInternalServerError {
			multiErr.Push(fmt.Sprintf("server error: %d", response.StatusCode))

			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			fmt.Println("R: ", response.StatusCode)
			continue
		}

		multiErr = &valkyrie.MultiError{} // Clear errors if any iteration succeeds
		break
	}

	return response, multiErr.HasError()
}
