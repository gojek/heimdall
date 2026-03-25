package hystrix

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/internal"
	"github.com/pkg/errors"
)

type fallbackFunc func(error) error

// Client is the hystrix client implementation
type Client struct {
	client *httpclient.Client

	timeout                time.Duration
	hystrixTimeout         time.Duration
	hystrixCommandName     string
	maxConcurrentRequests  int
	requestVolumeThreshold int
	sleepWindow            int
	errorPercentThreshold  int
	retryCount             int
	retrier                heimdall.Retriable
	fallbackFunc           func(ctx context.Context, err error) error
	statsD                 *plugins.StatsdCollectorConfig
}

const (
	defaultHystrixRetryCount      = 0
	defaultHTTPTimeout            = 30 * time.Second
	defaultHystrixTimeout         = 30 * time.Second
	defaultMaxConcurrentRequests  = 100
	defaultErrorPercentThreshold  = 25
	defaultSleepWindow            = 10
	defaultRequestVolumeThreshold = 10

	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
)

var _ heimdall.Client = (*Client)(nil)
var err5xx = errors.New("server returned 5xx status code")

// NewClient returns a new instance of hystrix Client
func NewClient(opts ...Option) *Client {
	client := Client{
		client:                 httpclient.NewClient(),
		timeout:                defaultHTTPTimeout,
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

	if client.statsD != nil {
		c, err := plugins.InitializeStatsdCollector(client.statsD)
		if err != nil {
			panic(err)
		}

		metricCollector.Registry.Register(c.NewStatsdCollector)
	}

	hystrix.ConfigureCommand(client.hystrixCommandName, hystrix.CommandConfig{
		Timeout:                durationToInt(client.hystrixTimeout, time.Millisecond),
		MaxConcurrentRequests:  client.maxConcurrentRequests,
		RequestVolumeThreshold: client.requestVolumeThreshold,
		SleepWindow:            client.sleepWindow,
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
	if origReqBody := request.Body; origReqBody != nil {
		defer func() {
			// close the original request body as internal.BuildReadSeekCloser wraps body with noop closer.
			_ = origReqBody.Close()
		}()
	}

	var reqBody io.ReadSeekCloser
	var err error
	if hhc.retryCount > 0 && request.Body != nil {
		reqBody, err = internal.BuildReadSeekCloser(request.Body)
		if err != nil {
			return nil, err
		}
		request.Body = reqBody
	}

	var response *http.Response
	for i := 0; i <= hhc.retryCount; i++ {
		if response != nil {
			_, _ = io.Copy(io.Discard, response.Body)
			_ = response.Body.Close()
		}

		if i > 0 {
			time.Sleep(hhc.retrier.NextInterval(i - 1)) // sleep after closing the previous response body
		}

		if reqBody != nil {
			// Reset the body reader before the request to ensure even in case of retries from hystrix timeout, the body is at correct offset.
			// Note that it's safe to ignore the error here since the 0,0 position is always valid
			_, _ = reqBody.Seek(0, 0)
		}
		response, err = hhc.hystrixDo(request)
		if err == nil {
			break
		}
	}

	if err != nil {
		if err == err5xx {
			return response, nil
		}

		return nil, err
	}

	return response, nil
}

func (hhc *Client) hystrixDo(request *http.Request) (response *http.Response, err error) {
	ctx, cancelFn := context.WithCancel(request.Context())
	defer cancelFn()

	err = hystrix.DoC(ctx, hhc.hystrixCommandName, func(_ context.Context) (err error) {
		resp, err := hhc.client.Do(request.WithContext(ctx))
		if err != nil {
			return err
		}
		response = resp

		if response.StatusCode >= http.StatusInternalServerError {
			return err5xx
		}

		return nil
	}, hhc.fallbackFunc)

	return response, err
}

// AddPlugin Adds plugin to client
func (hhc *Client) AddPlugin(p heimdall.Plugin) {
	hhc.client.AddPlugin(p)
}
