package pll_generator

import (
	"os"
	"text/template"

	l "github.com/komarovn654/embedded_configurator/utils/log"
)

type PllGenerator struct {
	tmpl *template.Template

	dstPath string
}

func New(tmplPath string, dstPath string) (*PllGenerator, error) {
	gnrt := PllGenerator{}
	gnrt.dstPath = dstPath
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

	l.Logger.Infof("pll generator instance: tmpl: %v", tmplPath)
	gnrt.tmpl = tmpl
	return nil
}

func (gnrt *PllGenerator) GenerateHeader(config any) error {
	if gnrt.dstPath != "" {

		f, err := os.OpenFile(gnrt.dstPath, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := gnrt.tmpl.Execute(f, config); err != nil {
			return err
		}
		l.Logger.Infof("header generated %v", gnrt.dstPath)
		return nil
	}
	l.Logger.Warn("no destination path, the generation will be in stdout")
	return gnrt.tmpl.Execute(os.Stdout, config)
}
