package stm32_target_pll

import (
	"testing"

	l "github.com/komarovn654/embedded_configurator/utils/log"
	"github.com/stretchr/testify/require"
)

func TestSetupPll(t *testing.T) {
	tests := []struct {
		name string
		src  PllSettings
		res  PllSettings
		err  bool
	}{
		{
			name: "no err",
			src:  PllSettings{PllSource: "HSE", HseFreq: 8_000_000, LseFreq: 16_000_000, RequireFreq: 180_000_000},
			res:  PllSettings{PllSource: "HSE", SrcFreq: 8_000_000, DivFactors: divFactors{M: 4, N: 180, P: 2}},
			err:  false,
		},
		{
			name: "assert err",
			src:  PllSettings{PllSource: "hse", HseFreq: 8_000_000, RequireFreq: 180_000_000},
			res:  PllSettings{},
			err:  true,
		},
		{
			name: "calculation err",
			src:  PllSettings{PllSource: "HSE", HseFreq: 8_000_000, LseFreq: 16_000_000, RequireFreq: 1},
			res:  PllSettings{},
			err:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l.InitializeLogger()
			err := tc.src.SetupPll()
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.res.PllSource, tc.src.PllSource)
			require.Equal(t, tc.res.SrcFreq, tc.src.SrcFreq)
			require.Equal(t, tc.res.DivFactors, tc.src.DivFactors)
		})
	}
}

func TestAssertFields(t *testing.T) {
	tests := []struct {
		name   string
		src    PllSettings
		errors bool
	}{
		{
			name:   "pll source error",
			src:    PllSettings{PllSource: "sda", HseFreq: 8_000_000, LseFreq: 16_000_000, RequireFreq: 50},
			errors: true,
		},
		{
			name:   "hse freq error",
			src:    PllSettings{PllSource: "HSE", HseFreq: 0, LseFreq: 16_000_000, RequireFreq: 50},
			errors: true,
		},
		{
			name:   "lse freq error",
			src:    PllSettings{PllSource: "LSE", HseFreq: 8_000_000, LseFreq: 0, RequireFreq: 50},
			errors: true,
		},
		{
			name:   "req freq error",
			src:    PllSettings{PllSource: "HSE", HseFreq: 8_000_000, LseFreq: 16_000_000, RequireFreq: 200_000_000},
			errors: true,
		},
		{
			name:   "no error",
			src:    PllSettings{PllSource: "HSE", HseFreq: 8_000_000, LseFreq: 16_000_000, RequireFreq: 50},
			errors: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.src.validateFields()
			if tc.errors {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestSetPllFreq(t *testing.T) {
	tests := []struct {
		name     string
		src      PllSettings
		expected int
	}{
		{
			name:     "set hse source",
			expected: 10,
			src:      PllSettings{HseFreq: 10, LseFreq: 5, PllSource: "HSE"},
		},
		{
			name:     "set lse source",
			expected: 5,
			src:      PllSettings{HseFreq: 10, LseFreq: 5, PllSource: "LSE"},
		},
		{
			name:     "set unknown source",
			expected: 0,
			src:      PllSettings{HseFreq: 10, LseFreq: 5, PllSource: "fsd"},
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
		src      PllSettings
		expected divFactors
		err      error
	}{
		{
			name:     "calculate for 180Mhz",
			src:      PllSettings{SrcFreq: 8_000_000, RequireFreq: 180_000_000},
			expected: divFactors{M: 4, N: 180, P: 2},
			err:      nil,
		},
		{
			name:     "calculate for 120Mhz",
			src:      PllSettings{SrcFreq: 8_000_000, RequireFreq: 120_000_000},
			expected: divFactors{M: 4, N: 120, P: 2},
			err:      nil,
		},
		{
			name:     "calculate for 60Mhz",
			src:      PllSettings{SrcFreq: 8_000_000, RequireFreq: 60_000_000},
			expected: divFactors{M: 4, N: 60, P: 2},
			err:      nil,
		},
		{
			name:     "calculate for 0Mhz",
			src:      PllSettings{SrcFreq: 8_000_000, RequireFreq: 0},
			expected: divFactors{},
			err:      ErrorDivFactorCalc,
		},
		{
			name:     "calculate for 181Mhz",
			src:      PllSettings{SrcFreq: 8_000_000, RequireFreq: 180_000_001},
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
