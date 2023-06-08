package stm32_pllconfig

import (
	"errors"

	"github.com/komarovn654/embedded_configurator/utils"
	validator "github.com/komarovn654/struct_validator"
)

var (
	FactorValuesP = [4]int{2, 4, 6, 8}
	FactorM       = Range{2, 63}
	FactorN       = Range{50, 432}
	FreqVcoIn     = Range{1_000_000, 2_000_000}
	FreqVcoOut    = Range{100_000_000, 432_000_000}

	ErrorDivFactorCalc = errors.New("error in the calculation of division factors")
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

	SrcFreq    int
	DivFactors divFactors
}

type divFactors struct {
	P int
	N int
	M int
	Q int
	R int
}

func NewPllSource() *PllSource {
	return &PllSource{}
}

func (this *PllSource) SetupPll() error {
	// if err := this.assertFields(); err != nil {
	// 	return err
	// }
	utils.Logger.Sugar().Info("success fields assert")

	this.setSrcFreq()
	utils.Logger.Sugar().Infof("pll source freq: %vHz, with %v source", this.SrcFreq, this.PllSource)

	if err := this.calculateDivisionFactors(); err != nil {
		return err
	}
	utils.Logger.Sugar().Infof("div factors: m: %v, n: %v, p: %v", this.DivFactors.M, this.DivFactors.N, this.DivFactors.P)

	return nil
}

func (this *PllSource) assertFields() error {
	return validator.Validate(this)
}

func (this *PllSource) setSrcFreq() {
	if this.PllSource == "HSE" {
		this.SrcFreq = this.HseFreq
	}

	if this.PllSource == "LSE" {
		this.SrcFreq = this.LseFreq
	}
}

func (this *PllSource) calculateDivisionFactors() error {
	for m := FactorM.min; m <= FactorM.max; m++ {
		vcoIn := this.SrcFreq / m
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
					this.DivFactors.M = m
					this.DivFactors.N = n
					this.DivFactors.P = p
					return nil
				}
			}
		}
	}

	return ErrorDivFactorCalc
}
