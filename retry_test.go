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

	constantBackoff := NewConstantBackoff(2, 1)
	constantRetrier := NewRetrier(constantBackoff)

	assert.True(t, 2*time.Millisecond <= constantRetrier.NextInterval(1))
}
