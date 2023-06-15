package config

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/komarovn654/embedded_configurator/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type paths struct {
	dir  string
	name string
}

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
	cnfg.Close()
	os.Rename(cnfg.Name(), cnfg.Name()+".yaml")

	return paths{dir: tmpDir, name: cnfg.Name()}, err
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
		err  bool
	}{
		{
			name: "supported mcu",
			mcu:  "stm32f4xx",
			err:  false,
		},
		{
			name: "unsupported mcu",
			mcu:  "lpc1778",
			err:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cnfg := Config{parser: *viper.New()}
			f, err := os.OpenFile(paths.name+".yaml", os.O_RDWR, 0777)
			require.NoError(t, err, "setup environment error")

			_, err = f.WriteString("MCU: " + tc.mcu)
			require.NoError(t, err, "setup environment error")
			cnfg.parser.SetConfigName(paths.name)
			cnfg.parser.AddConfigPath(".")
			err = cnfg.parser.ReadInConfig()
			require.NoError(t, err, "setup environment error")

			err = cnfg.setMCUType()
			if tc.err {
				require.EqualError(t, err, ErrorMCUType.Error())
				return
			}
			require.NoError(t, err)
			for _, mcu := range McuTypes {
				if mcu == cnfg.configName {
					require.True(t, mcu == cnfg.configName)
					return
				}
			}
			require.Fail(t, "mcu doesnt math with supported")
		})
	}

}
