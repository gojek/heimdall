package hystrix

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/gojek/heimdall/v8"
	"github.com/gojek/heimdall/v8/httpclient"
	"github.com/gojek/heimdall/v8/internal"
	"github.com/gojek/hystrix-go/hystrix"
)

type fallbackFunc func(error) error
type fallbackCtxFunc func(context.Context, error) error

// Client is the hystrix client implementation
type Client struct {
	client *httpclient.Client

	hystrixTimeout         time.Duration
	hystrixCommandName     string
	maxConcurrentRequests  int
	requestVolumeThreshold int
	sleepWindow            time.Duration
	errorPercentThreshold  int
	fallbackFunc           func(ctx context.Context, err error) error

	retrier          heimdall.Retriable
	retryCount       int
	retryableCodes   []int
	retryErrorBudget *internal.ErrorBudget
}

const (
	defaultHystrixRetryCount      = 0
	defaultHystrixTimeout         = 30 * time.Second
	defaultMaxConcurrentRequests  = 100
	defaultErrorPercentThreshold  = 25
	defaultSleepWindow            = 10 * time.Millisecond
	defaultRequestVolumeThreshold = 10

	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
)

var _ heimdall.Client = (*Client)(nil)
var errRetryableCode = errors.New("server returned status code to retry")

// NewClient returns a new instance of hystrix Client
func NewClient(opts ...Option) *Client {
	client := Client{
		client:                 httpclient.NewClient(),
		hystrixTimeout:         defaultHystrixTimeout,
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

	hystrix.ConfigureCommand(client.hystrixCommandName, hystrix.CommandConfig{
		Timeout:                durationToInt(client.hystrixTimeout, time.Millisecond),
		MaxConcurrentRequests:  client.maxConcurrentRequests,
		RequestVolumeThreshold: client.requestVolumeThreshold,
		SleepWindow:            durationToInt(client.sleepWindow, time.Millisecond),
		ErrorPercentThreshold:  client.errorPercentThreshold,
	})

	return &client
}

func durationToInt(duration, unit time.Duration) int {
	durationAsNumber := duration / unit

	if int64(durationAsNumber) > int64(maxInt) {
		// Returning max possible value seems like best possible solution here
		// the alternative is to panic as there is no way of returning an error
		// without changing the NewClient API
		return maxInt
	}
	return int(durationAsNumber)
}

// Get makes a HTTP GET request to provided URL
func (hhc *Client) Get(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, fmt.Errorf("GET %s - request creation failed: %w", hhc.hystrixCommandName, err)
	}

	request.Header = headers

	return hhc.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (hhc *Client) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response, fmt.Errorf("POST %s - request creation failed: %w", hhc.hystrixCommandName, err)
	}

	request.Header = headers

	return hhc.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (hhc *Client) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return response, fmt.Errorf("PUT %s - request creation failed: %w", hhc.hystrixCommandName, err)
	}

	request.Header = headers

	return hhc.Do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (hhc *Client) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return response, fmt.Errorf("PATCH %s - request creation failed: %w", hhc.hystrixCommandName, err)
	}

	request.Header = headers

	return hhc.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (hhc *Client) Delete(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, fmt.Errorf("DELETE %s - request creation failed: %w", hhc.hystrixCommandName, err)
	}

	request.Header = headers

	return hhc.Do(request)
}

// Do makes an HTTP request with the native `http.Do` interface
func (hhc *Client) Do(request *http.Request) (*http.Response, error) {
	if origReqBody := request.Body; origReqBody != nil {
		defer func() {
			// close the original request body as internal.SetRequestGetBody wraps body with noop closer.
			_ = origReqBody.Close()
		}()
	}

	var reqGetBody internal.RequestGetBody
	var err error
	// Only SetRequestGetBody if retry is enabled to avoid unnecessary overhead for non-retry requests
	if hhc.retryCount > 0 {
		if err := internal.SetRequestGetBody(request); err != nil {
			return nil, err
		}
		// keeping a local variable just in case request.GetBody gets overridden by some plugins/middlewares
		reqGetBody = request.GetBody
	}

	var response *http.Response
	for i := 0; i <= hhc.retryCount; i++ {
		if response != nil {
			_, _ = io.Copy(io.Discard, response.Body)
			_ = response.Body.Close()
		}

		if i > 0 {
			err = internal.SleepInterruptible(request.Context(), hhc.retrier.NextInterval(i-1))
			if err != nil {
				return nil, err
			}

			request, err = internal.CloneRequest(request, reqGetBody) // Clone the request to reset the body for retry
			if err != nil {
				return nil, err
			}
		}

		response, err = hhc.hystrixDo(request)
		if err == nil || internal.IsCtxDone(request.Context()) {
			_ = hhc.retryErrorBudget.Success()
			break
		}

		if hhc.retryErrorBudget.Failure() {
			break
		}
	}

	if err != nil {
		if errors.Is(err, errRetryableCode) {
			return response, nil
		}

		return nil, err
	}

	return response, nil
}

func (hhc *Client) hystrixDo(request *http.Request) (*http.Response, error) {
	var response *http.Response
	err := hystrix.DoC(request.Context(), hhc.hystrixCommandName, func(_ context.Context) error {
		resp, doErr := hhc.client.Do(request)
		if doErr != nil {
			return doErr
		}
		response = resp

		if _, ok := slices.BinarySearch(hhc.retryableCodes, response.StatusCode); ok ||
			response.StatusCode >= http.StatusInternalServerError {
			return errRetryableCode
		}

		return nil
	}, hhc.fallbackFunc)
	if err != nil && !errors.Is(err, errRetryableCode) { // Special handling to avoid data race conditions
		return nil, err
	}

	return response, err
}

// AddPlugin Adds plugin to client
func (hhc *Client) AddPlugin(p heimdall.Plugin) {
	hhc.client.AddPlugin(p)
}
