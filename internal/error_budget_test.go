package internal_test

import (
	"strconv"
	"testing"

	"github.com/gojek/heimdall/v7/internal"
	"github.com/stretchr/testify/require"
)

func Test_ErrorBudget_VaryingErrorRate(t *testing.T) {
	t.Parallel()

	const minFailureVolume = 1000
	const budgetFailurePercent = 10

	// Note: given budgetSuccessPerFailure is not normalised to the precision supported,
	// we expect few slightly off calculation and we will skip them.

	t.Run("100% failure rate", func(t *testing.T) {
		t.Parallel()
		// failureNeededToExceedBudget same as minFailureVolume
		testSuccessPerFailure(t, budgetFailurePercent, minFailureVolume, 0, 1000)
	})

	t.Run("50% failure rate", func(t *testing.T) {
		t.Parallel()
		// Default state: 9000 preceding success events (i.e. (100 - budgetFailurePercent) * minFailureVolume)
		// failureNeededToExceedBudget = 1125, new events needed = 2250
		// effective failure rate(including default state): 1125/(2250+9000) = 10%
		testSuccessPerFailure(t, budgetFailurePercent, minFailureVolume, 1, 1125)
	})

	t.Run("33% failure rate", func(t *testing.T) {
		t.Parallel()
		// Default state: 9000 preceding success events (i.e. (100 - budgetFailurePercent) * minFailureVolume)
		// failureNeededToExceedBudget = 1286, new events needed = 3858
		// effective failure rate(including default state): 1286/(3858+9000) = ~10%
		testSuccessPerFailure(t, budgetFailurePercent, minFailureVolume, 2, 1286)
	})

	t.Run("25% failure rate", func(t *testing.T) {
		t.Parallel()
		// Default state: 9000 preceding success events (i.e. (100 - budgetFailurePercent) * minFailureVolume)
		// failureNeededToExceedBudget = 1500, new events needed = 6000
		// effective failure rate(including default state): 1500/(6000+9000) = 10%
		testSuccessPerFailure(t, budgetFailurePercent, minFailureVolume, 3, 1500)
	})

	t.Run("20% failure rate", func(t *testing.T) {
		t.Parallel()
		// Default state: 9000 preceding success events (i.e. (100 - budgetFailurePercent) * minFailureVolume)
		// failureNeededToExceedBudget = 1800, new events needed = 9000
		// effective failure rate(including default state): 1800/(9000+9000) = 10%
		testSuccessPerFailure(t, budgetFailurePercent, minFailureVolume, 4, 1800)
	})
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

func testSuccessPerFailure(t *testing.T, budgetFailurePercent float32, volume int32, expectedSuccessPerFailure int32, failureNeededToExceedBudget int) {
	t.Helper()
	budget := internal.NewPercentErrorBudget(volume, budgetFailurePercent)

	for range failureNeededToExceedBudget - 1 {
		require.False(t, budget.Failure())
		for range expectedSuccessPerFailure {
			require.False(t, budget.Success())
		}
	}

	require.True(t, budget.Failure())
	require.True(t, budget.IsOverBudget())
	budget.Reset()
}
