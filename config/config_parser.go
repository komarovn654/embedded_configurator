package config

import (
	"errors"
	"fmt"

	"github.com/komarovn654/embedded_configurator/utils"
	"github.com/spf13/viper"
)

var (
	McuSTM32F4xx = "stm32f4xx"
	McuTypes     = []string{McuSTM32F4xx}

	ErrorMCUKey     = errors.New("MCU key not found")
	ErrorMCUType    = errors.New("unsupported mcu type")
	ErrorConfigType = errors.New("unsupported config type")
)

type Option = func(c *Config)

type Config struct {
	parser      viper.Viper
	configPaths []string
	configName  string

	mcu string
	Pll *PllConfig
}

func New(opts ...Option) (*Config, error) {
	cnfg := Config{Pll: &PllConfig{}, parser: *viper.New()}
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

func (cnfg *Config) applyOptions() {
	for _, path := range cnfg.configPaths {
		cnfg.parser.AddConfigPath(path)
	}

	cnfg.parser.SetConfigName(cnfg.configName)
}

func SetConfigPath(name string) Option {
	return func(c *Config) {
		c.configPaths = append(c.configPaths, name)
	}
}

func SetConfigName(name string) Option {
	return func(c *Config) {
		c.configName = name
	}
}

func (cnfg *Config) setMCUType() error {
	if !cnfg.parser.IsSet("MCU") {
		return ErrorMCUKey
	}

	mcu, ok := cnfg.parser.Get("MCU").(string)
	if !ok {
		return nil
	}

	for _, supMcu := range McuTypes {
		if supMcu == mcu {
			cnfg.mcu = mcu
			return nil
		}
	}

	return ErrorMCUType
}

func (cnfg *Config) GetMCUType() string {
	return cnfg.mcu
}

func (cnfg *Config) GetPllTmplPath() string {
	return cnfg.Pll.Paths.PllTemplate
}

func (cnfg *Config) GetPllDstPath() string {
	return cnfg.Pll.Paths.PllDstPath
}

func (cnfg *Config) ParseConfig(configs ConfigInterfaces) error {
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

func (cnfg *Config) parsePllConfig(config interface{}) error {
	fmt.Println(config)
	c, ok := config.(PllSourceIf)
	if !ok {
		utils.Logger.Sugar().Fatalf("cast error") // TODO
	}
	cnfg.Pll.PllSrc = c

	if err := cnfg.parser.Unmarshal(cnfg.Pll); err != nil {
		return err
	}

	utils.Logger.Sugar().Infof("pll src: %+v", cnfg.Pll.PllSrc)
	utils.Logger.Sugar().Infof("pll template: %v; tmpl dst: %v", cnfg.Pll.Paths.PllTemplate, cnfg.Pll.Paths.PllDstPath)

	return nil
}
