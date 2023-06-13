package main

import (
	"os"

	"github.com/komarovn654/embedded_configurator/config"
	"github.com/komarovn654/embedded_configurator/generator"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx/pll_config"
	"github.com/komarovn654/embedded_configurator/utils"
)

func GenerateHeadersPLL(cnfg *config.Config) error {
	if err := cnfg.Pll.PllSrc.SetupPll(); err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	gnrt, err := generator.New(cnfg.GetPllTmplPath(), os.Stdout)
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	return gnrt.GenerateHeader(cnfg.Pll.Paths.PllDstPath, cnfg.Pll)
}

func GenerateHeadersSTM32(cnfg *config.Config) error {
	err := cnfg.ParseConfig(config.Interfaces{"pll": stm32_pllconfig.NewPllSource()})
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	return GenerateHeadersPLL(cnfg)
}

func init() {
	utils.InitializeLogger()
}

func main() {
	var err error
	cnfg, err := config.New()
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
