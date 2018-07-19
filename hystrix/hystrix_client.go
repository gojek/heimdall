package hystrix

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gojektech/heimdall"
	"github.com/gojektech/valkyrie"
	"github.com/pkg/errors"
)

type fallbackFunc func(error) error

// Client is the hystrix client implementation
type Client struct {
	client                 heimdall.Doer
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

const (
	defaultHystrixRetryCount      = 0
	defaultHTTPTimeout            = 30 * time.Second
	defaultHystriximeout          = 30 * time.Second
	defaultMaxConcurrentRequests  = 100
	defaultErrorPercentThreshold  = 25
	defaultSleepWindow            = 10
	defaultRequestVolumeThreshold = 10
)

var _ heimdall.Client = (*Client)(nil)

// NewClient returns a new instance of hystrix Client
func NewClient(opts ...Option) *Client {
	client := Client{
		timeout:                defaultHTTPTimeout,
		hystrixTimeout:         defaultHystriximeout,
		maxConcurrentRequests:  defaultMaxConcurrentRequests,
		errorPercentThreshold:  defaultErrorPercentThreshold,
		sleepWindow:            defaultSleepWindow,
		requestVolumeThreshold: defaultRequestVolumeThreshold,
		retryCount:             defaultHystrixRetryCount,
		retrier:                heimdall.NewNoRetrier(),
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
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "GET - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (hhc *Client) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "POST - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (hhc *Client) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "PUT - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (hhc *Client) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "PATCH - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (hhc *Client) Delete(url string, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DELETE - request creation failed")
	}
	request.Header = headers
	return hhc.Do(request)
}

// Do makes an HTTP request with the native `http.Do` interface
func (hhc *Client) Do(request *http.Request) (*http.Response, error) {
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
