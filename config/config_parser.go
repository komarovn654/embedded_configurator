package config

import (
	"errors"

	"github.com/spf13/viper"
)

var (
	ErrorMCUKey  = errors.New("MCU key not found")
	ErrorMCUType = errors.New("unsupported mcu type")
)

type PllSourceIf interface {
	SetupPll() error
}

type PllConfig struct {
	PllSrc PllSourceIf `mapstructure:"PllConfig"`
}

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
	cnfg.parser.AddConfigPath("/Users/nikolajkomarov/stm32f4/tools/embedded_configurator/config")
	cnfg.parser.SetConfigFile("/Users/nikolajkomarov/stm32f4/tools/embedded_configurator/config/config.yaml")
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

	if mcu, ok := cnfg.parser.Get("MCU").(string); ok {
		cnfg.mcu = mcu
		return nil
	}
	return ErrorMCUType
}

func (cnfg *Config) GetMCUType() string {
	return cnfg.mcu
}

func (cnfg *Config) ParsePllConfig(pllSource PllSourceIf) (*PllConfig, error) {
	cnfg.Pll.PllSrc = pllSource
	if err := cnfg.parser.Unmarshal(&cnfg.Pll); err != nil {
		return nil, err
	}

	return &cnfg.Pll, nil
}
