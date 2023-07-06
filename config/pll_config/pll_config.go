package pllconfig

import (
	paths "github.com/komarovn654/embedded_configurator/config/common_paths"
)

var (
	ConfigName = "pll"
)

type PllTargetIf interface {
	SetupPll() error
}

type PllConfig struct {
	Target PllTargetIf  `mapstructure:"PllConfig"`
	Paths  *paths.Paths `mapstructure:"PllConfig"`
}

func New() *PllConfig {
	cnfg := new(PllConfig)
	cnfg.Paths = paths.New()
	return cnfg
}

func (pc *PllConfig) SetTarget(target PllTargetIf) {
	pc.Target = target
}

func (pc *PllConfig) GetTarget() PllTargetIf {
	return pc.Target
}

func (pc *PllConfig) SetPaths(path *paths.Paths) {
	pc.Paths = path
}

func (pc *PllConfig) GetPaths() *paths.Paths {
	return pc.Paths
}
