package main

import (
	"errors"
	"flag"
	"log"

	config "github.com/komarovn654/embedded_configurator/config"
	pllgenerator "github.com/komarovn654/embedded_configurator/generator/pll_generator"
	stm32targets "github.com/komarovn654/embedded_configurator/targets/stm32"
	configparser "github.com/komarovn654/embedded_configurator/utils/config_parser"
	logger "github.com/komarovn654/embedded_configurator/utils/log"
)

var (
	configName, configPath, loggerPath string

	ErrorUnknownMcuType = errors.New("unknown mcu type")
)

func GenerateHeadersPLL(cnfg *config.Configs) error {
	pllConfig := cnfg.GetPllConfig()
	if err := pllConfig.GetTarget().SetupPll(); err != nil {
		return err
	}

	gnrt, err := pllgenerator.New(
		pllConfig.GetPaths().GetTemplatePath(),
		pllConfig.GetPaths().GetDestinationPath(),
	)
	if err != nil {
		return err
	}

	return gnrt.GenerateHeader(pllConfig)
}

func GenerateHeadersSTM32(parser *configparser.ConfigParser) error {
	cnfg := config.New()
	if err := cnfg.SetConfigTargets(stm32targets.GetTargets()); err != nil {
		return err
	}

	if err := parser.ParseConfig(cnfg); err != nil {
		return err
	}

	return GenerateHeadersPLL(cnfg)
}

func initFlags() {
	flag.StringVar(&configName, "config", "", "yaml config name")
	flag.StringVar(&configPath, "config_path", ".", "config path")
	flag.StringVar(&loggerPath, "logger_path", ".log", "config path")
	flag.Parse()
}

func main() {
	initFlags()

	if err := logger.InitializeLogger(logger.SetLoggerPath(loggerPath)); err != nil {
		log.Fatal(err)
	}
	logger.Info("embedded configurator start")

	parser, err := configparser.New(
		configparser.SetConfigName(configName),
		configparser.SetConfigPath(configPath),
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
