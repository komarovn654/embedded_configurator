package main

import (
	"log"

	parser "github.com/komarovn654/embedded_configurator/config_parser"
	"github.com/komarovn654/embedded_configurator/generator"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx/pll_config"
	logger "github.com/komarovn654/embedded_configurator/utils/log"
)

func GenerateHeadersPLL(cnfg *parser.ConfigParser) error {
	if err := cnfg.Pll.PllSrc.SetupPll(); err != nil {
		return err
	}

	gnrt, err := generator.New(cnfg.GetPllTmplPath(), cnfg.GetPllDstPath())
	if err != nil {
		return err
	}

	return gnrt.GenerateHeader(cnfg.Pll)
}

func GenerateHeadersSTM32(cnfg *parser.ConfigParser) error {
	err := cnfg.ParseConfig(parser.ConfigInterfaces{parser.PllConfigName: stm32_pllconfig.NewPllSource()})
	if err != nil {
		return err
	}

	return GenerateHeadersPLL(cnfg)
}

func main() {
	if err := l.InitializeLogger(l.SetLoggerPath(".log")); err != nil {
		log.Fatal(err)
	}
	logger.Info("embedded configurator start")

	cnfg, err := parser.New(parser.SetConfigName("config"), parser.SetConfigPath("."))
	if err != nil {
		logger.Fatal(err)
	}

	switch cnfg.GetMCUType() {
	case parser.McuSTM32F4xx:
		err = GenerateHeadersSTM32(cnfg)
		if err != nil {
			logger.Fatal(err)
		}
	}
}
