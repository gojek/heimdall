package hystrix

import (
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
	)

	assert.Equal(t, 10*time.Second, c.timeout)
	assert.Equal(t, "test", c.hystrixCommandName)
	assert.Equal(t, time.Duration(1100), c.hystrixTimeout)
	assert.Equal(t, 10, c.maxConcurrentRequests)
	assert.Equal(t, 30, c.errorPercentThreshold)
	assert.Equal(t, 5, c.sleepWindow)
	assert.Equal(t, 5, c.requestVolumeThreshold)
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
}
