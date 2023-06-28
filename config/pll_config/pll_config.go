package pll_config

import (
	cnfgCommon "github.com/komarovn654/embedded_configurator/config/common"
)

var (
	ConfigName = "pll"
)

type PllSettingsIf interface {
	SetupPll() error
}

type PllConfig struct {
	Settings PllSettingsIf           `mapstructure:"PllConfig"`
	Paths    *cnfgCommon.CommonPaths `mapstructure:"PllConfig"`
}

func (pc *PllConfig) GetSettings() PllSettingsIf {
	return pc.Settings
}

func (pc *PllConfig) GetPath() *cnfgCommon.CommonPaths {
	return pc.Paths
}
