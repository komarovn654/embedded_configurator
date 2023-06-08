package generator

import (
	"io"
	"os"
	"text/template"

	"github.com/komarovn654/embedded_configurator/utils"
)

type PllGenerator struct {
	tmplPath string
	destOut  io.Writer

	tmpl *template.Template
}

func New(tmplPath string) (*PllGenerator, error) {
	gnrt := PllGenerator{}
	if err := gnrt.init(tmplPath); err != nil {
		return nil, err
	}
	utils.Logger.Info("Pllgenerator instance created")
	return &gnrt, nil
}

func (gnrt *PllGenerator) init(tmplPath string) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	gnrt.tmpl = tmpl
	return nil
}

func (gnrt *PllGenerator) GenerateHeader(config any) error {
	return gnrt.tmpl.Execute(os.Stdout, config)
}
