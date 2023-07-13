package configparser

import (
	"os"
	"strconv"
	"strings"
	"testing"

	logger "github.com/komarovn654/embedded_configurator/utils/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type confPaths struct {
	dir      string
	nameYaml string
	nameFull string
}

func setupTmp() (*confPaths, error) {
	logger.InitializeLogger()

	tmpDir, err := os.MkdirTemp("./", "TestCnfg")
	if err != nil {
		return nil, err
	}
	cnfg, err := os.CreateTemp(tmpDir, "cnfg")
	if err != nil {
		return nil, err
	}
	defer cnfg.Close()
	os.Rename(cnfg.Name(), cnfg.Name()+".yaml")

	cp := new(confPaths)
	cp.dir = tmpDir
	cp.nameYaml = strings.Split(cnfg.Name(), "/")[2]
	cp.nameFull = cnfg.Name() + ".yaml"
	return cp, err
}

func (cp *confPaths) writeConfigString(cnfgText string) error {
	f, err := os.OpenFile(cp.nameFull, os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	_, err = f.WriteString(cnfgText)

	return err
}

func TestNew(t *testing.T) {
	paths, err := setupTmp()
	require.NoError(t, err, "setup environment error")
	defer os.RemoveAll(paths.dir)

	t.Run("new config parser with opt", func(t *testing.T) {
		cp, err := New(SetConfigName(paths.nameYaml), SetConfigPath(paths.dir))
		require.NoError(t, err)
		require.NotNil(t, cp.parser)
		require.Equal(t, paths.nameYaml, cp.configName)
		require.Equal(t, []string{paths.dir}, cp.configPaths)
	})

	t.Run("new config parser without opt", func(t *testing.T) {
		_, err := New()
		require.Error(t, err)
	})
}

func TestSetConfigPath(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		cnfg := ConfigParser{}
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
		cnfg := ConfigParser{}
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
		cnfg := ConfigParser{parser: viper.New()}
		opt := SetConfigPath(paths.dir)
		opt(&cnfg)
		opt = SetConfigName(paths.nameYaml)
		opt(&cnfg)

		err := cnfg.parser.ReadInConfig()
		require.Error(t, err)

		cnfg.applyOptions()
		err = cnfg.parser.ReadInConfig()
		require.NoError(t, err)
	})
}

func TestReadMCUType(t *testing.T) {
	paths, err := setupTmp()
	require.NoError(t, err, "setup environment error")
	defer os.RemoveAll(paths.dir)

	t.Run("any mcu", func(t *testing.T) {
		err := paths.writeConfigString("MCU: " + "my mcu")
		require.NoError(t, err, "setup environment error")
		cnfg, err := New(SetConfigName(paths.nameYaml), SetConfigPath(paths.dir))
		require.NoError(t, err, "setup environment error")

		require.Equal(t, "my mcu", cnfg.ReadMcuType())
	})

	t.Run("no mcu", func(t *testing.T) {
		err := paths.writeConfigString(" ")
		cnfg, err := New(SetConfigName(paths.nameYaml), SetConfigPath(paths.dir))
		require.NoError(t, err, "setup environment error")

		require.Equal(t, "", cnfg.ReadMcuType())
	})
}
