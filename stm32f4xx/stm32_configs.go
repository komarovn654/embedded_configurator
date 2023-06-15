package stm32_configs

import (
	"github.com/komarovn654/embedded_configurator/config"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx/pll_config"
)

func GetSTM32Configs() config.Interfaces {
	return config.Interfaces{config.PllConfigName: stm32_pllconfig.NewPllSource()}
}
