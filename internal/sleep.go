package internal

import (
	"context"
	"time"
)

// SleepInterruptible sleeps until either the timer triggers or context is cancelled
func SleepInterruptible(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
	}
	return nil
}
