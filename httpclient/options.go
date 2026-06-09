package httpclient

import (
	"slices"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/internal"
)

// Option represents the client options
type Option func(*Client)

// WithHTTPTimeout sets timeout for http.Client
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = &timeout

		c.updateHTTPTimeout() // hystrix.WithHTTPTimeout relies on this
	}
}

// WithRetryCount sets the retry count for the hystrixHTTPClient
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

// WithHTTPClient sets a custom http client
func WithHTTPClient(client heimdall.Doer) Option {
	return func(c *Client) {
		c.client = client

		c.updateHTTPTimeout() // hystrix.WithHTTPTimeout relies on this
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
