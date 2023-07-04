package main

import (
	"errors"
	"log"

	config "github.com/komarovn654/embedded_configurator/config"
	stm32targets "github.com/komarovn654/embedded_configurator/targets/stm32"
	configparser "github.com/komarovn654/embedded_configurator/utils/config_parser"
	logger "github.com/komarovn654/embedded_configurator/utils/log"
)

var (
	ErrorUnknownMcuType = errors.New("unknown mcu type")
)

func GenerateHeadersPLL(cnfg *configparser.ConfigParser) error {
	// if err := cnfg.Pll.PllSrc.SetupPll(); err != nil {
	// 	return err
	// }

	// gnrt, err := pllgenerator.New(cnfg.GetPllTmplPath(), cnfg.GetPllDstPath())
	// if err != nil {
	// 	return err
	// }

	// return gnrt.GenerateHeader(cnfg.Pll)
	return nil
}

func GenerateHeadersSTM32(parser *configparser.ConfigParser) error {
	cnfg := config.NewConfig()
	if err := cnfg.SetConfigTargets(stm32targets.GetTargets()); err != nil {
		return err
	}

	err := parser.ParseConfig(cnfg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := logger.InitializeLogger(logger.SetLoggerPath("")); err != nil {
		log.Fatal(err)
	}
	logger.Info("embedded configurator start")

	parser, err := configparser.New(
		configparser.SetConfigName("config"),
		configparser.SetConfigPath("."),
	)
	if err != nil {
		logger.Fatal(err)
	}

	switch parser.ReadMcuType() {
	case config.McuStm32f4xx:
		err = GenerateHeadersSTM32(parser)
		if err != nil {
			logger.Fatal(err)
		}
	default:
		logger.Fatal(ErrorUnknownMcuType)
	}
}
