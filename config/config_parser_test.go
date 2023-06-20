package config

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestSetConfigPath(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		cnfg := Config{}
		paths := make([]string, 10)
		for i := 0; i < 10; i++ {
			path := "test_path" + strconv.Itoa(i)
			opt := SetConfigPath(path)
			opt(&cnfg)
			paths[i] = path
		}

		require.Equal(t, cnfg.configPaths, paths)
	})
}

func TestSetConfigName(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		cnfg := Config{}
		name := "test_name"
		opt := SetConfigName(name)
		opt(&cnfg)

		require.Equal(t, cnfg.configName, name)
	})
}

func TestApplyOptions(t *testing.T) {
	paths, err := setupTmp()
	require.NoError(t, err, "setup environment error")
	defer os.RemoveAll(paths.dir)

	t.Run("common", func(t *testing.T) {
		cnfg := Config{parser: *viper.New()}
		path := paths.dir
		name := strings.Split(paths.name, "/")[2] // get file name only without path and ".yaml"
		opt := SetConfigPath(path)
		opt(&cnfg)
		opt = SetConfigName(name)
		opt(&cnfg)

		err := cnfg.parser.ReadInConfig()
		require.Error(t, err)

		cnfg.applyOptions()
		err = cnfg.parser.ReadInConfig()
		require.NoError(t, err)
	})
}

func TestSetMCUType(t *testing.T) {
	paths, err := setupTmp()
	require.NoError(t, err, "setup environment error")
	defer os.RemoveAll(paths.dir)

	tests := []struct {
		name string
		mcu  string
		err  error
	}{
		{
			name: "supported mcu",
			mcu:  "MCU: stm32f4xx",
			err:  nil,
		},
		{
			name: "unsupported mcu",
			mcu:  "MCU: lpc1778",
			err:  ErrorMCUType,
		},
		{
			name: "config withou mcu",
			mcu:  " ",
			err:  ErrorMCUKey,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cnfg := Config{parser: *viper.New()}
			err := writeConfigString(paths.name+".yaml", tc.mcu)
			require.NoError(t, err, "setup environment error")
			err = setupConfig(&cnfg, paths.name)
			require.NoError(t, err, "setup environment error")

			err = cnfg.setMCUType()
			require.Equal(t, tc.err, err)
			if err == nil {
				require.True(t, assertMCUType(cnfg.mcu))
			}
		})
	}
}

func TestGetMCUType(t *testing.T) {
	tests := []struct {
		name string
		mcu  string
	}{
		{
			name: "any mcu",
			mcu:  "my mcu",
		},
		{
			name: "empty string",
			mcu:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cnfg := Config{parser: *viper.New(), mcu: tc.mcu}

			require.Equal(t, tc.mcu, cnfg.GetMCUType())
		})
	}
}

func TestGetPllTmplPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "any path",
			path: "my path",
		},
		{
			name: "empty string",
			path: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cnfg := Config{parser: *viper.New(), Pll: &PllConfig{Paths: Common{PllTemplate: tc.path}}}

			require.Equal(t, tc.path, cnfg.GetPllTmplPath())
		})
	}
}

func TestGetPllDstPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "any path",
			path: "my path",
		},
		{
			name: "empty string",
			path: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cnfg := Config{parser: *viper.New(), Pll: &PllConfig{Paths: Common{PllDstPath: tc.path}}}

			require.Equal(t, tc.path, cnfg.GetPllDstPath())
		})
	}
}

func TestParsePllConfig(t *testing.T) {
	paths, err := setupTmp()
	require.NoError(t, err, "setup environment error")
	defer os.RemoveAll(paths.dir)

	tests := []struct {
		name      string
		config    ConfigMap
		err       error
		assertRes bool
	}{
		{
			name:      "valid config",
			config:    ValidConfig,
			err:       nil,
			assertRes: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// setup environment
			cnfg := Config{parser: *viper.New(), Pll: &PllConfig{}}
			err = writeConfigMap(paths.name+".yaml", tc.config)
			require.NoError(t, err, "setup environment error")
			err = setupConfig(&cnfg, paths.name)
			require.NoError(t, err, "setup environment error")

			// parse config
			pll := TestPllSource{}
			err = cnfg.parsePllConfig(&pll)
			require.Equal(t, tc.err, err)

			// assert pll fields
			refPll := TestPllSource{}
			ref := TestConfig{PllSrc: &refPll}
			mapstructure.Decode(ValidConfig, &ref)
			require.Equal(t, tc.assertRes, assertPllFields(refPll, &cnfg))
		})
	}
}

func TestParseConfig(t *testing.T) {
	paths, err := setupTmp()
	require.NoError(t, err, "setup environment error")
	defer os.RemoveAll(paths.dir)

	tests := []struct {
		name      string
		intf      ConfigInterfaces
		err       error
		assertRes bool
	}{
		{
			name:      "valid interfaces",
			intf:      ConfigInterfaces{PllConfigName: &TestPllSource{}},
			err:       nil,
			assertRes: true,
		},
		{
			name:      "valid interfaces",
			intf:      ConfigInterfaces{"unsup interface": &TestPllSource{}},
			err:       ErrorConfigType,
			assertRes: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cnfg := Config{parser: *viper.New(), Pll: &PllConfig{}}
			err = writeConfigMap(paths.name+".yaml", ValidConfig)
			require.NoError(t, err, "setup environment error")
			err = setupConfig(&cnfg, paths.name)
			require.NoError(t, err, "setup environment error")

			require.Equal(t, tc.err, cnfg.ParseConfig(tc.intf))
		})
	}

}
