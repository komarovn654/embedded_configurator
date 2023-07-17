package stm32_target_pll

import (
	"testing"

	l "github.com/komarovn654/embedded_configurator/utils/log"
	"github.com/stretchr/testify/require"
)

func TestSetupPll(t *testing.T) {
	tests := []struct {
		name string
		src  PllTarget
		res  PllTarget
		err  bool
	}{
		{
			name: "no err",
			src:  PllTarget{PllSource: "HSE", HseFreq: 8_000_000, HsiFreq: 16_000_000, RequireFreq: 180_000_000},
			res:  PllTarget{PllSource: "HSE", SrcFreq: 8_000_000, DivFactors: divFactors{M: 4, N: 180, P: 2}},
			err:  false,
		},
		{
			name: "assert err",
			src:  PllTarget{PllSource: "hse", HseFreq: 8_000_000, RequireFreq: 180_000_000},
			res:  PllTarget{},
			err:  true,
		},
		{
			name: "calculation err",
			src:  PllTarget{PllSource: "HSE", HseFreq: 8_000_000, HsiFreq: 16_000_000, RequireFreq: 1},
			res:  PllTarget{},
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
		src    PllTarget
		errors bool
	}{
		{
			name:   "pll source error",
			src:    PllTarget{PllSource: "sda", HseFreq: 8_000_000, HsiFreq: 16_000_000, RequireFreq: 50},
			errors: true,
		},
		{
			name:   "hse freq error",
			src:    PllTarget{PllSource: "HSE", HseFreq: 0, HsiFreq: 16_000_000, RequireFreq: 50},
			errors: true,
		},
		{
			name:   "hsi freq error",
			src:    PllTarget{PllSource: "HSI", HseFreq: 8_000_000, HsiFreq: 0, RequireFreq: 50},
			errors: true,
		},
		{
			name:   "req freq error",
			src:    PllTarget{PllSource: "HSE", HseFreq: 8_000_000, HsiFreq: 16_000_000, RequireFreq: 200_000_000},
			errors: true,
		},
		{
			name:   "no error",
			src:    PllTarget{PllSource: "HSE", HseFreq: 8_000_000, HsiFreq: 16_000_000, RequireFreq: 50},
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
		src      PllTarget
		expected int
	}{
		{
			name:     "set hse source",
			expected: 10,
			src:      PllTarget{HseFreq: 10, HsiFreq: 5, PllSource: "HSE"},
		},
		{
			name:     "set hsi source",
			expected: 5,
			src:      PllTarget{HseFreq: 10, HsiFreq: 5, PllSource: "HSI"},
		},
		{
			name:     "set unknown source",
			expected: 0,
			src:      PllTarget{HseFreq: 10, HsiFreq: 5, PllSource: "fsd"},
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
	l.InitializeLogger()

	tests := []struct {
		name     string
		src      PllTarget
		expected divFactors
		err      error
	}{
		{
			name:     "calculate for 180Mhz",
			src:      PllTarget{SrcFreq: 8_000_000, RequireFreq: 180_000_000},
			expected: divFactors{M: 4, N: 180, P: 2, Q: 8},
			err:      nil,
		},
		{
			name:     "calculate for 120Mhz",
			src:      PllTarget{SrcFreq: 8_000_000, RequireFreq: 120_000_000},
			expected: divFactors{M: 4, N: 120, P: 2, Q: 5},
			err:      nil,
		},
		{
			name:     "calculate for 60Mhz",
			src:      PllTarget{SrcFreq: 8_000_000, RequireFreq: 60_000_000},
			expected: divFactors{M: 4, N: 60, P: 2, Q: 3},
			err:      nil,
		},
		{
			name:     "calculate for 0Mhz",
			src:      PllTarget{SrcFreq: 8_000_000, RequireFreq: 0},
			expected: divFactors{},
			err:      ErrorDivFactorCalc,
		},
		{
			name:     "calculate for 181Mhz",
			src:      PllTarget{SrcFreq: 8_000_000, RequireFreq: 180_000_001},
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

func TestSetQ(t *testing.T) {
	t.Run("usb prio", func(t *testing.T) {
		validFreqs := map[int]int{
			144000000: 3,
			192000000: 4,
			240000000: 5,
			288000000: 6,
			336000000: 7,
			384000000: 8,
			432000000: 9,
		}
		for vco := 144000000; vco <= FreqVcoOut.max; vco += 24_000_000 {
			target := NewPllTarget()
			target.UsbFreqPrio = true
			err := target.setQ(vco)
			if q, ok := validFreqs[vco]; ok {
				require.Equal(t, USBFreq, target.Pll48Freq)
				require.Equal(t, q, target.DivFactors.Q)
				require.NoError(t, err)
				continue
			}
			require.Equal(t, err, ErrorAssertQ)
		}
	})
	t.Run("without usb prio", func(t *testing.T) {
		for vco := 144000000; vco <= FreqVcoOut.max; vco += 10_000_000 {
			target := NewPllTarget()
			err := target.setQ(vco)
			require.LessOrEqual(t, target.Pll48Freq, USBFreq)
			require.GreaterOrEqual(t, target.DivFactors.Q, FactorQ.min)
			require.LessOrEqual(t, target.DivFactors.Q, FactorQ.max)
			require.NoError(t, err)
		}
	})
}

func TestAssertFrequency(t *testing.T) {
	tests := []struct {
		name  string
		value int
		rng   Range
		err   error
	}{
		{
			name:  "no error",
			value: 11,
			rng:   Range{0, 11},
			err:   nil,
		},
		{
			name:  "no error",
			value: 0,
			rng:   Range{0, 11},
			err:   nil,
		},
		{
			name:  "no error",
			value: 5,
			rng:   Range{0, 11},
			err:   nil,
		},
		{
			name:  "error",
			value: 144,
			rng:   Range{0, 11},
			err:   ErrorAssertVcoFreq,
		},
		{
			name:  "error",
			value: -1,
			rng:   Range{0, 11},
			err:   ErrorAssertVcoFreq,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := assertFrequency(tc.value, tc.rng.min, tc.rng.max)
			require.Equal(t, tc.err, err)
		})
	}
}
