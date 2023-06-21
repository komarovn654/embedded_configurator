package main

import (
	"github.com/komarovn654/embedded_configurator/config"
	"github.com/komarovn654/embedded_configurator/generator"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx/pll_config"
	"github.com/komarovn654/embedded_configurator/utils"
)

func GenerateHeadersPLL(cnfg *config.Config) error {
	if err := cnfg.Pll.PllSrc.SetupPll(); err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	gnrt, err := generator.New(cnfg.GetPllTmplPath(), cnfg.GetPllDstPath())
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	return gnrt.GenerateHeader(cnfg.Pll)
}

func GenerateHeadersSTM32(cnfg *config.Config) error {
	err := cnfg.ParseConfig(config.ConfigInterfaces{"pll": stm32_pllconfig.NewPllSource()})
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	return GenerateHeadersPLL(cnfg)
}

func init() {
	utils.InitializeLogger()
	utils.Logger.Sugar().Infof("embedded configurator start")
}

func main() {
	var err error
	cnfg, err := config.New(config.SetConfigName("config"), config.SetConfigPath("./config"))
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	switch cnfg.GetMCUType() {
	case config.McuSTM32F4xx:
		err = GenerateHeadersSTM32(cnfg)
		if err != nil {
			utils.Logger.Sugar().Fatal(err)
		}
	}
}
