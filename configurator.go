package main

import (
	"fmt"
	"log"

	stm32_pllconfig "github.com/komarovn654/embedded_configurator/stm32f4xx"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/Users/nikolajkomarov/stm32f4/tools/embedded_configurator")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(viper.AllKeys())

	var pllConfig stm32_pllconfig.PllConfig
	if err := viper.Unmarshal(&pllConfig); err != nil {
		log.Fatal(err)
	}

	fmt.Println(pllConfig)
	fmt.Println()
}
