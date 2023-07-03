package paths

type Paths struct {
	TemplatePath    string `mapstructure:"TemplatePath"`
	DestinationPath string `mapstructure:"DestinationPath"`
}

func New() *Paths {
	return new(Paths)
}

func (c *Paths) GetTemplatePath() string {
	return c.TemplatePath
}

func (c *Paths) GetDestinationPath() string {
	return c.DestinationPath
}
