package heimdall

import "time"

// Retriable defines contract for retriers to implement
type Retriable interface {
	NextInterval(retry int) time.Duration
}

// RetriableFunc is an adapter to allow the use of ordinary functions
// as a Retriable
type RetriableFunc func(retry int) time.Duration

// NextInterval calls f(retry)
func (f RetriableFunc) NextInterval(retry int) time.Duration {
	return f(retry)
}

type retrier struct {
	backoff Backoff
}

// NewRetrier returns retrier with some backoff strategy
func NewRetrier(backoff Backoff) Retriable {
	return &retrier{
		backoff: backoff,
	}
}

// NewRetrierFunc returns a retrier with a retry function defined
func NewRetrierFunc(f RetriableFunc) Retriable {
	return f
}

// NextInterval returns next retriable time
func (r *retrier) NextInterval(retry int) time.Duration {
	return r.backoff.Next(retry)
}

type noRetrier struct {
}

// NewNoRetrier returns a null object for retriable
func NewNoRetrier() Retriable {
	return &noRetrier{}
}

// NextInterval returns next retriable time, always 0
func (r *noRetrier) NextInterval(retry int) time.Duration {
	return 0 * time.Millisecond
}
