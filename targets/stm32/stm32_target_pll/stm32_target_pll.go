package stm32_target_pll

import (
	"errors"

	logger "github.com/komarovn654/embedded_configurator/utils/log"
	validator "github.com/komarovn654/struct_validator"
)

var (
	FactorValuesP = [4]int{2, 4, 6, 8}
	FactorM       = Range{2, 63}
	FactorN       = Range{50, 432}
	FreqVcoIn     = Range{1_000_000, 2_000_000}
	FreqVcoOut    = Range{100_000_000, 432_000_000}

	ErrorDivFactorCalc  = errors.New("error in the calculation of division factors")
	ErrorUnsupPllSource = errors.New("unsupported pll source")
)

type Range struct {
	min int
	max int
}

type PllTarget struct {
	PllSource   string `mapstructure:"PllSource" validate:"in:HSE,LSE"`
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

func NewPllTarget() *PllTarget {
	return new(PllTarget)
}

func (target *PllTarget) SetupPll() error {
	if err := target.validateFields(); err != nil {
		return err
	}
	logger.Info("success fields validate")

	if err := target.setSrcFreq(); err != nil {
		return err
	}
	logger.Infof("pll source freq: %vHz with %v source", target.SrcFreq, target.PllSource)

	if err := target.calculateDivisionFactors(); err != nil {
		return err
	}
	logger.Infof("div factors: m: %v, n: %v, p: %v", target.DivFactors.M, target.DivFactors.N, target.DivFactors.P)

	return nil
}

func (target *PllTarget) validateFields() error {
	return validator.Validate(*target)
}

func (target *PllTarget) setSrcFreq() error {
	switch target.PllSource {
	case "HSE":
		target.SrcFreq = target.HseFreq
	case "LSE":
		target.SrcFreq = target.LseFreq
	default:
		return ErrorUnsupPllSource
	}

	return nil
}

func (target *PllTarget) calculateDivisionFactors() error {
	for m := FactorM.min; m <= FactorM.max; m++ {
		vcoIn := target.SrcFreq / m
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
				if pllFreq == target.RequireFreq {
					target.DivFactors.M = m
					target.DivFactors.N = n
					target.DivFactors.P = p
					return nil
				}
			}
		}
	}

	return ErrorDivFactorCalc
}
