package paths

type Paths struct {
	TemplatePath    string `mapstructure:"TemplatePath"`
	DestinationPath string `mapstructure:"DestinationPath"`
}

type Option = func(*Paths)

func SetTemplatePath(path string) Option {
	return func(p *Paths) {
		p.TemplatePath = path
	}
}

func SetDestinationPath(path string) Option {
	return func(p *Paths) {
		p.DestinationPath = path
	}
}

func New(opt ...Option) *Paths {
	p := new(Paths)
	for _, option := range opt {
		option(p)
	}
	return p
}

func (p *Paths) SetTemplatePath(path string) {
	p.TemplatePath = path
}

func (p *Paths) SetDestinationPath(path string) {
	p.DestinationPath = path
}

func (p *Paths) GetTemplatePath() string {
	return p.TemplatePath
}

func (p *Paths) GetDestinationPath() string {
	return p.DestinationPath
}
