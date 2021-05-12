package hystrix

import (
	"time"

	"github.com/afex/hystrix-go/plugins"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
)

// Option represents the hystrix client options
type Option func(*Client)

// WithCommandName sets the hystrix command name
func WithCommandName(name string) Option {
	return func(c *Client) {
		c.hystrixCommandName = name
	}
}

// WithHTTPTimeout sets hystrix timeout
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithHystrixTimeout sets hystrix timeout
func WithHystrixTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.hystrixTimeout = timeout
	}
}

// WithMaxConcurrentRequests sets hystrix max concurrent requests
func WithMaxConcurrentRequests(maxConcurrentRequests int) Option {
	return func(c *Client) {
		c.maxConcurrentRequests = maxConcurrentRequests
	}
}

// WithRequestVolumeThreshold sets hystrix request volume threshold
func WithRequestVolumeThreshold(requestVolumeThreshold int) Option {
	return func(c *Client) {
		c.requestVolumeThreshold = requestVolumeThreshold
	}
}

// WithSleepWindow sets hystrix sleep window
func WithSleepWindow(sleepWindow int) Option {
	return func(c *Client) {
		c.sleepWindow = sleepWindow
	}
}

// WithErrorPercentThreshold sets hystrix error percent threshold
func WithErrorPercentThreshold(errorPercentThreshold int) Option {
	return func(c *Client) {
		c.errorPercentThreshold = errorPercentThreshold
	}
}

// WithFallbackFunc sets the fallback function
func WithFallbackFunc(fn fallbackFunc) Option {
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
