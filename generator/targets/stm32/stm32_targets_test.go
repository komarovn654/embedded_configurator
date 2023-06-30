package stm32_targets

import (
	"testing"

	cnfgCommon "github.com/komarovn654/embedded_configurator/config/common"
	cnfgPll "github.com/komarovn654/embedded_configurator/config/pll_config"
	stm32Pll "github.com/komarovn654/embedded_configurator/generator/targets/stm32/stm32_target_pll"
	"github.com/stretchr/testify/require"
)

func assertPllIf(t *testing.T, cnfgs cnfgCommon.ConfigInterfaces) {
	for key, value := range cnfgs {
		require.Equal(t, key, cnfgPll.ConfigName)
		v, ok := value.(cnfgPll.PllSettingsIf)
		require.True(t, ok)
		require.IsType(t, &stm32Pll.PllSettings{}, v)
	}
}

func TestGetSTM32Configs(t *testing.T) {
	t.Run("stm32 config interfaces", func(t *testing.T) {
		res := GetSTM32Targets()
		assertPllIf(t, res)
	})
}
