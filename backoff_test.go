package heimdall

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoffNextTime(t *testing.T) {

	exponentialBackoff := NewExponentialBackoff(2*time.Millisecond, 10*time.Millisecond, 2.0)

	assert.Equal(t, 4*time.Millisecond, exponentialBackoff.Next(1))
}

func TestExponentialBackoffMaxTimeoutCrossed(t *testing.T) {

	exponentialBackoff := NewExponentialBackoff(2*time.Millisecond, 9*time.Millisecond, 2.0)

	assert.Equal(t, 9*time.Millisecond, exponentialBackoff.Next(3))
}

func TestExponentialBackoffMaxTimeoutReached(t *testing.T) {

	exponentialBackoff := NewExponentialBackoff(2*time.Millisecond, 10*time.Millisecond, 2.0)

	assert.Equal(t, 10*time.Millisecond, exponentialBackoff.Next(3))
}

func TestExponentialBackoffWhenRetryIsZero(t *testing.T) {

	exponentialBackoff := NewExponentialBackoff(2*time.Millisecond, 10*time.Millisecond, 2.0)

	assert.Equal(t, 0*time.Millisecond, exponentialBackoff.Next(0))
}

func TestConstantBackoffNextTime(t *testing.T) {

	constantBackoff := NewConstantBackoff(100, 50)

	assert.NotEqual(t, 100*time.Millisecond, constantBackoff.Next(1))
}

func TestConstantBackoffWhenRetryIsZero(t *testing.T) {

	constantBackoff := NewConstantBackoff(100, 50)

	assert.Equal(t, 0*time.Millisecond, constantBackoff.Next(0))
}
