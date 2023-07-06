package config

import (
	"reflect"
	"testing"

	paths "github.com/komarovn654/embedded_configurator/config/common_paths"
	pllconfig "github.com/komarovn654/embedded_configurator/config/pll_config"
	logger "github.com/komarovn654/embedded_configurator/utils/log"
	"github.com/stretchr/testify/require"
)

type TestPllTarget struct {
}

func (t *TestPllTarget) SetupPll() error {
	return nil
}

func TestNew(t *testing.T) {
	logger.InitializeLogger()

	t.Run("new config", func(t *testing.T) {
		c := New()
		require.NotNil(t, c)

		typeC := reflect.TypeOf(*c)
		for i := 0; i < typeC.NumField(); i++ {
			require.NotNil(t, typeC.Field(i), typeC.Field(i).Name+"is nil")
		}
	})
}

func TestGetPllConfig(t *testing.T) {
	t.Run("get pll config", func(t *testing.T) {
		c := Configs{pll: pllconfig.New()}
		require.Equal(t, pllconfig.New(), c.GetPllConfig())
	})
}

func TestSetPllTarget(t *testing.T) {
	tests := []struct {
		name   string
		target interface{}
		expect interface{}
		err    error
	}{
		{
			name:   "pll target",
			target: &TestPllTarget{},
			expect: &TestPllTarget{},
			err:    nil,
		},
		{
			name:   "nil target",
			target: nil,
			expect: nil,
			err:    ErrorTargetType,
		},
		{
			name:   "not pll target",
			target: "string",
			expect: nil,
			err:    ErrorTargetType,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := Configs{pll: pllconfig.New()}
			err := c.setPllTarget(tc.target)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expect, c.pll.GetTarget())
		})
	}
}

func TestSetConfigTargets(t *testing.T) {
	logger.InitializeLogger()

	tests := []struct {
		name    string
		targets ConfigInterfaces
		expect  *pllconfig.PllConfig
		err     error
	}{
		{
			name:    "pll target",
			targets: ConfigInterfaces{pllconfig.ConfigName: &TestPllTarget{}},
			expect:  &pllconfig.PllConfig{Target: &TestPllTarget{}, Paths: paths.New()},
			err:     nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := Configs{pll: pllconfig.New()}
			err := c.SetConfigTargets(tc.targets)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expect, c.pll)
		})
	}
}
