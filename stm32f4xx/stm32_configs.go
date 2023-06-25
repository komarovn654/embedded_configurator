package stm32_configs

import (
	parser "github.com/komarovn654/embedded_configurator/config_parser"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx/pll_config"
)

func GetSTM32Configs() parser.ConfigInterfaces {
	return parser.ConfigInterfaces{parser.PllConfigName: stm32_pllconfig.NewPllSource()}
}
