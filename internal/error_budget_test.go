package internal_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/gojek/heimdall/v7/internal"
	"github.com/stretchr/testify/require"
)

func Test_ErrorBudget_VaryingErrorRate(t *testing.T) {
	t.Parallel()

	const volume = 1000
	const budgetFailurePercent = 10
	var budgetSuccessPerFailure = (100 - budgetFailurePercent) / budgetFailurePercent

	budget := internal.NewPercentErrorBudget(volume, budgetFailurePercent)

	// Note: given budgetSuccessPerFailure is not normalised to the precision supported,
	// we expect few slightly off calculation and we will skip them.

	// 100% failure rate
	testSuccessPerFailure(t, budgetSuccessPerFailure, volume, 0, budget)

	// 50% failure rate
	testSuccessPerFailure(t, budgetSuccessPerFailure, volume, 1, budget)

	// 33% error rate
	testSuccessPerFailure(t, budgetSuccessPerFailure, volume, 2, budget)

	// 25% error rate
	testSuccessPerFailure(t, budgetSuccessPerFailure, volume, 3, budget)

	// 20% error rate
	testSuccessPerFailure(t, budgetSuccessPerFailure, volume, 4, budget)

}

func Test_ErrorBudget_BudgetedFailurePercentZeroOrLower(t *testing.T) {
	t.Parallel()

	for _, failurePercent := range []float32{0.0, -0.1, -1.0, -10000} {
		t.Run(strconv.FormatFloat(float64(failurePercent), 'f', -1, 64), func(t *testing.T) {
			t.Parallel()

			eb := internal.NewPercentErrorBudget(1000, failurePercent)

			for range 999 {
				require.False(t, eb.Failure())
			}
			require.True(t, eb.Failure())
			require.True(t, eb.IsOverBudget())

			// success event foes not reset errorbudget
			for range 10000 {
				require.True(t, eb.Success())
			}
			require.True(t, eb.IsOverBudget())

			// only resetting the errorbudget will reset the state
			eb.Reset()
			require.False(t, eb.IsOverBudget())
		})
	}
}

func Test_ErrorBudget_BudgetedFailurePercentHundredOrHigher(t *testing.T) {
	t.Parallel()

	for _, failurePercent := range []float32{100.0, 100.1, 101.0, 10000} {
		t.Run(strconv.FormatFloat(float64(failurePercent), 'f', -1, 64), func(t *testing.T) {
			t.Parallel()

			eb := internal.NewPercentErrorBudget(1000, failurePercent)

			for range 999 {
				require.False(t, eb.Failure())
			}
			require.True(t, eb.Failure())
			require.True(t, eb.IsOverBudget())

			// success event foes not reset errorbudget
			require.False(t, eb.Success())
			require.False(t, eb.IsOverBudget())

			// verify the budget was reseted, and test with more failures
			for range 999 {
				require.False(t, eb.Failure())
			}
			for range 10000 {
				require.True(t, eb.Failure())
			}
			require.True(t, eb.IsOverBudget())

			// success event foes not reset errorbudget
			require.False(t, eb.Success())
			require.False(t, eb.IsOverBudget())

			// verify the budget was reseted, and test with minimum failures
			for range 999 {
				require.False(t, eb.Failure())
			}
			require.True(t, eb.Failure())
		})
	}
}

func Test_ErrorBudget_Nil(t *testing.T) {
	t.Parallel()

	var eb *internal.ErrorBudget
	for range 999 {
		require.False(t, eb.Failure())
		require.False(t, eb.IsOverBudget())
	}

	for range 999 {
		require.False(t, eb.Success())
		require.False(t, eb.IsOverBudget())
	}

	// Reset should not panic
	eb.Reset()
}

func testSuccessPerFailure(t *testing.T, budgetSuccessPerFailure int, volume int, expectedSuccessPerFailure int, budget *internal.ErrorBudget) {
	t.Helper()

	// budgetSuccessPerFailure is not normalised to the precision supported,
	// so there will be edge case in the calculation
	cnt := failureNeededToExceedBudget(budgetSuccessPerFailure, volume, expectedSuccessPerFailure)
	for range cnt - 1 {
		require.False(t, budget.Failure())
		for range expectedSuccessPerFailure {
			require.False(t, budget.Success())
		}
	}

	require.True(t, budget.Failure())
	require.True(t, budget.IsOverBudget())
	budget.Reset()
}

func failureNeededToExceedBudget(budgetSuccessPerFailure int, volume int, expectedSuccessPerFailure int) int64 {
	// Given default state(/state after long 100% success rate) of the budget tracks accounts for budgetSuccessPerFailure*volume success event
	// We need to adjust our calculation to account for the fact that we are not starting from threshold

	return int64(math.Ceil(float64(budgetSuccessPerFailure*volume) / float64(budgetSuccessPerFailure-expectedSuccessPerFailure)))
}
