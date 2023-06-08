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
			src:      PllSource{HseFreq: 10, LseFreq: 5, PllSource: "HSE"},
		},
		{
			name:     "set lse source",
			expected: 5,
			src:      PllSource{HseFreq: 10, LseFreq: 5, PllSource: "LSE"},
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
			require.Equal(t, tc.expected, tc.src.SrcFreq)
		})
	}

}

func TestCalculateDivisionFactors(t *testing.T) {
	tests := []struct {
		name     string
		src      PllSource
		expected divFactors
		err      error
	}{
		{
			name:     "calculate for 180Mhz",
			src:      PllSource{SrcFreq: 8_000_000, RequireFreq: 180_000_000},
			expected: divFactors{M: 4, N: 180, P: 2},
			err:      nil,
		},
		{
			name:     "calculate for 120Mhz",
			src:      PllSource{SrcFreq: 8_000_000, RequireFreq: 120_000_000},
			expected: divFactors{M: 4, N: 120, P: 2},
			err:      nil,
		},
		{
			name:     "calculate for 60Mhz",
			src:      PllSource{SrcFreq: 8_000_000, RequireFreq: 60_000_000},
			expected: divFactors{M: 4, N: 60, P: 2},
			err:      nil,
		},
		{
			name:     "calculate for 0Mhz",
			src:      PllSource{SrcFreq: 8_000_000, RequireFreq: 0},
			expected: divFactors{},
			err:      ErrorDivFactorCalc,
		},
		{
			name:     "calculate for 181Mhz",
			src:      PllSource{SrcFreq: 8_000_000, RequireFreq: 180_000_001},
			expected: divFactors{},
			err:      ErrorDivFactorCalc,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.src.calculateDivisionFactors()
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, tc.src.DivFactors)
		})
	}
}
