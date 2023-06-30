package pll_config

import (
	cnfgCommon "github.com/komarovn654/embedded_configurator/config/common"
)

var (
	ConfigName = "pll"
)

type PllTargetIf interface {
	SetupPll() error
}

type PllConfig struct {
	Target PllTargetIf             `mapstructure:"PllConfig"`
	Paths  *cnfgCommon.CommonPaths `mapstructure:"PllConfig"`
}

func New() *PllConfig {
	cnfg := new(PllConfig)
	cnfg.Paths = cnfgCommon.New()
	return cnfg
}

func (pc *PllConfig) GetTarget() PllTargetIf {
	return pc.Target
}

func (pc *PllConfig) GetPath() *cnfgCommon.CommonPaths {
	return pc.Paths
}
