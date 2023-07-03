package configparser

import (
	"errors"
	"reflect"

	"github.com/komarovn654/embedded_configurator/config"
	pllconfig "github.com/komarovn654/embedded_configurator/config/pll_config"
	logger "github.com/komarovn654/embedded_configurator/utils/log"
	"github.com/spf13/viper"
)

var (
	McuSTM32F4xx = "stm32f4xx"
	McuTypes     = []string{McuSTM32F4xx}

	ErrorMCUKey          = errors.New("MCU key not found")
	ErrorMCUType         = errors.New("unsupported mcu type")
	ErrorMCUConfigType   = errors.New("mcu type must be a string")
	ErrorConfigType      = errors.New("unsupported config type")
	ErrorConfigInterface = errors.New("config interface cast error")
)

type Option = func(c *ConfigParser)

type ConfigParser struct {
	parser      *viper.Viper
	configPaths []string
	configName  string
}

func New(opts ...Option) (*ConfigParser, error) {
	parser := new(ConfigParser)
	parser.parser = viper.New()
	for _, opt := range opts {
		opt(parser)
	}
	parser.applyOptions()

	err := parser.parser.ReadInConfig()

	return parser, err
}

func (parser *ConfigParser) applyOptions() {
	for _, path := range parser.configPaths {
		parser.parser.AddConfigPath(path)
	}

	parser.parser.SetConfigName(parser.configName)
}

func SetConfigPath(name string) Option {
	return func(c *ConfigParser) {
		c.configPaths = append(c.configPaths, name)
	}
}

func SetConfigName(name string) Option {
	return func(c *ConfigParser) {
		c.configName = name
	}
}

func (parser *ConfigParser) ReadMcuType() string {
	mcu, ok := parser.parser.Get("MCU").(string)
	if !ok {
		return ""
	}

	return mcu
}

func (parser *ConfigParser) ParseConfig(config *config.Configs) error {
	configType := reflect.TypeOf(config)
	for i := 0; i < configType.NumField(); i++ {
		switch configType.Name() {
		case pllconfig.ConfigName:
			if err := parser.parsePllConfig(config.GetPllConfig().GetTarget()); err != nil {
				return err
			}
		default:
			return ErrorConfigType
		}
	}

	return nil
}

func (parser *ConfigParser) parsePllConfig(target pllconfig.PllTargetIf) error {
	if err := parser.parser.Unmarshal(target); err != nil {
		return err
	}
	logger.Info("pll parse done")
	// logger.Infof("pll src: %+v", target.PllSrc)
	// logger.Infof("pll template: %v; tmpl dst: %v", cnfg.Pll.Paths.PllTemplate, cnfg.Pll.Paths.PllDstPath)

	return nil
}
