package internal

import (
	"fmt"
	"sync/atomic"
)

const normaliser = 10000
const maxAllowedToken = ((1 << 30) - 1) / normaliser // signed int32 + overflow protection

// ErrorBudget is used to track if defined error budget is exceeded with weighted tokens.
type ErrorBudget struct {
	token atomic.Int32 // Current token count

	maxToken     int32 // Maximum number of tokens.
	threshold    int32 // Threshold for determining if over/under budget
	successToken int32 // Tokens added on success event.
	failureToken int32 // Tokens added on failure event.
}

// NewTokenErrorBudget creates a weighted token ErrorBudget with the following token details.
//
//	maxToken: The maximum/initial token value which is used to calculate token threshold(i.e. maxToken/2)
//	tokenRatio: The allowed ratio of failure in comparison to success.
func NewTokenErrorBudget(maxToken int32, tokenRatio float32) *ErrorBudget {
	if maxToken > maxAllowedToken {
		panic(fmt.Errorf("max token exceeds allowed limit (%d)", maxAllowedToken))
	}
	normalisedMaxToken := maxToken * normaliser

	eb := &ErrorBudget{
		maxToken:     normalisedMaxToken,
		threshold:    normalisedMaxToken / 2,
		successToken: int32(tokenRatio * normaliser),
		failureToken: -normaliser,
	}
	eb.token.Store(normalisedMaxToken)
	return eb
}

// NewPercentErrorBudget creates a weighted token ErrorBudget with the following failure details.
//
//	minFailureVolume: The minimum failure required.
//	failurePercent: The failure percentage (0-100).
//
// Note: To determine if budget is exceeded we use recent event which satisfies following
//
//	failureEvent <= maxFailureEvent
//	successEvent = (maxFailureEvent-failureEvent) / allowedSuccessPerFailure
//	totalEvent = failureEvent + successEvent
//
// Where
//
//	maxFailureEvent = minFailureVolume * 2
//	allowedSuccessPerFailure = (100 - failurePercent) / failurePercent
func NewPercentErrorBudget(minFailureVolume int32, failurePercent float32) *ErrorBudget {
	var tokenRatio float32

	maxToken := minFailureVolume * 2

	switch {
	case failurePercent <= 0:
		// We will have over budget if lifetime failure count is higher than allowed failure volume
		tokenRatio = 0
	case failurePercent >= 100:
		// a single success event will effectively reset the budget tracker
		tokenRatio = float32(maxToken)
	default:
		tokenRatio = float32(failurePercent) / (100 - float32(failurePercent))
	}

	return NewTokenErrorBudget(maxToken, tokenRatio)
}

// Success registers a successful operation and returns whether the error budget is over the threshold.
// Returns true if the token count is below the threshold, indicating over budget.
func (eb *ErrorBudget) Success() (overbudget bool) {
	if eb == nil {
		return false
	}

	if eb.token.Load() >= eb.maxToken {
		return false
	}

	token := eb.token.Add(eb.successToken)
	if token > eb.maxToken {
		eb.token.CompareAndSwap(token, eb.maxToken) // best-effort clamp
	}

	return token <= eb.threshold
}

// Failure registers a failed operation and returns whether the error budget is over the threshold.
// Returns true if the token count is below the threshold, indicating overbudget.
func (eb *ErrorBudget) Failure() (overbudget bool) {
	if eb == nil {
		return false
	}

	if eb.token.Load() <= 0 {
		return true
	}

	token := eb.token.Add(eb.failureToken)
	if token < 0 {
		eb.token.CompareAndSwap(token, 0) // best-effort clamp
	}

	return token <= eb.threshold
}

// IsOverBudget checks if the error budget is over the threshold.
// Returns true if the token count is below the threshold, indicating overbudget.
func (eb *ErrorBudget) IsOverBudget() bool {
	if eb == nil {
		return false
	}

	return eb.token.Load() <= eb.threshold
}

// Reset resets the error budget to the initial state.
func (eb *ErrorBudget) Reset() {
	if eb == nil {
		return
	}

	eb.token.Store(eb.maxToken)
}
