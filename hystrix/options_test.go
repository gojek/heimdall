package hystrix

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptionsAreSet(t *testing.T) {
	c := NewClient(
		WithHTTPTimeout(10*time.Second),
		WithCommandName("test"),
		WithHystrixTimeout(1100),
		WithMaxConcurrentRequests(10),
		WithErrorPercentThreshold(30),
		WithSleepWindow(5),
		WithRequestVolumeThreshold(5),
		WithStatsDCollector("localhost:8125", "myapp.hystrix"),
	)

	assert.Equal(t, 10*time.Second, c.timeout)
	assert.Equal(t, "test", c.hystrixCommandName)
	assert.Equal(t, time.Duration(1100), c.hystrixTimeout)
	assert.Equal(t, 10, c.maxConcurrentRequests)
	assert.Equal(t, 30, c.errorPercentThreshold)
	assert.Equal(t, 5, c.sleepWindow)
	assert.Equal(t, 5, c.requestVolumeThreshold)
	assert.Equal(t, "localhost:8125", c.statsD.StatsdAddr)
	assert.Equal(t, "myapp.hystrix", c.statsD.Prefix)
}

func TestOptionsHaveDefaults(t *testing.T) {
	c := NewClient(WithCommandName("test-defaults"))

	assert.Equal(t, 30*time.Second, c.timeout)
	assert.Equal(t, "test-defaults", c.hystrixCommandName)
	assert.Equal(t, 30*time.Second, c.hystrixTimeout)
	assert.Equal(t, 100, c.maxConcurrentRequests)
	assert.Equal(t, 25, c.errorPercentThreshold)
	assert.Equal(t, 10, c.sleepWindow)
	assert.Equal(t, 10, c.requestVolumeThreshold)
	assert.Nil(t, c.statsD)
}

func ExampleWithHTTPTimeout() {
	c := NewClient(WithHTTPTimeout(5 * time.Second))
	req, err := http.NewRequest(http.MethodGet, "https://www.gojek.io/", nil)
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
	req, err := http.NewRequest(http.MethodGet, "https://www.gojek.io/", nil)
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
	req, err := http.NewRequest(http.MethodGet, "https://www.gojek.io/", nil)
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
	req, err := http.NewRequest(http.MethodGet, "https://www.gojek.io/", nil)
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
	req, err := http.NewRequest(http.MethodGet, "https://www.link.doesnt.exist.io/", nil)
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
	// retry attempt 3
	// error
}

func ExampleWithStatsDCollector() {
	c := NewClient(WithStatsDCollector("localhost:8125", "myapp.hystrix"))
	req, err := http.NewRequest(http.MethodGet, "https://www.gojek.io/", nil)
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
