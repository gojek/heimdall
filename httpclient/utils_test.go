package httpclient

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSleepInterruptible_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the context immediately
	cancel()

	err := SleepInterruptible(ctx, 10*time.Second) // Long duration to ensure cancellation is what stops it
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestSleepInterruptible_CompletesWithoutCancel(t *testing.T) {
	ctx := context.Background()
	start := time.Now()

	err := SleepInterruptible(ctx, 50*time.Millisecond) // Short sleep time
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.True(t, elapsed.Milliseconds() >= int64(50), "Sleep duration should be at least 50ms") // Ensure it slept approximately 50ms
}
