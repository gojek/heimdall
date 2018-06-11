package heimdall

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gojektech/valkyrie"

	"github.com/pkg/errors"
)

const defaultHystrixRetryCount int = 0

type hystrixHTTPClient struct {
	client             Doer
	hystrixCommandName string
	retryCount         int
	retrier            Retriable
	fallbackFunc       func(err error) error
}

// NewHystrixHTTPClient returns a new instance of HystrixHTTPClient
func NewHystrixHTTPClient(timeout time.Duration, hystrixConfig HystrixConfig) Client {
	httpClient := &http.Client{
		Timeout: timeout,
	}

	hystrix.ConfigureCommand(hystrixConfig.commandName, hystrixConfig.commandConfig)

	return &hystrixHTTPClient{
		client: httpClient,

		retryCount:         defaultHystrixRetryCount,
		retrier:            NewNoRetrier(),
		hystrixCommandName: hystrixConfig.commandName,
		fallbackFunc:       hystrixConfig.fallbackFunc,
	}
}

// SetRetryCount sets the retry count for the hystrixHTTPClient
func (hhc *hystrixHTTPClient) SetRetryCount(count int) {
	hhc.retryCount = count
}

// SetRetrier sets the strategy for retrying
func (hhc *hystrixHTTPClient) SetRetrier(retrier Retriable) {
	hhc.retrier = retrier
}

// SetCustomHTTPClient sets custom HTTP client
func (hhc *hystrixHTTPClient) SetCustomHTTPClient(customHTTPClient Doer) {
	hhc.client = customHTTPClient
}

// Get makes a HTTP GET request to provided URL
func (hhc *hystrixHTTPClient) Get(url string, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "GET - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (hhc *hystrixHTTPClient) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "POST - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (hhc *hystrixHTTPClient) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "PUT - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (hhc *hystrixHTTPClient) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "PATCH - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (hhc *hystrixHTTPClient) Delete(url string, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DELETE - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Do makes an HTTP request with the native `http.Do` interface
func (hhc *hystrixHTTPClient) Do(request *http.Request) (*http.Response, error) {
	request.Close = true

	multiErr := &valkyrie.MultiError{}
	for i := 0; i <= hhc.retryCount; i++ {
		var response *http.Response
		err := hystrix.Do(hhc.hystrixCommandName, func() error {
			var err error
			response, err = hhc.client.Do(request)
			if err != nil {
				return err
			}
			if response.StatusCode >= http.StatusInternalServerError {
				return fmt.Errorf("server error: %d", response.StatusCode)
			}
			return nil
		}, hhc.fallbackFunc)

		if err != nil {
			multiErr.Push(err.Error())

			backoffTime := hhc.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}
		return response, nil
	}
	return nil, multiErr
}
