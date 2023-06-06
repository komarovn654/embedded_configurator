package stm32_pllconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAssertFields(t *testing.T) {
	// tests := []struct {
	// 	name   string
	// 	src    PllSource
	// 	errors bool
	// }{
	// 	{
	// 		name:   "pll source error",
	// 		src:    PllSource{PllSource: "sda", HseFreq: 8_000_000, LseFreq: 16_000_000, RequireFreq: 50},
	// 		errors: true,
	// 	},
	// 	{
	// 		name:   "hse freq error",
	// 		src:    PllSource{PllSource: "hse", HseFreq: 0, LseFreq: 16_000_000, RequireFreq: 50},
	// 		errors: true,
	// 	},
	// 	{
	// 		name:   "lse freq error",
	// 		src:    PllSource{PllSource: "hse", HseFreq: 8_000_000, LseFreq: 0, RequireFreq: 50},
	// 		errors: true,
	// 	},
	// 	{
	// 		name:   "req freq error",
	// 		src:    PllSource{PllSource: "hse", HseFreq: 8_000_000, LseFreq: 16_000_000, RequireFreq: 200_000_000},
	// 		errors: true,
	// 	},
	// 	{
	// 		name:   "no error",
	// 		src:    PllSource{PllSource: "hse", HseFreq: 8_000_000, LseFreq: 16_000_000, RequireFreq: 50},
	// 		errors: false,
	// 	},
	// }

	// for _, tc := range tests {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		err := tc.src.AssertFields()
	// 		fmt.Println(err)
	// 		if tc.errors && err.Error() != "invalid type, expected struct" {
	// 			require.Error(t, err)
	// 			return
	// 		}
	// 		require.NoError(t, err)
	// 	})
	// }
}

func TestSetPllFreq(t *testing.T) {
	tests := []struct {
		name     string
		src      PllSource
		expected int
	}{
		{
			name:     "set hse source",
			expected: 10,
			src:      PllSource{HseFreq: 10, LseFreq: 5, PllSource: "hse"},
		},
		{
			name:     "set lse source",
			expected: 5,
			src:      PllSource{HseFreq: 10, LseFreq: 5, PllSource: "lse"},
		},
		{
			name:     "set unknown source",
			expected: 0,
			src:      PllSource{HseFreq: 10, LseFreq: 5, PllSource: "fsd"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.src.setSrcFreq()
			require.Equal(t, tc.expected, tc.src.srcFreq)
		})
	}

}

func TestCalculateDivisionFactors(t *testing.T) {
	tests := []struct {
		name     string
		src      PllSource
		expected DivFactors
		result   bool
	}{
		{
			name:     "calculate for 180Mhz",
			src:      PllSource{srcFreq: 8_000_000, RequireFreq: 180_000_000},
			expected: DivFactors{m: 4, n: 180, p: 2},
			result:   true,
		},
		{
			name:     "calculate for 120Mhz",
			src:      PllSource{srcFreq: 8_000_000, RequireFreq: 120_000_000},
			expected: DivFactors{m: 4, n: 120, p: 2},
			result:   true,
		},
		{
			name:     "calculate for 60Mhz",
			src:      PllSource{srcFreq: 8_000_000, RequireFreq: 60_000_000},
			expected: DivFactors{m: 4, n: 60, p: 2},
			result:   true,
		},
		{
			name:     "calculate for 0Mhz",
			src:      PllSource{srcFreq: 8_000_000, RequireFreq: 0},
			expected: DivFactors{},
			result:   false,
		},
		{
			name:     "calculate for 181Mhz",
			src:      PllSource{srcFreq: 8_000_000, RequireFreq: 180_000_001},
			expected: DivFactors{},
			result:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.src.CalculateDivisionFactors()
			require.Equal(t, tc.result, result)
			require.Equal(t, tc.expected, tc.src.divFactors)
		})
	}
}
