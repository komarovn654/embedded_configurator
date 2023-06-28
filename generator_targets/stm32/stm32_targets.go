package stm32_targets

import (
	cnfgCommon "github.com/komarovn654/embedded_configurator/config/common"
	cnfgPll "github.com/komarovn654/embedded_configurator/config/pll_config"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx/pll_config"
)

func GetSTM32Configs() cnfgCommon.ConfigInterfaces {
	return cnfgCommon.ConfigInterfaces{cnfgPll.ConfigName: stm32_pllconfig.NewPllSource()}
}
