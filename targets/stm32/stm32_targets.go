package stm32targets

import (
	cnfg "github.com/komarovn654/embedded_configurator/config"
	cnfgPll "github.com/komarovn654/embedded_configurator/config/pll_config"
	stm32Pll "github.com/komarovn654/embedded_configurator/targets/stm32/stm32_target_pll"
)

func GetTargets() cnfg.ConfigInterfaces {
	return cnfg.ConfigInterfaces{cnfgPll.ConfigName: stm32Pll.NewPllTarget()}
}
