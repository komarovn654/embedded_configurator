package main

import (
	"log"

	parser "github.com/komarovn654/embedded_configurator/config_parser"
	"github.com/komarovn654/embedded_configurator/generator"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx/pll_config"
	"github.com/komarovn654/embedded_configurator/utils"
)

func GenerateHeadersPLL(cnfg *parser.ConfigParser) error {
	if err := cnfg.Pll.PllSrc.SetupPll(); err != nil {
		utils.Logger.Fatal(err)
	}

	gnrt, err := generator.New(cnfg.GetPllTmplPath(), cnfg.GetPllDstPath())
	if err != nil {
		utils.Logger.Fatal(err)
	}

	return gnrt.GenerateHeader(cnfg.Pll)
}

func GenerateHeadersSTM32(cnfg *parser.ConfigParser) error {
	err := cnfg.ParseConfig(parser.ConfigInterfaces{parser.PllConfigName: stm32_pllconfig.NewPllSource()})
	if err != nil {
		utils.Logger.Fatal(err)
	}

	return GenerateHeadersPLL(cnfg)
}

func main() {
	if err := utils.InitializeLogger(utils.SetLoggerPath(".log")); err != nil {
		log.Fatal(err)
	}
	utils.Logger.Info("embedded configurator start")

	cnfg, err := parser.New(parser.SetConfigName("config"), parser.SetConfigPath("."))
	if err != nil {
		utils.Logger.Fatal(err)
	}

	switch cnfg.GetMCUType() {
	case parser.McuSTM32F4xx:
		err = GenerateHeadersSTM32(cnfg)
		if err != nil {
			utils.Logger.Fatal(err)
		}
	}
}
