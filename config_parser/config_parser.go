package config_parser

import (
	"errors"

	"github.com/komarovn654/embedded_configurator/utils"
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

	mcu string
	Pll *PllConfig
}

func New(opts ...Option) (*ConfigParser, error) {
	cnfg := ConfigParser{Pll: &PllConfig{}, parser: viper.New()}
	for _, opt := range opts {
		opt(&cnfg)
	}

	cnfg.applyOptions()

	if err := cnfg.parser.ReadInConfig(); err != nil {
		return nil, err
	}

	err := cnfg.setMCUType()
	return &cnfg, err
}

func (cnfg *ConfigParser) applyOptions() {
	for _, path := range cnfg.configPaths {
		cnfg.parser.AddConfigPath(path)
	}

	cnfg.parser.SetConfigName(cnfg.configName)
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

func (cnfg *ConfigParser) setMCUType() error {
	if !cnfg.parser.IsSet("MCU") {
		return ErrorMCUKey
	}

	mcu, ok := cnfg.parser.Get("MCU").(string)
	if !ok {
		return ErrorMCUConfigType
	}

	for _, supMcu := range McuTypes {
		if supMcu == mcu {
			cnfg.mcu = mcu
			return nil
		}
	}

	return ErrorMCUType
}

func (cnfg *ConfigParser) GetMCUType() string {
	return cnfg.mcu
}

func (cnfg *ConfigParser) GetPllTmplPath() string {
	return cnfg.Pll.Paths.PllTemplate
}

func (cnfg *ConfigParser) GetPllDstPath() string {
	return cnfg.Pll.Paths.PllDstPath
}

func (cnfg *ConfigParser) ParseConfig(configs ConfigInterfaces) error {
	for name, config := range configs {
		switch name {
		case PllConfigName:
			return cnfg.parsePllConfig(config)
		default:
			return ErrorConfigType
		}
	}

	return nil
}

func (cnfg *ConfigParser) parsePllConfig(config interface{}) error {
	c, ok := config.(PllSourceIf)
	if !ok {
		return ErrorConfigInterface
	}
	cnfg.Pll.PllSrc = c

	if err := cnfg.parser.Unmarshal(cnfg.Pll); err != nil {
		return err
	}

	utils.Logger.Infof("pll src: %+v", cnfg.Pll.PllSrc)
	utils.Logger.Infof("pll template: %v; tmpl dst: %v", cnfg.Pll.Paths.PllTemplate, cnfg.Pll.Paths.PllDstPath)

	return nil
}
