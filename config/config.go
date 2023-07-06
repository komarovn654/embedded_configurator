package config

import (
	"errors"

	pllconfig "github.com/komarovn654/embedded_configurator/config/pll_config"
	logger "github.com/komarovn654/embedded_configurator/utils/log"
)

type ConfigInterfaces map[string]interface{}

var (
	McuStm32f4xx = "stm32f4xx"

	ErrorTargetType = errors.New("target interface cast error")
)

type Configs struct {
	pll *pllconfig.PllConfig
}

func New() *Configs {
	cnfg := new(Configs)
	cnfg.pll = pllconfig.New()
	logger.Info("new config created")
	return cnfg
}

func (cnfg *Configs) GetPllConfig() *pllconfig.PllConfig {
	return cnfg.pll
}

func (cnfg *Configs) setPllTarget(target interface{}) error {
	t, ok := target.(pllconfig.PllTargetIf)
	if !ok {
		return ErrorTargetType
	}

	cnfg.pll.SetTarget(t)
	return nil
}

func (cnfg *Configs) SetConfigTargets(targets ConfigInterfaces) error {
	for name, target := range targets {
		logger.Infof("config target setup: %v", name)
		switch name {
		case pllconfig.ConfigName:
			if err := cnfg.setPllTarget(target); err != nil {
				return err
			}
		}
	}

	return nil
}
