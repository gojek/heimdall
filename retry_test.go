package heimdall

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetrierWithExponentialBackoff(t *testing.T) {

	exponentialBackoff := NewExponentialBackoff(2*time.Millisecond, 10*time.Millisecond, 2.0, 1)
	exponentialRetrier := NewRetrier(exponentialBackoff)

	assert.Equal(t, 4*time.Millisecond, exponentialRetrier.NextInterval(1))
}

func TestRetrierWithConstantBackoff(t *testing.T) {

	constantBackoff := NewConstantBackoff(2, 2)
	constantRetrier := NewRetrier(constantBackoff)

	assert.Equal(t, 2*time.Millisecond, constantRetrier.NextInterval(1))
}
