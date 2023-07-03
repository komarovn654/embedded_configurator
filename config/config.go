package config

import (
	pllconfig "github.com/komarovn654/embedded_configurator/config/pll_config"
)

type ConfigInterfaces map[string]interface{}

var (
	McuStm32f4xx = "stm32f4xx"
)

type Configs struct {
	pll *pllconfig.PllConfig
}

func (cnfg *Configs) GetPllConfig() *pllconfig.PllConfig {
	return cnfg.pll
}

func (cnfg *Configs) SetConfigTargets(targets ConfigInterfaces) error {
	for name, target := range targets {
		switch name {
		case pllconfig.ConfigName:
			if err := cnfg.pll.SetTarget(target); err != nil {
				return err
			}
		}
	}

	return nil
}

func NewConfig() *Configs {
	cnfg := new(Configs)
	cnfg.pll = pllconfig.New()
	return cnfg
}
