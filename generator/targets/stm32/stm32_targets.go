package stm32_targets

import (
	cnfgCommon "github.com/komarovn654/embedded_configurator/config/common"
	cnfgPll "github.com/komarovn654/embedded_configurator/config/pll_config"
	stm32Pll "github.com/komarovn654/embedded_configurator/generator/targets/stm32/stm32_target_pll"
)

func GetSTM32Targets() cnfgCommon.ConfigInterfaces {
	return cnfgCommon.ConfigInterfaces{cnfgPll.ConfigName: stm32Pll.NewPllSettings()}
}
