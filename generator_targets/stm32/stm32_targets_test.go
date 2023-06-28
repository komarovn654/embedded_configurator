package stm32_targets

import (
	"testing"

	cnfgCommon "github.com/komarovn654/embedded_configurator/config/common"
	cnfgPll "github.com/komarovn654/embedded_configurator/config/pll_config"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx/pll_config"
	"github.com/stretchr/testify/require"
)

func assertPllIf(t *testing.T, cnfgs cnfgCommon.ConfigInterfaces) {
	for key, value := range cnfgs {
		require.Equal(t, key, cnfgPll.ConfigName)
		v, ok := value.(cnfgPll.PllSettingsIf)
		require.True(t, ok)
		require.IsType(t, &stm32_pllconfig.PllSource{}, v)
	}
}

func TestGetSTM32Configs(t *testing.T) {
	t.Run("stm32 config interfaces", func(t *testing.T) {
		res := GetSTM32Configs()
		assertPllIf(t, res)
	})
}
