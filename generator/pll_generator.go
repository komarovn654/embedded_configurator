package generator

import (
	"io"
	"os"
	"text/template"

	"github.com/komarovn654/embedded_configurator/utils"
)

type PllGenerator struct {
	destPath io.Writer

	tmpl *template.Template
}

func New(tmplPath string, destPath io.Writer) (*PllGenerator, error) {
	gnrt := PllGenerator{}
	gnrt.destPath = destPath
	if err := gnrt.init(tmplPath); err != nil {
		return nil, err
	}
	return &gnrt, nil
}

func (gnrt *PllGenerator) init(tmplPath string) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}
	utils.Logger.Sugar().Infof("pllGenerator instance: tmpl: %v", tmplPath)
	gnrt.tmpl = tmpl
	return nil
}

func (gnrt *PllGenerator) GenerateHeader(dest string, config any) error {
	if dest != "" {
		f, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
		defer f.Close()
		return gnrt.tmpl.Execute(f, config)
	}
	return gnrt.tmpl.Execute(os.Stdout, config)
}
