package stm32targets

import (
	"testing"

	cnfg "github.com/komarovn654/embedded_configurator/config"
	cnfgPll "github.com/komarovn654/embedded_configurator/config/pll_config"
	stm32Pll "github.com/komarovn654/embedded_configurator/targets/stm32/stm32_target_pll"
	"github.com/stretchr/testify/require"
)

func assertPllIf(t *testing.T, cnfgs cnfg.ConfigInterfaces) {
	for key, value := range cnfgs {
		require.Equal(t, key, cnfgPll.ConfigName)
		v, ok := value.(cnfgPll.PllTargetIf)
		require.True(t, ok)
		require.IsType(t, new(stm32Pll.PllTarget), v)
	}
}

func TestGetSTM32Configs(t *testing.T) {
	t.Run("stm32 config interfaces", func(t *testing.T) {
		res := GetTargets()
		assertPllIf(t, res)
	})
}
