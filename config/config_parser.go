package config

import (
	"errors"

	"github.com/komarovn654/embedded_configurator/utils"
	"github.com/spf13/viper"
)

var (
	McuSTM32F4xx = "stm32f4xx"
	McuTypes     = []string{McuSTM32F4xx}

	ErrorMCUKey  = errors.New("MCU key not found")
	ErrorMCUType = errors.New("unsupported mcu type")
)

type Config struct {
	mcu    string
	parser viper.Viper

	Pll PllConfig
}

func New() (*Config, error) {
	cnfg := Config{}
	if err := cnfg.init(); err != nil {
		return nil, err
	}

	return &cnfg, nil
}

func (cnfg *Config) init() error {
	cnfg.parser = *viper.New()
	cnfg.parser.AddConfigPath(".")
	cnfg.parser.AddConfigPath("config")
	cnfg.parser.SetConfigName("config")
	if err := cnfg.parser.ReadInConfig(); err != nil {
		return err
	}

	err := cnfg.setMCUType()
	return err
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

func (cnfg *Config) ParseConfig(configs Interfaces) error {
	for name, config := range configs {
		if name == "pll" {
			if err := cnfg.parsePllConfig(config); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cnfg *Config) parsePllConfig(config interface{}) error {
	c, ok := config.(PllSourceIf)
	if !ok {
		utils.Logger.Sugar().Fatalf("cast error") // TODO
	}
	cnfg.Pll.PllSrc = c

	if err := cnfg.parser.Unmarshal(&cnfg.Pll); err != nil {
		return err
	}

	utils.Logger.Sugar().Infof("pll src: %+v", cnfg.Pll.PllSrc)
	utils.Logger.Sugar().Infof("pll template: %v; tmpl dst: %v", cnfg.Pll.Paths.PllTemplate, cnfg.Pll.Paths.PllDstPath)

	return nil
}
