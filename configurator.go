package main

import (
	"github.com/komarovn654/embedded_configurator/config"
	"github.com/komarovn654/embedded_configurator/generator"
	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx"
	"github.com/komarovn654/embedded_configurator/utils"
)

func SetupSTM32Pll() *config.PllConfig {
	cnfg, err := config.New()
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	pllCnfg, err := cnfg.ParsePllConfig(stm32_pllconfig.NewPllSource())
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	if err := pllCnfg.PllSrc.SetupPll(); err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	return pllCnfg
}

func init() {
	utils.InitializeLogger()
}

type StmPllConfig struct {
	PllSrc stm32_pllconfig.PllSource
}

func main() {
	pllConfig := SetupSTM32Pll()
	gnrt, err := generator.New("/Users/nikolajkomarov/stm32f4/tools/embedded_configurator/stm32f4xx/pll.template")
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}

	gnrt.GenerateHeader(pllConfig)

	// config := Config{}
	// src := &stm32_pllconfig.PllSource{}
	// fmt.Println(config.Init())
	// fmt.Println(config.GetMCUType())
	// pll, _ := config.GetPllConfig(src)
	// fmt.Println(pll.PllSrc)
	// var pllConfig PllConfig

	// var PllSource stm32_pllconfig.PllSource
	// pllConfig.Pll = &PllSource

	// viper.AddConfigPath(".")
	// viper.SetConfigFile("config.yaml")
	// viper.ReadInConfig()

	// viper.Unmarshal(&pllConfig)
	// fmt.Printf("%+v\n", pllConfig.Pll)

	// fmt.Printf("%+v\n", pllConfig)

	// f, _ := os.Open("pll.template")
	// pllConfig.Pll.SetSrcFreq()
	// pllConfig.Pll.CalculateDivisionFactors()
	// templ, _ := template.ParseFiles("pll.template")
	// templ.Execute(os.Stdout, pllConfig)

	// templ2, _ := template.New("test2").Parse("#define PLL_REQUIRED_FREQUENCY {{.Src.RequireFreq}}\n")
	// templ2.Execute(os.Stdout, nil)
}
