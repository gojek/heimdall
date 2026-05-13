package httpclient

import (
	"slices"
	"time"

	"github.com/gojek/heimdall/v7"
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
