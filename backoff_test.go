package heimdall

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoffNextTime(t *testing.T) {
	exponentialBackoff := NewExponentialBackoff(100*time.Millisecond, 1000*time.Millisecond, 2.0, 0*time.Millisecond)

	assert.Equal(t, 100*time.Millisecond, exponentialBackoff.Next(0))
	assert.Equal(t, 200*time.Millisecond, exponentialBackoff.Next(1))
	assert.Equal(t, 400*time.Millisecond, exponentialBackoff.Next(2))
	assert.Equal(t, 800*time.Millisecond, exponentialBackoff.Next(3))
}

func TestExponentialBackoffWithInvalidJitter(t *testing.T) {
	exponentialBackoff := NewExponentialBackoff(100*time.Millisecond, 1000*time.Millisecond, 2.0, -1*time.Millisecond)

	assert.Equal(t, 100*time.Millisecond, exponentialBackoff.Next(0))
	assert.Equal(t, 200*time.Millisecond, exponentialBackoff.Next(1))
	assert.Equal(t, 400*time.Millisecond, exponentialBackoff.Next(2))
	assert.Equal(t, 800*time.Millisecond, exponentialBackoff.Next(3))
}

func TestExponentialBackoffMaxTimeoutCrossed(t *testing.T) {
	exponentialBackoff := NewExponentialBackoff(100*time.Millisecond, 1000*time.Millisecond, 2.0, 0*time.Millisecond)

	assert.Equal(t, 1000*time.Millisecond, exponentialBackoff.Next(4))
}

func TestExponentialBackoffMaxTimeoutReached(t *testing.T) {
	exponentialBackoff := NewExponentialBackoff(100*time.Millisecond, 1600*time.Millisecond, 2.0, 0*time.Millisecond)

	assert.Equal(t, 1600*time.Millisecond, exponentialBackoff.Next(4))
}

func TestExponentialBackoffWhenRetryIsLessThanZero(t *testing.T) {
	exponentialBackoff := NewExponentialBackoff(100*time.Millisecond, 1000*time.Millisecond, 2.0, 0*time.Millisecond)

	assert.Equal(t, 100*time.Millisecond, exponentialBackoff.Next(-1))
}

func TestExponentialBackoffJitter0(t *testing.T) {
	exponentialBackoff := NewExponentialBackoff(100*time.Millisecond, 1000*time.Millisecond, 2.0, 0*time.Millisecond)
	for i := 0; i < 10000; i++ {
		assert.Equal(t, 200*time.Millisecond, exponentialBackoff.Next(1))
	}
}

func TestExponentialBackoffJitter1(t *testing.T) {
	exponentialBackoff := NewExponentialBackoff(100*time.Millisecond, 1000*time.Millisecond, 2.0, 1*time.Millisecond)
	for i := 0; i < 10000; i++ {
		assert.True(t, 200*time.Millisecond <= exponentialBackoff.Next(1) && exponentialBackoff.Next(1) <= 201*time.Millisecond)
	}
}

func TestExponentialBackoffJitter50(t *testing.T) {
	exponentialBackoff := NewExponentialBackoff(100*time.Millisecond, 1000*time.Millisecond, 2.0, 50*time.Millisecond)
	for i := 0; i < 10000; i++ {
		assert.True(t, 200*time.Millisecond <= exponentialBackoff.Next(1) && exponentialBackoff.Next(1) <= 250*time.Millisecond)
	}
}

func TestConstantBackoffNextTime(t *testing.T) {
	constantBackoff := NewConstantBackoff(100*time.Millisecond, 0*time.Millisecond)

	assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(0))
	assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(1))
	assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(2))
	assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(3))
}

func TestConstantBackoffWithInvalidJitter(t *testing.T) {
	constantBackoff := NewConstantBackoff(100*time.Millisecond, -1*time.Millisecond)

	assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(0))
	assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(1))
	assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(2))
	assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(3))
}

func TestConstantBackoffWhenRetryIsLessThanZero(t *testing.T) {
	constantBackoff := NewConstantBackoff(100*time.Millisecond, 0*time.Millisecond)

	assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(-1))
}

func TestConstantBackoffJitter0(t *testing.T) {
	constantBackoff := NewConstantBackoff(100*time.Millisecond, 0*time.Millisecond)
	for i := 0; i < 10000; i++ {
		assert.Equal(t, 100*time.Millisecond, constantBackoff.Next(i))
	}
}

func TestConstantBackoffJitter1(t *testing.T) {
	constantBackoff := NewConstantBackoff(100*time.Millisecond, 1*time.Millisecond)
	for i := 0; i < 10000; i++ {
		assert.True(t, 100*time.Millisecond <= constantBackoff.Next(i) && constantBackoff.Next(1) <= 101*time.Millisecond)
	}
}

func TestConstantBackoffJitter50(t *testing.T) {
	constantBackoff := NewConstantBackoff(100*time.Millisecond, 50*time.Millisecond)
	for i := 0; i < 10000; i++ {
		assert.True(t, 100*time.Millisecond <= constantBackoff.Next(i) && constantBackoff.Next(1) <= 150*time.Millisecond)
	}
}
