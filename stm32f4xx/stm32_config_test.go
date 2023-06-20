package stm32_configs

import (
	"testing"

	"github.com/komarovn654/embedded_configurator/config"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx/pll_config"
	"github.com/stretchr/testify/require"
)

func assertPllIf(t *testing.T, cnfgs config.ConfigInterfaces) {
	for key, value := range cnfgs {
		require.Equal(t, key, config.PllConfigName)
		v, ok := value.(config.PllSourceIf)
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
