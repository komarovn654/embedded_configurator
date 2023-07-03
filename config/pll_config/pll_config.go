package pllconfig

import (
	"errors"

	paths "github.com/komarovn654/embedded_configurator/config/common_paths"
)

var (
	ConfigName = "pll"

	ErrorTargetType = errors.New("target interface doesnt implement PllTargetIf methods")
)

type PllTargetIf interface {
	SetupPll() error
}

type PllConfig struct {
	name   string
	Target PllTargetIf  `mapstructure:"PllConfig"`
	Paths  *paths.Paths `mapstructure:"PllConfig"`
}

func New() *PllConfig {
	cnfg := new(PllConfig)
	cnfg.Paths = paths.New()
	cnfg.name = ConfigName
	return cnfg
}

func (pc *PllConfig) SetTarget(target interface{}) error {
	t, ok := target.(PllTargetIf)
	if !ok {
		return ErrorTargetType
	}
	pc.Target = t
	return nil
}

func (pc *PllConfig) GetTarget() PllTargetIf {
	return pc.Target
}

func (pc *PllConfig) GetPath() *paths.Paths {
	return pc.Paths
}

func (pc *PllConfig) GetName() string {
	return pc.name
}
