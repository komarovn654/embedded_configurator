package configparser

// import (
// 	"os"
// 	"reflect"

// 	l "github.com/komarovn654/embedded_configurator/utils/log"
// 	"gopkg.in/yaml.v2"
// )

// type ConfigMap map[string]map[string]interface{}

// type TestConfig struct {
// 	PllSrc PllSourceIf `mapstructure:"PllConfig"`

// 	Paths Common `mapstructure:"PllConfig"`
// }

// type TestPllSource struct {
// 	PllSource   string `mapstructure:"PllSource"`
// 	HseFreq     int    `mapstructure:"HseFrequency"`
// 	LseFreq     int    `mapstructure:"LseFrequency"`
// 	RequireFreq int    `mapstructure:"RequireFrequency"`
// }

// func (ps *TestPllSource) SetupPll() error {
// 	return nil
// }

// var (
// 	ValidConfig = ConfigMap{
// 		"PllConfig": {
// 			"PllSource":        "HSE",
// 			"HseFrequency":     8000000,
// 			"LseFrequency":     16000000,
// 			"RequireFrequency": 180000000,
// 		}}
// )

// func writeConfigMap(cnfgPath string, cnfgMap ConfigMap) error {
// 	f, err := os.OpenFile(cnfgPath, os.O_RDWR, 0777)
// 	if err != nil {
// 		return err
// 	}

// 	b, err := yaml.Marshal(cnfgMap)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = f.Write(b)

// 	return err
// }

// func setupConfig(cnfg *ConfigParser, cnfgName string) error {
// 	cnfg.parser.SetConfigName(cnfgName)
// 	cnfg.parser.AddConfigPath(".")
// 	return cnfg.parser.ReadInConfig()
// }

// func assertMCUType(mcu string) bool {
// 	for _, m := range McuTypes {
// 		if mcu == m {
// 			return true
// 		}
// 	}
// 	return false
// }

// func assertPllFields(cnfgRef TestPllSource, cnfg *ConfigParser) bool {
// 	srcCnfg := cnfg.Pll.PllSrc.(*TestPllSource)
// 	srcValue := reflect.ValueOf(*srcCnfg)
// 	srcType := reflect.TypeOf(*srcCnfg)

// 	refValue := reflect.ValueOf(cnfgRef)

// 	for i := 0; i < srcType.NumField(); i++ {
// 		_, ok := srcType.Field(i).Tag.Lookup("mapstructure")
// 		if ok {
// 			field := srcType.Field(i).Name
// 			if srcValue.FieldByName(field).Interface() != refValue.FieldByName(field).Interface() {
// 				return false
// 			}
// 		}
// 	}

// 	return true
// }
