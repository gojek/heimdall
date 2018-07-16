package httpclient

import (
	"net/http"
	"testing"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/stretchr/testify/assert"
)

func TestOptionsAreSet(t *testing.T) {
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond
	noOfRetries := 3
	httpTimeout := 10 * time.Second

	client := &myHTTPClient{client: http.Client{Timeout: 25 * time.Millisecond}}
	retrier := heimdall.NewRetrier(heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval))

	c := NewClient(
		WithHTTPClient(client),
		WithHTTPTimeout(httpTimeout),
		WithRetrier(retrier),
		WithRetryCount(noOfRetries),
	)

	assert.Equal(t, client, c.client)
	assert.Equal(t, httpTimeout, c.timeout)
	assert.Equal(t, retrier, c.retrier)
	assert.Equal(t, noOfRetries, c.retryCount)
}

func TestOptionsHaveDefaults(t *testing.T) {
	retrier := heimdall.NewNoRetrier()
	httpTimeout := 30 * time.Second
	http.DefaultClient.Timeout = httpTimeout
	noOfRetries := 0

	c := NewClient()

	assert.Equal(t, http.DefaultClient, c.client)
	assert.Equal(t, httpTimeout, c.timeout)
	assert.Equal(t, retrier, c.retrier)
	assert.Equal(t, noOfRetries, c.retryCount)
}
