package stm32_pllconfig

import (
	validator "github.com/komarovn654/struct_validator"
)

var (
	FactorValuesP = [4]int{2, 4, 6, 8}
	FactorM       = Range{2, 63}
	FactorN       = Range{50, 432}
	FreqVcoIn     = Range{1_000_000, 2_000_000}
	FreqVcoOut    = Range{100_000_000, 432_000_000}
)

type Range struct {
	min int
	max int
}

type PllSource struct {
	PllSource   string `mapstructure:"PllSource" validate:"in:hse,lse"`
	HseFreq     int    `mapstructure:"HseFrequency" validate:"min:4000000|max:24000000"`
	LseFreq     int    `mapstructure:"LseFrequency" validate:"in:16000000"`
	RequireFreq int    `mapstructure:"RequireFrequency" validate:"max:180000000"`

	srcFreq    int
	divFactors DivFactors
}

type PllConfig struct {
	Src PllSource `mapstructure:"PllConfig"`
}

type DivFactors struct {
	p int
	n int
	m int
	q int
	r int
}

func (this *PllSource) AssertFields() error {
	return validator.Validate(this)
}

func (this *PllSource) setSrcFreq() {
	if this.PllSource == "hse" {
		this.srcFreq = this.HseFreq
	}

	if this.PllSource == "lse" {
		this.srcFreq = this.LseFreq
	}
}

func (this *PllSource) CalculateDivisionFactors() bool {
	for m := FactorM.min; m <= FactorM.max; m++ {
		vcoIn := this.srcFreq / m
		if (vcoIn < FreqVcoIn.min) || (vcoIn > FreqVcoIn.max) {
			continue
		}

		for n := FactorN.min; n <= FactorN.max; n++ {
			vcoOut := vcoIn * n
			if (vcoOut < FreqVcoOut.min) || (vcoOut > FreqVcoOut.max) {
				continue
			}

			for _, p := range FactorValuesP {
				pllFreq := vcoOut / p
				if pllFreq == this.RequireFreq {
					this.divFactors.m = m
					this.divFactors.n = n
					this.divFactors.p = p
					return true
				}
			}
		}
	}

	return false
}
