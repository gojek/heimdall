package heimdall

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetrierWithExponentialBackoff(t *testing.T) {

	exponentialBackoff := NewExponentialBackoff(2*time.Millisecond, 10*time.Millisecond, 2.0, 1*time.Millisecond)
	exponentialRetrier := NewRetrier(exponentialBackoff)

	assert.True(t, 4*time.Millisecond <= exponentialRetrier.NextInterval(1))
}

func TestRetrierWithConstantBackoff(t *testing.T) {
	backoffInterval := 2 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	constantBackoff := NewConstantBackoff(backoffInterval, maximumJitterInterval)
	constantRetrier := NewRetrier(constantBackoff)

	assert.True(t, 2*time.Millisecond <= constantRetrier.NextInterval(1))
}

func TestRetrierFunc(t *testing.T) {
	linearRetrier := NewRetrierFunc(func(retry int) time.Duration {
		if retry <= 0 {
			return 0 * time.Millisecond
		}
		return time.Duration(retry) * time.Millisecond
	})

	assert.True(t, 3*time.Millisecond <= linearRetrier.NextInterval(4))
}

func TestNoRetrier(t *testing.T) {
	noRetrier := NewNoRetrier()
	nextInterval := noRetrier.NextInterval(1)
	assert.Equal(t, time.Duration(0), nextInterval)
}
