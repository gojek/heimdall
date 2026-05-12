package httpclient

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/stretchr/testify/assert"
)

func TestOptionsAreSet(t *testing.T) {
	t.Parallel()

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
		WithRetryableStatusCodes(400, 201, 424),
	)

	assert.Equal(t, client, c.client)
	assert.NotEqual(t, httpTimeout, client.client.Timeout) // can't override custom implementation
	assert.Equal(t, httpTimeout, *c.timeout)
	assert.Equal(t, retrier, c.retrier)
	assert.Equal(t, noOfRetries, c.retryCount)
	assert.Equal(t, []int{201, 400, 424}, c.retryableCodes)
}

func TestWithClientWihhoutHTTPTimeoutShouldNotOverrideUserHTTPClientTimeout(t *testing.T) {
	t.Parallel()

	client := &http.Client{Timeout: 25 * time.Millisecond}

	c := NewClient(
		WithHTTPClient(client),
	)

	assert.Equal(t, client, c.client)
	assert.Equal(t, 25*time.Millisecond, client.Timeout) // overrides user provided *http.Client
	assert.Nil(t, c.timeout)
}

func TestWithHTTPTimeoutOverridesUserHTTPClientTimeout(t *testing.T) {
	t.Parallel()

	httpTimeout := 10 * time.Second

	client := &http.Client{Timeout: 25 * time.Millisecond}

	c := NewClient(
		WithHTTPClient(client),
		WithHTTPTimeout(httpTimeout),
	)

	assert.Equal(t, client, c.client)
	assert.Equal(t, httpTimeout, client.Timeout) // overrides user provided *http.Client
	assert.Equal(t, httpTimeout, *c.timeout)
}

func TestWithHTTPTimeoutOverridesUserHTTPClientTimeout_InverseSeq(t *testing.T) {
	t.Parallel()

	httpTimeout := 10 * time.Second

	client := &http.Client{Timeout: 25 * time.Millisecond}

	c := NewClient(
		WithHTTPTimeout(httpTimeout),
		WithHTTPClient(client),
	)

	assert.Equal(t, client, c.client)
	assert.Equal(t, httpTimeout, client.Timeout) // overrides user provided *http.Client
	assert.Equal(t, httpTimeout, *c.timeout)
}

func TestOptionsHaveDefaults(t *testing.T) {
	t.Parallel()

	retrier := heimdall.NewNoRetrier()
	httpTimeout := 30 * time.Second
	http.DefaultClient.Timeout = httpTimeout
	noOfRetries := 0

	c := NewClient()

	assert.Equal(t, http.DefaultClient, c.client)
	assert.Nil(t, c.timeout)
	httpClient, ok := c.client.(*http.Client)
	assert.True(t, ok)
	assert.Equal(t, httpTimeout, httpClient.Timeout)
	assert.Equal(t, retrier, c.retrier)
	assert.Equal(t, noOfRetries, c.retryCount)
}

func ExampleWithHTTPTimeout() {
	c := NewClient(WithHTTPTimeout(5 * time.Second))
	req, err := http.NewRequest(http.MethodGet, "https://gojek.com/", nil)
	if err != nil {
		panic(err)
	}
	res, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response status : ", res.StatusCode)
	// Output: Response status :  200
}

func ExampleWithHTTPTimeout_expired() {
	c := NewClient(WithHTTPTimeout(1 * time.Millisecond))
	req, err := http.NewRequest(http.MethodGet, "https://gojek.com/", nil)
	if err != nil {
		panic(err)
	}
	res, err := c.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Response status : ", res.StatusCode)
}

func ExampleWithRetryCount() {
	c := NewClient(WithHTTPTimeout(1*time.Millisecond), WithRetryCount(3))
	req, err := http.NewRequest(http.MethodGet, "https://gojek.com/", nil)
	if err != nil {
		panic(err)
	}
	res, err := c.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Response status : ", res.StatusCode)
}

type mockClient struct{}

func (m *mockClient) Do(r *http.Request) (*http.Response, error) {
	fmt.Println("mock client called")
	return &http.Response{}, nil
}

func ExampleWithHTTPClient() {
	m := &mockClient{}
	c := NewClient(WithHTTPClient(m))
	req, err := http.NewRequest(http.MethodGet, "https://gojek.com/", nil)
	if err != nil {
		panic(err)
	}
	_, _ = c.Do(req)
	// Output: mock client called
}

type mockRetrier struct{}

func (m *mockRetrier) NextInterval(attempt int) time.Duration {
	fmt.Println("retry attempt", attempt)
	return time.Millisecond
}

func ExampleWithRetrier() {
	c := NewClient(WithHTTPTimeout(1*time.Millisecond), WithRetryCount(3), WithRetrier(&mockRetrier{}))
	req, err := http.NewRequest(http.MethodGet, "https://gojek.com/", nil)
	if err != nil {
		panic(err)
	}
	res, err := c.Do(req)
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println("Response status : ", res.StatusCode)
	// Output: retry attempt 0
	// retry attempt 1
	// retry attempt 2
	// error
}
