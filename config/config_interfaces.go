package config

var (
	PllConfigName = "pll"
)

type Interfaces map[string]interface{}

type PllSourceIf interface {
	SetupPll() error
}

type PllConfig struct {
	PllSrc PllSourceIf `mapstructure:"PllConfig"`

	Paths Common `mapstructure:"PllConfig"`
}

type Common struct {
	PllTemplate string `mapstructure:"PllTmpltPath"`
	PllDstPath  string `mapstructure:"PllDstPath"`
}
