package hystrix

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gojektech/heimdall"
	"github.com/pkg/errors"
)

type fallbackFunc func(error) error

type Client struct {
	client heimdall.Doer

	timeout                time.Duration
	hystrixTimeout         time.Duration
	hystrixCommandName     string
	maxConcurrentRequests  int
	requestVolumeThreshold int
	sleepWindow            int
	errorPercentThreshold  int
	retryCount             int
	retrier                heimdall.Retriable
	fallbackFunc           func(err error) error
}

const defaultHystrixRetryCount int = 0

var _ heimdall.Client = (*Client)(nil)

// NewClient returns a new instance of Client
func NewClient(opts ...Option) *Client {
	client := Client{
		retryCount: defaultHystrixRetryCount,
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

	hystrix.ConfigureCommand(client.hystrixCommandName, hystrix.CommandConfig{
		Timeout:                int(client.hystrixTimeout),
		MaxConcurrentRequests:  client.maxConcurrentRequests,
		RequestVolumeThreshold: client.requestVolumeThreshold,
		SleepWindow:            client.sleepWindow,
		ErrorPercentThreshold:  client.errorPercentThreshold,
	})

	return &client
}

// Get makes a HTTP GET request to provided URL
func (hhc *Client) Get(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}

	request.Header = headers

	return hhc.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (hhc *Client) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response, errors.Wrap(err, "POST - request creation failed")
	}

	request.Header = headers

	return hhc.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (hhc *Client) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PUT - request creation failed")
	}

	request.Header = headers

	return hhc.Do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (hhc *Client) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PATCH - request creation failed")
	}

	request.Header = headers

	return hhc.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (hhc *Client) Delete(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "DELETE - request creation failed")
	}

	request.Header = headers

	return hhc.Do(request)
}

// Do makes an HTTP request with the native `http.Do` interface
func (hhc *Client) Do(request *http.Request) (*http.Response, error) {
	request.Close = true

	var response *http.Response
	var err error

	for i := 0; i <= hhc.retryCount; i++ {
		err = hystrix.Do(hhc.hystrixCommandName, func() error {
			response, err = hhc.client.Do(request)
			if err != nil {
				return err
			}

			if response.StatusCode >= http.StatusInternalServerError {
				return fmt.Errorf("Server is down: returned status code: %d", response.StatusCode)
			}
			return nil
		}, hhc.fallbackFunc)

		if err != nil {
			backoffTime := hhc.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		break
	}

	return response, err
}
