package httpclient

import (
	"context"
	"time"
)

// SleepInterruptible sleeps until either the timer triggers or context is cancelled
func SleepInterruptible(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
	}
	return nil
}
