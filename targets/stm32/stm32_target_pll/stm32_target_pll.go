package stm32_target_pll

import (
	"errors"

	logger "github.com/komarovn654/embedded_configurator/utils/log"
	validator "github.com/komarovn654/struct_validator"
)

var (
	FactorValuesP      = [4]int{2, 4, 6, 8}
	FactorM            = Range{2, 63}
	FactorN            = Range{50, 432}
	FactorQ            = Range{2, 15}
	FreqVcoIn          = Range{1_000_000, 2_000_000}
	FreqVcoOut         = Range{100_000_000, 432_000_000}
	USBFreq            = 48_000_000
	RequireFreqsForUSB = []int{168_000_000, 144_000_000, 120_000_000, 108_000_000, 96_000_000,
		84_000_000, 72_000_000, 64_000_000, 60_000_000, 56_000_000, 54_000_000, 48_000_000, 42_000_000,
		40_000_000, 36_000_000, 32_000_000, 30_000_000, 24_000_000, 18_000_000}

	ErrorDivFactorCalc  = errors.New("error in the calculation of division factors")
	ErrorAssertQ        = errors.New("q division factor assert failed")
	ErrorAssertVcoFreq  = errors.New("frequency assert failed")
	ErrorUnsupPllSource = errors.New("unsupported pll source")
)

type Range struct {
	min int
	max int
}

type PllTarget struct {
	PllSource   string `mapstructure:"PllSource" validate:"in:HSE,HSI"`
	HseFreq     int    `mapstructure:"HseFrequency" validate:"min:4000000|max:24000000"`
	HsiFreq     int    `mapstructure:"HsiFrequency" validate:"in:16000000"`
	RequireFreq int    `mapstructure:"RequireFrequency" validate:"max:180000000"`
	UsbFreqPrio bool   `mapstructure:"UsbFrequencyPriority"`

	SrcFreq    int
	Pll48Freq  int
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
		logger.Warnf("RequireFrequinces for usb: %v", RequireFreqsForUSB)
		return err
	}
	logger.Infof("div factors: %+v", target.DivFactors)

	return nil
}

func (target *PllTarget) validateFields() error {
	return validator.Validate(*target)
}

func (target *PllTarget) setSrcFreq() error {
	switch target.PllSource {
	case "HSE":
		target.SrcFreq = target.HseFreq
	case "HSI":
		target.SrcFreq = target.HsiFreq
	default:
		return ErrorUnsupPllSource
	}

	return nil
}

func assertFrequency(value int, min int, max int) error {
	if value < min || value > max {
		return ErrorAssertVcoFreq
	}
	return nil
}

func (target *PllTarget) setQ(vcoFreq int) error {
	if target.UsbFreqPrio && vcoFreq%USBFreq != 0 {
		return ErrorAssertQ
	}

	for q := FactorQ.min; q <= FactorQ.max; q++ {
		pll48 := float64(vcoFreq) / float64(q)
		if pll48 <= float64(USBFreq) {
			target.Pll48Freq = vcoFreq / q
			target.DivFactors.Q = q
			return nil
		}
	}
	return ErrorAssertQ
}

func (target *PllTarget) calculateDivisionFactors() error {
	for m := FactorM.min; m <= FactorM.max; m++ {
		vcoIn := target.SrcFreq / m
		if err := assertFrequency(vcoIn, FreqVcoIn.min, FreqVcoIn.max); err != nil {
			continue
		}

		for n := FactorN.min; n <= FactorN.max; n++ {
			vcoOut := vcoIn * n
			if err := assertFrequency(vcoOut, FreqVcoOut.min, FreqVcoOut.max); err != nil {
				continue
			}

			if err := target.setQ(vcoOut); err != nil {
				continue
			}

			for _, p := range FactorValuesP {
				pllFreq := vcoOut / p
				if pllFreq == target.RequireFreq {
					target.DivFactors.M = m
					target.DivFactors.N = n
					target.DivFactors.P = p
					logger.Infof("VcoIn: %v, VcoOut: %v", vcoIn, vcoOut)
					return nil
				}
			}
		}
	}

	target.DivFactors = divFactors{}
	return ErrorDivFactorCalc
}
