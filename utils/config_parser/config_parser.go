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
	cp := new(ConfigParser)
	cp.parser = viper.New()
	for _, opt := range opts {
		opt(cp)
	}
	cp.applyOptions()

	if err := cp.parser.ReadInConfig(); err != nil {
		return nil, err
	}
	logger.Infof("config parser created, config paths: %v, config name: %v",
		cp.configPaths, cp.configName)
	return cp, nil
}

func (cp *ConfigParser) applyOptions() {
	for _, path := range cp.configPaths {
		cp.parser.AddConfigPath(path)
	}

	cp.parser.SetConfigName(cp.configName)
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

func (cp *ConfigParser) ReadMcuType() string {
	mcu, ok := cp.parser.Get("MCU").(string)
	if !ok {
		return ""
	}

	return mcu
}

func (cp *ConfigParser) ParseConfig(config *config.Configs) error {
	configType := reflect.TypeOf(*config)
	for i := 0; i < configType.NumField(); i++ {
		logger.Infof("parse target: %v", configType.Field(i).Name)
		switch configType.Field(i).Name {
		case pllconfig.ConfigName:
			if err := cp.parsePllConfig(config.GetPllConfig()); err != nil {
				return err
			}
		default:
			return ErrorConfigType
		}
	}

	return nil
}

func (cp *ConfigParser) parsePllConfig(target *pllconfig.PllConfig) error {
	if err := cp.parser.Unmarshal(target); err != nil {
		return err
	}
	logger.Infof("target paths: %+v", target.GetPath())
	logger.Infof("target: %+v", target.GetTarget())

	return nil
}
