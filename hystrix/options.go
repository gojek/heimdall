package hystrix

import (
	"context"
	"slices"
	"time"

	"github.com/afex/hystrix-go/plugins"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/internal"
)

// Option represents the hystrix client options
type Option func(*Client)

// WithCommandName sets the hystrix command name
func WithCommandName(name string) Option {
	return func(c *Client) {
		c.hystrixCommandName = name
	}
}

// WithHTTPTimeout sets timeout for http.Client
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		httpclient.WithHTTPTimeout(timeout)(c.client)
	}
}

// WithHystrixTimeout sets hystrix timeout
func WithHystrixTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.hystrixTimeout = timeout
	}
}

// WithMaxConcurrentRequests sets the maximum number of concurrent requests
// allowed to the command before further requests are short-circuited.
func WithMaxConcurrentRequests(maxConcurrentRequests int) Option {
	return func(c *Client) {
		c.maxConcurrentRequests = maxConcurrentRequests
	}
}

// WithRequestVolumeThreshold sets the minimum number of requests in a rolling
// window needed before the circuit breaker can trip.
func WithRequestVolumeThreshold(requestVolumeThreshold int) Option {
	return func(c *Client) {
		c.requestVolumeThreshold = requestVolumeThreshold
	}
}

// WithSleepWindow sets how long, in milliseconds, to wait after the circuit
// opens before allowing a single retry to test whether the backend has recovered.
// The value is passed through to hystrix unchanged (hystrix defaults to 5000ms).
func WithSleepWindow(sleepWindow int) Option {
	return func(c *Client) {
		c.sleepWindow = sleepWindow
	}
}

// WithErrorPercentThreshold sets the error percentage (0-100) at or above which,
// once the request volume threshold is met, the circuit breaker opens.
func WithErrorPercentThreshold(errorPercentThreshold int) Option {
	return func(c *Client) {
		c.errorPercentThreshold = errorPercentThreshold
	}
}

// WithFallbackFunc sets the fallback function
func WithFallbackFunc(fn fallbackFunc) Option {
	if fn == nil {
		return func(_ *Client) {}
	}
	return WithFallbackCtxFunc(func(_ context.Context, err error) error {
		return fn(err)
	})
}

// WithFallbackCtxFunc sets the fallback function with context support
func WithFallbackCtxFunc(fn fallbackCtxFunc) Option {
	return func(c *Client) {
		c.fallbackFunc = fn
	}
}

// WithRetryCount sets the retry count for the Client
func WithRetryCount(retryCount int) Option {
	return func(c *Client) {
		c.retryCount = retryCount
	}
}

// WithRetrier sets the strategy for retrying
func WithRetrier(retrier heimdall.Retriable) Option {
	return func(c *Client) {
		c.retrier = retrier
	}
}

// WithHTTPClient sets a custom http client for hystrix client
func WithHTTPClient(client heimdall.Doer) Option {
	return func(c *Client) {
		opt := httpclient.WithHTTPClient(client)
		opt(c.client)
	}
}

// WithStatsDCollector exports hystrix metrics to a statsD backend
func WithStatsDCollector(addr, prefix string) Option {
	return func(c *Client) {
		c.statsD = &plugins.StatsdCollectorConfig{StatsdAddr: addr, Prefix: prefix}
	}
}

// WithRetryableStatusCodes sets status codes to be retried
// Note: All 5xx status codes are always eligible for retry, thus not required for WithRetryableStatusCodes option.
func WithRetryableStatusCodes(statusCodes ...int) Option {
	return func(c *Client) {
		codes := append(c.retryableCodes, statusCodes...)
		slices.Sort(codes)

		c.retryableCodes = codes
	}
}

// WithRetryErrorBudgetToken creates a weighted token retry error budget with the following token details.
//
//	maxToken: The maximum/initial token value which is used to calculate token threshold(i.e. maxToken/2)
//	tokenRatio: The allowed ratio of failure in comparison to success.
func WithRetryErrorBudgetToken(maxToken int32, tokenRatio float32) Option {
	return func(c *Client) {
		c.retryErrorBudget = internal.NewTokenErrorBudget(maxToken, tokenRatio)
	}
}

// WithRetryErrorBudgetPercent creates a weighted token retry error budget with the following failure details.
//
//	minFailureVolume: The minimum failure required.
//	failurePercent: The failure percentage (0-100).
//
// Note: To determine if budget is exceeded we use recent event which satisfies following
//
//	failureEvent <= maxFailureEvent
//	successEvent = (maxFailureEvent-failureEvent) / allowedSuccessPerFailure
//	totalEvent = failureEvent + successEvent
//
// Where
//
//	maxFailureEvent = minFailureVolume * 2
//	allowedSuccessPerFailure = (100 - failurePercent) / failurePercent
func WithRetryErrorBudgetPercent(minFailureVolume int32, failurePercent float32) Option {
	return func(c *Client) {
		c.retryErrorBudget = internal.NewPercentErrorBudget(minFailureVolume, failurePercent)
	}
}
