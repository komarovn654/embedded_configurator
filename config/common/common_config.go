package common_config

type ConfigInterfaces map[string]interface{}

type CommonPaths struct {
	TemplatePath    string `mapstructure:"TemplatePath"`
	DestinationPath string `mapstructure:"DestinationPath"`
}

func New() *CommonPaths {
	return new(CommonPaths)
}

func (c *CommonPaths) GetTemplatePath() string {
	return c.TemplatePath
}

func (c *CommonPaths) GetDestinationPath() string {
	return c.DestinationPath
}
