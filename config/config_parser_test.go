package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/komarovn654/embedded_configurator/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type ConfigMap map[string]map[string]interface{}

type paths struct {
	dir  string
	name string
}

type TestConfig struct {
	PllSrc PllSourceIf `mapstructure:"PllConfig"`

	Paths Common `mapstructure:"PllConfig"`
}

type TestPllSource struct {
	PllSource   string `mapstructure:"PllSource"`
	HseFreq     int    `mapstructure:"HseFrequency"`
	LseFreq     int    `mapstructure:"LseFrequency"`
	RequireFreq int    `mapstructure:"RequireFrequency"`
}

func (ps *TestPllSource) SetupPll() error {
	return nil
}

var (
	ValidConfig = ConfigMap{
		"PllConfig": {
			"PllSource":        "HSE",
			"HseFrequency":     8000000,
			"LseFrequency":     16000000,
			"RequireFrequency": 180000000,
		}}
)

func setupTmp() (paths, error) {
	utils.InitializeLogger()

	tmpDir, err := os.MkdirTemp("./", "TestCnfg")
	if err != nil {
		return paths{}, err
	}
	cnfg, err := os.CreateTemp(tmpDir, "cnfg")
	if err != nil {
		return paths{}, err
	}
	defer cnfg.Close()

	os.Rename(cnfg.Name(), cnfg.Name()+".yaml")

	return paths{dir: tmpDir, name: cnfg.Name()}, err
}

func writeConfigString(cnfgPath string, cnfgText string) error {
	f, err := os.OpenFile(cnfgPath, os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	_, err = f.WriteString(cnfgText)

	return err
}

func writeConfigMap(cnfgPath string, cnfgMap ConfigMap) error {
	f, err := os.OpenFile(cnfgPath, os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(cnfgMap)
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func setupConfig(cnfg *Config, cnfgName string) error {
	cnfg.parser.SetConfigName(cnfgName)
	cnfg.parser.AddConfigPath(".")
	return cnfg.parser.ReadInConfig()
}

func assertMCUType(mcu string) bool {
	for _, m := range McuTypes {
		if mcu == m {
			return true
		}
	}
	return false
}

func assertPllFields(cnfgRef TestPllSource, cnfg *Config) bool {
	srcCnfg := cnfg.Pll.PllSrc.(*TestPllSource)
	srcValue := reflect.ValueOf(*srcCnfg)
	srcType := reflect.TypeOf(*srcCnfg)

	refValue := reflect.ValueOf(cnfgRef)

	for i := 0; i < srcType.NumField(); i++ {
		_, ok := srcType.Field(i).Tag.Lookup("mapstructure")
		if ok {
			field := srcType.Field(i).Name
			if srcValue.FieldByName(field).Interface() != refValue.FieldByName(field).Interface() {
				return false
			}
		}
	}

	return true
}

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
			cnfg := Config{parser: *viper.New(), Pll: PllConfig{Paths: Common{PllTemplate: tc.path}}}

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
			cnfg := Config{parser: *viper.New(), Pll: PllConfig{Paths: Common{PllDstPath: tc.path}}}

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
			cnfg := Config{parser: *viper.New()}
			err = writeConfigMap(paths.name+".yaml", ValidConfig)
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
			require.True(t, assertPllFields(refPll, &cnfg))
		})
	}

}
