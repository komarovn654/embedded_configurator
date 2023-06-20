package config

import (
	"os"
	"reflect"

	"github.com/komarovn654/embedded_configurator/utils"
	"gopkg.in/yaml.v2"
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